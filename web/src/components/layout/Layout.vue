<template>
  <el-container class="layout-container">
    <!-- Tech Grid 背景层 -->
    <div class="tech-grid-overlay" />

    <el-aside :width="sidebarWidth" class="sidebar">
      <AppSidebar />
    </el-aside>
    <el-container>
      <el-header class="header">
        <AppHeader />
      </el-header>
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
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

<style scoped lang="scss">
.layout-container {
  height: 100vh;
  position: relative;

  .tech-grid-overlay {
    position: fixed;
    inset: 0;
    background-image: radial-gradient(var(--tech-grid-color) 1px, transparent 1px);
    background-size: 32px 32px;
    mask-image: linear-gradient(to bottom, black 0%, transparent 80%);
    -webkit-mask-image: linear-gradient(to bottom, black 0%, transparent 80%);
    pointer-events: none;
    z-index: 0;
    opacity: 0.6;
  }

  .sidebar {
    transition: width 0.3s ease;
    background: var(--bg-card);
    border-right: 1px solid var(--c-border-accent);
    overflow: hidden;
    z-index: 10;
    position: relative;
    backdrop-filter: blur(8px);
  }

  .header {
    background: var(--bg-card);
    border-bottom: 1px solid var(--c-border-accent);
    padding: 0;
    height: 60px;
    display: flex;
    align-items: center;
    z-index: 10;
    position: relative;
    backdrop-filter: blur(8px);
  }

  .main-content {
    background-color: transparent;
    padding: 24px;
    overflow-y: auto;
    position: relative;
    z-index: 1;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
