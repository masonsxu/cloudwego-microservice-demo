<template>
  <div class="flex h-full flex-col">
    <!-- Logo 区：与 Topbar 等高（56px） -->
    <div
      class="flex h-[var(--layout-topbar-height)] flex-shrink-0 items-center border-b border-subtle px-4"
      :class="appStore.sidebarCollapsed ? 'justify-center px-0' : ''"
    >
      <img
        v-if="appStore.theme === 'dark'"
        src="/logo-dark.svg"
        alt="Logo"
        class="h-7 w-auto"
      />
      <img
        v-else
        src="/logo-light.svg"
        alt="Logo"
        class="h-7 w-auto"
      />
    </div>

    <!-- 菜单 -->
    <nav class="flex-1 overflow-y-auto py-3" :class="appStore.sidebarCollapsed ? 'px-2' : 'px-3'">
      <template v-for="item in menuList" :key="item.id">
        <!-- 父项（带子菜单） -->
        <div v-if="item.children && item.children.length > 0" class="mb-0.5">
          <button
            class="group flex w-full items-center gap-2.5 rounded-sm px-2.5 py-1.5 text-left text-[14px] transition-colors duration-[var(--duration-fast)]"
            :class="[
              appStore.sidebarCollapsed ? 'justify-center px-0' : '',
              isMenuBranchActive(item)
                ? 'bg-canvas text-[color:var(--color-primary)] font-semibold'
                : 'text-[color:var(--color-ink-muted)] hover:bg-canvas hover:text-ink'
            ]"
            @click="handleParentMenuClick(item)"
          >
            <component :is="getIconComponent(item.icon)" v-if="item.icon" class="h-4 w-4 flex-shrink-0" />
            <span v-show="!appStore.sidebarCollapsed" class="flex-1 truncate">{{ getMenuLabel(item.name) }}</span>
            <ChevronRight
              v-show="!appStore.sidebarCollapsed && hasNavigableChildren(item)"
              class="h-3.5 w-3.5 flex-shrink-0 text-[color:var(--color-ink-subtle)] transition-transform duration-[var(--duration-fast)]"
              :class="isMenuExpanded(item) ? 'rotate-90' : ''"
            />
          </button>

          <!-- 子菜单（展开时） -->
          <div
            v-show="isMenuExpanded(item) && !appStore.sidebarCollapsed"
            class="mt-0.5 space-y-0.5"
          >
            <template v-for="child in item.children" :key="child.id">
              <RouterLink
                v-if="getMenuRoute(child)"
                :to="getMenuRoute(child)!"
                class="flex h-8 items-center gap-2.5 rounded-sm pl-9 pr-2.5 text-[14px] transition-colors duration-[var(--duration-fast)] relative"
                :class="
                  isMenuActive(child)
                    ? 'bg-canvas text-[color:var(--color-primary)] font-semibold before:content-[\'\'] before:absolute before:left-3 before:top-2 before:bottom-2 before:w-0.5 before:rounded-pill before:bg-[color:var(--color-primary)]'
                    : 'text-[color:var(--color-ink-muted)] hover:bg-canvas hover:text-ink'
                "
              >
                <component :is="getIconComponent(child.icon)" v-if="child.icon" class="h-3.5 w-3.5 flex-shrink-0" />
                <span class="truncate">{{ getMenuLabel(child.name) }}</span>
              </RouterLink>
              <div
                v-else
                class="flex h-8 cursor-not-allowed items-center gap-2.5 pl-9 pr-2.5 text-[14px] text-[color:var(--color-ink-disabled)]"
              >
                <component :is="getIconComponent(child.icon)" v-if="child.icon" class="h-3.5 w-3.5 flex-shrink-0" />
                <span class="truncate">{{ getMenuLabel(child.name) }}</span>
              </div>
            </template>
          </div>
        </div>

        <!-- 单项（无子菜单） -->
        <RouterLink
          v-else-if="getMenuRoute(item)"
          :to="getMenuRoute(item)!"
          class="mb-0.5 flex h-8 items-center gap-2.5 rounded-sm px-2.5 text-[14px] transition-colors duration-[var(--duration-fast)] relative"
          :class="[
            appStore.sidebarCollapsed ? 'justify-center px-0' : '',
            isMenuActive(item)
              ? 'bg-canvas text-[color:var(--color-primary)] font-semibold' + (!appStore.sidebarCollapsed ? ' before:content-[\'\'] before:absolute before:left-0 before:top-2 before:bottom-2 before:w-0.5 before:rounded-pill before:bg-[color:var(--color-primary)]' : '')
              : 'text-[color:var(--color-ink-muted)] hover:bg-canvas hover:text-ink'
          ]"
        >
          <component :is="getIconComponent(item.icon)" v-if="item.icon" class="h-4 w-4 flex-shrink-0" />
          <span v-show="!appStore.sidebarCollapsed" class="truncate">{{ getMenuLabel(item.name) }}</span>
        </RouterLink>

        <!-- 不可访问项（占位） -->
        <div
          v-else
          class="mb-0.5 flex h-8 cursor-not-allowed items-center gap-2.5 rounded-sm px-2.5 text-[14px] text-[color:var(--color-ink-disabled)]"
          :class="appStore.sidebarCollapsed ? 'justify-center px-0' : ''"
        >
          <component :is="getIconComponent(item.icon)" v-if="item.icon" class="h-4 w-4 flex-shrink-0" />
          <span v-show="!appStore.sidebarCollapsed" class="truncate">{{ getMenuLabel(item.name) }}</span>
        </div>
      </template>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import {
  Settings,
  Building2,
  Key,
  User,
  FileText,
  LayoutDashboard,
  Network,
  ChevronRight
} from 'lucide-vue-next'
import type { MenuNodeDTO } from '@/api/role'
import { menuRouteMap, visibleRoutePaths } from '@/router/routes'
import type { AppRouteMeta } from '@/router/routes'

