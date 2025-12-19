# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

cloudwego-microservice-demo 是一个基于 Go 语言的微服务项目，采用 CloudWeGo 技术栈，使用 Kitex (RPC) 和 Hertz (HTTP) 框架。项目采用 Go Workspace 管理多个服务模块，遵循 IDL-First 开发模式。

当前包含以下服务：

- `gateway`: HTTP 网关服务（基于 Hertz，提供统一的 API 入口）
- `identity_srv`: 身份认证 RPC 服务（用户身份、认证授权、组织管理等）

## 核心技术栈

- **Go**: 1.24+
- **RPC 框架**: Kitex (CloudWeGo)
- **HTTP 框架**: Hertz (CloudWeGo)
- **接口协议**: Thrift
- **数据库**: PostgreSQL + GORM
- **依赖注入**: Google Wire
- **日志**: zerolog
- **代码检查**: golangci-lint

## 常用开发命令

### 服务启动

```bash
# Docker 环境（推荐）
cd docker
./deploy.sh up                  # 启动所有服务（基础设施 + 应用）
./deploy.sh up -d               # 后台启动所有服务
./deploy.sh up-base             # 仅启动基础设施（postgres, etcd, rustfs, redis）
./deploy.sh up-app              # 仅启动应用服务（identity_srv, gateway）
./deploy.sh down                # 停止所有服务
./deploy.sh logs                # 查看所有日志
./deploy.sh logs gateway        # 查看特定服务日志
./deploy.sh ps                  # 查看服务状态

# Docker 原生命令（更多功能）
docker logs -f backend-gateway-1    # 实时跟踪日志
docker exec -it backend-gateway-1 sh  # 进入容器

# 本地开发模式
# 1. 运行 identity_srv RPC 服务
cd rpc/identity_srv
sh build.sh && sh output/bootstrap.sh

# 2. 运行 gateway HTTP 服务
cd gateway
sh build.sh && sh output/bootstrap.sh
```

### 测试命令

```bash
# 运行所有测试
go test ./... -v

# 运行单个包的测试
go test -v ./biz/logic/...

# 生成测试覆盖率报告
go test ./... -coverprofile=coverage.out -v
go tool cover -html=coverage.out

# 运行集成测试
go test -v ./integration_test.go

# 运行性能测试
go test -bench=. -benchmem ./...
```

### 代码生成

```bash
# 生成 Kitex RPC 代码（在 RPC 服务目录下执行）
cd rpc/identity_srv
./script/gen_kitex_code.sh

# 生成 Hertz HTTP 代码（在 gateway 目录下执行）
cd gateway
./script/gen_hertz_code.sh               # 生成所有服务代码
./script/gen_hertz_code.sh identity      # 仅生成 identity 服务代码

# 生成 Wire 依赖注入代码
cd rpc/identity_srv/wire && wire         # RPC 服务
cd gateway/internal/wire && wire         # HTTP 网关
```

### 构建和部署

```bash
# 构建 RPC 服务
cd rpc/identity_srv && ./build.sh

# 构建 HTTP 网关（会自动生成 Swagger 文档）
cd gateway && ./build.sh

# Docker 镜像构建
cd docker
./deploy.sh up --build-only           # 仅构建所有镜像
./deploy.sh up-base --build-only      # 仅构建基础设施镜像
./deploy.sh up-app --build-only       # 仅构建应用镜像
./deploy.sh up --no-build             # 启动但不构建镜像

# 代码检查和格式化
golangci-lint run                    # 运行代码检查
golangci-lint run --fix              # 自动修复问题
```

## 架构设计

### 微服务架构

项目严格遵循 IDL-First 开发模式和分层架构：

1. **HTTP 网关** (`gateway/`)

   - 统一入口，负责请求路由、身份认证、协议转换
   - 基于 Hertz 框架，处理 HTTP 请求并转换为 RPC 调用
   - **⭐ 所有 RPC 调用的唯一发起点**（RPC 服务不互相调用，简化微服务拓扑）
   - 所有安全相关逻辑（JWT、权限校验）必须在网关层处理
   - 自动生成 Swagger API 文档

2. **RPC 服务** (`rpc/*/`)

   - `identity_srv`: 身份认证服务，封装用户身份、认证授权、组织管理等业务逻辑
   - 基于 Kitex 框架，提供高性能 RPC 服务
   - **⭐ RPC 服务不直接调用其他 RPC 服务，所有服务间通信通过网关协调**
   - 禁止在 RPC 服务中处理身份验证或权限逻辑

3. **基础设施服务**（通过 Docker 部署）
   - PostgreSQL: 关系型数据库
   - etcd: 服务注册与发现
   - RustFS: S3 兼容的对象存储服务

### RPC 服务分层结构

