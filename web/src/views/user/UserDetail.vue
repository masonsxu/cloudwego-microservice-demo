<template>
  <div class="user-detail">
    <el-page-header @back="handleBack" :title="t('common.back')">
      <template #content>
        <span class="page-title">
          {{ t('user.userDetail') }}
        </span>
      </template>
      <template #extra>
        <el-button-group>
          <el-button @click="handleEdit">
            <el-icon><Edit /></el-icon>
            {{ t('user.editUser') }}
          </el-button>
          <el-dropdown @command="handleAction">
            <el-button>
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
        </el-button-group>
      </template>
    </el-page-header>

    <div class="detail-content" v-loading="loading">
      <el-row :gutter="20">
        <!-- 基本信息 -->
        <el-col :xs="24" :lg="8">
          <el-card class="info-card">
            <template #header>
              <div class="card-header">
                <el-icon><User /></el-icon>
                <span>{{ t('user.userDetail') }}</span>
              </div>
            </template>
            <div class="user-avatar">
              <el-avatar :size="100" :src="userDetail?.avatar">
                {{ userDetail?.username?.charAt(0).toUpperCase() }}
              </el-avatar>
              <el-tag
                :type="getStatusType(userDetail?.status)"
                effect="dark"
                size="large"
                class="status-tag"
              >
                {{ t(`user.status.${getStatusKey(userDetail?.status)}`) }}
              </el-tag>
            </div>
            <el-descriptions :column="1" class="user-info">
              <el-descriptions-item :label="t('user.username')">
                {{ userDetail?.username }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.realName')">
                {{ userDetail?.real_name || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.email')">
                {{ userDetail?.email || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.phone')">
                {{ userDetail?.phone || '-' }}
              </el-descriptions-item>
            </el-descriptions>
          </el-card>
        </el-col>

        <!-- 详细信息 -->
        <el-col :xs="24" :lg="16">
          <el-card class="detail-card">
            <template #header>
              <div class="card-header">
                <el-icon><InfoFilled /></el-icon>
                <span>{{ t('common.description') }}</span>
              </div>
            </template>
            <el-descriptions :column="2" border>
              <el-descriptions-item :label="t('user.firstName')">
                {{ userDetail?.first_name || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.lastName')">
                {{ userDetail?.last_name || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.gender')">
                {{ t(`user.gender.${getGenderKey(userDetail?.gender)}`) }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.employeeId')">
                {{ userDetail?.employee_id || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.professionalTitle')">
                {{ userDetail?.professional_title || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.licenseNumber')">
                {{ userDetail?.license_number || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.specialties')" :span="2">
                <el-tag
                  v-for="specialty in userDetail?.specialties"
                  :key="specialty"
                  class="mr-1"
                >
                  {{ specialty }}
                </el-tag>
                <span v-if="!userDetail?.specialties || userDetail.specialties.length === 0">-</span>
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.organization')">
                {{ userDetail?.organization?.name || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.department')">
                {{ userDetail?.department?.name || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.roles')" :span="2">
                <el-tag
                  v-for="role in userDetail?.role_names"
                  :key="role"
                  type="success"
                  class="mr-1"
                >
                  {{ role }}
                </el-tag>
                <span v-if="!userDetail?.role_names || userDetail.role_names.length === 0">-</span>
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.accountExpiry')">
                {{ formatDate(userDetail?.account_expiry) }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('user.mustChangePassword')">
                <el-tag :type="userDetail?.must_change_password ? 'warning' : 'success'">
                  {{ userDetail?.must_change_password ? t('common.yes') : t('common.no') }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('common.createTime')">
                {{ formatDateTime(userDetail?.created_at) }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('common.updateTime')">
                {{ formatDateTime(userDetail?.updated_at) }}
              </el-descriptions-item>
            </el-descriptions>
          </el-card>

          <!-- 成员关系 -->
          <el-card class="membership-card" v-if="memberships.length > 0">
            <template #header>
              <div class="card-header">
                <el-icon><Connection /></el-icon>
                <span>{{ t('user.memberships') }}</span>
              </div>
            </template>
            <el-timeline>
              <el-timeline-item
                v-for="membership in memberships"
                :key="membership.id"
                :timestamp="formatDateTime(membership.created_at)"
                placement="top"
              >
                <el-card>
                  <h4>{{ membership.organization?.name }}</h4>
                  <p>{{ membership.department?.name || t('user.noDepartment') }}</p>
                  <el-tag v-if="membership.is_primary" type="warning" size="small">
                    {{ t('user.primary') }}
                  </el-tag>
                </el-card>
              </el-timeline-item>
            </el-timeline>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  User,
  Edit,
  Delete,
  Lock,
  RefreshRight,
  ArrowDown,
  InfoFilled,
  Connection
} from '@element-plus/icons-vue'
import { getUserDetail } from '@/api/user'
import { deleteUser } from '@/api/user'
import type { UserProfile } from '@/types/user'
import { formatOptionalTimestamp } from '@/utils/date'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const loading = ref(false)
const userDetail = ref<UserProfile | null>(null)
const memberships = ref<any[]>([])

onMounted(() => {
  fetchUserDetail()
})

async function fetchUserDetail() {
  loading.value = true
  try {
    const userId = route.params.id as string
    const response = await getUserDetail(userId)
    userDetail.value = response.user

    // TODO: 获取成员关系
    // const membershipsResponse = await getUserMemberships(userId)
    // memberships.value = membershipsResponse.memberships
  } catch (error: any) {
    console.error('Failed to fetch user detail:', error)
    ElMessage.error(error.message || t('common.operationFailed'))
  } finally {
    loading.value = false
  }
}

function handleBack() {
  router.back()
}

function handleEdit() {
  router.push(`/users/${route.params.id}/edit`)
}

async function handleAction(command: string) {
  switch (command) {
    case 'delete':
      await handleDelete()
      break
    case 'resetPassword':
      // TODO: 实现重置密码
      ElMessage.info('功能开发中...')
      break
    case 'changeStatus':
      // TODO: 实现变更状态
      ElMessage.info('功能开发中...')
      break
  }
}

async function handleDelete() {
  if (!userDetail.value) return

  try {
    await ElMessageBox.confirm(
      `${t('user.username')}: ${userDetail.value.username}`,
      t('common.deleteConfirm'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )

    await deleteUser(userDetail.value.id)
    ElMessage.success(t('common.deleteSuccess'))
    router.push('/users')
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to delete user:', error)
      ElMessage.error(error.message || t('common.operationFailed'))
    }
  }
}

function getStatusKey(status?: number): string {
  if (!status) return 'unknown'
  const statusMap: Record<number, string> = {
    1: 'active',
    2: 'inactive',
    3: 'suspended',
    4: 'locked'
  }
  return statusMap[status] || 'unknown'
}

function getStatusType(status?: number): string {
  if (!status) return 'info'
  const typeMap: Record<number, string> = {
    1: 'success',
    2: 'info',
    3: 'warning',
    4: 'danger'
  }
  return typeMap[status] || 'info'
}

function getGenderKey(gender?: number): string {
  if (!gender) return 'unknown'
  const genderMap: Record<number, string> = {
    0: 'unknown',
    1: 'male',
    2: 'female'
  }
  return genderMap[gender] || 'unknown'
}

function formatDate(timestamp?: number): string {
  // 后端返回的是毫秒时间戳，直接使用工具函数
  return formatOptionalTimestamp(timestamp, 'YYYY-MM-DD')
}

function formatDateTime(timestamp?: number): string {
  // 后端返回的是毫秒时间戳，直接使用工具函数
  return formatOptionalTimestamp(timestamp, 'YYYY-MM-DD HH:mm:ss')
}
</script>

<style scoped lang="scss">
.user-detail {
  .page-title {
    font-size: 18px;
    font-weight: 600;
    color: #D4AF37;
    font-family: 'Cinzel', serif;
  }

  .detail-content {
    margin-top: 20px;

    .info-card,
    .detail-card,
    .membership-card {
      background: linear-gradient(145deg, rgba(30, 32, 36, 0.9), rgba(20, 20, 22, 0.95));
      border: 1px solid rgba(255, 255, 255, 0.05);
      border-radius: 20px;
      box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5);
      margin-bottom: 20px;

      .card-header {
        display: flex;
        align-items: center;
        gap: 10px;
        color: #D4AF37;
        font-family: 'Cinzel', serif;
        font-size: 16px;
        font-weight: 600;
      }

      :deep(.el-card__header) {
        border-bottom: 1px solid rgba(212, 175, 55, 0.2);
      }
    }

    .user-avatar {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 20px;
      margin-bottom: 30px;

      .el-avatar {
        font-size: 40px;
        font-weight: 600;
        background: linear-gradient(135deg, #D4AF37 0%, #C4A963 100%);
      }

      .status-tag {
        font-size: 14px;
      }
    }

    .user-info {
      width: 100%;

      :deep(.el-descriptions__label) {
        color: #8B9bb4;
      }

      :deep(.el-descriptions__content) {
        color: #F2F0E4;
      }
    }

    .mr-1 {
      margin-right: 4px;
    }
  }
}
</style>
