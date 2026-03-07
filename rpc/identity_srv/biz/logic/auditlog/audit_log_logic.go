package auditlog

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// AuditLogLogic 审计日志业务逻辑接口
type AuditLogLogic interface {
	// CreateAuditLog 创建审计日志
	CreateAuditLog(ctx context.Context, req *identity_srv.CreateAuditLogRequest) error

	// ListAuditLogs 查询审计日志列表
	ListAuditLogs(
		ctx context.Context,
		req *identity_srv.ListAuditLogsRequest,
	) (*identity_srv.ListAuditLogsResponse, error)
}
