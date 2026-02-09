<template>
  <div class="login-container">
    <div class="login-background">
      <div class="bg-grid"></div>
      <div class="bg-gradient"></div>
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
      <div class="bg-shape shape-3"></div>
      <div class="bg-shape shape-4"></div>
    </div>
    <div class="login-box">
      <div class="login-header container">
        <!-- <div class="logo-wrapper left">
          <el-icon class="logo-icon">
            <Document />
          </el-icon>
        </div> -->
        <div class="right">
          <h1>博客管理系统</h1>
        </div>
      </div>

      <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" class="login-form">
        <el-form-item prop="username">
          <el-input v-model="loginForm.username" placeholder="请输入用户名" :prefix-icon="User" size="large" clearable />
        </el-form-item>
        <el-form-item prop="email">
          <el-input v-model="loginForm.email" placeholder="请输入邮箱" :prefix-icon="Message" size="large" clearable />
        </el-form-item>

        <el-form-item prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" :prefix-icon="Lock" size="large"
            show-password @keyup.enter="handleLogin" />
        </el-form-item>

        <el-form-item prop="code">
          <div class="code-input-group">
            <el-input v-model="loginForm.code" placeholder="请输入验证码" :prefix-icon="Key" size="large" maxlength="6"
              clearable />
            <el-button :disabled="codeCountdown > 0 || !loginForm.email" @click="sendCode" size="large"
              :loading="codeloading" class="code-btn">
              {{ codeCountdown > 0 ? `${codeCountdown}s` : '获取验证码' }}
            </el-button>
          </div>
        </el-form-item>

        <el-form-item>
          <div class="form-footer">
            <el-checkbox v-model="loginForm.remember" size="large">记住密码</el-checkbox>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" size="large" :loading="loading" @click="handleLogin" class="login-btn">
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Lock, Message, Key, User } from '@element-plus/icons-vue'
import { rbacLogin, sendEmailCode } from '@/api'
import { useRbacStore } from '@/store/rbac'

const router = useRouter()
const rbacStore = useRbacStore()
const loginFormRef = ref(null)
const loading = ref(false)
const codeloading = ref(false)
const codeCountdown = ref(0)

const loginForm = reactive({
  username: '',
  email: '',
  password: '',
  code: '',
  remember: false
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 4, max: 20, message: '请输入正确的用户名', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码长度为6位', trigger: 'blur' }
  ]
}

// 组件挂载时，从localStorage读取记住的密码
onMounted(() => {
  const savedLogin = localStorage.getItem('rememberedLogin')
  if (savedLogin) {
    const parsed = JSON.parse(savedLogin)
    if (parsed.remember) {
      loginForm.username = parsed.username
      loginForm.email = parsed.email
      loginForm.password = parsed.password
      loginForm.remember = true
    }
  }
})

