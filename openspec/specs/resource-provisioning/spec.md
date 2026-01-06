# 资源拨备 (Resource Provisioning)

## Purpose
实现应用依赖资源（如 MySQL 数据库、Redis 缓存）的即时自动化创建（JIT Provisioning）与安全凭据分发，简化开发者的基础设施配置工作。

## Requirements

### Requirement: MySQL 拨备 (MySQL Provisioning)
系统 MUST 能根据请求在宿主 MySQL 实例中创建数据库和用户。

#### Scenario: 应用申请 MySQL
- **WHEN** 应用请求资源类型为 `mysql`，名称为 `billing_db`
- **THEN** 系统应检查 Host 配置中 MySQL 服务是否可用
- **AND** 系统应连接 MySQL 创建数据库 `billing_db`
- **AND** 系统应创建专属用户并授权
- **AND** 系统应生成 DSN 并注入到应用配置中
- **AND** 系统应返回包含 DSN 的配置片段

### Requirement: 资源鉴权 (Resource Auth)
系统生成的资源凭据 MUST 仅对申请的应用可见。

#### Scenario: 凭据存储
- **WHEN** 资源创建成功
- **THEN** 凭据应被加密或安全地存储在 Server 的 SQLite 配置表中
- **AND** 仅能通过该应用的 Config 接口获取
