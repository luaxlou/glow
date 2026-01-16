# Glow CLI 用户手册

`glow` 是与 Glow Server 交互的命令行工具，采用类似 Kubernetes (kubectl) 的声明式与命令式相结合的设计风格，旨在提供高效、直观的本地开发环境管理体验。

## 1. 安装与配置

### 安装
（假设 `glow` 二进制文件已编译并放置于系统 PATH 路径下）

### 初始化配置
首次使用时，需要配置连接到的 Glow Server 地址及认证密钥。

```bash
# 交互式引导配置
glow context add local --url http://localhost:32102 --key <YOUR_API_KEY>

# 切换上下文
glow context use local
```

配置信息默认存储在 `~/.glow.json`。

## 2. 核心命令概览

Glow CLI 的命令结构主要由“动词 + 资源”构成。

### 资源查看与描述 (Read-only)

用于查看系统状态和资源详情。

*   **`glow get <resource>`**: 列出资源列表。
    *   `glow get apps` (或 `app`): 列出所有应用及其状态（CPU、内存、端口等）。
    *   `glow get ingress`: 列出所有网关路由规则。
    *   `glow get nodes` (或 `node`): 查看集群节点及其系统负载。
    *   `glow get resources`: 查看受管的基础设施资源（如 MySQL, Redis）。

*   **`glow describe <resource> <name>`**: 查看特定资源的详细元数据。
    *   `glow describe app my-app`: 显示应用的完整配置、环境变量及实时统计。
    *   `glow describe node localhost`: 显示节点系统详情。

*   **`glow logs <name>`**: 查看应用日志。
    *   `glow logs my-app`: 打印应用的标准输出日志。

### 资源生命周期管理 (Lifecycle)

直接控制应用进程的运行状态。

*   **`glow start app <name>`**: 启动已停止的应用。
*   **`glow stop app <name>`**: 优雅停止运行中的应用。
*   **`glow restart app <name>`**: 重启应用进程。
*   **`glow delete <resource> <name>`**: 删除资源。
    *   `glow delete app my-app`: 停止并移除应用。
    *   `glow delete ingress my-app`: 删除对应的路由规则。

### 配置与声明式操作 (Configuration)

*   **`glow apply -f <filename>`**: 核心声明式命令。
    *   支持应用 YAML/JSON 格式的 Manifest 文件，一次性部署或更新主机 (Host)、应用 (App)、配置 (Config) 及路由 (Ingress)。
    
*   **`glow config`**: 应用配置中心管理。
    *   `glow config view <app>`: 查看应用的当前 JSON 配置。
    *   `glow config edit <app>`: 调用系统编辑器（如 vim）交互式修改配置，保存后即热更新。

### 环境与认证管理 (System)

*   **`glow context`**: 多环境管理。
    *   `glow context list`: 列出所有环境。
    *   `glow context use <name>`: 切换当前环境。
*   **`glow auth`**: 认证管理。
    *   `glow auth view`: 查看当前连接信息。
    *   `glow auth reset`: 重置本地认证信息。

## 3. 常用场景示例

### 场景一：部署一个新应用
1. 编写 `app.yaml` 描述文件。
2. 执行部署：
   ```bash
   glow apply -f app.yaml
   ```
3. 查看状态：
   ```bash
   glow get apps
   ```

### 场景二：暴露应用 (Ingress)
1. 编写 `ingress.yaml`:
   ```yaml
   kind: Ingress
   metadata:
     name: my-app-ingress
   spec:
     domain: myapp.local
     service: my-app
     port: 8080
   ```
2. 应用配置：
   ```bash
   glow apply -f ingress.yaml
   ```

### 场景三：调试应用配置
1. 查看当前配置：
   ```bash
   glow config view my-service
   ```
2. 在线修改配置（例如开启调试模式）：
   ```bash
   glow config edit my-service
   ```
   或者通过文件更新：
   ```bash
   glow apply -f config.yaml
   ```
3. 重启应用以确保某些静态配置生效（如果需要）：
   ```bash
   glow restart app my-service
   ```

### 场景四：查看应用日志
```bash
# 实时查看日志
glow logs my-service
```
