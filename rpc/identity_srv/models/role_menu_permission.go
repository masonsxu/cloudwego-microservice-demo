package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MenuPermissionType 菜单权限类型枚举
type MenuPermissionType int8

const (
	PermissionNone   MenuPermissionType = 0 // 无权限
	PermissionView   MenuPermissionType = 1 // 查看
	PermissionEdit   MenuPermissionType = 2 // 编辑
	PermissionManage MenuPermissionType = 3 // 管理
	PermissionFull   MenuPermissionType = 4 // 完全控制
)

// String 返回权限类型的字符串表示
func (p MenuPermissionType) String() string {
	switch p {
	case PermissionNone:
		return "none"
	case PermissionView:
		return "view"
	case PermissionEdit:
		return "edit"
	case PermissionManage:
		return "manage"
	case PermissionFull:
		return "full"
	default:
		return "unknown"
	}
}

// PermissionLevel 返回前端期望的权限级别字符串
// 前端定义: "none" / "read" / "write" / "full"
func (p MenuPermissionType) PermissionLevel() string {
	switch p {
	case PermissionNone:
		return "none"
	case PermissionView:
		return "read"
	case PermissionEdit, PermissionManage:
		return "write"
	case PermissionFull:
		return "full"
	default:
		return "none"
	}
}

// ParsePermissionType 从字符串解析权限类型
// 支持内部值 (view/edit/manage/full) 和前端值 (read/write/full/none)
func ParsePermissionType(s string) MenuPermissionType {
	switch s {
	case "view", "read":
		return PermissionView
	case "edit", "write":
		return PermissionEdit
	case "manage":
		return PermissionManage
	case "full":
		return PermissionFull
	case "none":
		return PermissionNone
	default:
		return PermissionNone
	}
}

// DataScope 数据范围枚举
type DataScope int8

const (
	DataScopeNone    DataScope = 0 // 无数据范围
	DataScopeOwnOrg  DataScope = 1 // 所在组织
	DataScopeAllOrgs DataScope = 2 // 所有组织
)

// String 返回数据范围的字符串表示
func (d DataScope) String() string {
	switch d {
	case DataScopeNone:
		return "none"
	case DataScopeOwnOrg:
		return "own_org"
	case DataScopeAllOrgs:
		return "all_orgs"
	default:
		return "unknown"
	}
}

// ParseDataScope 从字符串解析数据范围
func ParseDataScope(s string) DataScope {
	switch s {
	case "own_org", "view_own_organization":
		return DataScopeOwnOrg
	case "all_orgs", "view_all_organizations":
		return DataScopeAllOrgs
	default:
		return DataScopeNone
	}
}

// RoleMenuPermission 角色菜单权限模型
// 存储角色与菜单的权限映射关系
type RoleMenuPermission struct {
	BaseModel

	RoleID         uuid.UUID          `gorm:"column:role_id;type:uuid;not null;index:idx_role_menu,unique;comment:角色ID"`
	MenuID         string             `gorm:"column:menu_id;size:100;not null;index:idx_role_menu,unique;comment:菜单ID"`
	PermissionType MenuPermissionType `gorm:"column:permission_type;not null;default:1;comment:权限类型:无/查看/编辑/管理/完全控制"`
	DataScope      DataScope          `gorm:"column:data_scope;not null;default:1;comment:数据范围:0-无,1-所在组织,2-所有组织"`

	// 关联
	Role *RoleDefinition `gorm:"foreignKey:RoleID;references:ID"`
}

// TableName 指定表名
func (RoleMenuPermission) TableName() string {
	return "role_menu_permissions"
}

// BeforeCreate GORM钩子，在创建记录前执行
func (r *RoleMenuPermission) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}

	return nil
}

// MergePermissionTypes 合并多个权限类型，取最高权限
func MergePermissionTypes(types ...MenuPermissionType) MenuPermissionType {
	max := PermissionNone
	for _, t := range types {
		if t > max {
			max = t
		}
	}

	return max
}

// MergeDataScopes 合并多个数据范围，取最广范围
func MergeDataScopes(scopes ...DataScope) DataScope {
	max := DataScopeNone
	for _, s := range scopes {
		if s > max {
			max = s
		}
	}

	return max
}

// MenuPermissionInfo 菜单权限信息（用于业务逻辑传递）
type MenuPermissionInfo struct {
	MenuID         string
	PermissionType MenuPermissionType
	DataScope      DataScope
}

// MergeMenuPermissions 合并多个菜单权限列表，对相同菜单取最高权限
func MergeMenuPermissions(permissionLists ...[]MenuPermissionInfo) []MenuPermissionInfo {
	merged := make(map[string]*MenuPermissionInfo)

	for _, list := range permissionLists {
		for _, perm := range list {
			if existing, ok := merged[perm.MenuID]; ok {
				existing.PermissionType = MergePermissionTypes(existing.PermissionType, perm.PermissionType)
				existing.DataScope = MergeDataScopes(existing.DataScope, perm.DataScope)
			} else {
				merged[perm.MenuID] = &MenuPermissionInfo{
					MenuID:         perm.MenuID,
					PermissionType: perm.PermissionType,
					DataScope:      perm.DataScope,
				}
			}
		}
	}

	result := make([]MenuPermissionInfo, 0, len(merged))
	for _, perm := range merged {
		result = append(result, *perm)
	}

	return result
}
