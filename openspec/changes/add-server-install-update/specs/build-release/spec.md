## ADDED Requirements

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

