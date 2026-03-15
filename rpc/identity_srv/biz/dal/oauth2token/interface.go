package oauth2token

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// OAuth2TokenRepository OAuth2 令牌存储数据访问接口
// 统一管理授权码、访问令牌、刷新令牌和 PKCE 会话。
type OAuth2TokenRepository interface {
	// 授权码
	CreateAuthorizationCode(ctx context.Context, code *models.OAuth2AuthorizationCode) error
	GetAuthorizationCode(ctx context.Context, signature string) (*models.OAuth2AuthorizationCode, error)
	InvalidateAuthorizationCode(ctx context.Context, signature string) error

	// 访问令牌
	CreateAccessToken(ctx context.Context, token *models.OAuth2AccessToken) error
	GetAccessToken(ctx context.Context, signature string) (*models.OAuth2AccessToken, error)
	DeleteAccessToken(ctx context.Context, signature string) error
	RevokeAccessTokenByRequestID(ctx context.Context, requestID string) error

	// 刷新令牌
	CreateRefreshToken(ctx context.Context, token *models.OAuth2RefreshToken) error
	GetRefreshToken(ctx context.Context, signature string) (*models.OAuth2RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, signature string) error
	RevokeRefreshTokenByRequestID(ctx context.Context, requestID string) error

	// PKCE
	CreatePKCESession(ctx context.Context, session *models.OAuth2PKCESession) error
	GetPKCESession(ctx context.Context, signature string) (*models.OAuth2PKCESession, error)
	DeletePKCESession(ctx context.Context, signature string) error
}
