package organization

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// OrganizationLogic 组织管理业务逻辑接口
// 负责机构和组织的层级结构管理，包括创建、更新、查询和关系维护
type OrganizationLogic interface {
	// ============================================================================
	// 组织基础操作
	// ============================================================================

	// CreateOrganization 创建组织
	CreateOrganization(
		ctx context.Context,
		req *identity_srv.CreateOrganizationRequest,
	) (*identity_srv.Organization, error)

	// GetOrganization 根据ID获取组织信息
	GetOrganization(
		ctx context.Context,
		req *identity_srv.GetOrganizationRequest,
	) (*identity_srv.Organization, error)

	// UpdateOrganization 更新组织信息
	UpdateOrganization(
		ctx context.Context,
		req *identity_srv.UpdateOrganizationRequest,
	) (*identity_srv.Organization, error)

	// DeleteOrganization 删除组织（软删除）
	DeleteOrganization(ctx context.Context, organizationID string) error

	// ============================================================================
	// 组织查询操作
	// ============================================================================

	// ListOrganizations 分页查询组织列表
	ListOrganizations(
		ctx context.Context,
		req *identity_srv.ListOrganizationsRequest,
	) (*identity_srv.ListOrganizationsResponse, error)
}
