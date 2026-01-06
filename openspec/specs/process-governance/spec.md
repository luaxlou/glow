# 进程治理 (Process Governance)

## Purpose
负责运行态进程的深度治理，包括实时监控进程存活状态、采集 CPU/内存/IO 资源指标、执行异常退出的自动保活重启策略以及日志轮转管理。

## Requirements

### Requirement: 进程监控 (Process Monitoring)
系统 MUST 周期性地检查受管进程的存活状态及资源使用情况。

#### Scenario: 采集资源指标
- **WHEN** 监控周期触发（如每5秒）
- **THEN** 系统应检查所有记录的 PID 是否存在
- **AND** 对于存活进程，采集 CPU 使用率、内存占用 (RSS) 及 IO 读写字节数
- **AND** 更新应用状态信息

### Requirement: 自动重启 (Auto Restart)
系统 MUST 具备应用崩溃后的自动恢复能力。

#### Scenario: 进程意外退出
- **WHEN** 监控检测到应用状态非 `STOPPED` 但进程不存在
- **THEN** 系统应标记应用为 `EXITED` 或 `ERROR`
- **AND** 如果开启了 `AutoRestart` 且重试次数未超限（如5次）
- **THEN** 系统应自动尝试重新启动该应用
- **AND** 增加重试计数器

### Requirement: 日志轮转 (Log Rotation)
系统 MUST 管理应用产生的日志文件，防止磁盘占满。

#### Scenario: 日志文件过大
- **WHEN** 应用日志文件超过阈值（如 10MB）
- **THEN** 系统应自动轮转日志（保留历史备份，如5个）
- **AND** 应用的标准输出继续写入新的日志文件
