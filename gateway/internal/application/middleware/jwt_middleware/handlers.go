package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"unicode"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	"github.com/rs/zerolog"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
	tracelog "github.com/masonsxu/cloudwego-microservice-demo/gateway/pkg/log"
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

func camelToSnake(s string) string {
	if s == "" {
		return s
	}

	buf := make([]rune, 0, len(s)+4)
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := runes[i-1]
				nextIsLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
				prevPrevIsLower := i > 1 && unicode.IsLower(runes[i-2])
				shouldSplitAcronymBoundary := nextIsLower && !prevPrevIsLower
				if unicode.IsLower(prev) || unicode.IsDigit(prev) || shouldSplitAcronymBoundary {
					buf = append(buf, '_')
				}
			}
			buf = append(buf, unicode.ToLower(r))
			continue
		}
		buf = append(buf, r)
	}

	return string(buf)
}

func toSnakeCaseValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(val))
		for k, item := range val {
			out[camelToSnake(k)] = toSnakeCaseValue(item)
		}
		return out
	case []interface{}:
		for i := range val {
			val[i] = toSnakeCaseValue(val[i])
		}
		return val
	default:
		return v
	}
}

func toSnakeCasePayload(v interface{}) (interface{}, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var decoded interface{}
	if err = json.Unmarshal(raw, &decoded); err != nil {
		return nil, err
	}

	return toSnakeCaseValue(decoded), nil
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

			if payload, err := toSnakeCasePayload(loginResp); err == nil {
				c.JSON(http.StatusOK, payload)
			} else {
				c.JSON(http.StatusOK, loginResp)
			}

			return
		}
	}

	// 如果没有找到登录响应,返回空的成功响应
	code := errors.ErrSuccess.Code()
	message := errors.ErrSuccess.Message()
	c.JSON(http.StatusOK, &http_base.BaseResponseDTO{
		Code:    &code,
		Message: &message,
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
			Code:    func() *int32 { v := errors.ErrSuccess.Code(); return &v }(),
			Message: func() *string { v := errors.ErrSuccess.Message(); return &v }(),
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
			Code:    func() *int32 { v := errors.ErrSuccess.Code(); return &v }(),
			Message: func() *string { v := errors.ErrSuccess.Message(); return &v }(),
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
	logger *zerolog.Logger,
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

		tracelog.Event(ctx, logger.Debug()).Err(e).Msg("Authentication failed")
	case jwt.ErrExpiredToken:
		// Token过期
		apiError = errors.ErrJWTTokenExpired

		tracelog.Event(ctx, logger.Debug()).Err(e).Msg("Token expired")
	case jwt.ErrFailedTokenCreation:
		// Token创建失败
		apiError = errors.ErrJWTCreationFail

		tracelog.Event(ctx, logger.Warn()).Err(e).Msg("Token creation failed")
	default:
		// 检查是否是项目内部的业务错误
		if bizErr, ok := e.(errors.APIError); ok {
			apiError = bizErr
			tracelog.Event(ctx, logger.Debug()).Err(bizErr).Msg("Business error")
		} else {
			// 其他未知错误，默认为JWT验证失败
			apiError = errors.ErrJWTValidationFail

			tracelog.Event(ctx, logger.Warn()).Err(e).Msg("Unknown JWT error")
		}
	}

	// 生成标准化的错误响应
	httpStatus := errors.GetHTTPStatus(apiError.Code())

	code := apiError.Code()
	message := apiError.Message()
	response := &http_base.OperationStatusResponseDTO{
		BaseResp: &http_base.BaseResponseDTO{
			Code:    &code,
			Message: &message,
		},
	}

	// 直接写入响应
	c.JSON(httpStatus, response)
	c.Abort()

	// 返回空字符串，因为我们已经处理了响应
	return ""
}
