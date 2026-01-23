---
name: glow-deploy
description: Glow 应用部署与运维指南。包含 glow CLI 命令、部署工作流、配置管理、多环境部署、CI/CD 集成、日志查看和故障排查。当用户需要：部署 glow 应用、管理应用生命周期、查看日志、配置多环境、集成 CI/CD 时使用。
---

# Glow Deploy

Glow 应用的部署、运维和管理指南。

## 快速开始

### 编译与部署

```bash
# 编译应用
go build -o my-app

# 部署到 Glow Server
glow deploy ./my-app

# 查看应用状态
glow get apps

# 查看日志
glow logs my-app
```

## Glow CLI 命令

### 应用部署

```bash
# 部署应用
glow deploy <binary_path>

# 指定应用名称
glow deploy ./myapp --name custom-name

# 部署脚本（由 glow init 生成）
./scripts/deploy.sh              # 单应用自动部署
./scripts/deploy.sh <app_name>   # 指定应用部署
./scripts/deploy.sh              # 多应用交互式选择
```

### 资源查看

```bash
# 列出所有应用
glow get apps

# 列出所有网关路由
glow get ingress

# 查看节点信息
glow get nodes

# 查看基础设施资源
glow get resources

# 查看应用详情
glow describe app my-app

# 查看节点详情
glow describe node localhost
```

### 日志查看

```bash
# 查看应用日志
glow logs my-app

# 实时跟踪日志
glow logs my-app -f
```

### 生命周期管理

```bash
# 启动应用
glow start app my-app

# 停止应用
glow stop app my-app

# 重启应用
glow restart app my-app

# 删除应用
glow delete app my-app

# 删除路由
glow delete ingress my-app
```

### 配置管理

```bash
# 查看应用配置
glow config view my-app

# 编辑配置（支持热更新）
glow config edit my-app
```

### 环境管理

```bash
# 列出所有环境
glow context list

# 切换环境
glow context use prod

# 添加环境
glow context add prod --url http://prod-server:32102 --key <api-key>

# 查看认证信息
glow auth view

# 重置认证
glow auth reset
```

## 部署工作流

### 标准流程

```bash
# 1. 编译
go build -o my-app

# 2. 部署
glow deploy ./my-app

# 3. 验证
glow get apps
glow logs my-app

# 4. 配置（如需要）
glow config edit my-app
```

### 多应用部署

```bash
# 使用部署脚本
./scripts/deploy.sh

# 脚本会自动：
# - 扫描 cmd/ 目录检测应用
# - 交互式选择应用
# - 自动构建和部署
```

## 配置管理

### 配置优先级

1. Glow Server 配置中心（推荐）
2. 本地 `local_config.json`（降级）

### 动态配置

```bash
# 编辑配置
glow config edit my-app

# 修改 JSON 配置后保存
# 配置会自动热更新到应用
```

### 环境特定配置

```json
{
  "debug": false,
  "db_name": "prod_db",
  "redis": {
    "host": "prod-redis.example.com",
    "port": 6379
  }
}
```

## 多环境管理

### 环境配置

```bash
# 开发环境
glow context add dev --url http://dev-server:32102 --key <dev-key>

# 生产环境
glow context add prod --url http://prod-server:32102 --key <prod-key>

# 测试环境
glow context add test --url http://test-server:32102 --key <test-key>
```

### 环境切换

```bash
# 查看当前环境
glow context list

# 切换到生产环境
glow context use prod

# 部署到生产环境
glow deploy ./my-app
```

## CI/CD 集成

### GitHub Actions

```yaml
name: Deploy to Glow

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Build
        run: go build -o app

      - name: Deploy
        env:
          GLOW_SERVER_URL: ${{ secrets.GLOW_SERVER_URL }}
          GLOW_API_KEY: ${{ secrets.GLOW_API_KEY }}
        run: |
          curl -fsSL "https://raw.githubusercontent.com/luaxlou/glow/main/scripts/install-local.sh" | bash
          glow context add production --url $GLOW_SERVER_URL --key $GLOW_API_KEY
          glow context use production
          glow deploy ./app
```

### GitLab CI

```yaml
deploy:
  stage: deploy
  image: golang:1.21
  script:
    - go build -o app
    - curl -fsSL "https://raw.githubusercontent.com/luaxlou/glow/main/scripts/install-local.sh" | bash
    - glow context add production --url $GLOW_SERVER_URL --key $GLOW_API_KEY
    - glow context use production
    - glow deploy ./app
  only:
    - main
```

## 域名配置

### 通过 API 配置

```bash
curl -H "Authorization: Bearer <api-key>" \
  -H "Content-Type: application/json" \
  -d '{"app_name":"my-app","domain":"myapp.example.com","port":8080}' \
  http://localhost:32102/ingress/update
```

### Nginx 配置

确保主配置包含：

```nginx
include /var/lib/glow-server/nginx/*.conf;
```

## 日志管理

### 查看日志

```bash
# 实时日志
glow logs my-app

# 日志位置
# Linux: /var/lib/glow-server/apps/<app-name>/logs/<app-name>.log
# macOS: ~/Library/Application Support/glow-server/apps/<app-name>/logs/
```

### 日志轮转

- 单文件最大 10MB
- 保留最近 5 个文件

## 故障排查

### 应用无法启动

```bash
# 1. 查看 glow-server 日志
glow-server info

# 2. 查看应用日志
glow logs my-app

# 3. 检查应用状态
glow describe app my-app

# 4. 检查端口占用
lsof -i :<port>
```

### 配置未生效

```bash
# 1. 确认配置已更新
glow config view my-app

# 2. 重启应用（某些配置需要重启）
glow restart app my-app

# 3. 检查应用日志
glow logs my-app
```

### 数据库连接失败

```bash
# 1. 检查资源
glow get resources

# 2. 查看 glow-server info
glow-server info

# 3. 检查数据库
mysql -u root -p -e "SHOW DATABASES;"
```

## 性能优化

### 减小二进制大小

```bash
# 编译优化
go build -ldflags="-s -w" -o my-app

# 使用 upx 压缩（可选）
upx --best --lzma my-app
```

### 回滚策略

```bash
# 部署前备份
cp my-app my-app.backup

# 部署新版本
glow deploy ./my-app

# 如有问题，回滚
glow delete app my-app
glow deploy ./my-app.backup --name my-app
```

## 最佳实践

1. **使用部署脚本**: 优先使用 `./scripts/deploy.sh`
2. **配置管理**: 通过 `glow config` 管理环境配置
3. **多环境隔离**: 使用 context 隔离不同环境
4. **日志监控**: 定期查看 `glow logs`
5. **优雅重启**: 使用 `glow restart` 而非 `stop` + `start`
6. **版本管理**: 部署前保留备份，便于回滚

## HTTP API 参考

详见 [API Reference](references/http-api.md)

## 常见问题

### Q: 如何查看应用的完整信息？
A: 使用 `glow describe app <name>`

### Q: 配置修改后需要重启吗？
A: 大部分配置支持热更新，某些需要重启，查看应用日志确认

### Q: 如何同时部署多个应用？
A: 使用 `./scripts/deploy.sh` 选择"全部"或逐个部署
