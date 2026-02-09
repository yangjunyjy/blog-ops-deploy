import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref(false)

  // 从 localStorage 加载主题
  const loadThemeFromStorage = () => {
    try {
      const savedTheme = localStorage.getItem('theme')
      if (savedTheme) {
        isDark.value = savedTheme === 'dark'
      } else {
        // 默认使用系统主题
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
        isDark.value = prefersDark
      }
    } catch (error) {
      console.error('加载主题失败:', error)
    }
  }

  // 初始化时加载
  loadThemeFromStorage()

  // 监听主题变化并应用到 DOM
  watch(isDark, (newVal) => {
    if (newVal) {
      document.documentElement.classList.add('dark')
      localStorage.setItem('theme', 'dark')
    } else {
      document.documentElement.classList.remove('dark')
      localStorage.setItem('theme', 'light')
    }
  }, { immediate: true })

  // 切换主题
  const toggleTheme = () => {
    isDark.value = !isDark.value
  }

  // 设置主题
  const setTheme = (dark) => {
    isDark.value = dark
  }

  return {
    isDark,
    toggleTheme,
    setTheme,
    loadThemeFromStorage
  }
})
