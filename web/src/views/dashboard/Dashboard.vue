<template>
  <div class="dashboard">
    <div class="grid gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-5">
      <Card class="stat-card">
        <CardContent class="p-6">
          <div class="flex items-center gap-5">
            <div class="stat-icon user-icon">
              <User class="w-7 h-7" />
            </div>
            <div>
              <p class="text-3xl font-bold text-[var(--c-text-main)] font-['Inter'] mb-1">{{ stats.userCount }}</p>
              <p class="text-sm text-[var(--c-text-sub)] m-0">{{ t('dashboard.userCount') }}</p>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card class="stat-card">
        <CardContent class="p-6">
          <div class="flex items-center gap-5">
            <div class="stat-icon org-icon">
              <Building2 class="w-7 h-7" />
            </div>
            <div>
              <p class="text-3xl font-bold text-[var(--c-text-main)] font-['Inter'] mb-1">{{ stats.orgCount }}</p>
              <p class="text-sm text-[var(--c-text-sub)] m-0">{{ t('dashboard.orgCount') }}</p>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card class="stat-card">
        <CardContent class="p-6">
          <div class="flex items-center gap-5">
            <div class="stat-icon role-icon">
              <Key class="w-7 h-7" />
            </div>
            <div>
              <p class="text-3xl font-bold text-[var(--c-text-main)] font-['Inter'] mb-1">{{ stats.roleCount }}</p>
              <p class="text-sm text-[var(--c-text-sub)] m-0">{{ t('dashboard.roleCount') }}</p>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card class="stat-card">
        <CardContent class="p-6">
          <div class="flex items-center gap-5">
            <div class="stat-icon activity-icon">
              <TrendingUp class="w-7 h-7" />
            </div>
            <div>
              <p class="text-3xl font-bold text-[var(--c-text-main)] font-['Inter'] mb-1">{{ stats.activityCount }}</p>
              <p class="text-sm text-[var(--c-text-sub)] m-0">{{ t('dashboard.recentActivity') }}</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="grid gap-5 lg:grid-cols-3">
      <div class="lg:col-span-2">
        <Card class="content-card">
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-[var(--c-primary)] font-['Inter'] text-lg font-semibold">{{ t('dashboard.welcome') }}</CardTitle>
            <span class="text-xs uppercase tracking-widest text-[var(--c-text-sub)]">Overview</span>
          </CardHeader>
          <CardContent>
            <div>
              <p class="text-2xl text-[var(--c-primary)] mb-2.5 font-['Inter']">{{ greeting }}</p>
              <p class="text-[32px] text-[var(--c-text-main)] mb-5 font-semibold">{{ authStore.username }}</p>
              <p class="text-sm text-[var(--c-text-sub)] leading-relaxed">
                {{ t('dashboard.welcomeDescription') }}
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
      <div>
        <Card class="content-card">
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-[var(--c-primary)] font-['Inter'] text-lg font-semibold">系统信息</CardTitle>
            <span class="text-xs uppercase tracking-widest text-[var(--c-text-sub)]">System</span>
          </CardHeader>
          <CardContent>
            <div>
              <div class="flex justify-between py-3.5 border-b border-[hsl(var(--border)/0.6)]">
                <span class="text-sm text-[var(--c-text-sub)]">框架版本</span>
                <span class="text-sm text-[var(--c-primary)] font-semibold">Vue 3.4+</span>
              </div>
              <div class="flex justify-between py-3.5 border-b border-[hsl(var(--border)/0.6)]">
                <span class="text-sm text-[var(--c-text-sub)]">UI 组件</span>
                <span class="text-sm text-[var(--c-primary)] font-semibold">shadcn-vue</span>
              </div>
              <div class="flex justify-between py-3.5 border-b border-[hsl(var(--border)/0.6)]">
                <span class="text-sm text-[var(--c-text-sub)]">状态管理</span>
                <span class="text-sm text-[var(--c-primary)] font-semibold">Pinia</span>
              </div>
              <div class="flex justify-between py-3.5">
                <span class="text-sm text-[var(--c-text-sub)]">主题</span>
                <span class="text-sm text-[var(--c-primary)] font-semibold">CareerCompass Inspired</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { getDashboardStats } from '@/api/dashboard'
import { toast } from 'vue-sonner'
import { User, Building2, Key, TrendingUp } from 'lucide-vue-next'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

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
    toast.error(error.message || t('common.operationFailed'))
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
  .stat-card {
    background: var(--bg-card);
    border: 1px solid hsl(var(--border) / 0.6);
    border-radius: 18px;
    box-shadow: var(--shadow-card);
    backdrop-filter: blur(12px);
    position: relative;
    overflow: hidden;
  }

  .stat-icon {
    width: 60px;
    height: 60px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #ffffff;
    box-shadow: inset 0 0 0 1px hsl(var(--border) / 0.4);
    flex-shrink: 0;

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

  .content-card {
    background: var(--bg-card);
    border: 1px solid hsl(var(--border) / 0.6);
    border-radius: 18px;
    box-shadow: var(--shadow-card);
    backdrop-filter: blur(12px);
  }
}
</style>
