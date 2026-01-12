package casbin_middleware

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"

	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
)

// PolicySyncService 策略同步服务
// 负责从 RPC 服务同步策略到本地 Casbin Enforcer
type PolicySyncService struct {
	enforcer       *CasbinEnforcer
	identityClient identitycli.IdentityClient
	logger         *zerolog.Logger
	syncInterval   time.Duration
	stopCh         chan struct{}
	mu             sync.Mutex
	lastSyncTime   time.Time
	syncCount      int64
}

// NewPolicySyncService 创建策略同步服务
func NewPolicySyncService(
	enforcer *CasbinEnforcer,
	identityClient identitycli.IdentityClient,
	logger *zerolog.Logger,
	syncIntervalSeconds int,
) *PolicySyncService {
	interval := time.Duration(syncIntervalSeconds) * time.Second
	if interval <= 0 {
		interval = 5 * time.Minute // 默认5分钟
	}

	return &PolicySyncService{
		enforcer:       enforcer,
		identityClient: identityClient,
		logger:         logger,
		syncInterval:   interval,
		stopCh:         make(chan struct{}),
	}
}

// Start 启动策略同步服务
func (s *PolicySyncService) Start(ctx context.Context) error {
	// 启动时立即同步一次
	if err := s.SyncPolicies(ctx); err != nil {
		s.logger.Warn().Err(err).Msg("Initial policy sync failed, will retry on next interval")
	}

	// 启动定时同步
	go s.runSyncLoop(ctx)

	s.logger.Info().
		Dur("interval", s.syncInterval).
		Msg("Policy sync service started")

	return nil
}

// Stop 停止策略同步服务
func (s *PolicySyncService) Stop() {
	close(s.stopCh)
	s.logger.Info().Msg("Policy sync service stopped")
}

// runSyncLoop 定时同步循环
func (s *PolicySyncService) runSyncLoop(ctx context.Context) {
	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.SyncPolicies(ctx); err != nil {
				s.logger.Error().Err(err).Msg("Periodic policy sync failed")
			}
		case <-s.stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

// SyncPolicies 从 RPC 服务同步策略
func (s *PolicySyncService) SyncPolicies(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.enforcer == nil {
		s.logger.Warn().Msg("Enforcer is nil, skipping policy sync")
		return nil
	}

	startTime := time.Now()
	s.logger.Info().Msg("Starting policy sync from RPC service")

	// 调用 RPC 服务获取策略
	resp, err := s.identityClient.SyncPolicies(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to call SyncPolicies RPC")
		return err
	}

	// 清空现有策略
	s.enforcer.ClearPolicy()

	var policyCount, groupingCount, inheritanceCount int

	// 从响应中获取统计信息
	if resp != nil {
		if resp.Success != nil && *resp.Success {
			if resp.RolePolicyCount != nil {
				policyCount = int(*resp.RolePolicyCount)
			}
			if resp.UserRoleBindingCount != nil {
				groupingCount = int(*resp.UserRoleBindingCount)
			}
			if resp.RoleInheritanceCount != nil {
				inheritanceCount = int(*resp.RoleInheritanceCount)
			}
		}
		if resp.Message != nil {
			s.logger.Debug().Str("message", *resp.Message).Msg("RPC response message")
		}
	}

	s.lastSyncTime = time.Now()
	s.syncCount++

	duration := time.Since(startTime)
	s.logger.Info().
		Int("policy_count", policyCount).
		Int("grouping_count", groupingCount).
		Int("inheritance_count", inheritanceCount).
		Dur("duration", duration).
		Int64("total_syncs", s.syncCount).
		Msg("Policy sync completed")

	return nil
}

// GetLastSyncTime 获取上次同步时间
func (s *PolicySyncService) GetLastSyncTime() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastSyncTime
}

// GetSyncCount 获取同步次数
func (s *PolicySyncService) GetSyncCount() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.syncCount
}

// ForceSyncPolicies 强制同步策略（用于权限变更时主动刷新）
func (s *PolicySyncService) ForceSyncPolicies(ctx context.Context) error {
	return s.SyncPolicies(ctx)
}

// LoadPoliciesFromData 从数据加载策略（用于直接加载策略数据）
func (s *PolicySyncService) LoadPoliciesFromData(
	policies [][]string,
	groupingPolicies [][]string,
	roleInheritance [][]string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.enforcer == nil {
		return nil
	}

	// 清空现有策略
	s.enforcer.ClearPolicy()

	// 加载权限策略
	if len(policies) > 0 {
		if _, err := s.enforcer.AddPolicies(policies); err != nil {
			s.logger.Error().Err(err).Msg("Failed to add policies")
			return err
		}
	}

	// 加载用户-角色绑定
	if len(groupingPolicies) > 0 {
		if _, err := s.enforcer.AddGroupingPolicies(groupingPolicies); err != nil {
			s.logger.Error().Err(err).Msg("Failed to add grouping policies")
			return err
		}
	}

	// 加载角色继承关系
	if len(roleInheritance) > 0 {
		if _, err := s.enforcer.AddNamedGroupingPolicies("g2", roleInheritance); err != nil {
			s.logger.Error().Err(err).Msg("Failed to add role inheritance")
			return err
		}
	}

	s.lastSyncTime = time.Now()
	s.syncCount++

	s.logger.Info().
		Int("policy_count", len(policies)).
		Int("grouping_count", len(groupingPolicies)).
		Int("inheritance_count", len(roleInheritance)).
		Msg("Policies loaded from data")

	return nil
}
