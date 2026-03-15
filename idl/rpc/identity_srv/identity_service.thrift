/**
 * 统一身份认证服务 (Identity Service)
 *
 * 定义了与身份、认证、授权、组织架构等相关的所有 RPC 接口。
 */
namespace go identity_srv

include "../../base/core.thrift"
include "../../base/enums.thrift"
include "../base/base.thrift"
include "./identity_model.thrift"
include "./oauth2_model.thrift"

// =================================================================
// 服务接口定义 (Service Definition)
// =================================================================

/**
 * 身份服务 (IdentityService)
 *
 * 提供用户身份认证、用户管理、组织架构管理的统一服务接口。
 * 按功能模块组织方法，保持高内聚低耦合的设计原则。
 */
service IdentityService {
    // -----------------------------------------------------------------
    // 身份认证模块 (Authentication)
    // -----------------------------------------------------------------

    /**
     * 用户登录认证。
     * @param req 登录请求，包含用户名和密码。
     * @return 登录成功后返回用户信息和成员关系列表。
     */
    LoginResponse Login(1: LoginRequest req),

    /**
     * 修改当前用户密码。
     * @param req 包含用户ID、旧密码和新密码。
     */
    void ChangePassword(1: ChangePasswordRequest req),

    /**
     * 重置用户密码（通常由管理员执行）。
     * @param req 包含用户ID和可选的新密码。如果新密码为空，服务端将生成一个随机密码。
     */
    void ResetPassword(1: ResetPasswordRequest req),

    /**
     * 强制用户在下次登录时必须修改密码。
     * @param req 包含需要强制修改密码的用户ID。
     */
    void ForcePasswordChange(1: ForcePasswordChangeRequest req),
    // -----------------------------------------------------------------
    // 用户管理模块 (User Management)
    // -----------------------------------------------------------------

    /**
     * 创建一个新用户。
     * @param req 包含新用户的基本信息，如用户名、密码、邮箱等。
     * @return 创建成功后的用户画像 (UserProfile)。
     */
    identity_model.UserProfile CreateUser(1: CreateUserRequest req),

    /**
     * 根据ID获取用户详细信息。
     * @param req 包含要查询的用户ID。
     * @return 查找到的用户画像 (UserProfile)。
     */
    identity_model.UserProfile GetUser(1: GetUserRequest req),

    /**
     * 更新用户信息。
     * @param req 包含用户ID和需要更新的字段信息。
     * @return 更新后的用户画像 (UserProfile)。
     */
    identity_model.UserProfile UpdateUser(1: UpdateUserRequest req),

    /**
     * 删除用户（逻辑删除）。
     * @param req 包含要删除的用户ID。
     */
    void DeleteUser(1: DeleteUserRequest req),

    /**
     * 分页列出用户信息。
     * @param req 查询条件，可根据组织、状态等进行筛选。
     * @return 用户画像列表及分页信息。
     */
    ListUsersResponse ListUsers(1: ListUsersRequest req),

    /**
     * 搜索用户。
     * @param req 包含搜索关键字和分页信息。
     * @return 符合条件的用户列表。
     */
    SearchUsersResponse SearchUsers(1: SearchUsersRequest req),

    /**
     * 改变用户的状态（如激活、禁用等）。
     * @param req 包含用户ID和新的状态。
     */
    void ChangeUserStatus(1: ChangeUserStatusRequest req),

    /**
     * 解锁因多次尝试登录失败而被锁定的用户。
     * @param req 包含要解锁的用户ID。
     */
    void UnlockUser(1: UnlockUserRequest req),
    // -----------------------------------------------------------------
    // 组织与成员关系管理模块 (Organization & Membership)
    // -----------------------------------------------------------------

    /**
     * 创建一个新组织。
     * @param req 包含组织的名称、父组织ID等信息。
     * @return 创建成功的组织信息。
     */
    identity_model.Organization CreateOrganization(1: CreateOrganizationRequest req),

    /**
     * 根据ID获取组织详细信息。
     * @param req 包含要查询的组织ID。
     * @return 查找到的组织信息。
     */
    identity_model.Organization GetOrganization(1: GetOrganizationRequest req),

    /**
     * 更新组织信息。
     * @param req 包含组织ID和需要更新的字段。
     * @return 更新后的组织信息。
     */
    identity_model.Organization UpdateOrganization(1: UpdateOrganizationRequest req),

    /**
     * 删除组织（逻辑删除）。
     * @param organizationID 要删除的组织ID。
     */
    void DeleteOrganization(1: core.UUID organizationID),

    /**
     * 分页列出组织信息。
     * @param req 查询条件，可根据父组织ID进行筛选。
     * @return 组织列表及分页信息。
     */
    ListOrganizationsResponse ListOrganizations(1: ListOrganizationsRequest req),

    /**
     * 为用户添加新的组织成员关系。
     * @param req 包含用户ID、组织ID、角色等信息。
     * @return 创建成功的成员关系信息。
     */
    identity_model.UserMembership AddMembership(1: AddMembershipRequest req),

    /**
     * 更新用户的组织成员关系。
     * @param req 包含成员关系ID和需要更新的字段。
     * @return 更新后的成员关系信息。
     */
    identity_model.UserMembership UpdateMembership(1: UpdateMembershipRequest req),

    /**
     * 移除用户的组织成员关系（逻辑删除）。
     * @param membershipID 要移除的成员关系ID。
     */
    void RemoveMembership(1: core.UUID membershipID),

    /**
     * 根据ID获取成员关系详情。
     * @param membershipID 要查询的成员关系ID。
     * @return 查找到的成员关系信息。
     */
    identity_model.UserMembership GetMembership(1: core.UUID membershipID),

    /**
     * 获取一个用户的所有成员关系列表。
     * @param req 查询条件，可根据用户ID、组织ID等筛选。
     * @return 成员关系列表及分页信息。
     */
    GetUserMembershipsResponse GetUserMemberships(1: GetUserMembershipsRequest req),

    /**
     * 获取用户的主要成员关系。
     * @param userID 要查询的用户ID。
     * @return 用户的主要成员关系信息。
     */
    identity_model.UserMembership GetPrimaryMembership(1: core.UUID userID),

    /**
     * 检查用户是否属于某个组织或部门。
     * @param req 包含用户ID、组织ID等检查信息。
     * @return 如果用户是该组织的成员，则返回 true，否则返回 false。
     */
    bool CheckMembership(1: CheckMembershipRequest req),
    // -----------------------------------------------------------------
    // 部门管理模块 (Department Management)
    // -----------------------------------------------------------------

    /**
     * 在指定组织下创建新部门。
     * @param req 包含组织ID和部门名称等信息。
     * @return 创建成功的部门信息。
     */
    identity_model.Department CreateDepartment(1: CreateDepartmentRequest req),

    /**
     * 根据ID获取部门详细信息。
     * @param req 包含要查询的部门ID。
     * @return 查找到的部门信息。
     */
    identity_model.Department GetDepartment(1: GetDepartmentRequest req),

    /**
     * 更新部门信息。
     * @param req 包含部门ID和需要更新的字段。
     * @return 更新后的部门信息。
     */
    identity_model.Department UpdateDepartment(1: UpdateDepartmentRequest req),

    /**
     * 删除部门（逻辑删除）。
     * @param departmentID 要删除的部门ID。
     */
    void DeleteDepartment(1: core.UUID departmentID),

    /**
     * 获取指定组织下的所有部门列表。
     * @param req 包含组织ID和分页信息。
     * @return 部门列表及分页信息。
     */
    GetOrganizationDepartmentsResponse GetOrganizationDepartments(1: GetOrganizationDepartmentsRequest req),
    // -----------------------------------------------------------------
    // 组织Logo管理模块 (Organization Logo Management)
    // -----------------------------------------------------------------

    /**
     * 上传临时组织Logo。
     * @param req 包含文件内容、文件名等信息。
     * @return 创建成功的Logo信息（临时状态，7天后过期）。
     */
    identity_model.OrganizationLogo UploadTemporaryLogo(1: UploadTemporaryLogoRequest req),

    /**
     * 根据ID获取Logo详细信息。
     * @param req 包含要查询的LogoID。
     * @return 查找到的Logo信息。
     */
    identity_model.OrganizationLogo GetOrganizationLogo(1: GetOrganizationLogoRequest req),

    /**
     * 删除Logo（逻辑删除并删除S3文件）。
     * @param req 包含要删除的LogoID。
     */
    void DeleteOrganizationLogo(1: DeleteOrganizationLogoRequest req),

    /**
     * 绑定Logo到组织（内部方法，将临时Logo转为永久保存）。
     * @param req 包含LogoID和组织ID。
     * @return 更新后的Logo信息。
     */
    identity_model.OrganizationLogo BindLogoToOrganization(1: BindLogoToOrganizationRequest req),
    // -----------------------------------------------------------------
    // 角色与权限管理模块 (Role & Permission Management)
    // -----------------------------------------------------------------

    /**
     * 创建一个新的角色定义。
     * @param req 包含角色名称、权限列表等信息。
     * @return 创建成功的角色定义信息。
     */
    identity_model.RoleDefinition CreateRoleDefinition(1: RoleDefinitionCreateRequest req),

    /**
     * 更新一个已有的角色定义。
     * @param req 包含角色ID和需要更新的字段。
     * @return 更新后的角色定义信息。
     */
    identity_model.RoleDefinition UpdateRoleDefinition(1: RoleDefinitionUpdateRequest req),

    /**
     * 删除一个角色定义。
     * @param roleID 要删除的角色ID。
     */
    void DeleteRoleDefinition(1: core.UUID roleID),

    /**
     * 根据ID获取角色定义。
     * @param roleID 要查询的角色ID。
     * @return 查找到的角色定义信息。
     */
    identity_model.RoleDefinition GetRoleDefinition(1: core.UUID roleID),

    /**
     * 分页列出角色定义。
     * @param req 查询条件，可根据分类、状态等筛选。
     * @return 角色定义列表及分页信息。
     */
    RoleDefinitionListResponse ListRoleDefinitions(1: RoleDefinitionQueryRequest req),

    /**
     * 为用户分配一个角色。
     * @param req 包含用户ID、组织ID和角色名称等信息。
     * @return 分配成功后返回的分配信息。
     */
    UserRoleAssignmentResponse AssignRoleToUser(1: AssignRoleToUserRequest req),

    /**
     * 更新用户的角色分配信息。
     * @param req 包含分配ID和需要更新的字段。
     */
    void UpdateUserRoleAssignment(1: UpdateUserRoleAssignmentRequest req),

    /**
     * 撤销用户的角色分配。
     * @param req 包含用户ID和角色ID。
     */
    void RevokeRoleFromUser(1: RevokeRoleFromUserRequest req),

    /**
     * 获取用户最后一次的角色分配信息。
     * @param userID 要查询的用户ID。
     */
    identity_model.UserRoleAssignment GetLastUserRoleAssignment(1: core.UUID userID),

    /**
     * 列出用户的角色分配记录。
     * @param req 查询条件，可根据用户、角色等筛选。
     * @return 角色分配记录列表及分页信息。
     */
    UserRoleListResponse ListUserRoleAssignments(1: UserRoleQueryRequest req),

    /**
     * 根据角色ID获取该角色下所有用户。
     * @param req 包含角色ID的请求。
     * @return 该角色下所有用户的ID列表（不分页）。
     */
    GetUsersByRoleResponse GetUsersByRole(1: GetUsersByRoleRequest req),

    /**
     * 批量绑定用户到角色。
     * @param req 包含角色ID、用户ID列表和操作者ID。
     * @return 批量绑定操作结果。
     */
    BatchBindUsersToRoleResponse BatchBindUsersToRole(1: BatchBindUsersToRoleRequest req),

    /**
     * 批量获取多个用户的角色分配。
     * @param req 包含用户ID列表的请求。
     * @return 包含每个用户的角色ID列表的响应。
     */
    BatchGetUserRolesResponse BatchGetUserRoles(1: BatchGetUserRolesRequest req),
    // -----------------------------------------------------------------
    // 菜单管理模块 (Menu Management)
    // -----------------------------------------------------------------

    /**
     * 上传并解析菜单配置文件 (menu.yaml)。
     * @param req 包含 YAML 文件内容的请求。
     */
    void UploadMenu(1: UploadMenuRequest req),

    /**
     * 获取指定用户的菜单树。
     * @return 用户可见的菜单树结构。
     */
    GetMenuTreeResponse GetMenuTree(),
    // -----------------------------------------------------------------
    // 菜单权限管理模块 (Menu Permission Management)
    // -----------------------------------------------------------------

    /**
     * 配置角色的菜单权限。
     * @param req 包含角色ID和菜单权限配置信息。
     * @return 配置成功响应。
     */
    ConfigureRoleMenusResponse ConfigureRoleMenus(1: ConfigureRoleMenusRequest req),

    /**
     * 获取角色的菜单树。
     * @param req 包含角色ID的请求。
     * @return 角色可访问的菜单树结构。
     */
    GetRoleMenuTreeResponse GetRoleMenuTree(1: GetRoleMenuTreeRequest req),

    /**
     * 获取用户的菜单树（基于角色合并）。
     * @param req 包含用户ID的请求。
     * @return 用户可访问的菜单树结构。
     */
    GetUserMenuTreeResponse GetUserMenuTree(1: GetUserMenuTreeRequest req),

    /**
     * 获取角色的菜单权限列表。
     * @param req 包含角色ID的请求。
     * @return 角色的菜单权限配置列表。
     */
    GetRoleMenuPermissionsResponse GetRoleMenuPermissions(1: GetRoleMenuPermissionsRequest req),

    /**
     * 检查角色是否具有指定菜单权限。
     * @param req 包含角色ID、菜单ID和权限类型的请求。
     * @return 权限检查结果。
     */
    HasMenuPermissionResponse HasMenuPermission(1: HasMenuPermissionRequest req),

    /**
     * 获取用户的菜单权限列表（基于所有活跃角色合并）。
     * @param req 包含用户ID的请求。
     * @return 用户的合并菜单权限列表（去重，取最高权限）。
     */
    GetUserMenuPermissionsResponse GetUserMenuPermissions(1: GetUserMenuPermissionsRequest req),
    // -----------------------------------------------------------------
    // Casbin 权限管理模块 (Casbin Permission Management)
    // -----------------------------------------------------------------

    /**
     * 检查用户权限（基于 Casbin RBAC）。
     * @param req 包含用户ID、资源、操作等信息。
     * @return 权限检查结果和数据范围。
     */
    CheckPermissionResponse CheckPermission(1: CheckPermissionRequest req),

    /**
     * 同步所有 Casbin 策略（从数据库生成 Casbin 规则）。
     * @return 同步结果。
     */
    SyncPoliciesResponse SyncPolicies(),

    /**
     * 获取用户的数据范围。
     * @param req 包含用户ID、资源、操作等信息。
     * @return 用户在指定资源上的数据范围。
     */
    GetUserDataScopeResponse GetUserDataScope(1: GetUserDataScopeRequest req),
    // -----------------------------------------------------------------
    // 审计日志模块 (Audit Log)
    // -----------------------------------------------------------------

    /**
     * 创建审计日志记录。
     * @param req 包含审计日志的详细信息。
     */
    void CreateAuditLog(1: CreateAuditLogRequest req),

    /**
     * 分页查询审计日志。
     * @param req 查询条件，支持按用户、操作类型、资源、时间范围等筛选。
     * @return 审计日志列表及分页信息。
     */
    ListAuditLogsResponse ListAuditLogs(1: ListAuditLogsRequest req),
    // -----------------------------------------------------------------
    // OAuth2 客户端管理模块 (OAuth2 Client Management)
    // -----------------------------------------------------------------

    /**
     * 创建 OAuth2 客户端（应用）。
     * @param req 包含客户端名称、类型、回调地址等信息。
     * @return 创建成功的客户端信息（含生成的 client_id 和 client_secret）。
     */
    CreateOAuth2ClientResponse CreateOAuth2Client(1: CreateOAuth2ClientRequest req),

    /**
     * 根据 client_id 获取 OAuth2 客户端信息。
     * @param req 包含 client_id。
     * @return 客户端详细信息。
     */
    oauth2_model.OAuth2Client GetOAuth2Client(1: GetOAuth2ClientRequest req),

    /**
     * 更新 OAuth2 客户端信息。
     * @param req 包含客户端ID和需要更新的字段。
     * @return 更新后的客户端信息。
     */
    oauth2_model.OAuth2Client UpdateOAuth2Client(1: UpdateOAuth2ClientRequest req),

    /**
     * 删除 OAuth2 客户端。
     * @param req 包含要删除的客户端ID。
     */
    void DeleteOAuth2Client(1: DeleteOAuth2ClientRequest req),

    /**
     * 分页列出 OAuth2 客户端。
     * @param req 查询条件。
     * @return 客户端列表及分页信息。
     */
    ListOAuth2ClientsResponse ListOAuth2Clients(1: ListOAuth2ClientsRequest req),

    /**
     * 轮换 OAuth2 客户端密钥。
     * @param req 包含客户端ID。
     * @return 新的 client_secret。
     */
    RotateOAuth2ClientSecretResponse RotateOAuth2ClientSecret(1: RotateOAuth2ClientSecretRequest req),

    /**
     * 根据 client_id 获取 OAuth2 客户端（供 fosite 内部存储层调用）。
     * 与 GetOAuth2Client 的区别：返回完整信息含密钥哈希，用于 fosite 验证。
     * @param clientID 客户端标识符。
     * @return 客户端完整信息。
     */
    GetOAuth2ClientForAuthResponse GetOAuth2ClientForAuth(1: string clientID),
    // -----------------------------------------------------------------
    // OAuth2 作用域管理模块 (OAuth2 Scope Management)
    // -----------------------------------------------------------------

    /**
     * 列出所有可用的 OAuth2 作用域。
     * @param req 查询条件。
     * @return 作用域列表。
     */
    ListOAuth2ScopesResponse ListOAuth2Scopes(1: ListOAuth2ScopesRequest req),
    // -----------------------------------------------------------------
    // OAuth2 用户授权同意管理模块 (OAuth2 Consent Management)
    // -----------------------------------------------------------------

    /**
     * 保存用户的授权同意记录。
     * @param req 包含用户ID、客户端ID、授权的 scopes。
     */
    void SaveOAuth2Consent(1: SaveOAuth2ConsentRequest req),

    /**
     * 获取用户对指定客户端的授权同意记录。
     * @param req 包含用户ID和客户端ID。
     * @return 授权同意记录（如存在）。
     */
    GetOAuth2ConsentResponse GetOAuth2Consent(1: GetOAuth2ConsentRequest req),

    /**
     * 列出用户的所有授权同意记录。
     * @param req 包含用户ID。
     * @return 授权同意记录列表。
     */
    ListOAuth2ConsentsResponse ListOAuth2Consents(1: ListOAuth2ConsentsRequest req),

    /**
     * 撤销用户对指定客户端的授权同意。
     * @param req 包含用户ID和客户端ID。
     */
    void RevokeOAuth2Consent(1: RevokeOAuth2ConsentRequest req),
    // -----------------------------------------------------------------
    // OAuth2 令牌存储模块 (OAuth2 Token Storage)
    // 供 fosite 存储层通过 RPC 调用，实现分布式令牌管理。
    // -----------------------------------------------------------------

    /**
     * 创建授权码会话。
     */
    void CreateOAuth2AuthorizeCodeSession(1: OAuth2TokenSession req),

    /**
     * 获取授权码会话。
     */
    OAuth2TokenSession GetOAuth2AuthorizeCodeSession(1: string signature),

    /**
     * 使授权码会话失效（标记已使用）。
     */
    void InvalidateOAuth2AuthorizeCodeSession(1: string signature),

    /**
     * 创建访问令牌会话。
     */
    void CreateOAuth2AccessTokenSession(1: OAuth2TokenSession req),

    /**
     * 获取访问令牌会话。
     */
    OAuth2TokenSession GetOAuth2AccessTokenSession(1: string signature),

    /**
     * 删除访问令牌会话。
     */
    void DeleteOAuth2AccessTokenSession(1: string signature),

    /**
     * 按请求ID吊销访问令牌。
     */
    void RevokeOAuth2AccessToken(1: string requestID),

    /**
     * 创建刷新令牌会话。
     */
    void CreateOAuth2RefreshTokenSession(1: OAuth2TokenSession req),

    /**
     * 获取刷新令牌会话。
     */
    OAuth2TokenSession GetOAuth2RefreshTokenSession(1: string signature),

    /**
     * 删除刷新令牌会话。
     */
    void DeleteOAuth2RefreshTokenSession(1: string signature),

    /**
     * 按请求ID吊销刷新令牌。
     */
    void RevokeOAuth2RefreshToken(1: string requestID),

    /**
     * 创建 PKCE 会话。
     */
    void CreateOAuth2PKCESession(1: OAuth2TokenSession req),

    /**
     * 获取 PKCE 会话。
     */
    OAuth2TokenSession GetOAuth2PKCESession(1: string signature),

    /**
     * 删除 PKCE 会话。
     */
    void DeleteOAuth2PKCESession(1: string signature),
}

