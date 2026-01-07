## ADDED Requirements
### Requirement: 多环境管理 (Context Management)
系统 MUST 支持多环境配置（Context），允许用户在不同 Glow Server 之间切换。

#### Scenario: 查看环境列表
- **WHEN** 用户执行 `glow context list`
- **THEN** CLI 应列出所有配置的环境，并标记当前使用的环境

#### Scenario: 切换环境
- **WHEN** 用户执行 `glow context use <name>`
- **THEN** CLI 应将当前上下文切换至指定环境

#### Scenario: 添加环境
- **WHEN** 用户执行 `glow context add <name> --url <url> --key <key>`
- **THEN** CLI 应保存新环境配置

#### Scenario: 删除环境
- **WHEN** 用户执行 `glow context delete <name>`
- **THEN** CLI 应删除指定环境配置

## MODIFIED Requirements
### Requirement: 交互式配置引导 (Interactive Config Bootstrap)
系统 MUST 在检测到配置缺失时自动触发交互式引导，创建默认环境。

#### Scenario: 自动引导
- **WHEN** 用户执行任意命令且本地无配置文件
- **THEN** CLI 应提示 "No context found..."
- **AND** CLI 应引导用户输入 URL 和 Key
- **AND** CLI 应将其保存为 `default` context 并设为 current
