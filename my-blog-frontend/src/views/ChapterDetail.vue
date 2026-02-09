<script setup>
import { ref, onMounted, watch, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getSeriesDetail, getArticleDetail } from '@/api'
import { markdownToHtml } from '@/utils/format'

const route = useRoute()
const router = useRouter()

const series = ref(null)
const chapter = ref(null)
const allArticles = ref([])
const article = ref(null)
const loading = ref(true)
const renderedContent = ref('')

const currentArticleIndex = computed(() => {
  return allArticles.value.findIndex(a => a.id === parseInt(route.query.articleId))
})

const prevArticle = computed(() => {
  if (currentArticleIndex.value > 0) {
    return allArticles.value[currentArticleIndex.value - 1]
  }
  return null
})

const nextArticle = computed(() => {
  if (currentArticleIndex.value < allArticles.value.length - 1) {
    return allArticles.value[currentArticleIndex.value + 1]
  }
  return null
})

const loadData = async () => {
  loading.value = true
  try {
    console.log('ChapterDetail - 开始加载数据')
    console.log('route.params:', route.params)
    console.log('route.query:', route.query)
    
    // 加载系列详情
    const seriesRes = await getSeriesDetail(route.params.seriesId)
    series.value = seriesRes.data
    console.log('系列数据:', series.value)
    console.log('sections:', series.value.sections)

    // 在所有章节标题的子章节中查找当前子章节
    let foundSubchapter = null
    if (series.value.sections && series.value.sections.length > 0) {
      for (const section of series.value.sections) {
        console.log('检查section:', section.name, 'subchapters:', section.subchapters)
        if (section.subchapters && section.subchapters.length > 0) {
          const subchapter = section.subchapters.find(sc => sc.id === parseInt(route.params.chapterId))
          console.log('查找 chapterId:', route.params.chapterId, '找到子章节:', subchapter)
          if (subchapter) {
            foundSubchapter = subchapter
            break
          }
        }
      }
    }
    
    if (!foundSubchapter) {
      console.error('未找到子章节, chapterId:', route.params.chapterId)
    }
    
    chapter.value = foundSubchapter
    console.log('当前章节:', chapter.value)

    // 获取子章节下的所有文章
    allArticles.value = foundSubchapter?.articles || []
    console.log('文章列表:', allArticles.value)

    // 默认加载第一篇文章
    if (allArticles.value.length > 0 && !route.query.articleId) {
      console.log('加载默认文章:', allArticles.value[0].id)
      router.replace({
        name: 'ChapterDetail',
        params: {
          seriesId: route.params.seriesId,
          chapterId: route.params.chapterId
        },
        query: { articleId: allArticles.value[0].id }
      })
    } else if (route.query.articleId) {
      await loadArticle(route.query.articleId)
    }
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

const loadArticle = async (articleId) => {
  try {
    const res = await getArticleDetail(articleId)
    article.value = res.data
    renderedContent.value = markdownToHtml(res.data.content || '')

    // 等待 DOM 更新后为代码块添加复制按钮
    await nextTick()
    enhanceCodeBlocks()
  } catch (error) {
    console.error('加载文章失败:', error)
  }
}

// 为 Markdown 生成的代码块添加复制按钮功能
const enhanceCodeBlocks = () => {
  const codeBlocks = document.querySelectorAll('.article-content .code-block-wrapper')
  codeBlocks.forEach((block) => {
    const copyBtn = block.querySelector('.code-copy-btn')
    if (copyBtn && !copyBtn.dataset.enhanced) {
      copyBtn.dataset.enhanced = 'true'
      copyBtn.addEventListener('click', async () => {
        const codeContent = block.querySelector('.code-content code')
        try {
          await navigator.clipboard.writeText(codeContent.textContent)
          const originalContent = copyBtn.innerHTML
          copyBtn.innerHTML = `
            <svg viewBox="0 0 24 24" width="14" height="14">
              <path fill="currentColor" d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
            </svg>
            已复制
          `
          copyBtn.classList.add('copied')
          setTimeout(() => {
            copyBtn.innerHTML = originalContent
            copyBtn.classList.remove('copied')
          }, 2000)
        } catch (error) {
          console.error('复制失败:', error)
        }
      })
    }
  })
}

const handleArticleClick = (articleId) => {
  router.push({
    name: 'ChapterDetail',
    params: {
      seriesId: route.params.seriesId,
      chapterId: route.params.chapterId
    },
    query: { articleId }
  })
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const handleNavClick = (articleId) => {
  handleArticleClick(articleId)
}

// 监听路由参数变化，重新加载文章
watch(() => route.query.articleId, (newArticleId) => {
  if (newArticleId) {
    loadArticle(newArticleId)
  }
})

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="chapter-detail-page">
    <div class="container">
      <!-- 左侧目录 -->
      <aside class="sidebar">
        <div class="sidebar-content">
          <div class="sidebar-header">
            <h3 class="sidebar-title">{{ chapter?.name }}</h3>
            <p class="sidebar-desc">{{ chapter?.description }}</p>
          </div>

          <div class="articles-list">
            <div
              v-for="article in allArticles"
              :key="article.id"
              class="article-item"
              :class="{ active: article.id === parseInt(route.query.articleId) }"
              @click="handleArticleClick(article.id)"
            >
              <span class="article-dot"></span>
              <span class="article-text">{{ article.title }}</span>
            </div>
          </div>
        </div>
      </aside>

      <!-- 右侧文章内容 -->
      <main class="main-content">
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
              <!-- 文章头部 -->
              <div class="article-header">
                <h1 class="article-title">{{ article.title }}</h1>
                <div class="article-meta">
                  <span class="meta-item">
                    <span>📅</span>
                    {{ new Date(article.createdAt).toLocaleDateString() }}
                  </span>
                  <span class="meta-item">
                    <span>👁</span>
                    {{ article.views || 0 }} 阅读
                  </span>
                  <span class="category" v-if="article.category">
                    分类: {{ article.category.name }}
                  </span>
                  <span class="meta-item" v-if="article.author">
                    作者: {{ article.author.name }}
                  </span>
                </div>

                <!-- 文章标签 -->
                <div class="article-tags" v-if="article.tags && article.tags.length">
                  <el-tag
                    v-for="tag in article.tags"
                    :key="tag.name"
                    type="info"
                    effect="plain"
                    size="small"
                  >
                    {{ tag.name }}
                  </el-tag>
                </div>
              </div>

              <!-- 文章封面 -->
              <div class="article-cover" v-if="article.cover">
                <img :src="article.cover" :alt="article.title" />
              </div>

              <!-- 文章内容 -->
              <div class="article-body">
                <div
                  class="article-content"
                  v-html="renderedContent"
                ></div>
              </div>

              <!-- 上一篇/下一篇 -->
              <div class="article-nav">
                <div
                  v-if="prevArticle"
                  class="prev-article"
                  @click="handleNavClick(prevArticle.id)"
                >
                  <span>←</span>
                  <span>{{ prevArticle.title }}</span>
                </div>
                <div
                  v-if="nextArticle"
                  class="next-article"
                  @click="handleNavClick(nextArticle.id)"
                >
                  <span>{{ nextArticle.title }}</span>
                  <span>→</span>
                </div>
              </div>
            </div>
            <el-empty v-else description="文章不存在" />
          </template>
        </el-skeleton>
      </main>
    </div>
  </div>
</template>

<style scoped>
.chapter-detail-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f8fafc 0%, #fff 100%);
}

html.dark .chapter-detail-page {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
}

.chapter-detail-page .container {
  display: flex;
  gap: 30px;
  padding: 40px 20px;
  max-width: 1600px;
}

/* 左侧目录 */
.sidebar {
  width: 320px;
  flex-shrink: 0;
}

.sidebar-content {
  position: sticky;
  top: 20px;
  background: #f8fafc;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.5);
  max-height: calc(100vh - 80px);
  overflow-y: auto;
}

html.dark .sidebar-content {
  background: #2d3748;
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.25);
}

.sidebar-content::-webkit-scrollbar {
  width: 6px;
}

.sidebar-content::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 3px;
}

