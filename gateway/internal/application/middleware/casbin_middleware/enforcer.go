package casbin_middleware

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/rs/zerolog"
)

// CasbinEnforcer Casbin Enforcer 封装（使用内存策略）
// 策略数据从 RPC 服务同步，不直接访问数据库
type CasbinEnforcer struct {
	enforcer *casbin.SyncedEnforcer
	logger   *zerolog.Logger
	mu       sync.RWMutex
}

// 内嵌模型定义
const defaultModelText = `
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

// NewCasbinEnforcer 创建新的 Casbin Enforcer（使用内存 Adapter）
// 策略从 RPC 服务同步，不依赖 GORM
func NewCasbinEnforcer(logger *zerolog.Logger) (*CasbinEnforcer, error) {
	m, err := model.NewModelFromString(defaultModelText)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin model: %w", err)
	}

	// 使用内存 Adapter（空策略，后续从 RPC 同步）
	enforcer, err := casbin.NewSyncedEnforcer(m)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// 启用自动保存（内存模式下无效，但保持接口一致）
	enforcer.EnableAutoSave(false)

	return &CasbinEnforcer{
		enforcer: enforcer,
		logger:   logger,
	}, nil
}

// NewCasbinEnforcerFromFile 从文件加载模型创建 Enforcer
func NewCasbinEnforcerFromFile(modelPath string, logger *zerolog.Logger) (*CasbinEnforcer, error) {
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load casbin model from %s: %w", modelPath, err)
	}

	enforcer, err := casbin.NewSyncedEnforcer(m)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	enforcer.EnableAutoSave(false)

	return &CasbinEnforcer{
		enforcer: enforcer,
		logger:   logger,
	}, nil
}

// Enforce 执行权限校验
func (e *CasbinEnforcer) Enforce(sub, dom, obj, act string) (bool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	allowed, err := e.enforcer.Enforce(sub, dom, obj, act)
	if err != nil {
		e.logger.Error().Err(err).
			Str("sub", sub).
			Str("dom", dom).
			Str("obj", obj).
			Str("act", act).
			Msg("Casbin enforce error")
		return false, err
	}

	e.logger.Debug().
		Str("sub", sub).
		Str("dom", dom).
		Str("obj", obj).
		Str("act", act).
		Bool("allowed", allowed).
		Msg("Casbin enforce result")

	return allowed, nil
}

// EnforceWithDataScope 执行权限校验并返回数据范围
func (e *CasbinEnforcer) EnforceWithDataScope(sub, dom, obj, act string) (bool, string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// 获取所有匹配的策略
	policies, err := e.enforcer.GetFilteredPolicy(0, sub, dom, obj, act)
	if err != nil {
		return false, "", err
	}
	if len(policies) == 0 {
		// 尝试使用通配符域匹配
		policies, err = e.enforcer.GetFilteredPolicy(0, sub, "*", obj, act)
		if err != nil {
			return false, "", err
		}
	}

	if len(policies) == 0 {
		return false, "", nil
	}

	// 取最大数据范围（org > dept > self）
	maxScope := ""
	for _, policy := range policies {
		if len(policy) >= 5 {
			scope := policy[4]
			if compareDataScope(scope, maxScope) > 0 {
				maxScope = scope
			}
		}
	}

	return true, maxScope, nil
}

// compareDataScope 比较数据范围大小
// 返回: 1 表示 a > b, -1 表示 a < b, 0 表示相等
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

// AddPolicy 添加策略
func (e *CasbinEnforcer) AddPolicy(sub, dom, obj, act, dataScope string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	added, err := e.enforcer.AddPolicy(sub, dom, obj, act, dataScope)
	if err != nil {
		e.logger.Error().Err(err).
			Str("sub", sub).
			Str("dom", dom).
			Str("obj", obj).
			Str("act", act).
			Str("data_scope", dataScope).
			Msg("Failed to add policy")
		return false, err
	}

	return added, nil
}

// AddPolicies 批量添加策略
func (e *CasbinEnforcer) AddPolicies(policies [][]string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	added, err := e.enforcer.AddPolicies(policies)
	if err != nil {
		e.logger.Error().Err(err).
			Int("count", len(policies)).
			Msg("Failed to add policies")
		return false, err
	}

	e.logger.Info().
		Int("count", len(policies)).
		Bool("added", added).
		Msg("Policies added")

	return added, nil
}

// RemovePolicy 移除策略
func (e *CasbinEnforcer) RemovePolicy(sub, dom, obj, act, dataScope string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	removed, err := e.enforcer.RemovePolicy(sub, dom, obj, act, dataScope)
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to remove policy")
		return false, err
	}

	return removed, nil
}

// ClearPolicy 清空所有策略
func (e *CasbinEnforcer) ClearPolicy() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.enforcer.ClearPolicy()
	e.logger.Info().Msg("All policies cleared")
}

// AddRoleForUserInDomain 添加用户-角色-域绑定
func (e *CasbinEnforcer) AddRoleForUserInDomain(user, role, domain string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	added, err := e.enforcer.AddRoleForUserInDomain(user, role, domain)
	if err != nil {
		e.logger.Error().Err(err).
			Str("user", user).
			Str("role", role).
			Str("domain", domain).
			Msg("Failed to add role for user in domain")
		return false, err
	}

	return added, nil
}

// AddGroupingPolicies 批量添加用户-角色绑定
func (e *CasbinEnforcer) AddGroupingPolicies(policies [][]string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	added, err := e.enforcer.AddGroupingPolicies(policies)
	if err != nil {
		e.logger.Error().Err(err).
			Int("count", len(policies)).
			Msg("Failed to add grouping policies")
		return false, err
	}

	e.logger.Info().
		Int("count", len(policies)).
		Bool("added", added).
		Msg("Grouping policies added")

	return added, nil
}

// RemoveRoleForUserInDomain 移除用户-角色-域绑定
func (e *CasbinEnforcer) RemoveRoleForUserInDomain(user, role, domain string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	removed, err := e.enforcer.DeleteRoleForUserInDomain(user, role, domain)
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to remove role for user in domain")
		return false, err
	}

	return removed, nil
}

// AddRoleInheritance 添加角色继承关系（使用 g2）
func (e *CasbinEnforcer) AddRoleInheritance(childRole, parentRole string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	added, err := e.enforcer.AddNamedGroupingPolicy("g2", childRole, parentRole)
	if err != nil {
		e.logger.Error().Err(err).
			Str("child_role", childRole).
			Str("parent_role", parentRole).
			Msg("Failed to add role inheritance")
		return false, err
	}

	return added, nil
}

// AddNamedGroupingPolicies 批量添加角色继承关系
func (e *CasbinEnforcer) AddNamedGroupingPolicies(ptype string, policies [][]string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	added, err := e.enforcer.AddNamedGroupingPolicies(ptype, policies)
	if err != nil {
		e.logger.Error().Err(err).
			Str("ptype", ptype).
			Int("count", len(policies)).
			Msg("Failed to add named grouping policies")
		return false, err
	}

	e.logger.Info().
		Str("ptype", ptype).
		Int("count", len(policies)).
		Bool("added", added).
		Msg("Named grouping policies added")

	return added, nil
}

// RemoveRoleInheritance 移除角色继承关系
func (e *CasbinEnforcer) RemoveRoleInheritance(childRole, parentRole string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	removed, err := e.enforcer.RemoveNamedGroupingPolicy("g2", childRole, parentRole)
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to remove role inheritance")
		return false, err
	}

	return removed, nil
}

// GetRolesForUserInDomain 获取用户在指定域下的角色
func (e *CasbinEnforcer) GetRolesForUserInDomain(user, domain string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.enforcer.GetRolesForUserInDomain(user, domain)
}

// GetPermissionsForUser 获取用户的所有权限
func (e *CasbinEnforcer) GetPermissionsForUser(user string) [][]string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	perms, _ := e.enforcer.GetPermissionsForUser(user)
	return perms
}

// GetAllPolicies 获取所有策略
func (e *CasbinEnforcer) GetAllPolicies() [][]string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	policies, _ := e.enforcer.GetPolicy()
	return policies
}

// GetAllGroupingPolicies 获取所有用户-角色绑定
func (e *CasbinEnforcer) GetAllGroupingPolicies() [][]string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	policies, _ := e.enforcer.GetGroupingPolicy()
	return policies
}

// GetPolicyCount 获取策略数量
func (e *CasbinEnforcer) GetPolicyCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()

	policies, _ := e.enforcer.GetPolicy()
	return len(policies)
}

// GetEnforcer 获取底层 enforcer（用于高级操作）
func (e *CasbinEnforcer) GetEnforcer() *casbin.SyncedEnforcer {
	return e.enforcer
}
