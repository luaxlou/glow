## MODIFIED Requirements

### Requirement: 进程监控 (Process Monitoring)
系统 MUST 在用户查询应用状态时按需检查受管进程的存活状态及资源使用情况；系统 MUST NOT 依赖周期性心跳/长连接来判断进程存活。

#### Scenario: 按需采集资源指标
- **WHEN** 用户执行 `glow get app` 或 `glow describe app <name>`
- **THEN** 系统应检查该应用记录的 PID 是否存在
- **AND** 对于存活进程，采集 CPU 使用率、内存占用 (RSS) 及 IO 读写字节数
- **AND** 将采集结果返回给客户端（并可选持久化更新应用状态信息）

### Requirement: 日志轮转 (Log Rotation)
系统 MUST 管理应用产生的日志文件，防止磁盘占满。

#### Scenario: 日志文件过大
- **WHEN** 应用日志文件超过阈值（如 10MB）
- **THEN** 系统应自动轮转日志（保留历史备份，如5个）
- **AND** 应用的标准输出继续写入新的日志文件

## REMOVED Requirements

### Requirement: 自动重启 (Auto Restart)
**Reason**: 为去除 app 与 glow-server 的运行时耦合并降低治理面误判，本变更移除 glow-server 的自动保活重启能力；应用是否需要保活由用户显式操作或外部进程管理器（systemd/容器编排）承担。

#### Scenario: 迁移指引
- **WHEN** 用户需要恢复应用运行
- **THEN** 用户应显式执行 `glow restart app <name>` 或使用外部守护机制确保进程常驻

