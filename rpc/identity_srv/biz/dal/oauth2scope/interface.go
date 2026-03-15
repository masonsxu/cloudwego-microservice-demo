package oauth2scope

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// OAuth2ScopeRepository OAuth2 作用域数据访问接口
type OAuth2ScopeRepository interface {
	ListAll(ctx context.Context) ([]*models.OAuth2Scope, error)
	ListDefaults(ctx context.Context) ([]*models.OAuth2Scope, error)
}
