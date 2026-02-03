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

	// 产品线标识 (用于多租户/多产品线隔离)
	ProductLine string     `gorm:"column:product_line;not null;size:50;index:idx_product_semantic_version,priority:1;comment:产品线标识"`
	SemanticID  string     `gorm:"column:semantic_id;not null;size:100;index:idx_product_semantic_version,priority:2;comment:语义化标识符"`
	Version     int        `gorm:"column:version;not null;index:idx_product_semantic_version,priority:3;default:1;comment:版本号(自增)"`
	ContentHash string     `gorm:"column:content_hash;size:64;index;comment:内容哈希(SHA256)"`
	Name        string     `gorm:"column:name;not null;size:100;comment:菜单显示名称"`
	Path        string     `gorm:"column:path;not null;size:255;comment:前端路由路径"`
	Component   string     `gorm:"column:component;size:255;comment:前端组件的路径"`
	Icon        string     `gorm:"column:icon;size:100;comment:菜单图标的标识符	"`
	ParentID    *uuid.UUID `gorm:"column:parent_id;index;type:uuid;comment:父菜单ID"`
	Sort        int        `gorm:"column:sort;not null;default:0;comment:排序字段"`
	CreatedBy   *uuid.UUID `gorm:"column:created_by;type:uuid;comment:创建者ID"`
	UpdatedBy   *uuid.UUID `gorm:"column:updated_by;type:uuid;comment:最后更新者ID"`

	// Casbin 权限扩展字段
	PermCode string `gorm:"column:perm_code;size:100;index;comment:权限编码,如 emr:create, patient:read"`
	ApiPath  string `gorm:"column:api_path;size:255;comment:关联的API路径,如 /api/v1/patients"`
	IsButton bool   `gorm:"column:is_button;default:false;comment:是否为按钮级权限(非菜单项)"`

	// 关联关系
	Parent   *Menu   `gorm:"foreignKey:ParentID;references:ID;comment:父菜单关联"`
	Children []*Menu `gorm:"foreignKey:ParentID;references:ID;comment:子菜单列表关联"`
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

	// 自动生成 PermCode（如果未设置且有 SemanticID）
	if m.PermCode == "" && m.SemanticID != "" {
		m.PermCode = "menu:" + m.SemanticID
	}

	return err
}

// GetCasbinObject 获取用于 Casbin 策略的资源对象标识
func (m *Menu) GetCasbinObject() string {
	if m.PermCode != "" {
		return m.PermCode
	}
	if m.SemanticID != "" {
		return "menu:" + m.SemanticID
	}
	return "menu:" + m.ID.String()
}

// GetCasbinApiObject 获取用于 API 路径匹配的 Casbin 资源对象
func (m *Menu) GetCasbinApiObject() string {
	if m.ApiPath != "" {
		return m.ApiPath
	}
	return ""
}
