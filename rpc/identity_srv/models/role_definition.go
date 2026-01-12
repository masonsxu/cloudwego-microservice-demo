package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

// DataScopeType 数据范围类型
type DataScopeType int8

const (
	// DataScopeSelf 仅本人数据
	DataScopeSelf DataScopeType = 1
	// DataScopeDept 本科室/部门数据
	DataScopeDept DataScopeType = 2
	// DataScopeOrg 全院/组织数据
	DataScopeOrg DataScopeType = 3
)

// String 返回数据范围类型的字符串表示
func (d DataScopeType) String() string {
	switch d {
	case DataScopeSelf:
		return "self"
	case DataScopeDept:
		return "dept"
	case DataScopeOrg:
		return "org"
	default:
		return "unknown"
	}
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

	// Casbin 权限扩展字段
	RoleCode     string        `gorm:"column:role_code;uniqueIndex;size:50;comment:角色编码,用于Casbin策略标识"`
	ParentRoleID *uuid.UUID    `gorm:"column:parent_role_id;type:uuid;index;comment:父角色ID,支持角色继承"`
	DepartmentID *uuid.UUID    `gorm:"column:department_id;type:uuid;index;comment:绑定科室ID,NULL表示全院通用角色"`
	DefaultScope DataScopeType `gorm:"column:default_scope;default:1;comment:默认数据范围:1-本人,2-本科室,3-全院"`

	// 关联关系
	ParentRole *RoleDefinition `gorm:"foreignKey:ParentRoleID;references:ID;comment:父角色关联"`

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

	// 自动生成 RoleCode（如果未设置）
	if r.RoleCode == "" && r.Name != "" {
		r.RoleCode = r.GenerateRoleCode()
	}

	// 默认数据范围为本人
	if r.DefaultScope == 0 {
		r.DefaultScope = DataScopeSelf
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

// GenerateRoleCode 根据角色名称生成角色编码
// 规则：将名称转换为 snake_case 格式，添加 role: 前缀
func (r *RoleDefinition) GenerateRoleCode() string {
	if r.Name == "" {
		return ""
	}

	// 移除变音符号
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalized, _, _ := transform.String(t, r.Name)

	// 转换为小写
	normalized = strings.ToLower(normalized)

	// 将中文和特殊字符替换为下划线
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	normalized = reg.ReplaceAllString(normalized, "_")

	// 移除首尾下划线
	normalized = strings.Trim(normalized, "_")

	// 压缩连续下划线
	reg = regexp.MustCompile(`_+`)
	normalized = reg.ReplaceAllString(normalized, "_")

	if normalized == "" {
		// 如果全是中文或特殊字符，使用 UUID 的前8位
		normalized = uuid.New().String()[:8]
	}

	return "role:" + normalized
}

// GetCasbinSubject 获取用于 Casbin 策略的主体标识
func (r *RoleDefinition) GetCasbinSubject() string {
	if r.RoleCode != "" {
		return r.RoleCode
	}
	return "role:" + r.ID.String()
}

// GetCasbinDomain 获取用于 Casbin 策略的域标识
func (r *RoleDefinition) GetCasbinDomain() string {
	if r.DepartmentID != nil {
		return "dept:" + r.DepartmentID.String()
	}
	return "*" // 全院通用
}
