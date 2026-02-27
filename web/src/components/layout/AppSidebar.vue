<template>
  <div class="app-sidebar">
    <div class="logo">
      <h1 class="title">
        <span class="gold-text">♑</span>
        <span v-if="!appStore.sidebarCollapsed" class="text">CloudWeGo</span>
      </h1>
    </div>
    <el-menu
      :default-active="activeMenu"
      :collapse="appStore.sidebarCollapsed"
      :unique-opened="true"
      router
      class="sidebar-menu"
    >
      <template v-for="item in menuList" :key="item.id">
        <el-sub-menu v-if="item.children && item.children.length > 0" :index="getMenuIndex(item)">
          <template #title>
            <el-icon v-if="item.icon">
              <component :is="getIconComponent(item.icon)" />
            </el-icon>
            <span>{{ t(item.name || '') }}</span>
          </template>
          <template v-for="child in item.children" :key="child.id">
            <el-menu-item :index="getMenuIndex(child)">
              <el-icon v-if="child.icon">
                <component :is="getIconComponent(child.icon)" />
              </el-icon>
              <span>{{ t(child.name || '') }}</span>
            </el-menu-item>
          </template>
        </el-sub-menu>
        <el-menu-item v-else :index="getMenuIndex(item)">
          <el-icon v-if="item.icon">
            <component :is="getIconComponent(item.icon)" />
          </el-icon>
          <span>{{ t(item.name || '') }}</span>
        </el-menu-item>
      </template>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import type { MenuNodeDTO } from '@/api/role'

const appStore = useAppStore()
const route = useRoute()
const authStore = useAuthStore()
const { t } = useI18n()

const activeMenu = computed(() => route.path)

// 从后端返回的菜单树生成侧边栏菜单
const menuList = computed(() => {
  const menuTree = authStore.menuTree || []
  console.log('[AppSidebar] menuTree:', JSON.stringify(menuTree, null, 2))
  return menuTree
})

// 路径映射：后端菜单路径 -> 前端路由路径
// 注意：子菜单的 path 是相对路径（如 "organization-management"）
// 需要和父级路径组合后映射
const pathMapping: Record<string, string> = {
  // 绝对路径（顶级菜单）
  '/system-settings': '/system-settings',
  // 相对路径映射（二级菜单）
  'organization-management': '/system-settings/organization',
  'role-permissions': '/system-settings/roles',
  'account-management': '/system-settings/accounts'
}

// 获取菜单索引（使用映射后的路径）
function getMenuIndex(menu: MenuNodeDTO): string {
  const originalPath = menu.path || ''

  // 如果是绝对路径（以 / 开头），直接映射
  if (originalPath.startsWith('/')) {
    return pathMapping[originalPath] || originalPath
  }

  // 如果是相对路径，先尝试直接映射，如果失败则与父路径组合
  const mapped = pathMapping[originalPath]
  if (mapped) {
    return mapped
  }

  // 如果没有映射，使用原始路径
  return originalPath
}

// 图标名称到 Element Plus 图标组件的映射
function getIconComponent(iconName: string) {
  const iconMap: Record<string, string> = {
    'IconSystemSettings': 'Setting',
    'IconOrganizationManagement': 'OfficeBuilding',
    'IconRolePermissions': 'Key',
    'IconAccountManagement': 'User',
    'Odometer': 'Odometer',
    'User': 'User',
    'OfficeBuilding': 'OfficeBuilding',
    'Key': 'Key',
    'Setting': 'Setting',
    'Menu': 'Menu'
  }
  return iconMap[iconName] || 'Menu'
}
</script>

<style scoped lang="scss">
.app-sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;

  .logo {
    padding: 20px;
    display: flex;
    justify-content: center;
    align-items: center;
    border-bottom: 1px solid rgba(212, 175, 55, 0.2);

    .title {
      font-family: 'Cinzel', serif;
      font-size: 24px;
      display: flex;
      align-items: center;
      gap: 10px;

      .gold-text {
        font-size: 32px;
        background: linear-gradient(to right, #D4AF37, #F2F0E4, #D4AF37);
        background-size: 200% auto;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        animation: shine 5s linear infinite;
      }

      .text {
        background: linear-gradient(to right, #D4AF37, #F2F0E4, #D4AF37);
        background-size: 200% auto;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        animation: shine 5s linear infinite;
      }
    }
  }

  .sidebar-menu {
    flex: 1;
    border: none;
    overflow-y: auto;

    &::-webkit-scrollbar {
      width: 4px;
    }

    &::-webkit-scrollbar-track {
      background: transparent;
    }

    &::-webkit-scrollbar-thumb {
      background: rgba(212, 175, 55, 0.3);
      border-radius: 2px;
    }
  }
}

@keyframes shine {
  to {
    background-position: 200% center;
  }
}
</style>
