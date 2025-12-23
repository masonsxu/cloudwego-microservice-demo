package logo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
)

// logoRepository 组织Logo仓储实现
type logoRepository struct {
	db *gorm.DB
	base.BaseRepository[models.OrganizationLogo]
}

// NewLogoRepository 创建Logo仓储实例
func NewLogoRepository(db *gorm.DB) LogoRepository {
	baseRepo := base.NewBaseRepository[models.OrganizationLogo](db)

	return &logoRepository{
		db:             db,
		BaseRepository: baseRepo,
	}
}

// DB 获取数据库实例
func (r *logoRepository) DB() *gorm.DB {
	return r.db
}

// GetByFileID 根据文件ID获取Logo
func (r *logoRepository) GetByFileID(
	ctx context.Context,
	fileID string,
) (*models.OrganizationLogo, error) {
	var logo models.OrganizationLogo

	err := r.DB().WithContext(ctx).
		Where("file_id = ?", fileID).
		First(&logo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrLogoNotFound
		}

		return nil, errno.WrapDatabaseError(err, "查询Logo失败")
	}

	return &logo, nil
}

// GetByOrganizationID 根据组织ID获取绑定的Logo
func (r *logoRepository) GetByOrganizationID(
	ctx context.Context,
	organizationID uuid.UUID,
) (*models.OrganizationLogo, error) {
	var logo models.OrganizationLogo

	err := r.DB().WithContext(ctx).
		Where("bound_organization_id = ? AND status = ?", organizationID, models.LogoStatusBound).
		First(&logo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrLogoNotFound
		}

		return nil, errno.WrapDatabaseError(err, "查询组织Logo失败")
	}

	return &logo, nil
}

// ListByStatus 根据状态获取Logo列表
func (r *logoRepository) ListByStatus(
	ctx context.Context,
	status models.OrganizationLogoStatus,
	opts *base.QueryOptions,
) ([]*models.OrganizationLogo, *models.PageResult, error) {
	if opts == nil {
		opts = base.NewQueryOptions()
	}

	query := r.DB().WithContext(ctx).Where("status = ?", status)

	// 应用排序
	if opts.OrderBy != "" {
		orderClause := opts.OrderBy
		if opts.OrderDesc {
			orderClause += " DESC"
		} else {
			orderClause += " ASC"
		}

		query = query.Order(orderClause)
	} else {
		// 默认按创建时间倒序
		query = query.Order("created_at DESC")
	}

	// 应用分页
	offset := (opts.Page - 1) * opts.PageSize
	query = query.Offset(int(offset)).Limit(int(opts.PageSize))

	var logos []*models.OrganizationLogo

	err := query.Find(&logos).Error
	if err != nil {
		return nil, nil, errno.WrapDatabaseError(err, "查询Logo列表失败")
	}

	// 获取总数
	var total int64

	countQuery := r.DB().
		WithContext(ctx).
		Model(&models.OrganizationLogo{}).
		Where("status = ?", status)

	err = countQuery.Count(&total).Error
	if err != nil {
		return nil, nil, errno.WrapDatabaseError(err, "查询Logo总数失败")
	}

	pageResult := models.NewPageResult(int32(total), opts.Page, opts.PageSize)

	return logos, pageResult, nil
}

// ListByUploader 根据上传者ID获取Logo列表
func (r *logoRepository) ListByUploader(
	ctx context.Context,
	uploaderID uuid.UUID,
	opts *base.QueryOptions,
) ([]*models.OrganizationLogo, *models.PageResult, error) {
	if opts == nil {
		opts = base.NewQueryOptions()
	}

	query := r.DB().WithContext(ctx).Where("uploaded_by = ?", uploaderID)

	// 应用排序
	if opts.OrderBy != "" {
		orderClause := opts.OrderBy
		if opts.OrderDesc {
			orderClause += " DESC"
		} else {
			orderClause += " ASC"
		}

		query = query.Order(orderClause)
	} else {
		// 默认按创建时间倒序
		query = query.Order("created_at DESC")
	}

	// 应用分页
	offset := (opts.Page - 1) * opts.PageSize
	query = query.Offset(int(offset)).Limit(int(opts.PageSize))

	var logos []*models.OrganizationLogo

	err := query.Find(&logos).Error
	if err != nil {
		return nil, nil, errno.WrapDatabaseError(err, "查询Logo列表失败")
	}

	// 获取总数
	var total int64

	countQuery := r.DB().WithContext(ctx).Model(&models.OrganizationLogo{}).
		Where("uploaded_by = ?", uploaderID)

	err = countQuery.Count(&total).Error
	if err != nil {
		return nil, nil, errno.WrapDatabaseError(err, "查询Logo总数失败")
	}

	pageResult := models.NewPageResult(int32(total), opts.Page, opts.PageSize)

	return logos, pageResult, nil
}

