package audit_middleware

import (
	"github.com/cloudwego/hertz/pkg/app"
)

// AuditMiddlewareService 审计中间件服务接口
type AuditMiddlewareService interface {
	// MiddlewareFunc 返回审计中间件函数
	MiddlewareFunc() app.HandlerFunc
}
