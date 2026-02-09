<template>
  <div class="series-page" v-if="!loading">
    <div class="container">
      <!-- 专栏头部 -->
      <div class="series-header">
        <div class="series-icon">{{ series?.name.charAt(0) }}</div>
        <div class="series-info">
          <h1 class="series-title">{{ series?.name }}</h1>
          <p class="series-description">{{ series?.description }}</p>
          <div class="series-stats">
            <span class="stat">
              <span class="stat-number">{{ chapters.length }}</span>
              <span class="stat-label">章节</span>
            </span>
          </div>
        </div>
      </div>

      <!-- 专栏导航 -->
      <div class="series-nav">
        <div
          v-for="chapter in chapters"
          :key="chapter.id"
          class="chapter-item"
          :class="{ active: chapter.id === activeChapterId }"
          @click="handleChapterClick(chapter.id)"
        >
          <div class="chapter-number">{{ chapter.id }}</div>
          <div class="chapter-content">
            <div class="chapter-name">{{ chapter.name }}</div>
            <div class="chapter-desc">{{ chapter.description || '暂无描述' }}</div>
            <div class="chapter-count">{{ chapter.articles?.length || 0 }} 篇文章</div>
          </div>
        </div>
        </div>
      </div>
      <!-- 文章列表 -->
      <div class="chapter-articles" v-if="activeChapter">
        <h2 class="chapter-title">{{ activeChapter.name }}</h2>
        <div class="article-grid">
          <div
            v-for="article in activeChapter.articles"
            :key="article.id"
            class="chapter-article-card"
            @click="handleArticleClick(article)"
          >
            <div class="article-icon">📄</div>
            <div class="article-content">
              <h3 class="article-title">{{ article.title }}</h3>
              <p class="article-summary">{{ article.summary }}</p>
              <div class="article-meta">
                <span class="article-views">👁 {{ article.views || 0 }}</span>
                <span class="article-date">{{ new Date(article.createdAt).toLocaleDateString() }}</span>
              </div>
            </div>
          </div>
        </div>
        <el-empty v-if="!activeChapter.articles?.length" description="该章节暂无文章" />
      </div>
    </div>

  <!-- 加载状态-->
  <div class="loading-container" v-else>
    <el-skeleton animated>
      <template #template>
        <div class="skeleton-header">
          <el-skeleton-item variant="circle" style="width: 80px; height: 80px;" />
          <el-skeleton-item variant="h3" style="width: 60%;" />
          <el-skeleton-item variant="text" style="width: 80%;" />
        </div>
        <div class="skeleton-nav">
          <el-skeleton-item variant="rect" style="height: 80px; margin-bottom: 16px;" />
          <el-skeleton-item variant="rect" style="height: 80px; margin-bottom: 16px;" />
          <el-skeleton-item variant="rect" style="height: 80px;" />
        </div>
        <div class="skeleton-articles">
          <el-skeleton-item variant="rect" style="height: 100px; margin-bottom: 16px;" />
          <el-skeleton-item variant="rect" style="height: 100px; margin-bottom: 16px;" />
          <el-skeleton-item variant="rect" style="height: 100px;" />
        </div>
      </template>
    </el-skeleton>
  </div>
</template>


<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getSeriesDetail } from '@/api'

const route = useRoute()
const router = useRouter()

const series = ref(null)
const chapters = ref([])
const loading = ref(true)
const activeChapterId = ref(null)

const activeChapter = computed(() => {
  if (!activeChapterId.value || !chapters.value.length) return null
  return chapters.value.find(c => c.id === activeChapterId.value)
})

const loadSeriesDetail = async () => {
  loading.value = true
  try {
    const res = await getSeriesDetail(route.params.id)
    series.value = res.data
    chapters.value = res.data.chapters || []
    // 默认选中第一个章
    if (chapters.value.length > 0 && !activeChapterId.value) {
      activeChapterId.value = chapters.value[0].id
    }
  } catch (error) {
    console.error('加载专栏详情失败:', error)
  } finally {
    loading.value = false
  }
}

const handleChapterClick = (chapterId) => {
  activeChapterId.value = chapterId
}

const handleArticleClick = (article) => {
  const url = router.resolve({
    name: 'ChapterArticle',
    params: { articleId: article.id }
  }).href
  window.open(url, '_blank')
}

onMounted(() => {
  loadSeriesDetail()
})
</script>

<style scoped>
.series-page {
  min-height: calc(100vh - 140px);
  padding: 40px 0;
}

html.dark .series-page {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
}

