# Redefine Glow Server CLI

## Background
Currently, `glow-server` has disjointed commands (`add`, `keygen`) and lacks a unified initialization experience. Users have to manually configure services and run multiple commands to set up a production-ready environment.

## Goal
To provide a unified `install` command that guides the user through setting up `glow-server` as a robust service, including key generation, resource integration (MySQL, Redis, Nginx), and self-hosting via Nginx. Additionally, an `info` command will be added to inspect the server's state.

## Proposal
1.  **Unified `install` Command**:
    -   Interactive CLI wizard.
    -   Idempotent execution (safe to re-run).
    -   Steps:
        1.  **Key Generation**: Ensure server identity keys exist.
        2.  **Service Installation**: Install `glow-server` as a system service (Systemd/Launchd).
        3.  **Resource Selection**: Multi-select prompt for resources (Nginx, MySQL, Redis).
        4.  **Resource Configuration**: Interactive setup for selected resources (reuse existing `add` logic).
        5.  **Self-Hosting**: If Nginx is available, offer to reverse-proxy `glow-server` (e.g., `glow.example.com`).

2. **新增 `info` 命令与命令优化**：
    - 展示当前配置（密钥、集成资源、服务状态）。
    - **命令重命名**：将 `server` 命令重命名为 `serve`，使其更符合 Go 惯例。
    - **禁用默认命令**：禁用 Cobra 默认生成的 `completion` 命令。

3. **重构现有命令**：
    -   Ensure `add`, `keygen` can still be run independently.
    -   Ensure `install` reuses these underlying implementations.

## Design
-   **CLI Framework**: Continue using `cobra`.
-   **Interactivity**: Use `survey` or `promptui` (or standard input/output) for the wizard.
-   **Service Management**: Detect OS/Init system (Systemd vs Launchd) to write appropriate service files.
-   **Idempotency**: Check for existing configuration/service files before overwriting.

## Outcome
A user can simply run `glow-server install` on a fresh server and end up with a fully configured, running `glow-server` service with database and ingress integrations.
