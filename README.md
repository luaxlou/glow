# Glow: Go 语言全生命周期应用治理框架

**Glow** 是一个专为 Go 语言打造的**全栈式应用框架**。它不仅仅是一个 Web 框架或运维工具，而是提供了一套从**开发、定义、部署到运行**的完整生命周期解决方案。

> **核心愿景**: 让 Go 应用的生命周期管理像编写代码一样简单、确定。

## 什么是 Glow？

传统的开发模式中，应用开发（Dev）与应用运维（Ops）往往是割裂的：开发者关注业务逻辑，运维关注端口、配置与进程保活。

**Glow** 将这两者通过**Code-First（代码优先）**理念深度融合。它将基础设施的定义权回归给代码，让开发者在编写业务逻辑的同时，自动完成资源的申请与连接。

## 核心能力：全生命周期管理

Glow 覆盖了应用从诞生到运行的两个关键阶段：

### 1. 开发态 (Development) - 标准化 Starter
Glow 提供了一套类似 Spring Boot 的 Starter 机制（`glow/starter`），确立了 Go 应用的标准骨架：
*   **配置中心集成**: Starter 启动时自动连接 `glow-server` 获取配置，支持动态热更新，彻底告别本地配置文件。
*   **内建可观测性**: 默认集成健康检查、版本元数据与基础监控指标。

### 2. 运行态 (Runtime) - 进程级治理
在运行时，Glow 持续守护应用状态：
*   **去容器化轻量运行**: 直接管理宿主机进程，零虚拟化损耗，适合高性能场景。
*   **自愈与保活**: 配合监控组件，确保应用进程的稳定性与高可用。
*   **动态配置**: 支持运行时配置变更，组件自动重载。

## 架构详解 (Architecture)

```text
+---------------------------------------------------------------------------------------+
|  HOST MACHINE                                                                         |
|                                                                                       |
|   +-------------------------------------------------------------------------------+   |
|   |  GLOW SERVER (Daemon)                                                         |   |
|   |                                                                               |   |
|   |   [ Process Mgr ] --------+               [ HTTP API ] <---> [ Config DB ]    |   |
|   |          |                |                    ^   |                              |
|   |          | (5) Monitor    | (6) Config         |   +---> [ Provisioner ]      |   |
|   |          v                v                    |               |              |   |
|   +----------|----------------|--------------------|---------------|--------------+   |
|              |                |                    |               | (2) Create       |
|              |                |              (1)   |               |     User/DB      |
|              |                |            Register|               |                  |
|   +----------|-------+   +----|----------+ & Req   |      +--------v-----------+      |
|   | USER APPLICATION |   | NGINX GATEWAY | --------+      | LOCAL INFRA        |      |
|   |                  |   |               |                |                    |      |
|   |  [ Biz Logic ]   |   |  [ Vhost ] <------- Traffic    |  [ MySQL / Redis ] |      |
|   |        |         |   +---------------+                |                    |      |
|   |        v         |                                    |                    |      |
|   |   [ Glow SDK ] <-----------------------------------------+                 |      |
|   |                  |        (4) Connect                 |                    |      |
|   +------------------+                                    +--------------------+      |
|                                                                                       |
+---------------------------------------------------------------------------------------+
```

### Glow Server
`glow-server` 是运行在宿主机上的核心守护进程，充当了**应用运行时 (Runtime)**、**配置中心 (Config Center)** 与 **基础设施管理器 (Infra Manager)** 的角色。

它具备以下核心能力：

1.  **进程托管 (Process Management)**
    *   替代 Systemd/Supervisord，直接接管应用进程。
    *   提供**自动重启**（Crash Loop Backoff）、**日志轮转**（Log Rotation）与**资源监控**（CPU/Mem/IO）。
    *   自动为应用分配空闲端口，并注入到环境变量 `OP_APP_PORT` 中。

2.  **自动化网关 (Ingress Automation)**
    *   根据应用声明的域名（Domain），自动生成并刷新 Nginx 配置文件（`upstream` & `server` block）。
    *   实现服务发现与负载均衡的自动化闭环，无需人工干预 Nginx。

3.  **嵌入式配置中心**
    *   基于 SQLite 的轻量级配置存储，无外部依赖。
    *   提供 Restful API 供 Starter 拉取配置，支持配置的版本化与持久化。

### 交互流程
1.  **启动**: 开发者执行 `go run main.go` 或 `glow start`。
2.  **注册**: Starter 向本地 `glow-server` 注册身份。
3.  **运行**: 应用启动 HTTP Server，Server 接管其 PID 并配置 Nginx 路由。

## 快速开始 (Quick Start)

### 1. 安装 Glow（只提供一键脚本）

安装与初始化统一通过脚本完成（不依赖 Go 工具链，不提供手动编译/拷贝安装方式）。

```bash
# Linux 服务器（安装 glow-server + glow，启用服务）
curl -fsSL "https://raw.githubusercontent.com/luaxlou/glow/main/scripts/install.sh" | sudo bash

# 本地安装（macOS/Linux，本地开发/仅客户端机器也可用；不常驻、不注册服务）
curl -fsSL "https://raw.githubusercontent.com/luaxlou/glow/main/scripts/install-local.sh" | bash

# 卸载（保留配置与数据库）
curl -fsSL "https://raw.githubusercontent.com/luaxlou/glow/main/scripts/uninstall.sh" | sudo bash
```

脚本会自动完成：
- 下载 release 预编译二进制（`glow-server` + `glow`）
- sha256 校验
- 安装到 PATH
- 执行 `glow-server keygen`（生成/复用 API Key）
- 写入 `glow` 默认 context（安装后可直接使用 `glow` 命令）
- Linux 场景下安装/启用服务（`install-local.sh` 不常驻、不注册服务）

### 2. 启动 Server

```bash
# 本地开发（macOS / 或任意不常驻场景）：前台启动
# HTTP API :32102, App Center :32101
glow-server serve
```

### 3. 运行示例应用
保持 Server 运行，打开一个新的终端窗口运行示例应用。

```bash
# 进入示例目录
cd examples/simple-app

# 运行应用
# SDK 会自动连接本地 Glow Server，注册应用并启动 HTTP 服务
go run main.go
```

**预期效果**:
1.  终端显示 `App simple-app started`。
2.  Server 端日志显示接收到连接与注册请求。
3.  访问 `http://localhost:8080` 可看到 "Hello from Implicit-wiring app with Gin!"。

## 📚 文档 (Documentation)

*   [Glow Server 用户手册](docs/server_manual.md): 详细介绍了 CLI 命令、API 接口与核心运行机制。
*   [Glow CLI 用户手册](docs/cli_manual.md): 命令行工具 `glow` 的详细使用说明。
*   [Glow SDK 用户手册](docs/sdk_manual.md): 指导开发者如何使用 Glow SDK (Starter) 编写应用。

## 为什么叫 Glow？

**Glow** 寓意着应用在黑暗的服务器丛林中发出清晰、可控的光芒。它让每一个 Go 应用不再是黑盒中的孤岛，而是处于一个可观测、可管理、可预测的生命周期闭环之中。
