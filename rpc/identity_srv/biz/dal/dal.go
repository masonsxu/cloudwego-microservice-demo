package dal

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/assignment"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/definition"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/department"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/logo"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/membership"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/menu"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/organization"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/user"
	"gorm.io/gorm"
)

// DAL 数据访问层统一接口
// 基于重构后的架构设计，聚合所有实体的仓储接口，提供统一的数据访问入口
// 支持事务管理、查询优化和完整的CRUD操作
type DAL interface {
	// ============================================================================
	// 实体仓储接口
	// ============================================================================

	// UserProfile 用户档案仓储
	UserProfile() user.UserProfileRepository

	// UserMembership 用户成员关系仓储
	UserMembership() membership.UserMembershipRepository

	// Organization 组织仓储
	Organization() organization.OrganizationRepository

	// Department 部门仓储
	Department() department.DepartmentRepository

	// Logo 组织Logo仓储
	Logo() logo.LogoRepository

	// Menu 菜单仓储
	Menu() menu.MenuRepository

	// RoleDefinition 角色定义仓储
	RoleDefinition() definition.RoleDefinitionRepository

	// UserRoleAssignment 用户角色分配仓储
	UserRoleAssignment() assignment.UserRoleAssignmentRepository

	// ============================================================================
	// 事务管理
	// ============================================================================

	// WithTransaction 在事务中执行操作（推荐使用）
	WithTransaction(ctx context.Context, fn func(ctx context.Context, dal DAL) error) error

	// BeginTx 开始事务（返回新的DAL实例）
	BeginTx(ctx context.Context) (DAL, error)

	// Commit 提交事务
	Commit() error

	// Rollback 回滚事务
	Rollback() error

	// ============================================================================
	// 数据库连接管理
	// ============================================================================

	// DB 获取数据库连接（用于复杂查询）
	DB() *gorm.DB

	// WithDB 使用指定数据库连接创建新的DAL实例
	WithDB(db *gorm.DB) DAL
}

// NewDAL 创建DAL实例
func NewDAL(db *gorm.DB) DAL {
	return newDALImpl(db)
}
