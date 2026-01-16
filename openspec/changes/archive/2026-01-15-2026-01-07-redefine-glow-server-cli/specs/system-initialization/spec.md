# 系统初始化 (System Initialization)

## ADDED Requirements

### Requirement: 交互式安装 (Interactive Install)
系统 MUST 提供 `install` 命令，引导用户完成服务器的初始化配置，且过程需具备幂等性。

#### Scenario: 首次安装
- **WHEN** 用户运行 `glow-server install` 且系统未配置
- **THEN** CLI 应交互式询问是否生成密钥
- **AND** CLI 应交互式询问是否安装为系统服务
- **AND** CLI 应提供资源（Nginx, MySQL, Redis）的多选列表
- **AND** 根据选择，依次引导配置各资源（如输入 MySQL 密码）
- **AND** 若检测到 Nginx，询问是否为 glow-server 配置反向代理域名

#### Scenario: 幂等执行
- **WHEN** 用户再次运行 `glow-server install`
- **THEN** 系统应检测已存在的配置（如密钥已存在、服务已安装）
- **AND** 对于已存在的项，默认跳过或显示当前状态，允许用户选择是否重新配置/覆盖

### Requirement: 服务注册 (Service Registration)
系统 MUST 支持将自身注册为操作系统服务（Systemd 或 Launchd）。

#### Scenario: 安装服务
- **WHEN** 用户在 `install` 过程中确认安装服务
- **THEN** 系统应根据 OS 生成对应的服务配置文件 (e.g., `/etc/systemd/system/glow-server.service`)
- **AND** 系统应设置服务开机自启并启动服务

### Requirement: 密钥生成集成 (Keygen Integration)
`install` 流程 MUST 集成密钥生成步骤。

#### Scenario: 密钥检查
- **WHEN** `install` 流程启动
- **THEN** 检查密钥文件是否存在
- **IF** 不存在，自动调用 Keygen 逻辑生成
- **IF** 存在，提示用户并保留原密钥（除非用户强制覆盖）

### Requirement: 命令优化 (CLI Optimization)
系统 MUST 提供简洁且符合惯例的命令结构。

#### Scenario: 命令重命名
- **WHEN** 用户启动 API 服务
- **THEN** 应使用 `glow-server serve` 命令（原 `server` 命令被移除或重命名）

#### Scenario: 移除冗余命令
- **WHEN** 用户查看帮助列表
- **THEN** `completion` 命令不应出现
