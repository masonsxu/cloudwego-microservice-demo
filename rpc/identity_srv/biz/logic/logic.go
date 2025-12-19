package logic

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/assignment"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/authentication"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/definition"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/department"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/logo"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/membership"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/menu"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/organization"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic/user"
)

// Logic 业务逻辑层统一接口
// 基于重构后的模块化架构，聚合所有功能模块的业务逻辑接口
// 提供完整的身份认证、用户管理、组织管理和审批流程能力
type Logic interface {
	// ============================================================================
	// 核心模块接口 - 基于新的DAL仓储架构
	// ============================================================================

	// Authentication 身份认证
	// 负责用户登录、密码修改、密码重置、密码强制修改等功能
	authentication.AuthenticationLogic

	// UserProfile 用户档案管理
	// 负责用户个人信息、认证状态、资质等核心档案数据的管理
	user.ProfileLogic

	// UserMembership 用户成员关系管理
	// 负责用户与组织、部门之间的成员关系管理，包括角色权限分配
	membership.MembershipLogic

	// Organization 组织管理
	// 负责机构和组织的层级结构管理，包括创建、更新、查询和关系维护
	organization.OrganizationLogic

	// Department 部门管理
	// 负责机构内部门的创建、管理和成员关系维护
	department.DepartmentLogic

	// Logo 组织Logo管理
	// 负责组织Logo的上传、绑定、查询和删除，支持临时Logo的生命周期管理
	logo.LogoLogic
	// ============================================================================
	// 角色与权限管理模块 - 基于新的角色权限架构
	// ============================================================================

	// RoleDefinition 角色定义管理
	// 负责角色定义的创建、更新、查询、删除等核心业务功能
	definition.RoleDefinitionLogic

	// UserRoleAssignment 用户角色分配管理
	// 负责用户角色分配的创建、更新、查询、删除等核心业务功能
	assignment.RoleAssignmentLogic

	// ============================================================================
	// 菜单管理模块 - 基于新的菜单权限架构
	// ============================================================================

	// Menu 菜单管理
	// 负责菜单配置的上传、解析、存储以及用户菜单树的构建和权限过滤
	menu.MenuLogic
}
