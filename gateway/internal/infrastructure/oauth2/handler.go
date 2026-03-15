// Package oauth2 实现 OAuth2 核心协议端点。
// 这些端点遵循 RFC 6749 规范，由 fosite 直接处理，不经过 IDL 代码生成。
//
// 端点：
//   - GET  /oauth2/authorize   - 授权端点
//   - POST /oauth2/token       - 令牌端点
//   - POST /oauth2/revoke      - 令牌吊销 (RFC 7009)
//   - POST /oauth2/introspect  - 令牌自省 (RFC 7662)
package oauth2

import (
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/ory/fosite"

	oauth2svc "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/oauth2"
)

// Handler OAuth2 核心协议端点处理器
type Handler struct {
	provider fosite.OAuth2Provider
}

// NewHandler 创建 OAuth2 Handler 实例。
func NewHandler(provider fosite.OAuth2Provider) *Handler {
	return &Handler{provider: provider}
}

// AuthorizeEndpoint 处理 GET /oauth2/authorize 请求。
// Authorization Code Grant 的第一步：用户浏览器重定向到此端点进行授权。
//
//	@Summary		OAuth2 授权端点
//	@Description	Authorization Code Grant 的第一步，用户浏览器重定向到此端点进行授权
//	@Tags			OAuth2 协议
//	@Accept			x-www-form-urlencoded
//	@Produce		html
//	@Param			response_type			query	string	true	"响应类型"	Enums(code)
//	@Param			client_id				query	string	true	"客户端ID"
//	@Param			redirect_uri			query	string	true	"重定向URI"
//	@Param			scope					query	string	false	"请求的作用域（空格分隔）"
//	@Param			state					query	string	false	"客户端状态值（防CSRF）"
//	@Param			code_challenge			query	string	false	"PKCE Code Challenge"
//	@Param			code_challenge_method	query	string	false	"PKCE 方法"	Enums(S256, plain)
//	@Success		302						"重定向到 redirect_uri，携带 authorization code"
//	@Failure		400						{object}	object	"请求参数错误"
//	@Router			/oauth2/authorize [GET]
func (h *Handler) AuthorizeEndpoint(_ context.Context, ctx *app.RequestContext) {
	rw := newHertzResponseWriter(ctx)
	r := hertzToHTTPRequest(ctx)

	ar, err := h.provider.NewAuthorizeRequest(r.Context(), r)
	if err != nil {
		log.Printf("[OAuth2] authorize request error: %v", err)
		h.provider.WriteAuthorizeError(r.Context(), rw, ar, err)
		flushResponse(ctx, rw)

		return
	}

	// TODO: 集成用户认证和授权同意流程
	// 当前简化实现：如果用户已登录（通过 JWT cookie），直接授权
	// 生产环境需要：
	// 1. 检查用户是否已登录，未登录则重定向到登录页
	// 2. 检查是否已有授权同意记录
	// 3. 没有则展示授权同意页面
	// 4. 用户确认后调用 h.provider.NewAuthorizeResponse

	session := oauth2svc.NewDefaultSession("user") // TODO: 从认证上下文获取

	for _, scope := range ar.GetRequestedScopes() {
		ar.GrantScope(scope)
	}

	response, err := h.provider.NewAuthorizeResponse(r.Context(), ar, session)
	if err != nil {
		log.Printf("[OAuth2] authorize response error: %v", err)
		h.provider.WriteAuthorizeError(r.Context(), rw, ar, err)
		flushResponse(ctx, rw)

		return
	}

	h.provider.WriteAuthorizeResponse(r.Context(), rw, ar, response)
	flushResponse(ctx, rw)
}

