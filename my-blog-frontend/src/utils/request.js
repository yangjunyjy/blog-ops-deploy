import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081/',
  timeout: 15000
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 可以在这里添加token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    console.log(response);
    const res = response.data
    if (res.code !== 200) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return res
  },
  error => {
    console.error('请求错误:', error)

    // 处理401未授权错误（token过期或无效）
    if (error.response?.status === 401) {
      // 清除本地存储的用户信息和token
      localStorage.removeItem('user')
      localStorage.removeItem('token')

      // 提示用户登录已过期
      ElMessage.warning('登录已过期，请重新登录')

      // 跳转到登录页
      router.push('/login')
      return Promise.reject(error)
    }

    // 处理其他错误
    ElMessage.error(error.response?.data?.message || error.message || '网络请求失败')
    return Promise.reject(error)
  }
)

export default request