const appStore = useAppStore()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { t } = useI18n()

const expandedMenus = ref<Set<string>>(new Set())

const menuList = computed(() => authStore.menuTree || [])

const activeMenuId = computed(() => {
  const meta = route.meta as AppRouteMeta
  return meta.activeMenu || meta.menuId || ''
})

const activeRoutePath = computed(() => {
  const meta = route.meta as AppRouteMeta

  if (meta.activeMenu) {
    return menuRouteMap.get(meta.activeMenu)?.path || route.path
  }

  return route.path
})

function toggleSubMenu(id: string) {
  if (expandedMenus.value.has(id)) {
    expandedMenus.value.delete(id)
    return
  }

  expandedMenus.value.add(id)
}

function getMenuRoute(menu: MenuNodeDTO): string | null {
  const mappedRoute = menuRouteMap.get(menu.id)?.path
  if (mappedRoute) {
    return mappedRoute
  }

  if (menu.path && visibleRoutePaths.has(menu.path)) {
    return menu.path
  }

  return null
}

function getMenuLabel(name?: string): string {
  if (!name) {
    return ''
  }

  return name.includes('.') ? t(name) : name
}

function hasNavigableChildren(menu: MenuNodeDTO): boolean {
  return !!menu.children?.some(child => getMenuRoute(child))
}

function isMenuActive(menu: MenuNodeDTO): boolean {
  const target = getMenuRoute(menu)

  if (menu.id && activeMenuId.value) {
    return menu.id === activeMenuId.value
  }

  return !!target && activeRoutePath.value === target
}

function isMenuBranchActive(menu: MenuNodeDTO): boolean {
  if (isMenuActive(menu)) {
    return true
  }

  return !!menu.children?.some(child => isMenuActive(child))
}

function isMenuExpanded(menu: MenuNodeDTO): boolean {
  return isMenuBranchActive(menu) || expandedMenus.value.has(menu.id)
}

async function handleParentMenuClick(menu: MenuNodeDTO) {
  if (appStore.sidebarCollapsed) {
    appStore.toggleSidebar()
  }

  const target = getMenuRoute(menu)

  if (target && target !== route.path) {
    await router.push(target)
    return
  }

  if (hasNavigableChildren(menu)) {
    toggleSubMenu(menu.id)
  }
}

function getIconComponent(iconName: string) {
  const iconMap: Record<string, any> = {
    IconSystemSettings: Settings,
    IconOrganizationManagement: Building2,
    IconRolePermissions: Key,
    IconAccountManagement: User,
    IconAuditLogs: FileText,
    Odometer: LayoutDashboard,
    User,
    OfficeBuilding: Building2,
    Key,
    Settings,
    Document: FileText,
    Connection: Network,
    Network,
    Menu: LayoutDashboard
  }

  return iconMap[iconName] || LayoutDashboard
}
</script>
