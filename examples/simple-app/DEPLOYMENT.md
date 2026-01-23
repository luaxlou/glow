# Simple App 发布指南

## 前提条件

1. **glow-server 正在运行**，且是最新编译版本
2. **MySQL 服务正在运行**
3. **CLI 已配置** 连接到 glow-server

## 发布步骤

### 1. 重新编译 glow-server（如果需要）

```bash
cd /path/to/glow
go build -o glow-server ./cmd/glow-server
```

### 2. 重启 glow-server

**重要**: 需要重启 glow-server 以加载新的 API 路由（`PUT /apps/:name`）

```bash
# 停止旧的 glow-server
sudo pkill glow-server

# 启动新的 glow-server
sudo ./glow-server serve
```

### 3. 确认 CLI 上下文

```bash
# 切换到本地上下文
./glow context use default

# 测试连接
./glow get apps
```

### 4. 发布应用

#### 方式一：使用自动脚本（推荐）

```bash
cd examples/simple-app
./deploy.sh
```

#### 方式二：手动执行

```bash
cd examples/simple-app

# Step 1: 应用配置（创建数据库、配置域名）
./glow apply -f app.yaml

# Step 2: 启动应用
./glow start app simple-app

# Step 3: 查看状态
./glow get app simple-app
```

### 5. 访问应用

```bash
# 测试 HTTP 访问
curl http://app.local/

# 或在浏览器打开
open http://app.local/
```

**注意**: 需要先在 `/etc/hosts` 添加域名解析：
```
127.0.0.1 app.local
```

## 验证部署

### 1. 检查应用状态

```bash
./glow get app simple-app
```

预期输出：
```
Name:       simple-app
Status:     RUNNING
PID:        <pid-number>
Port:       33203
Domain:     app.local
Restarts:   0
Command:    /var/lib/glow-server/apps/simple-app/glow_simple-app
```

### 2. 查看应用日志

```bash
./glow logs simple-app
```

预期输出包含：
```
App simple-app starting
MySQL is ready and connected!
[GIN-debug] Listening and serving HTTP on :33203
```

### 3. 检查生成的配置文件

```bash
cat /var/lib/glow-server/apps/simple-app/simple-app_local_config.json
```

预期输出：
```json
{
  "mysql": {
    "dsn": "glow_simple_app:<password>@tcp(localhost:3306)/glow_simple_app_db"
  }
}
```

## 常见问题

### Q1: glow apply 返回 404

**原因**: glow-server 版本过旧，缺少 `PUT /apps/:name` 路由

**解决**:
```bash
# 重新编译并重启 glow-server
go build -o glow-server ./cmd/glow-server
sudo pkill glow-server
sudo ./glow-server serve
```

### Q2: MySQL 连接失败

**原因**: MySQL 服务未运行或密码不正确

**解决**:
```bash
# 检查 MySQL 状态
sudo systemctl status mysql

# 如果需要密码，在 app.yaml 中添加
resources:
  mysql:
    - dbName: simple_app_db
      existingPassword: "your-mysql-root-password"
```

### Q3: 应用启动失败

**排查步骤**:
```bash
# 1. 查看详细日志
./glow logs simple-app

# 2. 检查进程状态
./glow get app simple-app

# 3. 手动测试配置文件
cat /var/lib/glow-server/apps/simple-app/simple-app_local_config.json
```

### Q4: 域名无法访问

**解决**:
```bash
# 添加 hosts 条目
echo "127.0.0.1 app.local" | sudo tee -a /etc/hosts

# 测试 DNS 解析
ping app.local

# 检查 Nginx 配置
cat /etc/nginx/sites-enabled/simple-app
```

## 清理

如果需要完全删除应用：

```bash
# 停止应用
./glow stop simple-app

# 删除应用（包括配置文件、二进制、Nginx配置）
./glow delete app simple-app

# 验证删除
./glow get apps
```

## 下一步

- 尝试修改 `app.yaml` 并重新 apply
- 添加更多资源（Redis、多数据库）
- 查看其他示例应用
- 阅读 [完整文档](../../docs/)

## 支持的应用操作

```bash
# 生命周期
./glow start app simple-app      # 启动
./glow stop simple-app           # 停止
./glow restart app simple-app    # 重启
./glow delete app simple-app     # 删除

# 查询
./glow get app simple-app        # 查看详情
./glow get apps                  # 列出所有
./glow logs simple-app           # 查看日志
./glow get ingress               # 查看域名绑定

# 配置
./glow apply -f app.yaml         # 应用/更新配置
./glow get config simple-app     # 查看配置
```
