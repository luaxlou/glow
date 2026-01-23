## MODIFIED Requirements

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

## ADDED Requirements

### Requirement: 配置分发模型 (Config Distribution Model)
系统 MUST 将应用运行态配置分发模型定义为“控制面生成本地配置文件”，并且 MUST NOT 依赖 AppCenter/TCP 推送机制完成配置分发。

#### Scenario: 配置变更的生效方式
- **WHEN** 用户通过 Glow CLI 更新配置并完成 render
- **THEN** 新配置应在应用下次启动时从本地配置文件生效
- **AND** 系统 MAY 提示用户通过显式重启（如 `glow restart app <appName>`）使变更生效

