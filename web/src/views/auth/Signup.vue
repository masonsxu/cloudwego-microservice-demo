<template>
  <div class="min-h-screen grid grid-cols-1 lg:grid-cols-[1.1fr_1fr] bg-canvas">
    <!-- 左侧 brand 区 -->
    <section class="hidden lg:flex flex-col justify-between bg-sunken px-12 py-10 border-r border-subtle">
      <div class="flex items-center">
        <img src="/logo-light.svg" alt="CloudWeGo" class="h-8 w-auto" />
      </div>
      <div class="max-w-[420px]">
        <h2 class="text-[36px] font-semibold leading-[1.15] tracking-[-0.02em] text-[color:var(--color-ink-strong)]">
          {{ t('auth.brandTitle') }}
        </h2>
        <p class="mt-4 text-[16px] leading-relaxed text-[color:var(--color-ink-muted)]">
          {{ t('auth.signupBrandSubtitle') }}
        </p>
      </div>
      <div class="flex gap-6 text-[12px] text-[color:var(--color-ink-subtle)]">
        <span>© {{ currentYear }} CloudWeGo</span>
        <a href="#" class="transition-colors duration-[var(--duration-fast)] hover:text-[color:var(--color-ink-muted)]">{{ t('auth.privacy') }}</a>
        <a href="#" class="transition-colors duration-[var(--duration-fast)] hover:text-[color:var(--color-ink-muted)]">{{ t('auth.terms') }}</a>
      </div>
    </section>

    <!-- 右侧表单 -->
    <section class="flex items-center justify-center px-6 py-12">
      <div class="w-full max-w-[380px]">
        <div class="mb-10 flex justify-center lg:hidden">
          <img src="/logo-light.svg" alt="CloudWeGo" class="h-8 w-auto" />
        </div>

        <header class="mb-8">
          <h1 class="text-[24px] font-semibold leading-tight tracking-[-0.012em] text-[color:var(--color-ink-strong)]">
            {{ t('auth.signupTitle') }}
          </h1>
          <p class="mt-1.5 text-[14px] text-[color:var(--color-ink-muted)]">
            {{ t('auth.signupSubtitle') }}
          </p>
        </header>

        <form class="space-y-5" @submit.prevent="handleSubmit">
          <div class="space-y-1.5">
            <label class="block text-[13px] font-semibold text-ink">{{ t('auth.username') }}</label>
            <Input v-model="form.username" :placeholder="t('auth.username')" autocomplete="username" />
            <p v-if="errors.username" class="text-[12px] text-[color:var(--color-danger-ink)]">{{ errors.username }}</p>
          </div>

          <div class="space-y-1.5">
            <label class="block text-[13px] font-semibold text-ink">{{ t('user.email') }}</label>
            <Input v-model="form.email" placeholder="you@example.com" autocomplete="email" />
            <p v-if="errors.email" class="text-[12px] text-[color:var(--color-danger-ink)]">{{ errors.email }}</p>
          </div>

          <div class="space-y-1.5">
            <label class="block text-[13px] font-semibold text-ink">{{ t('auth.password') }}</label>
            <div class="relative">
              <Input
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="••••••••"
                autocomplete="new-password"
                class="pr-10"
              />
              <button
                type="button"
                class="absolute right-2.5 top-1/2 -translate-y-1/2 flex h-6 w-6 items-center justify-center rounded-xs text-[color:var(--color-ink-subtle)] transition-colors duration-[var(--duration-fast)] hover:text-[color:var(--color-ink)]"
                @click="togglePassword"
              >
                <Eye v-if="!showPassword" class="h-4 w-4" />
                <EyeOff v-else class="h-4 w-4" />
              </button>
            </div>
            <p v-if="errors.password" class="text-[12px] text-[color:var(--color-danger-ink)]">{{ errors.password }}</p>
          </div>

          <div class="space-y-1.5">
            <label class="block text-[13px] font-semibold text-ink">{{ t('auth.confirmPassword') }}</label>
            <Input
              v-model="form.confirmPassword"
              :type="showPassword ? 'text' : 'password'"
              placeholder="••••••••"
              autocomplete="new-password"
            />
            <p v-if="errors.confirmPassword" class="text-[12px] text-[color:var(--color-danger-ink)]">{{ errors.confirmPassword }}</p>
          </div>

          <Button type="submit" class="w-full h-10" :disabled="loading">
            <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
            {{ loading ? t('auth.creatingAccount') : t('auth.createAccount') }}
          </Button>
        </form>

        <p class="mt-8 text-center text-[13px] text-[color:var(--color-ink-muted)]">
          {{ t('auth.haveAccount') }}
          <router-link
            to="/login"
            class="ml-1 text-[color:var(--color-primary)] font-semibold transition-colors duration-[var(--duration-fast)] hover:text-[color:var(--color-primary-hover)]"
          >
            {{ t('auth.login') }}
          </router-link>
        </p>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Eye, EyeOff, Loader2 } from 'lucide-vue-next'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { toast } from 'vue-sonner'

const router = useRouter()
const { t } = useI18n()

const loading = ref(false)
const showPassword = ref(false)
const currentYear = computed(() => new Date().getFullYear())

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const errors = reactive<Record<string, string>>({})

function validate(): boolean {
  Object.keys(errors).forEach((k) => delete errors[k])
  let valid = true

  if (!form.username) {
    errors.username = t('auth.usernameRequired')
    valid = false
  }

  if (!form.email) {
    errors.email = t('auth.emailRequired')
    valid = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = t('auth.emailInvalid')
    valid = false
  }

  if (!form.password) {
    errors.password = t('auth.passwordRequired')
    valid = false
  } else if (form.password.length < 6) {
    errors.password = t('auth.passwordMinLength')
    valid = false
  }

  if (!form.confirmPassword) {
    errors.confirmPassword = t('auth.confirmPasswordRequired')
    valid = false
  } else if (form.confirmPassword !== form.password) {
    errors.confirmPassword = t('auth.passwordNotMatch')
    valid = false
  }

  return valid
}

async function handleSubmit() {
  if (!validate()) return
  try {
    loading.value = true
    toast.info(t('auth.signupComingSoon'))
    router.push('/login')
  } finally {
    loading.value = false
  }
}

function togglePassword() {
  showPassword.value = !showPassword.value
}
</script>
