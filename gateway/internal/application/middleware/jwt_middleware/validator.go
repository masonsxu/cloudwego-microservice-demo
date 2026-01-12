package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/core"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/context/auth_context"
	authservice "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
)

// identityHandler 身份处理函数
// 从 JWT Claims 中提取用户信息并设置到认证上下文中
func identityHandler(ctx context.Context, c *app.RequestContext) interface{} {
	claims := jwt.ExtractClaims(ctx, c)
	if claims == nil {
		return nil
	}

	// 构造JWTClaimsDTO对象
	jwtClaims := &http_base.JWTClaimsDTO{}

	// 从claims中提取用户信息并设置到JWTClaimsDTO
	if userIDStr, ok := extractStringClaim(claims, IdentityKey); ok {
		jwtClaims.UserProfileID = &userIDStr
	}

	if usernameStr, ok := extractStringClaim(claims, Username); ok {
		jwtClaims.Username = &usernameStr
	}

	if statusInt, ok := extractIntClaim(claims, Status); ok {
		status := core.UserStatus(statusInt)
		jwtClaims.Status = &status
	}

	// 多角色模式：提取角色ID列表
	if roleIDs, ok := extractStringSliceClaim(claims, RoleIDs); ok && len(roleIDs) > 0 {
		jwtClaims.RoleIDs = roleIDs
	}

	if orgIDStr, ok := extractStringClaim(claims, OrganizationID); ok {
		jwtClaims.OrganizationID = &orgIDStr
	}

	// 多部门模式：提取部门ID列表
	if deptIDs, ok := extractStringSliceClaim(claims, DepartmentIDs); ok && len(deptIDs) > 0 {
		jwtClaims.DepartmentIDs = deptIDs
	}

	if permissionStr, ok := extractStringClaim(claims, CorePermission); ok {
		jwtClaims.Permission = &permissionStr
	}

	// 提取数据范围
	if dataScope, ok := extractStringClaim(claims, DataScope); ok {
		jwtClaims.DataScope = &dataScope
	}

	// 设置JWT标准时间戳声明
	if expTime, ok := extractInt64Claim(claims, "exp"); ok {
		jwtClaims.Exp = &expTime
	}

	if iatTime, ok := extractInt64Claim(claims, "iat"); ok {
		jwtClaims.Iat = &iatTime
	}

	// 创建统一的认证上下文
	authCtx := auth_context.NewAuthContext(jwtClaims)
	auth_context.SetAuthContext(c, authCtx)

	return jwtClaims
}

// checkUserStatusFromClaims 检查用户状态是否为激活状态（从JWT claims直接获取）
func checkUserStatusFromClaims(ctx context.Context, c *app.RequestContext) bool {
	claims := jwt.ExtractClaims(ctx, c)
	if claims != nil {
		if statusInt, ok := extractIntClaim(claims, Status); ok {
			return core.UserStatus(statusInt) == core.UserStatus_ACTIVE
		}
	}

	return false
}

// authorizator 授权函数
// 检查用户状态是否为激活状态，只有激活用户才能通过授权
func authorizator(data interface{}, ctx context.Context, c *app.RequestContext) bool {
	// 检查用户状态是否为激活状态
	if !checkUserStatusFromClaims(ctx, c) {
		return false
	}

	// 检查Token是否被吊销（需要TokenCacheService，这里通过中间件实例访问）
	// 注意：由于authorizator是独立函数，无法直接访问tokenCache
	// 吊销检查将在MiddlewareFunc中处理

	return true
}

// authenticatorWithoutAbort 认证函数（不调用AbortWithError）
// 让JWT中间件通过HTTPStatusMessageFunc统一处理错误响应
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
			// 直接返回错误，不调用AbortWithError
			return nil, errors.ErrInvalidParams.WithMessage(err.Error())
		}

		// 调用业务服务层进行身份验证
		resp, permission, err := authService.Login(ctx, &req)
		if err != nil {
			// 直接返回业务错误，不调用HandleServiceError
			return nil, err
		}

		// 构造用户信息map
		userData := buildUserDataMap(resp, string(permission))

		// 存储用户信息供LoginResponseHandler使用
		c.Set(LoginUserContextKey, resp)

		return userData, nil
	}
}
