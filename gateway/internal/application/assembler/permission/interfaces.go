package permission

import (
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	permissionModel "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/rpc_base"
)

// Assembler 权限管理聚合组装器接口 - 统一暴露给Service层
// 提供所有权限相关的协议转换能力，避免Service层直接依赖多个细分Assembler
type Assembler interface {
	// 获取各个业务领域的组装器
	Role() IRoleAssembler
	Permission() IPermissionAssembler
	UserRole() IUserRoleAssembler
	Menu() IMenuAssembler

	// 通用转换方法（避免重复代码）
	ToHTTPPageResponse(*rpc_base.PageResponse) *http_base.PageResponseDTO
	ToRPCPageRequest(*http_base.PageRequestDTO) *rpc_base.PageRequest
}

// IMenuAssembler 菜单相关组装器接口
type IMenuAssembler interface {
	// 菜单相关转换
	ToHTTPMenuNode(*identity_srv.MenuNode) *permissionModel.MenuNodeDTO
	ToHTTPMenuNodes([]*identity_srv.MenuNode) []*permissionModel.MenuNodeDTO
	ToRPCMenuNode(*permissionModel.MenuNodeDTO) *identity_srv.MenuNode
	ToRPCMenuNodes([]*permissionModel.MenuNodeDTO) []*identity_srv.MenuNode

	ToRPCUploadMenuRequest(*permissionModel.UploadMenuRequestDTO) *identity_srv.UploadMenuRequest

	ToHTTPGetMenuTreeResponse(
		*identity_srv.GetMenuTreeResponse,
	) *permissionModel.GetMenuTreeResponseDTO

	// 菜单权限管理相关转换
	// 菜单配置转换
	ToHTTPMenuConfig(*identity_srv.MenuConfig) *permissionModel.MenuConfigDTO
	ToHTTPMenuConfigs([]*identity_srv.MenuConfig) []*permissionModel.MenuConfigDTO
	ToRPCMenuConfig(*permissionModel.MenuConfigDTO) *identity_srv.MenuConfig
	ToRPCMenuConfigs([]*permissionModel.MenuConfigDTO) []*identity_srv.MenuConfig

	// 菜单权限转换
	ToHTTPMenuPermission(*identity_srv.MenuPermission) *permissionModel.MenuPermissionDTO
	ToHTTPMenuPermissions([]*identity_srv.MenuPermission) []*permissionModel.MenuPermissionDTO

	// 配置角色菜单权限转换
	ToRPCConfigureRoleMenusRequest(
		operatorID string,
		req *permissionModel.ConfigureRoleMenusRequestDTO,
	) *identity_srv.ConfigureRoleMenusRequest
	ToHTTPConfigureRoleMenusResponse(
		*identity_srv.ConfigureRoleMenusResponse,
	) *permissionModel.ConfigureRoleMenusResponseDTO

	// 获取角色菜单树转换
	ToRPCGetRoleMenuTreeRequest(
		req *permissionModel.GetRoleMenuTreeRequestDTO,
	) *identity_srv.GetRoleMenuTreeRequest
	ToHTTPGetRoleMenuTreeResponse(
		*identity_srv.GetRoleMenuTreeResponse,
	) *permissionModel.GetRoleMenuTreeResponseDTO

	// 获取用户菜单树转换
	ToRPCGetUserMenuTreeRequest(
		req *permissionModel.GetUserMenuTreeRequestDTO,
	) *identity_srv.GetUserMenuTreeRequest
	ToHTTPGetUserMenuTreeResponse(
		*identity_srv.GetUserMenuTreeResponse,
	) *permissionModel.GetUserMenuTreeResponseDTO

	// 获取角色菜单权限转换
	ToRPCGetRoleMenuPermissionsRequest(
		req *permissionModel.GetRoleMenuPermissionsRequestDTO,
	) *identity_srv.GetRoleMenuPermissionsRequest
	ToHTTPGetRoleMenuPermissionsResponse(
		*identity_srv.GetRoleMenuPermissionsResponse,
	) *permissionModel.GetRoleMenuPermissionsResponseDTO

	// 检查菜单权限转换
	ToRPCHasMenuPermissionRequest(
		req *permissionModel.HasMenuPermissionRequestDTO,
	) *identity_srv.HasMenuPermissionRequest
	ToHTTPHasMenuPermissionResponse(
		*identity_srv.HasMenuPermissionResponse,
	) *permissionModel.HasMenuPermissionResponseDTO
}

// IPermissionAssembler 权限相关组装器接口
type IPermissionAssembler interface {
	// 基础权限结构转换
	ToHTTPPermission(*identity_srv.Permission) *permissionModel.PermissionDTO
	ToHTTPPermissions([]*identity_srv.Permission) []*permissionModel.PermissionDTO
	ToRPCPermission(*permissionModel.PermissionDTO) *identity_srv.Permission
	ToRPCPermissions([]*permissionModel.PermissionDTO) []*identity_srv.Permission
}

