# system-initialization Specification

## Purpose
定义 glow-server 的安装、初始化和卸载规范，包括一键安装脚本、目录约定和服务注册。

## Requirements

### Requirement: 一键安装 (One-Click Install)
系统 MUST 提供基于 curl 的 `curl | sh` 一键安装脚本，支持 Linux/macOS (amd64/arm64)。

#### Scenario: 从 GitHub Releases 安装
- **WHEN** 用户执行一键安装脚本
- **THEN** 脚本 MUST 从 GitHub Releases 下载预编译二进制文件
- **AND** 脚本 MUST NOT 使用 `go install`（不依赖 Go 工具链）
- **AND** 脚本 MUST 同时安装 `glow-server` 和 `glow` 两个二进制文件
- **AND** 脚本 MUST 验证下载文件的 SHA256 校验和
- **AND** 安装路径 MUST 固定为 `/usr/local/bin`
- **AND** 数据目录 MUST 固定为 `/var/lib/glow-server`

#### Scenario: 重装前备份
- **WHEN** 用户在已有安装的情况下重新执行安装脚本
- **THEN** 脚本 MUST 在覆盖现有安装前先备份
- **AND** 备份 MUST 至少包含已安装的二进制文件
- **AND** 备份 SHOULD 包含服务配置文件
- **AND** 配置和数据库文件如将被修改也需备份

#### Scenario: 重装复用既有配置与数据库
- **WHEN** 用户再次执行安装脚本（重装场景）
- **THEN** 脚本 MUST 检查本地配置目录与数据库文件是否已存在（例如配置目录与 `<data-dir>/glow.db`）
- **AND** 若已存在，脚本 MUST NOT 重新创建或覆盖这些文件/目录
- **AND** 若用户希望全新初始化，用户需要手动删除配置目录与数据库文件后再执行安装
- **AND** 安装脚本 MUST 在输出中明确提示"已检测到并复用既有配置/数据库"
- **AND** MUST 告知用户若需重置需手动删除哪些路径（例如配置目录与 `<data-dir>/glow.db`）

#### Scenario: 安装期初始化
- **WHEN** 安装脚本执行
- **THEN** 脚本 MUST 自动执行 `glow-server keygen` 生成或复用 API Key
- **AND** 脚本 MUST 将生成的 API Key 配置为 `glow` CLI 的默认 context
- **AND** 确保 `glow` 安装后可直接连接本机服务

### Requirement: 一键卸载 (One-Click Uninstall)
系统 MUST 提供卸载脚本，用于移除二进制文件与服务注册，但保留配置与数据库。

#### Scenario: 卸载服务
- **WHEN** 用户执行卸载脚本
- **THEN** 脚本 MUST 停止并禁用系统服务
- **AND** 脚本 MUST 移除二进制文件（`glow-server` 和 `glow`）
- **AND** 脚本 MUST 移除服务定义文件（systemd/launchd）
- **AND** 脚本 MUST NOT 删除配置文件和数据库
- **AND** 脚本 MUST 保留 `/var/lib/glow-server` 目录下的所有数据

### Requirement: 目录约定 (Directory Convention)
系统 MUST 使用固定的目录结构以避免路径漂移。

#### Scenario: 目录结构
- **WHEN** glow-server 安装和运行
- **THEN** 二进制文件 MUST 位于 `/usr/local/bin/`
- **AND** 配置文件 MUST 位于 `/var/lib/glow-server/config/`
- **AND** 数据库文件 MUST 位于 `/var/lib/glow-server/db/`
- **AND** 日志文件 MUST 位于 `/var/lib/glow-server/logs/`
- **AND** 应用数据 MUST 位于 `/var/lib/glow-server/apps/`

### Requirement: 安装入口 (Installer Entry)
系统 MUST 提供可重复执行且幂等的安装入口，用于完成 Glow 运行环境的初始化配置（安装二进制、生成/复用密钥、安装并启用服务）。

#### Scenario: 移除 install 命令
- **WHEN** 用户运行 `glow-server --help`
- **THEN** 帮助列表中不应出现 `install` 命令
- **AND** 用户运行 `glow-server install` MUST 返回错误并退出为非 0 状态码

#### Scenario: 幂等安装
- **WHEN** 用户再次执行一键安装脚本
- **THEN** 系统应检测已存在的配置（如密钥已存在、服务已安装）
- **AND** 对于已存在的项，默认跳过或显示当前状态，允许用户选择是否重新配置/覆盖

### Requirement: 服务注册 (Service Registration)
安装脚本 MUST 支持将 glow-server 注册为操作系统服务（Systemd 或 Launchd）。

#### Scenario: 安装服务
- **WHEN** 安装脚本在 Linux 系统上执行
- **THEN** 脚本 MUST 生成 systemd 服务配置文件（例如 `/etc/systemd/system/glow-server.service`）
- **AND** 脚本 MUST 设置服务开机自启并启动服务
- **AND** 脚本 MUST 支持从环境文件读取配置（端口、data-dir 等）

