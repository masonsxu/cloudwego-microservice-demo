package audit_middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/requestid"
	"github.com/rs/zerolog"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/context/auth_context"
	middleware "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/jwt_middleware"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	tracelog "github.com/masonsxu/cloudwego-microservice-demo/gateway/pkg/log"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

const (
	// maxRequestBodySize 请求体最大保留长度（字节）
	maxRequestBodySize = 2048

	// truncatedMarker 截断标记
	truncatedMarker = "[truncated]"
)

// sensitiveFields 需要脱敏的字段名集合
var sensitiveFields = map[string]bool{
	"password":    true,
	"oldPassword": true,
	"newPassword": true,
	"token":       true,
	"secret":      true,
}

// AuditMiddlewareImpl 审计中间件实现
type AuditMiddlewareImpl struct {
	identityClient identitycli.IdentityClient
	logger         *zerolog.Logger
	cookieName     string // JWT token 的 Cookie 名称（HTTPOnly 模式下从 Cookie 提取）
}

// NewAuditMiddleware 创建审计中间件实例
func NewAuditMiddleware(
	identityClient identitycli.IdentityClient,
	logger *zerolog.Logger,
	cookieName string,
) AuditMiddlewareService {
	return &AuditMiddlewareImpl{
		identityClient: identityClient,
		logger:         logger,
		cookieName:     cookieName,
	}
}

// MiddlewareFunc 返回审计中间件函数
// 在请求处理完成后，异步记录写操作和认证事件的审计日志
func (am *AuditMiddlewareImpl) MiddlewareFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		startTime := time.Now()

		// 在 c.Next 之前捕获请求体（c.Next 后可能被消费）
		var requestBody string

		method := string(c.Method())
		path := string(c.Request.URI().Path())

		if isAuditableRequest(method, path) {
			requestBody = sanitizeRequestBody(c.Request.Body())
		}

		// 执行下游处理器
		c.Next(ctx)

		// 只记录写操作和认证事件
		if !isAuditableRequest(method, path) {
			return
		}

		// 构建审计日志条目
		auditReq := am.buildAuditEntry(ctx, c, method, path, requestBody, startTime)

		// 异步发送到 RPC（fire-and-forget）
		go am.sendAuditLog(auditReq)
	}
}

// isAuditableRequest 判断请求是否需要记录审计日志
func isAuditableRequest(method, path string) bool {
	// 记录所有写操作
	switch method {
	case "POST", "PUT", "PATCH", "DELETE":
		return true
	}

	return false
}

// buildAuditEntry 构建审计日志条目
func (am *AuditMiddlewareImpl) buildAuditEntry(
	ctx context.Context,
	c *app.RequestContext,
	method string,
	path string,
	requestBody string,
	startTime time.Time,
) *identity_srv.CreateAuditLogRequest {
	// 计算请求耗时
	durationMs := int32(time.Since(startTime).Milliseconds())

	// HTTP 状态码
	statusCode := int32(c.Response.StatusCode())

	// 判断是否成功（2xx 为成功）
	success := statusCode >= 200 && statusCode < 300

	// 提取请求 ID
	reqID := requestid.Get(c)

	// 提取 TraceID
	traceID := tracelog.GetTraceID(ctx)

	// 提取客户端信息
	clientIP := c.ClientIP()
	userAgent := string(c.Request.Header.UserAgent())

	// 确定操作类型
	action := resolveAction(method, path)

	// 从路径提取资源 ID（如果有）
	resourceID := extractResourceID(path)

	req := &identity_srv.CreateAuditLogRequest{
		RequestID:   &reqID,
		TraceID:     &traceID,
		Action:      &action,
		Resource:    &path,
		ResourceID:  &resourceID,
		StatusCode:  &statusCode,
		Success:     &success,
		ClientIP:    &clientIP,
		UserAgent:   &userAgent,
		RequestBody: &requestBody,
		DurationMs:  &durationMs,
	}

	// 从认证上下文提取用户信息（登录请求可能没有认证上下文）
	if userID, ok := auth_context.GetCurrentUserProfileID(c); ok {
		req.UserID = &userID
	}

	if username, ok := auth_context.GetCurrentUsername(c); ok {
		req.Username = &username
	}

	if orgID, ok := auth_context.GetCurrentOrganizationID(c); ok {
		req.OrganizationID = &orgID
	}

	// 回退逻辑：JWT 跳过路径（login/refresh）没有 AuthContext，从请求中提取用户信息
	if req.UserID == nil {
		switch {
		case strings.Contains(path, "/auth/login"):
			if username := extractUsernameFromBody(requestBody); username != "" {
				req.Username = &username
			}
		case strings.Contains(path, "/auth/refresh"):
			am.fillUserInfoFromJWT(c, req)
		}
	}

	return req
}

// sendAuditLog 异步发送审计日志到 RPC 服务
func (am *AuditMiddlewareImpl) sendAuditLog(req *identity_srv.CreateAuditLogRequest) {
	// 使用独立的 context 避免父 context 取消影响审计日志写入
	ctx := context.Background()

	if _, err := am.identityClient.CreateAuditLog(ctx, req); err != nil {
		am.logger.Warn().
			Err(err).
			Str("resource", req.GetResource()).
			Str("action", req.GetAction().String()).
			Msg("发送审计日志失败")
	}
}

