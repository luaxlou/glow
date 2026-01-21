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

