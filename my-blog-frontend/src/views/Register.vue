<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Message, Check, Key } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { registerWithEmailCode, sendEmailCode } from '@/api'

const router = useRouter()
const userStore = useUserStore()

const registerForm = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  code: ''
})

const loading = ref(false)
const codeLoading = ref(false)
const countdown = ref(0)

const handleSendCode = async () => {
  if (!registerForm.value.email) {
    ElMessage.warning('请输入邮箱地址')
    return
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(registerForm.value.email)) {
    ElMessage.warning('请输入有效的邮箱地址')
    return
  }

  if (countdown.value > 0) {
    return
  }

  codeLoading.value = true
  try {
    await sendEmailCode(registerForm.value.email, registerForm.value.username)
    ElMessage.success('验证码已发送到您的邮箱')

    // 开始倒计�?
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '验证码发送失败，请重试')
  } finally {
    codeLoading.value = false
  }
}

const handleRegister = async () => {
  if (!registerForm.value.username || !registerForm.value.email || !registerForm.value.password) {
    ElMessage.warning('请填写所有必填项')
    return
  }

  if (!registerForm.value.code) {
    ElMessage.warning('请输入邮箱验证码')
    return
  }

  if (registerForm.value.password !== registerForm.value.confirmPassword) {
    ElMessage.warning('两次密码输入不一致')
    return
  }

  if (registerForm.value.password.length < 6) {
    ElMessage.warning('密码长度至少6位')
    return
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(registerForm.value.email)) {
    ElMessage.warning('请输入有效的邮箱地址')
    return
  }

  loading.value = true
  try {
    const res = await registerWithEmailCode({
      username: registerForm.value.username,
      email: registerForm.value.email,
      password: registerForm.value.password,
      code: registerForm.value.code
    })

    if (res.code === 200) {
      userStore.login(res.data.user, res.data.token)
      ElMessage.success('注册成功！欢迎加入我们')
      await router.push('/login')
    }
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '注册失败，请重试')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-card">
      <div class="register-header">
        <h1 class="register-title">加入我们</h1>
        <p class="register-subtitle">开启你的博客之旅</p>
      </div>

        <el-form class="register-form" @submit.prevent="handleRegister">
          <div class="form-item">
            <el-icon class="form-icon"><User /></el-icon>
            <el-input
              v-model="registerForm.username"
              placeholder="用户名"
              size="large"
              clearable
            />
          </div>

          <div class="form-item">
            <el-icon class="form-icon"><Message /></el-icon>
            <el-input
              v-model="registerForm.email"
              type="email"
              placeholder="邮箱地址"
              size="large"
              clearable
            />
          </div>

          <div class="form-item">
            <el-icon class="form-icon"><Key /></el-icon>
            <el-input
              v-model="registerForm.code"
              placeholder="邮箱验证码"
              size="large"
              clearable
              maxlength="6"
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

          <div class="form-item">
            <el-icon class="form-icon"><Lock /></el-icon>
            <el-input
              v-model="registerForm.password"
              type="password"
              placeholder="密码（至少8位）"
              size="large"
              show-password
              clearable
            />
          </div>

          <div class="form-item">
            <el-icon class="form-icon"><Check /></el-icon>
            <el-input
              v-model="registerForm.confirmPassword"
              type="password"
              placeholder="确认密码"
              size="large"
              show-password
              clearable
            />
          </div>

          <el-button
            type="primary"
            size="large"
            class="register-button"
            :loading="loading"
            @click="handleRegister"
          >
            <span v-if="!loading">🎉 注册</span>
            <span v-else>注册中...</span>
          </el-button>

          <div class="register-footer">
            <p>已有账号?</p>
            <router-link to="/login" class="login-link">
              快去登录吧?
            </router-link>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: url('/images/background.jpeg') center center / cover no-repeat;
  position: relative;
  overflow: hidden;
  padding: 40px 20px;
}

.register-page::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  pointer-events: none;
}

.register-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 480px;
}

.register-card {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 24px;
  padding: 45px 35px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.register-header {
  text-align: center;
  margin-bottom: 35px;
}

.register-title {
  font-size: 32px;
  font-weight: 800;
  color: #1890ff;
  margin: 0 0 16px;
  letter-spacing: -0.5px;
}

.register-subtitle {
  color: #718096;
  font-size: 15px;
  margin: 0;
  opacity: 0.8;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
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

.register-button {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  border: none;
  background: #1890ff;
  border-radius: 12px;
  margin-top: 10px;
  transition: all 0.3s;
}

.register-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(24, 144, 255, 0.4);
}

.register-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  color: #718096;
}

.register-footer p {
  margin: 0 0 4px;
}

.login-link {
  color: #1890ff;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s;
}

.login-link:hover {
  color: #40a9ff;
  text-decoration: underline;
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

@media (max-width: 768px) {
  .register-page {
    padding: 20px 15px;
  }

  .register-card {
    padding: 35px 25px;
    border-radius: 20px;
  }

  .register-title {
    font-size: 26px;
  }
}

</style>
