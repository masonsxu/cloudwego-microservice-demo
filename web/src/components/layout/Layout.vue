<template>
  <el-container class="layout-container">
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

  .sidebar {
    transition: width 0.3s ease;
    background: linear-gradient(145deg, rgba(30, 32, 36, 0.9), rgba(20, 20, 22, 0.95));
    border-right: 1px solid rgba(212, 175, 55, 0.2);
    overflow: hidden;
  }

  .header {
    background: linear-gradient(145deg, rgba(30, 32, 36, 0.9), rgba(20, 20, 22, 0.95));
    border-bottom: 1px solid rgba(212, 175, 55, 0.2);
    padding: 0;
    height: 60px;
    display: flex;
    align-items: center;
  }

  .main-content {
    background-color: #141416;
    background-image: radial-gradient(circle at 50% 10%, #2a2d35 0%, #000000 100%);
    padding: 20px;
    overflow-y: auto;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
