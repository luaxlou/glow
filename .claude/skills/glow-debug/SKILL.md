---
name: glow-debug
description: Glow 应用本地开发与调试指南。包含 glow-server 设置、本地开发环境配置、热重载、IDE 集成、调试技巧、性能分析和测试。当用户需要：本地开发 glow 应用、配置开发环境、调试代码、设置 IDE、运行测试时使用。
---

# Glow Debug

Glow 应用的本地开发、调试和测试指南。

## 环境准备

### 安装 glow-server

```bash
# 本地安装（macOS/Linux）
curl -fsSL "https://raw.githubusercontent.com/luaxlou/glow/main/scripts/install-local.sh" | bash
```

### 启动 glow-server

```bash
# 前台启动（本地开发）
glow-server serve
```

输出：
```
Glow Server starting...
HTTP API: http://localhost:32102
App Center: localhost:32101
```

保持此终端运行。

## 开发工作流

### 方式一：直接运行（推荐）

```bash
# 在项目目录
cd my-glow-app

# 直接运行应用
go run main.go
```

SDK 会自动：
1. 连接到 glow-server (127.0.0.1:32101)
2. 注册应用身份
3. 获取配置
4. 申请资源（MySQL、Redis）

### 方式二：手动分离

```bash
# 终端 1: 启动 glow-server
glow-server serve

# 终端 2: 运行应用
cd my-glow-app
go run main.go
```

## 本地配置文件

当 glow-server 不可用时，SDK 降级使用本地配置。

创建 `local_config.json`：

```json
{
  "debug": true,
  "log_level": "debug",
  "database": {
    "host": "localhost",
    "port": 3306,
    "user": "root",
    "password": "",
    "name": "my_app_db"
  },
  "redis": {
    "host": "localhost",
    "port": 6379,
    "db": 0
  }
}
```

## 热重载开发

### 使用 air（推荐）

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 创建 .air.toml
cat > .air.toml << 'EOF'
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o tmp/main ."
bin = "tmp/main"
include_ext = ["go"]
exclude_dir = ["tmp", "vendor"]
delay = 1000
stop_on_error = true
EOF

# 启动热重载
air
```

### 使用 gow

```bash
go install github.com/mitranim/gow@latest
gow
```

## IDE 配置

### VS Code

创建 `.vscode/launch.json`：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Glow App",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}",
      "env": {
        "OP_APP_NAME": "my-app",
        "OP_APP_PORT": "8080",
        "OP_SERVER_URL": "127.0.0.1:32101"
      },
      "preLaunchTask": "start-glow-server"
    }
  ]
}
```

创建 `.vscode/tasks.json`：

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "start-glow-server",
      "type": "shell",
      "command": "glow-server serve",
      "isBackground": true,
      "problemMatcher": []
    }
  ]
}
```

### GoLand / IntelliJ IDEA

1. 打开 Run Configuration
2. 设置 Environment Variables：
   ```
   OP_APP_NAME=my-app
   OP_APP_PORT=8080
   OP_SERVER_URL=127.0.0.1:32101
   ```
3. 添加 Before Launch Task 运行 `glow-server serve`

## 调试技巧

### 1. 查看应用注册信息

```go
import "github.com/luaxlou/glow/starter/glowapp"

func main() {
    glowapp.Init("debug-app")

    // 输出应用信息
    fmt.Printf("App Name: %s\n", glowapp.Name())
    fmt.Printf("Server URL: %s\n", os.Getenv("OP_SERVER_URL"))
    fmt.Printf("App Port: %s\n", os.Getenv("OP_APP_PORT"))
}
```

### 2. 配置调试

```go
import "github.com/luaxlou/glow/starter/glowapp/config"

type Config struct {
    Debug bool `json:"debug"`
}

var cfg Config
if err := config.Get("app", &cfg); err != nil {
    log.Printf("Failed to load config: %v", err)
} else {
    log.Printf("Loaded config: %+v\n", cfg)
}
```

### 3. 数据库连接调试

```go
db, err := glowmysql.Gorm()
if err != nil {
    log.Printf("MySQL connection failed: %v", err)
} else {
    sqlDB, _ := db.DB()
    if err := sqlDB.Ping(); err != nil {
        log.Printf("MySQL ping failed: %v", err)
    } else {
        log.Println("MySQL connected successfully")
    }
}
```

### 4. HTTP 调试端点

```go
r := glowhttp.Router()
r.GET("/debug", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "app": glowapp.Name(),
        "env": os.Environ(),
        "config": cfg,
    })
})
```

## 测试

### 单元测试

```go
package main

