package logic

import (
	"context"
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// casbinModelText 内嵌的 Casbin RBAC with Domains 模型定义
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
m = (g(r.sub, p.sub, r.dom) || g(r.sub, p.sub, "*")) && \
    (r.dom == p.dom || p.dom == "*") && \
    keyMatch2(r.obj, p.obj) && \
    (r.act == p.act || p.act == "*")
`

// EnforcerService Casbin enforcer 封装
type EnforcerService struct {
	enforcer *casbin.SyncedEnforcer
	db       *gorm.DB
	logger   *zerolog.Logger
	mu       sync.RWMutex
}

// NewEnforcerService 创建 enforcer 服务
func NewEnforcerService(db *gorm.DB, logger *zerolog.Logger) (*EnforcerService, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("创建 GORM adapter 失败: %w", err)
	}

	m, err := model.NewModelFromString(casbinModelText)
	if err != nil {
		return nil, fmt.Errorf("创建 Casbin model 失败: %w", err)
	}

	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("创建 SyncedEnforcer 失败: %w", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("加载策略失败: %w", err)
	}

	logger.Info().Msg("Casbin enforcer 初始化完成")

	return &EnforcerService{
		enforcer: enforcer,
		db:       db,
		logger:   logger,
	}, nil
}

// Enforce 单点权限检查
func (s *EnforcerService) Enforce(sub, dom, obj, act string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.enforcer.Enforce(sub, dom, obj, act)
}

// EnforceWithDataScope 权限检查并返回数据范围
func (s *EnforcerService) EnforceWithDataScope(sub, dom, obj, act string) (bool, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.enforcerEx(sub, dom, obj, act)
}

// GetRolesForUserInDomain 获取用户在域中的角色
func (s *EnforcerService) GetRolesForUserInDomain(userID, domain string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.enforcer.GetRolesForUserInDomain(userID, domain)
}

// GetPermissionsForUser 获取用户所有权限
func (s *EnforcerService) GetPermissionsForUser(userID string, domains ...string) ([][]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.enforcer.GetPermissionsForUser(userID, domains...)
}

// ReloadPolicy 从 DB 重新加载策略
func (s *EnforcerService) ReloadPolicy(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("重新加载策略失败: %w", err)
	}

	s.logger.Info().Msg("策略重新加载完成")
	return nil
}

// GetPolicyCount 获取策略统计
func (s *EnforcerService) GetPolicyCount() (int, int, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	policies, _ := s.enforcer.GetPolicy()
	grouping, _ := s.enforcer.GetGroupingPolicy()
	inheritance, _ := s.enforcer.GetNamedGroupingPolicy("g2")

	return len(policies), len(grouping), len(inheritance)
}

// enforcerEx 扩展的 Enforce（含 data_scope 决策）
func (s *EnforcerService) enforcerEx(sub, dom, obj, act string) (bool, string, error) {
	// 尝试精确 action 匹配
	policies, err := s.enforcer.GetFilteredPolicy(0, "", dom, obj, act)
	if err != nil {
		return false, "", err
	}

	// 尝试通配 action 匹配
	wildcardPolicies, err := s.enforcer.GetFilteredPolicy(0, "", dom, obj, "*")
	if err != nil {
		return false, "", err
	}

	allPolicies := append(policies, wildcardPolicies...)
	if len(allPolicies) == 0 {
		return false, "", nil
	}

	// 检查用户是否有匹配的角色
	roles := s.enforcer.GetRolesForUserInDomain(sub, dom)
	roles = append(roles, s.enforcer.GetRolesForUserInDomain(sub, "*")...)

	if len(roles) == 0 {
		return false, "", nil
	}

	maxScope := ""
	allowed := false

	for _, policy := range allPolicies {
		if len(policy) < 5 {
			continue
		}

		policySub := policy[0]
		for _, role := range roles {
			if role == policySub {
				allowed = true
				scope := policy[4]
				if compareDataScope(scope, maxScope) > 0 {
					maxScope = scope
				}
				break
			}
		}
	}

	return allowed, maxScope, nil
}

// compareDataScope 比较数据范围大小 (self < dept < org)
func compareDataScope(a, b string) int {
	scopeOrder := map[string]int{"self": 1, "dept": 2, "org": 3}
	orderA := scopeOrder[a]
	orderB := scopeOrder[b]

	if orderA > orderB {
		return 1
	} else if orderA < orderB {
		return -1
	}
	return 0
}
