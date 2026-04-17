package oidcstore

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
)

const (
	keyPrefixAuthReq    = "oidc:auth_req:"
	keyPrefixAuthCode   = "oidc:auth_code:"
	keyPrefixAccessTok  = "oidc:access_token:"
	keyPrefixRefreshTok = "oidc:refresh_token:"
	keyPrefixClient     = "oidc:client:"
)

// Storage 实现 op.Storage 接口
type Storage struct {
	rdb            *redis.Client
	oidcConfig     *config.OIDCConfig
	identityClient identitycli.IdentityClient
	signingKey     *rsa.PrivateKey
	keyID          string
}

func NewStorage(rdb *redis.Client, cfg *config.OIDCConfig, identityClient identitycli.IdentityClient) *Storage {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	return &Storage{
		rdb:            rdb,
		oidcConfig:     cfg,
		identityClient: identityClient,
		signingKey:     key,
		keyID:          "oidc-rsa-key",
	}
}

// SetSigningKey 设置用于签名 OIDC Token 的 RSA 私钥
func (s *Storage) SetSigningKey(key *rsa.PrivateKey, keyID string) {
	s.signingKey = key
	s.keyID = keyID
}

// ===========================================================================
// OPStorage
// ===========================================================================

func (s *Storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	data, err := s.rdb.Get(ctx, keyPrefixClient+clientID).Bytes()
	if err == redis.Nil {
		return &defaultClient{ID: clientID}, nil
	}

	if err != nil {
		return nil, err
	}

	var c storedClient
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (s *Storage) AuthorizeClientIDSecret(ctx context.Context, clientID, clientSecret string) error {
	return nil
}

func (s *Storage) SetUserinfoFromScopes(
	ctx context.Context,
	userinfo *oidc.UserInfo,
	userID, clientID string,
	scopes []string,
) error {
	userinfo.Subject = userID
	userinfo.PreferredUsername = userID
	userinfo.Name = userID
	userinfo.Email = userID + "@example.com"
	userinfo.EmailVerified = true

	return nil
}

func (s *Storage) SetUserinfoFromToken(
	ctx context.Context,
	userinfo *oidc.UserInfo,
	tokenID, subject, origin string,
) error {
	userinfo.Subject = subject
	userinfo.PreferredUsername = subject

	return nil
}

func (s *Storage) SetIntrospectionFromToken(
	ctx context.Context,
	introspection *oidc.IntrospectionResponse,
	tokenID, subject, clientID string,
) error {
	introspection.Active = true
	introspection.Subject = subject
	introspection.ClientID = clientID
	introspection.Username = subject

	return nil
}

func (s *Storage) GetPrivateClaimsFromScopes(
	ctx context.Context,
	userID, clientID string,
	scopes []string,
) (map[string]any, error) {
	return map[string]any{}, nil
}

func (s *Storage) GetKeyByIDAndClientID(ctx context.Context, keyID, clientID string) (*jose.JSONWebKey, error) {
	return nil, nil
}

func (s *Storage) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	return scopes, nil
}

func (s *Storage) CheckUsernamePassword(username, password, id string) error {
	ar, err := s.AuthRequestByID(context.Background(), id)
	if err != nil {
		return err
	}

	authReq := ar.(*authRequest)
	authReq.IsDone = true
	authReq.UserID = username
	data, _ := json.Marshal(authReq)
	s.rdb.Set(context.Background(), keyPrefixAuthReq+id, data, s.oidcConfig.AuthCodeLifespan)

	return nil
}

// ===========================================================================
// AuthStorage
// ===========================================================================

