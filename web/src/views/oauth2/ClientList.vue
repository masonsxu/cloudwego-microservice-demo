<template>
  <div class="oauth2-client-list">
    <div class="page-header">
      <h2>{{ t('oauth2.client.title') }}</h2>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        {{ t('oauth2.client.createClient') }}
      </el-button>
    </div>

    <el-alert class="mvp-guide" type="info" :closable="false" show-icon>
      <template #title>
        {{ t('oauth2.guide.title') }}
      </template>
      <div class="guide-lines">
        <div>{{ t('oauth2.guide.step1') }}</div>
        <div>{{ t('oauth2.guide.step2') }}</div>
        <div>{{ t('oauth2.guide.step3') }}</div>
        <div class="guide-note">{{ t('oauth2.guide.unsupported') }}</div>
      </div>
    </el-alert>

    <div class="stats-row">
      <div class="stat-card spotlight-card">
        <div class="stat-value">{{ stats.total }}</div>
        <div class="stat-label text-sub">{{ t('common.total', { total: '' }) }}</div>
      </div>
      <div class="stat-card spotlight-card">
        <div class="stat-value">{{ stats.active }}</div>
        <div class="stat-label text-sub">{{ t('oauth2.client.enabled') }}</div>
      </div>
      <div class="stat-card spotlight-card">
        <div class="stat-value">{{ stats.inactive }}</div>
        <div class="stat-label text-sub">{{ t('oauth2.client.disabled') }}</div>
      </div>
    </div>

    <div class="table-card spotlight-card">
      <ListPageSkeleton v-if="initialLoading" :columns="6" :rows="8" />
      <template v-else>
        <el-table v-loading="loading" :data="clientList" height="100%" stripe>
          <el-table-column prop="client_name" :label="t('oauth2.client.clientName')" min-width="160">
            <template #default="{ row }">
              <router-link :to="`/system-settings/oauth2/clients/${row.id}`" class="text-link">
                {{ row.client_name }}
              </router-link>
            </template>
          </el-table-column>
          <el-table-column prop="client_id" :label="t('oauth2.client.clientId')" min-width="180">
            <template #default="{ row }">
              <code class="code-text">{{ row.client_id }}</code>
            </template>
          </el-table-column>
          <el-table-column prop="client_type" :label="t('oauth2.client.clientType')" width="130">
            <template #default="{ row }">
              <el-tag :type="row.client_type === 'confidential' ? 'primary' : 'warning'" size="small">
                {{ row.client_type === 'confidential' ? t('oauth2.client.confidential') : t('oauth2.client.public') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="grant_types" :label="t('oauth2.client.grantTypes')" min-width="200">
            <template #default="{ row }">
              <el-tag v-for="gt in row.grant_types" :key="gt" size="small" class="grant-tag">
                {{ formatGrantType(gt) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="is_active" :label="t('oauth2.client.isActive')" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="row.is_active ? 'success' : 'danger'" size="small" effect="dark">
                {{ row.is_active ? t('oauth2.client.enabled') : t('oauth2.client.disabled') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.actions')" width="120" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="router.push(`/system-settings/oauth2/clients/${row.id}`)">
                {{ t('common.edit') }}
              </el-button>
              <el-button link type="danger" @click="handleDelete(row)">
                {{ t('common.delete') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.size"
            :total="pagination.total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @size-change="loadClients"
            @current-change="loadClients"
          />
        </div>
      </template>
    </div>

    <!-- 创建客户端对话框 -->
    <el-dialog v-model="showCreateDialog" :title="t('oauth2.client.createClient')" width="600px" destroy-on-close>
      <el-form ref="formRef" :model="createForm" :rules="formRules" label-width="140px">
        <el-form-item :label="t('oauth2.client.clientName')" prop="client_name">
          <el-input v-model="createForm.client_name" />
        </el-form-item>
        <el-form-item :label="t('common.description')">
          <el-input v-model="createForm.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item :label="t('oauth2.client.clientType')" prop="client_type">
          <el-radio-group v-model="createForm.client_type">
            <el-radio value="confidential">{{ t('oauth2.client.confidential') }}</el-radio>
            <el-radio value="public">{{ t('oauth2.client.public') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="t('oauth2.client.grantTypes')" prop="grant_types">
          <el-checkbox-group v-model="createForm.grant_types">
            <el-checkbox value="authorization_code">{{ t('oauth2.client.grantType.authorizationCode') }}</el-checkbox>
            <el-checkbox value="refresh_token">{{ t('oauth2.client.grantType.refreshToken') }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item v-if="createForm.client_type === 'public'">
          <el-alert type="warning" :title="t('oauth2.client.pkceRequired')" :closable="false" show-icon />
        </el-form-item>
        <el-form-item :label="t('oauth2.client.redirectUris')">
          <div class="uri-list">
            <div v-for="(_, idx) in createForm.redirect_uris" :key="idx" class="uri-item">
              <el-input v-model="createForm.redirect_uris[idx]" placeholder="https://example.com/callback" />
              <el-button link type="danger" @click="createForm.redirect_uris.splice(idx, 1)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
            <el-button link type="primary" @click="createForm.redirect_uris.push('')">
              + {{ t('oauth2.client.redirectUris') }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleCreate">{{ t('common.submit') }}</el-button>
      </template>
    </el-dialog>

    <!-- 密钥展示对话框 -->
    <el-dialog v-model="showSecretDialog" title="Client Secret" width="500px" :close-on-click-modal="false">
      <el-alert :title="t('oauth2.client.secretWarning')" type="warning" show-icon :closable="false" class="secret-alert" />
      <div class="secret-display">
        <code>{{ clientSecret }}</code>
        <el-button size="small" @click="copySecret">{{ t('oauth2.client.copySecret') }}</el-button>
      </div>
      <template #footer>
        <el-button type="primary" @click="showSecretDialog = false">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import ListPageSkeleton from '@/components/skeleton/ListPageSkeleton.vue'
import { oauth2Api } from '@/api/oauth2'
import type { OAuth2Client, OAuth2ClientType, OAuth2GrantType } from '@/types/oauth2'

const router = useRouter()
const { t } = useI18n()

const initialLoading = ref(true)
const loading = ref(false)
const submitting = ref(false)
const clientList = ref<OAuth2Client[]>([])
const pagination = reactive({ page: 1, size: 20, total: 0 })

const showCreateDialog = ref(false)
const showSecretDialog = ref(false)
const clientSecret = ref('')
const formRef = ref<FormInstance>()

const stats = computed(() => {
  const total = clientList.value.length
  const active = clientList.value.filter(c => c.is_active).length
  return { total, active, inactive: total - active }
})

const createForm = reactive<{
  client_name: string
  description: string
  client_type: OAuth2ClientType
  grant_types: OAuth2GrantType[]
  redirect_uris: string[]
}>({
  client_name: '',
  description: '',
  client_type: 'confidential',
  grant_types: ['authorization_code', 'refresh_token'],
  redirect_uris: [''],
})

const formRules: FormRules = {
  client_name: [{ required: true, message: t('oauth2.client.clientName'), trigger: 'blur' }],
  client_type: [{ required: true, message: t('oauth2.client.clientType'), trigger: 'change' }],
  grant_types: [{ required: true, type: 'array', min: 1, message: t('oauth2.client.grantTypes'), trigger: 'change' }],
}

const formatGrantType = (gt: string) => {
  const map: Record<string, string> = {
    authorization_code: t('oauth2.client.grantType.authorizationCode'),
    refresh_token: t('oauth2.client.grantType.refreshToken'),
  }
  return map[gt] || gt
}

const loadClients = async () => {
  loading.value = true
  try {
    const response = await oauth2Api.getClients({
      page: pagination.page,
      limit: pagination.size,
    })
    clientList.value = response.clients || []
    pagination.total = response.page?.total || 0
  } finally {
    loading.value = false
    initialLoading.value = false
  }
}

const handleCreate = async () => {
  if (!formRef.value) return
  await formRef.value.validate()

  submitting.value = true
  try {
    const uris = createForm.redirect_uris.filter(u => u.trim() !== '')
    const response = await oauth2Api.createClient({
      client_name: createForm.client_name,
      description: createForm.description || undefined,
      client_type: createForm.client_type,
      grant_types: createForm.grant_types,
      redirect_uris: uris.length > 0 ? uris : undefined,
    })

    showCreateDialog.value = false
    clientSecret.value = response.client_secret
    showSecretDialog.value = true

    Object.assign(createForm, {
      client_name: '',
      description: '',
      client_type: 'confidential',
      grant_types: ['authorization_code', 'refresh_token'],
      redirect_uris: [''],
    })

    await loadClients()
    ElMessage.success(t('common.operationSuccess'))
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row: OAuth2Client) => {
  await ElMessageBox.confirm(t('oauth2.client.deleteConfirm'), t('common.confirm'), { type: 'warning' })
  await oauth2Api.deleteClient(row.id)
  ElMessage.success(t('common.operationSuccess'))
  await loadClients()
}

const copySecret = () => {
  navigator.clipboard.writeText(clientSecret.value)
  ElMessage.success('Copied!')
}

onMounted(() => {
  loadClients()
})
</script>

<style scoped lang="scss">
.oauth2-client-list {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 108px);
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
}

.mvp-guide {
  :deep(.el-alert__title) {
    font-weight: 600;
  }
}

.guide-lines {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 13px;

  .guide-note {
    color: var(--el-color-warning-dark-2);
    font-weight: 500;
  }
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.stat-card {
  padding: 20px;
  border-radius: 12px;
  background: var(--bg-card);
  border: 1px solid var(--c-border);
  text-align: center;

  .stat-value {
    font-size: 28px;
    font-weight: 700;
    color: var(--c-primary);
  }

  .stat-label {
    margin-top: 4px;
    font-size: 13px;
  }
}

.table-card {
  flex: 1;
  min-height: 0;
  padding: 16px;
  border-radius: 12px;
  background: var(--bg-card);
  border: 1px solid var(--c-border);
  display: flex;
  flex-direction: column;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 12px;
}

.code-text {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
  background: var(--bg-base);
  color: var(--c-accent);
}

.text-link {
  color: var(--c-primary);
  text-decoration: none;
  font-weight: 500;

  &:hover {
    text-decoration: underline;
  }
}

.grant-tag {
  margin-right: 4px;
  margin-bottom: 2px;
}

.uri-list {
  width: 100%;

  .uri-item {
    display: flex;
    gap: 8px;
    margin-bottom: 8px;
  }
}

.secret-alert {
  margin-bottom: 16px;
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
