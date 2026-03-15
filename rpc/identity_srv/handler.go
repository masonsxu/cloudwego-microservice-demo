package main

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic"
	identity_srv "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/wire"
)

// IdentityServiceImpl implements the last service interface defined in the IDL.
type IdentityServiceImpl struct {
	logic logic.Logic
}

// NewIdentityServiceImpl 从应用容器创建 IdentityServiceImpl 实例
func NewIdentityServiceImpl(container *wire.AppContainer) *IdentityServiceImpl {
	return &IdentityServiceImpl{
		logic: container.Logic,
	}
}

// ===========================================================================
// UserProfile
// ===========================================================================

// CreateUser implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CreateUser(
	ctx context.Context,
	req *identity_srv.CreateUserRequest,
) (resp *identity_srv.CreateUserResponse, err error) {
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
func (s *IdentityServiceImpl) ChangePassword(
	ctx context.Context,
	req *identity_srv.ChangePasswordRequest,
) (resp *identity_srv.ChangePasswordResponse, err error) {
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

// CheckPermission implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CheckPermission(
	ctx context.Context,
	req *identity_srv.CheckPermissionRequest,
) (resp *identity_srv.CheckPermissionResponse, err error) {
	userID := ""
	if req.UserID != nil {
		userID = *req.UserID
	}

	resource := ""
	if req.Resource != nil {
		resource = *req.Resource
	}

	action := ""
	if req.Action != nil {
		action = *req.Action
	}

	result, err := s.logic.CheckPermission(ctx, userID, req.RoleIDs, req.DepartmentIDs, resource, action)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	resp = &identity_srv.CheckPermissionResponse{
		Allowed:   &result.Allowed,
		DataScope: &result.DataScope,
		UserID:    req.UserID,
		Resource:  req.Resource,
		Action:    req.Action,
	}

	return resp, nil
}

// SyncPolicies implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) SyncPolicies(
	ctx context.Context,
	_ *identity_srv.SyncPoliciesRequest,
) (resp *identity_srv.SyncPoliciesResponse, err error) {
	result, err := s.logic.SyncPolicies(ctx)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	policyRules := make([]*identity_srv.CasbinPolicyRule, 0, len(result.Policies))
	for _, rule := range result.Policies {
		policyRules = append(policyRules, &identity_srv.CasbinPolicyRule{Values: rule})
	}

	groupingRules := make([]*identity_srv.CasbinPolicyRule, 0, len(result.GroupingPolicies))
	for _, rule := range result.GroupingPolicies {
		groupingRules = append(groupingRules, &identity_srv.CasbinPolicyRule{Values: rule})
	}

	inheritanceRules := make([]*identity_srv.CasbinPolicyRule, 0, len(result.RoleInheritancePolicies))
	for _, rule := range result.RoleInheritancePolicies {
		inheritanceRules = append(inheritanceRules, &identity_srv.CasbinPolicyRule{Values: rule})
	}

	resp = &identity_srv.SyncPoliciesResponse{
		Success:                 &result.Success,
		RolePolicyCount:         &result.RolePolicyCount,
		UserRoleBindingCount:    &result.UserRoleBindingCount,
		RoleInheritanceCount:    &result.RoleInheritanceCount,
		Message:                 &result.Message,
		Policies:                policyRules,
		GroupingPolicies:        groupingRules,
		RoleInheritancePolicies: inheritanceRules,
	}

	return resp, nil
}

