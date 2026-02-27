<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { initApiClient } from '@/api/config'

// 在应用启动时初始化 authStore 和 API 客户端配置
// 确保在任何 API 请求之前，token 已经从 localStorage 恢复
onMounted(() => {
  const authStore = useAuthStore()
  authStore.restoreState()

  // 初始化 API 客户端配置（设置 TOKEN）
  initApiClient()
})
</script>

<style>
#app {
  height: 100%;
}
</style>
