# CloudWeGo 微服务实践项目

<p align="center">
  一个基于 <a href="https://www.cloudwego.io/">CloudWeGo</a> 生态 (Kitex + Hertz) 构建的生产级微服务架构 Demo
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go" alt="Go Version"></a>
  <a href="https://github.com/cloudwego/kitex"><img src="https://img.shields.io/badge/Kitex-latest-00ADD8?style=flat" alt="Kitex"></a>
  <a href="https://github.com/cloudwego/hertz"><img src="https://img.shields.io/badge/Hertz-latest-00ADD8?style=flat" alt="Hertz"></a>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License"></a>
</p>

## 概览

本项目演示如何使用 CloudWeGo 的 **Kitex** (RPC) 和 **Hertz** (HTTP) 框架，构建一个以**网关为核心**的微服务体系，展示以下最佳实践：

- **API 网关**: Hertz 作为统一流量入口
- **RPC 微服务**: Kitex 构建高性能 RPC 服务
- **权限管理**: Casbin RBAC 权限引擎
- **整洁架构**: 业务逻辑、数据处理与框架分离
- **依赖注入**: Google Wire 编译时依赖注入
- **用户认证**: JWT Token 认证
- **可观测性**: OpenTelemetry 链路追踪

## 架构设计

```mermaid
%%{init: {'theme': 'base', 'themeVariables': { 'primaryColor': '#2C3E50', 'primaryTextColor': '#ECF0F1', 'primaryBorderColor': '#BDC3C7', 'lineColor': '#BDC3C7', 'secondaryColor': '#34495E', 'tertiaryColor': '#34495E', 'mainBkg': '#34495E', 'nodeBorder': '#BDC3C7', 'clusterBkg': '#34495E', 'clusterBorder': '#BDC3C7', 'defaultLinkColor': '#BDC3C7', 'fontFamily': 'arial'}}}%%
graph TD
    %% 样式定义
    classDef base fill:#2C3E50,stroke:#BDC3C7,stroke-width:1px,color:#ECF0F1;
    classDef highlight fill:#F39C12,stroke:#ECF0F1,stroke-width:2px,color:#2C3E50,font-weight:bold;
    classDef sub fill:#34495E,stroke:#7F8C8D,stroke-width:1px,color:#BDC3C7,stroke-dasharray: 5 5;
    classDef infra fill:#7F8C8D,stroke:#ECF0F1,stroke-width:1px,color:#ECF0F1;

    %% 客户端层
    Client[("📱 客户端<br/>Web / Mobile / API Client")]:::base

    %% API 网关层
    subgraph Gateway_Layer [API 网关层 - Hertz :8080]
        direction TB
        Gateway[("🚪 API Gateway")]:::highlight
        
        subgraph Middleware [中间件链]
            direction LR
            MW_CORS(CORS):::sub
            MW_Trace(Trace):::sub
            MW_Log(AccessLog):::sub
            MW_JWT(JWT Auth):::sub
            MW_Casbin(Casbin RBAC):::sub
            MW_Error(Error Handle):::sub
            MW_Resp(Response):::sub
            
            MW_CORS --> MW_Trace --> MW_Log --> MW_JWT --> MW_Casbin --> MW_Error --> MW_Resp
        end
        
        subgraph GW_Components [分层架构]
            direction LR
            GW_Handler(Handler<br/>biz/handler):::sub
            GW_Service(Domain Service<br/>internal/domain):::sub
            GW_Assembler(Assembler<br/>DTO Convert):::sub
            GW_Client(RPC Client<br/>infrastructure):::sub
            
            GW_Handler --> GW_Service --> GW_Assembler --> GW_Client
        end
    end

    %% RPC 服务层
    subgraph RPC_Layer [RPC 服务层 - Kitex :8891]
        direction TB
        IdentitySRV[("🛡️ Identity Service")]:::highlight
        
        subgraph Modules [业务模块]
            direction LR
            Mod_User(User):::sub
            Mod_Org(Org):::sub
            Mod_Role(Role):::sub
            Mod_Menu(Menu):::sub
            Mod_Logo(Logo):::sub
        end
        
        subgraph RPC_Components [分层架构]
            direction LR
            RPC_Handler(Handler<br/>RPC Adaptor):::sub
            RPC_Logic(Logic<br/>Business):::sub
            RPC_DAL(DAL<br/>Data Access):::sub
            RPC_Model(Models<br/>GORM):::sub
            
            RPC_Handler --> RPC_Logic --> RPC_DAL --> RPC_Model
        end
    end

    %% 基础设施层
    subgraph Infra_Layer [基础设施层]
        direction LR
        DB[("🐘 PostgreSQL<br/>:5432")]:::infra
        Redis[("🔴 Redis<br/>:6379")]:::infra
        Etcd[("🏗️ etcd<br/>:2379")]:::infra
        S3[("📦 RustFS (S3)<br/>:9000")]:::infra
        Jaeger[("🔍 Jaeger<br/>:16686")]:::infra
    end

    %% 连接关系
    Client ==>|HTTP/JSON| Gateway
    Gateway --> Middleware
    Middleware --> GW_Components
    GW_Components ==>|Thrift RPC| IdentitySRV
    
    IdentitySRV --> Modules
    Modules --> RPC_Components
    
    RPC_Components --> DB
    RPC_Components --> Redis
    RPC_Components --> S3
    
    Gateway -.->|服务发现| Etcd
    IdentitySRV -.->|服务注册| Etcd
    
    Gateway -.->|Trace| Jaeger
    IdentitySRV -.->|Trace| Jaeger
```

