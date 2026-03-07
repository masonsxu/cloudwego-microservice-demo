// Package identitycli 提供与身份服务交互的客户端实现
package identitycli

import (
	"context"
	"log/slog"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/fallback"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/transport"
	kitextracing "github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"

	conf "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv/identityservice"
)

const (
	defaultIdentityServiceName = "identity-service"
)

// NewIdentityClient 创建聚合的用户客户端，使用单一的 Kitex 客户端实例
func NewIdentityClient() (IdentityClient, error) {
	r, err := etcd.NewEtcdResolver([]string{conf.Config.Etcd.Address})
	if err != nil {
		slog.Error("Failed to create etcd resolver", "error", err)
		return nil, err
	}

	identityServiceName := defaultIdentityServiceName
	if service, exists := conf.Config.Client.Services["identity"]; exists &&
		service.Name != "" {
		identityServiceName = service.Name
	}

	slog.Info("Creating identity client", "service_name", identityServiceName)

	// ========== 熔断器配置 ==========
	// CBSuite 提供服务级别 + 实例级别双层熔断保护
	// 默认：10s 窗口内采样 >= 200 且错误率 >= 50% 时触发熔断
	cbs := circuitbreak.NewCBSuite(circuitbreak.RPCInfo2Key)

	// ========== 重试配置（与熔断器联动，共享统计减少开销）==========
	rc := retry.NewRetryContainerWithCB(cbs.ServiceControl(), cbs.ServicePanel())

	// ========== 连接池配置 ==========
	idleConfig := connpool.IdleConfig{
		// MaxIdlePerAddress 估算: QPS_per_host * avg_response_time_sec
		MaxIdlePerAddress: conf.Config.Client.Pool.MaxIdlePerAddress,
		MinIdlePerAddress: conf.Config.Client.Pool.MinIdlePerAddress,
		MaxIdleGlobal:     conf.Config.Client.Pool.MaxIdleGlobal,
		MaxIdleTimeout:    conf.Config.Client.Pool.MaxIdleTimeout,
	}

	// ========== Fallback 降级配置 ==========
	// 超时或熔断时触发降级，记录日志后返回原始错误
	fbPolicy := fallback.TimeoutAndCBFallback(func(
		ctx context.Context,
		_ utils.KitexArgs,
		_ utils.KitexResult,
		err error,
	) error {
		ri := rpcinfo.GetRPCInfo(ctx)
		slog.Warn("RPC fallback triggered",
			"method", ri.To().Method(),
			"service", ri.To().ServiceName(),
			"error", err,
		)

		return err
	})

	// ========== 组装 Client Options ==========
	opts := []client.Option{
		// 基础配置
		client.WithResolver(r),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: conf.Config.Server.Name,
		}),
		client.WithTransportProtocol(transport.TTHeader),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),

		// 连接池
		client.WithLongConnection(idleConfig),

		// 超时
		client.WithConnectTimeout(conf.Config.Client.ConnectionTimeout),
		client.WithRPCTimeout(conf.Config.Client.RequestTimeout),

		// 负载均衡：加权轮询，分布更均匀
		client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()),

		// 熔断器：服务级别 + 实例级别双层保护
		client.WithCircuitBreaker(cbs),
		client.WithCloseCallbacks(func() error {
			return cbs.Close()
		}),

		// 重试（与熔断器联动）
		client.WithRetryContainer(rc),
		// 方法级别差异化重试策略
		client.WithRetryMethodPolicies(buildRetryPolicies()),

		// Fallback 降级
		client.WithFallback(fbPolicy),
	}

	// 添加 OpenTelemetry tracing Suite (如果启用)
	if conf.Config.Tracing.Enabled {
		opts = append(opts, client.WithSuite(kitextracing.NewClientSuite()))

		slog.Debug("OpenTelemetry tracing enabled for identity client")
	}

	cli, err := identityservice.NewClient(
		identityServiceName,
		opts...,
	)
	if err != nil {
		slog.Error("Failed to create identity client", "error", err)
		return nil, err
	}

	slog.Info("Successfully created identity client",
		"circuit_breaker", "enabled",
		"load_balancer", "weighted_round_robin",
		"retry", "per_method",
		"fallback", "timeout_and_cb",
	)

	return cli, nil
}

