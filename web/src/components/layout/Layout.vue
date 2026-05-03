<template>
  <div class="flex h-screen w-full overflow-hidden bg-canvas text-ink">
    <aside
      class="flex-shrink-0 overflow-hidden bg-sunken border-r border-subtle transition-[width] duration-200 ease-out"
      :style="{ width: sidebarWidth }"
    >
      <AppSidebar />
    </aside>

    <div class="flex flex-col flex-1 min-w-0 overflow-hidden">
      <header
        class="flex-shrink-0 flex h-[var(--layout-topbar-height)] items-center bg-canvas border-b border-subtle"
      >
        <AppHeader />
      </header>

      <main class="flex-1 overflow-y-auto bg-canvas">
        <div class="mx-auto w-full max-w-[var(--layout-content-max)] px-[var(--layout-content-px)] py-[var(--layout-content-py)]">
          <router-view />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'
import AppHeader from './AppHeader.vue'
import AppSidebar from './AppSidebar.vue'

const appStore = useAppStore()

const sidebarWidth = computed(() =>
  appStore.sidebarCollapsed ? 'var(--layout-sidebar-collapsed)' : 'var(--layout-sidebar-width)'
)
</script>
