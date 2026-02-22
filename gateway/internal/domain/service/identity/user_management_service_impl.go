package identity

import (
	"context"
	"fmt"

	"github.com/cloudwego/kitex/pkg/kerrors"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	identityConv "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/common"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// userManagementServiceImpl 用户管理服务实现
type userManagementServiceImpl struct {
	*common.BaseService
	identityClient identitycli.IdentityClient
	assembler      identityConv.Assembler
}

// NewUserManagementService 创建新的用户管理服务实例
func NewUserManagementService(
	identityClient identitycli.IdentityClient,
	assembler identityConv.Assembler,
	logger *hertzZerolog.Logger,
) UserService {
	return &userManagementServiceImpl{
		BaseService:    common.NewBaseService(logger),
		identityClient: identityClient,
		assembler:      assembler,
	}
}

// =================================================================
// 2. 用户管理模块 (User Management)
// =================================================================

func (s *userManagementServiceImpl) CreateUser(
	ctx context.Context,
	req *identity.CreateUserRequestDTO,
	operatorID string,
) (*identity.UserProfileResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "创建用户",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.User().ToRPCCreateUserRequest(req)
			return s.identityClient.CreateUser(ctx, rpcReq)
		},
		"username", req.Username,
	)
	if err != nil {
		return nil, err
	}

	rpcUserProfile := result.(*identity_srv.UserProfile)
	httpUserProfile := s.assembler.User().ToHTTPUserProfile(rpcUserProfile)

	// 用户创建成功后，处理组织关系和角色分配
	if rpcUserProfile.ID != nil {
		// 1. 如果指定了组织ID，创建主成员关系
		if req.OrganizationID != nil && *req.OrganizationID != "" {
			httpUserProfile.PrimaryOrganizationID = req.OrganizationID
			s.assignOrganizationMembership(ctx, rpcUserProfile.ID, req.OrganizationID, operatorID)
		}

		// 2. 如果指定了角色ID列表，批量分配角色
		if len(req.RoleIDs) > 0 {
			httpUserProfile.RoleIDs = req.RoleIDs
			s.assignRolesToUser(ctx, rpcUserProfile.ID, req.RoleIDs, operatorID)
		}
	}

	httpResp := &identity.UserProfileResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		User:     httpUserProfile,
	}

	return httpResp, nil
}

func (s *userManagementServiceImpl) GetUser(
	ctx context.Context,
	req *identity.GetUserRequestDTO,
) (*identity.UserProfileResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "获取用户信息",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.User().ToRPCGetUserRequest(req)
			return s.identityClient.GetUser(ctx, rpcReq)
		},
		"user_id", req.UserID,
	)
	if err != nil {
		return nil, err
	}

	rpcUserProfile := result.(*identity_srv.UserProfile)
	httpUserProfile := s.assembler.User().ToHTTPUserProfile(rpcUserProfile)

	// 填充用户角色ID列表
	s.enrichUserProfileWithRoles(ctx, httpUserProfile)

	httpResp := &identity.UserProfileResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		User:     httpUserProfile,
	}

	return httpResp, nil
}

func (s *userManagementServiceImpl) GetMe(
	ctx context.Context,
	userID string,
) (*identity.UserProfileResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "获取当前用户信息",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.User().
				ToRPCGetUserRequest(&identity.GetUserRequestDTO{UserID: &userID})

			return s.identityClient.GetUser(ctx, rpcReq)
		},
		"user_id", userID,
	)
	if err != nil {
		return nil, err
	}

	rpcUserProfile := result.(*identity_srv.UserProfile)
	httpUserProfile := s.assembler.User().ToHTTPUserProfile(rpcUserProfile)

	// 填充用户角色ID列表
	s.enrichUserProfileWithRoles(ctx, httpUserProfile)

	httpResp := &identity.UserProfileResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		User:     httpUserProfile,
	}

	return httpResp, nil
}

