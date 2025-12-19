package models

import (
	"gorm.io/gorm"
)

// CasbinRule Casbin 策略规则表
// 使用 GORM Adapter 标准结构，扩展审计字段以支持权限管理需求
type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement;comment:主键"`
	Ptype string `gorm:"size:100;index;comment:策略类型：p, p2, p3, g, g2, g3"`
	V0    string `gorm:"size:256;index;comment:根据 ptype 不同含义不同"`
	V1    string `gorm:"size:256;index;comment:根据 ptype 不同含义不同"`
	V2    string `gorm:"size:256;index;comment:根据 ptype 不同含义不同"`
	V3    string `gorm:"size:256;index;comment:根据 ptype 不同含义不同"`
	V4    string `gorm:"size:256;index;comment:扩展字段1"`
	V5    string `gorm:"size:256;index;comment:扩展字段2"`

	// 审计字段
	CreatedAt int64          `gorm:"column:created_at;autoCreateTime:milli;index;comment:创建时间"`
	UpdatedAt int64          `gorm:"column:updated_at;autoUpdateTime:milli;index;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
	CreatedBy string         `gorm:"size:100;comment:创建者用户ID"`
	UpdatedBy string         `gorm:"size:100;comment:更新者用户ID"`
	Comment   string         `gorm:"size:500;comment:策略说明"`
}

// TableName 指定表名
func (CasbinRule) TableName() string {
	return "casbin_rule"
}

// 策略类型常量定义
const (
	// PolicyTypePermission 权限策略类型：角色 -> 资源 -> 动作
	PolicyTypePermission = "p"

	// PolicyTypeMenuMapping 菜单映射策略类型：角色 -> 菜单 -> 权限
	PolicyTypeMenuMapping = "p2"

	// PolicyTypeRoleInheritance 角色继承策略类型：子角色 -> 父角色
	PolicyTypeRoleInheritance = "g"

	// PolicyTypeUserRole 用户角色分配策略类型：用户 -> 角色
	PolicyTypeUserRole = "g2"
)

// 菜单权限类型常量定义
// 基于组织范围的权限控制
const (
	// MenuPermissionViewOwnOrganization 查看所在组织权限
	// 用户只能查看和操作自己所在组织的数据
	MenuPermissionViewOwnOrganization = "view_own_organization"

	// MenuPermissionViewAllOrganizations 查看所有组织权限
	// 用户可以查看和操作所有组织的数据（超管权限）
	MenuPermissionViewAllOrganizations = "view_all_organizations"
)

// IsValidMenuPermission 验证菜单权限类型是否有效
func IsValidMenuPermission(permission string) bool {
	switch permission {
	case MenuPermissionViewOwnOrganization, MenuPermissionViewAllOrganizations:
		return true
	default:
		return false
	}
}
