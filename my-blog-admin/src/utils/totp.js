/**
 * TOTP (基于时间的一次性密码) 生成器
 * 实现类似CMDB系统的客户端动态Token生成机制
 * 无需请求后端获取Token，客户端本地计算
 */

import CryptoJS from 'crypto-js'

class TOTPGenerator {
  /**
   * 创建TOTP生成器
   * @param {string} secretKey - 共享密钥（从后端配置获取）
   * @param {object} options - 配置选项
   */
  constructor(secretKey, options = {}) {
    if (!secretKey || secretKey.trim() === '') {
      throw new Error('TOTP密钥不能为空')
    }
    
    this.secretKey = secretKey
    this.timeStep = options.timeStep || 30 // 时间步长，默认30秒
    this.windowSize = options.windowSize || 1 // 验证窗口大小
    this.tokenLength = options.tokenLength || 6 // Token长度，默认6位
  }
  
  /**
   * 生成当前时间窗口的TOTP Token（6位数字）
   * @returns {string} 6位数字Token
   */
  generate() {
    const counter = this.getCounter()
    const token = this.hotp(counter)
    return token.toString().padStart(this.tokenLength, '0')
  }
  
  /**
   * 验证TOTP Token是否有效
   * @param {string} token - 待验证的Token
   * @returns {boolean} 是否有效
   */
  verify(token) {
    if (!token || typeof token !== 'string') {
      return false
    }
    
    // 清理Token
    token = token.trim()
    
    // 检查Token长度
    if (token.length !== this.tokenLength) {
      return false
    }
    
    // 检查是否为纯数字
    if (!/^\d+$/.test(token)) {
      return false
    }
    
    const currentCounter = this.getCounter()
    
    // 检查当前窗口及前后窗口
    for (let i = -this.windowSize; i <= this.windowSize; i++) {
      const testCounter = currentCounter + i
      const expectedToken = this.hotp(testCounter)
      const expectedTokenStr = expectedToken.toString().padStart(this.tokenLength, '0')
      
      if (token === expectedTokenStr) {
        return true
      }
    }
    
    return false
  }
  
  /**
   * 为API请求生成动态Token（类似CMDB的长Token格式）
   * @param {string} endpoint - API端点
   * @param {object} data - 请求数据
   * @returns {string} 动态Token
   */
  generateForAPI(endpoint, data = {}) {
    const timestamp = Date.now()
    const timeCounter = this.getCounter()
    
    // 创建请求摘要：时间戳 + 时间窗口 + 端点 + 数据摘要
    const requestDigest = `${timestamp}:${timeCounter}:${endpoint}:${JSON.stringify(data)}`
    
    // 使用HMAC-SHA256生成Token
    const hmac = CryptoJS.HmacSHA256(requestDigest, this.secretKey)
    
    // 转换为Base64，并替换URL不安全的字符
    return CryptoJS.enc.Base64.stringify(hmac)
      .replace(/\+/g, '-')
      .replace(/\//g, '_')
      .replace(/=/g, '')
  }
  
  /**
   * 获取当前时间计数器
   * @returns {number} 时间计数器
   */
  getCounter() {
    return Math.floor(Date.now() / 1000 / this.timeStep)
  }
  
  /**
   * HOTP算法：生成HMAC-based一次性密码
   * @param {number} counter - 计数器
   * @returns {number} 生成的Token（整数）
   */
  hotp(counter) {
    // 将计数器转换为8字节的大端序
    const counterBytes = new ArrayBuffer(8)
    const counterView = new DataView(counterBytes)
    counterView.setUint32(0, Math.floor(counter / 0x100000000), false) // 高32位
    counterView.setUint32(4, counter & 0xffffffff, false) // 低32位
    
    // 转换为WordArray（CryptoJS格式）
    const counterHex = Array.from(new Uint8Array(counterBytes))
      .map(b => b.toString(16).padStart(2, '0'))
      .join('')
    const counterWordArray = CryptoJS.enc.Hex.parse(counterHex)
    
    // 计算HMAC-SHA1
    const hmac = CryptoJS.HmacSHA1(counterWordArray, this.secretKey)
    const hash = hmac.toString(CryptoJS.enc.Hex)
    
    // 动态截取（标准HOTP算法）
    const offset = parseInt(hash.substr(-2, 2), 16) & 0x0f
    const binaryCode = ((parseInt(hash.substr(offset * 2, 2), 16) & 0x7f) << 24) |
      ((parseInt(hash.substr(offset * 2 + 2, 2), 16) & 0xff) << 16) |
      ((parseInt(hash.substr(offset * 2 + 4, 2), 16) & 0xff) << 8) |
      (parseInt(hash.substr(offset * 2 + 6, 2), 16) & 0xff)
    
    // 取模得到6位数字
    const mod = Math.pow(10, this.tokenLength)
    return binaryCode % mod
  }
}

/**
 * Token管理器：根据后端配置自动选择Token生成方式
 */
class TokenManager {
  constructor() {
    this.totpEnabled = false
    this.totpGenerator = null
    this.tokenConfig = null
    this.tokenCache = new Map() // 缓存生成的Token，避免短时间内重复计算
    this.cacheTimeout = 2000 // Token缓存时间（毫秒）
  }
  
