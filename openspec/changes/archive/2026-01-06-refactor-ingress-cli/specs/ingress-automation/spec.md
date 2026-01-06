## ADDED Requirements
### Requirement: Ingress 管理 (Ingress Management)
系统必须 (MUST) 提供基于 API 的接口来管理应用的外部访问入口（Ingress），并支持通过客户端 CLI 进行操作。

#### Scenario: 通过 Client 添加 Ingress
- **WHEN** 用户执行 `glow ingress apply --app <name> --domain <domain>`
- **THEN** 客户端向 Server 发送配置请求
- **AND** Server 验证应用状态并生成 Nginx 配置
- **AND** Server 重载 Nginx 并返回成功响应

#### Scenario: 通过 Client 移除 Ingress
- **WHEN** 用户执行 `glow ingress delete --app <name>`
- **THEN** 客户端向 Server 发送删除请求
- **AND** Server 删除 Nginx 配置并重载
- **AND** 返回操作结果

#### Scenario: 通过 Client 列出 Ingress
- **WHEN** 用户执行 `glow ingress list`
- **THEN** 客户端请求 Server 获取列表
- **AND** Server 返回所有已配置 Ingress 的应用及域名

## REMOVED Requirements
### Requirement: 生成 Nginx 配置 (Generate Config)
**Reason**: Replaced by explicit `Ingress Management` API.

### Requirement: 清理配置 (Cleanup Config)
**Reason**: Merged into `Ingress Management` API.
