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
          :password-length="form.password.length"
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
          <h1>Create account</h1>
          <p>快速创建你的管理账号</p>
        </div>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          class="login-form"
          @keyup.enter="handleSubmit"
        >
          <el-form-item prop="username">
            <label class="field-label">用户名</label>
            <el-input
              v-model="form.username"
              placeholder="请输入用户名"
              size="large"
              @focus="isTyping = true"
              @blur="isTyping = false"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item prop="email">
            <label class="field-label">邮箱</label>
            <el-input
              v-model="form.email"
              placeholder="you@example.com"
              size="large"
              @focus="isTyping = true"
              @blur="isTyping = false"
            >
              <template #prefix>
                <el-icon><Message /></el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item prop="password">
            <label class="field-label">密码</label>
            <el-input
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="••••••••"
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

          <el-form-item prop="confirmPassword">
            <label class="field-label">确认密码</label>
            <el-input
              v-model="form.confirmPassword"
              :type="showPassword ? 'text' : 'password'"
              placeholder="••••••••"
              size="large"
              @focus="isTyping = true"
              @blur="isTyping = false"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-button
            type="primary"
            size="large"
            class="login-button"
            :loading="loading"
            @click="handleSubmit"
          >
            {{ loading ? 'Creating...' : 'Create account' }}
          </el-button>
        </el-form>

        <div class="signup-hint">
          已有账号？
          <router-link class="link-text" to="/login">Log in</router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock, View, Hide, Message } from '@element-plus/icons-vue'
import AnimatedCharacters from '@/components/AnimatedCharacters.vue'

const router = useRouter()

const formRef = ref<FormInstance>()
const loading = ref(false)
const showPassword = ref(false)
const isTyping = ref(false)

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (!value || value !== form.password) callback(new Error('两次输入的密码不一致'))
        callback()
      },
      trigger: 'blur'
    }
  ]
}

async function handleSubmit() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    loading.value = true
    ElMessage.info('功能开发中...')
    router.push('/login')
  } finally {
    loading.value = false
  }
}

function togglePassword() {
  showPassword.value = !showPassword.value
}
</script>

<style scoped lang="scss" src="./login-shared.scss"></style>
