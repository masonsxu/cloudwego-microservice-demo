package oauth2scope

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

type oauth2ScopeRepositoryImpl struct {
	db *gorm.DB
}

func NewOAuth2ScopeRepository(db *gorm.DB) OAuth2ScopeRepository {
	return &oauth2ScopeRepositoryImpl{db: db}
}

func (r *oauth2ScopeRepositoryImpl) ListAll(ctx context.Context) ([]*models.OAuth2Scope, error) {
	var scopes []*models.OAuth2Scope
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&scopes).Error; err != nil {
		return nil, fmt.Errorf("查询作用域列表失败: %w", err)
	}

	return scopes, nil
}

func (r *oauth2ScopeRepositoryImpl) ListDefaults(ctx context.Context) ([]*models.OAuth2Scope, error) {
	var scopes []*models.OAuth2Scope
	if err := r.db.WithContext(ctx).Where("is_default = true").Order("name ASC").Find(&scopes).Error; err != nil {
		return nil, fmt.Errorf("查询默认作用域列表失败: %w", err)
	}

	return scopes, nil
}
