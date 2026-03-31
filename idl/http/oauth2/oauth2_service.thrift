/**
 * OAuth2 管理服务 HTTP 接口定义
 *
 * 提供 OAuth2 客户端管理、作用域查询、用户授权同意管理等 REST API。
 *
 * 注意：OAuth2 核心协议端点遵循 RFC 6749，由 fosite 直接处理：
 *   GET  /oauth2/authorize  - 授权端点
 *   POST /oauth2/token      - 令牌端点
 *   POST /oauth2/revoke     - 令牌吊销 (RFC 7009)
 *   POST /oauth2/introspect - 令牌自省 (RFC 7662)
 * 这些端点不通过 IDL 生成，而是手动实现 Hertz Handler。
 */
namespace go oauth2

include "../base/base.thrift"
include "oauth2_model.thrift"

service OAuth2ManagementService {
    // =================================================================
    // 1. OAuth2 客户端管理 (Client Management) - 管理员接口
    // =================================================================

    /**
     * 创建 OAuth2 客户端
     * 注册一个新的第三方或内部应用
     */
    oauth2_model.CreateOAuth2ClientResponseDTO createOAuth2Client(1: oauth2_model.CreateOAuth2ClientRequestDTO req) (api.post = "/api/v1/oauth2/clients"),

    /**
     * 获取 OAuth2 客户端详情
     */
    oauth2_model.OAuth2ClientResponseDTO getOAuth2Client(1: oauth2_model.GetOAuth2ClientRequestDTO req) (api.get = "/api/v1/oauth2/clients/:id"),

    /**
     * 更新 OAuth2 客户端
     */
    oauth2_model.OAuth2ClientResponseDTO updateOAuth2Client(1: oauth2_model.UpdateOAuth2ClientRequestDTO req) (api.put = "/api/v1/oauth2/clients/:id"),

    /**
     * 删除 OAuth2 客户端
     */
    base.OperationStatusResponseDTO deleteOAuth2Client(1: oauth2_model.DeleteOAuth2ClientRequestDTO req) (api.delete = "/api/v1/oauth2/clients/:id"),

    /**
     * 列出 OAuth2 客户端
     */
    oauth2_model.ListOAuth2ClientsResponseDTO listOAuth2Clients(1: oauth2_model.ListOAuth2ClientsRequestDTO req) (api.get = "/api/v1/oauth2/clients"),

    /**
     * 轮换客户端密钥
     */
    oauth2_model.RotateOAuth2ClientSecretResponseDTO rotateOAuth2ClientSecret(1: oauth2_model.RotateOAuth2ClientSecretRequestDTO req) (api.post = "/api/v1/oauth2/clients/:id/rotate-secret"),
    // =================================================================
    // 2. OAuth2 作用域查询 (Scope Query)
    // =================================================================

    /**
     * 列出所有可用的 OAuth2 作用域
     */
    oauth2_model.ListOAuth2ScopesResponseDTO listOAuth2Scopes(1: oauth2_model.ListOAuth2ScopesRequestDTO req) (api.get = "/api/v1/oauth2/scopes"),

    /**
     * 获取 OAuth2 运行时配置（只读）
     */
    oauth2_model.GetOAuth2ConfigResponseDTO getOAuth2Config() (api.get = "/api/v1/oauth2/config"),
    // =================================================================
    // 3. 用户授权同意管理 (Consent Management) - 用户接口
    // =================================================================

    /**
     * 列出当前用户的授权同意记录
     * 用户查看自己授权过的第三方应用
     */
    oauth2_model.ListOAuth2ConsentsResponseDTO listMyOAuth2Consents(1: oauth2_model.ListMyOAuth2ConsentsRequestDTO req) (api.get = "/api/v1/oauth2/consents"),

    /**
     * 撤销当前用户对指定客户端的授权同意
     * 撤销后该客户端将无法再使用用户的 token
     */
    base.OperationStatusResponseDTO revokeMyOAuth2Consent(1: oauth2_model.RevokeMyOAuth2ConsentRequestDTO req) (api.delete = "/api/v1/oauth2/consents/:clientID"),
}
