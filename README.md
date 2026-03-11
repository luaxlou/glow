# Glow

## 诞生理念

`glow` 的诞生，不是为了再造一个“大而全框架”，而是为了解决 Go 业务项目里一个反复出现的问题：

- 业务代码与基础接入代码（配置、HTTP、数据库、缓存）混在一起，项目越做越重。
- 团队在每个服务里重复造轮子，启动方式、配置约定、组件接入风格不一致。
- 业务开发被运维流程细节牵制，导致工程效率下降。

所以 `glow` 的定位非常克制：

- 只做应用侧框架能力（starter / sdk）
- 让业务代码保持聚焦和可读
- 通过稳定约定提升团队协作效率

## 为什么拆分为两个仓库

为彻底解耦职责，`glow` 与 `glow-ops` 已拆分为两个独立仓库：

- [`glow`](https://github.com/luaxlou/glow)：应用侧框架（starter / sdk）
- [`glow-ops`](https://github.com/luaxlou/glow-ops)：部署与运维基础设施（`glow-server` / `glow-cli`）

这不是组织形式调整，而是工程边界设计：

- 业务研发只关注应用构建与能力接入
- 运维平台只关注部署生命周期与治理编排
- 双方独立演进，避免相互牵制

## 设计原则

- 轻量优先：只提供高频、稳定、可复用的 starter 能力。
- 低心智负担：统一接入范式，降低新项目与新成员成本。
- 可组合：按需引入组件，不强绑定完整技术栈。
- 边界清晰：框架不承载运维控制面职责。

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

以下能力不在 `glow`，请到 [`glow-ops`](https://github.com/luaxlou/glow-ops)：

- `glow-server`
- `glow-cli`
- 部署编排、回滚、节点与入口治理

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
