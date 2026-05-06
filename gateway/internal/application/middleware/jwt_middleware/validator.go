package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/context/auth_context"
	authservice "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
)

// identityHandler 身份处理函数，从 JWT Claims 提取用户信息并设置认证上下文（目标态最小字段）
//
// 同时注入下游契约 header（提案 §5.2），供 authz/access_log 中间件以及未来业务系统读取，
// 业务侧禁止再次解析 JWT。
func identityHandler(ctx context.Context, c *app.RequestContext) interface{} {
	claims := jwt.ExtractClaims(ctx, c)
	if claims == nil {
		return nil
	}

	jwtClaims := &http_base.JWTClaimsDTO{}

	if userIDStr, ok := extractStringClaim(claims, IdentityKey); ok {
		jwtClaims.UserProfileID = &userIDStr
	}

	if usernameStr, ok := extractStringClaim(claims, Username); ok {
		jwtClaims.Username = &usernameStr
	}

	if tenantStr, ok := extractStringClaim(claims, Tenant); ok {
		jwtClaims.OrganizationID = &tenantStr
	}

	if roleCodes, ok := extractStringSliceClaim(claims, Roles); ok {
		jwtClaims.RoleIDs = roleCodes
	}

	if expTime, ok := extractInt64Claim(claims, "exp"); ok {
		jwtClaims.Exp = &expTime
	}

	if iatTime, ok := extractInt64Claim(claims, "iat"); ok {
		jwtClaims.Iat = &iatTime
	}

	authCtx := auth_context.NewAuthContext(jwtClaims)
	auth_context.SetAuthContext(c, authCtx)

	injectIdentityHeaders(c, jwtClaims)

	return jwtClaims
}

// injectIdentityHeaders 把验证通过的身份信息注入请求 header，
// 供同链路下游中间件 / 未来业务系统统一读取。
func injectIdentityHeaders(c *app.RequestContext, claims *http_base.JWTClaimsDTO) {
	if claims == nil {
		return
	}

	if claims.UserProfileID != nil && *claims.UserProfileID != "" {
		c.Request.Header.Set(HeaderUserID, *claims.UserProfileID)
	}

	if claims.Username != nil && *claims.Username != "" {
		c.Request.Header.Set(HeaderUserName, *claims.Username)
	}

	if claims.OrganizationID != nil && *claims.OrganizationID != "" {
		c.Request.Header.Set(HeaderTenantID, *claims.OrganizationID)
	}

	if len(claims.RoleIDs) > 0 {
		c.Request.Header.Set(HeaderUserRoles, strings.Join(claims.RoleIDs, ","))
	}
}

// authorizator 授权函数（目标态：不检查 Status，仅做基础校验）
func authorizator(data interface{}, ctx context.Context, c *app.RequestContext) bool {
	return true
}

// authenticatorWithoutAbort 认证函数（不调用 AbortWithError）
func authenticatorWithoutAbort(
	authService authservice.AuthService,
) func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
	return func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
		var (
			err error
			req identity.LoginRequestDTO
		)

		err = c.BindAndValidate(&req)
		if err != nil {
			return nil, errors.ErrInvalidParams.WithMessage(err.Error())
		}

		resp, _, err := authService.Login(ctx, &req)
		if err != nil {
			return nil, err
		}

		userData := buildUserDataMap(resp)

		c.Set(LoginUserContextKey, resp)

		return userData, nil
	}
}
