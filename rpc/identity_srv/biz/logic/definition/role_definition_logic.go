package definition

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// RoleDefinitionLogic 角色定义管理业务逻辑接口
// 负责角色定义的创建、更新、查询、删除等核心业务功能
type RoleDefinitionLogic interface {
	// ============================================================================
	// 角色定义生命周期管理
	// ============================================================================

	// CreateRoleDefinition 创建一个新的角色定义
	// 对应 IDL: CreateRoleDefinition(1: RoleDefinitionCreateRequest req)
	CreateRoleDefinition(
		ctx context.Context,
		req *identity_srv.RoleDefinitionCreateRequest,
	) (*identity_srv.RoleDefinition, error)

	// UpdateRoleDefinition 更新一个已有的角色定义
	// 对应 IDL: UpdateRoleDefinition(1: RoleDefinitionUpdateRequest req)
	UpdateRoleDefinition(
		ctx context.Context,
		req *identity_srv.RoleDefinitionUpdateRequest,
	) (*identity_srv.RoleDefinition, error)

	// DeleteRoleDefinition 删除一个角色定义
	// 对应 IDL: DeleteRoleDefinition(1: core.UUID roleID)
	DeleteRoleDefinition(ctx context.Context, roleID string) error

	// GetRoleDefinition 根据ID获取角色定义
	// 对应 IDL: GetRoleDefinition(1: core.UUID roleID)
	GetRoleDefinition(ctx context.Context, roleID string) (*identity_srv.RoleDefinition, error)

	// ListRoleDefinitions 分页列出角色定义
	// 对应 IDL: ListRoleDefinitions(1: RoleDefinitionQueryRequest req)
	ListRoleDefinitions(
		ctx context.Context,
		req *identity_srv.RoleDefinitionQueryRequest,
	) (*identity_srv.RoleDefinitionListResponse, error)
}
