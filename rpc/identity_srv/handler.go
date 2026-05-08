package main

import (
	"context"
	"errors"

	"github.com/masonsxu/cloudwego-microservice-demo/iamclient"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic"
	identity_srv "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/wire"
)

// IdentityServiceImpl implements the last service interface defined in the IDL.
type IdentityServiceImpl struct {
	logic logic.Logic
	iam   *iamclient.Client
}

// NewIdentityServiceImpl 从应用容器创建 IdentityServiceImpl 实例
func NewIdentityServiceImpl(container *wire.AppContainer) *IdentityServiceImpl {
	return &IdentityServiceImpl{
		logic: container.Logic,
		iam:   container.IAMClient,
	}
}

// requirePerm 从 ctx 还原 Subject 并请求 PDP 决策。
//
// 行为：
//   - Subject 还原失败 → ErrUnauthenticated（Kitex BizStatusError）
//   - PDP 拒绝 → ErrPermissionDenied
//   - 网络/RPC 错误 → ErrOperationFailed（fail-closed，绝不放行）
func (s *IdentityServiceImpl) requirePerm(
	ctx context.Context,
	action, resource string,
	opts ...iamclient.CheckOpt,
) error {
	subject, err := s.iam.SubjectFromContext(ctx)
	if err != nil {
		return errno.ToKitexError(errno.ErrUnauthenticated.WithMessage(err.Error()))
	}

	if err := subject.MustCheck(ctx, action, resource, opts...); err != nil {
		if errors.Is(err, iamclient.ErrPermissionDenied) {
			return errno.ToKitexError(errno.ErrPermissionDenied.WithMessage(err.Error()))
		}

		return errno.ToKitexError(errno.ErrOperationFailed.WithMessage("authz failed: " + err.Error()))
	}

	return nil
}

// requirePermSelfOrCheck 改自己资源（caller_id == targetUserID）直接放行，否则走 PDP。
//
// 用于 ChangePassword 等"自己改自己/管理员改他人"双语义场景（提案 §11 Q1 决策 A）。
func (s *IdentityServiceImpl) requirePermSelfOrCheck(
	ctx context.Context,
	targetUserID, action, resource string,
) error {
	subject, err := s.iam.SubjectFromContext(ctx)
	if err != nil {
		return errno.ToKitexError(errno.ErrUnauthenticated.WithMessage(err.Error()))
	}

	if subject.UserID == targetUserID {
		return nil
	}

	if err := subject.MustCheck(ctx, action, resource); err != nil {
		if errors.Is(err, iamclient.ErrPermissionDenied) {
			return errno.ToKitexError(errno.ErrPermissionDenied.WithMessage(err.Error()))
		}

		return errno.ToKitexError(errno.ErrOperationFailed.WithMessage("authz failed: " + err.Error()))
	}

	return nil
}

// derefStr 安全解引用 *string，nil 时返回空串。
func derefStr(p *string) string {
	if p == nil {
		return ""
	}

	return *p
}

// ===========================================================================
// UserProfile
// ===========================================================================

// CreateUser implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CreateUser(
	ctx context.Context,
	req *identity_srv.CreateUserRequest,
) (resp *identity_srv.CreateUserResponse, err error) {
	if err := s.requirePerm(ctx, "create", "user"); err != nil {
		return nil, err
	}

	userProfile, err := s.logic.CreateUser(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.CreateUserResponse{UserProfile: userProfile}, nil
}

// GetUser implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetUser(
	ctx context.Context,
	req *identity_srv.GetUserRequest,
) (resp *identity_srv.GetUserResponse, err error) {
	userProfile, err := s.logic.GetUser(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.GetUserResponse{UserProfile: userProfile}, nil
}

// UpdateUser implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UpdateUser(
	ctx context.Context,
	req *identity_srv.UpdateUserRequest,
) (resp *identity_srv.UpdateUserResponse, err error) {
	if err := s.requirePerm(ctx, "update", "user:"+derefStr(req.UserID)); err != nil {
		return nil, err
	}

	userProfile, err := s.logic.UpdateUser(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.UpdateUserResponse{UserProfile: userProfile}, nil
}

// DeleteUser implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteUser(
	ctx context.Context,
	req *identity_srv.DeleteUserRequest,
) (resp *identity_srv.DeleteUserResponse, err error) {
	if err := s.requirePerm(ctx, "delete", "user:"+derefStr(req.UserID)); err != nil {
		return nil, err
	}

	if err := s.logic.DeleteUser(ctx, req); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.DeleteUserResponse{}, nil
}

