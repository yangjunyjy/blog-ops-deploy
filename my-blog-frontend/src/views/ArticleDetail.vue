<script setup>
import { ref, onMounted, watch, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { addComment, getArticleDetail, getComments, likeComment, unlikeComment, deleteComment } from '@/api'
import { formatDate, markdownToHtml } from '@/utils/format'
import { Calendar, StarFilled,Star, ChatDotRound, ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
import ArticleActions from '@/components/ArticleActions.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import HaloButton from '@/components/HaloButton.vue'
import ReplayComment from '@/components/ReplayComment.vue'
import CommentItem from '@/components/CommentItem.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const article = ref(null)
const loading = ref(true)
const reployDialogVisible = ref(false)

const targetComment = ref({
  id:'',
  content:''
})
// 评论相关
const comments = ref([])
const commentsLoading = ref(false)
const newComment = ref('')
const commentItemRefs = ref({}) // 保存评论组件引用

// 计算属性：检查当前用户是否已登录
const currentUser = computed(() => userStore.user)
const currentUserId = computed(() => userStore.user?.id)

// 为 wangEditor 生成的代码块添加样式和复制按钮
const enhanceCodeBlocks = () => {
  const preElements = document.querySelectorAll('.article-content pre')
  preElements.forEach((pre) => {
    // 如果已经处理过，跳过
    if (pre.classList.contains('enhanced')) return

    const codeElement = pre.querySelector('code')
    if (!codeElement) return

    // 获取语言类型（从 class 中提取）
    const languageClass = Array.from(codeElement.classList).find(cls => cls.startsWith('language-'))
    const lang = languageClass ? languageClass.replace('language-', '') : 'plaintext'

    // 获取代码内容
    const codeText = codeElement.textContent

    // 生成行号
    const lines = codeText.split('\n')
    const lineNumbers = lines.map((_, idx) => `<span class="code-line-number">${idx + 1}</span>`).join('')

    // 创建包装器
    const wrapper = document.createElement('div')
    wrapper.className = 'code-block-wrapper'
    wrapper.innerHTML = `
      <div class="code-header">
        <span class="code-language">${lang}</span>
        <button class="code-copy-btn" onclick="copyWangEditorCode(this)">
          <svg viewBox="0 0 24 24" width="14" height="14">
            <path fill="currentColor" d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z"/>
          </svg>
          复制
        </button>
      </div>
      <div class="code-body">
        <div class="code-line-numbers">${lineNumbers}</div>
        <div class="code-content">
          <code class="hljs language-${lang}">${codeElement.innerHTML}</code>
        </div>
      </div>
    `

    // 替换原始 pre 标签
    pre.parentNode.replaceChild(wrapper, pre)
  })
}

// 复制 wangEditor 代码块的函数
window.copyWangEditorCode = async (button) => {
  const codeBlock = button.closest('.code-block-wrapper')
  const codeElement = codeBlock.querySelector('.code-content code')

  try {
    await navigator.clipboard.writeText(codeElement.textContent)

    const originalContent = button.innerHTML
    button.innerHTML = `
      <svg viewBox="0 0 24 24" width="14" height="14">
        <path fill="currentColor" d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
      </svg>
      已复制
    `
    button.classList.add('copied')

    setTimeout(() => {
      button.innerHTML = originalContent
      button.classList.remove('copied')
    }, 2000)
  } catch (error) {
    console.error('复制失败:', error)
  }
}

const loadArticleDetail = async () => {
  loading.value = true
  article.value = null
  comments.value = []
  try {
    const res = await getArticleDetail(route.params.id)
    console.log('后端返回的文章数据:', res.data);

    article.value = res.data

    // 等待 DOM 更新后增强代码块
    await nextTick()
    enhanceCodeBlocks()

    // 加载评论
    loadComments()
  } catch (error) {
    console.error('加载文章详情失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载评论
const loadComments = async () => {
  commentsLoading.value = true
  try {
    const res = await getComments(route.params.id)
    console.log('后端返回的评论数据:', res.data)

    // 转换后端数据为前端需要的格式
    if (res.data && res.data.items) {
      comments.value = res.data.items.map(item => ({
        id: item.id,
        userId: item.user_id,
        content: item.content,
        author: item.author ? {
          id: item.author.id,
          name: item.author.name,
          avatar: item.author.avatar || '/images/default.jpg'
        } : {
          id: item.user_id,
          name: '用户',
          avatar: '/images/default.jpg'
        },
        createdAt: item.created_at,
        likes: item.likes || 0,
        isLiked: item.is_liked || false
      }))
    } else {
      comments.value = []
    }

    console.log('转换后的评论数据:', comments.value)
  } catch (error) {
    console.error('加载评论失败:', error)
    comments.value = []
  } finally {
    commentsLoading.value = false
  }
}

// 点赞评论
const handleLikeComment = async (commentId) => {
  if (!currentUser.value) {
    ElMessage.warning('请先登录')
    return
  }

  try {
    const comment = comments.value.find(c => c.id === commentId)
    if (!comment) return

    if (comment.isLiked) {
      // 取消点赞
      await unlikeComment(commentId)
      comment.likes--
      comment.isLiked = false
      ElMessage.success('取消点赞成功')
    } else {
      // 点赞
      await likeComment(commentId)
      comment.likes++
      comment.isLiked = true
      ElMessage.success('点赞成功')
    }
  } catch (error) {
    console.error('点赞操作失败:', error)
    ElMessage.error(error.response?.data?.message || '操作失败，请重试')
  }
}

// 提交评论
const submitComment = async () => {
  if (!newComment.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }

  try {
    const data = await addComment({
      articleId: Number(route.params.id),
      content: newComment.value,
      parentId: null
    })

    console.log('评论提交成功:', data)
    console.log('当前用户信息:', currentUser.value)
    console.log('后端返回的user:', data.data?.user)
    console.log('后端返回的comment:', data.data?.comment)

    // 添加新评论到列表 - 使用后端返回的过滤后的内容
    const newCommentData = {
      id: data.data?.comment?.id || data.data?.items?.id,
      userId: data.data?.comment?.user_id || data.data?.comment?.userId || currentUser.value.id,
      content: data.data?.comment?.content || newComment.value, // 使用后端返回的过滤后的内容
      author: {
        // 优先使用后端返回的用户信息
        id: data.data?.user?.id || currentUser.value.id,
        name: data.data?.user?.nickname || data.data?.user?.username || currentUser.value.nickname || currentUser.value.username,
        avatar: data.data?.user?.avatar || currentUser.value.avatar || '/images/default.jpg'
      },
      createdAt: data.data?.comment?.created_at || new Date().toISOString(),
      likes: data.data?.comment?.likes || 0,
      status: data.data?.status || 1,
      isLiked: false
    }

    comments.value.unshift(newCommentData)
    newComment.value = ''

    // 根据评论状态显示不同提示
    if (data.data?.status === 0) {
      ElMessage.info('评论已提交，等待审核')
    } else {
      ElMessage.success('评论发表成功！')
    }
  } catch (error) {
    console.error('提交评论失败:', error)
    ElMessage.error(error.response?.data?.message || '提交评论失败，请重试')
  }
}

// 复制代码功能（全局函数，在 HTML 中调用）
window.copyCode = async (button) => {
  const codeBlock = button.closest('.code-block-wrapper')
  const codeElement = codeBlock.querySelector('.code-content code')

  try {
    await navigator.clipboard.writeText(codeElement.textContent)

    // 更新按钮状态
    const originalContent = button.innerHTML
    button.innerHTML = `
      <svg viewBox="0 0 24 24" width="14" height="14">
        <path fill="currentColor" d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
      </svg>
      已复制
    `
    button.classList.add('copied')

    setTimeout(() => {
      button.innerHTML = originalContent
      button.classList.remove('copied')
    }, 2000)
  } catch (error) {
    console.error('复制失败:', error)
  }
}

const replyComment = (comment) =>{
  console.log("回复的评论id为",comment.id);
  reployDialogVisible.value = !reployDialogVisible.value
  targetComment.value.id = comment.id
  targetComment.value.content = comment.content
}

// 设置评论组件引用
const setCommentRef = (el, commentId) => {
  if (el) {
    commentItemRefs.value[commentId] = el
  }
}

// 处理回复提交
const handleReplySubmit = async (replyData) => {
  console.log('提交的回复数据：', replyData)

  // 获取父评论ID
  const parentCommentId = Number(targetComment.value.id)

  // 尝试找到对应的评论组件并添加回复
  const commentRef = commentItemRefs.value[parentCommentId]
  if (commentRef && commentRef.addNewReply) {
    // 直接传递完整的后端返回数据
    commentRef.addNewReply(replyData)
  } else {
    // 如果找不到对应的组件，重新加载整个评论列表
    await loadComments()
  }
}

// 处理评论删除
const handleCommentDelete = async (comment) => {
  // 从列表中删除评论（删除操作已在子组件中完成）
  const index = comments.value.findIndex(c => c.id === comment.id)
  if (index > -1) {
    comments.value.splice(index, 1)
  }
}

onMounted(() => {
  loadArticleDetail()
})

watch(() => route.params.id, async () => {
  await loadArticleDetail()
  window.scrollTo({ top: 0, behavior: 'smooth' })
})
</script>

<template>
  <!-- 模板部分保持不变 -->
  <div class="article-detail-page">
    <div class="container">
      <el-skeleton :loading="loading" animated>
        <template #template>
          <div style="max-width: 800px; margin: 40px auto;">
            <el-skeleton-item variant="h1" style="width: 80%; margin-bottom: 20px;" />
            <el-skeleton-item variant="text" style="margin-bottom: 20px;" />
            <el-skeleton-item variant="p" style="width: 100%; height: 400px;" />
          </div>
        </template>
        <template #default>
          <div v-if="article" class="article-detail">
            <div class="article-header">
              <h1 class="article-title">{{ article.title }}</h1>

              <div class="article-meta">
                <span class="meta-item">
                  <el-icon><Calendar /></el-icon>
                  {{ formatDate(article.createdAt) }}
                </span>
                <span class="meta-item">
                  <el-icon><View /></el-icon>
                  {{ article.views || 0 }} 阅读
                </span>
                <span class="meta-item" v-if="article.comments">
                  <el-icon><ChatDotRound /></el-icon>
                  {{ article.comments }} 评论
                </span>
                <span class="category" v-if="article.category">
                  分类: {{ article.category.name }}
                </span>
                <span class="meta-item" v-if="article.author">
                  作者: {{ article.author.name }}
                </span>
              </div>

              <div class="article-tags" v-if="article.tags && article.tags.length">
                <el-tag
                  v-for="tag in article.tags"
                  :key="tag.name"
                  type="info"
                  effect="plain"
                >
                  {{ tag.name }}
                </el-tag>
              </div>
            </div>

            <div class="article-cover" v-if="article.cover">
              <img :src="article.cover" :alt="article.title" />
            </div>

            <div class="article-body">
              <div
                class="article-content"
                v-html="markdownToHtml(article.content)"
              ></div>
            </div>

            <div class="article-footer">
              <div class="article-nav">
                <router-link
                  v-if="article.prevArticle"
                  :to="`/article/${article.prevArticle.id}`"
                  class="prev-article"
                >
                  <el-icon><ArrowLeft /></el-icon>
                  {{ article.prevArticle.title }}
                </router-link>
                <router-link
                  v-if="article.nextArticle"
                  :to="`/article/${article.nextArticle.id}`"
                  class="next-article"
                >
                  {{ article.nextArticle.title }}
                  <el-icon><ArrowRight /></el-icon>
                </router-link>
              </div>
            </div>

            <!-- 评论区 -->
            <div class="comments-section">
              <div class="comments-header">
                <h3>
                  <el-icon><ChatDotRound /></el-icon>
                  评论 ({{ comments.length }})
                </h3>
              </div>

              <!-- 发表评论 -->
              <div class="comment-form">
                <el-input
                  v-model="newComment"
                  type="textarea"
                  :rows="4"
                  placeholder="写下你的评论..."
                  maxlength="500"
                  show-word-limit
                />
                <div class="comment-form-actions">
                  <!-- <el-button
                    v-if="userStore.isLoggedIn"
                    class="comment-submit-btn"
                    @click="submitComment"
                  >
                    发表评论
                  </el-button> -->
                  <HaloButton content="发表评论" 
                  size="large"  
                  v-if="userStore.isLoggedIn"
                  @click="submitComment"/>
                  <el-button
                    v-else
                    class="comment-login-btn"
                    @click="router.push('/login')"
                  >
                    去登录发表评论
                  </el-button>
                </div>
              </div>

              <!-- 评论列表 -->
              <div class="comments-list" v-loading="commentsLoading">
                <CommentItem
                  v-for="comment in comments"
                  :key="comment.id"
                  :ref="(el) => setCommentRef(el, comment.id)"
                  :comment="comment"
                  :article-id="Number(route.params.id)"
                  :current-user-id="currentUserId"
                  @like="handleLikeComment"
                  @reply="replyComment"
                  @delete="handleCommentDelete"
                />
                <el-empty v-if="!commentsLoading && !comments.length" description="暂无评论，快来抢沙发吧~" />
              </div>
            </div>
          </div>
          <el-empty v-else description="文章不存在" />

          <!-- 文章操作按钮组 -->
          <ArticleActions
            v-if="article"
            :article-id="article.id"
            :likes="article.likes || 0"
            :favorites="article.favorites || 0"
            :title="article.title"
            @update:likes="article.likes = $event"
            @update:favorites="article.favorites = $event"
          />

          <!-- 回复评论对话框-->
        </template>
      </el-skeleton>
    </div>
  </div>
  <ReplayComment
    v-model:reployDialogVisible="reployDialogVisible"
    :replyTargetContent="targetComment.content"
    :replyTargetId="targetComment.id"
    :articleId="Number(route.params.id)"
    @submit="handleReplySubmit" />
</template>

<style scoped>
/* 复制按钮样式 - 更紧凑的版本 */
.article-content :deep(.paper-copy-btn) {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 5px 10px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(255, 255, 255, 0.85) 100%);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 16px;
  cursor: pointer;
  font-size: 11px;
  font-weight: 500;
  color: #4a5568;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow:
    0 2px 8px rgba(0, 0, 0, 0.1),
    inset 0 0 0 1px rgba(255, 255, 255, 0.8);
  min-width: 0; /* 关键：移除最小宽度限制 */
  justify-content: center;
  white-space: nowrap; /* 防止文字换行 */
  position: absolute;
  top: 10px;
  right: 10px;
}

.article-content :deep(.paper-copy-btn:hover) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.95) 0%, rgba(118, 75, 162, 0.85) 100%);
  border-color: transparent;
  color: white;
  transform: translateY(-1px);
  box-shadow: 
    0 4px 12px rgba(102, 126, 234, 0.25),
    inset 0 0 0 1px rgba(255, 255, 255, 0.3);
}

.article-content :deep(.paper-copy-btn:active) {
  transform: translateY(0);
}

.article-content :deep(.paper-copy-btn.copied) {
  background: linear-gradient(135deg, rgba(72, 187, 120, 0.95) 0%, rgba(56, 161, 105, 0.85) 100%);
  color: white;
  border-color: transparent;
  box-shadow: 
    0 2px 8px rgba(72, 187, 120, 0.3),
    inset 0 0 0 1px rgba(255, 255, 255, 0.3);
}

/* 两页纸图标 - 缩小尺寸 */
.article-content :deep(.papers-icon) {
  position: relative;
  width: 14px;
  height: 18px;
  perspective: 600px;
  flex-shrink: 0; /* 防止图标被压缩 */
}

.article-content :deep(.paper) {
  position: absolute;
  width: 100%;
  height: 100%;
  background: #fff;
  border-radius: 1px;
  box-shadow: 
    0 1px 2px rgba(0, 0, 0, 0.1),
    inset 0 0 0 1px rgba(0, 0, 0, 0.05);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: bottom center;
}

.article-content :deep(.paper.front) {
  z-index: 2;
  transform: rotateX(0deg);
}

.article-content :deep(.paper.back) {
  z-index: 1;
  transform: rotateX(-10deg) translateY(-2px);
  opacity: 0.8;
  filter: blur(0.5px);
}

.article-content :deep(.paper-copy-btn:hover .paper.front) {
  transform: rotateX(5deg) translateY(-1px);
}

.article-content :deep(.paper-copy-btn:hover .paper.back) {
  transform: rotateX(-15deg) translateY(-3px);
  opacity: 0.6;
}

.article-content :deep(.paper-copy-btn.copied .paper.front) {
  transform: rotateX(20deg) translateY(-2px);
}

.article-content :deep(.paper-copy-btn.copied .paper.back) {
  transform: rotateX(-25deg) translateY(-4px);
  opacity: 0.4;
}

.article-content :deep(.paper-content) {
  padding: 2px;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.article-content :deep(.line) {
  height: 1.5px;
  background: rgba(102, 126, 234, 0.2);
  border-radius: 1px;
}

.article-content :deep(.line-1) { width: 100%; }
.article-content :deep(.line-2) { width: 80%; align-self: flex-end; }
.article-content :deep(.line-3) { width: 60%; }

.article-content :deep(.paper-copy-btn.copied .line) {
  background: rgba(255, 255, 255, 0.6);
}

.article-content :deep(.copy-text) {
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.2px;
  transition: color 0.3s;
  flex-shrink: 0; /* 防止文字被压缩 */
}

.article-content :deep(.paper-copy-btn.copied .copy-text) {
  color: white;
}

.article-content :deep(.success-check) {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 4px;
  animation: checkIn 0.3s ease-out;
}
/* 暗黑模式适配 */
html.dark .article-content :deep(.paper-copy-btn) {
  background: linear-gradient(135deg, rgba(45, 55, 72, 0.95) 0%, rgba(30, 41, 59, 0.85) 100%);
  color: #a0aec0;
  border: 1px solid rgba(160, 174, 192, 0.2);
  box-shadow:
    0 2px 8px rgba(0, 0, 0, 0.3),
    inset 0 0 0 1px rgba(255, 255, 255, 0.1);
}

html.dark .article-content :deep(.paper-copy-btn:hover) {
  background: linear-gradient(135deg, rgba(56, 178, 172, 0.95) 0%, rgba(49, 151, 149, 0.85) 100%);
  color: white;
}

html.dark .article-content :deep(.paper) {
  background: #2d3748;
  box-shadow:
    0 1px 2px rgba(0, 0, 0, 0.3),
    inset 0 0 0 1px rgba(255, 255, 255, 0.05);
}

html.dark .article-content :deep(.line) {
  background: rgba(160, 174, 192, 0.2);
}

html.dark .article-content :deep(.paper-copy-btn.copied) {
  background: linear-gradient(135deg, rgba(72, 187, 120, 0.95) 0%, rgba(56, 161, 105, 0.85) 100%);
}

/* 代码块上的复制按钮特殊样式 */
.article-content :deep(pre .copy-btn-wrapper) {
  position: absolute;
  top: 10px;
  right: 10px;
}
.article-detail-page {
  min-height: calc(100vh - 140px);
  padding: 50px 0;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eef5 100%);
}

html.dark .article-detail-page {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
}

.article-detail {
  max-width: 840px;
  margin: 0 auto;
  background: #fff;
  border-radius: 20px;
  padding: 50px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.5);
}

