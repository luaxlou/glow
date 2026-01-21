## ADDED Requirements

### Requirement: HTTP 管理面鉴权 (HTTP Management API Authentication)
glow-server 的 HTTP 管理 API MUST 要求客户端提供有效的 API Key，以限制仅 Glow CLI 可访问受保护的管理能力。

#### Scenario: 鉴权成功
- **WHEN** 客户端对受保护的 HTTP 路由发起请求并携带 `Authorization: Bearer <api_key>`
- **AND** `<api_key>` 与服务端持久化的 `system_config.api_key` 一致
- **THEN** 服务端 MUST 允许请求继续处理并返回业务响应

#### Scenario: 缺少或错误的 Authorization 头
- **WHEN** 客户端对受保护的 HTTP 路由发起请求但未携带 `Authorization` 头
- **OR** `Authorization` 头不是 `Bearer <token>` 形式
- **THEN** 服务端 MUST 返回 HTTP 401

#### Scenario: API Key 不匹配
- **WHEN** 客户端对受保护的 HTTP 路由发起请求并携带 `Authorization: Bearer <api_key>`
- **AND** `<api_key>` 与服务端持久化的 `system_config.api_key` 不一致
- **THEN** 服务端 MUST 返回 HTTP 403

#### Scenario: 健康检查不需要鉴权
- **WHEN** 客户端请求 `GET /health`
- **THEN** 服务端 MUST 返回 HTTP 200
- **AND** 该请求 MUST NOT 需要提供 API Key

