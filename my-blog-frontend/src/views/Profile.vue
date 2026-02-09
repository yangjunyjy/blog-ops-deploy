<template>
  <div class="profile-page">
    <div class="container">
      <div v-if="loading" class="loading-container">
        <el-skeleton :loading="loading" animated>
          <template #template>
            <div class="profile-skeleton">
              <el-skeleton-item variant="circle" style="width: 120px; height: 120px;" />
              <el-skeleton-item variant="h1" style="width: 200px; margin-top: 20px;" />
              <el-skeleton-item variant="text" style="width: 300px; margin-top: 10px;" />
            </div>
          </template>
        </el-skeleton>
      </div>

      <div v-else class="profile-content">
        <!-- 左侧个人信息卡片 -->
        <div class="profile-card">
          <div class="profile-header">
            <img :src="avatarUrl" :alt="userName" class="profile-avatar" />
            <div class="avatar-ring"></div>
          </div>

          <div class="profile-info">
            <h2 class="profile-name">{{ userName }}</h2>
            <p class="profile-bio">{{ userStore.user.bio || '这个人很懒，什么都没写~' }}</p>
            <div class="profile-email">
              <el-icon><Message /></el-icon>
              <span>{{ userEmail }}</span>
            </div>
          </div>

          <div class="profile-stats">
            <div class="stat-item">
              <el-icon><Star /></el-icon>
              <div class="stat-value">{{ stats.comments }}</div>
              <div class="stat-label">评论</div>
            </div>
            <div class="stat-item">
              <el-icon><Star /></el-icon>
              <div class="stat-value">{{ stats.likes }}</div>
              <div class="stat-label">获赞</div>
            </div>
            <div class="stat-item">
              <el-icon><Star /></el-icon>
              <div class="stat-value">{{ stats.favorites }}</div>
              <div class="stat-label">收藏</div>
            </div>
          </div>

          <div class="profile-actions" v-if="!editing">
            <el-button type="primary" :icon="Edit" @click="handleEdit">
              编辑资料
            </el-button>
            <el-button :icon="Setting" style="margin-left: 0px;">
              账号设置
            </el-button>
          </div>
        </div>

        <!-- 右侧内容区域 -->
        <div class="profile-main">
          <el-tabs v-model="activeTab" class="profile-tabs">
            <el-tab-pane label="我的收藏" name="favorites">
              <!-- 文件夹列表视图-->
              <div v-if="viewMode === 'folders'" class="favorites-view">
                <div class="view-header">
                  <h3 class="view-title">
                    <el-icon><Folder /></el-icon>
                    收藏文件夹
                  </h3>
                  <el-button
                    class="create-folder-btn"
                    type="primary"
                    :icon="FolderAdd"
                    @click="openCreateFolderDialog"
                    size="default"
                  >
                    新建文件夹
                  </el-button>
                </div>
                <div class="folders-grid">
                  <div
                    v-for="folder in folders"
                    :key="folder.id"
                    class="folder-card"
                    @click="switchFolder(folder.id)"
                  >
                    <div class="folder-card-header">
                      <el-icon class="folder-card-icon">
                        <Folder />
                      </el-icon>
                      <div class="folder-card-actions" @click.stop>
                        <el-button
                          type="primary"
                          link
                          size="small"
                          @click="openEditFolderDialog(folder)"
                        >
                          <el-icon><Edit /></el-icon>
                        </el-button>
                        <el-button
                          type="danger"
                          link
                          size="small"
                          @click="handleDeleteFolder(folder)"
                        >
                          <el-icon><Delete /></el-icon>
                        </el-button>
                      </div>
                    </div>
                    <div class="folder-card-content">
                      <h4 class="folder-card-title">{{ folder.name }}</h4>
                      <p class="folder-card-desc">{{ folder.description || '暂无描述' }}</p>
                      <div class="folder-card-stats">
                        <span class="stat">
                          <el-icon><Star /></el-icon>
                          {{ folder.articleCount || 0 }} 篇文章
                        </span>
                      </div>
                    </div>
                  </div>
                  <el-empty v-if="!folders.length" description="暂无收藏文件夹" />
                </div>
              </div>

              <!-- 文章列表视图 -->
              <div v-else class="articles-view">
                <div class="view-header">
                  <div class="header-left">
                    <el-button link @click="backToFolders">
                      <el-icon><Right style="transform: rotate(180deg)" /></el-icon>
                      返回文件夹
                    </el-button>
                    <h3 class="view-title">
                      <el-icon><Star /></el-icon>
                      {{ activeFolder?.name }}
                    </h3>
                  </div>
                  <el-button
                    type="primary"
                    link
                    @click="openEditFolderDialog(activeFolder)"
                  >
                    <el-icon><Edit /></el-icon>
                    编辑文件夹
                  </el-button>
                </div>

                <div v-loading="folderLoading" class="articles-list">
                  <div
                    v-for="article in folderArticles"
                    :key="article.id"
                    class="article-item"
                  >
                    <div class="article-main" @click="handleGoToArticle(article.id)">
                      <h4 class="article-title">{{ article.title }}</h4>
                      <div class="article-meta">
                        <span><el-icon><User /></el-icon> 阅读: {{ article.views }}</span>
                        <span><el-icon><Star /></el-icon> 点赞: {{ article.likes || 0 }}</span>
                        <span><el-icon><Message /></el-icon> {{ article.createdAt }}</span>
                      </div>
                    </div>
                    <div class="article-actions">
                      <el-dropdown trigger="click">
                        <el-button circle size="small">
                          <el-icon><Sort /></el-icon>
                        </el-button>
                        <template #dropdown>
                          <el-dropdown-menu>
                            <el-dropdown-item @click="openMoveArticleDialog(article)">
                              <el-icon><Folder /></el-icon>
                              移动到其他文件夹
                            </el-dropdown-item>
                            <el-dropdown-item @click="openDeleteArticleDialog(article)">
                              <el-icon><Delete /></el-icon>
                              移除收藏
                            </el-dropdown-item>
                          </el-dropdown-menu>
                        </template>
                      </el-dropdown>
                    </div>
                  </div>
                  <el-empty v-if="!folderArticles.length" description="该文件夹暂无收藏文章" />
                </div>

                <!-- 分页 -->
                <Pagination
                  v-if="articlePagination.total > 0"
                  :current-page="articlePagination.page"
                  :total="articlePagination.total"
                  :page-size="articlePagination.pageSize"
                  :total-pages="totalPages"
                  @change="handlePageChange"
                />
              </div>
            </el-tab-pane>

            <el-tab-pane label="互动记录" name="interactions">
              <ActivityRecord :user-id="userStore.user?.id" />
            </el-tab-pane>

            <el-tab-pane label="个人资料" name="profile">
              <div class="edit-profile">
                <div class="form-section">
                  <h3 class="section-title">
                    <el-icon><User /></el-icon>
                    基本信息
                  </h3>

                  <div class="form-group">
                    <label>用户名</label>
                    <el-input
                      v-model="editForm.name"
                      :disabled="!editing"
                      size="large"
                      clearable
                    />
                  </div>

                  <div class="form-group">
                    <label>邮箱地址</label>
                    <el-input
                      v-model="editForm.email"
                      type="email"
                      :disabled="!editing"
                      size="large"
                      clearable
                    />
                  </div>

                  <div class="form-group">
                    <label>个人简介</label>
                    <el-input
                      v-model="editForm.bio"
                      type="textarea"
                      :rows="4"
                      :disabled="!editing"
                      maxlength="200"
                      show-word-limit
                    />
                  </div>

                  <div class="form-actions" v-if="editing">
                    <el-button @click="handleCancel">取消</el-button>
                    <el-button type="primary" @click="handleSave">
                      <el-icon><Star /></el-icon>
                      保存修改
                    </el-button>
                  </div>
                </div>
              </div>
            </el-tab-pane>

            <el-tab-pane label="账号安全" name="security">
              <div class="security-section">
                <div class="security-item">
                  <div class="security-icon">
                    <el-icon><Lock /></el-icon>
                  </div>
                  <div class="security-info">
                    <h4>修改密码</h4>
                    <p>定期修改密码可以提高账号安全性</p>
                  </div>
                  <el-button type="primary" plain>修改</el-button>
                </div>

                <div class="security-item">
                  <div class="security-icon">
                    <el-icon><Message /></el-icon>
                  </div>
                  <div class="security-info">
                    <h4>绑定邮箱</h4>
                    <p>{{ userEmail }}</p>
                  </div>
                  <el-button type="primary" plain>更改</el-button>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </div>

    <!-- 文件夹对话框 -->
    <el-dialog
      v-model="folderDialogVisible"
      :title="folderDialogMode === 'create' ? '新建收藏文件夹' : '编辑收藏文件夹'"
      width="500px"
    >
      <el-form :model="folderForm" label-width="80px">
        <el-form-item label="名称" required>
          <el-input
            v-model="folderForm.name"
            placeholder="请输入文件夹名称"
            maxlength="20"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="序号">
          <el-input-number
            v-model="folderForm.sortOrder"
            :min="0"
            :max="999"
            controls-position="right"
            placeholder="请输入文件夹序号"
          />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="folderForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入描述（可选）"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="folderDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleFolderSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 文章操作对话框-->
    <el-dialog
      v-model="articleDialogVisible"
      :title="articleDialogMode === 'delete' ? '移除收藏' : '移动到其他文件夹'"
      width="500px"
    >
      <div v-if="articleDialogMode === 'delete'" class="article-action-content">
        <p>确定从当前文件夹移除这篇收藏文章吗？</p>
        <div class="article-preview">
          <strong>{{ currentArticle?.title }}</strong>
        </div>
      </div>
      <div v-else class="article-action-content">
        <p>选择目标文件夹</p>
        <el-select
          v-model="moveTargetFolderId"
          placeholder="请选择文件夹"
          style="width: 100%; margin-top: 12px;"
        >
          <el-option
            v-for="folder in folders.filter(f => f.id !== activeFolderId)"
            :key="folder.id"
            :label="`${folder.name} (${folder.articleCount || 0} 篇)`"
            :value="folder.id"
          />
        </el-select>
      </div>
      <template #footer>
        <el-button @click="articleDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          @click="articleDialogMode === 'delete' ? handleDeleteArticle() : handleMoveArticle()"
        >
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { User, Message, Edit, Lock, Star, Setting, Right, Folder, FolderAdd, Delete, Sort } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import Pagination from '@/components/Pagination.vue'
import ActivityRecord from '@/components/ActivityRecord.vue'
import {
  getFavoriteFolders,
  createFavoriteFolder,
  updateFavoriteFolder,
  deleteFavoriteFolder,
  getFolderArticles,
  removeArticleFromFolder,
  moveArticleToFolder
} from '@/api'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(true)
const editing = ref(false)
const editForm = ref({
  name: '',
  email: '',
  bio: ''
})

