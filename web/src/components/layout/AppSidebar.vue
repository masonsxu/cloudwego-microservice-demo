<template>
  <div class="flex flex-col h-full">
    <div class="flex items-center justify-center p-5 border-b border-[rgba(212,175,55,0.2)]">
      <img
        v-if="appStore.theme === 'dark'"
        src="/logo-dark.svg"
        alt="Logo"
        class="h-9 w-auto transition-all duration-300"
        :class="{ 'h-8': appStore.sidebarCollapsed }"
      />
      <img
        v-else
        src="/logo-light.svg"
        alt="Logo"
        class="h-9 w-auto transition-all duration-300"
        :class="{ 'h-8': appStore.sidebarCollapsed }"
      />
    </div>
    <nav class="flex-1 overflow-y-auto [&::-webkit-scrollbar]:w-1 [&::-webkit-scrollbar-track]:bg-transparent [&::-webkit-scrollbar-thumb]:bg-[rgba(212,175,55,0.3)] [&::-webkit-scrollbar-thumb]:rounded-sm">
      <template v-for="item in menuList" :key="item.id">
        <div v-if="item.children && item.children.length > 0" class="relative">
          <button
            class="flex w-full items-center gap-3 px-4 py-2.5 text-sm text-[var(--c-text-sub)] transition-colors hover:bg-[rgba(212,175,55,0.08)] hover:text-[var(--c-accent)]"
            :class="{
              'justify-center px-2': appStore.sidebarCollapsed,
              'bg-[rgba(212,175,55,0.12)] text-[var(--c-accent)]': isMenuBranchActive(item)
            }"
            @click="handleParentMenuClick(item)"
          >
            <component :is="getIconComponent(item.icon)" v-if="item.icon" class="h-4 w-4" />
            <span v-show="!appStore.sidebarCollapsed">{{ getMenuLabel(item.name) }}</span>
          </button>
          <div v-show="isMenuExpanded(item) && !appStore.sidebarCollapsed" class="bg-[rgba(0,0,0,0.1)]">
            <template v-for="child in item.children" :key="child.id">
              <RouterLink
                v-if="getMenuRoute(child)"
                :to="getMenuRoute(child)!"
                class="flex items-center gap-3 pl-12 pr-4 py-2 text-sm text-[var(--c-text-sub)] transition-colors hover:bg-[rgba(212,175,55,0.08)] hover:text-[var(--c-accent)]"
                :class="{ 'bg-[rgba(212,175,55,0.12)] text-[var(--c-accent)]': isMenuActive(child) }"
              >
                <component :is="getIconComponent(child.icon)" v-if="child.icon" class="h-4 w-4" />
                <span>{{ getMenuLabel(child.name) }}</span>
              </RouterLink>
              <div
                v-else
                class="flex items-center gap-3 pl-12 pr-4 py-2 text-sm text-[var(--c-text-sub)] opacity-60 cursor-not-allowed"
              >
                <component :is="getIconComponent(child.icon)" v-if="child.icon" class="h-4 w-4" />
                <span>{{ getMenuLabel(child.name) }}</span>
              </div>
            </template>
          </div>
        </div>
        <RouterLink
          v-else-if="getMenuRoute(item)"
          :to="getMenuRoute(item)!"
          class="flex items-center gap-3 px-4 py-2.5 text-sm text-[var(--c-text-sub)] transition-colors hover:bg-[rgba(212,175,55,0.08)] hover:text-[var(--c-accent)]"
          :class="{
            'justify-center px-2': appStore.sidebarCollapsed,
            'bg-[rgba(212,175,55,0.12)] text-[var(--c-accent)]': isMenuActive(item)
          }"
        >
          <component :is="getIconComponent(item.icon)" v-if="item.icon" class="h-4 w-4" />
          <span v-show="!appStore.sidebarCollapsed">{{ getMenuLabel(item.name) }}</span>
        </RouterLink>
        <div
          v-else
          class="flex items-center gap-3 px-4 py-2.5 text-sm text-[var(--c-text-sub)] opacity-60 cursor-not-allowed"
          :class="{ 'justify-center px-2': appStore.sidebarCollapsed }"
        >
          <component :is="getIconComponent(item.icon)" v-if="item.icon" class="h-4 w-4" />
          <span v-show="!appStore.sidebarCollapsed">{{ getMenuLabel(item.name) }}</span>
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
import { Settings, Building2, Key, User, FileText, Menu, Network } from 'lucide-vue-next'
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
    Odometer: Menu,
    User,
    OfficeBuilding: Building2,
    Key,
    Settings,
    Document: FileText,
    Connection: Network,
    Network,
    Menu
  }

  return iconMap[iconName] || Menu
}
</script>
