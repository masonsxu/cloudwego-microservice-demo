/**
 * 身份认证模块 HTTP DTO 定义
 *
 * 定义了身份认证、用户管理、组织架构等相关的所有 HTTP 请求/响应数据传输对象。
 * 遵循 IDL-First 开发模式，支持完整的身份认证和用户管理功能。
 */
namespace go identity

include "../../base/core.thrift"
include "../base/base.thrift"
include "../permission/permission_model.thrift"

// =================================================================
// 1. 身份认证模块 DTO (Authentication)
// =================================================================
// ---- 登录相关 ----

/**
 * 用户登录请求
 * 用户通过用户名密码进行身份认证的请求数据
 */
struct LoginRequestDTO {

    /** 用户名 */
    1: optional string username (api.body = "username", api.vd = "@:len($) > 0; msg:'用户名不能为空'", go.tag = "json:\"username\""),

    /** 密码 */
    2: optional string password (api.body = "password", api.vd = "@:len($) > 0; msg:'密码不能为空'", go.tag = "json:\"password\""),
}

/**
 * 用户登录响应
 * 登录成功后返回的用户信息、权限信息和令牌
 */
struct LoginResponseDTO {

    /** 基础响应信息 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 用户个人信息 */
    2: optional UserProfileDTO userProfile (go.tag = "json:\"user_profile\""),

    /** 用户菜单树（新增字段，推荐使用） */
    3: optional list<permission_model.MenuNodeDTO> menuTree (go.tag = "json:\"menu_tree\""),

    /** 访问令牌信息 */
    4: optional base.TokenInfoDTO tokenInfo (go.tag = "json:\"token_info\""),

    /** 用户成员关系列表 */
    5: optional list<UserMembershipDTO> memberships (go.tag = "json:\"memberships,omitempty\""),

    /** 用户角色ID列表 */
    6: optional list<string> roleIDs (go.tag = "json:\"role_ids,omitempty\""),
}

/**
 * 用户注销请求
 * 用户主动登出时的请求数据
 */
struct LogoutRequestDTO {

    /** 刷新令牌（可选，用于撤销） */
    1: optional string refreshToken (api.body = "refresh_token", api.vd = "@:len($) > 0; msg:'刷新令牌不能为空'", go.tag = "json:\"refresh_token\""),
}

// ---- 密码管理 ----

/**
 * 修改密码请求
 * 用户主动修改自己密码的请求数据
 */
struct ChangePasswordRequestDTO {

    /** 当前密码 */
    1: optional string oldPassword (api.body = "old_password", api.vd = "@:len($)>0; msg:'当前密码不能为空'", go.tag = "json:\"old_password\""),

    /** 新密码 */
    2: optional string newPassword (api.body = "new_password", api.vd = "@:len($)>0 && len($)>=6; msg:'新密码不能为空且长度至少为6位'", go.tag = "json:\"new_password\""),
}

/**
 * 重置密码请求
 * 管理员为用户重置密码的请求数据
 */
struct ResetPasswordRequestDTO {

    /** 目标用户ID */
    1: optional string userID (api.body = "user_id", api.vd = "@:len($)==36; msg:'用户ID格式不正确'", go.tag = "json:\"user_id\""),

    /** 新密码（可选，为空时系统生成） */
    2: optional string newPassword (api.body = "new_password", api.vd = "@:len($)==0 || len($)>=6; msg:'新密码长度至少为6位'", go.tag = "json:\"new_password,omitempty\""),

    /** 重置原因 */
    3: optional string resetReason (api.body = "reset_reason", api.vd = "@:len($)<=200; msg:'重置原因不能超过200个字符'", go.tag = "json:\"reset_reason,omitempty\""),
}

/**
 * 强制修改密码请求
 * 管理员强制用户在下次登录时修改密码
 */
struct ForcePasswordChangeRequestDTO {

    /** 目标用户ID */
    1: optional string userID (api.body = "user_id", api.vd = "@:len($)==36; msg:'用户ID格式不正确'", go.tag = "json:\"user_id\""),

    /** 强制修改原因 */
    2: optional string reason (api.body = "reason", api.vd = "@:len($)<=200; msg:'原因不能超过200个字符'", go.tag = "json:\"reason,omitempty\""),
}

// ---- 令牌管理 ----

/**
 * 刷新访问令牌请求
 * 使用刷新令牌获取新的访问令牌
 */
struct RefreshTokenRequestDTO {

    /** 刷新令牌 */
    1: optional string refreshToken (api.body = "refresh_token", api.vd = "@:len($)>0; msg:'刷新令牌不能为空'", go.tag = "json:\"refresh_token\""),
}

/**
 * 刷新令牌响应
 * 返回新的访问令牌和刷新令牌
 */
struct RefreshTokenResponseDTO {

    /** 基础响应信息 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 新的令牌信息 */
    2: optional base.TokenInfoDTO tokenInfo (go.tag = "json:\"token_info,omitempty\""),
}

