import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'
import type { UserProfile, LoginRequest } from '@/types/user'
import type { MenuNodeDTO, RoleInfoDTO, MenuPermissionDTO } from '@/api/generated'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<UserProfile | null>(null)
  const menuTree = ref<MenuNodeDTO[]>([])
  const menuPermissions = ref<MenuPermissionDTO[]>([])
  const roles = ref<RoleInfoDTO[]>([])

  // 计算属性
  const isAuthenticated = computed(() => !!token.value)
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

      token.value = response.token_info?.access_token || ''
      user.value = response.user_profile as unknown as UserProfile
      menuTree.value = (response.menu_tree || []) as unknown as MenuNodeDTO[]
      menuPermissions.value = response.permissions || []
      roles.value = response.roles || []

      // 保存到本地存储
      localStorage.setItem('token', token.value)
      localStorage.setItem('user', JSON.stringify(user.value))
      localStorage.setItem('menuTree', JSON.stringify(response.menu_tree))
      localStorage.setItem('menuPermissions', JSON.stringify(response.permissions))

      return response
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  // 登出
  async function logout() {
    try {
      // 后端不需要 refresh token，直接清除本地状态
      // await authApi.logout(refreshTokenValue.value)
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      // 清除状态
      token.value = ''
      user.value = null
      menuTree.value = []
      menuPermissions.value = []
      roles.value = []

      // 清除本地存储
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      localStorage.removeItem('menuTree')
      localStorage.removeItem('menuPermissions')
    }
  }

  // 刷新访问令牌（暂时不可用，因为后端没有提供 refresh endpoint）
  async function refreshAccessToken() {
    // TODO: 后端如果支持 token 刷新，在这里实现
    throw new Error('Token refresh not supported')
  }

  // 从本地存储恢复状态
  function restoreState() {
    const savedUser = localStorage.getItem('user')
    const savedMenuTree = localStorage.getItem('menuTree')
    const savedPermissions = localStorage.getItem('menuPermissions')

    if (savedUser) {
      try {
        user.value = JSON.parse(savedUser)
      } catch (error) {
        console.error('Failed to parse user from localStorage:', error)
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
    token,
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
    refreshAccessToken,
    restoreState
  }
})
