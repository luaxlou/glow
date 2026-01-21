# 配置管理 (Config Management)

## Purpose
提供集中式的应用配置存储与分发服务，支持配置的版本化管理、动态更新推送以及宿主基础设施配置的统一管理，消除本地配置文件的依赖。
## Requirements
### Requirement: 获取配置 (Get Config)
应用启动时 MUST 能从 Server 获取其运行时配置。

#### Scenario: 应用通过 AppCenter 获取初始配置
- **WHEN** 应用（通过 Starter）与 AppCenter 建立 TCP 连接并发送启动请求（例如 `ActionAppStart`）
- **THEN** 系统应返回该应用专属的 JSON 配置对象（若存在）

#### Scenario: 通过 HTTP 管理面读取配置（需要鉴权）
- **WHEN** Glow CLI 对 `/config/<appName>` 发起 GET 请求
- **THEN** 系统应返回该应用专属的 JSON 配置对象
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

### Requirement: 更新配置 (Update Config)
系统 MUST 支持动态更新应用配置。

#### Scenario: API 更新配置（需要鉴权）
- **WHEN** Glow CLI 以 PUT 方式提交新的 JSON 配置到 `/config/<appName>`
- **THEN** 系统应持久化存储该配置
- **AND** 配置变更应对下一次获取生效
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

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