const activeTab = ref('favorites')
const activeFolderId = ref(null)

const stats = ref({
  comments: 48,
  likes: 156,
  favorites: 23
})

// 收藏文件夹
const folders = ref([])
const folderArticles = ref([])
const folderLoading = ref(false)
const viewMode = ref('folders')

// 文章分页
const articlePagination = ref({
  page: 1,
  pageSize: 5,
  total: 0
})

const folderPagination = ref({
  page: 1,
  pageSize: 5,
  total: 0
})

// 对话框
const folderDialogVisible = ref(false)
const folderDialogMode = ref('create') // 'create' or 'edit'
const currentFolder = ref(null)
const folderForm = ref({
  name: '',
  description: '',
  sortOrder:0
})

const articleDialogVisible = ref(false)
const articleDialogMode = ref('delete') // 'delete' or 'move'
const currentArticle = ref(null)
const moveTargetFolderId = ref(null)

const avatarUrl = computed(() => userStore.user?.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=default')
const userName = computed(() => userStore.user?.name || userStore.user?.username || '用户')
const userEmail = computed(() => userStore.user?.email || '')

// 获取当前激活的文件夹
const activeFolder = computed(() => {
  if (!activeFolderId.value) return null
  return folders.value.find(f => f.id === activeFolderId.value)
})

// 计算总页
const totalPages = computed(() => {
  return Math.ceil(articlePagination.value.total / articlePagination.value.pageSize) || 1
})

const checkLoginStatus = () => {
  userStore.loadUserFromStorage()

  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return false
  }
  return true
}