// =================================================================
// 3. 用户管理模块 DTO (User Management)
// =================================================================
// ---- 用户基础信息 ----

/**
 * 用户个人信息数据传输对象
 * 包含用户的完整个人信息和状态数据
 */
struct UserProfileDTO {

    /** 用户唯一标识符 */
    1: optional string id (go.tag = "json:\"id\""),

    /** 用户名（登录名） */
    2: optional string username (go.tag = "json:\"username\""),

    /** 邮箱地址 */
    3: optional string email (go.tag = "json:\"email,omitempty\""),

    /** 手机号码 */
    4: optional string phone (go.tag = "json:\"phone,omitempty\""),

    /** 名 */
    5: optional string firstName (go.tag = "json:\"first_name,omitempty\""),

    /** 姓 */
    6: optional string lastName (go.tag = "json:\"last_name,omitempty\""),

    /** 真实姓名 */
    7: optional string realName (go.tag = "json:\"real_name,omitempty\""),

    /** 职业头衔 */
    8: optional string professionalTitle (go.tag = "json:\"professional_title,omitempty\""),

    /** 执业证书号 */
    9: optional string licenseNumber (go.tag = "json:\"license_number,omitempty\""),

    /** 专业特长列表 */
    10: optional list<string> specialties (go.tag = "json:\"specialties,omitempty\""),

    /** 员工工号 */
    11: optional string employeeID (go.tag = "json:\"employee_id,omitempty\""),

    /** 用户状态 */
    12: optional i32 status (go.tag = "json:\"status\""),

    /** 是否必须修改密码 */
    13: optional bool mustChangePassword (go.tag = "json:\"must_change_password,omitempty\""),

    /** 账户过期时间 */
    14: optional core.TimestampMS accountExpiry (go.tag = "json:\"account_expiry,omitempty\""),

    /** 性别 */
    15: optional i32 gender (go.tag = "json:\"gender,omitempty\""),

    /** 创建时间 */
    16: optional core.TimestampMS createdAt (go.tag = "json:\"created_at\""),

    /** 更新时间 */
    17: optional core.TimestampMS updatedAt (go.tag = "json:\"updated_at\""),

    /** 最后登录时间 */
    18: optional core.TimestampMS lastLoginTime (go.tag = "json:\"last_login_time,omitempty\""),

    /** 乐观锁版本号 */
    19: optional i32 version (go.tag = "json:\"version\""),

    /** 连续登录失败次数 */
    20: optional i32 loginAttempts (go.tag = "json:\"login_attempts,omitempty\""),

    /** 创建者用户ID */
    21: optional string createdBy (go.tag = "json:\"created_by,omitempty\""),

    /** 最后更新者用户ID */
    22: optional string updatedBy (go.tag = "json:\"updated_by,omitempty\""),

    /** 用户角色ID列表 */
    23: optional list<string> roleIDs (go.tag = "json:\"role_ids,omitempty\""),

    /** 主组织ID */
    24: optional string primaryOrganizationID (go.tag = "json:\"primary_organization_id,omitempty\""),

    /** 主部门ID */
    25: optional string primaryDepartmentID (go.tag = "json:\"primary_department_id,omitempty\""),
}

/**
 * 用户个人信息响应
 * 单个用户查询的响应数据
 */
struct UserProfileResponseDTO {

    /** 基础响应信息 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 用户个人信息 */
    2: optional UserProfileDTO user (go.tag = "json:\"user,omitempty\""),
}

// ---- 用户操作 (CRUD) ----

/**
 * 创建用户请求
 * 管理员创建新用户的请求数据
 */
struct CreateUserRequestDTO {

    /** 用户名 */
    1: optional string username (api.body = "username", api.vd = "@:len($)>0 && len($)>=3 && len($)<=20 && regexp('^[a-zA-Z0-9_-]+$',$); msg:'用户名必须是3-20位字母、数字、下划线或短横线'", go.tag = "json:\"username\""),

    /** 密码 */
    2: optional string password (api.body = "password", api.vd = "@:len($)>0 && len($)>=6; msg:'密码不能为空且长度至少为6位'", go.tag = "json:\"password\""),

    /** 邮箱地址 */
    3: optional string email (api.body = "email", api.vd = "@:len($)==0 || email($); msg:'邮箱格式不正确'", go.tag = "json:\"email,omitempty\""),

    /** 手机号码 */
    4: optional string phone (api.body = "phone", api.vd = "@:len($)==0 || phone($); msg:'手机号格式不正确'", go.tag = "json:\"phone,omitempty\""),

    /** 名 */
    5: optional string firstName (api.body = "first_name", api.vd = "@:len($)<=50; msg:'名字长度不能超过50个字符'", go.tag = "json:\"first_name,omitempty\""),

    /** 姓 */
    6: optional string lastName (api.body = "last_name", api.vd = "@:len($)<=50; msg:'姓氏长度不能超过50个字符'", go.tag = "json:\"last_name,omitempty\""),

