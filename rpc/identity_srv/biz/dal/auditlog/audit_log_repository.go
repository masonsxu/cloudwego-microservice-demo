package auditlog

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// AuditLogRepositoryImpl 审计日志仓储实现
type AuditLogRepositoryImpl struct {
	db *gorm.DB
}

// NewAuditLogRepository 创建审计日志仓储实例
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &AuditLogRepositoryImpl{
		db: db,
	}
}

// Create 创建审计日志记录
func (r *AuditLogRepositoryImpl) Create(ctx context.Context, log *models.AuditLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("创建审计日志失败: %w", err)
	}

	return nil
}

// FindWithConditions 根据组合查询条件查询审计日志列表
func (r *AuditLogRepositoryImpl) FindWithConditions(
	ctx context.Context,
	conditions *AuditLogQueryConditions,
) ([]*models.AuditLog, *models.PageResult, error) {
	opts := base.NewQueryOptions()
	if conditions != nil && conditions.Page != nil {
		opts = conditions.Page
	}

	// 构建查询
	query := r.db.WithContext(ctx).Model(&models.AuditLog{})
	query = r.applyFilterConditions(query, conditions)

	// 计算总数
	var total int64

	countQuery := r.db.WithContext(ctx).Model(&models.AuditLog{})
	countQuery = r.applyFilterConditions(countQuery, conditions)

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("统计审计日志总数失败: %w", err)
	}

	// 排序（默认按创建时间降序）
	orderClause := r.buildOrderClause(opts)
	query = query.Order(orderClause)

	// 分页查询
	var logs []*models.AuditLog

	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Offset(int(offset)).Limit(int(opts.PageSize)).Find(&logs).Error; err != nil {
		return nil, nil, fmt.Errorf("查询审计日志列表失败: %w", err)
	}

	pageResult := models.NewPageResult(int32(total), opts.Page, opts.PageSize)

	return logs, pageResult, nil
}

// applyFilterConditions 应用过滤条件到查询
func (r *AuditLogRepositoryImpl) applyFilterConditions(
	query *gorm.DB,
	conditions *AuditLogQueryConditions,
) *gorm.DB {
	if conditions == nil {
		return query
	}

	if conditions.UserID != nil && *conditions.UserID != "" {
		query = query.Where("user_id = ?", *conditions.UserID)
	}

	if conditions.Action != nil {
		query = query.Where("action = ?", *conditions.Action)
	}

	if conditions.Resource != nil && *conditions.Resource != "" {
		query = query.Where("resource = ?", *conditions.Resource)
	}

	if conditions.Success != nil {
		query = query.Where("success = ?", *conditions.Success)
	}

	if conditions.StartTime != nil {
		query = query.Where("created_at >= ?", *conditions.StartTime)
	}

	if conditions.EndTime != nil {
		query = query.Where("created_at <= ?", *conditions.EndTime)
	}

	return query
}

// buildOrderClause 构建排序子句
func (r *AuditLogRepositoryImpl) buildOrderClause(opts *base.QueryOptions) string {
	orderBy := "created_at"
	if opts.OrderBy != "" {
		orderBy = opts.OrderBy
	}

	if opts.OrderDesc {
		return orderBy + " DESC"
	}

	return orderBy + " ASC"
}

// GetStatsByConditions 根据筛选条件获取全局统计数据（不受分页限制）
func (r *AuditLogRepositoryImpl) GetStatsByConditions(
	ctx context.Context,
	conditions *AuditLogQueryConditions,
) (*AuditLogStats, error) {
	query := r.db.WithContext(ctx).Model(&models.AuditLog{})
	query = r.applyFilterConditions(query, conditions)

	var stats AuditLogStats

	err := query.Select(
		"COUNT(*) AS total_count, " +
			"SUM(CASE WHEN success = true THEN 1 ELSE 0 END) AS success_count, " +
			"COALESCE(AVG(duration_ms), 0) AS avg_duration_ms",
	).Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("统计审计日志失败: %w", err)
	}

	return &stats, nil
}
