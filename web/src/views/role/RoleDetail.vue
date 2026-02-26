<template>
  <div class="role-detail" v-loading="loading">
    <el-page-header @back="handleBack" :title="t('common.back')">
      <template #content>
        <span class="page-title">{{ roleDetail?.name || t('role.roleDetail') }}</span>
      </template>
      <template #extra>
        <el-button-group>
          <el-button @click="showEditDialog = true">
            <el-icon><Edit /></el-icon>
            {{ t('common.edit') }}
          </el-button>
          <el-button
            type="danger"
            :disabled="roleDetail?.is_system_role"
            :title="roleDetail?.is_system_role ? '系统角色不可删除' : ''"
            @click="handleDelete"
          >
            <el-icon><Delete /></el-icon>
            {{ t('common.delete') }}
          </el-button>
        </el-button-group>
      </template>
    </el-page-header>

    <el-row :gutter="20" class="content-row">
      <!-- 左列：基本信息 + API 权限 -->
      <el-col :span="16">
        <el-card class="detail-card" v-if="roleDetail">
          <template #header>
            <div class="card-header">
              <el-icon><Key /></el-icon>
              <span>基本信息</span>
              <el-tag v-if="roleDetail.is_system_role" type="danger" size="small">系统角色</el-tag>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item :label="t('role.roleName')">
              <strong class="role-name-text">{{ roleDetail.name }}</strong>
            </el-descriptions-item>
            <el-descriptions-item :label="t('common.status')">
              <el-tag :type="getStatusType(roleDetail.status)" size="small">
                {{ getStatusText(roleDetail.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="t('role.isSystemRole')">
              <el-tag :type="roleDetail.is_system_role ? 'danger' : 'success'" size="small">
                {{ roleDetail.is_system_role ? t('common.yes') : t('common.no') }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="t('role.userCount')">
              {{ roleDetail.user_count || 0 }} 人
            </el-descriptions-item>
            <el-descriptions-item :label="t('role.description')" :span="2">
              {{ roleDetail.description || '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('common.createTime')">
              {{ formatTimestamp(roleDetail.created_at) }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('common.updateTime')">
              {{ formatTimestamp(roleDetail.updated_at) }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card class="detail-card" v-if="roleDetail?.permissions?.length">
          <template #header>
            <div class="card-header">
              <el-icon><Lock /></el-icon>
              <span>{{ t('role.permissions') }}</span>
            </div>
          </template>
          <el-table :data="roleDetail!.permissions" border stripe max-height="300">
            <el-table-column prop="resource" :label="t('role.resource')" width="200" />
            <el-table-column prop="action" :label="t('role.action')" width="150" />
            <el-table-column prop="description" :label="t('common.description')" show-overflow-tooltip />
          </el-table>
        </el-card>
      </el-col>

      <!-- 右列：菜单权限 + 已绑定用户 -->
      <el-col :span="8">
        <el-card class="detail-card menu-card">
          <template #header>
            <div class="card-header">
              <el-icon><Menu /></el-icon>
              <span>{{ t('role.menuPermission') }}</span>
              <el-button type="primary" size="small" @click="openMenuConfigDialog">
                {{ t('role.configurePermission') }}
              </el-button>
            </div>
          </template>
          <el-tree
            v-if="menuTree.length"
            :data="menuTree"
            :props="{ label: 'name', children: 'children' }"
            node-key="id"
            default-expand-all
            :expand-on-click-node="false"
          >
            <template #default="{ node, data }">
              <div class="menu-tree-node">
                <span class="node-name">{{ node.label }}</span>
                <el-tag
                  v-if="data.permission_level !== undefined"
                  size="small"
                  :type="getPermTagType(data.permission_level)"
                >
                  {{ getPermissionText(data.permission_level) }}
                </el-tag>
              </div>
            </template>
          </el-tree>
          <el-empty v-else :description="t('common.noData')" :image-size="60" />
        </el-card>

        <el-card class="detail-card users-card">
          <template #header>
            <div class="card-header">
              <el-icon><User /></el-icon>
              <span>已绑定用户（{{ roleUsers.user_ids?.length || 0 }}）</span>
              <el-button type="primary" size="small" @click="openUserAssignDialog">
                管理用户
              </el-button>
            </div>
          </template>
          <div v-if="roleUsers.user_ids?.length" class="user-tags">
            <el-tooltip
              v-for="uid in roleUsers.user_ids"
              :key="uid"
              :content="`ID: ${uid}`"
              placement="top"
            >
              <el-tag
                closable
                class="user-tag"
                @close="handleRemoveUser(uid)"
              >
                {{ userDisplayName(uid) }}
              </el-tag>
            </el-tooltip>
          </div>
          <el-empty v-else :description="t('common.noData')" :image-size="60" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 编辑角色对话框 -->
    <el-dialog
      v-model="showEditDialog"
      :title="t('role.editRole')"
      width="560px"
      @close="editFormRef?.resetFields()"
    >
      <el-form ref="editFormRef" :model="editForm" :rules="editFormRules" label-width="100px">
        <el-form-item :label="t('role.roleName')">
          <el-input :value="roleDetail?.name" disabled />
        </el-form-item>
        <el-form-item :label="t('role.description')" prop="description">
          <el-input v-model="editForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item :label="t('common.status')" prop="status">
          <el-select v-model="editForm.status" style="width: 100%">
            <el-option :label="t('common.active')" :value="1" />
            <el-option :label="t('common.inactive')" :value="2" />
            <el-option :label="t('common.deprecated')" :value="3" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="saving" @click="handleUpdateRole">
          {{ t('common.save') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 菜单权限配置对话框 -->
    <el-dialog
      v-model="showMenuConfigDialog"
      title="配置菜单权限"
      width="700px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <div class="config-tips">
        <el-icon><InfoFilled /></el-icon>
        为角色配置各菜单的访问权限级别。设为"无权限"表示该角色不可访问该菜单，仅保存非零权限。
      </div>

      <div class="menu-config-body" v-loading="menuConfigLoading">
        <template v-if="fullMenuTree.length">
          <div class="menu-config-header">
            <span>菜单名称</span>
            <span>路径</span>
            <span>权限级别</span>
          </div>
          <el-tree
            :data="fullMenuTree"
            :props="{ label: 'name', children: 'children' }"
            node-key="id"
            default-expand-all
            :expand-on-click-node="false"
          >
            <template #default="{ data }">
              <div class="menu-config-node">
                <span class="node-label">{{ data.name }}</span>
                <span class="node-path">{{ data.path }}</span>
                <el-select
                  v-model="menuPermMap[data.id]"
                  size="small"
                  class="perm-select"
                  :class="{ 'has-perm': menuPermMap[data.id] > 0 }"
                >
                  <el-option :value="0" label="无权限" />
                  <el-option :value="1" label="只读" />
                  <el-option :value="2" label="读写" />
                  <el-option :value="3" label="完全权限" />
                </el-select>
              </div>
            </template>
          </el-tree>
        </template>
        <el-empty v-else-if="!menuConfigLoading" :description="t('common.noData')" />
      </div>

      <template #footer>
        <el-button @click="showMenuConfigDialog = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="savingMenuConfig" @click="handleSaveMenuConfig">
          {{ t('common.save') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 用户分配对话框 -->
    <el-dialog
      v-model="showUserAssignDialog"
      title="管理角色用户"
      width="800px"
      :close-on-click-modal="false"
      destroy-on-close
      @open="initUserAssignDialog"
    >
      <div class="user-assign-toolbar">
        <el-input
          v-model="userSearchKeyword"
          placeholder="搜索用户名或真实姓名"
          clearable
          style="width: 300px"
          @input="handleUserSearch"
        >
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <span class="assign-count">
          已选 <strong>{{ pendingUserIds.size }}</strong> 人
        </span>
      </div>

      <el-table
        ref="userTableRef"
        v-loading="userSearchLoading"
        :data="searchedUsers"
        row-key="id"
        max-height="380"
        class="user-assign-table"
        @select="handleUserSelect"
        @select-all="handleUserSelectAll"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column prop="username" :label="t('user.username')" min-width="120" />
        <el-table-column :label="t('user.realName')" width="110">
          <template #default="{ row }">{{ row.real_name || '-' }}</template>
        </el-table-column>
        <el-table-column prop="email" :label="t('user.email')" show-overflow-tooltip min-width="160" />
        <el-table-column :label="t('common.status')" width="90" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? t('common.active') : t('common.inactive') }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <el-button @click="showUserAssignDialog = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="savingUsers" @click="handleSaveUserAssign">
          保存（{{ pendingUserIds.size }} 人）
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Edit, Delete, Key, Lock, User, Menu, Search, InfoFilled } from '@element-plus/icons-vue'
import { roleApi } from '@/api/role'
import { menuApi } from '@/api/menu'
import { getUserList, getUserDetail } from '@/api/user'
import type { RoleDefinitionDTO, MenuNodeDTO } from '@/api/role'
import type { MenuItem, UserListItem, UserProfile } from '@/types/user'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const roleId = route.params.id as string

// ── 主数据 ──────────────────────────────────────────────────────────
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
const editFormRef = ref<FormInstance>()
const editForm = reactive({ description: '', status: 1 })
const editFormRules: FormRules = {
  description: [{ max: 500, message: '描述不超过 500 个字符', trigger: 'blur' }],
}

// ── 菜单权限配置 ──────────────────────────────────────────────────────
const menuConfigLoading = ref(false)
const savingMenuConfig = ref(false)
const fullMenuTree = ref<MenuItem[]>([])
const menuPermMap = reactive<Record<string, number>>({})

// ── 用户分配 ──────────────────────────────────────────────────────────
const userSearchLoading = ref(false)
const savingUsers = ref(false)
const userSearchKeyword = ref('')
const searchedUsers = ref<UserListItem[]>([])
const pendingUserIds = reactive(new Set<string>())
const userTableRef = ref<any>()
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
    ElMessage.error('获取角色详情失败')
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
    if (result.status === 'fulfilled') {
      map[userIds[i]] = result.value.user
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
      ElMessage.success(t('common.updateSuccess'))
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
    await ElMessageBox.confirm(
      `确定要删除角色 "${roleDetail.value.name}" 吗？此操作不可恢复！`,
      '警告',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      }
    )
    await roleApi.deleteRole(roleId)
    ElMessage.success(t('common.deleteSuccess'))
    router.back()
  } catch (error: any) {
    if (error !== 'cancel') ElMessage.error('删除角色失败')
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
      menuPermMap[node.id] = 0
    })

    // 填入当前角色已有的权限
    const perms = permRes.permissions || []
    perms.forEach(p => {
      menuPermMap[p.menu_id] = p.permission
    })
  } catch {
    ElMessage.error('加载菜单配置失败')
  } finally {
    menuConfigLoading.value = false
  }
}

const handleSaveMenuConfig = async () => {
  savingMenuConfig.value = true
  try {
    // 仅提交权限 > 0 的配置
    const configs = Object.entries(menuPermMap)
      .filter(([, perm]) => perm > 0)
      .map(([menu_id, permission]) => ({ menu_id, permission }))

    await roleApi.configureRoleMenus(roleId, configs)
    ElMessage.success('菜单权限配置已保存')
    showMenuConfigDialog.value = false
    loadMenuTree()
  } catch {
    ElMessage.error('保存菜单权限失败')
  } finally {
    savingMenuConfig.value = false
  }
}

function flattenTree(nodes: MenuItem[]): MenuItem[] {
  const result: MenuItem[] = []
  for (const node of nodes) {
    result.push(node)
    if (node.children?.length) result.push(...flattenTree(node.children))
  }
  return result
}

// ════════════════════ 用户分配 ════════════════════

const openUserAssignDialog = () => {
  showUserAssignDialog.value = true
}

const initUserAssignDialog = async () => {
  // 初始化待选集合为当前已绑定用户
  pendingUserIds.clear()
  ;(roleUsers.value.user_ids || []).forEach(id => pendingUserIds.add(id))
  userSearchKeyword.value = ''
  await loadUsers()
}

const loadUsers = async (search?: string) => {
  userSearchLoading.value = true
  try {
    const res = await getUserList({ page: 1, limit: 50, search })
    searchedUsers.value = res.users || []
    // 恢复选中状态
    await nextTick()
    if (userTableRef.value) {
      searchedUsers.value.forEach(user => {
        userTableRef.value.toggleRowSelection(user, pendingUserIds.has(user.id))
      })
    }
  } catch {
    ElMessage.error('加载用户列表失败')
  } finally {
    userSearchLoading.value = false
  }
}

const handleUserSearch = () => {
  if (userSearchTimer) clearTimeout(userSearchTimer)
  userSearchTimer = setTimeout(() => {
    loadUsers(userSearchKeyword.value || undefined)
  }, 400)
}

const handleUserSelect = (selection: UserListItem[], row: UserListItem) => {
  if (selection.some(u => u.id === row.id)) {
    pendingUserIds.add(row.id)
  } else {
    pendingUserIds.delete(row.id)
  }
}

const handleUserSelectAll = (selection: UserListItem[]) => {
  // 取消当前页所有行的选中
  searchedUsers.value.forEach(u => pendingUserIds.delete(u.id))
  // 重新加入选中的行
  selection.forEach(u => pendingUserIds.add(u.id))
}

const handleSaveUserAssign = async () => {
  savingUsers.value = true
  try {
    await roleApi.bindUsersToRole(roleId, [...pendingUserIds])
    ElMessage.success('用户分配已保存')
    showUserAssignDialog.value = false
    loadRoleUsers()
    loadRoleDetail()
  } catch {
    ElMessage.error('保存用户分配失败')
  } finally {
    savingUsers.value = false
  }
}

const handleRemoveUser = async (userId: string) => {
  try {
    await ElMessageBox.confirm(`确定要移除用户 ${userId}？`, '提示', {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning',
    })
    const remaining = (roleUsers.value.user_ids || []).filter(id => id !== userId)
    await roleApi.bindUsersToRole(roleId, remaining)
    ElMessage.success(t('common.removeSuccess'))
    loadRoleUsers()
    loadRoleDetail()
  } catch (error: any) {
    if (error !== 'cancel') ElMessage.error('移除用户失败')
  }
}

// ════════════════════ 工具函数 ════════════════════

const handleBack = () => router.back()

const userDisplayName = (uid: string) => {
  const u = userDetailMap.value[uid]
  if (!u) return uid
  return u.real_name || u.username
}

const getStatusType = (status: number) =>
  ({ 1: 'success', 2: 'info', 3: 'danger' } as Record<number, string>)[status] || ''

const getStatusText = (status: number) =>
  ({ 1: t('common.active'), 2: t('common.inactive'), 3: t('common.deprecated') } as Record<number, string>)[status] || ''

const getPermTagType = (level: number) =>
  ({ 0: 'info', 1: 'success', 2: 'warning', 3: 'danger' } as Record<number, string>)[level] || ''

const getPermissionText = (level: number) =>
  ({
    0: t('role.permissionLevel.none'),
    1: t('role.permissionLevel.read'),
    2: t('role.permissionLevel.write'),
    3: t('role.permissionLevel.full'),
  } as Record<number, string>)[level] || ''

const formatTimestamp = (ts?: number) => (ts ? new Date(ts).toLocaleString('zh-CN') : '-')

onMounted(() => {
  loadRoleDetail()
  loadMenuTree()
  loadRoleUsers()
})
</script>

<style scoped lang="scss">
.role-detail {
  padding: 20px;

  .page-title {
    font-size: 20px;
    font-weight: 700;
    font-family: 'Cinzel', serif;
    color: #D4AF37;
  }

  .content-row {
    margin-top: 20px;
  }

  .detail-card {
    margin-bottom: 20px;

    :deep(.el-card__header) {
      padding: 14px 20px;
    }

    .card-header {
      display: flex;
      align-items: center;
      gap: 8px;
      color: #D4AF37;
      font-weight: 600;

      > span:first-of-type {
        flex: 1;
      }
    }

    .role-name-text {
      color: var(--c-text-main);
    }
  }

  // 菜单权限树（只读展示）
  .menu-tree-node {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    padding-right: 4px;

    .node-name {
      flex: 1;
      font-size: 13px;
    }
  }

  .menu-card {
    :deep(.el-tree-node__content) {
      height: auto;
      padding: 3px 0;
    }
  }

  // 已绑定用户标签
  .user-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    max-height: 220px;
    overflow-y: auto;
    padding: 4px 0;

    .user-tag {
      font-size: 12px;
      font-family: 'JetBrains Mono', monospace;
    }
  }
}

// ── 菜单权限配置对话框 ────────────────────────────────────────────────
.config-tips {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  margin-bottom: 16px;
  background: rgba(212, 175, 55, 0.08);
  border: 1px solid rgba(212, 175, 55, 0.2);
  border-radius: 8px;
  font-size: 13px;
  color: var(--c-text-sub, #8b9bb4);
}

.menu-config-body {
  max-height: 460px;
  overflow-y: auto;

  .menu-config-header {
    display: flex;
    align-items: center;
    padding: 0 24px 8px 24px;
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.05em;
    text-transform: uppercase;
    color: var(--c-text-sub, #8b9bb4);
    border-bottom: 1px solid var(--c-border, rgba(255, 255, 255, 0.08));
    margin-bottom: 4px;

    span:first-child { flex: 1; }
    span:nth-child(2) { width: 140px; }
    span:last-child { width: 110px; flex-shrink: 0; }
  }

  :deep(.el-tree-node__content) {
    height: auto;
    padding: 2px 0;
  }

  .menu-config-node {
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
    padding: 3px 4px 3px 0;

    .node-label {
      flex: 1;
      font-size: 13px;
      color: var(--c-text-main);
    }

    .node-path {
      width: 140px;
      font-size: 11px;
      color: var(--c-text-sub, #8b9bb4);
      font-family: 'JetBrains Mono', monospace;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      flex-shrink: 0;
    }

    .perm-select {
      width: 110px;
      flex-shrink: 0;

      &.has-perm :deep(.el-input__inner) {
        color: #D4AF37;
        font-weight: 500;
      }
    }
  }
}

// ── 用户分配对话框 ────────────────────────────────────────────────────
.user-assign-toolbar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 14px;

  .assign-count {
    font-size: 13px;
    color: var(--c-text-sub, #8b9bb4);

    strong {
      color: #D4AF37;
      font-size: 16px;
    }
  }
}

.user-assign-table {
  width: 100%;
}
</style>
