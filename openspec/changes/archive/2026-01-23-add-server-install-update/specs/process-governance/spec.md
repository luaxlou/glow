## ADDED Requirements

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

