package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
)

// createTokenInfo 创建Token信息
func createTokenInfo(token string, expire time.Time) *http_base.TokenInfoDTO {
	// 计算token过期时间（秒）
	expiresInSeconds := int64(time.Until(expire).Seconds())

	// 构造Token信息
	tokenType := "Bearer"

	return &http_base.TokenInfoDTO{
		AccessToken: &token,
		TokenType:   &tokenType,
		ExpiresIn:   &expiresInSeconds,
	}
}

// loginResponseHandler 登录响应处理函数
func loginResponseHandler(
	_ context.Context,
	c *app.RequestContext,
	_ int,
	token string,
	expire time.Time,
) {
	// 构造Token信息
	tokenInfo := createTokenInfo(token, expire)

	// 从context中获取登录响应
	if userVal, exists := c.Get(LoginUserContextKey); exists {
		if loginResp, ok := userVal.(*identity.LoginResponseDTO); ok {
			loginResp.TokenInfo = tokenInfo
			c.JSON(http.StatusOK, loginResp)

			return
		}
	}

	// 如果没有找到登录响应,返回空的成功响应
	c.JSON(http.StatusOK, &http_base.BaseResponseDTO{
		Code:    errors.ErrSuccess.Code(),
		Message: errors.ErrSuccess.Message(),
	})
}

// logoutResponseHandler 登出响应处理函数
func logoutResponseHandler(
	_ context.Context,
	c *app.RequestContext,
	_ int,
) {
	// 构造统一的登出响应
	response := &http_base.OperationStatusResponseDTO{
		BaseResp: &http_base.BaseResponseDTO{
			Code:    errors.ErrSuccess.Code(),
			Message: errors.ErrSuccess.Message(),
		},
	}

	c.JSON(http.StatusOK, response)
}

// refreshResponseHandler 刷新Token响应处理函数
func refreshResponseHandler(
	_ context.Context,
	c *app.RequestContext,
	_ int,
	token string,
	expire time.Time,
) {
	// 构造新的Token信息
	tokenInfo := createTokenInfo(token, expire)

	// 构造刷新Token响应
	response := &identity.RefreshTokenResponseDTO{
		BaseResp: &http_base.BaseResponseDTO{
			Code:    errors.ErrSuccess.Code(),
			Message: errors.ErrSuccess.Message(),
		},
		TokenInfo: tokenInfo,
	}

	c.JSON(http.StatusOK, response)
}

// unauthorizedHandler 未认证处理函数
func unauthorizedHandler(
	_ context.Context,
	c *app.RequestContext,
	_ int,
	_ string,
) {
	// 检查响应是否已被写入，如果是则直接返回
	// 避免与 HTTPStatusMessageFunc 产生冲突
	if c.Response.IsBodyStream() || len(c.Response.Body()) > 0 {
		return
	}
}

// customHTTPStatusMessageFunc 自定义HTTP状态消息函数
// 这个函数会在JWT中间件需要返回错误响应时被调用
// 通过自定义这个函数，我们可以控制错误响应的格式，避免与AbortWithError冲突
// 注意：这个函数会被包装在 provider.go 中以适配 hertz-contrib/jwt 的接口
func customHTTPStatusMessageFunc(
	e error,
	ctx context.Context,
	c *app.RequestContext,
	logger *hertzZerolog.Logger,
) string {
	// 检查是否已经有响应被写入（避免重复写入）
	if c.Response.IsBodyStream() || len(c.Response.Body()) > 0 {
		return ""
	}

	var apiError errors.APIError

	// 根据错误类型映射到项目的业务错误码
	switch e {
	case jwt.ErrFailedAuthentication:
		// 认证失败（用户名密码错误等）
		apiError = errors.ErrInvalidCredentials

		logger.Debugf("Authentication failed: error=%v", e)
	case jwt.ErrExpiredToken:
		// Token过期
		apiError = errors.ErrJWTTokenExpired

		logger.Debugf("Token expired: error=%v", e)
	case jwt.ErrFailedTokenCreation:
		// Token创建失败
		apiError = errors.ErrJWTCreationFail

		logger.Warnf("Token creation failed: error=%v", e)
	default:
		// 检查是否是项目内部的业务错误
		if bizErr, ok := e.(errors.APIError); ok {
			apiError = bizErr
			logger.Debugf("Business error: error=%v", bizErr)
		} else {
			// 其他未知错误，默认为JWT验证失败
			apiError = errors.ErrJWTValidationFail

			logger.Warnf("Unknown JWT error: error=%v", e)
		}
	}

	// 生成标准化的错误响应
	httpStatus := errors.GetHTTPStatus(apiError.Code())

	response := &http_base.OperationStatusResponseDTO{
		BaseResp: &http_base.BaseResponseDTO{
			Code:    apiError.Code(),
			Message: apiError.Message(),
		},
	}

	// 直接写入响应
	c.JSON(httpStatus, response)
	c.Abort()

	// 返回空字符串，因为我们已经处理了响应
	return ""
}
