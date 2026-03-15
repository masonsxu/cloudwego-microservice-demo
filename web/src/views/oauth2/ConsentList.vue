<template>
  <div class="oauth2-consent-list">
    <div class="page-header">
      <h2>{{ t('oauth2.consent.title') }}</h2>
    </div>

    <div class="table-card spotlight-card">
      <ListPageSkeleton v-if="initialLoading" :columns="4" :rows="5" />
      <template v-else>
        <el-empty v-if="consentList.length === 0" :description="t('oauth2.consent.noConsents')" />
        <div v-else class="consent-grid">
          <div v-for="consent in consentList" :key="consent.id" class="consent-card spotlight-card">
            <div class="consent-header">
              <div class="consent-name">{{ consent.client_name || consent.client_id }}</div>
              <el-button type="danger" size="small" plain @click="handleRevoke(consent)">
                {{ t('oauth2.consent.revoke') }}
              </el-button>
            </div>
            <div class="consent-info">
              <div class="consent-field">
                <span class="field-label text-sub">{{ t('oauth2.consent.grantedAt') }}</span>
                <span class="field-value">{{ formatTime(consent.granted_at) }}</span>
              </div>
              <div class="consent-field">
                <span class="field-label text-sub">{{ t('oauth2.consent.authorizedScopes') }}</span>
                <div class="scope-tags">
                  <el-tag v-for="scope in consent.scopes" :key="scope" size="small" type="info">{{ scope }}</el-tag>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import ListPageSkeleton from '@/components/skeleton/ListPageSkeleton.vue'
import { oauth2Api } from '@/api/oauth2'
import { formatTimestamp } from '@/utils/date'
import type { OAuth2Consent } from '@/types/oauth2'

const { t } = useI18n()

const initialLoading = ref(true)
const consentList = ref<OAuth2Consent[]>([])

const formatTime = (ts?: number) => ts ? formatTimestamp(ts) : '-'

const loadConsents = async () => {
  try {
    const response = await oauth2Api.getMyConsents()
    consentList.value = response.consents || []
  } finally {
    initialLoading.value = false
  }
}

const handleRevoke = async (consent: OAuth2Consent) => {
  await ElMessageBox.confirm(t('oauth2.consent.revokeConfirm'), t('common.confirm'), { type: 'warning' })
  await oauth2Api.revokeMyConsent(consent.client_id)
  ElMessage.success(t('common.operationSuccess'))
  await loadConsents()
}

onMounted(() => {
  loadConsents()
})
</script>

<style scoped lang="scss">
.oauth2-consent-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--c-text-main);
}

.table-card {
  padding: 20px;
  border-radius: 12px;
  background: var(--bg-card);
  border: 1px solid var(--c-border);
}

.consent-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 16px;
}

.consent-card {
  padding: 20px;
  border-radius: 12px;
  background: var(--bg-base);
  border: 1px solid var(--c-border);
}

.consent-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  .consent-name {
    font-size: 16px;
    font-weight: 600;
    color: var(--c-text-main);
  }
}

.consent-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.consent-field {
  .field-label {
    display: block;
    font-size: 12px;
    margin-bottom: 4px;
  }

  .field-value {
    font-size: 14px;
    color: var(--c-text-main);
  }
}

.scope-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}
</style>
