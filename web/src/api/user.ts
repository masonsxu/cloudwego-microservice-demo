import request from './request'
import type {
  UserProfile,
  UserListItem,
  CreateUserRequest,
  UpdateUserRequest
} from '@/types/user'
import type { PageParams, PageResponse } from '@/types/api'

// 获取当前用户信息
export function getCurrentUser() {
  return request<{ user: UserProfile }>({
    url: '/api/v1/identity/users/me',
    method: 'GET'
  })
}

// 获取用户列表
export function getUserList(params: PageParams & { organization_id?: string; status?: number }) {
  return request<{ users: UserListItem[]; page: PageResponse<UserListItem> }>({
    url: '/api/v1/identity/users',
    method: 'GET',
    params
  })
}

// 获取用户详情
export function getUserDetail(userId: string) {
  return request<{ user: UserProfile }>({
    url: `/api/v1/identity/users/${userId}`,
    method: 'GET'
  })
}

// 创建用户
export function createUser(data: CreateUserRequest) {
  return request<{ user: UserProfile }>({
    url: '/api/v1/identity/users',
    method: 'POST',
    data
  })
}

// 更新用户信息
export function updateUser(userId: string, data: UpdateUserRequest) {
  return request<{ user: UserProfile }>({
    url: `/api/v1/identity/users/${userId}`,
    method: 'PUT',
    data
  })
}

// 更新当前用户信息
export function updateCurrentUser(data: UpdateUserRequest) {
  return request<{ user: UserProfile }>({
    url: '/api/v1/identity/users/me',
    method: 'PUT',
    data
  })
}

// 删除用户
export function deleteUser(userId: string, reason?: string) {
  return request({
    url: `/api/v1/identity/users/${userId}`,
    method: 'DELETE',
    data: { reason }
  })
}

// 变更用户状态
export function changeUserStatus(userId: string, newStatus: number, reason?: string) {
  return request({
    url: `/api/v1/identity/users/${userId}/status`,
    method: 'PUT',
    data: { new_status: newStatus, reason }
  })
}

// 解锁用户
export function unlockUser(userId: string) {
  return request({
    url: `/api/v1/identity/users/${userId}/unlock`,
    method: 'PUT'
  })
}

// 搜索用户
export function searchUsers(params: PageParams & { organization_id?: string }) {
  return request<{ users: UserListItem[]; page: PageResponse<UserListItem> }>({
    url: '/api/v1/identity/users/search',
    method: 'GET',
    params
  })
}
