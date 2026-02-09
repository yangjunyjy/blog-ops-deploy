<template>
  <div class="settings-container">
    <el-card class="settings-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon class="header-icon"><Setting /></el-icon>
            <span>系统设置</span>
          </div>
          <el-button type="primary" @click="handleSaveAll" :loading="saving">
            <el-icon><Check /></el-icon>
            保存全部
          </el-button>
        </div>
      </template>
      <el-tabs v-model="activeTab" class="settings-tabs">
        <el-tab-pane name="basic">
          <template #label>
            <span class="tab-label">
              <el-icon><House /></el-icon>
              网站信息
            </span>
          </template>
          <el-form
            ref="basicFormRef"
            :model="basicForm"
            :rules="basicRules"
            label-width="120px"
            class="settings-form"
          >
            <div class="form-section">
              <div class="section-title">
                <el-icon><InfoFilled /></el-icon>
                基本信息
              </div>
              <el-form-item label="网站标题" prop="siteTitle">
                <el-input v-model="basicForm.siteTitle" placeholder="请输入网站标题" clearable />
              </el-form-item>
              <el-form-item label="网站描述" prop="siteDescription">
                <el-input
                  v-model="basicForm.siteDescription"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入网站描述"
                  maxlength="200"
                  show-word-limit
                />
              </el-form-item>
              <el-form-item label="网站关键词" prop="siteKeywords">
                <el-input
                  v-model="basicForm.siteKeywords"
                  placeholder="多个关键词用逗号分隔"
                  clearable
                />
              </el-form-item>
            </div>
            <div class="form-section">
              <div class="section-title">
                <el-icon><Picture /></el-icon>
                视觉设置
              </div>
              <el-form-item label="网站Logo" prop="siteLogo">
                <el-input v-model="basicForm.siteLogo" placeholder="请输入Logo URL" clearable>
                  <template #prepend>
                    <el-icon><Link /></el-icon>
                  </template>
                </el-input>
              </el-form-item>
              <el-form-item label="ICP备案号" prop="icpNumber">
                <el-input v-model="basicForm.icpNumber" placeholder="请输入ICP备案号" clearable />
              </el-form-item>
            </div>
          </el-form>
        </el-tab-pane>

        <el-tab-pane name="seo">
          <template #label>
            <span class="tab-label">
              <el-icon><Search /></el-icon>
              SEO设置
            </span>
          </template>
          <el-form
            ref="seoFormRef"
            :model="seoForm"
            label-width="120px"
            class="settings-form"
          >
            <div class="form-section">
              <div class="section-title">
                <el-icon><MagicStick /></el-icon>
                优化开关
              </div>
              <div class="switch-group">
                <el-form-item label="启用SEO">
                  <el-switch v-model="seoForm.enableSeo" active-text="开启" inactive-text="关闭" />
                </el-form-item>
                <el-form-item label="结构化数据">
                  <el-switch v-model="seoForm.enableStructuredData" active-text="开启" inactive-text="关闭" />
                  <span class="form-tip">启用JSON-LD结构化数据</span>
                </el-form-item>
                <el-form-item label="Sitemap">
                  <el-switch v-model="seoForm.enableSitemap" active-text="开启" inactive-text="关闭" />
                  <span class="form-tip">自动生成Sitemap</span>
                </el-form-item>
              </div>
            </div>
            <div class="form-section">
              <div class="section-title">
                <el-icon><Document /></el-icon>
                首页SEO
              </div>
              <el-form-item label="首页标题" prop="homeTitle">
                <el-input v-model="seoForm.homeTitle" placeholder="默认使用网站标题" clearable />
              </el-form-item>
              <el-form-item label="首页描述" prop="homeDescription">
                <el-input
                  v-model="seoForm.homeDescription"
                  type="textarea"
                  :rows="3"
                  placeholder="默认使用网站描述"
                  maxlength="200"
                  show-word-limit
                />
              </el-form-item>
              <el-form-item label="首页关键词" prop="homeKeywords">
                <el-input
                  v-model="seoForm.homeKeywords"
                  placeholder="默认使用网站关键词"
                  clearable
                />
              </el-form-item>
            </div>
          </el-form>
        </el-tab-pane>

        <el-tab-pane name="system">
          <template #label>
            <span class="tab-label">
              <el-icon><Tools /></el-icon>
              系统配置
            </span>
          </template>
          <el-form
            ref="systemFormRef"
            :model="systemForm"
            label-width="120px"
            class="settings-form"
          >
            <div class="form-section">
              <div class="section-title">
                <el-icon><List /></el-icon>
                内容展示
              </div>
              <el-form-item label="每页文章数">
                <el-input-number v-model="systemForm.pageSize" :min="5" :max="50" controls-position="right" />
              </el-form-item>
            </div>
            <div class="form-section">
              <div class="section-title">
                <el-icon><ChatDotRound /></el-icon>
                互动设置
              </div>
              <div class="switch-group">
                <el-form-item label="评论审核">
                  <el-switch v-model="systemForm.commentModeration" active-text="开启" inactive-text="关闭" />
                  <span class="form-tip">新评论需要审核后显示</span>
                </el-form-item>
                <el-form-item label="允许评论">
                  <el-switch v-model="systemForm.enableComment" active-text="允许" inactive-text="禁止" />
                </el-form-item>
              </div>
            </div>
            <div class="form-section">
              <div class="section-title">
                <el-icon><Operation /></el-icon>
                功能开关
              </div>
              <div class="switch-group">
                <el-form-item label="文章审核">
                  <el-switch v-model="systemForm.articleModeration" active-text="开启" inactive-text="关闭" />
                  <span class="form-tip">新文章需要审核后发布</span>
                </el-form-item>
                <el-form-item label="注册开关">
                  <el-switch v-model="systemForm.enableRegister" active-text="允许" inactive-text="禁止" />
                  <span class="form-tip">允许用户注册</span>
                </el-form-item>
              </div>
            </div>
            <div class="form-section">
              <div class="section-title">
                <el-icon><Clock /></el-icon>
                时区设置
              </div>
              <el-form-item label="时区设置">
                <el-select v-model="systemForm.timezone" placeholder="选择时区" style="width: 100%">
                  <el-option label="Asia/Shanghai (UTC+8)" value="Asia/Shanghai" />
                  <el-option label="Asia/Tokyo (UTC+9)" value="Asia/Tokyo" />
                  <el-option label="America/New_York (UTC-5)" value="America/New_York" />
                  <el-option label="Europe/London (UTC+0)" value="Europe/London" />
                </el-select>
              </el-form-item>
            </div>
          </el-form>
        </el-tab-pane>

        <el-tab-pane name="email">
          <template #label>
            <span class="tab-label">
              <el-icon><Message /></el-icon>
              邮件设置
            </span>
          </template>
          <el-form
            ref="emailFormRef"
            :model="emailForm"
            label-width="120px"
            class="settings-form"
          >
            <div class="form-section">
              <div class="section-title">
                <el-icon><Connection /></el-icon>
                邮件配置
              </div>
              <div class="switch-group">
                <el-form-item label="启用邮件">
                  <el-switch v-model="emailForm.enableEmail" active-text="开启" inactive-text="关闭" />
                </el-form-item>
                <el-form-item label="使用SSL">
                  <el-switch v-model="emailForm.useSSL" active-text="开启" inactive-text="关闭" />
                </el-form-item>
              </div>
              <el-divider />
              <el-form-item label="SMTP主机" prop="smtpHost">
                <el-input v-model="emailForm.smtpHost" placeholder="例如: smtp.gmail.com" clearable>
                  <template #prepend>
                    <el-icon><Monitor /></el-icon>
                  </template>
                </el-input>
              </el-form-item>
              <el-form-item label="SMTP端口" prop="smtpPort">
                <el-input-number v-model="emailForm.smtpPort" :min="1" :max="65535" controls-position="right" />
              </el-form-item>
              <el-form-item label="发件人邮箱" prop="fromEmail">
                <el-input v-model="emailForm.fromEmail" placeholder="请输入发件人邮箱" clearable>
                  <template #prepend>
                    <el-icon><User /></el-icon>
                  </template>
                </el-input>
              </el-form-item>
              <el-form-item label="发件人名称" prop="fromName">
                <el-input v-model="emailForm.fromName" placeholder="请输入发件人名称" clearable />
              </el-form-item>
              <el-form-item label="用户名" prop="smtpUser">
                <el-input v-model="emailForm.smtpUser" placeholder="请输入SMTP用户名" clearable />
              </el-form-item>
              <el-form-item label="密码" prop="smtpPass">
                <el-input
                  v-model="emailForm.smtpPass"
                  type="password"
                  placeholder="请输入SMTP密码"
                  show-password
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleTestEmail">
                  <el-icon><Promotion /></el-icon>
                  发送测试邮件
                </el-button>
              </el-form-item>
            </div>
          </el-form>
        </el-tab-pane>

        <el-tab-pane name="storage">
          <template #label>
            <span class="tab-label">
              <el-icon><FolderOpened /></el-icon>
              存储设置
            </span>
          </template>
          <el-form
            ref="storageFormRef"
            :model="storageForm"
            label-width="120px"
            class="settings-form"
          >
            <div class="form-section">
              <div class="section-title">
                <el-icon><Coin /></el-icon>
                存储配置
              </div>
              <el-form-item label="存储类型">
                <el-radio-group v-model="storageForm.storageType" class="storage-type-group">
                  <el-radio label="local" border>
                    <el-icon><Files /></el-icon>
                    本地存储
                  </el-radio>
                  <el-radio label="oss" border>
                    <el-icon><Cloudy /></el-icon>
                    阿里云OSS
                  </el-radio>
                  <el-radio label="cos" border>
                    <el-icon><Cloudy /></el-icon>
                    腾讯云COS
                  </el-radio>
                  <el-radio label="s3" border>
                    <el-icon><Cloudy /></el-icon>
                    AWS S3
                  </el-radio>
                </el-radio-group>
              </el-form-item>
              <el-form-item label="上传路径" prop="uploadPath">
                <el-input v-model="storageForm.uploadPath" placeholder="例如: /uploads" clearable>
                  <template #prepend>
                    <el-icon><Folder /></el-icon>
                  </template>
                </el-input>
              </el-form-item>
              <el-form-item label="文件大小限制">
                <el-input-number
                  v-model="storageForm.maxFileSize"
                  :min="1"
                  :max="100"
                  controls-position="right"
                />
                <span class="form-tip">单位: MB</span>
              </el-form-item>
              <el-form-item label="允许的文件类型" prop="allowedTypes">
                <el-input
                  v-model="storageForm.allowedTypes"
                  placeholder="例如: jpg,png,gif,pdf,doc,docx"
                  clearable
                />
              </el-form-item>
              <el-form-item label="图片压缩">
                <el-switch v-model="storageForm.enableImageCompress" active-text="开启" inactive-text="关闭" />
              </el-form-item>
            </div>
            <div v-if="storageForm.storageType === 'oss'" class="form-section oss-section">
              <div class="section-title">
                <el-icon><Lock /></el-icon>
                OSS配置
              </div>
              <el-form-item label="Endpoint">
                <el-input v-model="storageForm.ossEndpoint" placeholder="OSS Endpoint" clearable />
              </el-form-item>
              <el-form-item label="Bucket">
                <el-input v-model="storageForm.ossBucket" placeholder="Bucket名称" clearable />
              </el-form-item>
              <el-form-item label="AccessKey">
                <el-input v-model="storageForm.ossAccessKey" placeholder="AccessKey" clearable />
              </el-form-item>
              <el-form-item label="SecretKey">
                <el-input
                  v-model="storageForm.ossSecretKey"
                  type="password"
                  placeholder="SecretKey"
                  show-password
                />
              </el-form-item>
            </div>
          </el-form>
        </el-tab-pane>

        <el-tab-pane name="backup">
          <template #label>
            <span class="tab-label">
              <el-icon><Files /></el-icon>
              备份恢复
            </span>
          </template>
          <div class="backup-section">
            <el-card shadow="never" class="backup-card">
              <template #header>
                <div class="section-header">
                  <el-icon><Download /></el-icon>
                  <span>数据备份</span>
                </div>
              </template>
              <el-form :model="backupForm" label-width="120px">
                <el-form-item label="备份类型">
                  <el-radio-group v-model="backupForm.backupType">
                    <el-radio label="full" border>完整备份</el-radio>
                    <el-radio label="incremental" border>增量备份</el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-form-item label="自动备份">
                  <el-switch v-model="backupForm.autoBackup" active-text="开启" inactive-text="关闭" />
                </el-form-item>
                <el-form-item v-if="backupForm.autoBackup" label="备份频率">
                  <el-select v-model="backupForm.backupFrequency" placeholder="选择频率">
                    <el-option label="每天" value="daily" />
                    <el-option label="每周" value="weekly" />
                    <el-option label="每月" value="monthly" />
                  </el-select>
                </el-form-item>
                <el-form-item label="保留数量">
                  <el-input-number v-model="backupForm.keepBackups" :min="1" :max="30" controls-position="right" />
                  <span class="form-tip">个备份文件</span>
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="handleCreateBackup" :loading="backingUp">
                    <el-icon><FolderAdd /></el-icon>
                    立即备份
                  </el-button>
                </el-form-item>
              </el-form>
            </el-card>

            <el-card shadow="never" class="backup-card">
              <template #header>
                <div class="section-header">
                  <el-icon><List /></el-icon>
                  <span>备份列表</span>
                </div>
              </template>
              <el-table :data="backupList" stripe style="width: 100%">
                <el-table-column prop="name" label="备份文件名" min-width="200">
                  <template #default="{ row }">
                    <el-icon class="file-icon"><Document /></el-icon>
                    {{ row.name }}
                  </template>
                </el-table-column>
                <el-table-column prop="size" label="大小" width="100" align="center" />
                <el-table-column prop="date" label="备份时间" width="180" align="center" />
                <el-table-column prop="type" label="类型" width="80" align="center">
                  <template #default="{ row }">
                    <el-tag :type="row.type === '完整' ? 'primary' : 'success'" size="small">
                      {{ row.type }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="160" align="center" fixed="right">
                  <template #default="{ row }">
                    <el-button type="primary" link size="small" @click="handleRestore(row)">
                      <el-icon><RefreshRight /></el-icon>
                      恢复
                    </el-button>
                    <el-button type="danger" link size="small" @click="handleDeleteBackup(row)">
                      <el-icon><Delete /></el-icon>
                      删除
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-card>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
defineOptions({
  name: 'Settings'
})

import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Setting,
  Check,
  House,
  InfoFilled,
  Picture,
  Link,
  Search,
  MagicStick,
  Document,
  Tools,
  List,
  ChatDotRound,
  Operation,
  Clock,
  Message,
  Connection,
  Monitor,
  User,
  Promotion,
  FolderOpened,
  Coin,
  Files,
  Cloudy,
  Folder,
  Lock,
  Download,
  FolderAdd,
  RefreshRight,
  Delete
} from '@element-plus/icons-vue'

