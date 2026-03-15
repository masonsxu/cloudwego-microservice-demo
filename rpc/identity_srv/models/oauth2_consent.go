package models

import "github.com/google/uuid"

// OAuth2Consent 用户授权同意记录模型
// 记录用户对 OAuth2 客户端的授权同意状态，避免重复弹出授权确认页。
type OAuth2Consent struct {
	BaseModel

	// 关联方（联合唯一索引）
	UserID   uuid.UUID `gorm:"column:user_id;type:uuid;not null;uniqueIndex:idx_consent_user_client;comment:授权用户ID"`
	ClientID string    `gorm:"column:client_id;size:64;not null;uniqueIndex:idx_consent_user_client;comment:客户端标识符"`

	// 授权信息
	Scopes    string `gorm:"column:scopes;size:1024;not null;comment:授权的作用域(空格分隔)"`
	GrantedAt int64  `gorm:"column:granted_at;not null;comment:授权时间(毫秒)"`
	ExpiresAt *int64 `gorm:"column:expires_at;comment:授权过期时间(毫秒,可选)"`

	// 状态
	IsRevoked bool `gorm:"column:is_revoked;default:false;comment:是否已撤销"`
}

func (OAuth2Consent) TableName() string {
	return "oauth2_consents"
}
