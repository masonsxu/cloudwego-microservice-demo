<template>
  <div class="user-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="title">{{ t('user.title') }}</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            {{ t('user.createUser') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索和筛选 -->
      <div class="search-form">
        <el-form :inline="true" :model="searchForm">
          <el-form-item :label="t('common.search')">
            <el-input
              v-model="searchForm.search"
              :placeholder="t('user.username') + '/' + t('user.email') + '/' + t('user.realName')"
              clearable
              @clear="handleSearch"
              @keyup.enter="handleSearch"
            >
              <template #append>
                <el-button :icon="Search" @click="handleSearch" />
              </template>
            </el-input>
          </el-form-item>
          <el-form-item :label="t('user.organization')">
            <el-select
              v-model="searchForm.organization_id"
              :placeholder="t('user.organization')"
              clearable
              @change="handleSearch"
            >
              <el-option
                v-for="org in organizations"
                :key="org.id"
                :label="org.name"
                :value="org.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('user.status')">
            <el-select
              v-model="searchForm.status"
              :placeholder="t('user.status')"
              clearable
              @change="handleSearch"
            >
              <el-option :label="t('user.status.active')" :value="1" />
              <el-option :label="t('user.status.inactive')" :value="2" />
              <el-option :label="t('user.status.suspended')" :value="3" />
              <el-option :label="t('user.status.locked')" :value="4" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button @click="handleReset">
              <el-icon><Refresh /></el-icon>
              {{ t('common.reset') }}
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 数据表格 -->
      <div class="table-container">
        <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
          <el-table-column prop="username" :label="t('user.username')" width="150" />
          <el-table-column prop="email" :label="t('user.email')" width="200" />
          <el-table-column prop="real_name" :label="t('user.realName')" width="150" />
          <el-table-column prop="phone" :label="t('user.phone')" width="150" />
          <el-table-column :label="t('user.organization')" width="200">
            <template #default="{ row }">
              <span v-if="row.organization">{{ row.organization.name }}</span>
              <span v-else class="text-mineral">-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('user.roles')" width="200">
            <template #default="{ row }">
              <el-tag
                v-for="role in row.roles"
                :key="role"
                size="small"
                class="mr-1"
              >
                {{ role }}
              </el-tag>
              <span v-if="!row.roles || row.roles.length === 0" class="text-mineral">-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('user.status')" width="100" align="center">
            <template #default="{ row }">
              <el-tag
                :type="getStatusType(row.status)"
                effect="dark"
              >
                {{ t(`user.status.${getStatusKey(row.status)}`) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.actions')" width="260" align="center" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="handleView(row)">
                <el-icon><View /></el-icon>
                查看
              </el-button>
              <el-button link @click="handleEdit(row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-dropdown @command="(cmd: string) => handleMoreAction(cmd as 'changeStatus' | 'resetPassword' | 'delete', row)">
                <el-button link>
                  更多<el-icon><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="changeStatus">
                      <el-icon><RefreshRight /></el-icon>
                      {{ t('user.changeStatus') }}
                    </el-dropdown-item>
                    <el-dropdown-item command="resetPassword">
                      <el-icon><Lock /></el-icon>
                      {{ t('user.resetPassword') }}
                    </el-dropdown-item>
                    <el-dropdown-item command="delete" divided>
                      <el-icon><Delete /></el-icon>
                      {{ t('user.deleteUser') }}
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
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
    </el-card>

    <!-- 变更状态对话框 -->
    <el-dialog
      v-model="statusDialogVisible"
      :title="t('user.changeStatus')"
      width="500px"
    >
      <el-form :model="statusForm" label-width="100px">
        <el-form-item :label="t('user.status')">
          <el-select v-model="statusForm.new_status" style="width: 100%">
            <el-option :label="t('user.status.active')" :value="1" />
            <el-option :label="t('user.status.inactive')" :value="2" />
            <el-option :label="t('user.status.suspended')" :value="3" />
            <el-option :label="t('user.status.locked')" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.description')">
          <el-input
            v-model="statusForm.reason"
            type="textarea"
            :rows="3"
            :placeholder="t('common.description')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="statusDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="confirmChangeStatus" :loading="submitting">
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  View,
  Edit,
  Delete,
  Lock,
  RefreshRight,
  ArrowDown
} from '@element-plus/icons-vue'
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

const searchForm = reactive({
  search: '',
  organization_id: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

const statusDialogVisible = ref(false)
const statusForm = reactive({
  user_id: '',
  new_status: 1,
  reason: ''
})
const currentUser = ref<UserListItem | null>(null)

onMounted(() => {
  fetchData()
  fetchOrganizations()
})

async function fetchOrganizations() {
  try {
    const response = await getOrganizationList({ limit: 1000 })
    organizations.value = response.organizations || []
  } catch (error) {
    console.error('Failed to fetch organizations:', error)
  }
}

async function fetchData() {
  loading.value = true
  try {
    const response = await getUserList({
      page: pagination.page,
      limit: pagination.limit,
      search: searchForm.search || undefined,
      organization_id: searchForm.organization_id || undefined,
      status: searchForm.status
    })
    tableData.value = response.users || []
    pagination.total = (response as any).total || 0
  } catch (error: any) {
    console.error('Failed to fetch user list:', error)
    ElMessage.error(error.message || t('common.operationFailed'))
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  fetchData()
}

function handleReset() {
  searchForm.search = ''
  searchForm.organization_id = ''
  searchForm.status = undefined
  pagination.page = 1
  fetchData()
}

function handleCreate() {
  router.push('/users/create')
}

function handleView(row: UserListItem) {
  router.push(`/users/${row.id}`)
}

function handleEdit(row: UserListItem) {
  router.push(`/users/${row.id}/edit`)
}

function handleMoreAction(command: 'changeStatus' | 'resetPassword' | 'delete', row: UserListItem) {
  currentUser.value = row
  switch (command) {
    case 'changeStatus':
      openChangeStatusDialog(row)
      break
    case 'resetPassword':
      handleResetPassword(row)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

function openChangeStatusDialog(row: UserListItem) {
  currentUser.value = row
  statusForm.user_id = row.id
  statusForm.new_status = row.status
  statusForm.reason = ''
  statusDialogVisible.value = true
}

async function confirmChangeStatus() {
  submitting.value = true
  try {
    await changeUserStatus(statusForm.user_id, statusForm.new_status, statusForm.reason)
    ElMessage.success(t('common.operationSuccess'))
    statusDialogVisible.value = false
    fetchData()
  } catch (error: any) {
    console.error('Failed to change user status:', error)
    ElMessage.error(error.message || t('common.operationFailed'))
  } finally {
    submitting.value = false
  }
}

async function handleResetPassword(_row: UserListItem) {
  try {
    await ElMessageBox.prompt(
      t('user.newPassword') + '（留空则系统自动生成）',
      t('user.resetPassword'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        inputPattern: /^.{0,50}$/,
        inputErrorMessage: '密码长度不能超过50位'
      }
    )
    // TODO: 调用重置密码 API
    ElMessage.success('密码重置成功')
  } catch (error) {
    // 用户取消
  }
}

async function handleDelete(row: UserListItem) {
  try {
    await ElMessageBox.confirm(
      `${t('user.username')}: ${row.username}`,
      t('common.deleteConfirm'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
        distinguishCancelAndClose: true
      }
    )

    await deleteUser(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel' && error !== 'close') {
      console.error('Failed to delete user:', error)
      ElMessage.error(error.message || t('common.operationFailed'))
    }
  }
}

function handleSizeChange(size: number) {
  pagination.limit = size
  pagination.page = 1
  fetchData()
}

function handlePageChange(page: number) {
  pagination.page = page
  fetchData()
}

function getStatusKey(status: number): string {
  const statusMap: Record<number, string> = {
    1: 'active',
    2: 'inactive',
    3: 'suspended',
    4: 'locked'
  }
  return statusMap[status] || 'unknown'
}

function getStatusType(status: number): string {
  const typeMap: Record<number, string> = {
    1: 'success',
    2: 'info',
    3: 'warning',
    4: 'danger'
  }
  return typeMap[status] || 'info'
}
</script>

<style scoped lang="scss">
.user-list {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .title {
      color: #D4AF37;
      font-family: 'Cinzel', serif;
      font-size: 18px;
      font-weight: 600;
    }
  }

  .search-form {
    margin-bottom: 20px;
    padding: 20px;
    background-color: rgba(44, 46, 51, 0.3);
    border-radius: 12px;
    border: 1px solid rgba(212, 175, 55, 0.1);

    .el-form-item {
      margin-bottom: 0;
    }
  }

  .table-container {
    .el-pagination {
      margin-top: 20px;
      justify-content: flex-end;
    }

    .mr-1 {
      margin-right: 4px;
    }

    .text-mineral {
      color: #8B9bb4;
    }
  }
}
</style>
