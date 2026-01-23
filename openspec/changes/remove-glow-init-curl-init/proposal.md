# Change: Remove glow init Command, Use Online Shell Script for Project Initialization

## Why
当前的 `glow init` 命令需要在安装 glow CLI 后才能使用，这增加了新项目接入手续：用户必须先安装 CLI，然后才能初始化项目。此外，该命令内嵌于二进制中，更新逻辑需要发布新版本。通过将项目初始化逻辑迁移到在线 shell 脚本（类似 glow-server 的安装脚本），可以实现：

1. **零依赖初始化**：用户无需预先安装 glow CLI，直接通过 `curl | sh` 即可初始化项目
2. **独立更新**：初始化逻辑可以随时更新，无需发布新版本
3. **轻量级 CLI**：移除 init 命令后，glow CLI 更加精简，专注于应用治理核心功能
4. **统一体验**：与 glow-server 的 `curl -fsSL https://... | bash` 安装方式保持一致

## What Changes
- **移除** `glow init` 命令（包括 `cmd/glow/cmd/init.go`）
- **新增** 在线 shell 脚本，使用 GitHub raw URL（`https://raw.githubusercontent.com/{owner}/{repo}/main/scripts/init-project.sh`）
- **迁移** 项目结构初始化逻辑到在线脚本（目录创建、deploy.sh 生成、AI 工具集成）
- **更新** 项目初始化文档，推荐使用 `curl -fsSL <url> | bash` 方式
- **保留** AI 工具集成能力（Claude Code skills 复制等），但通过脚本实现

## Impact
- Affected specs:
  - **project-initialization**: 移除 `glow init` 命令相关需求，改为在线脚本需求
- Affected code:
  - `cmd/glow/cmd/init.go` - 完全删除
  - `cmd/glow/cmd/root.go` - 移除 initCmd 注册
  - 新增在线脚本 `scripts/init-project.sh`（通过 GitHub raw URL 访问）
- User-facing changes:
  - 用户使用 `curl -fsSL https://raw.githubusercontent.com/{owner}/{repo}/main/scripts/init-project.sh | bash` 初始化项目
  - 不再需要先安装 glow CLI
  - 支持通过参数自定义初始化行为（如 `curl ... | bash -s -- --skip-ai`）

## Migration Path
- 已初始化的项目：不受影响，继续使用现有的 deploy.sh 和配置
- 文档更新：将 `glow init` 示例更新为在线脚本调用方式
- 立即移除：直接移除 `glow init` 命令，不保留过渡期，通过文档和 CHANGELOG 引导用户迁移
