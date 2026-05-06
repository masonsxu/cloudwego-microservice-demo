package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/rs/zerolog"

	jwtmw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/jwt_middleware"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
	tracelog "github.com/masonsxu/cloudwego-microservice-demo/gateway/pkg/log"
)

// Outcome 决策结果
type Outcome int

const (
	// OutcomeAllow 放行
	OutcomeAllow Outcome = iota
	// OutcomeUnauthorized 401（缺少身份）
	OutcomeUnauthorized
	// OutcomeForbidden 403（角色不足或默认拒绝）
	OutcomeForbidden
)

// Decision 一次决策的结构化输出，便于日志/审计
type Decision struct {
	Outcome      Outcome
	Reason       string   // 拒绝原因，仅在非 Allow 时有意义
	MatchedRule  string   // 命中的规则类别 / 规则文案，便于调试
	RequiredRole []string // roles 决策中要求的角色，仅 Forbidden 且来自 RolePrefix 时填充
}

// Decide 根据规则与请求属性返回决策结果，纯函数，无副作用
//
// 决策顺序：
//  1. public：放行（无需身份）
//  2. 无 userID：未认证（401）
//  3. authenticated：放行
//  4. roles 前缀：require 与 userRoles 交集判定
//  5. 默认：DefaultAllow → 放行；DefaultDeny → 403
func Decide(rules *Rules, method, path, userID string, userRoles []string) Decision {
	if rules.MatchPublic(method, path) {
		return Decision{Outcome: OutcomeAllow, MatchedRule: "public"}
	}

	if userID == "" {
		return Decision{
			Outcome: OutcomeUnauthorized,
			Reason:  "missing X-User-Id",
		}
	}

	if rules.MatchAuthenticated(method, path) {
		return Decision{Outcome: OutcomeAllow, MatchedRule: "authenticated"}
	}

	if rule, hit := rules.MatchRolePrefix(path); hit {
		if rule.HasAnyRole(userRoles) {
			return Decision{Outcome: OutcomeAllow, MatchedRule: "roles:" + rule.Prefix}
		}

		return Decision{
			Outcome:      OutcomeForbidden,
			Reason:       "required roles not satisfied",
			MatchedRule:  "roles:" + rule.Prefix,
			RequiredRole: append([]string{}, rule.Require...),
		}
	}

	if rules.Default == DefaultDeny {
		return Decision{
			Outcome:     OutcomeForbidden,
			Reason:      "no rule matched and default=deny",
			MatchedRule: "default:deny",
		}
	}

	return Decision{Outcome: OutcomeAllow, MatchedRule: "default:allow"}
}

// AuthZMiddlewareImpl 路由级 ACL 中间件实现
type AuthZMiddlewareImpl struct {
	rules  *Rules
	logger *zerolog.Logger
}

// NewAuthZMiddleware 创建 authz 中间件实例
func NewAuthZMiddleware(rules *Rules, logger *zerolog.Logger) *AuthZMiddlewareImpl {
	return &AuthZMiddlewareImpl{rules: rules, logger: logger}
}

// MiddlewareFunc 返回中间件函数
func (m *AuthZMiddlewareImpl) MiddlewareFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		method := string(c.Request.Method())
		path := string(c.Request.URI().Path())
		userID := c.Request.Header.Get(jwtmw.HeaderUserID)
		userRoles := splitHeader(c.Request.Header.Get(jwtmw.HeaderUserRoles))

		decision := Decide(m.rules, method, path, userID, userRoles)

		switch decision.Outcome {
		case OutcomeAllow:
			c.Next(ctx)
		case OutcomeUnauthorized:
			tracelog.Event(ctx, m.logger.Warn()).
				Str("component", "authz_middleware").
				Str("method", method).
				Str("path", path).
				Str("reason", decision.Reason).
				Msg("authz denied: not authenticated")
			errors.AbortWithError(c, errors.ErrUnauthorized)
		case OutcomeForbidden:
			tracelog.Event(ctx, m.logger.Warn()).
				Str("component", "authz_middleware").
				Str("method", method).
				Str("path", path).
				Str("user_id", userID).
				Strs("user_roles", userRoles).
				Strs("required_roles", decision.RequiredRole).
				Str("matched_rule", decision.MatchedRule).
				Str("reason", decision.Reason).
				Msg("authz denied")
			errors.AbortWithError(c, errors.ErrForbidden)
		}
	}
}

// splitHeader 将逗号分隔的 header 拆分为 trim 过的非空字符串切片
func splitHeader(value string) []string {
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")

	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}

	return out
}
