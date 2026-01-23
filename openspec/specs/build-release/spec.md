# build-release Specification

## Purpose
定义 glow 项目的构建和发布规范，包括二进制输出位置、产物命名和校验文件要求。
## Requirements
### Requirement: Binary Output Location
所有编译的二进制文件 MUST 输出到项目根目录的 `./bin` 目录。

#### Scenario: 构建组件
- **WHEN** 开发者或代理编译组件（如 `glow-server`）
- **THEN** 生成的二进制文件 MUST 放置在 `./bin/` 中
- **AND** 项目根目录 MUST 保持无二进制产物

### Requirement: Release 产物命名 (Release Artifact Naming)
所有发布的二进制文件 MUST 遵循统一的命名规范。

#### Scenario: 命名规范
- **WHEN** 构建 release 产物
- **THEN** 二进制文件 MUST 使用格式 `{binary_name}-{os}-{arch}`
- **AND** 支持的平台组合包括：
  - `glow-server-linux-amd64`
  - `glow-server-linux-arm64`
  - `glow-server-darwin-amd64`
  - `glow-server-darwin-arm64`
  - `glow-linux-amd64`
  - `glow-linux-arm64`
  - `glow-darwin-amd64`
  - `glow-darwin-arm64`

### Requirement: 校验文件 (Checksum Files)
每个发布的二进制文件 MUST 配备 SHA256 校验文件。

#### Scenario: 生成校验文件
- **WHEN** 构建 release 产物
- **THEN** MUST 为每个二进制文件生成对应的 `.sha256` 文件
- **AND** 校验文件 MUST 使用与二进制文件相同的命名（添加 `.sha256` 后缀）
- **AND** 校验文件内容 MUST 包含 SHA256 哈希值和文件名
- **AND** 格式为：`<sha256_hash>  <filename>`

#### Scenario: 校验文件示例
- **GIVEN** 二进制文件 `glow-server-linux-amd64`
- **THEN** 校验文件名 MUST 为 `glow-server-linux-amd64.sha256`
- **AND** 文件内容示例：`a1b2c3d4e5f6...  glow-server-linux-amd64`

### Requirement: Release 产物分发 (Release Artifacts)
系统 MUST 提供可被安装脚本与自我更新逻辑消费的标准化 release 产物，覆盖 `glow-server`（服务端）与 `glow`（客户端）。

#### Scenario: 发布多平台归档
- **WHEN** 项目发布一个版本
- **THEN** release MUST 包含针对 `linux/darwin` 与 `amd64/arm64` 的二进制归档文件
- **AND** 归档内 MUST 包含 `glow-server` 与 `glow` 可执行文件（或提供等价的分别归档，确保安装脚本可同时安装两者）

#### Scenario: 安装脚本无需源码与 Go 环境
- **WHEN** 用户使用一键安装脚本安装 `glow-server` 与 `glow`
- **THEN** 安装脚本 MUST 仅依赖 release 产物（归档 + 校验文件）完成安装
- **AND** MUST NOT 依赖源码检出或 `go install`

#### Scenario: 提供校验文件
- **WHEN** release 发布归档文件
- **THEN** release MUST 同步发布 sha256 校验文件（例如 `checksums.txt`）
- **AND** 校验文件 MUST 可用于验证每个归档文件的完整性

### Requirement: 版本信息注入 (Version Metadata)
`glow-server` 与 `glow` MUST 提供可查询的版本元信息，以支持安装脚本与自我更新的决策与可观测性。

#### Scenario: 查询当前版本
- **WHEN** 用户运行 `glow-server version` 或 `glow version`
- **THEN** 命令 MUST 输出 `version`、`commit` 与 `buildDate`（若不可用则明确标注 unknown）

