# Glow CLI 用户手册

`glow` 是与 Glow Server 交互的命令行工具，采用声明式配置设计，通过 `glow apply` 命令统一管理应用配置和资源绑定。

## 1. 安装与配置

### 安装

```bash
# 本地安装（安装 glow + glow-server）
curl -fsSL "https://raw.githubusercontent.com/luaxlou/glow/main/scripts/install-local.sh" | bash
```

### 初始化配置

安装脚本会自动执行 `glow-server keygen` 并配置默认 context。此时你可以直接使用 `glow get ...` 等命令。

如需连接到其他 Glow Server：

```bash
# 添加远端环境
glow context add prod --server-url http://<YOUR_SERVER>:32102 --api-key <YOUR_API_KEY>

# 切换上下文
glow context use prod
```

配置信息存储在 `~/.glow.json`。

## 2. 核心命令概览

### 应用配置 (glow apply) - **重要**

`glow apply` 是**唯一的资源配置方式**。所有资源（应用、端口、域名、MySQL、Redis）都在 YAML 文件中声明。

*   **`glow apply -f <app.yaml>`**: 应用配置文件。
    *   注册/更新应用元数据
    *   绑定资源（MySQL、Redis）
    *   配置 Ingress（域名绑定）
    *   生成应用配置文件

**详细文档**: 参见 [Glow Apply 手册](glow_apply_manual.md)

### 应用生命周期管理 (Lifecycle)

控制应用进程的启动、停止和删除。

*   **`glow start app <name>`**: 启动应用。
*   **`glow stop <name>`**: 停止应用。
*   **`glow restart app <name>`**: 重启应用。
*   **`glow delete app <name>`**: 删除应用及其配置。

### 资源查看 (Read-only)

查看系统状态和资源详情。

*   **`glow get <resource>`**: 列出资源。
    *   `glow get apps` (或 `app`): 列出所有应用及其状态。
    *   `glow get ingress`: 列出所有域名绑定规则。
    *   `glow get node`: 查看节点状态。

*   **`glow describe <resource> <name>`**: 查看详细信息。
    *   `glow describe app my-app`: 显示应用完整配置。

*   **`glow logs <name>`**: 查看应用日志。
    *   `glow logs my-app`: 打印应用标准输出日志。

### 配置管理 (Configuration)

Glow 采用**完全声明式的配置管理**。所有应用配置通过 `app.yaml` 声明，使用 `glow apply` 应用。

*   **配置声明**: 在 `app.yaml` 的 `spec.config` 字段中声明配置
*   **应用配置**: 执行 `glow apply -f app.yaml` 应用配置
*   **查看配置**: 读取生成的 `<data-dir>/apps/<app>/<app>_local_config.json`

**配置管理原则**:
- ✅ 配置即代码：所有配置在 YAML 文件中管理
- ✅ 版本控制：配置变更可通过 Git 追踪
- ✅ 不可变性：不支持运行时通过命令行修改配置
- ❌ 已移除：`glow config set/get/list/export/edit` 等命令

### 服务器管理 (Server Management)

查看和管理 glow-server 的运行状态。

*   **`glow server info`**: 显示服务器信息（人类可读格式）
    *   显示 PID、数据目录、日志目录、配置路径
    *   显示服务器版本和运行时长
*   **`glow server info --json`**: 以 JSON 格式输出服务器信息

### 环境管理 (Context)

多环境管理。

*   **`glow context list`**: 列出所有环境。
*   **`glow context use <name>`**: 切换当前环境。
*   **`glow context add <name> --server-url <url> --api-key <key>`**: 添加新环境。

## 3. 工作流说明

### 新架构工作流

**重要**: Glow 现在使用"配置即代码"模型。所有配置（包括数据库连接）都在 YAML 文件中声明，用户自行提供数据库连接等基础设施配置。

#### 配置阶段

```bash
# 1. 编写应用配置文件
cat > app.yaml <<EOF
apiVersion: v1
kind: App
metadata:
  name: my-app
spec:
  binary: ./my-app
  port: 8080
  domain: myapp.example.com
  config:
    mysql_dsn: "user:pass@tcp(localhost:3306)/myapp_db"
    log_level: "info"
EOF

# 2. 应用配置
glow apply -f app.yaml
```

**执行结果**:
- ✅ 应用元数据已注册
- ✅ Ingress（域名）已配置
- ✅ 配置文件已生成（包含用户提供的 mysql_dsn）

#### 运行阶段

```bash
# 3. 启动应用
glow start app my-app

# 4. 查看状态
glow get app my-app

# 5. 查看日志
glow logs my-app
```

### 架构对比

**旧工作流（已废弃）**:
```bash
./my-app  # ❌ 应用启动时连接服务器申请资源
```

**新工作流（当前）**:
```bash
glow apply -f app.yaml     # ✅ 声明配置（用户自行提供数据库连接等）
glow start app my-app      # ✅ 启动应用
```

## 4. 常用场景示例

### 场景一：部署新的 Web 应用

