# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 提供项目开发规范和技术约束。

## 项目概述

cloudwego-microservice-demo 是一个基于 Go 语言的微服务项目，采用 CloudWeGo 技术栈，使用 Kitex (RPC) 和 Hertz (HTTP) 框架。项目采用 Go Workspace 管理多个服务模块，遵循 IDL-First 开发模式。

**服务列表**：
- `gateway`: HTTP 网关服务（Hertz，端口 8080）
- `identity_srv`: 身份认证 RPC 服务（Kitex，端口 8891）

## 核心技术栈

- **Go**: 1.24+
- **RPC 框架**: Kitex (CloudWeGo)
- **HTTP 框架**: Hertz (CloudWeGo)
- **接口协议**: Thrift
- **数据库**: PostgreSQL + GORM
- **依赖注入**: Google Wire
- **日志**: zerolog
- **代码检查**: golangci-lint

## 关键架构约束

### 微服务架构

1. **星型拓扑**：网关为中心，RPC 服务不互相调用
2. **安全分层**：所有鉴权逻辑必须在网关层，RPC 服务不处理权限
3. **协议转换**：网关负责 HTTP → Thrift RPC 转换

### RPC 服务分层

```
handler.go           # RPC 接口实现（适配层）
biz/
├── logic/           # 业务逻辑实现
├── dal/             # 数据访问层
└── converter/       # DTO ↔ Model 转换
models/              # GORM 数据模型
wire/                # Wire 依赖注入
```

**分层职责**：
- **Handler**: 参数校验、调用转换器、委托业务逻辑
- **Logic**: 核心业务逻辑、编排 DAL 操作
- **DAL**: 数据持久化、封装 GORM 操作
- **Converter**: DTO 与 Model 纯函数转换

### HTTP 网关分层

```
gateway/
├── biz/handler/     # HTTP Handler（IDL 生成）
├── internal/
│   ├── application/ # 应用层（assembler, middleware）
│   ├── domain/      # 领域层（service）
│   ├── infrastructure/ # 基础设施层（client, config）
│   └── wire/        # 依赖注入
```

## 常用命令速查

### 服务启动

```bash
# 1. 启动基础设施（PostgreSQL、etcd、Redis、RustFS、Jaeger）
cd docker && ./deploy.sh up

# 2. 启动 RPC 服务（新终端）
cd rpc/identity_srv && sh build.sh && sh output/bootstrap.sh

# 3. 启动网关服务（新终端）
cd gateway && sh build.sh && sh output/bootstrap.sh

# 基础设施管理
./deploy.sh down                      # 停止基础设施
./deploy.sh ps                        # 查看状态
./deploy.sh logs                      # 查看日志
```

### 代码生成

```bash
# Kitex RPC 代码
cd rpc/identity_srv && ./script/gen_kitex_code.sh

# Hertz HTTP 代码
cd gateway && ./script/gen_hertz_code.sh

# Wire 依赖注入
cd rpc/identity_srv/wire && wire
cd gateway/internal/wire && wire
```

### 测试和构建

```bash
go test ./... -v                      # 运行测试
go test ./... -coverprofile=coverage.out  # 测试覆盖率
golangci-lint run                     # 代码检查
```

## 开发规范

### IDL-First 开发流程

1. 修改 `idl/` 目录下的 Thrift 文件
2. 使用 Kitex/Hertz 工具生成代码
3. 在 `biz/` 目录实现业务逻辑
4. 更新 Wire 依赖注入
5. 编写单元测试

### 错误处理规范

项目采用 6 位数字业务错误码：

```
DAL 层: GORM 错误 → ErrNo
Logic 层: 处理业务逻辑错误，返回 ErrNo
Handler 层: errno.ToKitexError() → BizStatusError
客户端: kerrors.FromBizStatusError() 解析
```

### Wire 依赖注入规范

- Provider 职责单一，每个函数只创建一个依赖
- 数据库初始化在 `config` 层，Wire 层只调用
- 需要清理的资源返回 cleanup 函数

### 配置管理约定

- **不使用 YAML 配置文件**：所有配置通过环境变量或 `.env` 文件
- **优先级**：系统环境变量 > `.env` 文件 > `config/defaults.go`
- **Duration 格式**：支持 `1h`、`30m`、`3600s` 或纯数字秒

### 文件存储（RustFS）

采用双端点配置解决容器化部署问题：
- `LOGO_STORAGE_S3_ENDPOINT`: 内部端点（容器间通信）
- `LOGO_STORAGE_S3_PUBLIC_ENDPOINT`: 公共端点（预签名 URL）

## Git 提交规范

### 提交消息格式

```
feat: 新功能
fix: 修复 bug
refactor: 重构代码
docs: 文档更新
test: 测试相关
chore: 构建/工具链更新
```

### Git Hooks

```bash
# 安装 pre-commit 钩子
ln -s -f ../../scripts/git-hooks/pre-commit .git/hooks/pre-commit
```

## 关键约定清单

1. **安全原则**: 所有鉴权逻辑必须在 API 网关层
2. **无状态设计**: RPC 服务设计为无状态，支持水平扩展
3. **接口稳定性**: Thrift 接口保持向后兼容
4. **代码生成**: 不要手动修改 `kitex_gen/` 目录
5. **Go Workspace**: 使用 `go.work` 管理多模块
6. **RPC 调用规则**: 所有 RPC 调用必须由网关层发起
7. **错误处理**: DAL 层转换 GORM 错误，Handler 层使用 `errno.ToKitexError()`
8. **RPC 客户端**: 使用 TTHeader MetaHandler 支持 BizStatusError

## 详细文档

详细的开发指南、配置说明、部署指南和故障排查请参考 `docs/` 目录：

- [快速开始](docs/01-getting-started.md)
- [架构设计](docs/02-architecture.md)
- [开发指南](docs/03-development.md)
- [配置参考](docs/04-configuration.md)
- [部署指南](docs/05-deployment.md)
- [故障排查](docs/06-troubleshooting.md)
