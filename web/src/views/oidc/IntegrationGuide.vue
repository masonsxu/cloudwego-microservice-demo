<template>
  <div class="oidc-integration">
    <div class="page-header">
      <h2>{{ t("oidc.integration.title") }}</h2>
    </div>

    <Alert>
      <AlertTitle>{{ t("oidc.integration.guideTitle") }}</AlertTitle>
      <AlertDescription>
        <div class="guide-steps">
          <div class="guide-step">
            <span class="step-num">1</span>
            <span>{{ t("oidc.integration.step1") }}</span>
          </div>
          <div class="guide-step">
            <span class="step-num">2</span>
            <span>{{ t("oidc.integration.step2") }}</span>
          </div>
          <div class="guide-step">
            <span class="step-num">3</span>
            <span>{{ t("oidc.integration.step3") }}</span>
          </div>
          <div class="guide-step">
            <span class="step-num">4</span>
            <span>{{ t("oidc.integration.step4") }}</span>
          </div>
        </div>
      </AlertDescription>
    </Alert>

    <div class="grid gap-4 lg:grid-cols-3">
      <div class="lg:col-span-2 space-y-4">
        <Card>
          <CardHeader>
            <CardTitle>{{ t("oidc.integration.testFlow") }}</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="space-y-4">
              <div class="space-y-2">
                <Label>{{ t("oidc.integration.clientId") }}</Label>
                <Input
                  v-model="testForm.clientId"
                  placeholder="your-client-id"
                />
              </div>
              <div class="space-y-2">
                <Label>{{ t("oidc.integration.redirectUri") }}</Label>
                <Input
                  v-model="testForm.redirectUri"
                  placeholder="http://localhost:3000/callback"
                />
              </div>
              <div class="space-y-2">
                <Label>{{ t("oidc.integration.scope") }}</Label>
                <div class="flex flex-wrap gap-4">
                  <div
                    class="flex items-center gap-2"
                    v-for="scope in [
                      'openid',
                      'profile',
                      'email',
                      'offline_access',
                    ]"
                    :key="scope"
                  >
                    <Checkbox
                      :id="scope"
                      :value="scope"
                      :checked="testForm.scopes.includes(scope)"
                      @update:checked="
                        (checked: boolean) => toggleScope(scope, checked)
                      "
                    />
                    <label :for="scope" class="text-sm">{{ scope }}</label>
                  </div>
                </div>
              </div>
              <div class="flex items-center gap-2">
                <Switch
                  :checked="testForm.enablePkce"
                  @update:checked="
                    (checked: boolean) => (testForm.enablePkce = checked)
                  "
                />
                <Label>{{ t("oidc.integration.enablePkce") }}</Label>
              </div>
              <Button
                variant="default"
                @click="startAuthFlow"
                :loading="flowLoading"
              >
                {{ t("oidc.integration.startFlow") }}
              </Button>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>{{ t("oidc.integration.flowLog") }}</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="flow-log" ref="logContainerRef">
              <div
                v-for="(log, idx) in flowLogs"
                :key="idx"
                :class="['log-entry', `log-${log.level}`]"
              >
                <span class="log-time">{{ log.time }}</span>
                <span class="log-message">{{ log.message }}</span>
              </div>
              <div v-if="flowLogs.length === 0" class="log-empty">
                {{ t("oidc.integration.noLogs") }}
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <div class="space-y-4">
        <Card>
          <CardHeader>
            <CardTitle>{{ t("oidc.integration.endpoints") }}</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="space-y-3">
              <div
                class="flex justify-between items-center py-2 border-b border-[hsl(var(--border))]"
              >
                <span class="text-sm text-[var(--c-text-sub)]">{{
                  t("oidc.config.discovery")
                }}</span>
                <code class="code-text">/.well-known/openid-configuration</code>
              </div>
              <div
                class="flex justify-between items-center py-2 border-b border-[hsl(var(--border))]"
              >
                <span class="text-sm text-[var(--c-text-sub)]">{{
                  t("oidc.config.jwksUri")
                }}</span>
                <code class="code-text">/keys</code>
              </div>
              <div
                class="flex justify-between items-center py-2 border-b border-[hsl(var(--border))]"
              >
                <span class="text-sm text-[var(--c-text-sub)]">{{
                  t("oidc.config.authorizationEndpoint")
                }}</span>
                <code class="code-text">/authorize</code>
              </div>
              <div
                class="flex justify-between items-center py-2 border-b border-[hsl(var(--border))]"
              >
                <span class="text-sm text-[var(--c-text-sub)]">{{
                  t("oidc.config.tokenEndpoint")
                }}</span>
                <code class="code-text">/oauth/token</code>
              </div>
              <div
                class="flex justify-between items-center py-2 border-b border-[hsl(var(--border))]"
              >
                <span class="text-sm text-[var(--c-text-sub)]">{{
                  t("oidc.config.userinfoEndpoint")
                }}</span>
                <code class="code-text">/userinfo</code>
              </div>
              <div
                class="flex justify-between items-center py-2 border-b border-[hsl(var(--border))]"
              >
                <span class="text-sm text-[var(--c-text-sub)]">{{
                  t("oidc.config.revocationEndpoint")
                }}</span>
                <code class="code-text">/revoke</code>
              </div>
              <div class="flex justify-between items-center py-2">
                <span class="text-sm text-[var(--c-text-sub)]">{{
                  t("oidc.config.introspectionEndpoint")
                }}</span>
                <code class="code-text">/oauth/introspect</code>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>{{ t("oidc.integration.codeExample") }}</CardTitle>
          </CardHeader>
          <CardContent>
            <pre class="code-block"><code>{{ codeExample }}</code></pre>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { toast } from "vue-sonner";
