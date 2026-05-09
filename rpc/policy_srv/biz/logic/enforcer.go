package logic

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
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

// SeedDefaultsIfEmpty 仅当策略表完全为空时写入默认策略：
//
//	p, role:superadmin, *, *, *, all
//
// 这条规则让 role code = role:superadmin 的主体在任意 tenant 对任意 resource 任意 action
// 都获得 data_scope=all。配合 identity_srv seeder 中已建立的 superadmin 角色
// （RoleCode 由 model.GenerateRoleCode 生成为 "role:superadmin"），
// 即可在首次启动时打通 PDP 决策链路（提案 §14 Phase 4b 后的前置遗漏 #2）。
//
// 行为说明：
//   - 仅在所有策略表（p / g / g2）都为空时写入，避免污染已有数据；
//   - 通过 SyncedEnforcer.AddPolicy 写入，会自动持久化到 GORM adapter；
//   - 单测路径（newEnforcerWithAdapter + 文件 adapter）不会调用本方法，
//     测试逻辑保持原样。
func (s *EnforcerService) SeedDefaultsIfEmpty(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pCount, gCount, g2Count := s.policyCountsLocked()
	if pCount > 0 || gCount > 0 || g2Count > 0 {
		s.logger.Debug().
			Int("p", pCount).
			Int("g", gCount).
			Int("g2", g2Count).
			Msg("策略表非空，跳过种子写入")

		return nil
	}

	rule := []any{"role:superadmin", "*", "*", "*", "all"}

	added, err := s.enforcer.AddPolicy(rule...)
	if err != nil {
		return fmt.Errorf("写入 superadmin 默认策略失败: %w", err)
	}

	s.logger.Info().
		Bool("added", added).
		Strs("rule", []string{"role:superadmin", "*", "*", "*", "all"}).
		Msg("已写入 superadmin 默认通配策略")

	return nil
}

// policyCountsLocked 调用方需持有 mu，返回 p/g/g2 三类策略的数量。
func (s *EnforcerService) policyCountsLocked() (int, int, int) {
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
//
// 实现说明：Casbin 内置 matcher 支持 keyMatch2 / dom 通配，但只能返回 bool；
// data_scope 需要"哪条策略命中、它的第 5 列是什么"——必须直接查策略表。
// 因此用 GetFilteredPolicy 8 种通配组合（dom × obj × act 各 ±通配）穷举，
// 然后按"sub 匹配（直接或通过 g 关系）"筛选，命中即取最大数据范围。
func (s *EnforcerService) enforcerEx(sub, dom, obj, act string) (bool, string, error) {
	allPolicies, err := s.collectMatchingPolicies(dom, obj, act)
	if err != nil {
		return false, "", err
	}

	if len(allPolicies) == 0 {
		return false, "", nil
	}

	// 主体集合：sub 自身 + 通过 g 关系映射到的角色
	// 容量预估为 sub 自身 + 单 dom 与全 dom 各取 4 个角色，避免线性扩容（prealloc）
	roles := make([]string, 0, 1+8)
	roles = append(roles, sub)
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

// collectMatchingPolicies 返回所有 dom/obj/act 三列与请求兼容（精确相等或为 "*"）的策略。
//
// 实现说明：早期版本调用 GetFilteredPolicy(0, "", d, o, a) 8 种组合，但 Casbin v3
// 的 GetFilteredPolicy 对 "*" 是字面比对（要求策略列也精确为 "*"），且组合容易掉
// 项。改成 GetPolicy() 拿全部策略后单次内存扫描，逻辑简单且与 Casbin model 中
// `r.dom == p.dom || p.dom == "*"` 等价。
func (s *EnforcerService) collectMatchingPolicies(dom, obj, act string) ([][]string, error) {
	all, err := s.enforcer.GetPolicy()
	if err != nil {
		return nil, err
	}

	out := make([][]string, 0, len(all))

	for _, p := range all {
		if len(p) < 5 {
			continue
		}

		if matchField(p[1], dom) && matchField(p[2], obj) && matchField(p[3], act) {
			out = append(out, p)
		}
	}

	return out, nil
}

// matchField 返回策略列值是否与请求列值兼容：精确相等或策略列为 "*"。
func matchField(policyVal, requestVal string) bool {
	return policyVal == requestVal || policyVal == "*"
}

// compareDataScope 比较数据范围大小 (self < dept < org < all)
//
// "all" 用于 superadmin 通配策略（提案 §14 Phase 4b 后置遗漏修复），
// 表示任意 tenant、任意资源都允许；任何具体 scope 都不应覆盖它。
func compareDataScope(a, b string) int {
	scopeOrder := map[string]int{"self": 1, "dept": 2, "org": 3, "all": 4}
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
