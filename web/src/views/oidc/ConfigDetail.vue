<template>
  <div class="space-y-6">
    <!-- 标题区（DESIGN.md 详情页范式） -->
    <header class="flex items-start justify-between gap-4 pb-5 border-b border-subtle">
      <div class="flex items-center gap-3 min-w-0">
        <span class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-md bg-[color:var(--color-primary-soft)] text-[color:var(--color-primary-active)]">
          <ShieldCheck class="h-5 w-5" />
        </span>
        <div class="min-w-0">
          <div class="flex items-center gap-2">
            <h1 class="text-[22px] font-semibold leading-tight tracking-[-0.012em] text-[color:var(--color-ink-strong)]">
              {{ t('oidc.config.title') }}
            </h1>
            <Badge :variant="discoveryConfig ? 'success' : 'destructive'">
              <span class="relative inline-flex h-1.5 w-1.5 mr-1 rounded-full" :class="discoveryConfig ? 'bg-[color:var(--color-success)]' : 'bg-[color:var(--color-danger)]'" />
              {{ discoveryConfig ? t('oidc.config.discovered') : t('oidc.config.unavailable') }}
            </Badge>
          </div>
          <p class="mt-1 text-[13px] text-[color:var(--color-ink-muted)] truncate">
            {{ discoveryConfig?.issuer || t('oidc.config.noIssuer') }}
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2 flex-shrink-0">
        <Button variant="outline" size="sm" @click="showJwksDialog = true">
          <LinkIcon class="h-3.5 w-3.5" />
          {{ t('oidc.config.viewJwks') }}
        </Button>
        <Button variant="outline" size="sm" @click="showDiscoveryDialog = true">
          <FileText class="h-3.5 w-3.5" />
          {{ t('oidc.config.viewDiscovery') }}
        </Button>
        <Button size="sm" @click="refreshData">
          <RefreshCw class="h-3.5 w-3.5" :class="loading ? 'animate-spin' : ''" />
          {{ t('common.refresh') }}
        </Button>
      </div>
    </header>

    <!-- 内容区 -->
    <DetailPageSkeleton
      v-if="initialLoading"
      :side-span="16"
      :main-span="8"
      :side-cards="2"
      :main-cards="1"
      :show-avatar="false"
      :side-item-counts="[8, 4]"
      :main-item-counts="[6]"
    />

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-5">
      <!-- 主区：Discovery + Capabilities -->
      <div class="lg:col-span-2 space-y-5">
        <section class="rounded-md border border-subtle bg-canvas">
          <header class="flex items-center justify-between border-b border-subtle px-5 py-3">
            <h2 class="text-[14px] font-semibold text-ink">{{ t('oidc.config.providerConfig') }}</h2>
            <span class="text-[11px] font-medium uppercase tracking-[0.06em] text-[color:var(--color-ink-subtle)]">Discovery</span>
          </header>
          <dl class="divide-y divide-[color:var(--color-border-subtle)]">
            <div v-for="row in discoveryRows" :key="row.key" class="flex items-start gap-4 px-5 py-3">
              <dt class="w-[200px] flex-shrink-0 text-[13px] font-medium text-[color:var(--color-ink-muted)]">
                {{ t(row.label) }}
              </dt>
              <dd class="flex-1 min-w-0">
                <code v-if="row.value" class="font-mono text-[12px] text-ink break-all">{{ row.value }}</code>
                <span v-else class="text-[color:var(--color-ink-subtle)]">—</span>
              </dd>
            </div>
          </dl>
        </section>

        <section class="rounded-md border border-subtle bg-canvas">
          <header class="flex items-center justify-between border-b border-subtle px-5 py-3">
            <h2 class="text-[14px] font-semibold text-ink">{{ t('oidc.config.capabilities') }}</h2>
            <span class="text-[11px] font-medium uppercase tracking-[0.06em] text-[color:var(--color-ink-subtle)]">Supports</span>
          </header>
          <dl class="divide-y divide-[color:var(--color-border-subtle)]">
            <div v-for="cap in capabilityRows" :key="cap.key" class="flex items-start gap-4 px-5 py-3">
              <dt class="w-[200px] flex-shrink-0 text-[13px] font-medium text-[color:var(--color-ink-muted)]">
                {{ t(cap.label) }}
              </dt>
              <dd class="flex-1 min-w-0">
                <div v-if="cap.values && cap.values.length" class="flex flex-wrap gap-1">
                  <Badge v-for="value in cap.values" :key="value" variant="default">{{ value }}</Badge>
                </div>
                <span v-else class="text-[color:var(--color-ink-subtle)]">—</span>
              </dd>
            </div>
          </dl>
        </section>
      </div>

      <!-- 侧栏：Runtime + Quick Links -->
      <div class="space-y-5">
        <section class="rounded-md border border-subtle bg-canvas">
          <header class="flex items-center justify-between border-b border-subtle px-5 py-3">
            <h2 class="text-[14px] font-semibold text-ink">{{ t('oidc.config.runtimeConfig') }}</h2>
            <span class="text-[11px] font-medium uppercase tracking-[0.06em] text-[color:var(--color-ink-subtle)]">Runtime</span>
          </header>
          <dl class="divide-y divide-[color:var(--color-border-subtle)]">
            <div class="flex items-center justify-between px-5 py-2.5">
              <dt class="text-[13px] text-[color:var(--color-ink-muted)]">{{ t('oidc.config.enabled') }}</dt>
              <Badge :variant="envConfig.enabled ? 'success' : 'destructive'">
                {{ envConfig.enabled ? t('common.yes') : t('common.no') }}
              </Badge>
            </div>
            <div class="flex items-center justify-between px-5 py-2.5">
              <dt class="text-[13px] text-[color:var(--color-ink-muted)]">{{ t('oidc.config.enforcePkce') }}</dt>
              <Badge :variant="envConfig.enforcePkce ? 'success' : 'default'">
                {{ envConfig.enforcePkce ? t('common.yes') : t('common.no') }}
              </Badge>
            </div>
            <div class="flex items-center justify-between px-5 py-2.5">
              <dt class="text-[13px] text-[color:var(--color-ink-muted)]">{{ t('oidc.config.accessTokenLifespan') }}</dt>
              <span class="font-mono tabular text-[13px] text-ink">{{ envConfig.accessTokenLifespan || '—' }}</span>
            </div>
            <div class="flex items-center justify-between px-5 py-2.5">
              <dt class="text-[13px] text-[color:var(--color-ink-muted)]">{{ t('oidc.config.refreshTokenLifespan') }}</dt>
              <span class="font-mono tabular text-[13px] text-ink">{{ envConfig.refreshTokenLifespan || '—' }}</span>
            </div>
            <div class="flex items-center justify-between px-5 py-2.5">
              <dt class="text-[13px] text-[color:var(--color-ink-muted)]">{{ t('oidc.config.authCodeLifespan') }}</dt>
              <span class="font-mono tabular text-[13px] text-ink">{{ envConfig.authCodeLifespan || '—' }}</span>
            </div>
            <div class="flex items-center justify-between px-5 py-2.5">
              <dt class="text-[13px] text-[color:var(--color-ink-muted)]">{{ t('oidc.config.idTokenLifespan') }}</dt>
              <span class="font-mono tabular text-[13px] text-ink">{{ envConfig.idTokenLifespan || '—' }}</span>
            </div>
          </dl>
        </section>

        <section class="rounded-md border border-subtle bg-canvas">
          <header class="border-b border-subtle px-5 py-3">
            <h2 class="text-[14px] font-semibold text-ink">{{ t('oidc.config.quickLinks') }}</h2>
          </header>
          <div class="p-3 space-y-1">
            <button
              class="flex w-full items-center gap-2 rounded-sm px-3 py-2 text-left text-[13px] text-[color:var(--color-ink-muted)] transition-colors hover:bg-sunken hover:text-ink"
              @click="router.push('/system-settings/oidc/integration')"
            >
              <FileText class="h-3.5 w-3.5 flex-shrink-0" />
              <span class="flex-1">{{ t('oidc.config.integrationGuide') }}</span>
              <ArrowUpRight class="h-3 w-3 text-[color:var(--color-ink-subtle)]" />
            </button>
          </div>
        </section>
      </div>
    </div>

    <!-- JWKS / Discovery JSON Dialog -->
    <Dialog v-model:open="showJwksDialog">
      <DialogContent class="max-w-[720px]">
        <DialogHeader>
          <DialogTitle>{{ t('oidc.config.viewJwks') }}</DialogTitle>
        </DialogHeader>
        <pre class="json-viewer" v-if="jwksLoading">Loading...</pre>
        <pre class="json-viewer" v-else>{{ jwksContent ? JSON.stringify(jwksContent, null, 2) : t('common.noData') }}</pre>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="showDiscoveryDialog">
      <DialogContent class="max-w-[720px]">
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
import { RefreshCw, Link as LinkIcon, FileText, ShieldCheck, ArrowUpRight } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
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

