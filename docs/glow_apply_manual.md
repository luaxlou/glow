# Glow Apply 用户手册

## 概述

`glow apply` 是 Glow 框架中**唯一的应用配置方式**。通过声明式 YAML 文件，你可以配置应用的所有方面（执行参数、环境变量、端口、域名、应用配置等）。

Glow 采用"配置即代码"的设计理念：所有配置在 `spec.config` 字段中声明，用户自行提供数据库连接等基础设施配置，系统只负责配置文件生成和管理。

## 核心概念

### 声明式配置

所有配置都在 YAML 文件中声明，而不是通过命令行参数或交互式配置：

```yaml
apiVersion: v1
kind: App
metadata:
  name: my-app
spec:
  port: 8080          # 端口
  domain: app.local   # 域名
  config:             # 应用配置 ⭐
    mysql_dsn: "user:pass@tcp(localhost:3306)/mydb"
    redis_addr: "localhost:6379"
    log_level: "info"
```

**配置即代码原则**：
- 所有配置（包括数据库连接、缓存配置等）都在 `spec.config` 中声明
- 用户自行提供数据库连接等基础设施配置
- 系统将配置写入本地 JSON 文件供应用读取

### 幂等性

可以重复执行 `glow apply`，每次执行都会：
- 更新应用元数据
- 重新生成配置文件
- 如果应用正在运行且配置变化，自动重启应用

### 原子性

所有配置在一次 `apply` 中完成，要么全部成功，要么全部失败。

## 命令语法

```bash
glow apply -f <filename>
```

**必需参数**:
- `-f, --file string`: YAML 配置文件路径

**示例**:
```bash
glow apply -f app.yaml
glow apply -f /path/to/config.yaml
```

## YAML 文件格式

### 基本结构

```yaml
apiVersion: v1        # API 版本（必需）
kind: App              # 资源类型（必需）
metadata:              # 元数据（必需）
  name: app-name       # 应用名称（必需）
spec:                  # 规格说明（必需）
  # 配置项...
```

### 完整示例

```yaml
apiVersion: v1
kind: App
metadata:
  name: my-web-app
spec:
  # 执行配置（可选，有约定值）
  binary: ./my-app              # 应用二进制路径
  workingDir: /path/to/work     # 工作目录
  args: ["--server", "--v"]     # 启动参数
  env:                          # 环境变量
    - name: ENV
      value: production
    - name: LOG_LEVEL
      value: info

  # 网络配置
  port: 8080                    # HTTP 端口（可选）
  domain: myapp.example.com     # 域名绑定（可选）

  # 应用配置（所有配置在这里声明）⭐
  config:
    # 数据库配置（用户自行提供）
    mysql_dsn: "user:pass@tcp(localhost:3306)/myapp_db"

    # 缓存配置（用户自行提供）
    redis_addr: "localhost:6379"
    redis_password: ""

    # 应用配置
    log_level: "info"
    max_connections: 100

    # 功能开关
    feature_new_ui: true
```

## 字段说明

### metadata.name（必需）

应用的唯一标识符。

**规则**:
- 只能包含小写字母、数字和连字符
- 必须以字母开头
- 长度 2-63 个字符

**示例**:
```yaml
metadata:
  name: my-app        # ✅ 正确
  name: my_app        # ❌ 错误（包含下划线）
  name: 123app        # ❌ 错误（以数字开头）
```

### spec.binary（可选）

应用二进制文件的路径。

**默认值**: `<data-dir>/apps/<app-name>/<app-name>`

**相对路径**: 相对于工作目录
**绝对路径**: 从根目录开始

**约定大于配置**: 大多数情况下可以省略此字段，系统会自动使用约定路径。

**示例**:
```yaml
spec:
  binary: ./my-app           # 相对路径
  binary: /usr/bin/my-app    # 绝对路径
  # binary:                  # 省略时使用约定值
```

### spec.workingDir（可选）

应用的工作目录。如果未指定，默认为 `/var/lib/glow-server/apps/<app-name>`。

**示例**:
```yaml
spec:
  workingDir: /var/lib/glow-server/apps/my-app
```

### spec.port（可选）

应用监听的 HTTP 端口。

**重要**:
- 如果不指定 `port`，应用不会对外开放端口
- 如果不指定 `port`，则不能指定 `domain`
- 端口号会通过 `OP_APP_PORT` 环境变量注入到应用

**示例**:
```yaml
spec:
  port: 8080         # 对外开放 8080 端口
  # port: 0         # 不对外开放端口
```

