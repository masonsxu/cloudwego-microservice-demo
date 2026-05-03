<template>
  <div class="space-y-6">
    <header>
      <h1 class="text-[22px] font-semibold leading-tight tracking-[-0.012em] text-[color:var(--color-ink-strong)]">
        {{ t('system.menuManage') }}
      </h1>
      <p class="mt-1 text-[13px] text-[color:var(--color-ink-muted)]">
        {{ t('system.menuManageSubtitle') }}
      </p>
    </header>

    <div class="rounded-md border border-subtle bg-canvas">
      <div class="flex flex-col items-center justify-center gap-2 py-16">
        <Menu class="h-8 w-8 text-[color:var(--color-ink-subtle)]" />
        <p class="text-[14px] font-medium text-ink">{{ t('system.menuComingSoon') }}</p>
        <p class="text-[12px] text-[color:var(--color-ink-muted)]">{{ t('system.menuComingSoonHint') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Menu } from 'lucide-vue-next'

const { t } = useI18n()
</script>
