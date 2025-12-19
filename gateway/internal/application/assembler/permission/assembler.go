package permission

import (
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/rpc_base"
)

// ============================================================================
// 聚合组装器实现 - 统一访问入口
// ============================================================================

// permissionAggregatedAssembler 权限管理聚合组装器实现
type permissionAggregatedAssembler struct {
	roleAssembler       IRoleAssembler
	permissionAssembler IPermissionAssembler
	userRoleAssembler   IUserRoleAssembler
	menuAssembler       IMenuAssembler
}

// NewPermissionAssembler 创建权限管理聚合组装器
func NewPermissionAggregateAssembler(
	roleAssembler IRoleAssembler,
	permissionAssembler IPermissionAssembler,
	userRoleAssembler IUserRoleAssembler,
	menuAssembler IMenuAssembler,
) Assembler {
	return &permissionAggregatedAssembler{
		roleAssembler:       roleAssembler,
		permissionAssembler: permissionAssembler,
		userRoleAssembler:   userRoleAssembler,
		menuAssembler:       menuAssembler,
	}
}

// 获取各个业务领域的组装器
func (a *permissionAggregatedAssembler) Role() IRoleAssembler { return a.roleAssembler }

func (a *permissionAggregatedAssembler) Permission() IPermissionAssembler {
	return a.permissionAssembler
}

func (a *permissionAggregatedAssembler) UserRole() IUserRoleAssembler { return a.userRoleAssembler }
func (a *permissionAggregatedAssembler) Menu() IMenuAssembler         { return a.menuAssembler }

// 通用转换方法
func (a *permissionAggregatedAssembler) ToHTTPPageResponse(
	rpc *rpc_base.PageResponse,
) *http_base.PageResponseDTO {
	return ToHTTPPageResponse(rpc)
}

func (a *permissionAggregatedAssembler) ToRPCPageRequest(
	http *http_base.PageRequestDTO,
) *rpc_base.PageRequest {
	return ToRPCPageRequest(http)
}
