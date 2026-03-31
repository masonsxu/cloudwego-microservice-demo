package oauth2

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/mock"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

type stubOAuth2TokenRepo struct {
	createAuthorizationCodeFn       func(ctx context.Context, code *models.OAuth2AuthorizationCode) error
	getAuthorizationCodeFn          func(ctx context.Context, signature string) (*models.OAuth2AuthorizationCode, error)
	invalidateAuthorizationCodeFn   func(ctx context.Context, signature string) error
	createAccessTokenFn             func(ctx context.Context, token *models.OAuth2AccessToken) error
	getAccessTokenFn                func(ctx context.Context, signature string) (*models.OAuth2AccessToken, error)
	deleteAccessTokenFn             func(ctx context.Context, signature string) error
	revokeAccessTokenByRequestIDFn  func(ctx context.Context, requestID string) error
	createRefreshTokenFn            func(ctx context.Context, token *models.OAuth2RefreshToken) error
	getRefreshTokenFn               func(ctx context.Context, signature string) (*models.OAuth2RefreshToken, error)
	deleteRefreshTokenFn            func(ctx context.Context, signature string) error
	revokeRefreshTokenByRequestIDFn func(ctx context.Context, requestID string) error
	createPKCESessionFn             func(ctx context.Context, session *models.OAuth2PKCESession) error
	getPKCESessionFn                func(ctx context.Context, signature string) (*models.OAuth2PKCESession, error)
	deletePKCESessionFn             func(ctx context.Context, signature string) error
}

func (s *stubOAuth2TokenRepo) CreateAuthorizationCode(ctx context.Context, code *models.OAuth2AuthorizationCode) error {
	if s.createAuthorizationCodeFn == nil {
		return nil
	}

	return s.createAuthorizationCodeFn(ctx, code)
}

func (s *stubOAuth2TokenRepo) GetAuthorizationCode(
	ctx context.Context,
	signature string,
) (*models.OAuth2AuthorizationCode, error) {
	if s.getAuthorizationCodeFn == nil {
		return nil, nil
	}

	return s.getAuthorizationCodeFn(ctx, signature)
}

func (s *stubOAuth2TokenRepo) InvalidateAuthorizationCode(ctx context.Context, signature string) error {
	if s.invalidateAuthorizationCodeFn == nil {
		return nil
	}

	return s.invalidateAuthorizationCodeFn(ctx, signature)
}

func (s *stubOAuth2TokenRepo) CreateAccessToken(ctx context.Context, token *models.OAuth2AccessToken) error {
	if s.createAccessTokenFn == nil {
		return nil
	}

	return s.createAccessTokenFn(ctx, token)
}

func (s *stubOAuth2TokenRepo) GetAccessToken(ctx context.Context, signature string) (*models.OAuth2AccessToken, error) {
	if s.getAccessTokenFn == nil {
		return nil, nil
	}

	return s.getAccessTokenFn(ctx, signature)
}

func (s *stubOAuth2TokenRepo) DeleteAccessToken(ctx context.Context, signature string) error {
	if s.deleteAccessTokenFn == nil {
		return nil
	}

	return s.deleteAccessTokenFn(ctx, signature)
}

func (s *stubOAuth2TokenRepo) RevokeAccessTokenByRequestID(ctx context.Context, requestID string) error {
	if s.revokeAccessTokenByRequestIDFn == nil {
		return nil
	}

	return s.revokeAccessTokenByRequestIDFn(ctx, requestID)
}

func (s *stubOAuth2TokenRepo) CreateRefreshToken(ctx context.Context, token *models.OAuth2RefreshToken) error {
	if s.createRefreshTokenFn == nil {
		return nil
	}

	return s.createRefreshTokenFn(ctx, token)
}

func (s *stubOAuth2TokenRepo) GetRefreshToken(
	ctx context.Context,
	signature string,
) (*models.OAuth2RefreshToken, error) {
	if s.getRefreshTokenFn == nil {
		return nil, nil
	}

	return s.getRefreshTokenFn(ctx, signature)
}

