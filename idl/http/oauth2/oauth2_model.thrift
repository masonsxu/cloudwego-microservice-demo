/**
 * OAuth2 管理模块 HTTP DTO 定义
 *
 * 定义了 OAuth2 客户端管理、作用域管理、用户授权同意管理等
 * HTTP 请求/响应数据传输对象。
 *
 * 注意：OAuth2 核心协议端点（/oauth2/authorize, /oauth2/token 等）
 * 遵循 RFC 6749 规范，由 fosite 直接处理，不经过 IDL 代码生成。
 */
namespace go oauth2

include "../../base/core.thrift"
include "../base/base.thrift"

// =================================================================
// 1. OAuth2 客户端管理 DTO (Client Management)
// =================================================================

/**
 * OAuth2 客户端信息 DTO
 */
struct OAuth2ClientDTO {

    /** 内部唯一ID */
    1: optional string id (go.tag = "json:\"id\""),

    /** 客户端标识符（client_id） */
    2: optional string clientID (go.tag = "json:\"client_id\""),

    /** 客户端名称 */
    3: optional string clientName (go.tag = "json:\"client_name\""),

    /** 客户端描述 */
    4: optional string description (go.tag = "json:\"description,omitempty\""),

    /** 客户端类型（confidential/public） */
    5: optional string clientType (go.tag = "json:\"client_type\""),

    /** 允许的授权类型 */
    6: optional list<string> grantTypes (go.tag = "json:\"grant_types\""),

    /** 允许的回调地址 */
    7: optional list<string> redirectURIs (go.tag = "json:\"redirect_uris,omitempty\""),

    /** 允许的作用域 */
    8: optional list<string> scopes (go.tag = "json:\"scopes,omitempty\""),

    /** Logo URL */
    9: optional string logoURI (go.tag = "json:\"logo_uri,omitempty\""),

    /** 主页 URL */
    10: optional string clientURI (go.tag = "json:\"client_uri,omitempty\""),

    /** Access Token 有效期（秒） */
    11: optional i32 accessTokenLifespan (go.tag = "json:\"access_token_lifespan\""),

    /** Refresh Token 有效期（秒） */
    12: optional i32 refreshTokenLifespan (go.tag = "json:\"refresh_token_lifespan\""),

    /** 是否启用 */
    13: optional bool isActive (go.tag = "json:\"is_active\""),

    /** 创建者用户ID */
    14: optional string ownerID (go.tag = "json:\"owner_id,omitempty\""),

    /** 创建时间 */
    15: optional core.TimestampMS createdAt (go.tag = "json:\"created_at\""),

    /** 更新时间 */
    16: optional core.TimestampMS updatedAt (go.tag = "json:\"updated_at\""),
}

/**
 * 创建 OAuth2 客户端请求
 */
struct CreateOAuth2ClientRequestDTO {

    /** 客户端名称 */
    1: optional string clientName (api.body = "client_name", api.vd = "@:len($) > 0 && len($) <= 128; msg:'客户端名称不能为空且不超过128字符'", go.tag = "json:\"client_name\""),

    /** 客户端描述 */
    2: optional string description (api.body = "description", api.vd = "@:len($) <= 512; msg:'描述不超过512字符'", go.tag = "json:\"description,omitempty\""),

    /** 客户端类型（confidential/public） */
    3: optional string clientType (api.body = "client_type", api.vd = "@:$ == 'confidential' || $ == 'public'; msg:'客户端类型必须为 confidential 或 public'", go.tag = "json:\"client_type\""),

    /** 允许的授权类型 */
    4: optional list<string> grantTypes (api.body = "grant_types", go.tag = "json:\"grant_types\""),

    /** 允许的回调地址 */
    5: optional list<string> redirectURIs (api.body = "redirect_uris", go.tag = "json:\"redirect_uris,omitempty\""),

    /** 允许的作用域 */
    6: optional list<string> scopes (api.body = "scopes", go.tag = "json:\"scopes,omitempty\""),

    /** Logo URL */
    7: optional string logoURI (api.body = "logo_uri", go.tag = "json:\"logo_uri,omitempty\""),

    /** 主页 URL */
    8: optional string clientURI (api.body = "client_uri", go.tag = "json:\"client_uri,omitempty\""),

    /** Access Token 有效期（秒） */
    9: optional i32 accessTokenLifespan (api.body = "access_token_lifespan", go.tag = "json:\"access_token_lifespan,omitempty\""),

    /** Refresh Token 有效期（秒） */
    10: optional i32 refreshTokenLifespan (api.body = "refresh_token_lifespan", go.tag = "json:\"refresh_token_lifespan,omitempty\""),
}