// ListUsers implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ListUsers(
	ctx context.Context,
	req *identity_srv.ListUsersRequest,
) (resp *identity_srv.ListUsersResponse, err error) {
	resp, err = s.logic.ListUsers(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// SearchUsers implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) SearchUsers(
	ctx context.Context,
	req *identity_srv.SearchUsersRequest,
) (resp *identity_srv.SearchUsersResponse, err error) {
	resp, err = s.logic.SearchUsers(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// ChangeUserStatus implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ChangeUserStatus(
	ctx context.Context,
	req *identity_srv.ChangeUserStatusRequest,
) (resp *identity_srv.ChangeUserStatusResponse, err error) {
	if err := s.requirePerm(ctx, "change_status", "user:"+derefStr(req.UserID)); err != nil {
		return nil, err
	}

	if err := s.logic.ChangeUserStatus(ctx, req); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.ChangeUserStatusResponse{}, nil
}

// UnlockUser implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UnlockUser(
	ctx context.Context,
	req *identity_srv.UnlockUserRequest,
) (resp *identity_srv.UnlockUserResponse, err error) {
	if err := s.requirePerm(ctx, "unlock", "user:"+derefStr(req.UserID)); err != nil {
		return nil, err
	}

	if err := s.logic.UnlockUser(ctx, req); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.UnlockUserResponse{}, nil
}

// ===========================================================================
// Authentication
// ===========================================================================
// Login implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) Login(
	ctx context.Context,
	req *identity_srv.LoginRequest,
) (resp *identity_srv.LoginResponse, err error) {
	resp, err = s.logic.Login(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// ChangePassword implements the IdentityServiceImpl interface.
//
// 自助改密码：caller 改自己直接放行；管理员代改他人需 PDP 决策
// `update_password` on `user:{id}`（提案 §11 Q1 决策 A）。
func (s *IdentityServiceImpl) ChangePassword(
	ctx context.Context,
	req *identity_srv.ChangePasswordRequest,
) (resp *identity_srv.ChangePasswordResponse, err error) {
	target := derefStr(req.UserID)
	if err := s.requirePermSelfOrCheck(ctx, target, "update_password", "user:"+target); err != nil {
		return nil, err
	}

	if err := s.logic.ChangePassword(ctx, req); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.ChangePasswordResponse{}, nil
}

// ResetPassword implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ResetPassword(
	ctx context.Context,
	req *identity_srv.ResetPasswordRequest,
) (resp *identity_srv.ResetPasswordResponse, err error) {
	if err := s.requirePerm(ctx, "reset_password", "user:"+derefStr(req.UserID)); err != nil {
		return nil, err
	}

	if err := s.logic.ResetPassword(ctx, req); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.ResetPasswordResponse{}, nil
}

// ForcePasswordChange implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ForcePasswordChange(
	ctx context.Context,
	req *identity_srv.ForcePasswordChangeRequest,
) (resp *identity_srv.ForcePasswordChangeResponse, err error) {
	if err := s.requirePerm(ctx, "force_password_change", "user:"+derefStr(req.UserID)); err != nil {
		return nil, err
	}

	if err := s.logic.ForcePasswordChange(ctx, req); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.ForcePasswordChangeResponse{}, nil
}

// ===========================================================================
// OrgManagement
// ===========================================================================

