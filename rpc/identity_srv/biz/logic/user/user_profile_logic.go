package user

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// ProfileLogic 用户档案管理业务逻辑接口
// 负责用户个人信息、认证状态、资质等核心档案数据的管理
type ProfileLogic interface {
	// ============================================================================
	// 用户生命周期
	// ============================================================================

	// CreateUser 创建用户
	CreateUser(
		ctx context.Context,
		req *identity_srv.CreateUserRequest,
	) (*identity_srv.UserProfile, error)

	// GetUser 根据用户ID获取用户
	GetUser(
		ctx context.Context,
		req *identity_srv.GetUserRequest,
	) (*identity_srv.UserProfile, error)

	// UpdateUser 更新用户信息
	UpdateUser(
		ctx context.Context,
		req *identity_srv.UpdateUserRequest,
	) (*identity_srv.UserProfile, error)

	// DeleteUser 删除用户（软删除）
	DeleteUser(ctx context.Context, req *identity_srv.DeleteUserRequest) error

	// ListUsers 分页查询用户列表
	ListUsers(
		ctx context.Context,
		req *identity_srv.ListUsersRequest,
	) (*identity_srv.ListUsersResponse, error)

	// SearchUsers 按条件搜索用户
	SearchUsers(
		ctx context.Context,
		req *identity_srv.SearchUsersRequest,
	) (*identity_srv.SearchUsersResponse, error)

	// ============================================================================
	// 用户状态管理
	// ============================================================================

	// ChangeUserStatus 激活、停用、锁定用户
	ChangeUserStatus(ctx context.Context, req *identity_srv.ChangeUserStatusRequest) error

	// UnlockUser 解锁用户
	UnlockUser(ctx context.Context, req *identity_srv.UnlockUserRequest) error
}
