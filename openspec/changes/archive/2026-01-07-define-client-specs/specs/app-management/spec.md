## REMOVED Requirements
### Requirement: CLI 交互 (CLI Interaction)
**Reason**: Replaced by K8s-style Deployment Management.

## ADDED Requirements
### Requirement: 部署管理 (Deployment Management)
系统 MUST 使用 "Deployment" (部署) 作为应用及其运行配置的资源定义，并提供类 K8s 的 CLI 操作接口。

#### Scenario: 获取部署列表
- **WHEN** 用户执行 `glow get deploy`
- **THEN** CLI 应列出所有 Deployment 的状态（NAME, READY, STATUS, RESTARTS, AGE）

#### Scenario: 查看部署详情
- **WHEN** 用户执行 `glow describe deploy <name>`
- **THEN** CLI 应显示指定 Deployment 的详细信息（Events, Config, Resources）

#### Scenario: 删除部署
- **WHEN** 用户执行 `glow delete deploy <name>`
- **THEN** CLI 应向服务端发送删除请求

#### Scenario: 重启部署
- **WHEN** 用户执行 `glow restart deploy <name>`
- **THEN** CLI 应触发滚动重启或原地重启

#### Scenario: 停止部署 (Scale down)
- **WHEN** 用户执行 `glow stop deploy <name>`
- **THEN** CLI 应停止该 Deployment 的运行进程

#### Scenario: 启动部署 (Scale up)
- **WHEN** 用户执行 `glow start deploy <name>`
- **THEN** CLI 应启动该 Deployment 的运行进程

#### Scenario: 查看日志
- **WHEN** 用户执行 `glow logs <name>`
- **THEN** CLI 应获取并打印日志（支持 -f 实时流式传输）