```bash
# 1. 编写配置文件
cat > my-app.yaml <<EOF
apiVersion: v1
kind: App
metadata:
  name: my-web-app
spec:
  binary: ./my-web-app
  port: 8080
  domain: myapp.example.com
  config:
    mysql_dsn: "user:pass@tcp(localhost:3306)/myapp_db"
    redis_addr: "localhost:6379"
    log_level: "info"
EOF

# 2. 应用配置
glow apply -f my-app.yaml

# 3. 启动应用
glow start app my-web-app

# 4. 验证
curl http://myapp.example.com/
```

### 场景二：更新应用配置

```bash
# 1. 修改配置文件
vim my-app.yaml

# 2. 重新应用配置
glow apply -f my-app.yaml

# 3. 重启应用使新配置生效
glow restart app my-app
```

### 场景三：查看应用状态和日志

```bash
# 查看所有应用
glow get apps

# 查看应用详情
glow describe app my-app

# 查看实时日志
glow logs my-app

# 查看生成的配置文件
cat /var/lib/glow-server/apps/my-app/my-app_local_config.json
```

### 场景四：多环境管理

```bash
# 添加生产环境
glow context add prod --server-url http://prod-server:32102 --api-key <prod-key>

# 添加开发环境
glow context add dev --server-url http://localhost:32102 --api-key <dev-key>

# 切换到生产环境
glow context use prod

# 在生产环境部署
glow apply -f app-prod.yaml
glow start app my-app

# 切换回开发环境
glow context use dev
```

### 场景五：后台 Worker（无端口）

```bash
# Worker 不需要对外开放端口
cat > worker.yaml <<EOF
apiVersion: v1
kind: App
metadata:
  name: my-worker
spec:
  binary: ./my-worker
  # 不指定 port
  config:
    mysql_dsn: "user:pass@tcp(localhost:3306)/worker_db"
    queue_name: "tasks"
EOF

# 应用并启动
glow apply -f worker.yaml
glow start app my-worker
```

### 场景六：微服务（共享数据库）

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
```

```yaml
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

```bash
# 部署两个服务
glow apply -f api-service.yaml
glow apply -f worker.yaml

# 启动服务
glow start app api-service
glow start app worker
```

## 5. 命令参考

### glow apply

**语法**: `glow apply -f <filename>`

**参数**:
- `-f, --file string`: YAML 配置文件路径（必需）

**示例**:
```bash
glow apply -f app.yaml
glow apply -f /path/to/config.yaml
```

**输出**:
```
Applying App 'my-app' from app.yaml...
✓ App 'my-app' registered successfully
→ Configuring Ingress for domain: myapp.example.com
✓ Ingress configured: http://myapp.example.com -> port 8080
→ Generating config file...
✓ Config file written to: /var/lib/glow-server/apps/my-app/my-app_local_config.json

Summary:
  App Name: my-app
  Port: 8080
  Domain: myapp.example.com

Next steps:
  1. Review the config file generated
  2. Start the app: glow start app my-app
  3. Check status: glow get app my-app
```

### glow get

**语法**: `glow get <resource>`

**资源类型**:
- `apps` (或 `app`, `applications`): 列出所有应用
- `ingress`: 列出所有域名绑定
- `node`: 查看节点状态

**示例**:
```bash
glow get apps           # 列出应用
glow get app            # 同上
glow get ingress        # 列出 Ingress
glow get node           # 查看节点
```

### glow describe

**语法**: `glow describe <resource> <name>`

**示例**:
```bash
glow describe app my-app     # 查看应用详情
```

**输出**:
```
Name:       my-app
Status:     RUNNING
PID:        12345
Port:       8080
Domain:     myapp.example.com
Restarts:   0
Command:    /var/lib/glow-server/apps/my-app/glow_my-app
Args:       []
WorkDir:    /var/lib/glow-server/apps/my-app
Age:        5m30s

Resources:
  CPU:      0.1%
  Memory:   8.5 MB

Config:
  {
    "mysql_dsn": "user:***@tcp(localhost:3306)/myapp_db",
    "redis_addr": "localhost:6379",
    "log_level": "info"
  }
```

### glow start

**语法**: `glow start app <name>`

**示例**:
```bash
glow start app my-app
```

**输出**: `app.apps/my-app started`

### glow stop

**语法**: `glow stop <name>`

**示例**:
```bash
glow stop my-app
```

**输出**: `app.apps/my-app stopped`

### glow restart

**语法**: `glow restart app <name>`

**示例**:
```bash
glow restart app my-app
```

**输出**: `app.apps/my-app restarted`

### glow delete

**语法**: `glow delete app <name>`

**行为**:
- 停止应用（如果正在运行）
- 删除应用配置
- 删除应用文件（二进制、日志等）
- 删除 Ingress 配置

**示例**:
```bash
glow delete app my-app
```

**输出**: `app.apps/my-app deleted`

### glow logs

**语法**: `glow logs <name>`

**示例**:
```bash
glow logs my-app
```

**输出**: 显示应用的标准输出和标准错误。

### glow context