```
rpc/<service_name>/
├── handler.go          # RPC 接口实现层 (适配层)
├── biz/                # 核心业务逻辑层
│   ├── converter/      # DTO <-> Model 转换
│   │   ├── base/       # 基础转换器
│   │   ├── user_profile/
│   │   ├── organization/
│   │   └── ...
│   ├── dal/            # 数据访问层
│   │   ├── base/       # 基础数据访问组件
│   │   ├── user_profile/
│   │   ├── organization/
│   │   └── ...
│   └── logic/          # 业务逻辑实现
│       ├── user_profile/
│       ├── organization/
│       └── ...
├── models/             # GORM 数据模型
├── kitex_gen/          # IDL 生成代码 (勿手动修改)
├── config/             # 服务配置
├── wire/               # Wire 依赖注入配置
├── pkg/                # 工具包和通用组件
└── internal/           # 内部实现
    └── middleware/     # RPC 中间件
        └── meta_middleware.go  # 追踪中间件
```

**架构约束**：

- ⭐ **服务调用规则**: RPC 服务只接收来自网关的调用，不调用其他 RPC 服务
- ⭐ **简化拓扑**: 星型微服务架构（网关为中心），避免服务间复杂依赖
- **通信方式**: 仅通过网关协调服务间通信，统一管理追踪、认证、降级等横切关注点

### RPC 追踪中间件设计

项目采用基于 **metainfo** 的链路追踪机制，使用 **request_id** 追踪请求在微服务调用链中的完整路径。

**核心特性**：

- ✅ **唯一追踪标识**: 使用 `request_id` 进行全链路追踪
- ✅ **自动生成追踪 ID**: 缺失的 request_id 自动生成（网关层或 RPC 服务层）
- ✅ **直接使用 metainfo**: 不使用 context.WithValue 重复存储（性能优化）
- ✅ **100% 可追踪性**: 确保每个请求都有完整的追踪信息

**中间件位置**：

```
rpc/<service_name>/internal/middleware/meta_middleware.go
```

**使用方式**：

```go
// 在业务代码中获取追踪 ID
requestID := middleware.GetRequestID(ctx)

// 在日志中添加追踪信息
logger.Info().
    Str("request_id", requestID).
    Msg("Processing request")
```

**设计原则**：

1. **性能优先**: 直接从 metainfo 读取，避免 context.WithValue 的开销
2. **防御性编程**: 自动生成缺失的 request_id，应对网关异常
3. **服务独立**: 每个 RPC 服务保持独立的中间件实现（便于定制）
4. **最佳实践**: 遵循 CloudWeGo 和 Google 的微服务追踪标准

**追踪链传播**：

```
HTTP 请求 → 网关
  ↓
1. requestid 中间件生成 request_id
2. trace_middleware 将 request_id 注入到 metainfo
  ↓
3. RPC 客户端调用（使用 TTHeader MetaHandler）
   metainfo 中的 request_id 通过 TTHeader 自动传递
  ↓
4. RPC 服务端接收
   meta_middleware 从 metainfo 提取 request_id
  ↓
5. 业务逻辑层
   使用 GetRequestID(ctx) 访问追踪信息
```

**注意**: 由于架构限制 RPC 服务不互相调用，追踪链简化为：`客户端 → 网关 → 单个 RPC 服务`

### HTTP 网关分层结构

```
gateway/
├── biz/                  # HTTP 业务层 (IDL生成)
│   ├── handler/          # HTTP Handler 实现
│   ├── model/            # HTTP DTO
│   └── router/           # 路由注册
├── internal/             # 内部实现
│   ├── application/      # 应用层
│   │   ├── assembler/    # 数据组装器
│   │   ├── middleware/   # 中间件 (JWT, CORS, 上下文等)
│   │   └── context/      # 上下文管理
│   ├── domain/           # 领域层
│   │   ├── service/      # 领域服务
│   │   └── common/       # 通用组件
│   ├── infrastructure/   # 基础设施层
│   │   ├── client/       # RPC 客户端封装
│   │   ├── config/       # 配置管理
│   │   └── errors/       # 统一错误处理
│   └── wire/             # 依赖注入配置
├── docs/                 # Swagger 文档（自动生成）
├── pkg/                  # 工具包
└── main.go
```

## 开发规范

### IDL-First 开发流程

1. **定义接口**: 修改 `idl/` 目录下的 Thrift 文件定义接口
2. **生成代码**: 使用 Kitex/Hertz 工具生成服务代码
3. **实现业务逻辑**: 在 `biz/` 目录下实现具体业务逻辑
4. **测试验证**: 编写单元测试和集成测试

### 分层职责

- **Handler 层**: 参数校验、调用转换器、委托业务逻辑层
- **Logic 层**: 核心业务逻辑实现，编排 DAL 层操作
- **DAL 层**: 数据持久化，封装 GORM 操作
- **Converter 层**: DTO 与 Model 之间的纯函数转换

### RPC 调用规范

⭐ **关键架构约束**：

- **RPC 服务不互相调用**: 所有 RPC 调用必须由网关层发起
- **服务间通信**: 如需多个 RPC 服务协作，必须通过网关编排
- **拓扑简化**: 星型架构（网关为中心）避免服务间复杂依赖
- **追踪链简化**: 调用路径固定为 `客户端 → 网关 → RPC 服务`

**设计优势**：
- ✅ 简化微服务拓扑，避免调用链复杂化
- ✅ 统一管理认证、追踪、降级等横切关注点
- ✅ 易于监控和故障排查
- ✅ 降低服务间耦合，提高系统可维护性

