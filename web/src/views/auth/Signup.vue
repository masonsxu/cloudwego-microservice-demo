<template>
  <div class="min-h-screen grid grid-cols-1 lg:grid-cols-[1.1fr_1fr] bg-[var(--bg-base)]">
    <section class="hidden lg:flex flex-col justify-between relative p-12 overflow-hidden bg-gradient-to-br from-[rgba(160,170,190,0.9)] to-[rgba(110,120,140,0.9)] text-white">
      <div class="brand relative z-20 flex items-center">
        <div class="brand-mark">CW</div>
        <span>CloudWeGo</span>
      </div>

      <div class="character-stage relative z-20 flex items-end justify-center h-[520px]">
        <AnimatedCharacters
          :is-typing="isTyping"
          :show-password="showPassword"
          :password-length="form.password.length"
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
          <div class="brand-mark small">CW</div>
          <span>CloudWeGo</span>
        </div>

        <div class="text-center mb-8">
          <h1 class="text-[28px] font-bold mb-1.5">Create account</h1>
          <p class="text-muted-foreground text-sm">快速创建你的管理账号</p>
        </div>

        <form class="flex flex-col gap-4.5" @submit.prevent="handleSubmit" @keyup.enter="handleSubmit">
          <div>
            <label class="block text-[13px] font-semibold text-muted-foreground mb-1.5">用户名</label>
            <Input
              v-model="form.username"
              placeholder="请输入用户名"
              class="h-12"
              @focus="isTyping = true"
              @blur="isTyping = false"
            />
            <p v-if="errors.username" class="text-sm text-destructive mt-1">{{ errors.username }}</p>
          </div>

          <div>
            <label class="block text-[13px] font-semibold text-muted-foreground mb-1.5">邮箱</label>
            <Input
              v-model="form.email"
              placeholder="you@example.com"
              class="h-12"
              @focus="isTyping = true"
              @blur="isTyping = false"
            />
            <p v-if="errors.email" class="text-sm text-destructive mt-1">{{ errors.email }}</p>
          </div>

          <div>
            <label class="block text-[13px] font-semibold text-muted-foreground mb-1.5">密码</label>
            <div class="relative">
              <Input
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="••••••••"
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

          <div>
            <label class="block text-[13px] font-semibold text-muted-foreground mb-1.5">确认密码</label>
            <Input
              v-model="form.confirmPassword"
              :type="showPassword ? 'text' : 'password'"
              placeholder="••••••••"
              class="h-12"
              @focus="isTyping = true"
              @blur="isTyping = false"
            />
            <p v-if="errors.confirmPassword" class="text-sm text-destructive mt-1">{{ errors.confirmPassword }}</p>
          </div>

          <Button type="button" class="w-full h-12 text-base font-semibold" :disabled="loading" @click="handleSubmit">
            <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
            {{ loading ? 'Creating...' : 'Create account' }}
          </Button>
        </form>

        <div class="text-center text-[13px] text-muted-foreground mt-4.5">
          已有账号？
          <router-link class="text-primary font-semibold" to="/login">Log in</router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { Eye, EyeOff, Loader2 } from 'lucide-vue-next'
import AnimatedCharacters from '@/components/AnimatedCharacters.vue'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { toast } from 'vue-sonner'

const router = useRouter()

const loading = ref(false)
const showPassword = ref(false)
const isTyping = ref(false)

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
    errors.username = '请输入用户名'
    valid = false
  }

  if (!form.email) {
    errors.email = '请输入邮箱'
    valid = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = '请输入正确的邮箱地址'
    valid = false
  }

  if (!form.password) {
    errors.password = '请输入密码'
    valid = false
  } else if (form.password.length < 6) {
    errors.password = '密码至少6位'
    valid = false
  }

  if (!form.confirmPassword) {
    errors.confirmPassword = '请确认密码'
    valid = false
  } else if (form.confirmPassword !== form.password) {
    errors.confirmPassword = '两次输入的密码不一致'
    valid = false
  }

  return valid
}

async function handleSubmit() {
  if (!validate()) return
  try {
    loading.value = true
    toast.info('功能开发中...')
    router.push('/login')
  } finally {
    loading.value = false
  }
}

function togglePassword() {
  showPassword.value = !showPassword.value
}
</script>