html.dark .sidebar-content::-webkit-scrollbar-thumb {
  background: #4a5568;
}

.sidebar-header {
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 2px solid #e2e8f0;
}

html.dark .sidebar-header {
  border-bottom-color: #4a5568;
}

.sidebar-title {
  font-size: 20px;
  font-weight: 700;
  color: #1a202c;
  margin: 0 0 8px;
}

html.dark .sidebar-title {
  color: #f7fafc;
}

.sidebar-desc {
  font-size: 13px;
  color: #718096;
  margin: 0;
  line-height: 1.5;
}

html.dark .sidebar-desc {
  color: #a0aec0;
}

.articles-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.article-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  background: #f8fafc;
}

html.dark .article-item {
  background: #1a202c;
}

.article-item:hover {
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
  transform: translateX(4px);
}

html.dark .article-item:hover {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.2) 0%, rgba(118, 75, 162, 0.2) 100%);
}

.article-item.active {
  background: #86e2ee;
  border-radius: 12px;
  color: #fff;
}

.article-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #cbd5e0;
  flex-shrink: 0;
  margin-top: 7px;
}

html.dark .article-dot {
  background: #4a5568;
}

.article-item.active .article-dot {
  background: #fff;
}

.article-text {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
  color: #4a5568;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

html.dark .article-text {
  color: #a0aec0;
}

.article-item.active .article-text {
  color: #fff;
}

/* 右侧内容 */
.main-content {
  flex: 1;
  min-width: 0;
}

/* 文章详情容器 */
.article-detail {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
  background: #f8fafc;
  border-radius: 20px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.5);
}

