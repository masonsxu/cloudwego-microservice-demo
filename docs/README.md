# CloudWeGo 微服务实践项目文档

欢迎阅读 CloudWeGo 微服务实践项目文档。本项目是一个基于 CloudWeGo 生态构建的生产级微服务架构 Demo，采用 Kitex (RPC) + Hertz (HTTP) 框架，实现了以网关为核心的星型微服务体系。

## 文档目录

| 文档 | 说明 | 适合读者 |
|------|------|----------|
| [01-快速开始](01-getting-started.md) | 环境配置、项目启动、首次运行 | 新用户 |
| [02-架构设计](02-architecture.md) | 微服务架构、分层设计、技术选型 | 架构师、高级开发者 |
| [03-开发指南](03-development.md) | IDL-First 流程、代码生成、分层规范 | 开发者 |
| [04-配置参考](04-configuration.md) | 环境变量、数据库、JWT、Redis 配置 | 运维、开发者 |
| [05-部署指南](05-deployment.md) | Docker 部署、生产环境配置 | 运维 |
| [06-故障排查](06-troubleshooting.md) | 常见问题、错误诊断、日志调试 | 所有人 |

## 核心功能

- **API 网关 (Hertz)**: 统一 API 入口、JWT 认证、OpenTelemetry 链路追踪
- **身份服务 (Kitex)**: 用户管理、组织架构、角色管理、JWT Token 管理、菜单配置
- **可观测性**: OpenTelemetry 分布式链路追踪、结构化日志
- **安全性**: JWT 认证、密码加密

## 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| RPC 框架 | Kitex | 字节跳动开源高性能 RPC 框架 |
| HTTP 框架 | Hertz | 字节跳动开源高性能 HTTP 框架 |
| 接口定义 | Thrift IDL | IDL-First 开发模式 |
| 数据库 | PostgreSQL + GORM | 关系型数据库 |
| 缓存 | Redis | 会话管理、热点数据缓存 |
| 服务发现 | etcd | 服务注册与发现 |
| 依赖注入 | Google Wire | 编译时依赖注入 |
| 认证 | JWT | 无状态 Token 认证 |
| 可观测性 | OpenTelemetry + Jaeger | 分布式链路追踪 |

## 阅读建议

### 新用户

1. 先阅读 [01-快速开始](01-getting-started.md)，5 分钟启动项目
2. 浏览 [02-架构设计](02-architecture.md) 了解整体设计
3. 需要开发时参考 [03-开发指南](03-development.md)

### 开发者

1. [03-开发指南](03-development.md) - 日常开发必读
2. [04-配置参考](04-configuration.md) - 环境变量参考
3. [06-故障排查](06-troubleshooting.md) - 遇到问题时查阅

### 运维人员

1. [05-部署指南](05-deployment.md) - 部署最佳实践
2. [04-配置参考](04-configuration.md) - 生产环境配置
3. [06-故障排查](06-troubleshooting.md) - 问题诊断

## 外部资源

### CloudWeGo 生态

- [Kitex 官方文档](https://www.cloudwego.io/zh/docs/kitex/)
- [Hertz 官方文档](https://www.cloudwego.io/zh/docs/hertz/)
- [CloudWeGo 官网](https://www.cloudwego.io/)

### 核心依赖

- [GORM 文档](https://gorm.io/zh_CN/docs/)
- [Wire 指南](https://github.com/google/wire/blob/main/docs/guide.md)
- [OpenTelemetry 文档](https://opentelemetry.io/docs/)

---

> **提示**：如果你是第一次接触本项目，建议从 [01-快速开始](01-getting-started.md) 开始阅读。
