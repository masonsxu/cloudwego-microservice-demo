package fositestore

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/pkce"

	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	identitysrv "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

type rpcStore struct {
	identityClient identitycli.IdentityClient
}

type sessionPayload struct {
	Session           *fosite.DefaultSession `json:"session"`
	RequestedScopes   []string               `json:"requested_scopes"`
	GrantedScopes     []string               `json:"granted_scopes"`
	RequestedAudience []string               `json:"requested_audience"`
	GrantedAudience   []string               `json:"granted_audience"`
}

func NewRPCStore(identityClient identitycli.IdentityClient) *rpcStore {
	return &rpcStore{identityClient: identityClient}
}

func (s *rpcStore) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	resp, err := s.identityClient.GetOAuth2ClientForAuth(ctx, &identitysrv.GetOAuth2ClientForAuthRequest{ClientID: &id})
	if err != nil {
		return nil, toFositeError(err)
	}

	return &DefaultClient{
		ID:            resp.GetClientID(),
		Secret:        []byte(resp.GetClientSecretHash()),
		RedirectURIs:  resp.GetRedirectURIs(),
		GrantTypes:    resp.GetGrantTypes(),
		ResponseTypes: []string{"code"},
		Scopes:        resp.GetScopes(),
		Public:        resp.GetClientType() == "public",
	}, nil
}

func (s *rpcStore) ClientAssertionJWTValid(_ context.Context, _ string) error {
	return nil
}

func (s *rpcStore) SetClientAssertionJWT(_ context.Context, _ string, _ time.Time) error {
	return nil
}

func (s *rpcStore) CreateAuthorizeCodeSession(ctx context.Context, code string, req fosite.Requester) error {
	tokenSession, err := requesterToTokenSession(code, req)
	if err != nil {
		return err
	}

	_, err = s.identityClient.CreateOAuth2AuthorizeCodeSession(ctx, tokenSession)
	if err != nil {
		return toFositeError(err)
	}

	return nil
}

func (s *rpcStore) GetAuthorizeCodeSession(
	ctx context.Context,
	code string,
	_ fosite.Session,
) (fosite.Requester, error) {
	resp, err := s.identityClient.GetOAuth2AuthorizeCodeSession(
		ctx,
		&identitysrv.GetOAuth2AuthorizeCodeSessionRequest{Signature: &code},
	)
	if err != nil {
		return nil, toFositeError(err)
	}

	return s.tokenSessionToRequester(ctx, resp)
}

func (s *rpcStore) InvalidateAuthorizeCodeSession(ctx context.Context, code string) error {
	_, err := s.identityClient.InvalidateOAuth2AuthorizeCodeSession(
		ctx,
		&identitysrv.InvalidateOAuth2AuthorizeCodeSessionRequest{Signature: &code},
	)
	if err != nil {
		return toFositeError(err)
	}

	return nil
}

func (s *rpcStore) CreateAccessTokenSession(ctx context.Context, signature string, req fosite.Requester) error {
	tokenSession, err := requesterToTokenSession(signature, req)
	if err != nil {
		return err
	}

	_, err = s.identityClient.CreateOAuth2AccessTokenSession(ctx, tokenSession)
	if err != nil {
		return toFositeError(err)
	}

	return nil
}

func (s *rpcStore) GetAccessTokenSession(
	ctx context.Context,
	signature string,
	_ fosite.Session,
) (fosite.Requester, error) {
	resp, err := s.identityClient.GetOAuth2AccessTokenSession(
		ctx,
		&identitysrv.GetOAuth2AccessTokenSessionRequest{Signature: &signature},
	)
	if err != nil {
		return nil, toFositeError(err)
	}

	return s.tokenSessionToRequester(ctx, resp)
}

func (s *rpcStore) DeleteAccessTokenSession(ctx context.Context, signature string) error {
	_, err := s.identityClient.DeleteOAuth2AccessTokenSession(
		ctx,
		&identitysrv.DeleteOAuth2AccessTokenSessionRequest{Signature: &signature},
	)
	if err != nil {
		return toFositeError(err)
	}

	return nil
}

func (s *rpcStore) CreateRefreshTokenSession(
	ctx context.Context,
	signature string,
	_ string,
	req fosite.Requester,
) error {
	tokenSession, err := requesterToTokenSession(signature, req)
	if err != nil {
		return err
	}

	_, err = s.identityClient.CreateOAuth2RefreshTokenSession(ctx, tokenSession)
	if err != nil {
		return toFositeError(err)
	}

	return nil
}

