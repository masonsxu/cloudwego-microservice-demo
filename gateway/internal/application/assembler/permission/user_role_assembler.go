package permission

import (
	permissionModel "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/common"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// userRoleAssembler 用户角色分配转换器实现
type userRoleAssembler struct{}

// NewUserRoleAssembler 创建用户角色分配转换器
func NewUserRoleAssembler() IUserRoleAssembler {
	return &userRoleAssembler{}
}

// ToHTTPUserRoleAssignment 将RPC用户角色分配转换为HTTP DTO
func (a *userRoleAssembler) ToHTTPUserRoleAssignment(
	rpc *identity_srv.UserRoleAssignment,
) *permissionModel.UserRoleAssignmentDTO {
	if rpc == nil {
		return nil
	}

	return &permissionModel.UserRoleAssignmentDTO{
		ID:        rpc.Id,
		UserID:    rpc.UserID,
		RoleID:    rpc.RoleID,
		CreatedBy: common.CopyStringPtr(rpc.CreatedBy),
		UpdatedBy: common.CopyStringPtr(rpc.UpdatedBy),
		CreatedAt: common.CopyInt64Ptr(rpc.CreatedAt),
		UpdatedAt: common.CopyInt64Ptr(rpc.UpdatedAt),
	}
}

// ToHTTPUserRoleAssignments 将RPC用户角色分配列表转换为HTTP DTO列表
func (a *userRoleAssembler) ToHTTPUserRoleAssignments(
	rpc []*identity_srv.UserRoleAssignment,
) []*permissionModel.UserRoleAssignmentDTO {
	if rpc == nil {
		return nil
	}

	result := make([]*permissionModel.UserRoleAssignmentDTO, len(rpc))
	for i, r := range rpc {
		result[i] = a.ToHTTPUserRoleAssignment(r)
	}

	return result
}

// ToRPCUserRoleAssignment 将HTTP用户角色分配DTO转换为RPC
func (a *userRoleAssembler) ToRPCUserRoleAssignment(
	http *permissionModel.UserRoleAssignmentDTO,
) *identity_srv.UserRoleAssignment {
	if http == nil {
		return nil
	}

	return &identity_srv.UserRoleAssignment{
		Id:        http.ID,
		UserID:    http.UserID,
		RoleID:    http.RoleID,
		CreatedBy: http.CreatedBy,
		UpdatedBy: http.UpdatedBy,
		CreatedAt: http.CreatedAt,
		UpdatedAt: http.UpdatedAt,
	}
}

// ToHTTPAssignRoleToUserResponse 将RPC分配响应转换为HTTP响应
func (a *userRoleAssembler) ToHTTPAssignRoleToUserResponse(
	rpc *identity_srv.UserRoleAssignmentResponse,
) *permissionModel.AssignRoleToUserResponseDTO {
	if rpc == nil {
		return nil
	}

	return &permissionModel.AssignRoleToUserResponseDTO{
		AssignmentID: rpc.AssignmentID,
	}
}

// ToHTTPGetLastUserRoleAssignmentResponse 将RPC用户角色分配转换为HTTP最后分配响应
func (a *userRoleAssembler) ToHTTPGetLastUserRoleAssignmentResponse(
	rpc *identity_srv.UserRoleAssignment,
) *permissionModel.GetLastUserRoleAssignmentResponseDTO {
	return &permissionModel.GetLastUserRoleAssignmentResponseDTO{
		Assignment: a.ToHTTPUserRoleAssignment(rpc),
	}
}

// ToRPCUserRoleQueryRequest 将HTTP用户角色查询请求转换为RPC请求
func (a *userRoleAssembler) ToRPCUserRoleQueryRequest(
	http *permissionModel.UserRoleQueryRequestDTO,
) *identity_srv.UserRoleQueryRequest {
	if http == nil {
		return nil
	}

	return &identity_srv.UserRoleQueryRequest{
		UserID: http.UserID,
		RoleID: http.RoleID,
		Page:   ToRPCPageRequest(http.Page),
	}
}

// ToHTTPUserRoleListResponse 将RPC用户角色列表响应转换为HTTP响应
func (a *userRoleAssembler) ToHTTPUserRoleListResponse(
	rpc *identity_srv.UserRoleListResponse,
) *permissionModel.UserRoleListResponseDTO {
	if rpc == nil {
		return nil
	}

	return &permissionModel.UserRoleListResponseDTO{
		Assignments: a.ToHTTPUserRoleAssignments(rpc.Assignments),
		Page:        ToHTTPPageResponse(rpc.Page),
	}
}

// ToRPCGetUsersByRoleRequest 将HTTP获取角色用户请求转换为RPC请求
func (a *userRoleAssembler) ToRPCGetUsersByRoleRequest(
	http *permissionModel.GetUsersByRoleRequestDTO,
) *identity_srv.GetUsersByRoleRequest {
	if http == nil {
		return nil
	}

	return &identity_srv.GetUsersByRoleRequest{
		RoleID: http.RoleID,
	}
}

// ToHTTPGetUsersByRoleResponse 将RPC获取角色用户响应转换为HTTP响应
func (a *userRoleAssembler) ToHTTPGetUsersByRoleResponse(
	rpc *identity_srv.GetUsersByRoleResponse,
) *permissionModel.GetUsersByRoleResponseDTO {
	if rpc == nil {
		return nil
	}

	return &permissionModel.GetUsersByRoleResponseDTO{
		RoleID:  rpc.RoleID,
		UserIDs: rpc.UserIDs,
	}
}

// ToRPCBatchBindUsersToRoleRequest 将HTTP批量绑定用户请求转换为RPC请求
func (a *userRoleAssembler) ToRPCBatchBindUsersToRoleRequest(
	operatorID string,
	http *permissionModel.BatchBindUsersToRoleRequestDTO,
) *identity_srv.BatchBindUsersToRoleRequest {
	if http == nil {
		return nil
	}

	return &identity_srv.BatchBindUsersToRoleRequest{
		RoleID:     http.RoleID,
		UserIDs:    http.UserIDs,
		OperatorID: &operatorID,
	}
}

// ToHTTPBatchBindUsersToRoleResponse 将RPC批量绑定用户响应转换为HTTP响应
func (a *userRoleAssembler) ToHTTPBatchBindUsersToRoleResponse(
	rpc *identity_srv.BatchBindUsersToRoleResponse,
) *permissionModel.BatchBindUsersToRoleResponseDTO {
	if rpc == nil {
		return nil
	}

	return &permissionModel.BatchBindUsersToRoleResponseDTO{
		Success:      rpc.Success,
		SuccessCount: rpc.SuccessCount,
		Message:      rpc.Message,
	}
}
