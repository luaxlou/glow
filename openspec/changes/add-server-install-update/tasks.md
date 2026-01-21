## 1. Specification
- [x] 1.1 增加 `system-initialization` 的 curl 安装脚本（同时安装 `glow-server` 与 `glow`）、固定目录与服务环境文件要求
- [x] 1.2 明确安装脚本不得使用 `go install`，必须仅通过下载 release 二进制归档完成安装（不依赖 Go 工具链）
- [x] 1.3 增加一键卸载脚本要求：卸载移除二进制与服务，但不删除配置与数据库；重装前先备份
- [x] 1.4 增加 `build-release` 的 release 产物命名与校验文件要求（sha256，覆盖 `glow-server` 与 `glow`）
- [x] 1.5 增加 `server-management` 的 `update`/`version` 命令要求（含回滚与服务重启行为）
- [x] 1.6 增加 `process-governance` 的 logs 目录与自清理要求（按年龄+按总量）
- [x] 1.7 增加 `authentication` 的"安装期写入默认 context（非交互）"要求，确保 `glow` 安装后可直接连接本机服务
- [x] 1.8 简化脚本：合并 `install-glow.sh` 与 `install-local-dev.sh` 为 `install-local.sh`（本地安装不常驻、不注册服务，安装 `glow` + `glow-server`）
- [x] 1.9 增加 macOS 本地开发（用户级安装/不常驻/默认目录）支持要求
- [x] 1.10 增加"本地开发一键安装脚本（不常驻）"要求（macOS 为主）
- [x] 1.11 增加"重装复用既有配置与数据库（不覆盖）+ 安装提示告知如何重置"的要求（install/install-local）

## 2. Implementation (after approval)
- [x] 2.1 增加 `glow-server version` 输出（version/commit/buildDate 由 ldflags 注入）
- [x] 2.2 实现 `glow-server update`：检查最新版本、下载、校验、原子替换、可回滚
- [x] 2.3 `serve` 支持显式 `--data-dir`（并确保 DB/日志/应用目录均从 data-dir 派生）
- [x] 2.4 更新 service 模板：固定 WorkingDirectory、EnvironmentFile、ExecStart 参数
- [x] 2.5 增加日志清理器：周期执行，清理 `*.log.N*`，支持 maxAgeDays 与 maxTotalMB
- [x] 2.6 提供安装脚本：平台识别、下载 release、校验、落盘、安装服务并启动，并在安装期执行 `glow-server keygen` 且写入 `glow` 默认 context
- [x] 2.7 移除 `glow-server install` 命令（删除/下线 `cmd/glow-server/cmd/install.go` 并更新相关文档与测试）
- [x] 2.8（可选）增加 `glow version` 输出与自我更新命令 `glow update`（若决定对齐"客户端也可自更新"的体验）
- [x] 2.9 提供一键卸载脚本：停止/禁用服务、移除二进制与服务定义文件，但不删除配置与数据库；重装前自动备份逻辑落地
- [x] 2.10 统一本地脚本为 `install-local.sh`：平台识别、下载 release、校验、落盘到 PATH（不注册/不启动常驻服务）
- [x] 2.11 `install-local.sh` 支持 macOS 本地开发：用户级落盘、keygen、写入默认 context（不常驻）
- [x] 2.12 落地重装幂等逻辑：检测并复用既有 config/db（存在则不覆盖），并在输出中提示检测结果与手动重置路径

## 3. Documentation
- [x] 3.1 更新 `docs/server_manual.md`：一键安装、服务管理、更新与日志策略
- [x] 3.2 更新 `README.md`：推荐安装路径与 quickstart
- [x] 3.3 补充"重装复用与如何重置（手动删除哪些路径）"的说明（面向 install/install-local）

## 4. Validation
- [x] 4.1 `openspec validate add-server-install-update --strict`