const loadUserData = () => {
  if (!checkLoginStatus()) return

  editForm.value = {
    name: userName.value,
    email: userEmail.value,
    bio: userStore.user?.bio || ''
  }

  loading.value = false
  loadFavoriteFolders()
}

// 加载收藏文件夹
const loadFavoriteFolders = async () => {
  try {
    const res = await getFavoriteFolders({
      page:folderPagination.value.page,
      pageSize:folderPagination.value.pageSize
    })
    if (res.code === 200 && res.data) {
      // 后端返回的是分页数据结构：{ list: [], total: 0, page: 1, pageSize: 5 }
      folders.value = res.data.list || []
      folderPagination.value.total = res.data.total || 0
      if (folders.value.length > 0 && !activeFolderId.value) {
        activeFolderId.value = folders.value[0].id
      }
      loadFolderArticles(activeFolderId.value)
    }
  } catch (error) {
    console.error('加载收藏文件夹失败', error)
    // 使用 Mock 数据作为后备
    folders.value = [
      { id: 1, name: '默认收藏夹', isDefault: true, articleCount: 2, description: '系统默认收藏夹' }
    ]
    if (!activeFolderId.value && folders.value.length > 0) {
      activeFolderId.value = folders.value[0].id
    }
    loadFolderArticlesMock(activeFolderId.value)
  }
}