func (s *Storage) CreateAuthRequest(
	ctx context.Context,
	authReq *oidc.AuthRequest,
	userID string,
) (op.AuthRequest, error) {
	ar := &authRequest{
		ID:            uuid.NewString(),
		CreationDate:  time.Now(),
		ClientID:      authReq.ClientID,
		CallbackURI:   authReq.RedirectURI,
		TransferState: authReq.State,
		Prompt:        []string(authReq.Prompt),
		LoginHint:     authReq.LoginHint,
		UserID:        userID,
		Scopes:        []string(authReq.Scopes),
		ResponseType:  string(authReq.ResponseType),
		Nonce:         authReq.Nonce,
		IsDone:        userID != "",
	}
	if authReq.CodeChallenge != "" {
		ar.CodeChallenge = &codeChallenge{
			Challenge: authReq.CodeChallenge,
			Method:    string(authReq.CodeChallengeMethod),
		}
	}

	data, _ := json.Marshal(ar)

	ttl := s.oidcConfig.AuthCodeLifespan
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}

	if err := s.rdb.Set(ctx, keyPrefixAuthReq+ar.ID, data, ttl).Err(); err != nil {
		return nil, err
	}

	return ar, nil
}

func (s *Storage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	data, err := s.rdb.Get(ctx, keyPrefixAuthReq+id).Bytes()
	if err == redis.Nil {
		return nil, oidc.ErrInvalidRequest()
	}

	if err != nil {
		return nil, err
	}

	var ar authRequest
	if err := json.Unmarshal(data, &ar); err != nil {
		return nil, err
	}

	return &ar, nil
}

func (s *Storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	id, err := s.rdb.Get(ctx, keyPrefixAuthCode+code).Result()
	if err == redis.Nil {
		return nil, oidc.ErrInvalidRequest()
	}

	if err != nil {
		return nil, err
	}

	return s.AuthRequestByID(ctx, id)
}

func (s *Storage) SaveAuthCode(ctx context.Context, id string, code string) error {
	return s.rdb.Set(ctx, keyPrefixAuthCode+code, id, s.oidcConfig.AuthCodeLifespan).Err()
}

func (s *Storage) DeleteAuthRequest(ctx context.Context, id string) error {
	return s.rdb.Del(ctx, keyPrefixAuthReq+id).Err()
}

func (s *Storage) CreateAccessToken(ctx context.Context, request op.TokenRequest) (string, time.Time, error) {
	tokenID := uuid.NewString()
	expiration := time.Now().Add(s.oidcConfig.AccessTokenLifespan)

	return tokenID, expiration, nil
}

func (s *Storage) CreateAccessAndRefreshTokens(
	ctx context.Context,
	request op.TokenRequest,
	currentRefreshToken string,
) (string, string, time.Time, error) {
	accessTokenID := uuid.NewString()
	refreshToken := uuid.NewString()
	expiration := time.Now().Add(s.oidcConfig.AccessTokenLifespan)

	return accessTokenID, refreshToken, expiration, nil
}

func (s *Storage) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (op.RefreshTokenRequest, error) {
	data, err := s.rdb.Get(ctx, keyPrefixRefreshTok+refreshToken).Bytes()
	if err == redis.Nil {
		return nil, oidc.ErrInvalidRequest()
	}

	if err != nil {
		return nil, err
	}

	var rrt refreshTokenRequest
	if err := json.Unmarshal(data, &rrt); err != nil {
		return nil, err
	}

	return &rrt, nil
}

func (s *Storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	return nil
}

func (s *Storage) RevokeToken(ctx context.Context, tokenIDOrToken string, userID string, clientID string) *oidc.Error {
	s.rdb.Del(ctx, keyPrefixAccessTok+tokenIDOrToken)
	s.rdb.Del(ctx, keyPrefixRefreshTok+tokenIDOrToken)

	return nil
}

func (s *Storage) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (string, string, error) {
	return "", "", oidc.ErrInvalidRequest()
}

// ===========================================================================
// Health
// ===========================================================================

func (s *Storage) Health(ctx context.Context) error {
	return s.rdb.Ping(ctx).Err()
}

// ===========================================================================
// Key 管理
// ===========================================================================

func (s *Storage) SigningKey(ctx context.Context) (op.SigningKey, error) {
	if s.signingKey == nil {
		return nil, oidc.ErrServerError()
	}

	return &signingKey{KeyData: s.signingKey, KeyID: s.keyID}, nil
}

