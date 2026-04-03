<template>
  <div class="audit-log-list">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('audit.title') }}</h1>
        <p class="page-subtitle">{{ t('audit.subtitle') }}</p>
      </div>
    </div>

    <div class="stats-row">
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon"><FileText class="w-5 h-5" /></div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.total_count }}</span>
          <span class="stat-label">{{ t('audit.totalRecords') }}</span>
        </div>
      </div>
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon success-icon"><CircleCheck class="w-5 h-5" /></div>
        <div class="stat-info">
          <span class="stat-value">{{ successRate }}%</span>
          <span class="stat-label">{{ t('audit.successRate') }}</span>
        </div>
      </div>
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon duration-icon"><Timer class="w-5 h-5" /></div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.avg_duration_ms }}ms</span>
          <span class="stat-label">{{ t('audit.avgDuration') }}</span>
        </div>
      </div>
    </div>

    <div class="search-section">
      <div class="search-form">
        <div class="form-field">
          <Select v-model="searchForm.action" @update:modelValue="handleSearch">
            <SelectTrigger class="w-[150px]">
              <SelectValue :placeholder="t('audit.allActions')" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="all">{{ t('audit.allActions') }}</SelectItem>
                <SelectItem value="1">{{ t('audit.actionType.create') }}</SelectItem>
                <SelectItem value="2">{{ t('audit.actionType.update') }}</SelectItem>
                <SelectItem value="3">{{ t('audit.actionType.delete') }}</SelectItem>
                <SelectItem value="4">{{ t('audit.actionType.login') }}</SelectItem>
                <SelectItem value="5">{{ t('audit.actionType.logout') }}</SelectItem>
                <SelectItem value="6">{{ t('audit.actionType.passwordChange') }}</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
        <div class="form-field">
          <Select v-model="searchForm.success" @update:modelValue="handleSearch">
            <SelectTrigger class="w-[130px]">
              <SelectValue :placeholder="t('audit.allStatus')" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="all">{{ t('audit.allStatus') }}</SelectItem>
                <SelectItem value="true">{{ t('audit.success') }}</SelectItem>
                <SelectItem value="false">{{ t('audit.failed') }}</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
        <div class="form-field">
          <Input
            v-model="dateRangeDisplay"
            :placeholder="t('audit.dateRange')"
            class="w-[220px]"
            readonly
            @click="showDatePicker = true"
          />
        </div>
        <div class="form-field">
          <Button variant="default" @click="handleSearch">
            <Search class="w-4 h-4 mr-1" />{{ t('common.search') }}
          </Button>
          <Button variant="outline" @click="handleReset">
            <RotateCcw class="w-4 h-4 mr-1" />{{ t('common.reset') }}
          </Button>
        </div>
      </div>
    </div>

    <div class="table-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
      <ListPageSkeleton v-if="initialLoading" :columns="7" :rows="10" />
      <template v-else>
        <div class="table-body">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead class="w-[170px]">{{ t('common.createTime') }}</TableHead>
                <TableHead class="min-w-[120px]">{{ t('audit.user') }}</TableHead>
                <TableHead class="w-[110px] text-center">{{ t('audit.action') }}</TableHead>
                <TableHead class="min-w-[200px]">{{ t('audit.resource') }}</TableHead>
                <TableHead class="w-[100px] text-center">{{ t('audit.statusCode') }}</TableHead>
                <TableHead class="w-[90px] text-center">{{ t('common.status') }}</TableHead>
                <TableHead class="w-[140px]">{{ t('audit.clientIP') }}</TableHead>
                <TableHead class="w-[100px] text-right">{{ t('audit.duration') }}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="row in auditLogs" :key="row.id">
                <TableCell class="text-[var(--c-text-sub)] text-xs font-mono">{{ formatTimestamp(row.created_at) }}</TableCell>
                <TableCell class="font-medium text-[var(--c-text-main)]">{{ row.username || row.user_id || '-' }}</TableCell>
                <TableCell class="text-center">
                  <Badge :variant="getActionBadgeVariant(row.action)" class="min-w-[56px] justify-center">
                    {{ getActionText(row.action) }}
                  </Badge>
                </TableCell>
                <TableCell class="font-mono text-xs text-[var(--c-text-sub)] truncate max-w-[200px]" :title="row.resource">{{ row.resource }}</TableCell>
                <TableCell class="text-center">
                  <span class="font-mono font-semibold text-sm" :class="getStatusCodeClass(row.status_code)">{{ row.status_code }}</span>
                </TableCell>
                <TableCell class="text-center">
                  <Badge :variant="row.success ? 'default' : 'destructive'" class="text-xs">
                    {{ row.success ? t('audit.success') : t('audit.failed') }}
                  </Badge>
                </TableCell>
                <TableCell class="text-[var(--c-text-sub)] text-xs">{{ row.client_ip }}</TableCell>
                <TableCell class="text-right">
                  <span class="font-mono text-xs text-[var(--c-text-sub)]" :class="{ 'duration-slow': row.duration_ms >= 1000 }">
                    {{ row.duration_ms }}ms
                  </span>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>

        <div class="pagination-wrapper">
          <div class="flex items-center gap-4">
            <span class="text-sm text-[var(--c-text-sub)]">
              {{ t('common.total') }}: {{ pagination.total }}
            </span>
            <Select v-model="paginationSizeString" @update:modelValue="loadAuditLogs">
              <SelectTrigger class="w-[80px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="10">10</SelectItem>
                <SelectItem value="20">20</SelectItem>
                <SelectItem value="50">50</SelectItem>
                <SelectItem value="100">100</SelectItem>
              </SelectContent>
            </Select>
            <div class="flex items-center gap-1">
              <Button variant="outline" size="sm" :disabled="pagination.page <= 1" @click="pagination.page = 1; loadAuditLogs()">
                <ChevronsLeft class="w-4 h-4" />
              </Button>
              <Button variant="outline" size="sm" :disabled="pagination.page <= 1" @click="pagination.page--; loadAuditLogs()">
                <ChevronLeft class="w-4 h-4" />
              </Button>
              <span class="px-3 text-sm">{{ pagination.page }}</span>
              <Button variant="outline" size="sm" :disabled="pagination.page >= Math.ceil(pagination.total / pagination.size)" @click="pagination.page++; loadAuditLogs()">
                <ChevronRight class="w-4 h-4" />
              </Button>
              <Button variant="outline" size="sm" :disabled="pagination.page >= Math.ceil(pagination.total / pagination.size)" @click="pagination.page = Math.ceil(pagination.total / pagination.size); loadAuditLogs()">
                <ChevronsRight class="w-4 h-4" />
              </Button>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { FileText, CircleCheck, Timer, Search, RotateCcw, ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from 'lucide-vue-next'
import { auditApi } from '@/api/audit'
import type { AuditLogDTO } from '@/api/audit'
import type { AuditLogStats } from '@/types/audit'
import ListPageSkeleton from '@/components/skeleton/ListPageSkeleton.vue'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'

const { t } = useI18n()

const initialLoading = ref(true)
const loading = ref(false)
const auditLogs = ref<AuditLogDTO[]>([])
const dateRange = ref<[number, number] | null>(null)
const showDatePicker = ref(false)
const stats = ref<AuditLogStats>({ total_count: 0, success_count: 0, avg_duration_ms: 0 })

const searchForm = reactive({
  action: 'all',
  success: 'all',
})

const pagination = reactive({ page: 1, size: 20, total: 0 })

const paginationSizeString = computed({
  get: () => String(pagination.size),
  set: (val: string) => { pagination.size = Number(val) },
})

const dateRangeDisplay = computed(() => {
  if (!dateRange.value) return ''
  const pad = (n: number) => String(n).padStart(2, '0')
  const start = new Date(dateRange.value[0])
  const end = new Date(dateRange.value[1])
  return `${start.getFullYear()}-${pad(start.getMonth() + 1)}-${pad(start.getDate())} ~ ${end.getFullYear()}-${pad(end.getMonth() + 1)}-${pad(end.getDate())}`
})

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
      action: searchForm.action !== 'all' ? Number(searchForm.action) : undefined,
      success: searchForm.success === 'true' ? true : searchForm.success === 'false' ? false : undefined,
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
    toast.error(t('common.operationFailed'))
  } finally {
    loading.value = false
    initialLoading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; loadAuditLogs() }
const handleReset = () => {
  searchForm.action = 'all'
  searchForm.success = 'all'
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

const getActionBadgeVariant = (action: number): 'default' | 'secondary' | 'destructive' | 'outline' => {
  const map: Record<number, 'default' | 'secondary' | 'destructive' | 'outline'> = {
    1: 'default',
    2: 'secondary',
    3: 'destructive',
    4: 'default',
    5: 'outline',
    6: 'default',
  }
  return map[action] || 'outline'
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
      gap: 12px;
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

    .duration-slow {
      color: #F56C6C;
      font-weight: 600;
    }

    .status-2xx { color: #67C23A; }
    .status-4xx { color: #E6A23C; }
    .status-5xx { color: #F56C6C; }

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