// resolveAction 根据 HTTP 方法和路径确定审计操作类型
func resolveAction(method, path string) identity_srv.AuditAction {
	// 特殊路径处理
	switch {
	case strings.Contains(path, "/auth/login"):
		return identity_srv.AuditAction_AUDIT_ACTION_LOGIN
	case strings.Contains(path, "/auth/refresh"):
		return identity_srv.AuditAction_AUDIT_ACTION_LOGIN
	case strings.Contains(path, "/auth/logout"):
		return identity_srv.AuditAction_AUDIT_ACTION_LOGOUT
	case strings.Contains(path, "/password"):
		return identity_srv.AuditAction_AUDIT_ACTION_PASSWORD_CHANGE
	}

	// 通用 HTTP 方法映射
	switch method {
	case "POST":
		return identity_srv.AuditAction_AUDIT_ACTION_CREATE
	case "PUT", "PATCH":
		return identity_srv.AuditAction_AUDIT_ACTION_UPDATE
	case "DELETE":
		return identity_srv.AuditAction_AUDIT_ACTION_DELETE
	default:
		return identity_srv.AuditAction_AUDIT_ACTION_CREATE
	}
}

// extractResourceID 从 URL 路径提取资源 ID
// 例如：/api/v1/users/550e8400-e29b-41d4-a716-446655440000 → 550e8400-e29b-41d4-a716-446655440000
func extractResourceID(path string) string {
	parts := strings.Split(strings.TrimRight(path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	// 取最后一段，如果看起来像 UUID 则作为资源 ID
	last := parts[len(parts)-1]
	if len(last) == 36 && strings.Count(last, "-") == 4 {
		return last
	}

	return ""
}

// sanitizeRequestBody 对请求体进行脱敏处理
func sanitizeRequestBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	// 尝试作为 JSON 解析并脱敏
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err == nil {
		sanitizeMap(data)

		sanitized, err := json.Marshal(data)
		if err == nil {
			result := string(sanitized)
			if len(result) > maxRequestBodySize {
				return result[:maxRequestBodySize] + truncatedMarker
			}

			return result
		}
	}

	// 非 JSON 请求体，直接截断
	result := string(body)
	if len(result) > maxRequestBodySize {
		return result[:maxRequestBodySize] + truncatedMarker
	}

	return result
}

// sanitizeMap 递归脱敏 map 中的敏感字段
func sanitizeMap(data map[string]interface{}) {
	for key, value := range data {
		if sensitiveFields[key] {
			data[key] = "***"
			continue
		}

		// 递归处理嵌套 map
		if nested, ok := value.(map[string]interface{}); ok {
			sanitizeMap(nested)
		}
	}
}

// extractUsernameFromBody 从已脱敏的请求体 JSON 中提取 username 字段
// login 请求体结构为 {"username":"xxx","password":"***"}
func extractUsernameFromBody(requestBody string) string {
	if requestBody == "" {
		return ""
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(requestBody), &data); err != nil {
		return ""
	}

	if username, ok := data["username"].(string); ok {
		return username
	}

	return ""
}

// fillUserInfoFromJWT 从 JWT token 提取用户信息并填充到审计日志请求
func (am *AuditMiddlewareImpl) fillUserInfoFromJWT(c *app.RequestContext, req *identity_srv.CreateAuditLogRequest) {
	claims := am.extractUserInfoFromJWT(c)
	if claims == nil {
		return
	}

	if v, ok := claims[middleware.IdentityKey].(string); ok {
		req.UserID = &v
	}

	if v, ok := claims[middleware.Username].(string); ok {
		req.Username = &v
	}

	if v, ok := claims[middleware.OrganizationID].(string); ok {
		req.OrganizationID = &v
	}
}

// extractUserInfoFromJWT 从请求中提取 JWT token 并解码用户信息
// 优先从 Authorization header 提取，若为空则尝试从 Cookie 提取（HTTPOnly 模式）
// 只做 base64 解码（不验证签名和过期），用于审计日志回退
func (am *AuditMiddlewareImpl) extractUserInfoFromJWT(c *app.RequestContext) map[string]interface{} {
	token := extractTokenFromRequest(c, am.cookieName)
	if token == "" {
		return nil
	}

	// JWT 格式: header.payload.signature，提取 payload 部分
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil
	}

	// base64url 解码 payload
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil
	}

	return claims
}

// extractTokenFromRequest 从请求中提取 JWT token
// 优先级：Authorization header > Cookie
func extractTokenFromRequest(c *app.RequestContext, cookieName string) string {
	// 1. 优先从 Authorization header 提取
	authHeader := string(c.GetHeader("Authorization"))
	if authHeader != "" {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != authHeader {
			return token
		}
	}

	// 2. 从 Cookie 提取（HTTPOnly 模式下 token 在 Cookie 中）
	if cookieName != "" {
		if cookieToken := string(c.Cookie(cookieName)); cookieToken != "" {
			return cookieToken
		}
	}

	return ""
}
