<template>
  <div class="login-shell">
    <section class="side-visual">
      <div class="brand">
        <div class="brand-mark">CW</div>
        <span>CloudWeGo</span>
      </div>

      <div class="character-stage">
        <AnimatedCharacters
          :is-typing="isTyping"
          :show-password="showPassword"
          :password-length="loginForm.password.length"
        />
      </div>

      <div class="legal">
        <span>Privacy Policy</span>
        <span>Terms of Service</span>
      </div>

      <div class="grid-overlay"></div>
      <div class="glow glow-top"></div>
      <div class="glow glow-bottom"></div>
    </section>

    <section class="login-panel">
      <div class="panel-inner">
        <div class="mobile-brand">
          <div class="brand-mark small">CW</div>
          <span>CloudWeGo</span>
        </div>

        <div class="panel-header">
          <h1>Welcome back!</h1>
          <p>请使用你的账号登录</p>
        </div>

        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          class="login-form"
          @keyup.enter="handleLogin"
        >
          <el-form-item prop="username">
            <label class="field-label">{{ t('auth.username') }}</label>
            <el-input
              v-model="loginForm.username"
              :placeholder="t('auth.username')"
              size="large"
              @focus="isTyping = true"
              @blur="isTyping = false"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item prop="password">
            <label class="field-label">{{ t('auth.password') }}</label>
            <el-input
              v-model="loginForm.password"
              :type="showPassword ? 'text' : 'password'"
              :placeholder="t('auth.password')"
              size="large"
              @focus="isTyping = true"
              @blur="isTyping = false"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
              <template #suffix>
                <button type="button" class="icon-btn" @click="togglePassword">
                  <el-icon><View v-if="!showPassword" /><Hide v-else /></el-icon>
                </button>
              </template>
            </el-input>
          </el-form-item>

          <div class="form-row">
            <el-checkbox v-model="rememberMe">Remember for 30 days</el-checkbox>
            <router-link class="link-text" to="/forgot-password">Forgot password?</router-link>
          </div>

          <el-button
            type="primary"
            size="large"
            class="login-button"
            :loading="loading"
            @click="handleLogin"
          >
            {{ loading ? 'Signing in...' : t('auth.login') }}
          </el-button>
        </el-form>

        <div class="divider">
          <span>or</span>
        </div>

        <el-button class="outline-button" size="large">
          Log in with Google
        </el-button>

        <div class="signup-hint">
          还没有账号？
          <router-link class="link-text" to="/signup">Sign Up</router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock, View, Hide } from '@element-plus/icons-vue'
import AnimatedCharacters from '@/components/AnimatedCharacters.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const { t } = useI18n()

const loginFormRef = ref<FormInstance>()
const loading = ref(false)
const showPassword = ref(false)
const isTyping = ref(false)
const rememberMe = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ]
}

async function handleLogin() {
  if (!loginFormRef.value) return

  try {
    await loginFormRef.value.validate()
    loading.value = true

    await authStore.login({
      username: loginForm.username,
      password: loginForm.password
    })

    ElMessage.success(t('auth.loginSuccess'))

    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect)
  } catch (error: any) {
    console.error('Login failed:', error)
    ElMessage.error(error.message || t('auth.loginFailed'))
  } finally {
    loading.value = false
  }
}

function togglePassword() {
  showPassword.value = !showPassword.value
}
</script>

<style scoped lang="scss" src="./login-shared.scss"></style>
