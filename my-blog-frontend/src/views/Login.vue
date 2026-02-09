<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Message, Key, RefreshRight } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { loginWithEmailCode, sendEmailCode, getCaptcha, login } from '@/api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loginForm = ref({
  email:'',
  user: '',
  password: '',
  code: '',
  captcha: ''
})

const loading = ref(false)
const codeLoading = ref(false)
const captchaLoading = ref(false)
const countdown = ref(0)
const useCodeLogin = ref(false)
const captchaUrl = ref('')
const captchaId = ref('')

const loadCaptcha = async () => {
  captchaLoading.value = true
  try {
    const res = await getCaptcha()
    console.log('验证码响应', res)
    if (res.code === 200) {
      captchaId.value = res.data.captchaId
      captchaUrl.value = res.data.image
      console.log('验证码ID:', captchaId.value)
      console.log('验证码图片URL:', captchaUrl.value)
    }
  } catch (error) {
    console.error('获取验证码失效?', error)
    ElMessage.error('获取验证码失效')
  } finally {
    captchaLoading.value = false
  }
}

const handleSendCode = async () => {
  if (!loginForm.value.email) {
    ElMessage.warning('请输入邮箱地址')
    return
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(loginForm.value.email)) {
    ElMessage.warning('请输入有效的邮箱地址')
    return
  }

  if (countdown.value > 0) {
    return
  }

  codeLoading.value = true
  try {
    await sendEmailCode(loginForm.value.email)
    ElMessage.success('验证码已发送到您的邮箱')

    // 开始倒计时
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '验证码发送失败，请重试')
    // 如果是验证码错误，刷新验证码
    if (error.response?.data?.code === 4000) {
      loadCaptcha()
    }
  } finally {
    codeLoading.value = false
  }
}

const handleLogin = async () => {
  if ( !useCodeLogin.value && !loginForm.value.user) {
    ElMessage.warning('请输入用户名')
    return
  }
  if (!useCodeLogin.value && !loginForm.value.captcha) {
    ElMessage.warning('请输入图形验证码')
    return
  }
  if (useCodeLogin.value && !loginForm.value.email){
    ElMessage.warning("请输入邮箱")
    return
  }
  if (useCodeLogin.value) {
    if (!loginForm.value.code) {
      ElMessage.warning('请输入验证码')
      return
    }
  } else {
    if (!loginForm.value.password) {
      ElMessage.warning('请输入密码')
      return
    }
  }

  loading.value = true
  try {
    let res
    if (useCodeLogin.value) {
      // 邮箱验证码登录
      res = await loginWithEmailCode({
        email: loginForm.value.email,
        code: loginForm.value.code
      })
    } else {
      // 密码登录
      res = await login({
        username: loginForm.value.user,
        password: loginForm.value.password,
        captchaId: captchaId.value,
        captchaAnswer: loginForm.value.captcha
      })
    }

    if (res.code === 200) {
      console.log('登录响应数据:', res.data)
      console.log('用户信息:', res.data.user)
      console.log('Token:', res.data.token)

      // 确保用户信息和token都存在
      if (res.data.user && res.data.user.id && res.data.token) {
        userStore.login(res.data.user, res.data.token)
        ElMessage.success('登录成功！欢迎回来')

        const redirect = route.query.redirect || '/'
        await router.push(redirect)
      } else {
        console.error('登录响应数据异常:', res.data)
        ElMessage.error('登录数据异常，请重试')
      }
    }
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '登录失败，请重试')
    // 刷新验证码
    if (!useCodeLogin.value) {
      loadCaptcha()
      loginForm.value.captcha = ''
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCaptcha()

  // 检查并清除可能的无效token
  const existingUser = localStorage.getItem('user')
  if (existingUser) {
    try {
      const parsedUser = JSON.parse(existingUser)
      if (!parsedUser || !parsedUser.id || !parsedUser.username) {
        console.log('发现无效的用户数据，清除中...')
        localStorage.removeItem('user')
        localStorage.removeItem('token')
      }
    } catch (error) {
      console.log('用户数据解析失败，清除中...', error)
      localStorage.removeItem('user')
      localStorage.removeItem('token')
    }
  }
})
</script>

