# Design: Remove glow init Command and Migrate to Online Script

## Context

当前 `glow init` 命令作为 CLI 的一部分存在，要求用户先安装 glow CLI 才能初始化项目。这增加了新用户的接入成本。Glow 项目已有使用在线 shell 脚本的成功经验（`scripts/install.sh` 用于安装 glow-server），因此可以采用相同模式来实现项目初始化。

### Stakeholders
- **新用户**：希望快速接入 Glow，无需预先安装工具
- **现有用户**：已使用 `glow init`，需要迁移到新方式
- **维护者**：希望简化 CLI，独立更新初始化逻辑

## Goals / Non-Goals

### Goals
1. **零依赖初始化**：用户无需安装 glow CLI 即可初始化项目
2. **独立更新**：初始化逻辑可以随时更新，无需发布新版本
3. **统一体验**：与 glow-server 的安装脚本保持一致的使用模式
4. **功能完整**：保留所有现有功能（目录创建、deploy.sh 生成、AI 工具集成）

### Non-Goals
1. 不改变项目结构标准（cmd/, bin/, scripts/）
2. 不改变 AI 工具集成逻辑（Claude Code skills 等）
3. 不在脚本中添加 Go 编译或二进制下载功能

## Decisions

### Decision 1: 使用在线 Shell 脚本替代 CLI 命令

**Rationale**:
- 减少新用户的接入步骤（无需先安装 CLI）
- 初始化逻辑可以独立更新，无需发布新版本
- 与 glow-server 的安装模式保持一致

**Alternatives Considered**:
1. **保留 `glow init` 命令**：
   - 优点：现有用户熟悉，无需改变文档
   - 缺点：仍需先安装 CLI，更新需要发版
   - **决策**：不采用

2. **使用 Go 模板或脚手架工具**（如 cookiecutter）：
   - 优点：功能更强大，支持复杂模板
   - 缺点：引入额外依赖，不符合 Glow 的轻量级理念
   - **决策**：不采用

3. **使用在线 shell 脚本**：
   - 优点：零依赖、易更新、与现有模式一致
   - 缺点：需要学习 shell 脚本维护
   - **决策**：采用此方案

### Decision 2: 脚本 URL 命名

**Choice**: 使用 GitHub raw URL 作为唯一源，格式为 `https://raw.githubusercontent.com/{owner}/{repo}/main/scripts/init-project.sh`

**Rationale**:
- 直接使用源代码仓库，无需额外维护 CDN 或域名
- 版本控制清晰，可通过分支或 tag 管理版本
- 减少依赖，降低维护成本
- 与项目代码同步，确保一致性

**Alternatives**:
1. 使用自定义域名（如 `get.glow.dev`）：需要额外维护，增加复杂度
2. 使用 CDN：依赖第三方服务，增加故障点

**决策**：仅使用 GitHub raw URL 作为唯一源

### Decision 3: 迁移策略

**Choice**: 立即移除 `glow init` 命令，不保留过渡期

**Rationale**:
- 避免维护两套逻辑（CLI + 脚本）
- 简化代码，减少维护成本
- 通过文档和 CHANGELOG 引导用户迁移

**Implementation**:
- 直接从 CLI 代码中移除 `init.go` 和命令注册
- 更新所有文档，使用在线脚本方式
- 在 CHANGELOG 中明确说明迁移方式

### Decision 4: AI Skills 复制逻辑

**Challenge**: 在线脚本需要找到 glow 安装目录以复制 skills

**Solution**:
1. 脚本首先尝试通过 `which glow` 找到 CLI 安装路径
2. 如果找不到，提示用户先安装 glow CLI 或手动下载 skills
3. 从 glow 安装目录的 `.claude/skills/` 复制到项目目录

**Error Handling**:
- 如果 glow CLI 未安装，跳过 skills 复制步骤，输出警告
- 提供手动下载 skills 的 URL 或说明

## Risks / Trade-offs

### Risk 1: 用户依赖 `glow init` 命令
- **影响**：现有用户的自动化脚本或文档可能失效
- **缓解**：
  - 更新所有官方文档，使用新的在线脚本方式
  - 在 CHANGELOG 中突出说明迁移方式
  - 提供清晰的迁移指南

