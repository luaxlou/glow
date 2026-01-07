# Glow SDK 用户手册

`glow/starter` 是 Glow 框架的 Go 语言 SDK，提供了开箱即用的应用骨架和组件。通过使用 SDK，开发者可以轻松地将应用接入 Glow 的全生命周期管理，享受自动化资源拨备、配置热更新和进程托管等特性。

## 1. 安装

在你的 Go 项目中引入 Glow SDK：

```bash
go get github.com/luaxlou/glow/starter
```

## 2. 快速开始 (Quick Start)

以下是一个最简单的 Web 应用示例，展示了 Glow 的核心用法：

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/luaxlou/glow/starter/glowapp"
    "github.com/luaxlou/glow/starter/glowhttp"
)

func main() {
    // 1. 初始化应用身份
    glowapp.Init("my-first-app")

    // 2. 初始化 HTTP 组件 (指定默认端口，实际运行端口由 Server 分配)
    glowhttp.Init(8080)

    // 3. 获取 Router 并注册路由
    r := glowhttp.Router()
    r.GET("/", func(c *gin.Context) {
        c.String(200, "Hello Glow!")
    })

    // 4. 启动服务
    // SDK 会自动处理端口监听、信号捕捉和优雅停机
    glowhttp.Run()
    glowapp.WaitForShutdown()
}
```

## 3. 核心组件

Glow SDK 提供了一系列标准化的 Starter 组件。

### 3.1 GlowApp (应用核心)

`glowapp` 是应用的入口，负责身份注册和生命周期管理。

*   **Init(appName string, opts ...Option)**: 初始化应用身份，并向本地 Glow Server 注册。
*   **WaitForShutdown()**: 阻塞主 goroutine，直到接收到系统信号 (SIGINT/SIGTERM)。收到信号后，它会通知所有已注册的组件执行清理操作 (Graceful Shutdown)。
*   **RegisterCleanup(name string, fn func())**: 注册自定义的清理函数，在停机时执行。

### 3.2 GlowHTTP (Web 服务)

`glowhttp` 基于 [Gin](https://github.com/gin-gonic/gin) 框架封装，集成了 Glow 的端口管理和生命周期。

*   **Init(defaultPort int)**: 初始化 HTTP 组件。如果 Glow Server 指定了 `OP_APP_PORT` 环境变量，则忽略 `defaultPort`。
*   **Router() *gin.Engine**: 获取 Gin 实例，用于注册路由和中间件。
*   **Run()**: 异步启动 HTTP Server。

### 3.3 GlowMySQL (数据库)

`glowmysql` 提供了 MySQL 数据库的自动连接管理。

*   **Init(dbName string)**: 初始化组件，指定需要的数据库名称。
*   **DB() (*sql.DB, error)**: 获取 `*sql.DB` 连接对象。
    *   **Lazy Loading**: 第一次调用时，SDK 会向 Glow Server 发起资源申请。
    *   **JIT Provisioning**: Server 会自动创建数据库和用户，并返回连接凭据 (DSN)。
    *   **Auto Config**: 开发者无需手动配置 DSN，一切由 SDK 和 Server 自动协商。

```go
glowmysql.Init("billing_db")
// ...
db, err := glowmysql.DB()
if err != nil {
    log.Fatal("Database not available")
}
```

### 3.4 GlowRedis (缓存)

`glowredis` 提供了 Redis 客户端的自动连接管理，基于 `go-redis/v9`。

*   **Client() (*redis.Client, error)**: 获取 Redis 客户端。
    *   同样支持 Lazy Loading 和 JIT Provisioning。
    *   默认申请名为 `<app-name>_redis` 的资源。

```go
rdb, err := glowredis.Client()
if err != nil {
    log.Fatal("Redis not available")
}
```

### 3.5 GlowConfig (配置管理)

`starter/glowapp/config` 提供了统一的配置接口，支持本地文件 (`local_config.json`) 和配置中心 (Glow Server) 双源加载。

*   **Get(key string, target interface{})**: 获取配置项。
*   **动态更新**: 当在 Server 端修改配置时，SDK 会通过 TCP 长连接收到推送，并自动更新内存中的配置。

## 4. 运行机制

### 4.1 隐式连接 (Implicit Wiring)

Glow 提倡 "Implicit Wiring" (隐式连接) 理念。开发者只需在代码中声明 "我需要一个叫 'users' 的数据库"，而不需要关心数据库的 IP、端口、用户名和密码。

1.  **Code**: `glowmysql.Init("users")`
2.  **SDK**: 向 Server 发送 `ProvisionRequest{Type: "mysql", Name: "users"}`
3.  **Server**: 检查资源池 -> 创建 DB/User -> 返回 DSN
4.  **SDK**: 使用 DSN 连接数据库 -> 返回 `*sql.DB`

### 4.2 本地调试 vs 托管运行

*   **托管运行 (推荐)**: 使用 `glow-server` 启动应用。Server 会自动注入环境变量和配置。
*   **本地调试**: 直接 `go run main.go`。
    *   如果本地运行了 `glow-server`，SDK 会尝试连接它进行注册和资源申请。
    *   如果连接失败，SDK 会降级读取当前目录下的 `local_config.json`。
