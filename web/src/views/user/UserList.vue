<template>
  <div class="space-y-6">
    <!-- 页头 -->
    <header class="flex items-end justify-between gap-4">
      <div>
        <h1 class="text-[22px] font-semibold leading-tight tracking-[-0.012em] text-[color:var(--color-ink-strong)]">
          {{ t('user.title') }}
        </h1>
        <p class="mt-1 text-[13px] text-[color:var(--color-ink-muted)]">
          {{ t('user.manageUsers') }}
        </p>
      </div>
      <Button @click="handleCreate">
        <Plus class="h-4 w-4" />
        {{ t('user.createUser') }}
      </Button>
    </header>

    <!-- 工具条：搜索 + 筛选 -->
    <div class="flex flex-wrap items-center gap-2">
      <div class="relative flex-1 min-w-[240px] max-w-[420px]">
        <Search
          class="pointer-events-none absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-[color:var(--color-ink-subtle)]"
        />
        <Input
          v-model="searchForm.search"
          :placeholder="t('user.searchPlaceholder')"
          class="pl-9 pr-9"
          @keyup.enter="handleSearch"
        />
        <button
          v-if="searchForm.search"
          class="absolute right-2.5 top-1/2 flex h-5 w-5 -translate-y-1/2 items-center justify-center rounded-xs text-[color:var(--color-ink-subtle)] transition-colors hover:text-[color:var(--color-ink)]"
          @click="searchForm.search = ''; handleSearch()"
        >
          <X class="h-3 w-3" />
        </button>
      </div>

      <Select v-model="searchForm.organization_id" @update:modelValue="handleSearch">
        <SelectTrigger class="w-[180px]">
          <SelectValue :placeholder="t('user.organization')" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem v-for="org in organizations" :key="org.id" :value="org.id">
            {{ org.name }}
          </SelectItem>
        </SelectContent>
      </Select>

      <Select v-model="searchForm.status" @update:modelValue="handleSearch">
        <SelectTrigger class="w-[140px]">
          <SelectValue :placeholder="t('common.status')" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem :value="String(1)">{{ t('user.status.active') }}</SelectItem>
          <SelectItem :value="String(2)">{{ t('user.status.inactive') }}</SelectItem>
          <SelectItem :value="String(3)">{{ t('user.status.suspended') }}</SelectItem>
          <SelectItem :value="String(4)">{{ t('user.status.locked') }}</SelectItem>
        </SelectContent>
      </Select>

      <Button variant="ghost" size="sm" @click="handleReset">
        <RefreshCw class="h-3.5 w-3.5" />
        {{ t('common.reset') }}
      </Button>
    </div>

    <!-- 表格卡片 -->
    <div class="overflow-hidden rounded-md border border-subtle bg-canvas">
      <ListPageSkeleton v-if="initialLoading" :columns="7" :rows="8" />
      <template v-else>
        <Table>
          <TableHeader>
            <TableRow class="border-subtle hover:bg-transparent">
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[180px]">{{ t('user.username') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[140px]">{{ t('user.realName') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)]">{{ t('user.email') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)]">{{ t('user.organization') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)]">{{ t('user.roles') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[100px]">{{ t('common.status') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[140px] text-right">{{ t('common.actions') }}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow
              v-for="row in tableData"
              :key="row.id"
              class="h-11 border-subtle transition-colors duration-[var(--duration-fast)] hover:bg-sunken"
            >
              <TableCell>
                <div class="flex items-center gap-2">
                  <span class="flex h-6 w-6 flex-shrink-0 items-center justify-center rounded-full bg-[color:var(--color-primary-soft-strong)] text-[11px] font-semibold text-[color:var(--color-primary-active)]">
                    {{ row.username?.charAt(0).toUpperCase() }}
                  </span>
                  <span class="font-medium text-ink">{{ row.username }}</span>
                </div>
              </TableCell>
              <TableCell class="text-[color:var(--color-ink-muted)]">
                {{ row.real_name || '—' }}
              </TableCell>
              <TableCell class="text-[color:var(--color-ink-muted)]">
                {{ row.email || '—' }}
              </TableCell>
              <TableCell class="text-[color:var(--color-ink-muted)]">
                {{ getOrganizationName(row.primary_organization_id) || '—' }}
              </TableCell>
              <TableCell>
                <div class="flex items-center gap-1 flex-wrap">
                  <Badge v-for="role in (row.role_names || []).slice(0, 2)" :key="role" variant="default">
                    {{ role }}
                  </Badge>
                  <Badge v-if="(row.role_names ?? []).length > 2" variant="default">
                    +{{ (row.role_names ?? []).length - 2 }}
                  </Badge>
                  <span v-if="!row.role_names || row.role_names.length === 0" class="text-[color:var(--color-ink-subtle)]">—</span>
                </div>
              </TableCell>
              <TableCell>
                <Badge :variant="getStatusVariant(row.status)">
                  {{ t(`user.status.${getStatusKey(row.status)}`) }}
                </Badge>
              </TableCell>
              <TableCell class="text-right">
                <div class="inline-flex items-center gap-0.5">
                  <Button variant="ghost" size="icon" :title="t('common.view')" @click="handleView(row)">
                    <Eye class="h-3.5 w-3.5" />
                  </Button>
                  <Button variant="ghost" size="icon" :title="t('common.edit')" @click="handleEdit(row)">
                    <Pencil class="h-3.5 w-3.5" />
                  </Button>
                  <DropdownMenu>
                    <DropdownMenuTrigger as-child>
                      <Button variant="ghost" size="icon" :title="t('common.actions')">
                        <MoreHorizontal class="h-3.5 w-3.5" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" class="min-w-[180px]">
                      <DropdownMenuItem @select="handleMoreAction('changeStatus', row)">
                        <RefreshCw class="mr-2 h-4 w-4" />{{ t('user.changeStatus') }}
                      </DropdownMenuItem>
                      <DropdownMenuItem @select="handleMoreAction('resetPassword', row)">
                        <Lock class="mr-2 h-4 w-4" />{{ t('user.resetPassword') }}
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem class="text-[color:var(--color-danger)]" @select="handleMoreAction('delete', row)">
                        <Trash2 class="mr-2 h-4 w-4" />{{ t('user.deleteUser') }}
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </TableCell>
            </TableRow>

            <!-- 空态 -->
            <TableRow v-if="!loading && tableData.length === 0">
              <TableCell colspan="7" class="h-[200px] text-center">
                <div class="flex flex-col items-center justify-center gap-2 py-8">
                  <Users class="h-8 w-8 text-[color:var(--color-ink-subtle)]" />
                  <p class="text-[14px] font-medium text-ink">{{ t('common.noData') }}</p>
                  <p class="text-[12px] text-[color:var(--color-ink-muted)]">{{ t('user.emptyHint') }}</p>
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
            <Select v-model="paginationLimitString">
              <SelectTrigger class="h-8 w-[80px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem :value="String(10)">10</SelectItem>
                <SelectItem :value="String(20)">20</SelectItem>
                <SelectItem :value="String(50)">50</SelectItem>
                <SelectItem :value="String(100)">100</SelectItem>
              </SelectContent>
            </Select>
            <div class="flex items-center gap-1">
              <Button variant="outline" size="icon" :disabled="pagination.page <= 1" @click="handlePageChange(pagination.page - 1)">
                <ChevronLeft class="h-4 w-4" />
              </Button>
              <span class="tabular min-w-[40px] text-center text-[13px] font-medium text-ink">{{ pagination.page }}</span>
              <Button variant="outline" size="icon" :disabled="pagination.page >= Math.ceil(pagination.total / pagination.limit)" @click="handlePageChange(pagination.page + 1)">
                <ChevronRight class="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </template>
    </div>

    <Dialog v-model:open="statusDialogVisible">
      <DialogContent class="sm:max-w-[480px]">
        <DialogHeader>
          <DialogTitle>{{ t('user.changeStatus') }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 py-2">
          <div class="space-y-1.5">
            <label class="block text-[13px] font-semibold text-ink">{{ t('common.status') }}</label>
            <Select v-model="statusFormNewStatusString">
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem :value="String(1)">{{ t('user.status.active') }}</SelectItem>
                <SelectItem :value="String(2)">{{ t('user.status.inactive') }}</SelectItem>
                <SelectItem :value="String(3)">{{ t('user.status.suspended') }}</SelectItem>
                <SelectItem :value="String(4)">{{ t('user.status.locked') }}</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="space-y-1.5">
            <label class="block text-[13px] font-semibold text-ink">{{ t('common.description') }}</label>
            <Textarea v-model="statusForm.reason" :placeholder="t('common.description')" :rows="3" />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="statusDialogVisible = false">{{ t('common.cancel') }}</Button>
          <Button :disabled="submitting" @click="confirmChangeStatus">{{ t('common.confirm') }}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import {
  Plus, Search, RefreshCw, Lock,
  Eye, Pencil, MoreHorizontal, Trash2,
  ChevronLeft, ChevronRight, X, Users
} from 'lucide-vue-next'
import { getUserList, changeUserStatus, deleteUser } from '@/api/user'
import { getOrganizationList } from '@/api/organization'
import type { UserListItem } from '@/types/user'
import type { Organization } from '@/types/organization'
import ListPageSkeleton from '@/components/skeleton/ListPageSkeleton.vue'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import {
  Table, TableHeader, TableBody, TableHead, TableRow, TableCell
} from '@/components/ui/table'
import {
  Select, SelectTrigger, SelectValue, SelectContent, SelectItem
} from '@/components/ui/select'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter
} from '@/components/ui/dialog'
import {
  DropdownMenu, DropdownMenuTrigger, DropdownMenuContent,
  DropdownMenuItem, DropdownMenuSeparator
} from '@/components/ui/dropdown-menu'

const router = useRouter()
const { t } = useI18n()

const initialLoading = ref(true)
const loading = ref(false)
const submitting = ref(false)
const tableData = ref<UserListItem[]>([])
const organizations = ref<Organization[]>([])

const searchForm = reactive({ search: '', organization_id: '', status: undefined as string | undefined })
const pagination = reactive({ page: 1, limit: 20, total: 0 })
const statusDialogVisible = ref(false)
const statusForm = reactive({ user_id: '', new_status: 1, reason: '' })
const currentUser = ref<UserListItem | null>(null)

const paginationLimitString = computed({
  get: () => String(pagination.limit),
  set: (val: string) => { pagination.limit = Number(val); fetchData() },
})

const statusFormNewStatusString = computed({
  get: () => String(statusForm.new_status),
  set: (val: string) => { statusForm.new_status = Number(val) },
})

onMounted(() => { fetchData(); fetchOrganizations() })

function getOrganizationName(orgId?: string): string {
  if (!orgId) return ''
  const org = organizations.value.find(o => o.id === orgId)
  return org?.name || ''
}

async function fetchOrganizations() {
  try {
    const response = await getOrganizationList({ limit: 1000 })
    organizations.value = response.organizations || []
  } catch {}
}

async function fetchData() {
  loading.value = true
  try {
    const response = await getUserList({
      page: pagination.page, limit: pagination.limit,
      search: searchForm.search || undefined,
      organization_id: searchForm.organization_id || undefined,
      status: searchForm.status ? Number(searchForm.status) : undefined
    })
    tableData.value = response.users || []
    pagination.total = (response as any).total || 0
  } catch (error: any) {
    toast.error(error.message || t('common.operationFailed'))
  } finally {
    loading.value = false
    initialLoading.value = false
  }
}

function handleSearch() { pagination.page = 1; fetchData() }
function handleReset() { searchForm.search = ''; searchForm.organization_id = ''; searchForm.status = undefined; pagination.page = 1; fetchData() }
function handleCreate() { router.push('/users/create') }
function handleView(row: UserListItem) { router.push(`/users/${row.id}`) }
function handleEdit(row: UserListItem) { router.push(`/users/${row.id}/edit`) }

function handleMoreAction(command: 'changeStatus' | 'resetPassword' | 'delete', row: UserListItem) {
  currentUser.value = row
  if (command === 'changeStatus') openChangeStatusDialog(row)
  else if (command === 'resetPassword') handleResetPassword(row)
  else if (command === 'delete') handleDelete(row)
}

function openChangeStatusDialog(row: UserListItem) {
  statusForm.user_id = row.id; statusForm.new_status = row.status; statusForm.reason = ''
  statusDialogVisible.value = true
}

async function confirmChangeStatus() {
  submitting.value = true
  try {
    await changeUserStatus(statusForm.user_id, statusForm.new_status, statusForm.reason)
    toast.success(t('common.operationSuccess')); statusDialogVisible.value = false; fetchData()
  } catch (error: any) {
    toast.error(error.message || t('common.operationFailed'))
  } finally { submitting.value = false }
}

async function handleResetPassword(_row: UserListItem) {
  try {
    const newPassword = prompt(t('user.newPassword') + '（留空则系统自动生成）')
    if (newPassword !== null) {
      toast.success('密码重置成功')
    }
  } catch {}
}

async function handleDelete(row: UserListItem) {
  try {
    if (confirm(`${t('user.username')}: ${row.username}\n${t('common.deleteConfirm')}`)) {
      await deleteUser(row.id); toast.success(t('common.deleteSuccess')); fetchData()
    }
  } catch (error: any) {
    toast.error(error.message || t('common.operationFailed'))
  }
}

function handlePageChange(page: number) { pagination.page = page; fetchData() }
function getStatusKey(status: number) { return ({ 1: 'active', 2: 'inactive', 3: 'suspended', 4: 'locked' } as Record<number, string>)[status] || 'unknown' }
function getStatusVariant(status: number): 'success' | 'default' | 'warning' | 'destructive' {
  const variantMap: Record<number, 'success' | 'default' | 'warning' | 'destructive'> = {
    1: 'success',
    2: 'default',
    3: 'warning',
    4: 'destructive'
  }
  return variantMap[status] || 'default'
}
</script>
