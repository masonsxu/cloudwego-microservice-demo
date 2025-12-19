package identity

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
)

// identityServiceImpl 身份管理聚合服务实现
// 实现所有子服务接口，提供统一的服务入口
type identityServiceImpl struct {
	authService       AuthService
	userService       UserService
	membershipService MembershipService
	orgService        OrganizationService
	deptService       DepartmentService
	logoService       LogoService
}

// NewService 创建身份管理聚合服务
func NewService(
	authService AuthService,
	userService UserService,
	membershipService MembershipService,
	orgService OrganizationService,
	deptService DepartmentService,
	logoService LogoService,
) Service {
	return &identityServiceImpl{
		authService:       authService,
		userService:       userService,
		membershipService: membershipService,
		orgService:        orgService,
		deptService:       deptService,
		logoService:       logoService,
	}
}

// =================================================================
// AuthService 接口实现 - 委托给 authService
// =================================================================

func (s *identityServiceImpl) Login(
	ctx context.Context,
	req *identity.LoginRequestDTO,
) (*identity.LoginResponseDTO, Permission, error) {
	return s.authService.Login(ctx, req)
}

func (s *identityServiceImpl) ChangePassword(
	ctx context.Context,
	req *identity.ChangePasswordRequestDTO,
	userID string,
) (*http_base.OperationStatusResponseDTO, error) {
	return s.authService.ChangePassword(ctx, req, userID)
}

func (s *identityServiceImpl) ResetPassword(
	ctx context.Context,
	req *identity.ResetPasswordRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	return s.authService.ResetPassword(ctx, req)
}

func (s *identityServiceImpl) ForcePasswordChange(
	ctx context.Context,
	req *identity.ForcePasswordChangeRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	return s.authService.ForcePasswordChange(ctx, req)
}

// =================================================================
// UserService 接口实现 - 委托给 userService
// =================================================================

func (s *identityServiceImpl) CreateUser(
	ctx context.Context,
	req *identity.CreateUserRequestDTO,
	operatorID string,
) (*identity.UserProfileResponseDTO, error) {
	return s.userService.CreateUser(ctx, req, operatorID)
}

func (s *identityServiceImpl) GetUser(
	ctx context.Context,
	req *identity.GetUserRequestDTO,
) (*identity.UserProfileResponseDTO, error) {
	return s.userService.GetUser(ctx, req)
}

func (s *identityServiceImpl) GetMe(
	ctx context.Context,
	userID string,
) (*identity.UserProfileResponseDTO, error) {
	return s.userService.GetMe(ctx, userID)
}

func (s *identityServiceImpl) UpdateUser(
	ctx context.Context,
	req *identity.UpdateUserRequestDTO,
	operatorID string,
) (*identity.UserProfileResponseDTO, error) {
	return s.userService.UpdateUser(ctx, req, operatorID)
}

func (s *identityServiceImpl) UpdateMe(
	ctx context.Context,
	req *identity.UpdateMeRequestDTO,
	userID string,
) (*identity.UserProfileResponseDTO, error) {
	return s.userService.UpdateMe(ctx, req, userID)
}

func (s *identityServiceImpl) DeleteUser(
	ctx context.Context,
	req *identity.DeleteUserRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	return s.userService.DeleteUser(ctx, req)
}

func (s *identityServiceImpl) ListUsers(
	ctx context.Context,
	req *identity.ListUsersRequestDTO,
) (*identity.ListUsersResponseDTO, error) {
	return s.userService.ListUsers(ctx, req)
}

func (s *identityServiceImpl) SearchUsers(
	ctx context.Context,
	req *identity.SearchUsersRequestDTO,
) (*identity.SearchUsersResponseDTO, error) {
	return s.userService.SearchUsers(ctx, req)
}

func (s *identityServiceImpl) ChangeUserStatus(
	ctx context.Context,
	req *identity.ChangeUserStatusRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	return s.userService.ChangeUserStatus(ctx, req)
}

func (s *identityServiceImpl) UnlockUser(
	ctx context.Context,
	req *identity.UnlockUserRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	return s.userService.UnlockUser(ctx, req)
}

// =================================================================
// MembershipService 接口实现 - 委托给 membershipService
// =================================================================

func (s *identityServiceImpl) GetUserMemberships(
	ctx context.Context,
	req *identity.GetUserMembershipsRequestDTO,
) (*identity.GetUserMembershipsResponseDTO, error) {
	return s.membershipService.GetUserMemberships(ctx, req)
}

