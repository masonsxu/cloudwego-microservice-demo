package logo

import (
	"context"

	"github.com/google/uuid"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// LogoRepository 组织Logo仓储接口
// 提供组织Logo的数据访问能力
type LogoRepository interface {
	// 嵌入基础仓储接口
	base.BaseRepository[models.OrganizationLogo]

	// ============================================================================
	// Logo查询
	// ============================================================================

	// GetByFileID 根据文件ID获取Logo
	GetByFileID(ctx context.Context, fileID string) (*models.OrganizationLogo, error)

	// GetByOrganizationID 根据组织ID获取绑定的Logo
	GetByOrganizationID(
		ctx context.Context,
		organizationID uuid.UUID,
	) (*models.OrganizationLogo, error)

	// ListByStatus 根据状态获取Logo列表
	ListByStatus(
		ctx context.Context,
		status models.OrganizationLogoStatus,
		opts *base.QueryOptions,
	) ([]*models.OrganizationLogo, *models.PageResult, error)

	// ListByUploader 根据上传者ID获取Logo列表
	ListByUploader(
		ctx context.Context,
		uploaderID uuid.UUID,
		opts *base.QueryOptions,
	) ([]*models.OrganizationLogo, *models.PageResult, error)

	// ============================================================================
	// 存在性检查
	// ============================================================================

	// ExistsByFileID 检查文件ID是否存在
	ExistsByFileID(ctx context.Context, fileID string) (bool, error)

	// ExistsByOrganizationID 检查组织是否已有Logo
	ExistsByOrganizationID(ctx context.Context, organizationID uuid.UUID) (bool, error)

	// ============================================================================
	// 状态管理
	// ============================================================================

	// UpdateStatus 更新Logo状态
	UpdateStatus(ctx context.Context, logoID uuid.UUID, status models.OrganizationLogoStatus) error

	// BindToOrganization 绑定Logo到组织（临时→永久）
	BindToOrganization(ctx context.Context, logoID uuid.UUID, organizationID uuid.UUID) error

	// ============================================================================
	// 生命周期管理
	// ============================================================================

	// ListExpiredTemporaryLogos 获取已过期的临时Logo列表
	ListExpiredTemporaryLogos(ctx context.Context) ([]*models.OrganizationLogo, error)

	// CleanupExpiredLogos 清理已过期的临时Logo（软删除）
	CleanupExpiredLogos(ctx context.Context) (int64, error)
}
