## 1. 基础架构与认证
- [x] 1.1 在 `cmd/glow` 中引入 Cobra 并重构基础结构。
- [x] 1.2 实现 `authentication` 规范：`glow auth view/reset` 及隐式交互引导。
- [x] 1.3 在 `pkg/api` 中定义 K8s 风格的资源结构（Deployment, Node, Ingress）。

## 2. 部署管理 (Deployment)
- [x] 2.1 Server: 确保 `/apps/list` 返回完整的 Deployment 状态信息。
- [x] 2.2 Client: 实现 `app-management` 规范：`glow get deploy`。
- [x] 2.3 Client: 实现 `glow describe deploy <name>`。
- [x] 2.4 Client: 实现 `glow logs <name>`。
- [x] 2.5 Client: 实现 `glow start/stop/restart/delete deploy <name>`。

## 3. 应用配置 (App Config)
- [x] 3.1 Server: 确认 `/config/<app>` 接口支持完整的 CRUD 操作。
- [x] 3.2 Client: 实现 `config-management` 规范：`glow config view/apply/edit`。

## 4. 节点与资源 (Node & Resources)
- [x] 4.1 Server: 实现资源注册表 (Resource Registry)，在拨备资源时持久化记录资源元数据（Type, Name, Binding App）。
- [x] 4.2 Server: 实现 `/node/status` 接口，集成 `gopsutil` 获取宿主机指标。
- [x] 4.3 Server: 实现 `/resources/list` 接口，返回注册表中的所有资源。
- [x] 4.4 Client: 实现 `glow get node`，展示节点列表及核心压力指标。
- [x] 4.5 Client: 实现 `glow describe node <name>`，展示系统详情及受管基础设施列表。
- [x] 4.6 Client: 实现 `glow get resources`，跨类型汇总展示所有托管资源。
- [x] 4.7 Client: 实现 `glow describe <name>`，支持不带类型直接查询资源详情。

## 5. Ingress 管理
- [x] 5.1 Server: 增强 `/ingress` 接口，支持独立的 CRUD 操作（不依赖 App 启停）。
- [x] 5.2 Client: 实现 `ingress-automation` 规范：`glow get/create/delete ingress`。

## 6. 声明式申请 (Manifest)
- [x] 6.1 Server: 增强 Manifest 解析器，支持 `kind: Deployment` 和 `kind: Node`。
- [x] 6.2 Client: 实现 `manifest-application` 规范：`glow apply -f`。
