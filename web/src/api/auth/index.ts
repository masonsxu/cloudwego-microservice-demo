import request from '../request'
import type { LoginRequest, LoginResponse, ChangePasswordRequest } from '@/types/user'

type LoginResponseData = Omit<LoginResponse, 'base_resp'>

// 登录
export function login(data: LoginRequest) {
  return request<LoginResponseData>({
    url: '/api/v1/identity/auth/login',
    method: 'POST',
    data,
  })
}

// 登出
export function logout() {
  return request({
    url: '/api/v1/identity/auth/logout',
    method: 'POST',
    data: {},
  })
}

// 修改密码
export function changePassword(data: ChangePasswordRequest) {
  return request({
    url: '/api/v1/identity/auth/password',
    method: 'PUT',
    data,
  })
}

// 强制下次登录修改密码（管理员操作）
export function forcePasswordChange(userId: string, reason?: string) {
  return request({
    url: '/api/v1/identity/auth/password/force-change',
    method: 'PUT',
    data: { user_id: userId, reason },
  })
}

// 重置密码（管理员操作）
export function resetPassword(userId: string, newPassword: string, resetReason?: string) {
  return request({
    url: '/api/v1/identity/auth/password/reset',
    method: 'POST',
    data: { user_id: userId, new_password: newPassword, reset_reason: resetReason },
  })
}

// 刷新访问令牌
export function refreshToken(refreshTokenValue: string) {
  return request<{ token_info: { access_token: string; expires_in: number; token_type: string } }>({
    url: '/api/v1/identity/auth/refresh',
    method: 'POST',
    data: { refresh_token: refreshTokenValue },
  })
}

export const authApi = {
  login,
  logout,
  changePassword,
  forcePasswordChange,
  resetPassword,
  refreshToken,
}
