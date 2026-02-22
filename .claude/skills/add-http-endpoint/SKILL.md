---
name: add-http-endpoint
description: 添加新 HTTP 端点（IDL→网关层全流程）
argument-hint: <METHOD /api/v1/path>
---

# 添加新 HTTP 端点

端到端引导在网关层添加新的 HTTP 端点，前提是对应的 RPC 接口已存在。

**用户输入**: $ARGUMENTS（HTTP 方法和路径，例如 `POST /api/v1/departments`）

---

## 前置检查

1. 确认对应的 RPC 接口已在 `rpc/identity_srv/` 中实现
2. 如果 RPC 接口不存在，提示用户先执行 `/add-rpc-method`

## 执行流程

### Step 1: 需求确认

根据用户输入 `$ARGUMENTS`，确认：
- HTTP 方法（GET/POST/PUT/DELETE）
- 路由路径
- 对应的 RPC 方法
- 所属服务领域（identity / permission）
- 是否需要认证（JWT）和权限控制（Casbin）

### Step 2: HTTP IDL 定义

在 `idl/http/<service>/` 下找到对应的 Thrift IDL 文件：
1. 定义请求/响应结构体
2. 添加 HTTP 方法和路由注解（`api.post`、`api.get` 等）

**参考已有 HTTP IDL 文件的注解风格**。

### Step 3: 代码生成

```bash
cd gateway && ./script/gen_hertz_code.sh <service>
```

支持的 service 参数：`identity`、`permission`，或不传参数生成全部。

生成后确认 `biz/handler/` 和 `biz/router/` 下新代码已生成。

### Step 4: Domain Service（领域服务）

- **参考目录**：`gateway/internal/domain/service/identity/` 或 `gateway/internal/domain/service/permission/`
- **接口文件**：参考 `identity_service.go` 或 `permission_service.go`
  - 在聚合接口中添加新方法
- **实现文件**：参考已有的 `*_service_impl.go`
  - 调用 RPC Client，编排业务流程
  - 错误处理：使用 `kerrors.FromBizStatusError()` 解析 RPC 错误

### Step 5: Assembler（DTO 转换）

- **参考目录**：`gateway/internal/application/assembler/identity/` 或 `gateway/internal/application/assembler/permission/`
- **接口文件**：参考 `interfaces.go`
  - 添加新的转换方法声明
- **实现文件**：参考已有的 `*_assembler.go`
  - RPC DTO → HTTP DTO 转换
- **聚合文件**：更新 `assembler.go`

### Step 6: Handler 完善

- **文件**：`gateway/biz/handler/<service>/<service>_service.go`
- **模式**：Hertz 生成的 Handler 只做参数绑定和响应，业务逻辑委托给 Domain Service
- 参考已有 Handler 的实现风格

### Step 7: Wire 依赖注入

如果新增了 Provider：
1. 更新 `gateway/internal/wire/provider.go`
2. 运行 `cd gateway/internal/wire && wire`

### Step 8: 权限配置

如果端点需要 Casbin 权限控制：
1. 先读取 `docs/04-权限管理/权限管理设计.md` 了解权限体系
2. 配置 Casbin 策略（角色 → 资源 → 操作的映射）
3. 确认中间件执行顺序：JWT → Casbin

### Step 9: 验证

```bash
cd gateway && go build ./...
cd gateway && go vet ./...
cd gateway && golangci-lint run
```

---

## 注意事项

- 先读取 `docs/02-开发规范/开发指南.md` 了解完整开发规范
- 中间件执行顺序：OpenTelemetry → RequestID → ResponseHeader → Trace → CORS → ErrorHandler → JWT → Casbin → ETag
- 不使用 YAML 配置文件，配置通过环境变量管理
- 导入顺序：标准库 → 第三方库 → 项目内部包
