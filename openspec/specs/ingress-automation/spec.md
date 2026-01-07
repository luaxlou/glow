# 接入自动化 (Ingress Automation)

## Purpose
自动化管理 Nginx 反向代理配置，动态生成虚拟主机（Server Block），实现从外部域名到本地应用端口的流量路由与负载均衡。
## Requirements
### Requirement: Ingress 管理 (Ingress Management)
系统 MUST 提供类 K8s 的 Ingress 资源管理接口。

#### Scenario: 创建 Ingress
- **WHEN** 用户执行 `glow create ingress <name> --domain <domain> --service <app_name>`
- **THEN** CLI 应发送创建请求

#### Scenario: 获取 Ingress 列表
- **WHEN** 用户执行 `glow get ingress`
- **THEN** CLI 应列出 Ingress 资源（NAME, HOST, SERVICE, PORT）

#### Scenario: 删除 Ingress
- **WHEN** 用户执行 `glow delete ingress <name>`
- **THEN** CLI 应发送删除请求

