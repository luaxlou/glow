# 项目背景 (Project Context)

## 目标 (Purpose)
Glow 是一个专为 Go 语言打造的全生命周期应用治理框架。旨在通过 Code-First（代码优先）理念，将基础设施定义回归代码，提供从开发、部署到运行的完整解决方案。它解决传统开发与运维割裂的问题，实现应用资源的自动申请、连接与进程级治理。

## 技术栈 (Tech Stack)
- **Language**: Go (Golang)
- **Web Framework**: Gin (HTTP API)
- **Storage**: SQLite (Config & State), MySQL/Redis (Managed Resources)
- **Ingress**: Nginx
- **Process Management**: Native `os/exec`, `syscall`
- **System Stats**: `gopsutil`

## 项目约定 (Project Conventions)

### 核心架构 (Architecture Patterns)
- **Server-Agent 模型**: `glow-server` 作为宿主机守护进程，负责所有治理任务。
- **Starter 机制**: 应用通过引入 `glow/starter` SDK 自动接入治理体系。
- **Lazy Provisioning**: 资源（如数据库）在应用启动请求时即时创建。

### 命名与路径 (Naming & Paths)
- **Managed Binaries**: 受管应用二进制重命名为 `glow_<app_name>`。
- **Env Variables**: 注入环境变量以 `OP_` 开头 (e.g., `OP_APP_PORT`, `OP_SERVER_URL`)。
- **Data Directory**: 默认数据存储在 `data_dir` (通常为运行目录或系统配置目录)，结构如下：
  ```
  apps/
    <app_name>/
      glow_<app_name>  # Binary
      logs/            # Rotated logs
  ```

### 错误处理 (Error Handling)
- 统一使用 `api.Response` 结构返回 JSON 格式结果。
- HTTP 状态码映射业务结果 (200 OK, 400 Bad Request, 500 Internal Error)。

## 领域上下文 (Domain Context)
- **App**: 一个独立运行的 Go 服务进程。
- **Provisioner**: 负责在基础设施（Host）上创建具体资源（如 MySQL Database）的组件。
- **Manifest**: 应用或宿主机的声明式配置文件（YAML/JSON）。
- **Ingress**: 负责将外部 HTTP 流量路由到具体应用端口的机制（当前实现为 Nginx）。

## 重要约束 (Important Constraints)
- **单机治理**: 当前主要针对单机多应用部署场景。
- **端口管理**: 应用必须监听由 Server 分配的 `OP_APP_PORT`。
- **权限**: Server 需要有权限操作 Nginx 配置文件及 Reload Nginx 进程，且需有权限切换用户（如 `su glow`）运行应用。