import (
    "testing"
    "github.com/luaxlou/glow/starter/glowapp"
)

func TestAppInit(t *testing.T) {
    glowapp.Init("test-app")
    if glowapp.Name() != "test-app" {
        t.Errorf("Expected test-app, got %s", glowapp.Name())
    }
}
```

### 集成测试

```go
func TestWithGlowServer(t *testing.T) {
    // 启动测试用的 glow-server（通过 exec.Command）
    // 或使用 mock
}
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test -run TestAppInit

# 测试覆盖率
go test -cover ./...
```

## 性能分析

### CPU 性能

```bash
# 启用 CPU profiling
go run main.go -cpuprofile=cpu.prof

# 分析
go tool pprof cpu.prof
```

### 内存性能

```bash
# 启用内存 profiling
go run main.go -memprofile=mem.prof

# 分析
go tool pprof mem.prof
```

### HTTP 基准测试

```bash
# 使用 wrk
wrk -t4 -c100 -d30s http://localhost:8080/api

# 使用 ab (Apache Bench)
ab -n 1000 -c 100 http://localhost:8080/api
```

## 常见问题

### glow-server 未启动

**错误：**
```
Failed to connect to Glow Server: dial tcp 127.0.0.1:32101: connect: connection refused
```

**解决：**
```bash
# 启动 glow-server
glow-server serve

# 或使用 local_config.json 降级运行
```

### 端口冲突

**错误：**
```
bind: address already in use
```

**解决：**
```bash
# 查找占用端口的进程
lsof -i :8080

# 或让 glow-server 自动分配端口
export OP_APP_PORT=""
```

### 配置未加载

**检查：**
```bash
# 查看 glow-server 中的配置
glow config view my-app

# 确认配置格式正确（JSON）
cat local_config.json | jq
```

### MySQL 连接失败

**检查：**
```bash
# 查看 glow-server info
glow-server info

# 检查 MySQL 资源
glow get resources

# 测试连接
mysql -u root -p -e "SHOW DATABASES;"
```

## 生产模拟

### 本地测试部署流程

```bash
# 1. 编译
go build -o my-app

# 2. 部署到本地 glow-server
glow deploy ./my-app

# 3. 查看状态
glow get apps

# 4. 测试 API
curl http://localhost:8080/health

# 5. 查看日志
glow logs my-app
```

### 压力测试

```bash
# 启动应用
glow deploy ./my-app

# 压力测试
wrk -t4 -c100 -d30s http://localhost:8080/api

# 查看性能指标
glow describe app my-app
```

## 资源清理

### 清理测试应用

```bash
# 删除所有测试应用
glow get apps
glow delete app test-app-1
glow delete app test-app-2
```

### 清理数据库

```bash
# 连接 MySQL
mysql -u root -p

# 删除测试数据库
DROP DATABASE test_db;
```

### 重置 glow-server

```bash
# 停止 glow-server
# 删除数据目录
rm -rf ~/Library/Application\ Support/glow-server/  # macOS
rm -rf ~/.local/share/glow-server/                  # Linux

# 重新启动
glow-server serve
```

## 最佳实践

1. **使用 glow-server**: 本地开发时保持 glow-server 运行
2. **热重载**: 使用 air 或 gow 提高开发效率
3. **环境隔离**: 开发使用本地配置，生产使用 glow-server
4. **调试端点**: 添加 `/debug` 端点查看应用状态
5. **日志输出**: 使用 `log.Printf()` 输出调试信息
6. **测试驱动**: 编写测试确保代码质量

## 调试清单

开发前检查：
- [ ] glow-server 是否运行
- [ ] 端口是否可用（8080 或由 glow-server 分配）
- [ ] MySQL/Redis 资源是否配置
- [ ] local_config.json 是否准备（降级用）

开发时：
- [ ] 使用 `glow logs` 查看日志
- [ ] 使用热重载加速开发
- [ ] 定期运行测试

部署前：
- [ ] 编译二进制文件
- [ ] 测试部署流程
- [ ] 验证配置正确
- [ ] 检查日志无错误