// 加载文件夹下的文章
const loadFolderArticles = async (folderId, page = 1, pageSize = 10) => {
  if (!folderId) return

  folderLoading.value = true
  try {
    const res = await getFolderArticles(folderId, { page, pageSize })
    if (res.code === 200) {
      folderArticles.value = res.data.list || []
      articlePagination.value.total = res.data.total || 0
    }
  } catch (error) {
    console.error('加载收藏文章失败', error)
    // 使用 Mock 数据作为后备
    loadFolderArticlesMock(folderId)
  } finally {
    folderLoading.value = false
  }
}

// Mock 数据加载
const loadFolderArticlesMock = () => {
  const page = articlePagination.value.page
  const pageSize = articlePagination.value.pageSize
  const allArticles = [
    { id: 2, title: 'JavaScript 异步编程深入理解', views: 987, likes: 28, createdAt: '2024-01-12' },
    { id: 3, title: 'React Hooks 完全指南', views: 856, likes: 22, createdAt: '2024-01-08' },
    { id: 4, title: 'TypeScript 高级类型实战', views: 765, likes: 25, createdAt: '2024-01-10' },
    { id: 5, title: 'CSS Grid 布局详解', views: 543, likes: 19, createdAt: '2024-01-05' },
    { id: 6, title: 'Vue 3 Composition API 最佳实践', views: 1234, likes: 35, createdAt: '2024-01-15' },
    { id: 7, title: 'Webpack 性能优化指南', views: 678, likes: 21, createdAt: '2024-01-09' },
    { id: 8, title: 'Node.js 事件循环机制', views: 892, likes: 27, createdAt: '2024-01-11' },
    { id: 9, title: '前端性能监控方案', views: 456, likes: 16, createdAt: '2024-01-03' },
    { id: 10, title: 'React 源码解析', views: 1567, likes: 42, createdAt: '2024-01-18' },
    { id: 11, title: 'JavaScript 设计模式', views: 1098, likes: 31, createdAt: '2024-01-14' }
  ]

  const startIndex = (page - 1) * pageSize
  const endIndex = startIndex + pageSize
  folderArticles.value = allArticles.slice(startIndex, endIndex)
  articlePagination.value.total = 10
}

// 切换文件夹
const switchFolder = (folderId) => {
  activeFolderId.value = folderId
  viewMode.value = 'articles'
  articlePagination.value.page = 1
  loadFolderArticles(folderId, 1, articlePagination.value.pageSize)
}