    /** 真实姓名 */
    7: optional string realName (api.body = "real_name", api.vd = "@:len($)<=100; msg:'真实姓名长度不能超过100个字符'", go.tag = "json:\"real_name,omitempty\""),

    /** 职业头衔 */
    8: optional string professionalTitle (api.body = "professional_title", api.vd = "@:len($)<=100; msg:'职业头衔长度不能超过100个字符'", go.tag = "json:\"professional_title,omitempty\""),

    /** 执业证书号 */
    9: optional string licenseNumber (api.body = "license_number", api.vd = "@:len($)<=100; msg:'执业证书号长度不能超过100个字符'", go.tag = "json:\"license_number,omitempty\""),

    /** 专业特长列表 */
    10: optional list<string> specialties (api.body = "specialties", go.tag = "json:\"specialties,omitempty\""),

    /** 员工工号 */
    11: optional string employeeID (api.body = "employee_id", api.vd = "@:len($)<=50; msg:'员工工号长度不能超过50个字符'", go.tag = "json:\"employee_id,omitempty\""),

    /** 是否必须在下次登录时修改密码 */
    12: optional bool mustChangePassword (api.body = "must_change_password", go.tag = "json:\"must_change_password,omitempty\""),

    /** 账户过期时间 */
    13: optional core.TimestampMS accountExpiry (api.body = "account_expiry", go.tag = "json:\"account_expiry,omitempty\""),

    /** 性别（非必填，0:未知, 1:男, 2:女） */
    14: optional i32 gender (api.body = "gender", api.vd = "@:$ == null || ($ >= 0 && $ <= 2); msg:'性别值必须为null或在0-2之间'", go.tag = "json:\"gender,omitempty\""),

    /** 角色ID列表 */
    15: optional list<string> roleIDs (api.body = "role_ids", go.tag = "json:\"role_ids,omitempty\""),

    /** 组织ID */
    16: optional string organizationID (api.body = "organization_id", api.vd = "@:len($)==0 || len($)==36; msg:'组织ID格式不正确'", go.tag = "json:\"organization_id,omitempty\""),
}

/**
 * 获取用户请求
 * 根据用户ID查询用户信息的请求
 */
struct GetUserRequestDTO {

    /** 用户ID */
    1: optional string userID (api.path = "userID", api.vd = "@:len($)==36; msg:'用户ID格式不正确'", go.tag = "json:\"-\""),
}

/**
 * 更新用户请求
 * 更新用户个人信息的请求数据
 */
struct UpdateUserRequestDTO {

    /** 用户ID */
    1: optional string userID (api.path = "userID", api.vd = "@:len($)==36; msg:'用户ID格式不正确'", go.tag = "json:\"-\""),

    /** 邮箱地址 */
    2: optional string email (api.body = "email", api.vd = "@:len($)==0 || email($); msg:'邮箱格式不正确'", go.tag = "json:\"email,omitempty\""),

    /** 手机号码 */
    3: optional string phone (api.body = "phone", api.vd = "@:len($)==0 || phone($); msg:'手机号格式不正确'", go.tag = "json:\"phone,omitempty\""),

    /** 名 */
    4: optional string firstName (api.body = "first_name", api.vd = "@:len($)<=50; msg:'名字长度不能超过50个字符'", go.tag = "json:\"first_name,omitempty\""),

    /** 姓 */
    5: optional string lastName (api.body = "last_name", api.vd = "@:len($)<=50; msg:'姓氏长度不能超过50个字符'", go.tag = "json:\"last_name,omitempty\""),

    /** 真实姓名 */
    6: optional string realName (api.body = "real_name", api.vd = "@:len($)<=100; msg:'真实姓名长度不能超过100个字符'", go.tag = "json:\"real_name,omitempty\""),

    /** 乐观锁版本号 */
    7: optional i32 version (api.body = "version", api.vd = "@:$>=0; msg:'版本号不能为负数'", go.tag = "json:\"version\""),

    /** 职业头衔 */
    8: optional string professionalTitle (api.body = "professional_title", api.vd = "@:len($)<=100; msg:'职业头衔长度不能超过100个字符'", go.tag = "json:\"professional_title,omitempty\""),

    /** 执业证书号 */
    9: optional string licenseNumber (api.body = "license_number", api.vd = "@:len($)<=100; msg:'执业证书号长度不能超过100个字符'", go.tag = "json:\"license_number,omitempty\""),

    /** 专业特长列表 */
    10: optional list<string> specialties (api.body = "specialties", go.tag = "json:\"specialties,omitempty\""),

    /** 员工工号 */
    11: optional string employeeID (api.body = "employee_id", api.vd = "@:len($)<=50; msg:'员工工号长度不能超过50个字符'", go.tag = "json:\"employee_id,omitempty\""),

    /** 账户过期时间 */
    12: optional core.TimestampMS accountExpiry (api.body = "account_expiry", go.tag = "json:\"account_expiry,omitempty\""),

