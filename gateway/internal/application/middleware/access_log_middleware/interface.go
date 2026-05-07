// Package middleware 提供轻量访问日志中间件，记录每个请求的最小要素：
// method/path/status/duration/user_id/request_id。
//
// 取代 audit_middleware：HTTP 层不再做"业务审计"，真正的合规审计应在
// 业务系统内部针对领域事件落地。
package middleware

import "github.com/cloudwego/hertz/pkg/app"

// AccessLogMiddlewareService 访问日志中间件接口
type AccessLogMiddlewareService interface {
	MiddlewareFunc() app.HandlerFunc
}
