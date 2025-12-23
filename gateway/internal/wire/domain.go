// Package wire 领域服务层依赖注入提供者
package wire

import (
	"github.com/google/wire"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"

	identityassembler "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/identity"
	permissionConv "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/permission"
	identityservice "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/identity"
	permissionservice "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/permission"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
)

// DomainServiceSet 领域服务层依赖注入集合
// 按照分层架构组织：各业务服务 -> 聚合服务
var DomainServiceSet = wire.NewSet(
	// 身份管理服务领域服务
	ProvideAuthService,
	ProvideUserService,
	ProvideMembershipService,
	ProvideOrganizationService,
	ProvideDepartmentService,
	ProvideLogoService,

	// 角色与权限管理领域服务
	ProvideRoleDefinitionService,
	ProvideUserRoleAssignmentService,
	ProvideMenuService,

	// 聚合服务
	ProvideIdentityService,
	ProvidePermissionService,
)

// ============================================================================
// 各业务领域服务提供者
// ============================================================================

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
	return identityservice.NewUserManagementService(
		identityClient,
		assembler,
		logger,
	)
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

// ============================================================================
// 聚合服务提供者
// ============================================================================

// ProvideIdentityService 提供统一身份管理服务
// 使用聚合设计模式，统一管理所有身份相关功能
func ProvideIdentityService(
	authService identityservice.AuthService,
	userService identityservice.UserService,
	membershipService identityservice.MembershipService,
	orgService identityservice.OrganizationService,
	deptService identityservice.DepartmentService,
	logoService identityservice.LogoService,
) identityservice.Service {
	return identityservice.NewService(
		authService,
		userService,
		membershipService,
		orgService,
		deptService,
		logoService,
	)
}

// ProvidePermissionService 提供统一权限管理服务
// 使用聚合设计模式，统一管理所有权限相关功能
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
