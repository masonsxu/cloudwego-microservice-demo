package membership

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// UserMembershipRepository 用户成员关系仓储接口
// 基于 models.UserMembership 和 IDL 设计，管理用户与组织/部门的成员关系
type UserMembershipRepository interface {
	// 嵌入基础仓储接口
	base.BaseRepository[models.UserMembership]

	// ============================================================================
	// 核心成员关系查询
	// ============================================================================

	// GetByUserAndOrganization 获取用户在指定组织的成员关系
	GetByUserAndOrganization(
		ctx context.Context,
		userID, organizationID string,
	) (*models.UserMembership, error)

	// GetByUserAndDepartment 获取用户在指定部门的成员关系
	GetByUserAndDepartment(
		ctx context.Context,
		userID, departmentID string,
	) (*models.UserMembership, error)

	// GetPrimaryMembership 获取用户的主要成员关系
	GetPrimaryMembership(ctx context.Context, userID string) (*models.UserMembership, error)

	// GetPrimaryMembershipsByUserIDs 批量获取多个用户的主成员关系
	// 返回: map[userID]*UserMembership，避免 N+1 查询问题
	GetPrimaryMembershipsByUserIDs(
		ctx context.Context,
		userIDs []string,
	) (map[string]*models.UserMembership, error)

	// CountByDepartmentID 统计部门的成员数量
	CountByDepartmentID(ctx context.Context, departmentID string) (int64, error)

	// ============================================================================
	// 组织层级查询
	// ============================================================================

	// FindUsersByOrganizationHierarchy 查询组织及其子组织的所有用户成员关系
	FindUsersByOrganizationHierarchy(
		ctx context.Context,
		rootOrganizationID string,
		opts *base.QueryOptions,
	) ([]*models.UserMembership, *models.PageResult, error)

	// GetUserOrganizations 获取用户所属的所有组织ID列表
	GetUserOrganizations(ctx context.Context, userID string) ([]string, error)

	// GetOrganizationUsers 获取组织的所有用户ID列表
	GetOrganizationUsers(ctx context.Context, organizationID string) ([]string, error)

	// ============================================================================
	// 成员关系管理操作
	// ============================================================================

	// UpdateStatus 更新成员关系状态
	UpdateStatus(ctx context.Context, membershipID string, status models.MembershipStatus) error

	// SetPrimaryMembership 设置主要成员关系（会将其他关系的primary标志置为false）
	SetPrimaryMembership(ctx context.Context, userID, membershipID string) error

	// UnsetPrimaryByUserID 取消用户的所有主要成员关系标志
	UnsetPrimaryByUserID(ctx context.Context, userID string) error

	// ExistsByID 检查成员关系是否存在
	ExistsByID(ctx context.Context, membershipID string) (bool, error)

	// ============================================================================
	// 批量操作
	// ============================================================================

	// BatchCreateMemberships 批量创建成员关系
	BatchCreateMemberships(ctx context.Context, memberships []*models.UserMembership) error

	// BatchUpdateStatus 批量更新成员关系状态
	BatchUpdateStatus(
		ctx context.Context,
		membershipIDs []string,
		status models.MembershipStatus,
	) error

	// BatchDeleteByUser 批量删除用户的所有成员关系
	BatchDeleteByUser(ctx context.Context, userID string) error

	// BatchDeleteByOrganization 批量删除组织的所有成员关系
	BatchDeleteByOrganization(ctx context.Context, organizationID string) error

	// ============================================================================
	// 数据完整性检查
	// ============================================================================

	// CheckMembershipExists 检查成员关系是否已存在
	CheckMembershipExists(
		ctx context.Context,
		userID, organizationID, departmentID string,
	) (bool, error)

	// ============================================================================
	// 统计分析
	// ============================================================================

	// CountByOrganization 统计组织的成员数量
	CountByOrganization(ctx context.Context, organizationID string) (int64, error)

	// GetMembershipStatistics 获取成员关系统计信息
	GetMembershipStatistics(
		ctx context.Context,
		organizationID string,
	) (*MembershipStatistics, error)

	// ============================================================================
	// 统一查询方法
	// ============================================================================

	// FindWithConditions 根据组合查询条件查询成员关系列表
	FindWithConditions(
		ctx context.Context,
		conditions *UserMembershipQueryConditions,
	) ([]*models.UserMembership, *models.PageResult, error)
}

// UserMembershipQueryConditions 成员关系查询条件
// 支持多条件组合查询，提供灵活的查询能力
type UserMembershipQueryConditions struct {
	UserID         *string                  // 用户ID
	OrganizationID *string                  // 组织ID
	DepartmentID   *string                  // 部门ID
	Status         *models.MembershipStatus // 成员关系状态
	IsPrimary      *bool                    // 是否为主要成员关系
	Page           *base.QueryOptions       // 分页、排序、搜索参数
}
