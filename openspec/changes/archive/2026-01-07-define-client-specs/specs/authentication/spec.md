## ADDED Requirements
### Requirement: 认证管理 (Auth Management)
系统 MUST 提供管理连接认证信息（Server URL & API Key）的命令，命令组命名为 `auth`。

#### Scenario: 查看认证信息
- **WHEN** 用户执行 `glow auth view`
- **THEN** CLI 应显示脱敏后的连接与认证信息

#### Scenario: 重置认证
- **WHEN** 用户执行 `glow auth reset`
- **THEN** CLI 应清除本地认证缓存
- **AND** 立即进入交互式引导流程

### Requirement: 交互式配置引导 (Interactive Config Bootstrap)
系统 MUST 在检测到配置缺失时自动触发交互式引导，而不是报错退出。

#### Scenario: 自动引导
- **WHEN** 用户执行任意命令（如 `glow get deploy`）且本地无配置文件
- **THEN** CLI 应提示 "Configuration not found, entering interactive setup..."
- **AND** CLI 应提示输入 Server URL（默认 `http://localhost:32102`）
- **AND** CLI 应提示输入 API Key
- **AND** CLI 应保存配置并自动继续执行原命令
