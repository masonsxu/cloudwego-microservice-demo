package authentication

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// AuthenticationLogic 用户身份认证业务逻辑接口
type AuthenticationLogic interface {
	// Login
	Login(
		ctx context.Context,
		req *identity_srv.LoginRequest,
	) (*identity_srv.LoginResponse, error)

	// ============================================================================
	// 认证和安全
	// ============================================================================

	// ChangePassword 修改用户密码
	ChangePassword(ctx context.Context, req *identity_srv.ChangePasswordRequest) error

	// ResetPassword 重置用户密码
	ResetPassword(ctx context.Context, req *identity_srv.ResetPasswordRequest) error

	// ForcePasswordChange 强制用户修改密码
	ForcePasswordChange(ctx context.Context, req *identity_srv.ForcePasswordChangeRequest) error
}
