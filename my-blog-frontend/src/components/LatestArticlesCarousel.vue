<script setup>
import { ref, onMounted, computed, onBeforeUnmount } from 'vue'
import { ArrowRight, ArrowLeft, Clock } from '@element-plus/icons-vue'
import { getLatestArticles } from '@/api'
import { useRouter } from 'vue-router'

const router = useRouter()
const articles = ref([])
const currentIndex = ref(0)
const loading = ref(false)
const autoPlayInterval = ref(null)
const AUTO_PLAY_DELAY = 3000

const canSlideLeft = computed(() => currentIndex.value > 0)
const canSlideRight = computed(() => currentIndex.value < articles.value.length - 1)

const loadLatestArticles = async () => {
  loading.value = true
  try {
    const res = await getLatestArticles()
    articles.value = (res.data || []).slice(0, 5)
  } catch (error) {
    console.error('加载最新文章失败:', error)
  } finally {
    loading.value = false
  }
}

const slideLeft = () => {
  if (canSlideLeft.value) {
    currentIndex.value--
  } else {
    currentIndex.value = articles.value.length - 1
  }
}

const slideRight = () => {
  if (canSlideRight.value) {
    currentIndex.value++
  } else {
    currentIndex.value = 0
  }
}

const goToArticle = (id) => {
  const url = router.resolve({
    name: 'ArticleDetail',
    params: { id: id }
  }).href
  window.open(url, '_blank')
}

const startAutoPlay = () => {
  stopAutoPlay()
  autoPlayInterval.value = setInterval(() => {
    slideRight()
  }, AUTO_PLAY_DELAY)
}

const stopAutoPlay = () => {
  if (autoPlayInterval.value) {
    clearInterval(autoPlayInterval.value)
    autoPlayInterval.value = null
  }
}

const handleMouseEnter = () => {
  stopAutoPlay()
}

const handleMouseLeave = () => {
  startAutoPlay()
}

const goToSlide = (index) => {
  currentIndex.value = index
}

onMounted(() => {
  loadLatestArticles()
  startAutoPlay()
})

onBeforeUnmount(() => {
  stopAutoPlay()
})
</script>

<template>
  <div
    class="latest-articles-carousel"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <div class="carousel-header">
      <h3 class="carousel-title">
        <el-icon><Clock /></el-icon>
        最新动态
      </h3>
      <div class="carousel-nav">
        <button
          class="nav-btn"
          :class="{ disabled: !canSlideLeft }"
          @click="slideLeft"
        >
          <el-icon><ArrowLeft /></el-icon>
        </button>
        <button
          class="nav-btn"
          :class="{ disabled: !canSlideRight }"
          @click="slideRight"
        >
          <el-icon><ArrowRight /></el-icon>
        </button>
      </div>
    </div>

    <div class="carousel-container">
      <div
        class="carousel-track"
        :style="{ transform: `translateX(-${currentIndex * 100}%)` }"
      >
        <div
          v-for="(article, index) in articles"
          :key="article.id"
          class="carousel-slide"
          @click="goToArticle(article.id)"
        >
          <div class="slide-content">
            <div class="slide-cover" v-if="article.cover">
              <img :src="article.cover" :alt="article.title" loading="lazy" />
              <div class="slide-overlay">
                <span class="slide-badge">{{ index + 1 }}</span>
              </div>
            </div>
            <div class="slide-info">
              <h4 class="slide-title">{{ article.title }}</h4>
              <p class="slide-summary">{{ article.summary }}</p>
              <div class="slide-meta">
                <span class="slide-views">
                  <el-icon><Clock /></el-icon>
                  {{ article.views || 0 }} 阅读
                </span>
                <span class="slide-date">
                  {{ article.createdAt }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="carousel-dots">
      <button
        v-for="(article, index) in articles"
        :key="article.id"
        class="dot"
        :class="{ active: index === currentIndex }"
        @click="goToSlide(index)"
      ></button>
    </div>

    <div v-if="loading" class="carousel-loading">
      <el-skeleton animated>
        <template #template>
          <div class="loading-skeleton">
            <el-skeleton-item variant="image" style="width: 100%; height: 300px;" />
            <el-skeleton-item variant="h3" style="width: 80%; margin: 20px 0;" />
            <el-skeleton-item variant="text" style="width: 100%;" />
          </div>
        </template>
      </el-skeleton>
    </div>

    <el-empty v-if="!loading && !articles.length" description="暂无最新文章" />
  </div>
</template>

<style scoped>
.latest-articles-carousel {
  background: #f8fafc;
  border-radius: 20px;
  padding: 24px;
  margin-bottom: 30px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.5);
  position: relative;
  overflow: hidden;
}

html.dark .latest-articles-carousel {
  background: #2d3748;
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.carousel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.carousel-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: 700;
  color: #1a202c;
  margin: 0;
}

html.dark .carousel-title {
  color: #f7fafc;
}

.carousel-nav {
  display: flex;
  gap: 8px;
}

.nav-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 2px solid #e2e8f0;
  background: #fff;
  color: #4a5568;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
}

html.dark .nav-btn {
  background: #1a202c;
  border-color: #4a5568;
  color: #cbd5e0;
}

.nav-btn:hover:not(.disabled) {
  border-color: #667eea;
  color: #667eea;
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.nav-btn.disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.carousel-container {
  position: relative;
  overflow: hidden;
  border-radius: 16px;
}

.carousel-track {
  display: flex;
  transition: transform 0.5s ease;
}

.carousel-slide {
  min-width: 100%;
  cursor: pointer;
}

.slide-content {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 16px;
  overflow: hidden;
  transition: all 0.3s;
}

html.dark .slide-content {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
}

.slide-content:hover {
  transform: scale(1.02);
  box-shadow: 0 8px 30px rgba(102, 126, 234, 0.2);
}

.slide-cover {
  position: relative;
  width: 100%;
  height: 300px;
  overflow: hidden;
}

.slide-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.slide-content:hover .slide-cover img {
  transform: scale(1.1);
}

.slide-overlay {
  position: absolute;
  top: 16px;
  right: 16px;
}

.slide-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: #3b82f6;
  color: #fff;
  font-size: 20px;
  font-weight: 800;
  border-radius: 50%;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.slide-info {
  padding: 24px;
}

.slide-title {
  font-size: 22px;
  font-weight: 700;
  color: #1a202c;
  margin: 0 0 12px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

html.dark .slide-title {
  color: #f7fafc;
}

.slide-summary {
  font-size: 15px;
  color: #718096;
  margin: 0 0 16px;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

html.dark .slide-summary {
  color: #a0aec0;
}

.slide-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
  color: #a0aec0;
}

.slide-views {
  display: flex;
  align-items: center;
  gap: 6px;
}

.slide-date {
  color: #a0aec0;
}

.carousel-dots {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-top: 20px;
}

.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #e2e8f0;
  border: none;
  cursor: pointer;
  transition: all 0.3s;
}

html.dark .dot {
  background: #4a5568;
}

.dot:hover {
  background: #667eea;
}

.dot.active {
  width: 32px;
  border-radius: 6px;
  background: #3b82f6;
}

.carousel-loading {
  min-height: 400px;
}

.loading-skeleton {
  padding: 0;
}

@media (max-width: 768px) {
  .latest-articles-carousel {
    padding: 16px;
  }

  .carousel-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .slide-cover {
    height: 200px;
  }

  .slide-title {
    font-size: 18px;
  }

  .slide-info {
    padding: 16px;
  }
}
</style>
