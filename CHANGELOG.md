# Changelog

本项目的所有重要变更都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本管理遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### Added
- 文档重组：新增 `docs/` 目录存放详细文档
- 新增 CONTRIBUTING.md 贡献指南
- 新增 SECURITY.md 安全政策
- 新增 CODE_OF_CONDUCT.md 行为准则
- 新增 ROADMAP.md 项目路线图

### Changed
- README.md 精简为快速入门指南

---

## [0.1.0] - 2025-12-02

### Added
- 初始项目结构
- Gateway HTTP 网关服务（基于 Hertz）
- Identity RPC 服务（基于 Kitex）
- 用户管理功能
- 组织管理功能
- JWT 认证中间件
- Docker Compose 部署配置
- Wire 依赖注入
- golangci-lint 代码检查配置
- Git Hooks（pre-commit）

### Infrastructure
- PostgreSQL 数据库支持
- etcd 服务注册发现
- RustFS 对象存储（S3 兼容）

---

## 版本说明

- **Major (X.0.0)**: 不兼容的 API 变更
- **Minor (0.X.0)**: 向后兼容的新功能
- **Patch (0.0.X)**: 向后兼容的 Bug 修复

[Unreleased]: https://github.com/masonsxu/cloudwego-scaffold/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/masonsxu/cloudwego-scaffold/releases/tag/v0.1.0
