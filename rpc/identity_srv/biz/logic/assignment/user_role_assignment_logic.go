package assignment

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// RoleAssignmentLogic 用户角色分配管理业务逻辑接口
// 负责用户角色分配的创建、更新、查询、删除等核心业务功能
type RoleAssignmentLogic interface {
	// ============================================================================
	// 用户角色分配管理
	// ============================================================================

	// AssignRoleToUser 为用户分配一个角色
	// 对应 IDL: AssignRoleToUser(1: AssignRoleToUserRequest req)
	AssignRoleToUser(
		ctx context.Context,
		req *identity_srv.AssignRoleToUserRequest,
	) (*identity_srv.UserRoleAssignmentResponse, error)

	// UpdateUserRoleAssignment 更新用户的角色分配信息
	// 对应 IDL: UpdateUserRoleAssignment(1: UpdateUserRoleAssignmentRequest req)
	UpdateUserRoleAssignment(
		ctx context.Context,
		req *identity_srv.UpdateUserRoleAssignmentRequest,
	) error

	// RevokeRoleFromUser 撤销用户的角色分配
	// 对应 IDL: RevokeRoleFromUser(1: RevokeRoleFromUserRequest req)
	RevokeRoleFromUser(
		ctx context.Context,
		req *identity_srv.RevokeRoleFromUserRequest,
	) error

	// GetLastUserRoleAssignment 获取用户最后一次的角色分配信息
	// 对应 IDL: GetLastUserRoleAssignment(1: core.UUID userID)
	GetLastUserRoleAssignment(
		ctx context.Context,
		userID string,
	) (*identity_srv.UserRoleAssignment, error)

	// ListUserRoleAssignments 列出用户的角色分配记录
	// 对应 IDL: ListUserRoleAssignments(1: UserRoleQueryRequest req)
	ListUserRoleAssignments(
		ctx context.Context,
		req *identity_srv.UserRoleQueryRequest,
	) (*identity_srv.UserRoleListResponse, error)

	// GetUsersByRole 根据角色ID获取该角色下所有用户
	// 对应 IDL: GetUsersByRole(1: GetUsersByRoleRequest req)
	GetUsersByRole(
		ctx context.Context,
		req *identity_srv.GetUsersByRoleRequest,
	) (*identity_srv.GetUsersByRoleResponse, error)

	// BatchBindUsersToRole 批量绑定用户到角色
	// 对应 IDL: BatchBindUsersToRole(1: BatchBindUsersToRoleRequest req)
	BatchBindUsersToRole(
		ctx context.Context,
		req *identity_srv.BatchBindUsersToRoleRequest,
	) (*identity_srv.BatchBindUsersToRoleResponse, error)

	// BatchGetUserRoles 批量获取多个用户的角色分配
	// 对应 IDL: BatchGetUserRoles(1: BatchGetUserRolesRequest req)
	BatchGetUserRoles(
		ctx context.Context,
		req *identity_srv.BatchGetUserRolesRequest,
	) (*identity_srv.BatchGetUserRolesResponse, error)
}
