/**
 * 权限管理模块 HTTP DTO 定义
 *
 * 定义了角色与权限管理相关的所有 HTTP 请求/响应数据传输对象。
 * 遵循 IDL-First 开发模式，支持完整的 RBAC 权限控制。
 */
namespace go permission

include "../../base/core.thrift"
include "../base/base.thrift"

// =================================================================
// 1. 角色与权限管理模块 DTO (RBAC - Role & Permission Management)
// =================================================================

/** 权限DTO */
struct PermissionDTO {

    /** 权限作用的资源 */
    1: optional string resource (go.tag = "json:\"resource\""),

    /** 对资源执行的操作 */
    2: optional string action (go.tag = "json:\"action\""),

    /** 权限描述 */
    3: optional string description (go.tag = "json:\"description,omitempty\""),
}

/** 角色定义DTO */
struct RoleDefinitionDTO {

    /** 角色唯一ID */
    1: optional string id (go.tag = "json:\"id\""),

    /** 角色唯一名称 (中文，页面展示使用) */
    2: optional string name (go.tag = "json:\"name\""),

    /** 角色详细描述 */
    3: optional string description (go.tag = "json:\"description\""),

    /** 角色状态 */
    4: optional i32 status (go.tag = "json:\"status\""),

    /** 该角色拥有的权限列表 */
    5: optional list<PermissionDTO> permissions (go.tag = "json:\"permissions\""),

    /** 是否为系统内置角色，不可删除 */
    6: optional bool isSystemRole (go.tag = "json:\"is_system_role,omitempty\""),

    /** 创建者用户ID */
    7: optional string createdBy (go.tag = "json:\"created_by,omitempty\""),

    /** 更新者用户ID */
    8: optional string updatedBy (go.tag = "json:\"updated_by,omitempty\""),

    /** 创建时间 */
    9: optional core.TimestampMS createdAt (go.tag = "json:\"created_at\""),

    /** 最后更新时间 */
    10: optional core.TimestampMS updatedAt (go.tag = "json:\"updated_at\""),

    /** 当前角色绑定的用户数量（非持久化字段，查询时动态计算） */
    11: optional i64 userCount (go.tag = "json:\"user_count,omitempty\""),
}

/** 用户角色分配DTO */
struct UserRoleAssignmentDTO {

    /** 分配记录的唯一ID */
    1: optional string id (go.tag = "json:\"id\""),

    /** 被分配角色的用户ID */
    2: optional string userID (go.tag = "json:\"user_id\""),

    /** 分配的角色ID (对应 RoleDefinition.id) */
    3: optional string roleID (go.tag = "json:\"role_id\""),

    /** 创建者用户ID */
    4: optional string createdBy (go.tag = "json:\"created_by,omitempty\""),

    /** 更新者用户ID */
    5: optional string updatedBy (go.tag = "json:\"updated_by,omitempty\""),

    /** 创建时间 */
    6: optional core.TimestampMS createdAt (go.tag = "json:\"created_at\""),

    /** 最后更新时间 */
    7: optional core.TimestampMS updatedAt (go.tag = "json:\"updated_at\""),
}

// 角色定义请求/响应 DTO

/** 角色定义创建请求DTO */
struct RoleDefinitionCreateRequestDTO {

    /** 角色唯一名称 (英文，用于程序识别) */
    1: optional string name (api.body = "name", api.vd = "@:len($) > 0; msg:'角色名称不能为空'", go.tag = "json:\"name\""),

    /** 角色描述 */
    2: optional string description (api.body = "description", go.tag = "json:\"description\""),

    /** 角色包含的权限列表 */
    3: optional list<PermissionDTO> permissions (api.body = "permissions", go.tag = "json:\"permissions\""),

    /** 是否为系统内置角色 */
    4: optional bool isSystemRole (api.body = "is_system_role", go.tag = "json:\"is_system_role,omitempty\""),
}

/** 角色定义创建响应DTO */
struct RoleDefinitionCreateResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 创建成功的角色定义信息 */
    2: optional RoleDefinitionDTO role (go.tag = "json:\"role\""),
}