### 配置管理

项目使用 Viper 进行配置管理，采用环境驱动配置模式：

**配置优先级**（从高到低）：

1. 系统环境变量
2. `.env` 文件（仅在环境变量未设置时加载）
3. `config/defaults.go` 中的默认值

**重要配置约定**：

- **不使用 YAML 配置文件**：所有配置通过环境变量或 `.env` 文件提供
- **环境变量映射**：`config/env.go` 定义环境变量到配置结构的映射
- **Duration 类型**：支持多种格式（`1h`, `30m`, `3600s` 或纯数字秒）
- **类型转换**：使用 `mapToViper` 和 `parseDurationWithDefault` 统一处理

**示例配置加载**：

```go
// config/env.go
mapToViper(v, "DB_CONN_MAX_LIFETIME", "database.conn_max_lifetime", func(value string) interface{} {
    return parseDurationWithDefault(value, 60*time.Minute)
})
```

### Wire 依赖注入

项目使用 Google Wire 管理依赖注入：

- 每个服务在 `wire/` 目录配置依赖注入
- 分层组织 Provider Sets (基础设施层、数据层、业务层等)
- 运行 `wire` 命令生成 `wire_gen.go` 文件
- 遵循接口化原则，便于测试和解耦

**Wire 最佳实践**：

- **Provider 职责单一**：每个 Provider 函数只负责创建一个依赖
- **委托 config 层**：数据库初始化等逻辑应在 `config` 层实现，Wire 层只调用
- **参数传递完整**：确保所有依赖参数（如 logger）正确传递
- **避免重复逻辑**：不在 Provider 中重复配置已在 config 层完成的初始化

**示例**：

```go
// wire/provider.go - 正确做法
func ProvideDB(cfg *config.Config, logger *slog.Logger) (*gorm.DB, error) {
    return config.InitDB(cfg, logger)  // 委托给 config 层
}

// config/database.go - 完整的初始化逻辑
func InitDB(cfg *Config, logger *slog.Logger) (*gorm.DB, error) {
    // 连接创建、连接池配置、自动迁移等所有逻辑
}
```

### 数据库管理

**数据库初始化**：

- 自动迁移在 `config/database.go` 的 `InitDB` 函数中执行
- 支持连接池配置（MaxIdleConns, MaxOpenConns, ConnMaxLifetime）
- 开发/生产环境使用不同的连接池参数
- 使用 GORM AutoMigrate 自动创建/更新表结构

**迁移策略**：

- 开发环境：自动迁移 + Debug 模式
- 生产环境：需要人工审核迁移 SQL

### 文件存储 (RustFS)

项目集成了 S3 兼容的文件存储服务用于组织 Logo 管理：

**AWS SDK 使用规范**：

- 使用 AWS SDK Go v2
- **废弃的 API**：不要使用 `WithEndpointResolverWithOptions`
- **推荐做法**：在 S3 客户端配置中使用 `o.BaseEndpoint`

**双端点配置（容器化部署最佳实践）**：

为了解决容器化部署中内外部访问的问题，项目采用双端点配置：

1. **内部端点** (`LOGO_STORAGE_S3_ENDPOINT`)：

   - 用于容器间通信（上传、删除等服务端操作）
   - 使用容器名称或内部网络地址
   - 示例：`http://rustfs:9000`

2. **公共端点** (`LOGO_STORAGE_S3_PUBLIC_ENDPOINT`)：
   - 用于生成预签名 URL（浏览器访问）
   - 使用外部可访问的地址（localhost、域名或公网地址）
   - 开发环境示例：`http://localhost:9000`
   - 生产环境示例：`https://s3.your-domain.com`

**配置示例**：

```bash
# 开发环境 (.env.dev.example)
LOGO_STORAGE_S3_ENDPOINT=http://rustfs:9000          # 容器内部访问
LOGO_STORAGE_S3_PUBLIC_ENDPOINT=http://localhost:9000  # 浏览器访问

# 生产环境 (.env.prod.example)
LOGO_STORAGE_S3_ENDPOINT=http://rustfs:9000          # 容器内部访问
LOGO_STORAGE_S3_PUBLIC_ENDPOINT=https://s3.your-domain.com  # 浏览器访问
```

**实现原理**：

Logo 存储客户端创建两个 S3 客户端实例：

- 内部客户端：用于所有服务端操作（上传、删除、标签管理等）
- 公共客户端：专门用于生成预签名 URL

```go
// 创建内部端点客户端（容器间通信）
s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
    o.BaseEndpoint = aws.String(cfg.S3Endpoint)
    o.UsePathStyle = true
})

// 创建公共端点客户端（生成预签名 URL）
s3PublicClient := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
    o.BaseEndpoint = aws.String(cfg.S3PublicEndpoint)
    o.UsePathStyle = true
})

// 生成预签名 URL 时使用公共客户端
presignClient := s3.NewPresignClient(c.s3PublicClient)
```

**核心配置项**：

