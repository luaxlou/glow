# Change: 在 app.yaml 中支持声明应用配置，并新增 glow-server info 命令

## Why
当前 Glow 存在以下问题：
1. **配置管理分散**：应用配置需要通过 `glow config set` 命令单独管理，无法通过 app.yaml 统一声明，不利于版本控制和部署一致性。
2. **服务器信息不透明**：缺少查看 glow-server 自身信息的命令（如安装路径、日志路径、配置路径等），不便于运维和故障排查。

## What Changes
- 在 `app.yaml` 的 `spec.config` 字段中支持声明应用配置（如 mysql.dsn、redis.addr 等）
- **配置完全采用声明式管理**，通过 app.yaml 声明，`glow apply` 应用
- **移除 `glow config` 命令**，所有配置变更通过修改 app.yaml 并重新 apply 来实现
- 更新 `glow apply` 逻辑，读取 `spec.config` 并生成 local_config.json
- 新增 `glow server info` 命令，显示 glow-server 的运行信息：
  - 安装路径（二进制文件位置）
  - 数据目录（data_dir）
  - 日志目录
  - 配置文件路径
  - 进程 PID
  - 运行状态
  - 版本信息

## Impact
- Affected specs:
  - `config-management` - 新增通过 app.yaml 声明配置的能力
  - `server-management` - 新增 server info 功能
- Affected code:
  - CLI: `cmd/glow/apply.go`（支持 spec.config 解析并写入配置文件）
  - CLI: ~~`cmd/glow/config.go`~~（已移除，配置完全通过 app.yaml 管理）
  - CLI: `cmd/glow/server.go`（新增 info 子命令）
  - Server API: `internal/apiserver/server.go`（新增 info 端点）
  - Types: `pkg/api/types.go`（AppSpec 增加 Config 字段，新增 ServerInfo 结构体）