html.dark .article-detail {
  background: #2d3748;
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3);
}

.article-header {
  margin-bottom: 40px;
  padding-bottom: 35px;
  border-bottom: 2px solid #f0f0f0;
  position: relative;
}

.article-header::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 80px;
  height: 2px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.article-title {
  font-size: 38px;
  font-weight: 800;
  color: #1a202c;
  margin-bottom: 24px;
  line-height: 1.3;
  letter-spacing: -0.5px;
}

html.dark .article-title {
  color: #f7fafc;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  margin-bottom: 24px;
  color: #a0aec0;
  font-size: 14px;
  align-items: center;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  transition: color 0.3s;
}

.meta-item:hover {
  color: #667eea;
}

.category {
  color: #667eea;
  font-weight: 600;
  padding: 6px 16px;
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
  border-radius: 20px;
  font-size: 13px;
}

.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.article-cover {
  width: 100%;
  margin-bottom: 40px;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
  transition: transform 0.4s;
}

.article-cover:hover {
  transform: scale(1.01);
}

.article-cover img {
  width: 100%;
  height: auto;
}

.article-body {
  margin-bottom: 50px;
}

/* Markdown 内容样式 */
.article-content {
  font-size: 16px;
  line-height: 1.8;
  color: #2d3748;
}

.article-content :deep(h1),
.article-content :deep(h2),
.article-content :deep(h3),
.article-content :deep(h4),
.article-content :deep(h5),
.article-content :deep(h6) {
  margin-top: 32px;
  margin-bottom: 16px;
  font-weight: 600;
  color: #1a202c;
  line-height: 1.4;
}

