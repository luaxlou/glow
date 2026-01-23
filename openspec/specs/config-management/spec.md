# 配置管理 (Config Management)

## Purpose
提供集中式的应用配置存储与分发服务，支持配置的版本化管理、配置文件渲染与落盘，以及宿主基础设施配置的统一管理。应用通过读取本地配置文件获取配置，无需主动连接 glow-server。
## Requirements
### Requirement: 获取配置 (Get Config)
应用启动时 MUST 能从“本地配置文件”获取其运行时配置；该配置文件由 Glow CLI（经由 glow-server 管理面能力）生成与更新，应用运行时 MUST NOT 依赖与 glow-server 的网络连接。

#### Scenario: 通过本地配置文件获取初始配置
- **GIVEN** glow-server 已为应用 `<appName>` 生成配置文件 `<data-dir>/apps/<appName>/<appName>_local_config.json`
- **WHEN** 应用启动（通过 Starter 加载配置）
- **THEN** Starter 应从该本地配置文件读取并加载 JSON 配置对象
- **AND** 应用启动流程 MUST NOT 需要主动连接 glow-server（无 AppCenter/TCP 依赖）

#### Scenario: 通过 HTTP 管理面读取配置（需要鉴权）
- **WHEN** Glow CLI 对 `/config/<appName>` 发起 GET 请求
- **THEN** 系统应返回该应用专属的 JSON 配置对象
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

### Requirement: 更新配置 (Update Config)
系统 MUST 支持通过管理面更新应用配置，并支持将最新配置落盘为应用可读取的本地配置文件。

#### Scenario: API 更新配置（需要鉴权）
- **WHEN** Glow CLI 以 PUT 方式提交新的 JSON 配置到 `/config/<appName>`
- **THEN** 系统应持久化存储该配置
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

#### Scenario: 生成/更新本地配置文件（需要鉴权）
- **WHEN** Glow CLI 对 `/config/<appName>/render` 发起 POST 请求
- **THEN** 系统应将服务端存储的应用配置写入 `<data-dir>/apps/<appName>/<appName>_local_config.json`
- **AND** 写入成功后应返回落盘路径与写入结果（如字节数/哈希）
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

### Requirement: 渲染与落盘配置 (Render and Materialize Config)
系统 MUST 支持将服务端存储的配置渲染为本地配置文件。

#### Scenario: CLI 触发配置落盘（需要鉴权）
- **WHEN** Glow CLI 以 POST 方式请求 `/config/<appName>/render`
- **THEN** 系统应读取服务端存储的该应用配置
- **AND** 将配置写入 `<data-dir>/apps/<appName>/<appName>_local_config.json`
- **AND** 返回写入路径、字节数、可选的配置哈希值
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

#### Scenario: 配置落盘目录不存在时自动创建
- **WHEN** 配置落盘目录 `<data-dir>/apps/<appName>` 不存在
- **THEN** 系统应自动创建该目录（包括必要的父目录）
- **AND** 确保目录权限正确（应用可读）

### Requirement: 宿主配置 (Host Config)
系统 MUST 管理宿主机的基础设施配置（如本地 MySQL/Redis 连接信息）。

#### Scenario: 设置 Host 配置
- **WHEN** 客户端提交 Host Manifest
- **THEN** 系统应解析并保存服务定义（如 MySQL root 账号、端口）供 Provisioner 使用

### Requirement: 声明式配置管理 (Declarative Config Management)
系统 MUST 支持通过 app.yaml 声明应用配置，配置变更通过修改 YAML 并重新 apply 实现。

#### Scenario: 通过 app.yaml 声明配置
- **GIVEN** 用户在 app.yaml 的 `spec.config` 字段中声明配置
- **WHEN** 用户执行 `glow apply -f app.yaml`
- **THEN** 系统应解析 `spec.config` 并保存到服务端存储
- **AND** 自动触发配置渲染，生成 `<data-dir>/apps/<appName>/<appName>_local_config.json`

#### Scenario: 配置文件自动生成
- **WHEN** 执行 `glow apply` 时，app.yaml 中包含 `spec.config` 字段
- **THEN** 系统应自动调用 `/config/<appName>/render` 端点
- **AND** 将配置写入本地 JSON 文件
- **AND** 返回配置文件路径、大小和哈希值

#### Scenario: 配置变更流程
- **WHEN** 用户需要修改应用配置
- **THEN** 用户 MUST 编辑 app.yaml 文件
- **AND** 执行 `glow apply -f app.yaml` 应用新配置
- **AND** 系统 MUST NOT 提供命令行式的配置修改工具（如 `glow config set`）

### Requirement: 配置分发模型 (Config Distribution Model)
系统 MUST 将应用运行态配置分发模型定义为"控制面生成本地配置文件"，并且 MUST NOT 依赖 AppCenter/TCP 推送机制完成配置分发。

#### Scenario: 配置变更的生效方式
- **WHEN** 用户通过 `glow apply` 更新配置
- **THEN** 新配置应自动渲染到本地配置文件
- **AND** 新配置应在应用下次启动时生效
- **AND** 如果应用正在运行，系统 MAY 提示用户重启应用

#### Scenario: 配置的单一真相来源
- **WHEN** 应用配置需要变更
- **THEN** app.yaml 是配置的唯一声明来源
- **AND** 所有配置变更 MUST 通过修改 app.yaml 并执行 `glow apply` 实现
- **AND** 系统 MUST NOT 支持运行时通过命令行直接修改配置
- **AND** 生成的本地 JSON 文件是只读的，仅供应用读取

