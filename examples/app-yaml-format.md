# Glow App.yaml 配置格式

## 约定大于配置原则

Glow 采用"约定大于配置"的设计理念，最小化必须配置的字段。

## 必填字段

只有 `metadata.name` 是必填的：

```yaml
apiVersion: v1
kind: App
metadata:
  name: my-app  # 必填：应用名称
spec: {}
```

## 可选字段与默认值

### 1. binary（可选）

**默认值**: `<data-dir>/apps/<app-name>/<app-name>`

**示例**:
- 指定: `binary: ./my-app`
- 省略: 自动使用 `/var/lib/glow-server/apps/my-app/my-app`

### 2. workingDir（可选）

**默认值**: `<data-dir>/apps/<app-name>`

**示例**:
- 指定: `workingDir: /path/to/dir`
- 省略: 自动使用 `/var/lib/glow-server/apps/my-app`

### 3. args（可选）

**默认值**: `[]`（空列表）

**示例**:
```yaml
args:
  - "--server"
  - "--port=8080"
```

### 4. env（可选）

**默认值**: `{}`（空 map）

**示例**:
```yaml
env:
  - name: ENV
    value: production
  - name: LOG_LEVEL
    value: info
```

### 5. port（可选）

**默认值**: `0`（不开放端口）

**说明**:
- 未指定时，应用不开放 HTTP 端口
- 指定后，glow-server 注入 `OP_APP_PORT` 环境变量

### 6. domain（可选）

**默认值**: 空（不绑定域名）

**约束**: 指定 `domain` 时必须同时指定 `port`

**说明**:
- 指定后，glow-server 自动配置 Nginx 反向代理

### 7. resources（可选）

**默认值**: 不绑定任何资源

**支持**:
- `mysql`: MySQL 数据库列表
- `redis`: Redis 缓存列表

## 最小化配置示例

### 最小配置

```yaml
apiVersion: v1
kind: App
metadata:
  name: my-app
spec:
  # binary 和 workingDir 会自动使用约定值
```

### 带端口的配置

```yaml
apiVersion: v1
kind: App
metadata:
  name: my-app
spec:
  port: 33203
```

### 带域名和资源的配置

```yaml
apiVersion: v1
kind: App
metadata:
  name: my-app
spec:
  port: 33203
  domain: myapp.example.com
  resources:
    mysql:
      - dbName: myapp_db
    redis:
      - db: 0
```

### 完整配置

```yaml
apiVersion: v1
kind: App
metadata:
  name: my-awesome-app
spec:
  # 执行配置（可选）
  binary: ./my-app          # 可选，默认: <data-dir>/apps/<name>/<name>
  workingDir: /var/app       # 可选，默认: <data-dir>/apps/<name>
  args:
    - "--server"
    - "--port=8080"

  # 环境变量（可选）
  env:
    - name: ENV
      value: production

  # HTTP 服务（可选）
  port: 33203               # 可选，默认: 0（不开放端口）

  # Ingress（可选）
  domain: myapp.example.com # 可选，需要同时指定 port

  # 资源（可选）
  resources:
    mysql:
      - dbName: myapp_db
    redis:
      - db: 0
```

## 部署流程

### 1. 准备二进制文件

```bash
# 编译应用
go build -o my-app main.go

# 放到约定位置（与 spec 中应用名同名）
mkdir -p /var/lib/glow-server/apps/my-app
cp my-app /var/lib/glow-server/apps/my-app/my-app
```

### 2. 应用配置

```bash
# 应用配置（不启动应用）
glow apply -f app.yaml
```

### 3. 启动应用

```bash
# 启动应用
glow start app my-app
```

### 4. 查看状态

```bash
# 查看应用状态
glow get app my-app
```

## 配置文件生成

执行 `glow apply` 后，会在以下位置生成配置文件：

**路径**: `<data-dir>/apps/<app-name>/<app-name>_local_config.json`

**内容示例**:
```json
{
  "mysql": {
    "dsn": "user:pass@tcp(localhost:3306)/myapp_db"
  },
  "redis": {
    "addr": "localhost:6379",
    "username": "",
    "password": "",
    "db": 0
  }
}
```

## 应用读取配置

应用使用 SDK 读取本地配置文件：

```go
import (
    "github.com/yourusername/glow/starter/glowapp/config"
    "github.com/yourusername/glow/starter/glowmysql"
    "github.com/yourusername/glow/starter/glowredis"
)

func main() {
    // 启动配置管理（只读本地文件）
    config.Start()

    // 初始化 MySQL
    glowmysql.Init("myapp_db")

    // 初始化 Redis
    glowredis.Init()

    // 启动 HTTP 服务
    glowhttp.Init(33203)

    // 应用逻辑...
}
```

## 约定总结

| 字段 | 默认值 | 说明 |
|------|--------|------|
| `spec.binary` | `<data-dir>/apps/<name>/<name>` | 二进制文件路径 |
| `spec.workingDir` | `<data-dir>/apps/<name>` | 工作目录 |
| `spec.args` | `[]` | 命令行参数 |
| `spec.env` | `{}` | 环境变量 |
| `spec.port` | `0` | HTTP 端口（0=不开放） |
| `spec.domain` | 空 | 域名绑定 |
| `spec.resources` | 不绑定 | 资源绑定 |

## 优势

1. **简化配置**: 大多数应用只需要 3-5 行配置
2. **一致性**: 所有应用遵循相同的目录结构
3. **可预测**: 配置路径和行为都是可预测的
4. **灵活性**: 需要时可以覆盖默认值
