# AGENTS.md

AI 辅助开发规范文件，为 AI 编程助手提供项目上下文和开发约束。

## 项目概述

基于 Go 1.24+ 的微服务项目，采用 CloudWeGo 技术栈和 Go Workspace 管理两个独立模块：

- `gateway/`：HTTP 网关服务（Hertz，端口 8080）- 鉴权、权限检查、HTTP→RPC 协议转换
- `rpc/identity_srv/`：身份认证 RPC 服务（Kitex，端口 8891）- 用户、组织、角色、权限、菜单管理

## 常用命令

### 基础设施和服务启动

```bash
# 启动基础设施（PostgreSQL、etcd、Redis、RustFS、Jaeger）
cd docker && podman-compose up -d

# 查看基础设施状态
cd docker && podman-compose ps

# 停止基础设施
cd docker && podman-compose down

# 查看日志
cd docker && podman-compose logs -f          # 全部日志
cd docker && podman-compose logs -f postgres  # 特定服务

# 清理环境（删除容器和数据卷）
cd docker && podman-compose down -v --remove-orphans

# 启动 RPC 服务
cd rpc/identity_srv && sh build.sh && sh output/bootstrap.sh

# 启动网关服务
cd gateway && sh build.sh && sh output/bootstrap.sh
```

### 测试

```bash
# 运行全部测试
go test ./... -v

# 运行单个模块的测试
cd rpc/identity_srv && go test ./... -v
cd gateway && go test ./... -v

# 运行单个测试函数
go test ./biz/logic/user -run TestUserLogic_CreateUser -v -count=1

# 测试覆盖率
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total
```

## AI 行为约束

### 🚫 禁止自动启动服务

**规则**：AI **严禁**自动启动前端开发服务器（`npm run dev`）或任何长期运行的服务。

**原因**：
- 开发服务器会占用端口（5173、5174、5175等），多次启动会导致端口冲突
- 用户可能已经在运行其他服务的实例
- 启动服务是用户的主动操作，不应由AI自动执行

**正确做法**：
- 当需要测试前端页面时，AI应该**提示用户**手动在终端执行启动命令
- AI只负责代码生成、文件编辑、配置修改等工作
- 示例提示："前端页面已实现完成，请在终端手动运行 `cd web && npm run dev` 启动开发服务器"

**适用范围**：
- ✅ 可以执行：代码生成、编译测试、lint检查、数据库迁移等一次性任务
- ❌ 禁止执行：`npm run dev`、`sh build.sh && sh output/bootstrap.sh` 等长期运行的服务

### 测试

```bash
# 自动格式化（gofumpt + golines + gci）
golangci-lint format

# 代码检查（Lint）
golangci-lint run

# 修复导入排序（独立使用）
cd rpc/identity_srv && gci write .
cd gateway && gci write .
```

### 📋 提交前检查规范（避免 CI 失败）

**问题**：CI 中 `golangci-lint run` 只检查不修复，导致 push 后 lint 失败。

**解决方案**：每次 commit 前必须执行以下步骤：

```bash
# 1. 确保本地 golangci-lint 版本与 CI 一致（v2.4.0）
golangci-lint version

# 2. 自动修复格式问题
golangci-lint format

# 3. 验证是否还有问题（需手动修复 unused 代码等）
golangci-lint run

# 4. 确保测试通过
go test ./... -v
```

**安装 pre-commit hook**（自动检查）：
```bash
ln -sf ../../scripts/git-hooks/pre-commit .git/hooks/pre-commit
```

**常见 CI 失败原因**：
- `gci` 导入顺序问题 → `golangci-lint format` 自动修复
- `wsl_v5` 空行问题 → `golangci-lint format` 自动修复
- `unused` 函数/变量 → 需手动删除代码
- 本地版本与 CI 不一致 → 重新安装 v2.4.0

### 代码生成

