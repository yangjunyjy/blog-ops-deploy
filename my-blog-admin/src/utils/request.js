import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { generateOnceToken } from '@/api/auth'
import tokenManager from './totp'

// TOTP配置获取状态
let totpConfigPromise = null
let isTOTPInitialized = false

// 获取TOTP配置信息
const fetchTOTPConfig = async () => {
  if (totpConfigPromise) {
    return totpConfigPromise
  }

  try {
    totpConfigPromise = axios.get('/rbac/auth/totp-config', {
      baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1'
    })
      .then(response => {
        const config = response.data?.data || response.data
        return config
      })
      .catch(error => {
        console.warn('获取TOTP配置失败，将使用传统Token模式:', error)
        return { enabled: false }
      })
      .finally(() => {
        totpConfigPromise = null
      })
    
    return totpConfigPromise
  } catch (error) {
    console.warn('获取TOTP配置失败，将使用传统Token模式:', error)
    return { enabled: false }
  }
}

// 初始化TOTP（如果需要）
const initializeTOTP = async () => {
  if (isTOTPInitialized) {
    return
  }

  try {
    const config = await fetchTOTPConfig()
    
    if (config.enabled) {
      // 注意：实际密钥应该从更安全的途径获取，这里使用配置中的密钥
      // 生产环境中，密钥应该通过安全的方式传输给客户端
      await tokenManager.initialize({
        enabled: true,
        secret: config.secret || 'default-totp-secret-change-in-production',
        timeStep: config.timeStep || 30,
        windowSize: config.windowSize || 1
      })
      console.log('TOTP动态Token已启用')
    } else {
      console.log('TOTP未启用，使用传统Token模式')
      await tokenManager.initialize(null)
    }
    
    isTOTPInitialized = true
  } catch (error) {
    console.error('TOTP初始化失败，使用传统Token模式:', error)
    isTOTPInitialized = false
  }
}

// Token获取状态管理（传统模式）
let tokenFetchPromise = null

// 获取一次性Token的辅助函数（传统模式）
const fetchOnceToken = async (forceNew = false) => {
  // 如果有正在进行的token获取请求且不强制新建，直接返回该Promise
  if (!forceNew && tokenFetchPromise) {
    return tokenFetchPromise
  }

  // 发起新的token获取请求
  tokenFetchPromise = generateOnceToken()
    .then(response => {
      const token = response.data?.token || response.data
      if (token) {
        sessionStorage.setItem('rbac_once_token', token)
        return token
      } else {
        throw new Error('无法获取一次性Token')
      }
    })
    .catch(error => {
      console.error('获取一次性Token失败:', error)
      throw error
    })
    .finally(() => {
      tokenFetchPromise = null
    })

  return tokenFetchPromise
}

// 创建axios实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
  withCredentials: true // 允许携带cookie
})

// 请求拦截器
request.interceptors.request.use(
  async config => {
    // 判断是否为后台管理API（以/rbac/开头）
    const isRBACAPI = config.url && config.url.startsWith('/rbac/')

    if (isRBACAPI) {
      // 排除公开接口和token获取请求本身，避免递归
      const publicPaths = [
        '/rbac/auth/login',
        '/rbac/auth/send-email-code',
        '/rbac/auth/totp-config',
        '/rbac/auth/totp-status'
      ]
      const isPublicAPI = publicPaths.some(path => config.url.startsWith(path))
      const isTokenRequest = config.url === '/rbac/auth/token'

      if (!isPublicAPI && !isTokenRequest) {
        try {
          // 确保TOTP已初始化
          if (!isTOTPInitialized) {
            await initializeTOTP()
          }
          
          // 根据TOTP是否启用选择不同的Token生成方式
          if (tokenManager.isTOTPEnabled()) {
            // TOTP模式：本地计算动态Token
            const dynamicToken = await tokenManager.generateToken(config.url, config.data || {})
            
            // 将Token添加到URL参数中（类似CMDB系统）
            // 注意：实际应用中可以根据后端要求选择添加方式（URL参数、请求头等）
            const separator = config.url.includes('?') ? '&' : '?'
            config.url = `${config.url}${separator}otnonce_token=${encodeURIComponent(dynamicToken.token_id)}`
            
            console.debug('使用TOTP动态Token:', dynamicToken.token_id.substring(0, 16) + '...')
          } else {
            // 传统模式：请求后端获取Token，通过URL参数传递避免CORS问题
            // 检查是否为文件上传接口，如果是则强制获取新token
            // 只精确匹配 SFTP 上传接口，避免误判其他接口
            const isUploadRequest = config.url.includes('/sftp/uploadFile')
            const onceToken = await fetchOnceToken(isUploadRequest)
            console.log("获取到的token",onceToken);

            if (onceToken) {
              const separator = config.url.includes('?') ? '&' : '?'
              config.url = `${config.url}${separator}once_token=${onceToken.token_id}`
              console.debug('使用一次性Token:', onceToken.token_id.substring(0, 16) + '...')
            }
          }
        } catch (error) {
          console.error('获取Token失败:', error)
          // 不抛出错误，让请求继续，后端会返回401
        }
      }
      // token获取请求本身不需要添加Token，只依赖Session Cookie
    } else {
      // 前台API：使用Bearer Token认证
      const token = localStorage.getItem('token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
    }
    return config
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    const { data } = response

    // 统一处理后端API响应格式 {code: 0, message: '', data: {}}
    if (data.hasOwnProperty('code')) {
      if (data.code === 0) {
        return data
      } else if (data.code === 200) {
        // 兼容Mock API的响应格式
        return data
      } else {
        ElMessage.error(data.message || '请求失败')
        return Promise.reject(new Error(data.message || '请求失败'))
      }
    }

    return data
  },
  error => {
    console.error('响应错误:', error)

    if (error.response) {
      const { status, data } = error.response
      const url = error.config?.url || ''
      const isRBACAPI = url.startsWith('/rbac/')

      switch (status) {
        case 401:
          if (isRBACAPI) {
            // 后台API：可能是Session过期或Token无效
            const message = data?.message || '认证失败'
            if (message.includes('Token') || message.includes('token')) {
              // Token相关错误，清除Token并提示重新获取
              ElMessage.error('安全令牌已失效，请重新获取')
            } else {
              // Session过期，需要重新登录
              ElMessage.error('会话已过期，请重新登录')
              sessionStorage.removeItem('token')
              router.push('/login')
            }
          } else {
            // 前台API：Token过期
            ElMessage.error('未授权,请重新登录')
            localStorage.removeItem('token')
            localStorage.removeItem('userInfo')
            router.push('/login')
          }
          break
        case 403:
          ElMessage.error('拒绝访问')
          break
        case 404:
          ElMessage.error('请求资源不存在')
          break
        case 500:
          ElMessage.error('服务器错误')
          break
        default:
          ElMessage.error(data?.message || '请求失败')
      }
    } else {
      ElMessage.error('网络错误,请检查网络连接')
    }

    return Promise.reject(error)
  }
)

export default request