// 返回文件夹列表
const backToFolders = () => {
  activeFolderId.value = null
  viewMode.value = 'folders'
  folderArticles.value = []
}

// 分页改变
const handlePageChange = (page) => {
  articlePagination.value.page = page
  loadFolderArticles(activeFolderId.value, page, articlePagination.value.pageSize)
}

// 打开创建文件夹对话框
const openCreateFolderDialog = () => {
  folderDialogMode.value = 'create'
  folderForm.value = {
    name: '',
    description: ''
  }
  folderDialogVisible.value = true
}

// 打开编辑文件夹对话框
const openEditFolderDialog = (folder) => {
  folderDialogMode.value = 'edit'
  currentFolder.value = folder
  folderForm.value = {
    name: folder.name,
    description: folder.description || '',
    sortOrder:folder.sortOrder
  }
  folderDialogVisible.value = true
}

const handleFolderSubmit = async () => {
  if (!folderForm.value.name.trim()) {
    ElMessage.warning('请输入文件夹名称')
    return
  }

  try {
    if (folderDialogMode.value === 'create') {
      const res = await createFavoriteFolder(folderForm.value)
      if (res.code === 200) {
        ElMessage.success('创建成功')
        folderDialogVisible.value = false
        loadFavoriteFolders()
      }
    } else {
      const res = await updateFavoriteFolder(currentFolder.value.id, folderForm.value)
      if (res.code === 200) {
        ElMessage.success('更新成功')
        folderDialogVisible.value = false
        loadFavoriteFolders()
      }
    }
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleDeleteFolder = async (folder) => {

  if (folder.articleCount > 0) {
    try {
      await ElMessageBox.confirm(
        `该文件夹包含 ${folder.articleCount} 篇文章，删除后这些文章将被移至默认文件夹，确定删除吗？`,
        '提示',
        { type: 'warning' }
      )
    } catch {
      return
    }
  } else {
    try {
      await ElMessageBox.confirm('确定删除该文件夹吗？', '提示', { type: 'warning' })
    } catch {
      return
    }
  }

  try {
    const res = await deleteFavoriteFolder(folder.id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      if (activeFolderId.value === folder.id) {
        activeFolderId.value = null
      }
      loadFavoriteFolders()
    }
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const openDeleteArticleDialog = (article) => {
  articleDialogMode.value = 'delete'
  currentArticle.value = article
  articleDialogVisible.value = true
}

const openMoveArticleDialog = (article) => {
  articleDialogMode.value = 'move'
  currentArticle.value = article
  moveTargetFolderId.value = null
  articleDialogVisible.value = true
}

// 处理删除文章
const handleDeleteArticle = async () => {
  try {
    const res = await removeArticleFromFolder(activeFolderId.value, currentArticle.value.id)
    if (res.code === 200) {
      articleDialogVisible.value = false
      loadFavoriteFolders()
      loadFolderArticles(activeFolderId.value)
    }
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 处理移动文章
const handleMoveArticle = async () => {
  if (!moveTargetFolderId.value) {
    return
  }

  if (moveTargetFolderId.value === activeFolderId.value) {
    ElMessage.warning('不能移动到当前文件夹')
    return
  }

  try {
    const res = await moveArticleToFolder(
      currentArticle.value.id,
      activeFolderId.value,
      moveTargetFolderId.value
    )
    if (res.code === 200) {
      ElMessage.success('移动成功')
      articleDialogVisible.value = false
      loadFavoriteFolders()
      loadFolderArticles(activeFolderId.value)
    }
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

onMounted(() => {
  setTimeout(() => {
    loadUserData()
    loadInteractions()
  }, 100)
})

watch(() => userStore.isLoggedIn, (isLoggedIn) => {
  if (isLoggedIn && loading.value) {
    loadUserData()
  }
}, { immediate: true })

const handleEdit = () => {
  editing.value = true
}

const handleCancel = () => {
  editing.value = false
  editForm.value = {
    name: userName.value,
    email: userEmail.value,
    bio: userStore.user?.bio || ''
  }
}

const handleSave = async () => {
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))

    userStore.updateUser({
      name: editForm.value.name,
      email: editForm.value.email,
      bio: editForm.value.bio
    })

    ElMessage.success('保存成功')
    editing.value = false
  } catch (error) {
    ElMessage.error('保存失败，请重试')
  }
}

const handleGoToArticle = (id) => {
  router.push(`/article/${id}`)
}
</script>

<style scoped>
.profile-page {
  min-height: calc(100vh - 140px);
  padding: 50px 0;
}

.profile-content {
  display: grid;
  grid-template-columns: 340px 1fr;
  gap: 40px;
  align-items: start;
}

.profile-card {
  background: #f8fafc;
  border-radius: 20px;
  padding: 0;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.5);
  overflow: hidden;
}

html.dark .profile-card {
  background: #2d3748;
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.profile-header {
  position: relative;
  display: flex;
  justify-content: center;
  padding: 40px 0;
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
}

.profile-avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: 4px solid #fff;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
  position: relative;
  z-index: 2;
}

.avatar-ring {
  position: absolute;
  top: 20px;
  left: 20px;
  right: 20px;
  bottom: 20px;
  border-radius: 50%;
  border: 3px solid rgba(255, 255, 255, 0.3);
}

.profile-info {
  padding: 30px;
  text-align: center;
}

.profile-name {
  font-size: 24px;
  font-weight: 700;
  color: #1a202c;
  margin: 0 0 12px;
}

html.dark .profile-name {
  color: #f7fafc;
}

.profile-bio {
  font-size: 14px;
  color: #718096;
  margin: 0 0 16px;
  line-height: 1.6;
}

html.dark .profile-bio {
  color: #a0aec0;
}

.profile-email {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 14px;
  color: #a0aec0;
}

html.dark .profile-email {
  color: #718096;
}

.profile-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1px;
  background: #f8fafc;
}

html.dark .profile-stats {
  background: #1a202c;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px 16px;
  gap: 8px;
  transition: all 0.3s;
}

.stat-item:hover {
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
}

html.dark .stat-item:hover {
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
}

.stat-item .el-icon {
  font-size: 24px;
  color: #1890ff;
}

.stat-value {
  font-size: 28px;
  font-weight: 800;
  color: #1890ff;
  line-height: 1;
}

.stat-label {
  font-size: 13px;
  color: #718096;
}

.profile-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 30px;
}

.profile-actions :deep(.el-button) {
  width: 100%;
  justify-content: center;
  gap: 8px;
}

.profile-main {
  background: #f8fafc;
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.5);
}

html.dark .profile-main {
  background: #2d3748;
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.profile-tabs :deep(.el-tabs__header) {
  margin: 0 0 30px;
}

.profile-tabs :deep(.el-tabs__item) {
  font-size: 16px;
  font-weight: 500;
}

.articles-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.article-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  background: #f8fafc;
  border-radius: 12px;
  transition: all 0.3s;
  cursor: pointer;
}

html.dark .article-item {
  background: #1a202c;
}

.article-item:hover {
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
  transform: translateX(8px);
}

html.dark .article-item:hover {
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
}

.article-info {
  flex: 1;
}

.article-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a202c;
  margin: 0 0 8px;
}

