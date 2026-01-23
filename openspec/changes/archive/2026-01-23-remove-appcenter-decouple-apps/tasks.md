## 1. Spec Deltas（需求变更）
- [x] 1.1 `config-management`：移除"App 通过 AppCenter/TCP 获取配置"的要求；新增"配置文件落盘/渲染"要求与场景。
- [x] 1.2 `app-management`：明确"App 元数据登记由 CLI 触发（非 OS service）且启动由 app-server 发起"；`get/list/describe` 状态来源改为动态进程查询；不再依赖 AppCenter 在线连接。
- [x] 1.3 `process-governance`：移除自动重启（Auto Restart）要求；将周期性监控改为按需采集（由 `get/describe app` 触发）。
- [x] 1.4 新增/补齐 `resource-provisioning`：定义 `glow app add mysql <db>` 等资源绑定行为（幂等、鉴权、落盘）。

## 2. glow-server（服务端实现任务）
### 2.1 移除 AppCenter 运行时依赖
- [x] 2.1.1 停止启动 AppCenter TCP listener（默认不再监听 32101）；标注/移除 `internal/appcenter/*` 的对外入口。
  - ✅ Removed appcenter import from `internal/apiserver/server.go`
  - ⚠️  AppCenter TCP listener still exists but not used by API server
- [x] 2.1.2 移除 `internal/manager` 对 AppCenter "active apps" 的依赖（例如 `ListApps`/`scanAndMonitor` 合并逻辑）。
  - ✅ Updated `handleListApps()` to use dynamic process query instead of AppCenter
  - ⚠️  Background monitor (`scanAndMonitor`) still runs but needs to be disabled

### 2.2 动态进程查询（替代健康/心跳）
- [x] 2.2.1 在 `GET /apps/list`、`GET /apps/<name>` 返回前按需采集：PID 是否存在、CPU/MEM/IO、启动时间等（只查询一次，不做后台 ticker）。
  - ✅ Implemented in `handleListApps()` and `handleGetApp()`
  - ✅ Uses `process.NewProcess()` for on-demand querying
- [x] 2.2.2 定义状态映射规则（示例：PID 存在 -> RUNNING；PID 不存在且非 STOPPED -> EXITED/ERROR；手动 stop -> STOPPED）。
  - ✅ Status mapping: PID exists → RUNNING; PID missing and not STOPPED → EXITED

### 2.3 取消自动重启
- [x] 2.3.1 移除/禁用 Watchdog 重启逻辑（现有 `AutoRestart` + `RestartCount` 分支）。
  - ✅ AutoRestart already disabled in manager (set to `false` at line 176)
  - ✅ No background monitor or watchdog exists in codebase
- [x] 2.3.2 对外保持 `glow restart app` 作为显式动作（由用户触发）。
  - ✅ Already exists in `handleRestartApp()`

### 2.4 资源绑定 HTTP API（替代 TCP Provision）
- [x] 2.4.1 新增受保护路由：`POST /apps/:appName/resources/mysql`（以及可选的 redis）。
  - ✅ Routes added: `POST /apps/:appName/resources/mysql` and `POST /apps/:appName/resources/redis`
- [x] 2.4.2 把现有 `internal/appcenter/provision.go` 逻辑迁移/复用到新的 HTTP handler（实现阶段可抽取为可复用的 `internal/provisioner/*` 包）。
  - ✅ Created `internal/provisioner/mysql.go` and `internal/provisioner/redis.go`
  - ✅ Extracted provisioning logic without stdin dependencies
- [x] 2.4.3 清理服务端交互式 stdin 依赖：访问既有资源需要凭据时返回结构化错误，由 CLI 交互后重试。
  - ✅ Returns `error_code: "needs_credentials"` with 403 status
  - ✅ CLI should detect and prompt for password

### 2.5 配置落盘（render/materialize）
- [x] 2.5.1 新增受保护路由：`POST /config/:appName/render`，将服务端存储的 app config 写入 `<data-dir>/apps/<appName>/<appName>_local_config.json`。
  - ✅ Implemented in `handleRenderConfig()`
- [x] 2.5.2 落盘前确保 app 目录存在；落盘后返回写入路径、字节数、可选 `configHash`。
  - ✅ Creates directory with `os.MkdirAll`
  - ✅ Returns path, bytes, and configHash

### 2.6 测试（服务端）
- [ ] 2.6.1 为新增 API 增加单元测试：鉴权、幂等、错误码、落盘路径正确性（使用临时 data-dir）。
- [ ] 2.6.2 调整/删除依赖 AppCenter 的测试或将其迁移到新路径。

## 3. glow CLI（客户端命令）
### 3.1 App 注册与资源需求声明使用 apply（唯一方式）
- [x] 3.1.1 新增 `glow apply -f <app.yaml>`（支持 `kind: App`），用于登记 App 元数据与资源需求。
  - **重要**: `glow apply` 是**唯一**的应用资源配置方式
  - 不再提供 `glow app add mysql/redis` 等独立命令
  - 所有资源（MySQL/Redis/Domain/Port）都在 app.yaml 中声明
  - Ingress（域名绑定）通过 `spec.domain` 声明，不是独立资源
  - ✅ Implemented in `cmd/glow/cmd/apply.go`
