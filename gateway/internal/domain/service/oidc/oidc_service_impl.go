package oidc

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"

	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/oidcstore"
)

// serviceImpl OIDC 领域服务实现
type serviceImpl struct {
	provider *op.Provider
	storage  *oidcstore.Storage
	issuer   string
}

func (s *serviceImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 处理 /login 自动认证：OIDC Provider 重定向到 /login?id=xxx
	// 自动完成用户认证并重定向到 /authorize/callback?id=xxx
	if r.URL.Path == "/login" {
		authReqID := r.URL.Query().Get("id")
		if authReqID == "" {
			http.Error(w, "missing id parameter", http.StatusBadRequest)
			return
		}
		// CheckUsernamePassword 始终成功，设置 IsDone=true 和 UserID
		if err := s.storage.CheckUsernamePassword("oidc-test-user", "", authReqID); err != nil {
			http.Error(w, "authentication failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// 重定向到 OIDC Provider 的回调端点以完成授权
		callbackURL := "/authorize/callback?id=" + authReqID
		http.Redirect(w, r, callbackURL, http.StatusFound)
		return
	}

	s.provider.ServeHTTP(w, r)
}

func (s *serviceImpl) Issuer() string {
	return s.issuer
}

func (s *serviceImpl) Storage() *oidcstore.Storage {
	return s.storage
}

// NewService 创建 OIDC 领域服务实例
func NewService(cfg *config.OIDCConfig, storage op.Storage) (Service, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("OIDC is disabled")
	}

	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	storageImpl, ok := storage.(*oidcstore.Storage)
	if !ok {
		return nil, fmt.Errorf("storage must be *oidcstore.Storage")
	}

	storageImpl.SetSigningKey(rsaKey, "oidc-rsa-key")

	// Generate a random 32-byte crypto key for zitadel internal encryption
	cryptoKey := make([]byte, 32)
	if _, err := rand.Read(cryptoKey); err != nil {
		return nil, fmt.Errorf("failed to generate crypto key: %w", err)
	}

	var keyArray [32]byte
	copy(keyArray[:], cryptoKey)

	zitadelConfig := &op.Config{
		CryptoKey:                keyArray,
		DefaultLogoutRedirectURI: "",
		CodeMethodS256:           cfg.EnforcePKCE,
		AuthMethodPost:           true,
		AuthMethodPrivateKeyJWT:  false,
		GrantTypeRefreshToken:    true,
		RequestObjectSupported:   false,
		SupportedUILocales:       nil,
		SupportedClaims:          []string{"sub", "name", "preferred_username", "email", "email_verified", "picture"},
		SupportedScopes:          []string{"openid", "profile", "email", "offline_access"},
	}

	provider, err := op.NewProvider(
		zitadelConfig,
		storage,
		op.StaticIssuer(cfg.Issuer),
		op.WithAllowInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	return &serviceImpl{
		provider: provider,
		storage:  storageImpl,
		issuer:   cfg.Issuer,
	}, nil
}
