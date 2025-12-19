package identity

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
)

type Permission string

// Service 身份管理聚合服务接口 - 统一暴露给Handler层
type Service interface {
	AuthService
	UserService
	MembershipService
	OrganizationService
	DepartmentService
	LogoService
}

// =================================================================
// 专门化服务接口定义 - 按业务领域划分
// =================================================================

// AuthService 身份认证服务接口
type AuthService interface {
	// Login 用户登录 - 验证用户凭据并返回访问令牌和用户信息
	Login(
		ctx context.Context,
		req *identity.LoginRequestDTO,
	) (*identity.LoginResponseDTO, Permission, error)

	// ChangePassword 修改密码 - 用户修改自己的密码（需要提供旧密码）
	ChangePassword(
		ctx context.Context,
		req *identity.ChangePasswordRequestDTO,
		userID string,
	) (*http_base.OperationStatusResponseDTO, error)

	// ResetPassword 重置密码 - 管理员重置用户密码（管理员权限）
	ResetPassword(
		ctx context.Context,
		req *identity.ResetPasswordRequestDTO,
	) (*http_base.OperationStatusResponseDTO, error)

	// ForcePasswordChange 强制下次登录修改密码 - 管理员标记用户需要在下次登录时强制修改密码
	ForcePasswordChange(
		ctx context.Context,
		req *identity.ForcePasswordChangeRequestDTO,
	) (*http_base.OperationStatusResponseDTO, error)
}

// UserService 用户管理服务接口
type UserService interface {
	// CreateUser 创建用户 - 管理员创建新用户账户
	CreateUser(
		ctx context.Context,
		req *identity.CreateUserRequestDTO,
		operatorID string,
	) (*identity.UserProfileResponseDTO, error)

	// GetUser 获取用户信息 - 根据用户ID获取用户详细信息
	GetUser(
		ctx context.Context,
		req *identity.GetUserRequestDTO,
	) (*identity.UserProfileResponseDTO, error)

	// GetMe 获取当前用户信息 - 获取当前登录用户的详细信息
	GetMe(ctx context.Context, userID string) (*identity.UserProfileResponseDTO, error)

	// UpdateUser 更新用户信息 - 更新指定用户的基本信息
	UpdateUser(
		ctx context.Context,
		req *identity.UpdateUserRequestDTO,
		operatorID string,
	) (*identity.UserProfileResponseDTO, error)

	// UpdateMe 更新当前用户信息 - 用户更新自己的基本信息（从认证上下文获取用户ID）
	UpdateMe(
		ctx context.Context,
		req *identity.UpdateMeRequestDTO,
		userID string,
	) (*identity.UserProfileResponseDTO, error)

	// DeleteUser 删除用户 - 软删除指定用户（管理员权限）
	DeleteUser(
		ctx context.Context,
		req *identity.DeleteUserRequestDTO,
	) (*http_base.OperationStatusResponseDTO, error)

	// ListUsers 获取用户列表 - 分页查询用户列表，支持按组织、状态等条件筛选
	ListUsers(
		ctx context.Context,
		req *identity.ListUsersRequestDTO,
	) (*identity.ListUsersResponseDTO, error)

	// SearchUsers 搜索用户 - 按关键词搜索用户
	SearchUsers(
		ctx context.Context,
		req *identity.SearchUsersRequestDTO,
	) (*identity.SearchUsersResponseDTO, error)

	// ChangeUserStatus 变更用户状态 - 管理员变更用户状态（激活、停用、锁定等）
	ChangeUserStatus(
		ctx context.Context,
		req *identity.ChangeUserStatusRequestDTO,
	) (*http_base.OperationStatusResponseDTO, error)

	// UnlockUser 解锁用户 - 管理员解锁被锁定的用户
	UnlockUser(
		ctx context.Context,
		req *identity.UnlockUserRequestDTO,
	) (*http_base.OperationStatusResponseDTO, error)
}

