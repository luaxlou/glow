# 移除 GlowApp Starter 并使用本地配置

## Change ID
`remove-glowapp-use-local-config`

## 背景 (Background)
目前，Glow 应用使用 `starter/glowapp` 包来初始化和管理其生命周期。该包包含加载配置、可能连接到 `glow-server`（已弃用）以及处理关闭信号的逻辑。
目标是通过移除沉重的 `glowapp` starter 并依赖简单、明确的配置契约来简化这种交互。

## 目标 (Purpose)
- **简化应用集成**：应用不需要“starter”框架，只需要一种读取配置的方法。
- **应用与 Glow 内部解耦**：应用只需读取标准的 JSON 文件。
- **标准化配置**：配置始终位于工作目录（或 `data_dir` 上下文）下的 `config.json` 中。

## 需求 (Requirements)
### Requirement: 本地配置文件 (Local Config File)
- **GIVEN** 一个 Glow 管理的应用
- **WHEN** 应用启动时
- **THEN** 它必须能够从其工作目录中的 `config.json` 文件读取其配置。
- **AND** 它可以使用 `pkg/glowconfig` 库来简化此读取过程。

### Requirement: 渲染配置为 config.json (Render Config as config.json)
- **GIVEN** 执行了 `glow apply -f app.yaml`
- **WHEN** 服务器为应用渲染配置时
- **THEN** 它必须将配置写入 `config.json`（而不是 `<appName>_local_config.json`）。
- **AND** 文件必须放置在应用的工作目录 `<data_dir>/apps/<appName>/` 中。

### Requirement: 移除 GlowApp Starter (Remove GlowApp Starter)
- **GIVEN** 代码库
- **WHEN** 重构完成时
- **THEN** 必须移除 `starter/glowapp` 包。
- **AND** `starter/glowapp/config` 包必须被 `pkg/glowconfig` 替换/迁移。

## 设计 (Design)
- **SDK**: 一个新的（或移动的）包 `pkg/glowconfig` 将提供 `Load() (*Config, error)` 或类似的函数。
- **Server**: `internal/apiserver` 的 `handleRenderConfig` 将被更新为使用 `config.json`。
- **Migration**: 现有的应用（如果有）需要更新其代码以移除 `glowapp.Init()` 并使用 `glowconfig.Load()`。