/**
 * 创建 OAuth2 客户端响应（含明文密钥，仅返回一次）
 */
struct CreateOAuth2ClientResponseDTO {

    /** 基础响应 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 客户端信息 */
    2: optional OAuth2ClientDTO client (go.tag = "json:\"client,omitempty\""),

    /** 明文 client_secret（仅创建时返回） */
    3: optional string clientSecret (go.tag = "json:\"client_secret,omitempty\""),
}

/**
 * 获取 OAuth2 客户端请求
 */
struct GetOAuth2ClientRequestDTO {

    /** 客户端内部ID（UUID） */
    1: optional string id (api.path = "id", api.vd = "@:len($)==36; msg:'客户端ID格式不正确'", go.tag = "json:\"-\""),
}

/**
 * OAuth2 客户端响应
 */
struct OAuth2ClientResponseDTO {

    /** 基础响应 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 客户端信息 */
    2: optional OAuth2ClientDTO client (go.tag = "json:\"client,omitempty\""),
}

/**
 * 更新 OAuth2 客户端请求
 */
struct UpdateOAuth2ClientRequestDTO {

    /** 客户端内部ID（UUID） */
    1: optional string id (api.path = "id", api.vd = "@:len($)==36; msg:'客户端ID格式不正确'", go.tag = "json:\"-\""),

    /** 客户端名称 */
    2: optional string clientName (api.body = "client_name", api.vd = "@:len($) == 0 || len($) <= 128; msg:'客户端名称不超过128字符'", go.tag = "json:\"client_name,omitempty\""),

    /** 客户端描述 */
    3: optional string description (api.body = "description", api.vd = "@:len($) <= 512; msg:'描述不超过512字符'", go.tag = "json:\"description,omitempty\""),

    /** 允许的授权类型 */
    4: optional list<string> grantTypes (api.body = "grant_types", go.tag = "json:\"grant_types,omitempty\""),

    /** 允许的回调地址 */
    5: optional list<string> redirectURIs (api.body = "redirect_uris", go.tag = "json:\"redirect_uris,omitempty\""),

    /** 允许的作用域 */
    6: optional list<string> scopes (api.body = "scopes", go.tag = "json:\"scopes,omitempty\""),

    /** Logo URL */
    7: optional string logoURI (api.body = "logo_uri", go.tag = "json:\"logo_uri,omitempty\""),

    /** 主页 URL */
    8: optional string clientURI (api.body = "client_uri", go.tag = "json:\"client_uri,omitempty\""),

    /** Access Token 有效期（秒） */
    9: optional i32 accessTokenLifespan (api.body = "access_token_lifespan", go.tag = "json:\"access_token_lifespan,omitempty\""),

    /** Refresh Token 有效期（秒） */
    10: optional i32 refreshTokenLifespan (api.body = "refresh_token_lifespan", go.tag = "json:\"refresh_token_lifespan,omitempty\""),

    /** 是否启用 */
    11: optional bool isActive (api.body = "is_active", go.tag = "json:\"is_active,omitempty\""),
}

/**
 * 删除 OAuth2 客户端请求
 */
struct DeleteOAuth2ClientRequestDTO {

    /** 客户端内部ID（UUID） */
    1: optional string id (api.path = "id", api.vd = "@:len($)==36; msg:'客户端ID格式不正确'", go.tag = "json:\"-\""),
}

/**
 * 列出 OAuth2 客户端请求
 */
struct ListOAuth2ClientsRequestDTO {

    /** 分页信息 */
    1: optional base.PageRequestDTO page (api.none = "true", go.tag = "json:\"page,omitempty\""),

    /** 按状态筛选 */
    2: optional bool isActive (api.query = "is_active", go.tag = "json:\"is_active,omitempty\""),
}

/**
 * 列出 OAuth2 客户端响应
 */
struct ListOAuth2ClientsResponseDTO {

    /** 基础响应 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 客户端列表 */
    2: optional list<OAuth2ClientDTO> clients (go.tag = "json:\"clients,omitempty\""),

    /** 分页信息 */
    3: optional base.PageResponseDTO page (go.tag = "json:\"page,omitempty\""),
}

/**
 * 轮换客户端密钥请求
 */
struct RotateOAuth2ClientSecretRequestDTO {

    /** 客户端内部ID（UUID） */
    1: optional string id (api.path = "id", api.vd = "@:len($)==36; msg:'客户端ID格式不正确'", go.tag = "json:\"-\""),
}

/**
 * 轮换客户端密钥响应
 */
