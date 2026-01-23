# 测试发布 - 执行清单

## ✅ 准备工作已完成

- [x] 重新编译 glow-server（`glow-server-new`）
- [x] 创建应用配置文件（`app.yaml`）
- [x] 创建自动化脚本（`deploy.sh`, `test-deploy.sh`）
- [x] 创建完整文档（手册、指南、README）
- [x] 准备测试命令

## 📋 执行清单（需要手动操作）

### 步骤 1: 重启 glow-server ⚠️ 需要sudo

**在终端执行以下命令**:

```bash
cd /Users/john/workspace/luaxlou/glow

# 1. 停止旧服务器
sudo pkill -f glow-server

# 2. 替换二进制文件
sudo mv glow-server-new glow-server

# 3. 启动新服务器（选择一种方式）

# 方式 A: 后台运行（推荐）
nohup sudo ./glow-server serve > /tmp/glow-server.log 2>&1 &

# 方式 B: 前台运行（可看到实时日志）
sudo ./glow-server serve
```

**验证服务器启动**:
```bash
# 等待 3 秒
sleep 3

# 测试健康检查
curl http://localhost:32102/health
```

预期输出: `ok`

---

### 步骤 2: 测试新的 API 路由

```bash
# 测试 PUT /apps/:name 路由是否存在
curl -X PUT http://localhost:32102/apps/test \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test" \
  -d '{"name":"test"}'
```

**检查**: 不应该返回 `404 page not found`

---

### 步骤 3: 发布应用

#### 方式 A: 使用自动化脚本（推荐）

```bash
cd /Users/john/workspace/luaxlou/glow/examples/simple-app
./deploy.sh
```

#### 方式 B: 手动执行

```bash
cd /Users/john/workspace/luaxlou/glow/examples/simple-app

# 3.1 应用配置
./glow apply -f app.yaml
```

**预期输出**:
```
Applying App 'simple-app' from app.yaml...
✓ App 'simple-app' registered successfully
→ Configuring Ingress for domain: app.local
✓ Ingress configured: http://app.local -> port 33203
→ Binding MySQL resources...
✓ MySQL 'simple_app_db' bound successfully (DSN: ****)
→ Generating config file...
✓ Config file written to: /var/lib/glow-server/apps/simple-app/simple-app_local_config.json
```

```bash
# 3.2 启动应用
./glow start app simple-app
```

预期输出: `app.apps/simple-app started`

---

### 步骤 4: 验证应用运行

```bash
# 4.1 查看应用状态
./glow get app simple-app
```

**检查点**:
- [ ] Status: `RUNNING`
- [ ] PID: 非 0
- [ ] Port: `33203`
- [ ] Domain: `app.local`

```bash
# 4.2 查看应用日志
./glow logs simple-app
```

**检查点**:
- [ ] 看到 "App simple-app starting"
- [ ] 看到 "MySQL is ready and connected!"
- [ ] 看到 "Listening and serving HTTP on :33203"

```bash
# 4.3 测试 HTTP 访问
curl http://app.local/
```

**检查点**:
- [ ] 返回 "Hello from Glow app with Gin!"

---

### 步骤 5: 验证配置文件

```bash
# 查看生成的配置文件
cat /var/lib/glow-server/apps/simple-app/simple-app_local_config.json
```

**检查点**:
- [ ] 文件存在
- [ ] 包含 `mysql.dsn` 字段
- [ ] DSN 格式正确

---

### 步骤 6: 验证 Ingress

```bash
# 检查 Nginx 配置
cat /etc/nginx/sites-available/simple-app
```

**检查点**:
- [ ] 配置文件存在
- [ ] 包含正确的 proxy_pass 配置

---

### 步骤 7: 完整测试流程（一键）

如果前面的步骤都完成了，可以使用测试脚本：

```bash
cd /Users/john/workspace/luaxlou/glow/examples/simple-app
./test-deploy.sh
```

这个脚本会自动执行所有验证步骤。

---

## 🔍 故障排查

### 问题 1: 服务器返回 404

**症状**: `glow apply` 返回 `server returned status 404`

**原因**: glow-server 没有加载新的 API 路由

**解决**:
```bash
# 确认使用了新编译的 glow-server
cd /Users/john/workspace/luaxlou/glow
ls -lh glow-server*

# 完全重启服务器
sudo pkill -9 -f glow-server
sudo mv glow-server-new glow-server
sudo ./glow-server serve
```

### 问题 2: MySQL 连接失败

**症状**: 日志显示 "MySQL not available"

**解决**:
```bash
# 检查 MySQL 状态
sudo systemctl status mysql

# 启动 MySQL（如果未运行）
sudo systemctl start mysql

# 验证 MySQL 可以连接
mysql -u root -p -e "SELECT 1;"
```

### 问题 3: 域名无法访问

**症状**: `curl http://app.local/` 失败

**解决**:
```bash
# 添加 hosts 条目
echo "127.0.0.1 app.local" | sudo tee -a /etc/hosts

# 测试 DNS 解析
ping app.local

# 检查 Nginx 配置
sudo nginx -t
```

### 问题 4: 应用启动失败

**症状**: 状态不是 RUNNING

**解决**:
```bash
# 查看详细日志
./glow logs simple-app

# 手动测试应用
cd /var/lib/glow-server/apps/simple-app
./simple-app --help
```

---

## ✅ 成功标准

所有检查点都通过表示测试成功：

- [ ] glow-server 运行正常
- [ ] PUT /apps/:name 路由存在
- [ ] `glow apply` 执行成功
- [ ] MySQL 自动创建并绑定
- [ ] Ingress 自动配置
- [ ] 配置文件正确生成
- [ ] 应用状态为 RUNNING
- [ ] HTTP 访问成功

---

## 🎉 完成后

测试成功后，你可以：

1. **探索更多功能**:
   - 修改 `app.yaml` 并重新 apply
   - 添加更多资源（多数据库、Redis）
   - 创建多个应用

2. **查看文档**:
   - [Glow Apply 手册](../docs/glow_apply_manual.md)
   - [快速开始](../QUICKSTART.md)
   - [示例总览](../examples/README.md)

3. **清理环境**（可选）:
   ```bash
   ./glow stop simple-app
   ./glow delete app simple-app
   ```

---

## 📞 需要帮助？

- 查看日志: `tail -f /tmp/glow-server.log`
- 查看应用日志: `./glow logs simple-app`
- 查看部署指南: `cat DEPLOYMENT.md`

**祝测试顺利！** 🚀
