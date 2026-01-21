## 1. Implementation
- [x] 1.1 在 `internal/apiserver/` 增加 Gin 鉴权中间件（Bearer API Key），并从 `system_config.api_key` 读取期望值
- [x] 1.2 将除 `/health` 外的 HTTP 管理路由迁移到受保护的 route group（apps/config/ingress/node）
- [x] 1.3 为鉴权行为补充最小测试（至少覆盖：missing header=401、invalid=403、valid=200、/health bypass）
- [x] 1.4 更新 `docs/server_manual.md`：说明 HTTP 管理面鉴权头部与常见错误码

## 2. Validation
- [x] 2.1 `openspec validate add-http-api-key-auth --strict`
- [x] 2.2 `go test ./...`（或至少确保 `glow-server` 可编译启动）