// buildRetryPolicies 构建方法级别的差异化重试策略
//
// 策略原则：
//   - 只读方法（Get/List/Search/Check）: BackupRequest — 降低尾部延迟
//   - 幂等写方法（Update/ChangeStatus）: FailureRetry — 提高成功率
//   - 非幂等写方法（Create/Delete）: 不配置重试 — 避免重复执行
func buildRetryPolicies() map[string]retry.Policy {
	// BackupRequest: 等待指定延迟后并发发送备份请求，谁先返回用谁
	backupPolicy := func(delayMS uint32) retry.Policy {
		bp := retry.NewBackupPolicy(delayMS)
		bp.WithMaxRetryTimes(1)
		bp.WithRetryBreaker(0.1) // 重试请求占比超 10% 自动停止重试

		return retry.BuildBackupRequest(bp)
	}

	// FailureRetry: 请求失败后重试，带随机退避
	failurePolicy := func() retry.Policy {
		fp := retry.NewFailurePolicy()
		fp.WithMaxRetryTimes(1)
		fp.WithMaxDurationMS(3000)
		fp.WithRetryBreaker(0.1)
		fp.WithRandomBackOff(50, 200)

		return retry.BuildFailurePolicy(fp)
	}

	return map[string]retry.Policy{
		// ===== 用户模块 =====
		"GetUser":     backupPolicy(200),
		"ListUsers":   backupPolicy(300),
		"SearchUsers": backupPolicy(300),
		"UpdateUser":  failurePolicy(),

		// ===== 认证模块 =====
		"Login": failurePolicy(),

		// ===== 成员关系模块 =====
		"GetMembership":        backupPolicy(200),
		"GetUserMemberships":   backupPolicy(200),
		"GetPrimaryMembership": backupPolicy(200),
		"CheckMembership":      backupPolicy(200),
		"UpdateMembership":     failurePolicy(),

		// ===== 组织模块 =====
		"GetOrganization":    backupPolicy(200),
		"ListOrganizations":  backupPolicy(300),
		"UpdateOrganization": failurePolicy(),

		// ===== 部门模块 =====
		"GetDepartment":              backupPolicy(200),
		"GetOrganizationDepartments": backupPolicy(300),
		"UpdateDepartment":           failurePolicy(),

		// ===== Logo 模块 =====
		"GetOrganizationLogo": backupPolicy(200),

		// ===== 角色模块 =====
		"GetRoleDefinition":    backupPolicy(200),
		"ListRoleDefinitions":  backupPolicy(300),
		"UpdateRoleDefinition": failurePolicy(),

		// ===== 用户角色分配模块 =====
		"GetLastUserRoleAssignment": backupPolicy(200),
		"ListUserRoleAssignments":   backupPolicy(200),
		"GetUsersByRole":            backupPolicy(300),
		"BatchGetUserRoles":         backupPolicy(300),
		"UpdateUserRoleAssignment":  failurePolicy(),

		// ===== 菜单与权限模块 =====
		"GetMenuTree":            backupPolicy(200),
		"GetRoleMenuTree":        backupPolicy(200),
		"GetUserMenuTree":        backupPolicy(200),
		"GetRoleMenuPermissions": backupPolicy(200),
		"HasMenuPermission":      backupPolicy(100),
		"GetUserMenuPermissions": backupPolicy(200),
		"CheckPermission":        backupPolicy(100),
		"GetUserDataScope":       backupPolicy(200),

		// 注意：以下方法为非幂等操作，不配置重试
		// CreateUser, DeleteUser, ChangeUserStatus, UnlockUser
		// ChangePassword, ResetPassword, ForcePasswordChange
		// AddMembership, RemoveMembership
		// CreateOrganization, DeleteOrganization
		// CreateDepartment, DeleteDepartment
		// UploadTemporaryLogo, DeleteOrganizationLogo, BindLogoToOrganization
		// CreateRoleDefinition, DeleteRoleDefinition
		// AssignRoleToUser, RevokeRoleFromUser, BatchBindUsersToRole
		// UploadMenu, ConfigureRoleMenus, SyncPolicies
	}
}
