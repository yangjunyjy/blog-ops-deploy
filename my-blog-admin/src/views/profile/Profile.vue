<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <span>个人中心</span>
        </div>
      </template>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="个人信息" name="info">
          <el-form ref="infoFormRef" :model="userInfo" :rules="infoRules" label-width="100px" class="profile-form">
            <el-form-item label="头像">
              <div class="avatar-uploader">
                <el-upload class="avatar-uploader-inner" :show-file-list="false" :before-upload="beforeAvatarUpload"
                  :on-change="handleAvatarChange" :auto-upload="false">
                  <el-avatar :size="100" :src="userInfo.avatar" class="avatar-preview" />
                  <template #tip>
                    <div class="el-upload__tip">支持 jpg/png 格式，文件大小不超过 2MB</div>
                  </template>
                </el-upload>
              </div>
            </el-form-item>
            <el-form-item label="用户名" prop="username">
              <el-input v-model="userInfo.username" disabled />
            </el-form-item>
            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="userInfo.nickname" placeholder="请输入昵称" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="userInfo.email" placeholder="请输入邮箱" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleUpdateInfo">保存修改</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="我的文章" name="articles">
          <div class="stats-cards">
            <el-row :gutter="20">
              <el-col :span="8">
                <div class="stat-card">
                  <div class="stat-value">{{ stats.articles.total }}</div>
                  <div class="stat-label">文章总数</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="stat-card">
                  <div class="stat-value">{{ stats.articles.published }}</div>
                  <div class="stat-label">已发布</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="stat-card">
                  <div class="stat-value">{{ stats.articles.draft }}</div>
                  <div class="stat-label">草稿</div>
                </div>
              </el-col>
            </el-row>
          </div>
          <el-divider />
          <div class="article-list">
            <el-empty v-if="myArticles.length === 0" description="暂无文章" />
            <div v-else class="article-items">
              <div v-for="article in myArticles" :key="article.id" class="article-item">
                <div class="article-info">
                  <h4>{{ article.title }}</h4>
                  <p>{{ article.summary }}</p>
                  <div class="article-meta">
                    <span class="meta-item">
                      <el-icon>
                        <View />
                      </el-icon>
                      {{ article.views || 0 }}
                    </span>
                    <span class="meta-item">
                      <el-icon>
                        <ChatDotRound />
                      </el-icon>
                      {{ article.comments || 0 }}
                    </span>
                    <span class="meta-item">
                      <el-icon>
                        <Clock />
                      </el-icon>
                      {{ formatDate(article.createdAt) }}
                    </span>
                    <el-tag v-if="article.status === 'published'" type="success" size="small">已发布</el-tag>
                    <el-tag v-else type="info" size="small">草稿</el-tag>
                  </div>
                </div>
                <div class="article-actions">
                  <el-button type="primary" link size="small" @click="editArticle(article.id)">编辑</el-button>
                </div>
              </div>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="我的系列" name="series">
          <div class="stats-cards">
            <el-row :gutter="20">
              <el-col :span="8">
                <div class="stat-card">
                  <div class="stat-value">{{ stats.series.total }}</div>
                  <div class="stat-label">系列总数</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="stat-card">
                  <div class="stat-value">{{ stats.series.articles }}</div>
                  <div class="stat-label">包含文章</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="stat-card">
                  <div class="stat-value">{{ stats.series.views }}</div>
                  <div class="stat-label">总浏览量</div>
                </div>
              </el-col>
            </el-row>
          </div>
          <el-divider />
          <div class="series-list">
            <el-empty v-if="mySeries.length === 0" description="暂无系列" />
            <div v-else class="series-items">
              <el-card v-for="series in mySeries" :key="series.id" class="series-item" shadow="hover">
                <div class="series-info">
                  <h4>{{ series.name }}</h4>
                  <p>{{ series.description }}</p>
                  <div class="series-meta">
                    <span class="meta-item">
                      <el-icon>
                        <Document />
                      </el-icon>
                      {{ series.articleCount || 0 }} 篇文章
                    </span>
                    <span class="meta-item">
                      <el-icon>
                        <View />
                      </el-icon>
                      {{ series.views || 0 }} 次浏览
                    </span>
                  </div>
                </div>
              </el-card>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="修改密码" name="password" v-model="authDialog">
          <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="100px"
            class="profile-form">
            <el-form-item label="原密码" prop="oldPassword">
              <el-input v-model="passwordForm.oldPassword" type="password" placeholder="请输入原密码" show-password />
            </el-form-item>
            <el-form-item label="新密码" prop="newPassword">
              <el-input v-model="passwordForm.newPassword" type="password" placeholder="请输入新密码（至少6位）" show-password />
            </el-form-item>
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="请再次输入新密码" show-password />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleUpdatePassword">修改密码</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
    <el-dialog v-model="authDialog" title="邮箱二次认证" width="500" class="auth-dialog">
      <el-form ref="authFormRef" :model="authForm" :rules="authRules" label-width="100px">
        <el-form-item label="验证码" prop="code">
          <el-input v-model="authForm.code" placeholder="请输入邮箱验证码">
            <template #append>
              <el-button 
                type="primary" 
                :disabled="countdown > 0"
                @click="handleSendCode"
                :loading="sendingCode">
                {{ countdown > 0 ? `${countdown}秒后重试` : '发送验证码' }}
              </el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="handleConfirmAuth" :loading="confirmingAuth">
            确认修改
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  View,
  ChatDotRound,
  Clock,
  Document
} from '@element-plus/icons-vue'
import { useRbacStore } from '../../store/rbac'

