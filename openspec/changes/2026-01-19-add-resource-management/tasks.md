## 1. Implementation
- [x] 1.1 在系统配置模块中增加删除 system_config 的接口
- [x] 1.2 新增 `glow-server remove` 命令组并实现交互确认与 `--yes`
- [x] 1.3 实现 `remove mysql|redis|nginx` 并删除对应存储键（幂等）

## 2. Validation
- [x] 2.1 增加单元测试覆盖：各资源的已配置/未配置/确认与 `--yes` 分支
- [x] 2.2 运行 `go test ./...`