const activeTab = ref('basic')
const saving = ref(false)
const backingUp = ref(false)
const basicFormRef = ref()
const seoFormRef = ref()
const systemFormRef = ref()
const emailFormRef = ref()
const storageFormRef = ref()

const basicForm = reactive({
  siteTitle: '我的博客',
  siteDescription: '这是一个基于Vue3的个人博客系统',
  siteKeywords: '博客,Vue3,技术分享',
  siteLogo: '',
  icpNumber: ''
})

const seoForm = reactive({
  enableSeo: true,
  homeTitle: '',
  homeDescription: '',
  homeKeywords: '',
  enableStructuredData: true,
  enableSitemap: true
})

const systemForm = reactive({
  pageSize: 10,
  commentModeration: false,
  enableComment: true,
  articleModeration: false,
  enableRegister: true,
  timezone: 'Asia/Shanghai'
})

const emailForm = reactive({
  enableEmail: false,
  smtpHost: '',
  smtpPort: 587,
  fromEmail: '',
  fromName: 'MyBlog',
  smtpUser: '',
  smtpPass: '',
  useSSL: true
})

const storageForm = reactive({
  storageType: 'local',
  uploadPath: '/uploads',
  maxFileSize: 10,
  allowedTypes: 'jpg,png,gif,pdf,doc,docx,txt',
  enableImageCompress: true,
  ossEndpoint: '',
  ossBucket: '',
  ossAccessKey: '',
  ossSecretKey: ''
})