func (s *userManagementServiceImpl) UpdateUser(
	ctx context.Context,
	req *identity.UpdateUserRequestDTO,
	operatorID string,
) (*identity.UserProfileResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "更新用户信息",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.User().ToRPCUpdateUserRequest(req)
			return s.identityClient.UpdateUser(ctx, rpcReq)
		},
		"user_id", req.UserID,
	)
	if err != nil {
		return nil, err
	}

	rpcUserProfile := result.(*identity_srv.UserProfile)
	httpUserProfile := s.assembler.User().ToHTTPUserProfile(rpcUserProfile)

	// 用户更新成功后，处理组织关系和角色更新
	if req.UserID != nil {
		// 1. 如果指定了组织ID，更新主成员关系
		if req.OrganizationID != nil && *req.OrganizationID != "" {
			httpUserProfile.PrimaryOrganizationID = req.OrganizationID
			s.updateOrganizationMembership(ctx, req.UserID, req.OrganizationID, operatorID)
		}

		// 2. 如果指定了角色ID列表，更新角色分配
		// 使用 req.RoleIDs != nil 来判断字段是否被提供，支持清空所有角色
		if req.RoleIDs != nil {
			httpUserProfile.RoleIDs = req.RoleIDs
			if err := s.updateUserRoles(ctx, req.UserID, req.RoleIDs, operatorID); err != nil {
				// 角色更新失败，记录错误并返回
				s.Log(ctx).Error().Err(err).Str("user_id", *req.UserID).Msg("更新用户角色失败")

				return nil, err
			}
		}
	}

	httpResp := &identity.UserProfileResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		User:     httpUserProfile,
	}

	return httpResp, nil
}

func (s *userManagementServiceImpl) UpdateMe(
	ctx context.Context,
	req *identity.UpdateMeRequestDTO,
	userID string,
) (*identity.UserProfileResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "更新当前用户信息",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.User().ToRPCUpdateMeRequest(req)
			// 从认证上下文设置 UserID
			rpcReq.UserID = &userID

			return s.identityClient.UpdateUser(ctx, rpcReq)
		},
		"user_id", userID,
	)
	if err != nil {
		return nil, err
	}

	rpcUserProfile := result.(*identity_srv.UserProfile)
	httpUserProfile := s.assembler.User().ToHTTPUserProfile(rpcUserProfile)

	httpResp := &identity.UserProfileResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		User:     httpUserProfile,
	}

	return httpResp, nil
}

func (s *userManagementServiceImpl) DeleteUser(
	ctx context.Context,
	req *identity.DeleteUserRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	err := s.ProcessRPCVoidCall(ctx, "删除用户",
		func(ctx context.Context) error {
			rpcReq := s.assembler.User().ToRPCDeleteUserRequest(req)
			return s.identityClient.DeleteUser(ctx, rpcReq)
		},
		"user_id", req.UserID,
	)
	if err != nil {
		return nil, err
	}

	return s.ResponseBuilder().BuildOperationStatusResponse(), nil
}

func (s *userManagementServiceImpl) ListUsers(
	ctx context.Context,
	req *identity.ListUsersRequestDTO,
) (*identity.ListUsersResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "获取用户列表",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.User().ToRPCListUsersRequest(req)
			return s.identityClient.ListUsers(ctx, rpcReq)
		},
		"page", req.Page, "organization_id", req.OrganizationID, "status", req.Status,
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.ListUsersResponse)
	httpResp := s.assembler.User().ToHTTPListUsersResponse(rpcResp)
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	// 批量填充用户角色信息
	if httpResp.Users != nil {
		s.enrichUserProfilesWithRolesBatch(ctx, httpResp.Users)
	}

	return httpResp, nil
}

func (s *userManagementServiceImpl) SearchUsers(
	ctx context.Context,
	req *identity.SearchUsersRequestDTO,
) (*identity.SearchUsersResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "搜索用户",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.User().ToRPCSearchUsersRequest(req)
			return s.identityClient.SearchUsers(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.SearchUsersResponse)
	httpResp := s.assembler.User().ToHTTPSearchUsersResponse(rpcResp)
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	// 批量填充用户角色信息
	if httpResp.Users != nil {
		s.enrichUserProfilesWithRolesBatch(ctx, httpResp.Users)
	}

	return httpResp, nil
}

func (s *userManagementServiceImpl) ChangeUserStatus(
	ctx context.Context,
	req *identity.ChangeUserStatusRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	err := s.ProcessRPCVoidCall(ctx, "变更用户状态",
		func(ctx context.Context) error {
			rpcReq := s.assembler.User().ToRPCChangeUserStatusRequest(req)
			return s.identityClient.ChangeUserStatus(ctx, rpcReq)
		},
		"user_id", req.UserID, "status", req.NewStatus,
	)
	if err != nil {
		return nil, err
	}

	return s.ResponseBuilder().BuildOperationStatusResponse(), nil
}

func (s *userManagementServiceImpl) UnlockUser(
	ctx context.Context,
	req *identity.UnlockUserRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	err := s.ProcessRPCVoidCall(ctx, "解锁用户",
		func(ctx context.Context) error {
			rpcReq := s.assembler.User().ToRPCUnlockUserRequest(req)
			return s.identityClient.UnlockUser(ctx, rpcReq)
		},
		"user_id", req.UserID,
	)
	if err != nil {
		return nil, err
	}

	return s.ResponseBuilder().BuildOperationStatusResponse(), nil
}

