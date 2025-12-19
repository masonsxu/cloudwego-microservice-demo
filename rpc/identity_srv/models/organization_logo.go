package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrganizationLogoStatus 组织Logo状态枚举
type OrganizationLogoStatus int8

const (
	// LogoStatusTemporary 临时状态 - 上传后7天内未绑定组织将自动删除
	LogoStatusTemporary OrganizationLogoStatus = 0
	// LogoStatusBound 已绑定状态 - 已绑定到组织，永久保存
	LogoStatusBound OrganizationLogoStatus = 1
	// LogoStatusDeleted 已删除状态 - 软删除标记
	LogoStatusDeleted OrganizationLogoStatus = 2
)

// String 实现 Stringer 接口
func (s OrganizationLogoStatus) String() string {
	switch s {
	case LogoStatusTemporary:
		return "TEMPORARY"
	case LogoStatusBound:
		return "BOUND"
	case LogoStatusDeleted:
		return "DELETED"
	default:
		return "UNKNOWN"
	}
}

// IsValid 验证状态是否有效
func (s OrganizationLogoStatus) IsValid() bool {
	return s >= LogoStatusTemporary && s <= LogoStatusDeleted
}

// OrganizationLogo 组织Logo模型
// 用于存储组织Logo文件的元数据信息
type OrganizationLogo struct {
	BaseModel

	// Logo状态
	Status OrganizationLogoStatus `gorm:"column:status;type:smallint;not null;default:0;index;comment:Logo状态"`

	// 组织绑定关系
	BoundOrganizationID *uuid.UUID `gorm:"column:bound_organization_id;type:uuid;index;comment:绑定的组织ID（临时状态时为NULL）"`

	// 文件存储信息
	FileID string `gorm:"column:file_id;size:500;not null;uniqueIndex;comment:S3存储路径: organization-logos/{uuid}.{ext}"`

	// 文件属性
	FileName string `gorm:"column:file_name;size:255;not null;comment:原始文件名"`
	FileSize int64  `gorm:"column:file_size;not null;comment:文件大小（字节）"`
	MimeType string `gorm:"column:mime_type;size:100;not null;comment:MIME类型（image/png, image/jpeg等）"`

	// 生命周期管理
	ExpiresAt *int64 `gorm:"column:expires_at;comment:过期时间（毫秒时间戳，临时状态必填）"`

	// 审计信息
	UploadedBy uuid.UUID `gorm:"column:uploaded_by;type:uuid;not null;index;comment:上传者用户ID"`
}

// TableName 指定表名
func (OrganizationLogo) TableName() string {
	return "organization_logos"
}

// BeforeCreate GORM 钩子 - 创建前验证
func (o *OrganizationLogo) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}

	return nil
}

// BeforeUpdate GORM 钩子 - 更新前验证
func (o *OrganizationLogo) BeforeUpdate(tx *gorm.DB) error {
	// 检查哪些字段正在被更新
	stmt := tx.Statement

	// 只在实际更新 Status 字段时验证状态
	if stmt.Changed("Status") {
		if !o.Status.IsValid() {
			return ErrLogoInvalidStatus
		}

		// 验证绑定状态必须设置组织ID
		if o.Status == LogoStatusBound {
			if o.BoundOrganizationID == nil || *o.BoundOrganizationID == uuid.Nil {
				return ErrLogoMissingOrganization
			}
		}
	}

	// 只在实际更新 FileSize 字段时验证文件大小
	// 注意：FileSize 在 BeforeCreate 已验证，通常不应在更新时修改
	if stmt.Changed("FileSize") && o.FileSize <= 0 {
		return ErrLogoInvalidFileSize
	}

	return nil
}

// IsTemporary 判断是否为临时Logo
func (o *OrganizationLogo) IsTemporary() bool {
	return o.Status == LogoStatusTemporary
}

// IsBound 判断是否已绑定
func (o *OrganizationLogo) IsBound() bool {
	return o.Status == LogoStatusBound
}

// IsExpired 判断是否已过期（仅对临时Logo有效）
func (o *OrganizationLogo) IsExpired() bool {
	if !o.IsTemporary() || o.ExpiresAt == nil {
		return false
	}

	return *o.ExpiresAt <= time.Now().UnixMilli()
}

// BindToOrganization 绑定到组织（将临时Logo转为永久）
func (o *OrganizationLogo) BindToOrganization(organizationID uuid.UUID) error {
	if !o.IsTemporary() {
		return fmt.Errorf("只能绑定临时状态的Logo")
	}

	if o.IsExpired() {
		return fmt.Errorf("Logo已过期，无法绑定")
	}

	o.Status = LogoStatusBound
	o.BoundOrganizationID = &organizationID
	o.ExpiresAt = nil // 清除过期时间

	return nil
}

// MarkAsDeleted 标记为删除状态
func (o *OrganizationLogo) MarkAsDeleted() {
	o.Status = LogoStatusDeleted
}
