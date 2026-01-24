# 实现任务清单

## 1. App.yaml 配置支持
- [x] 1.1 在 `pkg/api/types.go` 的 `AppSpec` 中确认 `Config map[string]any` 字段存在
- [x] 1.2 更新 `cmd/glow/apply.go`，解析 `spec.config` 字段
- [x] 1.3 实现配置文件生成逻辑：将 `spec.config` 写入 `<data-dir>/apps/<appName>/<appName>_local_config.json`
- [x] 1.4 更新 `examples/simple-app/app.yaml` 示例，展示 config 字段用法
- [x] 1.5 更新 `examples/app-yaml-format.md` 文档，说明 config 字段

## 2. 配置管理方式调整（完全声明式）
- [x] 2.1 **移除** `cmd/glow/config.go` 及所有 `glow config` 命令
- [x] 2.2 配置完全通过 app.yaml 的 `spec.config` 声明式管理
- [x] 2.3 配置变更通过修改 app.yaml 并执行 `glow apply` 实现

## 3. 服务器信息查询功能
- [x] 3.1 在 `pkg/api/types.go` 新增 `ServerInfo` 结构体
- [x] 3.2 在 `internal/apiserver/server.go` 新增端点 `GET /server/info`
- [x] 3.3 实现服务器信息收集逻辑（版本、路径、PID、运行时长等）
- [x] 3.4 在 `cmd/glow/server.go` 实现 `glow server info` 子命令
- [x] 3.5 支持 `--json` 参数输出 JSON 格式

## 4. 测试与文档
- [x] 4.1 手动测试 app.yaml config 字段声明与 apply 流程
- [x] 4.2 ~~手动测试 `glow config set/get/list/export` 命令~~（已移除，无需测试）
- [x] 4.3 手动测试 `glow server info` 命令（包括 --json 参数）
- [x] 4.4 更新 `docs/cli_manual.md`，说明 config 字段和 server info 用法
- [x] 4.5 更新 `docs/sdk_manual.md`，说明应用如何读取配置
- [x] 4.6 更新 `README.md`，说明配置管理方式和运维命令
