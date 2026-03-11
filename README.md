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

## AI Coding 快速引入（可直接复制）

### AI 提示词模板

以下模板面向“快速引入”，可直接复制到 AI IDE 使用。
AI 协作文档（完整链接）：https://github.com/luaxlou/glow/blob/main/docs/ai_coding_guide.md

#### 模板 1：新项目（基于 glow 初始化）

```text
请先阅读：https://github.com/luaxlou/glow/blob/main/docs/ai_coding_guide.md
你是 Go 工程师，请基于 glow 初始化一个新服务。
要求：使用 glowconfig + glowhttp，按需接入 mysql/redis/sqlite/websocket。
约束：仅改动本仓库职责范围；优先复用 starter；避免无关重构。
输出：先给实施计划，再给文件改动清单、启动命令、验证结果。
完成后执行并反馈：go test ./... && go vet ./...
```

#### 模板 2：现有项目（接入/改造 glow）

```text
请先阅读：https://github.com/luaxlou/glow/blob/main/docs/ai_coding_guide.md
你是 Go 重构工程师，请在现有项目中最小侵入接入 glow starter。
目标：识别并替换配置加载、HTTP 初始化、数据库连接等重复样板代码。
约束：保持对外行为兼容；分步改造、每步可回滚；不做任务外重构。
输出：受影响文件清单、关键差异说明、验证结果与回退方式。
完成后执行并反馈：go test ./... && go vet ./...
```

## 这个仓库包含什么

- [`starter/glowconfig`](./starter/glowconfig)：配置读取
- [`starter/glowhttp`](./starter/glowhttp)：HTTP 服务启动适配（Gin）
- [`starter/glowmysql`](./starter/glowmysql)：MySQL 初始化
- [`starter/glowredis`](./starter/glowredis)：Redis 初始化
- [`starter/glowsqlite`](./starter/glowsqlite)：SQLite 初始化
- [`starter/glowwebsocket`](./starter/glowwebsocket)：WebSocket 适配
- [`examples/`](./examples)：使用示例
- [`docs/sdk_manual.md`](./docs/sdk_manual.md)：SDK 手册

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

## 开发与验证

```bash
go test ./...
go vet ./...
```

## 相关项目（简述）

`glow-ops` 是配套的部署基础设施项目，聚焦部署与运维生命周期，不属于本仓库实现范围。
