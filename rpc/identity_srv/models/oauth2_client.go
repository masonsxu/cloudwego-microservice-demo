package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// OAuth2ClientType OAuth2 客户端类型
type OAuth2ClientType string

const (
	OAuth2ClientTypeConfidential OAuth2ClientType = "confidential"
	OAuth2ClientTypePublic       OAuth2ClientType = "public"
)

// OAuth2Client OAuth2 客户端（应用）模型
type OAuth2Client struct {
	BaseModel

	// 客户端标识
	ClientID     string           `gorm:"column:client_id;uniqueIndex;size:64;not null;comment:客户端标识符"`
	ClientSecret string           `gorm:"column:client_secret;size:256;not null;comment:客户端密钥(bcrypt哈希)"`
	ClientName   string           `gorm:"column:client_name;size:128;not null;comment:客户端显示名称"`
	Description  string           `gorm:"column:description;size:512;comment:客户端描述"`
	ClientType   OAuth2ClientType `gorm:"column:client_type;size:20;not null;default:'confidential';comment:客户端类型"`

	// 授权配置
	GrantTypes   StringSlice `gorm:"column:grant_types;type:jsonb;not null;comment:允许的授权类型"`
	RedirectURIs StringSlice `gorm:"column:redirect_uris;type:jsonb;comment:允许的回调地址"`
	Scopes       StringSlice `gorm:"column:scopes;type:jsonb;comment:允许的作用域"`

	// 展示信息
	LogoURI   string `gorm:"column:logo_uri;size:512;comment:客户端Logo URL"`
	ClientURI string `gorm:"column:client_uri;size:512;comment:客户端主页URL"`

	// Token 配置
	AccessTokenLifespan  int `gorm:"column:access_token_lifespan;default:3600;comment:Access Token有效期(秒)"`
	RefreshTokenLifespan int `gorm:"column:refresh_token_lifespan;default:2592000;comment:Refresh Token有效期(秒)"`

	// 归属和状态
	OwnerID  *uuid.UUID `gorm:"column:owner_id;type:uuid;comment:创建者用户ID"`
	IsActive bool       `gorm:"column:is_active;default:true;comment:是否启用"`
}

func (OAuth2Client) TableName() string {
	return "oauth2_clients"
}

func (c *OAuth2Client) BeforeCreate(tx *gorm.DB) error {
	if c.ClientID == "" {
		c.ClientID = generateClientID()
	}

	if c.ClientType == "" {
		c.ClientType = OAuth2ClientTypeConfidential
	}

	return c.validate()
}

func (c *OAuth2Client) BeforeUpdate(tx *gorm.DB) error {
	return c.validate()
}

func (c *OAuth2Client) validate() error {
	if c.ClientName == "" {
		return fmt.Errorf("客户端名称不能为空")
	}

	if len(c.ClientName) > 128 {
		return fmt.Errorf("客户端名称不能超过128字符")
	}

	if c.ClientType != OAuth2ClientTypeConfidential && c.ClientType != OAuth2ClientTypePublic {
		return fmt.Errorf("客户端类型必须为 confidential 或 public")
	}

	return nil
}

// GenerateSecret 生成新的客户端密钥，返回明文（仅返回一次）
func (c *OAuth2Client) GenerateSecret() (string, error) {
	secret := generateRandomString(48)

	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("生成密钥哈希失败: %w", err)
	}

	c.ClientSecret = string(hash)

	return secret, nil
}

// VerifySecret 验证客户端密钥
func (c *OAuth2Client) VerifySecret(secret string) bool {
	return bcrypt.CompareHashAndPassword([]byte(c.ClientSecret), []byte(secret)) == nil
}

// generateClientID 生成唯一的 client_id（UUID 格式）
func generateClientID() string {
	return uuid.New().String()
}

// generateRandomString 生成指定长度的 URL 安全随机字符串
func generateRandomString(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)

	return base64.RawURLEncoding.EncodeToString(b)[:length]
}
