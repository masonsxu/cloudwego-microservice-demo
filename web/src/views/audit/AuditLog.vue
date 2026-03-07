<template>
  <div class="audit-log-list">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('audit.title') }}</h1>
        <p class="page-subtitle">{{ t('audit.subtitle') }}</p>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-row">
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon"><el-icon size="20"><Document /></el-icon></div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.total_count }}</span>
          <span class="stat-label">{{ t('audit.totalRecords') }}</span>
        </div>
      </div>
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon success-icon"><el-icon size="20"><CircleCheck /></el-icon></div>
        <div class="stat-info">
          <span class="stat-value">{{ successRate }}%</span>
          <span class="stat-label">{{ t('audit.successRate') }}</span>
        </div>
      </div>
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon duration-icon"><el-icon size="20"><Timer /></el-icon></div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.avg_duration_ms }}ms</span>
          <span class="stat-label">{{ t('audit.avgDuration') }}</span>
        </div>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="search-section">
      <el-form :inline="true" class="search-form">
        <el-form-item>
          <el-select v-model="searchForm.action" :placeholder="t('audit.allActions')" clearable @clear="handleSearch" style="width: 150px">
            <el-option :label="t('audit.actionType.create')" :value="1" />
            <el-option :label="t('audit.actionType.update')" :value="2" />
            <el-option :label="t('audit.actionType.delete')" :value="3" />
            <el-option :label="t('audit.actionType.login')" :value="4" />
            <el-option :label="t('audit.actionType.logout')" :value="5" />
            <el-option :label="t('audit.actionType.passwordChange')" :value="6" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchForm.success" :placeholder="t('audit.allStatus')" clearable @clear="handleSearch" style="width: 130px">
            <el-option :label="t('audit.success')" :value="true" />
            <el-option :label="t('audit.failed')" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            :start-placeholder="t('audit.dateRange')"
            end-placeholder=""
            value-format="x"
            style="width: 340px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>{{ t('common.search') }}
          </el-button>
          <el-button @click="handleReset">
            <el-icon><RefreshLeft /></el-icon>{{ t('common.reset') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 审计日志表格 -->
    <div class="table-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
      <ListPageSkeleton v-if="initialLoading" :columns="7" :rows="10" />
      <template v-else>
        <div class="table-body">
        <el-table v-loading="loading" :data="auditLogs" class="modern-table" height="100%" style="width: 100%">
          <el-table-column :label="t('common.createTime')" width="170">
            <template #default="{ row }">
              <span class="text-sub time-text">{{ formatTimestamp(row.created_at) }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('audit.user')" min-width="120">
            <template #default="{ row }">
              <span class="user-text">{{ row.username || row.user_id || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('audit.action')" width="110" align="center">
            <template #default="{ row }">
              <el-tag :type="getActionTagType(row.action)" size="small" class="action-tag">
                {{ getActionText(row.action) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="resource" :label="t('audit.resource')" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">
              <span class="resource-text">{{ row.resource }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('audit.statusCode')" width="100" align="center">
            <template #default="{ row }">
              <span class="status-code" :class="getStatusCodeClass(row.status_code)">{{ row.status_code }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.status')" width="90" align="center">
            <template #default="{ row }">
              <el-tag :type="row.success ? 'success' : 'danger'" size="small">
                {{ row.success ? t('audit.success') : t('audit.failed') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('audit.clientIP')" width="140">
            <template #default="{ row }">
              <span class="text-sub">{{ row.client_ip }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('audit.duration')" width="100" align="right">
            <template #default="{ row }">
              <span class="duration-text" :class="{ 'duration-slow': row.duration_ms >= 1000 }">
                {{ row.duration_ms }}ms
              </span>
            </template>
          </el-table-column>
        </el-table>
        </div>

        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.size"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadAuditLogs"
            @current-change="loadAuditLogs"
          />
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Document, CircleCheck, Timer, Search, RefreshLeft } from '@element-plus/icons-vue'
import { auditApi } from '@/api/audit'
import type { AuditLogDTO } from '@/api/audit'
import type { AuditLogStats } from '@/types/audit'
import ListPageSkeleton from '@/components/skeleton/ListPageSkeleton.vue'

const { t } = useI18n()

const initialLoading = ref(true)
const loading = ref(false)
const auditLogs = ref<AuditLogDTO[]>([])
const dateRange = ref<[number, number] | null>(null)
const stats = ref<AuditLogStats>({ total_count: 0, success_count: 0, avg_duration_ms: 0 })

const searchForm = reactive({
  action: undefined as number | undefined,
  success: undefined as boolean | undefined,
})

const pagination = reactive({ page: 1, size: 20, total: 0 })

const successRate = computed(() => {
  if (stats.value.total_count === 0) return 0
  return Math.round((stats.value.success_count / stats.value.total_count) * 100)
})

function onMouseMove(e: MouseEvent) {
  const el = e.currentTarget as HTMLElement
  const rect = el.getBoundingClientRect()
  el.style.setProperty('--mouse-x', `${e.clientX - rect.left}px`)
  el.style.setProperty('--mouse-y', `${e.clientY - rect.top}px`)
}

function onMouseLeave(e: MouseEvent) {
  const el = e.currentTarget as HTMLElement
  el.style.removeProperty('--mouse-x')
  el.style.removeProperty('--mouse-y')
}

const loadAuditLogs = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: pagination.page,
      limit: pagination.size,
      action: searchForm.action,
      success: searchForm.success,
    }
    if (dateRange.value) {
      params.start_time = dateRange.value[0]
      params.end_time = dateRange.value[1]
    }
    const response = await auditApi.getAuditLogs(params)
    auditLogs.value = response.audit_logs || []
    pagination.total = response.page?.total || 0
    if (response.stats) {
      stats.value = response.stats
    }
  } catch {
    ElMessage.error(t('common.operationFailed'))
  } finally {
    loading.value = false
    initialLoading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; loadAuditLogs() }
const handleReset = () => {
  searchForm.action = undefined
  searchForm.success = undefined
  dateRange.value = null
  pagination.page = 1
  loadAuditLogs()
}

const getActionText = (action: number) => {
  const map: Record<number, string> = {
    1: t('audit.actionType.create'),
    2: t('audit.actionType.update'),
    3: t('audit.actionType.delete'),
    4: t('audit.actionType.login'),
    5: t('audit.actionType.logout'),
    6: t('audit.actionType.passwordChange'),
  }
  return map[action] || String(action)
}

const getActionTagType = (action: number) => {
  const map: Record<number, string> = {
    1: '',        // 创建 - 蓝色 (default)
    2: 'warning', // 更新 - 橙色
    3: 'danger',  // 删除 - 红色
    4: 'success', // 登录 - 绿色
    5: 'info',    // 登出 - 灰色
    6: '',        // 密码修改 - 蓝色
  }
  return map[action] || 'info'
}

const getStatusCodeClass = (code: number) => {
  if (code >= 200 && code < 300) return 'status-2xx'
  if (code >= 400 && code < 500) return 'status-4xx'
  if (code >= 500) return 'status-5xx'
  return ''
}

const formatTimestamp = (ts?: number) => {
  if (!ts) return '-'
  const d = new Date(ts)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

onMounted(() => { loadAuditLogs() })
</script>

<style scoped lang="scss">
.audit-log-list {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 108px);

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 28px;

    .header-left {
      .page-title {
        font-size: 26px;
        font-weight: 700;
        font-family: 'Cinzel', serif;
        background: linear-gradient(to right, #D4AF37, #F2D288, #D4AF37);
        background-clip: text;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        margin: 0 0 6px;
        line-height: 1.2;
      }

      .page-subtitle {
        color: var(--c-text-sub);
        font-size: 13px;
        margin: 0;
      }
    }
  }

  .stats-row {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
    margin-bottom: 24px;

    .stat-card {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 20px 24px;
      background: var(--bg-card);
      border: 1px solid var(--c-border-accent);
      border-radius: 14px;
      box-shadow: var(--shadow-card);

      .stat-icon {
        width: 44px;
        height: 44px;
        border-radius: 10px;
        background: rgba(212, 175, 55, 0.12);
        border: 1px solid var(--c-border-accent);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--c-accent);
        flex-shrink: 0;

        &.success-icon {
          background: rgba(103, 194, 58, 0.1);
          border-color: rgba(103, 194, 58, 0.25);
          color: #67C23A;
        }

        &.duration-icon {
          background: rgba(64, 158, 255, 0.1);
          border-color: rgba(64, 158, 255, 0.25);
          color: #409EFF;
        }
      }

      .stat-info {
        display: flex;
        flex-direction: column;
        gap: 2px;

        .stat-value {
          font-size: 24px;
          font-weight: 700;
          color: var(--c-text-main);
          font-family: 'JetBrains Mono', monospace;
          line-height: 1;
        }

        .stat-label {
          font-size: 12px;
          color: var(--c-text-sub);
        }
      }
    }
  }

  .search-section {
    background: var(--bg-card);
    border: 1px solid var(--c-border-accent);
    border-radius: 14px;
    padding: 16px 20px;
    margin-bottom: 20px;

    .search-form {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      :deep(.el-form-item) { margin-bottom: 0; margin-right: 12px; }
    }
  }

  .table-card {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
    background: var(--bg-card);
    border: 1px solid var(--c-border-accent);
    border-radius: 14px;
    overflow: hidden;
    box-shadow: var(--shadow-card);

    .table-body {
      flex: 1;
      min-height: 0;
    }

    .modern-table {
      :deep(th.el-table__cell) { padding: 14px 12px; font-size: 12px; letter-spacing: 0.05em; text-transform: uppercase; }
      :deep(td.el-table__cell) { padding: 14px 12px; }
    }

    .text-sub { color: var(--c-text-sub); font-size: 13px; }
    .time-text { font-size: 12px; font-family: 'JetBrains Mono', monospace; }
    .user-text { font-weight: 500; color: var(--c-text-main); }

    .resource-text {
      font-family: 'JetBrains Mono', monospace;
      font-size: 12px;
      color: var(--c-text-sub);
    }

    .action-tag {
      min-width: 56px;
      text-align: center;
    }

    .status-code {
      font-family: 'JetBrains Mono', monospace;
      font-weight: 600;
      font-size: 13px;

      &.status-2xx { color: #67C23A; }
      &.status-4xx { color: #E6A23C; }
      &.status-5xx { color: #F56C6C; }
    }

    .duration-text {
      font-family: 'JetBrains Mono', monospace;
      font-size: 12px;
      color: var(--c-text-sub);

      &.duration-slow {
        color: #F56C6C;
        font-weight: 600;
      }
    }

    .pagination-wrapper {
      display: flex;
      justify-content: flex-end;
      padding: 16px 20px;
      border-top: 1px solid var(--c-border);
    }
  }
}

.spotlight-card {
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(500px circle at var(--mouse-x, -100%) var(--mouse-y, -100%),
      rgba(212, 175, 55, 0.05), transparent 50%);
    pointer-events: none;
    z-index: 1;
  }
}
</style>
