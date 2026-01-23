## ADDED Requirements

### Requirement: 安装期写入默认连接 (Installer Seeded Auth)
系统 MUST 允许安装流程在非交互场景下为 `glow` 客户端写入默认连接配置，使安装后可直接使用客户端命令访问本机 `glow-server`。

#### Scenario: 由安装脚本写入 default context
- **WHEN** 安装脚本已获取 Glow Server URL 与 API Key
- **THEN** 脚本 MUST 将其写入 `glow` 客户端配置文件（例如 `~/.glow.json`）
- **AND** 配置 MUST 包含一个 `default`（或 `local`）context 并设为 current
- **AND** 用户首次运行 `glow get apps` 时不应进入交互式引导

