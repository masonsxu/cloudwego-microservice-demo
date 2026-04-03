<template>
  <div class="user-list">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('user.title') }}</h1>
        <p class="page-subtitle">{{ t('user.manageUsers') || '管理系统用户账号与访问权限' }}</p>
      </div>
      <button class="create-btn shimmer-btn" @click="handleCreate">
        <Plus class="h-4 w-4" />
        {{ t('user.createUser') }}
      </button>
    </div>

    <div class="stats-row">
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon"><User class="h-5 w-5" /></div>
        <div class="stat-info">
          <span class="stat-value">{{ pagination.total }}</span>
          <span class="stat-label">{{ t('user.title') }}</span>
        </div>
      </div>
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon active-icon"><CircleCheck class="h-5 w-5" /></div>
        <div class="stat-info">
          <span class="stat-value">{{ tableData.filter(u => u.status === 1).length }}</span>
          <span class="stat-label">{{ t('user.status.active') }}</span>
        </div>
      </div>
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon locked-icon"><Lock class="h-5 w-5" /></div>
        <div class="stat-info">
          <span class="stat-value">{{ tableData.filter(u => u.status === 4).length }}</span>
          <span class="stat-label">{{ t('user.status.locked') }}</span>
        </div>
      </div>
    </div>

    <div class="search-section">
      <div class="search-form">
        <div class="search-input-wrapper">
          <Search class="search-icon" />
          <input
            v-model="searchForm.search"
            :placeholder="t('user.username') + ' / ' + t('user.email') + ' / ' + t('user.realName')"
            class="search-input"
            @keyup.enter="handleSearch"
          />
          <button v-if="searchForm.search" class="clear-btn" @click="searchForm.search = ''; handleSearch()">
            <X class="h-3 w-3" />
          </button>
        </div>
        <div class="select-wrapper">
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
        </div>
        <div class="select-wrapper">
          <Select v-model="searchForm.status" @update:modelValue="handleSearch">
            <SelectTrigger class="w-[130px]">
              <SelectValue :placeholder="t('common.status')" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem :value="String(1)">{{ t('user.status.active') }}</SelectItem>
              <SelectItem :value="String(2)">{{ t('user.status.inactive') }}</SelectItem>
              <SelectItem :value="String(3)">{{ t('user.status.suspended') }}</SelectItem>
              <SelectItem :value="String(4)">{{ t('user.status.locked') }}</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div class="search-buttons">
          <Button variant="outline" @click="handleReset">
            <Refresh class="mr-1 h-4 w-4" />{{ t('common.reset') }}
          </Button>
          <Button @click="handleSearch">
            <Search class="mr-1 h-4 w-4" />{{ t('common.search') }}
          </Button>
        </div>
      </div>
    </div>

    <div class="table-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
      <ListPageSkeleton v-if="initialLoading" :columns="7" :rows="8" />
      <template v-else>
        <div class="table-body">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead class="w-[150px]">{{ t('user.username') }}</TableHead>
                <TableHead class="w-[120px]">{{ t('user.realName') }}</TableHead>
                <TableHead>{{ t('user.email') }}</TableHead>
                <TableHead>{{ t('user.organization') }}</TableHead>
                <TableHead>{{ t('user.roles') }}</TableHead>
                <TableHead class="w-[110px] text-center">{{ t('common.status') }}</TableHead>
                <TableHead class="w-[160px] text-right">{{ t('common.actions') }}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="row in tableData" :key="row.id">
                <TableCell>
                  <div class="user-cell">
                    <div class="user-avatar-mini">{{ row.username?.charAt(0).toUpperCase() }}</div>
                    <span class="username">{{ row.username }}</span>
                  </div>
                </TableCell>
                <TableCell>
                  <span class="text-sub">{{ row.real_name || '-' }}</span>
                </TableCell>
                <TableCell>
                  <span class="text-sub">{{ row.email || '-' }}</span>
                </TableCell>
                <TableCell>
                  <span v-if="getOrganizationName(row.primary_organization_id)" class="text-sub">{{ getOrganizationName(row.primary_organization_id) }}</span>
                  <span v-else class="text-muted">-</span>
                </TableCell>
                <TableCell>
                  <div class="roles-cell">
                    <Badge v-for="role in (row.role_names || []).slice(0, 2)" :key="role" variant="secondary" class="role-tag">{{ role }}</Badge>
                    <span v-if="(row.role_names ?? []).length > 2" class="more-tag">+{{ (row.role_names ?? []).length - 2 }}</span>
                    <span v-if="!row.role_names || row.role_names.length === 0" class="text-muted">-</span>
                  </div>
                </TableCell>
                <TableCell class="text-center">
                  <Badge :class="getStatusBadgeClass(row.status)">
                    {{ t(`user.status.${getStatusKey(row.status)}`) }}
                  </Badge>
                </TableCell>
                <TableCell>
                  <div class="action-group">
                    <button class="action-btn view-btn" @click="handleView(row)" :title="t('common.view')">
                      <Eye class="h-4 w-4" />
                    </button>
                    <button class="action-btn edit-btn" @click="handleEdit(row)" :title="t('common.edit')">
                      <Pencil class="h-4 w-4" />
                    </button>
                    <DropdownMenu>
                      <DropdownMenuTrigger as-child>
                        <button class="action-btn more-btn" :title="t('common.actions')">
                          <MoreHorizontal class="h-4 w-4" />
                        </button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem @select="handleMoreAction('changeStatus', row)">
                          <RefreshCw class="mr-2 h-4 w-4" />{{ t('user.changeStatus') }}
                        </DropdownMenuItem>
                        <DropdownMenuItem @select="handleMoreAction('resetPassword', row)">
                          <Lock class="mr-2 h-4 w-4" />{{ t('user.resetPassword') }}
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem class="text-destructive" @select="handleMoreAction('delete', row)">
                          <Trash2 class="mr-2 h-4 w-4" />{{ t('user.deleteUser') }}
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>

        <div class="pagination-wrapper">
          <div class="pagination-info">
            <span class="text-sm text-muted-foreground">
              {{ t('common.total') || '共' }} {{ pagination.total }} {{ t('common.items') || '条' }}
            </span>
          </div>
          <div class="pagination-controls">
            <Select v-model="paginationLimitString">
              <SelectTrigger class="w-[80px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem :value="String(10)">10</SelectItem>
                <SelectItem :value="String(20)">20</SelectItem>
                <SelectItem :value="String(50)">50</SelectItem>
                <SelectItem :value="String(100)">100</SelectItem>
              </SelectContent>
            </Select>
            <div class="page-buttons">
              <Button variant="outline" size="sm" :disabled="pagination.page <= 1" @click="handlePageChange(pagination.page - 1)">
                <ChevronLeft class="h-4 w-4" />
              </Button>
              <span class="page-info">{{ pagination.page }}</span>
              <Button variant="outline" size="sm" :disabled="pagination.page >= Math.ceil(pagination.total / pagination.limit)" @click="handlePageChange(pagination.page + 1)">
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
        <div class="grid gap-4 py-4">
          <div class="grid grid-cols-4 items-center gap-4">
            <label class="text-right text-sm">{{ t('common.status') }}</label>
            <Select v-model="statusFormNewStatusString" class="col-span-3">
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
          <div class="grid grid-cols-4 items-center gap-4">
            <label class="text-right text-sm">{{ t('common.descriptions') }}</label>
            <Textarea v-model="statusForm.reason" :placeholder="t('common.descriptions')" class="col-span-3" :rows="3" />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="statusDialogVisible = false">{{ t('common.cancel') }}</Button>
          <Button @click="confirmChangeStatus" :disabled="submitting">{{ t('common.confirm') }}</Button>
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
  Plus, Search, RefreshCw, User, CircleCheck, Lock,
  Eye, Pencil, MoreHorizontal, Trash2,
  ChevronLeft, ChevronRight, X
} from 'lucide-vue-next'
import { getUserList, changeUserStatus, deleteUser } from '@/api/user'
import { getOrganizationList } from '@/api/organization'
import type { UserListItem } from '@/types/user'
import type { Organization } from '@/types/organization'
import ListPageSkeleton from '@/components/skeleton/ListPageSkeleton.vue'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
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
  set: (val: string) => { pagination.limit = Number(val) },
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
function getStatusBadgeClass(status: number) {
  const classMap: Record<number, string> = {
    1: 'bg-emerald-500 text-white',
    2: 'bg-gray-400 text-white',
    3: 'bg-amber-500 text-white',
    4: 'bg-red-500 text-white'
  }
  return classMap[status] || 'bg-gray-400 text-white'
}
</script>

