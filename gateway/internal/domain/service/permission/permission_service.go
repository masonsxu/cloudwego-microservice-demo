package permission

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
)

// Service 权限管理聚合服务接口 - 统一暴露给Handler层
type Service interface {
	RoleDefinitionService
	UserRoleAssignmentService
	// PermissionSyncService
	MenuService
}

// =================================================================
// 专门化服务接口定义 - 按业务领域划分
// =================================================================

// RoleDefinitionService 角色定义管理服务接口
type RoleDefinitionService interface {
	// CreateRoleDefinition 创建角色定义 - 创建一个新的角色及其权限配置
	CreateRoleDefinition(
		ctx context.Context,
		req *permission.RoleDefinitionCreateRequestDTO,
	) (*permission.RoleDefinitionCreateResponseDTO, error)

	// UpdateRoleDefinition 更新角色定义 - 更新现有角色的权限配置
	UpdateRoleDefinition(
		ctx context.Context,
		req *permission.RoleDefinitionUpdateRequestDTO,
	) (*permission.RoleDefinitionUpdateResponseDTO, error)

	// DeleteRoleDefinition 删除角色定义 - 软删除指定角色
	DeleteRoleDefinition(
		ctx context.Context,
		roleID string,
	) (*http_base.OperationStatusResponseDTO, error)

	// GetRoleDefinition 获取角色定义 - 根据角色ID获取角色详细信息
	GetRoleDefinition(
		ctx context.Context,
		req *permission.RoleDefinitionGetRequestDTO,
	) (*permission.RoleDefinitionGetResponseDTO, error)

	// ListRoleDefinitions 列出角色定义 - 分页查询角色列表，支持多条件筛选
	ListRoleDefinitions(
		ctx context.Context,
		req *permission.RoleDefinitionQueryRequestDTO,
	) (*permission.RoleDefinitionListResponseDTO, error)
}

// UserRoleAssignmentService 用户角色分配管理服务接口
type UserRoleAssignmentService interface {
	// GetLastUserRoleAssignment 获取最新用户角色分配 - 获取用户最新的角色分配记录
	GetLastUserRoleAssignment(
		ctx context.Context,
		userID string,
	) (*permission.AssignRoleToUserResponseDTO, error)
	// ListUserRoleAssignments 列出用户角色分配 - 查询用户角色分配记录
	ListUserRoleAssignments(
		ctx context.Context,
		req *permission.UserRoleQueryRequestDTO,
	) (*permission.UserRoleListResponseDTO, error)

	// GetUsersByRole 根据角色ID获取所有用户 - 获取指定角色下所有用户的ID列表
	GetUsersByRole(
		ctx context.Context,
		req *permission.GetUsersByRoleRequestDTO,
	) (*permission.GetUsersByRoleResponseDTO, error)

	// BatchBindUsersToRole 批量绑定用户到角色 - 批量为角色绑定用户，覆盖旧的绑定关系
	BatchBindUsersToRole(
		ctx context.Context,
		operatorID string,
		req *permission.BatchBindUsersToRoleRequestDTO,
	) (*permission.BatchBindUsersToRoleResponseDTO, error)
}

// MenuService 菜单服务接口
type MenuService interface {
	// UploadMenu 上传菜单 - 上传菜单到权限引擎
	UploadMenu(
		ctx context.Context,
		req *permission.UploadMenuRequestDTO,
	) (*http_base.OperationStatusResponseDTO, error)

	// GetMenuTree 获取菜单树 - 获取权限引擎中的菜单树结构
	GetMenuTree(
		ctx context.Context,
	) (*permission.GetMenuTreeResponseDTO, error)

	// ConfigureRoleMenus 配置角色的菜单权限 - 为指定角色配置菜单权限
	ConfigureRoleMenus(
		ctx context.Context,
		operatorID string,
		req *permission.ConfigureRoleMenusRequestDTO,
	) (*permission.ConfigureRoleMenusResponseDTO, error)

	// GetRoleMenuTree 获取角色的菜单树 - 根据角色的权限配置返回菜单树
	GetRoleMenuTree(
		ctx context.Context,
		req *permission.GetRoleMenuTreeRequestDTO,
	) (*permission.GetRoleMenuTreeResponseDTO, error)

	// GetUserMenuTree 获取用户的菜单树 - 根据用户的最新角色绑定返回菜单树
	GetUserMenuTree(
		ctx context.Context,
		req *permission.GetUserMenuTreeRequestDTO,
	) (*permission.GetUserMenuTreeResponseDTO, error)

	// GetRoleMenuPermissions 获取角色的菜单权限列表 - 获取指定角色的所有菜单权限配置
	GetRoleMenuPermissions(
		ctx context.Context,
		req *permission.GetRoleMenuPermissionsRequestDTO,
	) (*permission.GetRoleMenuPermissionsResponseDTO, error)

	// HasMenuPermission 检查角色是否具有指定菜单权限 - 检查权限
	HasMenuPermission(
		ctx context.Context,
		req *permission.HasMenuPermissionRequestDTO,
	) (*permission.HasMenuPermissionResponseDTO, error)
}