// =================================================================
// 认证相关 (Authentication)
// =================================================================

/** 用户登录请求 */
struct LoginRequest {

    /** 用户名 */
    1: optional string username,

    /** 密码 (应在传输过程中加密) */
    2: optional string password,
}

/** 用户登录响应 */
struct LoginResponse {

    /** 登录用户的个人信息 */
    1: optional identity_model.UserProfile userProfile,

    /** 用户所属的组织成员关系列表 */
    2: optional list<identity_model.UserMembership> memberships,

    /** 用户可见的菜单树 */
    3: optional list<identity_model.MenuNode> menuTree,

    /** 用户拥有的角色ID列表 */
    4: optional list<string> roleIDs,

    /** 用户拥有的菜单权限列表（菜单ID -> 权限） */
    5: optional list<MenuPermission> permissions,

    /** 用户角色详情列表（包含角色编码、名称和数据范围） */
    6: optional list<identity_model.RoleDefinition> roleDetails,
}

/** 修改密码请求 */
struct ChangePasswordRequest {

    /** 用户ID */
    1: optional core.UUID userID,

    /** 旧密码 */
    2: optional string oldPassword,

    /** 新密码 */
    3: optional string newPassword,
}

/** 重置密码请求 */
struct ResetPasswordRequest {

    /** 用户ID */
    1: optional core.UUID userID,