- `LOGO_STORAGE_S3_ENDPOINT`: S3 内部端点地址（容器间通信）
- `LOGO_STORAGE_S3_PUBLIC_ENDPOINT`: S3 公共端点地址（生成预签名 URL）
- `LOGO_STORAGE_S3_REGION`: S3 区域
- `LOGO_STORAGE_S3_USE_SSL`: 是否使用 SSL
- `LOGO_STORAGE_ACCESS_KEY`: S3 访问密钥
- `LOGO_STORAGE_SECRET_KEY`: S3 私钥
- `LOGO_STORAGE_MAX_FILE_SIZE`: 最大文件大小（字节，默认 10MB）
- `LOGO_STORAGE_ALLOWED_FILE_TYPES`: 允许的文件类型（逗号分隔）

**优势**：

- ✅ 容器内部通信高效（使用 Docker 网络）
- ✅ 外部客户端可访问预签名 URL
- ✅ 支持不同部署环境（开发/生产）
- ✅ 符合容器化微服务最佳实践
- ✅ 无需额外的网络配置或反向代理

### 错误处理

- 使用 6 位数字业务错误码 (参考 `doc/01-开发规范/RPC业务错误处理规范.md`)
- RPC 错误与业务错误分离处理
- 业务错误使用 Kitex 的 `BizStatusError` 机制通过 TTHeader 传递
- Handler 层使用 `errno.ToKitexError()` 转换业务错误
- 网关层负责处理 BizStatusError 并统一错误响应格式
- DAL 层负责将 GORM 错误转换为业务错误码

### 代码检查配置

项目配置了 `.golangci.yml` 进行代码检查：

- 启用多种 Linters: gocyclo, nestif, prealloc, revive, staticcheck
- 排除生成代码 (`kitex_gen/`) 检查
- 行长度限制 120 字符，支持自动换行
- 使用 gofumpt 进行代码格式化

## 关键约定

1. **安全原则**: 所有鉴权逻辑必须在 API 网关层，RPC 服务不处理权限
2. **无状态设计**: RPC 服务设计为无状态，支持水平扩展
3. **接口稳定性**: Thrift 接口保持向后兼容
4. **测试覆盖**: 为 logic 和 dal 层编写充分的单元测试
5. **代码生成**: 不要手动修改 `kitex_gen/` 目录下的生成代码
6. **Go Workspace**: 使用 `go.work` 管理多模块项目
7. **架构约束**:
   - ⭐ RPC 服务不直接调用其他 RPC 服务
   - ⭐ 所有 RPC 调用必须由网关层发起
   - 服务间通信通过网关协调，简化微服务拓扑为星型结构
8. **错误处理规范**:
   - DAL 层: 将 GORM 错误转换为 ErrNo 业务错误
   - Logic 层: 处理业务逻辑错误，返回 ErrNo
   - Handler 层: 使用 `errno.ToKitexError()` 将 ErrNo 转换为 BizStatusError
   - 客户端: 使用 `kerrors.FromBizStatusError()` 解析业务错误
9. **RPC 客户端配置**: 使用 TTHeader MetaHandler 以支持 BizStatusError 传递
10. **配置文件约定**:
   - 不使用 YAML 配置文件
   - 所有配置通过环境变量或 `.env` 文件管理
   - `.env.example` 文件记录所有可用的环境变量及说明

## Git 提交规范

### Git Hooks

项目配置了 pre-commit 钩子，在提交时自动执行以下检查：

1. **文件名检查**: 确保文件名为 ASCII 字符（跨平台兼容）
2. **文件大小检查**: 阻止提交超过 10MB 的大文件

**安装 Git Hooks**：

```bash
# 在项目根目录执行
ln -s -f ../../scripts/git-hooks/pre-commit .git/hooks/pre-commit
```

### 提交消息规范

遵循常规的提交消息格式：

- `feat`: 新功能
- `fix`: 修复 bug
- `refactor`: 重构代码
- `docs`: 文档更新
- `test`: 测试相关
- `chore`: 构建/工具链更新

## 环境配置快速指南

### 配置文件位置

每个服务都有独立的 `.env.example` 配置模板：

```
项目根目录/
├── docker/
│   └── .env.dev.example          # Docker 开发环境配置模板
├── gateway/
│   └── .env.example              # HTTP 网关配置模板
└── rpc/identity_srv/
    └── .env.example              # 身份认证服务配置模板
```

### 首次配置步骤

#### 1. Docker 容器化部署（推荐）

```bash
# 进入 docker 目录
cd docker

# 复制环境配置文件
cp .env.dev.example .env

# 编辑配置（可选，默认配置通常可以直接使用）
vim .env  # 或使用你喜欢的编辑器

# 启动所有服务
./deploy.sh up
```

**Docker 部署注意事项**：

- 开发环境默认配置已优化，通常无需修改即可使用
- 数据库和其他基础设施服务通过 Docker 自动启动
- 所有端口已自动映射到 localhost

#### 2. 本地开发模式

如果需要在本地直接运行服务（不使用 Docker），需要为每个服务配置 `.env` 文件：

```bash
# 配置 identity_srv
cd rpc/identity_srv
cp .env.example .env
vim .env  # 修改数据库连接等配置

# 配置 HTTP 网关
cd ../../gateway
cp .env.example .env
vim .env
```

### 关键配置项说明

#### 数据库配置（所有服务）

