<script setup>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getArticles, getCategories } from '@/api'
import ArticleCard from '@/components/ArticleCard.vue'
import Pagination from '@/components/Pagination.vue'

const route = useRoute()

const articles = ref([])
const categories = ref([])
const loading = ref(true)
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const selectedCategory = ref('')
const selectedTag = ref('')
const selectValue = ref("")
const showBackToTop = ref(false)

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const getSelectValue = () => {
  categories.value.forEach(item => {
    if (item.id == selectedCategory.value) {
        selectValue.value = item.name
    }
  })
}

const loadArticles = async () => {
  loading.value = true
  articles.value = []
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      status:1
    }

    if (selectedCategory.value && selectedCategory.value !== '') {
      params.categoryId = selectedCategory.value
    }

    if (selectedTag.value && selectedTag.value !== '') {
      params.tag = selectedTag.value
    }

    const res = await getArticles(params)
    // console.log("获取的文章数量",res);
    articles.value = res.data?.items || []
    total.value = res.data?.total || 0
  } catch (error) {
    console.error('加载文章失败:', error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadArticles()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const handleCategoryChange = (categoryId) => {
  selectedCategory.value = categoryId
  currentPage.value = 1
  
  // 更新显示的分类名称
  if (categoryId === '') {
    selectValue.value = ''
  } else {
    const category = categories.value.find(c => c.id == categoryId)
    selectValue.value = category ? category.name : ''
  }
  
  loadArticles()
}

onMounted(async () => {
  const catRes = await getCategories()
  categories.value = catRes.data.items || []
  await loadArticles()
  
  // 如果有路由参数中的 categoryId，需要更新 selectValue
  if (route.query.categoryId) {
    const category = categories.value.find(c => c.id == route.query.categoryId)
    if (category) {
      selectValue.value = category.name
    }
  }
  
  // 监听滚动事件
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

const handleScroll = () => {
  showBackToTop.value = window.scrollY > 300
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 监听路由变化，从其他页面进入时重新加载数据
watch(() => route.path, (newPath) => {
  if (newPath === '/articles') {
    currentPage.value = 1
    selectedCategory.value = ''
    selectedTag.value = ''
    selectValue.value = ''
    loadArticles()
  }
})

watch(() => route.query.categoryId, (newVal) => {
  if (newVal) {
    const categoryId = Number(newVal)
    selectedCategory.value = categoryId
    
    // 更新显示的分类名称
    const category = categories.value.find(c => c.id == categoryId)
    selectValue.value = category ? category.name : ''
    
    currentPage.value = 1
    loadArticles()
  }
})
</script>

<template>
  <div class="articles-page">
    <div class="container">
      <div class="main-content">
        <div class="content-area">
          <div class="section-header">
            <h2>文章列表</h2>
            <div class="filter-bar">
              <el-select
                v-model="selectedCategory"
                placeholder="选择分类"
                clearable
                @change="handleCategoryChange"
                style="min-width: 100px;"
              >
                <el-option label="全部分类" :value="''" />
                <el-option
                  v-for="category in categories"
                  :key="category.id"
                  :label="category.name"
                  :value="category.id"
                />
              </el-select>
            </div>
          </div>

          <el-skeleton :loading="loading" animated :count="6">
            <template #template>
              <div style="margin-bottom: 20px;">
                <el-skeleton-item variant="image" style="width: 100%; height: 200px" />
                <div style="padding: 20px;">
                  <el-skeleton-item variant="h3" style="width: 70%" />
                  <el-skeleton-item variant="text" style="margin-top: 10px" />
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
              <el-empty v-else description="暂无文章" />
            </template>
          </el-skeleton>

          <Pagination
            :current-page="currentPage"
            :total="total"
            :page-size="pageSize"
            :total-pages="totalPages"
            @change="handlePageChange"
          />

          <button
            v-if="showBackToTop"
            @click="scrollToTop"
            class="back-to-top"
          >
            ↑
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.articles-page {
  min-height: calc(100vh - 140px);
  padding: 50px 0;
}
html.dark .articles-page {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
}

.content-area {
  width: 100%;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 30px;
}

.section-header h2 {
  font-size: 24px;
  color: #1a202c;
  margin: 0;
  font-weight: 600;
  letter-spacing: 0;
}

html.dark .section-header h2 {
  color: #f7fafc;
}

.filter-bar {
  display: flex;
  gap: 15px;
}

.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.article-grid :deep(.el-skeleton) {
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

@media (max-width: 992px) {
  .container {
    padding: 0 20px;
  }
}

@media (max-width: 768px) {
  .articles-page {
    padding: 30px 0;
  }

  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 20px;
  }

  .section-header h2 {
    font-size: 24px;
  }

  .filter-bar {
    width: 100%;
  }

  .filter-bar .el-select {
    width: 100%;
  }

  .article-grid {
    grid-template-columns: 1fr;
  }
}

.back-to-top {
  position: fixed;
  top: 50%;
  right: 40px;
  transform: translateY(-50%);
  width: 50px;
  height: 50px;
  background: linear-gradient(135deg, #38bdf8 0%, #0ea5e9 50%, #0284c7 100%);
  color: #fff;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  font-size: 24px;
  font-weight: bold;
  box-shadow: 0 4px 15px rgba(14, 165, 233, 0.4);
  transition: all 0.3s ease;
  z-index: 999;
  display: flex;
  align-items: center;
  justify-content: center;
}

.back-to-top:hover {
  transform: translateY(-50%) scale(1.1);
  box-shadow: 0 6px 20px rgba(14, 165, 233, 0.5);
  background: linear-gradient(135deg, #0ea5e9 0%, #0284c7 50%, #0369a1 100%);
}

@media (max-width: 768px) {
  .back-to-top {
    right: 20px;
    width: 45px;
    height: 45px;
    font-size: 20px;
  }
}
</style>
