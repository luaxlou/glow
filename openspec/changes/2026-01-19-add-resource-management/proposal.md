# Change: 增加 glow-server 资源管理能力（支持删除/重建已集成资源）

## Why
当前 `glow-server add <resource>` 会把资源管理信息写入系统配置存储，但缺少对已集成资源的“解绑/删除”能力。一旦资源需要重建（例如 MySQL root 密码变更、Redis ACL 变更、Nginx 配置迁移），只能手动清理配置存储，流程不透明且容易误删。

## What Changes
- 新增 `glow-server remove <resource>` 命令组，用于删除已保存的资源集成配置：
  - `glow-server remove mysql`：删除 `mysql_info` 与 `mysql_users`
  - `glow-server remove redis`：删除 `redis_info` 与 `redis_users`
  - `glow-server remove nginx`：删除 `nginx_info`
- 命令具备幂等性：当资源未配置时，执行也应成功并给出提示
- 默认进行交互式二次确认；提供 `--yes` 以跳过确认（便于自动化脚本）

## Impact
- Affected specs: `server-management`
- Affected code:
  - `cmd/glow-server/cmd/*`（新增 remove 命令与资源子命令）
  - `internal/configmanager/*`（提供删除 system_config 的能力）
  - `internal/appcenter/provision.go`（删除后应返回“未配置”类错误信息的行为保持一致）