### spec.domain（可选）

Ingress 域名绑定。

**规则**:
- 必须同时指定 `port`
- glow-server 会自动配置 Nginx 反向代理

**示例**:
```yaml
spec:
  port: 8080
  domain: myapp.example.com    # ✅ 正确
```

```yaml
spec:
  # port: 8080    # ❌ 错误：必须指定 port
  domain: myapp.example.com
```

### spec.args（可选）

应用启动参数。

**格式**: 字符串数组

**示例**:
```yaml
spec:
  args: ["--server", "--port=8080", "--v"]
```

等价于命令行:
```bash
./my-app --server --port=8080 --v
```

### spec.env（可选）

环境变量。

**格式**: 键值对列表

**示例**:
```yaml
spec:
  env:
    - name: DATABASE_URL
      value: "mysql://..."
    - name: LOG_LEVEL
      value: debug
    - name: PORT
      value: "8080"
```

### spec.config（可选）⭐

应用配置，包括数据库连接、缓存配置、应用参数等所有配置项。

**默认值**: `{}`（空 map）

**重要**: 这是声明所有应用配置的唯一方式。

**示例**:
```yaml
spec:
  config:
    # 数据库配置（用户自行提供）
    mysql_dsn: "user:pass@tcp(localhost:3306)/mydb"

    # Redis 配置（用户自行提供）
    redis_addr: "localhost:6379"
    redis_password: ""
    redis_db: 0

    # 应用配置
    log_level: "debug"
    max_connections: 100
    timeout: 30

    # 功能开关
    feature_new_ui: true
    feature_cache: false

    # 自定义配置
    app_name: "My App"
    admin_email: "admin@example.com"
```

**配置文件生成**:
执行 `glow apply` 后，配置会写入到：
```
<data-dir>/apps/<app-name>/config.json
```

**配置来源**:
1. 完全来自 `app.yaml` 中的 `spec.config` 字段
2. 用户自行提供所有配置值（包括数据库连接等）
3. 系统不提供自动化的配置生成或资源创建

**应用读取配置**:
应用使用 SDK 读取本地配置文件：
```go
import "github.com/luaxlou/glow/starter/glowconfig"

func main() {
    // 读取配置
    config, err := glowconfig.Get()
    if err != nil {
        // 处理错误
    }

    mysqlDSN := config.GetString("mysql_dsn")
    redisAddr := config.GetString("redis_addr")
    logLevel := config.GetString("log_level")

    // 使用配置...
}
```

## 执行流程

`glow apply` 按以下顺序执行：

### 1. 验证 YAML 文件

检查必需字段和格式。

### 2. 注册/更新应用

调用 `PUT /apps/:name` API，保存应用元数据到数据库。

### 3. 配置 Ingress（如果指定了 domain）

生成 Nginx 配置文件并重载 Nginx。

**配置文件位置**: `/etc/nginx/sites-available/<app-name>`

### 4. 生成配置文件

将 `spec.config` 中的所有配置写入本地配置文件。

**配置文件位置**: `<data-dir>/apps/<app-name>/config.json`

**示例内容**:
```json
{
  "mysql_dsn": "user:pass@tcp(localhost:3306)/mydb",
  "redis_addr": "localhost:6379",
  "redis_password": "",
  "log_level": "info",
  "max_connections": 100,
  "feature_new_ui": true
}
```

**注意**: 配置内容完全来自 `spec.config`，系统不会自动生成任何配置。

### 5. 自动重启（如果需要）

如果应用正在运行且配置文件发生变化，自动重启应用以使新配置生效。

### 6. 输出摘要

显示操作结果和下一步提示。

## 输出示例

### 成功输出

```
Applying App 'my-app' from app.yaml...
✓ App 'my-app' registered successfully
→ Configuring Ingress for domain: myapp.example.com
✓ Ingress configured: http://myapp.example.com -> port 8080
→ Generating config file...
✓ Config file written to: /var/lib/glow-server/apps/my-app/config.json (245 bytes)

Summary:
  App Name: my-app
  Port: 8080
  Domain: myapp.example.com
  Status: Config updated (no restart needed)

Next steps:
  1. Review the config file generated
  2. If the app is not running, start it: glow start app my-app
  3. Check status: glow get app my-app
```

### 配置变更后自动重启