const router = useRouter()
const rbacStore = useRbacStore()
const activeTab = ref('info')
const infoFormRef = ref()
const passwordFormRef = ref()
const authFormRef = ref()
const authDialog = ref(false)
const sendingCode = ref(false)
const confirmingAuth = ref(false)
const countdown = ref(0)

const userInfo = reactive({
  id: null,
  username: '',
  nickname: '',
  email: '',
  avatar: '',
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const authForm = reactive({
  code: ''
})

const stats = reactive({
  articles: {
    total: 0,
    published: 0,
    draft: 0
  },
  series: {
    total: 0,
    articles: 0,
    views: 0
  }
})

const myArticles = ref([])
const mySeries = ref([])


const infoRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const authRules = {
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码为6位数字', trigger: 'blur' }
  ]
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const loadUserInfo = () => {
  const storedInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')

  // 确保 users 存在且是数组
  if (!rbacStore.users.value || !Array.isArray(rbacStore.users.value)) {
    // 使用 sessionStorage 中的 userInfo 作为 fallback
    const sessionInfo = JSON.parse(sessionStorage.getItem('userInfo') || '{}')
    if (sessionInfo.email) {
      Object.assign(userInfo, {
        id: sessionInfo.id,
        username: sessionInfo.username,
        nickname: sessionInfo.nickname || sessionInfo.username,
        email: sessionInfo.email,
        avatar: sessionInfo.avatar || 'https://api.dicebear.com/7.x/initials/svg?seed=Admin&backgroundColor=1890ff',
        role: sessionInfo.role || 'visitor'
      })
    }
    return
  }

  const user = rbacStore.users.value.find(u => u.email === storedInfo.email) || rbacStore.users.value[0]
  if (user) {
    Object.assign(userInfo, {
      id: user.id,
      username: user.username,
      nickname: user.nickname,
      email: user.email,
      avatar: user.avatar,
      role: user.role || 'visitor'
    })
  }
}

const loadUserStats = () => {
  // 模拟数据 - 实际项目中应该从后端获取
  if (true) {
    stats.articles.total = 12
    stats.articles.published = 8
    stats.articles.draft = 4

    myArticles.value = [
      {
        id: 1,
        title: 'Vue3 组合式 API 入门指南',
        summary: '详细介绍 Vue3 组合式 API 的使用方法和最佳实践...',
        views: 1234,
        comments: 56,
        status: 'published',
        createdAt: new Date('2024-01-15')
      },
      {
        id: 2,
        title: 'TypeScript 进阶教程',
        summary: '深入学习 TypeScript 的高级特性和设计模式...',
        views: 892,
        comments: 23,
        status: 'published',
        createdAt: new Date('2024-01-10')
      },
      {
        id: 3,
        title: '前端性能优化实践',
        summary: '总结前端性能优化的各种技巧和工具...',
        views: 0,
        comments: 0,
        status: 'draft',
        createdAt: new Date('2024-01-18')
      }
    ]
  }

  if (true) {
    stats.series.total = 3
    stats.series.articles = 12
    stats.series.views = 5678

    mySeries.value = [
      {
        id: 1,
        name: 'Vue3 完全指南',
        description: '从零开始学习 Vue3，掌握现代前端开发',
        articleCount: 5,
        views: 2345
      },
      {
        id: 2,
        name: 'TypeScript 实战',
        description: 'TypeScript 实战项目案例分享',
        articleCount: 4,
        views: 1890
      },
      {
        id: 3,
        name: '前端工程化',
        description: '现代前端工程化实践',
        articleCount: 3,
        views: 1443
      }
    ]
  }
}

const beforeAvatarUpload = (file) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 3

  if (!isJPG) {
    ElMessage.error('头像图片只能是 JPG/PNG 格式!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('头像图片大小不能超过 3MB!')
    return false
  }
  return true
}

const handleAvatarChange = (file) => {
  if (beforeAvatarUpload(file.raw)) {
    const reader = new FileReader()
    reader.onload = (e) => {
      userInfo.avatar = e.target.result
    }
    reader.readAsDataURL(file.raw)
  }
}

const handleUpdateInfo = async () => {
  if (!infoFormRef.value) return

  await infoFormRef.value.validate((valid) => {
    if (valid) {
      const userIndex = rbacStore.users.value.findIndex(u => u.id === userInfo.id)
      if (userIndex !== -1) {
        rbacStore.users.value[userIndex].nickname = userInfo.nickname
        rbacStore.users.value[userIndex].email = userInfo.email
        rbacStore.users.value[userIndex].avatar = userInfo.avatar
        rbacStore.saveUsers()

        localStorage.setItem('userInfo', JSON.stringify({
          ...JSON.parse(localStorage.getItem('userInfo') || '{}'),
          nickname: userInfo.nickname,
          email: userInfo.email,
          avatar: userInfo.avatar
        }))

        ElMessage.success('个人信息修改成功')
      }
    }
  })
}

const handleUpdatePassword = async () => {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate((valid) => {
    if (valid) {
      authDialog.value = true
      // 自动发送验证码
      handleSendCode()
    }
  })
}

const handleSendCode = async () => {
  try {
    sendingCode.value = true

    // 模拟发送验证码（实际项目中应该调用后端API）
    await new Promise(resolve => setTimeout(resolve, 1000))

    ElMessage.success(`验证码已发送至 ${userInfo.email}`)

    // 开始倒计时
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)

  } catch (error) {
    ElMessage.error('发送验证码失败，请重试')
  } finally {
    sendingCode.value = false
  }
}

const handleConfirmAuth = async () => {
  if (!authFormRef.value) return

  await authFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        confirmingAuth.value = true

        // 模拟验证验证码并修改密码（实际项目中应该调用后端API）
        await new Promise(resolve => setTimeout(resolve, 1000))

        // 验证成功，更新用户密码（本地模拟）
        const userIndex = rbacStore.users.value.findIndex(u => u.id === userInfo.id)
        if (userIndex !== -1) {
          // 在实际项目中，这里应该调用后端API修改密码
          ElMessage.success('密码修改成功')
          authDialog.value = false

          // 重置表单
          passwordForm.oldPassword = ''
          passwordForm.newPassword = ''
          passwordForm.confirmPassword = ''
          authForm.code = ''
          countdown.value = 0
        }
      } catch (error) {
        ElMessage.error('验证码错误或密码修改失败')
      } finally {
        confirmingAuth.value = false
      }
    }
  })
}