    /** 性别 */
    13: optional i32 gender (api.body = "gender", api.vd = "@:len($)==0 || ($ >= 0 && $ <= 2); msg:'性别值必须在0-2之间'", go.tag = "json:\"gender,omitempty\""),

    /** 角色ID列表 */
    14: optional list<string> roleIDs (api.body = "role_ids", go.tag = "json:\"role_ids,omitempty\""),

    /** 组织ID */
    15: optional string organizationID (api.body = "organization_id", api.vd = "@:len($)==0 || len($)==36; msg:'组织ID格式不正确'", go.tag = "json:\"organization_id,omitempty\""),
}

/**
 * 更新当前用户信息请求
 * 用户更新自己的个人信息（不包含userID，从认证上下文获取）
 */
struct UpdateMeRequestDTO {

    /** 邮箱地址 */
    1: optional string email (api.body = "email", api.vd = "@:len($)==0 || email($); msg:'邮箱格式不正确'", go.tag = "json:\"email,omitempty\""),

    /** 手机号码 */
    2: optional string phone (api.body = "phone", api.vd = "@:len($)==0 || phone($); msg:'手机号格式不正确'", go.tag = "json:\"phone,omitempty\""),

    /** 名 */
    3: optional string firstName (api.body = "first_name", api.vd = "@:len($)<=50; msg:'名字长度不能超过50个字符'", go.tag = "json:\"first_name,omitempty\""),

    /** 姓 */
    4: optional string lastName (api.body = "last_name", api.vd = "@:len($)<=50; msg:'姓氏长度不能超过50个字符'", go.tag = "json:\"last_name,omitempty\""),

    /** 真实姓名 */
    5: optional string realName (api.body = "real_name", api.vd = "@:len($)<=100; msg:'真实姓名长度不能超过100个字符'", go.tag = "json:\"real_name,omitempty\""),

    /** 乐观锁版本号 */
    6: optional i32 version (api.body = "version", api.vd = "@:$>=0; msg:'版本号不能为负数'", go.tag = "json:\"version\""),

    /** 职业头衔 */
    7: optional string professionalTitle (api.body = "professional_title", api.vd = "@:len($)<=100; msg:'职业头衔长度不能超过100个字符'", go.tag = "json:\"professional_title,omitempty\""),

    /** 许可证号 */
    8: optional string licenseNumber (api.body = "medical_license_number", api.vd = "@:len($)<=100; msg:'许可证号长度不能超过100个字符'", go.tag = "json:\"medical_license_number,omitempty\""),

    /** 专业特长列表 */
    9: optional list<string> specialties (api.body = "specialties", go.tag = "json:\"specialties,omitempty\""),

    /** 员工工号 */
    10: optional string employeeID (api.body = "employee_id", api.vd = "@:len($)<=50; msg:'员工工号长度不能超过50个字符'", go.tag = "json:\"employee_id,omitempty\""),

    /** 账户过期时间 */
    11: optional core.TimestampMS accountExpiry (api.body = "account_expiry", go.tag = "json:\"account_expiry,omitempty\""),

    /** 性别 */
    12: optional i32 gender (api.body = "gender", api.vd = "@:$ == null || ($ >= 0 && $ <= 2); msg:'性别值必须为null或在0-2之间'", go.tag = "json:\"gender,omitempty\""),
}

/**
 * 删除用户请求
 * 逻辑删除用户的请求数据
 */
struct DeleteUserRequestDTO {

    /** 用户ID */
    1: optional string userID (api.path = "userID", api.vd = "@:len($)==36; msg:'用户ID格式不正确'", go.tag = "json:\"-\""),

    /** 删除原因 */
    2: optional string reason (api.body = "reason", api.vd = "@:len($)<=200; msg:'删除原因不能超过200个字符'", go.tag = "json:\"reason,omitempty\""),
}

// ---- 用户列表与搜索 ----

/**
 * 用户列表查询请求
 * 支持分页和条件筛选的用户列表查询
 */
struct ListUsersRequestDTO {

    /** 分页信息 */
    1: optional base.PageRequestDTO page (api.none = "true", go.tag = "json:\"page,omitempty\""),

    /** 按组织ID筛选 */
    2: optional string organizationID (api.query = "organization_id", go.tag = "json:\"organization_id,omitempty\""),

    /** 按用户状态筛选 */
    3: optional i32 status (api.query = "status", go.tag = "json:\"status,omitempty\""),
}

/**
 * 用户列表查询响应
 * 包含用户列表和分页信息
 */
struct ListUsersResponseDTO {

    /** 基础响应信息 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 用户列表 */
    2: optional list<UserProfileDTO> users (go.tag = "json:\"users,omitempty\""),

    /** 分页信息 */
    3: optional base.PageResponseDTO page (go.tag = "json:\"page,omitempty\""),
}

/**
 * 用户搜索请求
 * 支持关键字搜索的用户查询
 */
struct SearchUsersRequestDTO {

    /** 分页信息 */
    1: optional base.PageRequestDTO page (api.none = "true", go.tag = "json:\"page,omitempty\""),

