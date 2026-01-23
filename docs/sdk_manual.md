# Glow SDK 用户手册

`glow/starter` 是 Glow 框架的 Go 语言 SDK，提供了开箱即用的应用骨架和组件。通过使用 SDK，开发者可以轻松地将应用接入 Glow 的配置管理和资源绑定。

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

`glowapp` 是应用的入口，负责生命周期管理。

*   **Init(appName string, opts ...Option)**: 初始化应用身份。
*   **WaitForShutdown()**: 阻塞主 goroutine，直到接收到系统信号 (SIGINT/SIGTERM)。收到信号后，它会通知所有已注册的组件执行清理操作 (Graceful Shutdown)。
*   **RegisterCleanup(name string, fn func())**: 注册自定义的清理函数，在停机时执行。

### 3.2 GlowHTTP (Web 服务)

`glowhttp` 基于 [Gin](https://github.com/gin-gonic/gin) 框架封装，集成了 Glow 的端口管理和生命周期。

*   **Init(defaultPort int)**: 初始化 HTTP 组件。如果 Glow Server 指定了 `OP_APP_PORT` 环境变量，则忽略 `defaultPort`。
*   **Router() *gin.Engine**: 获取 Gin 实例，用于注册路由和中间件。
*   **Run()**: 异步启动 HTTP Server。

### 3.3 GlowConfig (配置管理)

`starter/glowapp/config` 提供了统一的配置接口，从本地配置文件读取配置。

*   **Get(key string, target interface{})**: 获取配置项（从 `<appName>_local_config.json` 读取）。
*   **IsSet(key string) bool**: 检查配置项是否存在。

**声明式配置管理**:
- 配置通过 `app.yaml` 的 `spec.config` 字段声明
- 执行 `glow apply -f app.yaml` 生成本地配置文件
- 配置文件包含：
  - 用户在 `spec.config` 中声明的自定义配置
  - 资源绑定自动生成的连接信息（MySQL DSN、Redis addr 等）

示例 - 在 app.yaml 中声明配置：

```yaml
apiVersion: v1
kind: App
metadata:
  name: my-app
spec:
  config:
    log_level: debug
    max_connections: 100
    feature_flags:
      new_ui: true

```

应用配置文件由 `glow apply` 命令生成，包含应用元数据和资源绑定信息（如 MySQL DSN、Redis 地址等）。

### 3.4 GlowMySQL (数据库 - 基于 GORM)

`glowmysql` 提供了 MySQL 数据库访问能力，底层基于 [GORM](https://gorm.io)。

*   **Init(dbName string)**: 声明应用希望使用的数据库名（用于配置查找）。
*   **Gorm() (*gorm.DB, error)**: 返回单例的 `*gorm.DB`，从本地配置读取 `mysql.dsn`。
*   **DB() (*sql.DB, error)**: 在需要原生 `*sql.DB` 的场景下使用，内部复用同一连接。

**重要**: 使用前需先通过 `glow apply -f app.yaml` 绑定 MySQL 资源，在 YAML 中声明：
```yaml
spec:
  resources:
    mysql:
      - dbName: myapp_db
```

示例：

```go
import (
    "log"

    "github.com/luaxlou/glow/starter/glowapp"
    "github.com/luaxlou/glow/starter/glowhttp"
    "github.com/luaxlou/glow/starter/glowmysql"
)

func main() {
    glowapp.Init("my-gorm-app")

    glowhttp.Init(8080)
    glowmysql.Init("my_gorm_app_db")  // 声明数据库名

    db, err := glowmysql.Gorm()  // 从配置读取 mysql.dsn
    if err != nil {
        log.Fatalf("init mysql via gorm failed: %v", err)
    }

    type User struct {
        ID   uint
        Name string
    }

    if err := db.AutoMigrate(&User{}); err != nil {
        log.Fatalf("auto migrate failed: %v", err)
    }

    // 继续使用 db 进行 CRUD ...
}
```

### 3.5 GlowRedis (缓存)

`glowredis` 提供了 Redis 缓存访问能力，底层基于 [go-redis/v9](https://github.com/redis/go-redis)。

*   **Client() (*redis.Client, error)**: 返回单例的 `*redis.Client`，从本地配置读取 Redis 连接信息。

**重要**: 使用前需先通过 `glow apply -f app.yaml` 绑定 Redis 资源，在 YAML 中声明：
```yaml
spec:
  resources:
    redis:
      - db: 0
```

### 3.6 GlowConfig (配置文件)

配置文件路径：`<data-dir>/apps/<appName>/<appName>_local_config.json`

示例配置结构：
```json
{
  "mysql": {
    "dsn": "user:pass@tcp(localhost:3306)/dbname"
  },
  "redis": {
    "addr": "localhost:6379",
    "password": "",
    "db": 0
  }
}
```

## 4. 运行机制

### 4.1 托管运行

应用由 `glow-server` 启动（`glow start app <name>`）：
* Server 设置环境变量（`OP_APP_NAME`, `OP_APP_PORT`, `OP_APP_DOMAIN`）
* Server 生成配置文件（`<appName>_local_config.json`）
* 应用启动时只读取本地配置，**不连接** Server

### 4.2 本地调试

直接 `go run main.go`：
* SDK 读取当前目录下的 `<appName>_local_config.json`
* 如需绑定资源，先运行 `glow apply -f app.yaml` 生成配置文件
