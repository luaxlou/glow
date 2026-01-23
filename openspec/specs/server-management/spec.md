# server-management Specification

## Purpose
TBD - created by archiving change 2026-01-07-redefine-glow-server-cli. Update Purpose after archive.
## Requirements
### Requirement: 服务器信息查看 (Server Info)
系统 MUST 提供命令以展示 glow-server 的运行信息和配置状态。

#### Scenario: 查看服务器基本信息
- **WHEN** 用户运行 `glow server info`
- **THEN** CLI 应显示服务器进程 PID
- **AND** 显示数据目录路径（data_dir）
- **AND** 显示日志目录路径（log_dir）
- **AND** 显示配置文件路径（config_path）
- **AND** 显示服务器版本信息
- **AND** 显示服务器运行时长（uptime）

#### Scenario: 以 JSON 格式查看服务器信息
- **WHEN** 用户运行 `glow server info --json`
- **THEN** CLI 应以 JSON 格式输出服务器信息
- **AND** JSON 应包含所有上述字段（pid, data_dir, log_dir, config_path, version, uptime）

#### Scenario: 服务器版本获取
- **WHEN** 查询服务器版本时
- **THEN** 系统 SHOULD 尝试从 git tag 获取版本信息
- **AND** 如果无法获取 git tag，系统 MUST 返回 fallback 版本格式（如 "dev (linux/amd64)"）

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

### Requirement: 自我更新 (Self Update)
系统 MUST 提供自我更新命令以升级 `glow-server` 可执行文件到最新或指定版本。

#### Scenario: 检查更新
- **WHEN** 用户运行 `glow-server update --check`
- **THEN** CLI MUST 检查是否存在更新版本
- **AND** 输出当前版本与可更新版本（若无更新则明确说明）

#### Scenario: 更新到最新版本并原子替换
- **WHEN** 用户运行 `glow-server update`（默认更新到最新版本）
- **THEN** CLI MUST 下载匹配当前 OS/Arch 的 release 产物
- **AND** MUST 进行 sha256 校验
- **AND** MUST 以原子方式替换当前 `glow-server` 二进制（避免产生半写入状态）
- **AND** 在替换前 SHOULD 备份旧版本以支持回滚

#### Scenario: 更新失败可回滚
- **WHEN** 更新流程在校验失败或替换失败时中断
- **THEN** CLI MUST 保持现有可执行文件可用
- **AND** 若已生成备份文件，CLI SHOULD 提供 `glow-server update --rollback` 以恢复

#### Scenario: 更新后重启服务
- **WHEN** `glow-server` 以系统服务方式运行且用户有权限管理服务
- **THEN** 更新完成后 CLI SHOULD 自动重启服务使新版本生效

### Requirement: 版本查看 (Version)
系统 MUST 提供 `version` 命令，便于运维确认运行版本。

#### Scenario: 输出版本信息
- **WHEN** 用户运行 `glow-server version`
- **THEN** CLI MUST 输出 `version`、`commit`、`buildDate`