// AddMembership implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) AddMembership(
	ctx context.Context,
	req *identity_srv.AddMembershipRequest,
) (resp *identity_srv.AddMembershipResponse, err error) {
	if err := s.requirePerm(ctx, "create", "membership"); err != nil {
		return nil, err
	}

	membership, err := s.logic.AddMembership(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.AddMembershipResponse{Membership: membership}, nil
}

// UpdateMembership implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UpdateMembership(
	ctx context.Context,
	req *identity_srv.UpdateMembershipRequest,
) (resp *identity_srv.UpdateMembershipResponse, err error) {
	if err := s.requirePerm(ctx, "update", "membership:"+derefStr(req.MembershipID)); err != nil {
		return nil, err
	}

	membership, err := s.logic.UpdateMembership(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.UpdateMembershipResponse{Membership: membership}, nil
}

// RemoveMembership implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) RemoveMembership(
	ctx context.Context,
	req *identity_srv.RemoveMembershipRequest,
) (resp *identity_srv.RemoveMembershipResponse, err error) {
	if err := s.requirePerm(ctx, "delete", "membership:"+derefStr(req.MembershipID)); err != nil {
		return nil, err
	}

	if err := s.logic.RemoveMembership(ctx, req.GetMembershipID()); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.RemoveMembershipResponse{}, nil
}

// GetUserMemberships implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetUserMemberships(
	ctx context.Context,
	req *identity_srv.GetUserMembershipsRequest,
) (resp *identity_srv.GetUserMembershipsResponse, err error) {
	resp, err = s.logic.GetUserMemberships(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// CreateOrganization implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CreateOrganization(
	ctx context.Context,
	req *identity_srv.CreateOrganizationRequest,
) (resp *identity_srv.CreateOrganizationResponse, err error) {
	if err := s.requirePerm(ctx, "create", "organization"); err != nil {
		return nil, err
	}

	organization, err := s.logic.CreateOrganization(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.CreateOrganizationResponse{Organization: organization}, nil
}

// GetOrganization implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetOrganization(
	ctx context.Context,
	req *identity_srv.GetOrganizationRequest,
) (resp *identity_srv.GetOrganizationResponse, err error) {
	organization, err := s.logic.GetOrganization(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.GetOrganizationResponse{Organization: organization}, nil
}

// UpdateOrganization implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UpdateOrganization(
	ctx context.Context,
	req *identity_srv.UpdateOrganizationRequest,
) (resp *identity_srv.UpdateOrganizationResponse, err error) {
	if err := s.requirePerm(ctx, "update", "organization:"+derefStr(req.OrganizationID)); err != nil {
		return nil, err
	}

	organization, err := s.logic.UpdateOrganization(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.UpdateOrganizationResponse{Organization: organization}, nil
}

// DeleteOrganization implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteOrganization(
	ctx context.Context,
	req *identity_srv.DeleteOrganizationRequest,
) (resp *identity_srv.DeleteOrganizationResponse, err error) {
	if err := s.requirePerm(ctx, "delete", "organization:"+req.GetOrganizationID()); err != nil {
		return nil, err
	}

	if err := s.logic.DeleteOrganization(ctx, req.GetOrganizationID()); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.DeleteOrganizationResponse{}, nil
}

// ListOrganizations implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ListOrganizations(
	ctx context.Context,
	req *identity_srv.ListOrganizationsRequest,
) (resp *identity_srv.ListOrganizationsResponse, err error) {
	resp, err = s.logic.ListOrganizations(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// CreateDepartment implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CreateDepartment(
	ctx context.Context,
	req *identity_srv.CreateDepartmentRequest,
) (resp *identity_srv.CreateDepartmentResponse, err error) {
	if err := s.requirePerm(ctx, "create", "department"); err != nil {
		return nil, err
	}

	department, err := s.logic.CreateDepartment(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.CreateDepartmentResponse{Department: department}, nil
}

// GetDepartment implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetDepartment(
	ctx context.Context,
	req *identity_srv.GetDepartmentRequest,
) (resp *identity_srv.GetDepartmentResponse, err error) {
	department, err := s.logic.GetDepartment(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.GetDepartmentResponse{Department: department}, nil
}

// UpdateDepartment implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UpdateDepartment(
	ctx context.Context,
	req *identity_srv.UpdateDepartmentRequest,
) (resp *identity_srv.UpdateDepartmentResponse, err error) {
	if err := s.requirePerm(ctx, "update", "department:"+derefStr(req.DepartmentID)); err != nil {
		return nil, err
	}

	department, err := s.logic.UpdateDepartment(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.UpdateDepartmentResponse{Department: department}, nil
}

// DeleteDepartment implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteDepartment(
	ctx context.Context,
	req *identity_srv.DeleteDepartmentRequest,
) (resp *identity_srv.DeleteDepartmentResponse, err error) {
	if err := s.requirePerm(ctx, "delete", "department:"+req.GetDepartmentID()); err != nil {
		return nil, err
	}

	if err := s.logic.DeleteDepartment(ctx, req.GetDepartmentID()); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.DeleteDepartmentResponse{}, nil
}

// GetOrganizationDepartments implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetOrganizationDepartments(
	ctx context.Context,
	req *identity_srv.GetOrganizationDepartmentsRequest,
) (resp *identity_srv.GetOrganizationDepartmentsResponse, err error) {
	resp, err = s.logic.GetDepartmentsByOrganization(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// GetMembership implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetMembership(
	ctx context.Context,
	req *identity_srv.GetMembershipRequest,
) (resp *identity_srv.GetMembershipResponse, err error) {
	membership, err := s.logic.GetMembership(ctx, req.GetMembershipID())
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.GetMembershipResponse{Membership: membership}, nil
}

// GetPrimaryMembership implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetPrimaryMembership(
	ctx context.Context,
	req *identity_srv.GetPrimaryMembershipRequest,
) (resp *identity_srv.GetPrimaryMembershipResponse, err error) {
	membership, err := s.logic.GetPrimaryMembership(ctx, req.GetUserID())
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.GetPrimaryMembershipResponse{Membership: membership}, nil
}

