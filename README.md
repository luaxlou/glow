# Glow

## 诞生理念

`glow` 的目标不是增加一个大而全框架，而是把 Go 项目里的高频基础接入沉淀成稳定约定，让团队持续获得这些收益：

- 边界清晰：业务逻辑与基础接入职责分离
- 可组合：按需引入 starter，避免绑定整套技术栈
- 低心智负担：统一初始化范式，降低协作与接手成本
- 稳定约定：减少重复样板，提升长期可维护性

## AI Coding 快速引入（可直接复制）

AI 协作文档（完整链接）：
https://github.com/luaxlou/glow/blob/main/docs/ai_coding_guide.md

新项目流程卡（5 步）：
https://github.com/luaxlou/glow/blob/main/docs/quickstart_new_project.md

存量项目流程卡（5 步）：
https://github.com/luaxlou/glow/blob/main/docs/quickstart_existing_project.md

AI 输出契约：
https://github.com/luaxlou/glow/blob/main/docs/ai_output_contract.md

### 模板 1：新项目引入 glow

```text
请基于 glow 初始化一个新服务。
请先阅读：https://github.com/luaxlou/glow/blob/main/docs/ai_coding_guide.md
目标：快速接入 glowconfig + glowhttp，按需接入 mysql/redis/sqlite/websocket，并形成统一初始化范式。
重点收益：通过 glow 的稳定约定、可组合能力和低心智负担，提升项目可维护性与团队交付效率。
需解决问题：如何让项目在引入后持续享受 glow 设计哲学带来的收益，而不是一次性接入。
输出要求：严格遵循 https://github.com/luaxlou/glow/blob/main/docs/ai_output_contract.md
完成后执行并反馈：go test ./... && go vet ./...
```

### 模板 2：现有项目引入 glow

```text
请在现有项目中引入 glow starter。
请先阅读：https://github.com/luaxlou/glow/blob/main/docs/ai_coding_guide.md
目标：迁移到 glow 的统一接入范式，减少重复样板并提升长期可维护性。
重点收益：让项目获得边界清晰、可组合复用、低心智负担和稳定约定。
需解决问题：如何把接入过程沉淀为可持续演进的工程基线。
输出要求：严格遵循 https://github.com/luaxlou/glow/blob/main/docs/ai_output_contract.md
完成后执行并反馈：go test ./... && go vet ./...
```

## 仓库优化方向（已落地）

- P0 接入体验产品化：
  - [新项目流程卡](./docs/quickstart_new_project.md)
  - [存量项目流程卡](./docs/quickstart_existing_project.md)
  - [AI 输出契约](./docs/ai_output_contract.md)
- P1 Starter 一致性治理：
  - [Starter 约定规范](./docs/starter_conventions.md)
- P2 可组合能力矩阵：
  - [Starter 组合矩阵](./docs/starter_composition_matrix.md)
- P3 AI 协作契约化：
  - [PR 模板](./.github/pull_request_template.md)
- P4 示例与验证强化：
  - [`examples/minimal-api`](./examples/minimal-api)
  - [`examples/api-db-cache`](./examples/api-db-cache)
  - [`examples/api-websocket`](./examples/api-websocket)

## 这个仓库包含什么

- [`starter/glowconfig`](./starter/glowconfig)：配置读取
- [`starter/glowhttp`](./starter/glowhttp)：HTTP 服务启动适配（Gin）
- [`starter/glowmysql`](./starter/glowmysql)：MySQL 初始化
- [`starter/glowredis`](./starter/glowredis)：Redis 初始化
- [`starter/glowsqlite`](./starter/glowsqlite)：SQLite 初始化
- [`starter/glowwebsocket`](./starter/glowwebsocket)：WebSocket 适配
- [`examples/`](./examples)：示例集合
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

## 开发与验证

```bash
go test ./...
go vet ./...
```

## 相关项目（简述）

`glow-ops` 是配套部署基础设施项目，聚焦部署与运维生命周期，不属于本仓库实现范围。
