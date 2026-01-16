# Glow Server 用户手册

`glow-server` 是 Glow 框架的核心组件，运行在宿主机上，负责应用的全生命周期管理、配置存储与基础设施对接。

## 1. 核心功能概览

*   **进程托管**: 替代 Systemd，管理应用进程的启动、停止、重启与日志轮转。
*   **配置中心**: 基于 SQLite 的 KV 存储，支持应用配置的热更新与持久化。
*   **网关自动化**: 基于应用域名自动生成 Nginx 反向代理配置。
*   **SDK 协同**: 通过 TCP (AppCenter) 与应用 SDK 保持长连接，实现服务注册与心跳保活。

## 2. CLI 命令参考

### `install`
交互式安装向导，引导完成密钥生成、资源集成、系统服务注册与反向代理配置。

```bash
./glow-server install
```

### `serve`
启动核心服务守护进程。同时开启 HTTP API Server (默认端口 32102) 与 TCP AppCenter (默认端口 32101)。

```bash
# 启动服务 (需先运行 install 或 keygen)
./glow-server serve [--port 32102] [--app-center-port 32101] [--dir .]
```

### `info`
显示当前服务器的配置状态、集成的资源信息及服务运行状态。

```bash
./glow-server info
```

## 3. HTTP API 参考

Base URL: `http://localhost:32102`

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
