# 接入自动化 (Ingress Automation)

## ADDED Requirements

### Requirement: 服务器自托管 (Server Self-Hosting)
系统 MUST 支持通过 Nginx 将 glow-server 自身的 API 服务暴露到外部域名。

#### Scenario: 配置反向代理
- **WHEN** 用户在 `install` 流程中选择配置 glow-server 域名（例如 `glow.example.com`）
- **THEN** 系统应生成 Nginx 配置文件，将该域名流量转发到 `glow-server` 的监听端口
- **AND** 重载 Nginx 使配置生效
- **AND** 后续 `info` 命令应显示该外部访问地址
