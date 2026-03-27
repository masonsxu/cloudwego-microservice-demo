import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'
import { clearApiClient } from '@/api/config'
import type { UserProfile, LoginRequest } from '@/types/user'
import type { MenuItem, MenuPermission, Role } from '@/types/user'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  // Cookie 方案：token 由浏览器自动管理，前端不存储
  const isLoggedIn = ref<boolean>(false)
  const user = ref<UserProfile | null>(null)
  const menuTree = ref<MenuItem[]>([])
  const menuPermissions = ref<MenuPermission[]>([])
  const roles = ref<Role[]>([])

  // 计算属性
  const isAuthenticated = computed(() => isLoggedIn.value)
  const username = computed(() => user.value?.username || '')
  const userId = computed(() => user.value?.id || '')

  // 检查菜单权限
  function hasMenuPermission(menuId: string, requiredLevel: 'none' | 'read' | 'write' | 'full' = 'read'): boolean {
    const perm = menuPermissions.value.find(p => p.menu_id === menuId)
    if (!perm) return false

    // perm.permission 现在是数字类型（0, 1, 2, 3）
    const permLevel = perm.permission as number
    const requiredLevels = { none: 0, read: 1, write: 2, full: 3 }

    return permLevel >= requiredLevels[requiredLevel]
  }

  // 检查是否有任何权限（兼容旧代码）
  function hasPermission(_permission: string): boolean {
    // 暂时返回 true，让所有请求通过
    // TODO: 根据实际业务逻辑实现基于菜单的权限检查
    return isAuthenticated.value
  }

  // 检查多个权限（满足任一即可）
  function hasAnyPermission(permissionList: string[]): boolean {
    return permissionList.some(permission => hasPermission(permission))
  }

  // 检查多个权限（必须全部满足）
  function hasAllPermissions(permissionList: string[]): boolean {
    return permissionList.every(permission => hasPermission(permission))
  }

  // 登录
  async function login(credentials: LoginRequest) {
    try {
      const response = await authApi.login(credentials)

      // Cookie 方案：后端自动设置 HttpOnly Cookie，前端不处理 token
      isLoggedIn.value = true

      const userProfile = (response.user_profile ?? null) as UserProfile | null
      const menuTreeData = (response.menu_tree ?? []) as MenuItem[]
      const permissionsData = (response.permissions ?? []) as MenuPermission[]
      const rolesData = (response.roles ?? []) as Role[]

      user.value = userProfile
      menuTree.value = menuTreeData
      menuPermissions.value = permissionsData
      roles.value = rolesData

      // 保存用户相关信息到本地存储（token 存在 Cookie 中）
      localStorage.setItem('user', JSON.stringify(userProfile))
      localStorage.setItem('menuTree', JSON.stringify(menuTreeData))
      localStorage.setItem('menuPermissions', JSON.stringify(permissionsData))

      return response
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  // 清除本地认证状态（不调用后端接口）
  // 用于 token 刷新失败等场景，避免发起新的 HTTP 请求造成循环
  function clearAuthState() {
    isLoggedIn.value = false
    user.value = null
    menuTree.value = []
    menuPermissions.value = []
    roles.value = []
    localStorage.removeItem('user')
    localStorage.removeItem('menuTree')
    localStorage.removeItem('menuPermissions')
    clearApiClient()
  }

  // 登出
  async function logout() {
    try {
      // Cookie 方案：调用后端登出接口清除 HttpOnly Cookie
      await authApi.logout()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      clearAuthState()
    }
  }

  // 刷新访问令牌（Cookie 方案）
  // 浏览器自动携带 HttpOnly Cookie，后端读取 Cookie 中的 token 进行刷新
  async function refreshAccessToken() {
    try {
      const response = await authApi.refreshToken()

      if (response.token_info) {
        // Cookie 方案：后端会设置新的 Cookie，前端不需要存储 token
        isLoggedIn.value = true
        return response.token_info.access_token
      }
      throw new Error('Token refresh failed: invalid response')
    } catch (error) {
      console.error('Token refresh failed:', error)
      // 仅清除本地状态，不调用登出接口，避免产生循环请求
      clearAuthState()
      throw error
    }
  }

  // 从本地存储恢复状态
  // Cookie 方案：token 存在 HttpOnly Cookie 中，isLoggedIn 会在首次 API 调用时确认
  function restoreState() {
    const savedUser = localStorage.getItem('user')
    const savedMenuTree = localStorage.getItem('menuTree')
    const savedPermissions = localStorage.getItem('menuPermissions')

    // 如果有保存的用户数据，假设已登录（实际状态会在首次 API 调用时验证）
    if (savedUser) {
      try {
        user.value = JSON.parse(savedUser)
        isLoggedIn.value = true
      } catch (error) {
        console.error('Failed to parse user from localStorage:', error)
        isLoggedIn.value = false
      }
    }

    if (savedMenuTree) {
      try {
        menuTree.value = JSON.parse(savedMenuTree)
      } catch (error) {
        console.error('Failed to parse menuTree from localStorage:', error)
      }
    }

    if (savedPermissions) {
      try {
        menuPermissions.value = JSON.parse(savedPermissions)
      } catch (error) {
        console.error('Failed to parse permissions from localStorage:', error)
      }
    }
  }

  // 初始化时恢复状态
  restoreState()

  return {
    // 状态
    isLoggedIn,
    user,
    menuTree,
    menuPermissions,
    roles,

    // 计算属性
    isAuthenticated,
    username,
    userId,

    // 方法
    hasPermission,
    hasMenuPermission,
    hasAnyPermission,
    hasAllPermissions,
    login,
    logout,
    clearAuthState,
    refreshAccessToken,
    restoreState
  }
})
