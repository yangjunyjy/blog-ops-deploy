<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Message, ArrowLeft } from '@element-plus/icons-vue'

const router = useRouter()

const email = ref('')
const loading = ref(false)
const isSent = ref(false)

const handleSendCode = async () => {
  if (!email.value) {
    ElMessage.warning('请输入邮箱地址')
    return
  }

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
    ElMessage.warning('请输入有效的邮箱地址')
    return
  }

  loading.value = true
  try {
    // 模拟发送验证码
    await new Promise(resolve => setTimeout(resolve, 1500))

    isSent.value = true
    ElMessage.success('验证码已发送到您的邮箱，请查收')
  } catch (error) {
    ElMessage.error('发送失败，请重�?)
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.push('/login')
}

const countdown = ref(60)
let timer = null
</script>

<template>
  <div class="forgot-page">
    <div class="forgot-container">
      <div class="forgot-card">
        <div class="forgot-header">
          <el-button text class="back-btn" @click="goBack">
            <el-icon><ArrowLeft /></el-icon>
            返回登录
          </el-button>
          <h1 class="forgot-title">忘记密码</h1>
          <p class="forgot-subtitle">
            {{ isSent ? '验证码已发送，请注意查收~' : '输入邮箱，我们帮您找回密�? }}
          </p>
        </div>

        <div class="forgot-form" v-if="!isSent">
          <div class="form-item">
            <el-icon class="form-icon"><Message /></el-icon>
            <el-input
              v-model="email"
              placeholder="请输入您的邮箱地址"
              size="large"
              clearable
            />
          </div>

          <el-button
            type="primary"
            size="large"
            class="send-button"
            :loading="loading"
            @click="handleSendCode"
          >
            <span v-if="!loading">发送验证码</span>
            <span v-else>发送中...</span>
          </el-button>
        </div>

        <div class="sent-success" v-else>
          <div class="success-icon">✉️</div>
          <div class="success-text">
            <h3>验证码已发送！</h3>
            <p>请前往 <strong>{{ email }}</strong> 查收验证�?/p>
            <p class="tip">验证码有效期�?5 分钟，请及时查看~</p>
          </div>

          <div class="resend-section">
            <el-button
              text
              :disabled="countdown > 0"
              @click="isSent = false"
            >
              {{ countdown > 0 ? `重新发�?(${countdown}s)` : '重新发�? }}
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.forgot-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: url('/images/background.jpeg') center center / cover no-repeat;
  position: relative;
  overflow: hidden;
  padding: 40px 20px;
}

.forgot-page::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  pointer-events: none;
}

.forgot-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 480px;
}

.forgot-card {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 24px;
  padding: 45px 35px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.forgot-header {
  text-align: center;
  margin-bottom: 35px;
}

.back-btn {
  position: absolute;
  left: 20px;
  top: 20px;
  padding: 8px 12px;
  color: #718096;
  font-weight: 500;
}

.back-btn:hover {
  color: #1890ff;
}

.forgot-title {
  font-size: 32px;
  font-weight: 800;
  color: #1890ff;
  margin: 0 0 16px;
  letter-spacing: -0.5px;
}

.forgot-subtitle {
  color: #718096;
  font-size: 15px;
  margin: 0;
  opacity: 0.8;
}

.forgot-form {
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

.send-button {
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

.send-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(24, 144, 255, 0.4);
}

.sent-success {
  text-align: center;
  padding: 20px 0;
}

.success-icon {
  font-size: 80px;
  margin-bottom: 24px;
}

.success-text h3 {
  font-size: 22px;
  font-weight: 700;
  color: #2d3748;
  margin: 0 0 12px;
}

.success-text p {
  color: #718096;
  font-size: 15px;
  margin: 0 0 8px;
  line-height: 1.6;
}

.success-text .tip {
  font-size: 13px;
  color: #a0aec0;
  margin-top: 16px;
}

.resend-section {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #e2e8f0;
}

.resend-section :deep(.el-button) {
  color: #1890ff;
  font-weight: 600;
  font-size: 15px;
}

.resend-section :deep(.el-button:disabled) {
  color: #a0aec0;
}

@media (max-width: 768px) {
  .forgot-page {
    padding: 20px 15px;
  }

  .forgot-card {
    padding: 35px 25px;
    border-radius: 20px;
  }

  .forgot-title {
    font-size: 26px;
  }

  .cloud,
  .raindrop {
    display: none;
  }

  .back-btn {
    position: static;
    display: block;
    margin-bottom: 16px;
    width: fit-content;
  }
}
</style>
