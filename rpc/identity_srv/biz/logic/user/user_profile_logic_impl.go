package user

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/user"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/rpc_base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
	tracelog "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/log"
)

// LogicImpl 用户档案业务逻辑实现
type LogicImpl struct {
	dal       dal.DAL
	converter converter.Converter
}

// NewLogic 创建用户档案业务逻辑实例
func NewLogic(
	dal dal.DAL,
	converter converter.Converter,
) ProfileLogic {
	return &LogicImpl{
		dal:       dal,
		converter: converter,
	}
}

// ============================================================================
// 用户生命周期
// ============================================================================

// CreateUser 创建用户
func (l *LogicImpl) CreateUser(
	ctx context.Context,
	req *identity_srv.CreateUserRequest,
) (*identity_srv.UserProfile, error) {
	// 参数验证
	if err := l.validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	// 检查用户名是否已存在
	exists, err := l.dal.UserProfile().CheckUsernameExists(ctx, *req.Username)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("检查用户名是否存在失败: " + err.Error())
	}

	if exists {
		return nil, errno.ErrUsernameAlreadyExists
	}

	// 检查邮箱是否已存在
	if req.Email != nil {
		exists, err := l.dal.UserProfile().CheckEmailExists(ctx, *req.Email)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage("检查邮箱是否存在失败: " + err.Error())
		}

		if exists {
			return nil, errno.ErrEmailAlreadyExists
		}
	}

	// 检查手机号是否已存在
	if req.Phone != nil {
		exists, err := l.dal.UserProfile().CheckPhoneExists(ctx, *req.Phone)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage("检查手机号是否存在失败: " + err.Error())
		}

		if exists {
			return nil, errno.ErrPhoneAlreadyExists
		}
	}

	// 转换请求为模型
	userProfile := l.converter.UserProfile().CreateUserRequestToModel(req)

	// 在事务中创建用户档案
	err = l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		if err := txDAL.UserProfile().Create(ctx, userProfile); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	// 转换为响应格式
	userProfileDTO := l.converter.UserProfile().ModelUserProfileToThrift(userProfile)

	// 填充关联字段（主组织、主部门）
	if err := l.enrichUserProfileWithRelations(ctx, userProfileDTO); err != nil {
		// 记录警告但不影响主要结果
		tracelog.Ctx(ctx).Warn().
			Err(err).
			Str("user_id", *userProfileDTO.ID).
			Msg("填充用户关联信息失败")
	}

	return userProfileDTO, nil
}

// GetUser 根据用户ID获取用户
func (l *LogicImpl) GetUser(
	ctx context.Context,
	req *identity_srv.GetUserRequest,
) (*identity_srv.UserProfile, error) {
	if req.UserID == nil {
		return nil, errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}

	profile, err := l.dal.UserProfile().GetByID(ctx, *req.UserID)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return nil, errno.ErrUserNotFound
		}

		return nil, errno.ErrOperationFailed.WithMessage("获取用户档案失败: " + err.Error())
	}

	// 转换为 Thrift DTO
	userProfile := l.converter.UserProfile().ModelUserProfileToThrift(profile)

	// 填充关联字段（主组织、主部门）
	if err := l.enrichUserProfileWithRelations(ctx, userProfile); err != nil {
		// 记录警告但不影响主要结果
		tracelog.Ctx(ctx).Warn().
			Err(err).
			Str("user_id", *req.UserID).
			Msg("填充用户关联信息失败")
	}

	return userProfile, nil
}

