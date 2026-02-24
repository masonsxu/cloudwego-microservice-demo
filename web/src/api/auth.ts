import request from './request'
import type { LoginRequest, ChangePasswordRequest } from '@/types/user'
import type { LoginResponseDTO } from '@/api/generated'

// 登录响应（去除 base_resp 后）
type LoginResponseData = Omit<LoginResponseDTO, 'base_resp'>

// 登录
export function login(data: LoginRequest) {
  return request<LoginResponseData>({
    url: '/api/v1/identity/auth/login',
    method: 'POST',
    data
  })
}

// 登出（后端不需要 refresh_token）
export function logout() {
  return request({
    url: '/api/v1/identity/auth/logout',
    method: 'POST',
    data: {}
  })
}

// 修改密码
export function changePassword(data: ChangePasswordRequest) {
  return request({
    url: '/api/v1/identity/auth/password',
    method: 'PUT',
    data
  })
}

// 导出 API 对象
export const authApi = {
  login,
  logout,
  changePassword
}
