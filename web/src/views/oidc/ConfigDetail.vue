<template>
  <div class="oidc-config">
    <div class="page-header">
      <h2>{{ t('oidc.config.title') }}</h2>
      <div class="header-actions">
        <Button @click="refreshData">
          <RefreshCw class="w-4 h-4 mr-1" />
          {{ t('common.refresh') }}
        </Button>
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
      <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-5">
        <div class="lg:col-span-2 space-y-5">
          <Card class="bg-card border-border/60 rounded-[18px] shadow-card">
            <CardHeader class="border-b border-border/60 px-5 py-3.5">
              <div class="flex items-center justify-between">
                <span class="text-primary font-semibold">{{ t('oidc.config.providerConfig') }}</span>
                <Badge :variant="discoveryConfig ? 'outline' : 'destructive'" class="text-xs" :class="discoveryConfig ? 'bg-green-500/10 border-green-500/30 text-green-500' : ''">
                  {{ discoveryConfig ? t('oidc.config.discovered') : t('oidc.config.unavailable') }}
                </Badge>
              </div>
            </CardHeader>
            <CardContent class="p-5">
              <div class="grid grid-cols-1 gap-3">
                <div class="flex justify-between items-start text-sm">
                  <span class="text-muted-foreground w-[180px] flex-shrink-0">{{ t('oidc.config.issuer') }}</span>
                  <code class="text-xs bg-accent/50 px-2 py-0.5 rounded font-mono break-all">{{ discoveryConfig?.issuer || '-' }}</code>
                </div>
                <div class="flex justify-between items-start text-sm">
                  <span class="text-muted-foreground w-[180px] flex-shrink-0">{{ t('oidc.config.authorizationEndpoint') }}</span>
                  <code class="text-xs bg-accent/50 px-2 py-0.5 rounded font-mono break-all">{{ discoveryConfig?.authorization_endpoint || '-' }}</code>
                </div>
                <div class="flex justify-between items-start text-sm">
                  <span class="text-muted-foreground w-[180px] flex-shrink-0">{{ t('oidc.config.tokenEndpoint') }}</span>
                  <code class="text-xs bg-accent/50 px-2 py-0.5 rounded font-mono break-all">{{ discoveryConfig?.token_endpoint || '-' }}</code>
                </div>
                <div class="flex justify-between items-start text-sm">
                  <span class="text-muted-foreground w-[180px] flex-shrink-0">{{ t('oidc.config.userinfoEndpoint') }}</span>
                  <code class="text-xs bg-accent/50 px-2 py-0.5 rounded font-mono break-all">{{ discoveryConfig?.userinfo_endpoint || '-' }}</code>
                </div>
                <div class="flex justify-between items-start text-sm">
                  <span class="text-muted-foreground w-[180px] flex-shrink-0">{{ t('oidc.config.jwksUri') }}</span>
                  <code class="text-xs bg-accent/50 px-2 py-0.5 rounded font-mono break-all">{{ discoveryConfig?.jwks_uri || '-' }}</code>
                </div>
                <div class="flex justify-between items-start text-sm">
                  <span class="text-muted-foreground w-[180px] flex-shrink-0">{{ t('oidc.config.revocationEndpoint') }}</span>
                  <code class="text-xs bg-accent/50 px-2 py-0.5 rounded font-mono break-all">{{ discoveryConfig?.revocation_endpoint || '-' }}</code>
                </div>
                <div class="flex justify-between items-start text-sm">
                  <span class="text-muted-foreground w-[180px] flex-shrink-0">{{ t('oidc.config.introspectionEndpoint') }}</span>
                  <code class="text-xs bg-accent/50 px-2 py-0.5 rounded font-mono break-all">{{ discoveryConfig?.introspection_endpoint || '-' }}</code>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card class="bg-card border-border/60 rounded-[18px] shadow-card">
            <CardHeader class="border-b border-border/60 px-5 py-3.5">
              <span class="text-primary font-semibold">{{ t('oidc.config.capabilities') }}</span>
            </CardHeader>
            <CardContent class="p-5 space-y-3">
              <div class="text-sm">
                <span class="text-muted-foreground block mb-1.5">{{ t('oidc.config.responseTypes') }}</span>
                <div class="flex flex-wrap gap-1.5">
                  <Badge v-for="type in discoveryConfig?.response_types_supported" :key="type" variant="secondary" class="text-xs">{{ type }}</Badge>
                </div>
              </div>
              <div class="text-sm">
                <span class="text-muted-foreground block mb-1.5">{{ t('oidc.config.scopes') }}</span>
                <div class="flex flex-wrap gap-1.5">
                  <Badge v-for="scope in discoveryConfig?.scopes_supported" :key="scope" variant="secondary" class="text-xs">{{ scope }}</Badge>
                </div>
              </div>
              <div class="text-sm">
                <span class="text-muted-foreground block mb-1.5">{{ t('oidc.config.signingAlgorithms') }}</span>
                <div class="flex flex-wrap gap-1.5">
                  <Badge v-for="alg in discoveryConfig?.id_token_signing_alg_values_supported" :key="alg" variant="secondary" class="text-xs">{{ alg }}</Badge>
                </div>
              </div>
              <div class="text-sm">
                <span class="text-muted-foreground block mb-1.5">{{ t('oidc.config.authMethods') }}</span>
                <div class="flex flex-wrap gap-1.5">
                  <Badge v-for="method in discoveryConfig?.token_endpoint_auth_methods_supported" :key="method" variant="secondary" class="text-xs">{{ method }}</Badge>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <div class="space-y-5">
          <Card class="bg-card border-border/60 rounded-[18px] shadow-card">
            <CardHeader class="border-b border-border/60 px-5 py-3.5">
              <span class="text-primary font-semibold">{{ t('oidc.config.runtimeConfig') }}</span>
            </CardHeader>
            <CardContent class="p-5 space-y-3">
              <div class="flex justify-between text-sm">
                <span class="text-muted-foreground">{{ t('oidc.config.enabled') }}</span>
                <Badge :variant="envConfig?.enabled ? 'outline' : 'destructive'" class="text-xs" :class="envConfig?.enabled ? 'bg-green-500/10 border-green-500/30 text-green-500' : ''">
                  {{ envConfig?.enabled ? t('common.yes') : t('common.no') }}
                </Badge>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-muted-foreground">{{ t('oidc.config.enforcePkce') }}</span>
                <Badge :variant="envConfig?.enforcePkce ? 'outline' : 'default'" class="text-xs" :class="envConfig?.enforcePkce ? 'bg-green-500/10 border-green-500/30 text-green-500' : ''">
                  {{ envConfig?.enforcePkce ? t('common.yes') : t('common.no') }}
                </Badge>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-muted-foreground">{{ t('oidc.config.accessTokenLifespan') }}</span>
                <span class="text-foreground">{{ envConfig?.accessTokenLifespan || '-' }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-muted-foreground">{{ t('oidc.config.refreshTokenLifespan') }}</span>
                <span class="text-foreground">{{ envConfig?.refreshTokenLifespan || '-' }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-muted-foreground">{{ t('oidc.config.authCodeLifespan') }}</span>
                <span class="text-foreground">{{ envConfig?.authCodeLifespan || '-' }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-muted-foreground">{{ t('oidc.config.idTokenLifespan') }}</span>
                <span class="text-foreground">{{ envConfig?.idTokenLifespan || '-' }}</span>
              </div>
            </CardContent>
          </Card>

          <Card class="bg-card border-border/60 rounded-[18px] shadow-card">
            <CardHeader class="border-b border-border/60 px-5 py-3.5">
              <span class="text-primary font-semibold">{{ t('oidc.config.quickLinks') }}</span>
            </CardHeader>
            <CardContent class="p-5 space-y-2.5">
              <Button variant="outline" class="w-full justify-start" @click="showJwksDialog = true">
                <LinkIcon class="w-4 h-4 mr-2" />{{ t('oidc.config.viewJwks') }}
              </Button>
              <Button variant="outline" class="w-full justify-start" @click="showDiscoveryDialog = true">
                <LinkIcon class="w-4 h-4 mr-2" />{{ t('oidc.config.viewDiscovery') }}
              </Button>
              <Button variant="outline" class="w-full justify-start" @click="router.push('/system-settings/oidc/integration')">
                <FileText class="w-4 h-4 mr-2" />{{ t('oidc.config.integrationGuide') }}
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>

    <Dialog v-model:open="showJwksDialog">
      <DialogContent class="max-w-[70%]">
        <DialogHeader>
          <DialogTitle>{{ t('oidc.config.viewJwks') }}</DialogTitle>
        </DialogHeader>
        <pre class="json-viewer" v-if="jwksLoading">Loading...</pre>
        <pre class="json-viewer" v-else>{{ jwksContent ? JSON.stringify(jwksContent, null, 2) : t('common.noData') }}</pre>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="showDiscoveryDialog">
      <DialogContent class="max-w-[70%]">
        <DialogHeader>
          <DialogTitle>{{ t('oidc.config.viewDiscovery') }}</DialogTitle>
        </DialogHeader>
        <pre class="json-viewer">{{ discoveryJson || t('common.noData') }}</pre>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { oidcProviderApi } from '@/api/oidc'
import DetailPageSkeleton from '@/components/skeleton/DetailPageSkeleton.vue'
import type { OIDCDiscoveryConfig, OIDCJWKSResponse } from '@/types/oidc'
import { RefreshCw, Link as LinkIcon, FileText } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { Card, CardHeader, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'

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

onMounted(async () => {
  try {
    await fetchDiscovery()
  } finally {
    initialLoading.value = false
  }
})
</script>

<style scoped>
.oidc-config {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--c-text-main);
}

.header-actions {
  display: flex;
  gap: 8px;
}

.content-row {
  flex: 1;
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
