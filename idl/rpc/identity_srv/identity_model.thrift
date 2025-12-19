/**
 * 统一身份认证服务 - 数据模型 (Identity Service - Data Models)
 *
 * 定义了与身份、组织、角色、权限等相关的核心数据结构。
 */
namespace go identity_srv

include "../../base/core.thrift"
include "../../base/enums.thrift"

/**
 * 用户画像 (UserProfile)
 * 代表一个“自然人”，存储用户的核心身份信息，与其在具体组织中的角色和职位无关。
 */
struct UserProfile {

    /** 用户唯一ID */
    1: optional core.UUID ID,

    /** 登录用户名，全局唯一 */
    2: optional string username,

    /** 电子邮箱 */
    3: optional string email,

    /** 手机号码 */
    4: optional string phone,
    // --- 个人信息 ---

    /** 名字 */
    5: optional string firstName,

    /** 姓氏 */
    6: optional string lastName,

    /** 真实姓名 */
    7: optional string realName,

    /** 性别 */
    26: optional enums.Gender gender,
    // --- 专业信息 (仅用于展示) ---

    /** 专业职称 */
    8: optional string professionalTitle,

    /** 执业许可证号 */
    9: optional string licenseNumber,

    /** 专业领域列表 */
    10: optional list<string> specialties,
    // --- 状态与安全 ---

    /** 员工工号 */
    11: optional string employeeID,

    /** 用户账户状态 */
    12: optional enums.UserStatus status,

    /** 连续登录失败次数 */
    13: optional i32 loginAttempts = 0,

    /** 是否必须在下次登录时修改密码 */
    14: optional bool mustChangePassword = false,

    /** 账户过期时间 */
    15: optional core.TimestampMS accountExpiry,

    /** 上次登录时间 */
    20: optional core.TimestampMS lastLoginTime,
    // --- 审计与版本控制 ---

    /** 创建时间 */
    16: optional core.TimestampMS createdAt,

    /** 最后更新时间 */
    17: optional core.TimestampMS updatedAt,

    /** 创建者用户ID */
    18: optional core.UUID createdBy,

    /** 最后更新者用户ID */
    19: optional core.UUID updatedBy,

    /** 乐观锁版本号 */
    21: optional i32 version = 1,

    /** 逻辑删除标记 */
    22: optional bool deleted = false,

    /** 用户角色ID列表 */
    23: optional list<core.UUID> roleIDs,

    /** 主组织ID */
    24: optional core.UUID primaryOrganizationID,

    /** 主部门ID */
    25: optional core.UUID primaryDepartmentID,
}

/**
 * 用户成员关系 (UserMembership)
 * 代表一个“自然人”与一个“组织”的隶属关系。一个用户可以有多个成员关系，对应不同的组织。
 */
struct UserMembership {

    /** 成员关系唯一ID */
    1: optional core.UUID ID,
    // --- 核心关系映射 ---

    /** 关联的用户ID */
    2: optional core.UUID userID,

    /** 关联的组织ID */
    3: optional core.UUID organizationID,

    /** 关联的部门ID (可选) */
    4: optional core.UUID departmentID,
    // --- 角色与状态 ---

    /** 是否为该用户的主要成员关系 */
    5: optional bool isPrimary,
    // --- 审计信息 ---

    /** 创建时间 */
    6: optional core.TimestampMS createdAt,

    /** 最后更新时间 */
    7: optional core.TimestampMS updatedAt,
}

/**
 * 组织 (Organization)
 * 代表一个法人实体或机构。
 */
struct Organization {

    /** 组织唯一ID */
    1: optional core.UUID ID,

    /** 组织编码 */
    2: optional string code,

    /** 组织名称 */
    3: optional string name,

    /** 父组织ID，用于构建层级关系 */
    4: optional core.UUID parentID,
    // --- 机构信息 ---

    /** 机构类型 */
    5: optional string facilityType,

    /** 认证状态 (如 JCI, CAP) */
    6: optional string accreditationStatus,

    /** 机构Logo的URL或路径 */
    7: optional string logo,

    /** 绑定的Logo ID */
    15: optional core.UUID logoID,

    /** 机构所在省市 */
    8: optional list<string> provinceCity,
    // --- 审计信息 ---

    /** 创建时间 */
    9: optional core.TimestampMS createdAt,

    /** 最后更新时间 */
    10: optional core.TimestampMS updatedAt,
}

/**
 * 部门 (Department)
 * 代表组织内部的一个部门。
 */
struct Department {

    /** 部门唯一ID */
    1: optional core.UUID ID,

    /** 部门编码 */
    2: optional string code,

    /** 部门名称 */
    3: optional string name,

    /** 所属组织ID */
    4: optional core.UUID organizationID,
    // --- 部门特有信息 ---

    /** 部门类型 (如临床、行政、放疗科) */
    5: optional string departmentType,