.article-content :deep(h1) {
  font-size: 26px;
  margin-top: 40px;
}

.article-content :deep(h2) {
  font-size: 22px;
}

.article-content :deep(h3) {
  font-size: 20px;
}

.article-content :deep(p) {
  margin-bottom: 16px;
}

/* 代码块样式 - 使用 Atom One Dark 主题 */
/* 自定义代码块样式（Markdown 生成的） */
.article-content :deep(.code-block-wrapper) {
  position: relative;
  margin: 28px 0;
  border-radius: 8px;
  overflow: hidden;
  background: #282c34;
  border: 1px solid #3e4451;
}

/* wangEditor 生成的代码块样式 */
.article-content :deep(pre) {
  position: relative;
  margin: 28px 0;
  border-radius: 8px;
  overflow: hidden;
  background: #282c34;
  border: 1px solid #3e4451;
  padding: 0;
}

.article-content :deep(pre code) {
  display: block;
  padding: 16px;
  background: transparent !important;
  color: #abb2bf;
  font-size: 14px;
  line-height: 1.6;
  font-family: 'Fira Code', 'Monaco', 'Consolas', monospace !important;
  white-space: pre-wrap !important;
  word-wrap: break-word !important;
  overflow-wrap: break-word !important;
}

/* 为 wangEditor 代码块添加复制按钮容器 */
.article-content :deep(pre:hover) .copy-btn-wrapper {
  opacity: 1;
}

