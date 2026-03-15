package oauth2consent

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

type oauth2ConsentRepositoryImpl struct {
	db *gorm.DB
}

func NewOAuth2ConsentRepository(db *gorm.DB) OAuth2ConsentRepository {
	return &oauth2ConsentRepositoryImpl{db: db}
}

func (r *oauth2ConsentRepositoryImpl) Save(ctx context.Context, consent *models.OAuth2Consent) error {
	consent.GrantedAt = time.Now().UnixMilli()

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "client_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"scopes", "granted_at", "is_revoked", "updated_at"}),
		}).
		Create(consent).Error
}

func (r *oauth2ConsentRepositoryImpl) GetByUserAndClient(
	ctx context.Context,
	userID, clientID string,
) (*models.OAuth2Consent, error) {
	var consent models.OAuth2Consent
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND client_id = ? AND is_revoked = false", userID, clientID).
		First(&consent).Error; err != nil {
		return nil, fmt.Errorf("查询授权同意记录失败: %w", err)
	}

	return &consent, nil
}

func (r *oauth2ConsentRepositoryImpl) ListByUserID(
	ctx context.Context,
	userID string,
	page, limit int,
) ([]*models.OAuth2Consent, int64, error) {
	var consents []*models.OAuth2Consent

	var total int64

	query := r.db.WithContext(ctx).Model(&models.OAuth2Consent{}).
		Where("user_id = ? AND is_revoked = false", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计授权同意记录失败: %w", err)
	}

	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Order("granted_at DESC").Find(&consents).Error; err != nil {
		return nil, 0, fmt.Errorf("查询授权同意列表失败: %w", err)
	}

	return consents, total, nil
}

func (r *oauth2ConsentRepositoryImpl) Revoke(ctx context.Context, userID, clientID string) error {
	return r.db.WithContext(ctx).Model(&models.OAuth2Consent{}).
		Where("user_id = ? AND client_id = ?", userID, clientID).
		Update("is_revoked", true).Error
}
