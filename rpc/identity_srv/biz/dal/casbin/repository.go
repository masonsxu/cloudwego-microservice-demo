package casbin

import (
	"context"
	"fmt"
	"sync"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// PolicySyncRepositoryImpl Casbin 策略同步仓储实现
type PolicySyncRepositoryImpl struct {
	db       *gorm.DB
	enforcer *casbin.SyncedEnforcer
	logger   *zerolog.Logger
	mu       sync.RWMutex
}

// casbinModelText 内嵌的 Casbin 模型定义
const casbinModelText = `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act, data_scope

[role_definition]
g = _, _, _
g2 = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (g(r.sub, p.sub, r.dom) || g(r.sub, p.sub, "*")) && (r.dom == p.dom || p.dom == "*") && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`

// NewPolicySyncRepository 创建新的策略同步仓储
func NewPolicySyncRepository(db *gorm.DB, logger *zerolog.Logger) (PolicySyncRepository, error) {
	// 创建 GORM adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm adapter: %w", err)
	}

	// 使用内嵌模型
	m, err := model.NewModelFromString(casbinModelText)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin model: %w", err)
	}

	// 创建同步 enforcer
	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	return &PolicySyncRepositoryImpl{
		db:       db,
		enforcer: enforcer,
		logger:   logger,
	}, nil
}

// GetEnforcer 获取 Casbin Enforcer 实例
func (r *PolicySyncRepositoryImpl) GetEnforcer() *casbin.SyncedEnforcer {
	return r.enforcer
}

// SyncRolePolicies 同步角色权限策略
func (r *PolicySyncRepositoryImpl) SyncRolePolicies(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 获取所有角色菜单权限
	var roleMenuPerms []models.RoleMenuPermission
	if err := r.db.WithContext(ctx).Preload("Role").Find(&roleMenuPerms).Error; err != nil {
		return fmt.Errorf("failed to query role menu permissions: %w", err)
	}

	// 清除现有策略
	r.enforcer.ClearPolicy()

	// 生成新策略
	for _, rmp := range roleMenuPerms {
		if rmp.Role == nil {
			continue
		}

		roleCode := rmp.Role.GetCasbinSubject()
		domain := rmp.Role.GetCasbinDomain()
		resource := "menu:" + rmp.MenuID
		action := permissionTypeToAction(rmp.PermissionType)
		dataScope := dataScopeToString(rmp.DataScope)

		if _, err := r.enforcer.AddPolicy(roleCode, domain, resource, action, dataScope); err != nil {
			r.logger.Warn().
				Err(err).
				Str("role_code", roleCode).
				Str("menu_id", rmp.MenuID).
				Msg("Failed to add policy")
		}
	}

	// 保存策略
	if err := r.enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("failed to save policy: %w", err)
	}

	r.logger.Info().Int("count", len(roleMenuPerms)).Msg("Role policies synced")
	return nil
}

// SyncUserRoleBindings 同步用户角色绑定
func (r *PolicySyncRepositoryImpl) SyncUserRoleBindings(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 获取所有用户角色分配
	var assignments []models.UserRoleAssignment
	if err := r.db.WithContext(ctx).Find(&assignments).Error; err != nil {
		return fmt.Errorf("failed to query user role assignments: %w", err)
	}

	// 获取角色信息
	roleMap := make(map[uuid.UUID]*models.RoleDefinition)
	var roles []models.RoleDefinition
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return fmt.Errorf("failed to query role definitions: %w", err)
	}
	for i := range roles {
		roleMap[roles[i].ID] = &roles[i]
	}

	// 生成用户角色绑定
	for _, assignment := range assignments {
		role, ok := roleMap[assignment.RoleID]
		if !ok {
			continue
		}

		userID := "user:" + assignment.UserID.String()
		roleCode := role.GetCasbinSubject()
		domain := role.GetCasbinDomain()

		if _, err := r.enforcer.AddRoleForUserInDomain(userID, roleCode, domain); err != nil {
			r.logger.Warn().
				Err(err).
				Str("user_id", userID).
				Str("role_code", roleCode).
				Msg("Failed to add role for user")
		}
	}

	// 保存策略
	if err := r.enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("failed to save policy: %w", err)
	}

	r.logger.Info().Int("count", len(assignments)).Msg("User role bindings synced")
	return nil
}

// SyncRoleInheritance 同步角色继承关系
func (r *PolicySyncRepositoryImpl) SyncRoleInheritance(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 获取所有有父角色的角色定义
	var roles []models.RoleDefinition
	if err := r.db.WithContext(ctx).Where("parent_role_id IS NOT NULL").Preload("ParentRole").Find(&roles).Error; err != nil {
		return fmt.Errorf("failed to query role definitions with parent: %w", err)
	}

	// 生成角色继承关系
	for _, role := range roles {
		if role.ParentRole == nil {
			continue
		}

		childRole := role.GetCasbinSubject()
		parentRole := role.ParentRole.GetCasbinSubject()

		if _, err := r.enforcer.AddNamedGroupingPolicy("g2", childRole, parentRole); err != nil {
			r.logger.Warn().
				Err(err).
				Str("child_role", childRole).
				Str("parent_role", parentRole).
				Msg("Failed to add role inheritance")
		}
	}

	// 保存策略
	if err := r.enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("failed to save policy: %w", err)
	}

	r.logger.Info().Int("count", len(roles)).Msg("Role inheritance synced")
	return nil
}

