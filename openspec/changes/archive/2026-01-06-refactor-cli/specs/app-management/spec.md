## ADDED Requirements
### Requirement: CLI 交互 (CLI Interaction)
系统 MUST 提供基于动词-名词结构的命令行工具，不再支持 declarative apply 模式。

#### Scenario: 启动应用命令
- **WHEN** 用户执行 `glow app start <name> --command <cmd>`
- **THEN** CLI 应发送启动请求到服务端

#### Scenario: 停止应用命令
- **WHEN** 用户执行 `glow app stop <name>`
- **THEN** CLI 应发送停止请求到服务端

#### Scenario: 删除应用命令
- **WHEN** 用户执行 `glow app delete <name>`
- **THEN** CLI 应发送删除请求到服务端
