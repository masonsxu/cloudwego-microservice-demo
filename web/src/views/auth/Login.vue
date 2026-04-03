<template>
  <div class="min-h-screen grid grid-cols-1 lg:grid-cols-[1.1fr_1fr] bg-[var(--bg-base)]">
    <section class="hidden lg:flex flex-col justify-between relative p-12 overflow-hidden bg-gradient-to-br from-[rgba(160,170,190,0.9)] to-[rgba(110,120,140,0.9)] text-white">
      <div class="brand relative z-20 flex items-center">
        <img src="/logo-light.svg" alt="Logo" class="h-10 w-auto" />
      </div>

      <div class="character-stage relative z-20 flex items-end justify-center h-[520px]">
        <AnimatedCharacters
          :is-typing="isTyping"
          :show-password="showPassword"
          :password-length="loginForm.password.length"
        />
      </div>

      <div class="legal relative z-20 flex gap-6 text-[13px] text-[rgba(15,23,42,0.7)]">
        <span>Privacy Policy</span>
        <span>Terms of Service</span>
      </div>

      <div class="absolute inset-0 bg-[radial-gradient(rgba(255,255,255,0.25)_1px,transparent_1px)] bg-[length:20px_20px] opacity-20 z-0" />
      <div class="absolute top-[10%] right-[20%] w-[260px] h-[260px] rounded-full bg-white/60 blur-[80px] opacity-50 z-10" />
      <div class="absolute bottom-[10%] left-[15%] w-[320px] h-[320px] rounded-full bg-white/50 blur-[80px] opacity-50 z-10" />
    </section>

    <section class="flex items-center justify-center p-8 lg:p-12">
      <div class="w-full max-w-[420px]">
        <div class="flex lg:hidden items-center justify-center gap-2.5 font-semibold mb-9">
          <img src="/logo-dark.svg" alt="Logo" class="h-[34px] w-auto" />
        </div>

        <div class="text-center mb-8">
          <h1 class="text-[28px] font-bold mb-1.5">Welcome back!</h1>
          <p class="text-muted-foreground text-sm">请使用你的账号登录</p>
        </div>

        <form class="flex flex-col gap-4.5" @submit.prevent="handleLogin" @keyup.enter="handleLogin">
          <div>
            <label class="block text-[13px] font-semibold text-muted-foreground mb-1.5">{{ t('auth.username') }}</label>
            <Input
              v-model="loginForm.username"
              :placeholder="t('auth.username')"
              class="h-12"
              @focus="isTyping = true"
              @blur="isTyping = false"
            />
            <p v-if="errors.username" class="text-sm text-destructive mt-1">{{ errors.username }}</p>
          </div>

          <div>
            <label class="block text-[13px] font-semibold text-muted-foreground mb-1.5">{{ t('auth.password') }}</label>
            <div class="relative">
              <Input
                v-model="loginForm.password"
                :type="showPassword ? 'text' : 'password'"
                :placeholder="t('auth.password')"
                class="h-12 pr-10"
                @focus="isTyping = true"
                @blur="isTyping = false"
              />
              <button type="button" class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground" @click="togglePassword">
                <Eye v-if="!showPassword" class="h-4 w-4" />
                <EyeOff v-else class="h-4 w-4" />
              </button>
            </div>
            <p v-if="errors.password" class="text-sm text-destructive mt-1">{{ errors.password }}</p>
          </div>

          <div class="flex items-center justify-between text-[13px] text-muted-foreground">
            <div class="flex items-center gap-2">
              <Checkbox id="remember" v-model:checked="rememberMe" />
              <label for="remember" class="text-sm cursor-pointer">Remember for 30 days</label>
            </div>
            <router-link class="text-primary font-semibold" to="/forgot-password">Forgot password?</router-link>
          </div>

          <Button type="button" class="w-full h-12 text-base font-semibold" :disabled="loading" @click="handleLogin">
            <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
            {{ loading ? 'Signing in...' : t('auth.login') }}
          </Button>
        </form>

        <div class="my-5 mb-4 flex items-center justify-center text-xs text-muted-foreground">
          <span>or</span>
        </div>

        <Button variant="outline" class="w-full h-12 border-border/60">
          Log in with Google
        </Button>

        <div class="text-center text-[13px] text-muted-foreground mt-4.5">
          还没有账号？
          <router-link class="text-primary font-semibold" to="/signup">Sign Up</router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import AnimatedCharacters from '@/components/AnimatedCharacters.vue'
import { useAuthStore } from '@/stores/auth'
import { Eye, EyeOff, Loader2 } from 'lucide-vue-next'
import { reactive, ref } from 'vue'
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
const isTyping = ref(false)
const rememberMe = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const errors = reactive<Record<string, string>>({})

const rules: Record<string, Array<{ required?: boolean; min?: number; message: string }>> = {
  username: [{ required: true, message: '请输入用户名' }],
  password: [
    { required: true, message: '请输入密码' },
    { min: 6, message: '密码至少6位' }
  ]
}

function validate(): boolean {
  Object.keys(errors).forEach((k) => delete errors[k])
  let valid = true
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
