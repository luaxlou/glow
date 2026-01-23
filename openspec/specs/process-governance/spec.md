# 进程治理 (Process Governance)

## Purpose
负责运行态进程的信息采集，包括按需查询进程存活状态、采集 CPU/内存/IO 资源指标以及日志轮转管理。系统不执行自动重启，由用户显式触发或外部守护进程管理。

## Requirements

### Requirement: 进程监控 (Process Monitoring)
系统 MUST 支持按需查询受管进程的存活状态及资源使用情况。

#### Scenario: 按需采集资源指标
- **WHEN** 用户请求应用列表或详情（`glow get app` 或 `glow describe app <name>`）
- **THEN** 系统应检查记录的 PID 是否存在
- **AND** 对于存活进程，采集 CPU 使用率、内存占用 (RSS) 及 IO 读写字节数
- **AND** 更新应用状态信息
- **AND** 系统不维护后台监控 ticker，仅按需查询一次

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