func (s *Storage) SignatureAlgorithms(ctx context.Context) ([]jose.SignatureAlgorithm, error) {
	return []jose.SignatureAlgorithm{jose.RS256}, nil
}

func (s *Storage) KeySet(ctx context.Context) ([]op.Key, error) {
	if s.signingKey == nil {
		return nil, oidc.ErrServerError()
	}

	return []op.Key{&keyImpl{PubKey: &s.signingKey.PublicKey, IDVal: s.keyID}}, nil
}

// ===========================================================================
// 内部类型
// ===========================================================================

type authRequest struct {
	ID            string         `json:"id"`
	CreationDate  time.Time      `json:"creation_date"`
	ClientID      string         `json:"client_id"`
	CallbackURI   string         `json:"callback_uri"`
	TransferState string         `json:"transfer_state"`
	Prompt        []string       `json:"prompt"`
	LoginHint     string         `json:"login_hint"`
	UserID        string         `json:"user_id"`
	Scopes        []string       `json:"scopes"`
	ResponseType  string         `json:"response_type"`
	Nonce         string         `json:"nonce"`
	CodeChallenge *codeChallenge `json:"code_challenge"`
	IsDone        bool           `json:"is_done"`
}

type codeChallenge struct {
	Challenge string `json:"challenge"`
	Method    string `json:"method"`
}

func (a *authRequest) GetID() string          { return a.ID }
func (a *authRequest) GetACR() string         { return "" }
func (a *authRequest) GetAMR() []string       { return []string{"pwd"} }
func (a *authRequest) GetAudience() []string  { return []string{a.ClientID} }
func (a *authRequest) GetAuthTime() time.Time { return a.CreationDate }
func (a *authRequest) GetClientID() string    { return a.ClientID }
func (a *authRequest) GetCodeChallenge() *oidc.CodeChallenge {
	if a.CodeChallenge == nil {
		return nil
	}

	return &oidc.CodeChallenge{
		Challenge: a.CodeChallenge.Challenge,
		Method:    oidc.CodeChallengeMethod(a.CodeChallenge.Method),
	}
}
func (a *authRequest) GetNonce() string                   { return a.Nonce }
func (a *authRequest) GetRedirectURI() string             { return a.CallbackURI }
func (a *authRequest) GetResponseType() oidc.ResponseType { return oidc.ResponseType(a.ResponseType) }
func (a *authRequest) GetResponseMode() oidc.ResponseMode { return oidc.ResponseModeQuery }
func (a *authRequest) GetScopes() []string                { return a.Scopes }
func (a *authRequest) GetState() string                   { return a.TransferState }
func (a *authRequest) GetSubject() string                 { return a.UserID }
func (a *authRequest) Done() bool                         { return a.IsDone }

type storedClient struct {
	IDVal           string   `json:"id"`
	SecretVal       string   `json:"secret"`
	RedirectURIsVal []string `json:"redirect_uris"`
}

func (c *storedClient) GetID() string                       { return c.IDVal }
func (c *storedClient) RedirectURIs() []string              { return c.RedirectURIsVal }
func (c *storedClient) PostLogoutRedirectURIs() []string    { return []string{} }
func (c *storedClient) ApplicationType() op.ApplicationType { return op.ApplicationTypeWeb }
func (c *storedClient) AuthMethod() oidc.AuthMethod         { return oidc.AuthMethodBasic }
func (c *storedClient) ResponseTypes() []oidc.ResponseType {
	return []oidc.ResponseType{oidc.ResponseTypeCode}
}

