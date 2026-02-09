<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getSeriesDetail, getArticleDetail } from '@/api'
import { marked } from 'marked'

const route = useRoute()
const router = useRouter()

const series = ref(null)
const allArticles = ref([])
const article = ref(null)
const loading = ref(true)
const renderedContent = ref('')

const currentArticleIndex = computed(() => {
  return allArticles.value.findIndex(a => a.id === parseInt(route.params.articleId))
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
    // 加载系列详情
    const seriesRes = await getSeriesDetail(route.params.id)
    series.value = seriesRes.data

    // 收集所有文�?
    allArticles.value = []
    series.value.chapters?.forEach(chapter => {
      chapter.articles?.forEach(article => {
        allArticles.value.push({
          ...article,
          chapterName: chapter.name,
          chapterId: chapter.id
        })
      })
    })

    // 加载当前文章详情
    const articleRes = await getArticleDetail(route.params.articleId)
    article.value = articleRes.data
    renderedContent.value = marked(articleRes.data.content || '')
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

const handleArticleClick = (articleId) => {
  router.replace({
    name: 'ChapterArticle',
    params: { id: route.params.id, articleId }
  })
  window.scrollTo({ top: 0, behavior: 'smooth' })
  // 重新加载文章
  getArticleDetail(articleId).then(res => {
    article.value = res.data
    renderedContent.value = marked(res.data.content || '')
  })
}

const handleNavClick = (articleId) => {
  handleArticleClick(articleId)
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="chapter-article-page">
    <div class="container">
      <!-- 左侧目录 -->
      <aside class="sidebar">
        <div class="sidebar-content">
          <div class="sidebar-header">
            <h3 class="sidebar-title">{{ series?.name }}</h3>
            <p class="sidebar-desc">{{ series?.description }}</p>
          </div>

          <div class="chapters-list">
            <div
              v-for="chapter in series?.chapters"
              :key="chapter.id"
              class="sidebar-chapter"
            >
              <h4 class="chapter-heading">{{ chapter.name }}</h4>
              <div
                v-for="article in chapter.articles"
                :key="article.id"
                class="article-item"
                :class="{ active: article.id === parseInt(route.params.articleId) }"
                @click="handleArticleClick(article.id)"
              >
                <span class="article-dot"></span>
                <span class="article-text">{{ article.title }}</span>
              </div>
            </div>
          </div>
        </div>
      </aside>

      <!-- 右侧文章内容 -->
      <main class="main-content">
        <div class="article-wrapper" v-if="article">
          <!-- 文章头部 -->
          <div class="article-header">
            <h1 class="article-title">{{ article.title }}</h1>
            <div class="article-meta">
              <span class="category">{{ article.category?.name || '未分类' }}</span>
              <span class="date">📅 {{ new Date(article.created_at).toLocaleDateString() }}</span>
              <span class="views">👁 {{ article.views || 0 }}</span>
            </div>
          </div>

          <!-- 文章内容 -->
          <div class="article-content">
            <div class="markdown-body" v-html="renderedContent"></div>
          </div>

          <!-- 上一�?下一�?-->
          <div class="article-navigation">
            <div
              v-if="prevArticle"
              class="nav-link prev"
              @click="handleNavClick(prevArticle.id)"
            >
              <span class="nav-label">上一篇</span>
              <span class="nav-title">{{ prevArticle.title }}</span>
            </div>
            <div
              v-if="nextArticle"
              class="nav-link next"
              @click="handleNavClick(nextArticle.id)"
            >
              <span class="nav-title">{{ nextArticle.title }}</span>
              <span class="nav-label">下一篇?</span>
            </div>
          </div>
        </div>

        <!-- 加载状�?-->
        <div class="loading-container" v-if="loading">
          <el-skeleton animated>
            <template #template>
              <el-skeleton-item variant="h1" style="width: 70%; margin-bottom: 24px;" />
              <el-skeleton-item variant="text" style="margin-bottom: 12px;" />
              <el-skeleton-item variant="text" style="margin-bottom: 12px;" />
              <el-skeleton-item variant="text" style="width: 60%; margin-bottom: 32px;" />
              <el-skeleton-item variant="rect" style="height: 300px; margin-bottom: 32px;" />
            </template>
          </el-skeleton>
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped>
.chapter-article-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f8fafc 0%, #fff 100%);
}

html.dark .chapter-article-page {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
}

.chapter-article-page .container {
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

.chapters-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.sidebar-chapter {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.chapter-heading {
  font-size: 14px;
  font-weight: 700;
  color: #667eea;
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.article-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 12px;
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

.sidebar-content .article-item.active {
  background: #60e1f1;
  border-radius: 12px;
  color: #fff;
  box-shadow: 0 2px 8px rgba(96, 225, 241, 0.3);
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
  font-size: 13px;
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

.article-wrapper {
  max-width: 900px;
  margin: 0 auto;
}

.article-header {
  background: linear-gradient(135deg, #fff 0%, #f8fafc 100%);
  border-radius: 20px;
  padding: 32px;
  margin-bottom: 30px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.5);
}

html.dark .article-header {
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.article-title {
  font-size: 36px;
  font-weight: 800;
  color: #1a202c;
  margin: 0 0 24px;
  line-height: 1.3;
  letter-spacing: -0.5px;
}

html.dark .article-title {
  color: #f7fafc;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 20px;
  padding-top: 20px;
  border-top: 2px solid #e2e8f0;
  font-size: 14px;
  color: #718096;
}

html.dark .article-meta {
  border-top-color: #4a5568;
  color: #a0aec0;
}

.category {
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
  color: #667eea;
  font-weight: 600;
  padding: 6px 14px;
  border-radius: 10px;
}

.article-content {
  background: #f8fafc;
  border-radius: 20px;
  padding: 40px;
  margin-bottom: 30px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.5);
}

html.dark .article-content {
  background: #2d3748;
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);
}

.markdown-body {
  color: #2d3748;
  line-height: 1.8;
  font-size: 16px;
}

html.dark .markdown-body {
  color: #e2e8f0;
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3) {
  margin-top: 32px;
  margin-bottom: 16px;
  font-weight: 700;
  color: #1a202c;
}

html.dark .markdown-body :deep(h1),
html.dark .markdown-body :deep(h2),
html.dark .markdown-body :deep(h3) {
  color: #f7fafc;
}

.markdown-body :deep(h1) {
  font-size: 28px;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 12px;
}

html.dark .markdown-body :deep(h1) {
  border-bottom-color: #4a5568;
}

.markdown-body :deep(h2) {
  font-size: 24px;
}

.markdown-body :deep(h3) {
  font-size: 20px;
}

.markdown-body :deep(p) {
  margin-bottom: 16px;
}

.markdown-body :deep(code) {
  background: #f1f5f9;
  color: #1e40af;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 14px;
}

html.dark .markdown-body :deep(code) {
  background: rgba(102, 126, 234, 0.2);
  color: #a5b4fc;
}

.markdown-body :deep(pre) {
  background: #1a202c;
  color: #e2e8f0;
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 20px 0;
}

.markdown-body :deep(pre code) {
  background: transparent;
  color: inherit;
  padding: 0;
}

.markdown-body :deep(blockquote) {
  border-left: 4px solid #667eea;
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
  padding: 16px 20px;
  margin: 20px 0;
  border-radius: 8px;
  color: #4a5568;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  margin: 16px 0;
  padding-left: 24px;
}

.markdown-body :deep(li) {
  margin-bottom: 8px;
  line-height: 1.6;
}

.markdown-body :deep(a) {
  color: #667eea;
  text-decoration: none;
  transition: all 0.3s;
}

.markdown-body :deep(a:hover) {
  color: #764ba2;
  text-decoration: underline;
}

.article-navigation {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 30px;
}

.nav-link {
  background: #f8fafc;
  border-radius: 16px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

html.dark .nav-link {
  background: #2d3748;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);
}

.nav-link:hover {
  transform: translateY(-4px);
  border-color: rgba(102, 126, 234, 0.3);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.2);
}

.nav-link.prev {
  text-align: left;
}

.nav-link.next {
  text-align: right;
}

.nav-label {
  font-size: 13px;
  color: #667eea;
  font-weight: 600;
}

.nav-title {
  font-size: 15px;
  font-weight: 600;
  color: #1a202c;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

html.dark .nav-title {
  color: #f7fafc;
}

.loading-container {
  padding: 40px 0;
}

@media (max-width: 1200px) {
  .sidebar {
    width: 280px;
  }
}

@media (max-width: 900px) {
  .chapter-article-page .container {
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
    padding: 24px;
  }

  .article-navigation {
    grid-template-columns: 1fr;
  }

  .nav-link {
    padding: 20px;
  }

  .article-meta {
    flex-wrap: wrap;
    gap: 12px;
  }
}
</style>
