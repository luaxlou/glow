# Glow AI Coding 协作指南

本文档用于让 AI Agent 在 `glow` 仓库内快速建立上下文、减少误改并提升交付稳定性。

## 仓库定位

`glow` 是应用侧框架仓库，只提供 starter / sdk 能力。

## 修改边界

允许改动：

- `starter/*` 下的框架能力代码
- `examples/*` 下的使用示例
- `docs/*` 下的说明文档

禁止越界：

- 不实现与 starter / sdk 无关的系统能力
- 不引入与当前任务目标无关的额外职责
- 不做与当前任务无关的大规模重构

## 目录快速索引

- `starter/glowconfig`：配置读取
- `starter/glowhttp`：HTTP 启动适配（Gin）
- `starter/glowmysql`：MySQL 初始化
- `starter/glowredis`：Redis 初始化
- `starter/glowsqlite`：SQLite 初始化
- `starter/glowwebsocket`：WebSocket 适配
- `examples/simple-app`：最小示例入口
- `docs/sdk_manual.md`：SDK 手册

## 推荐工作流程（AI Agent）

1. 先读任务目标与约束，确认是否属于 `glow` 仓库职责。
2. 在改动前列出受影响文件与最小实施路径。
3. 优先复用现有 starter 能力，避免重复造轮子。
4. 控制提交粒度，保证每次改动都可回滚。
5. 交付前执行统一验证命令并回报结果。

## 统一验证命令

```bash
go test ./...
go vet ./...
```

## 常用任务提示

新项目初始化类任务：

- 以 `examples/simple-app` 为最小可运行参考
- 先接入 `glowconfig` 与 `glowhttp`，再按需增加数据库/缓存组件

存量项目改造类任务：

- 先识别现有配置与启动路径，再逐步替换为对应 starter
- 保持对外行为兼容，避免一次性重写

## 交付输出建议

- 明确列出新增/修改文件路径
- 说明关键实现决策与边界控制
- 给出运行与验证步骤
