namespace go permission

include "../../base/core.thrift"
include "../base/base.thrift"
include "./permission_model.thrift"

/**
 * PermissionService：角色与权限管理服务
 *
 * 提供对角色定义、权限、用户角色分配的管理功能。
 */
service PermissionService {
    // -----------------------------------------------------------------
    // 角色与权限管理模块 (Role & Permission Management)
    // -----------------------------------------------------------------

    /**
     * CreateRoleDefinition：创建角色定义
     *
     * 创建一个新的角色定义。
     *
     * @param req (RoleDefinitionCreateRequestDTO) - 角色定义创建请求参数
     * @return resp (RoleDefinitionCreateResponseDTO) - 创建成功的角色定义信息
     */
    permission_model.RoleDefinitionCreateResponseDTO CreateRoleDefinition(1: permission_model.RoleDefinitionCreateRequestDTO req) (api.post = "/api/v1/permission/roles"),

    /**
     * UpdateRoleDefinition：更新角色定义
     *
     * 更新一个已有的角色定义。
     *
     * @param req (RoleDefinitionUpdateRequestDTO) - 角色定义更新请求参数
     * @return resp (RoleDefinitionUpdateResponseDTO) - 更新后的角色定义信息
     */
    permission_model.RoleDefinitionUpdateResponseDTO UpdateRoleDefinition(1: permission_model.RoleDefinitionUpdateRequestDTO req) (api.put = "/api/v1/permission/roles/:roleDefinitionID"),

    /**
     * DeleteRoleDefinition：删除角色定义
     *
     * 删除一个角色定义。
     *
     * @param req (RoleDefinitionDeleteRequestDTO) - 包含要删除的角色ID
     * @return resp (OperationStatusResponseDTO) - 基础操作响应
     */
    base.OperationStatusResponseDTO DeleteRoleDefinition(1: permission_model.RoleDefinitionDeleteRequestDTO req) (api.delete = "/api/v1/permission/roles/:roleID"),

    /**
     * GetRoleDefinition：获取角色定义
     *
     * 根据ID获取角色定义。
     *
     * @param req (RoleDefinitionGetRequestDTO) - 包含要查询的角色ID
     * @return resp (RoleDefinitionGetResponseDTO) - 查找到的角色定义信息
     */
    permission_model.RoleDefinitionGetResponseDTO GetRoleDefinition(1: permission_model.RoleDefinitionGetRequestDTO req) (api.get = "/api/v1/permission/roles/:roleID"),

    /**
     * ListRoleDefinitions：列出角色定义
     *
     * 分页列出角色定义。
     *
     * @param req (RoleDefinitionQueryRequestDTO) - 查询条件，可根据分类、状态等筛选
     * @return resp (RoleDefinitionListResponseDTO) - 角色定义列表及分页信息
     */
    permission_model.RoleDefinitionListResponseDTO ListRoleDefinitions(1: permission_model.RoleDefinitionQueryRequestDTO req) (api.get = "/api/v1/permission/roles"),

    /**
     * GetLastUserRoleAssignment：获取用户最后一次角色分配
     *
     * 获取用户最后一次的角色分配信息。
     *
     * @param req (GetLastUserRoleAssignmentRequestDTO) - 用户最后角色分配请求DTO
     * @return resp (GetLastUserRoleAssignmentResponseDTO) - 用户角色分配信息
     */
    permission_model.GetLastUserRoleAssignmentResponseDTO GetLastUserRoleAssignment(1: permission_model.GetLastUserRoleAssignmentRequestDTO req) (api.get = "/api/v1/permission/users/:userID/roles/latest"),

    /**
     * ListUserRoleAssignments：列出用户角色分配记录
     *
     * 列出用户的角色分配记录。
     *
     * @param req (UserRoleQueryRequestDTO) - 查询条件，可根据用户、角色等筛选
     * @return resp (UserRoleListResponseDTO) - 角色分配记录列表及分页信息
     */
    permission_model.UserRoleListResponseDTO ListUserRoleAssignments(1: permission_model.UserRoleQueryRequestDTO req) (api.get = "/api/v1/permission/user-roles"),

    /**
     * GetUsersByRole：根据角色ID获取所有用户
     *
     * 获取指定角色下所有用户的ID列表（不分页）。
     *
     * @param req (GetUsersByRoleRequestDTO) - 包含角色ID的请求
     * @return resp (GetUsersByRoleResponseDTO) - 该角色下所有用户的ID列表
     */
    permission_model.GetUsersByRoleResponseDTO GetUsersByRole(1: permission_model.GetUsersByRoleRequestDTO req) (api.get = "/api/v1/permission/roles/:roleID/users"),

    /**
     * BatchBindUsersToRole：批量绑定用户到角色
     *
     * 批量为角色绑定用户，覆盖旧的绑定关系。
     *
     * @param req (BatchBindUsersToRoleRequestDTO) - 包含角色ID和用户ID列表
     * @return resp (BatchBindUsersToRoleResponseDTO) - 批量绑定操作结果
     */
    permission_model.BatchBindUsersToRoleResponseDTO BatchBindUsersToRole(1: permission_model.BatchBindUsersToRoleRequestDTO req) (api.post = "/api/v1/permission/roles/:roleID/users/batch-bind"),
    // -----------------------------------------------------------------
    // 菜单管理模块 (Menu Management)
    // -----------------------------------------------------------------

    /**
	 * UploadMenu：上传菜单
	 *
	 * 上传菜单配置文件，用于定义系统的菜单结构和权限。
	 *
	 * @param req (UploadMenuRequestDTO) - 上传菜单请求参数
	 * @return err (core.Error) - 错误信息
	 */
    base.OperationStatusResponseDTO UploadMenu(1: permission_model.UploadMenuRequestDTO req) (api.post = "/api/v1/permission/menu/upload"),

    /**
	 * GetMenuTree：获取菜单树
	 *
	 * 根据当前用户的角色，返回其有权限访问的菜单树结构。
	 *
	 * @return resp (GetMenuTreeResponseDTO) - 菜单树响应数据
	 * @return err (core.Error) - 错误信息
	 */
    permission_model.GetMenuTreeResponseDTO GetMenuTree() (api.get = "/api/v1/permission/menu/tree"),
    // -----------------------------------------------------------------
    // 菜单权限管理模块 (Menu Permission Management)
    // -----------------------------------------------------------------

    /**
     * ConfigureRoleMenus：配置角色的菜单权限
     *
     * 为指定角色配置菜单权限，会清除角色的旧菜单映射，然后添加新的映射。
     *
     * @param req (ConfigureRoleMenusRequestDTO) - 配置角色菜单权限请求参数
     * @return resp (ConfigureRoleMenusResponseDTO) - 配置操作响应
     */
    permission_model.ConfigureRoleMenusResponseDTO ConfigureRoleMenus(1: permission_model.ConfigureRoleMenusRequestDTO req) (api.post = "/api/v1/permission/roles/:roleID/menus"),

    /**
     * GetRoleMenuTree：获取角色的菜单树
     *
     * 根据角色的权限配置，返回该角色可访问的菜单树结构。
     *
     * @param req (GetRoleMenuTreeRequestDTO) - 获取角色菜单树请求参数
     * @return resp (GetRoleMenuTreeResponseDTO) - 角色菜单树响应数据
     */
    permission_model.GetRoleMenuTreeResponseDTO GetRoleMenuTree(1: permission_model.GetRoleMenuTreeRequestDTO req) (api.get = "/api/v1/permission/roles/:roleID/menu-tree"),

    /**
     * GetUserMenuTree：获取用户的菜单树
     *
     * 根据用户的最新角色绑定，返回该用户可访问的菜单树结构。
     *
     * @param req (GetUserMenuTreeRequestDTO) - 获取用户菜单树请求参数
     * @return resp (GetUserMenuTreeResponseDTO) - 用户菜单树响应数据
     */
    permission_model.GetUserMenuTreeResponseDTO GetUserMenuTree(1: permission_model.GetUserMenuTreeRequestDTO req) (api.get = "/api/v1/permission/users/:userID/menu-tree"),

    /**
     * GetRoleMenuPermissions：获取角色的菜单权限列表
     *
     * 获取指定角色的所有菜单权限配置列表。
     *
     * @param req (GetRoleMenuPermissionsRequestDTO) - 获取角色菜单权限请求参数
     * @return resp (GetRoleMenuPermissionsResponseDTO) - 角色菜单权限列表响应数据
     */
    permission_model.GetRoleMenuPermissionsResponseDTO GetRoleMenuPermissions(1: permission_model.GetRoleMenuPermissionsRequestDTO req) (api.get = "/api/v1/permission/roles/:roleID/menu-permissions"),

    /**
     * HasMenuPermission：检查角色是否具有指定菜单权限
     *
     * 检查指定角色是否具有对指定菜单的特定权限。
     *
     * @param req (HasMenuPermissionRequestDTO) - 检查菜单权限请求参数
     * @return resp (HasMenuPermissionResponseDTO) - 权限检查结果响应数据
     */
    permission_model.HasMenuPermissionResponseDTO HasMenuPermission(1: permission_model.HasMenuPermissionRequestDTO req) (api.post = "/api/v1/permission/roles/:roleID/check-menu-permission"),
}