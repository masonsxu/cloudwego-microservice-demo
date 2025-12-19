package organization

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Converter 组织转换器接口，负责组织模型的转换
// 基于重构后的 IDL 设计，支持以下实体的双向转换：
// - Organization: 组织实体（支持层级结构和机构属性）
type Converter interface {
	// ============================================================================
	// Organization 组织转换
	// ============================================================================

	// 组织：models -> Thrift
	ModelOrganizationToThrift(*models.Organization) *identity_srv.Organization
	ThriftOrganizationToModel(*identity_srv.Organization) *models.Organization
	ModelOrganizationsToThrift([]*models.Organization) []*identity_srv.Organization

	// 新增的Logic层需要的方法
	ModelToThrift(org *models.Organization) *identity_srv.Organization
	CreateRequestToModel(req *identity_srv.CreateOrganizationRequest) *models.Organization
	ApplyUpdateToModel(
		existing *models.Organization,
		req *identity_srv.UpdateOrganizationRequest,
	) *models.Organization
}
