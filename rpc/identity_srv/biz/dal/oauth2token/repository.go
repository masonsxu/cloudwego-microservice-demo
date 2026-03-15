package oauth2token

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

type oauth2TokenRepositoryImpl struct {
	db *gorm.DB
}

func NewOAuth2TokenRepository(db *gorm.DB) OAuth2TokenRepository {
	return &oauth2TokenRepositoryImpl{db: db}
}

// --- 授权码 ---

func (r *oauth2TokenRepositoryImpl) CreateAuthorizationCode(
	ctx context.Context,
	code *models.OAuth2AuthorizationCode,
) error {
	return r.db.WithContext(ctx).Create(code).Error
}

func (r *oauth2TokenRepositoryImpl) GetAuthorizationCode(
	ctx context.Context,
	signature string,
) (*models.OAuth2AuthorizationCode, error) {
	var code models.OAuth2AuthorizationCode
	if err := r.db.WithContext(ctx).Where("signature = ? AND used = false", signature).First(&code).Error; err != nil {
		return nil, fmt.Errorf("查询授权码失败: %w", err)
	}

	return &code, nil
}

func (r *oauth2TokenRepositoryImpl) InvalidateAuthorizationCode(ctx context.Context, signature string) error {
	return r.db.WithContext(ctx).Model(&models.OAuth2AuthorizationCode{}).
		Where("signature = ?", signature).
		Update("used", true).Error
}

// --- 访问令牌 ---

func (r *oauth2TokenRepositoryImpl) CreateAccessToken(ctx context.Context, token *models.OAuth2AccessToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *oauth2TokenRepositoryImpl) GetAccessToken(
	ctx context.Context,
	signature string,
) (*models.OAuth2AccessToken, error) {
	var token models.OAuth2AccessToken
	if err := r.db.WithContext(ctx).
		Where("signature = ? AND is_revoked = false", signature).
		First(&token).
		Error; err != nil {
		return nil, fmt.Errorf("查询访问令牌失败: %w", err)
	}

	return &token, nil
}

func (r *oauth2TokenRepositoryImpl) DeleteAccessToken(ctx context.Context, signature string) error {
	return r.db.WithContext(ctx).Where("signature = ?", signature).Delete(&models.OAuth2AccessToken{}).Error
}

func (r *oauth2TokenRepositoryImpl) RevokeAccessTokenByRequestID(ctx context.Context, requestID string) error {
	return r.db.WithContext(ctx).Model(&models.OAuth2AccessToken{}).
		Where("request_id = ?", requestID).
		Update("is_revoked", true).Error
}

// --- 刷新令牌 ---

func (r *oauth2TokenRepositoryImpl) CreateRefreshToken(ctx context.Context, token *models.OAuth2RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *oauth2TokenRepositoryImpl) GetRefreshToken(
	ctx context.Context,
	signature string,
) (*models.OAuth2RefreshToken, error) {
	var token models.OAuth2RefreshToken
	if err := r.db.WithContext(ctx).
		Where("signature = ? AND is_revoked = false", signature).
		First(&token).
		Error; err != nil {
		return nil, fmt.Errorf("查询刷新令牌失败: %w", err)
	}

	return &token, nil
}

func (r *oauth2TokenRepositoryImpl) DeleteRefreshToken(ctx context.Context, signature string) error {
	return r.db.WithContext(ctx).Where("signature = ?", signature).Delete(&models.OAuth2RefreshToken{}).Error
}

func (r *oauth2TokenRepositoryImpl) RevokeRefreshTokenByRequestID(ctx context.Context, requestID string) error {
	return r.db.WithContext(ctx).Model(&models.OAuth2RefreshToken{}).
		Where("request_id = ?", requestID).
		Update("is_revoked", true).Error
}

// --- PKCE ---

func (r *oauth2TokenRepositoryImpl) CreatePKCESession(ctx context.Context, session *models.OAuth2PKCESession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *oauth2TokenRepositoryImpl) GetPKCESession(
	ctx context.Context,
	signature string,
) (*models.OAuth2PKCESession, error) {
	var session models.OAuth2PKCESession
	if err := r.db.WithContext(ctx).Where("signature = ?", signature).First(&session).Error; err != nil {
		return nil, fmt.Errorf("查询 PKCE 会话失败: %w", err)
	}

	return &session, nil
}

func (r *oauth2TokenRepositoryImpl) DeletePKCESession(ctx context.Context, signature string) error {
	return r.db.WithContext(ctx).Where("signature = ?", signature).Delete(&models.OAuth2PKCESession{}).Error
}