    /** 新密码。如果为空，则由服务端生成随机密码 */
    2: optional string newPassword,
}

/** 强制密码修改请求 */
struct ForcePasswordChangeRequest {

    /** 用户ID */
    1: optional core.UUID userID,
}

// =================================================================
// 用户管理相关 (User Management)
// =================================================================

/** 创建用户请求 */
struct CreateUserRequest {
    1: optional string username,
    2: optional string password,
    3: optional string email,
    4: optional string phone,
    5: optional string firstName,
    6: optional string lastName,
    7: optional string realName,
    8: optional enums.Gender gender,
    9: optional string professionalTitle,
    10: optional string employeeID,
    11: optional bool mustChangePassword,
    12: optional core.TimestampMS accountExpiry,
}

/** 获取用户请求 */
struct GetUserRequest {
    1: optional core.UUID userID,
}

/** 更新用户请求 */
struct UpdateUserRequest {
    1: optional core.UUID userID,
    2: optional string email,
    3: optional string phone,
    4: optional string firstName,
    5: optional string lastName,
    6: optional string realName,
    7: optional enums.Gender gender,
    8: optional string professionalTitle,
    9: optional string employeeID,
    10: optional core.TimestampMS accountExpiry,
}

/** 删除用户请求 */
struct DeleteUserRequest {
    1: optional core.UUID userID,
}

