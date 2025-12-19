package assignment

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// UserRoleAssignmentQueryConditions 用户角色分配查询条件
// 对应 IDL 中的 UserRoleQueryRequest，用于仓储层查询
type UserRoleAssignmentQueryConditions struct {
	// UserID 用户ID，用于查询指定用户的角色分配
	UserID *string `json:"user_id,omitempty"`

	// RoleID 角色ID，用于查询指定角色的分配情况
	RoleID *string `json:"role_id,omitempty"`

	// Page 分页查询选项
	Page *base.QueryOptions `json:"page,omitempty"`
}

// UserRoleAssignmentRepository 用户角色分配仓储接口
// 基于 models.UserRoleAssignment，管理用户在组织/部门中的角色分配关系
type UserRoleAssignmentRepository interface {
	// 嵌入基础仓储接口
	base.BaseRepository[models.UserRoleAssignment]

	// ============================================================================
	// 核心查询方法
	// ============================================================================

	// FindByUserID 查询指定用户的所有角色分配记录
	// 支持分页，用于查看用户的角色历史和当前角色
	FindByUserID(
		ctx context.Context,
		userID string,
		page *base.QueryOptions,
	) ([]*models.UserRoleAssignment, *models.PageResult, error)

	// FindByRoleID 查询指定角色的所有分配记录
	// 支持分页，用于查看哪些用户被分配了某个角色
	FindByRoleID(
		ctx context.Context,
		roleID string,
		page *base.QueryOptions,
	) ([]*models.UserRoleAssignment, *models.PageResult, error)

	// FindByUserAndRole 查询指定用户和角色的分配记录
	// 用于验证用户是否拥有某个角色，或获取具体的分配详情
	FindByUserAndRole(
		ctx context.Context,
		userID, roleID string,
	) (*models.UserRoleAssignment, error)

	// GetLastUserRoleAssignment 获取用户最后一次的角色分配信息
	// 对应 IDL 中的 GetLastUserRoleAssignment 接口
	GetLastUserRoleAssignment(
		ctx context.Context,
		userID string,
	) (*models.UserRoleAssignment, error)

	// GetActiveRolesByUserID 获取用户所有活跃的角色ID列表（未删除）
	// 用于多角色权限合并场景，返回用户当前拥有的所有角色ID
	GetActiveRolesByUserID(
		ctx context.Context,
		userID string,
	) ([]string, error)

	// GetActiveRoleIDsWithStatus 获取用户所有处于指定状态的角色ID列表
	// 此方法会联表查询 user_role_assignments 和 role_definitions
	// 只返回角色状态匹配的角色ID（通常用于获取 Active 状态的角色）
	// 专门用于登录等需要验证角色可用性的场景
	GetActiveRoleIDsWithStatus(
		ctx context.Context,
		userID string,
		status models.RoleStatus,
	) ([]string, error)

	// FindWithConditions 根据组合查询条件查询角色分配记录
	// 这是最灵活的查询方法，对应 IDL 中的 UserRoleQueryRequest
	FindWithConditions(
		ctx context.Context,
		conditions *UserRoleAssignmentQueryConditions,
	) ([]*models.UserRoleAssignment, *models.PageResult, error)

	// GetRolesByUserIDs 批量查询多个用户的角色分配
	// 返回: map[userID][]roleID，避免 N+1 查询问题
	GetRolesByUserIDs(
		ctx context.Context,
		userIDs []string,
	) (map[string][]string, error)

	// ============================================================================
	// 业务验证方法
	// ============================================================================

	// CheckUserRoleExists 检查用户是否已分配指定角色
	// 用于分配前的重复性验证，避免重复分配
	CheckUserRoleExists(ctx context.Context, userID, roleID string) (bool, error)

	// CountByUserID 统计指定用户的角色分配数量
	// 用于用户角色概览和权限分析
	CountByUserID(ctx context.Context, userID string) (int64, error)

	// CountByRoleID 统计指定角色的分配数量
	// 用于角色使用情况分析和权限管理
	CountByRoleID(ctx context.Context, roleID string) (int64, error)

	// ============================================================================
	// 批量操作方法
	// ============================================================================

	// BatchAssignRoleToUsers 批量为多个用户分配同一角色
	// 用于批量用户权限管理场景
	BatchAssignRoleToUsers(
		ctx context.Context,
		userIDs []string,
		roleID string,
		assignedBy string,
	) error

	// BatchRevokeUserRoles 批量撤销用户的多个角色分配
	// 用于用户离职或权限批量调整场景
	BatchRevokeUserRoles(ctx context.Context, assignmentIDs []string) error

	// GetAllUserIDsByRoleID 获取指定角色下所有用户ID（不分页）
	// 用于获取某个角色下所有用户的场景
	GetAllUserIDsByRoleID(ctx context.Context, roleID string) ([]string, error)

	// ReplaceRoleUsers 批量替换角色的用户绑定（事务操作）
	// 先删除该角色下所有旧的用户绑定，再创建新的用户绑定
	// 用于批量更新角色用户的场景，确保数据一致性
	ReplaceRoleUsers(
		ctx context.Context,
		roleID string,
		userIDs []string,
		operatorID string,
	) error
}