**关键设计决策**：
- **星型拓扑**: 所有 RPC 调用由网关发起，服务间不直接调用
- **IDL-First**: Thrift 作为接口定义语言
- **编译时依赖注入**: Google Wire 完成依赖注入

## 技术栈

| 组件 | 技术 |
|------|------|
| RPC 框架 | [Kitex](https://github.com/cloudwego/kitex) |
| HTTP 框架 | [Hertz](https://github.com/cloudwego/hertz) |
| 接口定义 | Thrift |
| 数据库 | PostgreSQL + [GORM](https://gorm.io/) |
| 服务发现 | etcd |
| 缓存 | Redis |
| 权限引擎 | [Casbin](https://casbin.org/) |
| 依赖注入 | [Google Wire](https://github.com/google/wire) |
| 可观测性 | [OpenTelemetry](https://opentelemetry.io/) |

## 快速开始

### 环境要求

- Go 1.24+
- Docker 20.10+ / Podman 4.0+
- Docker Compose 2.0+ / podman-compose

### 使用 Docker 运行

```bash
# 1. 克隆仓库
git clone https://github.com/masonsxu/cloudwego-microservice-demo.git
cd cloudwego-microservice-demo

# 2. 启动基础设施（PostgreSQL、etcd、Redis、RustFS、Jaeger）
cd docker && podman-compose up -d

# 3. 启动 RPC 服务（新终端）
cd rpc/identity_srv && sh build.sh && sh output/bootstrap.sh

# 4. 启动网关服务（新终端）
cd gateway && sh build.sh && sh output/bootstrap.sh

# 5. 验证
curl http://localhost:8080/ping
# 返回: {"message":"pong"}
```

### 访问入口

- **API 网关**: http://localhost:8080
- **Swagger 文档**: http://localhost:8080/swagger/index.html
- **Jaeger 链路追踪**: http://localhost:16686

## 项目结构

```
cloudwego-microservice-demo/
├── gateway/              # HTTP 网关 (Hertz)
├── rpc/
│   └── identity_srv/     # 身份认证 RPC 服务 (Kitex)
├── idl/                  # Thrift IDL 定义
├── docker/               # Docker 部署配置
└── docs/                 # 项目文档
```

## 文档

详细文档请查看 [docs/](docs/README.md)：

- [快速开始](docs/01-快速入门/快速开始.md)
- [架构设计](docs/00-项目概览/架构设计.md)
- [开发指南](docs/02-开发规范/开发指南.md)
- [配置参考](docs/01-快速入门/配置参考.md)
- [部署指南](docs/03-部署运维/部署指南.md)
- [故障排查](docs/03-部署运维/故障排查.md)
- [权限管理](docs/04-权限管理/权限管理设计.md)
- [奢华摩羯座配色规范](docs/05-UI设计/配色规范.md)
- [测试指南](docs/02-开发规范/测试指南.md)

## 测试

### 运行测试

```bash
# 运行所有测试
go test ./... -v

# 生成测试覆盖率报告
./scripts/generate-coverage-report.sh

# 或手动生成
cd rpc/identity_srv && go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 测试覆盖率

**RPC 服务 (identity_srv)**

| 模块 | 覆盖率 | 状态 |
|------|--------|------|
| pkg/errno | 100.0% | ✅ |
| biz/logic/assignment | 100.0% | ✅ |
| biz/logic/definition | 99.1% | ✅ |
| biz/logic/department | 97.3% | ✅ |
| biz/logic/membership | 93.1% | ✅ |
| biz/parser | 92.1% | ✅ |
| internal/middleware | 89.1% | ✅ |
| biz/logic/user | 86.3% | ✅ |
| biz/converter | 84.0% | ✅ |
| pkg/password | 83.3% | ✅ |
| biz/logic/organization | 71.0% | ✅ |
| pkg/log | 41.4% | ⚠️ |
| biz/dal | 0.0% | ❌ |

**Gateway 服务**

| 模块 | 覆盖率 | 状态 |
|------|--------|------|
| middleware/casbin_middleware | 9.2% | ⚠️ |
| infrastructure/redis | 6.7% | ⚠️ |
| 其他包 | 0.0% | ❌ |

详细的测试指南请参考 [测试文档](docs/02-开发规范/测试指南.md)。

## 许可证

本项目基于 [MIT License](./LICENSE) 开源。
