---
name: gateway-logging
description: API 网关日志编写规范。在域服务、中间件、基础设施层编写日志代码时自动加载。
---

# API 网关日志编写规范

在 API 网关层编写日志代码时，必须遵循以下规范。

## 核心原则

1. **统一使用 zerolog** + `tracelog` 包 (`pkg/log/trace_logger.go`)
2. **禁止** 使用 `s.Logger()` 的 printf 风格 API（已废弃）
3. **所有日志必须携带追踪上下文**（trace_id/request_id/span_id）
4. **追踪字段由 OTelHook 自动注入**，不需要手动处理，但必须正确传递 ctx

## 域服务层

使用 `s.Log(ctx)` 获取带追踪信息的 `*zerolog.Logger`：

```go
s.Log(ctx).Error().Err(err).Msg("调用 RPC 服务失败")
s.Log(ctx).Info().Str("patient_id", id).Msg("操作成功")
s.Log(ctx).Warn().Msg("参数缺失")
```

## 中间件层

使用 `*zerolog.Logger` + `tracelog.Event(ctx, ...)`：

```go
import tracelog "gitlab.manteia.com/radius/radius-backend/api/radius-api-gateway/pkg/log"

tracelog.Event(ctx, m.logger.Warn()).
    Str("component", "jwt_middleware").
    Msg("Token 已吊销")
```

## 日志级别与 Jaeger 联动

OTelHook (`pkg/log/otel_hook.go`) 会根据日志级别自动同步到 Jaeger Span：

- `Error` → `span.RecordError()` + `span.SetStatus(Error)` → **Jaeger Span 标红**
- `Warn` → `span.AddEvent()` → Jaeger Span 详情中可见（不标红）
- `Info`/`Debug` → 仅写日志文件，不影响 Span

因此：
- 参数校验失败 → `Warn`（不要用 Error，否则 Jaeger Span 会被误标红）
- RPC 调用/业务逻辑失败 → `Error`
- 正常操作完成 → `Info`
- 调试/连接信息 → `Debug`

**禁止**直接使用不带 ctx 的 logger 调用（如 `m.logger.Warn().Msg(...)`），OTelHook 依赖 `e.GetCtx()` 获取 context，不传 ctx 则追踪字段和 Jaeger 联动均失效。

## 字段规范

- 错误：`.Err(err)`（不用 `.Str("error", ...)`）
- 字符串：`.Str("key", val)`
- 整数：`.Int("key", val)`
- 组件标识：`.Str("component", "xxx_middleware")`

## 参考

- 规范文档：`doc/01-开发规范/API网关日志编写规范.md`
- OTelHook：`api/radius_api_gateway/pkg/log/otel_hook.go`
- tracelog 包：`api/radius_api_gateway/pkg/log/trace_logger.go`
- 标杆实现：`application/middleware/error_middleware/error_middleware.go`
