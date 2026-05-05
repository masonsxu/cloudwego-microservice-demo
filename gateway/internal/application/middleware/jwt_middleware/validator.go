package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/context/auth_context"
	authservice "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
)

// identityHandler 身份处理函数，从 JWT Claims 提取用户信息并设置认证上下文（目标态最小字段）
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

	return jwtClaims
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