```env
DB_HOST=127.0.0.1              # 数据库主机（Docker 内使用服务名 "postgres"）
DB_PORT=5432                    # 数据库端口
DB_USERNAME=postgres            # 数据库用户名
DB_PASSWORD=your-password       # 数据库密码（必须修改）
DB_NAME=identity_srv     # 数据库名称（每个服务不同）
DB_SSLMODE=disable              # SSL 模式（生产环境建议 require）
DB_TIMEZONE=Asia/Shanghai       # 时区设置

# 连接池配置
DB_MAX_IDLE_CONNS=10            # 最大空闲连接数
DB_MAX_OPEN_CONNS=100           # 最大打开连接数
DB_CONN_MAX_LIFETIME=1h         # 连接最大生命周期（支持 1h、60m、3600s 或纯数字秒）
DB_CONN_MAX_IDLE_TIME=5m        # 连接最大空闲时间
```

#### 服务注册发现配置（所有服务）

```env
ETCD_ADDRESS=127.0.0.1:2379     # etcd 地址（Docker 内使用 "etcd:2379"）
ETCD_USERNAME=                  # etcd 用户名（可选）
ETCD_PASSWORD=                  # etcd 密码（可选）
ETCD_TIMEOUT=5                  # 连接超时时间（秒）
```

#### RPC 服务端口配置

```env
# identity_srv
SERVER_ADDRESS=:8891            # RPC 服务监听地址
HEALTH_CHECK_PORT=10000         # 健康检查端口

# gateway
SERVER_HOST=0.0.0.0
SERVER_PORT=8080                # HTTP 服务端口
```

#### JWT 认证配置（gateway）

```env
JWT_ENABLED=true
JWT_SIGNING_KEY=your-jwt-secret-key    # ⚠️ 生产环境必须修改为强密钥
JWT_TIMEOUT=30m                         # Token 有效期
JWT_MAX_REFRESH=168h                    # 最大刷新时间（7天）

# 跳过认证的路径（逗号分隔）
JWT_SKIP_PATHS=/api/v1/identity/auth/login,/api/v1/identity/auth/refresh,/ping,/health

# Cookie 配置（生产环境必须启用安全选项）
JWT_COOKIE_SEND_COOKIE=true
JWT_COOKIE_HTTP_ONLY=true               # 防止 XSS 攻击
JWT_COOKIE_SECURE_COOKIE=false          # 生产环境改为 true（需要 HTTPS）
JWT_COOKIE_COOKIE_SAME_SITE=lax         # CSRF 防护
```

#### 文件存储配置（identity_srv - 组织 Logo）

```env
LOGO_STORAGE_S3_ENDPOINT=http://localhost:9000
LOGO_STORAGE_ACCESS_KEY=RustFSadmin
LOGO_STORAGE_SECRET_KEY=
LOGO_STORAGE_MAX_FILE_SIZE=10485760     # 10MB
LOGO_STORAGE_ALLOWED_FILE_TYPES=image/jpeg,image/png,image/gif,image/webp,image/svg+xml
```

### 配置验证清单

部署前请确认以下配置项：

- [ ] **数据库连接**：`DB_HOST`、`DB_PORT`、`DB_USERNAME`、`DB_PASSWORD` 已正确配置
- [ ] **数据库名称**：每个服务的 `DB_NAME` 不同
  - identity_srv: `identity_srv`
- [ ] **etcd 地址**：`ETCD_ADDRESS` 指向正确的服务注册中心
- [ ] **服务端口**：确保端口未被占用
  - identity_srv: 8891 (RPC), 10000 (Health)
  - gateway: 8080 (HTTP)
- [ ] **JWT 密钥**：生产环境必须修改 `JWT_SIGNING_KEY` 为强随机字符串
- [ ] **安全选项**：生产环境启用 `JWT_COOKIE_SECURE_COOKIE=true` 和 `JWT_COOKIE_HTTP_ONLY=true`
- [ ] **日志级别**：生产环境建议使用 `LOG_LEVEL=info` 或 `warn`

### 环境差异对照表

| 配置项                     | 开发环境     | 生产环境                               |
| -------------------------- | ------------ | -------------------------------------- |
| `SERVER_DEBUG`             | `true`       | `false`（临时排障时可设为 `true`）     |
| `LOG_LEVEL`                | `debug`      | `info` 或 `warn`                       |
| `DB_PASSWORD`              | 弱密码可接受 | **必须使用强密码**                     |
| `JWT_SIGNING_KEY`          | 简单字符串   | **必须使用强随机密钥（至少 32 字符）** |
| `JWT_COOKIE_SECURE_COOKIE` | `false`      | `true`（需要 HTTPS）                   |
| `DB_SSLMODE`               | `disable`    | `require` 或 `verify-full`             |
| `GORM SQL 日志`            | 启用（自动） | 禁用（仅记录错误）                     |

## 环境变量参考

### 通用环境变量

| 变量名         | 说明                                       | 示例值        | 必需 |
| -------------- | ------------------------------------------ | ------------- | ---- |
| `SERVER_DEBUG` | 调试模式（控制日志详细度、GORM SQL 日志） | `true`/`false` | 否   |

### 数据库配置

