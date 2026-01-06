## MODIFIED Requirements
### Requirement: 启动应用 (Start App)
系统 MUST 能够启动一个新的应用实例。如果应用已在运行，操作 MUST 幂等。

#### Scenario: 成功启动新应用
- **WHEN** 客户端发送合法的启动请求（包含名称、命令、参数等）
- **THEN** 系统应检查应用是否已在运行
- **AND** 如果未运行，系统应准备运行环境并启动进程
- **AND** 如果应用配置中未指定端口，系统应分配一个空闲端口
- **AND** 系统应返回成功响应

#### Scenario: 启动已运行的应用 (幂等性)
- **WHEN** 客户端请求启动一个状态为 RUNNING 的应用
- **THEN** 系统应直接返回成功响应，不执行任何操作

### Requirement: 停止应用 (Stop App)
系统 MUST 能够优雅地停止正在运行的应用，并标记为手动停止状态。

#### Scenario: 停止运行中的应用
- **WHEN** 客户端发送停止请求指定应用名称
- **THEN** 系统应发送 `SIGTERM` 信号
- **AND** 系统应更新应用状态为 `STOPPED` (Manual Stop)
- **AND** Watchdog 不应自动重启状态为 `STOPPED` 的应用

## ADDED Requirements
### Requirement: 删除应用 (Delete App)
系统 MUST 能够彻底删除应用及其相关资源。

#### Scenario: 删除应用
- **WHEN** 客户端发送删除请求指定应用名称
- **THEN** 系统应停止该应用（如果正在运行）
- **AND** 系统应清理该应用的运行目录、日志和配置信息
- **AND** 系统应从应用列表中移除该应用

### Requirement: 重启应用 (Restart App)
系统 MUST 能够重启应用。

#### Scenario: 重启应用
- **WHEN** 客户端发送重启请求
- **THEN** 系统应先停止应用
- **AND** 系统应重新启动应用

### Requirement: 状态监控 (Status Monitoring)
系统 MUST 监控应用进程状态，并区分正常退出和异常退出。

#### Scenario: 异常退出
- **WHEN** 应用进程退出 (无论 Exit Code 为何) 且状态非 `STOPPED`
- **THEN** 系统应更新状态为 `ERROR`
- **AND** Watchdog 应尝试自动重启
