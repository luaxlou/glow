## REMOVED Requirements

### Requirement: 项目结构初始化 (Project Structure Initialization)
**Reason**: `glow init` 命令已移除，改用在线 shell 脚本初始化项目结构。新的实现方式允许用户无需预先安装 glow CLI 即可初始化项目。

**Migration**: 用户应使用 `curl -fsSL https://get.glow.dev/init.sh | bash` 代替 `glow init` 命令。

### Requirement: AI 工具集成 (AI Tool Integration)
**Reason**: `glow init` 命令已移除，AI 工具集成功能迁移到在线脚本。

**Migration**: 在线脚本支持通过 `--skip-ai` 参数控制 AI 工具集成，用户可在执行初始化时选择是否集成。

### Requirement: 部署脚本模板 (Deploy Script Template)
**Reason**: 该需求已合并到新的在线脚本需求中。

**Migration**: 在线脚本会自动生成 deploy.sh 和其他必要的文件模板。

### Requirement: 幂等性 (Idempotency)
**Reason**: 该需求已合并到新的在线脚本需求中。

**Migration**: 在线脚本继承相同的幂等性特性，支持 `--force` 参数覆盖现有文件。

## ADDED Requirements

### Requirement: 在线项目初始化脚本 (Online Project Initialization Script)
系统 MUST 提供可通过 curl 执行的在线 shell 脚本，用于初始化项目结构以支持 Glow 治理。

#### Scenario: 基本初始化流程
- **WHEN** 用户在项目根目录执行 `curl -fsSL <init-url> | bash`
- **THEN** 脚本 MUST 分析当前项目结构
- **AND** 脚本 MUST 检测是否存在 `cmd`、`bin`、`scripts` 目录
- **AND** 若目录不存在，脚本 MUST 创建这些目录
- **AND** 脚本 MUST 在 `scripts` 目录下创建 `deploy.sh` 脚本模板
- **AND** 脚本 MUST 输出清晰的进度信息和完成提示

#### Scenario: 交互式 AI 工具集成
- **WHEN** 用户执行在线初始化脚本且未使用 `--skip-ai` 参数
- **THEN** 脚本 MUST 提示用户是否配置 AI 工具集成
- **AND** 脚本 MUST 支持 Claude Code skills 复制
- **AND** 脚本 MUST 从 glow 安装目录复制相关 skills 到项目 `.claude/skills/`
- **AND** 若用户跳过或使用 `--skip-ai`，脚本 MUST 跳过 AI 工具配置步骤

#### Scenario: 命令行参数支持
- **WHEN** 用户执行 `curl -fsSL <url> | bash -s -- [options]`
- **THEN** 脚本 MUST 支持以下参数：
  - `--skip-ai`: 跳过 AI 工具集成提示
  - `--force`: 强制覆盖已存在的文件
  - `--help`: 显示帮助信息
- **AND** 脚本 MUST 正确解析并应用这些参数

#### Scenario: 幂等性保证
- **WHEN** 用户在已初始化的项目中再次执行在线脚本
- **THEN** 脚本 MUST 检测已存在的文件和目录
- **AND** 脚本 MUST 显示当前状态（哪些文件已存在）
- **AND** 脚本 MUST 默认跳过已存在的文件（不覆盖）
- **AND** 若使用 `--force` 参数，脚本 MUST 覆盖现有文件

#### Scenario: 脚本可访问性
- **WHEN** 用户访问初始化脚本 URL
- **THEN** 脚本 MUST 可通过 HTTPS 访问
- **AND** 推荐的 URL 格式为 `https://get.glow.dev/init.sh` 或 `https://raw.githubusercontent.com/luaxlou/glow/main/scripts/init-project.sh`
- **AND** 脚本 MUST 使用稳定的版本或版本化 URL

#### Scenario: 错误处理
- **WHEN** 脚本执行过程中发生错误
- **THEN** 脚本 MUST 输出清晰的错误信息
- **AND** 脚本 MUST 返回非零退出码
- **AND** 脚本 SHOULD 提供修复建议或相关文档链接

#### Scenario: 平台兼容性
- **WHEN** 用户在 Linux 或 macOS 系统上执行脚本
- **THEN** 脚本 MUST 在两个平台上正常工作
- **AND** 脚本 SHOULD 检测操作系统类型并输出相应提示
- **AND** 脚本 MUST 使用 POSIX 兼容的 shell 语法（避免 bashisms）

### Requirement: CLI 命令清理 (CLI Command Cleanup)
系统 MUST 移除 `glow init` 命令，以引导用户使用在线初始化脚本。

#### Scenario: 移除 init 命令
- **WHEN** 用户运行 `glow --help`
- **THEN** 帮助列表中不应出现 `init` 命令
- **AND** 用户运行 `glow init` MUST 返回错误并提示使用在线脚本
- **AND** 错误信息 MUST 包含在线脚本的 URL 或命令示例

#### Scenario: 迁移提示（可选过渡期）
- **WHEN** 用户运行旧的 `glow init` 命令（过渡期内）
- **THEN** 系统 MUST 输出迁移提示信息
- **AND** 提示信息 MUST 说明命令已废弃
- **AND** 提示信息 MUST 提供在线脚本的使用方法
- **AND** 系统 MUST 返回非零退出码

## MODIFIED Requirements

*No modified requirements in this change.*
