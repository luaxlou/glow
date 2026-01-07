## REMOVED Requirements
### Requirement: 部署管理 (Deployment Management)
**Reason**: Renamed to App Management.

## ADDED Requirements
### Requirement: 应用管理 (App Management)
系统 MUST 使用 "App" (应用) 作为核心资源定义，并提供类 K8s 的 CLI 操作接口。

#### Scenario: 获取应用列表
- **WHEN** 用户执行 `glow get app`
- **THEN** CLI 应列出所有 App 的状态（NAME, STATUS, RESTARTS, AGE, CPU, MEM, PID, PORT, DOMAIN）

#### Scenario: 查看应用详情
- **WHEN** 用户执行 `glow describe app <name>`
- **THEN** CLI 应显示指定 App 的详细信息（Events, Config, Resources）

#### Scenario: 删除应用
- **WHEN** 用户执行 `glow delete app <name>`
- **THEN** CLI 应向服务端发送删除请求

#### Scenario: 重启应用
- **WHEN** 用户执行 `glow restart app <name>`
- **THEN** CLI 应触发滚动重启或原地重启

#### Scenario: 停止应用 (Stop)
- **WHEN** 用户执行 `glow stop app <name>`
- **THEN** CLI 应停止该 App 的运行进程

#### Scenario: 启动应用 (Start)
- **WHEN** 用户执行 `glow start app <name>`
- **THEN** CLI 应启动该 App 的运行进程

#### Scenario: 查看日志
- **WHEN** 用户执行 `glow logs <name>`
- **THEN** CLI 应获取并打印日志（支持 -f 实时流式传输）