func (c *storedClient) GrantTypes() []oidc.GrantType {
	return []oidc.GrantType{oidc.GrantTypeCode, oidc.GrantTypeRefreshToken}
}
func (c *storedClient) LoginURL(authReqID string) string     { return "/login?id=" + authReqID }
func (c *storedClient) AccessTokenType() op.AccessTokenType  { return op.AccessTokenTypeBearer }
func (c *storedClient) IDTokenLifetime() time.Duration       { return 1 * time.Hour }
func (c *storedClient) DevMode() bool                        { return false }
func (c *storedClient) IsScopeAllowed(string) bool           { return true }
func (c *storedClient) IDTokenUserinfoClaimsAssertion() bool { return false }
func (c *storedClient) ClockSkew() time.Duration             { return 0 }
func (c *storedClient) RestrictAdditionalIdTokenScopes() func([]string) []string {
	return func(s []string) []string { return s }
}

func (c *storedClient) RestrictAdditionalAccessTokenScopes() func([]string) []string {
	return func(s []string) []string { return s }
}

type defaultClient struct {
	ID string
}

func (c *defaultClient) GetID() string { return c.ID }
func (c *defaultClient) RedirectURIs() []string {
	return []string{
		"http://localhost",
		"http://localhost:5173",
		"http://localhost:5173/oidc/callback",
	}
}
func (c *defaultClient) PostLogoutRedirectURIs() []string    { return []string{} }
func (c *defaultClient) ApplicationType() op.ApplicationType { return op.ApplicationTypeNative }
func (c *defaultClient) AuthMethod() oidc.AuthMethod         { return oidc.AuthMethodNone }
func (c *defaultClient) ResponseTypes() []oidc.ResponseType {
	return []oidc.ResponseType{oidc.ResponseTypeCode}
}

func (c *defaultClient) GrantTypes() []oidc.GrantType {
	return []oidc.GrantType{oidc.GrantTypeCode, oidc.GrantTypeRefreshToken}
}
func (c *defaultClient) LoginURL(authReqID string) string     { return "/login?id=" + authReqID }
func (c *defaultClient) AccessTokenType() op.AccessTokenType  { return op.AccessTokenTypeBearer }
func (c *defaultClient) IDTokenLifetime() time.Duration       { return 1 * time.Hour }
func (c *defaultClient) DevMode() bool                        { return true }
func (c *defaultClient) IsScopeAllowed(string) bool           { return true }
func (c *defaultClient) IDTokenUserinfoClaimsAssertion() bool { return false }
func (c *defaultClient) ClockSkew() time.Duration             { return 0 }
func (c *defaultClient) RestrictAdditionalIdTokenScopes() func([]string) []string {
	return func(s []string) []string { return s }
}

func (c *defaultClient) RestrictAdditionalAccessTokenScopes() func([]string) []string {
	return func(s []string) []string { return s }
}

type refreshTokenRequest struct {
	SubjectVal string   `json:"subject"`
	ScopesVal  []string `json:"scopes"`
}

func (r *refreshTokenRequest) GetSubject() string               { return r.SubjectVal }
func (r *refreshTokenRequest) GetScopes() []string              { return r.ScopesVal }
func (r *refreshTokenRequest) GetAudience() []string            { return nil }
func (r *refreshTokenRequest) GetAuthTime() time.Time           { return time.Now() }
func (r *refreshTokenRequest) GetAMR() []string                 { return []string{"pwd"} }
func (r *refreshTokenRequest) GetClientID() string              { return "" }
func (r *refreshTokenRequest) SetCurrentScopes(scopes []string) {}

type signingKey struct {
	KeyData *rsa.PrivateKey
	KeyID   string
}

func (s *signingKey) SignatureAlgorithm() jose.SignatureAlgorithm { return jose.RS256 }
func (s *signingKey) Key() interface{}                            { return s.KeyData }
func (s *signingKey) ID() string                                  { return s.KeyID }

type keyImpl struct {
	PubKey *rsa.PublicKey
	IDVal  string
}

func (k *keyImpl) ID() string                         { return k.IDVal }
func (k *keyImpl) Algorithm() jose.SignatureAlgorithm { return jose.RS256 }
func (k *keyImpl) Use() string                        { return "sig" }
func (k *keyImpl) Key() interface{}                   { return k.PubKey }
