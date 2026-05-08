# iamclient SDK 使用指南

> **状态**：随 Phase 4 第一个 PR 落地（仅 SDK + 网关透传），下一 PR 处理
> identity_srv 接入与内部 Casbin 清理。
> **代码位置**：仓库顶层 `iamclient/`（独立 Go module，参与 `go.work`）
> **依赖**：`policy_srv`（PDP 决策服务，Phase 2 已落地）

## 1. 设计目标

业务系统接入鉴权时，应该只看到这一组 API：

| API | 用途 |
|---|---|
| `iamclient.New(cfg)` | 应用启动时创建 Client（goroutine-safe，复用一次） |
| `cli.SubjectFromHeader(h)` | HTTP 请求入口从 header 还原 Subject |
| `cli.SubjectFromContext(ctx)` | RPC 请求入口从 metainfo 还原 Subject |
| `subject.Check(ctx, action, resource)` | 询问 PDP 是否放行（返回 `*Decision`） |
| `subject.MustCheck(ctx, ...)` | 同上，拒绝时返回 `*PermissionDeniedError` |

业务侧**禁止**出现：JWT 解析、Casbin Enforcer、直连 identity_srv 数据库、自行
拼接权限码列表。所有这些都是 PDP 的职责。

## 2. 接入三步法

### 第 1 步：声明依赖

业务模块的 `go.mod`：

```go
require github.com/masonsxu/cloudwego-microservice-demo/iamclient v0.0.0-00010101000000-000000000000

replace github.com/masonsxu/cloudwego-microservice-demo/iamclient => ../iamclient
```

如果业务模块要传递依赖到 `policy_srv` 的 kitex_gen，还需加：

```go
replace github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv => ../rpc/policy_srv
```

`go.work` 中确保两者都已 `use`：

```go
use (
    ./gateway
    ./iamclient
    ./rpc/identity_srv
    ./rpc/policy_srv
    ./your-business-module
)
```

### 第 2 步：Wire Provider（应用启动时一次性创建）

```go
// internal/wire/provider.go
package wire

import (
    "github.com/masonsxu/cloudwego-microservice-demo/iamclient"
)

// ProvideIAMClient 创建并初始化 IAM 客户端
func ProvideIAMClient(cfg *config.Config) (*iamclient.Client, func(), error) {
    cli, err := iamclient.New(iamclient.Config{
        EtcdEndpoints: []string{cfg.Etcd.Address},
        PolicyService: "policy-service", // 默认值，可省略
        CallerService: cfg.Server.Name,  // 当前服务名
        // 缓存默认 10000 条 / 30s TTL，按需调整：
        // CacheSize: 5000,
        // CacheTTL:  10 * time.Second,
    })
    if err != nil {
        return nil, nil, err
    }
    return cli, func() { _ = cli.Close() }, nil
}
```

### 第 3 步：在 Handler 中鉴权

#### 3a. HTTP 业务（Hertz）

```go
import (
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/masonsxu/cloudwego-microservice-demo/iamclient"
)

type PatientHandler struct {
    iam *iamclient.Client
}

func (h *PatientHandler) Get(ctx context.Context, c *app.RequestContext) {
    subject, err := h.iam.SubjectFromHeader(&c.Request.Header)
    if err != nil {
        c.AbortWithStatus(consts.StatusUnauthorized)
        return
    }

    patientID := c.Param("id")
    if err := subject.MustCheck(ctx, "read", "patient:"+patientID); err != nil {
        if errors.Is(err, iamclient.ErrPermissionDenied) {
            c.AbortWithStatus(consts.StatusForbidden)
        } else {
            c.AbortWithError(consts.StatusInternalServerError, err)
        }
        return
    }

    // 业务逻辑...
}
```

> **注**：Hertz 的 `c.Request.Header` 是 `*protocol.RequestHeader`，提供
> `Get(name string) string`，满足 SDK 的 `HeaderGetter` 接口。net/http.Header
> 也兼容。

#### 3b. Kitex RPC 业务

