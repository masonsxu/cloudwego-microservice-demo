<template>
  <div class="oidc-config">
    <div class="page-header">
      <h2>{{ t('oidc.config.title') }}</h2>
      <div class="header-actions">
        <el-button type="primary" @click="refreshData">
          <el-icon><Refresh /></el-icon>
          {{ t('common.refresh') }}
        </el-button>
      </div>
    </div>

    <div class="content-row">
      <DetailPageSkeleton
        v-if="initialLoading"
        :side-span="16"
        :main-span="8"
        :side-cards="2"
        :main-cards="1"
        :show-avatar="false"
        :side-item-counts="[8, 4]"
        :main-item-counts="[1]"
      />
      <el-row v-else v-loading="loading" :gutter="20">
        <el-col :span="16">
          <el-card shadow="never" class="config-card">
            <template #header>
              <div class="card-header">
                <span>{{ t('oidc.config.providerConfig') }}</span>
                <el-tag :type="discoveryConfig ? 'success' : 'danger'" size="small">
                  {{ discoveryConfig ? t('oidc.config.discovered') : t('oidc.config.unavailable') }}
                </el-tag>
              </div>
            </template>
            <el-descriptions :column="2" border>
              <el-descriptions-item :label="t('oidc.config.issuer')" :span="2">
                <code class="code-text">{{ discoveryConfig?.issuer || '-' }}</code>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.authorizationEndpoint')" :span="2">
                <code class="code-text">{{ discoveryConfig?.authorization_endpoint || '-' }}</code>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.tokenEndpoint')" :span="2">
                <code class="code-text">{{ discoveryConfig?.token_endpoint || '-' }}</code>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.userinfoEndpoint')" :span="2">
                <code class="code-text">{{ discoveryConfig?.userinfo_endpoint || '-' }}</code>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.jwksUri')" :span="2">
                <code class="code-text">{{ discoveryConfig?.jwks_uri || '-' }}</code>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.revocationEndpoint')" :span="2">
                <code class="code-text">{{ discoveryConfig?.revocation_endpoint || '-' }}</code>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.introspectionEndpoint')" :span="2">
                <code class="code-text">{{ discoveryConfig?.introspection_endpoint || '-' }}</code>
              </el-descriptions-item>
            </el-descriptions>
          </el-card>

          <el-card shadow="never" class="config-card" style="margin-top: 16px;">
            <template #header>{{ t('oidc.config.capabilities') }}</template>
            <el-descriptions :column="1" border>
              <el-descriptions-item :label="t('oidc.config.responseTypes')">
                <el-tag
                  v-for="type in discoveryConfig?.response_types_supported"
                  :key="type"
                  size="small"
                  style="margin: 2px 4px;"
                >
                  {{ type }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.scopes')">
                <el-tag
                  v-for="scope in discoveryConfig?.scopes_supported"
                  :key="scope"
                  size="small"
                  style="margin: 2px 4px;"
                >
                  {{ scope }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.signingAlgorithms')">
                <el-tag
                  v-for="alg in discoveryConfig?.id_token_signing_alg_values_supported"
                  :key="alg"
                  size="small"
                  style="margin: 2px 4px;"
                >
                  {{ alg }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.authMethods')">
                <el-tag
                  v-for="method in discoveryConfig?.token_endpoint_auth_methods_supported"
                  :key="method"
                  size="small"
                  style="margin: 2px 4px;"
                >
                  {{ method }}
                </el-tag>
              </el-descriptions-item>
            </el-descriptions>
          </el-card>
        </el-col>

        <el-col :span="8">
          <el-card shadow="never">
            <template #header>{{ t('oidc.config.runtimeConfig') }}</template>
            <el-descriptions :column="1" border>
              <el-descriptions-item :label="t('oidc.config.enabled')">
                <el-tag :type="envConfig?.enabled ? 'success' : 'danger'" size="small" effect="dark">
                  {{ envConfig?.enabled ? t('common.yes') : t('common.no') }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.enforcePkce')">
                <el-tag :type="envConfig?.enforcePkce ? 'success' : 'warning'" size="small" effect="dark">
                  {{ envConfig?.enforcePkce ? t('common.yes') : t('common.no') }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.accessTokenLifespan')">
                {{ envConfig?.accessTokenLifespan || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.refreshTokenLifespan')">
                {{ envConfig?.refreshTokenLifespan || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.authCodeLifespan')">
                {{ envConfig?.authCodeLifespan || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('oidc.config.idTokenLifespan')">
                {{ envConfig?.idTokenLifespan || '-' }}
              </el-descriptions-item>
            </el-descriptions>
          </el-card>

          <el-card shadow="never" style="margin-top: 16px;">
            <template #header>{{ t('oidc.config.quickLinks') }}</template>
            <div class="quick-links">
              <el-button class="link-btn" @click="showJwksDialog = true">
                <el-icon><Link /></el-icon>{{ t('oidc.config.viewJwks') }}
              </el-button>
              <el-button class="link-btn" @click="showDiscoveryDialog = true">
                <el-icon><Link /></el-icon>{{ t('oidc.config.viewDiscovery') }}
              </el-button>
              <el-button class="link-btn" @click="router.push('/system-settings/oidc/integration')">
                <el-icon><Document /></el-icon>{{ t('oidc.config.integrationGuide') }}
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <el-dialog v-model="showJwksDialog" :title="t('oidc.config.viewJwks')" width="70%" @open="fetchJwks">
      <pre class="json-viewer" v-loading="jwksLoading">{{ jwksContent ? JSON.stringify(jwksContent, null, 2) : t('common.noData') }}</pre>
    </el-dialog>

    <el-dialog v-model="showDiscoveryDialog" :title="t('oidc.config.viewDiscovery')" width="70%">
      <pre class="json-viewer">{{ discoveryJson || t('common.noData') }}</pre>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { oidcProviderApi } from '@/api/oidc'
import DetailPageSkeleton from '@/components/skeleton/DetailPageSkeleton.vue'
import type { OIDCDiscoveryConfig, OIDCJWKSResponse } from '@/types/oidc'
import { Document, Link, Refresh } from '@element-plus/icons-vue'
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const router = useRouter()
const { t } = useI18n()

const initialLoading = ref(true)
const loading = ref(false)
const discoveryConfig = ref<OIDCDiscoveryConfig | null>(null)
const showJwksDialog = ref(false)
const showDiscoveryDialog = ref(false)
const jwksContent = ref<OIDCJWKSResponse | null>(null)
const jwksLoading = ref(false)

const envConfig = ref({
  enabled: true,
  enforcePkce: true,
  accessTokenLifespan: '30m',
  refreshTokenLifespan: '7d',
  authCodeLifespan: '10m',
  idTokenLifespan: '30m',
})

const discoveryJson = computed(() => {
  if (!discoveryConfig.value) return ''
  return JSON.stringify(discoveryConfig.value, null, 2)
})

const fetchDiscovery = async () => {
  loading.value = true
  try {
    const response = await oidcProviderApi.getDiscovery()
    discoveryConfig.value = response
  } catch {
    discoveryConfig.value = null
  } finally {
    loading.value = false
  }
}

const refreshData = async () => {
  await fetchDiscovery()
}

const fetchJwks = async () => {
  jwksLoading.value = true
  try {
    const response = await oidcProviderApi.getJWKS()
    jwksContent.value = response
  } catch {
    jwksContent.value = null
  } finally {
    jwksLoading.value = false
  }
}

onMounted(async () => {
  try {
    await fetchDiscovery()
  } finally {
    initialLoading.value = false
  }
})
</script>

<style scoped lang="scss">
.oidc-config {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;

  h2 {
    font-size: 20px;
    font-weight: 600;
    color: var(--c-text-main);
  }

  .header-actions {
    display: flex;
    gap: 8px;
  }
}

.content-row {
  flex: 1;
}

.config-card {
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.code-text {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
  background: var(--bg-base);
  color: var(--c-accent);
}

.quick-links {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.link-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  text-align: left;
}

.json-viewer {
  background: var(--bg-base);
  border-radius: 8px;
  padding: 16px;
  overflow: auto;
  max-height: 60vh;
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: var(--c-text-main);
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
