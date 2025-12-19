package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Menu 菜单模型
// 用于在数据库中存储和管理层级菜单结构。
// 注意：此模型使用原生 UUID 作为主键，这是为了遵循 PostgreSQL 的最佳实践。
type Menu struct {
	BaseModel

	SemanticID string     `gorm:"column:semantic_id;not null;size:100;index:idx_semantic_version,unique;comment:语义化标识符"`
	Version    string     `gorm:"column:version;not null;index;size:50;index:idx_semantic_version,unique;comment:版本标识"`
	Name       string     `gorm:"column:name;not null;size:100;comment:菜单显示名称"`
	Path       string     `gorm:"column:path;not null;size:255;comment:前端路由路径"`
	Component  string     `gorm:"column:component;size:255;comment:前端组件的路径"`
	Icon       string     `gorm:"column:icon;size:100;comment:菜单图标的标识符	"`
	ParentID   *uuid.UUID `gorm:"column:parent_id;index;type:uuid;comment:父菜单ID"`
	Sort       int        `gorm:"column:sort;not null;default:0;comment:排序字段"`
	CreatedBy  *uuid.UUID `gorm:"column:created_by;type:uuid;comment:创建者ID"`
	UpdatedBy  *uuid.UUID `gorm:"column:updated_by;type:uuid;comment:最后更新者ID"`
	Parent     *Menu      `gorm:"foreignKey:ParentID;references:ID;comment:父菜单关联"`
	Children   []*Menu    `gorm:"foreignKey:ParentID;references:ID;comment:子菜单列表关联"`
}

// TableName 指定 Menu 模型对应的数据库表名。
func (Menu) TableName() string {
	return "menus"
}

// BeforeCreate 是一个 GORM 钩子，在创建新记录之前调用。
// 如果 ID 尚未被设置，它将自动生成一个新的 UUID。
func (m *Menu) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}

	return err
}
