package logic

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

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

// 合法的 ptype 取值
const (
	PTypePolicy           = "p"
	PTypeGroupingPolicy   = "g"
	PTypeRoleInheritance  = "g2"
	defaultReloadInterval = 30 * time.Second
)

// ErrInvalidPType ptype 取值非法
var ErrInvalidPType = errors.New("invalid ptype, must be one of: p, g, g2")

// EnforcerService Casbin enforcer 封装
type EnforcerService struct {
	enforcer *casbin.SyncedEnforcer
	db       *gorm.DB
	logger   *zerolog.Logger
	mu       sync.RWMutex

	stopCh   chan struct{}
	stopOnce sync.Once
	interval time.Duration
}

// NewEnforcerService 创建 enforcer 服务（生产路径，使用 GORM adapter）
func NewEnforcerService(db *gorm.DB, logger *zerolog.Logger) (*EnforcerService, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("创建 GORM adapter 失败: %w", err)
	}

	return newEnforcerWithAdapter(db, logger, adapter)
}

// newEnforcerWithAdapter 接受任意 adapter 构造 enforcer，便于单测注入内存 adapter
func newEnforcerWithAdapter(db *gorm.DB, logger *zerolog.Logger, adapter any) (*EnforcerService, error) {
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
		stopCh:   make(chan struct{}),
		interval: defaultReloadInterval,
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

// AddPolicyRule 按 ptype 增加一条规则；ptype ∈ {p, g, g2}
func (s *EnforcerService) AddPolicyRule(ptype string, rule []string) (bool, error) {
	if err := validatePType(ptype); err != nil {
		return false, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	switch ptype {
	case PTypePolicy:
		return s.enforcer.AddPolicy(toAny(rule)...)
	case PTypeGroupingPolicy:
		return s.enforcer.AddGroupingPolicy(toAny(rule)...)
	case PTypeRoleInheritance:
		return s.enforcer.AddNamedGroupingPolicy("g2", toAny(rule)...)
	}
	return false, ErrInvalidPType
}

// RemovePolicyRule 按 ptype 删除一条规则
func (s *EnforcerService) RemovePolicyRule(ptype string, rule []string) (bool, error) {
	if err := validatePType(ptype); err != nil {
		return false, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	switch ptype {
	case PTypePolicy:
		return s.enforcer.RemovePolicy(toAny(rule)...)
	case PTypeGroupingPolicy:
		return s.enforcer.RemoveGroupingPolicy(toAny(rule)...)
	case PTypeRoleInheritance:
		return s.enforcer.RemoveNamedGroupingPolicy("g2", toAny(rule)...)
	}
	return false, ErrInvalidPType
}

// StartAutoReload 启动后台轮询，从 DB 重载策略
// TODO(phase 2.1): 切换为 etcd watch 推送策略变更，移除轮询
func (s *EnforcerService) StartAutoReload(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		s.logger.Info().Dur("interval", s.interval).Msg("策略自动同步已启动")

		for {
			select {
			case <-ctx.Done():
				return
			case <-s.stopCh:
				return
			case <-ticker.C:
				if err := s.ReloadPolicy(ctx); err != nil {
					s.logger.Error().Err(err).Msg("策略自动同步失败")
				}
			}
		}
	}()
}

// Stop 停止后台同步
func (s *EnforcerService) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopCh)
		s.logger.Info().Msg("策略自动同步已停止")
	})
}

// enforcerEx 扩展的 Enforce（含 data_scope 决策）
func (s *EnforcerService) enforcerEx(sub, dom, obj, act string) (bool, string, error) {
	policies, err := s.enforcer.GetFilteredPolicy(0, "", dom, obj, act)
	if err != nil {
		return false, "", err
	}

	wildcardPolicies, err := s.enforcer.GetFilteredPolicy(0, "", dom, obj, "*")
	if err != nil {
		return false, "", err
	}

	allPolicies := append(policies, wildcardPolicies...)
	if len(allPolicies) == 0 {
		return false, "", nil
	}

	roles := []string{sub}
	roles = append(roles, s.enforcer.GetRolesForUserInDomain(sub, dom)...)
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

func validatePType(ptype string) error {
	switch ptype {
	case PTypePolicy, PTypeGroupingPolicy, PTypeRoleInheritance:
		return nil
	default:
		return ErrInvalidPType
	}
}

func toAny(in []string) []any {
	out := make([]any, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}
