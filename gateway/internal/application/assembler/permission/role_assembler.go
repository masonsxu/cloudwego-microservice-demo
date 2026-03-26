package permission

import (
	permissionModel "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/common"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// roleAssembler 角色定义转换器实现
type roleAssembler struct {
	permissionAssembler IPermissionAssembler
}

// NewRoleAssembler 创建角色定义转换器
func NewRoleAssembler(permissionAssembler IPermissionAssembler) IRoleAssembler {
	return &roleAssembler{
		permissionAssembler: permissionAssembler,
	}
}

// ToHTTPRoleDefinition 将RPC角色定义转换为HTTP角色定义DTO
func (a *roleAssembler) ToHTTPRoleDefinition(
	rpc *identity_srv.RoleDefinition,
) *permissionModel.RoleDefinitionDTO {
	if rpc == nil {
		return nil
	}

	return &permissionModel.RoleDefinitionDTO{
		Id:           rpc.Id,
		Name:         rpc.Name,
		Description:  common.CopyStringPtr(rpc.Description),
		Status:       common.ConvertRoleStatusPtrToHTTPPtr(rpc.Status),
		Permissions:  a.permissionAssembler.ToHTTPPermissions(rpc.Permissions),
		IsSystemRole: common.CopyBoolPtr(rpc.IsSystemRole),
		CreatedBy:    common.CopyStringPtr(rpc.CreatedBy),
		UpdatedBy:    common.CopyStringPtr(rpc.UpdatedBy),
		CreatedAt:    common.CopyInt64Ptr(rpc.CreatedAt),
		UpdatedAt:    common.CopyInt64Ptr(rpc.UpdatedAt),
		UserCount:    common.CopyInt64Ptr(rpc.UserCount),
	}
}

// ToHTTPRoleDefinitions 将RPC角色定义列表转换为HTTP角色定义DTO列表
func (a *roleAssembler) ToHTTPRoleDefinitions(
	rpc []*identity_srv.RoleDefinition,
) []*permissionModel.RoleDefinitionDTO {
	if rpc == nil {
		return nil
	}

	result := make([]*permissionModel.RoleDefinitionDTO, len(rpc))
	for i, r := range rpc {
		result[i] = a.ToHTTPRoleDefinition(r)
	}

	return result
}

// ToRPCRoleDefinition 将HTTP角色定义DTO转换为RPC角色定义
func (a *roleAssembler) ToRPCRoleDefinition(
	http *permissionModel.RoleDefinitionDTO,
) *identity_srv.RoleDefinition {
	if http == nil {
		return nil
	}

	rpc := &identity_srv.RoleDefinition{
		Id:          http.Id,
		Name:        http.Name,
		Description: http.Description,
		Permissions: a.permissionAssembler.ToRPCPermissions(http.Permissions),
		CreatedBy:   http.CreatedBy,
		UpdatedBy:   http.UpdatedBy,
		CreatedAt:   http.CreatedAt,
		UpdatedAt:   http.UpdatedAt,
	}
	if http.IsSystemRole != nil {
		rpc.IsSystemRole = http.IsSystemRole
	}

	return rpc
}

// ToRPCRoleDefinitionCreateRequest 将HTTP角色定义创建请求转换为RPC请求
func (a *roleAssembler) ToRPCRoleDefinitionCreateRequest(
	http *permissionModel.RoleDefinitionCreateRequestDTO,
) *identity_srv.RoleDefinitionCreateRequest {
	if http == nil {
		return nil
	}

	req := &identity_srv.RoleDefinitionCreateRequest{
		Name:        http.Name,
		Description: http.Description,
		Permissions: a.permissionAssembler.ToRPCPermissions(http.Permissions),
	}
	if http.IsSystemRole != nil {
		req.IsSystemRole = http.IsSystemRole
	}

	return req
}

// ToHTTPRoleDefinitionCreateResponse 将RPC角色定义转换为HTTP创建响应
func (a *roleAssembler) ToHTTPRoleDefinitionCreateResponse(
	rpc *identity_srv.RoleDefinition,
) *permissionModel.RoleDefinitionCreateResponseDTO {
	return &permissionModel.RoleDefinitionCreateResponseDTO{
		Role: a.ToHTTPRoleDefinition(rpc),
	}
}

// ToRPCRoleDefinitionUpdateRequest 将HTTP角色定义更新请求转换为RPC请求
func (a *roleAssembler) ToRPCRoleDefinitionUpdateRequest(
	http *permissionModel.RoleDefinitionUpdateRequestDTO,
) *identity_srv.RoleDefinitionUpdateRequest {
	if http == nil {
		return nil
	}

	return &identity_srv.RoleDefinitionUpdateRequest{
		RoleDefinitionID: http.RoleDefinitionID,
		Description:      http.Description,
		Status:           common.ConvertRoleStatusPtrToRPCPtr(http.Status),
		Permissions:      a.permissionAssembler.ToRPCPermissionListValue(http.Permissions),
		Name:             http.Name,
	}
}

// ToHTTPRoleDefinitionUpdateResponse 将RPC角色定义转换为HTTP更新响应
func (a *roleAssembler) ToHTTPRoleDefinitionUpdateResponse(
	rpc *identity_srv.RoleDefinition,
) *permissionModel.RoleDefinitionUpdateResponseDTO {
	return &permissionModel.RoleDefinitionUpdateResponseDTO{
		Role: a.ToHTTPRoleDefinition(rpc),
	}
}

// ToHTTPRoleDefinitionGetResponse 将RPC角色定义转换为HTTP获取响应
func (a *roleAssembler) ToHTTPRoleDefinitionGetResponse(
	rpc *identity_srv.RoleDefinition,
) *permissionModel.RoleDefinitionGetResponseDTO {
	return &permissionModel.RoleDefinitionGetResponseDTO{
		Role: a.ToHTTPRoleDefinition(rpc),
	}
}

// ToRPCRoleDefinitionQueryRequest 将HTTP角色定义查询请求转换为RPC请求
func (a *roleAssembler) ToRPCRoleDefinitionQueryRequest(
	http *permissionModel.RoleDefinitionQueryRequestDTO,
) *identity_srv.RoleDefinitionQueryRequest {
	if http == nil {
		return nil
	}

	return &identity_srv.RoleDefinitionQueryRequest{
		Name:         http.Name,
		Status:       common.ConvertRoleStatusPtrToRPCPtr(http.Status),
		IsSystemRole: http.IsSystemRole,
		Page:         ToRPCPageRequest(http.Page),
	}
}

// ToHTTPRoleDefinitionListResponse 将RPC角色定义列表响应转换为HTTP响应
func (a *roleAssembler) ToHTTPRoleDefinitionListResponse(
	rpc *identity_srv.RoleDefinitionListResponse,
) *permissionModel.RoleDefinitionListResponseDTO {
	if rpc == nil {
		return nil
	}

	return &permissionModel.RoleDefinitionListResponseDTO{
		Roles: a.ToHTTPRoleDefinitions(rpc.Roles),
		Page:  ToHTTPPageResponse(rpc.Page),
	}
}
