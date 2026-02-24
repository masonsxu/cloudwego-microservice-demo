import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 侧边栏折叠状态
  const sidebarCollapsed = ref(false)

  // 页面加载状态
  const loading = ref(false)

  // 当前语言
  const language = ref<string>(localStorage.getItem('language') || 'zh-CN')

  // 切换侧边栏
  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  // 设置加载状态
  function setLoading(value: boolean) {
    loading.value = value
  }

  // 设置语言
  function setLanguage(lang: string) {
    language.value = lang
    localStorage.setItem('language', lang)
  }

  return {
    sidebarCollapsed,
    loading,
    language,
    toggleSidebar,
    setLoading,
    setLanguage
  }
})
