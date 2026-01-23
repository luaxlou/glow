## ADDED Requirements

### Requirement: Curl 一键安装 (One-line Install)
系统 MUST 提供可通过 `curl` 一键执行的安装入口，用于在目标主机上安装 `glow-server` 二进制并完成基础初始化。

#### Scenario: 系统级一键安装
- **WHEN** 用户运行 `curl -fsSL <install-script-url> | sudo bash`
- **THEN** 脚本应安装 `glow-server` 与 `glow` 到系统 PATH（例如 `/usr/local/bin/glow-server` 与 `/usr/local/bin/glow`）
- **AND** 脚本应创建并初始化固定的配置目录与数据目录（例如 `/etc/glow-server/` 与 `/var/lib/glow-server/`）
- **AND** 脚本应以幂等方式安装并启动系统服务（Systemd 或 Launchd）
- **AND** 脚本 MUST 在安装过程中执行 `glow-server keygen`（生成或复用 API Key）

#### Scenario: 安装过程不依赖 Go 工具链
- **WHEN** 用户在一台未安装 Go（`go` 命令不存在）的主机上执行一键安装脚本
- **THEN** 安装脚本 MUST 仍可成功完成安装
- **AND** 脚本 MUST 通过下载预编译二进制归档并解压安装来完成安装
- **AND** 脚本 MUST NOT 使用 `go install` 或任何需要 Go 工具链的安装方式

#### Scenario: 重复安装前先备份
- **WHEN** 用户在已安装过 `glow-server`/`glow` 的主机上再次执行一键安装脚本
- **THEN** 脚本 MUST 在覆盖现有文件前先执行备份
- **AND** 备份 MUST 至少包含既有的 `glow-server` 与 `glow` 二进制文件
- **AND** 若脚本将修改或覆盖服务定义文件与环境文件，则这些文件也 MUST 被备份

#### Scenario: 重装检查并复用既有配置与数据库
- **WHEN** 用户再次执行一键安装脚本（重装场景）
- **THEN** 脚本 MUST 检查本地配置目录与数据库文件是否已存在（例如配置目录与 `<data-dir>/glow.db`）
- **AND** 若已存在，脚本 MUST NOT 重新创建或覆盖这些文件/目录
- **AND** 若用户希望全新初始化，用户需要手动删除配置目录与数据库文件后再执行安装

#### Scenario: 重装提示（告知已复用状态）
- **WHEN** 安装脚本检测到本地已存在配置目录或数据库文件
- **THEN** 安装脚本 MUST 在输出中明确提示“已检测到并复用既有配置/数据库”
- **AND** MUST 告知用户若需重置需手动删除哪些路径（例如配置目录与 `<data-dir>/glow.db`）

### Requirement: Curl 一键卸载 (One-line Uninstall)
系统 MUST 提供可通过 `curl` 一键执行的卸载入口，用于移除 `glow-server`/`glow` 二进制与系统服务注册。

#### Scenario: 系统级一键卸载
- **WHEN** 用户运行 `curl -fsSL <uninstall-script-url> | sudo bash`
- **THEN** 脚本 MUST 停止并禁用 `glow-server` 系统服务（若存在）
- **AND** 脚本 MUST 移除 `glow-server` 与 `glow` 二进制文件（从系统 PATH 位置）
- **AND** 脚本 MUST 移除服务定义文件（Systemd unit 或 Launchd plist）

#### Scenario: 卸载不删除配置与数据库
- **WHEN** 用户执行一键卸载脚本
- **THEN** 脚本 MUST NOT 删除既有的配置文件目录（例如 `/etc/glow-server/` 或等价路径）
- **AND** 脚本 MUST NOT 删除既有的数据库文件（例如 `<data-dir>/glow.db` 或等价路径）

### Requirement: 本地安装脚本 (Local Install Script)
系统 MUST 提供一个面向本地使用的一键安装入口（例如 `install-local.sh`），用于安装 `glow` 与 `glow-server`，但不将 `glow-server` 安装为常驻服务。

