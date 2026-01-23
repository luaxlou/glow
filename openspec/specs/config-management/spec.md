# 配置管理 (Config Management)

## Purpose
提供集中式的应用配置存储与分发服务，支持配置的版本化管理、配置文件渲染与落盘，以及宿主基础设施配置的统一管理。应用通过读取本地配置文件获取配置，无需主动连接 glow-server。
## Requirements
### Requirement: 读取配置 (Read Config)
应用启动时 MUST 能从本地配置文件读取其运行时配置。

#### Scenario: 应用从本地配置文件读取初始配置
- **WHEN** 应用（通过 Starter）启动时读取本地配置文件 `<data-dir>/apps/<appName>/<appName>_local_config.json`
- **THEN** 系统应返回该应用专属的 JSON 配置对象（若存在）
- **AND** 配置文件包含应用运行所需的所有配置（包括资源绑定结果如 MySQL DSN、Redis 连接信息等）

#### Scenario: 通过 HTTP 管理面读取配置（需要鉴权）
- **WHEN** Glow CLI 对 `/config/<appName>` 发起 GET 请求
- **THEN** 系统应返回该应用专属的 JSON 配置对象
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

### Requirement: 更新配置 (Update Config)
系统 MUST 支持更新应用配置。

#### Scenario: API 更新配置（需要鉴权）
- **WHEN** Glow CLI 以 PUT 方式提交新的 JSON 配置到 `/config/<appName>`
- **THEN** 系统应持久化存储该配置
- **AND** 配置变更应对下一次读取生效
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

### Requirement: 渲染与落盘配置 (Render and Materialize Config)
系统 MUST 支持将服务端存储的配置渲染为本地配置文件。

#### Scenario: CLI 触发配置落盘（需要鉴权）
- **WHEN** Glow CLI 以 POST 方式请求 `/config/<appName>/render`
- **THEN** 系统应读取服务端存储的该应用配置
- **AND** 将配置写入 `<data-dir>/apps/<appName>/<appName>_local_config.json`
- **AND** 返回写入路径、字节数、可选的配置哈希值
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

#### Scenario: 配置落盘目录不存在时自动创建
- **WHEN** 配置落盘目录 `<data-dir>/apps/<appName>` 不存在
- **THEN** 系统应自动创建该目录（包括必要的父目录）
- **AND** 确保目录权限正确（应用可读）

### Requirement: 宿主配置 (Host Config)
系统 MUST 管理宿主机的基础设施配置（如本地 MySQL/Redis 连接信息）。

#### Scenario: 设置 Host 配置
- **WHEN** 客户端提交 Host Manifest
- **THEN** 系统应解析并保存服务定义（如 MySQL root 账号、端口）供 Provisioner 使用

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

