import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { requiresAuth: false, title: 'auth.login' }
  },
  {
    path: '/signup',
    name: 'Signup',
    component: () => import('@/views/auth/Signup.vue'),
    meta: { requiresAuth: false, title: 'auth.signup' }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/auth/ForgotPassword.vue'),
    meta: { requiresAuth: false, title: 'auth.forgotPassword' }
  },
  {
    path: '/',
    component: () => import('@/components/layout/Layout.vue'),
    meta: { requiresAuth: true },
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Dashboard.vue'),
        meta: { title: 'dashboard.title', icon: 'Odometer' }
      },
      // 系统设置菜单（根据后端实际菜单结构）
      {
        path: 'system-settings',
        name: 'SystemSettings',
        meta: { title: 'system.title', icon: 'Setting', menuId: 'system_settings' },
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
            }
          },
          {
            path: 'roles',
            name: 'RoleManagement',
            component: () => import('@/views/role/RoleList.vue'),
            meta: {
              title: 'role.title',
              icon: 'Key',
              menuId: 'role_permissions'
            }
          },
          {
            path: 'accounts',
            name: 'AccountManagement',
            component: () => import('@/views/user/UserList.vue'),
            meta: {
              title: 'user.title',
              icon: 'User',
              menuId: 'account_management'
            }
          },
          {
            path: 'audit-logs',
            name: 'AuditLogs',
            component: () => import('@/views/audit/AuditLog.vue'),
            meta: {
              title: 'audit.title',
              icon: 'Document',
              menuId: 'audit_logs'
            }
          },
          {
            path: 'oauth2',
            name: 'OAuth2Management',
            meta: {
              title: 'oauth2.title',
              icon: 'Connection',
              menuId: 'oauth2_management'
            },
            redirect: '/system-settings/oauth2/config',
            children: [
              {
                path: 'config',
                name: 'OAuth2Config',
                component: () => import('@/views/oauth2/ConfigDetail.vue'),
                meta: {
                  title: 'oauth2.config.title',
                  menuId: 'oauth2_management'
                }
              },
              {
                path: 'clients',
                name: 'OAuth2ClientList',
                component: () => import('@/views/oauth2/ClientList.vue'),
                meta: {
                  title: 'oauth2.client.title',
                  menuId: 'oauth2_management'
                }
              },
              {
                path: 'clients/:id',
                name: 'OAuth2ClientDetail',
                component: () => import('@/views/oauth2/ClientDetail.vue'),
                meta: { title: 'oauth2.client.title', hidden: true, menuId: 'oauth2_management' }
              },
              {
                path: 'consents',
                name: 'OAuth2Consents',
                component: () => import('@/views/oauth2/ConsentList.vue'),
                meta: { title: 'oauth2.consent.title', hidden: true, menuId: 'oauth2_management' }
              }
            ]
          }
        ]
      },
      // 用户管理（保留作为二级页面）
      {
        path: 'users',
        name: 'UserList',
        component: () => import('@/views/user/UserList.vue'),
        meta: { title: 'user.title', icon: 'User', hidden: true }
      },
      {
        path: 'users/:id',
        name: 'UserDetail',
        component: () => import('@/views/user/UserDetail.vue'),
        meta: { title: 'user.userDetail', hidden: true }
      },
      {
        path: 'users/create',
        name: 'UserCreate',
        component: () => import('@/views/user/UserCreate.vue'),
        meta: { title: 'user.createUser', hidden: true }
      },
      {
        path: 'users/:id/edit',
        name: 'UserEdit',
        component: () => import('@/views/user/UserEdit.vue'),
        meta: { title: 'user.editUser', hidden: true }
      },
      {
        path: 'organizations/:id',
        name: 'OrgDetail',
        component: () => import('@/views/organization/OrgDetail.vue'),
        meta: { title: 'organization.detail', hidden: true }
      },
      {
        path: 'roles/:id',
        name: 'RoleDetail',
        component: () => import('@/views/role/RoleDetail.vue'),
        meta: { title: 'role.roleDetail', hidden: true }
      },
      {
        path: 'system',
        name: 'System',
        meta: { title: 'system.title', icon: 'Setting', hidden: true },
        children: [
          {
            path: 'menus',
            name: 'MenuManage',
            component: () => import('@/views/system/MenuManage.vue'),
            meta: { title: 'system.menuManage', hidden: true }
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

export default routes
