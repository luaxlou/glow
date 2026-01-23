# 配置管理规范变更

## ADDED Requirements

### Requirement: App.yaml 配置声明 (YAML Config Declaration)
系统 MUST 支持在 `app.yaml` 的 `spec.config` 字段中声明应用配置。

#### Scenario: 通过 app.yaml 声明应用配置
- **WHEN** 用户在 `app.yaml` 中声明 `spec.config` 字段，例如：
  ```yaml
  spec:
    config:
      app:
        timeout: 30
        maxConnections: 100
      mysql:
        dsn: "user:pass@tcp(localhost:3306)/mydb"
        charset: utf8mb4
  ```
- **THEN** CLI 应在 `glow apply` 时读取该配置
- **AND** 将该配置直接写入应用配置文件

#### Scenario: 空配置声明
- **WHEN** `app.yaml` 未声明 `spec.config` 字段或该字段为空
- **THEN** CLI 应生成空配置文件
- **AND** 不影响 `glow apply` 的其他功能

### Requirement: 配置导出 (Config Export)
系统 MUST 支持导出应用当前配置为 YAML 格式，便于迁移到 `app.yaml`。

#### Scenario: 导出应用配置
- **WHEN** 用户执行 `glow config export <app>`
- **THEN** CLI 应读取 `<data-dir>/apps/<app>/<app>_local_config.json`
- **AND** 将配置转换为 YAML 格式输出
- **AND** 可直接粘贴到 `app.yaml` 的 `spec.config` 字段

#### Scenario: 导出完整配置
- **WHEN** 导出应用配置时
- **THEN** CLI 应导出配置文件中的所有字段
- **AND** 以 YAML 格式输出，可直接复制到 app.yaml

## MODIFIED Requirements

### Requirement: 配置文件生成 (Config File Generation)
系统 MUST 支持将 `spec.config` 直接写入应用本地配置文件。

#### Scenario: CLI 触发配置落盘
- **WHEN** Glow CLI 在 `glow apply` 过程中触发配置落盘
- **THEN** CLI 应读取 `spec.config` 字段（如有）
- **AND** 将配置直接序列化为 JSON 格式
- **AND** 写入 `<data-dir>/apps/<appName>/<appName>_local_config.json`
- **AND** 显示写入路径和配置摘要

#### Scenario: 配置落盘目录不存在时自动创建
- **WHEN** 配置落盘目录 `<data-dir>/apps/<appName>` 不存在
- **THEN** CLI 应自动创建该目录（包括必要的父目录）
- **AND** 确保目录权限正确（755）

#### Scenario: spec.config 为空时的处理
- **WHEN** `spec.config` 字段为空或未声明
- **THEN** CLI 应生成空的配置文件（`{}`）
- **AND** 显示警告："未声明配置，生成空配置文件"

### Requirement: CLI 配置管理 (CLI Config Management)
系统 MUST 提供命令行工具来对应用配置进行增删改查操作。

#### Scenario: 设置配置项
- **WHEN** 用户执行 `glow config set <app> <key> <value>`
- **THEN** CLI 应读取 `<data-dir>/apps/<app>/<app>_local_config.json`
- **AND** 将 key-value 合并到配置中
- **AND** 写回配置文件
- **AND** 显示提示："配置已更新，建议同步到 app.yaml 的 spec.config 中"

#### Scenario: 获取配置项
- **WHEN** 用户执行 `glow config get <app> <key>`
- **THEN** CLI 应读取 `<data-dir>/apps/<app>/<app>_local_config.json`
- **AND** 返回该配置项的值

#### Scenario: 列出配置
- **WHEN** 用户执行 `glow config list <app>`
- **THEN** CLI 应读取 `<data-dir>/apps/<app>/<app>_local_config.json`
- **AND** 以表格或 JSON 形式显示该应用的所有配置

#### Scenario: 配置文件不存在时的处理
- **WHEN** 执行配置命令但配置文件不存在
- **THEN** CLI 应显示错误："配置文件不存在，请先执行 'glow apply -f app.yaml'"
