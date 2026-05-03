<template>
  <div class="space-y-6">
    <!-- 页头 -->
    <header>
      <h1 class="text-[22px] font-semibold leading-tight tracking-[-0.012em] text-[color:var(--color-ink-strong)]">
        {{ t('audit.title') }}
      </h1>
      <p class="mt-1 text-[13px] text-[color:var(--color-ink-muted)]">
        {{ t('audit.subtitle') }}
      </p>
    </header>

    <!-- 统计概览 -->
    <section>
      <div class="grid gap-4 sm:grid-cols-3">
        <div class="rounded-md border border-subtle bg-canvas p-4">
          <div class="flex items-center gap-2 text-[color:var(--color-ink-muted)]">
            <FileText class="h-4 w-4" />
            <span class="text-[13px] font-medium">{{ t('audit.totalRecords') }}</span>
          </div>
          <p class="mt-3 tabular text-[24px] font-semibold leading-none text-[color:var(--color-ink-strong)]">
            {{ formatNumber(stats.total_count) }}
          </p>
        </div>
        <div class="rounded-md border border-subtle bg-canvas p-4">
          <div class="flex items-center gap-2 text-[color:var(--color-ink-muted)]">
            <CircleCheck class="h-4 w-4" />
            <span class="text-[13px] font-medium">{{ t('audit.successRate') }}</span>
          </div>
          <p class="mt-3 tabular text-[24px] font-semibold leading-none" :class="successRate >= 95 ? 'text-[color:var(--color-success-ink)]' : 'text-[color:var(--color-warning-ink)]'">
            {{ successRate }}%
          </p>
        </div>
        <div class="rounded-md border border-subtle bg-canvas p-4">
          <div class="flex items-center gap-2 text-[color:var(--color-ink-muted)]">
            <Timer class="h-4 w-4" />
            <span class="text-[13px] font-medium">{{ t('audit.avgDuration') }}</span>
          </div>
          <p class="mt-3 tabular text-[24px] font-semibold leading-none text-[color:var(--color-ink-strong)]">
            {{ stats.avg_duration_ms }}<span class="ml-1 text-[14px] font-medium text-[color:var(--color-ink-subtle)]">ms</span>
          </p>
        </div>
      </div>
    </section>

    <!-- 工具条 -->
    <div class="flex flex-wrap items-center gap-2">
      <Select v-model="searchForm.action" @update:modelValue="handleSearch">
        <SelectTrigger class="w-[160px]">
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

      <Select v-model="searchForm.success" @update:modelValue="handleSearch">
        <SelectTrigger class="w-[140px]">
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

      <Button variant="ghost" size="sm" @click="handleReset">
        <RotateCcw class="h-3.5 w-3.5" />
        {{ t('common.reset') }}
      </Button>
    </div>

    <!-- 表格卡片 -->
    <div class="overflow-hidden rounded-md border border-subtle bg-canvas">
      <ListPageSkeleton v-if="initialLoading" :columns="7" :rows="10" />
      <template v-else>
        <Table>
          <TableHeader>
            <TableRow class="border-subtle hover:bg-transparent">
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[170px]">{{ t('common.createTime') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[140px]">{{ t('audit.user') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[100px]">{{ t('audit.action') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)]">{{ t('audit.resource') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[80px] text-right">{{ t('audit.statusCode') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[90px]">{{ t('common.status') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[140px]">{{ t('audit.clientIP') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[90px] text-right">{{ t('audit.duration') }}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow
              v-for="row in auditLogs"
              :key="row.id"
              class="h-11 border-subtle transition-colors duration-[var(--duration-fast)] hover:bg-sunken"
            >
              <TableCell class="font-mono tabular text-[12px] text-[color:var(--color-ink-muted)]">{{ formatTimestamp(row.created_at) }}</TableCell>
              <TableCell class="font-medium text-ink">{{ row.username || row.user_id || '—' }}</TableCell>
              <TableCell>
                <Badge :variant="getActionBadgeVariant(row.action)">
                  {{ getActionText(row.action) }}
                </Badge>
              </TableCell>
              <TableCell class="font-mono text-[12px] text-[color:var(--color-ink-muted)] truncate max-w-[260px]" :title="row.resource">
                {{ row.resource }}
              </TableCell>
              <TableCell>
                <span class="font-mono tabular text-[13px] font-semibold text-right block" :class="getStatusCodeColorClass(row.status_code)">
                  {{ row.status_code }}
                </span>
              </TableCell>
              <TableCell>
                <Badge :variant="row.success ? 'success' : 'destructive'">
                  {{ row.success ? t('audit.success') : t('audit.failed') }}
                </Badge>
              </TableCell>
              <TableCell class="font-mono text-[12px] text-[color:var(--color-ink-muted)]">{{ row.client_ip }}</TableCell>
              <TableCell class="font-mono tabular text-[12px] text-right" :class="row.duration_ms >= 1000 ? 'font-semibold text-[color:var(--color-danger-ink)]' : 'text-[color:var(--color-ink-muted)]'">
                {{ row.duration_ms }}ms
              </TableCell>
            </TableRow>

            <TableRow v-if="!loading && auditLogs.length === 0">
              <TableCell colspan="8" class="h-[200px] text-center">
                <div class="flex flex-col items-center justify-center gap-2 py-8">
                  <FileText class="h-8 w-8 text-[color:var(--color-ink-subtle)]" />
                  <p class="text-[14px] font-medium text-ink">{{ t('common.noData') }}</p>
                  <p class="text-[12px] text-[color:var(--color-ink-muted)]">{{ t('audit.emptyHint') }}</p>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>

        <!-- 分页 -->
        <div class="flex items-center justify-between border-t border-subtle px-4 py-2.5">
          <span class="tabular text-[13px] text-[color:var(--color-ink-muted)]">
            {{ t('common.total', { total: pagination.total }) }}
          </span>
          <div class="flex items-center gap-2">
            <Select v-model="paginationSizeString">
              <SelectTrigger class="h-8 w-[80px]">
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
              <Button variant="outline" size="icon" :disabled="pagination.page <= 1" @click="pagination.page = 1; loadAuditLogs()">
                <ChevronsLeft class="h-4 w-4" />
              </Button>
              <Button variant="outline" size="icon" :disabled="pagination.page <= 1" @click="pagination.page--; loadAuditLogs()">
                <ChevronLeft class="h-4 w-4" />
              </Button>
              <span class="tabular min-w-[40px] text-center text-[13px] font-medium text-ink">{{ pagination.page }}</span>
              <Button variant="outline" size="icon" :disabled="pagination.page >= Math.ceil(pagination.total / pagination.size)" @click="pagination.page++; loadAuditLogs()">
                <ChevronRight class="h-4 w-4" />
              </Button>
              <Button variant="outline" size="icon" :disabled="pagination.page >= Math.ceil(pagination.total / pagination.size)" @click="pagination.page = Math.ceil(pagination.total / pagination.size); loadAuditLogs()">
                <ChevronsRight class="h-4 w-4" />
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
import { FileText, CircleCheck, Timer, RotateCcw, ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from 'lucide-vue-next'
import { auditApi } from '@/api/audit'
import type { AuditLogDTO } from '@/api/audit'
import type { AuditLogStats } from '@/types/audit'
import ListPageSkeleton from '@/components/skeleton/ListPageSkeleton.vue'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'

const { t } = useI18n()

const initialLoading = ref(true)
const loading = ref(false)
const auditLogs = ref<AuditLogDTO[]>([])
const stats = ref<AuditLogStats>({ total_count: 0, success_count: 0, avg_duration_ms: 0 })

const searchForm = reactive({
  action: 'all',
  success: 'all',
})

const pagination = reactive({ page: 1, size: 20, total: 0 })

const paginationSizeString = computed({
  get: () => String(pagination.size),
  set: (val: string) => { pagination.size = Number(val); loadAuditLogs() },
})

const successRate = computed(() => {
  if (stats.value.total_count === 0) return 0
  return Math.round((stats.value.success_count / stats.value.total_count) * 100)
})

function formatNumber(value: number): string {
  return new Intl.NumberFormat('en-US').format(value)
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

function getActionBadgeVariant(action: number): 'success' | 'info' | 'destructive' | 'default' | 'warning' {
  const map: Record<number, 'success' | 'info' | 'destructive' | 'default' | 'warning'> = {
    1: 'success',
    2: 'info',
    3: 'destructive',
    4: 'default',
    5: 'default',
    6: 'warning',
  }
  return map[action] || 'default'
}

function getStatusCodeColorClass(code: number): string {
  if (code >= 200 && code < 300) return 'text-[color:var(--color-success-ink)]'
  if (code >= 400 && code < 500) return 'text-[color:var(--color-warning-ink)]'
  if (code >= 500) return 'text-[color:var(--color-danger-ink)]'
  return 'text-[color:var(--color-ink-muted)]'
}

const formatTimestamp = (ts?: number) => {
  if (!ts) return '—'
  const d = new Date(ts)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

onMounted(() => { loadAuditLogs() })
</script>