// MembershipService 成员关系管理服务接口
type MembershipService interface {
	// GetUserMemberships 获取用户成员关系 - 获取指定用户的所有组织成员关系
	GetUserMemberships(
		ctx context.Context,
		req *identity.GetUserMembershipsRequestDTO,
	) (*identity.GetUserMembershipsResponseDTO, error)

	// GetPrimaryMembership 获取用户主成员关系 - 获取用户在当前组织中的主成员关系
	GetPrimaryMembership(
		ctx context.Context,
		req *identity.GetPrimaryMembershipRequestDTO,
	) (*identity.UserMembershipResponseDTO, error)
}

// OrganizationService 组织架构管理服务接口
type OrganizationService interface {
	// CreateOrganization 创建组织 - 创建新的组织机构
	CreateOrganization(
		ctx context.Context,
		req *identity.CreateOrganizationRequestDTO,
	) (*identity.OrganizationResponseDTO, error)

	// GetOrganization 获取组织信息 - 根据组织ID获取组织详细信息
	GetOrganization(
		ctx context.Context,
		req *identity.GetOrganizationRequestDTO,
	) (*identity.OrganizationResponseDTO, error)

	// UpdateOrganization 更新组织信息 - 更新指定组织的信息
	UpdateOrganization(
		ctx context.Context,
		req *identity.UpdateOrganizationRequestDTO,
	) (*identity.OrganizationResponseDTO, error)

	// DeleteOrganization 删除组织 - 软删除指定组织
	DeleteOrganization(
		ctx context.Context,
		req *identity.DeleteOrganizationRequestDTO,
	) (*http_base.OperationStatusResponseDTO, error)

	// ListOrganizations 获取组织列表 - 分页查询组织列表，支持按父组织筛选
	ListOrganizations(
		ctx context.Context,
		req *identity.ListOrganizationsRequestDTO,
	) (*identity.ListOrganizationsResponseDTO, error)
}

// DepartmentService 部门管理服务接口
type DepartmentService interface {
	// CreateDepartment 创建部门 - 在指定组织下创建新部门
	CreateDepartment(
		ctx context.Context,
		req *identity.CreateDepartmentRequestDTO,
	) (*identity.DepartmentResponseDTO, error)

	// GetDepartment 获取部门信息 - 根据部门ID获取部门详细信息
	GetDepartment(
		ctx context.Context,
		req *identity.GetDepartmentRequestDTO,
	) (*identity.DepartmentResponseDTO, error)

	// UpdateDepartment 更新部门信息 - 更新指定部门的信息
	UpdateDepartment(
		ctx context.Context,
		req *identity.UpdateDepartmentRequestDTO,
	) (*identity.DepartmentResponseDTO, error)

	// DeleteDepartment 删除部门 - 软删除指定部门
	DeleteDepartment(
		ctx context.Context,
		req *identity.DeleteDepartmentRequestDTO,
	) (*http_base.OperationStatusResponseDTO, error)

	// GetOrganizationDepartments 获取组织部门列表 - 获取指定组织下的所有部门
	GetOrganizationDepartments(
		ctx context.Context,
		req *identity.GetOrganizationDepartmentsRequestDTO,
	) (*identity.GetOrganizationDepartmentsResponseDTO, error)
}

// LogoService 组织Logo管理服务接口
type LogoService interface {
	// UploadTemporaryLogo 上传临时Logo - 创建临时Logo资源（7天后过期）
	UploadTemporaryLogo(
		ctx context.Context,
		req *identity.UploadTemporaryLogoRequestDTO,
		userID string,
	) (*identity.OrganizationLogoResponseDTO, error)

	// GetOrganizationLogo 获取Logo信息 - 根据LogoID获取Logo详细信息和预签名下载URL
	GetOrganizationLogo(
		ctx context.Context,
		req *identity.GetOrganizationLogoRequestDTO,
	) (*identity.OrganizationLogoResponseDTO, error)

	// DeleteOrganizationLogo 删除Logo - 逻辑删除Logo并删除S3文件
	DeleteOrganizationLogo(
		ctx context.Context,
		req *identity.DeleteOrganizationLogoRequestDTO,
	) (*http_base.OperationStatusResponseDTO, error)

	// BindLogoToOrganization 绑定Logo到组织 - 将临时Logo绑定到组织并转为永久保存
	BindLogoToOrganization(
		ctx context.Context,
		req *identity.BindLogoToOrganizationRequestDTO,
	) (*identity.OrganizationResponseDTO, error)
}
