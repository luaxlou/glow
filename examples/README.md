# Glow Examples

本目录包含使用 Glow 框架的应用示例。

## 架构说明

**重要变更**: Glow 现在使用声明式配置模型。应用不再在运行时连接 glow-server，而是通过 `glow apply` 命令预先配置所有资源。

### 新架构特点

1. **声明式配置**: 所有资源配置（MySQL、Redis、Ingress）在 YAML 文件中声明
2. **本地运行**: 应用启动时只读取本地配置文件，无需连接服务器
3. **统一接口**: `glow apply` 是唯一的资源配置方式
4. **幂等操作**: 可以重复执行 apply 更新配置

## 示例应用

### simple-app

一个简单的 Web 应用，演示了：
- HTTP 服务（Gin）
- MySQL 数据库集成
- WebSocket 支持
- 域名绑定（Ingress）

#### 快速开始

```bash
# 1. 进入示例目录
cd examples/simple-app

# 2. 应用配置（创建数据库、绑定域名等）
glow apply -f app.yaml

# 3. 启动应用
glow start app simple-app

# 4. 访问应用
curl http://app.local/

# 5. 查看应用状态
glow get app simple-app

# 6. 查看日志
glow logs simple-app
```

#### 配置文件说明

`app.yaml` 包含以下配置：

```yaml
spec:
  binary: ./simple-app        # 应用二进制路径
  port: 33203                 # HTTP 服务端口
  domain: app.local           # 域名绑定（Ingress）
  resources:
    mysql:
      - dbName: simple_app_db # MySQL 数据库名
```

#### 应用代码说明

```go
func main() {
    // 初始化应用身份
    glowapp.Init("simple-app")

    // 初始化 HTTP 服务（端口 33203）
    glowhttp.Init(33203)

    // 声明需要使用的数据库
    glowmysql.Init("simple_app_db")

    // 使用 MySQL（从本地配置读取 DSN）
    db, err := glowmysql.Gorm()

    // 启动 HTTP 服务
    glowhttp.Run()

    // 等待关闭信号
    glowapp.WaitForShutdown()
}
```

**关键点**:
- 应用不再调用 `glowmysql.ProvisionResource()`
- `glowmysql.Gorm()` 直接从本地配置文件读取 `mysql.dsn`
- 配置文件由 `glow apply` 命令生成

#### 本地开发

如果不使用 glow-server 运行，需要手动创建配置文件：

```bash
# 创建本地配置文件
cat > simple-app_local_config.json <<EOF
{
  "mysql": {
    "dsn": "root:password@tcp(localhost:3306)/simple_app_db"
  }
}
EOF

# 直接运行应用
./simple-app
```

## 工作流对比

### 旧工作流（已废弃）

```bash
# ❌ 旧方式：应用运行时连接服务器申请资源
./simple-app  # 应用启动时连接 glow-server
```

### 新工作流（当前）

```bash
# ✅ 新方式：预先配置，应用只读本地配置
glow apply -f app.yaml     # 配置资源
glow start app simple-app  # 启动应用
```

## 常见问题

### Q: 如何更新资源配置？

```bash
# 编辑 app.yaml
vim app.yaml

# 重新应用配置
glow apply -f app.yaml

# 重启应用使新配置生效
glow restart app simple-app
```

### Q: 如何绑定多个数据库？

```yaml
resources:
  mysql:
    - dbName: myapp_db
    - dbName: myapp_cache_db
```

### Q: 不使用 Ingress（域名）？

删除 `spec.domain` 字段即可：

```yaml
spec:
  port: 8080
  # domain: xxx  # 删除此行
```

### Q: 如何查看生成的配置文件？

```bash
cat /var/lib/glow-server/apps/simple-app/simple-app_local_config.json
```

## 更多资源

- [Glow SDK 文档](../../docs/sdk_manual.md)
- [Glow CLI 文档](../../docs/cli_manual.md)
- [glow apply 命令](../../README.md)
