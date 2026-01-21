## ADDED Requirements

### Requirement: 项目结构初始化 (Project Structure Initialization)
系统 MUST 提供 `glow init` 命令，用于分析并初始化项目结构以支持 Glow 治理。

#### Scenario: 标准目录结构创建
- **WHEN** 用户在项目根目录运行 `glow init` 且缺少标准目录
- **THEN** 系统应检测是否存在 `cmd`、`bin`、`scripts` 目录
- **AND** 若目录不存在，系统应创建这些目录
- **AND** 系统应在 `scripts` 目录下创建 `deploy.sh` 脚本模板
- **AND** 系统应检查项目根目录是否存在 `Makefile`
- **AND** 若 `Makefile` 不存在，系统应创建包含 `make deploy` 目标的 Makefile
- **AND** 若 `Makefile` 已存在，系统应添加 `make deploy` 目标（不覆盖现有内容）

#### Scenario: 交互式应用名称配置
- **WHEN** 用户执行 `make deploy`
- **THEN** 系统应提示用户输入应用名称
- **AND** 系统应验证应用名称格式（字母、数字、连字符，不以连字符开头/结尾）
- **AND** 系统应将应用名称传递给 `scripts/deploy.sh`

### Requirement: AI 工具集成 (AI Tool Integration)
系统 MUST 支持在初始化时选择并配置 AI 工具集成。

#### Scenario: AI 工具选择
- **WHEN** 用户运行 `glow init`
- **THEN** 系统应提示 "是否要配置 AI 工具集成？"
- **AND** 系统应提供选项列表：
  - Claude Code
  - Gemini CLI
  - None (跳过)
- **AND** 用户可以选择多个工具或选择 None

#### Scenario: Claude Code Skill 配置
- **WHEN** 用户选择配置 Claude Code
- **THEN** 系统应在项目根目录创建 `.claude/commands/` 目录
- **AND** 系统应创建 `glow:deploy` command 文件
- **AND** command 文件应包含部署应用的 prompt
- **AND** 系统应创建 `glow:logs` command 文件用于查看日志
- **AND** 系统应创建 `glow:status` command 文件用于查看应用状态
- **AND** 每个 command 应提供清晰的说明和用法示例

#### Scenario: Gemini CLI 配置
- **WHEN** 用户选择配置 Gemini CLI
- **THEN** 系统应根据 Gemini CLI 的约定创建相应的配置文件
- **AND** 系统应添加 `glow:deploy` 等 alias 或 command

#### Scenario: 跳过 AI 工具配置
- **WHEN** 用户选择 None
- **THEN** 系统应跳过 AI 工具配置步骤
- **AND** 系统应完成项目结构初始化

### Requirement: 部署脚本模板 (Deploy Script Template)
系统 MUST 提供可执行的部署脚本模板。

#### Scenario: deploy.sh 模板生成
- **WHEN** 系统创建 `scripts/deploy.sh`
- **THEN** 脚本应包含 glow CLI 的基本部署命令
- **AND** 脚本应接收应用名称作为参数
- **AND** 脚本应包含错误处理和日志输出
- **AND** 脚本应具有可执行权限（chmod +x）

#### Scenario: Makefile 集成
- **WHEN** 系统创建或更新 Makefile
- **THEN** `make deploy` 目标应调用 `scripts/deploy.sh`
- **AND** 应包含构建和部署的完整流程
- **AND** 应提供 `help` 目标列出所有可用的 make 命令

### Requirement: 幂等性 (Idempotency)
`glow init` 命令 MUST 具备幂等性，可安全地重复执行。

#### Scenario: 重复执行 init
- **WHEN** 用户在已初始化的项目中再次运行 `glow init`
- **THEN** 系统应检测已存在的文件和目录
- **AND** 系统应提示用户哪些项目已配置
- **AND** 系统应询问是否要覆盖现有配置
- **AND** 系统应默认保留现有配置，仅在用户确认时覆盖

#### Scenario: 部分配置更新
- **WHEN** 用户选择重新配置 AI 工具
- **THEN** 系统应仅更新 AI 工具相关配置
- **AND** 系统应保留项目结构配置（目录、Makefile 等）