<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-decoration">
        <div class="circle circle-1"></div>
        <div class="circle circle-2"></div>
        <div class="circle circle-3"></div>
      </div>

      <div class="login-card">
        <div class="login-header">
          <h1 class="login-title">🎉 欢迎回来</h1>
          <p class="login-subtitle">登录到你的梦想空间</p>
        </div>

        <div class="login-tabs">
          <div
            class="tab-item"
            :class="{ active: !useCodeLogin }"
            @click="useCodeLogin = false"
          >
            密码登录
          </div>
          <div
            class="tab-item"
            :class="{ active: useCodeLogin }"
            @click="useCodeLogin = true"
          >
            邮箱登录
          </div>
        </div>

        <el-form class="login-form" @submit.prevent="handleLogin">
          <div class="form-item" v-if="!useCodeLogin">
            <el-icon class="form-icon"><User /></el-icon>
            <el-input
              v-model="loginForm.user"
              type="user"
              placeholder="用户"
              size="large"
              clearable
            />
          </div>
          <div class="form-item" v-else>
            <el-icon class="form-icon"><Message /></el-icon>
            <el-input
              v-model="loginForm.email"
              type="email"
              placeholder="邮箱"
              size="large"
              clearable
            />
          </div>

          <div class="form-item" v-if="!useCodeLogin">
            <el-icon class="form-icon"><Lock /></el-icon>
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="密码"
              size="large"
              show-password
              clearable
            />
          </div>

          <!-- 图形验证码- 仅密码登录显示-->
          <div class="form-item captcha-item" v-if="!useCodeLogin">
            <el-icon class="form-icon"><Key /></el-icon>
            <el-input
              v-model="loginForm.captcha"
              placeholder="图形验证码"
              size="large"
              clearable
              maxlength="4"
              style="padding-right: 120px;"
            />
            <div class="captcha-container" @click="loadCaptcha">
              <img
                v-if="captchaUrl"
                :src="captchaUrl"
                alt="验证码"
                class="captcha-image"
              />
              <el-icon v-else class="captcha-icon" :class="{ 'loading': captchaLoading }">
                <RefreshRight v-if="!captchaLoading" />
              </el-icon>
              <span class="refresh-tip" v-if="!captchaLoading"></span>
            </div>
          </div>

          <div class="form-item" v-else>
            <el-icon class="form-icon"><Key /></el-icon>
            <el-input
              v-model="loginForm.code"
              placeholder="邮箱验证码"
              size="large"
              clearable
              maxlength="6"
              style="padding-right: 120px;"
            />
            <el-button
              class="code-button"
              :disabled="countdown > 0 || codeLoading"
              :loading="codeLoading"
              @click="handleSendCode"
            >
              {{ countdown > 0 ? `${countdown}秒` : '获取验证码' }}
            </el-button>
          </div>

          <div class="form-options" v-if="!useCodeLogin">
            <el-checkbox>记住我</el-checkbox>
            <router-link to="/forgot" class="forgot-link">忘记密码？</router-link>
          </div>

          <el-button
            type="primary"
            size="large"
            class="login-button"
            :loading="loading"
            @click="handleLogin"
          >
            <span v-if="!loading">登录</span>
            <span v-else>登录中..</span>
          </el-button>

          <div class="login-footer">
            <p>还没有账号？</p>
            <router-link to="/register" class="register-link">
              快去注册吧~ 
            </router-link>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: url('/images/background.jpeg') center center / cover no-repeat;
  position: relative;
  overflow: hidden;
}

.login-page::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  pointer-events: none;
}

.login-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 480px;
  padding: 20px;
}

.login-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
}

.circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  will-change: transform;
}

