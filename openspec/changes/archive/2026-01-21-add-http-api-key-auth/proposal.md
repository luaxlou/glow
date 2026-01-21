# Change: 为 glow-server HTTP 管理面增加 API Key 鉴权（Gin Middleware）

## Why
当前 Glow CLI（`glow`）在调用 `glow-server` 的 HTTP 管理 API 时已经发送 `Authorization: Bearer <apiKey>`，但服务端在 Gin HTTP 层缺少统一的鉴权校验，导致只要管理端口可达就存在未授权访问风险（应用上传/启动/停止、配置读写、Ingress 修改等）。

本变更将把 “API Key 只给 Glow 客户端使用” 的边界落地到服务端 HTTP 层，实现默认安全的管理面访问控制。

## What Changes
- 在 `glow-server` 的 Gin HTTP 层新增统一中间件：校验 `Authorization: Bearer <api_key>`。
- API Key 来源：读取 SQLite `system_config` 中的 `api_key`（由 `glow-server keygen` 或安装流程生成/复用）。
- 保护范围：除 `/health` 外的所有 HTTP 管理 API（`/apps/*`、`/config/*`、`/ingress/*`、`/node/*`）。
- 失败行为：
  - 缺少/格式错误的 `Authorization`：返回 HTTP 401。
  - API Key 不匹配：返回 HTTP 403。
  - 服务端未配置 `api_key`：返回 HTTP 500（配置错误）。

## Impact
- Affected specs:
  - `openspec/specs/authentication/spec.md`（新增“HTTP 管理面鉴权”要求）
  - `openspec/specs/config-management/spec.md`（明确 `/config/*` 属于管理面并要求鉴权；应用配置获取以 AppCenter 为主）
- Affected code:
  - `internal/apiserver/`（新增 middleware，并将路由挂到受保护 group）
  - （可选）`docs/server_manual.md`（补充鉴权说明与示例）

## Non-Goals
- 不为 AppCenter TCP 通道引入 API Key 鉴权（应用侧不使用该 Key）。
- 不在本变更中改造 Provision 的交互式流程（“Provision 去交互/可自动化”不在本次范围）。

