# Glow SDK 用户手册

`glow/starter` 是 Glow 框架的 Go 语言 SDK，提供了开箱即用的应用骨架和组件。通过使用 SDK，开发者可以轻松地将应用接入 Glow 的配置管理和资源绑定。

## 1. 安装

在你的 Go 项目中引入 Glow SDK：

```bash
go get github.com/luaxlou/glow/starter
```

## 2. 快速开始 (Quick Start)

以下是一个最简单的 Web 应用示例，展示了如何使用 Glow 配置管理：

```go
package main

import (
    "fmt"
    "github.com/luaxlou/glow/starter/glowconfig"
)

func main() {
    // 1. 加载配置
    config, err := glowconfig.Load()
    if err != nil {
        fmt.Printf("Failed to load config: %v\n", err)
        return
    }

    // 2. 使用配置
    logLevel := config.GetString("log_level")
    maxConnections := config.GetInt("max_connections")

    fmt.Printf("Log Level: %s\n", logLevel)
    fmt.Printf("Max Connections: %d\n", maxConnections)

    // 3. 应用逻辑...
}
```

## 3. 核心组件

Glow 提供了简单的配置管理功能。

### 3.1 GlowConfig (配置管理)

`pkg/glowconfig` 提供了统一的配置接口，从本地配置文件读取配置。

*   **Load() (Config, error)**: 从工作目录加载 `config.json` 文件。
*   **LoadFromFile(path string) (Config, error)**: 从指定路径加载配置文件。
*   **Get(key string) interface{}**: 获取配置项值（支持点号分隔的嵌套键）。
*   **GetString(key string) string**: 获取字符串值。
*   **GetInt(key string) int**: 获取整数值。

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

### 3.6 GlowConfig (配置文件)

配置文件路径：`<data-dir>/apps/<appName>/config.json`

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
* Server 生成配置文件（`config.json`）
* 应用启动时只读取本地配置，**不连接** Server

### 4.2 本地调试

直接 `go run main.go`：
* 应用读取当前工作目录下的 `config.json`
* 如需绑定资源，先运行 `glow apply -f app.yaml` 生成配置文件