// =================================================================
// 私有辅助方法 (Private Helper Methods)
// =================================================================

// enrichUserProfileWithRoles 为用户档案填充角色ID列表
// 调用 identity_srv 获取用户的角色分配信息
func (s *userManagementServiceImpl) enrichUserProfileWithRoles(
	ctx context.Context,
	userProfile *identity.UserProfileDTO,
) {
	// 如果用户ID为空，直接返回
	if userProfile == nil || userProfile.ID == nil {
		return
	}

	// 调用 identity_srv 获取用户角色
	roleResp, err := s.identityClient.ListUserRoleAssignments(
		ctx,
		&identity_srv.UserRoleQueryRequest{
			UserID: userProfile.ID,
		},
	)
	if err != nil {
		// 角色获取失败不影响主流程，仅记录警告
		s.Log(ctx).Warn().Err(err).Str("user_id", *userProfile.ID).Msg("获取用户角色失败")

		return
	}

	// 提取角色ID列表
	if roleResp != nil && roleResp.Assignments != nil {
		roleIDs := make([]string, 0, len(roleResp.Assignments))
		for _, assignment := range roleResp.Assignments {
			if assignment.RoleID != nil {
				roleIDs = append(roleIDs, *assignment.RoleID)
			}
		}

		userProfile.RoleIDs = roleIDs
	}
}

// enrichUserProfilesWithRolesBatch 批量为用户档案填充角色ID列表
// 使用批量查询避免 N+1 查询问题
func (s *userManagementServiceImpl) enrichUserProfilesWithRolesBatch(
	ctx context.Context,
	userProfiles []*identity.UserProfileDTO,
) {
	// 如果用户列表为空，直接返回
	if len(userProfiles) == 0 {
		return
	}

	// 1. 提取所有用户ID
	userIDs := make([]string, 0, len(userProfiles))
	for _, profile := range userProfiles {
		if profile != nil && profile.ID != nil {
			userIDs = append(userIDs, *profile.ID)
		}
	}

	if len(userIDs) == 0 {
		return
	}

	// 2. 批量调用 BatchGetUserRoles RPC
	batchResp, err := s.identityClient.BatchGetUserRoles(ctx,
		&identity_srv.BatchGetUserRolesRequest{
			UserIDs: userIDs,
		})
	if err != nil {
		// 角色获取失败不影响主流程，仅记录警告
		s.Log(ctx).Warn().Err(err).Msg("批量获取用户角色失败")
		return
	}

	// 3. 构建用户ID到角色ID列表的映射
	rolesMap := make(map[string][]string)

	if batchResp != nil && batchResp.UserRoles != nil {
		for _, userRole := range batchResp.UserRoles {
			if userRole.UserID != nil {
				rolesMap[*userRole.UserID] = userRole.RoleIDs
			}
		}
	}

	// 4. 为每个用户填充角色信息
	for _, profile := range userProfiles {
		if profile != nil && profile.ID != nil {
			if roleIDs, exists := rolesMap[*profile.ID]; exists {
				profile.RoleIDs = roleIDs
			}
		}
	}
}

// assignOrganizationMembership 为用户分配组织成员关系
// 创建用户时，如果指定了 organizationID，则创建主成员关系
func (s *userManagementServiceImpl) assignOrganizationMembership(
	ctx context.Context,
	userID *string,
	organizationID *string,
	operatorID string,
) {
	if userID == nil || organizationID == nil {
		return
	}

	isPrimary := true

	_, err := s.identityClient.AddMembership(ctx, &identity_srv.AddMembershipRequest{
		UserID:         userID,
		OrganizationID: organizationID,
		IsPrimary:      isPrimary,
	})
	if err != nil {
		// 成员关系创建失败不影响主流程，仅记录警告
		s.Log(ctx).Warn().Err(err).
			Str("user_id", *userID).
			Str("organization_id", *organizationID).
			Msg("创建主成员关系失败")
	} else {
		s.Log(ctx).Info().
			Str("user_id", *userID).
			Str("organization_id", *organizationID).
			Msg("成功创建主成员关系")
	}
}