/* 代码块头部 */
.article-content :deep(.code-header) {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  background: #21252b;
  border-bottom: 1px solid #3e4451;
  font-size: 13px;
}

.article-content :deep(.code-language) {
  color: #abb2bf;
  font-weight: 500;
  text-transform: capitalize;
}

/* 复制按钮 */
.article-content :deep(.code-copy-btn) {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  background: transparent;
  border: 1px solid #4b5363;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  color: #abb2bf;
  transition: all 0.2s;
}

.article-content :deep(.code-copy-btn:hover) {
  background: #3e4451;
  border-color: #5c6370;
  color: #fff;
}

.article-content :deep(.code-copy-btn.copied) {
  background: #4caf50;
  border-color: #4caf50;
  color: #fff;
}

.article-content :deep(.code-copy-btn svg) {
  fill: currentColor;
}

/* 代码块主体 */
.article-content :deep(.code-body) {
  display: flex;
  overflow-x: auto;
}

/* 行号区域 */
.article-content :deep(.code-line-numbers) {
  display: flex;
  flex-direction: column;
  padding: 16px 8px;
  padding-right: 16px;
  background: #21252b;
  border-right: 1px solid #3e4451;
  user-select: none;
  min-width: 50px;
  text-align: right;
}

.article-content :deep(.code-line-number) {
  display: block;
  line-height: 1.6;
  font-size: 14px;
  color: #636d83;
  font-family: 'Fira Code', 'Monaco', 'Consolas', monospace;
}