import { nextTick, reactive, ref } from "vue";
import { useI18n } from "vue-i18n";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Switch } from "@/components/ui/switch";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";

const { t } = useI18n();

interface FlowLog {
  time: string;
  message: string;
  level: "info" | "success" | "error" | "warn";
}

const testForm = reactive({
  clientId: "demo-client",
  redirectUri: `${window.location.origin}/oidc/callback`,
  scopes: ["openid", "profile", "email"],
  enablePkce: true,
});

const flowLoading = ref(false);
const flowLogs = ref<FlowLog[]>([]);
const logContainerRef = ref<HTMLElement>();

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
const userinfo = await userinfoResponse.json();`;

const addLog = (message: string, level: FlowLog["level"] = "info") => {
  const now = new Date();
  flowLogs.value.push({
    time: now.toLocaleTimeString(),
    message,
    level,
  });
  nextTick(() => {
    if (logContainerRef.value) {
      logContainerRef.value.scrollTop = logContainerRef.value.scrollHeight;
    }
  });
};

const generateRandomString = (length: number) => {
  const chars =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~";
  return Array.from(
    { length },
    () => chars[Math.floor(Math.random() * chars.length)],
  ).join("");
};

const generateCodeVerifier = () => {
  const array = new Uint8Array(32);
  crypto.getRandomValues(array);
  return btoa(String.fromCharCode(...array))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");
};

const generateCodeChallenge = async (verifier: string) => {
  const encoder = new TextEncoder();
  const data = encoder.encode(verifier);
  const digest = await crypto.subtle.digest("SHA-256", data);
  return btoa(String.fromCharCode(...new Uint8Array(digest)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");
};

const toggleScope = (scope: string, checked: boolean) => {
  if (checked) {
    if (!testForm.scopes.includes(scope)) {
      testForm.scopes.push(scope);
    }
  } else {
    testForm.scopes = testForm.scopes.filter((s) => s !== scope);
  }
};

const startAuthFlow = async () => {
  flowLoading.value = true;
  flowLogs.value = [];

  try {
    addLog("正在构建授权 URL...", "info");

    const state = generateRandomString(32);
    // 授权请求必须直接发往后端（浏览器导航不走 Vite 代理）
    const apiBaseURL = import.meta.env.DEV
      ? "http://localhost:8088"
      : window.location.origin;
    const authUrl = new URL("/authorize", apiBaseURL);
    authUrl.searchParams.set("response_type", "code");
    authUrl.searchParams.set("client_id", testForm.clientId);
    authUrl.searchParams.set("redirect_uri", testForm.redirectUri);
    authUrl.searchParams.set("scope", testForm.scopes.join(" "));
    authUrl.searchParams.set("state", state);

    if (testForm.enablePkce) {
      const codeVerifier = generateCodeVerifier();
      const codeChallenge = await generateCodeChallenge(codeVerifier);
      authUrl.searchParams.set("code_challenge", codeChallenge);
      authUrl.searchParams.set("code_challenge_method", "S256");
      // 存储 PKCE 参数到 sessionStorage，回调页面需要用
      sessionStorage.setItem("oidc_code_verifier", codeVerifier);
      addLog("PKCE 已启用，code_verifier 已存储到 sessionStorage", "info");
    }

    // 存储 state 用于回调验证
    sessionStorage.setItem("oidc_state", state);
    sessionStorage.setItem("oidc_redirect_uri", testForm.redirectUri);

    addLog("授权 URL 构建完成", "success");
    addLog(`URL: ${authUrl.toString()}`, "info");
    addLog("即将跳转到授权页面...", "info");

    // 延迟跳转，让用户看到日志
    setTimeout(() => {
      window.location.href = authUrl.toString();
    }, 500);
  } catch (error) {
    addLog(`错误: ${error}`, "error");
    toast.error("授权流程启动失败");
    flowLoading.value = false;
  }
};
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

.flow-log {
  height: 300px;
  overflow-y: auto;
  background: var(--bg-base);
  border-radius: 8px;
  padding: 12px;
  font-family: "JetBrains Mono", monospace;
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
    color: #67c23a;
  }
}

.log-error {
  .log-message {
    color: #f56c6c;
  }
}

.log-warn {
  .log-message {
    color: #e6a23c;
  }
}

.log-empty {
  text-align: center;
  color: var(--c-text-sub);
  padding: 40px 0;
}

.code-text {
  font-family: "JetBrains Mono", monospace;
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