```go
import (
    "github.com/masonsxu/cloudwego-microservice-demo/iamclient"
)

type PatientServiceImpl struct {
    iam *iamclient.Client
}

func (s *PatientServiceImpl) GetPatient(ctx context.Context, req *pb.GetPatientRequest) (*pb.PatientResponse, error) {
    subject, err := s.iam.SubjectFromContext(ctx)
    if err != nil {
        return nil, errno.ErrUnauthenticated
    }

    if err := subject.MustCheck(ctx, "read", "patient:"+req.GetId()); err != nil {
        if errors.Is(err, iamclient.ErrPermissionDenied) {
            return nil, errno.ErrPermissionDenied
        }
        return nil, errno.ErrInternal.WithError(err)
    }

    // 业务逻辑...
}
```

> 这条路径依赖网关 `identity_propagation_middleware` 把 HTTP header 中的
> 身份字段镜像到 Kitex metainfo 持久值，由 TTHeader 透传到下游。

## 3. ABAC：注入资源属性

PDP 决策需要部门、所有者等资源属性时，用 `WithResourceAttr`：

```go
// 业务侧从自己的 DB 读出资源属性后，传给 PDP
patient := patientRepo.Get(ctx, id)
if err := subject.MustCheck(ctx, "read", "patient:"+id,
    iamclient.WithResourceAttr("department_id", patient.DepartmentID),
    iamclient.WithResourceAttr("owner_id", patient.OwnerID),
); err != nil {
    // ...
}
```

约定：**资源属性永远由业务侧从自己的领域 DB 读取后传入**，不在 token / Subject
中携带，也不让网关来填——网关无知于业务领域。

## 4. 缓存语义

- **默认开启**：LRU + TTL，单 key 30s 内复用同一条决策。
- **Key**：`(jti | user_id+roles, tenant, action, resource, sorted resource_attrs)`。
  jti 优先，roles 切片排序后参与（顺序不影响命中）。
- **跳过缓存**：`subject.Check(ctx, "x", "y", iamclient.WithoutCache())`。
  调试或对最新策略生效时间敏感时使用。
- **不主动失效**：当前版本不订阅策略变更通知，TTL 内可能命中策略变更前的旧决策。
  可接受范围：30 秒。后续如引入 etcd watch，可在 SDK 内部 invalidate。

## 5. 错误处理建议

```go
err := subject.MustCheck(ctx, action, resource)

switch {
case err == nil:
    // 通过
case errors.Is(err, iamclient.ErrPermissionDenied):
    // 业务上的"无权限"，转 403
    var pde *iamclient.PermissionDeniedError
    errors.As(err, &pde)
    log.Warn().Str("reason", pde.Reason).Msg("denied")
    return ErrForbidden
default:
    // 网络 / RPC 错误，转 500（不要兜底为放行）
    return fmt.Errorf("authz failed: %w", err)
}
```

**铁律**：PDP 不可达时，业务侧绝不允许"宽松降级"放行。失败即拒绝（fail-closed），
让 SRE 看到错误而非让权限漏洞静默生效。

## 6. 测试 Subject 不调真实 PDP

`Subject.client` 是私有字段，不能手工构造。测试中用 `iamclient.NewForTest` 风格
（暂未提供）目前的做法是：

- 单元测试：直接 mock 业务接口（不经过 Check）；
- 集成测试：起 policy_srv（podman pod），让真实链路跑通。

> 后续若有需要，会补 `iamclient.NewWithPolicyClient(cli policyservice.Client)`
> 测试构造器，让业务可以注入 fake policy client。

## 7. 反模式（提案 §3）

| ❌ 反例 | ✅ 正例 |
|---|---|
| 业务 handler 里 `casbin.Enforce(...)` | `subject.MustCheck(...)` |
| 业务里 `jwt.Parse(c.GetHeader("Authorization"))` | `cli.SubjectFromHeader(&c.Request.Header)` |
| token 里塞 `dataScope`，handler 读 token 决策 | PDP 返回 `Decision.DataScopeHint`，业务自取 |
| handler 里 `if username == "superadmin"` 白名单 | 在 PDP 策略里写 `p, role:superadmin, *, *, *, all` |
| 每次都先 `ListPermissions` 自己判断 | 每个具体动作 `Check(action, resource)` |

## 8. 相关文档

- 总设计：`docs/04-权限管理/重构提案-网关边界与权限模型.md`
- PDP 服务：`rpc/policy_srv/`
- 网关授权（路由级 ACL）：`docs/04-权限管理/重构提案...md` §5.4
- 网关身份透传：`gateway/internal/application/middleware/identity_propagation_middleware/`