/** 列出用户请求 */
struct ListUsersRequest {
    1: optional base.PageRequest page,
    2: optional core.UUID organizationID,
    3: optional enums.UserStatus status,
}

/** 列出用户响应 */
struct ListUsersResponse {
    1: optional list<identity_model.UserProfile> users,
    2: optional base.PageResponse page,
}

/** 搜索用户请求 */
struct SearchUsersRequest {
    1: optional base.PageRequest page,
    2: optional core.UUID organizationID,
// 可以根据需要添加更多搜索字段，如 name, email 等
}

/** 搜索用户响应 */
struct SearchUsersResponse {
    1: optional list<identity_model.UserProfile> users,
    2: optional base.PageResponse page,
}

/** 更改用户状态请求 */
struct ChangeUserStatusRequest {
    1: optional core.UUID userID,
    2: optional enums.UserStatus newStatus,
}

/** 解锁用户请求 */
struct UnlockUserRequest {
    1: optional core.UUID userID,
}

// =================================================================
// 组织与成员关系 (Organization & Membership)
// =================================================================

/** 创建组织请求 */
struct CreateOrganizationRequest {
    1: optional string name,
    2: optional core.UUID parentID,
    3: optional string facilityType,
    4: optional string accreditationStatus,

    /** 组织所在省市 */
    5: optional list<string> provinceCity,

