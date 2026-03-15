package main

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic"
	core "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
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
) (resp *identity_srv.UserProfile, err error) {
	resp, err = s.logic.CreateUser(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// GetUser implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetUser(
	ctx context.Context,
	req *identity_srv.GetUserRequest,
) (resp *identity_srv.UserProfile, err error) {
	resp, err = s.logic.GetUser(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// UpdateUser implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UpdateUser(
	ctx context.Context,
	req *identity_srv.UpdateUserRequest,
) (resp *identity_srv.UserProfile, err error) {
	resp, err = s.logic.UpdateUser(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// DeleteUser implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteUser(
	ctx context.Context,
	req *identity_srv.DeleteUserRequest,
) (err error) {
	err = s.logic.DeleteUser(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
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
) (err error) {
	err = s.logic.ChangeUserStatus(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
}

// UnlockUser implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UnlockUser(
	ctx context.Context,
	req *identity_srv.UnlockUserRequest,
) (err error) {
	err = s.logic.UnlockUser(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
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
) (err error) {
	err = s.logic.ChangePassword(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
}

// ResetPassword implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ResetPassword(
	ctx context.Context,
	req *identity_srv.ResetPasswordRequest,
) (err error) {
	err = s.logic.ResetPassword(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
}

// ForcePasswordChange implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) ForcePasswordChange(
	ctx context.Context,
	req *identity_srv.ForcePasswordChangeRequest,
) (err error) {
	err = s.logic.ForcePasswordChange(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
}

// ===========================================================================
// OrgManagement
// ===========================================================================

// AddMembership implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) AddMembership(
	ctx context.Context,
	req *identity_srv.AddMembershipRequest,
) (resp *identity_srv.UserMembership, err error) {
	resp, err = s.logic.AddMembership(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// UpdateMembership implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UpdateMembership(
	ctx context.Context,
	req *identity_srv.UpdateMembershipRequest,
) (resp *identity_srv.UserMembership, err error) {
	resp, err = s.logic.UpdateMembership(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// RemoveMembership implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) RemoveMembership(
	ctx context.Context,
	membershipID string,
) (err error) {
	err = s.logic.RemoveMembership(ctx, membershipID)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
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
) (resp *identity_srv.Organization, err error) {
	resp, err = s.logic.CreateOrganization(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// GetOrganization implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetOrganization(
	ctx context.Context,
	req *identity_srv.GetOrganizationRequest,
) (resp *identity_srv.Organization, err error) {
	resp, err = s.logic.GetOrganization(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// UpdateOrganization implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UpdateOrganization(
	ctx context.Context,
	req *identity_srv.UpdateOrganizationRequest,
) (resp *identity_srv.Organization, err error) {
	resp, err = s.logic.UpdateOrganization(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// DeleteOrganization implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteOrganization(
	ctx context.Context,
	organizationID string,
) (err error) {
	err = s.logic.DeleteOrganization(ctx, organizationID)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
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
) (resp *identity_srv.Department, err error) {
	resp, err = s.logic.CreateDepartment(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// GetDepartment implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetDepartment(
	ctx context.Context,
	req *identity_srv.GetDepartmentRequest,
) (resp *identity_srv.Department, err error) {
	resp, err = s.logic.GetDepartment(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// UpdateDepartment implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UpdateDepartment(
	ctx context.Context,
	req *identity_srv.UpdateDepartmentRequest,
) (resp *identity_srv.Department, err error) {
	resp, err = s.logic.UpdateDepartment(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// DeleteDepartment implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteDepartment(
	ctx context.Context,
	departmentID string,
) (err error) {
	err = s.logic.DeleteDepartment(ctx, departmentID)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
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
	membershipID core.UUID,
) (resp *identity_srv.UserMembership, err error) {
	resp, err = s.logic.GetMembership(ctx, membershipID)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// GetPrimaryMembership implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetPrimaryMembership(
	ctx context.Context,
	userID core.UUID,
) (resp *identity_srv.UserMembership, err error) {
	resp, err = s.logic.GetPrimaryMembership(ctx, userID)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// CheckMembership implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CheckMembership(
	ctx context.Context,
	req *identity_srv.CheckMembershipRequest,
) (resp bool, err error) {
	resp, err = s.logic.CheckMembership(ctx, req)
	if err != nil {
		return false, errno.ToKitexError(err)
	}

	return resp, nil
}

// UploadTemporaryLogo implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UploadTemporaryLogo(
	ctx context.Context,
	req *identity_srv.UploadTemporaryLogoRequest,
) (resp *identity_srv.OrganizationLogo, err error) {
	resp, err = s.logic.UploadTemporaryLogo(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// GetOrganizationLogo implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetOrganizationLogo(
	ctx context.Context,
	req *identity_srv.GetOrganizationLogoRequest,
) (resp *identity_srv.OrganizationLogo, err error) {
	resp, err = s.logic.GetOrganizationLogo(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// DeleteOrganizationLogo implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteOrganizationLogo(
	ctx context.Context,
	req *identity_srv.DeleteOrganizationLogoRequest,
) (err error) {
	err = s.logic.DeleteOrganizationLogo(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
}

// BindLogoToOrganization implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) BindLogoToOrganization(
	ctx context.Context,
	req *identity_srv.BindLogoToOrganizationRequest,
) (resp *identity_srv.OrganizationLogo, err error) {
	resp, err = s.logic.BindLogoToOrganization(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// CreateRoleDefinition implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) CreateRoleDefinition(
	ctx context.Context,
	req *identity_srv.RoleDefinitionCreateRequest,
) (resp *identity_srv.RoleDefinition, err error) {
	resp, err = s.logic.CreateRoleDefinition(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// UpdateRoleDefinition implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) UpdateRoleDefinition(
	ctx context.Context,
	req *identity_srv.RoleDefinitionUpdateRequest,
) (resp *identity_srv.RoleDefinition, err error) {
	resp, err = s.logic.UpdateRoleDefinition(ctx, req)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
}

// DeleteRoleDefinition implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) DeleteRoleDefinition(
	ctx context.Context,
	roleID core.UUID,
) (err error) {
	err = s.logic.DeleteRoleDefinition(ctx, roleID)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
}

// GetRoleDefinition implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetRoleDefinition(
	ctx context.Context,
	roleID core.UUID,
) (resp *identity_srv.RoleDefinition, err error) {
	resp, err = s.logic.GetRoleDefinition(ctx, roleID)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
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
) (err error) {
	err = s.logic.UpdateUserRoleAssignment(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
}

// RevokeRoleFromUser implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) RevokeRoleFromUser(
	ctx context.Context,
	req *identity_srv.RevokeRoleFromUserRequest,
) (err error) {
	err = s.logic.RevokeRoleFromUser(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
}

// GetLastUserRoleAssignment implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetLastUserRoleAssignment(
	ctx context.Context,
	userID core.UUID,
) (resp *identity_srv.UserRoleAssignment, err error) {
	resp, err = s.logic.GetLastUserRoleAssignment(ctx, userID)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	return resp, nil
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
) (err error) {
	err = s.logic.UploadMenu(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
}

// GetMenuTree implements the IdentityServiceImpl interface.
func (s *IdentityServiceImpl) GetMenuTree(
	ctx context.Context,
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
func (s *IdentityServiceImpl) SyncPolicies(ctx context.Context) (resp *identity_srv.SyncPoliciesResponse, err error) {
	result, err := s.logic.SyncPolicies(ctx)
	if err != nil {
		return nil, errno.ToKitexError(err)
	}

	resp = &identity_srv.SyncPoliciesResponse{
		Success:              &result.Success,
		RolePolicyCount:      &result.RolePolicyCount,
		UserRoleBindingCount: &result.UserRoleBindingCount,
		RoleInheritanceCount: &result.RoleInheritanceCount,
		Message:              &result.Message,
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
func (s *IdentityServiceImpl) CreateAuditLog(ctx context.Context, req *identity_srv.CreateAuditLogRequest) (err error) {
	err = s.logic.CreateAuditLog(ctx, req)
	if err != nil {
		return errno.ToKitexError(err)
	}

	return nil
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