// CheckMembership implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CheckMembership(
	ctx context.Context,
	req *identity_srv.CheckMembershipRequest,
) (resp *identity_srv.CheckMembershipResponse, err error) {
	isMember, err := s.logic.CheckMembership(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.CheckMembershipResponse{IsMember: &isMember}, nil
}

// UploadTemporaryLogo implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UploadTemporaryLogo(
	ctx context.Context,
	req *identity_srv.UploadTemporaryLogoRequest,
) (resp *identity_srv.UploadTemporaryLogoResponse, err error) {
	if err := s.requirePerm(ctx, "upload", "logo"); err != nil {
		return nil, err
	}

	organizationLogo, err := s.logic.UploadTemporaryLogo(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.UploadTemporaryLogoResponse{OrganizationLogo: organizationLogo}, nil
}

// GetOrganizationLogo implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetOrganizationLogo(
	ctx context.Context,
	req *identity_srv.GetOrganizationLogoRequest,
) (resp *identity_srv.GetOrganizationLogoResponse, err error) {
	organizationLogo, err := s.logic.GetOrganizationLogo(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.GetOrganizationLogoResponse{OrganizationLogo: organizationLogo}, nil
}

// DeleteOrganizationLogo implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteOrganizationLogo(
	ctx context.Context,
	req *identity_srv.DeleteOrganizationLogoRequest,
) (resp *identity_srv.DeleteOrganizationLogoResponse, err error) {
	if err := s.requirePerm(ctx, "delete", "logo:"+derefStr(req.LogoID)); err != nil {
		return nil, err
	}

	if err := s.logic.DeleteOrganizationLogo(ctx, req); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.DeleteOrganizationLogoResponse{}, nil
}

