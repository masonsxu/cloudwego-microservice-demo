<template>
  <div class="oidc-integration">
    <div class="page-header">
      <h2>{{ t('oidc.integration.title') }}</h2>
    </div>

    <el-alert type="info" :closable="false" show-icon class="guide-alert">
      <template #title>{{ t('oidc.integration.guideTitle') }}</template>
      <div class="guide-steps">
        <div class="guide-step">
          <span class="step-num">1</span>
          <span>{{ t('oidc.integration.step1') }}</span>
        </div>
        <div class="guide-step">
          <span class="step-num">2</span>
          <span>{{ t('oidc.integration.step2') }}</span>
        </div>
        <div class="guide-step">
          <span class="step-num">3</span>
          <span>{{ t('oidc.integration.step3') }}</span>
        </div>
        <div class="guide-step">
          <span class="step-num">4</span>
          <span>{{ t('oidc.integration.step4') }}</span>
        </div>
      </div>
    </el-alert>

    <el-row :gutter="20">
      <el-col :span="16">
        <el-card shadow="never">
          <template #header>{{ t('oidc.integration.testFlow') }}</template>

          <el-form :model="testForm" label-width="160px" class="test-form">
            <el-form-item :label="t('oidc.integration.clientId')">
              <el-input v-model="testForm.clientId" placeholder="your-client-id" />
            </el-form-item>
            <el-form-item :label="t('oidc.integration.redirectUri')">
              <el-input v-model="testForm.redirectUri" placeholder="http://localhost:3000/callback" />
            </el-form-item>
            <el-form-item :label="t('oidc.integration.scope')">
              <el-checkbox-group v-model="testForm.scopes">
                <el-checkbox value="openid">openid</el-checkbox>
                <el-checkbox value="profile">profile</el-checkbox>
                <el-checkbox value="email">email</el-checkbox>
                <el-checkbox value="offline_access">offline_access</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-form-item :label="t('oidc.integration.enablePkce')">
              <el-switch v-model="testForm.enablePkce" />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="startAuthFlow" :loading="flowLoading">
                {{ t('oidc.integration.startFlow') }}
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card shadow="never" style="margin-top: 16px;">
          <template #header>{{ t('oidc.integration.flowLog') }}</template>
          <div class="flow-log" ref="logContainerRef">
            <div v-for="(log, idx) in flowLogs" :key="idx" :class="['log-entry', `log-${log.level}`]">
              <span class="log-time">{{ log.time }}</span>
              <span class="log-message">{{ log.message }}</span>
            </div>
            <div v-if="flowLogs.length === 0" class="log-empty">{{ t('oidc.integration.noLogs') }}</div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card shadow="never">
          <template #header>{{ t('oidc.integration.endpoints') }}</template>
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item :label="t('oidc.config.discovery')">
              <code class="code-text">/.well-known/openid-configuration</code>
            </el-descriptions-item>
            <el-descriptions-item :label="t('oidc.config.jwksUri')">
              <code class="code-text">/keys</code>
            </el-descriptions-item>
            <el-descriptions-item :label="t('oidc.config.authorizationEndpoint')">
              <code class="code-text">/authorize</code>
            </el-descriptions-item>
            <el-descriptions-item :label="t('oidc.config.tokenEndpoint')">
              <code class="code-text">/oauth/token</code>
            </el-descriptions-item>
            <el-descriptions-item :label="t('oidc.config.userinfoEndpoint')">
              <code class="code-text">/userinfo</code>
            </el-descriptions-item>
            <el-descriptions-item :label="t('oidc.config.revocationEndpoint')">
              <code class="code-text">/revoke</code>
            </el-descriptions-item>
            <el-descriptions-item :label="t('oidc.config.introspectionEndpoint')">
              <code class="code-text">/oauth/introspect</code>
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card shadow="never" style="margin-top: 16px;">
          <template #header>{{ t('oidc.integration.codeExample') }}</template>
          <pre class="code-block"><code>{{ codeExample }}</code></pre>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { nextTick, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

interface FlowLog {
  time: string
  message: string
  level: 'info' | 'success' | 'error' | 'warn'
}

const testForm = reactive({
  clientId: 'demo-client',
  redirectUri: 'http://localhost:3000/callback',
  scopes: ['openid', 'profile', 'email'],
  enablePkce: true,
})

const flowLoading = ref(false)
const flowLogs = ref<FlowLog[]>([])
const logContainerRef = ref<HTMLElement>()

const codeExample = `// 1. 生成 PKCE code_verifier 和 code_challenge
const codeVerifier = generateRandomString(64);
const codeChallenge = await generateCodeChallenge(codeVerifier);

// 2. 构建授权 URL
const authUrl = new URL('/authorize', location.origin);
authUrl.searchParams.set('response_type', 'code');
authUrl.searchParams.set('client_id', 'your-client-id');
authUrl.searchParams.set('redirect_uri', 'http://localhost:3000/callback');
authUrl.searchParams.set('scope', 'openid profile email');
authUrl.searchParams.set('code_challenge', codeChallenge);
authUrl.searchParams.set('code_challenge_method', 'S256');
authUrl.searchParams.set('state', generateRandomString(32));

// 3. 重定向到授权页面
window.location.href = authUrl.toString();

// 4. 回调页面处理授权码
const params = new URLSearchParams(window.location.search);
const code = params.get('code');
const state = params.get('state');

// 5. 用授权码换取 Token
const tokenResponse = await fetch('/oauth/token', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    grant_type: 'authorization_code',
    code,
    redirect_uri: 'http://localhost:3000/callback',
    code_verifier: codeVerifier,
  }),
});

const { access_token, id_token, refresh_token } = await tokenResponse.json();

// 6. 用 Access Token 获取用户信息
const userinfoResponse = await fetch('/userinfo', {
  headers: { 'Authorization': \`Bearer \${access_token}\` },
});
const userinfo = await userinfoResponse.json();`

const addLog = (message: string, level: FlowLog['level'] = 'info') => {
  const now = new Date()
  flowLogs.value.push({
    time: now.toLocaleTimeString(),
    message,
    level,
  })
  nextTick(() => {
    if (logContainerRef.value) {
      logContainerRef.value.scrollTop = logContainerRef.value.scrollHeight
    }
  })
}

const generateRandomString = (length: number) => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~'
  return Array.from({ length }, () => chars[Math.floor(Math.random() * chars.length)]).join('')
}