struct RotateOAuth2ClientSecretResponseDTO {

    /** 基础响应 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 新的明文 client_secret */
    2: optional string clientSecret (go.tag = "json:\"client_secret\""),
}

// =================================================================
// 2. OAuth2 作用域管理 DTO (Scope Management)
// =================================================================

/**
 * OAuth2 作用域 DTO
 */
struct OAuth2ScopeDTO {

    /** 唯一ID */
    1: optional string id (go.tag = "json:\"id\""),

    /** 作用域名称 */
    2: optional string name (go.tag = "json:\"name\""),

    /** 显示名称 */
    3: optional string displayName (go.tag = "json:\"display_name\""),

    /** 描述 */
    4: optional string description (go.tag = "json:\"description,omitempty\""),

    /** 是否为默认作用域 */
    5: optional bool isDefault (go.tag = "json:\"is_default\""),

    /** 是否为系统内置 */
    6: optional bool isSystem (go.tag = "json:\"is_system\""),
}

/**
 * 列出作用域请求
 */
struct ListOAuth2ScopesRequestDTO {
    // 暂无筛选条件
}

/**
 * 列出作用域响应
 */
struct ListOAuth2ScopesResponseDTO {

    /** 基础响应 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 作用域列表 */
    2: optional list<OAuth2ScopeDTO> scopes (go.tag = "json:\"scopes,omitempty\""),
}

/**
 * OAuth2 运行时配置 DTO（只读）
 */
struct OAuth2ConfigDTO {

    /** 是否启用 OAuth2 */
    1: optional bool enabled (go.tag = "json:\"enabled\""),

    /** Issuer */
    2: optional string issuer (go.tag = "json:\"issuer\""),

    /** Access Token 有效期（秒） */
    3: optional i64 accessTokenLifespan (go.tag = "json:\"access_token_lifespan\""),

    /** Refresh Token 有效期（秒） */
    4: optional i64 refreshTokenLifespan (go.tag = "json:\"refresh_token_lifespan\""),

    /** 授权码有效期（秒） */
    5: optional i64 authCodeLifespan (go.tag = "json:\"auth_code_lifespan\""),

    /** 是否强制 PKCE */
    6: optional bool enforcePKCE (go.tag = "json:\"enforce_pkce\""),

    /** 同意页 URL */
    7: optional string consentPageURL (go.tag = "json:\"consent_page_url\""),
}

/**
 * 查询 OAuth2 配置响应
 */
struct GetOAuth2ConfigResponseDTO {

    /** 基础响应 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** OAuth2 配置 */
    2: optional OAuth2ConfigDTO config (go.tag = "json:\"config,omitempty\""),
}

// =================================================================
// 3. OAuth2 用户授权同意管理 DTO (Consent Management)
// =================================================================

/**
 * OAuth2 授权同意 DTO
 */
struct OAuth2ConsentDTO {

    /** 唯一ID */
    1: optional string id (go.tag = "json:\"id\""),

    /** 客户端标识符 */
    2: optional string clientID (go.tag = "json:\"client_id\""),

    /** 客户端名称 */
    3: optional string clientName (go.tag = "json:\"client_name\""),

    /** 授权的作用域 */
    4: optional list<string> scopes (go.tag = "json:\"scopes\""),

    /** 授权时间 */
    5: optional core.TimestampMS grantedAt (go.tag = "json:\"granted_at\""),

    /** 是否已撤销 */
    6: optional bool isRevoked (go.tag = "json:\"is_revoked\""),
}

/**
 * 列出当前用户的授权同意请求
 */
struct ListMyOAuth2ConsentsRequestDTO {

    /** 分页信息 */
    1: optional base.PageRequestDTO page (api.none = "true", go.tag = "json:\"page,omitempty\""),
}

/**
 * 列出授权同意响应
 */
struct ListOAuth2ConsentsResponseDTO {

    /** 基础响应 */
    1: optional base.BaseResponseDTO baseResp (go.tag = "json:\"base_resp\""),

    /** 授权同意记录列表 */
    2: optional list<OAuth2ConsentDTO> consents (go.tag = "json:\"consents,omitempty\""),

    /** 分页信息 */
    3: optional base.PageResponseDTO page (go.tag = "json:\"page,omitempty\""),
}

/**
 * 撤销授权同意请求
 */
struct RevokeMyOAuth2ConsentRequestDTO {

    /** OAuth2 客户端标识符（client_id） */
    1: optional string clientID (api.path = "clientID", api.vd = "@:len($) > 0 && len($) <= 64; msg:'客户端标识符格式不正确'", go.tag = "json:\"-\""),
}