const discoveryRows = computed(() => [
  { key: 'issuer', label: 'oidc.config.issuer', value: discoveryConfig.value?.issuer },
  { key: 'auth', label: 'oidc.config.authorizationEndpoint', value: discoveryConfig.value?.authorization_endpoint },
  { key: 'token', label: 'oidc.config.tokenEndpoint', value: discoveryConfig.value?.token_endpoint },
  { key: 'userinfo', label: 'oidc.config.userinfoEndpoint', value: discoveryConfig.value?.userinfo_endpoint },
  { key: 'jwks', label: 'oidc.config.jwksUri', value: discoveryConfig.value?.jwks_uri },
  { key: 'revoke', label: 'oidc.config.revocationEndpoint', value: discoveryConfig.value?.revocation_endpoint },
  { key: 'introspect', label: 'oidc.config.introspectionEndpoint', value: discoveryConfig.value?.introspection_endpoint },
])

const capabilityRows = computed(() => [
  { key: 'response', label: 'oidc.config.responseTypes', values: discoveryConfig.value?.response_types_supported },
  { key: 'scopes', label: 'oidc.config.scopes', values: discoveryConfig.value?.scopes_supported },
  { key: 'signing', label: 'oidc.config.signingAlgorithms', values: discoveryConfig.value?.id_token_signing_alg_values_supported },
  { key: 'auth-methods', label: 'oidc.config.authMethods', values: discoveryConfig.value?.token_endpoint_auth_methods_supported },
])

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
.json-viewer {
  background: var(--color-sunken);
  border: 1px solid var(--color-border-subtle);
  border-radius: 8px;
  padding: 16px;
  overflow: auto;
  max-height: 60vh;
  font-family: var(--font-family-mono);
  font-size: 12px;
  line-height: 1.6;
  color: var(--color-ink);
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
