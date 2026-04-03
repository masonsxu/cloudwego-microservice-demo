<template>
  <div class="flex h-[60px] w-full items-center justify-between px-5">
    <div class="flex items-center gap-5">
      <button
        class="text-[var(--c-accent)] transition-transform hover:scale-110 cursor-pointer"
        @click="toggleSidebar"
      >
        <PanelLeftClose v-if="!appStore.sidebarCollapsed" class="h-5 w-5" />
        <PanelLeftOpen v-else class="h-5 w-5" />
      </button>
      <AppBreadcrumb />
    </div>
    <div class="flex items-center gap-5">
      <!-- 主题切换 -->
      <button class="flex h-8 w-8 items-center justify-center rounded-lg border border-[var(--c-border-accent)] bg-transparent text-[var(--c-text-sub)] transition-all hover:border-[var(--c-accent)] hover:bg-[rgba(212,175,55,0.08)] hover:text-[var(--c-accent)] cursor-pointer p-0" @click="appStore.toggleTheme" :title="appStore.theme === 'dark' ? '切换浅色模式' : '切换深色模式'">
        <svg v-if="appStore.theme === 'dark'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="h-4 w-4">
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
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="h-4 w-4">
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
        </svg>
      </button>

      <!-- 语言切换 -->
      <DropdownMenu>
        <DropdownMenuTrigger class="flex items-center gap-2 text-[var(--c-text-sub)] transition-colors hover:text-[var(--c-accent)] cursor-pointer">
          <Languages class="h-4 w-4" />
          <span>{{ currentLanguage }}</span>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem @select="handleLanguageCommand('zh-CN')">
            简体中文
          </DropdownMenuItem>
          <DropdownMenuItem @select="handleLanguageCommand('en-US')">
            English
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <!-- 用户菜单 -->
      <DropdownMenu>
        <DropdownMenuTrigger class="flex items-center gap-2 text-[var(--c-text-sub)] transition-colors hover:text-[var(--c-accent)] cursor-pointer">
          <User class="h-4 w-4" />
          <span>{{ authStore.username }}</span>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem @select="handleUserCommand('profile')">
            <User class="h-4 w-4" />
            {{ t('common.edit') }}
          </DropdownMenuItem>
          <DropdownMenuItem @select="handleUserCommand('password')">
            <Lock class="h-4 w-4" />
            {{ t('auth.changePassword') }}
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem @select="handleUserCommand('logout')">
            <LogOut class="h-4 w-4" />
            {{ t('auth.logout') }}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { PanelLeftClose, PanelLeftOpen, Languages, User, Lock, LogOut } from 'lucide-vue-next'
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator
} from '@/components/ui/dropdown-menu'
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

async function handleUserCommand(command: string) {
  switch (command) {
    case 'profile':
      if (!authStore.userId) {
        return
      }
      await router.push({ name: 'UserDetail', params: { id: authStore.userId } })
      break
    case 'password':
      break
    case 'logout':
      if (window.confirm(t('auth.logoutConfirm'))) {
        await authStore.logout()
        await router.push({ name: 'Login' })
      }
      break
  }
}
</script>