// TokenEndpoint 处理 POST /oauth2/token 请求。
// 支持的 grant_type：authorization_code, client_credentials, refresh_token。
//
//	@Summary		OAuth2 令牌端点
//	@Description	通过授权码、客户端凭证或刷新令牌获取访问令牌
//	@Tags			OAuth2 协议
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			grant_type		formData	string	true	"授权类型"	Enums(authorization_code, client_credentials, refresh_token)
//	@Param			code			formData	string	false	"授权码（authorization_code 模式必填）"
//	@Param			redirect_uri	formData	string	false	"重定向URI（authorization_code 模式必填）"
//	@Param			client_id		formData	string	false	"客户端ID（未使用 Basic Auth 时必填）"
//	@Param			client_secret	formData	string	false	"客户端密钥（未使用 Basic Auth 时必填）"
//	@Param			refresh_token	formData	string	false	"刷新令牌（refresh_token 模式必填）"
//	@Param			code_verifier	formData	string	false	"PKCE Code Verifier"
//	@Success		200				{object}	object	"返回 access_token、token_type、expires_in 等"
//	@Failure		400				{object}	object	"请求参数错误"
//	@Failure		401				{object}	object	"客户端认证失败"
//	@Router			/oauth2/token [POST]
func (h *Handler) TokenEndpoint(_ context.Context, ctx *app.RequestContext) {
	rw := newHertzResponseWriter(ctx)
	r := hertzToHTTPRequest(ctx)

	session := oauth2svc.NewDefaultSession("")

	ar, err := h.provider.NewAccessRequest(r.Context(), r, session)
	if err != nil {
		log.Printf("[OAuth2] token request error: %v", err)
		h.provider.WriteAccessError(r.Context(), rw, ar, err)
		flushResponse(ctx, rw)

		return
	}

	// Client Credentials: 授予所有请求的 scope
	if ar.GetGrantTypes().ExactOne("client_credentials") {
		for _, scope := range ar.GetRequestedScopes() {
			ar.GrantScope(scope)
		}
	}

	response, err := h.provider.NewAccessResponse(r.Context(), ar)
	if err != nil {
		log.Printf("[OAuth2] access response error: %v", err)
		h.provider.WriteAccessError(r.Context(), rw, ar, err)
		flushResponse(ctx, rw)

		return
	}

	h.provider.WriteAccessResponse(r.Context(), rw, ar, response)
	flushResponse(ctx, rw)
}

// RevokeEndpoint 处理 POST /oauth2/revoke 请求 (RFC 7009)。
//
//	@Summary		OAuth2 令牌吊销端点
//	@Description	吊销已颁发的访问令牌或刷新令牌 (RFC 7009)
//	@Tags			OAuth2 协议
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			token			formData	string	true	"要吊销的令牌"
//	@Param			token_type_hint	formData	string	false	"令牌类型提示"	Enums(access_token, refresh_token)
//	@Success		200				"吊销成功"
//	@Failure		401				{object}	object	"客户端认证失败"
//	@Router			/oauth2/revoke [POST]
func (h *Handler) RevokeEndpoint(_ context.Context, ctx *app.RequestContext) {
	rw := newHertzResponseWriter(ctx)
	r := hertzToHTTPRequest(ctx)

	err := h.provider.NewRevocationRequest(r.Context(), r)
	h.provider.WriteRevocationResponse(r.Context(), rw, err)
	flushResponse(ctx, rw)
}

// IntrospectEndpoint 处理 POST /oauth2/introspect 请求 (RFC 7662)。
//
//	@Summary		OAuth2 令牌自省端点
//	@Description	查询令牌的当前状态和元信息 (RFC 7662)
//	@Tags			OAuth2 协议
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			token			formData	string	true	"要自省的令牌"
//	@Param			token_type_hint	formData	string	false	"令牌类型提示"	Enums(access_token, refresh_token)
//	@Success		200				{object}	object	"返回 active、scope、client_id、exp 等"
//	@Failure		401				{object}	object	"客户端认证失败"
//	@Router			/oauth2/introspect [POST]
func (h *Handler) IntrospectEndpoint(_ context.Context, ctx *app.RequestContext) {
	rw := newHertzResponseWriter(ctx)
	r := hertzToHTTPRequest(ctx)

	session := oauth2svc.NewDefaultSession("")

	ir, err := h.provider.NewIntrospectionRequest(r.Context(), r, session)
	if err != nil {
		log.Printf("[OAuth2] introspect error: %v", err)
		h.provider.WriteIntrospectionError(r.Context(), rw, err)
		flushResponse(ctx, rw)

		return
	}

	h.provider.WriteIntrospectionResponse(r.Context(), rw, ir)
	flushResponse(ctx, rw)
}
