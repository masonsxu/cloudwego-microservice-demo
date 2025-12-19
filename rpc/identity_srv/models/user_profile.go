package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserProfile 用户档案模型 (与 IDL identity_model.UserProfile 对应)
// 代表一个"自然人"，与其在组织中的角色和职位无关
type UserProfile struct {
	BaseModel

	// 用户标识
	Username     string `gorm:"column:username;uniqueIndex;not null;size:20;comment:用户名，唯一索引"`
	PasswordHash string `gorm:"column:password_hash;not null;size:255;comment:密码哈希"`
	Email        string `gorm:"column:email;index;size:255;comment:邮箱，索引"`
	Phone        string `gorm:"column:phone;index;size:20;comment:手机号，索引"`

	// 系统用户标识（与 RoleDefinition.IsSystemRole 对应）
	IsSystemUser bool `gorm:"column:is_system_user;not null;default:false;index;comment:是否为系统内置用户"`

	// 个人信息
	FirstName string `gorm:"column:first_name;size:50;comment:名"`
	LastName  string `gorm:"column:last_name;size:50;comment:姓"`
	RealName  string `gorm:"column:real_name;size:100;index;comment:真实姓名，索引"`
	Gender    Gender `gorm:"column:gender;default:0;comment:性别"`

	// 专业信息（仅用于展示，不用于权限判断）
	ProfessionalTitle string `gorm:"column:professional_title;size:100;comment:专业标题"`
	LicenseNumber     string `gorm:"column:license_number;size:100;index;comment:许可证号，索引"`
	Specialties       string `gorm:"column:specialties;type:text;comment:专业专长"`
	EmployeeID        string `gorm:"column:employee_id;size:50;index;comment:员工ID，索引"`

	// 状态管理
	Status             UserStatus `gorm:"column:status;not null;default:2;index;comment:用户状态"`
	LoginAttempts      int32      `gorm:"column:login_attempts;not null;default:0;comment:登录尝试次数"`
	MustChangePassword bool       `gorm:"column:must_change_password;not null;default:false;comment:是否必须修改密码"`
	AccountExpiry      *int64     `gorm:"column:account_expiry;comment:账户过期时间"`

	// 审计信息
	CreatedBy     *uuid.UUID `gorm:"column:created_by;type:uuid;comment:创建者ID"`
	UpdatedBy     *uuid.UUID `gorm:"column:updated_by;type:uuid;comment:更新者ID"`
	LastLoginTime *int64     `gorm:"column:last_login_time;comment:最后登录时间"`

	// 版本控制
	Version int32 `gorm:"column:version;not null;default:1;comment:版本号"`
}

// TableName 指定表名
func (UserProfile) TableName() string {
	return "user_profiles"
}

// BeforeCreate GORM钩子
func (u *UserProfile) BeforeCreate(tx *gorm.DB) error {
	if u.Status == 0 {
		u.Status = UserStatusInactive // 默认未激活状态
	}

	return u.validateFields(true)
}

// BeforeUpdate GORM钩子
func (u *UserProfile) BeforeUpdate(tx *gorm.DB) error {
	// 如果是系统用户，检查是否尝试修改 is_system_user 标识
	if u.IsSystemUser {
		// 从数据库获取原始值
		var original UserProfile
		if err := tx.Where("id = ?", u.ID).First(&original).Error; err == nil {
			if original.IsSystemUser && !u.IsSystemUser {
				return fmt.Errorf("系统用户的 is_system_user 标识不能被修改")
			}
		}
	}

	return u.validateFields(false)
}

// validateFields 验证字段
// isCreate: true表示创建操作，false表示更新操作
func (u *UserProfile) validateFields(isCreate bool) error {
	if isCreate {
		if u.Username == "" {
			return fmt.Errorf("用户名不能为空")
		}

		// 验证用户名格式 (与 IDL 的 validation.pattern 一致)
		if len(u.Username) < 3 || len(u.Username) > 20 {
			return fmt.Errorf("用户名长度必须在3-20个字符之间")
		}
	} else {
		// 更新操作只在用户名不为空时验证格式
		if u.Username != "" && (len(u.Username) < 3 || len(u.Username) > 20) {
			return fmt.Errorf("用户名长度必须在3-20个字符之间")
		}
	}

	return nil
}

// IsActive 检查用户是否处于活跃状态
func (u *UserProfile) IsActive() bool {
	if u.Status != UserStatusActive {
		return false
	}

	// 检查账户是否过期
	if u.AccountExpiry != nil && time.Now().UnixMilli() > *u.AccountExpiry {
		return false
	}

	return true
}

// IsLocked 检查用户是否被锁定
func (u *UserProfile) IsLocked() bool {
	return u.Status == UserStatusLocked
}

// ShouldChangePassword 检查是否需要强制修改密码
func (u *UserProfile) ShouldChangePassword() bool {
	return u.MustChangePassword
}

// IsSystem 判断是否为系统用户
func (u *UserProfile) IsSystem() bool {
	return u.IsSystemUser
}

// CanDelete 判断用户是否可以被删除
func (u *UserProfile) CanDelete() bool {
	return !u.IsSystemUser
}

// CanModifyUsername 判断用户名是否可以被修改
func (u *UserProfile) CanModifyUsername() bool {
	return !u.IsSystemUser
}