| 变量名                  | 说明             | 示例值         | 必需 |
| ----------------------- | ---------------- | -------------- | ---- |
| `DB_HOST`               | 数据库主机       | `localhost`    | 是   |
| `DB_PORT`               | 数据库端口       | `5432`         | 是   |
| `DB_USERNAME`           | 数据库用户名     | `postgres`     | 是   |
| `DB_PASSWORD`           | 数据库密码       | `password`     | 是   |
| `DB_NAME`               | 数据库名称       | `identity_srv` | 是   |
| `DB_CONN_MAX_LIFETIME`  | 连接最大生命周期 | `1h`/`3600`    | 否   |
| `DB_CONN_MAX_IDLE_TIME` | 连接最大空闲时间 | `5m`/`300`     | 否   |

### 日志配置

| 变量名       | 说明     | 示例值                        | 必需 |
| ------------ | -------- | ----------------------------- | ---- |
| `LOG_LEVEL`  | 日志级别 | `debug`/`info`/`warn`/`error` | 否   |
| `LOG_FORMAT` | 日志格式 | `json`/`text`                 | 否   |
| `LOG_OUTPUT` | 日志输出 | `stdout`/`file`               | 否   |

### 服务注册发现

| 变量名         | 说明         | 示例值           | 必需 |
| -------------- | ------------ | ---------------- | ---- |
| `ETCD_ADDRESS` | etcd 地址    | `127.0.0.1:2379` | 是   |
| `ETCD_TIMEOUT` | 连接超时时间 | `5` (秒) 或 `5s` | 否   |

## 常见问题

### 配置相关

**Q: 为什么不使用 config.yaml 文件？**
A: 项目已重构为纯环境变量驱动配置，便于容器化部署和环境隔离。所有配置通过 `.env` 文件或系统环境变量提供。

**Q: Duration 类型的环境变量如何设置？**
A: 支持多种格式：

- 带单位: `1h`, `30m`, `3600s`
- 纯数字: `3600` (默认按秒解析)

**Q: 如何添加新的环境变量？**
A:

1. 在 `config/types.go` 中添加配置字段
2. 在 `config/defaults.go` 中设置默认值
3. 在 `config/env.go` 中添加环境变量映射
4. 在 `.env.example` 中添加说明和示例

### Wire 依赖注入

**Q: 修改 Provider 后如何重新生成？**
A: 在 `wire/` 目录下运行 `wire` 命令，自动生成 `wire_gen.go`

**Q: Wire 生成失败怎么办？**
A:

1. 检查 Provider 函数参数和返回值类型是否匹配
2. 确保所有依赖都有对应的 Provider
3. 查看错误信息中的具体提示

### 数据库相关

**Q: 如何添加新的数据模型？**
A:

1. 在 `models/` 目录创建 GORM 模型
2. 在 `config/database.go` 的 `AutoMigrate` 中添加模型
3. 重启服务，自动创建/更新表结构

**Q: 生产环境如何管理数据库迁移？**
A: 建议关闭自动迁移，使用专门的迁移工具（如 golang-migrate）进行版本化管理

### 开发调试

**Q: 如何启用数据库 SQL 日志？**
A: 在开发环境中，GORM 已配置为 Debug 模式，会自动打印 SQL 语句

**Q: 如何调试 RPC 调用？**
A:

1. 检查 RPC 服务是否正常启动（端口监听）
2. 查看日志中的 BizStatusError 信息
3. 使用 Kitex 提供的调试工具

## 故障排查

### 服务启动问题

#### 问题：端口已被占用

**错误信息**：

```
bind: address already in use
listen tcp :8891: bind: address already in use
```

**解决方法**：

1. **停止占用端口的进程**：

   ```bash
   # 查找并终止进程
   kill -9 $(lsof -t -i:8891)
   ```

2. **或者修改服务端口**（在 `.env` 文件中）：

   ```env
   SERVER_ADDRESS=:8893  # 改用其他端口
   ```

3. **Docker 环境清理**：
   ```bash
   cd docker
   ./deploy.sh down
   ./deploy.sh up
   ```

#### 问题：数据库连接失败

**错误信息**：

```
failed to connect to database
dial tcp 127.0.0.1:5432: connect: connection refused
pq: password authentication failed for user "postgres"
```

**诊断方法**：

```bash
# 1. 检查 PostgreSQL 是否运行（Docker 环境）
cd docker && ./deploy.sh ps

# 2. 检查端口监听
netstat -tuln | grep 5432
lsof -i :5432

# 3. 测试数据库连接
psql -h localhost -p 5432 -U postgres -d identity_srv

# 4. Docker 环境查看日志
cd docker && ./deploy.sh logs postgres
```

**解决方法**：

1. **启动 PostgreSQL 服务**：

   ```bash
   # Docker 环境
   cd docker && ./deploy.sh up-base
   ```

2. **检查 `.env` 配置**：

   ```env
   DB_HOST=127.0.0.1           # 本地开发
   # DB_HOST=postgres          # Docker 内部使用服务名
   DB_PORT=5432
   DB_USERNAME=postgres        # 确认用户名正确
   DB_PASSWORD=your-password   # 确认密码正确
   DB_NAME=identity_srv # 确认数据库名称
   ```

