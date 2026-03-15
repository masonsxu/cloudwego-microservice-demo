package oauth2client

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

type oauth2ClientRepositoryImpl struct {
	db *gorm.DB
}

func NewOAuth2ClientRepository(db *gorm.DB) OAuth2ClientRepository {
	return &oauth2ClientRepositoryImpl{db: db}
}

func (r *oauth2ClientRepositoryImpl) Create(ctx context.Context, client *models.OAuth2Client) error {
	return r.db.WithContext(ctx).Create(client).Error
}

func (r *oauth2ClientRepositoryImpl) GetByID(ctx context.Context, id string) (*models.OAuth2Client, error) {
	var client models.OAuth2Client
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&client).Error; err != nil {
		return nil, fmt.Errorf("查询 OAuth2 客户端失败: %w", err)
	}

	return &client, nil
}

func (r *oauth2ClientRepositoryImpl) GetByClientID(ctx context.Context, clientID string) (*models.OAuth2Client, error) {
	var client models.OAuth2Client
	if err := r.db.WithContext(ctx).Where("client_id = ?", clientID).First(&client).Error; err != nil {
		return nil, fmt.Errorf("查询 OAuth2 客户端失败: %w", err)
	}

	return &client, nil
}

func (r *oauth2ClientRepositoryImpl) Update(ctx context.Context, client *models.OAuth2Client) error {
	return r.db.WithContext(ctx).Save(client).Error
}

func (r *oauth2ClientRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.OAuth2Client{}).Error
}

func (r *oauth2ClientRepositoryImpl) List(
	ctx context.Context,
	ownerID *string,
	isActive *bool,
	page, limit int,
) ([]*models.OAuth2Client, int64, error) {
	var clients []*models.OAuth2Client

	var total int64

	query := r.db.WithContext(ctx).Model(&models.OAuth2Client{})

	if ownerID != nil {
		query = query.Where("owner_id = ?", *ownerID)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计 OAuth2 客户端数量失败: %w", err)
	}

	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&clients).Error; err != nil {
		return nil, 0, fmt.Errorf("查询 OAuth2 客户端列表失败: %w", err)
	}

	return clients, total, nil
}
