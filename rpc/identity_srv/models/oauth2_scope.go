package models

// OAuth2Scope OAuth2 作用域定义模型
// 定义 OAuth2 中可授权的权限范围，与 Casbin 权限体系映射。
type OAuth2Scope struct {
	BaseModel

	// 作用域标识
	Name        string `gorm:"column:name;uniqueIndex;size:64;not null;comment:作用域名称(如 user:read)"`
	DisplayName string `gorm:"column:display_name;size:128;not null;comment:显示名称"`
	Description string `gorm:"column:description;size:512;comment:作用域描述"`

	// 分类
	IsDefault bool `gorm:"column:is_default;default:false;comment:是否为默认授予的作用域"`
	IsSystem  bool `gorm:"column:is_system;default:false;comment:是否为系统内置(不可删除)"`
}

func (OAuth2Scope) TableName() string {
	return "oauth2_scopes"
}
