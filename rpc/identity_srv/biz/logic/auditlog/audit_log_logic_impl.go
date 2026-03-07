package auditlog

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	auditlogDAL "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/auditlog"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
)

// LogicImpl 审计日志业务逻辑实现
type LogicImpl struct {
	dal       dal.DAL
	converter converter.Converter
}

// NewLogic 创建审计日志业务逻辑实例
func NewLogic(dal dal.DAL, converter converter.Converter) AuditLogLogic {
	return &LogicImpl{
		dal:       dal,
		converter: converter,
	}
}

// CreateAuditLog 创建审计日志
func (l *LogicImpl) CreateAuditLog(
	ctx context.Context,
	req *identity_srv.CreateAuditLogRequest,
) error {
	if req == nil {
		return errno.ErrInvalidParams.WithMessage("请求不能为空")
	}

	// 转换请求为模型
	auditLog := l.converter.AuditLog().CreateRequestToModel(req)

	// 直接写入，不使用事务（审计日志是追加操作，无需事务保护）
	if err := l.dal.AuditLog().Create(ctx, auditLog); err != nil {
		return errno.ErrOperationFailed.WithMessage("创建审计日志失败: " + err.Error())
	}

	return nil
}

// ListAuditLogs 查询审计日志列表
func (l *LogicImpl) ListAuditLogs(
	ctx context.Context,
	req *identity_srv.ListAuditLogsRequest,
) (*identity_srv.ListAuditLogsResponse, error) {
	if req == nil {
		return nil, errno.ErrInvalidParams.WithMessage("请求不能为空")
	}

	// 转换分页参数
	opts := l.converter.Base().PageRequestToQueryOptions(req.Page)

	// 构建查询条件
	conditions := &auditlogDAL.AuditLogQueryConditions{
		Page: opts,
	}

	if req.UserID != nil {
		conditions.UserID = req.UserID
	}

	if req.Action != nil {
		action := models.AuditAction(*req.Action)
		conditions.Action = &action
	}

	if req.Resource != nil {
		conditions.Resource = req.Resource
	}

	if req.Success != nil {
		conditions.Success = req.Success
	}

	if req.StartTime != nil {
		conditions.StartTime = req.StartTime
	}

	if req.EndTime != nil {
		conditions.EndTime = req.EndTime
	}

	// 查询
	logs, pageResult, err := l.dal.AuditLog().FindWithConditions(ctx, conditions)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询审计日志失败: " + err.Error())
	}

	// 查询全局统计数据
	stats, err := l.dal.AuditLog().GetStatsByConditions(ctx, conditions)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询审计日志统计失败: " + err.Error())
	}

	// 转换结果
	return &identity_srv.ListAuditLogsResponse{
		AuditLogs: l.converter.AuditLog().ModelsToThrift(logs),
		Page:      l.converter.Base().PageResponseToThrift(pageResult),
		Stats:     l.converter.AuditLog().StatsToThrift(stats),
	}, nil
}
