import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { toast } from 'vue-sonner'
import { i18n } from '@/locales'
import type { AppRouteMeta } from './routes'

const DEFAULT_AUTHED_ROUTE = '/dashboard'

export function setupRouterGuard(router: Router) {
  router.beforeEach((to, _from, next) => {
    const meta = to.meta as AppRouteMeta & { permissionLevel?: 'none' | 'read' | 'write' | 'full' }
    const title = meta.title || 'CloudWeGo'
    const translatedTitle = typeof title === 'string' && title.includes('.')
      ? i18n.global.t(title)
      : title
    document.title = `${translatedTitle} - ${import.meta.env.VITE_APP_TITLE}`

    const authStore = useAuthStore()

    if (to.path === '/login' && authStore.isAuthenticated) {
      next(DEFAULT_AUTHED_ROUTE)
      return
    }

    if (meta.requiresAuth === false) {
      next()
      return
    }

    if (!authStore.isAuthenticated) {
      toast.warning('请先登录')
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
      return
    }

    if (meta.menuId && !authStore.hasMenuPermission(meta.menuId, meta.permissionLevel ?? 'read')) {
      toast.error('没有访问权限')
      next(DEFAULT_AUTHED_ROUTE)
      return
    }

    next()
  })

}