html.dark .article-title {
  color: #f7fafc;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 13px;
  color: #a0aec0;
}

html.dark .article-meta {
  color: #718096;
}

.article-meta .el-icon {
  font-size: 14px;
  margin-right: 4px;
}

.article-arrow {
  font-size: 18px;
  color: #1890ff;
}

.article-main {
  flex: 1;
  cursor: pointer;
}

.article-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.favorites-view,
.articles-view {
  min-height: 400px;
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e2e8f0;
}

html.dark .view-header {
  border-bottom-color: #4a5568;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.view-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #1a202c;
}

html.dark .view-title {
  color: #f7fafc;
}

.view-title .el-icon {
  color: #1890ff;
}

.create-folder-btn {
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
  border: none;
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 500;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.25);
  transition: all 0.2s ease;
}

.create-folder-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.35);
  background: linear-gradient(135deg, #40a9ff 0%, #1890ff 100%);
}

.create-folder-btn:active {
  transform: translateY(0);
}

/* 文件夹网格*/
.folders-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.folder-card {
  /* background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%); */
  background-color: #e3e03e;
  border-radius: 10px;
  padding: 16px;
  border: 1px solid #e2e8f0;
  transition: all 0.2s ease;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

html.dark .folder-card {
  background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
  border-color: #334155;
}

.folder-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 3px;
  height: 100%;
  background: linear-gradient(180deg, #c4c9cd 0%, #096dd9 100%);
  opacity: 0;
  transition: opacity 0.2s;
}

.folder-card:hover::before {
  opacity: 1;
}

.folder-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.15);
  border-color: #54585c;
}