func (s *rpcStore) GetRefreshTokenSession(
	ctx context.Context,
	signature string,
	_ fosite.Session,
) (fosite.Requester, error) {
	resp, err := s.identityClient.GetOAuth2RefreshTokenSession(
		ctx,
		&identitysrv.GetOAuth2RefreshTokenSessionRequest{Signature: &signature},
	)
	if err != nil {
		return nil, toFositeError(err)
	}

	return s.tokenSessionToRequester(ctx, resp)
}

func (s *rpcStore) DeleteRefreshTokenSession(ctx context.Context, signature string) error {
	_, err := s.identityClient.DeleteOAuth2RefreshTokenSession(
		ctx,
		&identitysrv.DeleteOAuth2RefreshTokenSessionRequest{Signature: &signature},
	)
	if err != nil {
		return toFositeError(err)
	}

	return nil
}

func (s *rpcStore) RotateRefreshToken(ctx context.Context, requestID string, refreshTokenSignature string) error {
	if err := s.RevokeRefreshToken(ctx, requestID); err != nil {
		return err
	}

	return nil
}

func (s *rpcStore) RevokeRefreshToken(ctx context.Context, requestID string) error {
	_, err := s.identityClient.RevokeOAuth2RefreshToken(
		ctx,
		&identitysrv.RevokeOAuth2RefreshTokenRequest{RequestID: &requestID},
	)
	if err != nil {
		return toFositeError(err)
	}

	return nil
}

func (s *rpcStore) RevokeRefreshTokenMaybeGracePeriod(ctx context.Context, requestID string, _ string) error {
	return s.RevokeRefreshToken(ctx, requestID)
}

func (s *rpcStore) RevokeAccessToken(ctx context.Context, requestID string) error {
	_, err := s.identityClient.RevokeOAuth2AccessToken(
		ctx,
		&identitysrv.RevokeOAuth2AccessTokenRequest{RequestID: &requestID},
	)
	if err != nil {
		return toFositeError(err)
	}

	return nil
}

func (s *rpcStore) CreatePKCERequestSession(ctx context.Context, signature string, req fosite.Requester) error {
	tokenSession, err := requesterToTokenSession(signature, req)
	if err != nil {
		return err
	}

	_, err = s.identityClient.CreateOAuth2PKCESession(ctx, tokenSession)
	if err != nil {
		return toFositeError(err)
	}

	return nil
}

func (s *rpcStore) GetPKCERequestSession(
	ctx context.Context,
	signature string,
	_ fosite.Session,
) (fosite.Requester, error) {
	resp, err := s.identityClient.GetOAuth2PKCESession(
		ctx,
		&identitysrv.GetOAuth2PKCESessionRequest{Signature: &signature},
	)
	if err != nil {
		return nil, toFositeError(err)
	}

	return s.tokenSessionToRequester(ctx, resp)
}

func (s *rpcStore) DeletePKCERequestSession(ctx context.Context, signature string) error {
	_, err := s.identityClient.DeleteOAuth2PKCESession(
		ctx,
		&identitysrv.DeleteOAuth2PKCESessionRequest{Signature: &signature},
	)
	if err != nil {
		return toFositeError(err)
	}

	return nil
}

func requesterToTokenSession(signature string, req fosite.Requester) (*identitysrv.OAuth2TokenSession, error) {
	payload := sessionPayload{
		RequestedScopes:   append([]string{}, req.GetRequestedScopes()...),
		GrantedScopes:     append([]string{}, req.GetGrantedScopes()...),
		RequestedAudience: append([]string{}, req.GetRequestedAudience()...),
		GrantedAudience:   append([]string{}, req.GetGrantedAudience()...),
	}

	if session, ok := req.GetSession().(*fosite.DefaultSession); ok {
		payload.Session = session
	}

	sessionData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	formData, err := json.Marshal(req.GetRequestForm())
	if err != nil {
		return nil, err
	}

	requestedAt := req.GetRequestedAt().UnixMilli()
	clientID := req.GetClient().GetID()
	requestID := req.GetID()

	tokenSession := &identitysrv.OAuth2TokenSession{
		Signature:       &signature,
		RequestID:       &requestID,
		ClientID:        &clientID,
		Scopes:          optionalJoined(req.GetGrantedScopes()),
		GrantedAudience: optionalJoined(req.GetGrantedAudience()),
		SessionData:     sessionData,
		FormData:        formData,
		RequestedAt:     &requestedAt,
	}

	if userID := extractUserID(req.GetSession()); userID != "" {
		tokenSession.UserID = &userID
	}

	if redirectURI := req.GetRequestForm().Get("redirect_uri"); redirectURI != "" {
		tokenSession.RedirectURI = &redirectURI
	}
	if codeChallenge := req.GetRequestForm().Get("code_challenge"); codeChallenge != "" {
		tokenSession.CodeChallenge = &codeChallenge
	}
	if method := req.GetRequestForm().Get("code_challenge_method"); method != "" {
		tokenSession.CodeChallengeMethod = &method
	}

	if req.GetSession() != nil {
		expiresAt := req.GetSession().GetExpiresAt(fosite.AccessToken)
		if expiresAt.IsZero() {
			expiresAt = req.GetSession().GetExpiresAt(fosite.RefreshToken)
		}
		if expiresAt.IsZero() {
			expiresAt = req.GetSession().GetExpiresAt(fosite.AuthorizeCode)
		}
		if !expiresAt.IsZero() {
			expiresAtMS := expiresAt.UnixMilli()
			tokenSession.ExpiresAt = &expiresAtMS
		}
	}

	return tokenSession, nil
}

