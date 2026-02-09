<script setup>
import { ref, onMounted } from 'vue'
import { getHotArticles, getCategories, getTags } from '@/api'
import { TrendCharts, Folder, PriceTag } from '@element-plus/icons-vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const hotArticles = ref([])
const categories = ref([])
const tags = ref([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const [hotRes, catRes, tagRes] = await Promise.all([
      getHotArticles(),
      getCategories(),
      getTags()
    ])
    hotArticles.value = hotRes.data || []
    categories.value = catRes.data.items || []
    tags.value = tagRes.data.items || []
  } catch (error) {
    console.error('加载侧边栏数据失败', error)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="sidebar" :class="{ 'dark-mode': themeStore.isDark }">
    <!-- -->
    <div class="sidebar-section">
      <div class="section-title">
        <el-icon><TrendCharts /></el-icon>
        热门文章
      </div>
      <div class="hot-articles" v-loading="loading">
        <router-link
          v-for="article in hotArticles.slice(0, 5)"
          :key="article.id"
          :to="`/article/${article.id}`"
          class="hot-article-item"
        >
          {{ article.title }}
        </router-link>
      </div>
    </div>

    <!-- 分类 -->
    <div class="sidebar-section">
      <div class="section-title">
        <el-icon><Folder /></el-icon>
        分类
      </div>
      <div class="category-list" v-loading="loading">
        <router-link
          v-for="category in categories"
          :key="category.id"
          :to="`/category/${category.id}`"
          class="category-item"
        >
          <span>{{ category.name }}</span>
          <span class="count">{{ category.articleCount || 0 }}</span>
        </router-link>
      </div>
    </div>

    <!-- 标签 -->
    <div class="sidebar-section">
      <div class="section-title">
        <el-icon><PriceTag /></el-icon>
        标签
      </div>
      <div class="tag-cloud" v-loading="loading">
        <router-link
          v-for="tag in tags.slice(0, 15)"
          :key="tag.name"
          :to="`/tag/${tag.name}`"
          class="tag-item"
        >
          {{ tag.name }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.sidebar {
  position: sticky;
  top: 100px;
}

.sidebar-section {
  background: #f8fafc;
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.5);
  transition: all 0.3s;
}

.dark-mode .sidebar-section {
  background: rgba(45, 55, 72, 0.95);
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
}

.sidebar-section:hover {
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.08);
}

.dark-mode .sidebar-section:hover {
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.4);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 700;
  color: #1a202c;
  margin-bottom: 20px;
  padding-bottom: 14px;
  border-bottom: 2px solid #f0f0f0;
  position: relative;
  letter-spacing: -0.3px;
  transition: color 0.3s ease;
}

.dark-mode .section-title {
  color: #f7fafc;
  border-bottom-color: #4a5568;
}

.section-title::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 40px;
  height: 2px;
  background: #3b82f6;
  border-radius: 2px;
}

.section-title .el-icon {
  color: #667eea;
  font-size: 20px;
}

.hot-articles {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.hot-article-item {
  color: #4a5568;
  font-size: 14px;
  line-height: 1.6;
  padding: 10px;
  border-radius: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: all 0.3s;
  font-weight: 500;
}

.dark-mode .hot-article-item {
  color: #e2e8f0;
}

.hot-article-item:hover {
  color: #667eea;
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
  transform: translateX(4px);
}

.dark-mode .hot-article-item:hover {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.2) 0%, rgba(118, 75, 162, 0.2) 100%);
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.category-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 10px;
  color: #4a5568;
  font-size: 14px;
  transition: all 0.3s;
  font-weight: 500;
  border: 1px solid rgba(255, 255, 255, 0.5);
}

.dark-mode .category-item {
  background: linear-gradient(135deg, rgba(74, 85, 104, 0.3) 0%, rgba(45, 55, 72, 0.3) 100%);
  color: #e2e8f0;
  border-color: rgba(255, 255, 255, 0.1);
}

.category-item:hover {
  background: #3b82f6;
  color: #fff;
  transform: translateX(4px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.25);
}

.count {
  background: #667eea;
  color: #fff;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  transition: all 0.3s;
}

.dark-mode .count {
  background: #3b82f6;
}

.category-item:hover .count {
  background: rgba(255, 255, 255, 0.3);
  color: #fff;
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.tag-item {
  padding: 8px 16px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 20px;
  color: #4a5568;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.3s;
  border: 1px solid rgba(255, 255, 255, 0.5);
}

.dark-mode .tag-item {
  background: linear-gradient(135deg, rgba(74, 85, 104, 0.3) 0%, rgba(45, 55, 72, 0.3) 100%);
  color: #e2e8f0;
  border-color: rgba(255, 255, 255, 0.1);
}

.tag-item:hover {
  background: #3b82f6;
  color: #fff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.25);
}
</style>
