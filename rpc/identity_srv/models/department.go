package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Department 部门/部门模型 (与 IDL identity_model.Department 对应)
type Department struct {
	BaseModel

	Name           string    `gorm:"column:name;not null;size:100;index;comment:部门名称，用于搜索"`
	OrganizationID uuid.UUID `gorm:"column:organization_id;not null;type:uuid;index:idx_org_departments;comment:组织ID"`

	// 部门属性
	DepartmentType     string `gorm:"column:department_type;size:100;index;comment:部门类型"`
	AvailableEquipment string `gorm:"column:available_equipment;type:text;comment:JSON 存储 list<ULID>"`
}

// TableName 指定表名
func (Department) TableName() string {
	return "departments"
}

// BeforeCreate GORM钩子
func (d *Department) BeforeCreate(tx *gorm.DB) error {
	return d.validateFields(tx)
}

// BeforeUpdate GORM钩子
func (d *Department) BeforeUpdate(tx *gorm.DB) error {
	return d.validateFields(tx)
}

// validateFields 验证字段
func (d *Department) validateFields(tx *gorm.DB) error {
	if d.Name == "" {
		return fmt.Errorf("部门名称不能为空")
	}

	if d.OrganizationID == uuid.Nil {
		return fmt.Errorf("部门必须归属于一个组织")
	}

	// 验证组织是否存在
	var orgCount int64
	if err := tx.Model(&Organization{}).Where("id = ?", d.OrganizationID).Count(&orgCount).Error; err != nil {
		return fmt.Errorf("验证组织引用失败: %v", err)
	}

	if orgCount == 0 {
		return fmt.Errorf("引用的组织ID不存在: %s", d.OrganizationID)
	}

	// 验证部门名称在同一组织内的唯一性
	var existingCount int64

	query := tx.Model(&Department{}).
		Where("name = ? AND organization_id = ?", d.Name, d.OrganizationID)

	// 排除当前记录（更新场景）
	if d.ID != uuid.Nil {
		query = query.Where("id != ?", d.ID)
	}

	if err := query.Count(&existingCount).Error; err != nil {
		return fmt.Errorf("检查部门名称唯一性失败: %v", err)
	}

	if existingCount > 0 {
		return fmt.Errorf("该组织下已存在同名部门: %s", d.Name)
	}

	return nil
}