/* 代码内容区域 */
.article-content :deep(.code-content) {
  flex: 1;
  padding: 16px;
  overflow-x: auto;
}

.article-content :deep(.code-content code.hljs) {
  background: transparent !important;
  padding: 0 !important;
  margin: 0 !important;
  font-size: 14px;
  line-height: 1.6;
  font-family: 'Fira Code', 'Monaco', 'Consolas', monospace !important;
  color: #abb2bf;
  white-space: pre-wrap !important;
  word-wrap: break-word !important;
  overflow-wrap: break-word !important;
}

/* 行内代码样式 - 排除代码块内的代码（hljs 类） */
.article-content :deep(code):not(.hljs) {
  background: #f1f5f9;
  color: #e53e3e;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Fira Code', 'Monaco', 'Consolas', monospace;
  font-size: 14px;
}

html.dark .article-content :deep(code):not(.hljs) {
  background: #3d4759;
  color: #f87171;
}

.article-content :deep(blockquote) {
  border-left: 4px solid #667eea;
  padding-left: 20px;
  margin: 24px 0;
  color: #4a5568;
  font-style: italic;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  padding: 16px 20px;
  border-radius: 0 12px 12px 0;
}

.article-content :deep(ul),
.article-content :deep(ol) {
  margin: 16px 0;
  padding-left: 24px;
}

