<template>
  <div class="space-y-6">
    <!-- 标题区（DESIGN.md 详情页范式） -->
    <header class="space-y-4 pb-5 border-b border-subtle">
      <button
        class="inline-flex items-center gap-1 text-[13px] text-[color:var(--color-ink-muted)] transition-colors hover:text-[color:var(--color-ink)]"
        @click="handleBack"
      >
        <ArrowLeft class="h-3.5 w-3.5" />
        {{ t('common.back') }}
      </button>

      <div class="flex items-start justify-between gap-4">
        <div class="flex items-center gap-3 min-w-0">
          <span
            class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-md"
            :class="roleDetail?.is_system_role
              ? 'bg-[color:var(--color-warning-soft)] text-[color:var(--color-warning-ink)]'
              : 'bg-[color:var(--color-primary-soft)] text-[color:var(--color-primary-active)]'"
          >
            <KeyRound class="h-5 w-5" />
          </span>
          <div class="min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <h1 class="text-[22px] font-semibold leading-tight tracking-[-0.012em] text-[color:var(--color-ink-strong)] truncate">
                {{ roleDetail?.name || t('role.roleDetail') }}
              </h1>
              <Badge v-if="roleDetail?.is_system_role" variant="warning">{{ t('role.systemBadge') }}</Badge>
              <Badge v-if="roleDetail" :variant="getStatusVariant(roleDetail.status)">{{ getStatusText(roleDetail.status) }}</Badge>
            </div>
            <p class="mt-1 text-[13px] text-[color:var(--color-ink-muted)] truncate">
              <span v-if="roleDetail?.description">{{ roleDetail.description }}</span>
              <span v-else class="text-[color:var(--color-ink-subtle)]">{{ t('role.noDescription') }}</span>
              <span class="mx-2 text-[color:var(--color-ink-subtle)]">·</span>
              <span class="tabular">{{ roleDetail?.user_count || 0 }} {{ t('role.usersUnit') }}</span>
            </p>
          </div>
        </div>

        <div class="flex items-center gap-2 flex-shrink-0">
          <Button variant="outline" size="sm" @click="showEditDialog = true">
            <Pencil class="h-3.5 w-3.5" />
            {{ t('common.edit') }}
          </Button>
          <Button
            variant="destructive"
            size="sm"
            :disabled="roleDetail?.is_system_role"
            :title="roleDetail?.is_system_role ? t('role.systemRoleNoDelete') : ''"
            @click="handleDelete"
          >
            <Trash2 class="h-3.5 w-3.5" />
            {{ t('common.delete') }}
          </Button>
        </div>
      </div>
    </header>

    <div class="min-h-[300px]">
      <DetailPageSkeleton
        v-if="initialLoading"
        :side-span="16"
        :main-span="8"
        :side-cards="2"
        :main-cards="2"
        :show-avatar="false"
        :side-item-counts="[6, 4]"
        :main-item-counts="[5, 4]"
      />
      <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-5">
        <!-- Left Column: Basic Info + API Permissions -->
        <div class="lg:col-span-2 space-y-5">
          <Card v-if="roleDetail" class="bg-card border-border/60 rounded-[18px] shadow-card">
            <CardHeader class="border-b border-border/60 px-5 py-3.5">
              <div class="flex items-center gap-2 text-primary font-semibold">
                <KeyRound class="w-4 h-4" />
                <span>{{ t('role.basicInfo') }}</span>
              </div>
            </CardHeader>
            <CardContent class="p-5">
              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-1">
                  <span class="text-xs text-muted-foreground">{{ t('role.roleName') }}</span>
                  <p class="text-sm font-semibold text-foreground">{{ roleDetail.name }}</p>
                </div>
                <div class="space-y-1">
                  <span class="text-xs text-muted-foreground">{{ t('role.userCount') }}</span>
                  <p class="text-sm text-foreground tabular">{{ roleDetail.user_count || 0 }}</p>
                </div>
                <div class="space-y-1 col-span-2">
                  <span class="text-xs text-muted-foreground">{{ t('role.description') }}</span>
                  <p class="text-sm text-foreground">{{ roleDetail.description || '-' }}</p>
                </div>
                <div class="space-y-1">
                  <span class="text-xs text-muted-foreground">{{ t('common.createTime') }}</span>
                  <p class="text-sm text-foreground">{{ formatTimestamp(roleDetail.created_at) }}</p>
                </div>
                <div class="space-y-1">
                  <span class="text-xs text-muted-foreground">{{ t('common.updateTime') }}</span>
                  <p class="text-sm text-foreground">{{ formatTimestamp(roleDetail.updated_at) }}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card v-if="roleDetail?.permissions?.length" class="bg-card border-border/60 rounded-[18px] shadow-card">
            <CardHeader class="border-b border-border/60 px-5 py-3.5">
              <div class="flex items-center gap-2 text-primary font-semibold">
                <Lock class="w-4 h-4" />
                <span>{{ t('role.permissions') }}</span>
              </div>
            </CardHeader>
            <CardContent class="p-0">
              <div class="overflow-auto max-h-[300px]">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead class="w-[200px]">{{ t('role.resource') }}</TableHead>
                      <TableHead class="w-[150px]">{{ t('role.action') }}</TableHead>
                      <TableHead>{{ t('common.description') }}</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow v-for="(perm, i) in roleDetail!.permissions" :key="i">
                      <TableCell>{{ perm.resource }}</TableCell>
                      <TableCell>{{ perm.action }}</TableCell>
                      <TableCell class="text-muted-foreground">{{ perm.description || '-' }}</TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Right Column: Menu Permissions + Bound Users -->
        <div class="space-y-5">
          <Card class="bg-card border-border/60 rounded-[18px] shadow-card">
            <CardHeader class="border-b border-border/60 px-5 py-3.5">
              <div class="flex items-center gap-2 text-primary font-semibold">
                <LayoutList class="w-4 h-4" />
                <span class="flex-1">{{ t('role.menuPermission') }}</span>
                <Button size="sm" @click="openMenuConfigDialog">
                  {{ t('role.configurePermission') }}
                </Button>
              </div>
            </CardHeader>
            <CardContent class="p-5">
              <div v-if="menuTree.length" class="space-y-1">
                <template v-for="node in flattenMenuTree(menuTree)" :key="node.id">
                  <div class="flex items-center gap-2.5 py-1.5 px-2 rounded-md hover:bg-accent/50" :style="{ paddingLeft: `${(node._depth || 0) * 16 + 8}px` }">
                    <span class="flex-1 text-sm">{{ node.name }}</span>
                    <Badge v-if="node.permission_level !== undefined" size="sm" :variant="getPermBadgeVariant(node.permission_level)" class="text-xs">
                      {{ getPermissionText(node.permission_level) }}
                    </Badge>
                  </div>
                </template>
              </div>
              <div v-else class="text-center py-6 text-muted-foreground text-sm">
                {{ t('common.noData') }}
              </div>
            </CardContent>
          </Card>

          <Card class="bg-card border-border/60 rounded-[18px] shadow-card">
            <CardHeader class="border-b border-border/60 px-5 py-3.5">
              <div class="flex items-center gap-2 text-primary font-semibold">
                <User class="w-4 h-4" />
                <span class="flex-1">已绑定用户（{{ roleUsers.user_ids?.length || 0 }}）</span>
                <Button size="sm" @click="openUserAssignDialog">
                  管理用户
                </Button>
              </div>
            </CardHeader>
            <CardContent class="p-5">
              <div v-if="roleUsers.user_ids?.length" class="flex flex-wrap gap-2 max-h-[220px] overflow-y-auto">
                <div
                  v-for="uid in roleUsers.user_ids"
                  :key="uid"
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-md border bg-secondary text-xs font-mono group cursor-default"
                >
                  <span>{{ userDisplayName(uid) }}</span>
                  <button class="opacity-50 hover:opacity-100 transition-opacity" @click="handleRemoveUser(uid)">
                    <X class="w-3 h-3" />
                  </button>
                </div>
              </div>
              <div v-else class="text-center py-6 text-muted-foreground text-sm">
                {{ t('common.noData') }}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>

    <!-- Edit Role Dialog -->
    <Dialog v-model:open="showEditDialog">
      <DialogContent class="max-w-[560px]">
        <DialogHeader>
          <DialogTitle>{{ t('role.editRole') }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label>{{ t('role.roleName') }}</Label>
            <Input :value="roleDetail?.name" disabled />
          </div>
          <div class="space-y-2">
            <Label>{{ t('role.description') }}</Label>
            <Textarea v-model="editForm.description" :rows="3" />
          </div>
          <div class="space-y-2">
            <Label>{{ t('common.status') }}</Label>
            <Select v-model="editFormStatusString">
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
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showEditDialog = false">{{ t('common.cancel') }}</Button>
          <Button :disabled="saving" @click="handleUpdateRole">
            {{ t('common.save') }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Menu Permission Config Dialog -->
    <Dialog v-model:open="showMenuConfigDialog">
      <DialogContent class="max-w-[700px]">
        <DialogHeader>
          <DialogTitle>配置菜单权限</DialogTitle>
        </DialogHeader>
        <div class="py-4">
          <div class="flex items-center gap-2 p-2.5 mb-4 bg-primary/8 border border-primary/20 rounded-lg text-sm text-muted-foreground">
            <Info class="w-4 h-4 flex-shrink-0" />
            为角色配置各菜单的访问权限级别。设为"无权限"表示该角色不可访问该菜单，仅保存非零权限。
          </div>

          <div class="max-h-[460px] overflow-y-auto">
            <template v-if="fullMenuTree.length">
              <div class="flex items-center px-6 pb-2 text-[11px] font-semibold tracking-wider uppercase text-muted-foreground border-b border-border mb-1">
                <span class="flex-1">菜单名称</span>
                <span class="w-[140px]">路径</span>
                <span class="w-[110px] flex-shrink-0">权限级别</span>
              </div>
              <div v-for="node in flattenMenuTree(fullMenuTree)" :key="node.id" class="flex items-center gap-3 py-1 px-2 rounded-md hover:bg-accent/50" :style="{ paddingLeft: `${(node._depth || 0) * 16 + 24}px` }">
                <span class="flex-1 text-sm text-foreground">{{ node.name }}</span>
                <span class="w-[140px] text-xs text-muted-foreground font-mono truncate flex-shrink-0">{{ node.path }}</span>
                <Select v-model="menuPermMap[node.id]">
                  <SelectTrigger class="w-[110px]" :class="{ 'text-primary font-medium': Number(menuPermMap[node.id] ?? '0') > 0 }">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem :value="String(0)">无权限</SelectItem>
                    <SelectItem :value="String(1)">只读</SelectItem>
                    <SelectItem :value="String(2)">读写</SelectItem>
                    <SelectItem :value="String(3)">完全权限</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </template>
            <div v-else-if="!menuConfigLoading" class="text-center py-6 text-muted-foreground text-sm">
              {{ t('common.noData') }}
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showMenuConfigDialog = false">{{ t('common.cancel') }}</Button>
          <Button :disabled="savingMenuConfig" @click="handleSaveMenuConfig">
            {{ t('common.save') }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- User Assign Dialog -->
    <Dialog v-model:open="showUserAssignDialog">
      <DialogContent class="max-w-[800px]">
        <DialogHeader>
          <DialogTitle>管理角色用户</DialogTitle>
        </DialogHeader>
        <div class="py-4">
          <div class="flex items-center gap-4 mb-3.5">
            <div class="relative w-[300px]">
              <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
              <Input
                v-model="userSearchKeyword"
                placeholder="搜索用户名或真实姓名"
                class="pl-9"
              />
              <button v-if="userSearchKeyword" @click="userSearchKeyword = ''; handleUserSearch()" class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground">
                <X class="w-4 h-4" />
              </button>
            </div>
            <span class="text-sm text-muted-foreground">
              已选 <strong class="text-primary text-base">{{ pendingUserIds.size }}</strong> 人
            </span>
          </div>

          <div class="overflow-auto max-h-[380px]">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead class="w-[50px]"></TableHead>
                  <TableHead>{{ t('user.username') }}</TableHead>
                  <TableHead class="w-[110px]">{{ t('user.realName') }}</TableHead>
                  <TableHead>{{ t('user.email') }}</TableHead>
                  <TableHead class="w-[90px] text-center">{{ t('common.status') }}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="user in searchedUsers" :key="user.id">
                  <TableCell>
                    <Checkbox :checked="pendingUserIds.has(user.id)" @update:checked="(val: boolean) => toggleUserSelect(user, val)" />
                  </TableCell>
                  <TableCell class="font-medium">{{ user.username }}</TableCell>
                  <TableCell class="text-muted-foreground">{{ user.real_name || '-' }}</TableCell>
                  <TableCell class="text-muted-foreground">{{ user.email || '-' }}</TableCell>
                  <TableCell class="text-center">
                    <Badge variant="outline" class="text-xs" :class="user.status === 1 ? 'bg-green-500/10 border-green-500/30 text-green-500' : ''">
                      {{ user.status === 1 ? t('common.active') : t('common.inactive') }}
                    </Badge>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showUserAssignDialog = false">{{ t('common.cancel') }}</Button>
          <Button :disabled="savingUsers" @click="handleSaveUserAssign">
            保存（{{ pendingUserIds.size }} 人）
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { ArrowLeft, Pencil, Trash2, KeyRound, Lock, User, LayoutList, Search, Info, X } from 'lucide-vue-next'
import { roleApi } from '@/api/role'
import { menuApi } from '@/api/menu'
import { getUserList, getUserDetail } from '@/api/user'
import type { RoleDefinitionDTO, MenuNodeDTO } from '@/api/role'
import type { MenuItem, UserListItem, UserProfile } from '@/types/user'
import DetailPageSkeleton from '@/components/skeleton/DetailPageSkeleton.vue'
import { Card, CardHeader, CardContent } from '@/components/ui/card'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const roleId = route.params.id as string

// ── 主数据 ──────────────────────────────────────────────────────────
const initialLoading = ref(true)
const loading = ref(false)
const saving = ref(false)
const roleDetail = ref<RoleDefinitionDTO | null>(null)
const menuTree = ref<MenuNodeDTO[]>([])
const roleUsers = ref<{ role_id?: string; user_ids?: string[] }>({})
const userDetailMap = ref<Record<string, UserProfile>>({})

// ── 对话框控制 ───────────────────────────────────────────────────────
const showEditDialog = ref(false)
const showMenuConfigDialog = ref(false)
const showUserAssignDialog = ref(false)

// ── 编辑表单 ─────────────────────────────────────────────────────────
const editFormRef = ref()
const editForm = reactive({ description: '', status: 1 })

const editFormStatusString = computed({
  get: () => String(editForm.status),
  set: (val: string) => { editForm.status = Number(val) },
})

// ── 菜单权限配置 ──────────────────────────────────────────────────────
const menuConfigLoading = ref(false)
const savingMenuConfig = ref(false)
const fullMenuTree = ref<MenuItem[]>([])
const menuPermMap = reactive<Record<string, string>>({})

// ── 用户分配 ──────────────────────────────────────────────────────────
const userSearchLoading = ref(false)
const savingUsers = ref(false)
const userSearchKeyword = ref('')
const searchedUsers = ref<UserListItem[]>([])
const pendingUserIds = reactive(new Set<string>())
let userSearchTimer: ReturnType<typeof setTimeout> | null = null

// ════════════════════ 数据加载 ════════════════════

const loadRoleDetail = async () => {
  loading.value = true
  try {
    const res = await roleApi.getRole(roleId)
    roleDetail.value = res.role || null
    editForm.description = res.role?.description || ''
    editForm.status = res.role?.status || 1
  } catch {
    toast.error('获取角色详情失败')
  } finally {
    loading.value = false
  }
}

const loadMenuTree = async () => {
  try {
    const res = await roleApi.getRoleMenuTree(roleId)
    menuTree.value = res.menu_tree || []
  } catch {
    // 静默失败，菜单树显示空
  }
}

const loadRoleUsers = async () => {
  try {
    const res = await roleApi.getRoleUsers(roleId)
    roleUsers.value = res
    await loadRoleUserDetails(res.user_ids || [])
  } catch {
    // 静默失败
  }
}

const loadRoleUserDetails = async (userIds: string[]) => {
  if (!userIds.length) return
  const results = await Promise.allSettled(userIds.map(id => getUserDetail(id)))
  const map: Record<string, UserProfile> = {}
  results.forEach((result, i) => {
    const userId = userIds[i]
    if (result.status === 'fulfilled' && userId) {
      map[userId] = result.value.user
    }
  })
  userDetailMap.value = map
}

// ════════════════════ 编辑角色 ════════════════════

const handleUpdateRole = async () => {
  if (!editFormRef.value) return
  try {
    await editFormRef.value.validate()
    saving.value = true
    const res = await roleApi.updateRole(roleId, editForm)
    if (res.role) {
      toast.success(t('common.updateSuccess'))
      showEditDialog.value = false
      loadRoleDetail()
    }
  } catch {
    // 表单校验失败
  } finally {
    saving.value = false
  }
}

// ════════════════════ 删除角色 ════════════════════

const handleDelete = async () => {
  if (!roleDetail.value) return
  try {
    await roleApi.deleteRole(roleId)
    toast.success(t('common.deleteSuccess'))
    router.back()
  } catch (error: any) {
    if (error !== 'cancel') toast.error('删除角色失败')
  }
}

// ════════════════════ 菜单权限配置 ════════════════════

const openMenuConfigDialog = async () => {
  showMenuConfigDialog.value = true
  menuConfigLoading.value = true
  try {
    const [menuRes, permRes] = await Promise.all([
      menuApi.getMenuTree(),
      roleApi.getRoleMenuPermissions(roleId),
    ])
    fullMenuTree.value = menuRes.menu_tree || []

    // 初始化所有节点为 0（无权限）
    Object.keys(menuPermMap).forEach(k => delete menuPermMap[k])
    flattenTree(fullMenuTree.value).forEach(node => {
      menuPermMap[node.id] = '0'
    })

    // 填入当前角色已有的权限
    const perms = permRes.permissions || []
    perms.forEach(p => {
      menuPermMap[p.menu_id] = String(p.permission)
    })
  } catch {
    toast.error('加载菜单配置失败')
  } finally {
    menuConfigLoading.value = false
  }
}

const handleSaveMenuConfig = async () => {
  savingMenuConfig.value = true
  try {
    // 仅提交权限 > 0 的配置
    const configs = Object.entries(menuPermMap)
      .filter(([, perm]) => Number(perm) > 0)
      .map(([menu_id, permission]) => ({ menu_id, permission: Number(permission) }))

    await roleApi.configureRoleMenus(roleId, configs)
    toast.success('菜单权限配置已保存')
    showMenuConfigDialog.value = false
    loadMenuTree()
  } catch {
    toast.error('保存菜单权限失败')
  } finally {
    savingMenuConfig.value = false
  }
}

function flattenTree(nodes: MenuItem[], depth = 0): (MenuItem & { _depth: number })[] {
  const result: (MenuItem & { _depth: number })[] = []
  for (const node of nodes) {
    result.push({ ...node, _depth: depth })
    if (node.children?.length) result.push(...flattenTree(node.children, depth + 1))
  }
  return result
}

function flattenMenuTree(nodes: any[], depth = 0): any[] {
  const result: any[] = []
  for (const node of nodes) {
    result.push({ ...node, _depth: depth })
    if (node.children?.length) result.push(...flattenMenuTree(node.children, depth + 1))
  }
  return result
}

// ════════════════════ 用户分配 ════════════════════

const openUserAssignDialog = () => {
  showUserAssignDialog.value = true
  pendingUserIds.clear()
  ;(roleUsers.value.user_ids || []).forEach(id => pendingUserIds.add(id))
  userSearchKeyword.value = ''
  loadUsers()
}

const loadUsers = async (search?: string) => {
  userSearchLoading.value = true
  try {
    const res = await getUserList({ page: 1, limit: 50, search })
    searchedUsers.value = res.users || []
  } catch {
    toast.error('加载用户列表失败')
  } finally {
    userSearchLoading.value = false
  }
}

const toggleUserSelect = (user: UserListItem, checked: boolean) => {
  if (checked) {
    pendingUserIds.add(user.id)
  } else {
    pendingUserIds.delete(user.id)
  }
}

const handleUserSearch = () => {
  if (userSearchTimer) clearTimeout(userSearchTimer)
  userSearchTimer = setTimeout(() => {
    loadUsers(userSearchKeyword.value || undefined)
  }, 400)
}

const handleSaveUserAssign = async () => {
  savingUsers.value = true
  try {
    await roleApi.bindUsersToRole(roleId, [...pendingUserIds])
    toast.success('用户分配已保存')
    showUserAssignDialog.value = false
    loadRoleUsers()
    loadRoleDetail()
  } catch {
    toast.error('保存用户分配失败')
  } finally {
    savingUsers.value = false
  }
}

const handleRemoveUser = async (userId: string) => {
  try {
    const remaining = (roleUsers.value.user_ids || []).filter(id => id !== userId)
    await roleApi.bindUsersToRole(roleId, remaining)
    toast.success(t('common.removeSuccess'))
    loadRoleUsers()
    loadRoleDetail()
  } catch (error: any) {
    if (error !== 'cancel') toast.error('移除用户失败')
  }
}

// ════════════════════ 工具函数 ════════════════════

const handleBack = () => router.back()

const userDisplayName = (uid: string) => {
  const u = userDetailMap.value[uid]
  if (!u) return uid
  return u.real_name || u.username
}

const getStatusText = (status: number) =>
  ({ 1: t('common.active'), 2: t('common.inactive'), 3: t('common.deprecated') } as Record<number, string>)[status] || ''

function getStatusVariant(status: number): 'success' | 'default' | 'destructive' {
  const variantMap: Record<number, 'success' | 'default' | 'destructive'> = {
    1: 'success', 2: 'default', 3: 'destructive'
  }
  return variantMap[status] || 'default'
}

const getPermBadgeVariant = (level: number): 'default' | 'secondary' | 'destructive' | 'outline' => {
  if (level === 0) return 'secondary'
  if (level === 1) return 'outline'
  if (level === 2) return 'default'
  return 'destructive'
}

const getPermissionText = (level: number) =>
  ({
    0: t('role.permissionLevel.none'),
    1: t('role.permissionLevel.read'),
    2: t('role.permissionLevel.write'),
    3: t('role.permissionLevel.full'),
  } as Record<number, string>)[level] || ''

const formatTimestamp = (ts?: number) => (ts ? new Date(ts).toLocaleString('zh-CN') : '-')

onMounted(async () => {
  try {
    await Promise.all([loadRoleDetail(), loadMenuTree(), loadRoleUsers()])
  } finally {
    initialLoading.value = false
  }
})

onBeforeUnmount(() => {
  showEditDialog.value = false
  showMenuConfigDialog.value = false
  showUserAssignDialog.value = false

  if (userSearchTimer) {
    clearTimeout(userSearchTimer)
    userSearchTimer = null
  }
})
</script>

<style scoped>
</style>