/** 角色定义更新请求DTO */
struct RoleDefinitionUpdateRequestDTO {

    /** 角色定义ID */
    1: optional string roleDefinitionID (api.path = "roleDefinitionID", api.vd = "@:len($) > 0; msg:'角色定义ID不能为空'", go.tag = "json:\"-\""),

    /** 角色描述 */
    2: optional string description (api.body = "description", go.tag = "json:\"description,omitempty\""),

    /** 角色状态 */
    3: optional i32 status (api.body = "status", go.tag = "json:\"status,omitempty\""),

    /** 权限列表 */
    4: optional list<PermissionDTO> permissions (api.body = "permissions", go.tag = "json:\"permissions,omitempty\""),

    /** 角色名称 */
    5: optional string name (api.body = "name", go.tag = "json:\"name,omitempty\""),
}

/** 角色定义更新响应DTO */
struct RoleDefinitionUpdateResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 更新后的角色定义信息 */
    2: optional RoleDefinitionDTO role (go.tag = "json:\"role\""),
}

/** 角色定义删除请求DTO */
struct RoleDefinitionDeleteRequestDTO {

    /** 角色ID */
    1: optional string roleID (api.path = "roleID", api.vd = "@:len($) > 0; msg:'角色ID不能为空'", go.tag = "json:\"-\""),
}

/** 角色定义详情请求DTO */
struct RoleDefinitionGetRequestDTO {

    /** 角色ID */
    1: optional string roleID (api.path = "roleID", api.vd = "@:len($) > 0; msg:'角色ID不能为空'", go.tag = "json:\"-\""),
}

/** 角色定义详情响应DTO */
struct RoleDefinitionGetResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 角色定义信息 */
    2: optional RoleDefinitionDTO role (go.tag = "json:\"role\""),
}

/** 角色定义查询请求DTO */
struct RoleDefinitionQueryRequestDTO {

    /** 角色名称 */
    1: optional string name (api.query = "name", go.tag = "json:\"name,omitempty\""),

    /** 角色状态 */
    2: optional i32 status (api.query = "status", go.tag = "json:\"status,omitempty\""),

    /** 是否为系统角色 */
    3: optional bool isSystemRole (api.query = "isSystemRole", go.tag = "json:\"is_system_role,omitempty\""),

    /** 分页请求参数 */
    4: optional base.PageRequestDTO page (api.none = "true", go.tag = "json:\"page,omitempty\""),
}

/** 角色定义列表响应DTO */
struct RoleDefinitionListResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 角色定义列表 */
    2: optional list<RoleDefinitionDTO> roles (go.tag = "json:\"roles\""),

    /** 分页响应参数 */
    3: optional base.PageResponseDTO page (go.tag = "json:\"page\""),
}

/** 用户角色分配响应DTO */
struct AssignRoleToUserResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 创建的分配记录ID */
    2: optional string assignmentID (go.tag = "json:\"assignment_id\""),
}

/** 用户最后角色分配请求DTO */
struct GetLastUserRoleAssignmentRequestDTO {

    /** 用户ID */
    1: optional string userID (api.path = "userID", api.vd = "@:len($) > 0; msg:'用户ID不能为空'", go.tag = "json:\"user_id\""),
}

/** 用户最后角色分配响应DTO */
struct GetLastUserRoleAssignmentResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 用户角色分配信息 */
    2: optional UserRoleAssignmentDTO assignment (go.tag = "json:\"assignment\""),
}

/** 用户角色查询请求DTO */
struct UserRoleQueryRequestDTO {

    /** 用户ID */
    1: optional string userID (api.query = "userID", go.tag = "json:\"user_id,omitempty\""),

    /** 角色ID */
    2: optional string roleID (api.query = "roleID", go.tag = "json:\"role_id,omitempty\""),

    /** 分页请求参数 */
    3: optional base.PageRequestDTO page (api.none = "true", go.tag = "json:\"page,omitempty\""),
}

