<template>
  <div class="oidc-callback">
    <Card class="max-w-3xl mx-auto mt-12">
      <CardHeader class="border-b border-border/60">
        <CardTitle>OIDC 授权回调</CardTitle>
      </CardHeader>
      <CardContent class="p-5">
        <!-- 加载中 -->
        <div v-if="loading" class="flex items-center justify-center py-12">
          <Loader2 class="w-6 h-6 animate-spin mr-2" />
          <span>正在处理授权回调...</span>
        </div>

        <!-- 错误 -->
        <div v-else-if="error" class="space-y-4">
          <Alert variant="destructive">
            <AlertCircle class="h-4 w-4" />
            <AlertTitle>授权失败</AlertTitle>
            <AlertDescription>{{ error }}</AlertDescription>
          </Alert>
          <Button variant="outline" @click="router.push('/system-settings/oidc/integration')">
            返回测试页
          </Button>
        </div>

        <!-- 成功 -->
        <div v-else class="space-y-5">
          <Alert>
            <CheckCircle2 class="h-4 w-4 text-green-500" />
            <AlertTitle>授权成功</AlertTitle>
            <AlertDescription>OIDC Authorization Code Flow 已完成</AlertDescription>
          </Alert>

          <!-- 授权码 -->
          <div class="space-y-2">
            <h4 class="text-sm font-semibold text-muted-foreground">Step 1: 授权码 (Code)</h4>
            <div class="result-block">
              <code>{{ result.code }}</code>
            </div>
          </div>

          <!-- State 验证 -->
          <div class="space-y-2">
            <h4 class="text-sm font-semibold text-muted-foreground">Step 2: State 验证</h4>
            <div class="result-block">
              <Badge :variant="result.stateValid ? 'outline' : 'destructive'" class="text-xs"
                :class="result.stateValid ? 'bg-green-500/10 border-green-500/30 text-green-500' : ''">
                {{ result.stateValid ? 'State 验证通过' : 'State 验证失败' }}
              </Badge>
              <code class="ml-2 text-xs">{{ result.state }}</code>
            </div>
          </div>

          <!-- Token 响应 -->
          <div class="space-y-2">
            <h4 class="text-sm font-semibold text-muted-foreground">Step 3: Token 响应</h4>
            <div class="result-block">
              <Badge v-if="result.tokenResponse" variant="outline"
                class="text-xs bg-green-500/10 border-green-500/30 text-green-500">
                Token 获取成功
              </Badge>
              <Badge v-else variant="destructive" class="text-xs">
                Token 获取失败
              </Badge>
            </div>
            <pre v-if="result.tokenResponse" class="code-block">{{ JSON.stringify(result.tokenResponse, null, 2) }}</pre>
            <div v-if="result.tokenError" class="result-block text-red-500">
              <code>{{ result.tokenError }}</code>
            </div>
          </div>

          <!-- Userinfo -->
          <div v-if="result.userinfo" class="space-y-2">
            <h4 class="text-sm font-semibold text-muted-foreground">Step 4: 用户信息</h4>
            <pre class="code-block">{{ JSON.stringify(result.userinfo, null, 2) }}</pre>
          </div>
          <div v-else-if="result.userinfoError" class="space-y-2">
            <h4 class="text-sm font-semibold text-muted-foreground">Step 4: 用户信息（失败）</h4>
            <div class="result-block text-red-500">
              <code>{{ result.userinfoError }}</code>
            </div>
          </div>

          <div class="flex gap-3 pt-2">
            <Button variant="outline" @click="router.push('/system-settings/oidc/integration')">
              再次测试
            </Button>
            <Button variant="outline" @click="router.push('/system-settings/oidc/config')">
              查看配置
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert'
import { Loader2, AlertCircle, CheckCircle2 } from 'lucide-vue-next'

const router = useRouter()

const loading = ref(true)
const error = ref('')

const result = reactive({
  code: '',
  state: '',
  stateValid: false,
  tokenResponse: null as Record<string, unknown> | null,
  tokenError: '',
  userinfo: null as Record<string, unknown> | null,
  userinfoError: '',
})

// 获取后端 API 的 base URL（开发环境通过 Vite 代理，生产环境直接请求）
const getApiBaseURL = () => {
  return import.meta.env.DEV ? '' : (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080')
}

onMounted(async () => {
  const params = new URLSearchParams(window.location.search)
  const code = params.get('code')
  const state = params.get('state')
  const errorParam = params.get('error')

  if (errorParam) {
    error.value = `授权错误: ${errorParam} - ${params.get('error_description') || ''}`
    loading.value = false
    return
  }

  if (!code) {
    error.value = '回调 URL 中未找到授权码 (code)'
    loading.value = false
    return
  }

  result.code = code
  result.state = state || ''

  // 验证 state
  const savedState = sessionStorage.getItem('oidc_state')
  result.stateValid = !!savedState && savedState === state

  if (!result.stateValid) {
    error.value = 'State 验证失败，可能存在 CSRF 攻击'
    loading.value = false
    return
  }

  // 获取 PKCE code_verifier
  const codeVerifier = sessionStorage.getItem('oidc_code_verifier')
  const redirectUri = sessionStorage.getItem('oidc_redirect_uri')
    || `${window.location.origin}/oidc/callback`

  // 清理 sessionStorage
  sessionStorage.removeItem('oidc_state')
  sessionStorage.removeItem('oidc_code_verifier')
  sessionStorage.removeItem('oidc_redirect_uri')

  try {
    // 用授权码换取 Token
    const baseURL = getApiBaseURL()
    const tokenParams = new URLSearchParams()
    tokenParams.set('grant_type', 'authorization_code')
    tokenParams.set('code', code)
    tokenParams.set('redirect_uri', redirectUri)
    tokenParams.set('client_id', 'demo-client')

    if (codeVerifier) {
      tokenParams.set('code_verifier', codeVerifier)
    }

    const tokenResp = await fetch(`${baseURL}/oauth/token`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: tokenParams.toString(),
    })

    if (!tokenResp.ok) {
      const errText = await tokenResp.text()
      result.tokenError = `HTTP ${tokenResp.status}: ${errText}`
    } else {
      result.tokenResponse = await tokenResp.json()

      // 用 Access Token 获取用户信息
      if (result.tokenResponse?.access_token) {
        try {
          const userinfoResp = await fetch(`${baseURL}/userinfo`, {
            headers: { Authorization: `Bearer ${result.tokenResponse.access_token}` },
          })
          if (userinfoResp.ok) {
            result.userinfo = await userinfoResp.json()
          } else {
            result.userinfoError = `HTTP ${userinfoResp.status}`
          }
        } catch (e) {
          result.userinfoError = String(e)
        }
      }
    }
  } catch (e) {
    result.tokenError = String(e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.oidc-callback {
  padding: 0 24px 24px;
}

.result-block {
  background: var(--bg-base);
  border-radius: 8px;
  padding: 12px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
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
  font-family: 'JetBrains Mono', monospace;
  color: var(--c-text-main);
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
