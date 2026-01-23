## Context
现状采用 **Server-Agent + Starter** 模式：应用引入 `glow/starter` 后会在启动阶段主动连接 glow-server 的 AppCenter(TCP)，执行：
- 上报/注册（runtime state）
- 获取初始配置（并在 Server 热更新时通过长连接接收推送）
- 按需申请资源（`ActionProvision` 触发 MySQL/Redis 创建与凭据下发）

服务端（`glow-server`）再通过 AppCenter 的在线连接来推断存活、采集指标，并执行自动重启（Watchdog）。

这带来两类耦合：
- **运行时耦合**：app 不可避免依赖 glow-server 可达（至少启动期/资源期/配置期）。
- **治理耦合**：进程存活判定、状态采集、重启策略等逻辑与 AppCenter 在线状态绑定，导致“网络/连接问题”与“进程故障”难以区分。

本变更的核心意图：把“运行时依赖”从 app 侧移到 CLI/运维侧，让 app 在运行时只依赖 **本地配置文件 + 环境变量**。

## Goals / Non-Goals
### Goals
- App 运行时 **不再主动连接** glow-server（去除 AppCenter/TCP 依赖）。
- app **不注册为 service**（既不注册为 OS service，也不进行在线心跳注册）；应用进程由 **app-server（glow-server）发起启动**，并且后续不保活。
- app 的“元数据登记/资源需求声明/配置生成”由 **glow CLI 的 `apply`** 驱动完成，并落盘为 app 可读取的配置文件。
- glow-server **不再负责保活**：不自动重启、不依赖在线连接判断健康。
- `glow get app`/`glow describe app` 的状态与指标来源改为 **按需动态查询 PID/进程信息**。

### Non-Goals
- 不引入分布式服务发现（DNS/Consul/etcd 等）。
- 不保证“任意方式启动的进程都能被发现”。本提案的动态查询聚焦于 glow 已登记（registered）的 App。
- 不在 proposal 阶段实现/落地代码（仅给出可执行设计与 specs/tasks）。

## High-Level Design
### 新的责任边界
- **App（Runtime）**
  - 只负责：读取本地配置文件、启动自身服务、优雅退出。
  - 不负责：注册为 OS service、向 glow-server 在线注册/心跳、向 glow-server 拉取配置、向 glow-server 申请资源。
- **glow CLI（Control Plane Client）**
  - 负责：通过 `glow apply -f app.yaml` 登记 App 元数据（name/command/workdir/env/port/domain 等）与资源需求（mysql/redis...）、并触发配置文件生成（render/materialize）。
  - 运行在用户侧，可交互，可决定何时做资源申请与配置更新。
- **glow-server（Host Daemon / Control Plane Server）**
  - 继续负责：配置存储（管理面）、资源实际创建（MySQL/Redis on host）、应用启动/停止（如果保留现有 app lifecycle API）。
  - 不再负责：AppCenter 在线连接管理、周期性健康监控、自动重启。

### 启动模型（app-server 发起）
- app 的启动由 glow-server 侧的管理面动作触发（例如 `glow start app <name>` → Server 启动进程）。
- app 启动后不需要也不应进行“注册为 service/心跳保活”的回连动作；Server 侧不依赖连接状态判断存活。

### 配置中心与“服务发现”
本变更把服务发现收敛为 **配置文件分发**：
- 资源的“发现结果”（如 MySQL DSN、Redis Addr/Password）以 **JSON config** 写入 app 的配置文件。
- app 通过 `starter/glowapp/config` 读取这些字段，不需要任何网络发现机制。

### 配置文件落盘位置与命名
为复用现有 SDK 行为（`starter/glowapp/config/loader.go`），配置文件采用：
- **文件名**：`<appName>_local_config.json`（优先），否则 `local_config.json`
- **落盘目录**：`<data-dir>/apps/<appName>/`（当 app 由 glow-server 托管启动时，该目录也是默认工作目录）

该选择保证：
- 远程运行（CLI 在本地、Server 在远端）时，配置文件可以由 Server 侧落盘；
- app runtime 只需读取本地文件即可启动，无需连接 Server。

