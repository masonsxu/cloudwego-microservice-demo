package models

// OAuth2PKCESession PKCE 会话存储模型
// 存储 PKCE (Proof Key for Code Exchange) 的验证信息。
type OAuth2PKCESession struct {
	BaseModel

	// 签名（与授权码关联）
	Signature string `gorm:"column:signature;uniqueIndex;size:128;not null;comment:签名(与授权码关联)"`

	// 关联的请求ID
	RequestID string `gorm:"column:request_id;size:128;not null;index;comment:关联请求ID"`

	// 序列化数据
	SessionData []byte `gorm:"column:session_data;type:jsonb;comment:fosite Session 序列化数据"`
	FormData    []byte `gorm:"column:form_data;type:jsonb;comment:请求表单序列化数据"`

	// 生命周期
	RequestedAt int64 `gorm:"column:requested_at;not null;comment:请求时间(毫秒)"`
	ExpiresAt   int64 `gorm:"column:expires_at;not null;index;comment:过期时间(毫秒)"`
}

func (OAuth2PKCESession) TableName() string {
	return "oauth2_pkce_sessions"
}
