# 认证与连接 (Authentication & Connection)

## Purpose
管理客户端与 Glow Server 之间的安全认证与连接配置，确保 CLI 能够合法访问服务端资源。
## Requirements
### Requirement: 认证管理 (Auth Management)
系统 MUST 提供管理连接认证信息（Server URL & API Key）的命令，命令组命名为 `auth`。

#### Scenario: 查看认证信息
- **WHEN** 用户执行 `glow auth view`
- **THEN** CLI 应显示脱敏后的连接与认证信息

#### Scenario: 重置认证
- **WHEN** 用户执行 `glow auth reset`
- **THEN** CLI 应清除本地认证缓存
- **AND** 立即进入交互式引导流程

### Requirement: 交互式配置引导 (Interactive Config Bootstrap)
系统 MUST 在检测到配置缺失时自动触发交互式引导，创建默认环境。

#### Scenario: 自动引导
- **WHEN** 用户执行任意命令且本地无配置文件
- **THEN** CLI 应提示 "No context found..."
- **AND** CLI 应引导用户输入 URL 和 Key
- **AND** CLI 应将其保存为 `default` context 并设为 current

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

### Requirement: HTTP 管理面鉴权 (HTTP Management API Authentication)
glow-server 的 HTTP 管理 API MUST 要求客户端提供有效的 API Key，以限制仅 Glow CLI 可访问受保护的管理能力。

#### Scenario: 鉴权成功
- **WHEN** 客户端对受保护的 HTTP 路由发起请求并携带 `Authorization: Bearer <api_key>`
- **AND** `<api_key>` 与服务端持久化的 `system_config.api_key` 一致
- **THEN** 服务端 MUST 允许请求继续处理并返回业务响应

#### Scenario: 缺少或错误的 Authorization 头
- **WHEN** 客户端对受保护的 HTTP 路由发起请求但未携带 `Authorization` 头
- **OR** `Authorization` 头不是 `Bearer <token>` 形式
- **THEN** 服务端 MUST 返回 HTTP 401

#### Scenario: API Key 不匹配
- **WHEN** 客户端对受保护的 HTTP 路由发起请求并携带 `Authorization: Bearer <api_key>`
- **AND** `<api_key>` 与服务端持久化的 `system_config.api_key` 不一致
- **THEN** 服务端 MUST 返回 HTTP 403

#### Scenario: 健康检查不需要鉴权
- **WHEN** 客户端请求 `GET /health`
- **THEN** 服务端 MUST 返回 HTTP 200
- **AND** 该请求 MUST NOT 需要提供 API Key

### Requirement: 安装期写入默认连接 (Installer Seeded Auth)
系统 MUST 允许安装流程在非交互场景下为 `glow` 客户端写入默认连接配置，使安装后可直接使用客户端命令访问本机 `glow-server`。

#### Scenario: 由安装脚本写入 default context
- **WHEN** 安装脚本已获取 Glow Server URL 与 API Key
- **THEN** 脚本 MUST 将其写入 `glow` 客户端配置文件（例如 `~/.glow.json`）
- **AND** 配置 MUST 包含一个 `default`（或 `local`）context 并设为 current
- **AND** 用户首次运行 `glow get apps` 时不应进入交互式引导