    /** 组织Logo ID（临时Logo的ID，创建组织后自动绑定） */
    6: optional core.UUID logoID,
}

/** 获取组织请求 */
struct GetOrganizationRequest {
    1: optional core.UUID organizationID,
}

/** 更新组织请求 */
struct UpdateOrganizationRequest {
    1: optional core.UUID organizationID,
    2: optional string name,
    3: optional core.UUID parentID,
    4: optional string facilityType,
    5: optional string accreditationStatus,
    6: optional list<string> provinceCity,

    /** 组织Logo ID（新的Logo ID，更新时会删除旧Logo并绑定新Logo） */
    7: optional core.UUID logoID,
}

/** 列出组织请求 */
struct ListOrganizationsRequest {
    1: optional core.UUID parentID,
    2: optional base.PageRequest page,
}

/** 列出组织响应 */
struct ListOrganizationsResponse {
    1: optional list<identity_model.Organization> organizations,
    2: optional base.PageResponse page,
}

/** 添加成员关系请求 */
struct AddMembershipRequest {
    1: optional core.UUID userID,
    2: optional core.UUID organizationID,
    3: optional core.UUID departmentID,

    /** 是否为主要成员关系 */
    4: optional bool isPrimary = false,
}

/** 更新成员关系请求 */
struct UpdateMembershipRequest {
    1: optional core.UUID membershipID,
    2: optional core.UUID organizationID,
    3: optional core.UUID departmentID,
    4: optional bool isPrimary,
}

/** 获取用户成员关系请求 */
struct GetUserMembershipsRequest {
    1: optional core.UUID userID,
    2: optional core.UUID organizationID,
    3: optional core.UUID departmentID,
    4: optional base.PageRequest page,
}

/** 获取用户成员关系响应 */
struct GetUserMembershipsResponse {
    1: optional list<identity_model.UserMembership> memberships,
    2: optional base.PageResponse page,
}

/** 检查成员关系请求 */
struct CheckMembershipRequest {
    1: optional core.UUID userID,
    2: optional core.UUID organizationID,
    3: optional core.UUID departmentID,
}

// =================================================================
// 部门管理 (Department)
// =================================================================

/** 创建部门请求 */
struct CreateDepartmentRequest {
    1: optional core.UUID organizationID,
    2: optional string name,
    3: optional string departmentType,
}

/** 获取部门请求 */
struct GetDepartmentRequest {
    1: optional core.UUID departmentID,
}

/** 更新部门请求 */
struct UpdateDepartmentRequest {
    1: optional core.UUID departmentID,
    2: optional string name,
    3: optional string departmentType,
}

/** 获取组织下所有部门请求 */
struct GetOrganizationDepartmentsRequest {
    1: optional core.UUID organizationID,
    2: optional base.PageRequest page,
}

/** 获取组织下所有部门响应 */
struct GetOrganizationDepartmentsResponse {
    1: optional list<identity_model.Department> departments,
    2: optional base.PageResponse page,
}

// =================================================================
// 组织Logo管理 (Organization Logo Management)
// =================================================================

/** 上传临时Logo请求 */
struct UploadTemporaryLogoRequest {

    /** 文件内容 (二进制) */
    1: optional binary fileContent,

    /** 文件名 */
    2: optional string fileName,

    /** MIME 类型 (如 image/png, image/jpeg) */
    3: optional string mimeType,

    /** 上传者用户ID */
    4: optional core.UUID uploadedBy,
}

/** 获取Logo请求 */
struct GetOrganizationLogoRequest {

    /** LogoID */
    1: optional core.UUID logoID,
}

/** 删除Logo请求 */
struct DeleteOrganizationLogoRequest {

    /** LogoID */
    1: optional core.UUID logoID,
}

/** 绑定Logo到组织请求 */
struct BindLogoToOrganizationRequest {

    /** LogoID */
    1: optional core.UUID logoID,

    /** 组织ID */
    2: optional core.UUID organizationID,
}

// =================================================================
// 角色与权限 (Role & Permission)
// =================================================================

/** 角色定义创建请求 */
struct RoleDefinitionCreateRequest {

    /** 角色唯一名称 (英文，用于程序识别) */
    1: optional string name,

    /** 角色描述 */
    2: optional string description,

    /** 角色包含的权限列表 */
    3: optional list<identity_model.Permission> permissions,

    /** 是否为系统内置角色 */
    4: optional bool isSystemRole = false,
}

/** 角色定义更新请求 */
struct RoleDefinitionUpdateRequest {
    1: optional core.UUID roleDefinitionID,
    2: optional string description,
    3: optional enums.RoleStatus status,
    4: optional list<identity_model.Permission> permissions,

    /** 角色名称 */
    5: optional string name,
}

/** 角色定义查询请求 */
struct RoleDefinitionQueryRequest {
    1: optional string name,
    2: optional enums.RoleStatus status,
    3: optional bool isSystemRole,
    4: optional base.PageRequest page,
}

/** 角色定义列表响应 */
struct RoleDefinitionListResponse {
    1: optional list<identity_model.RoleDefinition> roles,
    2: optional base.PageResponse page,
}

/** 用户角色分配请求 */
struct AssignRoleToUserRequest {

    /** 用户ID */
    1: optional core.UUID userID,

    /** 角色ID，对应 RoleDefinition.id */
    2: optional core.UUID roleID,

    /** 分配者用户ID */
    3: optional core.UUID assignedBy,
}

/** 更新用户角色分配请求 */
struct UpdateUserRoleAssignmentRequest {

    /** 分配ID */
    1: optional core.UUID assignmentID,

    /** 用户ID */
    2: optional core.UUID userID,

    /** 角色ID */
    3: optional core.UUID roleID,

