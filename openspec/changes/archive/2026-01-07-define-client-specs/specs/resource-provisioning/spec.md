## ADDED Requirements
### Requirement: 托管资源管理 (Managed Resource Management)
系统 MUST 提供对托管资源（如 MySQL, Redis）的统一 CLI 管理接口。

#### Scenario: 获取资源列表
- **WHEN** 用户执行 `glow get <resource_type>` (支持 `mysql`, `redis`)
- **THEN** CLI 应列出该类型的所有活跃资源实例及其归属应用

#### Scenario: 获取所有托管资源
- **WHEN** 用户执行 `glow get resources`
- **THEN** CLI 应汇总并列出所有类型的托管资源（MySQL, Redis 等）

#### Scenario: 查看资源详情
- **WHEN** 用户执行 `glow describe <name>`
- **THEN** CLI 应自动识别资源类型并显示详细信息（如 DSN、连接参数、状态）

### Requirement: 节点管理 (Node Management)
系统 MUST 提供 `node` 资源的管理与监控能力，支持查看宿主机状态及基础设施。

#### Scenario: 获取节点列表
- **WHEN** 用户执行 `glow get node`
- **THEN** CLI 应列出节点列表及核心指标（NAME, STATUS, CPU%, MEM%）

#### Scenario: 查看节点详情
- **WHEN** 用户执行 `glow describe node <name>`
- **THEN** CLI 应显示节点的基础信息（OS, Arch, Kernel）
- **AND** CLI 应显示实时系统负载（CPU, Memory, IO, Disk Usage）
- **AND** CLI 应列出该节点上已注册的基础设施资源（如 MySQL, Redis, MQ 等服务及其端口/状态）