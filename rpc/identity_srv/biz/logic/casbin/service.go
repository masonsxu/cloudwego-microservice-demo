package casbin

import (
	"context"

	"github.com/rs/zerolog"

	casbindal "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/casbin"
)

// Service Casbin 权限服务接口
type Service interface {
	// SyncPolicies 同步所有策略
	// 返回同步的策略统计信息
	SyncPolicies(ctx context.Context) (*SyncResult, error)

	// CheckPermission 检查权限
	// userID: 用户ID
	// roleIDs: 角色ID列表（多角色模式）
	// departmentIDs: 部门ID列表
	// resource: 资源标识
	// action: 操作类型
	CheckPermission(ctx context.Context, userID string, roleIDs []string, departmentIDs []string, resource, action string) (*CheckResult, error)

	// GetUserDataScope 获取用户在指定资源上的数据范围
	GetUserDataScope(ctx context.Context, userID string, resource, action string) (string, error)
}

// SyncResult 策略同步结果
type SyncResult struct {
	Success              bool
	RolePolicyCount      int32
	UserRoleBindingCount int32
	RoleInheritanceCount int32
	Message              string
}

// CheckResult 权限检查结果
type CheckResult struct {
	Allowed   bool
	DataScope string
}

// ServiceImpl Casbin 权限服务实现
type ServiceImpl struct {
	policySyncRepo casbindal.PolicySyncRepository
	logger         *zerolog.Logger
}

// NewService 创建 Casbin 服务
func NewService(policySyncRepo casbindal.PolicySyncRepository, logger *zerolog.Logger) Service {
	return &ServiceImpl{
		policySyncRepo: policySyncRepo,
		logger:         logger,
	}
}

// SyncPolicies 同步所有策略
func (s *ServiceImpl) SyncPolicies(ctx context.Context) (*SyncResult, error) {
	s.logger.Info().Msg("Starting policy sync")

	result := &SyncResult{
		Success: false,
		Message: "",
	}

	// 同步角色策略
	if err := s.policySyncRepo.SyncRolePolicies(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Failed to sync role policies")
		result.Message = "同步角色策略失败: " + err.Error()
		return result, err
	}

	// 获取策略数量
	enforcer := s.policySyncRepo.GetEnforcer()
	if enforcer != nil {
		policies, _ := enforcer.GetPolicy()
		result.RolePolicyCount = int32(len(policies))
	}

	// 同步用户角色绑定
	if err := s.policySyncRepo.SyncUserRoleBindings(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Failed to sync user role bindings")
		result.Message = "同步用户角色绑定失败: " + err.Error()
		return result, err
	}

	// 获取用户角色绑定数量
	if enforcer != nil {
		groupingPolicies, _ := enforcer.GetGroupingPolicy()
		result.UserRoleBindingCount = int32(len(groupingPolicies))
	}

	// 同步角色继承
	if err := s.policySyncRepo.SyncRoleInheritance(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Failed to sync role inheritance")
		result.Message = "同步角色继承失败: " + err.Error()
		return result, err
	}

	// 获取角色继承数量
	if enforcer != nil {
		g2Policies, _ := enforcer.GetNamedGroupingPolicy("g2")
		result.RoleInheritanceCount = int32(len(g2Policies))
	}

	result.Success = true
	result.Message = "策略同步成功"

	s.logger.Info().
		Int32("role_policy_count", result.RolePolicyCount).
		Int32("user_role_binding_count", result.UserRoleBindingCount).
		Int32("role_inheritance_count", result.RoleInheritanceCount).
		Msg("Policy sync completed successfully")

	return result, nil
}

// CheckPermission 检查权限
func (s *ServiceImpl) CheckPermission(
	ctx context.Context,
	userID string,
	roleIDs []string,
	departmentIDs []string,
	resource, action string,
) (*CheckResult, error) {
	result := &CheckResult{
		Allowed:   false,
		DataScope: "",
	}

	// 对于每个角色检查权限（多角色取并集）
	maxDataScope := ""

	for _, roleID := range roleIDs {
		roleSubject := "role:" + roleID

		// 检查全局权限
		allowed, err := s.policySyncRepo.CheckPermission(ctx, roleSubject, "*", resource, action)
		if err != nil {
			s.logger.Warn().Err(err).Str("role_id", roleID).Msg("Failed to check permission")
			continue
		}

		if allowed {
			result.Allowed = true
			// 获取数据范围
			dataScope, _ := s.policySyncRepo.GetUserDataScope(ctx, roleSubject, "*", resource, action)
			if compareDataScope(dataScope, maxDataScope) > 0 {
				maxDataScope = dataScope
			}
		}

		// 检查部门级权限
		for _, deptID := range departmentIDs {
			domain := "dept:" + deptID
			allowed, err := s.policySyncRepo.CheckPermission(ctx, roleSubject, domain, resource, action)
			if err != nil {
				continue
			}

			if allowed {
				result.Allowed = true
				dataScope, _ := s.policySyncRepo.GetUserDataScope(ctx, roleSubject, domain, resource, action)
				if compareDataScope(dataScope, maxDataScope) > 0 {
					maxDataScope = dataScope
				}
			}
		}
	}

	// 也检查用户直接权限
	userSubject := "user:" + userID
	for _, deptID := range departmentIDs {
		domain := "dept:" + deptID
		allowed, err := s.policySyncRepo.CheckPermission(ctx, userSubject, domain, resource, action)
		if err != nil {
			continue
		}

		if allowed {
			result.Allowed = true
			dataScope, _ := s.policySyncRepo.GetUserDataScope(ctx, userSubject, domain, resource, action)
			if compareDataScope(dataScope, maxDataScope) > 0 {
				maxDataScope = dataScope
			}
		}
	}

	result.DataScope = maxDataScope

	s.logger.Debug().
		Str("user_id", userID).
		Strs("role_ids", roleIDs).
		Str("resource", resource).
		Str("action", action).
		Bool("allowed", result.Allowed).
		Str("data_scope", result.DataScope).
		Msg("Permission check result")

	return result, nil
}

// GetUserDataScope 获取用户在指定资源上的数据范围
func (s *ServiceImpl) GetUserDataScope(ctx context.Context, userID string, resource, action string) (string, error) {
	userSubject := "user:" + userID
	dataScope, err := s.policySyncRepo.GetUserDataScope(ctx, userSubject, "*", resource, action)
	if err != nil {
		return "", err
	}
	return dataScope, nil
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
