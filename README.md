
# CloudWeGo 微服务实践项目

<p align="center">
  <img src="https://github.com/cloudwego/hertz/raw/main/doc/logo.png" height="150" />
  &nbsp;&nbsp;&nbsp;&nbsp;
  <img src="https://github.com/cloudwego/kitex/raw/main/doc/logo.png" height="150" />
</p>
<p align="center">
  一个基于 <a href="https://www.cloudwego.io/">CloudWeGo</a> 生态 (Kitex + Hertz) 构建的，可用于生产环境的微服务架构 Demo。
<p>
<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go" alt="Go Version"></a>
  <a href="https://github.com/cloudwego/kitex"><img src="https://img.shields.io/badge/Kitex-latest-00ADD8?style=flat" alt="Kitex"></a>
  <a href="https://github.com/cloudwego/hertz"><img src="https://img.shields.io/badge/Hertz-latest-00ADD8?style=flat" alt="Hertz"></a>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License"></a>
</p>

## 概览

本项目旨在演示如何利用 CloudWeGo 的 **Kitex** (RPC) 和 **Hertz** (HTTP) 框架，构建一个以**网关为核心**的微服务体系。它提供了一个参考实现，展示了以下领域的最佳实践：

-   **API 网关**: 使用 Hertz 作为统一流量入口。
-   **RPC 微服务**: 使用 Kitex 构建高性能 RPC 服务。
-   **整洁架构**: 业务逻辑、数据处理与应用框架分离。
-   **依赖注入**: 使用 Google Wire 实现编译时依赖注入。
-   **认证与授权**: 集成 JWT (JSON Web Token) 和 Casbin 权限管理。
-   **可观测性**: 集成 OpenTelemetry 进行链路追踪。

## 核心功能

-   **API 网关 (Hertz)**
    -   统一的 API 入口
    -   JWT 中间件实现身份认证
    -   Casbin 中间件实现 RBAC 授权
    -   集成了 OpenTelemetry 的分布式链路追踪
    -   Swagger API 文档生成

-   **身份服务 (Kitex)**
    -   用户管理 (增删改查、登录、密码管理)
    -   组织架构管理
    -   角色与权限管理 (基于 Casbin)

## 架构设计

项目采用网关服务 + RPC 服务的经典微服务架构。

```mermaid
graph TD
    subgraph Client [客户端]
        A[用户/浏览器]
    end

    subgraph Gateway [API 网关 (Hertz)]
        direction LR
        B(JWT 认证) --> C(Casbin 授权) --> D(链路追踪) --> E(服务代理)
    end

    subgraph Services [RPC 服务 (Kitex)]
        F[身份服务]
    end

    subgraph Infrastructure [基础设施]
        G[PostgreSQL]
        H[etcd]
        I[Redis]
    end

    A -- HTTP/HTTPS --> Gateway
    Gateway -- RPC (Thrift) --> Services
    Services -- connects to --> G
    Services -- connects to --> I
    Gateway -- service discovery --> H
    Services -- service discovery & config --> H
```

### 关键决策

-   **星型拓扑**: 所有 RPC 调用均由 API 网关发起，服务之间原则上不直接相互调用，简化服务依赖。
-   **IDL-First**: 使用 Thrift 作为接口定义语言 (IDL)，确保前后端与服务间的接口规范。
-   **编译时依赖注入**: 借助 Google Wire 在编译期完成依赖注入，降低运行时风险。

## 技术栈

| 组件               | 技术选型                                                     |
| ------------------ | ------------------------------------------------------------ |
| RPC 框架           | [Kitex](https://github.com/cloudwego/kitex)                  |
| HTTP 框架          | [Hertz](https://github.com/cloudwego/hertz)                  |
| 接口定义           | Thrift                                                       |
| 数据库             | PostgreSQL + [GORM](https://gorm.io/)                        |
| 服务注册与发现     | etcd                                                         |
| 缓存               | Redis                                                        |
| 依赖注入           | [Google Wire](https://github.com/google/wire)                |
| 授权               | [Casbin](https://casbin.org/)                                |
| 可观测性           | [OpenTelemetry](https://opentelemetry.io/)                   |
| 部署               | Docker / Docker Compose                                      |

## 快速开始

### 环境要求

-   Go 1.24+
-   Docker 20.10+
-   Docker Compose 2.0+

### 使用 Docker 运行

```bash
# 1. 克隆仓库
git clone https://github.com/masonsxu/cloudwego-microservice-demo.git
cd cloudwego-microservice-demo

# 2. 启动所有服务
cd docker
./deploy.sh up

# 3. 验证服务是否启动成功
curl http://localhost:8080/ping
# 预期返回: {"message":"pong"}
```

### 访问入口

-   **API 网关**: `http://localhost:8080`
-   **Swagger 文档**: `http://localhost:8080/swagger/index.html`

## 本地开发

### 代码生成

项目采用 IDL-First 的开发模式，当 `idl` 目录下的 Thrift 文件变更后，需要重新生成代码。

```bash
# 生成 Kitex RPC 代码 (位于 rpc/identity_srv 目录)
./rpc/identity_srv/script/gen_kitex_code.sh

# 生成 Hertz HTTP 代码 (位于 gateway 目录)
./gateway/script/gen_hertz_code.sh

# 生成 Wire DI 代码
# (分别在 gateway/internal/wire 和 rpc/identity_srv/wire 目录执行)
cd gateway/internal/wire && wire
cd ../../../rpc/identity_srv/wire && wire
```

### 运行服务

```bash
# 1. 仅启动数据库等基础设施
cd docker
./deploy.sh up-base

# 2. 启动 identity_srv (RPC 服务)
cd rpc/identity_srv
sh build.sh && sh output/bootstrap.sh

# 3. 启动 gateway (API 网关)
cd gateway
sh build.sh && sh output/bootstrap.sh
```

## 项目结构

```
cloudwego-microservice-demo/
├── gateway/              # HTTP 网关 (Hertz)
│   ├── biz/              # 业务逻辑 (Handler, Model)
│   └── internal/         # 内部实现 (DI, 中间件等)
├── rpc/
│   └── identity_srv/     # 身份认证 RPC 服务 (Kitex)
│       ├── biz/          # 核心业务逻辑
│       ├── models/       # GORM 数据库模型
│       └── wire/         # 依赖注入配置
├── idl/                  # Thrift IDL 定义
├── docker/               # Docker 部署相关
└── docs/                 # 项目文档
```

## 贡献代码

欢迎任何形式的贡献！请在提交 Pull Request 前阅读 [CONTRIBUTING.md](./CONTRIBUTING.md)。

## 许可证

本项目基于 [MIT License](./LICENSE) 开源。
