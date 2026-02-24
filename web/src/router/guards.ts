import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

export function setupRouterGuard(router: Router) {
  // 前置守卫
  router.beforeEach((to, from, next) => {
    // 设置页面标题
    document.title = `${to.meta.title || 'CloudWeGo'} - ${import.meta.env.VITE_APP_TITLE}`

    const authStore = useAuthStore()

    // 检查是否需要认证
    if (to.meta.requiresAuth !== false) {
      if (!authStore.isAuthenticated) {
        ElMessage.warning('请先登录')
        next({
          path: '/login',
          query: { redirect: to.fullPath }
        })
        return
      }

      // 检查菜单权限（如果路由定义了 menuId）
      if (to.meta.menuId && !authStore.hasMenuPermission(to.meta.menuId as string, to.meta.permissionLevel as any)) {
        ElMessage.error('没有访问权限')
        next(from.path)
        return
      }

      // 检查普通权限（兼容旧代码）
      if (to.meta.permission && !authStore.hasPermission(to.meta.permission as string)) {
        ElMessage.error('没有访问权限')
        next(from.path)
        return
      }
    }

    // 如果已登录，访问登录页则跳转到首页
    if (to.path === '/login' && authStore.isAuthenticated) {
      next('/dashboard')
      return
    }

    next()
  })

  // 后置守卫
  router.afterEach((to, from) => {
    // 可以在这里添加页面访问统计等逻辑
    console.log(`Navigation: ${from.path} -> ${to.path}`)
  })
}
