---
name: glow-sdk
description: Glow Go SDK 开发指南。包含 GlowApp、GlowHTTP、GlowMySQL、GlowRedis、GlowConfig、GlowWebSocket 等 SDK 组件的使用方法、API 参考和代码示例。当用户需要：使用 Glow SDK 开发 Go 应用、集成数据库、配置管理、使用 HTTP/WebSocket 组件时使用。
---

# Glow SDK

Glow SDK 是 Go 应用的核心开发框架，提供开箱即用的组件和标准化的应用骨架。

## 安装

```bash
go get github.com/luaxlou/glow/starter
```

## 核心组件

### GlowApp - 应用核心

应用身份注册和生命周期管理。

```go
import "github.com/luaxlou/glow/starter/glowapp"

// 初始化应用
glowapp.Init("app-name")

// 获取应用名称
name := glowapp.Name()

// 注册清理函数（停机时执行）
glowapp.RegisterCleanup("cleanup-name", func() {
    // 清理逻辑
})

// 等待停机信号
glowapp.WaitForShutdown()
```

**环境变量（自动注入）：**
- `OP_APP_NAME`: 应用名称
- `OP_APP_PORT`: 分配的端口
- `OP_SERVER_URL`: Glow Server 地址

### GlowHTTP - Web 服务

基于 Gin 框架的 HTTP 服务。

```go
import "github.com/luaxlou/glow/starter/glowhttp"

// 初始化（默认端口 8080，会被 OP_APP_PORT 覆盖）
glowhttp.Init(8080)

// 获取 Router
r := glowhttp.Router()

// 注册路由
r.GET("/", func(c *gin.Context) {
    c.String(200, "Hello Glow!")
})

// 启动服务（异步）
glowhttp.Run()

// 等待停机
glowapp.WaitForShutdown()
```

### GlowConfig - 配置管理

统一配置接口，支持本地文件和配置中心双源加载。

```go
import "github.com/luaxlou/glow/starter/glowapp/config"

// 定义配置结构
type AppConfig struct {
    Debug     bool   `json:"debug"`
    LogPath   string `json:"log_path"`
    Database  struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    } `json:"database"`
}

// 加载配置
var cfg AppConfig
if err := config.Get("app", &cfg); err != nil {
    log.Fatalf("Failed to load config: %v", err)
}

// 使用配置
fmt.Printf("Debug mode: %v\n", cfg.Debug)
```

**配置优先级：**
1. Glow Server 配置中心（动态热更新）
2. 本地 `local_config.json`

### GlowMySQL - 数据库（GORM）

基于 GORM 的 MySQL 访问，自动资源申请。

```go
import "github.com/luaxlou/glow/starter/glowmysql"

// 初始化（声明数据库名）
glowmysql.Init("my_app_db")

// 获取 GORM 实例
db, err := glowmysql.Gorm()
if err != nil {
    log.Fatalf("Failed to init MySQL: %v", err)
}

// 使用 GORM
type User struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"size:255"`
    Email string `gorm:"size:255"`
}

// 自动迁移
db.AutoMigrate(&User{})

// CRUD 操作
user := User{Name: "John", Email: "john@example.com"}
db.Create(&user)

var users []User
db.Find(&users)
```

### GlowRedis - Redis 客户端

```go
import "github.com/luaxlou/glow/starter/glowredis"

// 初始化（使用默认数据库 0）
glowredis.Init()

// 或者指定数据库
glowredis.InitDB(1)

// 获取客户端
client, err := glowredis.Client()
if err != nil {
    log.Fatalf("Failed to init Redis: %v", err)
}

// 使用 Redis
ctx := context.Background()
err = client.Set(ctx, "key", "value", 0).Err()
if err != nil {
    panic(err)
}

val, err := client.Get(ctx, "key").Result()
fmt.Println("key", val)
```

### GlowWebSocket - WebSocket 支持

```go
import "github.com/luaxlou/glow/starter/glowwebsocket"

r := glowhttp.Router()
r.GET("/ws", func(c *gin.Context) {
    glowwebsocket.Handle(c, func(conn *websocket.Conn) {
        for {
            // 读取消息
            messageType, p, err := conn.ReadMessage()
            if err != nil {
                return
            }

            // 处理消息...

            // 回复消息
            if err := conn.WriteMessage(messageType, p); err != nil {
                return
            }
        }
    })
})
```

## 完整示例

### 最小化应用

```go
package main

