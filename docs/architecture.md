# 架构设计

本文档介绍 CloudWeGo Scaffold 的整体架构设计和技术选型。

## 目录

- [微服务架构](#微服务架构)
- [服务列表](#服务列表)
- [RPC 服务分层架构](#rpc-服务分层架构)
- [HTTP 网关分层架构](#http-网关分层架构)
- [技术栈](#技术栈)
- [关键设计决策](#关键设计决策)

---

## 微服务架构

```mermaid
flowchart TB
    subgraph ClientLayer["客户端层"]
        WebClient[Web 客户端]
        MobileClient[移动客户端]
    end

    subgraph GatewayLayer["API 网关层"]
        Gateway[API Gateway<br/>Hertz Framework<br/>Port :8080]

        subgraph GatewayMiddleware["网关中间件"]
            JWT[JWT 认证]
            CORS[CORS 跨域]
            RateLimit[限流]
            Protocol[协议转换<br/>HTTP → RPC]
            Tracing[日志追踪]
            ErrorHandler[统一错误处理]
        end
    end

    subgraph ServiceLayer["RPC 服务层"]
        subgraph IdentityService["Identity Service"]
            IdentitySrv[Identity Service<br/>Kitex Framework<br/>Port :8891]

            subgraph IdentityFeatures["服务功能"]
                UserMgmt[用户管理]
                Auth[身份认证]
                OrgMgmt[组织管理]
            end
        end

        FutureSrv1[Future Service 1<br/>Port :889X]
        FutureSrv2[Future Service N<br/>Port :889Y]
    end

    subgraph InfrastructureLayer["基础设施层"]
        PostgreSQL[(PostgreSQL 16<br/>Port :5432)]
        Etcd[(etcd<br/>Port :2379)]
        RustFS[(RustFS<br/>Port :9000)]
    end

    %% 客户端到网关
    WebClient -->|HTTP/HTTPS| Gateway
    MobileClient -->|HTTP/HTTPS| Gateway

    %% 网关内部中间件流程
    Gateway --> JWT --> CORS --> RateLimit --> Protocol --> Tracing --> ErrorHandler

    %% 网关到服务层
    ErrorHandler -->|RPC Thrift| Etcd
    Etcd -.->|服务注册| IdentitySrv
    ErrorHandler -->|RPC Thrift| IdentitySrv
    ErrorHandler -->|RPC Thrift| FutureSrv1
    ErrorHandler -->|RPC Thrift| FutureSrv2

    %% 服务到基础设施
    IdentitySrv -->|GORM| PostgreSQL
    IdentitySrv -->|服务注册| Etcd
    IdentitySrv -->|S3 API| RustFS
```

### 架构特点

- **星型拓扑**：网关为中心，RPC 服务不互相调用
- **协议转换**：HTTP → Thrift RPC
- **服务发现**：基于 etcd 的动态服务发现
- **统一入口**：所有外部请求通过网关

---

## 服务列表

| 服务名称 | 框架 | 端口 | 描述 |
|----------|------|------|------|
| **gateway** | Hertz | 8080 | HTTP 网关，统一 API 入口 |
| **identity_srv** | Kitex | 8891 | 身份认证服务 |

---

## RPC 服务分层架构

```mermaid
flowchart TB
    subgraph External["外部接口层"]
        RPCRequest[RPC 请求<br/>Thrift Protocol]
        RPCResponse[RPC 响应]
    end

    subgraph HandlerLayer["Handler Layer"]
        HandlerImpl[Handler 实现]
        subgraph HandlerOps["Handler 职责"]
            ParamValidation[参数校验]
            ErrorConvert[错误转换]
            ResponseBuild[响应构建]
        end
    end

    subgraph BusinessLayer["Business Layer"]
        subgraph LogicModule["Logic 模块"]
            UserLogic[UserLogic]
            OrgLogic[OrgLogic]
            AuthLogic[AuthLogic]
        end

        subgraph ConverterModule["Converter 模块"]
            UserConverter[UserConverter]
            OrgConverter[OrgConverter]
        end

        subgraph DALModule["DAL 模块"]
            UserDAL[UserDAL]
            OrgDAL[OrgDAL]
        end
    end

    subgraph DataLayer["Data Layer"]
        UserModel[User Model]
        OrgModel[Organization Model]
    end

    subgraph Infrastructure["基础设施"]
        GORM[(GORM ORM)]
        DB[(PostgreSQL)]
        Wire[Wire DI]
    end

    RPCRequest --> HandlerImpl
    HandlerImpl --> LogicModule
    LogicModule --> ConverterModule
    LogicModule --> DALModule
    DALModule --> DataLayer
    DataLayer --> GORM --> DB
    Wire -.->|注入| HandlerImpl
    Wire -.->|注入| LogicModule
    Wire -.->|注入| DALModule
    LogicModule --> RPCResponse
```

### 目录结构

```
rpc/<service_name>/
├── handler.go           # RPC 接口实现（适配层）
├── biz/                 # 核心业务逻辑层
│   ├── converter/       # DTO ↔ Model 转换
│   ├── dal/             # 数据访问层
│   └── logic/           # 业务逻辑实现
├── models/              # GORM 数据模型
├── kitex_gen/           # IDL 生成代码（勿修改）
├── config/              # 服务配置
├── wire/                # Wire 依赖注入
└── internal/            # 内部实现
    └── middleware/      # RPC 中间件
```

### 分层职责

| 层 | 职责 | 示例 |
|----|------|------|
| **Handler** | 参数校验、调用转换器、委托业务逻辑 | `handler.go` |
| **Logic** | 核心业务逻辑、编排 DAL 操作 | `biz/logic/user_profile/` |
| **DAL** | 数据持久化、封装 GORM 操作 | `biz/dal/user_profile/` |
| **Converter** | DTO 与 Model 纯函数转换 | `biz/converter/user_profile/` |

---

## HTTP 网关分层架构

```
gateway/
├── biz/                  # HTTP 业务层（IDL 生成）
│   ├── handler/          # HTTP Handler
│   ├── model/            # HTTP DTO
│   └── router/           # 路由注册
├── internal/
│   ├── application/      # 应用层
│   │   ├── assembler/    # 数据组装器
│   │   └── middleware/   # 中间件
│   ├── domain/           # 领域层
│   │   └── service/      # 领域服务
│   └── infrastructure/   # 基础设施层
│       ├── client/       # RPC 客户端
│       ├── config/       # 配置管理
│       └── errors/       # 统一错误处理
└── docs/                 # Swagger 文档
```

---

## 技术栈

### 核心框架

| 组件 | 技术 | 说明 |
|------|------|------|
| RPC 框架 | [Kitex](https://github.com/cloudwego/kitex) | CloudWeGo 高性能 RPC |
| HTTP 框架 | [Hertz](https://github.com/cloudwego/hertz) | CloudWeGo 高性能 HTTP |
| 接口协议 | Thrift | IDL 定义 |

### 基础设施

| 组件 | 技术 | 说明 |
|------|------|------|
| 数据库 | PostgreSQL 16 + GORM | 关系型数据库 |
| 服务发现 | etcd | 服务注册与发现 |
| 对象存储 | RustFS | S3 兼容存储 |
| 依赖注入 | Google Wire | 编译时 DI |

---

## 关键设计决策

### 1. RPC 服务不互相调用

**决策**：所有 RPC 调用必须由网关发起，RPC 服务之间不直接调用。

**原因**：
- 简化微服务拓扑（星型架构）
- 统一管理认证、追踪、降级
- 易于监控和故障排查
- 降低服务间耦合

### 2. 安全逻辑集中在网关

**决策**：所有鉴权逻辑（JWT、权限校验）在网关层处理，RPC 服务不处理权限。

**原因**：
- 单点控制安全策略
- RPC 服务保持简单
- 便于安全审计

### 3. 环境变量驱动配置

**决策**：不使用 YAML 配置文件，所有配置通过环境变量或 `.env` 文件提供。

**原因**：
- 便于容器化部署
- 符合 12-Factor App 原则
- 环境隔离更清晰

### 4. Wire 编译时依赖注入

**决策**：使用 Google Wire 进行编译时依赖注入。

**原因**：
- 类型安全，编译时检查
- 无运行时反射开销
- 依赖关系清晰可追踪

---

## 下一步

- [开发指南](development-guide.md) - 了解如何开发新功能
- [配置说明](configuration.md) - 详细的配置参考