### 资源绑定（以 MySQL 为例）
#### 目标用户体验
- `glow apply -f app.yaml`（`kind: App`）- **唯一的资源配置方式**
  - 触发：在 YAML 中声明该 App 需要 MySQL（仅包含 db name），由控制面为该 App 创建/复用 DB 与凭据
  - 结果：更新服务端存储的 app config，并在 Server 端为该 app 生成/更新 `<appName>_local_config.json`
  - **重要**: 不提供 `glow app add mysql/redis` 等独立命令，所有资源统一通过 `glow apply` 声明式配置

### App Apply 文件（Manifest）模型
App 的声明式文件（`kind: App`）应至少包含：
- `metadata.name`: 应用名
- `spec.binary`: 应用二进制路径
- `spec.port`（可选）：开放端口（**若缺省则视为不开放端口**，服务端 MUST NOT 为其分配端口，也不应注入 `OP_APP_PORT`）
- `spec.domain`（可选）：绑定域名（用于 ingress/反向代理；**仅当指定 `spec.port` 时才有意义**）
  - **重要**: Ingress 不是独立资源，通过 App YAML 的 `spec.domain` 直接声明
- `spec.args`（可选）：执行参数
- `spec.resources`（可选）：资源需求
  - `spec.resources.mysql[]`: 每项仅包含 `dbName`
  - `spec.resources.redis[]`: 包含 `db` 等配置

该文件用于"登记/更新"应用期望状态与资源需求，本变更不要求 App 运行时参与任何注册/发现流程。

**示例 App Manifest**（包含 Ingress 配置）:
```yaml
apiVersion: v1
kind: App
metadata:
  name: my-app
spec:
  binary: ./bin/my-app
  workingDir: /var/lib/glow-server/apps/my-app
  port: 8080
  domain: myapp.example.com  # Ingress 直接在 App 中声明
  args:
    - "--server"
  env:
    - name: ENV
      value: production
  resources:
    mysql:
      - dbName: myapp_db
    redis:
      - db: 0
```

#### 交互/鉴权
- 资源绑定与配置落盘属于“管理面能力”，沿用现有 `Authorization: Bearer <api_key>` 认证。

#### 关于交互式密码输入
现有实现把“访问既有 DB/用户时的密码输入”放在 glow-server 侧 stdin，这在 daemon 模式不可用。
本提案建议：
- **交互只发生在 CLI**（例如 `glow app add mysql ...` 需要密码时在 CLI 询问）
- Server API 返回可操作错误（如 `needs_credentials`）提示 CLI 补齐参数后重试

## API / Data Model Sketch
> 具体路径可在实现阶段微调，但需要满足 specs 的行为。

- **资源绑定（MySQL）**
  - `POST /apps/:appName/resources/mysql`
  - Body: `{ "dbName": "xxx", "mode": "create_or_use", "existingPassword": "..." }`
  - Response: `{"mysql":{"dsn":"..."}}`（与当前 `provisionMySQL` 返回结构一致）

- **配置落盘**
  - `POST /config/:appName/render`
  - 行为：读取服务端存储的 app config，写入 `<data-dir>/apps/<appName>/<appName>_local_config.json`
  - Response: `{ "path": "...", "bytes": 123, "configHash": "..." }`

- **应用列表/详情（动态进程查询）**
  - `GET /apps/list` 与 `GET /apps/:name`
  - 行为：以 DB 中的 app 记录为主，按需用 PID（或可推导的 binary path）动态查询进程存在性与 CPU/MEM/IO。

## Migration Plan
- Step 1: CLI 与 Server 先支持“资源绑定 + 配置落盘”的新路径（App 仍可暂时保留旧 AppCenter 行为以兼容）。
- Step 2: Starter 移除对 AppCenter/TCP 的连接与资源申请代码路径（仅保留本地配置文件读取）。
- Step 3: Server 移除 AppCenter 依赖与 Watchdog（停止基于连接的存活判定、移除自动重启逻辑）。
- Step 4: 文档与示例迁移到“CLI 驱动资源绑定 + 本地配置文件运行”的新范式。

## Risks / Trade-offs
- 失去配置热更新推送：新模型默认需要 CLI/运维主动更新并重启或由应用自行实现 reload。
- 失去自动保活：需要外部守护（systemd/容器编排）或显式 `glow restart app` 管理。
- 进程状态更“真实”：从连接状态切到 PID/进程查询后，可减少误判，但也需要对 PID 复用/重启场景做谨慎处理（实现阶段需定义策略）。

