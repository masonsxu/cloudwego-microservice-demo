<template>
  <div class="flex flex-col h-[calc(100vh-108px)]">
    <!-- Page Header -->
    <div class="flex justify-between items-end mb-7">
      <div>
        <h1 class="text-[26px] font-bold font-[Inter] bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent mb-1.5 leading-tight">{{ t('role.title') }}</h1>
        <p class="text-[13px] text-muted-foreground m-0">配置系统角色与权限策略</p>
      </div>
      <button class="flex items-center gap-2 px-5 py-2.5 bg-gradient-to-r from-primary to-accent text-white font-semibold text-sm border-none rounded-[10px] cursor-pointer transition-all duration-300 hover:shadow-[0_0_24px_rgba(63,81,181,0.35)] hover:-translate-y-0.5 shimmer-btn" @click="showCreateDialog = true">
        <Plus class="w-4 h-4" />
        {{ t('role.createRole') }}
      </button>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-3 gap-4 mb-6">
      <div class="spotlight-card flex items-center gap-4 p-5 bg-card border border-border/60 rounded-[14px] shadow-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="w-11 h-11 rounded-[10px] bg-primary/12 border border-primary/20 flex items-center justify-center text-primary flex-shrink-0">
          <KeyRound class="w-5 h-5" />
        </div>
        <div class="flex flex-col gap-0.5">
          <span class="text-2xl font-bold text-foreground font-[Source_Code_Pro] leading-none">{{ pagination.total }}</span>
          <span class="text-xs text-muted-foreground">{{ t('role.title') }}</span>
        </div>
      </div>
      <div class="spotlight-card flex items-center gap-4 p-5 bg-card border border-border/60 rounded-[14px] shadow-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="w-11 h-11 rounded-[10px] bg-red-500/12 border border-red-500/25 flex items-center justify-center text-red-500 flex-shrink-0">
          <Lock class="w-5 h-5" />
        </div>
        <div class="flex flex-col gap-0.5">
          <span class="text-2xl font-bold text-foreground font-[Source_Code_Pro] leading-none">{{ roleList.filter(r => r.is_system_role).length }}</span>
          <span class="text-xs text-muted-foreground">{{ t('role.isSystemRole') }}</span>
        </div>
      </div>
      <div class="spotlight-card flex items-center gap-4 p-5 bg-card border border-border/60 rounded-[14px] shadow-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="w-11 h-11 rounded-[10px] bg-green-500/12 border border-green-500/25 flex items-center justify-center text-green-600 flex-shrink-0">
          <User class="w-5 h-5" />
        </div>
        <div class="flex flex-col gap-0.5">
          <span class="text-2xl font-bold text-foreground font-[Source_Code_Pro] leading-none">{{ roleList.reduce((s, r) => s + (r.user_count || 0), 0) }}</span>
          <span class="text-xs text-muted-foreground">已分配用户</span>
        </div>
      </div>
    </div>

    <!-- Search Bar -->
    <div class="bg-card border border-border/60 rounded-[14px] px-5 py-4 mb-5">
      <div class="flex items-center flex-wrap gap-3">
        <div class="relative w-[240px]">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input
            v-model="searchForm.name"
            :placeholder="t('role.roleName')"
            class="pl-9"
          />
          <button v-if="searchForm.name" @click="searchForm.name = ''; handleSearch()" class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground">
            <X class="w-4 h-4" />
          </button>
        </div>
        <div class="w-[140px]">
          <Select v-model="searchForm.status">
            <SelectTrigger>
              <SelectValue :placeholder="t('common.status')" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1">{{ t('common.active') }}</SelectItem>
              <SelectItem value="2">{{ t('common.inactive') }}</SelectItem>
              <SelectItem value="3">{{ t('common.deprecated') }}</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div class="w-[140px]">
          <Select v-model="searchForm.isSystemRole">
            <SelectTrigger>
              <SelectValue :placeholder="t('role.isSystemRole')" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="true">{{ t('common.yes') }}</SelectItem>
              <SelectItem value="false">{{ t('common.no') }}</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div class="flex gap-2">
          <Button @click="handleSearch">
            <Search class="w-4 h-4" />{{ t('common.search') }}
          </Button>
          <Button variant="outline" @click="handleReset">
            <RotateCcw class="w-4 h-4" />{{ t('common.reset') }}
          </Button>
        </div>
      </div>
    </div>

    <!-- Role List Table -->
    <div class="spotlight-card flex-1 min-h-0 flex flex-col bg-card border border-border/60 rounded-[14px] overflow-hidden shadow-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
      <ListPageSkeleton v-if="initialLoading" :columns="6" :rows="8" />
      <template v-else>
        <div class="flex-1 min-h-0 overflow-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead class="uppercase text-xs tracking-wider">{{ t('role.roleName') }}</TableHead>
                <TableHead class="uppercase text-xs tracking-wider">{{ t('role.description') }}</TableHead>
                <TableHead class="uppercase text-xs tracking-wider w-[110px] text-center">{{ t('common.status') }}</TableHead>
                <TableHead class="uppercase text-xs tracking-wider w-[100px] text-center">{{ t('role.userCount') }}</TableHead>
                <TableHead class="uppercase text-xs tracking-wider w-[160px]">{{ t('common.createTime') }}</TableHead>
                <TableHead class="uppercase text-xs tracking-wider w-[160px] text-right">{{ t('common.actions') }}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="row in roleList" :key="row.id">
                <TableCell>
                  <div class="flex items-center gap-2.5">
                    <div class="w-7 h-7 rounded-md bg-primary/10 flex items-center justify-center text-primary flex-shrink-0" :class="{ 'bg-red-500/12 text-red-500': row.is_system_role }">
                      <KeyRound class="w-3.5 h-3.5" />
                    </div>
                    <span class="font-medium text-foreground">{{ row.name }}</span>
                    <Badge v-if="row.is_system_role" variant="outline" class="bg-red-500/12 border-red-500/30 text-red-500 text-[10px] px-1.5 py-0">系统</Badge>
                  </div>
                </TableCell>
                <TableCell>
                  <span class="text-muted-foreground text-[13px]">{{ row.description || '-' }}</span>
                </TableCell>
                <TableCell class="text-center">
                  <Badge :variant="row.status === 1 ? 'outline' : row.status === 2 ? 'secondary' : 'destructive'" class="text-xs" :class="row.status === 1 ? 'bg-green-500/10 border-green-500/30 text-green-500' : ''">
                    {{ getStatusText(row.status) }}
                  </Badge>
                </TableCell>
                <TableCell class="text-center">
                  <span class="inline-block px-2.5 py-0.5 bg-primary/8 rounded-full text-sm font-semibold text-foreground">{{ row.user_count || 0 }}</span>
                </TableCell>
                <TableCell>
                  <span class="text-muted-foreground text-xs">{{ formatTimestamp(row.created_at) }}</span>
                </TableCell>
                <TableCell class="text-right">
                  <div class="flex items-center justify-end gap-1.5">
                    <button class="flex items-center justify-center w-[30px] h-[30px] rounded-lg border border-border bg-input cursor-pointer transition-all duration-200 text-muted-foreground hover:text-primary hover:border-primary/30 hover:bg-primary/8 hover:-translate-y-0.5" @click="handleView(row)" :title="t('common.view')">
                      <Eye class="w-3.5 h-3.5" />
                    </button>
                    <button class="flex items-center justify-center w-[30px] h-[30px] rounded-lg border border-border bg-input cursor-pointer transition-all duration-200 text-muted-foreground hover:text-amber-600 hover:border-amber-500/40 hover:bg-amber-500/12 hover:-translate-y-0.5" @click="handleEdit(row)" :title="t('common.edit')">
                      <Pencil class="w-3.5 h-3.5" />
                    </button>
                    <button v-if="!row.is_system_role" class="flex items-center justify-center w-[30px] h-[30px] rounded-lg border border-border bg-input cursor-pointer transition-all duration-200 text-muted-foreground hover:text-red-500 hover:border-red-500/40 hover:bg-red-500/12 hover:-translate-y-0.5" @click="handleDelete(row)" :title="t('common.delete')">
                      <Trash2 class="w-3.5 h-3.5" />
                    </button>
                  </div>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>

        <!-- Pagination -->
        <div class="flex justify-end items-center px-5 py-4 border-t border-border">
          <div class="flex items-center gap-4 text-sm">
            <span class="text-muted-foreground">共 {{ pagination.total }} 条</span>
            <Select v-model="paginationSizeString" @update:model-value="loadRoles">
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
              <Button variant="outline" size="sm" :disabled="pagination.page <= 1" @click="pagination.page--; loadRoles()">
                <ChevronLeft class="w-4 h-4" />
              </Button>
              <span class="px-3 text-sm">{{ pagination.page }}</span>
              <Button variant="outline" size="sm" :disabled="pagination.page * pagination.size >= pagination.total" @click="pagination.page++; loadRoles()">
                <ChevronRight class="w-4 h-4" />
              </Button>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:open="showCreateDialog">
      <DialogContent class="max-w-[560px]">
        <DialogHeader>
          <DialogTitle>{{ editingRole ? t('role.editRole') : t('role.createRole') }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label>{{ t('role.roleName') }}</Label>
            <Input v-model="roleForm.name" :placeholder="t('role.roleName')" :disabled="!!editingRole" />
          </div>
          <div class="space-y-2">
            <Label>{{ t('role.description') }}</Label>
            <Textarea v-model="roleForm.description" :placeholder="t('role.description')" :rows="3" />
          </div>
          <div v-if="editingRole" class="space-y-2">
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
          <div class="flex items-center gap-2">
            <Switch v-model:checked="roleForm.is_system_role" :disabled="!!editingRole" />
            <Label>{{ t('role.isSystemRole') }}</Label>
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
const roleFormRef = ref()

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
  set: (val: string) => { pagination.size = Number(val) },
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
  } catch (error) {
    toast.error('获取角色列表失败')
  } finally {
    loading.value = false
    initialLoading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; loadRoles() }