const backupForm = reactive({
  backupType: 'full',
  autoBackup: false,
  backupFrequency: 'weekly',
  keepBackups: 7
})

const backupList = ref([
  {
    id: 1,
    name: 'backup_20240115_full.zip',
    size: '15.2 MB',
    date: '2024-01-15 10:30:00',
    type: '完整'
  },
  {
    id: 2,
    name: 'backup_20240114_full.zip',
    size: '15.0 MB',
    date: '2024-01-14 10:30:00',
    type: '完整'
  }
])

const basicRules = {
  siteTitle: [
    { required: true, message: '请输入网站标题', trigger: 'blur' }
  ],
  siteDescription: [
    { required: true, message: '请输入网站描述', trigger: 'blur' }
  ],
  siteKeywords: [
    { required: true, message: '请输入网站关键词', trigger: 'blur' }
  ]
}

const loadSettings = () => {
  const stored = localStorage.getItem('system_settings')
  if (stored) {
    const settings = JSON.parse(stored)
    Object.assign(basicForm, settings.basic || {})
    Object.assign(seoForm, settings.seo || {})
    Object.assign(systemForm, settings.system || {})
    Object.assign(emailForm, settings.email || {})
    Object.assign(storageForm, settings.storage || {})
    Object.assign(backupForm, settings.backup || {})
  }
}

