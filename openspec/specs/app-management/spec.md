# 应用管理 (App Management)

## Purpose
负责应用全生命周期的核心管理操作，包括应用实例的启动、优雅停止、故障重启以及实时状态查询，确保应用的可控性。
## Requirements
### Requirement: 启动应用 (Start App)
系统 MUST 能够启动一个新的应用实例。如果应用已在运行，操作 MUST 幂等。

#### Scenario: 成功启动新应用
- **WHEN** 客户端发送合法的启动请求（包含名称、命令、参数等）
- **THEN** 系统应检查应用是否已在运行
- **AND** 如果未运行，系统应准备运行环境并启动进程
- **AND** 如果应用配置中未指定端口，系统应分配一个空闲端口
- **AND** 系统应返回成功响应

#### Scenario: 启动已运行的应用 (幂等性)
- **WHEN** 客户端请求启动一个状态为 RUNNING 的应用
- **THEN** 系统应直接返回成功响应，不执行任何操作

### Requirement: 停止应用 (Stop App)
系统 MUST 能够优雅地停止正在运行的应用，并标记为手动停止状态。

#### Scenario: 停止运行中的应用
- **WHEN** 客户端发送停止请求指定应用名称
- **THEN** 系统应发送 `SIGTERM` 信号
- **AND** 系统应更新应用状态为 `STOPPED` (Manual Stop)
- **AND** Watchdog 不应自动重启状态为 `STOPPED` 的应用

### Requirement: 应用列表 (List Apps)
系统 MUST 能够列出所有受管应用及其当前状态。

#### Scenario: 获取列表
- **WHEN** 客户端请求应用列表
- **THEN** 系统应返回所有应用的名称、PID、状态 (RUNNING/STOPPED/ERROR)、端口及资源使用统计

### Requirement: 查看日志 (App Logs)
系统 MUST 提供访问应用标准输出/错误日志的能力。

#### Scenario: 读取日志
- **WHEN** 客户端请求指定应用的日志
- **THEN** 系统应从 `apps/<name>/logs` 读取并返回日志内容

### Requirement: 删除应用 (Delete App)
系统 MUST 能够彻底删除应用及其相关资源。

#### Scenario: 删除应用
- **WHEN** 客户端发送删除请求指定应用名称
- **THEN** 系统应停止该应用（如果正在运行）
- **AND** 系统应清理该应用的运行目录、日志和配置信息
- **AND** 系统应从应用列表中移除该应用

### Requirement: 重启应用 (Restart App)
系统 MUST 能够重启应用。

#### Scenario: 重启应用
- **WHEN** 客户端发送重启请求
- **THEN** 系统应先停止应用
- **AND** 系统应重新启动应用

### Requirement: 状态监控 (Status Monitoring)
系统 MUST 监控应用进程状态，并区分正常退出和异常退出。

#### Scenario: 异常退出
- **WHEN** 应用进程退出 (无论 Exit Code 为何) 且状态非 `STOPPED`
- **THEN** 系统应更新状态为 `ERROR`
- **AND** Watchdog 应尝试自动重启

### Requirement: 应用管理 (App Management)
系统 MUST 使用 "App" (应用) 作为核心资源定义，并提供类 K8s 的 CLI 操作接口。

#### Scenario: 获取应用列表
- **WHEN** 用户执行 `glow get app`
- **THEN** CLI 应列出所有 App 的状态（NAME, STATUS, RESTARTS, AGE, CPU, MEM, PID, PORT, DOMAIN）

#### Scenario: 查看应用详情
- **WHEN** 用户执行 `glow describe app <name>`
- **THEN** CLI 应显示指定 App 的详细信息（Events, Config, Resources）

#### Scenario: 删除应用
- **WHEN** 用户执行 `glow delete app <name>`
- **THEN** CLI 应向服务端发送删除请求

#### Scenario: 重启应用
- **WHEN** 用户执行 `glow restart app <name>`
- **THEN** CLI 应触发滚动重启或原地重启

#### Scenario: 停止应用 (Stop)
- **WHEN** 用户执行 `glow stop app <name>`
- **THEN** CLI 应停止该 App 的运行进程

#### Scenario: 启动应用 (Start)
- **WHEN** 用户执行 `glow start app <name>`
- **THEN** CLI 应启动该 App 的运行进程

#### Scenario: 查看日志
- **WHEN** 用户执行 `glow logs <name>`
- **THEN** CLI 应获取并打印日志（支持 -f 实时流式传输）

### Requirement: Application State Validation
The system MUST provide an endpoint to query the configuration and binary state of a deployed application.

#### Scenario: Get App State
- **GIVEN** an application "my-app" is deployed
- **WHEN** the client sends `GET /apps/my-app/state`
- **THEN** the server MUST return a JSON object containing the current `configHash` and `binaryHash`.

### Requirement: Application Binary Upload
The system MUST provide an endpoint for clients to upload application binaries.

#### Scenario: Upload a New Binary
- **WHEN** the client sends a `POST` request to `/apps/my-app/binary` with a multipart/form-data payload containing the binary file
- **THEN** the server MUST securely save the binary to a managed location.
- **AND** the server MUST update the application's stored `binaryHash`.