import (
    "github.com/luaxlou/glow/starter/glowapp"
    "github.com/luaxlou/glow/starter/glowhttp"
)

func main() {
    glowapp.Init("minimal-app")
    glowhttp.Init(8080)

    r := glowhttp.Router()
    r.GET("/", func(c *gin.Context) {
        c.String(200, "Hello Glow!")
    })

    glowhttp.Run()
    glowapp.WaitForShutdown()
}
```

### 完整应用（HTTP + MySQL + Redis）

```go
package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/luaxlou/glow/starter/glowapp"
    "github.com/luaxlou/glow/starter/glowhttp"
    "github.com/luaxlou/glow/starter/glowmysql"
    "github.com/luaxlou/glow/starter/glowredis"
)

func main() {
    // 初始化应用
    glowapp.Init("fullstack-app")

    // 初始化组件
    glowhttp.Init(8080)
    glowmysql.Init("app_db")
    glowredis.Init()

    // 获取 MySQL 连接
    db, err := glowmysql.Gorm()
    if err != nil {
        log.Printf("MySQL not available: %v", err)
    } else {
        log.Println("MySQL connected!")
    }

    // 获取 Redis 连接
    client, err := glowredis.Client()
    if err != nil {
        log.Printf("Redis not available: %v", err)
    } else {
        log.Println("Redis connected!")
        _ = client
    }

    // 注册路由
    r := glowhttp.Router()

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "app":    glowapp.Name(),
        })
    })

    r.GET("/users", func(c *gin.Context) {
        if db != nil {
            // 查询用户...
            c.JSON(200, gin.H{"users": []string{}})
        } else {
            c.JSON(503, gin.H{"error": "Database not available"})
        }
    })

    // 启动服务
    glowhttp.Run()
    glowapp.WaitForShutdown()
}
```

## 最佳实践

### 1. 组件初始化顺序

```go
func main() {
    // 1. 先初始化应用身份
    glowapp.Init("my-app")

    // 2. 初始化基础组件（HTTP、数据库等）
    glowhttp.Init(8080)
    glowmysql.Init("my_db")
    glowredis.Init()

    // 3. 获取连接并使用
    db, _ := glowmysql.Gorm()
    // ...

    // 4. 注册路由
    r := glowhttp.Router()
    // ...

    // 5. 启动服务并等待
    glowhttp.Run()
    glowapp.WaitForShutdown()
}
```

### 2. 错误处理

```go
// MySQL 可选
db, err := glowmysql.Gorm()
if err != nil {
    log.Printf("Warning: MySQL not available: %v. Running without DB.", err)
    // 继续运行，不致命
}

// 或者必需
db, err := glowmysql.Gorm()
if err != nil {
    log.Fatalf("Failed to init MySQL: %v", err)
    // 退出
}
```

### 3. 配置管理

```go
type Config struct {
    Debug bool   `json:"debug"`
    DBName string `json:"db_name"`
}

var cfg Config
if err := config.Get("app", &cfg); err != nil {
    // 使用默认值
    cfg = Config{Debug: false, DBName: "default_db"}
}

// 使用配置
if cfg.Debug {
    gin.SetMode(gin.DebugMode)
}
```

### 4. 优雅停机

```go
// 注册清理函数
glowapp.RegisterCleanup("close-db", func() {
    if db != nil {
        sqlDB, _ := db.DB()
        sqlDB.Close()
        log.Println("Database closed")
    }
})

glowapp.RegisterCleanup("flush-cache", func() {
    // 清理缓存...
    log.Println("Cache flushed")
})

// 正常停机
glowapp.WaitForShutdown()
```

## API 参考

详见 [API Reference](references/api-reference.md)

## 常见问题

### Q: 如何在本地调试？
A: 确保 glow-server 运行，然后 `go run main.go`。SDK 会自动连接。

### Q: 配置如何热更新？
A: 通过 `glow config edit` 修改配置，SDK 会自动收到推送并更新。

### Q: 数据库连接失败怎么办？
A: 检查 glow-server 的 MySQL 资源配置，使用 `glow-server info` 查看。
