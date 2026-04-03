// Package wire 领域服务层依赖注入提供者
package wire

import (
	"github.com/google/wire"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/zitadel/oidc/v3/pkg/op"

	identityassembler "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/identity"
	permissionConv "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/permission"
	identityservice "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/identity"
	oidcservice "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/oidc"
	permissionservice "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/permission"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/oidcstore"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/redis"
)

// DomainServiceSet 领域服务层依赖注入集合
var DomainServiceSet = wire.NewSet(
	ProvideAuthService,
	ProvideUserService,
	ProvideMembershipService,
	ProvideOrganizationService,
	ProvideDepartmentService,
	ProvideLogoService,
	ProvideAuditLogService,

	ProvideRoleDefinitionService,
	ProvideUserRoleAssignmentService,
	ProvideMenuService,

	ProvideOIDCStorage,
	ProvideOIDCService,

	ProvideIdentityService,
	ProvidePermissionService,
)

// ProvideAuthService 提供身份认证服务
func ProvideAuthService(
	identityClient identitycli.IdentityClient,
	assembler identityassembler.Assembler,
	logger *hertzZerolog.Logger,
) identityservice.AuthService {
	return identityservice.NewAuthService(identityClient, assembler, logger)
}

// ProvideUserService 提供用户管理服务
func ProvideUserService(
	identityClient identitycli.IdentityClient,
	assembler identityassembler.Assembler,
	logger *hertzZerolog.Logger,
) identityservice.UserService {
	return identityservice.NewUserManagementService(identityClient, assembler, logger)
}

// ProvideMembershipService 提供成员关系管理服务
func ProvideMembershipService(
	identityClient identitycli.IdentityClient,
	assembler identityassembler.Assembler,
	logger *hertzZerolog.Logger,
) identityservice.MembershipService {
	return identityservice.NewMembershipService(identityClient, assembler, logger)
}

// ProvideOrganizationService 提供组织管理服务
func ProvideOrganizationService(
	identityClient identitycli.IdentityClient,
	assembler identityassembler.Assembler,
	logger *hertzZerolog.Logger,
) identityservice.OrganizationService {
	return identityservice.NewOrganizationService(identityClient, assembler, logger)
}

// ProvideDepartmentService 提供部门管理服务
func ProvideDepartmentService(
	identityClient identitycli.IdentityClient,
	assembler identityassembler.Assembler,
	logger *hertzZerolog.Logger,
) identityservice.DepartmentService {
	return identityservice.NewDepartmentService(identityClient, assembler, logger)
}

// ProvideLogoService 提供Logo管理服务
func ProvideLogoService(
	identityClient identitycli.IdentityClient,
	assembler identityassembler.Assembler,
	logger *hertzZerolog.Logger,
) identityservice.LogoService {
	return identityservice.NewLogoService(identityClient, assembler, logger)
}

// ProvideAuditLogService 提供审计日志查询服务
func ProvideAuditLogService(
	identityClient identitycli.IdentityClient,
	assembler identityassembler.Assembler,
	logger *hertzZerolog.Logger,
) identityservice.AuditLogService {
	return identityservice.NewAuditLogService(identityClient, assembler, logger)
}

// ProvideRoleDefinitionService 提供角色定义服务
func ProvideRoleDefinitionService(
	identityClient identitycli.IdentityClient,
	assembler permissionConv.Assembler,
	logger *hertzZerolog.Logger,
) permissionservice.RoleDefinitionService {
	return permissionservice.NewRoleDefinitionService(identityClient, assembler, logger)
}

// ProvideUserRoleAssignmentService 提供用户角色分配服务
func ProvideUserRoleAssignmentService(
	identityClient identitycli.IdentityClient,
	assembler permissionConv.Assembler,
	logger *hertzZerolog.Logger,
) permissionservice.UserRoleAssignmentService {
	return permissionservice.NewUserRoleAssignmentService(identityClient, assembler, logger)
}

func ProvideMenuService(
	identityClient identitycli.IdentityClient,
	assembler permissionConv.Assembler,
	logger *hertzZerolog.Logger,
) permissionservice.MenuService {
	return permissionservice.NewMenuService(identityClient, assembler, logger)
}

// ProvideIdentityService 提供统一身份管理服务
func ProvideIdentityService(
	authService identityservice.AuthService,
	userService identityservice.UserService,
	membershipService identityservice.MembershipService,
	orgService identityservice.OrganizationService,
	deptService identityservice.DepartmentService,
	logoService identityservice.LogoService,
	auditLogService identityservice.AuditLogService,
) identityservice.Service {
	return identityservice.NewService(
		authService,
		userService,
		membershipService,
		orgService,
		deptService,
		logoService,
		auditLogService,
	)
}

// ProvidePermissionService 提供统一权限管理服务
func ProvidePermissionService(
	roleDefinitionService permissionservice.RoleDefinitionService,
	userRoleAssignmentService permissionservice.UserRoleAssignmentService,
	menuService permissionservice.MenuService,
) permissionservice.Service {
	return permissionservice.NewService(
		roleDefinitionService,
		userRoleAssignmentService,
		menuService,
	)
}

// ============================================================================
// OIDC 领域服务提供者
// ============================================================================

// ProvideOIDCStorage 提供 OIDC Storage（实现 op.Storage 接口）
func ProvideOIDCStorage(
	redisClient *redis.Client,
	oidcConfig *config.OIDCConfig,
	identityClient identitycli.IdentityClient,
) op.Storage {
	return oidcstore.NewStorage(redisClient.GetClient(), oidcConfig, identityClient)
}

// ProvideOIDCService 提供 OIDC 领域服务
func ProvideOIDCService(
	oidcConfig *config.OIDCConfig,
	storage op.Storage,
) oidcservice.Service {
	svc, err := oidcservice.NewService(oidcConfig, storage)
	if err != nil {
		panic(err)
	}

	return svc
}