.circle-1 {
  width: 300px;
  height: 300px;
  top: -150px;
  left: -150px;
  animation: float1 8s ease-in-out infinite;
}

.circle-2 {
  width: 200px;
  height: 200px;
  top: 20%;
  right: -100px;
  animation: float2 10s ease-in-out infinite;
}

.circle-3 {
  width: 150px;
  height: 150px;
  bottom: -75px;
  left: 30%;
  animation: float3 12s ease-in-out infinite;
}

@keyframes float1 {
  0%, 100% {
    transform: translateY(0) scale(1);
  }
  50% {
    transform: translateY(-30px) scale(1.1);
  }
}

@keyframes float2 {
  0%, 100% {
    transform: translateY(0) scale(1);
  }
  50% {
    transform: translateY(-20px) scale(0.9);
  }
}

@keyframes float3 {
  0%, 100% {
    transform: translateY(0) scale(1);
  }
  50% {
    transform: translateY(-25px) scale(1.2);
  }
}

.login-card {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 24px;
  padding: 50px 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-title {
  font-size: 32px;
  font-weight: 800;
  background: linear-gradient(135deg, #667eea 0%, #f093fb 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 16px;
  letter-spacing: -0.5px;
}

.login-subtitle {
  color: #718096;
  font-size: 15px;
  margin: 0;
  opacity: 0.8;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-item {
  position: relative;
}

.form-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 20px;
  color: #a0aec0;
  z-index: 10;
}

.form-item :deep(.el-input__wrapper) {
  padding-left: 48px;
  border-radius: 12px;
  border: 2px solid #e2e8f0;
  transition: all 0.3s;
}

.form-item :deep(.el-input__wrapper:hover) {
  border-color: #1890ff;
  box-shadow: 0 0 0 3px rgba(24, 144, 255, 0.1);
}

.form-item :deep(.el-input__wrapper.is-focus) {
  border-color: #1890ff;
  box-shadow: 0 0 0 3px rgba(24, 144, 255, 0.2);
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
}

.forgot-link {
  color: #1890ff;
  text-decoration: none;
  transition: all 0.3s;
}

.forgot-link:hover {
  color: #40a9ff;
  text-decoration: underline;
}

.login-button {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  border: none;
  background: #1890ff;
  border-radius: 12px;
  transition: all 0.3s;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(24, 144, 255, 0.4);
}

.login-footer {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: #718096;
}

.login-footer p {
  margin: 0 0 4px;
}

.register-link {
  color: #667eea;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s;
}

.register-link:hover {
  color: #f093fb;
  text-decoration: underline;
}

.login-tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 30px;
  padding: 4px;
  background: #f7fafc;
  border-radius: 12px;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 10px;
  font-size: 15px;
  font-weight: 600;
  color: #718096;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.tab-item.active {
  background: #1890ff;
  color: #fff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
}

.code-button {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  padding: 8px 16px;
  font-size: 14px;
  height: 36px;
  border-radius: 8px;
}

.captcha-item {
  position: relative;
}

.captcha-item :deep(.el-input__wrapper) {
  padding-right: 120px;
}

.captcha-container {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  width: 100px;
  height: 36px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f7fafc;
  border: 2px solid #e2e8f0;
  transition: all 0.3s;
}

.captcha-container:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.15);
}

.captcha-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
}

.captcha-icon {
  font-size: 20px;
  color: #a0aec0;
}

.captcha-icon.loading {
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.refresh-tip {
  position: absolute;
  bottom: 2px;
  font-size: 10px;
  color: #718096;
}

@media (max-width: 768px) {
  .login-container {
    padding: 15px;
  }

  .login-card {
    padding: 35px 25px;
    border-radius: 20px;
  }

  .login-title {
    font-size: 26px;
  }

  .circle-1,
  .circle-2,
  .circle-3 {
    display: none;
  }

  .captcha-container {
    width: 80px;
    height: 36px;
  }
}
</style>
