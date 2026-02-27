import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

// 创建 axios 实例
const instance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
instance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const authStore = useAuthStore()
    // 优先使用 authStore 中的 token，如果为空则从 localStorage 读取（备用方案）
    const token = authStore.token || localStorage.getItem('token')
    
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
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
        ElMessage.error(data.base_resp.message || '请求失败')
        return Promise.reject(new Error(data.base_resp.message))
      }
      // 返回去除 base_resp 后的数据
      const { base_resp, ...rest } = data
      return rest as any
    }

    return data
  },
  async (error: AxiosError) => {
    if (error.response) {
      const status = error.response.status

      switch (status) {
        case 401:
          // Token 过期或无效，直接跳转登录
          const authStore = useAuthStore()
          authStore.logout()
          window.location.href = '/login'
          ElMessage.error('登录已过期，请重新登录')
          break

        case 403:
          ElMessage.error('没有权限访问')
          break

        case 404:
          ElMessage.error('请求的资源不存在')
          break

        case 500:
          ElMessage.error('服务器内部错误')
          break

        default:
          ElMessage.error(error.message || '请求失败')
      }
    } else if (error.request) {
      ElMessage.error('网络连接失败，请检查网络')
    } else {
      ElMessage.error('请求配置错误')
    }

    return Promise.reject(error)
  }
)

// 请求函数
export default function request<T = any>(config: any): Promise<T> {
  return instance(config) as Promise<T>
}