// UpdateUser 更新用户信息
func (l *LogicImpl) UpdateUser(
	ctx context.Context,
	req *identity_srv.UpdateUserRequest,
) (*identity_srv.UserProfile, error) {
	if req.UserID == nil {
		return nil, errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}

	// 获取现有档案
	existingProfile, err := l.dal.UserProfile().GetByID(ctx, *req.UserID)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return nil, errno.ErrUserNotFound
		}

		return nil, errno.ErrOperationFailed.WithMessage("获取用户档案失败: " + err.Error())
	}

	// 系统用户保护检查
	if existingProfile.IsSystemUser {
		// 记录系统用户修改操作（用于审计）
		// 注意：UpdateUserRequest 本身不包含 username 和 status 字段，
		// 这些关键属性无法通过此接口修改，已从设计上保护
		tracelog.Ctx(ctx).Info().
			Str("user_id", *req.UserID).
			Str("username", existingProfile.Username).
			Strs("modified_fields", getModifiedFields(req)).
			Msg("修改系统用户信息")
	}

	// 检查唯一性约束
	if err := l.checkUniqueConstraints(ctx, req, existingProfile.ID.String()); err != nil {
		return nil, err
	}

	// 应用更新
	updatedProfile := l.converter.UserProfile().ApplyUpdateUserToModel(existingProfile, req)

	// 在事务中更新
	err = l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		if err := txDAL.UserProfile().Update(ctx, updatedProfile); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	userProfileDTO := l.converter.UserProfile().ModelUserProfileToThrift(updatedProfile)

	// 填充关联字段（主组织、主部门）
	if err := l.enrichUserProfileWithRelations(ctx, userProfileDTO); err != nil {
		// 记录警告但不影响主要结果
		tracelog.Ctx(ctx).Warn().
			Err(err).
			Str("user_id", *req.UserID).
			Msg("填充用户关联信息失败")
	}

	return userProfileDTO, nil
}

// DeleteUser 删除用户（软删除）
func (l *LogicImpl) DeleteUser(
	ctx context.Context,
	req *identity_srv.DeleteUserRequest,
) error {
	if req.UserID == nil {
		return errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}

	// 1. 获取用户信息
	user, err := l.dal.UserProfile().GetByID(ctx, *req.UserID)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return errno.ErrUserNotFound
		}

		return errno.ErrOperationFailed.WithMessage("获取用户信息失败: " + err.Error())
	}

	// 2. 系统用户保护检查
	if user.IsSystemUser {
		tracelog.Ctx(ctx).Warn().
			Str("user_id", *req.UserID).
			Str("username", user.Username).
			Msg("尝试删除系统用户被拒绝")

		return errno.ErrSystemUserCannotDelete
	}

	// 3. 执行软删除
	err = l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		return txDAL.UserProfile().SoftDelete(ctx, *req.UserID)
	})
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("删除用户档案失败: " + err.Error())
	}

	// 4. 审计日志
	tracelog.Ctx(ctx).Info().
		Str("user_id", *req.UserID).
		Str("username", user.Username).
		Msg("用户删除成功")

	return nil
}

// ListUsers 分页查询用户列表
func (l *LogicImpl) ListUsers(
	ctx context.Context,
	req *identity_srv.ListUsersRequest,
) (*identity_srv.ListUsersResponse, error) {
	// 使用 Base Converter 转换分页参数
	var pageReq *rpc_base.PageRequest
	if req != nil {
		pageReq = req.Page
	}

	opts := l.converter.Base().PageRequestToQueryOptions(pageReq)

	// 构建查询条件
	conditions := &user.UserProfileQueryConditions{
		Page: opts, // 设置分页参数
	}

	// 设置业务过滤条件
	if req != nil {
		if req.Status != nil {
			status := models.UserStatus(*req.Status)
			conditions.Status = &status
		}

		if req.OrganizationID != nil && *req.OrganizationID != "" {
			conditions.OrgID = req.OrganizationID
		}
	}

	// 查询用户档案列表
	profiles, pageResult, err := l.dal.UserProfile().FindWithConditions(ctx, conditions)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询用户档案列表失败: " + err.Error())
	}

	// 使用统一的转换方法
	userProfiles := l.convertProfilesToThrift(profiles)

	// 批量填充组织和部门信息
	if err := l.enrichUserProfilesWithRelationsBatch(ctx, userProfiles); err != nil {
		// 记录警告但不影响主要结果
		tracelog.Ctx(ctx).Warn().
			Err(err).
			Msg("批量填充用户关联信息失败")
	}

	return &identity_srv.ListUsersResponse{
		Users: userProfiles,
		Page:  l.converter.Base().PageResponseToThrift(pageResult),
	}, nil
}