#### Scenario: install-local 安装内容
- **WHEN** 用户运行 `curl -fsSL <install-local-script-url> | bash`
- **THEN** 脚本 MUST 安装 `glow` 与 `glow-server` 两个可执行文件
- **AND** 脚本 MUST NOT 注册/启用/启动任何系统服务（systemd/launchd）

#### Scenario: install-local 重装复用既有数据
- **WHEN** 用户再次执行 `install-local.sh`（重装场景）
- **THEN** 脚本 MUST 检查本地数据目录与数据库文件是否已存在（例如 `<data-dir>/glow.db`）
- **AND** 若已存在，脚本 MUST NOT 重新创建或覆盖该数据库
- **AND** 脚本 MUST 在输出中提示“已复用既有数据库/配置”，并提示如何手动删除以重置

#### Scenario: install-local 放置到 PATH（按环境选择）
- **WHEN** 脚本在本地安装模式下写入可执行文件
- **THEN** 脚本 MUST 根据系统环境选择合适的安装目录（例如优先使用可写的系统级目录，否则使用用户级目录）
- **AND** 若用户环境的 PATH 未包含该目录，脚本 MUST 输出明确指引用于添加 PATH

#### Scenario: install-local 不依赖 Go 工具链
- **WHEN** 用户在一台未安装 Go（`go` 命令不存在）的主机上执行 `install-local.sh`
- **THEN** 脚本 MUST 通过下载预编译二进制归档并解压安装来完成安装
- **AND** 脚本 MUST 使用发布的校验文件进行 sha256 校验
- **AND** 脚本 MUST NOT 使用 `go install` 或任何需要 Go 工具链的安装方式

#### Scenario: 安装后即可使用（类似 MySQL 成对安装）
- **WHEN** 安装脚本成功完成
- **THEN** 用户应可直接执行 `glow-server` 与 `glow` 两条命令（无需手动编译或移动二进制）
- **AND** `glow-server` 服务应处于可用状态（已启动或可被服务管理器启动）

#### Scenario: 安装前的安全校验
- **WHEN** 安装脚本从 release 下载二进制归档文件
- **THEN** 脚本 MUST 使用发布的校验文件（如 `checksums.txt`）进行 sha256 校验
- **AND** 校验失败 MUST 终止安装并返回非 0 退出码

#### Scenario: 写入 glow 默认连接配置（非交互）
- **WHEN** 安装脚本已获得 glow-server 的 Server URL 与 API Key（例如本机 `http://localhost:32102` + keygen 输出）
- **THEN** 脚本 SHOULD 为“执行安装的目标用户”（例如 `SUDO_USER`）写入 `glow` 的默认配置文件（例如 `~/.glow.json`）
- **AND** 配置 MUST 包含一个可用的默认 context（例如 `default` 或 `local`）
- **AND** 使用户首次执行 `glow get ...` 不需要进入交互式引导即可连通本机 `glow-server`

### Requirement: 固定目录约定 (Fixed Layout)
系统 MUST 采用固定且可预测的目录结构存放配置、数据与日志，避免依赖当前工作目录或不确定的 home 目录。

#### Scenario: 系统级默认目录
- **WHEN** `glow-server` 以系统服务方式运行
- **THEN** 默认数据目录 SHOULD 为 `/var/lib/glow-server`
- **AND** 默认配置目录 SHOULD 为 `/etc/glow-server`
- **AND** 默认日志目录 SHOULD 为 `<data-dir>/logs`

#### Scenario: macOS 本地开发默认目录（用户级）
- **WHEN** 用户在 macOS 上以本地开发方式安装并运行 `glow-server`（用户级服务）
- **THEN** 默认数据目录 SHOULD 位于用户目录下（例如 `~/Library/Application Support/glow-server`）
- **AND** 默认日志目录 SHOULD 为 `<data-dir>/logs`
- **AND** 默认配置（如有）SHOULD 位于用户目录下（避免写入 `/etc`）

#### Scenario: 显式指定 data-dir
- **WHEN** 用户通过 CLI 参数或环境文件指定 `data-dir`
- **THEN** `glow-server` MUST 使用该目录作为运行态根目录（apps/nginx/logs/db 等均从此派生）

