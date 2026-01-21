# Change: Add Project Initialization Command

## Why
当前开发者想要将项目接入 Glow 治理体系时，需要手动创建多个目录和配置文件（cmd、bin、scripts、deploy.sh、Makefile 等），且缺少对 AI 工具（Claude Code、Gemini CLI）的标准化集成方式。这导致新项目接入成本较高，且不同项目的配置可能不一致。通过提供 `glow init` 命令，可以自动化项目初始化流程，确保项目结构标准化，并简化 AI 工具集成。

## What Changes
- 添加 `glow init` 命令，用于初始化项目结构
- 自动创建标准目录结构（cmd、bin、scripts）和部署脚本
- 自动生成或更新 Makefile，添加 `make deploy` 目标
- 支持交互式选择 AI 工具集成（Claude Code、Gemini CLI）
- 为 Claude Code 创建 glow:deploy、glow:logs、glow:status 等 commands
- 确保命令具备幂等性，可安全重复执行

## Impact
- Affected specs:
  - **新增**: `project-initialization`（项目初始化能力）
- Affected code:
  - `cmd/glow/cmd/` - 添加新的 `init.go` 文件实现 init 命令
  - 可能需要添加模板文件用于生成 deploy.sh 和 Makefile
- User-facing changes:
  - 用户可以通过 `glow init` 快速初始化项目
  - 用户可以选择性地集成 AI 工具，提升开发体验
  - 项目结构更加标准化，降低学习成本
