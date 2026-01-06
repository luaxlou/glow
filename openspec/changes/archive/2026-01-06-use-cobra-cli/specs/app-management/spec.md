## MODIFIED Requirements
### Requirement: CLI 交互 (CLI Interaction)
系统 MUST 提供基于 Cobra 框架的命令行工具，支持标准化的子命令、Flag 解析和自动帮助生成。

#### Scenario: 启动应用命令
- **WHEN** 用户执行 `glow app start <name> --command <cmd>`
- **THEN** CLI 应解析参数并发送启动请求到服务端

#### Scenario: 停止应用命令
- **WHEN** 用户执行 `glow app stop <name>`
- **THEN** CLI 应发送停止请求到服务端

#### Scenario: 删除应用命令
- **WHEN** 用户执行 `glow app delete <name>`
- **THEN** CLI 应发送删除请求到服务端

#### Scenario: 获取帮助
- **WHEN** 用户执行 `glow app --help` 或 `glow --help`
- **THEN** CLI 应显示自动生成的帮助信息和子命令列表