    /** 按组织ID筛选 */
    2: optional string organizationID (api.query = "organization_id", go.tag = "json:\"organization_id,omitempty\""),
}

/**
 * 用户搜索响应
 * 包含搜索结果和分页信息
 */
struct SearchUsersResponseDTO {

    /** 基础响应信息 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 用户列表 */
    2: optional list<UserProfileDTO> users (go.tag = "json:\"users,omitempty\""),

    /** 分页信息 */
    3: optional base.PageResponseDTO page (go.tag = "json:\"page,omitempty\""),
}

// ---- 用户状态管理 ----

/**
 * 修改用户状态请求
 * 管理员修改用户状态（激活、禁用等）的请求
 */
struct ChangeUserStatusRequestDTO {

    /** 用户ID */
    1: optional string userID (api.path = "userID", go.tag = "json:\"-\""),

    /** 新状态 */
    2: optional i32 newStatus (api.body = "new_status", go.tag = "json:\"new_status\""),

    /** 修改原因 */
    3: optional string reason (api.body = "reason", go.tag = "json:\"reason,omitempty\""),
}

/**
 * 解锁用户请求
 * 解锁被锁定用户的请求
 */
struct UnlockUserRequestDTO {

    /** 用户ID */
    1: optional string userID (api.path = "userID", go.tag = "json:\"-\""),
}

// =================================================================
// 4. 成员关系管理模块 DTO (Membership Management)
// =================================================================

/**
 * 用户成员关系数据传输对象
 * 描述用户在组织/部门中的成员关系信息
 */
struct UserMembershipDTO {

    /** 成员关系唯一标识符 */
    1: optional string id (go.tag = "json:\"id\""),

    /** 用户ID */
    2: optional string userID (go.tag = "json:\"user_id\""),

    /** 组织ID */
    3: optional string organizationID (go.tag = "json:\"organization_id\""),

    /** 部门ID（可选） */
    4: optional string departmentID (go.tag = "json:\"department_id,omitempty\""),

    /** 是否为主要成员关系 */
    5: optional bool isPrimary (go.tag = "json:\"is_primary,omitempty\""),

    /** 创建时间 */
    6: optional core.TimestampMS createdAt (go.tag = "json:\"created_at\""),

    /** 更新时间 */
    7: optional core.TimestampMS updatedAt (go.tag = "json:\"updated_at\""),
    // 关联信息

    /** 所属组织信息 */
    8: optional OrganizationDTO organization (go.tag = "json:\"organization,omitempty\""),

    /** 所属部门信息 */
    9: optional DepartmentDTO department (go.tag = "json:\"department,omitempty\""),
}

/**
 * 用户成员关系响应
 * 单个成员关系操作的响应数据
 */
struct UserMembershipResponseDTO {

    /** 基础响应信息 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 成员关系信息 */
    2: optional UserMembershipDTO membership (go.tag = "json:\"membership,omitempty\""),
}

/**
 * 获取用户成员关系请求
 * 查询用户的所有成员关系列表
 */
struct GetUserMembershipsRequestDTO {

    /** 用户ID */
    1: optional string userID (api.path = "userID", api.vd = "@:len($)==36; msg:'用户ID格式不正确'", go.tag = "json:\"-\""),

    /** 分页信息 */
    2: optional base.PageRequestDTO page (api.none = "true", go.tag = "json:\"page,omitempty\""),
}

/**
 * 获取用户成员关系响应
 * 包含成员关系列表和分页信息
 */
struct GetUserMembershipsResponseDTO {

    /** 基础响应信息 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 成员关系列表 */
    2: optional list<UserMembershipDTO> memberships (go.tag = "json:\"memberships,omitempty\""),

    /** 分页信息 */
    3: optional base.PageResponseDTO page (go.tag = "json:\"page,omitempty\""),
}

/**
 * 获取用户的主要成员关系。
 * 检查用户是否为某个组织/部门的成员
 */
struct GetPrimaryMembershipRequestDTO {

    /** 用户ID */
    1: optional string userID (api.path = "userID", api.vd = "@:len($)==36; msg:'用户ID格式不正确'", go.tag = "json:\"user_id,omitempty\""),
}

/**
 * 检查成员关系请求
 * 检查用户是否为某个组织/部门的成员
 */
struct CheckMembershipRequestDTO {

    /** 用户ID */
    1: optional string userID (api.query = "user_id", api.vd = "@:len($)==36; msg:'用户ID格式不正确'", go.tag = "json:\"user_id,omitempty\""),

    /** 组织ID */
    2: optional string organizationID (api.query = "organization_id", api.vd = "@:len($)==36; msg:'组织ID格式不正确'", go.tag = "json:\"organization_id,omitempty\""),

    /** 部门ID（可选） */
    3: optional string departmentID (api.query = "department_id", api.vd = "@:len($)==0 || len($)==36; msg:'部门ID格式不正确'", go.tag = "json:\"department_id,omitempty\""),
}

