## ADDED Requirements
### Requirement: 应用配置管理 (App Config Management)
系统 MUST 提供 `config` 命令组用于管理应用的运行时配置。

#### Scenario: 查看应用配置
- **WHEN** 用户执行 `glow config view <app_name>`
- **THEN** CLI 应显示该应用的完整 JSON 配置

#### Scenario: 更新应用配置 (文件)
- **WHEN** 用户执行 `glow config apply <app_name> -f <config.json>`
- **THEN** CLI 应读取文件并上传更新配置

#### Scenario: 编辑应用配置 (交互式)
- **WHEN** 用户执行 `glow config edit <app_name>`
- **THEN** CLI 应获取当前配置并打开系统默认编辑器
- **AND** 用户保存退出后，CLI 应验证并上传新配置