// SearchUsers 按条件搜索用户
func (l *LogicImpl) SearchUsers(
	ctx context.Context,
	req *identity_srv.SearchUsersRequest,
) (*identity_srv.SearchUsersResponse, error) {
	// 参数验证
	if req == nil {
		return nil, errno.ErrInvalidParams.WithMessage("请求参数不能为空")
	}

	// 使用 Base Converter 转换分页参数
	opts := l.converter.Base().PageRequestToQueryOptions(req.Page)

	// 构建查询条件
	conditions := &user.UserProfileQueryConditions{
		Page: opts, // 设置分页参数
	}

	// 设置业务过滤条件
	if req.OrganizationID != nil && *req.OrganizationID != "" {
		conditions.OrgID = req.OrganizationID
	}

	// 查询用户档案列表
	profiles, pageResult, err := l.dal.UserProfile().FindWithConditions(ctx, conditions)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("搜索用户档案失败: " + err.Error())
	}

	// 使用统一的转换方法
	userProfiles := l.convertProfilesToThrift(profiles)

	// 批量填充组织和部门信息
	if err := l.enrichUserProfilesWithRelationsBatch(ctx, userProfiles); err != nil {
		// 记录警告但不影响主要结果
		tracelog.Ctx(ctx).Warn().
			Err(err).
			Msg("批量填充用户关联信息失败")
	}

	return &identity_srv.SearchUsersResponse{
		Users: userProfiles,
		Page:  l.converter.Base().PageResponseToThrift(pageResult),
	}, nil
}

// ============================================================================
// 用户状态管理
// ============================================================================

// ChangeUserStatus 激活、停用、锁定用户
func (l *LogicImpl) ChangeUserStatus(
	ctx context.Context,
	req *identity_srv.ChangeUserStatusRequest,
) error {
	_, err := l.updateUserStatus(ctx, *req.UserID, models.UserStatus(*req.NewStatus_))
	return err
}

// UnlockUser 解锁用户
func (l *LogicImpl) UnlockUser(
	ctx context.Context,
	req *identity_srv.UnlockUserRequest,
) error {
	_, err := l.updateUserStatus(ctx, *req.UserID, models.UserStatusActive)
	return err
}

// ============================================================================
// 私有辅助方法
// ============================================================================

// convertProfilesToThrift 统一的转换方法，避免重复代码
func (l *LogicImpl) convertProfilesToThrift(
	profiles []*models.UserProfile,
) []*identity_srv.UserProfile {
	userProfiles := make([]*identity_srv.UserProfile, len(profiles))
	for i, profile := range profiles {
		userProfiles[i] = l.converter.UserProfile().ModelUserProfileToThrift(profile)
	}

	return userProfiles
}

// validateCreateUserRequest 验证创建用户请求
func (l *LogicImpl) validateCreateUserRequest(
	req *identity_srv.CreateUserRequest,
) error {
	if req.Username == nil || *req.Username == "" {
		return errno.ErrInvalidParams.WithMessage("用户名不能为空")
	}

	if req.Password == nil || *req.Password == "" {
		return errno.ErrInvalidParams.WithMessage("密码不能为空")
	}

	return nil
}

// checkUniqueConstraints 检查唯一性约束
func (l *LogicImpl) checkUniqueConstraints(
	ctx context.Context,
	req *identity_srv.UpdateUserRequest,
	excludeID string,
) error {
	if req.Email != nil {
		exists, err := l.dal.UserProfile().CheckEmailExists(ctx, *req.Email, excludeID)
		if err != nil {
			return errno.ErrOperationFailed.WithMessage("检查邮箱唯一性失败: " + err.Error())
		}

		if exists {
			return errno.ErrEmailAlreadyExists
		}
	}

	if req.Phone != nil {
		exists, err := l.dal.UserProfile().CheckPhoneExists(ctx, *req.Phone, excludeID)
		if err != nil {
			return errno.ErrOperationFailed.WithMessage("检查手机号唯一性失败: " + err.Error())
		}

		if exists {
			return errno.ErrPhoneAlreadyExists
		}
	}

	return nil
}

