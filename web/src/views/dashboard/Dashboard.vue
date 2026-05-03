<template>
  <div class="space-y-8">
    <!-- 页头 -->
    <header>
      <h1 class="text-[28px] font-semibold leading-tight tracking-[-0.015em] text-[color:var(--color-ink-strong)]">
        {{ greeting }}{{ authStore.username ? `，${authStore.username}` : '' }}
      </h1>
      <p class="mt-1.5 text-[14px] text-[color:var(--color-ink-muted)]">
        {{ t('dashboard.welcomeDescription') }}
      </p>
    </header>

    <!-- 概览 -->
    <section>
      <h2 class="mb-3 text-[11px] font-semibold uppercase tracking-[0.08em] text-[color:var(--color-ink-subtle)]">
        {{ t('dashboard.overview') }}
      </h2>

      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <RouterLink
          v-for="stat in statCards"
          :key="stat.key"
          :to="stat.to"
          class="group block rounded-md border border-subtle bg-canvas p-5 transition-colors duration-[var(--duration-fast)] hover:border-strong"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 text-[color:var(--color-ink-muted)]">
              <component :is="stat.icon" class="h-4 w-4" />
              <span class="text-[13px] font-medium">{{ t(stat.label) }}</span>
            </div>
            <ArrowUpRight class="h-3.5 w-3.5 text-[color:var(--color-ink-subtle)] opacity-0 transition-opacity duration-[var(--duration-fast)] group-hover:opacity-100" />
          </div>

          <div class="mt-4 flex items-baseline gap-2">
            <span
              v-if="!loading"
              class="tabular text-[28px] font-semibold leading-none tracking-[-0.015em] text-[color:var(--color-ink-strong)]"
            >
              {{ formatNumber(stat.value) }}
            </span>
            <span v-else class="wb-skeleton inline-block h-7 w-20" />
          </div>

          <p class="mt-1.5 text-[12px] text-[color:var(--color-ink-subtle)]">
            {{ t(stat.subtitle) }}
          </p>
        </RouterLink>
      </div>
    </section>

    <!-- 系统状态 -->
    <section>
      <h2 class="mb-3 text-[11px] font-semibold uppercase tracking-[0.08em] text-[color:var(--color-ink-subtle)]">
        {{ t('dashboard.systemStatus') }}
      </h2>
      <div class="rounded-md border border-subtle bg-canvas p-5">
        <div class="flex items-center gap-3">
          <span class="relative flex h-2 w-2">
            <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-[color:var(--color-success)] opacity-60" />
            <span class="relative inline-flex h-2 w-2 rounded-full bg-[color:var(--color-success)]" />
          </span>
          <span class="text-[14px] font-semibold text-ink">{{ t('dashboard.statusHealthy') }}</span>
          <Badge variant="default" class="ml-auto">{{ t('dashboard.statusVersion', { version: 'v1.0' }) }}</Badge>
        </div>
        <p class="mt-3 text-[13px] text-[color:var(--color-ink-muted)]">
          {{ t('dashboard.statusDescription') }}
        </p>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { getDashboardStats } from '@/api/dashboard'
import { toast } from 'vue-sonner'
import { Users, Building2, KeyRound, ArrowUpRight } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'

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

const statCards = computed(() => [
  {
    key: 'users',
    label: 'dashboard.userCount',
    subtitle: 'dashboard.userCountSubtitle',
    icon: Users,
    value: stats.value.userCount,
    to: '/system-settings/accounts'
  },
  {
    key: 'orgs',
    label: 'dashboard.orgCount',
    subtitle: 'dashboard.orgCountSubtitle',
    icon: Building2,
    value: stats.value.orgCount,
    to: '/system-settings/organization'
  },
  {
    key: 'roles',
    label: 'dashboard.roleCount',
    subtitle: 'dashboard.roleCountSubtitle',
    icon: KeyRound,
    value: stats.value.roleCount,
    to: '/system-settings/roles'
  }
])

function formatNumber(value: number): string {
  return new Intl.NumberFormat('en-US').format(value)
}

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
