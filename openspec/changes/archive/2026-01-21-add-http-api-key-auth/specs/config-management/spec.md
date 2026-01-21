## MODIFIED Requirements

### Requirement: 获取配置 (Get Config)
应用启动时 MUST 能从 Server 获取其运行时配置。

#### Scenario: 应用通过 AppCenter 获取初始配置
- **WHEN** 应用（通过 Starter）与 AppCenter 建立 TCP 连接并发送启动请求（例如 `ActionAppStart`）
- **THEN** 系统应返回该应用专属的 JSON 配置对象（若存在）

#### Scenario: 通过 HTTP 管理面读取配置（需要鉴权）
- **WHEN** Glow CLI 对 `/config/<appName>` 发起 GET 请求
- **THEN** 系统应返回该应用专属的 JSON 配置对象
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

### Requirement: 更新配置 (Update Config)
系统 MUST 支持动态更新应用配置。

#### Scenario: API 更新配置（需要鉴权）
- **WHEN** Glow CLI 以 PUT 方式提交新的 JSON 配置到 `/config/<appName>`
- **THEN** 系统应持久化存储该配置
- **AND** 配置变更应对下一次获取生效
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

