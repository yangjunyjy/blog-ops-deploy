import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  const user = ref(null)
  const token = ref(null)

  // 检查用户信息是否有效
  const isValidUser = (userData) => {
    return userData && typeof userData === 'object' && userData.id && userData.username
  }

  // 从 localStorage 加载用户信息
  const loadUserFromStorage = () => {
    try {
      const userData = localStorage.getItem('user')
      const tokenData = localStorage.getItem('token')
      if (userData) {
        const parsedUser = JSON.parse(userData)
        // 只有当用户信息有效时才加载
        if (isValidUser(parsedUser)) {
          user.value = parsedUser
        } else {
          // 清除无效的用户数据
          localStorage.removeItem('user')
        }
      }
      if (tokenData) {
        token.value = tokenData
      }
    } catch (error) {
      console.error('加载用户信息失败:', error)
      // 解析失败，清除数据
      localStorage.removeItem('user')
    }
  }

  // 初始化时加载
  loadUserFromStorage()

  // 计算属性
  const isLoggedIn = computed(() => !!user.value && !!token.value && isValidUser(user.value))

  // 登录
  const login = (userData, userToken) => {
    // 只保存有效的用户信息
    if (isValidUser(userData)) {
      user.value = userData
      token.value = userToken
      localStorage.setItem('user', JSON.stringify(userData))
      localStorage.setItem('token', userToken)
    } else {
      console.error('无效的用户信息:', userData)
    }
  }

  // 登出
  const logout = () => {
    user.value = null
    token.value = null
    localStorage.removeItem('user')
    localStorage.removeItem('token')
  }

  // 更新用户信息
  const updateUser = (userData) => {
    user.value = { ...user.value, ...userData }
    localStorage.setItem('user', JSON.stringify(user.value))
  }

  return {
    user,
    token,
    isLoggedIn,
    login,
    logout,
    updateUser,
    loadUserFromStorage,
    isValidUser
  }
})