3. **创建数据库**（如果不存在）：

   ```bash
   # 连接到 PostgreSQL
   psql -h localhost -U postgres

   # 创建数据库
   CREATE DATABASE identity_srv;
   ```

#### 问题：etcd 连接失败

**错误信息**：

```
failed to connect to etcd
context deadline exceeded
etcd client: no endpoints available
```

**诊断方法**：

```bash
# 检查 etcd 服务状态
cd docker && ./deploy.sh ps

# 检查端口
lsof -i :2379

# 测试 etcd 连接
curl http://localhost:2379/version

# Docker 环境查看日志
cd docker && ./deploy.sh logs etcd
```

**解决方法**：

1. **启动 etcd 服务**：

   ```bash
   # Docker 环境
   cd docker && ./deploy.sh up-base
   ```

2. **检查配置**：

   ```env
   # 本地开发
   ETCD_ADDRESS=127.0.0.1:2379

   # Docker 内部
   ETCD_ADDRESS=etcd:2379
   ```

3. **验证 etcd 健康状态**：
   ```bash
   # Docker 环境
   docker exec etcd etcdctl endpoint health
   ```

### 代码生成问题

#### 问题：kitex/hz 命令未找到

**错误信息**：

```
command not found: kitex
command not found: hz
command not found: thriftgo
```

**解决方法**：

```bash
# 安装 Kitex 工具链
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
go install github.com/cloudwego/thriftgo@latest

# 安装 Hertz 工具
go install github.com/cloudwego/hertz/cmd/hz@latest

# 确保 $GOPATH/bin 在 PATH 中
export PATH=$PATH:$(go env GOPATH)/bin

# 验证安装
kitex --version
hz --version
thriftgo --version
```

#### 问题：IDL 文件未找到

**错误信息**：

```
open ../../idl/rpc/identity_srv/identity_service.thrift: no such file or directory
```

**解决方法**：

1. **确认在正确的目录执行**：

   ```bash
   # 必须在服务根目录执行
   cd rpc/identity_srv
   ./script/gen_kitex_code.sh
   ```

2. **检查 IDL 文件路径**：
   ```bash
   # 确认文件存在
   ls -la ../../idl/rpc/identity_srv/identity_service.thrift
   ```

#### 问题：Wire 生成失败

**错误信息**：

```
wire: no provider found for *gorm.DB
wire: cycle detected in provider set
```

**诊断方法**：

```bash
# 查看详细错误
cd wire
wire
```

**解决方法**：

1. **缺少 Provider**：

   ```go
   // 在 wire/wire.go 中添加缺失的 Provider
   var ProviderSet = wire.NewSet(
       config.ProvideConfig,
       config.ProvideDB,        // 确保包含数据库 Provider
       config.ProvideLogger,    // 确保包含日志 Provider
       // ... 其他 Providers
   )
   ```

2. **循环依赖**：

   - 检查 Provider 之间的依赖关系
   - 使用接口打破循环依赖
   - 重构 Provider 函数

3. **重新生成**：
   ```bash
   cd wire
   rm wire_gen.go  # 删除旧文件
   wire            # 重新生成
   ```

### 运行时问题

#### 问题：自动迁移失败

**错误信息**：

```
AutoMigrate failed: ERROR: permission denied for schema public
AutoMigrate failed: relation "users" already exists
```

**解决方法**：

1. **权限问题**：

   ```sql
   -- 授予权限
   GRANT ALL PRIVILEGES ON DATABASE identity_srv TO postgres;
   GRANT ALL ON SCHEMA public TO postgres;
   ```

2. **表已存在但结构不同**：

   ```bash
   # 开发环境：删除并重新创建
   psql -h localhost -U postgres -c "DROP DATABASE identity_srv;"
   psql -h localhost -U postgres -c "CREATE DATABASE identity_srv;"

   # 生产环境：使用迁移工具手动管理
   ```

3. **禁用自动迁移**（如果需要手动管理）：
   - 修改 `config/database.go` 中的 `InitDB` 函数
   - 注释掉 `AutoMigrate` 调用

#### 问题：RPC 调用超时

**错误信息**：

```
rpc timeout: deadline exceeded
context deadline exceeded
```

**诊断方法**：

```bash
# 1. 检查目标服务是否运行
cd docker && ./deploy.sh ps

# 2. 检查网络连接
ping localhost
telnet localhost 8891

# 3. 查看服务日志
cd docker && ./deploy.sh logs identity_srv
# 或使用 Docker 原生命令实时跟踪
docker logs -f backend-identity_srv-1
```

**解决方法**：

1. **增加超时时间**（gateway `.env`）：

   ```env
   CLIENT_REQUEST_TIMEOUT=60s     # 从 30s 增加到 60s
   CLIENT_CONNECTION_TIMEOUT=5s   # 连接超时
   ```

2. **检查服务注册**：

   ```bash
   # 确认服务已注册到 etcd
   docker exec etcd etcdctl get --prefix /kitex
   ```

3. **检查防火墙**：
   ```bash
   # 临时关闭防火墙测试
   sudo ufw status
   sudo ufw disable  # 仅用于测试
   ```