**子命令**:
- `glow context list`: 列出所有环境
- `glow context use <name>`: 切换环境
- `glow context add <name> --server-url <url> --api-key <key>`: 添加环境

**示例**:
```bash
# 列出环境
glow context list

# 切换环境
glow context use prod

# 添加新环境
glow context add staging --server-url http://staging:32102 --api-key xyz
```

## 6. 配置文件格式

### 应用配置 (app.yaml)

```yaml
apiVersion: v1
kind: App
metadata:
  name: app-name
spec:
  binary: ./app                # 必需：应用二进制路径
  workingDir: /path/to/dir     # 可选：工作目录
  port: 8080                    # 可选：HTTP 端口
  domain: app.local            # 可选：域名（需要 port）
  args: ["--v"]                # 可选：启动参数
  env:                          # 可选：环境变量
    - name: ENV
      value: production
  config:                       # 应用配置（所有配置在这里声明）⭐
    mysql_dsn: "user:pass@tcp(localhost:3306)/myapp_db"
    redis_addr: "localhost:6379"
    log_level: "info"
```

**详细说明**: 参见 [Glow Apply 手册](glow_apply_manual.md)

## 7. 故障排查

### 问题 1: apply 返回 404

**症状**: `glow apply -f app.yaml` 返回 `server returned status 404`

**原因**: glow-server 版本过旧，缺少新的 API 路由

**解决**:
```bash
# 重新编译并重启 glow-server
cd /path/to/glow
go build -o glow-server ./cmd/glow-server
sudo pkill -f glow-server
sudo ./glow-server serve
```

### 问题 2: 数据库连接失败

**症状**: 应用启动失败，日志中显示数据库连接错误

**排查步骤**:
```bash
# 1. 检查 MySQL 服务
sudo systemctl status mysql

# 2. 检查配置文件中的 DSN
cat /var/lib/glow-server/apps/my-app/my-app_local_config.json | grep mysql_dsn

# 3. 测试数据库连接
mysql -u user -p -h localhost -e "SELECT 1;"
```

**解决**: 确保 `spec.config.mysql_dsn` 配置正确

### 问题 3: 应用启动失败

**排查步骤**:
```bash
# 1. 查看详细日志
glow logs my-app

# 2. 检查配置文件
cat /var/lib/glow-server/apps/my-app/my-app_local_config.json

# 3. 检查进程状态
glow get app my-app
```

### 问题 4: 域名无法访问

**排查步骤**:
```bash
# 1. 检查 /etc/hosts
cat /etc/hosts | grep app.local

# 2. 添加 hosts 条目（如果缺失）
echo "127.0.0.1 app.local" | sudo tee -a /etc/hosts

# 3. 检查 Nginx 配置
cat /etc/nginx/sites-available/my-app

# 4. 测试 Nginx 配置
sudo nginx -t
```

## 8. 高级用法

### 环境变量注入

```yaml
spec:
  env:
    - name: DATABASE_URL
      value: "mysql://user:pass@localhost/db"
    - name: API_KEY
      value: "sk-xxxxxx"
```

### 多数据库配置

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

### 条件配置（使用注释）

```yaml
spec:
  port: 8080
  # domain: myapp.local    # 取消注释以启用
  config:
    mysql_dsn: "user:pass@tcp(localhost:3306)/myapp_db"
    # redis_addr: "localhost:6379"  # 取消注释以启用 Redis
```

## 9. 最佳实践

### 1. 版本控制

将所有 `app.yaml` 纳入 Git：

```bash
git add app.yaml
git commit -m "Add app configuration"
git push
```

### 2. 环境分离

为不同环境创建不同的配置文件：

```
app.yaml              # 开发环境
app-staging.yaml      # 测试环境
app-production.yaml   # 生产环境
```

### 3. 配置验证

应用前验证 YAML 语法：

```bash
# 使用 yamllint
yamllint app.yaml

# 或使用 Python
python3 -c "import yaml; yaml.safe_load(open('app.yaml'))"
```

### 4. 渐进式更新

先更新配置，再重启应用：

```bash
# Step 1: 更新配置（不影响运行中的应用）
glow apply -f app.yaml

# Step 2: 检查生成的配置
cat /var/lib/glow-server/apps/my-app/my-app_local_config.json

# Step 3: 重启应用
glow restart app my-app
```

## 10. 相关资源

- **Glow Apply 手册**: [docs/glow_apply_manual.md](glow_apply_manual.md)
- **快速开始**: [QUICKSTART.md](../QUICKSTART.md)
- **SDK 文档**: [docs/sdk_manual.md](sdk_manual.md)
- **示例应用**: [examples/README.md](../examples/README.md)

## 总结

Glow CLI 提供了声明式的应用管理方式：

✅ **统一配置**: `glow apply` 是唯一资源配置方式
✅ **简单易用**: YAML 文件定义所有资源
✅ **幂等操作**: 可以重复执行 apply
✅ **完整工具链**: 从配置到部署的完整支持

记住：**所有资源配置都在 YAML 文件中完成**。
