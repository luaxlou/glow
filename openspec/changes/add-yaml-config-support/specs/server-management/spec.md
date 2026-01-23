# 服务器管理规范变更

## ADDED Requirements

### Requirement: 服务器信息查询 (Server Info Query)
系统 MUST 支持查询 glow-server 的运行信息和配置。

#### Scenario: 查询 glow-server 信息
- **WHEN** 管理员执行 `glow server info`
- **THEN** CLI 应调用服务端 API（`GET /server/info`）
- **AND** 服务端返回 glow-server 的运行信息
- **AND** CLI 以易读格式显示信息
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）

#### Scenario: 服务器信息包含的字段
- **WHEN** 查询服务器信息时
- **THEN** 响应应包含以下字段：
  - `version`：glow-server 版本号
  - `binary_path`：glow-server 二进制文件路径
  - `data_dir`：数据目录路径
  - `log_dir`：日志目录路径
  - `config_file`：配置文件路径
  - `pid`：进程 PID
  - `status`：运行状态（running、stopped 等）
  - `uptime`：运行时长（秒）
  - `start_time`：启动时间（Unix 时间戳）

#### Scenario: 服务器信息输出格式
- **WHEN** CLI 显示服务器信息时
- **THEN** 应以表格或键值对形式展示
- **AND** 支持 `--json` 参数输出 JSON 格式（便于脚本解析）

#### Scenario: 服务器不可达时的处理
- **WHEN** glow-server 未运行或网络不可达
- **THEN** CLI 应显示错误："无法连接到 glow-server，请检查服务是否运行"
- **AND** 返回非零退出码

### Requirement: 服务器信息 HTTP API (Server Info API)
系统 MUST 提供 HTTP API 供 CLI 调用服务器信息查询功能。

#### Scenario: 服务器信息 API
- **WHEN** CLI 发送 `GET /server/info` 请求
- **THEN** 服务端应返回服务器信息
- **AND** 响应格式：
  ```json
  {
    "version": "1.0.0",
    "binary_path": "/usr/local/bin/glow-server",
    "data_dir": "/var/lib/glow-server",
    "log_dir": "/var/log/glow-server",
    "config_file": "/etc/glow-server/config.json",
    "pid": 12345,
    "status": "running",
    "uptime": 86400,
    "start_time": 1706000000
  }
  ```
- **AND** 该请求 MUST 通过 HTTP 管理面鉴权（见 `authentication` 规范）