// updateUserStatus 更新用户状态的通用方法
func (l *LogicImpl) updateUserStatus(
	ctx context.Context,
	userID string,
	status models.UserStatus,
) (*identity_srv.UserProfile, error) {
	if userID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}

	// 获取用户档案
	profile, err := l.dal.UserProfile().GetByID(ctx, userID)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return nil, errno.ErrUserNotFound
		}

		return nil, errno.ErrOperationFailed.WithMessage("获取用户档案失败: " + err.Error())
	}

	// 更新状态
	profile.Status = status

	// 保存更新
	err = l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		if err := txDAL.UserProfile().Update(ctx, profile); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return l.converter.UserProfile().ModelUserProfileToThrift(profile), nil
}

// ============================================================================
// 辅助方法 - 关联数据填充
// ============================================================================

// enrichUserProfileWithRelations 填充用户档案的关联字段
// 包括：主组织ID、主部门ID、角色ID列表
func (l *LogicImpl) enrichUserProfileWithRelations(
	ctx context.Context,
	profile *identity_srv.UserProfile,
) error {
	if profile == nil || profile.ID == nil {
		return nil
	}

	// 查询用户的主成员关系
	primaryMembership, err := l.dal.UserMembership().GetPrimaryMembership(ctx, *profile.ID)
	if err != nil {
		// 如果没有主成员关系，不返回错误（用户可能还未分配组织）
		if errno.IsRecordNotFound(err) {
			return nil
		}

		return err
	}

	// 填充主组织ID
	if primaryMembership.OrganizationID.String() != "00000000-0000-0000-0000-000000000000" {
		orgID := primaryMembership.OrganizationID.String()
		profile.PrimaryOrganizationID = &orgID
	}

	// 填充主部门ID（可选）
	if primaryMembership.DepartmentID.String() != "00000000-0000-0000-0000-000000000000" {
		deptID := primaryMembership.DepartmentID.String()
		profile.PrimaryDepartmentID = &deptID
	}

	return nil
}

// enrichUserProfilesWithRelationsBatch 批量填充用户档案的关联字段
// 使用批量查询避免 N+1 查询问题
func (l *LogicImpl) enrichUserProfilesWithRelationsBatch(
	ctx context.Context,
	profiles []*identity_srv.UserProfile,
) error {
	if len(profiles) == 0 {
		return nil
	}

	// 1. 提取所有用户ID
	userIDs := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		if profile != nil && profile.ID != nil {
			userIDs = append(userIDs, *profile.ID)
		}
	}

	if len(userIDs) == 0 {
		return nil
	}

	// 2. 批量查询主成员关系
	membershipsMap, err := l.dal.UserMembership().GetPrimaryMembershipsByUserIDs(ctx, userIDs)
	if err != nil {
		return err
	}

	// 3. 填充每个用户的组织/部门信息
	for _, profile := range profiles {
		if profile == nil || profile.ID == nil {
			continue
		}

		membership, exists := membershipsMap[*profile.ID]
		if !exists {
			// 用户没有主成员关系，跳过
			continue
		}

		// 填充主组织ID
		if membership.OrganizationID.String() != "00000000-0000-0000-0000-000000000000" {
			orgID := membership.OrganizationID.String()
			profile.PrimaryOrganizationID = &orgID
		}

		// 填充主部门ID（可选）
		if membership.DepartmentID.String() != "00000000-0000-0000-0000-000000000000" {
			deptID := membership.DepartmentID.String()
			profile.PrimaryDepartmentID = &deptID
		}
	}

	return nil
}

// getModifiedFields 获取修改的字段列表（用于审计日志）
func getModifiedFields(req *identity_srv.UpdateUserRequest) []string {
	fields := []string{}

	// 注意：UpdateUserRequest 不包含 username 和 status 字段
	// 这些关键属性通过专门的接口修改
	if req.Email != nil {
		fields = append(fields, "email")
	}

	if req.Phone != nil {
		fields = append(fields, "phone")
	}

	if req.RealName != nil {
		fields = append(fields, "real_name")
	}

	if req.FirstName != nil {
		fields = append(fields, "first_name")
	}

	if req.LastName != nil {
		fields = append(fields, "last_name")
	}

	if req.ProfessionalTitle != nil {
		fields = append(fields, "professional_title")
	}

	return fields
}
