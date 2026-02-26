<template>
  <div class="user-list">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('user.title') }}</h1>
        <p class="page-subtitle">管理系统用户账号与访问权限</p>
      </div>
      <button class="create-btn shimmer-btn" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        {{ t('user.createUser') }}
      </button>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-row">
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon"><el-icon size="20"><User /></el-icon></div>
        <div class="stat-info">
          <span class="stat-value">{{ pagination.total }}</span>
          <span class="stat-label">{{ t('user.title') }}</span>
        </div>
      </div>
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon active-icon"><el-icon size="20"><CircleCheck /></el-icon></div>
        <div class="stat-info">
          <span class="stat-value">{{ tableData.filter(u => u.status === 1).length }}</span>
          <span class="stat-label">{{ t('user.status.active') }}</span>
        </div>
      </div>
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon locked-icon"><el-icon size="20"><Lock /></el-icon></div>
        <div class="stat-info">
          <span class="stat-value">{{ tableData.filter(u => u.status === 4).length }}</span>
          <span class="stat-label">{{ t('user.status.locked') }}</span>
        </div>
      </div>
    </div>

    <!-- 搜索筛选 -->
    <div class="search-section">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item>
          <el-input
            v-model="searchForm.search"
            :placeholder="t('user.username') + ' / ' + t('user.email') + ' / ' + t('user.realName')"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
            style="width: 320px"
          >
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchForm.organization_id" :placeholder="t('user.organization')" clearable @change="handleSearch" style="width: 180px">
            <el-option v-for="org in organizations" :key="org.id" :label="org.name" :value="org.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchForm.status" :placeholder="t('common.status')" clearable @change="handleSearch" style="width: 130px">
            <el-option :label="t('user.status.active')" :value="1" />
            <el-option :label="t('user.status.inactive')" :value="2" />
            <el-option :label="t('user.status.suspended')" :value="3" />
            <el-option :label="t('user.status.locked')" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>{{ t('common.reset') }}
          </el-button>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>{{ t('common.search') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 数据表格 -->
    <div class="table-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
      <el-table :data="tableData" v-loading="loading" class="modern-table" style="width: 100%">
        <el-table-column prop="username" :label="t('user.username')" width="150">
          <template #default="{ row }">
            <div class="user-cell">
              <div class="user-avatar-mini">{{ row.username?.charAt(0).toUpperCase() }}</div>
              <span class="username">{{ row.username }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="real_name" :label="t('user.realName')" width="120">
          <template #default="{ row }">
            <span class="text-sub">{{ row.real_name || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="email" :label="t('user.email')" min-width="180">
          <template #default="{ row }">
            <span class="text-sub">{{ row.email || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('user.organization')" min-width="160">
          <template #default="{ row }">
            <span v-if="row.organization" class="text-sub">{{ row.organization.name }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('user.roles')" min-width="160">
          <template #default="{ row }">
            <div class="roles-cell">
              <el-tag v-for="role in (row.roles || []).slice(0, 2)" :key="role" size="small" class="role-tag">{{ role }}</el-tag>
              <span v-if="(row.roles || []).length > 2" class="more-tag">+{{ row.roles.length - 2 }}</span>
              <span v-if="!row.roles || row.roles.length === 0" class="text-muted">-</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.status')" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small" effect="dark">
              {{ t(`user.status.${getStatusKey(row.status)}`) }}
            </el-tag>
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
              <el-dropdown @command="(cmd: string) => handleMoreAction(cmd as any, row)" trigger="click">
                <button class="action-btn more-btn" :title="t('common.actions')">
                  <el-icon><MoreFilled /></el-icon>
                </button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="changeStatus">
                      <el-icon><RefreshRight /></el-icon>{{ t('user.changeStatus') }}
                    </el-dropdown-item>
                    <el-dropdown-item command="resetPassword">
                      <el-icon><Lock /></el-icon>{{ t('user.resetPassword') }}
                    </el-dropdown-item>
                    <el-dropdown-item command="delete" divided>
                      <el-icon><Delete /></el-icon>{{ t('user.deleteUser') }}
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.limit"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          :background="true"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 变更状态对话框 -->
    <el-dialog v-model="statusDialogVisible" :title="t('user.changeStatus')" width="480px">
      <el-form :model="statusForm" label-width="100px">
        <el-form-item :label="t('common.status')">
          <el-select v-model="statusForm.new_status" style="width: 100%">
            <el-option :label="t('user.status.active')" :value="1" />
            <el-option :label="t('user.status.inactive')" :value="2" />
            <el-option :label="t('user.status.suspended')" :value="3" />
            <el-option :label="t('user.status.locked')" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.description')">
          <el-input v-model="statusForm.reason" type="textarea" :rows="3" :placeholder="t('common.description')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="statusDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="confirmChangeStatus" :loading="submitting">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, View, Edit, Delete, Lock, RefreshRight, User, CircleCheck, MoreFilled } from '@element-plus/icons-vue'
import { getUserList, changeUserStatus, deleteUser } from '@/api/user'
import { getOrganizationList } from '@/api/organization'
import type { UserListItem } from '@/types/user'
import type { Organization } from '@/types/organization'

const router = useRouter()
const { t } = useI18n()

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<UserListItem[]>([])
const organizations = ref<Organization[]>([])

const searchForm = reactive({ search: '', organization_id: '', status: undefined as number | undefined })
const pagination = reactive({ page: 1, limit: 20, total: 0 })
const statusDialogVisible = ref(false)
const statusForm = reactive({ user_id: '', new_status: 1, reason: '' })
const currentUser = ref<UserListItem | null>(null)

onMounted(() => { fetchData(); fetchOrganizations() })

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
      status: searchForm.status
    })
    tableData.value = response.users || []
    pagination.total = (response as any).total || 0
  } catch (error: any) {
    ElMessage.error(error.message || t('common.operationFailed'))
  } finally {
    loading.value = false
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
    ElMessage.success(t('common.operationSuccess')); statusDialogVisible.value = false; fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || t('common.operationFailed'))
  } finally { submitting.value = false }
}

async function handleResetPassword(_row: UserListItem) {
  try {
    await ElMessageBox.prompt(t('user.newPassword') + '（留空则系统自动生成）', t('user.resetPassword'), {
      confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'),
      inputPattern: /^.{0,50}$/, inputErrorMessage: '密码长度不能超过50位'
    })
    ElMessage.success('密码重置成功')
  } catch {}
}

async function handleDelete(row: UserListItem) {
  try {
    await ElMessageBox.confirm(`${t('user.username')}: ${row.username}`, t('common.deleteConfirm'), {
      confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning', distinguishCancelAndClose: true
    })
    await deleteUser(row.id); ElMessage.success(t('common.deleteSuccess')); fetchData()
  } catch (error: any) {
    if (error !== 'cancel' && error !== 'close') ElMessage.error(error.message || t('common.operationFailed'))
  }
}

function handleSizeChange(size: number) { pagination.limit = size; pagination.page = 1; fetchData() }
function handlePageChange(page: number) { pagination.page = page; fetchData() }
function getStatusKey(status: number) { return ({ 1: 'active', 2: 'inactive', 3: 'suspended', 4: 'locked' } as Record<number, string>)[status] || 'unknown' }
function getStatusType(status: number) { return ({ 1: 'success', 2: 'info', 3: 'warning', 4: 'danger' } as Record<number, string>)[status] || 'info' }
</script>

<style scoped lang="scss">
.user-list {
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
      .page-subtitle { color: var(--c-text-sub); font-size: 13px; margin: 0; }
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

      &:hover { box-shadow: 0 0 24px rgba(212, 175, 55, 0.4); transform: translateY(-1px); }
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

        &.active-icon { background: rgba(103, 194, 58, 0.1); border-color: rgba(103, 194, 58, 0.25); color: #67C23A; }
        &.locked-icon { background: rgba(245, 108, 108, 0.1); border-color: rgba(245, 108, 108, 0.25); color: #F56C6C; }
      }

      .stat-info {
        display: flex;
        flex-direction: column;
        gap: 2px;
        .stat-value { font-size: 24px; font-weight: 700; color: var(--c-text-main); font-family: 'JetBrains Mono', monospace; line-height: 1; }
        .stat-label { font-size: 12px; color: var(--c-text-sub); }
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
      :deep(td.el-table__cell) { padding: 12px; }
    }

    .user-cell {
      display: flex;
      align-items: center;
      gap: 10px;

      .user-avatar-mini {
        width: 30px;
        height: 30px;
        border-radius: 8px;
        background: linear-gradient(135deg, #D4AF37 0%, #C4A963 100%);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 13px;
        font-weight: 700;
        color: #000;
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
        padding: 1px 8px;
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
        font-size: 14px;

        &:hover { transform: translateY(-1px); }
        &.view-btn:hover { color: var(--c-accent); border-color: var(--c-border-accent); background: rgba(212, 175, 55, 0.08); }
        &.edit-btn:hover { color: #E6A23C; border-color: rgba(230, 162, 60, 0.4); background: rgba(230, 162, 60, 0.08); }
        &.more-btn:hover { color: var(--c-text-main); border-color: var(--c-border-accent); background: var(--bg-input); }
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