const handleReset = () => { searchForm.name = ''; searchForm.status = undefined; searchForm.isSystemRole = undefined; pagination.page = 1; loadRoles() }
const handleView = (role: RoleDefinitionDTO) => { router.push({ name: 'RoleDetail', params: { id: role.id } }) }
const handleEdit = (role: RoleDefinitionDTO) => {
  editingRole.value = role
  roleForm.name = role.name || ''; roleForm.description = role.description || ''; roleForm.status = role.status || 1; roleForm.is_system_role = role.is_system_role || false
  showCreateDialog.value = true
}

const handleDelete = async (role: RoleDefinitionDTO) => {
  try {
    await roleApi.deleteRole(role.id!)
    toast.success('删除成功'); loadRoles()
  } catch (error: any) {
    toast.error('删除角色失败')
  }
}

const handleSaveRole = async () => {
  if (!roleFormRef.value) return
  try {
    await roleFormRef.value.validate()
    const data = { name: roleForm.name, description: roleForm.description, is_system_role: roleForm.is_system_role, permissions: roleForm.permissions }
    let response: any
    if (editingRole.value) {
      response = await roleApi.updateRole(editingRole.value.id!, { ...data, status: roleForm.status })
    } else {
      response = await roleApi.createRole(data)
    }
    if (response.role) {
      toast.success(editingRole.value ? '更新成功' : '创建成功')
      showCreateDialog.value = false; loadRoles()
    }
  } catch {}
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = { 1: t('common.active'), 2: t('common.inactive'), 3: t('common.deprecated') }
  return map[status] || ''
}
const formatTimestamp = (ts?: number) => ts ? new Date(ts).toLocaleDateString('zh-CN') : '-'

onMounted(() => { loadRoles() })
</script>

<style scoped>
.spotlight-card {
  position: relative;
  overflow: hidden;
}

.spotlight-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(520px circle at var(--mouse-x, -100%) var(--mouse-y, -100%),
    rgba(63, 81, 181, 0.12), transparent 55%);
  pointer-events: none;
  z-index: 1;
}
</style>