func (s *stubOAuth2TokenRepo) DeleteRefreshToken(ctx context.Context, signature string) error {
	if s.deleteRefreshTokenFn == nil {
		return nil
	}

	return s.deleteRefreshTokenFn(ctx, signature)
}

func (s *stubOAuth2TokenRepo) RevokeRefreshTokenByRequestID(ctx context.Context, requestID string) error {
	if s.revokeRefreshTokenByRequestIDFn == nil {
		return nil
	}

	return s.revokeRefreshTokenByRequestIDFn(ctx, requestID)
}

func (s *stubOAuth2TokenRepo) CreatePKCESession(ctx context.Context, session *models.OAuth2PKCESession) error {
	if s.createPKCESessionFn == nil {
		return nil
	}

	return s.createPKCESessionFn(ctx, session)
}

func (s *stubOAuth2TokenRepo) GetPKCESession(ctx context.Context, signature string) (*models.OAuth2PKCESession, error) {
	if s.getPKCESessionFn == nil {
		return nil, nil
	}

	return s.getPKCESessionFn(ctx, signature)
}

func (s *stubOAuth2TokenRepo) DeletePKCESession(ctx context.Context, signature string) error {
	if s.deletePKCESessionFn == nil {
		return nil
	}

	return s.deletePKCESessionFn(ctx, signature)
}

func TestLogicImpl_CreateOAuth2AuthorizeCodeSession(t *testing.T) {
	t.Run("成功创建授权码会话", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		m := mock.NewTestMocks(ctrl)
		repo := &stubOAuth2TokenRepo{}

		m.DAL.EXPECT().OAuth2Token().Return(repo).Times(1)

		logic := &logicImpl{dal: m.DAL, conv: m.Converter}
		ctx := context.Background()
		userID := uuid.NewString()
		signature := "sig-auth-code"
		requestID := "req-1"
		clientID := "client-1"
		scopes := "openid profile"
		requestedAt := int64(1700000000000)
		expiresAt := int64(1700000600000)
		called := false

		repo.createAuthorizationCodeFn = func(_ context.Context, code *models.OAuth2AuthorizationCode) error {
			called = true
			assert.Equal(t, signature, code.Signature)
			assert.Equal(t, requestID, code.RequestID)
			assert.Equal(t, clientID, code.ClientID)
			assert.Equal(t, scopes, code.Scopes)
			assert.Equal(t, requestedAt, code.RequestedAt)
			assert.Equal(t, expiresAt, code.ExpiresAt)
			assert.Equal(t, uuid.MustParse(userID), code.UserID)

			return nil
		}

		err := logic.CreateOAuth2AuthorizeCodeSession(ctx, &identity_srv.OAuth2TokenSession{
			Signature:   &signature,
			RequestID:   &requestID,
			ClientID:    &clientID,
			UserID:      &userID,
			Scopes:      &scopes,
			RequestedAt: &requestedAt,
			ExpiresAt:   &expiresAt,
		})

		require.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("userID 非法时返回错误", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		m := mock.NewTestMocks(ctrl)
		logic := &logicImpl{dal: m.DAL, conv: m.Converter}
		ctx := context.Background()

		signature := "sig-auth-code"
		requestID := "req-1"
		clientID := "client-1"
		invalidUserID := "not-uuid"

		err := logic.CreateOAuth2AuthorizeCodeSession(ctx, &identity_srv.OAuth2TokenSession{
			Signature: &signature,
			RequestID: &requestID,
			ClientID:  &clientID,
			UserID:    &invalidUserID,
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "无效的 userID")
	})
}

func TestLogicImpl_GetOAuth2AuthorizeCodeSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock.NewTestMocks(ctrl)
	repo := &stubOAuth2TokenRepo{}
	m.DAL.EXPECT().OAuth2Token().Return(repo).Times(1)

	logic := &logicImpl{dal: m.DAL, conv: m.Converter}
	ctx := context.Background()
	userID := uuid.New()
	signature := "sig-auth-code"
	reqID := "req-1"
	clientID := "client-1"
	scopes := "openid profile"
	requestedAt := int64(1700000000000)
	expiresAt := int64(1700000600000)
	used := false

	repo.getAuthorizationCodeFn = func(_ context.Context, gotSignature string) (*models.OAuth2AuthorizationCode, error) {
		assert.Equal(t, signature, gotSignature)

		return &models.OAuth2AuthorizationCode{
			Signature:   signature,
			RequestID:   reqID,
			ClientID:    clientID,
			UserID:      userID,
			Scopes:      scopes,
			RequestedAt: requestedAt,
			ExpiresAt:   expiresAt,
			Used:        used,
		}, nil
	}

	resp, err := logic.GetOAuth2AuthorizeCodeSession(ctx, signature)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, signature, resp.GetSignature())
	assert.Equal(t, reqID, resp.GetRequestID())
	assert.Equal(t, clientID, resp.GetClientID())
	assert.Equal(t, userID.String(), resp.GetUserID())
	assert.Equal(t, scopes, resp.GetScopes())
	assert.Equal(t, requestedAt, resp.GetRequestedAt())
	assert.Equal(t, expiresAt, resp.GetExpiresAt())
	assert.Equal(t, used, resp.GetUsed())
}

