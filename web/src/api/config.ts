/**
 * API 客户端配置（已迁移到 axios，此文件保留以避免破坏现有引用）
 * 实际的 baseURL 和 token 注入在 src/api/request.ts 中通过 axios 拦截器处理
 */
export function initApiClient() {
  // no-op：axios 拦截器已统一处理 baseURL 和 Authorization header
}

export function clearApiClient() {
  // no-op：token 清除由 authStore.logout() 和 localStorage.removeItem 负责
}
