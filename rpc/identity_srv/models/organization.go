package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StringSlice 自定义切片类型，用于JSON存储
type StringSlice []string

// Value 实现 driver.Valuer 接口，将切片转换为JSON存储到数据库
func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}

	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口，从数据库JSON读取并解析为切片
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = StringSlice{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return fmt.Errorf("cannot scan %T into StringSlice", value)
	}
}

// Organization 组织模型 (与 IDL identity_model.Organization 对应)
type Organization struct {
	BaseModel

	Code     string    `gorm:"column:code;not null;size:50;uniqueIndex;comment:组织代码，必须唯一"`
	Name     string    `gorm:"column:name;not null;size:100;index;comment:组织名称，用于搜索"`
	ParentID uuid.UUID `gorm:"column:parent_id;type:uuid;index:idx_parent_org;comment:支持层级组织结构"`

	// 组织属性
	FacilityType        string      `gorm:"column:facility_type;size:100;index;comment:组织类型"`
	AccreditationStatus string      `gorm:"column:accreditation_status;size:100;comment:认证状态"`
	ProvinceCity        StringSlice `gorm:"column:province_city;type:json;comment:组织所在省市列表"`
}

// TableName 指定表名
func (Organization) TableName() string {
	return "organizations"
}

// BeforeCreate GORM钩子
func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	// 如果code为空，自动生成
	if o.Code == "" {
		code, err := o.generateCode(tx)
		if err != nil {
			return fmt.Errorf("生成组织代码失败: %v", err)
		}

		o.Code = code
	}

	return o.validateFields(tx)
}

// BeforeUpdate GORM钩子
func (o *Organization) BeforeUpdate(tx *gorm.DB) error {
	return o.validateFields(tx)
}

// validateFields 验证字段（简化版，支持2级层级）
func (o *Organization) validateFields(tx *gorm.DB) error {
	if o.Name == "" {
		return fmt.Errorf("组织名称不能为空")
	}

	if o.Code == "" {
		return fmt.Errorf("组织代码不能为空")
	}

	// 限制为2级层级结构
	if o.ParentID != uuid.Nil {
		return o.validateTwoLevelHierarchy(tx)
	}

	return nil
}

// validateTwoLevelHierarchy 验证2级层级限制
func (o *Organization) validateTwoLevelHierarchy(tx *gorm.DB) error {
	// 防止自引用
	if o.ParentID == o.ID {
		return fmt.Errorf("组织不能将自己设置为父组织")
	}

	// 检查父组织是否存在
	var parent Organization
	if err := tx.Where("id = ?", o.ParentID).First(&parent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("引用的父组织ID不存在: %s", o.ParentID)
		}

		return fmt.Errorf("验证父组织引用失败: %v", err)
	}

	// 限制为2级层级：父组织不能有父组织
	if parent.ParentID != uuid.Nil {
		return fmt.Errorf("不支持超过2级的组织层级，父组织不能为子组织")
	}

	return nil
}

// IsRootOrganization 检查是否为根组织
func (o *Organization) IsRootOrganization() bool {
	return o.ParentID == uuid.Nil
}

// IsSubOrganization 检查是否为子组织
func (o *Organization) IsSubOrganization() bool {
	return o.ParentID != uuid.Nil
}

// generateCode 生成唯一的组织代码
func (o *Organization) generateCode(tx *gorm.DB) (string, error) {
	// 基于组织名称生成代码前缀
	baseCode := o.generateCodeFromName(o.Name)

	// 检查代码唯一性，如果重复则添加时间戳后缀
	code := baseCode

	for range 10 { // 最多尝试10次
		var count int64
		if err := tx.Model(&Organization{}).Where("code = ?", code).Count(&count).Error; err != nil {
			return "", err
		}

		if count == 0 {
			return code, nil
		}

		// 如果重复，添加时间戳后缀
		timestamp := time.Now().Format("060102150405") // YYMMDDHHMMSS
		code = fmt.Sprintf("%s_%s", baseCode, timestamp)
	}

	return "", fmt.Errorf("无法生成唯一的组织代码")
}

// generateCodeFromName 从组织名称生成代码
func (o *Organization) generateCodeFromName(name string) string {
	// 去除特殊字符，只保留字母和数字
	reg := regexp.MustCompile(`[^a-zA-Z0-9\x{4e00}-\x{9fa5}]`)
	cleaned := reg.ReplaceAllString(name, "")

	// 如果是中文，取拼音首字母（简化处理：这里暂时使用前几个字符）
	if len(cleaned) > 8 {
		cleaned = cleaned[:8]
	}

	// 转为大写
	code := strings.ToUpper(cleaned)

	// 如果代码为空或太短，使用默认前缀
	if len(code) < 2 {
		code = "ORG"
	}

	return code
}