// IRoleAssembler 角色定义相关组装器接口
type IRoleAssembler interface {
	// 基础角色定义结构转换
	ToHTTPRoleDefinition(*identity_srv.RoleDefinition) *permissionModel.RoleDefinitionDTO
	ToHTTPRoleDefinitions([]*identity_srv.RoleDefinition) []*permissionModel.RoleDefinitionDTO
	ToRPCRoleDefinition(*permissionModel.RoleDefinitionDTO) *identity_srv.RoleDefinition

	// 创建角色定义转换
	ToRPCRoleDefinitionCreateRequest(
		*permissionModel.RoleDefinitionCreateRequestDTO,
	) *identity_srv.RoleDefinitionCreateRequest
	ToHTTPRoleDefinitionCreateResponse(
		*identity_srv.RoleDefinition,
	) *permissionModel.RoleDefinitionCreateResponseDTO

	// 更新角色定义转换
	ToRPCRoleDefinitionUpdateRequest(
		*permissionModel.RoleDefinitionUpdateRequestDTO,
	) *identity_srv.RoleDefinitionUpdateRequest
	ToHTTPRoleDefinitionUpdateResponse(
		*identity_srv.RoleDefinition,
	) *permissionModel.RoleDefinitionUpdateResponseDTO

	// 获取角色定义转换
	ToHTTPRoleDefinitionGetResponse(
		*identity_srv.RoleDefinition,
	) *permissionModel.RoleDefinitionGetResponseDTO

	// 查询角色定义转换
	ToRPCRoleDefinitionQueryRequest(
		*permissionModel.RoleDefinitionQueryRequestDTO,
	) *identity_srv.RoleDefinitionQueryRequest
	ToHTTPRoleDefinitionListResponse(
		*identity_srv.RoleDefinitionListResponse,
	) *permissionModel.RoleDefinitionListResponseDTO
}

// IUserRoleAssembler 用户角色分配相关组装器接口
type IUserRoleAssembler interface {
	// 基础用户角色分配结构转换
	ToHTTPUserRoleAssignment(
		*identity_srv.UserRoleAssignment,
	) *permissionModel.UserRoleAssignmentDTO
	ToHTTPUserRoleAssignments(
		[]*identity_srv.UserRoleAssignment,
	) []*permissionModel.UserRoleAssignmentDTO
	ToRPCUserRoleAssignment(
		*permissionModel.UserRoleAssignmentDTO,
	) *identity_srv.UserRoleAssignment
	ToHTTPAssignRoleToUserResponse(
		*identity_srv.UserRoleAssignmentResponse,
	) *permissionModel.AssignRoleToUserResponseDTO

	// 获取用户最后角色分配转换
	ToHTTPGetLastUserRoleAssignmentResponse(
		*identity_srv.UserRoleAssignment,
	) *permissionModel.GetLastUserRoleAssignmentResponseDTO

	// 查询用户角色分配转换
	ToRPCUserRoleQueryRequest(
		*permissionModel.UserRoleQueryRequestDTO,
	) *identity_srv.UserRoleQueryRequest
	ToHTTPUserRoleListResponse(
		*identity_srv.UserRoleListResponse,
	) *permissionModel.UserRoleListResponseDTO

	// 根据角色ID获取用户转换
	ToRPCGetUsersByRoleRequest(
		*permissionModel.GetUsersByRoleRequestDTO,
	) *identity_srv.GetUsersByRoleRequest
	ToHTTPGetUsersByRoleResponse(
		*identity_srv.GetUsersByRoleResponse,
	) *permissionModel.GetUsersByRoleResponseDTO

	// 批量绑定用户到角色转换
	ToRPCBatchBindUsersToRoleRequest(
		operatorID string,
		req *permissionModel.BatchBindUsersToRoleRequestDTO,
	) *identity_srv.BatchBindUsersToRoleRequest
	ToHTTPBatchBindUsersToRoleResponse(
		*identity_srv.BatchBindUsersToRoleResponse,
	) *permissionModel.BatchBindUsersToRoleResponseDTO
}

// ToHTTPPageResponse is a generic function to convert RPC PageResponse to HTTP PageResponseDTO.
func ToHTTPPageResponse(rpc *rpc_base.PageResponse) *http_base.PageResponseDTO {
	if rpc == nil {
		return nil
	}

	return &http_base.PageResponseDTO{
		Total:      rpc.Total,
		Page:       rpc.Page,
		Limit:      rpc.Limit,
		TotalPages: rpc.TotalPages,
		HasNext:    rpc.HasNext,
		HasPrev:    rpc.HasPrev,
	}
}

// ToRPCPageRequest is a generic function to convert HTTP PageRequestDTO to RPC PageRequest.
func ToRPCPageRequest(http *http_base.PageRequestDTO) *rpc_base.PageRequest {
	if http == nil {
		return nil
	}

	return &rpc_base.PageRequest{
		Page:         http.Page,
		Limit:        http.Limit,
		Search:       http.Search,
		Filter:       http.Filter,
		Sort:         http.Sort,
		Fields:       http.Fields,
		IncludeTotal: http.IncludeTotal,
		FetchAll:     http.FetchAll,
	}
}
