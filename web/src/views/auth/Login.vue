<template>
  <div class="login-container">
    <div class="login-card">
      <div class="logo-section">
        <h1 class="title">
          <span class="zodiac">♑</span>
          <span class="text">CloudWeGo</span>
        </h1>
        <p class="subtitle">{{ t('auth.login') }}</p>
      </div>
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            :placeholder="t('auth.username')"
            prefix-icon="User"
            size="large"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            :placeholder="t('auth.password')"
            prefix-icon="Lock"
            size="large"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-button"
            @click="handleLogin"
          >
            {{ t('auth.login') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="stars">
      <span v-for="i in 50" :key="i" class="star" :style="getStarStyle(i)"></span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const { t } = useI18n()

const loginFormRef = ref<FormInstance>()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
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

function getStarStyle(_i: number) {
  const size = Math.random() * 3 + 1
  const x = Math.random() * 100
  const y = Math.random() * 100
  const delay = Math.random() * 5
  const duration = Math.random() * 3 + 2

  return {
    width: `${size}px`,
    height: `${size}px`,
    left: `${x}%`,
    top: `${y}%`,
    animationDelay: `${delay}s`,
    animationDuration: `${duration}s`
  }
}
</script>

<style scoped lang="scss">
.login-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  overflow: hidden;

  .login-card {
    background: linear-gradient(145deg, rgba(30, 32, 36, 0.9), rgba(20, 20, 22, 0.95));
    border: 1px solid rgba(255, 255, 255, 0.05);
    border-radius: 20px;
    box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5);
    padding: 50px 40px;
    width: 400px;
    z-index: 10;
    backdrop-filter: blur(10px);

    .logo-section {
      text-align: center;
      margin-bottom: 40px;

      .title {
        font-family: 'Cinzel', serif;
        font-size: 48px;
        margin-bottom: 10px;
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 15px;

        .zodiac {
          font-size: 56px;
          background: linear-gradient(to right, #D4AF37, #F2F0E4, #D4AF37);
          background-size: 200% auto;
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          animation: shine 5s linear infinite;
        }

        .text {
          background: linear-gradient(to right, #D4AF37, #F2F0E4, #D4AF37);
          background-size: 200% auto;
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          animation: shine 5s linear infinite;
        }
      }

      .subtitle {
        color: #8B9bb4;
        font-size: 16px;
        letter-spacing: 2px;
      }
    }

    .login-form {
      .el-form-item {
        margin-bottom: 25px;
      }

      .login-button {
        width: 100%;
        height: 50px;
        font-size: 16px;
        font-weight: 600;
        letter-spacing: 2px;
      }
    }
  }

  .stars {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;

    .star {
      position: absolute;
      background: #D4AF37;
      border-radius: 50%;
      animation: twinkle ease-in-out infinite;
    }
  }
}

@keyframes shine {
  to {
    background-position: 200% center;
  }
}

@keyframes twinkle {
  0%, 100% {
    opacity: 1;
    filter: brightness(1);
  }
  50% {
    opacity: 0.6;
    filter: brightness(0.7);
  }
}
</style>
