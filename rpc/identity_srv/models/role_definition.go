package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permissions 是一个自定义类型，用于处理权限列表的 JSONB 存储和读取。
type Permissions []*Permission

// Value 实现了 driver.Valuer 接口，用于将 Permissions 类型写入数据库。
func (p Permissions) Value() (driver.Value, error) {
	if len(p) == 0 {
		return nil, nil
	}

	return json.Marshal(p)
}

// Scan 实现了 sql.Scanner 接口，用于从数据库读取数据到 Permissions 类型。
func (p *Permissions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &p)
}

// RoleDefinition 角色定义模型
type RoleDefinition struct {
	BaseModel

	Name         string      `gorm:"column:name;uniqueIndex;not null;size:50;comment:角色唯一名称"`
	Description  string      `gorm:"column:description;type:text;comment:角色详细描述"`
	Status       RoleStatus  `gorm:"column:status;not null;comment:角色状态:1-活跃,2-未激活,3-已弃用"`
	Permissions  Permissions `gorm:"column:permissions;type:jsonb;comment:角色拥有的权限列表"`
	IsSystemRole bool        `gorm:"column:is_system_role;not null;default:false;comment:是否为系统内置角色"`
	CreatedBy    *uuid.UUID  `gorm:"column:created_by;type:uuid;comment:创建者ID"`
	UpdatedBy    *uuid.UUID  `gorm:"column:updated_by;type:uuid;comment:最后更新者ID"`

	// 当前角色绑定的用户数量（非数据库字段，用于业务逻辑传递）
	UserCount int64 `gorm:"-" json:"user_count,omitempty"`
}

// TableName 指定表名
func (RoleDefinition) TableName() string {
	return "role_definitions"
}

// BeforeCreate GORM钩子，在创建记录前执行。
func (r *RoleDefinition) BeforeCreate(tx *gorm.DB) error {
	// ID 由数据库默认生成，不再需要应用程序处理。
	if r.Status == 0 {
		r.Status = RoleStatusInactive // 默认未激活状态
	}

	return r.validateFields(true)
}

// BeforeUpdate GORM钩子，在更新记录前执行。
func (r *RoleDefinition) BeforeUpdate(tx *gorm.DB) error {
	return r.validateFields(false)
}

// validateFields 验证核心字段的业务规则。
// isCreate: true表示创建操作，false表示更新操作
func (r *RoleDefinition) validateFields(isCreate bool) error {
	if isCreate {
		if r.Name == "" {
			return fmt.Errorf("角色名称不能为空")
		}

		// 验证角色名称格式
		if len(r.Name) < 2 || len(r.Name) > 50 {
			return fmt.Errorf("角色名称长度必须在2-50个字符之间")
		}
	} else {
		// 更新操作只在名称不为空时验证格式
		if r.Name != "" && (len(r.Name) < 2 || len(r.Name) > 50) {
			return fmt.Errorf("角色名称长度必须在2-50个字符之间")
		}
	}

	return nil
}
