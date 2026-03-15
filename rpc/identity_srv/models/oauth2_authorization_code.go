package models

import "github.com/google/uuid"

// OAuth2AuthorizationCode 授权码模型
// 存储 Authorization Code Grant 流程中的授权码会话信息。
type OAuth2AuthorizationCode struct {
	BaseModel

	// 授权码（SHA256 哈希存储）
	Signature string `gorm:"column:signature;uniqueIndex;size:128;not null;comment:授权码签名(哈希)"`

	// 关联的请求ID（fosite 内部使用）
	RequestID string `gorm:"column:request_id;size:128;not null;index;comment:关联请求ID"`

	// 关联方
	ClientID string    `gorm:"column:client_id;size:64;not null;index;comment:客户端标识符"`
	UserID   uuid.UUID `gorm:"column:user_id;type:uuid;not null;index;comment:授权用户ID"`

	// 授权信息
	Scopes          string `gorm:"column:scopes;size:1024;comment:授权的作用域(空格分隔)"`
	GrantedAudience string `gorm:"column:granted_audience;size:1024;comment:授权的受众(空格分隔)"`
	RedirectURI     string `gorm:"column:redirect_uri;size:512;not null;comment:回调地址"`

	// PKCE
	CodeChallenge       string `gorm:"column:code_challenge;size:256;comment:PKCE Code Challenge"`
	CodeChallengeMethod string `gorm:"column:code_challenge_method;size:10;comment:PKCE 方法(S256)"`

	// 序列化数据
	SessionData []byte `gorm:"column:session_data;type:jsonb;comment:fosite Session 序列化数据"`
	FormData    []byte `gorm:"column:form_data;type:jsonb;comment:请求表单序列化数据"`

	// 生命周期
	RequestedAt int64 `gorm:"column:requested_at;not null;comment:请求时间(毫秒)"`
	ExpiresAt   int64 `gorm:"column:expires_at;not null;index;comment:过期时间(毫秒)"`
	Used        bool  `gorm:"column:used;default:false;comment:是否已使用(一次性)"`
}

func (OAuth2AuthorizationCode) TableName() string {
	return "oauth2_authorization_codes"
}
