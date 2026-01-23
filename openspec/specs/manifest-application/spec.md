# manifest-application Specification

## Purpose
TBD - created by archiving change remove-appcenter-decouple-apps. Update Purpose after archive.
## Requirements
### Requirement: App Manifest Support (glow apply)
系统 MUST 支持通过 `glow apply -f <file>` 应用 `kind: App` 的声明式资源文件，用于登记/更新应用元数据与资源需求。

#### Scenario: Apply App Manifest
- **WHEN** 用户执行 `glow apply -f app.yaml`
- **AND** `app.yaml` 的 `kind` 为 `App`
- **THEN** CLI 应解析该文件并将应用元数据与资源需求同步至服务端

### Requirement: App Apply 文件字段 (App Apply File Fields)
`kind: App` 的 manifest MUST 支持以下字段：

#### Scenario: App 字段示例与语义
- **GIVEN** 一个 `kind: App` manifest
- **THEN** manifest 应支持：
  - `metadata.name`: 应用名（必填）
  - `spec.port`: 开放端口（可选；**缺省时视为不开放端口**）
  - `spec.args`: 执行参数数组（可选）
  - `spec.domain`: 绑定域名（可选）
  - `spec.config`: 应用配置 map（可选；用户可声明所有配置项，包括数据库连接等）

#### Scenario: 声明式应用配置
- **GIVEN** 一个 `kind: App` manifest 包含 `spec.config` 字段
- **WHEN** 用户执行 `glow apply -f app.yaml`
- **THEN** 系统应将 `spec.config` 中的配置保存到服务端存储
- **AND** 系统应自动调用 `/config/<appName>/render` 生成本地配置文件
- **AND** 配置文件 MUST 包含 `spec.config` 中声明的所有配置项
- **AND** 应用可以在 `spec.config` 中声明任意配置，包括 MySQL DSN、Redis addr 等

#### Scenario: 未指定 port 时不开放端口
- **GIVEN** 一个 `kind: App` manifest 未设置 `spec.port`
- **WHEN** 用户执行 `glow apply -f app.yaml`
- **THEN** 系统 MUST 将该应用视为“不开放端口”
- **AND** 系统 MUST NOT 为该应用分配端口
- **AND** 系统 MUST NOT 注入 `OP_APP_PORT`（或等价端口注入机制）

### Requirement: 配置即代码 (Configuration as Code)
系统 MUST 支持将所有应用配置（包括数据库连接、缓存配置等）在 `spec.config` 中声明。

#### Scenario: 声明所有配置
- **GIVEN** 一个 `kind: App` manifest
- **WHEN** 用户在 `spec.config` 中声明所有需要的配置
  - 例如：MySQL DSN、Redis 地址、日志级别等
- **THEN** 系统应将这些配置写入 `<appName>_local_config.json`
- **AND** 应用启动时从本地配置文件读取所有配置
- **AND** 系统 MUST NOT 提供独立的资源绑定机制

