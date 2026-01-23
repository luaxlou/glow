## ADDED Requirements

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