html.dark .folder-card:hover {
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.2);
}

.folder-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.folder-card-icon {
  font-size: 28px;
  color: #1890ff;
  filter: drop-shadow(0 2px 4px rgba(24, 144, 255, 0.2));
}

.folder-card-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.25s;
}

.folder-card:hover .folder-card-actions {
  opacity: 1;
}

.folder-card-actions :deep(.el-button) {
  padding: 4px;
  width: 28px;
  height: 28px;
}

.folder-card-actions :deep(.el-button--primary.is-link) {
  color: #1890ff;
}

.folder-card-actions :deep(.el-button--danger.is-link) {
  color: #f56565;
}

.folder-card-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.folder-card-title {
  font-size: 15px;
  font-weight: 600;
  color: #1a202c;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

html.dark .folder-card-title {
  color: #f1f5f9;
}

.folder-card-desc {
  font-size: 12px;
  color: #94a3b8;
  margin: 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  min-height: 34px;
}

html.dark .folder-card-desc {
  color: #64748b;
}

.folder-card-stats {
  display: flex;
  gap: 12px;
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid rgba(148, 163, 184, 0.2);
}

html.dark .folder-card-stats {
  border-top-color: rgba(100, 116, 139, 0.2);
}

.folder-card-stats .stat {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #64748b;
}

.folder-card-stats .stat .el-icon {
  color: #1890ff;
  font-size: 14px;
}

/* 文章视图 */
.articles-view {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.articles-view .pagination {
  margin-top: auto;
  padding-top: 20px;
  border-top: 1px solid #e2e8f0;
}

html.dark .articles-view .pagination {
  border-top-color: #4a5568;
}

.article-action-content {
  padding: 12px 0;
}

.article-action-content p {
  margin: 0 0 16px;
  color: #4a5568;
}

html.dark .article-action-content p {
  color: #cbd5e0;
}

.article-preview {
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  border-left: 3px solid #1890ff;
}

html.dark .article-preview {
  background: #1a202c;
}

.article-preview strong {
  color: #1a202c;
  font-size: 14px;
}

html.dark .article-preview strong {
  color: #e2e8f0;
}

.edit-profile {
  max-width: 600px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 600;
  color: #1a202c;
  margin: 0 0 24px;
  padding-bottom: 16px;
  border-bottom: 2px solid #e2e8f0;
}

html.dark .section-title {
  color: #f7fafc;
  border-bottom-color: #4a5568;
}

.section-title .el-icon {
  font-size: 22px;
  color: #1890ff;
}

.form-group {
  margin-bottom: 24px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #4a5568;
  margin-bottom: 8px;
}

html.dark .form-group label {
  color: #cbd5e0;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 30px;
}

/* 互动记录样式 */
.interactions-view {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.interaction-filters {
  display: flex;
  gap: 8px;
  padding: 8px;
  background: #f8fafc;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

html.dark .interaction-filters {
  background: #1a202c;
  border-color: #4a5568;
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  border: none;
  background: transparent;
  color: #718096;
  font-size: 14px;
  font-weight: 500;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

html.dark .filter-btn {
  color: #a0aec0;
}

.filter-btn:hover {
  background: rgba(24, 144, 255, 0.1);
  color: #1890ff;
}

.filter-btn.active {
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
  color: #fff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.25);
}

.interactions-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.interaction-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 20px;
  background: #f8fafc;
  border-radius: 12px;
  transition: all 0.2s ease;
  cursor: pointer;
}

html.dark .interaction-item {
  background: #1a202c;
}

.interaction-item:hover {
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
  transform: translateX(4px);
}

html.dark .interaction-item:hover {
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
}

.interaction-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  font-size: 20px;
  flex-shrink: 0;
}

.interaction-like {
  background: linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%);
  color: #ed8936;
}

.interaction-comment {
  background: linear-gradient(135deg, #d4fc79 0%, #96e6a1 100%);
  color: #48bb78;
}

.interaction-share {
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  color: #ed64a6;
}

.interaction-content {
  flex: 1;
  min-width: 0;
}

.interaction-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  gap: 12px;
}

.interaction-type {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #718096;
  font-weight: 500;
}

html.dark .interaction-type {
  color: #a0aec0;
}

.interaction-time {
  font-size: 12px;
  color: #a0aec0;
}

html.dark .interaction-time {
  color: #718096;
}

.interaction-article-title {
  font-size: 15px;
  font-weight: 600;
  color: #1a202c;
  margin: 0 0 8px;
  line-height: 1.4;
}

html.dark .interaction-article-title {
  color: #f7fafc;
}

.interaction-comment {
  font-size: 13px;
  color: #718096;
  margin: 0 0 4px;
  padding: 8px 12px;
  background: rgba(24, 144, 255, 0.08);
  border-radius: 8px;
  border-left: 3px solid #1890ff;
}

html.dark .interaction-comment {
  color: #a0aec0;
  background: rgba(24, 144, 255, 0.15);
}

.interaction-platform {
  font-size: 12px;
  color: #ed64a6;
  margin: 0;
  font-weight: 500;
}

.interaction-arrow {
  color: #1890ff;
  flex-shrink: 0;
  margin-top: 4px;
}

.interactions-view .pagination {
  margin-top: auto;
  padding-top: 20px;
  border-top: 1px solid #e2e8f0;
}

html.dark .interactions-view .pagination {
  border-top-color: #4a5568;
}

.security-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.security-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px;
  background: #f8fafc;
  border-radius: 12px;
  gap: 20px;
}

html.dark .security-item {
  background: #1a202c;
}

.security-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
  border-radius: 12px;
  flex-shrink: 0;
}

