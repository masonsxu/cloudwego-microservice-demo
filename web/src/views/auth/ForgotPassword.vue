<template>
  <div class="min-h-screen grid grid-cols-1 lg:grid-cols-[1.1fr_1fr] bg-canvas">
    <!-- 左侧 brand 区 -->
    <section class="hidden lg:flex flex-col justify-between bg-sunken px-12 py-10 border-r border-subtle">
      <div class="flex items-center">
        <img src="/logo-light.svg" alt="CloudWeGo" class="h-8 w-auto" />
      </div>
      <div class="max-w-[420px]">
        <h2 class="text-[36px] font-semibold leading-[1.15] tracking-[-0.02em] text-[color:var(--color-ink-strong)]">
          {{ t('auth.forgotBrandTitle') }}
        </h2>
        <p class="mt-4 text-[16px] leading-relaxed text-[color:var(--color-ink-muted)]">
          {{ t('auth.forgotBrandSubtitle') }}
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
            {{ t('auth.forgotTitle') }}
          </h1>
          <p class="mt-1.5 text-[14px] text-[color:var(--color-ink-muted)]">
            {{ t('auth.forgotSubtitle') }}
          </p>
        </header>

        <form class="space-y-5" @submit.prevent="handleSubmit">
          <div class="space-y-1.5">
            <label class="block text-[13px] font-semibold text-ink">{{ t('user.email') }}</label>
            <Input v-model="form.email" placeholder="you@example.com" autocomplete="email" />
            <p v-if="errors.email" class="text-[12px] text-[color:var(--color-danger-ink)]">{{ errors.email }}</p>
          </div>

          <Button type="submit" class="w-full h-10" :disabled="loading">
            <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
            {{ loading ? t('auth.sending') : t('auth.sendResetLink') }}
          </Button>
        </form>

        <p class="mt-8 text-center text-[13px] text-[color:var(--color-ink-muted)]">
          <router-link
            to="/login"
            class="text-[color:var(--color-primary)] font-semibold transition-colors duration-[var(--duration-fast)] hover:text-[color:var(--color-primary-hover)]"
          >
            <ArrowLeft class="inline h-3.5 w-3.5 mr-1" />{{ t('auth.backToLogin') }}
          </router-link>
        </p>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Loader2, ArrowLeft } from 'lucide-vue-next'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { toast } from 'vue-sonner'

const { t } = useI18n()

const loading = ref(false)
const currentYear = computed(() => new Date().getFullYear())

const form = reactive({ email: '' })
const errors = reactive<Record<string, string>>({})

function validate(): boolean {
  Object.keys(errors).forEach((k) => delete errors[k])
  let valid = true

  if (!form.email) {
    errors.email = t('auth.emailRequired')
    valid = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = t('auth.emailInvalid')
    valid = false
  }

  return valid
}

async function handleSubmit() {
  if (!validate()) return
  try {
    loading.value = true
    toast.info(t('auth.forgotComingSoon'))
  } finally {
    loading.value = false
  }
}
</script>
