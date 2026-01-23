# Glow Skills for Claude Code

This directory contains three specialized skills that provide Claude Code with comprehensive knowledge about the Glow framework.

## Skills Overview

### 1. glow-sdk
**Focus**: SDK 开发指南

Contains:
- GlowApp, GlowHTTP, GlowMySQL, GlowRedis, GlowConfig, GlowWebSocket 组件使用
- API 参考和代码示例
- 最佳实践和模式
- Advanced topics (custom starters, middleware, metrics)

**When triggered**: 用户需要开发 Go 应用、集成 SDK 组件、处理配置和数据库时

### 2. glow-deploy
**Focus**: 部署与运维

Contains:
- Glow CLI 命令参考
- 部署工作流
- 多环境管理
- CI/CD 集成
- Architecture overview
- 配置管理和日志

**When triggered**: 用户需要部署应用、管理生命周期、查看日志、配置环境时

### 3. glow-debug
**Focus**: 本地开发与调试

Contains:
- glow-server 设置
- 本地开发环境配置
- 热重载和 IDE 集成
- 调试技巧和性能分析
- 测试方法
- Troubleshooting guide

**When triggered**: 用户需要本地开发、调试代码、解决技术问题时

## How It Works

When users run `glow init` and choose to configure Claude Code integration:

1. All three skills are automatically copied from `.claude/skills/` to the project's `.claude/skills/` directory
2. Claude Code will automatically load the appropriate skill based on the task
3. Each skill provides focused, domain-specific knowledge

## Skill Structure

```
.claude/skills/
├── README.md                     # This file
├── glow-sdk/
│   ├── SKILL.md                  # Main SDK guide
│   └── references/
│       └── advanced-topics.md    # Advanced SDK patterns
├── glow-deploy/
│   ├── SKILL.md                  # Main deployment guide
│   └── references/
│       └── architecture.md       # Architecture and design
└── glow-debug/
    ├── SKILL.md                  # Main debugging guide
    └── references/
        └── troubleshooting.md    # Common issues and solutions
```

## Development

If you update any of the skill files in this directory, users will get the updated version the next time they run `glow init`.

## Testing

To test the skill integration:

```bash
# Create a test project
mkdir test-project && cd test-project

# Initialize with Claude Code support
glow init
# Choose "y" when asked about Claude Code integration

# Verify the skills were copied
ls -la .claude/skills/
# Should show: glow-sdk/, glow-deploy/, glow-debug/
```

## Design Principles

### Why 3 Separate Skills?

1. **Focused expertise**: Each skill targets a specific domain (SDK/Deploy/Debug)
2. **Context efficiency**: Claude only loads relevant knowledge for the current task
3. **Easier maintenance**: Updates can be made to specific domains
4. **Clearer triggers**: Each skill has specific triggering conditions

### Skill Triggering Examples

- **glow-sdk**: "How do I use GlowMySQL?" / "Create a Glow app with Redis"
- **glow-deploy**: "Deploy my app" / "Check application logs" / "Setup CI/CD"
- **glow-debug**: "App won't start" / "How to debug locally?" / "Port conflict"

## Related Files

- `cmd/glow/cmd/init.go`: Contains the logic to copy skills during `glow init`
- `.gitignore`: Updated to include `.claude/` (skills are part of the repository)
- `docs/`: Official documentation (more detailed, user-facing)

## Future Enhancements

Possible future skills:
- **glow-monitoring**: Metrics, alerting, observability
- **glow-security**: Security best practices, authentication
- **glow-scaling**: Horizontal/vertical scaling strategies