.article-content :deep(li) {
  margin-bottom: 8px;
  line-height: 1.8;
}

.article-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 24px 0;
  overflow-x: auto;
}

.article-content :deep(th),
.article-content :deep(td) {
  padding: 12px 16px;
  border: 1px solid #e2e8f0;
  text-align: left;
}

.article-content :deep(th) {
  background: #f7fafc;
  font-weight: 600;
}

.article-content :deep(tr:hover) {
  background: #f7fafc;
}

.article-content :deep(a) {
  color: #667eea;
  text-decoration: none;
  transition: color 0.3s;
}

.article-content :deep(a:hover) {
  color: #764ba2;
  text-decoration: underline;
}

.article-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 24px 0;
}

/* 复制按钮容器 - 代码块内部右上角 */
.article-content :deep(.copy-btn-wrapper) {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 10;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.article-content :deep(pre:hover .copy-btn-wrapper) {
  opacity: 1;
}

/* 复制按钮样式 */
.article-content :deep(.paper-copy-button) {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(255, 255, 255, 0.85) 100%);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 16px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  color: #4a5568;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.1),
    inset 0 0 0 1px rgba(255, 255, 255, 0.8);
  min-width: 70px;
  justify-content: center;
}

.article-content :deep(.paper-copy-button:hover) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.95) 0%, rgba(118, 75, 162, 0.85) 100%);
  border-color: transparent;
  color: white;
  transform: translateY(-1px);
  box-shadow: 
    0 4px 12px rgba(102, 126, 234, 0.25),
    inset 0 0 0 1px rgba(255, 255, 255, 0.3);
}

.article-content :deep(.paper-copy-button:active) {
  transform: translateY(0);
}

