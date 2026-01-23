## ADDED Requirements

### Requirement: App Manifest Support (glow apply)
系统 MUST 支持通过 `glow apply -f <file>` 应用 `kind: App` 的声明式资源文件，用于登记/更新应用元数据与资源需求。

#### Scenario: Apply App Manifest
- **WHEN** 用户执行 `glow apply -f app.yaml`
- **AND** `app.yaml` 的 `kind` 为 `App`
- **THEN** CLI 应解析该文件并将应用元数据与资源需求同步至服务端

### Requirement: App Apply 文件字段 (App Apply File Fields)
`kind: App` 的 manifest MUST 支持以下字段，以满足“开放端口、执行参数、绑定域名、资源需求”的声明：

#### Scenario: App 字段示例与语义
- **GIVEN** 一个 `kind: App` manifest
- **THEN** manifest 应支持：
  - `metadata.name`: 应用名（必填）
  - `spec.port`: 开放端口（可选；**缺省时视为不开放端口**）
  - `spec.args`: 执行参数数组（可选）
  - `spec.domain`: 绑定域名（可选）
  - `spec.resources.mysql[]`: MySQL 资源需求数组（可选；每项 MUST 仅包含 `dbName`）

#### Scenario: 未指定 port 时不开放端口
- **GIVEN** 一个 `kind: App` manifest 未设置 `spec.port`
- **WHEN** 用户执行 `glow apply -f app.yaml`
- **THEN** 系统 MUST 将该应用视为“不开放端口”
- **AND** 系统 MUST NOT 为该应用分配端口
- **AND** 系统 MUST NOT 注入 `OP_APP_PORT`（或等价端口注入机制）

### Requirement: App 资源需求到配置的映射 (Resources-to-Config Mapping)
当 App manifest 声明资源需求时，系统 MUST 将资源供给结果映射为应用配置，并支持后续将配置落盘为本地配置文件。

#### Scenario: MySQL 资源需求映射到 `mysql.dsn`
- **WHEN** `kind: App` manifest 包含 `spec.resources.mysql[].dbName`
- **THEN** 系统应为该 dbName 创建/复用 MySQL 资源与凭据
- **AND** 系统应将结果写入应用配置中的 `mysql.dsn`
- **AND** 系统应支持通过 `/config/<appName>/render` 将最新配置落盘为 `<appName>_local_config.json`