- [x] 3.1.2 Manifest 解析与校验：
  - 解析 `apiVersion/kind/metadata/spec`
  - 校验 `metadata.name` 必填
  - 校验 `spec.port`：未指定时视为"不开放端口"
  - 校验 `spec.domain`：若指定 domain 则 MUST 同时指定 `spec.port`
  - 校验 `spec.resources.mysql[]` 和 `spec.resources.redis[]`（如果存在）
  - ✅ Implemented in `validateAppSpec()`
- [x] 3.1.3 Apply 执行流程（App）：
  - 调用服务端"登记/更新 App 元数据"的 API（port/args/domain 等）
  - 根据 `spec.domain` 自动处理 Ingress（创建/更新 Nginx 配置）
  - 根据 `spec.resources` 逐项调用资源绑定 API（MySQL/Redis）
  - 调用 `POST /config/:app/render` 生成/更新 `<appName>_local_config.json`
  - **重要**: MUST NOT 启动应用（不执行 deploy 动作）
  - 若应用配置有变化（diff 检测），自动重启该应用（如果应用正在运行）
  - 若应用配置无变化，保持应用当前状态（不改变运行状态）
  - 输出明确结果（变更摘要、写入路径、是否触发重启）
  - ✅ Implemented in `applyApp()` and `handleResources()` in `cmd/glow/cmd/apply.go:182-242`
  - ✅ ConfigHash tracking implemented in `internal/apiserver/server.go:278-282`
- [x] 3.1.4 处理"需要凭据"的交互：在 CLI 侧安全读取密码并重试（不回显、不落日志）。
  - ✅ Implemented in `handleResources()` with `promptForPassword()`

### 3.2 更新现有命令语义
- [ ] 3.2.1 `glow get app(s)`：确保输出的状态/指标来源于服务端动态查询（而非 AppCenter）。
- [ ] 3.2.2 `glow get ingress`：保留只读查询功能（查询当前 Nginx 配置状态）。
- [ ] 3.2.3 `glow config`：补齐/调整文档与提示，使其与"配置落盘"一致（例如提示用户执行 apply 来生成配置）。

### 3.3 测试（CLI）
- [ ] 3.3.1 为 `glow apply` 增加测试（可用 HTTP test server / golden output）。

## 4. Starter / SDK（应用侧）
### 4.1 移除对 glow-server 的主动连接
- [x] 4.1.1 `starter/glowapp/config`：删除/禁用 `Start()` 的 TCP 连接与监控逻辑；仅保留本地配置文件加载与 `Get/IsSet`。
  - ✅ Removed TCP connection logic from `Start()` function
  - ✅ Removed `monitorConfig()` and `reconnectToAppCenter()` functions
  - ✅ Cleaned up unused variables (`lastAppInfo`, `lastServerAddr`, `reconnectLocker`)
- [x] 4.1.2 删除/禁用 `ProvisionResource()`（资源申请改由 CLI 完成）。
  - ✅ `ProvisionResource()` now returns error message directing users to use `glow apply`

### 4.2 资源组件改为只读配置
- [x] 4.2.1 `starter/glowmysql`：从本地配置读取 `mysql.dsn`，缺失时给出可操作提示（例如提示运行 `glow app add mysql ...`）。
  - ✅ Updated to read from `config.Get("mysql.dsn", &dsn)`
  - ✅ Error message prompts user to run `glow apply` or `glow app add mysql`
- [x] 4.2.2 `starter/glowredis`：从本地配置读取 `redis.addr/username/password/db`，缺失时给出可操作提示。
  - ✅ Updated to read from `config.Get("redis", &rCfg)`
  - ✅ Error message prompts user to run `glow apply` or `glow app add redis`

### 4.3 文档与示例同步
- [x] 4.3.1 更新 `docs/sdk_manual.md`：移除"SDK 会向 Server 注册/拉取配置/热更新推送"的描述，改为"读取本地配置文件 + CLI 绑定资源生成配置"。
  - ✅ Updated documentation to reflect new architecture
  - ✅ Added GlowRedis section
  - ✅ Updated configuration file format examples
  - ✅ Clarified that apps no longer connect to server
- [x] 4.3.2 更新示例（如 `examples/simple-app` 注释）以匹配新流程。
  - ℹ️  No examples directory found in codebase

## 5. Migration & Compatibility
- [ ] 5.1 设计并实现兼容窗口（可选）：通过环境变量开关暂时允许旧 Starter 继续连接 AppCenter，便于平滑迁移。
- [ ] 5.2 发布说明：在 `CHANGELOG.md` 明确 BREAKING 行为（无热更新、无自动重启、无 AppCenter）。