func (s *rpcStore) tokenSessionToRequester(ctx context.Context, token *identitysrv.OAuth2TokenSession) (fosite.Requester, error) {
	client, err := s.GetClient(ctx, token.GetClientID())
	if err != nil {
		return nil, err
	}

	request := fosite.NewRequest()
	request.SetID(token.GetRequestID())
	request.Client = client

	if token.GetRequestedAt() > 0 {
		request.RequestedAt = time.UnixMilli(token.GetRequestedAt()).UTC()
	}

	payload := &sessionPayload{}
	if len(token.GetSessionData()) > 0 {
		_ = json.Unmarshal(token.GetSessionData(), payload)
	}

	if payload.Session != nil {
		request.Session = payload.Session
	} else {
		request.Session = &fosite.DefaultSession{}
	}

	if payload.RequestedScopes != nil {
		request.SetRequestedScopes(payload.RequestedScopes)
	} else {
		for _, scope := range strings.Fields(token.GetScopes()) {
			request.AppendRequestedScope(scope)
		}
	}

	if payload.GrantedScopes != nil {
		for _, scope := range payload.GrantedScopes {
			request.GrantScope(scope)
		}
	} else {
		for _, scope := range strings.Fields(token.GetScopes()) {
			request.GrantScope(scope)
		}
	}

	if payload.RequestedAudience != nil {
		request.SetRequestedAudience(payload.RequestedAudience)
	}
	if payload.GrantedAudience != nil {
		for _, aud := range payload.GrantedAudience {
			request.GrantAudience(aud)
		}
	} else {
		for _, aud := range strings.Fields(token.GetGrantedAudience()) {
			request.GrantAudience(aud)
		}
	}

	form := url.Values{}
	if len(token.GetFormData()) > 0 {
		_ = json.Unmarshal(token.GetFormData(), &form)
	}
	request.Form = form

	if token.GetRedirectURI() != "" {
		request.Form.Set("redirect_uri", token.GetRedirectURI())
	}
	if token.GetCodeChallenge() != "" {
		request.Form.Set("code_challenge", token.GetCodeChallenge())
	}
	if token.GetCodeChallengeMethod() != "" {
		request.Form.Set("code_challenge_method", token.GetCodeChallengeMethod())
	}

	if token.GetUserID() != "" {
		if defaultSession, ok := request.Session.(*fosite.DefaultSession); ok {
			defaultSession.Subject = token.GetUserID()
		}
	}

	if token.GetExpiresAt() > 0 {
		exp := time.UnixMilli(token.GetExpiresAt()).UTC()
		request.Session.SetExpiresAt(fosite.AccessToken, exp)
		request.Session.SetExpiresAt(fosite.RefreshToken, exp)
		request.Session.SetExpiresAt(fosite.AuthorizeCode, exp)
	}

	return request, nil
}

func extractUserID(session fosite.Session) string {
	if session == nil {
		return ""
	}

	return session.GetSubject()
}

func optionalJoined(values []string) *string {
	if len(values) == 0 {
		return nil
	}

	joined := strings.Join(values, " ")
	return &joined
}

func toFositeError(err error) error {
	if err == nil {
		return nil
	}

	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "not found") || strings.Contains(msg, "记录") || strings.Contains(msg, "不存在") {
		return fosite.ErrNotFound.WithWrap(err)
	}

	return fosite.ErrServerError.WithWrap(err)
}

var (
	_ fosite.ClientManager          = (*rpcStore)(nil)
	_ oauth2.AuthorizeCodeStorage   = (*rpcStore)(nil)
	_ oauth2.AccessTokenStorage     = (*rpcStore)(nil)
	_ oauth2.RefreshTokenStorage    = (*rpcStore)(nil)
	_ oauth2.TokenRevocationStorage = (*rpcStore)(nil)
	_ pkce.PKCERequestStorage       = (*rpcStore)(nil)
)