### Risk 2: 在线脚本可用性
- **影响**：如果 GitHub raw URL 不可用，新用户无法初始化项目
- **缓解**：
  - GitHub 服务稳定性高，raw URL 访问可靠
  - 在 README 中提供脚本的完整路径和替代访问方式
  - 考虑将脚本打包到 release 中作为备用方案

### Risk 3: Shell 脚本维护复杂度
- **影响**：Shell 脚本比 Go 代码更难测试和维护
- **缓解**：
  - 保持脚本逻辑简单（仅文件操作）
  - 添加充分的注释
  - 在多个平台（Linux/macOS）上测试
  - 使用 shellcheck 验证脚本语法

### Trade-off: 功能复杂度 vs 易用性
- **选择**：保持脚本功能简单，不添加过多特性
- **理由**：脚本越复杂，维护成本越高，出错概率越大
- **边界**：仅实现核心初始化功能，高级功能通过配置文件或额外脚本实现

## Migration Plan

### Phase 1: Preparation (实现阶段)
1. 创建 `scripts/init-project.sh` 脚本
2. 实现所有核心功能（目录创建、deploy.sh 生成、AI 工具集成）
3. 本地测试脚本功能

### Phase 2: Deployment (部署阶段)
1. 将脚本提交到 GitHub 仓库
2. 测试在线 URL 可访问性
3. 更新 README 和文档

### Phase 3: CLI Cleanup (清理阶段)
1. 从 `cmd/glow/cmd/` 删除 `init.go`
2. 更新 `root.go` 移除 initCmd 注册

### Phase 4: Release (发布阶段)
1. 发布新版本 glow CLI
2. 更新 CHANGELOG
3. 发布公告说明迁移方式

### Rollback Plan
如果新方式出现问题，可以：
1. 恢复 `cmd/glow/cmd/init.go` 文件
2. 在 initCmd 中添加弃用警告，引导用户使用在线脚本
3. 逐步引导用户迁移

## Open Questions

1. **脚本版本管理**：
   - Q: 是否需要版本化的脚本 URL（如 `init-v1.sh`）？
   - A: 暂不需要，使用 main 分支的脚本即可。如果需要重大变更，再考虑版本化。

2. **CDN 选择**：
   - Q: 是否使用 CDN（如 Cloudflare）托管脚本？
   - A: 不使用 CDN，直接使用 GitHub raw URL 作为唯一源。如果未来有访问速度问题，再考虑优化方案。

3. **Skills 下载**：
   - Q: 如果用户未安装 glow CLI，如何提供 AI skills？
   - A: 提供 skills 仓库的 URL 或单独的下载脚本，用户可手动下载。

4. **非交互模式**：
   - Q: 脚本是否需要完全支持 CI/CD 场景？
   - A: 是的，通过 `--skip-ai --force` 参数支持非交互模式。

## Implementation Notes

### Script Structure
```bash
#!/bin/bash
set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Functions
detect_platform() { ... }
analyze_project() { ... }
create_directories() { ... }
create_deploy_script() { ... }
setup_ai_tools() { ... }

# Main flow
main() {
  parse_arguments "$@"
  detect_platform
  analyze_project
  create_directories
  create_deploy_script
  setup_ai_tools
}

main "$@"
```

### Testing Strategy
1. 本地测试：在空目录、已初始化目录、部分配置目录中测试
2. 平台测试：在 Linux（Ubuntu、CentOS）和 macOS 上测试
3. 参数测试：测试各种参数组合（--skip-ai, --force）
4. 集成测试：测试生成的 deploy.sh 是否正常工作

### Documentation Updates
1. README.md：将 `glow init` 替换为 `curl -fsSL https://raw.githubusercontent.com/{owner}/{repo}/main/scripts/init-project.sh | bash`
2. docs/getting-started.md：更新快速开始指南
3. docs/cli_manual.md：移除 init 命令说明
4. CHANGELOG.md：添加迁移说明