    /** 更新者用户ID */
    4: optional core.UUID updatedBy,
}

/** 撤销用户角色分配请求 */
struct RevokeRoleFromUserRequest {

    /** 用户ID */
    1: optional core.UUID userID,

    /** 角色ID */
    2: optional core.UUID roleID,

    /** 操作者用户ID */
    3: optional core.UUID revokedBy,
}

/** 用户角色分配响应 */
struct UserRoleAssignmentResponse {

    /** 创建的分配记录ID */
    1: optional core.UUID assignmentID,
}

/** 用户角色查询请求 */
struct UserRoleQueryRequest {
    1: optional core.UUID userID,
    2: optional core.UUID roleID,
    3: optional base.PageRequest page,
}

/** 用户角色列表响应 */
struct UserRoleListResponse {
    1: optional list<identity_model.UserRoleAssignment> assignments,
    2: optional base.PageResponse page,
}

/** 根据角色ID获取用户请求 */
struct GetUsersByRoleRequest {

    /** 角色ID */
    1: optional core.UUID roleID,
}

/** 根据角色ID获取用户响应 */
struct GetUsersByRoleResponse {

    /** 角色ID */
    1: optional core.UUID roleID,

    /** 该角色下所有用户的ID列表 */
    2: optional list<core.UUID> userIDs,
}

/** 批量绑定用户到角色请求 */
struct BatchBindUsersToRoleRequest {

    /** 角色ID */
    1: optional core.UUID roleID,

    /** 用户ID列表 */
    2: optional list<core.UUID> userIDs,

    /** 操作者用户ID */
    3: optional core.UUID operatorID,
}

/** 批量绑定用户到角色响应 */
struct BatchBindUsersToRoleResponse {

    /** 操作是否成功 */
    1: optional bool success,

    /** 成功绑定的用户数量 */
    2: optional i32 successCount,

    /** 响应消息 */
    3: optional string message,
}

/** 单个用户的角色信息 */
struct UserRoles {

    /** 用户ID */
    1: optional core.UUID userID,

    /** 角色ID列表 */
    2: optional list<core.UUID> roleIDs,
}

/** 批量获取用户角色请求 */
struct BatchGetUserRolesRequest {

    /** 用户ID列表 */
    1: required list<core.UUID> userIDs,
}

/** 批量获取用户角色响应 */
struct BatchGetUserRolesResponse {

    /** 用户角色列表 */
    1: optional list<UserRoles> userRoles,
}

// =================================================================
// 菜单管理 (Menu Management)
// =================================================================

/** 菜单上传请求 */
struct UploadMenuRequest {

    /** 产品线标识 */
    1: optional string productLine,

    /** YAML 格式的菜单配置内容 */
    2: optional string yamlContent,
}

/** 菜单树获取响应 */
struct GetMenuTreeResponse {

    /** 完整的菜单树结构 */
    1: optional list<identity_model.MenuNode> menuTree,
}

// =================================================================
// 菜单权限管理 (Menu Permission Management)
// =================================================================

/** 菜单配置项 */
struct MenuConfig {

    /** 菜单ID */
    1: optional string menuID,

    /** 权限级别 */
    2: optional identity_model.PermissionLevel permission,
}

/** 菜单权限信息 */
struct MenuPermission {

    /** 菜单ID */
    1: optional string menuID,

    /** 权限级别 */
    2: optional identity_model.PermissionLevel permission,
}

/** 配置角色菜单权限请求 */
struct ConfigureRoleMenusRequest {

    /** 角色ID */
    1: optional core.UUID roleID,

    /** 菜单权限配置列表 */
    2: optional list<MenuConfig> menuConfigs,

    /** 操作者用户ID */
    3: optional core.UUID operatorID,
}

/** 配置角色菜单权限响应 */
struct ConfigureRoleMenusResponse {

    /** 配置成功标志 */
    1: optional bool success,

    /** 响应消息 */
    2: optional string message,
}

/** 获取角色菜单树请求 */
struct GetRoleMenuTreeRequest {

    /** 角色ID */
    1: optional core.UUID roleID,
}

/** 获取角色菜单树响应 */
struct GetRoleMenuTreeResponse {

    /** 角色可访问的菜单树 */
    1: optional list<identity_model.MenuNode> menuTree,

    /** 角色ID */
    2: optional core.UUID roleID,
}

/** 获取用户菜单树请求 */
struct GetUserMenuTreeRequest {

    /** 用户ID */
    1: optional core.UUID userID,
}

/** 获取用户菜单树响应 */
struct GetUserMenuTreeResponse {

    /** 用户可访问的菜单树 */
    1: optional list<identity_model.MenuNode> menuTree,

    /** 用户ID */
    2: optional core.UUID userID,

    /** 用户拥有的角色列表 */
    3: optional list<core.UUID> roleIDs,
}

/** 获取角色菜单权限请求 */
struct GetRoleMenuPermissionsRequest {

    /** 角色ID */
    1: optional core.UUID roleID,
}

/** 获取角色菜单权限响应 */
struct GetRoleMenuPermissionsResponse {

    /** 角色的菜单权限列表 */
    1: optional list<MenuPermission> permissions,

    /** 角色ID */
    2: optional core.UUID roleID,
}

/** 检查菜单权限请求 */
struct HasMenuPermissionRequest {

    /** 角色ID */
    1: optional core.UUID roleID,

    /** 菜单ID */
    2: optional string menuID,

    /** 权限级别 */
    3: optional identity_model.PermissionLevel permission,
}

/** 检查菜单权限响应 */
struct HasMenuPermissionResponse {

    /** 是否具有权限 */
    1: optional bool hasPermission,

    /** 角色ID */
    2: optional core.UUID roleID,

    /** 菜单ID */
    3: optional string menuID,

    /** 权限级别 */
    4: optional identity_model.PermissionLevel permission,
}

/** 获取用户菜单权限请求 */
struct GetUserMenuPermissionsRequest {

    /** 用户ID */
    1: optional core.UUID userID,
}

/** 获取用户菜单权限响应 */
struct GetUserMenuPermissionsResponse {

    /** 用户的菜单权限列表 */
    1: optional list<MenuPermission> permissions,

    /** 用户ID */
    2: optional core.UUID userID,

