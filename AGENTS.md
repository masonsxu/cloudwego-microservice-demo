# AGENTS.md

此文件供在代码库中工作的智能编码代理使用。

## 快速导航

- [项目概述](#项目概述) - 服务清单和技术栈
- [开发环境检查](#开发环境检查) - 快速验证环境
- [常用命令速查](#常用命令速查) - 常见开发命令
- [代码风格指南](#代码风格指南) - 命名约定、导入顺序、格式化
- [UI配色设计规范](#UI配色设计规范) - 官方配色方案和使用指南
- [架构约束](#架构约束) - 关键设计原则和分层职责
- [开发规范](#开发规范) - IDL-First 流程、Wire、配置
- [测试指南](#测试指南) - 各层测试编写规范和示例
- [常见开发流程](#常见开发流程) - 端到端开发步骤
- [故障排查](#故障排查) - 常见问题和解决方案
- [关键约定清单](#关键约定清单) - 必须遵守的约定
- [详细文档](#详细文档) - 完整开发指南

## 项目概述

cloudwego-microservice-demo 是一个基于 Go 语言的微服务项目，采用 CloudWeGo 技术栈，使用 Kitex (RPC) 和 Hertz (HTTP) 框架。项目采用 Go Workspace 管理多个服务模块，遵循 IDL-First 开发模式。

**服务列表**：
- `gateway`: HTTP 网关服务（Hertz，端口 8080）- 处理客户端请求、认证、权限、RPC 协议转换
- `identity_srv`: 身份认证 RPC 服务（Kitex，端口 8891）- 用户、认证、角色、权限、菜单管理

## 核心技术栈

- **Go**: 1.24+
- **RPC 框架**: Kitex (CloudWeGo)
- **HTTP 框架**: Hertz (CloudWeGo)
- **接口协议**: Thrift (IDL-First)
- **数据库**: PostgreSQL + GORM
- **依赖注入**: Google Wire (编译时)
- **日志**: zerolog
- **代码检查**: golangci-lint
- **权限引擎**: Casbin (RBAC)
- **可观测性**: OpenTelemetry + Jaeger
- **服务发现**: etcd
- **缓存**: Redis
- **文件存储**: RustFS (S3 兼容)

## 开发环境检查

快速验证开发环境是否完整配置：

```bash
# 检查 Go 版本（需要 1.24+）
go version

# 检查必要工具
which kitex && kitex version
which hertz && hertz version
which wire
which golangci-lint && golangci-lint version
which gci && gci version
which gofumpt && gofumpt -version
which golines && golines --version

# 检查 Docker 和 Docker Compose
docker --version && docker-compose --version

# 检查 Go Workspace 配置
cat go.work

# 检查两个模块都存在
ls -d gateway rpc/identity_srv
```

**期望输出**：
- Go 1.24 或更高
- Kitex、Hertz、Wire、golangci-lint 都可用
- Docker 20.10+, Docker Compose 2.0+
- `go.work` 存在且包含两个模块
- 两个模块目录都存在

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
cd docker
./deploy.sh down                      # 停止基础设施
./deploy.sh ps                        # 查看运行状态
./deploy.sh logs                      # 查看实时日志
./deploy.sh logs postgres             # 查看特定服务日志

# 访问入口
# API 网关: http://localhost:8080
# Swagger 文档: http://localhost:8080/swagger/index.html
# Jaeger 链路追踪: http://localhost:16686
```

### 构建和运行

```bash
# RPC 服务 - 完整流程
cd rpc/identity_srv
sh build.sh                           # 编译
sh output/bootstrap.sh                # 运行（默认端口 8891）
# 或者 sh output/bootstrap.sh -addr=:8891

# 网关服务 - 完整流程
cd gateway
sh build.sh                           # 编译
sh output/bootstrap.sh                # 运行（默认端口 8080）
# 或者 sh output/bootstrap.sh -addr=:8080
```

### 测试

```bash
# 运行所有测试（从项目根目录，测试所有模块）
go test ./... -v

# 运行特定模块的测试
cd rpc/identity_srv && go test ./... -v
cd gateway && go test ./... -v

# 运行特定包的测试
cd rpc/identity_srv && go test ./biz/logic/user -v

# 运行单个测试函数
cd rpc/identity_srv && go test ./biz/converter/permission -run TestConverterImpl_ModelToThrift -v

# 运行带日志输出的测试
cd rpc/identity_srv && go test ./biz/logic/user -v -count=1 2>&1 | grep -A 10 "TestUserLogic_CreateUser"

# 测试覆盖率
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out          # 生成 HTML 覆盖率报告
go tool cover -func=coverage.out | grep total  # 查看覆盖率总和

# 性能基准测试
cd rpc/identity_srv && go test -bench=. -benchmem ./biz/logic/user
```

### 代码检查和格式化

```bash
# 检查所有问题
golangci-lint run

# 检查特定目录
golangci-lint run ./gateway/...
golangci-lint run ./rpc/identity_srv/...

# 查看详细信息
golangci-lint run --verbose

# 自动格式化代码（gofumpt, golines, goimports, gci）
golangci-lint format

# 仅显示格式问题（不修改）
golangci-lint run --out-format=colored-line-number

# 修复导入顺序
cd rpc/identity_srv && gci write .
cd gateway && gci write .
```

### 代码生成

```bash
# Kitex RPC 代码生成（从 identity_srv 目录）
cd rpc/identity_srv && ./script/gen_kitex_code.sh

# Hertz HTTP 代码生成（从 gateway 目录）
cd gateway && ./script/gen_hertz_code.sh

# Wire 依赖注入生成（生成 wire_gen.go）
cd rpc/identity_srv/wire && wire
cd gateway/internal/wire && wire

# 注意：生成后需要重新构建
cd rpc/identity_srv && go mod tidy && sh build.sh
cd gateway && go mod tidy && sh build.sh
```

### 依赖管理

```bash
# 清理和整理依赖
go mod tidy

# 下载依赖
go mod download

# 查看依赖树
go mod graph

# 更新依赖
go get -u ./...
```

## 代码风格指南

### 导入顺序
1. 标准库（standard）
2. 第三方库（default）
3. 项目内部包（prefix: github.com/masonsxu/cloudwego-microservice-demo/）

使用 `gci` 格式化器自动排序。

### 命名约定
- **接口**：使用描述性名称，无 `I` 前缀（如 `DAL`, `Converter`）
- **实现**：接口名 + `Impl` 后缀（如 `DALImpl`, `ConverterImpl`）
- **测试文件**：`<name>_test.go`
- **常量**：大驼峰或大写下划线分隔（如 `ErrorCodeUserNotFound`）
- **错误变量**：`Err` + 描述（如 `ErrUserNotFound`）

### 错误处理规范
1. **错误码格式**：6位数字 A-BB-CCC
   - A：错误级别（1=系统级，2=业务级）
   - BB：服务/模块编码（00=通用，01-07=各领域）
   - CCC：具体错误编码
2. **分层错误处理**：
   - DAL层：转换 GORM 错误 → ErrNo
   - Logic层：处理业务逻辑，返回 ErrNo
   - Handler层：使用 `errno.ToKitexError()` → BizStatusError
   - 客户端：使用 `kerrors.FromBizStatusError()` 解析
3. **检查错误**：使用 `errno.IsRecordNotFound()` 判断记录不存在

### 函数和结构体
- **Provider函数**：每个函数只创建一个依赖，职责单一
- **接口方法**：返回接口类型，具体实现私有
- **构造函数**：使用 `New` 前缀，返回接口
- **注释**：导出的类型、函数、方法必须有注释

### 代码格式化
- **最大行长度**：120字符
- **缩进**：使用 tab
- **格式化工具**：`gofumpt`, `golines`, `goimports`

## UI配色设计规范

### 设计理念

**奢华摩羯座配色方案** (Luxury Capricorn Color Scheme)

设计核心：稳重、优雅、追求卓越 | Steady, Elegant, Excellence

### 快速参考

所有 UI 设计、前端开发、文档编写必须严格遵循 **[奢华摩羯座配色规范](docs/CAPRICORN-THEME-GUIDE.md)**。

#### 核心配色速览

| 颜色 | 名称 | 用途 | Tailwind 类名 |
|------|------|------|--------------|
| #141416 | 深岩灰 | 背景基础 | `bg-bg-base` |
| #D4AF37 | 香槟金 | 核心高亮/图标 | `text-gold`, `border-gold` |
| #F2F0E4 | 羊皮纸白 | 主标题/文本 | `text-vellum` |
| #8B9bb4 | 矿石灰 | 副文本 | `text-mineral` |
| #2C2E33 | 青铜褐 | 按钮/卡片背景 | `bg-bronze` |
| #FFF8E7 | 亮光色 | 极亮部 | `text-highlight` |

#### CSS 变量定义

```css
:root {
    --bg-base: #141416;              /* 深岩灰 - 背景基础 */
    --c-accent: #D4AF37;             /* 香槟金 - 核心高亮/图标 */
    --c-text-main: #F2F0E4;         /* 羊皮纸白 - 主标题 */
    --c-text-sub: #8B9bb4;          /* 矿石灰 - 副文本 */
    --c-btn-bg: #2C2E33;            /* 青铜褐 - 按钮底色 */
    --c-highlight: #FFF8E7;          /* 亮光色 - 极亮部 */
}
```

#### 常用组件示例

```html
<!-- 按钮 -->
<button class="bg-bronze text-vellum border border-gold px-7 py-3 hover:bg-gold hover:text-black transition-all">
  主按钮
</button>

<!-- 卡片 -->
<div class="bg-gradient-to-br from-[rgba(30,32,36,0.9)] to-[rgba(20,20,22,0.95)] border border-white/5 rounded-2xl hover:border-gold/30">
  内容
</div>

<!-- 标签 -->
<span class="bg-bronze/60 text-gold border border-gold px-4 py-2 rounded-full">
  ♑ 摩羯座
</span>

<!-- 标题 -->
<h1 class="font-cinzel text-4xl text-gold bg-gold-gradient bg-clip-text text-transparent animate-shine">
  摩羯座
</h1>
```

### 完整规范文档

详细的配色规范、组件示例、动画效果、测试指南等内容，请参考：

- **[奢华摩羯座配色规范](docs/CAPRICORN-THEME-GUIDE.md)** - 完整的设计文档

该文档包含：
- 官方配色定义（CSS 和 Tailwind 配置）
- 配色详解（每种颜色的设计意图和使用场景）
- 字体系统
- UI 组件配色规范（按钮、卡片、标签、输入框、进度条、提示框）
- 动画效果规范
- 图标颜色规范
- 配色使用原则
- 实施检查清单
- 常见配色场景
- 配色测试指南

---

## 架构约束

### 关键架构原则
1. **星型拓扑**：所有RPC调用由网关发起，服务间不直接调用
2. **安全分层**：鉴权逻辑必须在网关层，RPC服务不处理权限
3. **协议转换**：网关负责 HTTP → Thrift RPC 转换
4. **无状态设计**：RPC服务设计为无状态，支持水平扩展
5. **IDL-First**：Thrift接口保持向后兼容，不手动修改 `kitex_gen/` 目录

### RPC 服务分层职责

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
└── internal/
    ├── application/ # 应用层（assembler, middleware）
    ├── domain/      # 领域层（service）
    ├── infrastructure/ # 基础设施层（client, config）
    └── wire/        # 依赖注入
```

## 开发规范

### IDL-First 开发流程

1. 修改 `idl/` 目录下的 Thrift 文件
2. 使用 Kitex/Hertz 工具生成代码
3. 在 `biz/` 目录实现业务逻辑
4. 更新 Wire 依赖注入
5. 编写单元测试

### Wire 依赖注入规范

- Provider 职责单一，每个函数只创建一个依赖
- 数据库初始化在 `config` 层，Wire 层只调用
- 需要清理的资源返回 cleanup 函数

### 配置管理
- 不使用 YAML 配置文件，通过环境变量或 `.env` 文件配置
- 优先级：系统环境变量 > `.env` 文件 > `config/defaults.go`
- Duration 格式支持：`1h`, `30m`, `3600s` 或纯数字秒

### 文件存储（RustFS）

采用双端点配置解决容器化部署问题：
- `LOGO_STORAGE_S3_ENDPOINT`: 内部端点（容器间通信）
- `LOGO_STORAGE_S3_PUBLIC_ENDPOINT`: 公共端点（预签名 URL）

## Git 规范

### Git 提交消息格式

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

### Shell 脚本
- 使用 `set -e` 在错误时立即退出
- 使用 `set -ou pipefail` 进行严格错误检查
- 提供清晰的输出信息（INFO, WARN, ERROR）

### 注释和文档
- 导出的代码元素必须有注释
- 使用分隔注释区分代码区域：
  ```go
  // ============================================================================
  // UserProfile
  // ============================================================================
  ```
- 代码生成文件包含 `// Code generated by ... DO NOT EDIT.` 标记

## 测试指南

### 各层测试编写规范

#### DAL 层测试

DAL 层需要真实数据库。推荐使用 testcontainers 或真实 PostgreSQL 实例。

```go
// user_dal_test.go
package dal

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "gorm.io/gorm"

    "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity_srv/models"
)

// 假设已设置好测试数据库
var testDB *gorm.DB

func TestUserDALImpl_CreateUser(t *testing.T) {
    // 准备
    dal := NewUserDAL(testDB)
    user := &models.User{
        Name:        "test_user",
        Email:       "test@example.com",
        PhoneNumber: "1234567890",
    }

    // 执行
    err := dal.CreateUser(context.Background(), user)

    // 断言
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}

func TestUserDALImpl_GetUserByID(t *testing.T) {
    // 准备 - 先创建用户
    dal := NewUserDAL(testDB)
    user := &models.User{Name: "test_user", Email: "test@example.com"}
    require.NoError(t, dal.CreateUser(context.Background(), user))

    // 执行 - 查询
    retrieved, err := dal.GetUserByID(context.Background(), user.ID)

    // 断言
    assert.NoError(t, err)
    assert.Equal(t, user.Name, retrieved.Name)
    assert.Equal(t, user.Email, retrieved.Email)
}

func TestUserDALImpl_GetUserByID_NotFound(t *testing.T) {
    dal := NewUserDAL(testDB)

    _, err := dal.GetUserByID(context.Background(), 99999)

    assert.Error(t, err)
    assert.True(t, errno.IsRecordNotFound(err))
}
```

#### Logic 层测试

Logic 层需要 Mock DAL 依赖。使用 mockgen 或手写 mock。

```go
// user_logic_test.go
package logic

import (
    "context"
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"

    "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity_srv/biz"
    "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity_srv/models"
    "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity_srv/pkg/errno"
)

// 假设已生成 mock
// mockgen -source=dal.go -destination=mocks/mock_user_dal.go -package=mocks

func TestUserLogicImpl_CreateUser(t *testing.T) {
    // 准备
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockDAL := mocks.NewMockUserDAL(ctrl)
    mockDAL.EXPECT().
        CreateUser(gomock.Any(), gomock.Any()).
        DoAndReturn(func(ctx context.Context, user *models.User) error {
            user.ID = 1  // 模拟数据库设置 ID
            return nil
        }).
        Times(1)

    logic := NewUserLogic(mockDAL)

    // 执行
    user := &models.User{
        Name:  "test_user",
        Email: "test@example.com",
    }
    err := logic.CreateUser(context.Background(), user)

    // 断言
    assert.NoError(t, err)
    assert.Equal(t, int64(1), user.ID)
}

func TestUserLogicImpl_CreateUser_ValidationError(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockDAL := mocks.NewMockUserDAL(ctrl)
    // DAL 不应该被调用
    mockDAL.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)

    logic := NewUserLogic(mockDAL)

    // 执行 - 无效的邮箱
    user := &models.User{
        Name:  "test_user",
        Email: "invalid-email",
    }
    err := logic.CreateUser(context.Background(), user)

    // 断言 - 应该返回验证错误
    assert.Error(t, err)
    assert.Equal(t, errno.ErrorCodeInvalidEmail, err.(*errno.ErrNo).Code)
}
```

#### Handler 层测试

RPC Handler 层需要 Mock Logic 依赖。

```go
// handler_test.go
package identity_srv

import (
    "context"
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"

    "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity_srv/biz"
    "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity_srv/kitex_gen/base"
    "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity_srv/models"
)

func TestHandler_CreateUser(t *testing.T) {
    // 准备
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockLogic := mocks.NewMockUserLogic(ctrl)
    mockLogic.EXPECT().
        CreateUser(gomock.Any(), gomock.Any()).
        Return(&models.User{
            ID:    1,
            Name:  "test_user",
            Email: "test@example.com",
        }, nil).
        Times(1)

    handler := NewIdentitySrvImpl(mockLogic)

    // 执行
    req := &base.CreateUserRequest{
        Name:  "test_user",
        Email: "test@example.com",
    }
    resp, err := handler.CreateUser(context.Background(), req)

    // 断言
    assert.NoError(t, err)
    assert.Equal(t, int64(1), resp.ID)
    assert.Equal(t, "test_user", resp.Name)
}
```

### 测试命名约定

- **文件**: `<component>_test.go`（如 `user_dal_test.go`）
- **函数**: `Test<InterfaceName>_<MethodName>`
  - 正常情况: `TestUserDALImpl_CreateUser`
  - 错误情况: `TestUserDALImpl_CreateUser_Error`、`TestUserDALImpl_CreateUser_NotFound`

### 测试技巧

```bash
# 运行并检查是否有覆盖率漏洞
go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out | grep -v "100.0%"

# 只运行特定的测试，并显示打印输出
go test -v -run "TestUserLogic" -count=1 ./rpc/identity_srv/biz/logic/

# 运行并在失败时停止
go test -v -failfast ./...
```

## 常见开发流程

### 1. 添加新的 RPC 接口和实现

```bash
# 1. 修改 IDL 定义
vim idl/rpc/identity_srv/xxx.thrift

# 2. 生成 Kitex 代码
cd rpc/identity_srv && ./script/gen_kitex_code.sh

# 3. 实现 Handler
vim rpc/identity_srv/biz/handler.go
# 添加参数校验和业务逻辑委托

# 4. 实现 Logic 和 DAL
vim rpc/identity_srv/biz/logic/xxx_logic.go
vim rpc/identity_srv/biz/dal/xxx_dal.go

# 5. 实现 Converter
vim rpc/identity_srv/biz/converter/xxx_converter.go

# 6. 更新 Wire 依赖注入
vim rpc/identity_srv/wire/provider.go
# 添加 NewXxxLogic, NewXxxDAL 等 provider 函数

# 7. 编写单元测试
vim rpc/identity_srv/biz/logic/xxx_logic_test.go
vim rpc/identity_srv/biz/dal/xxx_dal_test.go

# 8. 验证
cd rpc/identity_srv
go mod tidy
sh build.sh
go test ./... -v
golangci-lint run
```

### 2. 添加新的 HTTP 端点

```bash
# 1. 修改 HTTP IDL
vim idl/http/xxx/xxx.thrift

# 2. 生成 Hertz 代码
cd gateway && ./script/gen_hertz_code.sh

# 3. 实现 HTTP Handler
vim gateway/biz/handler/xxx/xxx_handler.go

# 4. 实现领域服务
vim gateway/internal/domain/service/xxx_service.go
# 负责调用 RPC 客户端并进行业务编排

# 5. 实现 Assembler（可选）
vim gateway/internal/application/assembler/xxx_assembler.go
# 转换 RPC DTO 到 HTTP DTO

# 6. 更新 RPC 客户端（如需要）
vim gateway/internal/infrastructure/client/identity_client.go

# 7. 更新 Wire
vim gateway/internal/wire/provider.go

# 8. 编写单元测试
vim gateway/biz/handler/xxx/xxx_handler_test.go
vim gateway/internal/domain/service/xxx_service_test.go

# 9. 验证
cd gateway
go mod tidy
sh build.sh
go test ./... -v
golangci-lint run
```

### 3. 修改数据模型和数据库

```bash
# 1. 更新 GORM 模型
vim rpc/identity_srv/models/xxx.go

# 2. 如果需要数据库迁移
vim docker/init-scripts/02-tables.sql

# 3. 重启基础设施使迁移生效
cd docker && ./deploy.sh down && ./deploy.sh up

# 4. 更新相关的 DAL 和 Converter
vim rpc/identity_srv/biz/dal/xxx_dal.go
vim rpc/identity_srv/biz/converter/xxx_converter.go

# 5. 更新 Logic 层
vim rpc/identity_srv/biz/logic/xxx_logic.go

# 6. 更新 RPC IDL（如需要）
vim idl/rpc/identity_srv/xxx.thrift

# 7. 生成代码
cd rpc/identity_srv && ./script/gen_kitex_code.sh

# 8. 运行测试
go test ./... -v

# 9. 重新构建和测试
sh build.sh && go test ./... -v
```

## 故障排查

### 常见问题和解决方案

#### 问题：Go Workspace 中的模块找不到

**症状**：`cannot find module` 或 `missing go.mod`

**解决方案**：
```bash
# 检查 go.work 是否正确
cat go.work

# 确保 go.work 包含两个模块
# use (
#    ./gateway
#    ./rpc/identity_srv
# )

# 重新初始化 workspace
go work sync

# 清理缓存
go clean -modcache
```

#### 问题：Kitex 代码生成失败

**症状**：`./script/gen_kitex_code.sh` 返回错误

**解决方案**：
```bash
# 检查 Thrift 文件语法
cd rpc/identity_srv

# 查看生成脚本
cat script/gen_kitex_code.sh

# 手动运行 kitex 命令调试
kitex -module github.com/masonsxu/cloudwego-microservice-demo/rpc/identity_srv \
      -service identity_srv \
      ../../../idl/rpc/identity_srv/*.thrift

# 清理旧的生成文件
rm -rf kitex_gen/

# 重新生成
./script/gen_kitex_code.sh
```

#### 问题：Wire 依赖注入生成失败

**症状**：`wire_gen.go: no providers found` 或循环依赖

**解决方案**：
```bash
# 进入 wire 目录
cd rpc/identity_srv/wire

# 检查 provider 声明
cat provider.go

# 检查 wire.go 的 wire.Build() 声明
cat wire.go

# 删除旧的生成文件
rm -f wire_gen.go

# 运行 wire
wire

# 如果仍然失败，检查 provider 函数签名
# 确保每个 provider：
# 1. 只创建一个依赖
# 2. 返回类型是接口而非具体实现
# 3. 没有循环依赖
```

#### 问题：数据库连接失败

**症状**：`connection refused` 或 `GORM connection error`

**解决方案**：
```bash
# 检查基础设施是否运行
cd docker && ./deploy.sh ps

# 如果 postgres 不运行，启动它
./deploy.sh up

# 检查数据库连接信息
cat docker/docker-compose.yml | grep -A 5 "postgres"

# 手动测试连接
psql -h 127.0.0.1 -U postgres -d identity_db -c "SELECT 1"

# 检查环境变量
echo $DB_DSN

# 查看服务日志
./deploy.sh logs postgres
```

#### 问题：测试失败 - record not found

**症状**：`ErrRecordNotFound` 在测试中返回

**解决方案**：
```bash
# 确保测试数据库已初始化
cd docker && ./deploy.sh up

# 检查测试使用的数据库
grep "testdb\|test_db" rpc/identity_srv/config/*.go

# 查看测试设置代码
grep -r "gorm.Open\|testDB" rpc/identity_srv/**/*_test.go | head -5

# 运行单个测试并查看详细输出
go test -v -run "TestUserDALImpl_GetUserByID" ./rpc/identity_srv/biz/dal/ -count=1

# 检查测试前是否创建了数据
grep -B 10 "GetUserByID" rpc/identity_srv/biz/dal/xxx_test.go
```

#### 问题：导入顺序错误，golangci-lint 报错

**症状**：`gci: file is not gci formatted`

**解决方案**：
```bash
# 自动修复导入顺序
golangci-lint format

# 或手动运行 gci
cd rpc/identity_srv && gci write .
cd gateway && gci write .

# 检查配置
grep -A 10 "gci:" .golangci.yml
```

#### 问题：RPC 调用返回 BizStatusError 但无法解析

**症状**：客户端收到错误但无法获取错误码

**解决方案**：
```bash
# 检查客户端是否配置了 TTHeader MetaHandler
grep -r "TTHeaderMetaHandler\|NewBizStatusError" gateway/internal/infrastructure/client/

# 检查错误处理
grep -r "kerrors.FromBizStatusError" gateway/

# 验证 errno 包中的错误转换函数
cat rpc/identity_srv/pkg/errno/error.go | grep -A 5 "ToKitexError\|FromBizStatusError"
```

#### 问题：Swagger 文档生成失败

**症状**：`/swagger/index.html` 404

**解决方案**：
```bash
# 检查 Swagger 生成配置
grep -r "swag init\|swagger" gateway/script/

# 手动生成 Swagger
cd gateway && swag init -g biz/handler/xxx/xxx_handler.go

# 检查 swag 配置
cat gateway/.swag.yml 2>/dev/null || echo "No .swag.yml found"

# 验证生成的 Swagger 文件
ls -la gateway/docs/
```

### 调试技巧

```bash
# 启用详细日志
export LOG_LEVEL=debug
sh output/bootstrap.sh

# 附加调试器（使用 dlv）
dlv debug ./cmd/main.go -- -addr=:8891

# 查看环境变量
env | grep -E "IDENTITY|DB_|LOGO"

# 查看端口占用
lsof -i :8080 | grep -v COMMAND
lsof -i :8891 | grep -v COMMAND

# 查看进程日志
ps aux | grep bootstrap
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
9. **UI配色规范**: 所有UI设计必须遵循"奢华摩羯座配色方案"，参考 [UI配色设计规范](#UI配色设计规范)

## Go Workspace 开发注意事项

### 模块结构

项目使用 Go Workspace 管理两个独立模块：

```
cloudwego-microservice-demo/
├── go.work                 # Workspace 定义
├── go.work.sum             # Workspace 依赖锁
├── gateway/
│   ├── go.mod              # Gateway 模块定义
│   ├── go.sum              # Gateway 依赖锁
│   └── ...
└── rpc/
    └── identity_srv/
        ├── go.mod          # Identity 服务模块定义
        ├── go.sum          # Identity 服务依赖锁
        └── ...
```

### Workspace 关键点

- **测试**: `go test ./...` 从项目根目录会测试所有模块
- **构建**: 两个模块分别构建，使用各自的 `go.mod`
- **依赖**: 每个模块有独立的依赖，可以独立升级
- **同步**: 运行 `go work sync` 同步模块元数据

### 常见 Workspace 命令

```bash
# 查看 workspace 配置
cat go.work

# 同步 workspace（解决模块间的依赖问题）
go work sync

# 在 workspace 中编辑模块
go work edit -use=./new_module

# 从 workspace 移除模块
go work edit -dropuse=./module_to_remove
```

### 跨模块导入

网关导入 RPC 服务中的公共包：

```go
// gateway 中导入 identity_srv 的 errno 包
import (
    "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity_srv/pkg/errno"
)
```

## 详细文档

详细的开发指南、配置说明、部署指南和故障排查请参考 `docs/` 目录：

- [快速开始](docs/01-getting-started.md) - 环境配置和首次运行
- [架构设计](docs/02-architecture.md) - 系统架构和设计原理
- [开发指南](docs/03-development.md) - 详细的开发流程和最佳实践
- [配置参考](docs/04-configuration.md) - 所有配置项说明
- [部署指南](docs/05-deployment.md) - 生产环境部署步骤
- [故障排查](docs/06-troubleshooting.md) - 常见问题和解决方案
- [权限管理](docs/07-permission-management.md) - Casbin RBAC 配置和使用
- [奢华摩羯座配色规范](docs/CAPRICORN-THEME-GUIDE.md) - 官方配色方案和使用指南

## 关键文件速查

| 文件 | 用途 |
|------|------|
| `go.work` | Go Workspace 配置 |
| `.golangci.yml` | golangci-lint 配置（代码检查和格式化） |
| `scripts/git-hooks/pre-commit` | Git 预提交钩子 |
| `docker/docker-compose.yml` | 本地开发基础设施配置 |
| `docker/deploy.sh` | 基础设施启动脚本 |
| `rpc/identity_srv/script/gen_kitex_code.sh` | Kitex 代码生成脚本 |
| `gateway/script/gen_hertz_code.sh` | Hertz 代码生成脚本 |
| `idl/` | Thrift 接口定义目录 |

## 快速参考

### 错误码体系

6位数字格式：`A-BB-CCC`
- `A`: 级别（1=系统、2=业务）
- `BB`: 模块（00=通用、01=用户、02=认证、03=权限）
- `CCC`: 具体错误编码

示例：
- `100000` - 系统错误（通用）
- `200101` - 业务错误，用户模块，用户不存在

### 导入顺序

```go
import (
    // 标准库
    "context"
    "fmt"
    "os"

    // 第三方库
    "github.com/cloudwego/hertz/pkg/app"
    "gorm.io/gorm"

    // 项目包
    "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain"
    "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity_srv/pkg/errno"
)
```

### 常见命令速记

```bash
# 一键启动全部
cd docker && ./deploy.sh up &
cd rpc/identity_srv && sh build.sh && sh output/bootstrap.sh &
cd gateway && sh build.sh && sh output/bootstrap.sh

# 一键测试和检查
go test ./... -v && golangci-lint run

# 一键生成代码和构建
cd rpc/identity_srv && ./script/gen_kitex_code.sh && sh build.sh
cd gateway && ./script/gen_hertz_code.sh && sh build.sh
```

## 更新日志

最近的项目更新：
- 引入 Casbin 权限管理模块支持多角色模式
- 添加权限管理模块文档
- 添加 AGENTS.md 文档供智能编码代理使用
- 新增 UI配色设计规范 - 奢华摩羯座配色方案，所有 UI 设计必须遵循
- 将配色规范整理为独立文档 [docs/CAPRICORN-THEME-GUIDE.md](docs/CAPRICORN-THEME-GUIDE.md)，提供完整的 CSS 和 Tailwind 配置示例

详见 [Git 历史](https://github.com/masonsxu/cloudwego-microservice-demo/commits/main)

