# Glow Apply 用户手册

## 概述

`glow apply` 是 Glow 框架中**唯一的资源配置方式**。通过声明式 YAML 文件，你可以一次性配置应用的所有资源（端口、域名、MySQL、Redis 等）。

## 核心概念

### 声明式配置

所有资源配置都在 YAML 文件中声明，而不是通过命令行参数或交互式配置：

```yaml
apiVersion: v1
kind: App
metadata:
  name: my-app
spec:
  port: 8080          # 端口
  domain: app.local   # 域名
  resources:         # 资源
    mysql:
      - dbName: mydb
```

### 幂等性

可以重复执行 `glow apply`，每次执行都会：
- 更新应用配置
- 重新绑定资源（如果需要）
- 重新生成配置文件

### 原子性

所有资源在一次 `apply` 中配置，要么全部成功，要么全部失败。

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
  # 执行配置
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

  # 资源声明
  resources:
    mysql:
      - dbName: myapp_db        # MySQL 数据库
    redis:
      - db: 0                   # Redis 实例
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

### spec.binary（必需）

应用二进制文件的路径。

**相对路径**: 相对于工作目录
**绝对路径**: 从根目录开始

**示例**:
```yaml
spec:
  binary: ./my-app           # 相对路径
  binary: /usr/bin/my-app    # 绝对路径
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

### spec.resources（可选）

资源声明（MySQL、Redis 等）。

#### MySQL 资源

```yaml
resources:
  mysql:
    - dbName: myapp_db           # 必需：数据库名
      existingPassword: "..."     # 可选：现有数据库密码
```

**行为**:
- 如果数据库不存在，自动创建并生成随机密码
- 如果数据库存在且 `existingPassword` 为空，返回错误
- 如果数据库存在且提供 `existingPassword`，使用现有凭据

**多数据库**:
```yaml
resources:
  mysql:
    - dbName: main_db
    - dbName: cache_db
    - dbName: log_db
```

#### Redis 资源

```yaml
resources:
  redis:
    - db: 0                    # Redis 数据库编号（默认 0）
      password: "..."           # 可选：密码（不指定则不使用）
```

**行为**:
- 连接到本地 Redis（localhost:6379）
- 使用指定的数据库编号

## 执行流程

`glow apply` 按以下顺序执行：

### 1. 验证 YAML 文件

检查必需字段和格式。

### 2. 注册/更新应用

调用 `PUT /apps/:name` API，保存应用元数据到数据库。

### 3. 配置 Ingress（如果指定了 domain）

生成 Nginx 配置文件并重载 Nginx。

**配置文件位置**: `/etc/nginx/sites-available/<app-name>`

### 4. 绑定资源

依次处理 `resources` 中声明的所有资源：

- **MySQL**: 创建数据库（如果不存在），生成或验证凭据
- **Redis**: 验证 Redis 连接

### 5. 生成配置文件

将所有资源配置写入本地配置文件。

**配置文件位置**: `/var/lib/glow-server/apps/<app-name>/<app-name>_local_config.json`

**示例内容**:
```json
{
  "mysql": {
    "dsn": "glow_myapp:random@tcp(localhost:3306)/glow_myapp_db"
  },
  "redis": {
    "addr": "localhost:6379",
    "username": "",
    "password": "",
    "db": 0
  }
}
```

### 6. 输出摘要

显示操作结果和下一步提示。

## 输出示例

### 成功输出

```
Applying App 'my-app' from app.yaml...
✓ App 'my-app' registered successfully
→ Configuring Ingress for domain: myapp.example.com
✓ Ingress configured: http://myapp.example.com -> port 8080
→ Binding MySQL resources...
✓ MySQL 'myapp_db' bound successfully (DSN: ****)
→ Binding Redis resources...
✓ Redis (db 0) bound successfully (Addr: localhost:6379)
→ Generating config file...
✓ Config file written to: /var/lib/glow-server/apps/my-app/my-app_local_config.json (245 bytes)

Summary:
  App Name: my-app
  Port: 8080
  Domain: myapp.example.com
  MySQL: 1 database(s)
  Redis: 1 instance(s)

Next steps:
  1. Review the config file generated
  2. Start the app: glow start app my-app
  3. Check status: glow get app my-app
```

### 需要凭据的输出

```
→ Binding MySQL resources...
→ MySQL 'existing_db' requires credentials
Enter existing MySQL password: *****
✓ MySQL 'existing_db' bound successfully (DSN: ****)
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

#### 4. MySQL 需要凭据

```
⚠ Warning: MySQL 'mydb' binding failed: needs_credentials
```

**解决**: 在 YAML 中添加 `existingPassword` 或在 CLI 交互时输入密码。

#### 5. API 路由不存在 (404)

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
  resources:
    mysql:
      - dbName: webapp_db
    redis:
      - db: 0
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
  resources:
    mysql:
      - dbName: worker_db
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
  resources:
    mysql:
      - dbName: shared_db
    redis:
      - db: 0

# worker.yaml
apiVersion: v1
kind: App
metadata:
  name: worker
spec:
  binary: ./worker
  resources:
    mysql:
      - dbName: shared_db    # 共享同一数据库
    redis:
      - db: 1                # 不同的 Redis DB
```

### 场景 4: 更新应用配置

```bash
# 1. 编辑 YAML
vim app.yaml

# 2. 应用新配置
glow apply -f app.yaml

# 3. 重启应用使新配置生效
glow restart app my-app
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

先更新资源，再重启应用：

```bash
# Step 1: 更新配置（不影响运行中的应用）
glow apply -f app.yaml

# Step 2: 检查生成的配置
cat /var/lib/glow-server/apps/my-app/my-app_local_config.json

# Step 3: 重启应用
glow restart app my-app
```

### 5. 配置审查

应用前审查摘要输出，确认：

- [ ] 端口正确
- [ ] 域名正确
- [ ] 资源绑定成功
- [ ] 配置文件已生成

## 高级用法

### 1. 多个 MySQL 数据库

```yaml
resources:
  mysql:
    - dbName: main_db
    - dbName: cache_db
    - dbName: log_db
```

应用中访问：
```go
glowmysql.Init("main_db")   // 使用 main_db
glowmysql.Init("cache_db")  // 使用 cache_db
```

### 2. 环境变量注入

```yaml
spec:
  env:
    - name: MYSQL_DSN
      value: "mysql://user:pass@localhost/dbname"
```

### 3. 条件配置（使用注释）

```yaml
spec:
  port: 8080
  # domain: myapp.local    # 取消注释以启用域名
  resources:
    mysql:
      - dbName: myapp_db
  # redis:                  # 取消注释以启用 Redis
  #   - db: 0
```

### 4. 组合配置（包含 args 和 env）

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
  resources:
    mysql:
      - dbName: myapp_db
```

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
cat /var/lib/glow-server/apps/my-app/my-app_local_config.json

# 手动测试应用
cd /var/lib/glow-server/apps/my-app
./my-app --help
```

### 问题 2: MySQL 绑定失败

**排查步骤**:
```bash
# 检查 MySQL 服务
sudo systemctl status mysql

# 测试 MySQL 连接
mysql -u root -p -e "SHOW DATABASES;"

# 检查生成的 DSN
cat /var/lib/glow-server/apps/my-app/my-app_local_config.json | grep dsn
```

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

✅ 一次性配置所有资源
✅ 版本控制配置文件
✅ 幂等更新应用配置
✅ 自动绑定 MySQL/Redis
✅ 自动配置 Ingress

记住：**所有资源配置都在 YAML 中完成，不需要独立的 `glow app add mysql` 等命令**。
