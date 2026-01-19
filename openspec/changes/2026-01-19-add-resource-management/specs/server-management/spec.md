## ADDED Requirements
### Requirement: 删除已集成资源 (Remove Integrated Resource)
系统 MUST 提供命令行能力以删除已保存的资源集成信息，用于支持资源重建与重新集成。

#### Scenario: 删除 MySQL 集成
- **WHEN** 用户运行 `glow-server remove mysql`
- **THEN** 系统应从配置存储中删除 `mysql_info`
- **AND** 系统应从配置存储中删除 `mysql_users`（若存在）
- **AND** CLI 应输出成功结果

#### Scenario: 删除 Redis 集成
- **WHEN** 用户运行 `glow-server remove redis`
- **THEN** 系统应从配置存储中删除 `redis_info`
- **AND** 系统应从配置存储中删除 `redis_users`（若存在）
- **AND** CLI 应输出成功结果

#### Scenario: 删除 Nginx 集成
- **WHEN** 用户运行 `glow-server remove nginx`
- **THEN** 系统应从配置存储中删除 `nginx_info`
- **AND** CLI 应输出成功结果

#### Scenario: 删除未集成资源 (幂等性)
- **WHEN** 用户运行 `glow-server remove <resource>` 且系统未配置该资源
- **THEN** CLI 应返回成功
- **AND** CLI 应提示未发现可删除的资源配置（或等价信息）

#### Scenario: 跳过确认 (Non-Interactive)
- **WHEN** 用户运行 `glow-server remove <resource> --yes`
- **THEN** 系统应跳过交互确认并直接执行删除
