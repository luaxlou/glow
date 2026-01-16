## ADDED Requirements
### Requirement: 应用部署 (App Deployment)
系统 MUST 提供 `glow deploy` 命令用于部署或更新应用，并具备智能差异检测功能。

#### Scenario: 部署新应用
- **WHEN** 用户执行 `glow deploy <binary> --name myapp`，且应用不存在
- **THEN** CLI 应计算本地二进制 Hash
- **AND** CLI 应上传二进制文件到服务器
- **AND** CLI 应请求启动应用
- **AND** 服务器应保存二进制并记录其 Hash

#### Scenario: 更新应用 (二进制变更)
- **WHEN** 用户执行 `glow deploy <binary> --name myapp`，且应用已存在但二进制内容不同
- **THEN** CLI 检测到 Hash 不一致
- **AND** CLI 上传新二进制
- **AND** CLI 请求重启应用
- **AND** 服务器更新二进制和 Hash

#### Scenario: 跳过更新 (二进制未变)
- **WHEN** 用户执行 `glow deploy <binary> --name myapp`，且应用已存在且 Hash 一致
- **THEN** CLI 检测到 Hash 一致
- **AND** CLI 输出提示信息告知用户无需更新
- **AND** CLI 不执行上传和重启操作

### Requirement: 二进制上传 (Binary Upload)
系统 MUST 提供 API 接口支持二进制文件上传。

#### Scenario: 上传文件
- **WHEN** 客户端向 `/apps/upload` 发送 Multipart 请求
- **THEN** 服务器应接收并保存文件到临时或指定位置
- **AND** 返回文件保存路径或 ID