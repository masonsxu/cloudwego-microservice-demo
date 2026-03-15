/**
 * OAuth2 数据模型 (OAuth2 Data Models)
 *
 * 定义 OAuth2 授权服务相关的核心数据结构，包括客户端、授权码、令牌、
 * 作用域和用户授权同意等实体。
 *
 * 设计决策：OAuth2 的 client_type、grant_type 等使用字符串而非 Thrift 枚举，
 * 原因：（1）与 RFC 6749 协议字符串值保持一致（2）减少 HTTP→RPC→Model 间的
 * 枚举转换代码（3）便于扩展新的 grant_type 无需修改 IDL。
 */
namespace go identity_srv

include "../../base/core.thrift"

// =================================================================
// 核心数据模型 (Core Data Models)
// =================================================================

/**
 * OAuth2 客户端（应用）
 * 代表一个注册在本系统中的第三方或内部应用。
 *
 * clientType 取值: "confidential" | "public"
 * grantTypes 取值: "authorization_code" | "client_credentials" | "refresh_token"
 */
struct OAuth2Client {
    /** 内部唯一ID */
    1: optional core.UUID id,

    /** 客户端标识符（对外暴露的 client_id） */
    2: optional string clientID,

    /** 客户端名称（显示用） */
    3: optional string clientName,

    /** 客户端描述 */
    4: optional string description,

    /** 客户端类型: "confidential" | "public" */
    5: optional string clientType,

    /** 允许的授权类型列表: ["authorization_code", "client_credentials", "refresh_token"] */
    6: optional list<string> grantTypes,

    /** 允许的回调地址列表 */
    7: optional list<string> redirectURIs,

    /** 允许的作用域列表 */
    8: optional list<string> scopes,

    /** 客户端 Logo URL */
    9: optional string logoURI,

    /** 客户端主页 URL */
    10: optional string clientURI,

    /** Access Token 有效期（秒） */
    11: optional i32 accessTokenLifespan,

    /** Refresh Token 有效期（秒） */
    12: optional i32 refreshTokenLifespan,

    /** 是否启用 */
    13: optional bool isActive,

    /** 创建者用户ID */
    14: optional core.UUID ownerID,

    /** 创建时间 */
    15: optional core.TimestampMS createdAt,

    /** 更新时间 */
    16: optional core.TimestampMS updatedAt,
}

/**
 * OAuth2 作用域定义
 * 定义 OAuth2 中可授权的权限范围。
 */
struct OAuth2Scope {
    /** 唯一ID */
    1: optional core.UUID id,

    /** 作用域名称（如 user:read） */
    2: optional string name,

    /** 显示名称 */
    3: optional string displayName,

    /** 作用域描述 */
    4: optional string description,

    /** 是否为默认授予的作用域 */
    5: optional bool isDefault,

    /** 是否为系统内置（不可删除） */
    6: optional bool isSystem,

    /** 创建时间 */
    7: optional core.TimestampMS createdAt,

    /** 更新时间 */
    8: optional core.TimestampMS updatedAt,
}

/**
 * OAuth2 用户授权同意记录
 * 记录用户对第三方应用的授权同意状态。
 */
struct OAuth2Consent {
    /** 唯一ID */
    1: optional core.UUID id,

    /** 授权用户ID */
    2: optional core.UUID userID,

    /** 授权的客户端ID（client_id） */
    3: optional string clientID,

    /** 客户端名称（冗余，方便显示） */
    4: optional string clientName,

    /** 授权的作用域列表 */
    5: optional list<string> scopes,

    /** 授权时间 */
    6: optional core.TimestampMS grantedAt,

    /** 授权过期时间 */
    7: optional core.TimestampMS expiresAt,

    /** 是否已撤销 */
    8: optional bool isRevoked,

    /** 创建时间 */
    9: optional core.TimestampMS createdAt,

    /** 更新时间 */
    10: optional core.TimestampMS updatedAt,
}
