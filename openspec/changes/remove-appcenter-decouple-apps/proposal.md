# Change: 去除 App 对 glow-server 的主动依赖（移除 AppCenter 注册/发现）

## Why
当前 Glow 的运行时依赖以 **App(Starter) → glow-server(AppCenter/TCP)** 为中心：应用启动时需要主动连接并上报/注册、获取配置与资源（Provision），服务端再依赖该连接来判断存活与执行自动重启。这导致应用与 glow-server 强耦合，降低可移植性与可运维性（例如：离线/本地运行、分层权限、网络隔离、只读运行环境等）。

## What Changes
- **BREAKING**：应用侧（`glow/starter`）不再主动连接 glow-server（不再依赖 AppCenter/TCP）。
- 应用不再"注册为 service"（既不注册为 OS service，也不通过 AppCenter 进行在线注册/心跳）；应用进程由 **app-server（glow-server）发起启动**。
- 应用"元数据登记/资源需求声明/服务发现(配置分发)"由 **glow CLI 的 `apply`** **唯一驱动**完成；App 的 `apply` 文件包含开放端口、执行参数、绑定域名、资源需求（如 MySQL 仅包含 db name）等声明。
  - **重要**: `glow apply` 是**唯一**的应用资源配置方式，不提供 `glow app add mysql/redis` 等独立命令
  - Ingress（域名绑定）通过 App YAML 中的 `spec.domain` 声明，不是独立资源
- 配置中心从"TCP 长连接推送 + 应用侧写本地文件"改为"**CLI/Server 生成配置文件** → 应用只读加载"，应用侧不再承担注册/发现职责。
- **BREAKING**：应用不再作为"可保活服务"被 glow-server 自动重启；glow-server 不再管理应用健康与保活。
- `glow get app`/`glow describe app` 的状态与资源指标改为 **动态查询 PID/进程信息**（按需采集），不依赖 AppCenter 在线连接。

## Impact
- Affected specs:
  - `app-management`
  - `config-management`
  - `process-governance`
  - `resource-provisioning`（新增/补齐：由 CLI 为 app 绑定资源并生成配置）
- Affected code (implementation stage):
  - Starter: `starter/glowapp/config/*`, `starter/glowapp/*`, `starter/glowmysql/*`, `starter/glowredis/*`
  - Server: `internal/appcenter/*`, `internal/manager/*`, `internal/apiserver/*`
  - CLI: `cmd/glow/*`（新增 `glow app add mysql ...` 等命令与调用链）
  - Docs: `docs/sdk_manual.md`, `docs/server_manual.md`