    /** 用户拥有的角色列表 */
    3: optional list<core.UUID> roleIDs,
}

// =================================================================
// Casbin 权限管理 (Casbin Permission Management)
// =================================================================

/** 检查权限请求 */
struct CheckPermissionRequest {

    /** 用户ID */
    1: optional core.UUID userID,

    /** 角色ID列表（多角色模式） */
    2: optional list<string> roleIDs,

    /** 部门ID列表（多部门模式） */
    3: optional list<string> departmentIDs,

    /** 资源标识（菜单SemanticID 或 API路径） */
    4: optional string resource,

    /** 操作类型（read/write/manage） */
    5: optional string action,
}

/** 检查权限响应 */
struct CheckPermissionResponse {

    /** 是否具有权限 */
    1: optional bool allowed,

    /** 数据范围（self/dept/org） */
    2: optional string dataScope,

    /** 用户ID */
    3: optional core.UUID userID,

    /** 资源标识 */
    4: optional string resource,

    /** 操作类型 */
    5: optional string action,
}

/** 同步策略响应 */
struct SyncPoliciesResponse {

    /** 同步是否成功 */
    1: optional bool success,

    /** 同步的角色策略数量 */
    2: optional i32 rolePolicyCount,

    /** 同步的用户角色绑定数量 */
    3: optional i32 userRoleBindingCount,

    /** 同步的角色继承数量 */
    4: optional i32 roleInheritanceCount,

    /** 响应消息 */
    5: optional string message,
}

/** 获取用户数据范围请求 */
struct GetUserDataScopeRequest {

    /** 用户ID */
    1: optional core.UUID userID,

    /** 资源标识 */
    2: optional string resource,

    /** 操作类型 */
    3: optional string action,
}

/** 获取用户数据范围响应 */
struct GetUserDataScopeResponse {

    /** 数据范围（self/dept/org） */
    1: optional string dataScope,

    /** 用户ID */
    2: optional core.UUID userID,

    /** 资源标识 */
    3: optional string resource,

    /** 操作类型 */
    4: optional string action,
}

// =================================================================
// 审计日志 (Audit Log)
// =================================================================

/** 创建审计日志请求 */
struct CreateAuditLogRequest {

    /** 请求ID */
    1: optional string requestID,

    /** 链路追踪ID */
    2: optional string traceID,

    /** 操作者用户ID */
    3: optional core.UUID userID,

    /** 操作者用户名 */
    4: optional string username,

    /** 操作者所属组织ID */
    5: optional core.UUID organizationID,

    /** 审计操作类型 */
    6: optional identity_model.AuditAction action,

    /** 操作的资源路径 */
    7: optional string resource,

    /** 操作的资源ID */
    8: optional string resourceID,

    /** HTTP 状态码 */
    9: optional i32 statusCode,

    /** 操作是否成功 */
    10: optional bool success,

    /** 客户端IP地址 */
    11: optional string clientIP,

    /** 客户端User-Agent */
    12: optional string userAgent,

    /** 请求体（已脱敏） */
    13: optional string requestBody,

    /** 请求耗时（毫秒） */
    14: optional i32 durationMs,
}

/** 查询审计日志请求 */
struct ListAuditLogsRequest {

    /** 分页信息 */
    1: optional base.PageRequest page,

    /** 按用户ID筛选 */
    2: optional core.UUID userID,

    /** 按操作类型筛选 */
    3: optional identity_model.AuditAction action,

    /** 按资源路径筛选 */
    4: optional string resource,

    /** 按成功/失败筛选 */
    5: optional bool success,

    /** 开始时间（毫秒时间戳） */
    6: optional core.TimestampMS startTime,

    /** 结束时间（毫秒时间戳） */
    7: optional core.TimestampMS endTime,
}

/** 审计日志统计信息 */
struct AuditLogStats {

    /** 符合筛选条件的总记录数 */
    1: optional i64 totalCount,

    /** 符合筛选条件的成功记录数 */
    2: optional i64 successCount,

    /** 符合筛选条件的平均耗时（毫秒） */
    3: optional i32 avgDurationMs,
}

/** 查询审计日志响应 */
struct ListAuditLogsResponse {

    /** 审计日志列表 */
    1: optional list<identity_model.AuditLog> auditLogs,

    /** 分页信息 */
    2: optional base.PageResponse page,

    /** 全局统计信息（基于相同筛选条件，不受分页限制） */
    3: optional AuditLogStats stats,
}

// =================================================================
// OAuth2 客户端管理 (OAuth2 Client Management)
// =================================================================

/** 创建 OAuth2 客户端请求 */
struct CreateOAuth2ClientRequest {

    /** 客户端名称 */
    1: optional string clientName,

    /** 客户端描述 */
    2: optional string description,

    /** 客户端类型 */
    3: optional oauth2_model.OAuth2ClientType clientType,

    /** 允许的授权类型 */
    4: optional list<oauth2_model.OAuth2GrantType> grantTypes,

    /** 允许的回调地址 */
    5: optional list<string> redirectURIs,

    /** 允许的作用域 */
    6: optional list<string> scopes,

    /** 客户端 Logo URL */
    7: optional string logoURI,

    /** 客户端主页 URL */
    8: optional string clientURI,

    /** Access Token 有效期（秒） */
    9: optional i32 accessTokenLifespan,

    /** Refresh Token 有效期（秒） */
    10: optional i32 refreshTokenLifespan,

    /** 创建者用户ID */
    11: optional core.UUID ownerID,
}

/** 创建 OAuth2 客户端响应（含生成的凭证） */
struct CreateOAuth2ClientResponse {

    /** 客户端信息 */
    1: optional oauth2_model.OAuth2Client client,

    /** 明文 client_secret（仅在创建时返回一次） */
    2: optional string clientSecret,
}

/** 获取 OAuth2 客户端请求 */
struct GetOAuth2ClientRequest {

    /** 内部ID */
    1: optional core.UUID id,
}

/** 更新 OAuth2 客户端请求 */
struct UpdateOAuth2ClientRequest {

    /** 内部ID */
    1: optional core.UUID id,

    /** 客户端名称 */
    2: optional string clientName,

    /** 客户端描述 */
    3: optional string description,

    /** 允许的授权类型 */
    4: optional list<oauth2_model.OAuth2GrantType> grantTypes,