    /** 部门可用的设备ID列表 */
    6: optional list<core.UUID> availableEquipment,
    // --- 审计信息 ---

    /** 创建时间 */
    7: optional core.TimestampMS createdAt,

    /** 最后更新时间 */
    8: optional core.TimestampMS updatedAt,
}

/**
 * 组织Logo资源状态枚举
 */
enum OrganizationLogoStatus {

    /** 临时状态 - 已上传但未绑定到组织，7天后自动删除 */
    TEMPORARY = 0,

    /** 已绑定 - 已关联到组织，永久保存 */
    BOUND = 1,

    /** 已删除 - 标记为删除 */
    DELETED = 2,
}

/**
 * 组织Logo (OrganizationLogo)
 * 代表组织的Logo图片资源，支持临时上传和绑定机制。
 */
struct OrganizationLogo {

    /** Logo唯一ID */
    1: optional core.UUID ID,

    /** 文件ID (S3对象键，格式: organization-logos/{uuid}.{ext}) */
    2: optional string fileID,

    /** Logo状态 */
    3: optional OrganizationLogoStatus status,

    /** 绑定的组织ID (null表示临时状态) */
    4: optional core.UUID boundOrganizationID,

    /** 文件名 */
    5: optional string fileName,

    /** 文件大小 (字节) */
    6: optional i64 fileSize,

    /** MIME 类型 */
    7: optional string mimeType,

    /** 下载URL (临时访问链接，仅在需要时生成) */
    8: optional string downloadUrl,

    /** 过期时间 (仅临时状态有效，7天后过期) */
    9: optional core.TimestampMS expiresAt,

    /** 上传者用户ID */
    10: optional core.UUID uploadedBy,

    /** 创建时间 */
    11: optional core.TimestampMS createdAt,

    /** 最后更新时间 */
    12: optional core.TimestampMS updatedAt,
}

/**
 * 权限 (Permission)
 * 定义了一个具体的操作权限，由“资源+动作”构成，并可附加约束条件。
 */
struct Permission {

    /** 权限作用的资源 */
    1: optional string resource,

    /** 对资源执行的操作 */
    2: optional string action,

    /** 权限描述 */
    3: optional string description,
}

/**
 * 角色定义 (RoleDefinition)
 * 定义了一个角色及其拥有的权限集合。
 */
struct RoleDefinition {

    /** 角色唯一ID */
    1: optional core.UUID id,

    /** 角色唯一名称 (中文，页面展示使用) */
    2: optional string name,

    /** 角色详细描述 */
    3: optional string description,

    /** 角色状态 */
    5: optional enums.RoleStatus status,

    /** 该角色拥有的权限列表 */
    6: optional list<Permission> permissions,

    /** 是否为系统内置角色，不可删除 */
    7: optional bool isSystemRole = false,
    // --- 审计信息 ---

    /** 创建者用户ID */
    8: optional core.UUID createdBy,

    /** 更新者用户ID */
    9: optional core.UUID updatedBy,

    /** 创建时间 */
    10: optional core.TimestampMS createdAt,

    /** 最后更新时间 */
    11: optional core.TimestampMS updatedAt,

    /** 当前角色绑定的用户数量（非持久化字段，查询时动态计算） */
    12: optional i64 userCount,
}

/**
 * 用户角色分配 (UserRoleAssignment)
 * 记录了将哪个角色分配给了哪个用户，以及分配的上下文（如组织、部门等）。
 */
struct UserRoleAssignment {

    /** 分配记录的唯一ID */
    1: optional core.UUID id,

    /** 被分配角色的用户ID */
    2: optional core.UUID userID,

    /** 分配的角色ID (对应 RoleDefinition.id) */
    3: optional core.UUID roleID,
    // --- 审计信息 ---

    /** 创建者用户ID */
    11: optional core.UUID createdBy,

    /** 更新者用户ID */
    12: optional core.UUID updatedBy,

    /** 创建时间 */
    13: optional core.TimestampMS createdAt,

    /** 最后更新时间 */
    14: optional core.TimestampMS updatedAt,
}

/**
 * 菜单节点 (Menu Node)
 * 定义了前端动态菜单的单个节点结构，可嵌套形成树。
 */
struct MenuNode {

    /** 菜单名称 (用于显示) */
    1: optional string name,

    /** 菜单唯一标识符 */
    2: optional string id,

    /** 路由路径 */
    3: optional string path,

    /** 菜单图标 (可选) */
    4: optional string icon,

    /** 前端组件路径 (可选) */
    5: optional string component,

    /** 子菜单列表 (可选) */
    6: optional list<MenuNode> children,

    /** 是否有权限访问此菜单 (可选, 用于权限标记) */
    7: optional bool hasPermission,

    /** 权限级别 (可选): read, write, full, none */
    8: optional string permissionLevel,
}