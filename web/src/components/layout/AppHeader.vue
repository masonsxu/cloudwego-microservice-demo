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

    .language-selector,
    .user-info {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      color: #8B9bb4;
      transition: color 0.3s ease;

      &:hover {
        color: #D4AF37;
      }
    }
  }
}
</style>