// =================================================================
// 5. 组织架构管理模块 DTO (Organization Management)
// =================================================================

/**
 * 组织信息数据传输对象
 * 完整的组织机构信息，包含层级关系和统计数据
 */
struct OrganizationDTO {

    /** 组织唯一标识符 */
    1: optional string id (go.tag = "json:\"id\""),

    /** 组织代码 */
    2: optional string code (go.tag = "json:\"code,omitempty\""),

    /** 组织名称 */
    3: optional string name (go.tag = "json:\"name\""),

    /** 父组织ID */
    4: optional string parentID (go.tag = "json:\"parent_id,omitempty\""),

    /** 机构类型 */
    5: optional string facilityType (go.tag = "json:\"facility_type,omitempty\""),

    /** 认证状态 */
    6: optional string accreditationStatus (go.tag = "json:\"accreditation_status,omitempty\""),

    /** 组织Logo地址 */
    7: optional string logo (go.tag = "json:\"logo,omitempty\""),

    /** 绑定的Logo ID */
    15: optional string logoID (go.tag = "json:\"logo_id,omitempty\""),

    /** 所在省市列表 */
    8: optional list<string> provinceCity (go.tag = "json:\"province_city,omitempty\""),

    /** 创建时间 */
    9: optional core.TimestampMS createdAt (go.tag = "json:\"created_at\""),

    /** 更新时间 */
    10: optional core.TimestampMS updatedAt (go.tag = "json:\"updated_at\""),
    // 关联信息

    /** 父组织信息 */
    11: optional OrganizationDTO parent (go.tag = "json:\"parent,omitempty\""),

    /** 子组织列表 */
    12: optional list<OrganizationDTO> children (go.tag = "json:\"children,omitempty\""),

    /** 成员数量 */
    13: optional i32 memberCount (go.tag = "json:\"member_count,omitempty\""),

    /** 部门数量 */
    14: optional i32 departmentCount (go.tag = "json:\"department_count,omitempty\""),
}

/**
 * 组织信息响应
 * 单个组织查询的响应数据
 */
struct OrganizationResponseDTO {

    /** 基础响应信息 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 组织信息 */
    2: optional OrganizationDTO organization (go.tag = "json:\"organization,omitempty\""),
}

/**
 * 创建组织请求
 * 创建新组织的请求数据
 */
struct CreateOrganizationRequestDTO {

    /** 组织名称 */
    1: optional string name (api.body = "name", api.vd = "@:len($)>0 && len($)>=2 && len($)<=100; msg:'组织名称长度必须在2-100个字符之间'", go.tag = "json:\"name\""),

    /** 父组织ID（可选） */
    2: optional string parentID (api.body = "parent_id", api.vd = "@:len($)==0 || len($)==36; msg:'父组织ID格式不正确'", go.tag = "json:\"parent_id,omitempty\""),

    /** 机构类型 */
    3: optional string facilityType (api.body = "facility_type", api.vd = "@:len($)<=100; msg:'机构类型长度不能超过100个字符'", go.tag = "json:\"facility_type,omitempty\""),

    /** 认证状态 */
    4: optional string accreditationStatus (api.body = "accreditation_status", api.vd = "@:len($)<=100; msg:'认证状态长度不能超过100个字符'", go.tag = "json:\"accreditation_status,omitempty\""),

    /** 所在省市列表 */
    5: optional list<string> provinceCity (api.body = "province_city", go.tag = "json:\"province_city,omitempty\""),
}

/**
 * 获取组织请求
 * 根据组织ID查询组织信息
 */
struct GetOrganizationRequestDTO {

    /** 组织ID */
    1: optional string organizationID (api.path = "organizationID", api.vd = "@:len($)==36; msg:'组织ID格式不正确'", go.tag = "json:\"-\""),
}

/**
 * 更新组织请求
 * 修改现有组织信息的请求数据
 */
struct UpdateOrganizationRequestDTO {

    /** 组织ID */
    1: optional string organizationID (api.path = "organizationID", api.vd = "@:len($)==36; msg:'组织ID格式不正确'", go.tag = "json:\"-\""),

    /** 组织名称 */
    2: optional string name (api.body = "name", api.vd = "@:len($)==0 || (len($)>=2 && len($)<=100); msg:'名称长度必须在2-100个字符之间'", go.tag = "json:\"name,omitempty\""),

    /** 父组织ID */
    3: optional string parentID (api.body = "parent_id", api.vd = "@:len($)==0 || len($)==36; msg:'父组织ID格式不正确'", go.tag = "json:\"parent_id,omitempty\""),

    /** 机构类型 */
    4: optional string facilityType (api.body = "facility_type", api.vd = "@:len($)<=100; msg:'机构类型长度不能超过100个字符'", go.tag = "json:\"facility_type,omitempty\""),

    /** 认证状态 */
    5: optional string accreditationStatus (api.body = "accreditation_status", api.vd = "@:len($)<=100; msg:'认证状态长度不能超过100个字符'", go.tag = "json:\"accreditation_status,omitempty\""),

