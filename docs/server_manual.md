# Glow Server 用户手册

`glow-server` 是 Glow 框架的核心组件，运行在宿主机上，负责应用的全生命周期管理、配置存储与基础设施对接。

## 1. 核心功能概览

*   **进程托管**: 替代 Systemd，管理应用进程的启动、停止、重启与日志轮转。
*   **配置中心**: 基于 SQLite 的 KV 存储，支持应用配置的热更新与持久化。
*   **网关自动化**: 基于应用域名自动生成 Nginx 反向代理配置。
*   **SDK 协同**: 通过 TCP (AppCenter) 与应用 SDK 保持长连接，实现服务注册与心跳保活。

## 2. 安装与初始化

### 一键安装（推荐）

安装与初始化统一通过脚本完成（不依赖 Go 工具链，不提供手动编译/拷贝安装方式）。

```bash
# Linux 服务器（安装 glow-server + glow，启用服务）
curl -fsSL "https://raw.githubusercontent.com/luaxlou/glow/main/scripts/install.sh" | sudo bash

# 本地安装（macOS/Linux，不常驻；需要时前台启动）
curl -fsSL "https://raw.githubusercontent.com/luaxlou/glow/main/scripts/install-local.sh" | bash

# 卸载（保留配置与数据库）
curl -fsSL "https://raw.githubusercontent.com/luaxlou/glow/main/scripts/uninstall.sh" | sudo bash
```

安装脚本会自动完成：
- 下载最新版本二进制文件
- sha256 校验
- 安装到 PATH
- 执行 `glow-server keygen`（生成/复用 API Key）
- 配置 `glow` CLI 默认 context
- Linux 场景下安装并启动系统服务（`install-local.sh` 不常驻、不注册服务）

### 重装与升级

重复执行安装脚本会自动检测并复用已有的配置和数据库：

- **已检测到现有安装**：脚本会提示"已检测到并复用既有配置/数据库"
- **不会覆盖数据**：现有的数据库文件（`glow.db`）和配置目录会被保留
- **重装前备份**：脚本会在覆盖二进制文件前自动备份到 `/tmp/glow-server-backup-<timestamp>`

如需完全重置（清除所有数据），请手动删除以下目录后重新执行安装：

```bash
# Linux 系统安装
sudo rm -rf /var/lib/glow-server/db/
sudo rm -rf /var/lib/glow-server/config/

# 本地开发安装
rm -rf ~/Library/Application Support/glow-server/db/  # macOS
rm -rf ~/.local/share/glow-server/db/                  # Linux
```

## 3. CLI 命令参考

### `keygen`
生成或查看 API Key。

```bash
glow-server keygen
```

### `serve`
启动核心服务守护进程。

```bash
# 本地开发：前台启动
glow-server serve
```

### `add`
添加系统资源（MySQL、Redis、Nginx）。

```bash
# 添加 MySQL
glow-server add mysql

# 添加 Redis
glow-server add redis

# 添加/发现 Nginx
glow-server add nginx
```

### `info`
显示当前服务器的配置状态、集成的资源信息及服务运行状态。

```bash
glow-server info
```

示例输出（节选）：

```text
Glow Server Information
-----------------------
PID: 12345

Managed Resources
-----------------
MySQL:
  Host: 127.0.0.1
  Port: 3306
  Root User: root
  Root Password: <PASSWORD>
  Databases:
    - app_db (charset=utf8mb4)
  Raw Config:
  {
    "host": "127.0.0.1",
    "port": 3306,
    "user": "root",
    "password": "<PASSWORD>",
    "databases": [
      {
        "name": "app_db",
        "charset": "utf8mb4"
      }
    ]
  }
Nginx: /usr/sbin/nginx (Version: 1.24.0)
Service: [STATUS CHECK NOT IMPLEMENTED]
```

## 3. HTTP API 参考

Base URL: `http://localhost:32102`

### API 认证

除 `/health` 端点外，所有 HTTP 管理 API 均需要 API Key 认证。客户端必须在请求头中包含有效的 API Key：

```http
Authorization: Bearer <api_key>
```

**认证错误码**：

| HTTP 状态码 | 说明 | 示例场景 |
|------------|------|---------|
| 200 | 认证成功 | API Key 匹配 |
| 401 | 未授权 | 缺少 `Authorization` 头或格式错误 |
| 403 | 禁止访问 | API Key 不匹配 |
| 500 | 服务器错误 | 服务端未配置 `api_key` |

**示例请求**：

```bash
# 使用 curl 访问受保护的 API
curl -H "Authorization: Bearer your-api-key" \
  http://localhost:32102/apps/list

# 使用 Glow CLI（CLI 会自动添加认证头）
glow app list
```

**获取 API Key**：

```bash
# 查看服务器生成的 API Key
glow-server keygen
```

**注意**：
- API Key 在安装时通过 `glow-server keygen` 自动生成并存储在 SQLite `system_config` 表中
- Glow CLI 会自动读取并使用该 API Key，无需手动配置
- API Key 应妥善保管，建议仅通过本地或受信任的网络访问管理端口

### 应用管理
| Method | Endpoint | Description | Payload |
|--------|----------|-------------|---------|
| POST | `/apps/upload` | 上传应用二进制 | Multipart Form: `file` |
| POST | `/apps/start` | 启动应用 | `{ "name": "app1", "command": "./bin", "args": [], "port": 8080 }` |
| POST | `/apps/stop` | 停止应用 | `{ "name": "app1" }` |
| POST | `/apps/restart`| 重启应用 | `{ "name": "app1" }` |
| POST | `/apps/delete` | 删除应用 | `{ "name": "app1" }` |
| GET | `/apps/list` | 获取应用列表 | - |
| GET | `/apps/logs` | 获取应用日志 | `?name=app1` |

### 配置管理
| Method | Endpoint | Description | Payload |
|--------|----------|-------------|---------|
| GET | `/config/:appName` | 获取应用配置 | - |
| PUT | `/config/:appName` | 更新应用配置 | `{ "key": "value" }` (JSON) |

### 网关与域名 (Ingress)
| Method | Endpoint | Description | Payload |
|--------|----------|-------------|---------|
| POST | `/ingress/update` | 更新/创建 Nginx 路由 | `{ "app_name": "app1", "domain": "app1.com", "port": 8080 }` |
| POST | `/ingress/delete` | 删除 Nginx 路由 | `{ "app_name": "app1" }` |
| GET | `/ingress/list` | 列出所有路由 | - |

## 4. 核心机制详解

### 4.1 进程运行环境
Glow Server 在启动应用时，会自动注入以下环境变量：

*   `OP_APP_NAME`: 应用名称 (e.g., `billing-service`)
*   `OP_APP_PORT`: 分配的 HTTP 监听端口 (e.g., `54321`)
*   `OP_SERVER_URL`: Glow Server 地址 (e.g., `127.0.0.1:32101`)

**文件结构**:
应用运行时文件存放于 `data-dir/apps/<app-name>/`:
*   `glow_<app-name>`: 重命名后的二进制文件
*   `logs/<app-name>.log`: 标准输出日志 (自动轮转)

### 4.2 Nginx 自动化
如果启动参数中包含 `domain`，Server 会在 `data-dir/nginx/` 下生成 `<app-name>.conf`：

```nginx
upstream myapp {
    server 127.0.0.1:54321;
}
server {
    listen 80;
    server_name myapp.local;
    location / {
        proxy_pass http://myapp;
        ...
    }
}
```
*需确保主 Nginx 配置包含 `include data-dir/nginx/*.conf;`*
