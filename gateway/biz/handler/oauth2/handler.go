// Package oauth2handler 实现 OAuth2 核心协议端点。
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
	"net/http"
	"net/url"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
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
func (h *Handler) RevokeEndpoint(_ context.Context, ctx *app.RequestContext) {
	rw := newHertzResponseWriter(ctx)
	r := hertzToHTTPRequest(ctx)

	err := h.provider.NewRevocationRequest(r.Context(), r)
	h.provider.WriteRevocationResponse(r.Context(), rw, err)
	flushResponse(ctx, rw)
}

// IntrospectEndpoint 处理 POST /oauth2/introspect 请求 (RFC 7662)。
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

// hertzToHTTPRequest 将 Hertz RequestContext 转换为标准 http.Request。
// fosite 需要标准 http.Request 来解析 OAuth2 请求参数。
func hertzToHTTPRequest(ctx *app.RequestContext) *http.Request {
	u, _ := url.Parse(string(ctx.Request.URI().FullURI()))

	r := &http.Request{
		Method: string(ctx.Method()),
		URL:    u,
		Header: make(http.Header),
		Form:   make(url.Values),
	}

	// 复制请求头
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		r.Header.Set(string(key), string(value))
	})

	// 解析表单数据（POST body）
	ctx.Request.PostArgs().VisitAll(func(key, value []byte) {
		r.Form.Set(string(key), string(value))
	})

	// 也复制 URL query 参数到 Form
	ctx.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		r.Form.Set(string(key), string(value))
	})

	return r
}

// hertzResponseWriter 适配 Hertz RequestContext 为 http.ResponseWriter。
type hertzResponseWriter struct {
	ctx        *app.RequestContext
	statusCode int
	header     http.Header
	body       []byte
}

func newHertzResponseWriter(ctx *app.RequestContext) *hertzResponseWriter {
	return &hertzResponseWriter{
		ctx:        ctx,
		statusCode: consts.StatusOK,
		header:     make(http.Header),
	}
}

func (w *hertzResponseWriter) Header() http.Header {
	return w.header
}

func (w *hertzResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return len(data), nil
}

func (w *hertzResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// flushResponse 将 hertzResponseWriter 中缓冲的响应写回 Hertz RequestContext。
func flushResponse(ctx *app.RequestContext, rw *hertzResponseWriter) {
	for key, values := range rw.header {
		for _, v := range values {
			ctx.Response.Header.Set(key, v)
		}
	}

	ctx.Response.SetStatusCode(rw.statusCode)

	if len(rw.body) > 0 {
		ctx.Response.SetBody(rw.body)
	}
}