    /** 允许的回调地址 */
    5: optional list<string> redirectURIs,

    /** 允许的作用域 */
    6: optional list<string> scopes,

    /** 客户端 Logo URL */
    7: optional string logoURI,

    /** 客户端主页 URL */
    8: optional string clientURI,

    /** Access Token 有效期（秒） */
    9: optional i32 accessTokenLifespan,

    /** Refresh Token 有效期（秒） */
    10: optional i32 refreshTokenLifespan,

    /** 是否启用 */
    11: optional bool isActive,
}

/** 删除 OAuth2 客户端请求 */
struct DeleteOAuth2ClientRequest {

    /** 内部ID */
    1: optional core.UUID id,
}

/** 列出 OAuth2 客户端请求 */
struct ListOAuth2ClientsRequest {

    /** 分页信息 */
    1: optional base.PageRequest page,

    /** 按创建者筛选 */
    2: optional core.UUID ownerID,

    /** 按状态筛选 */
    3: optional bool isActive,
}

/** 列出 OAuth2 客户端响应 */
struct ListOAuth2ClientsResponse {

    /** 客户端列表 */
    1: optional list<oauth2_model.OAuth2Client> clients,

    /** 分页信息 */
    2: optional base.PageResponse page,
}

/** 轮换 OAuth2 客户端密钥请求 */
struct RotateOAuth2ClientSecretRequest {

    /** 内部ID */
    1: optional core.UUID id,
}

/** 轮换 OAuth2 客户端密钥响应 */
struct RotateOAuth2ClientSecretResponse {

    /** 新的明文 client_secret（仅返回一次） */
    1: optional string clientSecret,
}

/** 获取 OAuth2 客户端用于认证请求（含密钥哈希） */
struct GetOAuth2ClientForAuthResponse {

    /** 内部ID */
    1: optional core.UUID id,

    /** 客户端标识符 */
    2: optional string clientID,

    /** 密钥哈希（bcrypt） */
    3: optional string clientSecretHash,

    /** 客户端名称 */
    4: optional string clientName,

    /** 客户端类型 */
    5: optional oauth2_model.OAuth2ClientType clientType,

    /** 允许的授权类型（字符串列表） */
    6: optional list<string> grantTypes,

    /** 允许的回调地址 */
    7: optional list<string> redirectURIs,

    /** 允许的作用域 */
    8: optional list<string> scopes,

    /** Access Token 有效期（秒） */
    9: optional i32 accessTokenLifespan,

    /** Refresh Token 有效期（秒） */
    10: optional i32 refreshTokenLifespan,

    /** 是否启用 */
    11: optional bool isActive,
}

// =================================================================
// OAuth2 作用域管理 (OAuth2 Scope Management)
// =================================================================

/** 列出 OAuth2 作用域请求 */
struct ListOAuth2ScopesRequest {

    /** 是否只返回默认作用域 */
    1: optional bool defaultOnly,
}

/** 列出 OAuth2 作用域响应 */
struct ListOAuth2ScopesResponse {

    /** 作用域列表 */
    1: optional list<oauth2_model.OAuth2Scope> scopes,
}

// =================================================================
// OAuth2 用户授权同意管理 (OAuth2 Consent Management)
// =================================================================

/** 保存 OAuth2 用户授权同意请求 */
struct SaveOAuth2ConsentRequest {

    /** 授权用户ID */
    1: optional core.UUID userID,

    /** 客户端标识符（client_id） */
    2: optional string clientID,

    /** 授权的作用域列表 */
    3: optional list<string> scopes,
}

/** 获取 OAuth2 用户授权同意请求 */
struct GetOAuth2ConsentRequest {

    /** 用户ID */
    1: optional core.UUID userID,

    /** 客户端标识符 */
    2: optional string clientID,
}

/** 获取 OAuth2 用户授权同意响应 */
struct GetOAuth2ConsentResponse {

    /** 授权同意记录（如不存在则为空） */
    1: optional oauth2_model.OAuth2Consent consent,

    /** 是否存在有效的授权同意 */
    2: optional bool found,
}

/** 列出 OAuth2 用户授权同意请求 */
struct ListOAuth2ConsentsRequest {

    /** 用户ID */
    1: optional core.UUID userID,

    /** 分页信息 */
    2: optional base.PageRequest page,
}

/** 列出 OAuth2 用户授权同意响应 */
struct ListOAuth2ConsentsResponse {

    /** 授权同意记录列表 */
    1: optional list<oauth2_model.OAuth2Consent> consents,

    /** 分页信息 */
    2: optional base.PageResponse page,
}

/** 撤销 OAuth2 用户授权同意请求 */
struct RevokeOAuth2ConsentRequest {

    /** 用户ID */
    1: optional core.UUID userID,

    /** 客户端标识符 */
    2: optional string clientID,
}

// =================================================================
// OAuth2 令牌存储 (OAuth2 Token Storage)
// 供 fosite 存储层使用的通用令牌会话结构
// =================================================================

/** OAuth2 令牌会话（通用结构，用于授权码/访问令牌/刷新令牌/PKCE） */
struct OAuth2TokenSession {

    /** 签名/哈希（用作唯一键） */
    1: optional string signature,

    /** 关联的请求ID */
    2: optional string requestID,

    /** 客户端标识符 */
    3: optional string clientID,

    /** 用户ID（Client Credentials 模式可为空） */
    4: optional core.UUID userID,

    /** 授权的作用域（空格分隔） */
    5: optional string scopes,

    /** 授权的受众（空格分隔） */
    6: optional string grantedAudience,

    /** fosite Session 序列化数据（JSON） */
    7: optional binary sessionData,

    /** 请求表单数据（JSON 序列化） */
    8: optional binary formData,

    /** 回调地址 */
    9: optional string redirectURI,

    /** PKCE Code Challenge */
    10: optional string codeChallenge,

    /** PKCE Code Challenge Method */
    11: optional string codeChallengeMethod,

    /** 请求时间 */
    12: optional core.TimestampMS requestedAt,

    /** 过期时间 */
    13: optional core.TimestampMS expiresAt,

    /** 是否已使用（仅授权码） */
    14: optional bool used,

    /** 是否已吊销 */
    15: optional bool revoked,
}