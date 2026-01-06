## ADDED Requirements
### Requirement: CLI 配置管理 (CLI Config Management)
系统 MUST 提供命令行工具来对应用配置进行增删改查操作。

#### Scenario: 设置配置项
- **WHEN** 用户执行 `glow config set <app> <key> <value>`
- **THEN** CLI 应发送更新请求，将该 key-value 合并到应用配置中

#### Scenario: 获取配置项
- **WHEN** 用户执行 `glow config get <app> <key>`
- **THEN** CLI 应返回该配置项的值

#### Scenario: 列出配置
- **WHEN** 用户执行 `glow config list <app>`
- **THEN** CLI 应以表格或 JSON 形式显示该应用的所有配置

#### Scenario: 客户端设置 (Client Setup)
- **WHEN** 用户执行 `glow setup --url <url> --key <key>`
- **THEN** CLI 应保存本地连接信息 (原 `glow config`)
