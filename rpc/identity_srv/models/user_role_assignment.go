package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRoleAssignment 用户角色分配模型
type UserRoleAssignment struct {
	BaseModel

	UserID    uuid.UUID  `gorm:"column:user_id;not null;index;type:uuid;comment:用户ID"`
	RoleID    uuid.UUID  `gorm:"column:role_id;not null;index;type:uuid;comment:角色ID"`
	CreatedBy *uuid.UUID `gorm:"column:created_by;type:uuid;comment:创建者ID"`
	UpdatedBy *uuid.UUID `gorm:"column:updated_by;type:uuid;comment:最后更新者ID"`
}

// TableName 指定表名
func (UserRoleAssignment) TableName() string {
	return "user_role_assignments"
}

// BeforeCreate GORM钩子，在创建记录前执行。
func (u *UserRoleAssignment) BeforeCreate(tx *gorm.DB) error {
	// ID 由数据库默认生成，不再需要应用程序处理。
	return u.validateFields(tx)
}

// BeforeUpdate GORM钩子，在更新记录前执行。
func (u *UserRoleAssignment) BeforeUpdate(tx *gorm.DB) error {
	return u.validateFields(tx)
}

// validateFields 验证字段的业务规则。
func (u *UserRoleAssignment) validateFields(tx *gorm.DB) error {
	if u.RoleID == uuid.Nil {
		return fmt.Errorf("角色ID不能为空")
	}

	// 验证引用的角色是否存在
	var roleCount int64
	if err := tx.Model(&RoleDefinition{}).Where("id = ?", u.RoleID).Count(&roleCount).Error; err != nil {
		return fmt.Errorf("验证角色引用失败: %v", err)
	}

	if roleCount == 0 {
		return fmt.Errorf("引用的角色ID不存在: %s", u.RoleID)
	}

	return nil
}