// assignRolesToUser 批量为用户分配角色
// 创建用户时，如果指定了 roleIDs，则批量分配角色
func (s *userManagementServiceImpl) assignRolesToUser(
	ctx context.Context,
	userID *string,
	roleIDs []string,
	operatorID string,
) {
	if userID == nil || len(roleIDs) == 0 {
		return
	}

	successCount := 0

	for _, roleID := range roleIDs {
		if roleID == "" {
			continue
		}

		_, err := s.identityClient.AssignRoleToUser(ctx, &identity_srv.AssignRoleToUserRequest{
			UserID:     userID,
			RoleID:     &roleID,
			AssignedBy: &operatorID,
		})
		if err != nil {
			s.Log(ctx).Warn().Err(err).
				Str("user_id", *userID).
				Str("role_id", roleID).
				Msg("分配角色失败")
		} else {
			successCount++
		}
	}

	s.Log(ctx).Info().
		Str("user_id", *userID).
		Int("total", len(roleIDs)).
		Int("success", successCount).
		Msg("批量分配角色完成")
}

// updateOrganizationMembership 更新用户的组织成员关系
// 更新用户时，如果指定了 organizationID，则更新或创建主成员关系
func (s *userManagementServiceImpl) updateOrganizationMembership(
	ctx context.Context,
	userID *string,
	organizationID *string,
	operatorID string,
) {
	if userID == nil || organizationID == nil {
		return
	}

	// 获取用户的主成员关系
	primaryMembership, err := s.identityClient.GetPrimaryMembership(ctx, *userID)
	if err != nil || primaryMembership == nil {
		// 如果没有主成员关系，创建一个新的
		s.assignOrganizationMembership(ctx, userID, organizationID, operatorID)
		return
	}

	// 如果主成员关系的组织ID不同，更新它
	if primaryMembership.OrganizationID != nil &&
		*primaryMembership.OrganizationID != *organizationID {
		isPrimary := true

		_, err := s.identityClient.UpdateMembership(ctx, &identity_srv.UpdateMembershipRequest{
			MembershipID:   primaryMembership.ID,
			OrganizationID: organizationID,
			IsPrimary:      &isPrimary,
		})
		if err != nil {
			s.Log(ctx).Warn().Err(err).
				Str("user_id", *userID).
				Str("membership_id", *primaryMembership.ID).
				Str("new_organization_id", *organizationID).
				Msg("更新主成员关系失败")
		} else {
			s.Log(ctx).Info().
				Str("user_id", *userID).
				Str("membership_id", *primaryMembership.ID).
				Str("new_organization_id", *organizationID).
				Msg("成功更新主成员关系")
		}
	}
}

// updateUserRoles 更新用户的角色分配
// 更新用户时，如果指定了 roleIDs，则替换用户的所有角色
func (s *userManagementServiceImpl) updateUserRoles(
	ctx context.Context,
	userID *string,
	newRoleIDs []string,
	operatorID string,
) error {
	if userID == nil {
		return nil
	}

	// 1. 获取用户当前的所有角色
	roleResp, err := s.identityClient.ListUserRoleAssignments(
		ctx,
		&identity_srv.UserRoleQueryRequest{
			UserID: userID,
		},
	)
	if err != nil {
		s.Log(ctx).Warn().Err(err).Str("user_id", *userID).Msg("获取用户当前角色失败")
		// 即使获取失败，仍然尝试分配新角色
		s.assignRolesToUser(ctx, userID, newRoleIDs, operatorID)

		return nil
	}

	// 2. 构建当前角色ID集合
	currentRoleIDs := make(map[string]bool)

	if roleResp != nil && roleResp.Assignments != nil {
		for _, assignment := range roleResp.Assignments {
			if assignment.RoleID != nil {
				currentRoleIDs[*assignment.RoleID] = true
			}
		}
	}

	// 3. 构建新角色ID集合
	newRoleIDSet := make(map[string]bool)

	for _, roleID := range newRoleIDs {
		if roleID != "" {
			newRoleIDSet[roleID] = true
		}
	}

	// 4. 计算需要添加的角色（新有，旧没有）
	var rolesToAdd []string

	for roleID := range newRoleIDSet {
		if !currentRoleIDs[roleID] {
			rolesToAdd = append(rolesToAdd, roleID)
		}
	}

	// 5. 计算需要删除的角色（旧有，新没有）
	var rolesToRemove []string

	for roleID := range currentRoleIDs {
		if !newRoleIDSet[roleID] {
			rolesToRemove = append(rolesToRemove, roleID)
		}
	}

	// 6. 删除旧角色（增强版 - 带错误检测和回滚）
	if len(rolesToRemove) > 0 {
		if err := s.removeRolesWithRollback(ctx, userID, rolesToRemove, operatorID); err != nil {
			return err
		}
	}

	// 7. 分配新角色
	if len(rolesToAdd) > 0 {
		s.assignRolesToUser(ctx, userID, rolesToAdd, operatorID)
	}

	s.Log(ctx).Info().
		Str("user_id", *userID).
		Int("current_count", len(currentRoleIDs)).
		Int("new_count", len(newRoleIDSet)).
		Int("added", len(rolesToAdd)).
		Int("removed", len(rolesToRemove)).
		Msg("更新用户角色完成")

	return nil
}