const editArticle = (id) => {
  router.push(`/editor/${id}`)
}

onMounted(() => {
  loadUserInfo()
  loadUserStats()
})
</script>

<style scoped lang="scss">
.profile-container {
  position: relative;
  padding: 0px;

  .profile-card {
    width: 100%;
    margin: 0;
    height: 100%;
    min-height: 100vh;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-size: 18px;
      font-weight: 600;

      .role-tag {
        font-size: 14px;
      }
    }

    .profile-form {
      max-width: 600px;
      margin-top: 20px;

      .avatar-uploader {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 12px;

        .avatar-uploader-inner {
          display: flex;
          flex-direction: column;
          align-items: center;

          .avatar-preview {
            border: 2px dashed #d9d9d9;
            cursor: pointer;
            transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);

            &:hover {
              border-color: #1890ff;
              transform: scale(1.02);
              box-shadow: 0 2px 8px rgba(24, 144, 255, 0.15);
            }
          }

          .el-upload__tip {
            font-size: 12px;
            color: #8c8c8c;
            margin-top: 8px;
          }
        }
      }
    }

    .stats-cards {
      margin-bottom: 20px;

      .stat-card {
        background: #f5f7fa;
        border-radius: 2px;
        padding: 24px;
        text-align: center;
        transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
        }

        .stat-value {
          font-size: 32px;
          font-weight: 600;
          color: #1890ff;
          margin-bottom: 8px;
        }

        .stat-label {
          font-size: 14px;
          color: #606266;
        }
      }
    }

    .article-list,
    .series-list {
      min-height: 400px;

      .article-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 16px;
        border-bottom: 1px solid #f0f0f0;
        transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);

        &:hover {
          background: #f5f7fa;
        }

        &:last-child {
          border-bottom: none;
        }

        .article-info {
          flex: 1;

          h4 {
            margin: 0 0 8px;
            font-size: 16px;
            font-weight: 500;
            color: #303133;
          }

          p {
            margin: 0 0 12px;
            font-size: 14px;
            color: #909399;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            max-width: 600px;
          }

          .article-meta {
            display: flex;
            align-items: center;
            gap: 16px;

            .meta-item {
              display: flex;
              align-items: center;
              gap: 4px;
              font-size: 12px;
              color: #909399;
            }
          }
        }

        .article-actions {
          flex-shrink: 0;
        }
      }

      .series-items {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 16px;

        .series-item {
          cursor: pointer;
          transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);

          &:hover {
            transform: translateY(-2px);
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
          }

          .series-info {
            h4 {
              margin: 0 0 8px;
              font-size: 16px;
              font-weight: 500;
              color: #303133;
            }

            p {
              margin: 0 0 12px;
              font-size: 14px;
              color: #909399;
              min-height: 40px;
            }

            .series-meta {
              display: flex;
              align-items: center;
              gap: 16px;

              .meta-item {
                display: flex;
                align-items: center;
                gap: 4px;
                font-size: 12px;
                color: #909399;
              }
            }
          }
        }
      }
    }
  }
}
</style>

<style lang="scss">
// 对话框居中样式（不使用 scoped）
.auth-dialog.el-dialog {
  display: flex !important;
  flex-direction: column !important;
  margin: 0 !important;
  position: absolute !important;
  top: 50% !important;
  left: 50% !important;
  transform: translate(-50%, -50%) !important;
  max-height: calc(100vh - 30px) !important;
  max-width: calc(100vw - 30px) !important;

  .el-dialog__body {
    flex: 1;
    overflow: auto;
  }
}
</style>
