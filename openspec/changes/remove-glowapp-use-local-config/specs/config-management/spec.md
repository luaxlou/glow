# Config Management Spec Delta

## MODIFIED Requirements

### Requirement: 获取配置 (Get Config)
应用启动时 MUST 能从“本地配置文件”获取其运行时配置；该配置文件由 Glow CLI（经由 glow-server 管理面能力）生成与更新。

#### Scenario: 通过本地配置文件获取初始配置
- **GIVEN** glow-server 已为应用 `<appName>` 生成配置文件 `<data-dir>/apps/<appName>/config.json`
- **WHEN** 应用启动
- **THEN** 应用应从该本地配置文件读取并加载 JSON 配置对象
- **AND** 应用可使用 `pkg/glowconfig` SDK 简化读取过程

### Requirement: 渲染与落盘配置 (Render and Materialize Config)
系统 MUST 支持将服务端存储的配置渲染为本地配置文件 `config.json`。

#### Scenario: CLI 触发配置落盘（需要鉴权）
- **WHEN** Glow CLI 以 POST 方式请求 `/config/<appName>/render`
- **THEN** 系统应读取服务端存储的该应用配置
- **AND** 将配置写入 `<data-dir>/apps/<appName>/config.json`
- **AND** 返回写入路径、字节数、可选的配置哈希值

### Requirement: 声明式配置管理 (Declarative Config Management)
系统 MUST 支持通过 app.yaml 声明应用配置，配置变更通过修改 YAML 并重新 apply 实现。

#### Scenario: 通过 app.yaml 声明配置
- **GIVEN** 用户在 app.yaml 的 `spec.config` 字段中声明配置
- **WHEN** 用户执行 `glow apply -f app.yaml`
- **THEN** 系统应解析 `spec.config` 并保存到服务端存储
- **AND** 自动触发配置渲染，生成 `<data-dir>/apps/<appName>/config.json`