.article-content :deep(.paper-copy-button.copied) {
  background: linear-gradient(135deg, rgba(72, 187, 120, 0.95) 0%, rgba(56, 161, 105, 0.85) 100%);
  color: white;
  border-color: transparent;
  box-shadow: 
    0 2px 8px rgba(72, 187, 120, 0.3),
    inset 0 0 0 1px rgba(255, 255, 255, 0.3);
}

/* 两页纸图标 */
.article-content :deep(.papers-icon) {
  position: relative;
  width: 16px;
  height: 20px;
  perspective: 600px;
}

.article-content :deep(.paper) {
  position: absolute;
  width: 100%;
  height: 100%;
  background: #fff;
  border-radius: 1px;
  box-shadow: 
    0 1px 2px rgba(0, 0, 0, 0.1),
    inset 0 0 0 1px rgba(0, 0, 0, 0.05);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: bottom center;
}

.article-content :deep(.paper.front) {
  z-index: 2;
  transform: rotateX(0deg);
}

.article-content :deep(.paper.back) {
  z-index: 1;
  transform: rotateX(-10deg) translateY(-2px);
  opacity: 0.8;
  filter: blur(0.5px);
}

.article-content :deep(.paper-copy-button:hover .paper.front) {
  transform: rotateX(5deg) translateY(-1px);
}

.article-content :deep(.paper-copy-button:hover .paper.back) {
  transform: rotateX(-15deg) translateY(-3px);
  opacity: 0.6;
}

.article-content :deep(.paper-copy-button.copied .paper.front) {
  transform: rotateX(20deg) translateY(-2px);
}

.article-content :deep(.paper-copy-button.copied .paper.back) {
  transform: rotateX(-25deg) translateY(-4px);
  opacity: 0.4;
}

