# 资源供给 (Resource Provisioning)

## Purpose
为应用提供基础设施资源的绑定与配置生成服务，支持 MySQL、Redis 等资源的声明式绑定、凭据管理、连接信息生成，并将资源配置落盘到应用配置文件中。资源绑定由 CLI 驱动，应用运行时无需主动申请资源。
## Requirements
### Requirement: 资源绑定 (Resource Binding)
系统 MUST 支持为应用绑定基础设施资源（MySQL、Redis 等），并生成连接配置。

#### Scenario: 通过 CLI 绑定 MySQL 资源
- **WHEN** 用户在 App manifest 中声明 `spec.resources.mysql[]` 或执行 `glow app add mysql <appName> <dbName>`
- **THEN** CLI 应调用服务端资源绑定 API（`POST /apps/:appName/resources/mysql`）
- **AND** 服务端应创建或复用 MySQL 数据库与用户
- **AND** 服务端应生成 DSN 连接字符串并写入应用配置
- **AND** 操作 MUST 幂等（重复绑定同一资源应返回相同结果）

#### Scenario: 通过 CLI 绑定 Redis 资源
- **WHEN** 用户在 App manifest 中声明 `spec.resources.redis[]` 或执行 `glow app add redis <appName>`
- **THEN** CLI 应调用服务端资源绑定 API（`POST /apps/:appName/resources/redis`）
- **AND** 服务端应生成或复用 Redis 连接信息（addr、password、db）
- **AND** 服务端应将连接信息写入应用配置
- **AND** 操作 MUST 幂等

### Requirement: 资源鉴权与凭据管理 (Resource Authentication)
系统 MUST 安全管理资源的访问凭据，并支持既有资源的访问授权。

#### Scenario: 新建资源时自动生成凭据
- **WHEN** 绑定新资源时（如新建 MySQL 数据库）
- **THEN** 系统应自动生成强随机密码
- **AND** 凭据应安全存储（如加密存储在服务端配置中）
- **AND** 凭据应写入应用配置文件供应用读取

#### Scenario: 访问既有资源时需要鉴权
- **WHEN** 绑定既有资源需要访问凭据（如访问已存在的 MySQL root 账号）
- **THEN** 服务端应返回结构化错误（如 `{"error": "needs_credentials", "message": "需要提供 MySQL root 密码"}`）
- **AND** CLI 应在侧安全读取密码（不回显、不落日志）
- **AND** CLI 应使用凭据重试绑定请求
- **AND** 服务端不应依赖 stdin 交互（daemon 模式不可用）

### Requirement: 资源配置落盘 (Resource Config Materialization)
系统 MUST 将资源绑定结果写入应用可读取的本地配置文件。

#### Scenario: 资源绑定后自动落盘配置
- **WHEN** 资源绑定成功（MySQL DSN、Redis 连接信息等）
- **THEN** CLI 应触发配置落盘（`POST /config/:appName/render`）
- **AND** 服务端应将资源连接信息写入 `<data-dir>/apps/<appName>/<appName>_local_config.json`
- **AND** 配置文件应包含应用所有资源与运行时配置
- **AND** 应用启动时应能从该配置文件读取资源连接信息

#### Scenario: 配置文件格式
- **WHEN** 配置文件落盘
- **THEN** 文件 MUST 为 JSON 格式
- **AND** 应包含 `mysql.dsn`、`redis.addr`、`redis.password`、`redis.db` 等字段
- **AND** 应与 Starter SDK 的配置加载器兼容

### Requirement: 资源解绑 (Resource Unbinding)
系统 MUST 支持解除应用与资源的绑定关系。

#### Scenario: 删除资源绑定
- **WHEN** 用户执行 `glow app remove mysql <appName> <dbName>`
- **THEN** CLI 应调用服务端 API 删除该资源绑定
- **AND** 服务端应从应用配置中移除该资源连接信息
- **AND** CLI 应触发配置重新落盘（`POST /config/:appName/render`）
- **AND** 系统不删除实际的物理资源（如不删除数据库，仅解除绑定）

### Requirement: 资源绑定 HTTP API (Resource Binding API)
系统 MUST 提供 HTTP API 供 CLI 调用资源绑定功能，所有请求 MUST 通过鉴权。

#### Scenario: MySQL 资源绑定 API
- **WHEN** CLI 发送 `POST /apps/:appName/resources/mysql` 请求
- **Body**: `{ "dbName": "mydb", "mode": "create_or_use", "existingPassword": "..." }`
- **THEN** 服务端应创建或使用指定数据库
- **AND** 返回 `{"mysql": {"dsn": "user:pass@tcp(host:port)/dbName"}}`
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

#### Scenario: Redis 资源绑定 API
- **WHEN** CLI 发送 `POST /apps/:appName/resources/redis` 请求
- **Body**: `{ "mode": "create_or_use" }`
- **THEN** 服务端应生成或使用 Redis 连接配置
- **AND** 返回 `{"redis": {"addr": "localhost:6379", "password": "...", "db": 0}}`
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

### Requirement: CLI Apply 命令集成 (CLI Apply Integration)
系统 MUST 通过 `glow apply` 命令统一处理应用元数据登记与资源绑定。

#### Scenario: Apply 时自动绑定资源
- **WHEN** 用户执行 `glow apply -f app.yaml`，manifest 包含 `spec.resources.mysql[]`
- **THEN** CLI 应先登记/更新应用元数据
- **AND** 逐项调用资源绑定 API（如 MySQL、Redis）
- **AND** 所有资源绑定成功后触发配置落盘
- **AND** 输出变更摘要与下一步建议（如 `glow start app <name>`）

#### Scenario: Apply 失败回滚
- **WHEN** 资源绑定失败（如凭据错误、资源不足）
- **THEN** CLI 应停止后续操作
- **AND** 输出明确错误信息与可操作的修复建议
- **AND** 不更新应用元数据（或回滚已部分执行的变更）

### Requirement: 应用资源绑定（MySQL）(App Resource Binding - MySQL)
系统 MUST 提供由 Glow CLI 触发的“为应用绑定 MySQL 资源”的能力，并将结果写入应用配置与本地配置文件，以便应用运行时无需连接 glow-server 即可获得 DSN。

#### Scenario: 绑定新的 MySQL 数据库并生成配置文件
- **GIVEN** glow-server 已集成 MySQL（系统存在 `mysql_info` 或等价配置）
- **WHEN** 用户执行 `glow apply -f app.yaml` 且该文件声明 `kind: App` 并包含 `spec.resources.mysql[].dbName: <db_name>`
- **THEN** 系统应为该应用创建（或复用）名为 `<db_name>` 的数据库与最小可用访问凭据
- **AND** 系统应将 `mysql.dsn` 写入该应用的配置存储
- **AND** 系统应为该应用生成/更新 `<data-dir>/apps/<appName>/<appName>_local_config.json`

#### Scenario: 幂等绑定
- **GIVEN** 应用 `<appName>` 已绑定 MySQL 数据库 `<db_name>`
- **WHEN** 用户再次执行 `glow apply -f app.yaml`（声明相同的 MySQL dbName 需求）
- **THEN** 系统应返回成功
- **AND** 生成的配置文件内容应保持一致（除非凭据/host 发生变化）
- **AND** 若配置有变化且应用正在运行，系统应自动重启应用

#### Scenario: MySQL 未集成时返回可操作错误
- **GIVEN** glow-server 未集成 MySQL
- **WHEN** 用户执行 `glow apply -f app.yaml` 且声明 MySQL dbName 需求
- **THEN** 系统应返回明确错误并提示用户先执行 `glow-server add mysql`（或等价集成流程）

