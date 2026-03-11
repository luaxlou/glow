# Glow

## 诞生理念

`glow` 的诞生，不是为了再造一个“大而全框架”，而是为了解决 Go 业务项目里一个反复出现的问题：

- 业务代码与基础接入代码（配置、HTTP、数据库、缓存）混在一起，项目越做越重。
- 团队在每个服务里重复造轮子，启动方式、配置约定、组件接入风格不一致。
- 业务开发被基础能力接入细节牵制，导致工程效率下降。

所以 `glow` 的定位非常克制：

- 只做应用侧框架能力（starter / sdk）
- 让业务代码保持聚焦和可读
- 通过稳定约定提升团队协作效率

## 职责边界

`glow` 的边界是应用侧框架能力本身：

- 提供统一 starter 接入范式
- 保持业务代码聚焦在领域逻辑
- 避免在框架层承载与应用接入无关的能力

## 设计原则

- 轻量优先：只提供高频、稳定、可复用的 starter 能力。
- 低心智负担：统一接入范式，降低新项目与新成员成本。
- 可组合：按需引入组件，不强绑定完整技术栈。
- 边界清晰：框架只承载应用接入相关职责。

## 这个仓库包含什么

- [`starter/glowconfig`](./starter/glowconfig)：配置读取
- [`starter/glowhttp`](./starter/glowhttp)：HTTP 服务启动适配（Gin）
- [`starter/glowmysql`](./starter/glowmysql)：MySQL 初始化
- [`starter/glowredis`](./starter/glowredis)：Redis 初始化
- [`starter/glowsqlite`](./starter/glowsqlite)：SQLite 初始化
- [`starter/glowwebsocket`](./starter/glowwebsocket)：WebSocket 适配
- [`examples/`](./examples)：使用示例
- [`docs/sdk_manual.md`](./docs/sdk_manual.md)：SDK 手册

## 这个仓库不包含什么

- 与应用 starter / sdk 无关的能力
- 与当前任务目标无关的大规模重构

## 快速开始

```bash
go get github.com/luaxlou/glow/starter
```

示例：

```go
package main

import (
	"fmt"

	"github.com/luaxlou/glow/starter/glowconfig"
)

func main() {
	fmt.Println(glowconfig.GetString("log_level"))
}
```

在应用目录准备 `config.json` 后直接运行即可。可参考 [`examples/simple-app`](./examples/simple-app)。

## 面向 AI Coding 的设计

`glow` 在设计上尽量让 AI Agent 能低成本理解并稳定改动：

- 约定优先：starter 按功能拆包，目录语义稳定，降低检索成本。
- 边界清晰：只围绕应用侧能力改动，避免任务外扩。
- 最小可运行路径：从 `examples/` 到 starter 包形成可追踪调用链。
- 验证标准统一：默认使用 `go test ./...` 与 `go vet ./...` 做交付前校验。

### AI 快速导航

- 应用接入入口：[`examples/simple-app`](./examples/simple-app)
- 配置能力：[`starter/glowconfig`](./starter/glowconfig)
- HTTP 能力：[`starter/glowhttp`](./starter/glowhttp)
- 数据库与缓存：[`starter/glowmysql`](./starter/glowmysql)、[`starter/glowredis`](./starter/glowredis)、[`starter/glowsqlite`](./starter/glowsqlite)
- 实时通信：[`starter/glowwebsocket`](./starter/glowwebsocket)
- 使用说明：[`docs/sdk_manual.md`](./docs/sdk_manual.md)

### AI 提示词模板

以下提示词可直接用于 AI Coding 场景。
建议在提示词中先要求 AI 阅读 [`docs/ai_coding_guide.md`](./docs/ai_coding_guide.md) 再执行。

#### 模板 1：新项目（基于 glow 初始化）

```text
你是 Go 架构师，请基于 glow 为我初始化一个新服务。

执行前请先阅读：docs/ai_coding_guide.md

目标：
1. 使用 glow 的 starter 完成配置读取、HTTP 启动与必要的数据组件接入。
2. 保持业务代码与基础设施代码解耦。
3. 提供最小可运行示例与目录说明。

约束：
1. 只在 glow 应用侧范畴内实现，不引入任务边界外能力。
2. 代码风格保持简洁，优先使用现有 starter，不重复造轮子。
3. 完成后必须执行并反馈 go test ./... 与 go vet ./... 的结果。

输出要求：
1. 先给出改动计划，再逐步实施。
2. 明确新增/修改的文件路径与原因。
3. 提供启动命令与验证步骤。
```

#### 模板 2：现有项目（接入/改造 glow）

```text
你是 Go 重构工程师，请在现有项目中引入 glow starter 并做最小侵入改造。

执行前请先阅读：docs/ai_coding_guide.md

当前目标：
1. 识别现有配置加载、HTTP 初始化、数据库连接代码。
2. 用 glow 对应 starter 替换重复样板代码。
3. 保持对外行为兼容，不做与目标无关的重构。

约束：
1. 先审阅项目结构并给出分步改造方案，不直接大面积改动。
2. 每次改动控制在可回滚粒度，说明风险点与回退方式。
3. 改造完成后必须执行 go test ./... 与 go vet ./...，并汇总结果。

输出要求：
1. 列出受影响模块、文件与接口。
2. 给出改造前后关键差异（初始化路径、配置入口、依赖注入方式）。
3. 说明后续可选优化项，但不要在本次直接实现。
```

## 开发与验证

```bash
go test ./...
go vet ./...
```
