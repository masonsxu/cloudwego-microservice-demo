package oauth2

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/token/jwt"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
)

// NewOAuth2Provider 创建并配置 fosite OAuth2Provider 实例。
// store 参数需要实现 fosite 的多个存储接口。
func NewOAuth2Provider(cfg config.OAuth2Config, store interface{}) fosite.OAuth2Provider {
	secret := []byte("fosite-hmac-secret-change-me-in-production") // TODO: 从配置读取

	rsaKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	fositeConfig := &fosite.Config{
		AccessTokenLifespan:         cfg.AccessTokenLifespan,
		RefreshTokenLifespan:        cfg.RefreshTokenLifespan,
		AuthorizeCodeLifespan:       cfg.AuthCodeLifespan,
		TokenURL:                    cfg.Issuer + "/oauth2/token",
		AccessTokenIssuer:           cfg.Issuer,
		ScopeStrategy:               fosite.WildcardScopeStrategy,
		AudienceMatchingStrategy:    fosite.DefaultAudienceMatchingStrategy,
		EnforcePKCE:                 cfg.EnforcePKCE,
		EnforcePKCEForPublicClients: true,
		GlobalSecret:                secret,
		RefreshTokenScopes:          []string{"offline"},
		MinParameterEntropy:         fosite.MinParameterEntropy,
	}

	keyGetter := func(_ context.Context) (interface{}, error) {
		return rsaKey, nil
	}

	hmacStrategy := compose.NewOAuth2HMACStrategy(fositeConfig)

	return compose.Compose(
		fositeConfig,
		store,
		&compose.CommonStrategy{
			CoreStrategy: hmacStrategy,
			Signer: &jwt.DefaultSigner{
				GetPrivateKey: keyGetter,
			},
		},
		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2RefreshTokenGrantFactory,
		compose.OAuth2PKCEFactory,
	)
}

// NewDefaultSession 创建一个默认的 fosite Session。
func NewDefaultSession(subject string) *fosite.DefaultSession {
	return &fosite.DefaultSession{
		Subject: subject,
		ExpiresAt: map[fosite.TokenType]time.Time{
			fosite.AccessToken:  time.Now().Add(1 * time.Hour),
			fosite.RefreshToken: time.Now().Add(30 * 24 * time.Hour),
		},
	}
}
