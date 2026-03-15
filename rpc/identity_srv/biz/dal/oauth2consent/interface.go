package oauth2consent

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// OAuth2ConsentRepository OAuth2 用户授权同意数据访问接口
type OAuth2ConsentRepository interface {
	Save(ctx context.Context, consent *models.OAuth2Consent) error
	GetByUserAndClient(ctx context.Context, userID, clientID string) (*models.OAuth2Consent, error)
	ListByUserID(ctx context.Context, userID string, page, limit int) ([]*models.OAuth2Consent, int64, error)
	Revoke(ctx context.Context, userID, clientID string) error
}
