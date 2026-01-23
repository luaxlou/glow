# 应用管理 (App Management)

## Purpose
负责应用全生命周期的核心管理操作，包括应用元数据登记、应用实例的启动、优雅停止以及状态查询。应用不注册为 OS service，不依赖在线连接，应用进程由 glow-server 发起启动。
## Requirements
### Requirement: 启动应用 (Start App)
系统 MUST 能够启动一个新的应用实例。如果应用已在运行，操作 MUST 幂等。应用启动流程 MUST NOT 要求应用侧主动连接 glow-server（无 AppCenter/TCP 依赖）。

#### Scenario: 成功启动新应用
- **WHEN** 客户端发送合法的启动请求（包含名称、命令、参数等）
- **THEN** 系统应检查应用是否已在运行
- **AND** 如果未运行，系统应准备运行环境并启动进程
- **AND** 若应用声明了需要开放端口（例如通过 App manifest 指定 `spec.port`），系统应为其配置该端口（并注入 `OP_APP_PORT`）
- **AND** 若应用未声明端口，则系统 MUST 将其视为“不开放端口”，并 MUST NOT 为其分配端口
- **AND** 系统应返回成功响应

#### Scenario: 由 app-server 发起启动且应用不注册为 service
- **WHEN** 用户执行 `glow start app <name>`（或等价命令）触发 app-server（glow-server）启动应用进程
- **THEN** 应用启动后 MUST NOT 依赖向 glow-server 发起在线注册/心跳连接来完成启动流程
- **AND** 系统 MUST NOT 将该应用注册为 OS service

#### Scenario: 启动已运行的应用 (幂等性)
- **WHEN** 客户端请求启动一个状态为 RUNNING 的应用
- **THEN** 系统应直接返回成功响应，不执行任何操作

### Requirement: 应用元数据登记 (Register App Metadata)
系统 MUST 支持通过 CLI 登记/更新应用元数据，应用不主动注册。

#### Scenario: 通过 CLI 登记 App 元数据
- **WHEN** 用户执行 `glow apply -f app.yaml`（`kind: App`）
- **THEN** CLI 应解析应用元数据（name/command/workdir/env/port/domain 等）
- **AND** 调用服务端 API 登记/更新该应用的元数据
- **AND** 应用不注册为 OS service，不进行在线心跳注册

### Requirement: 停止应用 (Stop App)
系统 MUST 能够优雅地停止正在运行的应用，并标记为手动停止状态。

#### Scenario: 停止运行中的应用
- **WHEN** 客户端发送停止请求指定应用名称
- **THEN** 系统应发送 `SIGTERM` 信号
- **AND** 系统应更新应用状态为 `STOPPED` (Manual Stop)
- **AND** 系统 MUST NOT 对该应用执行自动重启（自动重启能力在本变更中被移除）

### Requirement: 应用列表 (List Apps)
系统 MUST 能够列出所有受管应用及其当前状态；状态与资源指标 MUST 通过动态查询进程信息得出，而不是依赖 AppCenter 在线连接。

#### Scenario: 获取列表（动态进程查询）
- **WHEN** 客户端请求应用列表
- **THEN** 系统应返回所有应用的名称、PID、状态 (RUNNING/STOPPED/ERROR)、端口及资源使用统计
- **AND** 若应用记录存在 PID，则系统应在返回前按需检查该 PID 对应进程是否存在并采集 CPU/内存/IO
- **AND** 系统 MUST NOT 通过“是否连接 AppCenter”来判断应用是否 RUNNING

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

### Requirement: 应用详情查询 (Get App Details)
系统 MUST 能够查询单个应用的详细信息（通过动态进程查询）。

#### Scenario: 获取应用详情
- **WHEN** 客户端请求应用详情（`GET /apps/<name>`）
- **THEN** 系统应按需动态查询该应用的 PID 存在性与进程信息
- **AND** 返回应用的元数据、状态、资源配置、PID、启动时间、资源使用等
- **AND** 状态映射规则：PID 存在 → RUNNING；PID 不存在且非 STOPPED → EXITED；手动 stop → STOPPED

### Requirement: 状态监控 (Status Monitoring)
系统 MUST 提供应用状态的查询能力，并基于进程信息区分正常退出、异常退出与手动停止；系统 MUST NOT 执行自动保活重启。

#### Scenario: 进程不存在（按需查询）
- **GIVEN** 应用在存储中状态非 `STOPPED` 且记录了 PID
- **WHEN** 用户执行 `glow get app` 或 `glow describe app <name>` 触发状态查询
- **THEN** 若 PID 对应进程不存在，系统应将该应用状态返回为 `ERROR` 或 `EXITED`（实现阶段定义具体枚举）
- **AND** 系统 MUST NOT 自动尝试重启该应用

### Requirement: 应用管理 (App Management)
系统 MUST 使用 "App" (应用) 作为核心资源定义，并提供类 K8s 的 CLI 操作接口。

#### Scenario: 获取应用列表
- **WHEN** 用户执行 `glow get app`
- **THEN** CLI 应从服务端获取应用列表
- **AND** CLI 应列出所有 App 的状态（NAME, STATUS, AGE, CPU, MEM, PID, PORT, DOMAIN）
- **AND** 状态与指标来源于服务端动态进程查询（而非 AppCenter 在线连接）

#### Scenario: 查看应用详情
- **WHEN** 用户执行 `glow describe app <name>`
- **THEN** CLI 应从服务端获取应用详情
- **AND** CLI 应显示指定 App 的详细信息（Events, Config, Resources）
- **AND** 状态与指标来源于服务端动态进程查询（而非 AppCenter 在线连接）

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

### Requirement: 应用元数据登记 (CLI-driven App Metadata Registration)
系统 MUST 提供由 Glow CLI 的 `apply` 触发的应用元数据登记能力，用于持久化保存应用元数据（如 name、command、workingDir、env、port/domain、资源需求等），以支持后续资源供给与配置落盘；该登记 MUST NOT 视为将应用注册为 OS service。

#### Scenario: 注册应用元数据
- **WHEN** 用户通过 Glow CLI 执行 `glow apply -f app.yaml`（`kind: App`）
- **THEN** 系统应创建/更新该应用的元数据记录
- **AND** 后续 `glow get app` 应能展示该应用（即使应用当前未运行）
- **AND** `glow apply` MUST NOT 启动应用（不执行 deploy 动作）
- **AND** 若应用配置有变化（diff 检测），系统应自动重启该应用
- **AND** 若应用配置无变化，系统应保持应用当前状态（不改变运行状态）