// GetPrimaryMembership 获取用户主成员关系 - 获取用户在当前组织中的主成员关系
func (s *identityServiceImpl) GetPrimaryMembership(
	ctx context.Context,
	req *identity.GetPrimaryMembershipRequestDTO,
) (*identity.UserMembershipResponseDTO, error) {
	return s.membershipService.GetPrimaryMembership(ctx, req)
}

// =================================================================
// OrganizationService 接口实现 - 委托给 orgService
// =================================================================

func (s *identityServiceImpl) CreateOrganization(
	ctx context.Context,
	req *identity.CreateOrganizationRequestDTO,
) (*identity.OrganizationResponseDTO, error) {
	return s.orgService.CreateOrganization(ctx, req)
}

func (s *identityServiceImpl) GetOrganization(
	ctx context.Context,
	req *identity.GetOrganizationRequestDTO,
) (*identity.OrganizationResponseDTO, error) {
	return s.orgService.GetOrganization(ctx, req)
}

func (s *identityServiceImpl) UpdateOrganization(
	ctx context.Context,
	req *identity.UpdateOrganizationRequestDTO,
) (*identity.OrganizationResponseDTO, error) {
	return s.orgService.UpdateOrganization(ctx, req)
}

func (s *identityServiceImpl) DeleteOrganization(
	ctx context.Context,
	req *identity.DeleteOrganizationRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	return s.orgService.DeleteOrganization(ctx, req)
}

func (s *identityServiceImpl) ListOrganizations(
	ctx context.Context,
	req *identity.ListOrganizationsRequestDTO,
) (*identity.ListOrganizationsResponseDTO, error) {
	return s.orgService.ListOrganizations(ctx, req)
}

// =================================================================
// DepartmentService 接口实现 - 委托给 deptService
// =================================================================

func (s *identityServiceImpl) CreateDepartment(
	ctx context.Context,
	req *identity.CreateDepartmentRequestDTO,
) (*identity.DepartmentResponseDTO, error) {
	return s.deptService.CreateDepartment(ctx, req)
}

func (s *identityServiceImpl) GetDepartment(
	ctx context.Context,
	req *identity.GetDepartmentRequestDTO,
) (*identity.DepartmentResponseDTO, error) {
	return s.deptService.GetDepartment(ctx, req)
}

func (s *identityServiceImpl) UpdateDepartment(
	ctx context.Context,
	req *identity.UpdateDepartmentRequestDTO,
) (*identity.DepartmentResponseDTO, error) {
	return s.deptService.UpdateDepartment(ctx, req)
}

func (s *identityServiceImpl) DeleteDepartment(
	ctx context.Context,
	req *identity.DeleteDepartmentRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	return s.deptService.DeleteDepartment(ctx, req)
}

func (s *identityServiceImpl) GetOrganizationDepartments(
	ctx context.Context,
	req *identity.GetOrganizationDepartmentsRequestDTO,
) (*identity.GetOrganizationDepartmentsResponseDTO, error) {
	return s.deptService.GetOrganizationDepartments(ctx, req)
}

// =================================================================
// LogoService 接口实现 - 委托给 logoService
// =================================================================

func (s *identityServiceImpl) UploadTemporaryLogo(
	ctx context.Context,
	req *identity.UploadTemporaryLogoRequestDTO,
	userID string,
) (*identity.OrganizationLogoResponseDTO, error) {
	return s.logoService.UploadTemporaryLogo(ctx, req, userID)
}

func (s *identityServiceImpl) GetOrganizationLogo(
	ctx context.Context,
	req *identity.GetOrganizationLogoRequestDTO,
) (*identity.OrganizationLogoResponseDTO, error) {
	return s.logoService.GetOrganizationLogo(ctx, req)
}

func (s *identityServiceImpl) DeleteOrganizationLogo(
	ctx context.Context,
	req *identity.DeleteOrganizationLogoRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	return s.logoService.DeleteOrganizationLogo(ctx, req)
}

func (s *identityServiceImpl) BindLogoToOrganization(
	ctx context.Context,
	req *identity.BindLogoToOrganizationRequestDTO,
) (*identity.OrganizationResponseDTO, error) {
	return s.logoService.BindLogoToOrganization(ctx, req)
}