// BindLogoToOrganization implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) BindLogoToOrganization(
	ctx context.Context,
	req *identity_srv.BindLogoToOrganizationRequest,
) (resp *identity_srv.BindLogoToOrganizationResponse, err error) {
	if err := s.requirePerm(ctx, "bind", "logo:"+derefStr(req.OrganizationID)); err != nil {
		return nil, err
	}

	organizationLogo, err := s.logic.BindLogoToOrganization(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.BindLogoToOrganizationResponse{OrganizationLogo: organizationLogo}, nil
}

// CreateRoleDefinition implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CreateRoleDefinition(
	ctx context.Context,
	req *identity_srv.RoleDefinitionCreateRequest,
) (resp *identity_srv.CreateRoleDefinitionResponse, err error) {
	if err := s.requirePerm(ctx, "create", "role"); err != nil {
		return nil, err
	}

	roleDefinition, err := s.logic.CreateRoleDefinition(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.CreateRoleDefinitionResponse{RoleDefinition: roleDefinition}, nil
}

// UpdateRoleDefinition implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UpdateRoleDefinition(
	ctx context.Context,
	req *identity_srv.RoleDefinitionUpdateRequest,
) (resp *identity_srv.UpdateRoleDefinitionResponse, err error) {
	if err := s.requirePerm(ctx, "update", "role:"+derefStr(req.RoleDefinitionID)); err != nil {
		return nil, err
	}

	roleDefinition, err := s.logic.UpdateRoleDefinition(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.UpdateRoleDefinitionResponse{RoleDefinition: roleDefinition}, nil
}

// DeleteRoleDefinition implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteRoleDefinition(
	ctx context.Context,
	req *identity_srv.DeleteRoleDefinitionRequest,
) (resp *identity_srv.DeleteRoleDefinitionResponse, err error) {
	if err := s.requirePerm(ctx, "delete", "role:"+req.GetRoleID()); err != nil {
		return nil, err
	}

	if err := s.logic.DeleteRoleDefinition(ctx, req.GetRoleID()); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.DeleteRoleDefinitionResponse{}, nil
}

// GetRoleDefinition implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetRoleDefinition(
	ctx context.Context,
	req *identity_srv.GetRoleDefinitionRequest,
) (resp *identity_srv.GetRoleDefinitionResponse, err error) {
	roleDefinition, err := s.logic.GetRoleDefinition(ctx, req.GetRoleID())
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.GetRoleDefinitionResponse{RoleDefinition: roleDefinition}, nil
}

// ListRoleDefinitions implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ListRoleDefinitions(
	ctx context.Context,
	req *identity_srv.RoleDefinitionQueryRequest,
) (resp *identity_srv.RoleDefinitionListResponse, err error) {
	resp, err = s.logic.ListRoleDefinitions(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// AssignRoleToUser implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) AssignRoleToUser(
	ctx context.Context,
	req *identity_srv.AssignRoleToUserRequest,
) (resp *identity_srv.UserRoleAssignmentResponse, err error) {
	if err := s.requirePerm(ctx, "assign", "role_assignment"); err != nil {
		return nil, err
	}

	resp, err = s.logic.AssignRoleToUser(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// UpdateUserRoleAssignment implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UpdateUserRoleAssignment(
	ctx context.Context,
	req *identity_srv.UpdateUserRoleAssignmentRequest,
) (resp *identity_srv.UpdateUserRoleAssignmentResponse, err error) {
	resource := "role_assignment:" + derefStr(req.UserID) + ":" + derefStr(req.RoleID)
	if err := s.requirePerm(ctx, "update", resource); err != nil {
		return nil, err
	}

	if err := s.logic.UpdateUserRoleAssignment(ctx, req); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.UpdateUserRoleAssignmentResponse{}, nil
}

// RevokeRoleFromUser implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) RevokeRoleFromUser(
	ctx context.Context,
	req *identity_srv.RevokeRoleFromUserRequest,
) (resp *identity_srv.RevokeRoleFromUserResponse, err error) {
	resource := "role_assignment:" + derefStr(req.UserID) + ":" + derefStr(req.RoleID)
	if err := s.requirePerm(ctx, "revoke", resource); err != nil {
		return nil, err
	}

	if err := s.logic.RevokeRoleFromUser(ctx, req); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.RevokeRoleFromUserResponse{}, nil
}

// GetLastUserRoleAssignment implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetLastUserRoleAssignment(
	ctx context.Context,
	req *identity_srv.GetLastUserRoleAssignmentRequest,
) (resp *identity_srv.GetLastUserRoleAssignmentResponse, err error) {
	assignment, err := s.logic.GetLastUserRoleAssignment(ctx, req.GetUserID())
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.GetLastUserRoleAssignmentResponse{UserRoleAssignment: assignment}, nil
}