.article-content :deep(.paper-content) {
  padding: 3px;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.article-content :deep(.line) {
  height: 1.5px;
  background: rgba(102, 126, 234, 0.2);
  border-radius: 1px;
}

.article-content :deep(.line-1) { width: 100%; }
.article-content :deep(.line-2) { width: 80%; align-self: flex-end; }
.article-content :deep(.line-3) { width: 60%; }

.article-content :deep(.paper-copy-button.copied .line) {
  background: rgba(255, 255, 255, 0.6);
}

.article-content :deep(.copy-text) {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.3px;
  transition: color 0.3s;
}

.article-content :deep(.paper-copy-button.copied .copy-text) {
  color: white;
}

.article-content :deep(.success-check) {
  display: flex;
  align-items: center;
  justify-content: center;
  animation: checkIn 0.3s ease-out;
}

@keyframes checkIn {
  0% {
    opacity: 0;
    transform: scale(0);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

/* 暗黑模式适配 */
html.dark .article-content :deep(.paper-copy-button) {
  background: linear-gradient(135deg, rgba(45, 55, 72, 0.95) 0%, rgba(30, 41, 59, 0.85) 100%);
  color: #a0aec0;
  border: 1px solid rgba(160, 174, 192, 0.2);
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.3),
    inset 0 0 0 1px rgba(255, 255, 255, 0.1);
}

html.dark .article-content :deep(.paper-copy-button:hover) {
  background: linear-gradient(135deg, rgba(56, 178, 172, 0.95) 0%, rgba(49, 151, 149, 0.85) 100%);
  color: white;
}

html.dark .article-content :deep(.paper) {
  background: #2d3748;
  box-shadow: 
    0 1px 2px rgba(0, 0, 0, 0.3),
    inset 0 0 0 1px rgba(255, 255, 255, 0.05);
}

html.dark .article-content :deep(.line) {
  background: rgba(160, 174, 192, 0.2);
}

html.dark .article-content :deep(.paper-copy-button.copied) {
  background: linear-gradient(135deg, rgba(72, 187, 120, 0.95) 0%, rgba(56, 161, 105, 0.85) 100%);
}

/* 移动端调整 */
@media (max-width: 768px) {
     .article-content :deep(.copy-btn-wrapper) {
    opacity: 1;
    top: 8px;
    right: 8px;
  }
  
  .article-content :deep(.paper-copy-btn) {
    padding: 4px 8px;
    font-size: 10px;
  }
  
  .article-content :deep(.papers-icon) {
    width: 12px;
    height: 16px;
  }
  
  .article-content :deep(.copy-text) {
    font-size: 9px;
  }

  .article-content :deep(.copy-button-container) {
    opacity: 1;
    top: 8px;
    right: 8px;
  }

  .article-content :deep(.paper-copy-button) {
    padding: 5px 10px;
    min-width: 60px;
    font-size: 11px;
  }

  .article-content :deep(.papers-icon) {
    width: 14px;
    height: 18px;
  }

  .article-content :deep(.copy-text) {
    font-size: 10px;
  }
}

/* 暗黑模式 */
html.dark .article-content {
  color: #e2e8f0;
}

html.dark .article-content :deep(h1),
html.dark .article-content :deep(h2),
html.dark .article-content :deep(h3),
html.dark .article-content :deep(h4),
html.dark .article-content :deep(h5),
html.dark .article-content :deep(h6) {
  color: #f7fafc;
}

html.dark .article-content :deep(h1),
html.dark .article-content :deep(h2) {
  border-bottom-color: #4a5568;
}

html.dark .article-content :deep(p) {
  color: #cbd5e0;
}

html.dark .article-content :deep(code) {
  background: #2d3748;
  color: #fc8181;
}

html.dark .article-content :deep(pre) {
  background: #1a202c;
}

html.dark .article-content :deep(blockquote) {
  border-left-color: #764ba2;
  color: #a0aec0;
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
}

html.dark .article-content :deep(th),
html.dark .article-content :deep(td) {
  border-color: #4a5568;
}

html.dark .article-content :deep(th) {
  background: #2d3748;
}

html.dark .article-content :deep(tr:hover) {
  background: #2d3748;
}

html.dark .article-content :deep(table) {
  border-color: #4a5568;
}

/* 评论区样式 */
.comments-section {
  margin-top: 50px;
  padding-top: 40px;
  border-top: 2px solid #f0f0f0;
}

html.dark .comments-section {
  border-top-color: #4a5568;
}

.comments-header {
  margin-bottom: 30px;
}

.comments-header h3 {
  font-size: 24px;
  font-weight: 700;
  color: #1a202c;
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0;
}

html.dark .comments-header h3 {
  color: #f7fafc;
}

.comment-form {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  padding: 24px;
  border-radius: 12px;
  margin-bottom: 30px;
}

html.dark .comment-form {
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
}

.comment-form-actions {
  margin-top: 16px;
  text-align: right;
}

/* 评论表单按钮 - 淡蓝色系 */
.comment-submit-btn {
  background: linear-gradient(135deg, #60a5fa 0%, #5b8ff9 100%);
  border: none;
  color: #fff;
  padding: 10px 24px;
  border-radius: 6px;
  font-weight: 500;
  transition: all 0.3s;
}

.comment-submit-btn:hover {
  background: linear-gradient(135deg, #5b8ff9 0%, #4a7df7 100%);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(96, 165, 250, 0.35);
}

html.dark .comment-submit-btn {
  background: linear-gradient(135deg, #60a5fa 0%, #5b8ff9 100%);
}

html.dark .comment-submit-btn:hover {
  background: linear-gradient(135deg, #5b8ff9 0%, #4a7df7 100%);
  box-shadow: 0 4px 12px rgba(96, 165, 250, 0.45);
}

.comment-login-btn {
  background: linear-gradient(135deg, #a8dadc 0%, #92c5f7 100%);
  border: none;
  color: #fff;
  padding: 10px 24px;
  border-radius: 6px;
  font-weight: 500;
  transition: all 0.3s;
}

.comment-login-btn:hover {
  background: linear-gradient(135deg, #92c5f7 0%, #7dd3fc 100%);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(168, 173, 220, 0.35);
}

html.dark .comment-login-btn {
  background: linear-gradient(135deg, #a8dadc 0%, #92c5f7 100%);
}

html.dark .comment-login-btn:hover {
  background: linear-gradient(135deg, #92c5f7 0%, #7dd3fc 100%);
  box-shadow: 0 4px 12px rgba(168, 173, 220, 0.45);
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.article-footer {
  padding-top: 35px;
  border-top: 2px solid #f0f0f0;
  position: relative;
}

.article-nav {
  display: flex;
  justify-content: space-between;
  gap: 24px;
}

.article-nav a {
  flex: 1;
  padding: 18px 24px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 12px;
  color: #4a5568;
  font-size: 14px;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(255, 255, 255, 0.5);
  font-weight: 500;
}

.article-nav a:hover {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
}

.prev-article {
  text-align: left !important;
  justify-content: flex-start !important;
}

.next-article {
  text-align: right !important;
  justify-content: flex-end !important;
}

@media (max-width: 768px) {
  .article-detail-page {
    padding: 30px 0;
  }

  .article-detail {
    padding: 30px 20px;
    border-radius: 16px;
  }

  .article-title {
    font-size: 28px;
  }

  .article-meta {
    gap: 16px;
  }

  .article-nav {
    flex-direction: column;
  }

  .article-nav a {
    text-align: center !important;
    justify-content: center !important;
  }
}

</style>