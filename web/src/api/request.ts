import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'
import { toast } from 'vue-sonner'
import { useAuthStore } from '@/stores/auth'

// 创建 axios 实例
// Cookie 方案：启用 withCredentials 允许跨域携带 Cookie
// 注意：开发环境下使用相对路径 ''，API 调用已经包含 /api 前缀，通过 Vite 代理转发到后端
// 生产环境下使用完整的 API URL
const isDev = import.meta.env.DEV
const apiBaseURL = isDev ? '' : (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080')
console.log('[Request] Environment:', isDev ? 'development' : 'production', 'API Base URL:', apiBaseURL || '(relative)')

const instance = axios.create({
  baseURL: apiBaseURL,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true // 允许跨域携带 Cookie（包括 HttpOnly Cookie）
})

// 标记是否正在刷新 token
let isRefreshing = false
// 存储等待刷新完成的请求队列
let refreshSubscribers: Array<(success: boolean) => void> = []

// 将等待中的请求添加到队列
function subscribeTokenRefresh(callback: (success: boolean) => void) {
  refreshSubscribers.push(callback)
}

// 通知所有等待中的请求
function onTokenRefreshed(success: boolean) {
  refreshSubscribers.forEach(callback => callback(success))
  refreshSubscribers = []
}

// 请求拦截器
// Cookie 方案：不需要手动设置 Authorization Header，浏览器自动携带 Cookie
instance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Cookie 方案：浏览器会自动携带 HttpOnly Cookie
    // 不需要手动设置 Authorization Header
    return config
  },
  (error: AxiosError) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
instance.interceptors.response.use(
  (response) => {
    const data = response.data

    // 处理统一响应格式
    if (data.base_resp) {
      if (data.base_resp.code !== 0) {
        toast.error(data.base_resp.message || '请求失败')
        return Promise.reject(new Error(data.base_resp.message))
      }
      // 返回去除 base_resp 后的数据
      const { base_resp, ...rest } = data
      return rest as any
    }

    return data
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    if (error.response) {
      const status = error.response.status

      switch (status) {
        case 401:
          // 避免重复刷新导致的循环
          if (originalRequest._retry) {
            const authStore = useAuthStore()
            authStore.clearAuthState()
            window.location.href = '/login'
            return Promise.reject(error)
          }

          // Cookie 方案：检查是否正在刷新
          if (isRefreshing) {
            return new Promise((resolve, reject) => {
              subscribeTokenRefresh((success) => {
                if (success) {
                  // 刷新成功，重试原请求（浏览器会自动携带新 Cookie）
                  resolve(instance(originalRequest))
                } else {
                  // 刷新失败，拒绝请求
                  reject(error)
                }
              })
            })
          }

          // 开始刷新 token
          isRefreshing = true
          originalRequest._retry = true

          try {
            const authStore = useAuthStore()
            await authStore.refreshAccessToken()

            // 通知所有等待中的请求刷新成功
            onTokenRefreshed(true)
            isRefreshing = false

            // 重试原请求（浏览器会自动携带新 Cookie）
            return instance(originalRequest)
          } catch (refreshError) {
            isRefreshing = false
            onTokenRefreshed(false)
            // clearAuthState 已在 refreshAccessToken 内部调用，直接跳转登录页
            window.location.href = '/login'
            return Promise.reject(refreshError)
          }

        case 403:
          toast.error('没有权限访问')
          break

        case 404:
          toast.error('请求的资源不存在')
          break

        case 500:
          toast.error('服务器内部错误')
          break

        default:
          toast.error(error.message || '请求失败')
      }
    } else if (error.request) {
      toast.error('网络连接失败，请检查网络')
    } else {
      toast.error('请求配置错误')
    }

    return Promise.reject(error)
  }
)

// 请求函数
export default function request<T = any>(config: any): Promise<T> {
  return instance(config) as Promise<T>
}
