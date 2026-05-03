<template>
  <div class="user-detail">
    <div class="page-header">
      <button class="back-btn" @click="handleBack">
        <ArrowLeft class="h-4 w-4" />
        {{ t('common.back') }}
      </button>
      <div class="header-content">
        <h1 class="page-title">{{ t('user.userDetail') }}</h1>
      </div>
      <div class="header-actions">
        <Button variant="outline" @click="handleEdit">
          <Pencil class="mr-1 h-4 w-4" />
          {{ t('user.editUser') }}
        </Button>
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button variant="outline">
              {{ t('common.actions') }}<ChevronDown class="ml-1 h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem @select="handleAction('changeStatus')">
              <RefreshCw class="mr-2 h-4 w-4" />
              {{ t('user.changeStatus') }}
            </DropdownMenuItem>
            <DropdownMenuItem @select="handleAction('resetPassword')">
              <Lock class="mr-2 h-4 w-4" />
              {{ t('user.resetPassword') }}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem class="text-destructive" @select="handleAction('delete')">
              <Trash2 class="mr-2 h-4 w-4" />
              {{ t('user.deleteUser') }}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>

    <div class="detail-content">
      <DetailPageSkeleton
        v-if="initialLoading"
        :side-span="8"
        :main-span="16"
        :side-cards="1"
        :main-cards="1"
        :show-avatar="true"
        :side-item-counts="[4]"
        :main-item-counts="[12]"
      />
      <div v-else class="detail-grid">
        <div class="info-card">
          <div class="card-header">
            <User class="h-5 w-5" />
            <span>{{ t('user.userDetail') }}</span>
          </div>
          <div class="user-avatar">
            <Avatar class="avatar-size">
              <AvatarImage :src="userDetail?.avatar ?? ''" />
              <AvatarFallback>
                {{ userDetail?.username?.charAt(0).toUpperCase() }}
              </AvatarFallback>
            </Avatar>
            <Badge :class="getStatusBadgeClass(userDetail?.status)" class="status-tag">
              {{ t(`user.status.${getStatusKey(userDetail?.status)}`) }}
            </Badge>
          </div>
          <div class="user-info">
            <div class="info-row">
              <span class="info-label">{{ t('user.username') }}</span>
              <span class="info-value">{{ userDetail?.username }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">{{ t('user.realName') }}</span>
              <span class="info-value">{{ userDetail?.real_name || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">{{ t('user.email') }}</span>
              <span class="info-value">{{ userDetail?.email || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">{{ t('user.phone') }}</span>
              <span class="info-value">{{ userDetail?.phone || '-' }}</span>
            </div>
          </div>
        </div>

        <div class="detail-card">
          <div class="card-header">
            <Info class="h-5 w-5" />
            <span>{{ t('common.descriptions') }}</span>
          </div>
          <div class="detail-grid-2col">
            <div class="detail-item">
              <span class="detail-label">{{ t('user.firstName') }}</span>
              <span class="detail-value">{{ userDetail?.first_name || '-' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">{{ t('user.lastName') }}</span>
              <span class="detail-value">{{ userDetail?.last_name || '-' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">{{ t('user.gender.label') }}</span>
              <span class="detail-value">{{ t(`user.gender.${getGenderKey(userDetail?.gender)}`) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">{{ t('user.employeeId') }}</span>
              <span class="detail-value">{{ userDetail?.employee_id || '-' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">{{ t('user.professionalTitle') }}</span>
              <span class="detail-value">{{ userDetail?.professional_title || '-' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">{{ t('user.organization') }}</span>
              <span class="detail-value">{{ userDetail?.organization?.name || '-' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">{{ t('user.department') }}</span>
              <span class="detail-value">{{ userDetail?.department?.name || '-' }}</span>
            </div>
            <div class="detail-item col-span-2">
              <span class="detail-label">{{ t('user.roles') }}</span>
              <span class="detail-value">
                <Badge v-for="role in userDetail?.role_names" :key="role" variant="secondary" class="mr-1">
                  {{ role }}
                </Badge>
                <span v-if="!userDetail?.role_names || userDetail.role_names.length === 0">-</span>
              </span>
            </div>
            <div class="detail-item">
              <span class="detail-label">{{ t('user.accountExpiry') }}</span>
              <span class="detail-value">{{ formatDate(userDetail?.account_expiry) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">{{ t('user.mustChangePassword') }}</span>
              <span class="detail-value">
                <Badge :variant="userDetail?.must_change_password ? 'default' : 'secondary'">
                  {{ userDetail?.must_change_password ? t('common.yes') : t('common.no') }}
                </Badge>
              </span>
            </div>
            <div class="detail-item">
              <span class="detail-label">{{ t('common.createTime') }}</span>
              <span class="detail-value">{{ formatDateTime(userDetail?.created_at) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">{{ t('common.updateTime') }}</span>
              <span class="detail-value">{{ formatDateTime(userDetail?.updated_at) }}</span>
            </div>
          </div>
        </div>

        <div class="membership-card" v-if="memberships.length > 0">
          <div class="card-header">
            <Link class="h-5 w-5" />
            <span>{{ t('user.memberships') }}</span>
          </div>
          <div class="timeline">
            <div v-for="membership in memberships" :key="membership.id" class="timeline-item">
              <div class="timeline-dot"></div>
              <div class="timeline-content">
                <div class="timeline-timestamp">{{ formatDateTime(membership.created_at) }}</div>
                <div class="timeline-card">
                  <h4>{{ membership.organization?.name }}</h4>
                  <p>{{ membership.department?.name || t('user.noDepartment') }}</p>
                  <Badge v-if="membership.is_primary" variant="outline" class="mt-2">
                    {{ t('user.primary') }}
                  </Badge>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import {
  User, Pencil, Trash2, Lock, RefreshCw, ChevronDown,
  Info, Link, ArrowLeft
} from 'lucide-vue-next'
import { getUserDetail } from '@/api/user'
import { deleteUser } from '@/api/user'
import { organizationApi } from '@/api/organization'
import { departmentApi } from '@/api/department'
import type { UserProfile } from '@/types/user'
import { formatOptionalTimestamp } from '@/utils/date'
import DetailPageSkeleton from '@/components/skeleton/DetailPageSkeleton.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import {
  DropdownMenu, DropdownMenuTrigger, DropdownMenuContent,
  DropdownMenuItem, DropdownMenuSeparator
} from '@/components/ui/dropdown-menu'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const initialLoading = ref(true)
const userDetail = ref<UserProfile | null>(null)
const memberships = ref<any[]>([])

onMounted(async () => {
  try {
    await fetchUserDetail()
  } finally {
    initialLoading.value = false
  }
})

async function fetchUserDetail() {
  try {
    const userId = route.params.id as string
    const response = await getUserDetail(userId)
    const user = response.user

    if (user.primary_organization_id) {
      try {
        const orgResponse = await organizationApi.getOrganization(user.primary_organization_id)
        user.organization = {
          id: orgResponse.organization.id,
          name: orgResponse.organization.name
        }
      } catch (e) {
        console.error('Failed to fetch organization:', e)
      }
    }

    if (user.primary_department_id) {
      try {
        const deptResponse = await departmentApi.getDepartment(user.primary_department_id)
        user.department = {
          id: deptResponse.department.id,
          name: deptResponse.department.name
        }
      } catch (e) {
        console.error('Failed to fetch department:', e)
      }
    }

    userDetail.value = user
  } catch (error: any) {
    console.error('Failed to fetch user detail:', error)
    toast.error(error.message || t('common.operationFailed'))
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
      toast.info('功能开发中...')
      break
    case 'changeStatus':
      toast.info('功能开发中...')
      break
  }
}

async function handleDelete() {
  if (!userDetail.value) return

  try {
    if (confirm(`${t('user.username')}: ${userDetail.value.username}\n${t('common.deleteConfirm')}`)) {
      await deleteUser(userDetail.value.id)
      toast.success(t('common.deleteSuccess'))
      router.push('/users')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to delete user:', error)
      toast.error(error.message || t('common.operationFailed'))
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

function getStatusBadgeClass(status?: number): string {
  if (!status) return 'bg-gray-400 text-white'
  const classMap: Record<number, string> = {
    1: 'bg-emerald-500 text-white',
    2: 'bg-gray-400 text-white',
    3: 'bg-amber-500 text-white',
    4: 'bg-red-500 text-white'
  }
  return classMap[status] || 'bg-gray-400 text-white'
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
  return formatOptionalTimestamp(timestamp, 'YYYY-MM-DD')
}

function formatDateTime(timestamp?: number): string {
  return formatOptionalTimestamp(timestamp, 'YYYY-MM-DD HH:mm:ss')
}
</script>

<style scoped lang="scss">
.user-detail {
  padding: 20px;

  .page-header {
    display: flex;
    align-items: center;
    gap: 16px;

    .back-btn {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 8px 12px;
      border: 1px solid hsl(var(--border));
      border-radius: 8px;
      background: var(--bg-card);
      color: var(--c-text-sub);
      font-size: 14px;
      cursor: pointer;
      transition: all 0.2s ease;

      &:hover {
        background: var(--bg-input);
        color: var(--c-text-main);
      }
    }

    .header-content {
      flex: 1;

      .page-title {
        font-size: 20px;
        font-weight: 600;
        color: var(--color-ink-strong);
        margin: 0;
      }
    }

    .header-actions {
      display: flex;
      gap: 8px;
    }
  }

  .detail-content {
    margin-top: 20px;

    .detail-grid {
      display: grid;
      grid-template-columns: 1fr 2fr;
      gap: 20px;
    }

    .info-card,
    .detail-card,
    .membership-card {
      background: var(--color-canvas);
      border: 1px solid var(--color-border-subtle);
      border-radius: 8px;

      .card-header {
        display: flex;
        align-items: center;
        gap: 10px;
        color: var(--color-ink-strong);
        font-size: 14px;
        font-weight: 600;
        padding: 14px 20px;
        border-bottom: 1px solid var(--color-border-subtle);
      }
    }

    .info-card {
      padding-bottom: 20px;

      .user-avatar {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 20px;
        margin-bottom: 30px;
        padding-top: 20px;

        .avatar-size {
          width: 100px;
          height: 100px;
          font-size: 40px;
          font-weight: 600;

          :deep(.avatar-fallback) {
            background: var(--color-primary-soft-strong);
            color: var(--color-primary-active);
          }
        }

        .status-tag {
          font-size: 14px;
        }
      }

      .user-info {
        padding: 0 20px;

        .info-row {
          display: flex;
          justify-content: space-between;
          padding: 10px 0;
          border-bottom: 1px solid hsl(var(--border) / 0.3);

          &:last-child {
            border-bottom: none;
          }

          .info-label {
            color: var(--c-text-sub);
            font-size: 13px;
          }

          .info-value {
            color: var(--c-text-main);
            font-size: 13px;
          }
        }
      }
    }

    .detail-card {
      padding-bottom: 20px;

      .detail-grid-2col {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 16px;
        padding: 20px;

        .detail-item {
          display: flex;
          flex-direction: column;
          gap: 4px;

          &.col-span-2 {
            grid-column: span 2;
          }

          .detail-label {
            color: var(--c-text-sub);
            font-size: 12px;
          }

          .detail-value {
            color: var(--c-text-main);
            font-size: 13px;
          }
        }
      }
    }

    .membership-card {
      margin-top: 20px;
      padding-bottom: 20px;

      .timeline {
        padding: 20px;

        .timeline-item {
          position: relative;
          padding-left: 24px;
          padding-bottom: 20px;

          &:last-child {
            padding-bottom: 0;
          }

          .timeline-dot {
            position: absolute;
            left: 0;
            top: 4px;
            width: 10px;
            height: 10px;
            border-radius: 50%;
            background: var(--c-primary);
          }

          .timeline-content {
            .timeline-timestamp {
              font-size: 12px;
              color: var(--c-text-muted);
              margin-bottom: 8px;
            }

            .timeline-card {
              padding: 12px 16px;
              background: var(--color-sunken);
              border: 1px solid var(--color-border-subtle);
              border-radius: 8px;

              h4 {
                margin: 0 0 4px;
                font-size: 14px;
                font-weight: 600;
                color: var(--color-ink);
              }

              p {
                margin: 0;
                font-size: 13px;
                color: var(--color-ink-muted);
              }
            }
          }
        }
      }
    }
  }
}
</style>
