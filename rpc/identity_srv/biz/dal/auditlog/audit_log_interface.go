package auditlog

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// AuditLogStats 审计日志统计结果
type AuditLogStats struct {
	TotalCount    int64
	SuccessCount  int64
	AvgDurationMs float64
}

// AuditLogRepository 审计日志仓储接口
type AuditLogRepository interface {
	// Create 创建审计日志记录
	Create(ctx context.Context, log *models.AuditLog) error

	// FindWithConditions 根据组合查询条件查询审计日志列表
	FindWithConditions(
		ctx context.Context,
		conditions *AuditLogQueryConditions,
	) ([]*models.AuditLog, *models.PageResult, error)

	// GetStatsByConditions 根据筛选条件获取全局统计数据（不受分页限制）
	GetStatsByConditions(
		ctx context.Context,
		conditions *AuditLogQueryConditions,
	) (*AuditLogStats, error)
}

// AuditLogQueryConditions 审计日志查询条件
type AuditLogQueryConditions struct {
	UserID    *string
	Action    *models.AuditAction
	Resource  *string
	Success   *bool
	StartTime *int64
	EndTime   *int64
	Page      *base.QueryOptions
}
