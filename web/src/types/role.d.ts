import type { BaseResponse, PermissionLevel } from './api'

// 角色定义
export interface RoleDefinition {
  id: string
  name: string
  description?: string
  status: number
  permissions: Permission[]
  is_system_role: boolean
  user_count?: number
  created_at: number
  updated_at?: number
}

// 权限
export interface Permission {
  resource: string
  action: string
  description?: string
}

// 创建角色请求
export interface CreateRoleRequest {
  name: string
  description?: string
  permissions?: Permission[]
  is_system_role?: boolean
}

// 更新角色请求
export interface UpdateRoleRequest {
  name?: string
  description?: string
  status?: number
  permissions?: Permission[]
}

// 菜单权限配置
export interface MenuPermission {
  menu_id: string
  permission: PermissionLevel
}

// 配置角色菜单权限请求
export interface ConfigureRoleMenusRequest {
  menu_configs: MenuPermission[]
}

// 角色菜单权限响应
export interface RoleMenuPermissionsResponse extends BaseResponse {
  permissions: MenuPermission[]
  role_id: string
}

// 用户角色分配
export interface UserRoleAssignment {
  id: string
  user_id: string
  role_id: string
  assigned_by: string
  assigned_at: number
  reason?: string
  user?: {
    id: string
    username: string
    real_name?: string
  }
  role?: {
    id: string
    name: string
  }
}
