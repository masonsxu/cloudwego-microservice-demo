import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 侧边栏折叠状态
  const sidebarCollapsed = ref(false)

  // 页面加载状态
  const loading = ref(false)

  // 当前语言
  const language = ref<string>(localStorage.getItem('language') || 'zh-CN')

  // 主题：dark | light
  const theme = ref<'dark' | 'light'>(
    (localStorage.getItem('theme') as 'dark' | 'light') || 'dark'
  )

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

  // 切换主题
  function toggleTheme() {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
    localStorage.setItem('theme', theme.value)
    applyTheme(theme.value)
  }

  // 应用主题到 html 元素
  function applyTheme(t: 'dark' | 'light') {
    document.documentElement.setAttribute('data-theme', t)
  }

  // 初始化时应用已保存的主题
  applyTheme(theme.value)

  return {
    sidebarCollapsed,
    loading,
    language,
    theme,
    toggleSidebar,
    setLoading,
    setLanguage,
    toggleTheme
  }
})