```
Applying App 'my-app' from app.yaml...
✓ App 'my-app' registered successfully
→ Configuring Ingress for domain: myapp.example.com
✓ Ingress configured: http://myapp.example.com -> port 8080
→ Generating config file...
✓ Config file written to: /var/lib/glow-server/apps/my-app/config.json (256 bytes)
→ Config changed, restarting app 'my-app'...
✓ App 'my-app' restarted successfully

Summary:
  App Name: my-app
  Port: 8080
  Domain: myapp.example.com
  Status: Restarted (config changed)
```

## 错误处理

### 常见错误

#### 1. YAML 格式错误

```
Error parsing YAML: yaml: line 10: mapping values are not allowed in this context
```

**解决**: 检查 YAML 语法，使用 YAML 验证工具。

#### 2. 必需字段缺失

```
Error: metadata.name is required
```

**解决**: 添加必需的字段。

#### 3. Domain 但没有 Port

```
Validation error: spec.port is required when spec.domain is specified
```

**解决**: 删除 `domain` 或添加 `port`。

#### 4. API 路由不存在 (404)

```
Error: server returned status 404
```

**解决**: 确认 glow-server 是最新版本并已重启。

## 使用场景

### 场景 1: 新建 Web 应用

```yaml
apiVersion: v1
kind: App
metadata:
  name: web-app
spec:
  binary: ./web-app
  port: 8080
  domain: myapp.example.com
  config:
    mysql_dsn: "user:pass@tcp(localhost:3306)/webapp_db"
    redis_addr: "localhost:6379"
    log_level: "info"
```

### 场景 2: 后台 Worker

```yaml
apiVersion: v1
kind: App
metadata:
  name: worker
spec:
  binary: ./worker
  # 不指定 port，不对外开放
  config:
    mysql_dsn: "user:pass@tcp(localhost:3306)/worker_db"
    queue_name: "tasks"
```

### 场景 3: 微服务（多应用共享数据库）

```yaml
# api-service.yaml
apiVersion: v1
kind: App
metadata:
  name: api-service
spec:
  binary: ./api-service
  port: 8080
  config:
    mysql_dsn: "user:pass@tcp(localhost:3306)/shared_db"
    redis_addr: "localhost:6379"
    redis_db: 0

# worker.yaml
apiVersion: v1
kind: App
metadata:
  name: worker
spec:
  binary: ./worker
  config:
    mysql_dsn: "user:pass@tcp(localhost:3306)/shared_db"  # 共享同一数据库
    redis_addr: "localhost:6379"
    redis_db: 1  # 不同的 Redis DB
```

### 场景 4: 更新应用配置

```bash
# 1. 编辑 YAML
vim app.yaml

# 2. 应用新配置（如果应用运行且配置变化，会自动重启）
glow apply -f app.yaml

# 3. 检查应用状态
glow get app my-app
```

## 最佳实践

### 1. 版本控制

将 `app.yaml` 纳入 Git 版本控制：

```bash
git add app.yaml
git commit -m "Add app configuration"
```

### 2. 环境分离

为不同环境创建不同的 YAML 文件：

```bash
app.yaml              # 开发环境
app-production.yaml   # 生产环境
app-staging.yaml      # 测试环境
```

### 3. 配置验证

应用前先验证 YAML：

```bash
# 使用 yamllint
yamllint app.yaml

# 或简单的语法检查
python3 -c "import yaml; yaml.safe_load(open('app.yaml'))"
```

### 4. 渐进式更新

更新配置时会自动处理重启：

```bash
# 更新配置（如果应用运行且配置变化，会自动重启）
glow apply -f app.yaml

# 检查生成的配置
cat /var/lib/glow-server/apps/my-app/config.json
```

### 5. 配置审查

应用前审查摘要输出，确认：

- [ ] 端口正确
- [ ] 域名正确
- [ ] 配置文件已生成
- [ ] 应用状态正确

## 高级用法

### 1. 多数据库配置

在 `spec.config` 中配置多个数据库连接：

```yaml
spec:
  config:
    # 主数据库
    mysql_dsn: "user:pass@tcp(localhost:3306)/main_db"

    # 缓存数据库
    mysql_cache_dsn: "user:pass@tcp(localhost:3306)/cache_db"

    # 日志数据库
    mysql_log_dsn: "user:pass@tcp(localhost:3306)/log_db"
```

应用中访问：
```go
config.GetString("mysql_dsn")       // 主数据库
config.GetString("mysql_cache_dsn") // 缓存数据库
config.GetString("mysql_log_dsn")   // 日志数据库
```

