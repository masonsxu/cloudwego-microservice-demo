<template>
  <div class="relative flex h-screen overflow-hidden">
    <div class="fixed inset-0 z-0 pointer-events-none opacity-60" style="background-image: radial-gradient(var(--tech-grid-color) 1px, transparent 1px); background-size: 32px 32px; mask-image: linear-gradient(to bottom, black 0%, transparent 80%); -webkit-mask-image: linear-gradient(to bottom, black 0%, transparent 80%);" />

    <div
      class="relative z-10 overflow-hidden transition-all duration-300 ease-in-out"
      :style="{ width: sidebarWidth }"
      style="background: var(--bg-card); border-right: 1px solid var(--c-border-accent); backdrop-filter: blur(8px);"
    >
      <AppSidebar />
    </div>
    <div class="flex flex-col flex-1 overflow-hidden">
      <header
        class="relative z-10 flex h-[60px] items-center"
        style="background: var(--bg-card); border-bottom: 1px solid var(--c-border-accent); backdrop-filter: blur(8px);"
      >
        <AppHeader />
      </header>
      <main class="relative z-10 flex-1 overflow-y-auto p-6 bg-transparent">
        <router-view />
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

const sidebarWidth = computed(() => {
  return appStore.sidebarCollapsed ? '64px' : '240px'
})
</script>