    /** 所在省市列表 */
    6: optional list<string> provinceCity (api.body = "province_city", go.tag = "json:\"province_city,omitempty\""),
}

/**
 * 删除组织请求
 * 逻辑删除组织的请求
 */
struct DeleteOrganizationRequestDTO {

    /** 组织ID */
    1: optional string organizationID (api.path = "organizationID", api.vd = "@:len($)==36; msg:'组织ID格式不正确'", go.tag = "json:\"-\""),
}

/**
 * 组织列表查询请求
 * 支持分页和父组织筛选的组织列表查询
 */
struct ListOrganizationsRequestDTO {

    /** 按父组织ID筛选 */
    1: optional string parentID (api.query = "parent_id", go.tag = "query:\"parent_id\""),

    /** 分页信息 */
    2: optional base.PageRequestDTO page (api.none = "true", go.tag = "json:\"page,omitempty\""),
}

/**
 * 组织列表查询响应
 * 包含组织列表和分页信息
 */
struct ListOrganizationsResponseDTO {

    /** 基础响应信息 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 组织列表 */
    2: optional list<OrganizationDTO> organizations (go.tag = "json:\"organizations,omitempty\""),

    /** 分页信息 */
    3: optional base.PageResponseDTO page (go.tag = "json:\"page,omitempty\""),
}

// =================================================================
// 6. 部门管理模块 DTO (Department Management)
// =================================================================

/**
 * 部门信息数据传输对象
 * 完整的部门信息，包含所属组织和统计数据
 */
struct DepartmentDTO {

    /** 部门唯一标识符 */
    1: optional string id (go.tag = "json:\"id\""),

    /** 部门代码 */
    2: optional string code (go.tag = "json:\"code,omitempty\""),

    /** 部门名称 */
    3: optional string name (go.tag = "json:\"name\""),

    /** 所属组织ID */
    4: optional string organizationID (go.tag = "json:\"organization_id\""),

    /** 部门类型 */
    5: optional string departmentType (go.tag = "json:\"department_type,omitempty\""),

    /** 可用设备列表 */
    6: optional list<string> availableEquipment (go.tag = "json:\"available_equipment,omitempty\""),

    /** 创建时间 */
    7: optional core.TimestampMS createdAt (go.tag = "json:\"created_at\""),

    /** 更新时间 */
    8: optional core.TimestampMS updatedAt (go.tag = "json:\"updated_at\""),
    // 关联信息

    /** 所属组织信息 */
    9: optional OrganizationDTO organization (go.tag = "json:\"organization,omitempty\""),

    /** 成员数量 */
    10: optional i32 memberCount (go.tag = "json:\"member_count,omitempty\""),
}

/**
 * 部门信息响应
 * 单个部门查询的响应数据
 */
struct DepartmentResponseDTO {

    /** 基础响应信息 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 部门信息 */
    2: optional DepartmentDTO department (go.tag = "json:\"department,omitempty\""),
}

/**
 * 创建部门请求
 * 在指定组织下创建新部门的请求数据
 */
struct CreateDepartmentRequestDTO {

    /** 所属组织ID */
    1: optional string organizationID (api.body = "organization_id", api.vd = "@:len($)==36; msg:'组织ID格式不正确'", go.tag = "json:\"organization_id\""),

    /** 部门名称 */
    2: optional string name (api.body = "name", api.vd = "@:len($)>0 && len($)>=2 && len($)<=100; msg:'部门名称长度必须在2-100个字符之间'", go.tag = "json:\"name\""),

    /** 部门类型 */
    3: optional string departmentType (api.body = "department_type", api.vd = "@:len($)<=50; msg:'部门类型长度不能超过50个字符'", go.tag = "json:\"department_type,omitempty\""),
}

/**
 * 获取部门请求
 * 根据部门ID查询部门信息
 */
struct GetDepartmentRequestDTO {

    /** 部门ID */
    1: optional string departmentID (api.path = "departmentID", api.vd = "@:len($)==36; msg:'部门ID格式不正确'", go.tag = "json:\"-\""),
}

/**
 * 更新部门请求
 * 修改现有部门信息的请求数据
 */
struct UpdateDepartmentRequestDTO {

    /** 部门ID */
    1: optional string departmentID (api.path = "departmentID", api.vd = "@:len($)==36; msg:'部门ID格式不正确'", go.tag = "json:\"-\""),

    /** 部门名称 */
    2: optional string name (api.body = "name", api.vd = "@:len($)==0 || (len($)>=2 && len($)<=100); msg:'名称长度必须在2-100个字符之间'", go.tag = "json:\"name,omitempty\""),

    /** 部门类型 */
    3: optional string departmentType (api.body = "department_type", api.vd = "@:len($)<=50; msg:'部门类型长度不能超过50个字符'", go.tag = "json:\"department_type,omitempty\""),
}

/**
 * 删除部门请求
 * 逻辑删除部门的请求
 */
struct DeleteDepartmentRequestDTO {

