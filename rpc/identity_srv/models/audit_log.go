package models

import (
	"github.com/google/uuid"
)

// AuditAction 审计操作类型
type AuditAction int8

const (
	AuditActionCreate         AuditAction = 1
	AuditActionUpdate         AuditAction = 2
	AuditActionDelete         AuditAction = 3
	AuditActionLogin          AuditAction = 4
	AuditActionLogout         AuditAction = 5
	AuditActionPasswordChange AuditAction = 6
)

// AuditLog 审计日志模型
// 审计日志是只追加的不可变记录，不使用 BaseModel（无 UpdatedAt、DeletedAt）
type AuditLog struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RequestID      string    `gorm:"index"`
	TraceID        string
	UserID         *uuid.UUID `gorm:"index"`
	Username       string
	OrganizationID *uuid.UUID  `gorm:"index"`
	Action         AuditAction `gorm:"index"`
	Resource       string      `gorm:"index"`
	ResourceID     string
	StatusCode     int32
	Success        bool `gorm:"index"`
	ClientIP       string
	UserAgent      string
	RequestBody    string `gorm:"type:text"`
	DurationMs     int32
	CreatedAt      int64 `gorm:"autoCreateTime:milli;index"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}
