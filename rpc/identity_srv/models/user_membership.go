package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserMembership 用户-组织成员关系模型 (与 IDL identity_model.UserMembership 对应)
// 用于构建 Casbin 角色关系
type UserMembership struct {
	BaseModel

	// 核心关系映射
	UserID         uuid.UUID `gorm:"column:user_id;not null;type:uuid;index:idx_user_memberships;comment:用户ID"`
	OrganizationID uuid.UUID `gorm:"column:organization_id;not null;type:uuid;index:idx_org_memberships;comment:组织ID"`
	DepartmentID   uuid.UUID `gorm:"column:department_id;type:uuid;index:idx_dept_memberships;comment:部门ID"`

	// 关系状态
	Status    MembershipStatus `gorm:"column:status;not null;default:1;index;comment:成员状态"`
	IsPrimary bool             `gorm:"column:is_primary;not null;default:false;index;comment:是否主要成员"`
}

// TableName 指定表名
func (UserMembership) TableName() string {
	return "user_memberships"
}

// BeforeCreate GORM钩子
func (um *UserMembership) BeforeCreate(tx *gorm.DB) error {
	if um.Status == 0 {
		um.Status = MembershipStatusActive // 默认活跃状态
	}

	return um.validateForCreate(tx)
}

// BeforeUpdate GORM钩子
func (um *UserMembership) BeforeUpdate(tx *gorm.DB) error {
	return um.validateForUpdate(tx)
}

// validateForCreate 创建时的完整验证
func (um *UserMembership) validateForCreate(tx *gorm.DB) error {
	// 必填字段验证
	if um.UserID == uuid.Nil {
		return fmt.Errorf("用户ID不能为空")
	}

	if um.OrganizationID == uuid.Nil {
		return fmt.Errorf("组织ID不能为空")
	}

	// 验证引用完整性
	if err := um.validateReferences(tx); err != nil {
		return err
	}

	// 验证唯一性约束
	return um.validateUniqueness(tx)
}

// validateForUpdate 更新时的智能验证
func (um *UserMembership) validateForUpdate(tx *gorm.DB) error {
	// 检测是否为部分更新（核心字段为零值时）
	if um.UserID == uuid.Nil || um.OrganizationID == uuid.Nil {
		var existing UserMembership
		if err := tx.Where("id = ?", um.ID).First(&existing).Error; err != nil {
			return fmt.Errorf("加载现有记录失败: %v", err)
		}

		// 使用现有值填充零值字段，确保验证逻辑正常工作
		if um.UserID == uuid.Nil {
			um.UserID = existing.UserID
		}

		if um.OrganizationID == uuid.Nil {
			um.OrganizationID = existing.OrganizationID
		}
		// 如果部门ID也是零值，使用现有值
		if um.DepartmentID == uuid.Nil {
			um.DepartmentID = existing.DepartmentID
		}
	}

	// 验证引用完整性（仅对非零值字段）
	if err := um.validateReferences(tx); err != nil {
		return err
	}

	// 验证唯一性约束
	return um.validateUniqueness(tx)
}

// validateReferences 验证引用完整性
func (um *UserMembership) validateReferences(tx *gorm.DB) error {
	// 验证用户存在（仅当UserID非零时）
	if um.UserID != uuid.Nil {
		var userCount int64
		if err := tx.Model(&UserProfile{}).Where("id = ?", um.UserID).Count(&userCount).Error; err != nil {
			return fmt.Errorf("验证用户引用失败: %v", err)
		}

		if userCount == 0 {
			return fmt.Errorf("引用的用户ID不存在: %s", um.UserID)
		}
	}

	// 验证组织存在（仅当OrganizationID非零时）
	if um.OrganizationID != uuid.Nil {
		var orgCount int64
		if err := tx.Model(&Organization{}).Where("id = ?", um.OrganizationID).Count(&orgCount).Error; err != nil {
			return fmt.Errorf("验证组织引用失败: %v", err)
		}

		if orgCount == 0 {
			return fmt.Errorf("引用的组织ID不存在: %s", um.OrganizationID)
		}
	}

	// 验证部门（如果指定且非零）
	if um.DepartmentID != uuid.Nil {
		var deptCount int64
		if err := tx.Model(&Department{}).
			Where("id = ? AND organization_id = ?", um.DepartmentID, um.OrganizationID).
			Count(&deptCount).Error; err != nil {
			return fmt.Errorf("验证部门引用失败: %v", err)
		}

		if deptCount == 0 {
			return fmt.Errorf("引用的部门ID不存在或不属于指定组织: %s", um.DepartmentID)
		}
	}

	return nil
}

// validateUniqueness 验证唯一性约束
func (um *UserMembership) validateUniqueness(tx *gorm.DB) error {
	var existingCount int64

	query := tx.Model(&UserMembership{}).Where(
		"user_id = ? AND organization_id = ? AND status = ?",
		um.UserID, um.OrganizationID, MembershipStatusActive,
	)

	// 处理部门条件
	if um.DepartmentID != uuid.Nil {
		query = query.Where("department_id = ?", um.DepartmentID)
	} else {
		query = query.Where("department_id IS NULL")
	}

	// 排除当前记录（更新场景）
	if um.ID != uuid.Nil {
		query = query.Where("id != ?", um.ID)
	}

	if err := query.Count(&existingCount).Error; err != nil {
		return fmt.Errorf("检查成员关系唯一性失败: %v", err)
	}

	if existingCount > 0 {
		return fmt.Errorf("用户在该组织部门已存在活跃的成员关系")
	}

	return nil
}

// IsActive 检查成员关系是否活跃
func (um *UserMembership) IsActive() bool {
	return um.Status == MembershipStatusActive
}
