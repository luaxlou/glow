# 服务器管理 (Server Management)

## ADDED Requirements

### Requirement: 服务器信息查看 (Server Info)
系统 MUST 提供 `info` 命令以展示当前服务器的配置和状态。

#### Scenario: 查看信息
- **WHEN** 用户运行 `glow-server info`
- **THEN** CLI 应显示服务器版本
- **AND** 显示密钥状态（是否存在，路径）
- **AND** 显示已集成的资源列表（MySQL, Redis, Nginx 及其连接状态）
- **AND** 显示服务状态（是否作为 Service 运行，PID 等）
