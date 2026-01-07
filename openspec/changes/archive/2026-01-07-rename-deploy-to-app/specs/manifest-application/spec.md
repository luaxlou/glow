## MODIFIED Requirements
### Requirement: 声明式资源 (Declarative Resources)
系统 MUST 支持 K8s 风格的 YAML 资源定义，核心资源包括 `Node` 和 `App`。

#### Scenario: 应用 Manifest
- **WHEN** 用户执行 `glow apply -f deploy.yaml`
- **THEN** CLI 应解析 YAML 中的 `kind: App` 和 `kind: Node`
- **AND** CLI 应将资源配置同步至服务端

#### Scenario: 资源定义
- `Node`: 定义宿主机基础设施（数据库服务、端口池等）
- `App`: 定义应用及其配置（Command, Env, Replicas, Ingress）
