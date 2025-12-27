# 架构设计

本文档介绍项目的整体架构设计和技术选型。

## 目录

- [微服务架构](#微服务架构)
- [服务列表](#服务列表)
- [API 网关架构](#api-网关架构)
- [RPC 服务架构](#rpc-服务架构)
- [请求追踪链](#请求追踪链)
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
            Casbin[Casbin 授权]
            Tracing[链路追踪]
            ErrorHandler[错误处理]
        end
    end

    subgraph ServiceLayer["RPC 服务层"]
        IdentitySrv[Identity Service<br/>Kitex Framework<br/>Port :8891]
    end

    subgraph InfrastructureLayer["基础设施层"]
        PostgreSQL[(PostgreSQL<br/>Port :5432)]
        Redis[(Redis<br/>Port :6379)]
        Etcd[(etcd<br/>Port :2379)]
        RustFS[(RustFS<br/>Port :9000)]
        Jaeger[(Jaeger<br/>Port :16686)]
    end

    WebClient -->|HTTP/HTTPS| Gateway
    MobileClient -->|HTTP/HTTPS| Gateway

    Gateway --> JWT --> Casbin --> Tracing --> ErrorHandler
    ErrorHandler -->|RPC Thrift| IdentitySrv

    IdentitySrv -->|GORM| PostgreSQL
    IdentitySrv -->|S3 API| RustFS
    Gateway -->|缓存| Redis
    Gateway & IdentitySrv -->|服务注册| Etcd
    Gateway & IdentitySrv -.->|链路追踪| Jaeger
```

### 架构特点

| 特点 | 说明 |
|------|------|
| **星型拓扑** | 网关为中心，RPC 服务不互相调用 |
| **协议转换** | HTTP → Thrift RPC |
| **服务发现** | 基于 etcd 的动态服务发现 |
| **统一入口** | 所有外部请求通过网关 |
| **安全分层** | 认证授权集中在网关层 |

---

## 服务列表

| 服务 | 框架 | 端口 | 描述 |
|------|------|------|------|
| **gateway** | Hertz | 8080 | HTTP 网关，统一 API 入口 |
| **identity_srv** | Kitex | 8891 | 身份认证 RPC 服务 |

---

## API 网关架构

```mermaid
flowchart TB
    subgraph Gateway["API 网关 (Hertz)"]
        subgraph BizLayer["业务层 biz/"]
            Handler[HTTP Handler]
            Router[路由注册]
        end

        subgraph AppLayer["应用层 internal/application/"]
            Assembler[数据组装器]
            Middleware[中间件]
            Context[上下文管理]
        end

        subgraph DomainLayer["领域层 internal/domain/"]
            Service[领域服务]
        end

        subgraph InfraLayer["基础设施层 internal/infrastructure/"]
            Client[RPC 客户端]
            Config[配置管理]
            Otel[OpenTelemetry]
            RedisClient[Redis 客户端]
        end

        subgraph WireLayer["依赖注入 internal/wire/"]
            Wire[Google Wire]
        end
    end

    Handler --> Assembler
    Assembler --> Service
    Service --> Client
    Wire -.->|注入| Handler
    Wire -.->|注入| Service
    Wire -.->|注入| Client
```

### 网关目录结构

```
gateway/
├── biz/                      # 业务层（IDL 生成）
│   ├── handler/              # HTTP Handler
│   ├── model/                # HTTP DTO
│   └── router/               # 路由注册
├── internal/
│   ├── application/          # 应用层
│   │   ├── assembler/        # 数据组装器
│   │   ├── middleware/       # 中间件（JWT、Casbin、Trace）
│   │   └── context/          # 上下文管理
│   ├── domain/               # 领域层
│   │   └── service/          # 领域服务
│   ├── infrastructure/       # 基础设施层
│   │   ├── client/           # RPC 客户端
│   │   ├── config/           # 配置管理
│   │   ├── otel/             # OpenTelemetry
│   │   └── redis/            # Redis 客户端
│   └── wire/                 # 依赖注入
├── docs/                     # Swagger 文档
└── main.go                   # 入口
```

---

## RPC 服务架构

```mermaid
flowchart TB
    subgraph RPC["RPC 服务 (Kitex)"]
        subgraph HandlerLayer["Handler 层"]
            RPCHandler[RPC Handler<br/>handler.go]
        end

        subgraph LogicLayer["Logic 层 biz/logic/"]
            UserLogic[UserLogic]
            OrgLogic[OrgLogic]
            AuthLogic[AuthLogic]
        end

        subgraph ConverterLayer["Converter 层 biz/converter/"]
            Converter[DTO ↔ Model 转换]
        end

        subgraph DALLayer["DAL 层 biz/dal/"]
            UserDAL[UserDAL]
            OrgDAL[OrgDAL]
        end

        subgraph ModelLayer["Model 层 models/"]
            Models[GORM 数据模型]
        end
    end

    subgraph DB["数据库"]
        PostgreSQL[(PostgreSQL)]
    end

    RPCHandler --> LogicLayer
    LogicLayer --> ConverterLayer
    LogicLayer --> DALLayer
    DALLayer --> Models
    Models --> PostgreSQL
```

### RPC 服务目录结构

```
rpc/identity_srv/
├── handler.go                # RPC 接口实现
├── biz/                      # 业务逻辑层
│   ├── logic/                # 业务逻辑
│   │   ├── user/
│   │   ├── organization/
│   │   └── authentication/
│   ├── dal/                  # 数据访问层
│   │   ├── user/
│   │   └── organization/
│   └── converter/            # DTO ↔ Model 转换
│       ├── user/
│       └── organization/
├── models/                   # GORM 数据模型
├── kitex_gen/                # IDL 生成代码（勿修改）
├── config/                   # 配置管理
├── wire/                     # 依赖注入
└── internal/
    └── middleware/           # RPC 中间件
```

### 分层职责

| 层 | 职责 | 位置 |
|----|------|------|
| **Handler** | 参数校验、错误转换、响应构建 | `handler.go` |
| **Logic** | 核心业务逻辑、编排 DAL 操作 | `biz/logic/` |
| **DAL** | 数据持久化、封装 GORM 操作 | `biz/dal/` |
| **Converter** | DTO 与 Model 纯函数转换 | `biz/converter/` |

### 数据流向

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant Gateway as API 网关
    participant Handler as RPC Handler
    participant Logic as Logic 层
    participant DAL as DAL 层
    participant DB as PostgreSQL

    Client->>Gateway: HTTP 请求
    Gateway->>Gateway: JWT 认证
    Gateway->>Gateway: Casbin 授权
    Gateway->>Handler: RPC 调用
    Handler->>Handler: 参数校验
    Handler->>Logic: 调用业务逻辑
    Logic->>DAL: 数据操作
    DAL->>DB: SQL 查询
    DB-->>DAL: 返回数据
    DAL-->>Logic: 返回 Model
    Logic-->>Handler: 返回结果
    Handler-->>Gateway: RPC 响应
    Gateway-->>Client: HTTP 响应
```

---

## 请求追踪链

项目采用基于 **metainfo** 的链路追踪机制，使用 `request_id` 追踪请求在微服务调用链中的完整路径。

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant Gateway as API 网关
    participant RPC as RPC 服务

    Client->>Gateway: HTTP 请求
    Note over Gateway: 1. requestid 中间件<br/>生成 request_id
    Note over Gateway: 2. trace_middleware<br/>注入到 metainfo
    Gateway->>RPC: RPC 调用<br/>(TTHeader 传递 request_id)
    Note over RPC: 3. meta_middleware<br/>提取 request_id
    Note over RPC: 4. 业务逻辑处理<br/>GetRequestID(ctx)
    RPC-->>Gateway: RPC 响应
    Gateway-->>Client: HTTP 响应
```

### 追踪特性

- **唯一追踪标识**：使用 `request_id` 进行全链路追踪
- **自动生成**：缺失的 request_id 自动生成
- **直接使用 metainfo**：不使用 context.WithValue，性能更优
- **100% 可追踪**：确保每个请求都有完整的追踪信息

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
| 数据库 | PostgreSQL + GORM | 关系型数据库 |
| 缓存 | Redis | 会话管理、热点数据缓存 |
| 服务发现 | etcd | 服务注册与发现 |
| 对象存储 | RustFS | S3 兼容存储 |
| 链路追踪 | Jaeger + OpenTelemetry | 分布式链路追踪 |
| 依赖注入 | Google Wire | 编译时依赖注入 |

---

## 关键设计决策

### 1. 星型拓扑架构

**决策**：所有 RPC 调用由网关发起，RPC 服务之间不直接调用。

```mermaid
flowchart LR
    Gateway((API 网关))
    S1[服务 A]
    S2[服务 B]
    S3[服务 N]

    Gateway --> S1
    Gateway --> S2
    Gateway --> S3
```

**优势**：
- 简化微服务拓扑
- 统一管理认证、追踪、降级
- 易于监控和故障排查
- 降低服务间耦合

### 2. 安全分层设计

**决策**：所有鉴权逻辑（JWT、权限校验）在网关层处理，RPC 服务不处理权限。

**优势**：
- 单点控制安全策略
- RPC 服务保持简单
- 便于安全审计

### 3. 环境变量驱动配置

**决策**：不使用 YAML 配置文件，所有配置通过环境变量或 `.env` 文件提供。

**优势**：
- 便于容器化部署
- 符合 12-Factor App 原则
- 环境隔离更清晰

### 4. Wire 编译时依赖注入

**决策**：使用 Google Wire 进行编译时依赖注入。

**优势**：
- 类型安全，编译时检查
- 无运行时反射开销
- 依赖关系清晰可追踪

### 5. Redis 缓存层

**决策**：在网关层引入 Redis 缓存，用于会话管理和热点数据缓存。

**用途**：
- JWT Token 黑名单
- 用户会话状态
- 权限数据缓存
- API 限流计数

### 6. OpenTelemetry 链路追踪

**决策**：集成 OpenTelemetry 和 Jaeger 实现端到端链路追踪。

**优势**：
- 分布式系统调用链可视化
- 性能瓶颈快速定位
- 错误根因分析

---

## 下一步

- [03-开发指南](03-development.md) - 了解如何开发新功能
- [04-配置参考](04-configuration.md) - 详细的配置参考
