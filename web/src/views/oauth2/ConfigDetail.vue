<template>
  <div class="oauth2-config-detail">
    <div class="page-header">
      <h2>{{ t('oauth2.config.title') }}</h2>
      <div class="header-actions">
        <el-button @click="router.push('/system-settings/oauth2/clients')">{{ t('oauth2.client.title') }}</el-button>
        <el-button @click="router.push('/system-settings/oauth2/consents')">{{ t('oauth2.consent.title') }}</el-button>
      </div>
    </div>

    <div class="content-row">
      <DetailPageSkeleton
        v-if="initialLoading"
        :side-span="16"
        :main-span="8"
        :side-cards="1"
        :main-cards="1"
        :show-avatar="false"
        :side-item-counts="[7]"
        :main-item-counts="[1]"
      />
      <el-row v-else v-loading="loading" :gutter="20">
        <el-col :span="16">
          <el-card shadow="never">
            <template #header>{{ t('oauth2.config.title') }}</template>
            <el-descriptions :column="2" border>
              <el-descriptions-item :label="t('oauth2.config.enabled')">
                <el-tag :type="config?.enabled ? 'success' : 'danger'" size="small" effect="dark">
                  {{ config?.enabled ? t('oauth2.client.enabled') : t('oauth2.client.disabled') }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oauth2.config.enforcePkce')">
                <el-tag :type="config?.enforce_pkce ? 'success' : 'warning'" size="small" effect="dark">
                  {{ config?.enforce_pkce ? t('common.yes') : t('common.no') }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oauth2.config.issuer')" :span="2">
                <code class="code-text">{{ config?.issuer || '-' }}</code>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oauth2.config.accessTokenLifespan')">
                {{ config?.access_token_lifespan ?? 0 }}{{ t('oauth2.client.seconds') }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('oauth2.config.refreshTokenLifespan')">
                {{ config?.refresh_token_lifespan ?? 0 }}{{ t('oauth2.client.seconds') }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('oauth2.config.authCodeLifespan')">
                {{ config?.auth_code_lifespan ?? 0 }}{{ t('oauth2.client.seconds') }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('oauth2.config.consentPageUrl')" :span="2">
                <code class="code-text">{{ config?.consent_page_url || '-' }}</code>
              </el-descriptions-item>
            </el-descriptions>
          </el-card>
        </el-col>

        <el-col :span="8">
          <el-card shadow="never">
            <template #header>{{ t('oauth2.config.availableScopes') }}</template>
            <div v-if="scopeList.length">
              <el-tag
                v-for="scope in scopeList"
                :key="scope.name"
                size="small"
                style="margin: 2px 4px;"
              >
                {{ scope.name }}
              </el-tag>
            </div>
            <el-empty v-else :description="t('common.noData')" :image-size="60" />
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import DetailPageSkeleton from '@/components/skeleton/DetailPageSkeleton.vue'
import { oauth2Api } from '@/api/oauth2'
import type { OAuth2Config, OAuth2Scope } from '@/types/oauth2'

const router = useRouter()
const { t } = useI18n()

const initialLoading = ref(true)
const loading = ref(false)
const config = ref<OAuth2Config | null>(null)
const scopeList = ref<OAuth2Scope[]>([])

const fetchData = async () => {
  loading.value = true
  try {
    const [configResp, scopeResp] = await Promise.all([
      oauth2Api.getConfig(),
      oauth2Api.getScopes(),
    ])
    config.value = configResp.config
    scopeList.value = scopeResp.scopes || []
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    await fetchData()
  } finally {
    initialLoading.value = false
  }
})
</script>

<style scoped lang="scss">
.oauth2-config-detail {
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

.code-text {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
  background: var(--bg-base);
  color: var(--c-accent);
}
</style>
