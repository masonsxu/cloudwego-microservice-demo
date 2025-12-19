package user

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// UserProfileRepository 用户档案仓储接口
// 基于 models.UserProfile 和 IDL 设计，提供用户档案的完整数据访问能力
type UserProfileRepository interface {
	// 嵌入基础仓储接口
	base.BaseRepository[models.UserProfile]

	// ============================================================================
	// 唯一性查询（用于认证和去重）
	// ============================================================================

	// GetByUsername 根据用户名获取用户档案
	GetByUsername(ctx context.Context, username string) (*models.UserProfile, error)

	// GetByEmail 根据邮箱获取用户档案
	GetByEmail(ctx context.Context, email string) (*models.UserProfile, error)

	// GetByPhone 根据手机号获取用户档案
	GetByPhone(ctx context.Context, phone string) (*models.UserProfile, error)

	// ExistsByID 检查用户ID是否存在
	ExistsByID(ctx context.Context, userID string) (bool, error)

	// ============================================================================
	// 数据完整性检查
	// ============================================================================

	// CheckUsernameExists 检查用户名是否已存在
	CheckUsernameExists(ctx context.Context, username string, excludeID ...string) (bool, error)

	// CheckEmailExists 检查邮箱是否已存在
	CheckEmailExists(ctx context.Context, email string, excludeID ...string) (bool, error)

	// CheckPhoneExists 检查手机号是否已存在
	CheckPhoneExists(ctx context.Context, phone string, excludeID ...string) (bool, error)

	// ============================================================================
	// 状态管理操作
	// ============================================================================

	// UpdatePassword 更新用户密码
	UpdatePassword(ctx context.Context, userID string, passwordHash string) error

	// UpdateLoginAttempts 更新登录尝试次数
	UpdateLoginAttempts(ctx context.Context, userID string, attempts int32) error

	// IncrementLoginAttempts 增加登录失败次数
	IncrementLoginAttempts(ctx context.Context, userID string) error

	// ResetLoginAttempts 重置登录尝试次数
	ResetLoginAttempts(ctx context.Context, userID string) error

	// UpdateLastLoginTime 更新最后登录时间
	UpdateLastLoginTime(ctx context.Context, userID string) error

	// SetMustChangePassword 设置强制修改密码标志
	SetMustChangePassword(ctx context.Context, userID string, mustChange bool) error

	// ============================================================================
	// 专业信息查询
	// ============================================================================

	// FindByMedicalLicense 根据执照号查询用户档案
	FindByMedicalLicense(ctx context.Context, licenseNumber string) (*models.UserProfile, error)

	// FindBySpecialty 根据专业领域查询用户档案列表
	FindBySpecialty(
		ctx context.Context,
		specialty string,
		opts *base.QueryOptions,
	) ([]*models.UserProfile, *models.PageResult, error)

	// ============================================================================
	// 统一查询方法（新增）
	// ============================================================================

	// FindWithConditions 根据组合查询条件查询用户档案列表
	FindWithConditions(
		ctx context.Context,
		conditions *UserProfileQueryConditions,
	) ([]*models.UserProfile, *models.PageResult, error)

	// ============================================================================
	// 系统用户管理
	// ============================================================================

	// FindSystemUsers 查询所有系统用户
	FindSystemUsers(ctx context.Context) ([]*models.UserProfile, error)

	// IsSystemUser 判断用户是否为系统用户
	IsSystemUser(ctx context.Context, userID string) (bool, error)
}

// UserProfileQueryConditions 用户档案查询条件
// 支持多条件组合查询，提供灵活的查询能力
type UserProfileQueryConditions struct {
	Username       *string            // 用户名（精确匹配）
	Email          *string            // 邮箱（精确匹配）
	Phone          *string            // 手机号（精确匹配）
	Status         *models.UserStatus // 用户状态
	OrgID          *string            // 组织ID（通过成员关系查询）
	MedicalLicense *string            // 执照号
	Specialty      *string            // 专业领域
	Page           *base.QueryOptions // 分页、排序、搜索参数
}
