namespace go identity

include "../base/base.thrift"
include "identity_model.thrift"

// =================================================================
//                        统一身份服务 (Identity Service)
// =================================================================

/**
 * IdentityService：统一身份管理服务
 *
 * 基于 RPC 层的 IdentityService 统一定义，提供核心的身份管理功能：
 * - 身份认证模块 (Authentication)
 * - 用户管理模块 (User Management)
 * - 成员关系管理模块 (Membership Management)
 * - 组织架构管理模块 (Organization Management)
 * - 部门管理模块 (Department Management)
 *
 * 注意：所有接口的权限控制在 API 网关层处理，遵循职责分离原则
 */
service IdentityService {
    // =================================================================
    // 1. 身份认证模块 (Authentication)
    // =================================================================

    /**
     * 用户登录
     * 验证用户凭据并返回访问令牌和用户信息
     */
    identity_model.LoginResponseDTO login(1: identity_model.LoginRequestDTO req) (api.post = "/api/v1/identity/auth/login"),

    /**
     * 用户登出
     * 注销当前会话并使令牌失效
     */
    base.OperationStatusResponseDTO logout(1: identity_model.LogoutRequestDTO req) (api.post = "/api/v1/identity/auth/logout"),

    /**
     * 修改密码
     * 用户修改自己的密码（需要提供旧密码）
     */
    base.OperationStatusResponseDTO changePassword(1: identity_model.ChangePasswordRequestDTO req) (api.put = "/api/v1/identity/auth/password"),

    /**
     * 重置密码
     * 管理员重置用户密码（管理员权限）
     */
    base.OperationStatusResponseDTO resetPassword(1: identity_model.ResetPasswordRequestDTO req) (api.post = "/api/v1/identity/auth/password/reset"),

    /**
     * 强制下次登录修改密码
     * 管理员标记用户需要在下次登录时强制修改密码
     */
    base.OperationStatusResponseDTO forcePasswordChange(1: identity_model.ForcePasswordChangeRequestDTO req) (api.put = "/api/v1/identity/auth/password/force-change"),

    /**
     * 刷新访问令牌
     * 使用刷新令牌获取新的访问令牌
     */
    identity_model.RefreshTokenResponseDTO refreshToken(1: identity_model.RefreshTokenRequestDTO req) (api.post = "/api/v1/identity/auth/refresh"),
    // =================================================================
    // 2. 用户管理模块 (User Management)
    // =================================================================

    /**
     * 创建用户
     * 管理员创建新用户账户
     */
    identity_model.UserProfileResponseDTO createUser(1: identity_model.CreateUserRequestDTO req) (api.post = "/api/v1/identity/users"),

    /**
     * 获取用户信息
     * 根据用户ID获取用户详细信息
     */
    identity_model.UserProfileResponseDTO getUser(1: identity_model.GetUserRequestDTO req) (api.get = "/api/v1/identity/users/:userID"),

    /**
     * 获取当前用户信息
     * 获取当前登录用户的详细信息
     */
    identity_model.UserProfileResponseDTO getMe() (api.get = "/api/v1/identity/users/me"),

    /**
     * 更新用户信息
     * 更新指定用户的基本信息
     */
    identity_model.UserProfileResponseDTO updateUser(1: identity_model.UpdateUserRequestDTO req) (api.put = "/api/v1/identity/users/:userID"),

    /**
     * 更新当前用户信息
     * 用户更新自己的基本信息（从认证上下文获取用户ID）
     */
    identity_model.UserProfileResponseDTO updateMe(1: identity_model.UpdateMeRequestDTO req) (api.put = "/api/v1/identity/users/me"),

    /**
     * 删除用户
     * 软删除指定用户（管理员权限）
     */
    base.OperationStatusResponseDTO deleteUser(1: identity_model.DeleteUserRequestDTO req) (api.delete = "/api/v1/identity/users/:userID"),

    /**
     * 获取用户列表
     * 分页查询用户列表，支持按组织、状态等条件筛选
     */
    identity_model.ListUsersResponseDTO listUsers(1: identity_model.ListUsersRequestDTO req) (api.get = "/api/v1/identity/users"),

    /**
     * 搜索用户
     * 按关键词搜索用户
     */
    identity_model.SearchUsersResponseDTO searchUsers(1: identity_model.SearchUsersRequestDTO req) (api.get = "/api/v1/identity/users/search"),

    /**
     * 变更用户状态
     * 管理员变更用户状态（激活、停用、锁定等）
     */
    base.OperationStatusResponseDTO changeUserStatus(1: identity_model.ChangeUserStatusRequestDTO req) (api.put = "/api/v1/identity/users/:userID/status"),

    /**
     * 解锁用户
     * 管理员解锁被锁定的用户
     */
    base.OperationStatusResponseDTO unlockUser(1: identity_model.UnlockUserRequestDTO req) (api.put = "/api/v1/identity/users/:userID/unlock"),
    // =================================================================
    // 3. 成员关系管理模块 (Membership Management)
    // =================================================================

    /**
     * 获取用户成员关系
     * 获取指定用户的所有组织成员关系
     */
    identity_model.GetUserMembershipsResponseDTO getUserMemberships(1: identity_model.GetUserMembershipsRequestDTO req) (api.get = "/api/v1/identity/users/:userID/memberships"),

    /**
     * 获取用户的主要成员关系。
     * @param userID 要查询的用户ID。
     * @return 用户的主要成员关系信息。
     */
    identity_model.UserMembershipResponseDTO getPrimaryMembership(1: identity_model.GetPrimaryMembershipRequestDTO req) (api.get = "/api/v1/identity/users/:userID/primary-membership"),

    /**
     * 检查用户是否属于某个组织或部门。
     * @param req 包含用户ID、组织ID等检查信息。
     * @return 如果用户是该组织的成员，则返回 true，否则返回 false。
     */
    base.OperationStatusResponseDTO checkMembership(1: identity_model.CheckMembershipRequestDTO req),
    // =================================================================
    // 4. 组织架构管理模块 (Organization Management)
    // =================================================================

    /**
     * 创建组织
     * 创建新的组织机构
     */
    identity_model.OrganizationResponseDTO createOrganization(1: identity_model.CreateOrganizationRequestDTO req) (api.post = "/api/v1/identity/organizations"),

    /**
     * 获取组织信息
     * 根据组织ID获取组织详细信息
     */
    identity_model.OrganizationResponseDTO getOrganization(1: identity_model.GetOrganizationRequestDTO req) (api.get = "/api/v1/identity/organizations/:organizationID"),

    /**
     * 更新组织信息
     * 更新指定组织的信息
     */
    identity_model.OrganizationResponseDTO updateOrganization(1: identity_model.UpdateOrganizationRequestDTO req) (api.put = "/api/v1/identity/organizations/:organizationID"),

    /**
     * 删除组织
     * 软删除指定组织
     */
    base.OperationStatusResponseDTO deleteOrganization(1: identity_model.DeleteOrganizationRequestDTO req) (api.delete = "/api/v1/identity/organizations/:organizationID"),

    /**
     * 获取组织列表
     * 分页查询组织列表，支持按父组织筛选
     */
    identity_model.ListOrganizationsResponseDTO listOrganizations(1: identity_model.ListOrganizationsRequestDTO req) (api.get = "/api/v1/identity/organizations"),
    // =================================================================
    // 5. 部门管理模块 (Department Management)
    // =================================================================

    /**
     * 创建部门
     * 在指定组织下创建新部门
     */
    identity_model.DepartmentResponseDTO createDepartment(1: identity_model.CreateDepartmentRequestDTO req) (api.post = "/api/v1/identity/departments"),

    /**
     * 获取部门信息
     * 根据部门ID获取部门详细信息
     */
    identity_model.DepartmentResponseDTO getDepartment(1: identity_model.GetDepartmentRequestDTO req) (api.get = "/api/v1/identity/departments/:departmentID"),

    /**
     * 更新部门信息
     * 更新指定部门的信息
     */
    identity_model.DepartmentResponseDTO updateDepartment(1: identity_model.UpdateDepartmentRequestDTO req) (api.put = "/api/v1/identity/departments/:departmentID"),

    /**
     * 删除部门
     * 软删除指定部门
     */
    base.OperationStatusResponseDTO deleteDepartment(1: identity_model.DeleteDepartmentRequestDTO req) (api.delete = "/api/v1/identity/departments/:departmentID"),

    /**
     * 获取组织部门列表
     * 获取指定组织下的所有部门
     */
    identity_model.GetOrganizationDepartmentsResponseDTO getOrganizationDepartments(1: identity_model.GetOrganizationDepartmentsRequestDTO req) (api.get = "/api/v1/identity/organizations/:organizationID/departments"),
    // =================================================================
    // 6. 组织Logo管理模块 (Organization Logo Management)
    // =================================================================

    /**
     * 上传临时Logo
     * 上传组织Logo文件到临时存储（7天过期）
     */
    identity_model.OrganizationLogoResponseDTO uploadTemporaryLogo(1: identity_model.UploadTemporaryLogoRequestDTO req) (api.post = "/api/v1/identity/organization-logos/temporary"),

    /**
     * 获取Logo详情
     * 根据Logo ID获取Logo元数据
     */
    identity_model.OrganizationLogoResponseDTO getOrganizationLogo(1: identity_model.GetOrganizationLogoRequestDTO req) (api.get = "/api/v1/identity/organization-logos/:logoID"),

    /**
     * 删除Logo
     * 删除Logo文件和数据库记录
     */
    base.OperationStatusResponseDTO deleteOrganizationLogo(1: identity_model.DeleteOrganizationLogoRequestDTO req) (api.delete = "/api/v1/identity/organization-logos/:logoID"),

    /**
     * 绑定Logo到组织
     * 将临时Logo绑定到组织（永久保存）
     */
    identity_model.OrganizationResponseDTO bindLogoToOrganization(1: identity_model.BindLogoToOrganizationRequestDTO req) (api.put = "/api/v1/identity/organizations/:organizationID/logo"),
}