html.dark .article-detail {
  background: #2d3748;
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3);
}

/* 文章头部 */
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

/* 文章封面 */
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

/* 文章主体 */
.article-body {
  margin-bottom: 50px;
}

/* Markdown 内容样式 */
.article-content {
  font-size: 16px;
  line-height: 1.8;
  color: #2d3748;
}

html.dark .article-content {
  color: #e2e8f0;
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

html.dark .article-content :deep(h1),
html.dark .article-content :deep(h2),
html.dark .article-content :deep(h3),
html.dark .article-content :deep(h4),
html.dark .article-content :deep(h5),
html.dark .article-content :deep(h6) {
  color: #f7fafc;
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
.article-content :deep(.code-block-wrapper) {
  position: relative;
  margin: 28px 0;
  border-radius: 8px;
  overflow: hidden;
  background: #282c34;
  border: 1px solid #3e4451;
}

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

/* 暗黑模式 */
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

/* 上一篇/下一篇导航 */
.article-nav {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 30px;
  padding-top: 30px;
  border-top: 2px solid #e2e8f0;
}

.prev-article,
.next-article {
  flex: 1;
  padding: 20px;
  background: linear-gradient(135deg, #f8fafc 0%, #fff 100%);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #1a202c;
  font-weight: 600;
}

html.dark .prev-article,
html.dark .next-article {
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
  color: #f7fafc;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);
}

.prev-article:hover,
.next-article:hover {
  transform: translateY(-4px);
  border-color: rgba(102, 126, 234, 0.3);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.2);
}

.prev-article {
  justify-content: flex-start;
}

.next-article {
  justify-content: flex-end;
}

@media (max-width: 1200px) {
  .sidebar {
    width: 280px;
  }
}

@media (max-width: 900px) {
  .chapter-detail-page .container {
    flex-direction: column;
    padding: 20px;
  }

  .sidebar {
    width: 100%;
  }

  .sidebar-content {
    max-height: 400px;
    position: relative;
    top: 0;
  }

  .article-header {
    padding: 24px;
  }

  .article-title {
    font-size: 28px;
  }

  .article-content {
    font-size: 15px;
  }

  .article-nav {
    flex-direction: column;
  }

  .article-meta {
    flex-wrap: wrap;
    gap: 12px;
  }
}
</style>