// ListUserRoleAssignments implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ListUserRoleAssignments(
	ctx context.Context,
	req *identity_srv.UserRoleQueryRequest,
) (resp *identity_srv.UserRoleListResponse, err error) {
	resp, err = s.logic.ListUserRoleAssignments(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// GetUsersByRole implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetUsersByRole(
	ctx context.Context,
	req *identity_srv.GetUsersByRoleRequest,
) (resp *identity_srv.GetUsersByRoleResponse, err error) {
	resp, err = s.logic.GetUsersByRole(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// BatchBindUsersToRole implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) BatchBindUsersToRole(
	ctx context.Context,
	req *identity_srv.BatchBindUsersToRoleRequest,
) (resp *identity_srv.BatchBindUsersToRoleResponse, err error) {
	if err := s.requirePerm(ctx, "assign", "role_assignment"); err != nil {
		return nil, err
	}

	resp, err = s.logic.BatchBindUsersToRole(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// BatchGetUserRoles implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) BatchGetUserRoles(
	ctx context.Context,
	req *identity_srv.BatchGetUserRolesRequest,
) (resp *identity_srv.BatchGetUserRolesResponse, err error) {
	resp, err = s.logic.BatchGetUserRoles(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// UploadMenu implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UploadMenu(
	ctx context.Context,
	req *identity_srv.UploadMenuRequest,
) (resp *identity_srv.UploadMenuResponse, err error) {
	if err := s.requirePerm(ctx, "upload", "menu"); err != nil {
		return nil, err
	}

	if err := s.logic.UploadMenu(ctx, req); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.UploadMenuResponse{}, nil
}

// GetMenuTree implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetMenuTree(
	ctx context.Context,
	_ *identity_srv.GetMenuTreeRequest,
) (resp *identity_srv.GetMenuTreeResponse, err error) {
	resp, err = s.logic.GetMenuTree(ctx)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// ConfigureRoleMenus implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ConfigureRoleMenus(
	ctx context.Context,
	req *identity_srv.ConfigureRoleMenusRequest,
) (resp *identity_srv.ConfigureRoleMenusResponse, err error) {
	if err := s.requirePerm(ctx, "configure", "menu:role:"+derefStr(req.RoleID)); err != nil {
		return nil, err
	}

	resp, err = s.logic.ConfigureRoleMenus(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// GetRoleMenuTree implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetRoleMenuTree(
	ctx context.Context,
	req *identity_srv.GetRoleMenuTreeRequest,
) (resp *identity_srv.GetRoleMenuTreeResponse, err error) {
	resp, err = s.logic.GetRoleMenuTree(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// GetUserMenuTree implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetUserMenuTree(
	ctx context.Context,
	req *identity_srv.GetUserMenuTreeRequest,
) (resp *identity_srv.GetUserMenuTreeResponse, err error) {
	resp, err = s.logic.GetUserMenuTree(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// GetRoleMenuPermissions implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetRoleMenuPermissions(
	ctx context.Context,
	req *identity_srv.GetRoleMenuPermissionsRequest,
) (resp *identity_srv.GetRoleMenuPermissionsResponse, err error) {
	resp, err = s.logic.GetRoleMenuPermissions(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// HasMenuPermission implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) HasMenuPermission(
	ctx context.Context,
	req *identity_srv.HasMenuPermissionRequest,
) (resp *identity_srv.HasMenuPermissionResponse, err error) {
	resp, err = s.logic.HasMenuPermission(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// GetUserMenuPermissions implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetUserMenuPermissions(
	ctx context.Context,
	req *identity_srv.GetUserMenuPermissionsRequest,
) (resp *identity_srv.GetUserMenuPermissionsResponse, err error) {
	resp, err = s.logic.GetUserMenuPermissions(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// CreateAuditLog implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CreateAuditLog(
	ctx context.Context,
	req *identity_srv.CreateAuditLogRequest,
) (resp *identity_srv.CreateAuditLogResponse, err error) {
	if err := s.logic.CreateAuditLog(ctx, req); err != nil {
		return nil, errno.ToKitexError(err)
	}

	return &identity_srv.CreateAuditLogResponse{}, nil
}

// ListAuditLogs implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ListAuditLogs(
	ctx context.Context,
	req *identity_srv.ListAuditLogsRequest,
) (resp *identity_srv.ListAuditLogsResponse, err error) {
	if err := s.requirePerm(ctx, "read", "audit_log"); err != nil {
		return nil, err
	}

	resp, err = s.logic.ListAuditLogs(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}
