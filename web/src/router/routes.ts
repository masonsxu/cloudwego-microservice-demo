import type { RouteRecordRaw } from 'vue-router'

export interface AppRouteMeta {
  title?: string
  requiresAuth?: boolean
  icon?: string
  menuId?: string
  hidden?: boolean
  activeMenu?: string
  breadcrumb?: boolean
}

export interface MenuRouteDefinition {
  menuId: string
  path: string
  name?: string | symbol
  title?: string
  icon?: string
  parentMenuId?: string
}

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { requiresAuth: false, title: 'auth.login' } satisfies AppRouteMeta
  },
  {
    path: '/signup',
    name: 'Signup',
    component: () => import('@/views/auth/Signup.vue'),
    meta: { requiresAuth: false, title: 'auth.signup' } satisfies AppRouteMeta
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/auth/ForgotPassword.vue'),
    meta: { requiresAuth: false, title: 'auth.forgotPassword' } satisfies AppRouteMeta
  },
  {
    path: '/',
    component: () => import('@/components/layout/Layout.vue'),
    meta: { requiresAuth: true, breadcrumb: false } satisfies AppRouteMeta,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Dashboard.vue'),
        meta: { title: 'dashboard.title', icon: 'Odometer', breadcrumb: false } satisfies AppRouteMeta
      },
      {
        path: 'system-settings',
        name: 'SystemSettings',
        meta: {
          title: 'system.title',
          icon: 'Setting',
          menuId: 'system_settings'
        } satisfies AppRouteMeta,
        redirect: '/system-settings/organization',
        children: [
          {
            path: 'organization',
            name: 'OrganizationManagement',
            component: () => import('@/views/organization/OrgList.vue'),
            meta: {
              title: 'organization.title',
              icon: 'OfficeBuilding',
              menuId: 'organization_management'
            } satisfies AppRouteMeta
          },
          {
            path: 'roles',
            name: 'RoleManagement',
            component: () => import('@/views/role/RoleList.vue'),
            meta: {
              title: 'role.title',
              icon: 'Key',
              menuId: 'role_permissions'
            } satisfies AppRouteMeta
          },
          {
            path: 'accounts',
            name: 'AccountManagement',
            component: () => import('@/views/user/UserList.vue'),
            meta: {
              title: 'user.title',
              icon: 'User',
              menuId: 'account_management'
            } satisfies AppRouteMeta
          },
          {
            path: 'audit-logs',
            name: 'AuditLogs',
            component: () => import('@/views/audit/AuditLog.vue'),
            meta: {
              title: 'audit.title',
              icon: 'Document',
              menuId: 'audit_logs'
            } satisfies AppRouteMeta
          },
          {
            path: 'oidc',
            name: 'OIDCManagement',
            meta: {
              title: 'oidc.title',
              icon: 'Connection',
              menuId: 'oidc_management'
            } satisfies AppRouteMeta,
            redirect: '/system-settings/oidc/config',
            children: [
              {
                path: 'config',
                name: 'OIDCConfig',
                component: () => import('@/views/oidc/ConfigDetail.vue'),
                meta: {
                  title: 'oidc.config.title',
                  menuId: 'oidc_management'
                } satisfies AppRouteMeta
              },
              {
                path: 'integration',
                name: 'OIDCIntegration',
                component: () => import('@/views/oidc/IntegrationGuide.vue'),
                meta: {
                  title: 'oidc.integration.title',
                  menuId: 'oidc_management'
                } satisfies AppRouteMeta
              }
            ]
          }
        ]
      },
      {
        path: 'users',
        name: 'UserList',
        component: () => import('@/views/user/UserList.vue'),
        meta: {
          title: 'user.title',
          hidden: true,
          activeMenu: 'account_management'
        } satisfies AppRouteMeta
      },
      {
        path: 'users/create',
        name: 'UserCreate',
        component: () => import('@/views/user/UserCreate.vue'),
        meta: {
          title: 'user.createUser',
          hidden: true,
          activeMenu: 'account_management'
        } satisfies AppRouteMeta
      },
      {
        path: 'users/:id',
        name: 'UserDetail',
        component: () => import('@/views/user/UserDetail.vue'),
        meta: {
          title: 'user.userDetail',
          hidden: true,
          activeMenu: 'account_management'
        } satisfies AppRouteMeta
      },
      {
        path: 'users/:id/edit',
        name: 'UserEdit',
        component: () => import('@/views/user/UserEdit.vue'),
        meta: {
          title: 'user.editUser',
          hidden: true,
          activeMenu: 'account_management'
        } satisfies AppRouteMeta
      },
      {
        path: 'organizations/:id',
        name: 'OrgDetail',
        component: () => import('@/views/organization/OrgDetail.vue'),
        meta: {
          title: 'organization.detail',
          hidden: true,
          activeMenu: 'organization_management'
        } satisfies AppRouteMeta
      },
      {
        path: 'roles/:id',
        name: 'RoleDetail',
        component: () => import('@/views/role/RoleDetail.vue'),
        meta: {
          title: 'role.roleDetail',
          hidden: true,
          activeMenu: 'role_permissions'
        } satisfies AppRouteMeta
      },
      {
        path: 'system',
        name: 'System',
        meta: {
          title: 'system.title',
          hidden: true,
          activeMenu: 'system_settings'
        } satisfies AppRouteMeta,
        children: [
          {
            path: 'menus',
            name: 'MenuManage',
            component: () => import('@/views/system/MenuManage.vue'),
            meta: {
              title: 'system.menuManage',
              hidden: true,
              activeMenu: 'system_settings'
            } satisfies AppRouteMeta
          }
        ]
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/dashboard'
  }
]

function normalizePath(path: string, parentPath = ''): string {
  if (path.startsWith('/')) {
    return path
  }

  const base = parentPath === '/' ? '' : parentPath.replace(/\/$/, '')
  const current = path.replace(/^\//, '')

  return `${base}/${current}` || '/'
}

function collectMenuRoutes(routeRecords: RouteRecordRaw[], parentPath = '', parentMenuId?: string): MenuRouteDefinition[] {
  return routeRecords.flatMap(route => {
    const fullPath = normalizePath(route.path, parentPath)
    const meta = (route.meta ?? {}) as AppRouteMeta
    const currentParentMenuId = meta.menuId ?? parentMenuId
    const currentRoute: MenuRouteDefinition[] = meta.menuId && !meta.hidden && route.component
      ? [{
          menuId: meta.menuId,
          path: fullPath,
          name: typeof route.name === 'string' || typeof route.name === 'symbol' ? route.name : undefined,
          title: meta.title,
          icon: meta.icon,
          parentMenuId
        }]
      : []

    const children: MenuRouteDefinition[] = route.children
      ? collectMenuRoutes(route.children, fullPath, currentParentMenuId)
      : []

    return currentRoute.concat(children)
  })
}

export const menuRouteDefinitions = collectMenuRoutes(routes)

export const menuRouteMap = new Map(menuRouteDefinitions.map(route => [route.menuId, route]))

function collectVisiblePaths(routeRecords: RouteRecordRaw[], parentPath = ''): string[] {
  return routeRecords.flatMap(route => {
    const fullPath = normalizePath(route.path, parentPath)
    const meta = (route.meta ?? {}) as AppRouteMeta
    const currentRoute = !meta.hidden && route.component ? [fullPath] : []
    const children = route.children ? collectVisiblePaths(route.children, fullPath) : []

    return currentRoute.concat(children)
  })
}

export const visibleRoutePaths = new Set(collectVisiblePaths(routes))

export default routes
