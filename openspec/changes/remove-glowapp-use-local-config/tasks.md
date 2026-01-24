# 任务: 移除 GlowApp Starter 并使用本地配置

## 实现 (Implementation)
- [x] 创建 `pkg/glowconfig` 包，包含简单的 JSON 配置加载器。 <!-- id: create-glowconfig -->
- [x] 更新 `internal/apiserver` 以写入 `config.json` 而不是 `<app>_local_config.json`。 <!-- id: update-server-render -->
- [x] 重构内部代码库（如有）或示例中的 `starter/glowapp` 依赖。 <!-- id: refactor-deps -->
- [x] 清理 `starter/glowapp` 目录及相关依赖。 <!-- id: delete-glowapp -->

## 验证 (Validation)
- [x] 验证 `glow apply` 生成 `config.json`。 <!-- id: verify-apply -->
- [x] 验证示例应用可以使用 `pkg/glowconfig` 读取 `config.json`。 <!-- id: verify-app-read -->