```bash
# Kitex RPC 代码生成（修改 IDL 后必须运行）
cd rpc/identity_srv && ./script/gen_kitex_code.sh

# Hertz HTTP 代码生成（修改 IDL 后必须运行）
cd gateway && ./script/gen_hertz_code.sh

# Wire 依赖注入生成（修改 provider 后必须运行）
cd rpc/identity_srv/wire && wire
cd gateway/internal/wire && wire
```

## 架构

### 代码分层（identity_srv）

```
handler.go               # RPC 接口实现（适配层，生成代码之上）
biz/
├── logic/               # 业务逻辑（logic.go 聚合所有子领域接口）
├── dal/                 # 数据访问（dal.go 聚合所有 Repository 接口）
└── converter/           # 纯函数 DTO ↔ Model 转换
models/                  # GORM 数据模型
wire/                    # Wire 依赖注入
config/                  # 配置（env vars → defaults.go）
```

**DAL 接口设计**：`biz/dal/dal.go` 是单一聚合接口，通过方法返回各子 Repository（如 `UserProfile()`、`Organization()`），同时提供事务管理（`WithTransaction`、`BeginTx`、`Commit`、`Rollback`）。

**Logic 接口设计**：`biz/logic/logic.go` 是单一聚合接口，通过嵌入各子领域接口（如 `user.ProfileLogic`、`organization.OrganizationLogic`）组成，对应 handler 统一注入一个 `Logic` 实例。

### 代码分层（gateway）

```
biz/handler/             # HTTP Handler（Hertz 生成，只做参数绑定和响应）
internal/
├── application/         # 应用层
│   ├── assembler/       # RPC DTO → HTTP DTO 转换（每个领域独立）
│   └── middleware/      # 中间件实现
├── domain/service/      # 领域服务（编排 RPC 调用，实现业务流程）
├── infrastructure/      # 基础设施
│   ├── client/          # Kitex RPC 客户端
│   ├── redis/           # Redis 缓存（JWT Token、Casbin 策略）
│   └── otel/            # OpenTelemetry 配置
└── wire/                # 依赖注入
```

**中间件执行顺序**：OpenTelemetry → RequestID → ResponseHeader → Trace → CORS → ErrorHandler → JWT → Casbin → ETag

### IDL-First 开发流程

1. 修改 `idl/` 下的 Thrift 文件
2. 运行代码生成脚本（Kitex 或 Hertz）
3. 在 `biz/` 中实现业务逻辑
4. 更新 `wire/provider.go`，运行 `wire`

**禁止手动修改** `kitex_gen/` 目录下的任何文件。

### 错误处理

错误码格式为 6 位数字 `A-BB-CCC`：
- `A`：1=系统级，2=业务级
- `BB`：00=通用，01=用户，02=组织，03=部门，04=级联删除，05=姿态资源，06=Logo，07=角色
- `CCC`：具体错误编码

**分层错误转换流程**：
1. DAL 层：GORM 错误 → `errno.ErrNo`（使用 `errno.WrapDatabaseError`）
2. Logic 层：返回 `errno.ErrNo`
3. Handler 层：`errno.ToKitexError()` → Kitex `BizStatusError`
4. 网关客户端：`kerrors.FromBizStatusError()` 解析

检查记录不存在：使用 `errno.IsRecordNotFound(err)`，不要直接比较 `gorm.ErrRecordNotFound`。

### 配置管理

**不使用 YAML 配置文件**，仅通过环境变量或 `.env` 文件配置。

优先级：系统环境变量 > `.env` 文件 > `config/defaults.go`

Logo 存储有两个端点：
- `LOGO_STORAGE_S3_ENDPOINT`：内部端点（容器间通信，上传用）
- `LOGO_STORAGE_S3_PUBLIC_ENDPOINT`：公共端点（生成预签名 URL 用）

### Wire 依赖注入规范