const sendCode = () => {
  if (!loginForm.email) {
    ElMessage.warning('请先输入邮箱')
    return
  }

  loginFormRef.value.validateField('email', async (valid) => {
    if (!valid) return
    try {
      codeloading.value = true
      const data = await sendEmailCode(loginForm.email)
      if (data.code !== 200) {
        ElMessage.error(data.message)
        return
      }
      ElMessage.success('验证码已发送至您的邮箱')
    } catch (error) {
      ElMessage.error('发送验证码失败')
    } finally {
      codeloading.value = false
    }
    codeCountdown.value = 60
    const timer = setInterval(() => {
      codeCountdown.value--
      if (codeCountdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  })
}

const handleLogin = async () => {
  if (!loginFormRef.value) return

  await loginFormRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true

    try {
      const res = await rbacLogin({
        username: loginForm.username,
        password: loginForm.password,
        email: loginForm.email,
        code: loginForm.code
      })
      console.log("获取的数据为", res.data)

      const { token, user } = res.data

      // 记住密码
      if (loginForm.remember) {
        localStorage.setItem('rememberedLogin', JSON.stringify({
          username: loginForm.username,
          email: loginForm.email,
          password: loginForm.password,
          remember: true
        }))
      } else {
        localStorage.removeItem('rememberedLogin')
      }

      // 清理旧用户的菜单和权限数据
      const rbacStore = useRbacStore()
      console.log('登录成功，开始清理旧数据...')
      rbacStore.resetAllData()

      // 强制重置路由初始化标志,让路由守卫重新加载菜单
      sessionStorage.setItem('routeInitialized', 'false')

      sessionStorage.setItem('token', token)
      sessionStorage.setItem('userInfo', JSON.stringify(user))
      console.log("userInfo", user);


      // 注意：不要在这里加载菜单，让路由守卫来处理
      console.log('登录成功，准备跳转到 dashboard')

      ElMessage.success('登录成功')

      router.push('/dashboard')
    } catch (error) {
      console.error('登录失败:', error)
      ElMessage.error(error.response?.data?.message || '登录失败，请检查用户名、密码和验证码')
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped lang="scss">
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: url(../../public/images/b3.jpg) center center/contain no-repeat;
  background-size: cover;
  image-rendering: -webkit-optimize-contrast;
  image-rendering: crisp-edges;
  image-rendering: pixelated;

  .login-background {
    position: absolute;
    width: 100%;
    height: 100%;
    overflow: hidden;

    .bg-grid {
      display: none;
    }

    .bg-gradient {
      position: absolute;
      width: 100%;
      height: 100%;
      background: radial-gradient(circle at 30% 40%, rgba(24, 144, 255, 0.1) 0%, transparent 50%),
        radial-gradient(circle at 70% 60%, rgba(64, 169, 255, 0.08) 0%, transparent 50%),
        radial-gradient(circle at 50% 50%, rgba(24, 144, 255, 0.05) 0%, transparent 70%);
    }

    .bg-shape {
      position: absolute;
      border-radius: 50%;
      opacity: 0.6;
      animation: float 25s infinite ease-in-out;
      filter: blur(60px);

      &.shape-1 {
        width: 500px;
        height: 500px;
        background: linear-gradient(135deg, rgba(24, 144, 255, 0.2), rgba(64, 169, 255, 0.15));
        top: -150px;
        left: -150px;
        animation-delay: 0s;
        animation-duration: 30s;
      }

      &.shape-2 {
        width: 400px;
        height: 400px;
        background: linear-gradient(135deg, rgba(64, 169, 255, 0.15), rgba(102, 126, 234, 0.1));
        bottom: -100px;
        right: -100px;
        animation-delay: -5s;
        animation-duration: 25s;
      }

      &.shape-3 {
        width: 300px;
        height: 300px;
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.12), rgba(24, 144, 255, 0.08));
        top: 50%;
        left: 60%;
        animation-delay: -10s;
        animation-duration: 20s;
      }

      &.shape-4 {
        width: 250px;
        height: 250px;
        background: linear-gradient(135deg, rgba(24, 144, 255, 0.15), rgba(64, 169, 255, 0.1));
        top: 30%;
        left: 20%;
        animation-delay: -15s;
        animation-duration: 35s;
      }
    }
  }
}

@keyframes float {

  0%,
  100% {
    transform: translate(0, 0) rotate(0deg) scale(1);
  }

  25% {
    transform: translate(40px, -40px) rotate(90deg) scale(1.05);
  }

  50% {
    transform: translate(20px, 40px) rotate(180deg) scale(0.95);
  }

  75% {
    transform: translate(-30px, 20px) rotate(270deg) scale(1.02);
  }
}

.login-box {
  width: 420px;
  padding: 32px 40px;
  // background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(5px);
  -webkit-backdrop-filter: blur(5px);
  border-radius: 15px;
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.25),
    0 8px 16px -8px rgba(0, 0, 0, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.6);
  position: relative;
  z-index: 8;
  border: 1px solid rgba(255, 255, 255, 0.4);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);

  &:hover {
    box-shadow:
      0 30px 60px -12px rgba(0, 0, 0, 0.3),
      0 10px 20px -8px rgba(0, 0, 0, 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.8);
    transform: translateY(-2px);
  }

  :deep(.el-input__prefix) {
    margin-left: 8px;
  }

  .login-header.container {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 20px;
    margin-bottom: 28px;

    .logo-wrapper.left {
      flex-shrink: 0;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 68px;
      height: 68px;
      background: rgb(227, 223, 223);
      border-radius: 20px;
      box-shadow: none;
      transition: all 0.3s ease;

      &:hover {
        transform: translateY(-4px) rotate(5deg);
        box-shadow: none;
      }

      .logo-icon {
        font-size: 34px;
        color: #fafafa;
        filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
      }
    }

    .right {
      display: flex;
      flex-direction: column;
      justify-content: center;
    }

    h1 {
      font-size: 28px;
      color: #ffffff;
      margin-bottom: 6px;
      font-weight: 700;
      letter-spacing: 0.5px;
    }
  }

  .login-form {
    :deep(.el-form-item) {
      margin-bottom: 17px;
    }

    .form-footer {
      display: flex;
      justify-content: flex-start;
      align-items: center;
      width: 100%;

      :deep(.el-checkbox) {
        .el-checkbox__label {
          font-size: 14px;
          color: #595959;
        }
      }
    }

    .code-input-group {
      display: flex;
      gap: 12px;
      width: 100%;
      height: 40px;

      .el-input {
        flex: 1;
      }

      .el-button.el-button--large.code-btn {
        width: 120px;
        height: 40px;
        font-size: 14px;
        border: 1px solid #d9d9d9;
        color: #595959;
        background: #fff;
        padding: 0;
        border-radius: 5px;

        &:hover:not(:disabled) {
          color: #1890ff;
          border-color: #1890ff;
          background: #f0f7ff;
        }

        &:disabled {
          color: #bfbfbf;
          background: #f5f5f5;
          border-color: #e8e8e8;
        }
      }
    }

    .login-btn {
      width: 100%;
      height: 44px;
      font-size: 16px;
      font-weight: 500;
      background: white;
      border: none;
      color: black;
      border-radius: 22px;
      transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
      box-shadow: 0 4px 12px rgba(24, 144, 255, 0.25);

      &:hover:not(:disabled) {
        background: linear-gradient(135deg, #6285a2, #5d768d);
        transform: translateY(-2px);
        box-shadow: 0 6px 16px rgba(24, 144, 255, 0.35);
      }

      &:active {
        transform: translateY(0);
        box-shadow: 0 2px 8px rgba(24, 144, 255, 0.2);
      }
    }
  }
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  padding: 2px;
  transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  background: transparent;
  border: 1px solid #ffffff;

  &:hover {
    border-color: #ffffff;
    box-shadow: 0 2px 6px rgba(24, 144, 255, 0.1);
  }

  &.is-focus {
    border-color: #ffffff;
    box-shadow: 0 0 0 3px rgba(24, 144, 255, 0.1), 0 2px 6px rgba(24, 144, 255, 0.15);
  }

  .el-input__inner {
    font-size: 14px;
    color: #fafafa;
    background: transparent;

    &::placeholder {
      color: rgba(255, 255, 255, 0.6);
    }
  }
}

:deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: #1890ff;
  border-color: #1890ff;
}

:deep(.el-checkbox__input.is-indeterminate .el-checkbox__inner) {
  background-color: #1890ff;
  border-color: #1890ff;
}
</style>