// ExistsByFileID 检查文件ID是否存在
func (r *logoRepository) ExistsByFileID(ctx context.Context, fileID string) (bool, error) {
	var count int64

	err := r.DB().WithContext(ctx).
		Model(&models.OrganizationLogo{}).
		Where("file_id = ?", fileID).
		Count(&count).Error
	if err != nil {
		return false, errno.WrapDatabaseError(err, "检查Logo是否存在失败")
	}

	return count > 0, nil
}

// ExistsByOrganizationID 检查组织是否已有Logo
func (r *logoRepository) ExistsByOrganizationID(
	ctx context.Context,
	organizationID uuid.UUID,
) (bool, error) {
	var count int64

	err := r.DB().WithContext(ctx).
		Model(&models.OrganizationLogo{}).
		Where("bound_organization_id = ? AND status = ?", organizationID, models.LogoStatusBound).
		Count(&count).Error
	if err != nil {
		return false, errno.WrapDatabaseError(err, "检查组织Logo是否存在失败")
	}

	return count > 0, nil
}

// UpdateStatus 更新Logo状态
func (r *logoRepository) UpdateStatus(
	ctx context.Context,
	logoID uuid.UUID,
	status models.OrganizationLogoStatus,
) error {
	err := r.DB().WithContext(ctx).
		Model(&models.OrganizationLogo{}).
		Where("id = ?", logoID).
		Update("status", status).Error
	if err != nil {
		return errno.WrapDatabaseError(err, "更新Logo状态失败")
	}

	return nil
}

// BindToOrganization 绑定Logo到组织（临时→永久）
func (r *logoRepository) BindToOrganization(
	ctx context.Context,
	logoID uuid.UUID,
	organizationID uuid.UUID,
) error {
	updates := map[string]interface{}{
		"status":                models.LogoStatusBound,
		"bound_organization_id": organizationID,
		"expires_at":            nil, // 清除过期时间
	}

	err := r.DB().WithContext(ctx).
		Model(&models.OrganizationLogo{}).
		Where("id = ?", logoID).
		Updates(updates).Error
	if err != nil {
		return errno.WrapDatabaseError(err, "绑定Logo到组织失败")
	}

	return nil
}

// ListExpiredTemporaryLogos 获取已过期的临时Logo列表
func (r *logoRepository) ListExpiredTemporaryLogos(
	ctx context.Context,
) ([]*models.OrganizationLogo, error) {
	var logos []*models.OrganizationLogo

	now := time.Now().UnixMilli()

	err := r.DB().WithContext(ctx).
		Where("status = ? AND expires_at IS NOT NULL AND expires_at <= ?", models.LogoStatusTemporary, now).
		Find(&logos).Error
	if err != nil {
		return nil, errno.WrapDatabaseError(err, "查询过期Logo失败")
	}

	return logos, nil
}

// CleanupExpiredLogos 清理已过期的临时Logo（软删除）
func (r *logoRepository) CleanupExpiredLogos(ctx context.Context) (int64, error) {
	now := time.Now().UnixMilli()

	result := r.DB().WithContext(ctx).
		Where("status = ? AND expires_at IS NOT NULL AND expires_at <= ?", models.LogoStatusTemporary, now).
		Delete(&models.OrganizationLogo{})

	if result.Error != nil {
		return 0, errno.WrapDatabaseError(result.Error, "清理过期Logo失败")
	}

	return result.RowsAffected, nil
}
