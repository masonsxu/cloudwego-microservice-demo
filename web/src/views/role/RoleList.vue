<template>
  <div class="space-y-6">
    <!-- 页头 -->
    <header class="flex items-end justify-between gap-4">
      <div>
        <h1 class="text-[22px] font-semibold leading-tight tracking-[-0.012em] text-[color:var(--color-ink-strong)]">
          {{ t('role.title') }}
        </h1>
        <p class="mt-1 text-[13px] text-[color:var(--color-ink-muted)]">
          {{ t('role.manageRoles') }}
        </p>
      </div>
      <Button @click="openCreateDialog">
        <Plus class="h-4 w-4" />
        {{ t('role.createRole') }}
      </Button>
    </header>

    <!-- 工具条 -->
    <div class="flex flex-wrap items-center gap-2">
      <div class="relative flex-1 min-w-[240px] max-w-[320px]">
        <Search class="pointer-events-none absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-[color:var(--color-ink-subtle)]" />
        <Input
          v-model="searchForm.name"
          :placeholder="t('role.searchPlaceholder')"
          class="pl-9 pr-9"
          @keyup.enter="handleSearch"
        />
        <button
          v-if="searchForm.name"
          class="absolute right-2.5 top-1/2 flex h-5 w-5 -translate-y-1/2 items-center justify-center rounded-xs text-[color:var(--color-ink-subtle)] transition-colors hover:text-[color:var(--color-ink)]"
          @click="searchForm.name = ''; handleSearch()"
        >
          <X class="h-3 w-3" />
        </button>
      </div>

      <Select v-model="searchForm.status" @update:modelValue="handleSearch">
        <SelectTrigger class="w-[140px]">
          <SelectValue :placeholder="t('common.status')" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="1">{{ t('common.active') }}</SelectItem>
          <SelectItem value="2">{{ t('common.inactive') }}</SelectItem>
          <SelectItem value="3">{{ t('common.deprecated') }}</SelectItem>
        </SelectContent>
      </Select>

      <Select v-model="searchForm.isSystemRole" @update:modelValue="handleSearch">
        <SelectTrigger class="w-[140px]">
          <SelectValue :placeholder="t('role.isSystemRole')" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="true">{{ t('common.yes') }}</SelectItem>
          <SelectItem value="false">{{ t('common.no') }}</SelectItem>
        </SelectContent>
      </Select>

      <Button variant="ghost" size="sm" @click="handleReset">
        <RotateCcw class="h-3.5 w-3.5" />
        {{ t('common.reset') }}
      </Button>
    </div>

    <!-- 表格卡片 -->
    <div class="overflow-hidden rounded-md border border-subtle bg-canvas">
      <ListPageSkeleton v-if="initialLoading" :columns="6" :rows="8" />
      <template v-else>
        <Table>
          <TableHeader>
            <TableRow class="border-subtle hover:bg-transparent">
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[220px]">{{ t('role.roleName') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)]">{{ t('role.description') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[100px]">{{ t('common.status') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[100px] text-right">{{ t('role.userCount') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[140px]">{{ t('common.createTime') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[120px] text-right">{{ t('common.actions') }}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow
              v-for="row in roleList"
              :key="row.id"
              class="h-11 border-subtle transition-colors duration-[var(--duration-fast)] hover:bg-sunken"
            >
              <TableCell>
                <div class="flex items-center gap-2">
                  <span
                    class="flex h-6 w-6 flex-shrink-0 items-center justify-center rounded-sm"
                    :class="row.is_system_role
                      ? 'bg-[color:var(--color-warning-soft)] text-[color:var(--color-warning-ink)]'
                      : 'bg-[color:var(--color-primary-soft)] text-[color:var(--color-primary-active)]'"
                  >
                    <KeyRound class="h-3 w-3" />
                  </span>
                  <span class="font-medium text-ink">{{ row.name }}</span>
                  <Badge v-if="row.is_system_role" variant="warning">{{ t('role.systemBadge') }}</Badge>
                </div>
              </TableCell>
              <TableCell class="text-[color:var(--color-ink-muted)]">
                {{ row.description || '—' }}
              </TableCell>
              <TableCell>
                <Badge :variant="getStatusVariant(row.status)">
                  {{ getStatusText(row.status) }}
                </Badge>
              </TableCell>
              <TableCell class="tabular text-right text-ink">
                {{ row.user_count || 0 }}
              </TableCell>
              <TableCell class="tabular text-[color:var(--color-ink-muted)]">
                {{ formatTimestamp(row.created_at) }}
              </TableCell>
              <TableCell class="text-right">
                <div class="inline-flex items-center gap-0.5">
                  <Button variant="ghost" size="icon" :title="t('common.view')" @click="handleView(row)">
                    <Eye class="h-3.5 w-3.5" />
                  </Button>
                  <Button variant="ghost" size="icon" :title="t('common.edit')" @click="handleEdit(row)">
                    <Pencil class="h-3.5 w-3.5" />
                  </Button>
                  <Button
                    v-if="!row.is_system_role"
                    variant="ghost"
                    size="icon"
                    :title="t('common.delete')"
                    class="hover:text-[color:var(--color-danger)]"
                    @click="handleDelete(row)"
                  >
                    <Trash2 class="h-3.5 w-3.5" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>

            <TableRow v-if="!loading && roleList.length === 0">
              <TableCell colspan="6" class="h-[200px] text-center">
                <div class="flex flex-col items-center justify-center gap-2 py-8">
                  <KeyRound class="h-8 w-8 text-[color:var(--color-ink-subtle)]" />
                  <p class="text-[14px] font-medium text-ink">{{ t('common.noData') }}</p>
                  <p class="text-[12px] text-[color:var(--color-ink-muted)]">{{ t('role.emptyHint') }}</p>
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
              <Button variant="outline" size="icon" :disabled="pagination.page <= 1" @click="pagination.page--; loadRoles()">
                <ChevronLeft class="h-4 w-4" />
              </Button>
              <span class="tabular min-w-[40px] text-center text-[13px] font-medium text-ink">{{ pagination.page }}</span>
              <Button variant="outline" size="icon" :disabled="pagination.page * pagination.size >= pagination.total" @click="pagination.page++; loadRoles()">
                <ChevronRight class="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- 创建/编辑对话框 -->
    <Dialog v-model:open="showCreateDialog">
      <DialogContent class="max-w-[480px]">
        <DialogHeader>
          <DialogTitle>{{ editingRole ? t('role.editRole') : t('role.createRole') }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 py-2">
          <div class="space-y-1.5">
            <Label>{{ t('role.roleName') }}</Label>
            <Input v-model="roleForm.name" :placeholder="t('role.roleName')" :disabled="!!editingRole" />
          </div>
          <div class="space-y-1.5">
            <Label>{{ t('role.description') }}</Label>
            <Textarea v-model="roleForm.description" :placeholder="t('role.description')" :rows="3" />
          </div>
          <div v-if="editingRole" class="space-y-1.5">
            <Label>{{ t('common.status') }}</Label>
            <Select v-model="roleFormStatusString">
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="1">{{ t('common.active') }}</SelectItem>
                <SelectItem value="2">{{ t('common.inactive') }}</SelectItem>
                <SelectItem value="3">{{ t('common.deprecated') }}</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="flex items-center gap-2 pt-1">
            <Switch v-model:checked="roleForm.is_system_role" :disabled="!!editingRole" />
            <Label class="cursor-pointer">{{ t('role.isSystemRole') }}</Label>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showCreateDialog = false">{{ t('common.cancel') }}</Button>
          <Button @click="handleSaveRole">{{ t('common.save') }}</Button>
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
import { KeyRound, Plus, Search, Eye, Pencil, Trash2, ChevronLeft, ChevronRight, X, RotateCcw } from 'lucide-vue-next'
import { roleApi } from '@/api/role'
import type { RoleDefinitionDTO } from '@/api/role'
import ListPageSkeleton from '@/components/skeleton/ListPageSkeleton.vue'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'

const { t } = useI18n()
const router = useRouter()

const initialLoading = ref(true)
const loading = ref(false)
const roleList = ref<RoleDefinitionDTO[]>([])
const showCreateDialog = ref(false)
const editingRole = ref<RoleDefinitionDTO | null>(null)

const searchForm = reactive({
  name: '',
  status: undefined as string | undefined,
  isSystemRole: undefined as string | undefined,
})

const pagination = reactive({ page: 1, size: 20, total: 0 })

const roleForm = reactive({
  name: '', description: '', status: 1, is_system_role: false, permissions: [],
})

const roleFormStatusString = computed({
  get: () => String(roleForm.status),
  set: (val: string) => { roleForm.status = Number(val) },
})

const paginationSizeString = computed({
  get: () => String(pagination.size),
  set: (val: string) => { pagination.size = Number(val); loadRoles() },
})

const loadRoles = async () => {
  loading.value = true
  try {
    const response = await roleApi.getRoles({
      name: searchForm.name || undefined,
      status: searchForm.status ? Number(searchForm.status) : undefined,
      is_system_role: searchForm.isSystemRole ? searchForm.isSystemRole === 'true' : undefined,
      page: pagination.page,
      limit: pagination.size,
    })
    roleList.value = response.roles || []
    pagination.total = response.page?.total || 0
  } catch (error: any) {
    toast.error(error?.message || t('common.operationFailed'))
  } finally {
    loading.value = false
    initialLoading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; loadRoles() }
const handleReset = () => { searchForm.name = ''; searchForm.status = undefined; searchForm.isSystemRole = undefined; pagination.page = 1; loadRoles() }
const handleView = (role: RoleDefinitionDTO) => { router.push({ name: 'RoleDetail', params: { id: role.id } }) }

const openCreateDialog = () => {
  editingRole.value = null
  roleForm.name = ''
  roleForm.description = ''
  roleForm.status = 1
  roleForm.is_system_role = false
  showCreateDialog.value = true
}

const handleEdit = (role: RoleDefinitionDTO) => {
  editingRole.value = role
  roleForm.name = role.name || ''
  roleForm.description = role.description || ''
  roleForm.status = role.status || 1
  roleForm.is_system_role = role.is_system_role || false
  showCreateDialog.value = true
}

const handleDelete = async (role: RoleDefinitionDTO) => {
  if (!confirm(`${t('role.roleName')}: ${role.name}\n${t('common.deleteConfirm')}`)) return
  try {
    await roleApi.deleteRole(role.id!)
    toast.success(t('common.deleteSuccess'))
    loadRoles()
  } catch (error: any) {
    toast.error(error?.message || t('common.operationFailed'))
  }
}

const handleSaveRole = async () => {
  try {
    const data = {
      name: roleForm.name,
      description: roleForm.description,
      is_system_role: roleForm.is_system_role,
      permissions: roleForm.permissions
    }
    let response: any
    if (editingRole.value) {
      response = await roleApi.updateRole(editingRole.value.id!, { ...data, status: roleForm.status })
    } else {
      response = await roleApi.createRole(data)
    }
    if (response.role) {
      toast.success(editingRole.value ? t('common.updateSuccess') : t('common.createSuccess'))
      showCreateDialog.value = false
      loadRoles()
    }
  } catch (error: any) {
    toast.error(error?.message || t('common.operationFailed'))
  }
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = { 1: t('common.active'), 2: t('common.inactive'), 3: t('common.deprecated') }
  return map[status] || ''
}

function getStatusVariant(status: number): 'success' | 'default' | 'destructive' {
  const variantMap: Record<number, 'success' | 'default' | 'destructive'> = {
    1: 'success', 2: 'default', 3: 'destructive'
  }
  return variantMap[status] || 'default'
}

const formatTimestamp = (ts?: number) => ts ? new Date(ts).toLocaleDateString('zh-CN') : '—'

onMounted(() => { loadRoles() })
</script>