### Requirement: 服务环境文件 (Service Environment File)
系统级服务 MUST 支持从独立环境文件读取配置（例如端口、data-dir），以便在不修改 service 文件的情况下调整参数。

#### Scenario: 修改端口无需重装服务
- **WHEN** 用户更新环境文件（例如 `/etc/glow-server/glow-server.env`）中的端口配置
- **THEN** 用户仅需重启服务即可使配置生效
- **AND** 无需重新生成/安装 service 文件

### Requirement: macOS 本地开发支持 (macOS Local Dev Support)
系统 MUST 支持在 macOS 上以“用户级（无需 sudo）”方式安装与运行，以便本地开发。该模式 MUST NOT 安装为常驻系统服务。

#### Scenario: macOS 用户级安装（不常驻）
- **WHEN** 用户在 macOS 上执行本地开发一键安装脚本（不使用 sudo）
- **THEN** 脚本 MUST 将 `glow` 与 `glow-server` 安装到用户级可执行目录（例如 `~/.local/bin/` 或等价目录）
- **AND** 脚本 MUST NOT 注册 LaunchAgent 或任何开机自启/常驻服务
- **AND** 脚本 MUST 输出明确指引：
  - 如何确保 PATH 可找到安装目录（若该目录未在 PATH 中）
  - 如何以前台方式启动服务（例如 `glow-server serve`）

#### Scenario: macOS 用户级卸载
- **WHEN** 用户在 macOS 上执行一键卸载脚本（不使用 sudo）
- **THEN** 脚本 MUST 移除用户级安装的二进制文件
- **AND** 脚本 MUST NOT 删除用户数据目录中的数据库与配置（例如 `~/Library/Application Support/glow-server/`）

### Requirement: 本地开发一键安装 (Local Dev One-line Install)
系统 MUST 提供一个面向本地开发的一键安装入口，用于安装 `glow` 与 `glow-server`，但不将 `glow-server` 安装为常驻服务。

#### Scenario: macOS 本地开发一键安装（无需 sudo，不常驻）
- **WHEN** 用户运行 `curl -fsSL <install-local-script-url> | bash`
- **THEN** 脚本 MUST 通过下载预编译二进制归档并解压安装来完成安装
- **AND** 脚本 MUST 使用发布的校验文件进行 sha256 校验
- **AND** 脚本 MUST NOT 使用 `go install`
- **AND** 脚本 MUST NOT 注册或启动常驻服务
- **AND** 脚本 MUST 执行 `glow-server keygen` 以生成或复用 API Key
- **AND** 脚本 MUST 为当前用户写入 `glow` 默认 context（指向 `http://localhost:32102`）

#### Scenario: 本地开发启动方式
- **WHEN** 用户希望在本地使用 Glow（开发模式）
- **THEN** 用户应以前台方式执行 `glow-server serve` 启动服务
- **AND** `glow` 命令应可使用已写入的默认 context 访问本机服务

## MODIFIED Requirements

### Requirement: 安装入口 (Installer Entry)
系统 MUST 提供可重复执行且幂等的安装入口，用于完成 Glow 运行环境的初始化配置（安装二进制、生成/复用密钥、安装并启用服务）。
系统 MUST 以“一键安装脚本”为默认与推荐入口，并且 MUST NOT 提供 `glow-server install` 命令以避免与脚本职责重叠。

#### Scenario: 帮助列表不包含 install 命令
- **WHEN** 用户运行 `glow-server --help`
- **THEN** 帮助列表中不应出现 `install` 命令

#### Scenario: 试图执行 install 命令
- **WHEN** 用户运行 `glow-server install`
- **THEN** CLI MUST 返回错误并退出为非 0 状态码

### Requirement: 密钥生成集成 (Keygen Integration)
安装入口 MUST 集成密钥生成步骤。

#### Scenario: 安装期密钥检查与生成
- **WHEN** 安装脚本执行初始化流程
- **THEN** 脚本 MUST 直接执行 `glow-server keygen` 以生成或复用 API Key
- **AND** 生成的 API Key MUST 可用于后续客户端连接配置与服务鉴权