### 2. 环境变量注入

```yaml
spec:
  env:
    - name: MYSQL_DSN
      value: "mysql://user:pass@localhost/dbname"
    - name: REDIS_ADDR
      value: "localhost:6379"
```

### 3. 条件配置（使用注释）

```yaml
spec:
  port: 8080
  # domain: myapp.local    # 取消注释以启用域名

  config:
    mysql_dsn: "user:pass@tcp(localhost:3306)/myapp_db"
    # redis_addr: "localhost:6379"  # 取消注释以启用 Redis
```

### 4. 组合配置（包含 args、env 和 config）

```yaml
spec:
  binary: ./my-app
  args: ["--server", "--port=8080"]
  env:
    - name: ENV
      value: production
    - name: PORT
      value: "8080"

  port: 8080
  domain: myapp.local

  config:
    mysql_dsn: "user:pass@tcp(localhost:3306)/myapp_db"
    redis_addr: "localhost:6379"
    log_level: "info"
```

### 5. 配置优先级

配置值的优先级（从高到低）：

1. **环境变量** (spec.env) - 直接注入到进程环境
2. **配置文件** (spec.config) - 写入本地 JSON 配置文件
3. **命令行参数** (spec.args) - 传递给应用的启动参数

示例：
```yaml
spec:
  args: ["--port=8080"]
  env:
    - name: APP_PORT
      value: "8080"
  config:
    port: 8080
```

应用可以这样读取：
- 命令行参数：`os.Args`
- 环境变量：`os.Getenv("APP_PORT")`
- 配置文件：`config.GetInt("port")`

## 与其他命令的配合

### 完整工作流

```bash
# 1. 配置应用
glow apply -f app.yaml

# 2. 启动应用
glow start app my-app

# 3. 查看状态
glow get app my-app

# 4. 查看日志
glow logs my-app

# 5. 更新配置
vim app.yaml
glow apply -f app.yaml
glow restart app my-app

# 6. 停止应用
glow stop my-app

# 7. 删除应用
glow delete app my-app
```

## 故障排查

### 问题 1: apply 成功但应用启动失败

**排查步骤**:
```bash
# 查看应用日志
glow logs my-app

# 检查配置文件
cat /var/lib/glow-server/apps/my-app/config.json

# 手动测试应用
cd /var/lib/glow-server/apps/my-app
./my-app --help
```

**常见原因**:
- 配置文件中的数据库连接字符串不正确
- 端口被占用
- 二进制文件不存在或无执行权限

### 问题 2: 数据库连接失败

**排查步骤**:
```bash
# 检查 MySQL 服务
sudo systemctl status mysql

# 测试 MySQL 连接（使用配置文件中的 DSN）
mysql -u user -p -e "SHOW DATABASES;"

# 检查配置文件中的 DSN
cat /var/lib/glow-server/apps/my-app/config.json | grep mysql_dsn
```

**解决方案**:
- 确认 MySQL 服务运行正常
- 验证 `spec.config.mysql_dsn` 配置正确
- 确认数据库已创建

### 问题 3: Ingress 不工作

**排查步骤**:
```bash
# 检查 Nginx 配置
cat /etc/nginx/sites-available/my-app

# 测试 Nginx 配置
sudo nginx -t

# 重载 Nginx
sudo systemctl reload nginx

# 检查 DNS
ping myapp.local
```

## 相关资源

- **快速开始**: [QUICKSTART.md](../QUICKSTART.md)
- **示例应用**: [examples/README.md](../examples/README.md)
- **SDK 文档**: [docs/sdk_manual.md](../docs/sdk_manual.md)
- **CLI 文档**: [docs/cli_manual.md](../docs/cli_manual.md)

## 总结

`glow apply` 是配置 Glow 应用的**唯一方式**。通过声明式 YAML 文件，你可以：

✅ 配置应用的所有方面（执行、环境变量、端口、域名、应用配置等）
✅ 版本控制配置文件
✅ 幂等更新应用配置
✅ 自动配置 Ingress（域名绑定）
✅ 自动重启应用（配置变更时）

**配置即代码原则**：
- 所有配置在 `spec.config` 字段中声明
- 用户自行提供数据库连接等基础设施配置
- 系统将配置写入本地 JSON 文件供应用读取
- 不提供独立的资源绑定或配置生成机制

记住：**所有应用配置都在 YAML 的 `spec.config` 中完成，不需要独立的配置命令或资源绑定命令**。