### Requirement: 密钥生成集成 (Keygen Integration)
安装流程 MUST 集成密钥生成步骤。

#### Scenario: 安装期密钥生成
- **WHEN** 安装脚本执行初始化流程
- **THEN** 脚本 MUST 直接执行 `glow-server keygen` 以生成或复用 API Key
- **AND** 生成的 API Key MUST 可用于后续客户端连接配置与服务鉴权

### Requirement: 命令优化 (CLI Optimization)
系统 MUST 提供简洁且符合惯例的命令结构。

#### Scenario: 命令重命名
- **WHEN** 用户启动 API 服务
- **THEN** 应使用 `glow-server serve` 命令（原 `server` 命令被移除或重命名）

#### Scenario: 移除冗余命令
- **WHEN** 用户查看帮助列表
- **THEN** `completion` 命令不应出现

### Requirement: 独立客户端安装 (Client-Only Installation)
系统 MUST 提供独立的安装脚本，用于仅安装 `glow` 客户端，不安装或修改 `glow-server` 服务端。

#### Scenario: 仅安装客户端
- **WHEN** 用户执行客户端安装脚本 `curl -fsSL <client-install-url> | bash`
- **THEN** 脚本 MUST 仅安装 `glow` 可执行文件到用户 PATH
- **AND** 脚本 MUST NOT 安装或修改 `glow-server` 的系统服务
- **AND** 脚本 MUST NOT 创建或修改 `glow-server` 的配置目录与数据库
- **AND** 脚本 MUST 通过下载预编译二进制归档完成安装（不使用 `go install`）

#### Scenario: 客户端安装校验
- **WHEN** 脚本下载二进制归档
- **THEN** 脚本 MUST 使用发布的校验文件进行 SHA256 校验
- **AND** 校验失败 MUST 终止安装并返回非 0 退出码

### Requirement: macOS 本地开发支持 (macOS Local Development Support)
系统 MUST 支持在 macOS 上以用户级（无需 sudo）方式安装与运行，该模式 MUST NOT 安装为常驻服务。

#### Scenario: 用户级安装
- **WHEN** 用户在 macOS 上执行本地开发安装脚本（不使用 sudo）
- **THEN** 脚本 MUST 将 `glow` 与 `glow-server` 安装到用户级可执行目录（例如 `~/.local/bin/`）
- **AND** 脚本 MUST NOT 注册 LaunchAgent 或任何开机自启/常驻服务
- **AND** 脚本 MUST 输出明确指引说明如何启动服务和配置 PATH

#### Scenario: 用户级目录约定
- **WHEN** 用户在 macOS 上以本地开发方式安装并运行
- **THEN** 默认数据目录 SHOULD 位于用户目录下（例如 `~/Library/Application Support/glow-server`）
- **AND** 默认日志目录 SHOULD 为 `<data-dir>/logs`
- **AND** 默认配置 SHOULD 位于用户目录下（避免写入 `/etc`）

#### Scenario: 用户级卸载
- **WHEN** 用户在 macOS 上执行卸载脚本（不使用 sudo）
- **THEN** 脚本 MUST 移除用户级安装的二进制文件
- **AND** 脚本 MUST NOT 删除用户数据目录中的数据库与配置

### Requirement: 本地开发安装脚本 (Local Development Installation)
系统 MUST 提供面向本地开发的一键安装入口，安装 `glow` 与 `glow-server`，但不注册常驻服务。

#### Scenario: 本地开发安装
- **WHEN** 用户运行 `curl -fsSL <local-dev-install-url> | bash`
- **THEN** 脚本 MUST 通过下载预编译二进制归档完成安装
- **AND** 脚本 MUST 使用发布的校验文件进行 SHA256 校验
- **AND** 脚本 MUST NOT 使用 `go install`
- **AND** 脚本 MUST NOT 注册或启动常驻服务
- **AND** 脚本 MUST 执行 `glow-server keygen` 以生成或复用 API Key
- **AND** 脚本 MUST 为当前用户写入 `glow` 默认 context（指向 `http://localhost:32102`）

#### Scenario: 本地开发重装复用既有配置与数据库
- **WHEN** 用户再次执行本地开发安装脚本（重装场景）
- **THEN** 脚本 MUST 检查本地数据目录与数据库文件是否已存在（例如 `<data-dir>/glow.db`）
- **AND** 若已存在，脚本 MUST NOT 重新创建或覆盖该数据库
- **AND** 脚本 MUST 在输出中提示"已复用既有数据库/配置"，并提示如何手动删除以重置

#### Scenario: 本地开发使用
- **WHEN** 用户希望在本地使用 Glow（开发模式）
- **THEN** 用户应以前台方式执行 `glow-server serve` 启动服务
- **AND** `glow` 命令应可使用已写入的默认 context 访问本机服务

