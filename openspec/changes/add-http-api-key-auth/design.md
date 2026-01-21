## Context
Glow 的 API Key 设计定位为“只给 Glow CLI 使用”的管理面认证凭据。客户端已按 `Authorization: Bearer <apiKey>` 发送，但服务端需要在 Gin HTTP 层落地统一校验，以避免管理端口暴露时被未授权调用。

## Goals / Non-Goals
- Goals:
  - 在 `glow-server` HTTP 管理面实现统一、可复用的 API Key 鉴权
  - 明确鉴权覆盖范围（除 `/health` 外）
  - 返回明确的 HTTP 状态码与错误信息（401/403/500）
- Non-Goals:
  - 不改变 AppCenter(TCP) 的认证边界（应用侧不使用该 Key）
  - 不改造 Provision 的交互式行为

## Decisions
- Decision: 采用 Gin Middleware + Route Group
  - Why: 最小侵入、覆盖面清晰、不会遗漏新增路由
- Decision: Bearer Scheme
  - 规则：`Authorization: Bearer <token>`（大小写不敏感）
- Decision: API Key 存储来源
  - 从 SQLite `system_config.api_key` 读取期望值；使用常量时间比较避免时序侧信道

## Risks / Trade-offs
- 每次请求读取 `system_config` 会产生额外 DB 访问
  - Mitigation: 可在后续增强中加入内存缓存/热加载，但本变更优先保证正确性与简单性
- 老版本 CLI 若未带 header，将收到 401
  - Mitigation: 文档明确说明；错误信息可提示缺少 Authorization

## Migration Plan
- 部署前确保 `glow-server keygen` 已生成/复用 `api_key`
- 更新客户端配置（已有安装脚本/默认 context 写入时应包含 key）
- 回滚：移除 middleware 或仅将服务绑定到 localhost（不在本变更范围）

## Open Questions
- 是否需要对 `/health` 以外的只读路由（如 `/apps/list`）提供“可选匿名只读”模式？（默认否）