func TestLogicImpl_CreateOAuth2AccessTokenSession(t *testing.T) {
	t.Run("成功创建访问令牌会话（无 userID）", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		m := mock.NewTestMocks(ctrl)
		repo := &stubOAuth2TokenRepo{}
		m.DAL.EXPECT().OAuth2Token().Return(repo).Times(1)

		logic := &logicImpl{dal: m.DAL, conv: m.Converter}
		ctx := context.Background()
		signature := "sig-access"
		requestID := "req-access"
		clientID := "client-1"
		called := false

		repo.createAccessTokenFn = func(_ context.Context, token *models.OAuth2AccessToken) error {
			called = true
			assert.Equal(t, signature, token.Signature)
			assert.Equal(t, requestID, token.RequestID)
			assert.Equal(t, clientID, token.ClientID)
			assert.Nil(t, token.UserID)

			return nil
		}

		err := logic.CreateOAuth2AccessTokenSession(ctx, &identity_srv.OAuth2TokenSession{
			Signature: &signature,
			RequestID: &requestID,
			ClientID:  &clientID,
		})

		require.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("userID 非法时返回错误", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		m := mock.NewTestMocks(ctrl)
		logic := &logicImpl{dal: m.DAL, conv: m.Converter}
		ctx := context.Background()

		signature := "sig-access"
		requestID := "req-access"
		clientID := "client-1"
		invalidUserID := "bad-user-id"

		err := logic.CreateOAuth2AccessTokenSession(ctx, &identity_srv.OAuth2TokenSession{
			Signature: &signature,
			RequestID: &requestID,
			ClientID:  &clientID,
			UserID:    &invalidUserID,
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "无效的 userID")
	})
}

func TestLogicImpl_RevokeOAuth2RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock.NewTestMocks(ctrl)
	repo := &stubOAuth2TokenRepo{}
	m.DAL.EXPECT().OAuth2Token().Return(repo).Times(1)

	logic := &logicImpl{dal: m.DAL, conv: m.Converter}
	ctx := context.Background()
	requestID := "req-refresh"
	expectedErr := errors.New("revoke failed")

	repo.revokeRefreshTokenByRequestIDFn = func(_ context.Context, gotRequestID string) error {
		assert.Equal(t, requestID, gotRequestID)
		return expectedErr
	}

	err := logic.RevokeOAuth2RefreshToken(ctx, requestID)
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}
