<template>
  <div class="login-shell">
    <section class="side-visual">
      <div class="brand">
        <div class="brand-mark">CW</div>
        <span>CloudWeGo</span>
      </div>

      <div class="character-stage">
        <AnimatedCharacters :is-typing="isTyping" :show-password="false" :password-length="0" />
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
          <h1>Reset password</h1>
          <p>我们会发送重置链接到你的邮箱</p>
        </div>

        <el-form ref="formRef" :model="form" :rules="rules" class="login-form" @keyup.enter="handleSubmit">
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

          <el-button type="primary" size="large" class="login-button" :loading="loading" @click="handleSubmit">
            {{ loading ? 'Sending...' : 'Send reset link' }}
          </el-button>
        </el-form>

        <div class="signup-hint">
          返回登录
          <router-link class="link-text" to="/login">Log in</router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Message } from '@element-plus/icons-vue'
import AnimatedCharacters from '@/components/AnimatedCharacters.vue'

const formRef = ref<FormInstance>()
const loading = ref(false)
const isTyping = ref(false)

const form = reactive({ email: '' })

const rules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

async function handleSubmit() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    loading.value = true
    ElMessage.info('功能开发中...')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
@import './login-shared.scss';
</style>
