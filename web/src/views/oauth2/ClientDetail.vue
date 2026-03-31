<template>
  <div class="oauth2-client-detail">
    <el-page-header @back="router.back()">
      <template #content>
        <span v-if="clientDetail">{{ clientDetail.client_name }}</span>
      </template>
      <template #extra>
        <el-button-group v-if="clientDetail">
          <el-button @click="showEditDialog = true">{{ t('common.edit') }}</el-button>
          <el-button type="warning" @click="handleRotateSecret">{{ t('oauth2.client.rotateSecret') }}</el-button>
          <el-button type="danger" @click="handleDelete">{{ t('common.delete') }}</el-button>
        </el-button-group>
      </template>
    </el-page-header>

    <div class="content-row">
      <DetailPageSkeleton v-if="initialLoading" :side-span="16" :main-span="8" :side-cards="1" :main-cards="1" :show-avatar="false" :side-item-counts="[8]" :main-item-counts="[4]" />
      <el-row v-else v-loading="loading" :gutter="20">
        <el-col :span="16">
          <el-card shadow="never">
            <template #header>{{ t('oauth2.client.title') }}</template>
            <el-descriptions :column="2" border>
              <el-descriptions-item :label="t('oauth2.client.clientName')">{{ clientDetail?.client_name }}</el-descriptions-item>
              <el-descriptions-item :label="t('oauth2.client.isActive')">
                <el-tag :type="clientDetail?.is_active ? 'success' : 'danger'" size="small" effect="dark">
                  {{ clientDetail?.is_active ? t('oauth2.client.enabled') : t('oauth2.client.disabled') }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oauth2.client.clientId')" :span="2">
                <code class="code-text">{{ clientDetail?.client_id }}</code>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oauth2.client.clientType')">
                <el-tag :type="clientDetail?.client_type === 'confidential' ? 'primary' : 'warning'" size="small">
                  {{ clientDetail?.client_type === 'confidential' ? t('oauth2.client.confidential') : t('oauth2.client.public') }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('common.description')">{{ clientDetail?.description || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('oauth2.client.grantTypes')" :span="2">
                <el-tag v-for="gt in clientDetail?.grant_types" :key="gt" size="small" style="margin-right: 4px;">{{ formatGrantType(gt) }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('oauth2.client.accessTokenLifespan')">{{ clientDetail?.access_token_lifespan }}{{ t('oauth2.client.seconds') }}</el-descriptions-item>
              <el-descriptions-item :label="t('oauth2.client.refreshTokenLifespan')">{{ clientDetail?.refresh_token_lifespan }}{{ t('oauth2.client.seconds') }}</el-descriptions-item>
              <el-descriptions-item :label="t('common.createTime')">{{ formatTime(clientDetail?.created_at) }}</el-descriptions-item>
              <el-descriptions-item :label="t('common.updateTime')">{{ formatTime(clientDetail?.updated_at) }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="never">
            <template #header>{{ t('oauth2.client.redirectUris') }}</template>
            <div v-if="clientDetail?.redirect_uris?.length" class="uri-list-detail">
              <div v-for="uri in clientDetail.redirect_uris" :key="uri" class="uri-item-detail">
                <el-icon><Link /></el-icon>
                <code>{{ uri }}</code>
              </div>
            </div>
            <el-empty v-else :description="t('common.noData')" :image-size="60" />
          </el-card>
          <el-card shadow="never" style="margin-top: 16px;">
            <template #header>{{ t('oauth2.client.scopes') }}</template>
            <div v-if="clientDetail?.scopes?.length">
              <el-tag v-for="scope in clientDetail.scopes" :key="scope" size="small" style="margin: 2px 4px;">{{ scope }}</el-tag>
            </div>
            <el-empty v-else :description="t('common.noData')" :image-size="60" />
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 编辑对话框 -->
    <el-dialog v-model="showEditDialog" :title="t('oauth2.client.editClient')" width="600px" destroy-on-close>
      <el-form ref="editFormRef" :model="editForm" label-width="140px">
        <el-form-item :label="t('oauth2.client.clientName')">
          <el-input v-model="editForm.client_name" />
        </el-form-item>
        <el-form-item :label="t('common.description')">
          <el-input v-model="editForm.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item :label="t('oauth2.client.isActive')">
          <el-switch v-model="editForm.is_active" />
        </el-form-item>
        <el-form-item :label="t('oauth2.client.redirectUris')">
          <div class="uri-list">
            <div v-for="(_, idx) in editForm.redirect_uris" :key="idx" class="uri-item">
              <el-input v-model="editForm.redirect_uris[idx]" />
              <el-button link type="danger" @click="editForm.redirect_uris.splice(idx, 1)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
            <el-button link type="primary" @click="editForm.redirect_uris.push('')">+ URI</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleUpdate">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <!-- 密钥展示对话框 -->
    <el-dialog v-model="showSecretDialog" title="New Client Secret" width="500px" :close-on-click-modal="false">
      <el-alert :title="t('oauth2.client.secretWarning')" type="warning" show-icon :closable="false" style="margin-bottom: 16px;" />
      <div class="secret-display">
        <code>{{ newSecret }}</code>
        <el-button size="small" @click="copySecret">{{ t('oauth2.client.copySecret') }}</el-button>
      </div>
      <template #footer>
        <el-button type="primary" @click="showSecretDialog = false">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Link, Delete } from '@element-plus/icons-vue'
import DetailPageSkeleton from '@/components/skeleton/DetailPageSkeleton.vue'
import { oauth2Api } from '@/api/oauth2'
import { formatTimestamp } from '@/utils/date'
import type { OAuth2Client } from '@/types/oauth2'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const initialLoading = ref(true)
const loading = ref(false)
const submitting = ref(false)
const clientDetail = ref<OAuth2Client | null>(null)
const showEditDialog = ref(false)
const showSecretDialog = ref(false)
const newSecret = ref('')

const editForm = reactive({
  client_name: '',
  description: '',
  is_active: true,
  redirect_uris: [''] as string[],
})

const formatGrantType = (gt: string) => {
  const map: Record<string, string> = {
    authorization_code: t('oauth2.client.grantType.authorizationCode'),
    refresh_token: t('oauth2.client.grantType.refreshToken'),
  }
  return map[gt] || gt
}

const formatTime = (ts?: number) => ts ? formatTimestamp(ts) : '-'

const fetchData = async () => {
  loading.value = true
  try {
    const response = await oauth2Api.getClient(route.params.id as string)
    clientDetail.value = response.client
  } finally {
    loading.value = false
  }
}

const handleUpdate = async () => {
  submitting.value = true
  try {
    const uris = editForm.redirect_uris.filter(u => u.trim() !== '')
    await oauth2Api.updateClient(route.params.id as string, {
      client_name: editForm.client_name,
      description: editForm.description,
      is_active: editForm.is_active,
      redirect_uris: uris,
    })
    showEditDialog.value = false
    ElMessage.success(t('common.operationSuccess'))
    await fetchData()
  } finally {
    submitting.value = false
  }
}

const handleRotateSecret = async () => {
  await ElMessageBox.confirm(t('oauth2.client.rotateSecretConfirm'), t('common.confirm'), { type: 'warning' })
  const response = await oauth2Api.rotateClientSecret(route.params.id as string)
  newSecret.value = response.client_secret
  showSecretDialog.value = true
}

const handleDelete = async () => {
  await ElMessageBox.confirm(t('oauth2.client.deleteConfirm'), t('common.confirm'), { type: 'warning' })
  await oauth2Api.deleteClient(route.params.id as string)
  ElMessage.success(t('common.operationSuccess'))
  router.push('/system-settings/oauth2/clients')
}

const copySecret = () => {
  navigator.clipboard.writeText(newSecret.value)
  ElMessage.success('Copied!')
}

onMounted(async () => {
  try {
    await fetchData()
    if (clientDetail.value) {
      Object.assign(editForm, {
        client_name: clientDetail.value.client_name,
        description: clientDetail.value.description || '',
        is_active: clientDetail.value.is_active,
        redirect_uris: clientDetail.value.redirect_uris?.length ? [...clientDetail.value.redirect_uris] : [''],
      })
    }
  } finally {
    initialLoading.value = false
  }
})
</script>

<style scoped lang="scss">
.oauth2-client-detail {
  display: flex;
  flex-direction: column;
  gap: 16px;
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

.uri-list-detail {
  .uri-item-detail {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 0;
    border-bottom: 1px solid var(--c-border);

    &:last-child {
      border-bottom: none;
    }

    code {
      font-size: 13px;
      color: var(--c-text-main);
      word-break: break-all;
    }
  }
}

.uri-list {
  width: 100%;

  .uri-item {
    display: flex;
    gap: 8px;
    margin-bottom: 8px;
  }
}

.secret-display {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--bg-base);
  border-radius: 8px;
  border: 1px solid var(--c-border);

  code {
    flex: 1;
    word-break: break-all;
    font-family: 'JetBrains Mono', monospace;
    font-size: 13px;
    color: var(--c-accent);
  }
}
</style>
