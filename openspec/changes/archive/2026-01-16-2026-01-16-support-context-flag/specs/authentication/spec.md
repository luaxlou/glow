## MODIFIED Requirements
### Requirement: 多环境管理 (Context Management)
系统 MUST 支持多环境配置（Context），允许用户在不同 Glow Server 之间切换或临时指定。

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

#### Scenario: 临时指定环境 (Command Flag)
- **WHEN** 用户执行命令时附加 `--context <name>` 参数（如 `glow get app --context prod`）
- **THEN** CLI 应使用指定 `<name>` 环境的连接信息执行该次命令
- **AND** 此次执行不应改变全局默认的 Current Context