const startAuthFlow = async () => {
  flowLoading.value = true
  flowLogs.value = []

  try {
    addLog(t('oidc.integration.log.buildingAuthUrl'), 'info')

    const baseUrl = window.location.origin
    const authUrl = new URL('/authorize', baseUrl)
    authUrl.searchParams.set('response_type', 'code')
    authUrl.searchParams.set('client_id', testForm.clientId)
    authUrl.searchParams.set('redirect_uri', testForm.redirectUri)
    authUrl.searchParams.set('scope', testForm.scopes.join(' '))
    authUrl.searchParams.set('state', generateRandomString(32))

    if (testForm.enablePkce) {
      const codeVerifier = generateRandomString(64)
      authUrl.searchParams.set('code_challenge_method', 'S256')
      addLog(t('oidc.integration.log.pkceEnabled'), 'info')
    }

    addLog(t('oidc.integration.log.authUrlReady'), 'success')
    addLog(`${t('oidc.integration.log.url')}: ${authUrl.toString()}`, 'info')

    ElMessage.success(t('oidc.integration.log.authUrlGenerated'))
  } catch (error) {
    addLog(`${t('oidc.integration.log.error')}: ${error}`, 'error')
    ElMessage.error(t('oidc.integration.log.flowFailed'))
  } finally {
    flowLoading.value = false
  }
}
</script>

<style scoped lang="scss">
.oidc-integration {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--c-text-main);
}

.guide-alert {
  :deep(.el-alert__title) {
    font-weight: 600;
  }
}

.guide-steps {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 13px;
  margin-top: 8px;
}

.guide-step {
  display: flex;
  align-items: center;
  gap: 8px;
}

.step-num {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--c-primary);
  color: white;
  font-size: 11px;
  font-weight: 600;
  flex-shrink: 0;
}

.test-form {
  :deep(.el-form-item__label) {
    font-weight: 500;
  }
}

.flow-log {
  height: 300px;
  overflow-y: auto;
  background: var(--bg-base);
  border-radius: 8px;
  padding: 12px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
}

.log-entry {
  display: flex;
  gap: 12px;
  padding: 4px 0;
  border-bottom: 1px solid var(--c-border);

  &:last-child {
    border-bottom: none;
  }
}

.log-time {
  color: var(--c-text-sub);
  flex-shrink: 0;
}

.log-info {
  .log-message {
    color: var(--c-text-main);
  }
}

.log-success {
  .log-message {
    color: var(--el-color-success);
  }
}

.log-error {
  .log-message {
    color: var(--el-color-danger);
  }
}

.log-warn {
  .log-message {
    color: var(--el-color-warning);
  }
}

.log-empty {
  text-align: center;
  color: var(--c-text-sub);
  padding: 40px 0;
}

.code-text {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 4px;
  background: var(--bg-base);
  color: var(--c-accent);
}

.code-block {
  background: var(--bg-base);
  border-radius: 8px;
  padding: 16px;
  overflow-x: auto;
  font-size: 12px;
  line-height: 1.6;
  max-height: 400px;
  overflow-y: auto;

  code {
    color: var(--c-text-main);
  }
}
</style>
