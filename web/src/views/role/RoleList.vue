<template>
  <div class="role-list">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('role.title') }}</h1>
        <p class="page-subtitle">配置系统角色与权限策略</p>
      </div>
      <button class="create-btn shimmer-btn" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        {{ t('role.createRole') }}
      </button>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-row">
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon"><el-icon size="20"><Key /></el-icon></div>
        <div class="stat-info">
          <span class="stat-value">{{ pagination.total }}</span>
          <span class="stat-label">{{ t('role.title') }}</span>
        </div>
      </div>
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon system-icon"><el-icon size="20"><Lock /></el-icon></div>
        <div class="stat-info">
          <span class="stat-value">{{ roleList.filter(r => r.is_system_role).length }}</span>
          <span class="stat-label">{{ t('role.isSystemRole') }}</span>
        </div>
      </div>
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon active-icon"><el-icon size="20"><User /></el-icon></div>
        <div class="stat-info">
          <span class="stat-value">{{ roleList.reduce((s, r) => s + (r.user_count || 0), 0) }}</span>
          <span class="stat-label">已分配用户</span>
        </div>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="search-section">
      <el-form :inline="true" class="search-form">
        <el-form-item>
          <el-input
            v-model="searchForm.name"
            :placeholder="t('role.roleName')"
            clearable
            @clear="handleSearch"
            style="width: 240px"
          >
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchForm.status" :placeholder="t('common.status')" clearable @clear="handleSearch" style="width: 140px">
            <el-option :label="t('common.active')" :value="1" />
            <el-option :label="t('common.inactive')" :value="2" />
            <el-option :label="t('common.deprecated')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchForm.isSystemRole" :placeholder="t('role.isSystemRole')" clearable @clear="handleSearch" style="width: 140px">
            <el-option :label="t('common.yes')" :value="true" />
            <el-option :label="t('common.no')" :value="false" />
          </el-select>
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

    <!-- 角色列表 -->
    <div class="table-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
      <el-table v-loading="loading" :data="roleList" class="modern-table" style="width: 100%">
        <el-table-column prop="name" :label="t('role.roleName')" min-width="180">
          <template #default="{ row }">
            <div class="role-name-cell">
              <div class="role-icon" :class="{ 'system-role': row.is_system_role }">
                <el-icon size="13"><Key /></el-icon>
              </div>
              <span class="role-name">{{ row.name }}</span>
              <el-tag v-if="row.is_system_role" size="small" class="sys-badge">系统</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="description" :label="t('role.description')" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="text-sub">{{ row.description || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.status')" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user_count" :label="t('role.userCount')" width="100" align="center">
          <template #default="{ row }">
            <span class="count-badge">{{ row.user_count || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.createTime')" width="160">
          <template #default="{ row }">
            <span class="text-sub time-text">{{ formatTimestamp(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="160" align="right" fixed="right">
          <template #default="{ row }">
            <div class="action-group">
              <button class="action-btn view-btn" @click="handleView(row)" :title="t('common.view')">
                <el-icon><View /></el-icon>
              </button>
              <button class="action-btn edit-btn" @click="handleEdit(row)" :title="t('common.edit')">
                <el-icon><Edit /></el-icon>
              </button>
              <button v-if="!row.is_system_role" class="action-btn delete-btn" @click="handleDelete(row)" :title="t('common.delete')">
                <el-icon><Delete /></el-icon>
              </button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadRoles"
          @current-change="loadRoles"
        />
      </div>
    </div>

    <!-- 创建/编辑角色对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="editingRole ? t('role.editRole') : t('role.createRole')"
      width="560px"
      @close="handleDialogClose"
    >
      <el-form ref="roleFormRef" :model="roleForm" :rules="roleFormRules" label-width="120px">
        <el-form-item :label="t('role.roleName')" prop="name">
          <el-input v-model="roleForm.name" :placeholder="t('role.roleName')" :disabled="!!editingRole" />
        </el-form-item>
        <el-form-item :label="t('role.description')" prop="description">
          <el-input v-model="roleForm.description" type="textarea" :rows="3" :placeholder="t('role.description')" />
        </el-form-item>
        <el-form-item v-if="editingRole" :label="t('common.status')" prop="status">
          <el-select v-model="roleForm.status" style="width: 100%">
            <el-option :label="t('common.active')" :value="1" />
            <el-option :label="t('common.inactive')" :value="2" />
            <el-option :label="t('common.deprecated')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('role.isSystemRole')" prop="is_system_role">
          <el-switch v-model="roleForm.is_system_role" :disabled="!!editingRole" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveRole">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Search, RefreshLeft, Key, Lock, User, View, Edit, Delete } from '@element-plus/icons-vue'
import { roleApi } from '@/api/role'
import type { RoleDefinitionDTO } from '@/api/role'

const { t } = useI18n()
const router = useRouter()

const loading = ref(false)
const roleList = ref<RoleDefinitionDTO[]>([])
const showCreateDialog = ref(false)
const editingRole = ref<RoleDefinitionDTO | null>(null)
const roleFormRef = ref<FormInstance>()

const searchForm = reactive({
  name: '',
  status: undefined as number | undefined,
  isSystemRole: undefined as boolean | undefined,
})

const pagination = reactive({ page: 1, size: 20, total: 0 })

const roleForm = reactive({
  name: '', description: '', status: 1, is_system_role: false, permissions: [],
})

const roleFormRules: FormRules = {
  name: [
    { required: true, message: t('role.roleName') + '不能为空', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' },
  ],
  description: [{ max: 500, message: '描述长度不能超过 500 个字符', trigger: 'blur' }],
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

const loadRoles = async () => {
  loading.value = true
  try {
    const response = await roleApi.getRoles({
      name: searchForm.name || undefined,
      status: searchForm.status,
      is_system_role: searchForm.isSystemRole,
      page: pagination.page,
      limit: pagination.size,
    })
    roleList.value = response.roles || []
    pagination.total = response.page?.total || 0
  } catch (error) {
    ElMessage.error('获取角色列表失败')
  } finally {
    loading.value = false
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
    await ElMessageBox.confirm(`确定要删除角色 "${role.name}" 吗？`, '提示', {
      confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning',
    })
    await roleApi.deleteRole(role.id!)
    ElMessage.success('删除成功'); loadRoles()
  } catch (error: any) {
    if (error !== 'cancel') ElMessage.error('删除角色失败')
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
      ElMessage.success(editingRole.value ? '更新成功' : '创建成功')
      showCreateDialog.value = false; loadRoles()
    }
  } catch {}
}

const handleDialogClose = () => { editingRole.value = null; roleFormRef.value?.resetFields() }

const getStatusType = (status: number) => ({ 1: 'success', 2: 'info', 3: 'danger' }[status] || '')
const getStatusText = (status: number) => {
  const map: Record<number, string> = { 1: t('common.active'), 2: t('common.inactive'), 3: t('common.deprecated') }
  return map[status] || ''
}
const formatTimestamp = (ts?: number) => ts ? new Date(ts).toLocaleDateString('zh-CN') : '-'

onMounted(() => { loadRoles() })
</script>

<style scoped lang="scss">
.role-list {
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

    .create-btn {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 10px 20px;
      background: linear-gradient(135deg, #D4AF37 0%, #C4A033 100%);
      color: #000;
      font-weight: 600;
      font-size: 14px;
      border: none;
      border-radius: 10px;
      cursor: pointer;
      transition: all 0.3s ease;

      &:hover {
        box-shadow: 0 0 24px rgba(212, 175, 55, 0.4);
        transform: translateY(-1px);
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

        &.system-icon {
          background: rgba(245, 108, 108, 0.1);
          border-color: rgba(245, 108, 108, 0.25);
          color: #F56C6C;
        }

        &.active-icon {
          background: rgba(103, 194, 58, 0.1);
          border-color: rgba(103, 194, 58, 0.25);
          color: #67C23A;
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
    background: var(--bg-card);
    border: 1px solid var(--c-border-accent);
    border-radius: 14px;
    overflow: hidden;
    box-shadow: var(--shadow-card);

    .modern-table {
      :deep(th.el-table__cell) { padding: 14px 12px; font-size: 12px; letter-spacing: 0.05em; text-transform: uppercase; }
      :deep(td.el-table__cell) { padding: 14px 12px; }
    }

    .role-name-cell {
      display: flex;
      align-items: center;
      gap: 10px;

      .role-icon {
        width: 28px;
        height: 28px;
        border-radius: 6px;
        background: rgba(212, 175, 55, 0.1);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--c-accent);
        flex-shrink: 0;

        &.system-role {
          background: rgba(245, 108, 108, 0.1);
          color: #F56C6C;
        }
      }

      .role-name { font-weight: 500; color: var(--c-text-main); }

      .sys-badge {
        background: rgba(245, 108, 108, 0.1);
        border-color: rgba(245, 108, 108, 0.3);
        color: #F56C6C;
        font-size: 10px;
        padding: 1px 6px;
      }
    }

    .text-sub { color: var(--c-text-sub); font-size: 13px; }
    .time-text { font-size: 12px; }

    .count-badge {
      display: inline-block;
      padding: 2px 10px;
      background: rgba(212, 175, 55, 0.08);
      border-radius: 9999px;
      font-size: 13px;
      font-weight: 600;
      color: var(--c-text-main);
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
        &.view-btn:hover { color: var(--c-accent); border-color: var(--c-border-accent); background: rgba(212, 175, 55, 0.08); }
        &.edit-btn:hover { color: #E6A23C; border-color: rgba(230, 162, 60, 0.4); background: rgba(230, 162, 60, 0.08); }
        &.delete-btn:hover { color: #F56C6C; border-color: rgba(245, 108, 108, 0.4); background: rgba(245, 108, 108, 0.08); }
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