const saveSettings = () => {
  const settings = {
    basic: { ...basicForm },
    seo: { ...seoForm },
    system: { ...systemForm },
    email: { ...emailForm },
    storage: { ...storageForm },
    backup: { ...backupForm }
  }
  localStorage.setItem('system_settings', JSON.stringify(settings))
}

const handleSaveAll = async () => {
  try {
    saving.value = true
    await new Promise(resolve => setTimeout(resolve, 500))
    if (basicFormRef.value) {
      await basicFormRef.value.validate()
    }
    saveSettings()
    ElMessage.success('设置保存成功')
  } catch (error) {
    ElMessage.error('请检查表单填写是否正确')
  } finally {
    saving.value = false
  }
}

const handleTestEmail = () => {
  ElMessage.success('测试邮件发送成功，请查收')
}

const handleCreateBackup = () => {
  backingUp.value = true
  setTimeout(() => {
    backingUp.value = false
    ElMessage.success('备份任务已创建，正在后台执行...')
    const newBackup = {
      id: Date.now(),
      name: `backup_${new Date().toISOString().slice(0, 10).replace(/-/g, '')}_${backupForm.backupType}.zip`,
      size: '15.5 MB',
      date: new Date().toLocaleString('zh-CN'),
      type: backupForm.backupType === 'full' ? '完整' : '增量'
    }
    backupList.value.unshift(newBackup)
  }, 1500)
}

