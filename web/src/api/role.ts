import request from './request'
import type {
  RoleDefinition,
  CreateRoleRequest,
  UpdateRoleRequest,
  ConfigureRoleMenusRequest,
  RoleMenuPermissionsResponse,
  UserRoleAssignment,
  MenuPermission
} from '@/types/role'
import type { MenuItem } from '@/types/user'
import type { PageParams, PageResponse, PermissionLevel } from '@/types/api'

// 获取角色列表
export function getRoleList(params?: PageParams & { name?: string; status?: number; isSystemRole?: boolean }) {
  return request<{ roles: RoleDefinition[]; page: PageResponse<RoleDefinition> }>({
    url: '/api/v1/permission/roles',
    method: 'GET',
    params
  })
}

// 获取角色详情
export function getRoleDetail(roleId: string) {
  return request<{ role: RoleDefinition }>({
    url: `/api/v1/permission/roles/${roleId}`,
    method: 'GET'
  })
}

// 创建角色
export function createRole(data: CreateRoleRequest) {
  return request<{ role: RoleDefinition }>({
    url: '/api/v1/permission/roles',
    method: 'POST',
    data
  })
}

// 更新角色
export function updateRole(roleId: string, data: UpdateRoleRequest) {
  return request<{ role: RoleDefinition }>({
    url: `/api/v1/permission/roles/${roleId}`,
    method: 'PUT',
    data
  })
}

// 删除角色
export function deleteRole(roleId: string) {
  return request({
    url: `/api/v1/permission/roles/${roleId}`,
    method: 'DELETE'
  })
}

// 批量绑定用户到角色
export function batchBindUsersToRole(roleId: string, userIds: string[]) {
  return request({
    url: `/api/v1/permission/roles/${roleId}/users/batch-bind`,
    method: 'POST',
    data: { user_ids: userIds }
  })
}

// 获取角色的用户列表
export function getRoleUsers(roleId: string) {
  return request<{ role_id: string; user_ids: string[] }>({
    url: `/api/v1/permission/roles/${roleId}/users`,
    method: 'GET'
  })
}

// 获取用户角色分配记录
export function getUserRoleAssignments(params?: PageParams & { userId?: string; roleId?: string }) {
  return request<{ assignments: UserRoleAssignment[]; page: PageResponse<UserRoleAssignment> }>({
    url: '/api/v1/permission/user-roles',
    method: 'GET',
    params
  })
}

// 配置角色菜单权限
export function configureRoleMenus(roleId: string, data: ConfigureRoleMenusRequest) {
  return request({
    url: `/api/v1/permission/roles/${roleId}/menus`,
    method: 'POST',
    data
  })
}

// 获取角色菜单树
export function getRoleMenuTree(roleId: string) {
  return request<{ menu_tree: MenuItem[]; role_id: string }>({
    url: `/api/v1/permission/roles/${roleId}/menu-tree`,
    method: 'GET'
  })
}

// 获取角色菜单权限列表
export function getRoleMenuPermissions(roleId: string) {
  return request<RoleMenuPermissionsResponse>({
    url: `/api/v1/permission/roles/${roleId}/menu-permissions`,
    method: 'GET'
  })
}

// 检查角色菜单权限
export function checkRoleMenuPermission(roleId: string, menuId: string, permission: PermissionLevel) {
  return request<{ has_permission: boolean; role_id: string; menu_id: string; permission: PermissionLevel }>({
    url: `/api/v1/permission/roles/${roleId}/check-menu-permission`,
    method: 'POST',
    params: { menuID: menuId, permission }
  })
}

// 获取用户菜单树
export function getUserMenuTree(userId: string) {
  return request<{ menu_tree: MenuItem[]; user_id: string; role_ids: string[]; permissions: MenuPermission[] }>({
    url: `/api/v1/permission/users/${userId}/menu-tree`,
    method: 'GET'
  })
}

// 获取当前用户的菜单树
export function getCurrentUserMenuTree() {
  return request<{ menu_tree: MenuItem[] }>({
    url: '/api/v1/permission/menu/tree',
    method: 'GET'
  })
}
