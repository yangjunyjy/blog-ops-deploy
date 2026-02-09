<script setup>
import { ref, onMounted } from 'vue'
import { getLatestArticles } from '@/api'
import ArticleCard from '@/components/ArticleCard.vue'
import Sidebar from '@/components/Sidebar.vue'
import LatestArticlesCarousel from '@/components/LatestArticlesCarousel.vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const articles = ref([])
const loading = ref(true)

onMounted(async () => {
  loading.value = true
  articles.value = []
  try {
    const res = await getLatestArticles()
    console.log("获取的最新文章数量",res.data);
    articles.value = (res.data || []).slice(4)
  } catch (error) {
    console.error('加载文章失败:', error)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="home-page">
    <div class="hero-section">
      <div class="hero-content">
        <div class="hero-text">
          <h1 class="hero-title">博客空间</h1>
          <p class="hero-subtitle">分享技术?· 记录生活 · 传递价值</p>
          <p class="hero-motto">" Code like you mean it "</p>
        </div>

        <div class="hero-image">
          <div class="avatar-ring"></div>
          <div class="avatar-glow"></div>
        </div>
      </div>
    </div>

    <div class="container">
      <div class="main-content">
        <div class="content-left">
          <LatestArticlesCarousel />

          <div class="section-header">
            <h2>最新文章</h2>
            <div class="header-decoration">
              <div class="dot"></div>
              <div class="dot"></div>
              <div class="dot"></div>
            </div>
          </div>

          <el-skeleton :loading="loading" animated :count="6">
            <template #template>
              <div style="margin-bottom: 24px;">
                <el-skeleton-item variant="image" style="width: 100%; height: 220px" />
                <div style="padding: 24px;">
                  <el-skeleton-item variant="h3" style="width: 70%" />
                  <el-skeleton-item variant="text" style="margin-top: 12px" />
                  <el-skeleton-item variant="text" style="width: 60%" />
                </div>
              </div>
            </template>
            <template #default>
              <div class="article-grid" v-if="articles.length">
                <ArticleCard
                  v-for="article in articles"
                  :key="article.id"
                  :article="article"
                />
              </div>
              <el-empty v-else description="暂无文章">
                <template #description>
                  <div class="empty-content">
                    <p>还没有文章，快去创作吧</p>
                  </div>
                </template>
              </el-empty>
            </template>
          </el-skeleton>
        </div>

        <div class="content-right">
          <Sidebar />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-page {
  min-height: calc(100vh - 140px);
  background: #f5f7fa;
  transition: background 0.3s ease;
}

html.dark .home-page {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
}

.hero-section {
  background: url('/images/background.jpeg') center center / cover no-repeat;
  padding: 100px 0 80px;
  position: relative;
  overflow: hidden;
  margin-bottom: 60px;
}

html.dark .hero-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  pointer-events: none;
}

.hero-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  pointer-events: none;
}

.hero-content {
  position: relative;
  z-index: 1;
}

.hero-text {
  text-align: center;
  color: #fff;
  position: relative;
  z-index: 2;
}

.hero-title {
  font-size: 56px;
  font-weight: 800;
  margin-bottom: 24px;
  text-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  letter-spacing: -1px;
}

.hero-subtitle {
  font-size: 24px;
  margin-bottom: 16px;
  opacity: 0.95;
  font-weight: 300;
  letter-spacing: 2px;
}

.hero-motto {
  font-size: 18px;
  font-style: italic;
  opacity: 0.9;
  font-family: 'Georgia', serif;
  position: relative;
  display: inline-block;
  padding: 0 20px;
}

.hero-motto::before,
.hero-motto::after {
  content: '"';
  font-size: 24px;
  color: rgba(255, 255, 255, 0.5);
}

.hero-motto::before {
  margin-right: 8px;
}

.hero-motto::after {
  margin-left: 8px;
}

.hero-image {
  position: relative;
  width: 200px;
  height: 200px;
  margin: 50px auto 0;
}

.avatar-ring {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 50%;
  border: 4px solid rgba(255, 255, 255, 0.3);
}

.avatar-glow {
  position: absolute;
  top: -20px;
  left: -20px;
  right: -20px;
  bottom: -20px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.3) 0%, transparent 70%);
}

.main-content {
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: 40px;
  padding-bottom: 60px;
}

.content-left {
  min-width: 0;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 35px;
  padding-bottom: 18px;
  border-bottom: 3px solid #e2e8f0;
  position: relative;
}

html.dark .section-header {
  border-bottom-color: #4a5568;
}

.section-header::after {
  content: '';
  position: absolute;
  bottom: -3px;
  left: 0;
  width: 80px;
  height: 3px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 3px;
}

.section-header h2 {
  font-size: 32px;
  color: #1a202c;
  margin: 0;
  font-weight: 800;
  letter-spacing: -0.5px;
  transition: color 0.3s ease;
}

html.dark .section-header h2 {
  color: #f7fafc;
}

.header-decoration {
  display: flex;
  gap: 8px;
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #3b82f6;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; transform: scale(0.8); }
  50% { opacity: 1; transform: scale(1); }
}

.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

.empty-content {
  text-align: center;
}

.empty-content p {
  font-size: 16px;
  color: #718096;
  margin: 0;
  transition: color 0.3s ease;
}

html.dark .empty-content p {
  color: #a0aec0;
}

@media (max-width: 992px) {
  .main-content {
    grid-template-columns: 1fr;
  }

  .content-right {
    display: none;
  }
}

@media (max-width: 768px) {
  .hero-section {
    padding: 70px 0 60px;
  }

  .hero-title {
    font-size: 36px;
  }

  .hero-subtitle {
    font-size: 18px;
  }

  .hero-motto {
    font-size: 15px;
  }

  .hero-image {
    width: 160px;
    height: 160px;
  }

  .section-header h2 {
    font-size: 24px;
  }

  .article-grid {
    grid-template-columns: 1fr;
  }
}
</style>
