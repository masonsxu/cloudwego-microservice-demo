package oauth2client

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// OAuth2ClientRepository OAuth2 客户端数据访问接口
type OAuth2ClientRepository interface {
	Create(ctx context.Context, client *models.OAuth2Client) error
	GetByID(ctx context.Context, id string) (*models.OAuth2Client, error)
	GetByClientID(ctx context.Context, clientID string) (*models.OAuth2Client, error)
	Update(ctx context.Context, client *models.OAuth2Client) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, ownerID *string, isActive *bool, page, limit int) ([]*models.OAuth2Client, int64, error)
}
