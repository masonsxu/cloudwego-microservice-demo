<template>
  <div class="app-header">
    <div class="left">
      <el-icon class="collapse-icon" @click="toggleSidebar">
        <Fold v-if="!appStore.sidebarCollapsed" />
        <Expand v-else />
      </el-icon>
      <AppBreadcrumb />
    </div>
    <div class="right">
      <!-- 主题切换 -->
      <button class="theme-toggle" @click="appStore.toggleTheme" :title="appStore.theme === 'dark' ? '切换浅色模式' : '切换深色模式'">
        <svg v-if="appStore.theme === 'dark'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="5"/>
          <line x1="12" y1="1" x2="12" y2="3"/>
          <line x1="12" y1="21" x2="12" y2="23"/>
          <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
          <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
          <line x1="1" y1="12" x2="3" y2="12"/>
          <line x1="21" y1="12" x2="23" y2="12"/>
          <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
          <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
        </svg>
      </button>

      <el-dropdown @command="handleLanguageCommand">
        <span class="language-selector">
          <el-icon><Connection /></el-icon>
          <span>{{ currentLanguage }}</span>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="zh-CN">简体中文</el-dropdown-item>
            <el-dropdown-item command="en-US">English</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <el-dropdown @command="handleUserCommand">
        <span class="user-info">
          <el-icon><User /></el-icon>
          <span>{{ authStore.username }}</span>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">
              <el-icon><User /></el-icon>
              {{ t('common.edit') }}
            </el-dropdown-item>
            <el-dropdown-item command="password">
              <el-icon><Lock /></el-icon>
              {{ t('auth.changePassword') }}
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>
              {{ t('auth.logout') }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessageBox } from 'element-plus'
import AppBreadcrumb from './AppBreadcrumb.vue'

const appStore = useAppStore()
const authStore = useAuthStore()
const router = useRouter()
const { t, locale } = useI18n()

const currentLanguage = computed(() => {
  return locale.value === 'zh-CN' ? '简体中文' : 'English'
})

function toggleSidebar() {
  appStore.toggleSidebar()
}

function handleLanguageCommand(command: string) {
  appStore.setLanguage(command)
  locale.value = command
}

function handleUserCommand(command: string) {
  switch (command) {
    case 'profile':
      router.push(`/users/${authStore.userId}`)
      break
    case 'password':
      // TODO: 打开修改密码对话框
      break
    case 'logout':
      ElMessageBox.confirm(t('auth.logoutConfirm'), t('common.confirm'), {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }).then(async () => {
        await authStore.logout()
        router.push('/login')
      })
      break
  }
}
</script>

<style scoped lang="scss">
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  width: 100%;

  .left {
    display: flex;
    align-items: center;
    gap: 20px;

    .collapse-icon {
      font-size: 20px;
      color: #D4AF37;
      cursor: pointer;
      transition: transform 0.3s ease;

      &:hover {
        transform: scale(1.1);
      }
    }
  }

  .right {
    display: flex;
    align-items: center;
    gap: 20px;

    .theme-toggle {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 32px;
      height: 32px;
      border-radius: 8px;
      border: 1px solid var(--c-border-accent);
      background: transparent;
      cursor: pointer;
      color: var(--c-text-sub);
      transition: all 0.3s ease;
      padding: 0;

      svg {
        width: 16px;
        height: 16px;
      }

      &:hover {
        color: var(--c-accent);
        border-color: var(--c-accent);
        background: rgba(212, 175, 55, 0.08);
      }
    }

    .language-selector,
    .user-info {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      color: var(--c-text-sub);
      transition: color 0.3s ease;

      &:hover {
        color: var(--c-accent);
      }
    }
  }
}
</style>