// buildSystemRoleProtectionError 构造系统角色保护的友好错误消息
func (s *userManagementServiceImpl) buildSystemRoleProtectionError(
	ctx context.Context,
	roleID string,
	originalMessage string,
) error {
	// 构造详细的错误消息
	// RPC 返回的 originalMessage 已经包含了角色名称信息
	friendlyMessage := fmt.Sprintf(
		"无法更新用户角色：%s\n\n"+
			"角色 ID: %s\n\n"+
			"提示: 系统角色是系统内置角色，用于保证系统正常运行，不能被撤销。",
		originalMessage,
		roleID,
	)

	// 返回业务错误，保持 207017 错误码以便客户端识别
	return errors.NewAPIError(207017, friendlyMessage)
}

// removeRolesWithRollback 批量撤销角色，失败时自动回滚已撤销的角色
func (s *userManagementServiceImpl) removeRolesWithRollback(
	ctx context.Context,
	userID *string,
	rolesToRemove []string,
	operatorID string,
) error {
	revokedRoles := make([]string, 0, len(rolesToRemove)) // 记录已成功撤销的角色，用于可能的回滚

	for _, roleID := range rolesToRemove {
		err := s.identityClient.RevokeRoleFromUser(
			ctx,
			&identity_srv.RevokeRoleFromUserRequest{
				UserID:    userID,
				RoleID:    &roleID,
				RevokedBy: &operatorID,
			},
		)
		if err != nil {
			return s.handleRevokeError(ctx, err, userID, roleID, revokedRoles, operatorID)
		}

		// 记录成功撤销的角色，用于可能的回滚
		revokedRoles = append(revokedRoles, roleID)
	}

	s.Log(ctx).Info().
		Str("user_id", *userID).
		Int("count", len(revokedRoles)).
		Msg("批量撤销角色成功")

	return nil
}

// handleRevokeError 处理角色撤销错误：检查系统角色保护并回滚
func (s *userManagementServiceImpl) handleRevokeError(
	ctx context.Context,
	err error,
	userID *string,
	roleID string,
	revokedRoles []string,
	operatorID string,
) error {
	// 检查是否为系统角色保护错误 (207017)
	if bizErr, isBizErr := kerrors.FromBizStatusError(err); isBizErr {
		if bizErr.BizStatusCode() == 207017 {
			s.rollbackRevokedRoles(ctx, userID, revokedRoles, operatorID)

			return s.buildSystemRoleProtectionError(ctx, roleID, bizErr.BizMessage())
		}
	}

	// 其他错误也应该中断并回滚（保证原子性）
	s.Log(ctx).Error().Err(err).
		Str("user_id", *userID).
		Str("role_id", roleID).
		Msg("撤销角色失败，回滚操作")
	s.rollbackRevokedRoles(ctx, userID, revokedRoles, operatorID)

	return errors.ProcessRPCError(err, "撤销角色失败")
}

// rollbackRevokedRoles 回滚已撤销的角色，恢复用户角色状态
func (s *userManagementServiceImpl) rollbackRevokedRoles(
	ctx context.Context,
	userID *string,
	revokedRoles []string,
	operatorID string,
) {
	if len(revokedRoles) == 0 {
		return
	}

	s.Log(ctx).Warn().
		Str("user_id", *userID).
		Int("count", len(revokedRoles)).
		Msg("开始回滚已撤销的角色")

	successCount := 0

	for _, roleID := range revokedRoles {
		_, err := s.identityClient.AssignRoleToUser(ctx, &identity_srv.AssignRoleToUserRequest{
			UserID:     userID,
			RoleID:     &roleID,
			AssignedBy: &operatorID,
		})
		if err != nil {
			s.Log(ctx).Error().Err(err).
				Str("user_id", *userID).
				Str("role_id", roleID).
				Msg("回滚角色失败")
		} else {
			successCount++
		}
	}

	s.Log(ctx).Info().
		Str("user_id", *userID).
		Int("total", len(revokedRoles)).
		Int("success", successCount).
		Msg("角色回滚完成")
}