/** 用户角色列表响应DTO */
struct UserRoleListResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 用户角色分配列表 */
    2: optional list<UserRoleAssignmentDTO> assignments (go.tag = "json:\"assignments\""),

    /** 分页响应参数 */
    3: optional base.PageResponseDTO page (go.tag = "json:\"page\""),
}

/** 根据角色ID获取用户请求DTO */
struct GetUsersByRoleRequestDTO {

    /** 角色ID */
    1: optional string roleID (api.path = "roleID", api.vd = "@:len($) > 0; msg:'角色ID不能为空'", go.tag = "json:\"-\""),
}

/** 根据角色ID获取用户响应DTO */
struct GetUsersByRoleResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 角色ID */
    2: optional string roleID (go.tag = "json:\"role_id\""),

    /** 该角色下所有用户的ID列表 */
    3: optional list<string> userIDs (go.tag = "json:\"user_ids\""),
}

/** 批量绑定用户到角色请求DTO */
struct BatchBindUsersToRoleRequestDTO {

    /** 角色ID */
    1: optional string roleID (api.path = "roleID", api.vd = "@:len($) > 0; msg:'角色ID不能为空'", go.tag = "json:\"-\""),

    /** 用户ID列表 */
    2: optional list<string> userIDs (api.body = "userIDs", api.vd = "@:len($) > 0; msg:'用户ID列表不能为空'", go.tag = "json:\"user_ids\""),
}

/** 批量绑定用户到角色响应DTO */
struct BatchBindUsersToRoleResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 操作是否成功 */
    2: optional bool success (go.tag = "json:\"success\""),

    /** 成功绑定的用户数量 */
    3: optional i32 successCount (go.tag = "json:\"success_count\""),

    /** 响应消息 */
    4: optional string message (go.tag = "json:\"message,omitempty\""),
}

// =================================================================
// 2. 菜单管理模块 DTO (Menu Management)
// =================================================================

/** 菜单上传请求DTO */
struct UploadMenuRequestDTO {

    /** YAML 格式的菜单配置文件 */
    1: optional binary menuFile (api.form = "menu_file", go.tag = "form:\"menu_file\""),
}

/** 菜单节点DTO */
struct MenuNodeDTO {

    /** 菜单名称 (用于显示) */
    1: optional string name (go.tag = "json:\"name\""),

    /** 菜单唯一标识符 */
    2: optional string id (go.tag = "json:\"id\""),

    /** 路由路径 */
    3: optional string path (go.tag = "json:\"path\""),

    /** 菜单图标 (可选) */
    4: optional string icon (go.tag = "json:\"icon,omitempty\""),

    /** 前端组件路径 (可选) */
    5: optional string component (go.tag = "json:\"component,omitempty\""),

    /** 子菜单列表 (可选) */
    6: optional list<MenuNodeDTO> children (go.tag = "json:\"children,omitempty\""),

    /** 是否有权限访问此菜单 (可选, 用于权限标记) */
    7: optional bool hasPermission (go.tag = "json:\"has_permission,omitempty\""),

    /** 权限级别 (可选): read, write, full, none */
    8: optional string permissionLevel (go.tag = "json:\"permission_level,omitempty\""),
}

/** 菜单树获取响应DTO */
struct GetMenuTreeResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 完整的菜单树结构 */
    2: optional list<MenuNodeDTO> menuTree (go.tag = "json:\"menu_tree\""),
}

// =================================================================
// 3. 菜单权限管理模块 DTO (Menu Permission Management)
// =================================================================

/** 菜单配置项DTO */
struct MenuConfigDTO {

    /** 菜单ID */
    1: optional string menuID (go.tag = "json:\"menu_id\""),

    /** 权限类型: read, write, full, none */
    2: optional string permission (go.tag = "json:\"permission\""),
}

/** 菜单权限信息DTO */
struct MenuPermissionDTO {

    /** 菜单ID */
    1: optional string menuID (go.tag = "json:\"menu_id\""),

    /** 权限类型: read, write, full, none */
    2: optional string permission (go.tag = "json:\"permission\""),
}

/** 配置角色菜单权限请求DTO */
struct ConfigureRoleMenusRequestDTO {

