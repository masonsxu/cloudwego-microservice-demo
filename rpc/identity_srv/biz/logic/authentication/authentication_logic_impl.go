package authentication

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/convutil"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	membershipDAL "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/membership"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/menu"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
	tracelog "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/log"
)

// LogicImpl 用户认证逻辑实现
type LogicImpl struct {
	dal       dal.DAL
	converter converter.Converter
	menuLogic menu.MenuLogic
}

// NewLogic 创建用户认证逻辑实现
func NewLogic(
	dal dal.DAL,
	converter converter.Converter,
	menuLogic menu.MenuLogic,
) AuthenticationLogic {
	return &LogicImpl{
		dal:       dal,
		converter: converter,
		menuLogic: menuLogic,
	}
}

// ============================================================================
// 认证和安全
// ============================================================================

// Login 用户登录
func (l *LogicImpl) Login(
	ctx context.Context,
	req *identity_srv.LoginRequest,
) (*identity_srv.LoginResponse, error) {
	// 根据用户名获取用户档案
	userProfile, err := l.dal.UserProfile().GetByUsername(ctx, *req.Username)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return nil, errno.ErrUserNotFound
		}

		return nil, errno.ErrOperationFailed.WithMessage("获取用户档案失败: " + err.Error())
	}

	// 验证密码
	if !convutil.VerifyPassword(*req.Password, userProfile.PasswordHash) {
		// 增加登录失败次数
		_ = l.dal.UserProfile().IncrementLoginAttempts(ctx, userProfile.ID.String())
		return nil, errno.ErrInvalidCredentials
	}

	// 检查账户状态
	if userProfile.Status == models.UserStatusInactive {
		return nil, errno.ErrUserInactive
	}

	if userProfile.Status == models.UserStatusSuspended {
		return nil, errno.ErrUserSuspended
	}

	// 检查是否需要强制修改密码
	if userProfile.MustChangePassword {
		return nil, errno.ErrMustChangePassword
	}

	// 获取用户的成员关系
	userID := userProfile.ID.String()
	// 获取用户的活跃成员关系
	activeStatus := models.MembershipStatusActive
	conditions := &membershipDAL.UserMembershipQueryConditions{
		UserID: &userID,
		Status: &activeStatus,
	}

	memberships, _, err := l.dal.UserMembership().FindWithConditions(ctx, conditions)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("获取用户成员关系失败: " + err.Error())
	}

	// 更新最后登录时间并重置登录失败次数
	err = l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		if err := txDAL.UserProfile().UpdateLastLoginTime(ctx, userProfile.ID.String()); err != nil {
			return err
		}

		if err := txDAL.UserProfile().ResetLoginAttempts(ctx, userProfile.ID.String()); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		// 记录错误但不影响登录流程
		tracelog.Ctx(ctx).Warn().
			Err(err).
			Str("user_id", userProfile.ID.String()).
			Str("username", userProfile.Username).
			Msg("更新登录时间或重置登录失败次数失败")
	}

	// 构建登录响应
	resp := l.converter.BuildLoginResponse(userProfile, memberships)

	// 获取用户菜单树和权限信息
	menuResp, err := l.menuLogic.GetUserMenuTree(ctx, &identity_srv.GetUserMenuTreeRequest{
		UserID: &userID,
	})
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("获取用户权限失败: " + err.Error())
	}

	// 检查用户是否有活跃角色
	if len(menuResp.RoleIDs) == 0 {
		return nil, errno.ErrNoActiveRoles.WithMessage("用户没有可用的角色，无法登录")
	}

	// 获取用户菜单权限
	permissions, err := l.menuLogic.GetUserMenuPermissions(
		ctx,
		&identity_srv.GetUserMenuPermissionsRequest{
			UserID: &userID,
		},
	)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("获取用户权限失败: " + err.Error())
	}

	resp.MenuTree = menuResp.MenuTree
	resp.RoleIDs = menuResp.RoleIDs
	resp.Permissions = permissions.Permissions

	return resp, nil
}

// ChangePassword 修改用户密码
func (l *LogicImpl) ChangePassword(
	ctx context.Context,
	req *identity_srv.ChangePasswordRequest,
) error {
	if req.UserID == nil {
		return errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}

	if req.NewPassword_ == nil || *req.NewPassword_ == "" {
		return errno.ErrInvalidParams.WithMessage("新密码不能为空")
	}

	// 获取用户档案
	profile, err := l.dal.UserProfile().GetByID(ctx, *req.UserID)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return errno.ErrUserNotFound
		}

		return errno.ErrOperationFailed.WithMessage("获取用户档案失败: " + err.Error())
	}

	// 验证旧密码
	if !convutil.VerifyPassword(*req.OldPassword, profile.PasswordHash) {
		return errno.ErrInvalidPassword
	}

	// 生成新密码哈希
	newPasswordHash, err := convutil.HashPassword(*req.NewPassword_)
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("密码哈希生成失败: " + err.Error())
	}

	// 更新密码
	err = l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		return txDAL.UserProfile().UpdatePassword(ctx, *req.UserID, newPasswordHash)
	})
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("更新密码失败: " + err.Error())
	}

	return nil
}

// ResetPassword 重置用户密码
func (l *LogicImpl) ResetPassword(
	ctx context.Context,
	req *identity_srv.ResetPasswordRequest,
) error {
	if req.UserID == nil {
		return errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}

	if req.NewPassword_ == nil || *req.NewPassword_ == "" {
		return errno.ErrInvalidParams.WithMessage("新密码不能为空")
	}

	// 生成新密码哈希
	newPasswordHash, err := convutil.HashPassword(*req.NewPassword_)
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("密码哈希生成失败: " + err.Error())
	}

	// 重置密码
	err = l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		return txDAL.UserProfile().UpdatePassword(ctx, *req.UserID, newPasswordHash)
	})
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("重置密码失败: " + err.Error())
	}

	return nil
}

// ForcePasswordChange 强制用户修改密码
func (l *LogicImpl) ForcePasswordChange(
	ctx context.Context,
	req *identity_srv.ForcePasswordChangeRequest,
) error {
	if req.UserID == nil {
		return errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}

	return l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		return txDAL.UserProfile().SetMustChangePassword(ctx, *req.UserID, true)
	})
}
