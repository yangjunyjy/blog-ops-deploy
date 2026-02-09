import request from '@/utils/request'
import { mockApi } from '@/mock'

// 是否使用Mock数据
const USE_MOCK = false


// 登出
export const logout = () => {
  if (USE_MOCK) {
    return mockApi.logout()
  }
  return request.post('/rbac/auth/logout')
}

// 获取当前用户信息
export const getUserInfo = () => {
  if (USE_MOCK) {
    return mockApi.getUserInfo()
  }
  return request.get('/rbac/auth/info')
}

// 刷新token
export const refreshToken = () => {
  return request.post('/auth/refresh')
}

export const rbacLogin = (data) => {
  return request.post('/rbac/auth/login',data)
}

// 发送邮箱验证码
export const sendEmailCode = (email) => {
  return request.post('/rbac/auth/send-email-code',{
    email
  })
}

// 邮箱验证码注册
export const registerWithEmailCode = (data) => {
  return request({
    url: `${API_PREFIX}/public/auth/register`,
    method: 'post',
    data
  })
}

// 验证邮箱验证码
export const verifyEmailCaptcha = (data) => {
  return request({
    url: `${API_PREFIX}/public/auth/verify`,
    method: 'post',
    data
  })
}

// 生成一次性Token（后台管理）
export const generateOnceToken = () => {
  if (USE_MOCK) {
    return mockApi.generateOnceToken()
  }
  return request.post('/rbac/auth/token')
}

// 获取TOTP配置信息
export const getTOTPConfig = () => {
  return request.get('/rbac/auth/totp-config')
}

// 检查TOTP状态
export const checkTOTPStatus = () => {
  return request.get('/rbac/auth/totp-status')
}

// 获取当前用户菜单（后台管理）
export const getRBACMenu = () => {
  if (USE_MOCK) {
    return mockApi.getRBACMenu()
  }
  return request.get('/rbac/auth/menu')
}