.series-header {
  background: linear-gradient(135deg, #fff 0%, #f8fafc 100%);
  border-radius: 20px;
  padding: 32px;
  margin-bottom: 30px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.5);
  display: flex;
  align-items: center;
  gap: 24px;
}

html.dark .series-header {
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.series-icon {
  width: 80px;
  height: 80px;
  font-size: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 20px;
  background: #3b82f6;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
  flex-shrink: 0;
}

.series-info {
  flex: 1;
}

.series-title {
  font-size: 32px;
  font-weight: 700;
  color: #1a202c;
  margin: 0 0 12px;
}

html.dark .series-title {
  color: #f7fafc;
}

.series-description {
  font-size: 16px;
  color: #718096;
  margin: 0 0 20px;
  line-height: 1.6;
}

html.dark .series-description {
  color: #a0aec0;
}

.series-stats {
  display: flex;
  gap: 24px;
}

.stat {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-number {
  font-size: 24px;
  font-weight: 700;
  color: #667eea;
}

.stat-label {
  font-size: 12px;
  color: #a0aec0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.series-nav {
  margin-bottom: 30px;
}

.chapter-item {
  background: #f8fafc;
  border-radius: 16px;
  padding: 20px 24px;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 20px;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

html.dark .chapter-item {
  background: #2d3748;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);
}

.chapter-item:hover {
  transform: translateX(8px);
  border-color: rgba(102, 126, 234, 0.3);
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.2);
}

.chapter-item.active {
  border-color: #66ead4;
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.3);
}

html.dark .chapter-item.active {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.2) 0%, rgba(118, 75, 162, 0.2) 100%);
}

.chapter-number {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: #3b82f6;
  color: #fff;
  font-size: 20px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.chapter-content {
  flex: 1;
}

.chapter-name {
  font-size: 18px;
  font-weight: 700;
  color: #1a202c;
  margin-bottom: 6px;
}

html.dark .chapter-name {
  color: #f7fafc;
}

.chapter-desc {
  font-size: 14px;
  color: #718096;
  margin-bottom: 8px;
}

html.dark .chapter-desc {
  color: #a0aec0;
}

.chapter-count {
  font-size: 13px;
  color: #667eea;
  font-weight: 600;
}

.chapter-articles {
  background: #f8fafc;
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

html.dark .chapter-articles {
  background: #2d3748;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);
}

.chapter-title {
  font-size: 24px;
  font-weight: 700;
  color: #1a202c;
  margin: 0 0 24px;
  padding-bottom: 16px;
  border-bottom: 3px solid #e2e8f0;
}

html.dark .chapter-title {
  color: #f7fafc;
  border-bottom-color: #4a5568;
}

.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.chapter-article-card {
  background: linear-gradient(135deg, #f8fafc 0%, #fff 100%);
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

html.dark .chapter-article-card {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
}

.chapter-article-card:hover {
  transform: translateY(-4px);
  border-color: rgba(102, 126, 234, 0.3);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.2);
}

.article-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: #3b82f6;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.article-content {
  flex: 1;
  min-width: 0;
}

.article-title {
  font-size: 16px;
  font-weight: 700;
  color: #1a202c;
  margin: 0 0 8px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

html.dark .article-title {
  color: #f7fafc;
}

.article-summary {
  font-size: 13px;
  color: #718096;
  margin: 0 0 12px;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

html.dark .article-summary {
  color: #a0aec0;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 12px;
  color: #a0aec0;
}

.article-views {
  display: flex;
  align-items: center;
  gap: 4px;
}

.article-date {
  color: #a0aec0;
}

.loading-container {
  padding: 40px 0;
}

.skeleton-header {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 32px;
  background: linear-gradient(135deg, #f8fafc 0%, #f0f4f8 100%);
  border-radius: 20px;
  margin-bottom: 30px;
}

.skeleton-nav {
  background: #f8fafc;
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 30px;
}

.skeleton-articles {
  background: #f8fafc;
  border-radius: 16px;
  padding: 32px;
}

@media (max-width: 768px) {
  .series-page {
    padding: 20px 0;
  }

  .series-header {
    flex-direction: column;
    text-align: center;
    padding: 24px;
  }

  .series-icon {
    width: 64px;
    height: 64px;
    font-size: 36px;
  }

  .series-title {
    font-size: 24px;
  }

  .chapter-item {
    padding: 16px;
  }

  .chapter-number {
    width: 40px;
    height: 40px;
    font-size: 16px;
  }

  .chapter-articles {
    padding: 20px;
  }

  .article-grid {
    grid-template-columns: 1fr;
  }
}
</style>
