<template>
  <div class="flex h-full w-full items-center justify-between gap-4 px-[var(--layout-content-px)]">
    <!-- 左：折叠按钮 + 面包屑 -->
    <div class="flex items-center gap-3 min-w-0 flex-1">
      <button
        class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-sm text-[color:var(--color-ink-muted)] transition-colors duration-[var(--duration-fast)] hover:bg-sunken hover:text-ink"
        :title="appStore.sidebarCollapsed ? t('common.expandSidebar') : t('common.collapseSidebar')"
        @click="toggleSidebar"
      >
        <PanelLeftClose v-if="!appStore.sidebarCollapsed" class="h-4 w-4" />
        <PanelLeftOpen v-else class="h-4 w-4" />
      </button>
      <AppBreadcrumb />
    </div>

    <!-- 中：全局搜索 omnibar（⌘K 占位，搜索功能后续实装） -->
    <button
      class="hidden md:flex h-8 w-[280px] items-center gap-2 rounded-sm bg-sunken px-3 text-[13px] text-[color:var(--color-ink-subtle)] transition-colors duration-[var(--duration-fast)] hover:text-[color:var(--color-ink-muted)]"
      :title="t('common.searchPlaceholder')"
      @click="handleSearchClick"
    >
      <Search class="h-3.5 w-3.5 flex-shrink-0" />
      <span class="flex-1 text-left truncate">{{ t('common.searchPlaceholder') }}</span>
      <kbd
        class="inline-flex h-5 items-center gap-0.5 rounded-xs border border-default bg-canvas px-1.5 text-[11px] font-medium text-[color:var(--color-ink-muted)]"
      >
        <span>⌘</span>
        <span>K</span>
      </kbd>
    </button>

    <!-- 右：主题 / 语言 / 用户 -->
    <div class="flex items-center gap-1 flex-shrink-0">
      <button
        class="flex h-8 w-8 items-center justify-center rounded-sm text-[color:var(--color-ink-muted)] transition-colors duration-[var(--duration-fast)] hover:bg-sunken hover:text-ink"
        :title="appStore.theme === 'dark' ? t('common.switchToLight') : t('common.switchToDark')"
        @click="appStore.toggleTheme"
      >
        <Sun v-if="appStore.theme === 'dark'" class="h-4 w-4" />
        <Moon v-else class="h-4 w-4" />
      </button>

      <DropdownMenu>
        <DropdownMenuTrigger
          class="flex h-8 items-center gap-1.5 rounded-sm px-2 text-[13px] text-[color:var(--color-ink-muted)] transition-colors duration-[var(--duration-fast)] hover:bg-sunken hover:text-ink"
        >
          <Languages class="h-4 w-4" />
          <span class="hidden sm:inline">{{ currentLanguage }}</span>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" class="min-w-[140px]">
          <DropdownMenuItem @select="handleLanguageCommand('zh-CN')">
            简体中文
          </DropdownMenuItem>
          <DropdownMenuItem @select="handleLanguageCommand('en-US')">
            English
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <DropdownMenu>
        <DropdownMenuTrigger
          class="flex h-8 items-center gap-2 rounded-sm pl-1.5 pr-2 text-[13px] text-[color:var(--color-ink)] transition-colors duration-[var(--duration-fast)] hover:bg-sunken"
        >
          <span
            class="flex h-6 w-6 items-center justify-center rounded-full bg-[color:var(--color-primary-soft-strong)] text-[11px] font-semibold text-[color:var(--color-primary-active)]"
          >
            {{ usernameInitial }}
          </span>
          <span class="hidden sm:inline max-w-[120px] truncate">{{ authStore.username }}</span>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" class="min-w-[180px]">
          <DropdownMenuItem @select="handleUserCommand('profile')">
            <User class="mr-2 h-4 w-4" />
            {{ t('common.edit') }}
          </DropdownMenuItem>
          <DropdownMenuItem @select="handleUserCommand('password')">
            <Lock class="mr-2 h-4 w-4" />
            {{ t('auth.changePassword') }}
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem class="text-[color:var(--color-danger)]" @select="handleUserCommand('logout')">
            <LogOut class="mr-2 h-4 w-4" />
            {{ t('auth.logout') }}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  </div>

  <Dialog v-model:open="passwordDialogOpen">
    <DialogContent class="max-w-[420px]">
      <DialogHeader>
        <DialogTitle>{{ t('auth.changePassword') }}</DialogTitle>
      </DialogHeader>
      <div class="space-y-4 py-4">
        <div class="space-y-1.5">
          <Label>{{ t('auth.oldPassword') }}</Label>
          <Input v-model="passwordForm.oldPassword" type="password" :placeholder="t('auth.oldPassword')" />
        </div>
        <div class="space-y-1.5">
          <Label>{{ t('auth.newPassword') }}</Label>
          <Input v-model="passwordForm.newPassword" type="password" :placeholder="t('auth.newPassword')" />
        </div>
        <div class="space-y-1.5">
          <Label>{{ t('auth.confirmPassword') }}</Label>
          <Input v-model="passwordForm.confirmPassword" type="password" :placeholder="t('auth.confirmPassword')" />
        </div>
      </div>
      <DialogFooter>
        <Button variant="outline" @click="handlePasswordDialogClose">{{ t('common.cancel') }}</Button>
        <Button :disabled="passwordSubmitting" @click="handlePasswordSubmit">{{ t('common.confirm') }}</Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { toast } from 'vue-sonner'
import { authApi } from '@/api/auth'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  PanelLeftClose,
  PanelLeftOpen,
  Languages,
  User,
  Lock,
  LogOut,
  Sun,
  Moon,
  Search
} from 'lucide-vue-next'
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator
} from '@/components/ui/dropdown-menu'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import AppBreadcrumb from './AppBreadcrumb.vue'

const appStore = useAppStore()
const authStore = useAuthStore()
const router = useRouter()
const { t, locale } = useI18n()

const passwordDialogOpen = ref(false)
const passwordSubmitting = ref(false)
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const currentLanguage = computed(() => {
  return locale.value === 'zh-CN' ? '简中' : 'EN'
})

const usernameInitial = computed(() => {
  const name = authStore.username || ''
  return name.charAt(0).toUpperCase() || '?'
})

function toggleSidebar() {
  appStore.toggleSidebar()
}

function handleLanguageCommand(command: string) {
  appStore.setLanguage(command)
  locale.value = command
}

function handleSearchClick() {
  toast.info(t('common.searchComingSoon'))
}

function resetPasswordForm() {
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
}

function handlePasswordDialogClose() {
  passwordDialogOpen.value = false
  resetPasswordForm()
}

async function handlePasswordSubmit() {
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    toast.error(t('auth.passwordNotMatch'))
    return
  }

  passwordSubmitting.value = true
  try {
    await authApi.changePassword({
      old_password: passwordForm.oldPassword,
      new_password: passwordForm.newPassword
    })
    toast.success(t('auth.passwordChanged'))
    handlePasswordDialogClose()
  } catch (error: any) {
    toast.error(error?.message || t('common.operationFailed'))
  } finally {
    passwordSubmitting.value = false
  }
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
      passwordDialogOpen.value = true
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