// GetUserDataScope implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetUserDataScope(
	ctx context.Context,
	req *identity_srv.GetUserDataScopeRequest,
) (resp *identity_srv.GetUserDataScopeResponse, err error) {
	userID := ""
	if req.UserID != nil {
		userID = *req.UserID
	}

	resource := ""
	if req.Resource != nil {
		resource = *req.Resource
	}

	action := ""
	if req.Action != nil {
		action = *req.Action
	}

	dataScope, err := s.logic.GetUserDataScope(ctx, userID, resource, action)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	resp = &identity_srv.GetUserDataScopeResponse{
		DataScope: &dataScope,
		UserID:    req.UserID,
		Resource:  req.Resource,
		Action:    req.Action,
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
	resp, err = s.logic.ListAuditLogs(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// CreateOAuth2Client implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CreateOAuth2Client(
	ctx context.Context,
	req *identity_srv.CreateOAuth2ClientRequest,
) (resp *identity_srv.CreateOAuth2ClientResponse, err error) {
	resp, err = s.logic.CreateOAuth2Client(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// GetOAuth2Client implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetOAuth2Client(
	ctx context.Context,
	req *identity_srv.GetOAuth2ClientRequest,
) (resp *identity_srv.OAuth2Client, err error) {
	resp, err = s.logic.GetOAuth2Client(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// UpdateOAuth2Client implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UpdateOAuth2Client(
	ctx context.Context,
	req *identity_srv.UpdateOAuth2ClientRequest,
) (resp *identity_srv.OAuth2Client, err error) {
	resp, err = s.logic.UpdateOAuth2Client(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// DeleteOAuth2Client implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteOAuth2Client(
	ctx context.Context,
	req *identity_srv.DeleteOAuth2ClientRequest,
) (err error) {
	err = s.logic.DeleteOAuth2Client(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
}

// ListOAuth2Clients implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ListOAuth2Clients(
	ctx context.Context,
	req *identity_srv.ListOAuth2ClientsRequest,
) (resp *identity_srv.ListOAuth2ClientsResponse, err error) {
	resp, err = s.logic.ListOAuth2Clients(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// RotateOAuth2ClientSecret implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) RotateOAuth2ClientSecret(
	ctx context.Context,
	req *identity_srv.RotateOAuth2ClientSecretRequest,
) (resp *identity_srv.RotateOAuth2ClientSecretResponse, err error) {
	resp, err = s.logic.RotateOAuth2ClientSecret(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// GetOAuth2ClientForAuth implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetOAuth2ClientForAuth(
	ctx context.Context,
	clientID string,
) (resp *identity_srv.GetOAuth2ClientForAuthResponse, err error) {
	resp, err = s.logic.GetOAuth2ClientForAuth(ctx, clientID)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// ListOAuth2Scopes implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ListOAuth2Scopes(
	ctx context.Context,
	req *identity_srv.ListOAuth2ScopesRequest,
) (resp *identity_srv.ListOAuth2ScopesResponse, err error) {
	resp, err = s.logic.ListOAuth2Scopes(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// SaveOAuth2Consent implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) SaveOAuth2Consent(
	ctx context.Context,
	req *identity_srv.SaveOAuth2ConsentRequest,
) (err error) {
	err = s.logic.SaveOAuth2Consent(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
}

// GetOAuth2Consent implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetOAuth2Consent(
	ctx context.Context,
	req *identity_srv.GetOAuth2ConsentRequest,
) (resp *identity_srv.GetOAuth2ConsentResponse, err error) {
	resp, err = s.logic.GetOAuth2Consent(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// ListOAuth2Consents implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ListOAuth2Consents(
	ctx context.Context,
	req *identity_srv.ListOAuth2ConsentsRequest,
) (resp *identity_srv.ListOAuth2ConsentsResponse, err error) {
	resp, err = s.logic.ListOAuth2Consents(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// RevokeOAuth2Consent implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) RevokeOAuth2Consent(
	ctx context.Context,
	req *identity_srv.RevokeOAuth2ConsentRequest,
) (err error) {
	err = s.logic.RevokeOAuth2Consent(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
}

// ===========================================================================
// OAuth2 Token Storage (供 fosite 存储层通过 RPC 调用)
// ===========================================================================

// CreateOAuth2AuthorizeCodeSession implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CreateOAuth2AuthorizeCodeSession(
	ctx context.Context,
	req *identity_srv.OAuth2TokenSession,
) (err error) {
	// TODO: 实现 Token 存储 RPC 接口（阶段三：fosite 存储持久化）
	return nil
}

// GetOAuth2AuthorizeCodeSession implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetOAuth2AuthorizeCodeSession(
	ctx context.Context,
	signature string,
) (resp *identity_srv.OAuth2TokenSession, err error) {
	return nil, nil
}

// InvalidateOAuth2AuthorizeCodeSession implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) InvalidateOAuth2AuthorizeCodeSession(ctx context.Context, signature string) (err error) {
	return nil
}

// CreateOAuth2AccessTokenSession implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CreateOAuth2AccessTokenSession(
	ctx context.Context,
	req *identity_srv.OAuth2TokenSession,
) (err error) {
	return nil
}

// GetOAuth2AccessTokenSession implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetOAuth2AccessTokenSession(
	ctx context.Context,
	signature string,
) (resp *identity_srv.OAuth2TokenSession, err error) {
	return nil, nil
}

// DeleteOAuth2AccessTokenSession implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteOAuth2AccessTokenSession(ctx context.Context, signature string) (err error) {
	return nil
}

// RevokeOAuth2AccessToken implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) RevokeOAuth2AccessToken(ctx context.Context, requestID string) (err error) {
	return nil
}

// CreateOAuth2RefreshTokenSession implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CreateOAuth2RefreshTokenSession(
	ctx context.Context,
	req *identity_srv.OAuth2TokenSession,
) (err error) {
	return nil
}

// GetOAuth2RefreshTokenSession implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetOAuth2RefreshTokenSession(
	ctx context.Context,
	signature string,
) (resp *identity_srv.OAuth2TokenSession, err error) {
	return nil, nil
}

// DeleteOAuth2RefreshTokenSession implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteOAuth2RefreshTokenSession(ctx context.Context, signature string) (err error) {
	return nil
}

// RevokeOAuth2RefreshToken implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) RevokeOAuth2RefreshToken(ctx context.Context, requestID string) (err error) {
	return nil
}

// CreateOAuth2PKCESession implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CreateOAuth2PKCESession(
	ctx context.Context,
	req *identity_srv.OAuth2TokenSession,
) (err error) {
	return nil
}

// GetOAuth2PKCESession implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetOAuth2PKCESession(
	ctx context.Context,
	signature string,
) (resp *identity_srv.OAuth2TokenSession, err error) {
	return nil, nil
}

// DeleteOAuth2PKCESession implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteOAuth2PKCESession(ctx context.Context, signature string) (err error) {
	return nil
}