#### 问题：JWT 认证失败

**错误信息**：

```
401 Unauthorized
token is expired
invalid token
```

**解决方法**：

1. **Token 已过期** - 刷新 Token：

   ```bash
   # 调用刷新接口
   curl -X POST http://localhost:8080/api/v1/identity/auth/refresh \
     -H "Authorization: Bearer YOUR_OLD_TOKEN"
   ```

2. **签名密钥不匹配**：

   - 确保 API 网关的 `JWT_SIGNING_KEY` 一致
   - 重启服务后需要重新登录

3. **Token 格式错误**：

   ```bash
   # 正确格式
   Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

   # 错误：缺少 "Bearer " 前缀
   Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
   ```

### Docker 环境问题

#### 问题：容器启动后立即退出

**诊断方法**：

```bash
# 查看容器状态
docker ps -a | grep cloudwego-scaffold

# 查看退出码
docker inspect identity-srv | grep ExitCode

# 查看容器日志
./deploy.sh logs identity_srv
docker logs identity-srv
```

**常见原因**：

1. **配置错误** - 检查 `.env` 文件
2. **依赖服务未就绪** - 先启动基础服务：`./deploy.sh up-base`
3. **端口冲突** - 修改 `docker-compose` 中的端口映射
4. **镜像构建失败** - 清理并重新构建：`./deploy.sh down && ./deploy.sh up`

#### 问题：无法连接到 Docker 内的服务

**解决方法**：

```bash
# 1. 检查容器网络
docker network ls
docker network inspect network

# 2. 使用服务名而非 localhost
# 在 Docker 内部：
DB_HOST=postgres       # 而不是 127.0.0.1
ETCD_ADDRESS=etcd:2379 # 而不是 localhost:2379

# 3. 从宿主机连接使用 localhost
DB_HOST=127.0.0.1      # 宿主机访问
```

### 性能问题

#### 问题：数据库查询慢

**诊断方法**：

```bash
# 启用 SQL 日志（开发环境自动启用）
LOG_LEVEL=debug
SERVER_DEBUG=true

# 查看慢查询
# 修改 postgresql.conf:
# log_min_duration_statement = 100  # 记录超过100ms的查询
```

**优化方法**：

1. 添加数据库索引
2. 使用 `Preload` 避免 N+1 查询
3. 调整连接池配置：
   ```env
   DB_MAX_IDLE_CONNS=20
   DB_MAX_OPEN_CONNS=200
   ```

#### 问题：内存占用过高

**诊断方法**：

```bash
# 查看容器资源使用
docker stats

# 查看 Go 内存分析
# 在代码中启用 pprof
import _ "net/http/pprof"
go tool pprof http://localhost:6060/debug/pprof/heap
```

**解决方法**：

1. 检查是否有内存泄漏
2. 调整 GORM 连接池
3. 使用分页查询，避免一次性加载大量数据
4. 设置 Docker 内存限制（生产环境）

### 日志和调试

#### 查看实时日志

```bash
# 查看所有服务日志
cd docker && ./deploy.sh logs

# 查看特定服务日志
cd docker && ./deploy.sh logs identity_srv
cd docker && ./deploy.sh logs gateway

# 实时跟踪日志（推荐）
# 或使用 Docker 原生命令实时跟踪
docker logs -f backend-identity_srv-1    # 按 Ctrl+C 退出

# 查看最近的日志
cd docker && ./deploy.sh logs identity_srv | tail -n 100

# 过滤错误日志
cd docker && ./deploy.sh logs identity_srv | grep -i error
```

#### 启用详细日志

```env
# .env 配置
LOG_LEVEL=debug              # 最详细的日志
LOG_FORMAT=text              # 更易读的格式（开发环境）
SERVER_DEBUG=true            # 启用调试模式（自动设置日志级别和 SQL 日志）
```

#### 追踪 RPC 调用链

```bash
# 1. 在请求中查找 request_id
# API 响应头会包含：X-Request-ID

# 2. 使用 request_id 搜索日志（主要方式）
cd docker && ./deploy.sh logs | grep "request_id=abc123"

# 3. 查看完整调用链
# 由于 RPC 服务不互相调用，调用链简化为：网关 → 单个 RPC 服务
cd docker && ./deploy.sh logs gateway | grep "request_id=abc123"
cd docker && ./deploy.sh logs identity_srv | grep "request_id=abc123"
```

## 参考文档

项目文档位于 `doc/` 目录（如果存在）。关键开发规范和指南请参考项目内部文档。

## 重要提醒

- ⚠️ **永远不要手动修改** `kitex_gen/` 和生成的 `wire_gen.go` 文件
- ⚠️ **所有鉴权逻辑** 必须在 API Gateway 层实现，RPC 服务不处理权限
- ⚠️ **配置文件已废弃**，使用环境变量或 `.env` 文件替代 `config.yaml`
- ⚠️ **AWS SDK 端点配置** 使用 `BaseEndpoint` 而非废弃的 `WithEndpointResolverWithOptions`
- ⚠️ **数据库初始化逻辑** 应在 `config` 层实现，Wire 层仅负责依赖声明