.security-icon .el-icon {
  font-size: 24px;
  color: #1890ff;
}

.security-info {
  flex: 1;
}

.security-info h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1a202c;
  margin: 0 0 8px;
}

html.dark .security-info h4 {
  color: #f7fafc;
}

.security-info p {
  font-size: 14px;
  color: #718096;
  margin: 0;
}

html.dark .security-info p {
  color: #a0aec0;
}

@media (max-width: 992px) {
  .profile-content {
    grid-template-columns: 1fr;
  }

  .profile-card {
    order: 2;
  }
}

@media (max-width: 768px) {
  .profile-page {
    padding: 30px 0;
  }

  .profile-content {
    gap: 30px;
  }

  .profile-stats {
    grid-template-columns: repeat(2, 1fr);
  }

  .security-item {
    flex-direction: column;
    text-align: center;
    gap: 16px;
  }

  .folders-grid {
    grid-template-columns: 1fr;
  }

  .view-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .header-left {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .folders-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }

  .interaction-filters {
    flex-wrap: wrap;
  }

  .filter-btn {
    padding: 8px 16px;
    font-size: 13px;
  }

  .interaction-item {
    padding: 16px;
  }

  .interaction-icon {
    width: 40px;
    height: 40px;
    font-size: 18px;
  }

  .interaction-article-title {
    font-size: 14px;
  }

  .folder-card {
    padding: 14px;
  }

  .folder-card-icon {
    font-size: 24px;
  }

  .folder-card-title {
    font-size: 14px;
  }

  .folder-card-desc {
    font-size: 11px;
    min-height: 30px;
  }

  .create-folder-btn {
    padding: 8px 16px;
    font-size: 13px;
  }
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.profile-skeleton {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 40px;
}

@media (max-width: 992px) {
  .profile-content {
    grid-template-columns: 1fr;
  }

  .profile-card {
    order: 2;
  }
}

@media (max-width: 480px) {
  .folders-grid {
    grid-template-columns: 1fr;
  }
}
</style>
