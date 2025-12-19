package logic

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/casbin"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	roleAssignLogic "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/assignment"
	authenticationLogic "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/authentication"
	roleDefLogic "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/definition"
	departmentLogic "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/department"
	logoLogic "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/logo"
	membershipLogic "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/membership"
	menuLogic "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/menu"
	orgLogic "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/organization"
	userLogic "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/user"
	rustfsclient "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/rustfs_client"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
)

// Impl 业务逻辑层统一实现
// 基于重构后的架构，整合DAL和新模块依赖
type Impl struct {
	dal dal.DAL
	cfg *config.Config

	// ============================================================================
	// 核心模块 - 基于新的DAL仓储架构
	// ============================================================================

	// 用户认证管理
	authenticationLogic.AuthenticationLogic

	// 用户档案管理
	userLogic.ProfileLogic

	// 用户成员关系管理
	membershipLogic.MembershipLogic

	// 组织管理
	orgLogic.OrganizationLogic

	// 部门管理
	departmentLogic.DepartmentLogic

	// 组织Logo管理
	logoLogic.LogoLogic
	// ============================================================================
	// 角色与权限管理 - 基于新的角色权限架构
	// ============================================================================

	// 角色定义管理
	roleDefLogic.RoleDefinitionLogic

	// 用户角色分配管理
	roleAssignLogic.RoleAssignmentLogic

	// ============================================================================
	// 菜单管理 - 基于新的菜单权限架构
	// ============================================================================

	// 菜单管理
	menuLogic.MenuLogic
}

// NewLogicImpl 创建业务逻辑层实例
// 基于新的DAL架构和模块化设计，初始化所有业务逻辑模块
func NewLogicImpl(dal dal.DAL, cfg *config.Config, casbinManager *casbin.CasbinManager) Logic {
	// 创建转换器实例
	conv := converter.NewConverter()

	// 创建 Logo 存储客户端
	logoStorageClient, err := rustfsclient.NewLogoStorageClient(&cfg.LogoStorage)

	// 创建 OrganizationLogic（需要 logoStorageClient 生成 logo URL）
	var orgLogicImpl orgLogic.OrganizationLogic
	if err != nil || logoStorageClient == nil {
		// 如果 Logo 存储客户端初始化失败，仍然创建 OrganizationLogic 但传入 nil
		// OrganizationLogic 内部会处理 logoStorageClient 为 nil 的情况
		orgLogicImpl = orgLogic.NewLogic(dal, conv, nil)
	} else {
		orgLogicImpl = orgLogic.NewLogic(dal, conv, logoStorageClient)
	}

	// 创建 LogoLogic（需要 logoStorageClient 上传文件到 S3）
	var logoLogicImpl logoLogic.LogoLogic
	if err != nil || logoStorageClient == nil {
		// 如果 Logo 存储客户端初始化失败，设置为nil
		// 实际调用时会返回友好的错误信息
		logoLogicImpl = nil
	} else {
		logoLogicImpl = logoLogic.NewLogic(
			dal.Logo(),
			conv,
			logoStorageClient,
		)
	}

	// 创建菜单逻辑（提前初始化以供AuthenticationLogic使用）
	menuLogicImpl := menuLogic.NewLogic(
		dal,
		conv,
		casbinManager,
		dal.UserRoleAssignment(),
		cfg,
	)

	return &Impl{
		dal: dal,
		cfg: cfg,

		// ============================================================================
		// 核心模块初始化 - 使用新的仓储架构
		// ============================================================================

		AuthenticationLogic: authenticationLogic.NewLogic(
			dal,
			conv,
			menuLogicImpl,
		),

		// 用户档案逻辑（替代传统的user模块）
		ProfileLogic: userLogic.NewLogic(dal, conv),

		// 用户成员关系逻辑（新增模块）
		MembershipLogic: membershipLogic.NewLogic(dal, conv),

		// 组织管理逻辑（重构现有organization模块）
		OrganizationLogic: orgLogicImpl,

		// 部门管理逻辑（新增模块）
		DepartmentLogic: departmentLogic.NewLogic(dal, conv),

		// 组织Logo管理逻辑
		LogoLogic: logoLogicImpl,
		// ============================================================================
		// 角色与权限管理初始化 - 使用新的角色权限架构
		// ============================================================================

		// 角色定义逻辑
		RoleDefinitionLogic: roleDefLogic.NewLogic(dal, conv),

		// 用户角色分配逻辑
		RoleAssignmentLogic: roleAssignLogic.NewLogic(dal, conv),

		// ============================================================================
		// 菜单管理初始化 - 使用新的菜单权限架构
		// ============================================================================

		// 菜单逻辑
		MenuLogic: menuLogicImpl,
	}
}

// NewLogic 创建业务逻辑层实例（工厂函数）
func NewLogic(dal dal.DAL, cfg *config.Config, casbinManager *casbin.CasbinManager) Logic {
	return NewLogicImpl(dal, cfg, casbinManager)
}