// SyncAll 同步所有策略数据
func (r *PolicySyncRepositoryImpl) SyncAll(ctx context.Context) error {
	if err := r.SyncRolePolicies(ctx); err != nil {
		return err
	}
	if err := r.SyncUserRoleBindings(ctx); err != nil {
		return err
	}
	if err := r.SyncRoleInheritance(ctx); err != nil {
		return err
	}
	return nil
}

// AddRolePolicy 添加角色策略
func (r *PolicySyncRepositoryImpl) AddRolePolicy(ctx context.Context, roleCode, domain, resource, action, dataScope string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.enforcer.AddPolicy(roleCode, domain, resource, action, dataScope)
	if err != nil {
		return fmt.Errorf("failed to add role policy: %w", err)
	}
	return r.enforcer.SavePolicy()
}

// RemoveRolePolicy 移除角色策略
func (r *PolicySyncRepositoryImpl) RemoveRolePolicy(ctx context.Context, roleCode, domain, resource, action, dataScope string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.enforcer.RemovePolicy(roleCode, domain, resource, action, dataScope)
	if err != nil {
		return fmt.Errorf("failed to remove role policy: %w", err)
	}
	return r.enforcer.SavePolicy()
}

// AddUserRoleBinding 添加用户角色绑定
func (r *PolicySyncRepositoryImpl) AddUserRoleBinding(ctx context.Context, userID, roleCode, domain string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.enforcer.AddRoleForUserInDomain(userID, roleCode, domain)
	if err != nil {
		return fmt.Errorf("failed to add user role binding: %w", err)
	}
	return r.enforcer.SavePolicy()
}

// RemoveUserRoleBinding 移除用户角色绑定
func (r *PolicySyncRepositoryImpl) RemoveUserRoleBinding(ctx context.Context, userID, roleCode, domain string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.enforcer.DeleteRoleForUserInDomain(userID, roleCode, domain)
	if err != nil {
		return fmt.Errorf("failed to remove user role binding: %w", err)
	}
	return r.enforcer.SavePolicy()
}

// AddRoleInheritance 添加角色继承
func (r *PolicySyncRepositoryImpl) AddRoleInheritance(ctx context.Context, childRole, parentRole string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.enforcer.AddNamedGroupingPolicy("g2", childRole, parentRole)
	if err != nil {
		return fmt.Errorf("failed to add role inheritance: %w", err)
	}
	return r.enforcer.SavePolicy()
}

// RemoveRoleInheritance 移除角色继承
func (r *PolicySyncRepositoryImpl) RemoveRoleInheritance(ctx context.Context, childRole, parentRole string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.enforcer.RemoveNamedGroupingPolicy("g2", childRole, parentRole)
	if err != nil {
		return fmt.Errorf("failed to remove role inheritance: %w", err)
	}
	return r.enforcer.SavePolicy()
}

// CheckPermission 检查权限
func (r *PolicySyncRepositoryImpl) CheckPermission(ctx context.Context, userID, domain, resource, action string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	allowed, err := r.enforcer.Enforce(userID, domain, resource, action)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}
	return allowed, nil
}

// GetUserDataScope 获取用户在指定资源上的数据范围
func (r *PolicySyncRepositoryImpl) GetUserDataScope(ctx context.Context, userID, domain, resource, action string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 获取用户的所有角色
	roles := r.enforcer.GetRolesForUserInDomain(userID, domain)

	maxScope := ""
	for _, role := range roles {
		// 获取角色的策略
		policies, err := r.enforcer.GetFilteredPolicy(0, role, domain, resource, action)
		if err != nil {
			continue
		}

		for _, policy := range policies {
			if len(policy) >= 5 {
				scope := policy[4]
				if compareDataScope(scope, maxScope) > 0 {
					maxScope = scope
				}
			}
		}
	}

	return maxScope, nil
}

// permissionTypeToAction 将权限类型转换为操作
func permissionTypeToAction(pt models.MenuPermissionType) string {
	switch pt {
	case models.PermissionView:
		return "read"
	case models.PermissionEdit:
		return "write"
	case models.PermissionManage:
		return "manage"
	case models.PermissionFull:
		return "*"
	default:
		return ""
	}
}

// dataScopeToString 将数据范围转换为字符串
func dataScopeToString(ds models.DataScope) string {
	switch ds {
	case models.DataScopeOwnOrg:
		return "dept"
	case models.DataScopeAllOrgs:
		return "org"
	default:
		return "self"
	}
}

// compareDataScope 比较数据范围大小
func compareDataScope(a, b string) int {
	scopeOrder := map[string]int{
		"self": 1,
		"dept": 2,
		"org":  3,
	}

	orderA := scopeOrder[a]
	orderB := scopeOrder[b]

	if orderA > orderB {
		return 1
	} else if orderA < orderB {
		return -1
	}
	return 0
}