  /**
   * 初始化Token管理器
   * @param {object} config - TOTP配置
   */
  async initialize(config) {
    if (!config) {
      this.totpEnabled = false
      this.totpGenerator = null
      return
    }
    
    this.totpEnabled = config.enabled || false
    
    if (this.totpEnabled && config.secret) {
      try {
        this.totpGenerator = new TOTPGenerator(config.secret, {
          timeStep: config.timeStep || 30,
          windowSize: config.windowSize || 1
        })
        this.tokenConfig = config
        console.log('TOTP已启用，时间步长:', config.timeStep, '秒')
      } catch (error) {
        console.error('TOTP初始化失败:', error)
        this.totpEnabled = false
        this.totpGenerator = null
      }
    } else {
      console.log('TOTP未启用，使用传统Token机制')
    }
  }
  
  /**
   * 生成Token（根据配置选择方式）
   * @param {string} endpoint - API端点
   * @param {object} data - 请求数据
   * @returns {Promise<string>} 生成的Token
   */
  async generateToken(endpoint, data = {}) {
    // 检查缓存
    const cacheKey = `${endpoint}:${JSON.stringify(data)}`
    const cached = this.tokenCache.get(cacheKey)
    if (cached && Date.now() - cached.timestamp < this.cacheTimeout) {
      return cached.token
    }
    
    let token
    if (this.totpEnabled && this.totpGenerator) {
      // 使用TOTP生成动态Token
      token = this.totpGenerator.generateForAPI(endpoint, data)
    } else {
      // 传统方式：调用后端API获取Token
      token = await this.fetchOnceTokenFromServer()
    }
    
    // 更新缓存
    this.tokenCache.set(cacheKey, {
      token,
      timestamp: Date.now()
    })
    
    return token
  }
  
  /**
   * 从服务器获取一次性Token（传统方式）
   * @returns {Promise<string>} 获取的Token
   */
  async fetchOnceTokenFromServer() {
    try {
      // 这里应该调用后端API生成Token
      // 简化示例：返回时间戳+随机数
      return Date.now().toString() + Math.random().toString(36).substring(2, 15)
    } catch (error) {
      console.error('获取Token失败:', error)
      throw error
    }
  }
  
  /**
   * 验证Token是否有效
   * @param {string} token - 待验证的Token
   * @returns {boolean} 是否有效
   */
  verifyToken(token) {
    if (this.totpEnabled && this.totpGenerator) {
      // 对于TOTP Token，需要特殊的验证逻辑
      // 这里简化处理，真实场景应该根据Token格式判断
      if (token.length === this.totpGenerator.tokenLength && /^\d+$/.test(token)) {
        return this.totpGenerator.verify(token)
      }
    }
    // 传统Token：假设总是有效（实际需要后端验证）
    return true
  }
  
  /**
   * 检查TOTP是否启用
   * @returns {boolean} 是否启用
   */
  isTOTPEnabled() {
    return this.totpEnabled && this.totpGenerator !== null
  }
  
  /**
   * 清理缓存
   */
  clearCache() {
    this.tokenCache.clear()
  }
}

// 创建全局Token管理器实例
const tokenManager = new TokenManager()

export { TOTPGenerator, TokenManager, tokenManager }
export default tokenManager