const handleRestore = (row) => {
  ElMessageBox.confirm(
    `确定要恢复备份【${row.name}】吗？此操作不可撤销！`,
    '恢复备份',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    ElMessage.success('备份恢复成功')
  }).catch(() => {})
}

const handleDeleteBackup = (row) => {
  ElMessageBox.confirm(
    `确定要删除备份【${row.name}】吗？`,
    '删除备份',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    const index = backupList.value.findIndex(b => b.id === row.id)
    if (index !== -1) {
      backupList.value.splice(index, 1)
      ElMessage.success('备份删除成功')
    }
  }).catch(() => {})
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped lang="scss">
.settings-container {
  .settings-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-size: 18px;
      font-weight: 600;

      .header-left {
        display: flex;
        align-items: center;
        gap: 10px;

        .header-icon {
          font-size: 20px;
          color: #1890ff;
        }
      }
    }

    :deep(.settings-tabs) {
      .el-tabs__header {
        margin: 0 0 24px 0;
        border-bottom: 1px solid #f0f0f0;

        .el-tabs__nav-wrap::after {
          display: none;
        }

        .el-tabs__item {
          padding: 12px 24px;
          font-size: 14px;
          color: #595959;

          .tab-label {
            display: flex;
            align-items: center;
            gap: 6px;

            .el-icon {
              font-size: 16px;
            }
          }

          &:hover {
            color: #1890ff;
          }

          &.is-active {
            color: #1890ff;
            font-weight: 500;
          }
        }

        .el-tabs__active-bar {
          background-color: #1890ff;
          height: 3px;
        }
      }

      .el-tabs__content {
        padding: 0;
      }
    }

    .settings-form {
      max-width: 700px;

      .form-section {
        margin-bottom: 32px;
        padding: 24px;
        background: #fafafa;
        border-radius: 8px;
        border: 1px solid #f0f0f0;

        .section-title {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-bottom: 20px;
          padding-bottom: 12px;
          font-size: 15px;
          font-weight: 600;
          color: #262626;
          border-bottom: 2px solid #1890ff;

          .el-icon {
            font-size: 18px;
            color: #1890ff;
          }
        }

        .switch-group {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }
      }

      .oss-section {
        animation: fadeIn 0.3s ease;
      }

      .storage-type-group {
        width: 100%;
        display: flex;
        flex-wrap: wrap;
        gap: 12px;

        .el-radio {
          margin-right: 0;
          padding: 12px 20px;

          .el-radio__label {
            display: flex;
            align-items: center;
            gap: 6px;
          }
        }
      }

      .form-tip {
        margin-left: 12px;
        font-size: 12px;
        color: #8c8c8c;
      }
    }

    .backup-section {
      .backup-card {
        margin-bottom: 20px;
        border: 1px solid #ebeef5;
        border-radius: 4px;

        :deep(.el-card__header) {
          background: #f5f7fa;
          border-bottom: 1px solid #ebeef5;
          padding: 14px 20px;

          .section-header {
            display: flex;
            align-items: center;
            gap: 8px;
            font-weight: 500;
            font-size: 14px;
            color: #303133;
          }
        }

        :deep(.el-card__body) {
          padding: 20px;
        }

        .file-icon {
          margin-right: 6px;
          color: #909399;
        }
      }
    }
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
