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
          {{ t('auth.brandSubtitle') }}
        </p>

        <ul class="mt-10 space-y-4">
          <li v-for="point in brandPoints" :key="point.key" class="flex items-start gap-3">
            <span class="mt-0.5 flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-sm bg-canvas text-[color:var(--color-primary)]">
              <component :is="point.icon" class="h-3.5 w-3.5" />
            </span>
            <div>
              <p class="text-[14px] font-semibold text-ink">{{ t(point.title) }}</p>
              <p class="mt-0.5 text-[13px] text-[color:var(--color-ink-muted)]">{{ t(point.desc) }}</p>
            </div>
          </li>
        </ul>
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
        <!-- 移动端 logo（lg 以下显示） -->
        <div class="mb-10 flex justify-center lg:hidden">
          <img src="/logo-light.svg" alt="CloudWeGo" class="h-8 w-auto" />
        </div>

        <header class="mb-8">
          <h1 class="text-[24px] font-semibold leading-tight tracking-[-0.012em] text-[color:var(--color-ink-strong)]">
            {{ t('auth.loginTitle') }}
          </h1>
          <p class="mt-1.5 text-[14px] text-[color:var(--color-ink-muted)]">
            {{ t('auth.loginSubtitle') }}
          </p>
        </header>

        <form class="space-y-5" @submit.prevent="handleLogin">
          <div class="space-y-1.5">
            <label class="block text-[13px] font-semibold text-ink" for="username">
              {{ t('auth.username') }}
            </label>
            <Input
              id="username"
              v-model="loginForm.username"
              :placeholder="t('auth.username')"
              autocomplete="username"
            />
            <p v-if="errors.username" class="text-[12px] text-[color:var(--color-danger-ink)]">{{ errors.username }}</p>
          </div>

          <div class="space-y-1.5">
            <div class="flex items-center justify-between">
              <label class="block text-[13px] font-semibold text-ink" for="password">
                {{ t('auth.password') }}
              </label>
              <router-link
                to="/forgot-password"
                class="text-[12px] text-[color:var(--color-primary)] transition-colors duration-[var(--duration-fast)] hover:text-[color:var(--color-primary-hover)]"
              >
                {{ t('auth.forgotPassword') }}
              </router-link>
            </div>
            <div class="relative">
              <Input
                id="password"
                v-model="loginForm.password"
                :type="showPassword ? 'text' : 'password'"
                :placeholder="t('auth.password')"
                autocomplete="current-password"
                class="pr-10"
              />
              <button
                type="button"
                class="absolute right-2.5 top-1/2 -translate-y-1/2 flex h-6 w-6 items-center justify-center rounded-xs text-[color:var(--color-ink-subtle)] transition-colors duration-[var(--duration-fast)] hover:text-[color:var(--color-ink)]"
                :title="showPassword ? t('auth.hidePassword') : t('auth.showPassword')"
                @click="togglePassword"
              >
                <Eye v-if="!showPassword" class="h-4 w-4" />
                <EyeOff v-else class="h-4 w-4" />
              </button>
            </div>
            <p v-if="errors.password" class="text-[12px] text-[color:var(--color-danger-ink)]">{{ errors.password }}</p>
          </div>

          <label class="flex items-center gap-2 text-[13px] text-[color:var(--color-ink-muted)] cursor-pointer select-none">
            <Checkbox v-model:checked="rememberMe" />
            <span>{{ t('auth.rememberFor30Days') }}</span>
          </label>

          <Button type="submit" class="w-full h-10" :disabled="loading">
            <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
            {{ loading ? t('auth.signingIn') : t('auth.signIn') }}
          </Button>
        </form>

        <p class="mt-8 text-center text-[13px] text-[color:var(--color-ink-muted)]">
          {{ t('auth.noAccount') }}
          <router-link
            to="/signup"
            class="ml-1 text-[color:var(--color-primary)] font-semibold transition-colors duration-[var(--duration-fast)] hover:text-[color:var(--color-primary-hover)]"
          >
            {{ t('auth.signup') }}
          </router-link>
        </p>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { Eye, EyeOff, Loader2, ShieldCheck, Layers, Zap } from 'lucide-vue-next'
import { computed, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const { t } = useI18n()

const loading = ref(false)
const showPassword = ref(false)
const rememberMe = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const errors = reactive<Record<string, string>>({})
const currentYear = computed(() => new Date().getFullYear())

const brandPoints = [
  { key: 'identity', icon: ShieldCheck, title: 'auth.brandPoint1Title', desc: 'auth.brandPoint1Desc' },
  { key: 'microservice', icon: Layers, title: 'auth.brandPoint2Title', desc: 'auth.brandPoint2Desc' },
  { key: 'performance', icon: Zap, title: 'auth.brandPoint3Title', desc: 'auth.brandPoint3Desc' }
]

const rules: Record<string, Array<{ required?: boolean; min?: number; message: string }>> = {
  username: [{ required: true, message: '' }],
  password: [
    { required: true, message: '' },
    { min: 6, message: '' }
  ]
}

function validate(): boolean {
  Object.keys(errors).forEach((k) => delete errors[k])
  let valid = true
  rules.username[0].message = t('auth.usernameRequired')
  rules.password[0].message = t('auth.passwordRequired')
  rules.password[1].message = t('auth.passwordMinLength')

  for (const [field, fieldRules] of Object.entries(rules)) {
    for (const rule of fieldRules) {
      const value = loginForm[field as keyof typeof loginForm]
      if (rule.required && !value) {
        errors[field] = rule.message
        valid = false
        break
      }
      if (rule.min && value && (value as string).length < rule.min) {
        errors[field] = rule.message
        valid = false
        break
      }
    }
  }
  return valid
}

async function handleLogin() {
  if (!validate()) return

  try {
    loading.value = true
    await authStore.login({
      username: loginForm.username,
      password: loginForm.password
    })
    toast.success(t('auth.loginSuccess'))
    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect)
  } catch (error: any) {
    console.error('Login failed:', error)
    toast.error(error.message || t('auth.loginFailed'))
  } finally {
    loading.value = false
  }
}

function togglePassword() {
  showPassword.value = !showPassword.value
}
</script>