<style scoped lang="scss">
.user-list {
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
        font-family: 'Inter', sans-serif;
        background: linear-gradient(120deg, var(--c-primary), var(--c-accent));
        background-clip: text;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        margin: 0 0 6px;
        line-height: 1.2;
      }
      .page-subtitle { color: var(--c-text-sub); font-size: 13px; margin: 0; }
    }

    .create-btn {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 10px 20px;
      background: linear-gradient(135deg, var(--c-primary) 0%, var(--c-accent) 100%);
      color: #fff;
      font-weight: 600;
      font-size: 14px;
      border: none;
      border-radius: 10px;
      cursor: pointer;
      transition: all 0.3s ease;

      &:hover { box-shadow: 0 0 24px rgba(63, 81, 181, 0.35); transform: translateY(-1px); }
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
      border: 1px solid hsl(var(--border) / 0.6);
      border-radius: 14px;
      box-shadow: var(--shadow-card);

      .stat-icon {
        width: 44px;
        height: 44px;
        border-radius: 10px;
        background: rgba(63, 81, 181, 0.12);
        border: 1px solid rgba(63, 81, 181, 0.2);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--c-primary);
        flex-shrink: 0;

        &.active-icon { background: rgba(67, 160, 71, 0.12); border-color: rgba(67, 160, 71, 0.25); color: #43a047; }
        &.locked-icon { background: rgba(239, 83, 80, 0.12); border-color: rgba(239, 83, 80, 0.25); color: #ef5350; }
      }

      .stat-info {
        display: flex;
        flex-direction: column;
        gap: 2px;
        .stat-value { font-size: 24px; font-weight: 700; color: var(--c-text-main); font-family: 'Source Code Pro', monospace; line-height: 1; }
        .stat-label { font-size: 12px; color: var(--c-text-sub); }
      }
    }
  }

  .search-section {
    background: var(--bg-card);
    border: 1px solid hsl(var(--border) / 0.6);
    border-radius: 14px;
    padding: 16px 20px;
    margin-bottom: 20px;

    .search-form {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: 12px;
    }

    .search-input-wrapper {
      position: relative;
      display: flex;
      align-items: center;
      width: 320px;

      .search-icon {
        position: absolute;
        left: 10px;
        height: 16px;
        width: 16px;
        color: var(--c-text-muted);
      }

      .search-input {
        width: 100%;
        height: 36px;
        padding: 0 32px 0 34px;
        border: 1px solid hsl(var(--border));
        border-radius: 8px;
        background: var(--bg-input);
        color: var(--c-text-main);
        font-size: 14px;

        &:focus {
          outline: none;
          border-color: var(--c-primary);
          box-shadow: 0 0 0 2px rgba(63, 81, 181, 0.15);
        }

        &::placeholder {
          color: var(--c-text-muted);
        }
      }

      .clear-btn {
        position: absolute;
        right: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 20px;
        height: 20px;
        border: none;
        background: transparent;
        color: var(--c-text-muted);
        cursor: pointer;
        border-radius: 4px;

        &:hover {
          background: var(--bg-card);
          color: var(--c-text-main);
        }
      }
    }

    .select-wrapper {
      :deep(.select-trigger) {
        height: 36px;
      }
    }

    .search-buttons {
      display: flex;
      gap: 8px;
    }
  }

  .table-card {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
    background: var(--bg-card);
    border: 1px solid hsl(var(--border) / 0.6);
    border-radius: 14px;
    overflow: hidden;
    box-shadow: var(--shadow-card);

    .table-body {
      flex: 1;
      min-height: 0;
      overflow: auto;
    }

    .user-cell {
      display: flex;
      align-items: center;
      gap: 10px;

      .user-avatar-mini {
        width: 30px;
        height: 30px;
        border-radius: 8px;
        background: linear-gradient(135deg, var(--c-primary) 0%, var(--c-accent) 100%);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 13px;
        font-weight: 700;
        color: #fff;
        flex-shrink: 0;
      }

      .username { font-weight: 500; color: var(--c-text-main); }
    }

    .text-sub { color: var(--c-text-sub); font-size: 13px; }
    .text-muted { color: var(--c-text-muted); font-size: 13px; }

    .roles-cell {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: 4px;

      .role-tag {
        font-size: 11px;
      }

      .more-tag {
        font-size: 11px;
        color: var(--c-text-sub);
        padding: 1px 6px;
        background: var(--bg-input);
        border-radius: 9999px;
      }
    }

    .action-group {
      display: flex;
      align-items: center;
      justify-content: flex-end;
      gap: 6px;

      .action-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 30px;
        height: 30px;
        border-radius: 7px;
        border: 1px solid var(--c-border);
        background: var(--bg-input);
        cursor: pointer;
        transition: all 0.2s ease;
        color: var(--c-text-sub);

        &:hover { transform: translateY(-1px); }
        &.view-btn:hover { color: var(--c-primary); border-color: rgba(63, 81, 181, 0.3); background: rgba(63, 81, 181, 0.08); }
        &.edit-btn:hover { color: var(--c-accent-dark); border-color: rgba(255, 152, 0, 0.4); background: rgba(255, 152, 0, 0.12); }
        &.more-btn:hover { color: var(--c-text-main); border-color: rgba(63, 81, 181, 0.3); background: var(--bg-input); }
      }
    }

    .pagination-wrapper {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 16px 20px;
      border-top: 1px solid var(--c-border);

      .pagination-info {
        display: flex;
        align-items: center;
      }

      .pagination-controls {
        display: flex;
        align-items: center;
        gap: 12px;

        .page-buttons {
          display: flex;
          align-items: center;
          gap: 8px;

          .page-info {
            font-size: 14px;
            font-weight: 500;
            color: var(--c-text-main);
            min-width: 32px;
            text-align: center;
          }
        }
      }
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
    background: radial-gradient(520px circle at var(--mouse-x, -100%) var(--mouse-y, -100%),
      rgba(63, 81, 181, 0.12), transparent 55%);
    pointer-events: none;
    z-index: 1;
  }
}
</style>
