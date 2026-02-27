import request from '../request'
import type { RoleDefinition, CreateRoleRequest, UpdateRoleRequest } from '@/types/role'
import type { MenuItem } from '@/types/user'
import type { PageParams } from '@/types/api'

// 供 view 直接使用的类型别名（保持命名一致性）
export type RoleDefinitionDTO = RoleDefinition
export type MenuNodeDTO = MenuItem

export interface MenuConfig {
  menu_id: string
  permission: number
}

export interface ListRolesParams extends PageParams {
  name?: string
  status?: number
  is_system_role?: boolean
  fetch_all?: boolean
}

export const roleApi = {
  // 获取角色列表
  getRoles: (params?: ListRolesParams) =>
    request<{ roles: RoleDefinition[]; page: any }>({
      method: 'GET',
      url: '/api/v1/permission/roles',
      params,
    }),

  // 获取角色详情
  getRole: (roleId: string) =>
    request<{ role: RoleDefinition }>({
      method: 'GET',
      url: `/api/v1/permission/roles/${roleId}`,
    }),

  // 创建角色
  createRole: (data: CreateRoleRequest) =>
    request<{ role: RoleDefinition }>({
      method: 'POST',
      url: '/api/v1/permission/roles',
      data,
    }),

  // 更新角色
  updateRole: (roleId: string, data: UpdateRoleRequest) =>
    request<{ role: RoleDefinition }>({
      method: 'PUT',
      url: `/api/v1/permission/roles/${roleId}`,
      data,
    }),

  // 删除角色
  deleteRole: (roleId: string) =>
    request({
      method: 'DELETE',
      url: `/api/v1/permission/roles/${roleId}`,
    }),

  // 配置角色菜单权限（全量覆盖）
  configureRoleMenus: (roleId: string, menuConfigs: MenuConfig[]) =>
    request<{ success: boolean; message: string }>({
      method: 'POST',
      url: `/api/v1/permission/roles/${roleId}/menus`,
      data: { menu_configs: menuConfigs },
    }),

  // 获取角色的菜单权限列表
  getRoleMenuPermissions: (roleId: string) =>
    request<{ permissions: { menu_id: string; permission: number }[]; role_id: string }>({
      method: 'GET',
      url: `/api/v1/permission/roles/${roleId}/menu-permissions`,
    }),

  // 获取角色的菜单树
  getRoleMenuTree: (roleId: string) =>
    request<{ menu_tree: MenuItem[]; role_id: string }>({
      method: 'GET',
      url: `/api/v1/permission/roles/${roleId}/menu-tree`,
    }),

  // 获取角色下的用户 ID 列表
  getRoleUsers: (roleId: string) =>
    request<{ role_id: string; user_ids: string[] }>({
      method: 'GET',
      url: `/api/v1/permission/roles/${roleId}/users`,
    }),

  // 批量绑定用户到角色
  bindUsersToRole: (roleId: string, userIds: string[]) =>
    request<{ success: boolean; success_count: number; message: string }>({
      method: 'POST',
      url: `/api/v1/permission/roles/${roleId}/users/batch-bind`,
      data: { user_ids: userIds },
    }),

  // 列出用户的角色分配记录
  listUserRoles: (params?: PageParams & { user_id?: string; role_id?: string }) =>
    request<{ assignments: any[]; page: any }>({
      method: 'GET',
      url: '/api/v1/permission/user-roles',
      params,
    }),

  // 获取用户最后一次角色分配
  getLatestUserRole: (userId: string) =>
    request<{ assignment_id: string }>({
      method: 'GET',
      url: `/api/v1/permission/users/${userId}/roles/latest`,
    }),
}

// 命名导出（供用户相关 view 使用）
export function getRoleList(params?: ListRolesParams) {
  return roleApi.getRoles(params)
}
