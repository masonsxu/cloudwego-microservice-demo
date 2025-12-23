package dal

import (
	"context"

	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/assignment"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/definition"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/department"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/logo"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/membership"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/menu"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/organization"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/user"
)

// DALImpl DAL统一实现
// 聚合所有仓储实现，提供统一的数据访问服务
type DALImpl struct {
	// 数据库连接
	db *gorm.DB

	// 事务管理器
	txManager base.TransactionManager

	// 各实体仓储
	userProfileRepo        user.UserProfileRepository
	userMembershipRepo     membership.UserMembershipRepository
	organizationRepo       organization.OrganizationRepository
	departmentRepo         department.DepartmentRepository
	logoRepo               logo.LogoRepository
	menuRepo               menu.MenuRepository
	roleDefinitionRepo     definition.RoleDefinitionRepository
	userRoleAssignmentRepo assignment.UserRoleAssignmentRepository

	// 事务状态
	isTransaction bool
}

// newDALImpl 创建DAL实现实例（内部使用）
func newDALImpl(db *gorm.DB) DAL {
	return &DALImpl{
		db:                     db,
		txManager:              base.NewTransactionManager(db),
		userProfileRepo:        user.NewUserProfileRepository(db),
		userMembershipRepo:     membership.NewUserMembershipRepository(db),
		organizationRepo:       organization.NewOrganizationRepository(db),
		departmentRepo:         department.NewDepartmentRepository(db),
		logoRepo:               logo.NewLogoRepository(db),
		menuRepo:               menu.NewMenuRepository(db),
		roleDefinitionRepo:     definition.NewRoleDefinitionRepository(db),
		userRoleAssignmentRepo: assignment.NewUserRoleAssignmentRepository(db),
		isTransaction:          false,
	}
}

// ============================================================================
// 实体仓储接口实现
// ============================================================================

// UserProfile 获取用户档案仓储
func (dal *DALImpl) UserProfile() user.UserProfileRepository {
	return dal.userProfileRepo
}

// UserMembership 获取用户成员关系仓储
func (dal *DALImpl) UserMembership() membership.UserMembershipRepository {
	return dal.userMembershipRepo
}

// Organization 获取组织仓储
func (dal *DALImpl) Organization() organization.OrganizationRepository {
	return dal.organizationRepo
}

// Department 获取部门仓储
func (dal *DALImpl) Department() department.DepartmentRepository {
	return dal.departmentRepo
}

// Logo 获取组织Logo仓储
func (dal *DALImpl) Logo() logo.LogoRepository {
	return dal.logoRepo
}

// Menu 获取菜单仓储
func (dal *DALImpl) Menu() menu.MenuRepository {
	return dal.menuRepo
}

// RoleDefinition 获取角色定义仓储
func (dal *DALImpl) RoleDefinition() definition.RoleDefinitionRepository {
	return dal.roleDefinitionRepo
}

// UserRoleAssignment 获取用户角色分配仓储
func (dal *DALImpl) UserRoleAssignment() assignment.UserRoleAssignmentRepository {
	return dal.userRoleAssignmentRepo
}

// ============================================================================
// 事务管理实现
// ============================================================================

// WithTransaction 在事务中执行操作（推荐使用）
func (dal *DALImpl) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context, dal DAL) error,
) error {
	return dal.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// 创建事务DAL实例
		txDAL := dal.WithDB(tx).(*DALImpl)
		txDAL.isTransaction = true

		// 执行业务逻辑
		return fn(ctx, txDAL)
	})
}

// BeginTx 开始事务
func (dal *DALImpl) BeginTx(ctx context.Context) (DAL, error) {
	tx, err := dal.txManager.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	txDAL := dal.WithDB(tx).(*DALImpl)
	txDAL.isTransaction = true

	return txDAL, nil
}

// Commit 提交事务
func (dal *DALImpl) Commit() error {
	if !dal.isTransaction {
		return nil // 非事务状态，忽略提交操作
	}

	return dal.txManager.CommitTx(dal.db)
}

// Rollback 回滚事务
func (dal *DALImpl) Rollback() error {
	if !dal.isTransaction {
		return nil // 非事务状态，忽略回滚操作
	}

	return dal.txManager.RollbackTx(dal.db)
}

// ============================================================================
// 数据库连接管理实现
// ============================================================================

// DB 获取数据库连接
func (dal *DALImpl) DB() *gorm.DB {
	return dal.db
}

// WithDB 使用指定数据库连接创建新的DAL实例
func (dal *DALImpl) WithDB(db *gorm.DB) DAL {
	return &DALImpl{
		db:                     db,
		txManager:              base.NewTransactionManager(db),
		userProfileRepo:        user.NewUserProfileRepository(db),
		userMembershipRepo:     membership.NewUserMembershipRepository(db),
		organizationRepo:       organization.NewOrganizationRepository(db),
		departmentRepo:         department.NewDepartmentRepository(db),
		logoRepo:               logo.NewLogoRepository(db),
		menuRepo:               menu.NewMenuRepository(db),
		roleDefinitionRepo:     definition.NewRoleDefinitionRepository(db),
		userRoleAssignmentRepo: assignment.NewUserRoleAssignmentRepository(db),
		isTransaction:          dal.isTransaction,
	}
}

// NewDALImpl 创建DAL实现实例
func NewDALImpl(db *gorm.DB) DAL {
	return newDALImpl(db)
}
