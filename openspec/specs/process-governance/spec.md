# 进程治理 (Process Governance)

## Purpose
负责运行态进程的信息采集，包括按需查询进程存活状态、采集 CPU/内存/IO 资源指标以及日志轮转管理。系统不执行自动重启，由用户显式触发或外部守护进程管理。
## Requirements
### Requirement: 进程监控 (Process Monitoring)
系统 MUST 在用户查询应用状态时按需检查受管进程的存活状态及资源使用情况；系统 MUST NOT 依赖周期性心跳/长连接来判断进程存活。

#### Scenario: 按需采集资源指标
- **WHEN** 用户执行 `glow get app` 或 `glow describe app <name>`
- **THEN** 系统应检查该应用记录的 PID 是否存在
- **AND** 对于存活进程，采集 CPU 使用率、内存占用 (RSS) 及 IO 读写字节数
- **AND** 将采集结果返回给客户端（并可选持久化更新应用状态信息）

### Requirement: 重启应用 (Restart App)
系统 MUST 支持用户显式触发应用重启。

#### Scenario: 用户手动重启应用
- **WHEN** 用户执行 `glow restart app <name>`
- **THEN** 系统应停止该应用（如果正在运行）
- **AND** 系统应重新启动该应用
- **AND** 系统不自动执行重启，仅响应用户显式请求

### Requirement: 日志轮转 (Log Rotation)
系统 MUST 管理应用产生的日志文件，防止磁盘占满。

#### Scenario: 日志文件过大
- **WHEN** 应用日志文件超过阈值（如 10MB）
- **THEN** 系统应自动轮转日志（保留历史备份，如5个）
- **AND** 应用的标准输出继续写入新的日志文件

### Requirement: 日志目录规范 (Log Directory Layout)
系统 MUST 为 `glow-server` 自身日志与托管应用日志提供统一且可预测的目录布局。

#### Scenario: glow-server 自身日志目录
- **WHEN** `glow-server` 以 data-dir 运行
- **THEN** `glow-server` 自身日志 MUST 写入 `<data-dir>/logs/glow-server.log`（或等价路径）
- **AND** 该目录 MUST 在启动时自动创建（如不存在）

#### Scenario: 托管应用日志目录
- **WHEN** `glow-server` 托管并启动应用 `<app-name>`
- **THEN** 应用标准输出/错误 MUST 写入 `<data-dir>/apps/<app-name>/logs/<app-name>.log`
- **AND** 该目录 MUST 在启动时自动创建（如不存在）

### Requirement: 日志自我清理 (Log Retention Cleanup)
系统 MUST 提供日志自我清理能力，防止日志长期累计导致磁盘占满。该能力 MUST 作为轮转之外的兜底机制存在。

#### Scenario: 仅清理历史文件
- **WHEN** 日志清理任务运行
- **THEN** 系统 MUST 仅删除历史轮转文件（例如 `*.log.N` 或 `*.log.N.gz`）
- **AND** MUST NOT 删除正在写入的当前日志文件（例如 `*.log`）

#### Scenario: 按年龄清理
- **WHEN** 历史日志文件的最后修改时间超过保留天数阈值（例如 14 天）
- **THEN** 系统 MUST 删除该历史日志文件

#### Scenario: 按总量兜底清理
- **WHEN** 指定日志范围（如 `<data-dir>/logs` 与 `<data-dir>/apps/*/logs`）的历史日志总大小超过阈值（例如 1024MB）
- **THEN** 系统 MUST 按“最旧优先”删除历史日志文件，直到总大小低于阈值

#### Scenario: 周期性执行
- **WHEN** `glow-server serve` 启动完成
- **THEN** 系统 SHOULD 周期性执行日志清理（例如每 1 小时一次）

