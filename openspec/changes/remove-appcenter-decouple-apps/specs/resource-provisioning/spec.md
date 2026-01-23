## ADDED Requirements

### Requirement: 应用资源绑定（MySQL）(App Resource Binding - MySQL)
系统 MUST 提供由 Glow CLI 触发的“为应用绑定 MySQL 资源”的能力，并将结果写入应用配置与本地配置文件，以便应用运行时无需连接 glow-server 即可获得 DSN。

#### Scenario: 绑定新的 MySQL 数据库并生成配置文件
- **GIVEN** glow-server 已集成 MySQL（系统存在 `mysql_info` 或等价配置）
- **WHEN** 用户执行 `glow apply -f app.yaml` 且该文件声明 `kind: App` 并包含 `spec.resources.mysql[].dbName: <db_name>`
- **THEN** 系统应为该应用创建（或复用）名为 `<db_name>` 的数据库与最小可用访问凭据
- **AND** 系统应将 `mysql.dsn` 写入该应用的配置存储
- **AND** 系统应为该应用生成/更新 `<data-dir>/apps/<appName>/<appName>_local_config.json`

#### Scenario: 幂等绑定
- **GIVEN** 应用 `<appName>` 已绑定 MySQL 数据库 `<db_name>`
- **WHEN** 用户再次执行 `glow apply -f app.yaml`（声明相同的 MySQL dbName 需求）
- **THEN** 系统应返回成功
- **AND** 生成的配置文件内容应保持一致（除非凭据/host 发生变化）
- **AND** 若配置有变化且应用正在运行，系统应自动重启应用

#### Scenario: MySQL 未集成时返回可操作错误
- **GIVEN** glow-server 未集成 MySQL
- **WHEN** 用户执行 `glow apply -f app.yaml` 且声明 MySQL dbName 需求
- **THEN** 系统应返回明确错误并提示用户先执行 `glow-server add mysql`（或等价集成流程）

