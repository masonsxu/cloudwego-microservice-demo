package models

import "github.com/google/uuid"

// OAuth2AccessToken 访问令牌记录模型
// 存储 fosite 签发的 Access Token 会话信息。
type OAuth2AccessToken struct {
	BaseModel

	// 令牌签名（SHA256 哈希存储）
	Signature string `gorm:"column:signature;uniqueIndex;size:128;not null;comment:令牌签名(哈希)"`

	// 关联的请求ID
	RequestID string `gorm:"column:request_id;size:128;not null;index;comment:关联请求ID"`

	// 关联方
	ClientID string     `gorm:"column:client_id;size:64;not null;index;comment:客户端标识符"`
	UserID   *uuid.UUID `gorm:"column:user_id;type:uuid;index;comment:用户ID(M2M模式可为空)"`

	// 授权信息
	Scopes          string `gorm:"column:scopes;size:1024;comment:授权的作用域(空格分隔)"`
	GrantedAudience string `gorm:"column:granted_audience;size:1024;comment:授权的受众(空格分隔)"`

	// 序列化数据
	SessionData []byte `gorm:"column:session_data;type:jsonb;comment:fosite Session 序列化数据"`
	FormData    []byte `gorm:"column:form_data;type:jsonb;comment:请求表单序列化数据"`

	// 生命周期
	RequestedAt int64 `gorm:"column:requested_at;not null;comment:请求时间(毫秒)"`
	ExpiresAt   int64 `gorm:"column:expires_at;not null;index;comment:过期时间(毫秒)"`
	IsRevoked   bool  `gorm:"column:is_revoked;default:false;index;comment:是否已吊销"`
}

func (OAuth2AccessToken) TableName() string {
	return "oauth2_access_tokens"
}
