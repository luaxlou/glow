# Change: 为 glow-server 增加 curl 一键安装、开机自启、自我更新与日志清理方案

## Why
当前 `glow-server` 虽已提供 `install` 与 Service 注册能力，但缺少标准化的发布分发（curl 一键安装）、自我更新命令，以及对 logs 目录/日志累计占用的统一治理。同时，`glow`（客户端 CLI）与 `glow-server`（服务端）未形成类似 MySQL 的“成对安装/即装即用”体验：用户安装后往往还需要手动编译、手动 keygen、手动配置连接信息，难以规模化运维与自动化。

## What Changes
- **发布与分发**：提供基于 GitHub Releases 的二进制分发规范，并提供 `curl | sh` 一键安装脚本（支持 Linux/macOS，amd64/arm64）。
- **成对安装（类似 MySQL）**：安装脚本同时安装 `glow-server` 与 `glow`，确保两条命令均在 PATH 中可直接执行。
- **独立安装 glow 客户端**：提供单独的一键安装脚本，仅安装 `glow` 客户端（不安装/不修改 `glow-server` 服务端），便于在仅需要客户端的机器上使用。
- **真正的一键（不依赖 Go）**：安装脚本通过下载预编译二进制归档完成安装，不依赖源码与 Go 工具链（MUST NOT `go install`）。
- **按需支持 macOS 本地开发**：为 macOS 提供用户级安装与 LaunchAgent（无需 sudo）的默认路径与行为，使开发者可在本机快速安装/启动/卸载。
- **按需支持 macOS 本地开发**：为 macOS 提供用户级（无需 sudo）的一键安装与默认目录；本地开发模式不安装为常驻服务，开发者以前台方式启动/停止即可。
- **一键卸载（保留数据）**：提供 `curl | sh` 一键卸载脚本，用于移除二进制与服务注册，但不删除已有配置与数据库。
- **重装前备份**：重复执行安装脚本时，安装脚本 MUST 在覆盖现有安装前先备份（至少包含已安装的二进制与服务配置；配置/数据库如将被修改也需备份）。
- **安装路径与目录约定**：固定系统级安装与运行目录（binary/config/data/logs），避免 `cwd`/`$HOME` 漂移导致服务不稳定。
- **开机自启**：完善 systemd/launchd 服务模板，支持从环境文件读取端口与 data-dir，并幂等安装/重装。
- **自我更新**：新增 `glow-server update`（检查新版本、下载校验、原子替换、可回滚、必要时重启服务）。
- **日志目录与自清理**：补齐 `glow-server` 自身日志目录与托管应用日志目录的统一规范，并增加按“年龄+总量”两阶段清理的要求（轮转之外的兜底）。
- **安装期自动 keygen**：安装过程 MUST 直接执行 `glow-server keygen`，生成/复用 API Key，并将其用于客户端的默认连接配置（使 `glow` 安装后可直接使用）。
- **移除 `glow-server install` 命令**：为避免与一键安装脚本职责冲突，`glow-server` CLI 不再提供 `install` 命令；安装与初始化由脚本完成。

## Impact
- **Affected specs**:
  - `openspec/specs/system-initialization/spec.md`（安装/服务注册/目录约定/脚本）
  - `openspec/specs/build-release/spec.md`（release 产物命名与校验文件）
  - `openspec/specs/server-management/spec.md`（新增 update/version 等管理命令）
  - `openspec/specs/process-governance/spec.md`（日志目录、轮转与自清理）
- **Affected code (planned)**:
  - `cmd/glow-server/cmd/service.go`（服务模板与固定目录/环境文件）
  - `cmd/glow-server/cmd/install.go`（install 幂等与非交互支持）
  - `cmd/glow-server/cmd/`（新增 `update.go`/`version.go`）
  - `internal/manager/`（日志清理器实现，统一 logs 目录约定）
  - `scripts/install-glow-server.sh`（新增）

## Non-Goals
- 不在本变更中引入容器化（Docker/K8s）安装方式。
- 不强制绑定某个包管理器（Homebrew/apt/yum）；脚本以“下载 release 产物”为主。
- 不在本变更中实现日志压缩（可作为后续增强）。

