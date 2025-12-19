package menu

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// MenuLogic 菜单管理逻辑接口
// 负责菜单配置的上传、解析、存储以及用户菜单树的构建和权限过滤
type MenuLogic interface {
	// UploadMenu 上传并解析菜单配置文件 (menu.yaml)
	//	@param	ctx	上下文
	//	@param	req	包含	YAML	文件内容的请求
	UploadMenu(ctx context.Context, req *identity_srv.UploadMenuRequest) error

	// GetMenuTree 获取指定用户的菜单树
	//	@param	ctx	上下文
	//	@param	req	包含用户标识等信息的请求
	//	@return	用户可见的菜单树结构
	GetMenuTree(
		ctx context.Context,
	) (*identity_srv.GetMenuTreeResponse, error)

	// ConfigureRoleMenus 配置角色的菜单权限
	//	@param	ctx	上下文
	//	@param	req	包含角色ID和菜单权限配置信息
	//	@return	配置成功响应
	ConfigureRoleMenus(
		ctx context.Context,
		req *identity_srv.ConfigureRoleMenusRequest,
	) (*identity_srv.ConfigureRoleMenusResponse, error)

	// GetRoleMenuTree 获取角色的菜单树
	//	@param	ctx	上下文
	//	@param	req	包含角色ID的请求
	//	@return	角色可访问的菜单树结构
	GetRoleMenuTree(
		ctx context.Context,
		req *identity_srv.GetRoleMenuTreeRequest,
	) (*identity_srv.GetRoleMenuTreeResponse, error)

	// GetUserMenuTree 获取用户的菜单树（基于最新角色绑定）
	//	@param	ctx	上下文
	//	@param	req	包含用户ID的请求
	//	@return	用户可访问的菜单树结构
	GetUserMenuTree(
		ctx context.Context,
		req *identity_srv.GetUserMenuTreeRequest,
	) (*identity_srv.GetUserMenuTreeResponse, error)

	// GetRoleMenuPermissions 获取角色的菜单权限列表
	//	@param	ctx	上下文
	//	@param	req	包含角色ID的请求
	//	@return	角色的菜单权限配置列表
	GetRoleMenuPermissions(
		ctx context.Context,
		req *identity_srv.GetRoleMenuPermissionsRequest,
	) (*identity_srv.GetRoleMenuPermissionsResponse, error)

	// HasMenuPermission 检查角色是否具有指定菜单权限
	//	@param	ctx	上下文
	//	@param	req	包含角色ID、菜单ID和权限类型的请求
	//	@return	权限检查结果
	HasMenuPermission(
		ctx context.Context,
		req *identity_srv.HasMenuPermissionRequest,
	) (*identity_srv.HasMenuPermissionResponse, error)

	// GetUserMenuPermissions 获取用户的菜单权限列表（基于所有活跃角色合并）
	//	@param	ctx	上下文
	//	@param	req	包含用户ID的请求
	//	@return	用户的合并菜单权限列表（去重，取最高权限）
	GetUserMenuPermissions(
		ctx context.Context,
		req *identity_srv.GetUserMenuPermissionsRequest,
	) (*identity_srv.GetUserMenuPermissionsResponse, error)
}
