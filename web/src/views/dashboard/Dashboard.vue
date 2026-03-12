<template>
  <div class="dashboard" v-loading="loading">
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon user-icon">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.userCount }}</p>
              <p class="stat-label">{{ t('dashboard.userCount') }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon org-icon">
              <el-icon><OfficeBuilding /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.orgCount }}</p>
              <p class="stat-label">{{ t('dashboard.orgCount') }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon role-icon">
              <el-icon><Key /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.roleCount }}</p>
              <p class="stat-label">{{ t('dashboard.roleCount') }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon activity-icon">
              <el-icon><TrendCharts /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.activityCount }}</p>
              <p class="stat-label">{{ t('dashboard.recentActivity') }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="content-row">
      <el-col :xs="24" :lg="16">
        <el-card class="welcome-card">
          <template #header>
            <div class="card-header">
              <span class="title">{{ t('dashboard.welcome') }}</span>
              <span class="meta">Overview</span>
            </div>
          </template>
          <div class="welcome-content">
            <p class="greeting">{{ greeting }}</p>
            <p class="username">{{ authStore.username }}</p>
            <p class="description">
              {{ t('dashboard.welcomeDescription') }}
            </p>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="8">
        <el-card class="info-card">
          <template #header>
            <div class="card-header">
              <span class="title">系统信息</span>
              <span class="meta">System</span>
            </div>
          </template>
          <div class="info-content">
            <div class="info-item">
              <span class="label">框架版本</span>
              <span class="value">Vue 3.4+</span>
            </div>
            <div class="info-item">
              <span class="label">UI 组件</span>
              <span class="value">Element Plus</span>
            </div>
            <div class="info-item">
              <span class="label">状态管理</span>
              <span class="value">Pinia</span>
            </div>
            <div class="info-item">
              <span class="label">主题</span>
              <span class="value">CareerCompass Inspired</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { getDashboardStats } from '@/api/dashboard'
import { ElMessage } from 'element-plus'

const authStore = useAuthStore()
const { t } = useI18n()

const loading = ref(true)
const stats = ref({
  userCount: 0,
  orgCount: 0,
  roleCount: 0,
  activityCount: 0
})

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return t('dashboard.greetingMorning')
  if (hour < 18) return t('dashboard.greetingAfternoon')
  return t('dashboard.greetingEvening')
})

async function fetchStats() {
  loading.value = true
  try {
    const data = await getDashboardStats()
    stats.value = data
  } catch (error: any) {
    console.error('Failed to fetch dashboard stats:', error)
    ElMessage.error(error.message || t('common.operationFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped lang="scss">
.dashboard {
  .stats-row,
  .content-row {
    position: relative;
    z-index: 1;
  }

  .stats-row {
    margin-bottom: 20px;

    .stat-card {
      background: var(--bg-card);
      border: 1px solid hsl(var(--border) / 0.6);
      border-radius: 18px;
      box-shadow: var(--shadow-card);
      backdrop-filter: blur(12px);
      position: relative;
      overflow: hidden;

      :deep(.el-card__body) {
        padding: 30px;
      }

      .stat-content {
        display: flex;
        align-items: center;
        gap: 20px;

        .stat-icon {
          width: 60px;
          height: 60px;
          border-radius: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 28px;
          color: #ffffff;
          box-shadow: inset 0 0 0 1px hsl(var(--border) / 0.4);

          &.user-icon {
            background: linear-gradient(135deg, var(--c-primary) 0%, #6f7fe0 100%);
          }

          &.org-icon {
            background: linear-gradient(135deg, #ffb74d 0%, var(--c-accent) 100%);
          }

          &.role-icon {
            background: linear-gradient(135deg, #5c6bc0 0%, var(--c-primary) 100%);
          }

          &.activity-icon {
            background: linear-gradient(135deg, #26a69a 0%, #4dd0e1 100%);
          }
        }

        .stat-info {
          .stat-value {
            font-size: 32px;
            font-weight: 700;
            color: var(--c-text-main);
            margin-bottom: 5px;
            font-family: 'Inter', sans-serif;
          }

          .stat-label {
            font-size: 14px;
            color: var(--c-text-sub);
            margin: 0;
          }
        }
      }
    }
  }

  .content-row {
    .welcome-card,
    .info-card {
      background: var(--bg-card);
      border: 1px solid hsl(var(--border) / 0.6);
      border-radius: 18px;
      box-shadow: var(--shadow-card);
      backdrop-filter: blur(12px);

      .card-header {
        display: flex;
        align-items: center;
        justify-content: space-between;

        .title {
          color: var(--c-primary);
          font-family: 'Inter', sans-serif;
          font-size: 18px;
          font-weight: 600;
        }

        .meta {
          font-size: 12px;
          text-transform: uppercase;
          letter-spacing: 0.12em;
          color: var(--c-text-sub);
        }
      }

      :deep(.el-card__header) {
        border-bottom: 1px solid hsl(var(--border) / 0.6);
      }
    }

    .welcome-card {
      .welcome-content {
        .greeting {
          font-size: 24px;
          color: var(--c-primary);
          margin-bottom: 10px;
          font-family: 'Inter', sans-serif;
        }

        .username {
          font-size: 32px;
          color: var(--c-text-main);
          margin-bottom: 20px;
          font-weight: 600;
        }

        .description {
          font-size: 14px;
          color: var(--c-text-sub);
          line-height: 1.6;
        }
      }
    }

    .info-card {
      .info-content {
        .info-item {
          display: flex;
          justify-content: space-between;
          padding: 15px 0;
          border-bottom: 1px solid hsl(var(--border) / 0.6);

          &:last-child {
            border-bottom: none;
          }

          .label {
            color: var(--c-text-sub);
            font-size: 14px;
          }

          .value {
            color: var(--c-primary);
            font-size: 14px;
            font-weight: 600;
          }
        }
      }
    }
  }
}
</style>
