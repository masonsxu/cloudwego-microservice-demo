import request from '../request'
import type { MenuItem, MenuPermission } from '@/types/user'

export const menuApi = {
  // 获取完整菜单树（用于管理后台展示）
  getMenuTree: () =>
    request<{ menu_tree: MenuItem[] }>({
      method: 'GET',
      url: '/api/v1/permission/menu/tree',
    }),

  // 上传菜单配置文件（YAML）
  uploadMenuConfig: (menuFile: File) => {
    const formData = new FormData()
    formData.append('menu_file', menuFile)
    return request<{ menu_tree: MenuItem[] }>({
      method: 'POST',
      url: '/api/v1/permission/menu/upload',
      data: formData,
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  // 获取角色的菜单权限列表
  getRoleMenuPermissions: (roleId: string) =>
    request<{ permissions: MenuPermission[]; role_id: string }>({
      method: 'GET',
      url: `/api/v1/permission/roles/${roleId}/menu-permissions`,
    }),

  // 获取角色的菜单树
  getRoleMenuTree: (roleId: string) =>
    request<{ menu_tree: MenuItem[]; role_id: string }>({
      method: 'GET',
      url: `/api/v1/permission/roles/${roleId}/menu-tree`,
    }),

  // 配置角色菜单权限（全量覆盖）
  configureRoleMenus: (roleId: string, menuConfigs: { menu_id: string; permission: number }[]) =>
    request<{ success: boolean; message: string }>({
      method: 'POST',
      url: `/api/v1/permission/roles/${roleId}/menus`,
      data: { menu_configs: menuConfigs },
    }),

  // 检查角色是否具有指定菜单权限
  checkRoleMenuPermission: (roleId: string, menuId: string, permission: number) =>
    request<{ has_permission: boolean; menu_id: string; permission: number; role_id: string }>({
      method: 'POST',
      url: `/api/v1/permission/roles/${roleId}/check-menu-permission`,
      data: { menu_id: menuId, permission },
    }),

  // 获取用户的菜单树
  getUserMenuTree: (userId: string) =>
    request<{ menu_tree: MenuItem[]; permissions: MenuPermission[]; role_ids: string[]; user_id: string }>({
      method: 'GET',
      url: `/api/v1/permission/users/${userId}/menu-tree`,
    }),
}