- 每个 Provider 函数只创建一个依赖，职责单一
- 返回类型必须是接口，不暴露具体实现
- 需要清理的资源（如数据库连接）返回 `cleanup func()`

### 代码风格

**导入顺序**（gci 自动排序）：
1. 标准库
2. 第三方库
3. 项目内部包（前缀 `github.com/masonsxu/cloudwego-microservice-demo/`）

**命名**：接口无 `I` 前缀（如 `DAL`、`Logic`），实现加 `Impl` 后缀（如 `DALImpl`、`LogicImpl`）。

**最大行长度**：120 字符（golines 自动处理）。

**测试命名**：`Test<InterfaceName>_<MethodName>`（如 `TestUserDALImpl_CreateUser_NotFound`）。

## 关键文件速查

| 文件 | 用途 |
|------|------|
| `go.work` | Go Workspace 配置（两个模块） |
| `.golangci.yml` | Lint + 格式化配置（行长 120，导入顺序） |
| `scripts/git-hooks/pre-commit` | 安装：`ln -sf ../../scripts/git-hooks/pre-commit .git/hooks/pre-commit` |
| `docker/docker-compose.yml` | 基础设施容器编排（podman-compose） |
| `idl/` | Thrift IDL 接口定义（修改此处触发代码生成） |
| `rpc/identity_srv/pkg/errno/` | 错误码定义和转换工具 |
| `gateway/config/casbin_model.conf` | Casbin RBAC 模型配置 |

## 文档索引

> **规则**：遇到具体场景时，先根据下表读取对应文档，不确定时优先读文档再动手。

| 文档 | 适用场景 |
|------|----------|
| `docs/01-快速入门/快速开始.md` | 首次搭建环境、启动服务、验证基础设施 |
| `docs/00-项目概览/架构设计.md` | 理解整体架构、服务间关系、数据模型设计、技术选型决策 |
| `docs/02-开发规范/开发指南.md` | IDL-First 开发流程、代码生成、分层实现规范、Wire 依赖注入、错误处理 |
| `docs/01-快速入门/配置参考.md` | 环境变量配置、数据库/Redis/JWT/S3 参数、多环境配置差异 |
| `docs/03-部署运维/部署指南.md` | Docker 部署、生产环境配置、服务端口、部署检查清单 |
| `docs/03-部署运维/故障排查.md` | 服务启动失败、代码生成报错、运行时异常、Docker 问题、性能排查 |
| `docs/04-权限管理/权限管理设计.md` | Casbin RBAC 策略、角色权限设计、API 权限配置、策略同步机制 |
| `docs/05-UI设计/配色规范.md` | UI 配色规范、CSS 变量定义、设计理念 |
| `docs/02-开发规范/测试指南.md` | 分层测试策略、覆盖率目标、测试编写规范、测试运行方式 |
| `docs/03-部署运维/CI-CD指南.md` | CI/CD 流水线配置、自动化检查项、Workflow 调试 |

## 可用 Skills

通过 `/command` 触发，封装高频多步骤开发操作。

| Command | 用途 | 典型用法 |
|---------|------|----------|
| `/add-rpc-method` | 端到端添加新 RPC 接口（IDL→全层实现） | `/add-rpc-method CreateDepartment` |
| `/add-http-endpoint` | 添加新 HTTP 端点（IDL→网关层） | `/add-http-endpoint POST /api/v1/departments` |
| `/add-domain` | 创建新业务领域模块骨架 | `/add-domain notification` |
| `/write-tests` | 按项目规范编写测试 | `/write-tests rpc/identity_srv/biz/converter/department/` |
| `/codegen` | 执行代码生成（Kitex/Hertz/Wire） | `/codegen kitex` 或 `/codegen`（自动检测） |
| `/diagnose` | 根据错误信息进行故障诊断 | `/diagnose` 后粘贴错误日志 |
| `/check-quality` | 全套代码质量检查 | `/check-quality` |