    /** 角色ID（通过路径参数传递，此结构体保留用于扩展） */
    1: optional string roleID (api.path = "roleID", api.vd = "@:len($) > 0; msg:'角色ID不能为空'", go.tag = "json:\"-\""),

    /** 菜单权限配置列表 */
    2: optional list<MenuConfigDTO> menuConfigs (api.body = "menuConfigs", api.vd = "@:len($) > 0; msg:'菜单配置列表不能为空'", go.tag = "json:\"menu_configs\""),
}

/** 配置角色菜单权限响应DTO */
struct ConfigureRoleMenusResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 配置成功标志 */
    2: optional bool success (go.tag = "json:\"success\""),

    /** 响应消息 */
    3: optional string message (go.tag = "json:\"message,omitempty\""),
}

/** 获取角色菜单树请求DTO */
struct GetRoleMenuTreeRequestDTO {

    /** 角色ID（通过路径参数传递，此结构体保留用于扩展） */
    1: optional string roleID (api.path = "roleID", api.vd = "@:len($) > 0; msg:'角色ID不能为空'", go.tag = "json:\"-\""),
}

/** 获取角色菜单树响应DTO */
struct GetRoleMenuTreeResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 角色可访问的菜单树 */
    2: optional list<MenuNodeDTO> menuTree (go.tag = "json:\"menu_tree\""),

    /** 角色ID */
    3: optional string roleID (go.tag = "json:\"role_id\""),
}

/** 获取用户菜单树请求DTO */
struct GetUserMenuTreeRequestDTO {

    /** 用户ID（通过路径参数传递，此结构体保留用于扩展） */
    1: optional string userID (api.path = "userID", api.vd = "@:len($) > 0; msg:'用户ID不能为空'", go.tag = "json:\"-\""),
}

/** 获取用户菜单树响应DTO */
struct GetUserMenuTreeResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 用户可访问的菜单树 */
    2: optional list<MenuNodeDTO> menuTree (go.tag = "json:\"menu_tree\""),

    /** 用户ID */
    3: optional string userID (go.tag = "json:\"user_id\""),

    /** 用户拥有的角色列表 */
    4: optional list<string> roleIDs (go.tag = "json:\"role_ids\""),
}

/** 获取角色菜单权限请求DTO */
struct GetRoleMenuPermissionsRequestDTO {

    /** 角色ID（通过路径参数传递，此结构体保留用于扩展） */
    1: optional string roleID (api.path = "roleID", api.vd = "@:len($) > 0; msg:'角色ID不能为空'", go.tag = "json:\"-\""),
}

/** 获取角色菜单权限响应DTO */
struct GetRoleMenuPermissionsResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 角色的菜单权限列表 */
    2: optional list<MenuPermissionDTO> permissions (go.tag = "json:\"permissions\""),

    /** 角色ID */
    3: optional string roleID (go.tag = "json:\"role_id\""),
}

/** 检查菜单权限请求DTO */
struct HasMenuPermissionRequestDTO {

    /** 角色ID（通过路径参数传递，此结构体保留用于扩展） */
    1: optional string roleID (api.path = "roleID", api.vd = "@:len($) > 0; msg:'角色ID不能为空'", go.tag = "json:\"-\""),

    /** 菜单ID */
    2: optional string menuID (api.query = "menuID", api.vd = "@:len($) > 0; msg:'菜单ID不能为空'", go.tag = "json:\"menu_id\""),

    /** 权限类型 */
    3: optional string permission (api.query = "permission", api.vd = "@:len($) > 0; msg:'权限类型不能为空'", go.tag = "json:\"permission\""),
}

/** 检查菜单权限响应DTO */
struct HasMenuPermissionResponseDTO {

    /** 响应状态码 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 是否具有权限 */
    2: optional bool hasPermission (go.tag = "json:\"has_permission\""),

    /** 角色ID */
    3: optional string roleID (go.tag = "json:\"role_id\""),

    /** 菜单ID */
    4: optional string menuID (go.tag = "json:\"menu_id\""),

    /** 权限类型 */
    5: optional string permission (go.tag = "json:\"permission\""),
}