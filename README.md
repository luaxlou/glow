# Glow（Framework）

`glow` 是应用侧框架仓库，提供 Go Starter 与 SDK 使用范式。

## 这个仓库解决什么问题

如果你是业务开发者，想在 Go 应用里快速接入：

- 配置读取
- HTTP 服务启动（Gin）
- MySQL / Redis / SQLite 客户端初始化
- WebSocket 能力

那你应该使用这个仓库。

## 你不该在这里找什么

这个仓库**不包含**运维控制面能力：

- 不包含 `glow-server`
- 不包含 `glow-cli`
- 不包含部署/回滚/进程托管/节点与 ingress 编排

这些能力已经拆分到 [`glow-ops`](https://github.com/luaxlou/glow-ops) 仓库。

## 仓库结构

- [`starter/glowconfig`](./starter/glowconfig)：本地配置读取
- [`starter/glowhttp`](./starter/glowhttp)：HTTP 启动适配
- [`starter/glowmysql`](./starter/glowmysql)：MySQL 初始化
- [`starter/glowredis`](./starter/glowredis)：Redis 初始化
- [`starter/glowsqlite`](./starter/glowsqlite)：SQLite 初始化
- [`starter/glowwebsocket`](./starter/glowwebsocket)：WebSocket 适配
- [`examples/`](./examples)：纯 SDK 示例
- [`docs/sdk_manual.md`](./docs/sdk_manual.md)：SDK 手册

## 3 分钟上手

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

## 相关仓库

- 运维控制面：[`glow-ops`](https://github.com/luaxlou/glow-ops)
