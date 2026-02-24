import type { BaseResponse } from './api'

// 导出枚举类型
export enum UserStatus {
  ACTIVE = 1,
  INACTIVE = 2,
  SUSPENDED = 3,
  LOCKED = 4
}

export enum Gender {
  UNKNOWN = 0,
  MALE = 1,
  FEMALE = 2
}

// 用户资料
export interface UserProfile {
  id: string
  username: string
  email?: string
  phone?: string
  first_name?: string
  last_name?: string
  real_name?: string
  professional_title?: string
  license_number?: string
  specialties?: string[]
  employee_id?: string
  gender?: Gender
  status: UserStatus
  created_at: number
  updated_at?: number
  account_expiry?: number
  must_change_password?: boolean
  // 扩展字段
  avatar?: string
  organization?: {
    id: string
    name: string
  }
  department?: {
    id: string
    name: string
  }
  role_ids?: string[]
  role_names?: string[]
}

// 用户列表项
export interface UserListItem extends UserProfile {
  organization?: {
    id: string
    name: string
  }
  department?: {
    id: string
    name: string
  }
  role_names?: string[]
}

// 登录请求
export interface LoginRequest {
  username: string
  password: string
}

// 登录响应
export interface LoginResponse extends BaseResponse {
  user_profile: UserProfile
  token_info: {
    access_token: string
    expires_in: number
    token_type: string
  }
  menu_tree: MenuItem[]
  role_ids: string[]
  permissions: MenuPermission[]
  roles: Role[]
}

// 菜单权限
export interface MenuPermission {
  menu_id: string
  permission: string  // "none", "read", "write", "full"
}

// 修改密码请求
export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

// 创建用户请求
export interface CreateUserRequest {
  username: string
  password: string
  email?: string
  phone?: string
  first_name?: string
  last_name?: string
  real_name?: string
  professional_title?: string
  license_number?: string
  specialties?: string[]
  employee_id?: string
  must_change_password?: boolean
  account_expiry?: number
  gender?: Gender
  role_ids?: string[]
  organization_id?: string
}

// 更新用户请求
export interface UpdateUserRequest {
  email?: string
  phone?: string
  first_name?: string
  last_name?: string
  real_name?: string
  professional_title?: string
  license_number?: string
  specialties?: string[]
  employee_id?: string
  account_expiry?: number
  gender?: Gender
  role_ids?: string[]
  organization_id?: string
}

// 角色
export interface Role {
  id: string
  name: string
  code?: string
  description?: string
  data_scope?: number
}

// 菜单项
export interface MenuItem {
  name: string
  id: string
  path: string
  icon?: string
  component?: string
  children?: MenuItem[]
  has_permission?: boolean
  permission_level?: number
  hidden?: boolean
}
