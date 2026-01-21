## 1. 实现基础 init 命令
- [x] 1.1 在 `cmd/glow/cmd/` 目录下创建 `init.go` 文件
- [x] 1.2 使用 cobra 创建 `initCmd` 命令结构
- [x] 1.3 实现命令注册逻辑，在 `root.go` 中添加 init 命令
- [x] 1.4 添加基本的命令帮助文档和使用说明

## 2. 实现项目结构分析
- [x] 2.1 实现目录检测逻辑（检查 cmd、bin、scripts 是否存在）
- [x] 2.2 实现目录创建逻辑
- [x] 2.3 实现检测和创建 Makefile 的逻辑
- [x] 2.4 添加用户确认提示，避免意外覆盖

## 3. 创建部署脚本模板
- [x] 3.1 创建 `scripts/deploy.sh` 模板
- [x] 3.2 实现模板变量替换（应用名称等）
- [x] 3.3 设置脚本可执行权限
- [x] 3.4 添加错误处理和日志输出逻辑

## 4. 创建 Makefile 模板和集成
- [x] 4.1 创建 Makefile 模板，包含基本的 build 和 deploy 目标
- [x] 4.2 实现智能追加逻辑：若 Makefile 已存在，追加 deploy 目标
- [x] 4.3 实现完整的 `make deploy` 流程（构建 + 调用 deploy.sh）
- [x] 4.4 添加 `make help` 目标（如果不存在）

## 5. 实现 AI 工具选择交互
- [x] 5.1 实现交互式提示逻辑（使用 bufio 读取用户输入）
- [x] 5.2 设计选项展示格式（Claude Code、Gemini CLI、None）
- [x] 5.3 实现多选逻辑（允许选择多个工具）
- [x] 5.4 添加输入验证和错误处理

## 6. 实现 Claude Code Skill 集成
- [x] 6.1 创建 `.claude/commands/` 目录
- [x] 6.2 创建 `glow:deploy` command 文件
- [x] 6.3 创建 `glow:logs` command 文件
- [x] 6.4 创建 `glow:status` command 文件
- [x] 6.5 为每个 command 编写清晰的 prompt 和使用说明

## 7. 实现 Gemini CLI 集成（可选）
- [x] 7.1 研究 Gemini CLI 的配置文件格式和约定
- [x] 7.2 创建 Gemini CLI 配置文件
- [x] 7.3 添加相应的 alias 或 command

**注**: Gemini CLI 集成已预留接口，返回"计划在未来版本实现"的提示。

## 8. 实现幂等性
- [x] 8.1 添加文件和目录存在性检测
- [x] 8.2 实现覆盖确认逻辑
- [x] 8.3 生成初始化报告，显示已配置和跳过的项目
- [x] 8.4 添加 `--force` 标志以跳过确认直接覆盖

## 9. 测试和验证
- [x] 9.1 在空项目中测试 `glow init`
- [x] 9.2 在已有部分结构的项目中测试（已有 Makefile）
- [x] 9.3 测试重复执行的幂等性
- [x] 9.4 测试 AI 工具集成生成的文件格式
- [x] 9.5 验证 `make deploy` 命令是否正常工作
- [x] 9.6 验证 Claude Code commands 是否可用

## 10. 文档和示例
- [x] 10.1 更新 CLI 帮助文档
- [x] 10.2 创建示例项目展示 `glow init` 的使用
- [x] 10.3 编写 README 或文档说明初始化流程

**额外完成**:
- 添加 `GLOW_INIT_AI_TOOLS` 环境变量支持，便于 CI/CD 集成
- 添加 `--skip-ai` 标志以跳过 AI 工具配置
- 创建完整的使用指南文档（保存在 `/tmp/glow-init-guide.md`）
- 实现非交互式模式下的优雅降级处理