    /** 部门ID */
    1: optional string departmentID (api.path = "departmentID", api.vd = "@:len($)==36; msg:'部门ID格式不正确'", go.tag = "json:\"-\""),
}

/**
 * 获取组织下部门列表请求
 * 查询指定组织下的所有部门列表
 */
struct GetOrganizationDepartmentsRequestDTO {

    /** 组织ID */
    1: optional string organizationID (api.path = "organizationID", api.vd = "@:len($)==36; msg:'组织ID格式不正确'", go.tag = "json:\"-\""),

    /** 分页信息 */
    2: optional base.PageRequestDTO page (api.none = "true", go.tag = "json:\"page,omitempty\""),
}

/**
 * 获取组织下部门列表响应
 * 包含部门列表和分页信息
 */
struct GetOrganizationDepartmentsResponseDTO {

    /** 基础响应信息 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 部门列表 */
    2: optional list<DepartmentDTO> departments (go.tag = "json:\"departments,omitempty\""),

    /** 分页信息 */
    3: optional base.PageResponseDTO page (go.tag = "json:\"page,omitempty\""),
}

// =================================================================
// 7. 组织Logo管理模块 DTO (Organization Logo Management)
// =================================================================

/**
 * 组织Logo数据传输对象
 * 存储组织Logo文件的元数据信息
 */
struct OrganizationLogoDTO {

    /** Logo唯一标识符 */
    1: optional string id (go.tag = "json:\"id\""),

    /** Logo状态（TEMPORARY=临时, BOUND=已绑定, DELETED=已删除） */
    2: optional string status (go.tag = "json:\"status\""),

    /** 绑定的组织ID（临时状态时为空） */
    3: optional string boundOrganizationID (go.tag = "json:\"bound_organization_id,omitempty\""),

    /** 文件存储ID（S3路径: bucket/key） */
    4: optional string fileID (go.tag = "json:\"file_id\""),

    /** 原始文件名 */
    5: optional string fileName (go.tag = "json:\"file_name\""),

    /** 文件大小（字节） */
    6: optional i64 fileSize (go.tag = "json:\"file_size\""),

    /** MIME类型 */
    7: optional string mimeType (go.tag = "json:\"mime_type\""),

    /** 过期时间（临时状态） */
    8: optional core.TimestampMS expiresAt (go.tag = "json:\"expires_at,omitempty\""),

    /** 下载URL（预签名URL） */
    9: optional string downloadUrl (go.tag = "json:\"download_url,omitempty\""),

    /** 上传者ID */
    10: optional string uploadedBy (go.tag = "json:\"uploaded_by\""),

    /** 创建时间 */
    11: optional core.TimestampMS createdAt (go.tag = "json:\"created_at\""),

    /** 更新时间 */
    12: optional core.TimestampMS updatedAt (go.tag = "json:\"updated_at\""),
}

/**
 * Logo响应DTO
 * 单个Logo操作的响应数据
 */
struct OrganizationLogoResponseDTO {

    /** 基础响应信息 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** Logo信息 */
    2: optional OrganizationLogoDTO logo (go.tag = "json:\"logo,omitempty\""),
}

/**
 * 上传临时Logo请求
 * 上传组织Logo文件到临时存储（7天过期）
 */
struct UploadTemporaryLogoRequestDTO {

    /** 文件名 */
    1: optional string fileName (api.body = "file_name", api.vd = "@:len($)>0; msg:'文件名不能为空'", go.tag = "json:\"file_name\""),

    /** 文件内容（二进制） */
    2: optional binary fileContent (api.body = "file_content", go.tag = "json:\"file_content\""),

    /** MIME类型 */
    3: optional string mimeType (api.body = "mime_type", go.tag = "json:\"mime_type,omitempty\""),
}

/**
 * 获取Logo请求
 * 根据Logo ID查询Logo元数据
 */
struct GetOrganizationLogoRequestDTO {

    /** Logo ID */
    1: optional string logoID (api.path = "logoID", api.vd = "@:len($)==36; msg:'Logo ID格式不正确'", go.tag = "json:\"-\""),
}

/**
 * 删除Logo请求
 * 删除Logo文件和数据库记录
 */
struct DeleteOrganizationLogoRequestDTO {

    /** Logo ID */
    1: optional string logoID (api.path = "logoID", api.vd = "@:len($)==36; msg:'Logo ID格式不正确'", go.tag = "json:\"-\""),
}

/**
 * 绑定Logo到组织请求
 * 将临时Logo绑定到组织（永久保存）
 */
struct BindLogoToOrganizationRequestDTO {

    /** 组织ID */
    1: optional string organizationID (api.path = "organizationID", api.vd = "@:len($)==36; msg:'组织ID格式不正确'", go.tag = "json:\"-\""),

    /** Logo ID */
    2: optional string logoID (api.body = "logo_id", api.vd = "@:len($)==36; msg:'Logo ID格式不正确'", go.tag = "json:\"logo_id\""),
}