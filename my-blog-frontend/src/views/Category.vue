<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getCategoryDetail } from '@/api'
import ArticleCard from '@/components/ArticleCard.vue'
import Pagination from '@/components/Pagination.vue'
import { FolderOpened } from '@element-plus/icons-vue'
import { useThemeStore } from '@/stores/theme'

const route = useRoute()
const themeStore = useThemeStore()

const category = ref(null)
const articles = ref([])
const loading = ref(true)
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const loadCategoryDetail = async () => {
  loading.value = true
  category.value = null
  articles.value = []
  try {
    const res = await getCategoryDetail({
      id: route.params.id,
      page: currentPage.value,
      pageSize: pageSize.value
    })
    console.log("分类内容",res.data);
    category.value = res.data?.category
    articles.value = res.data?.articles?.list || []
    total.value = res.data?.articles?.total || 0
  } catch (error) {
    console.error('加载分类详情失败:', error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadCategoryDetail()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  loadCategoryDetail()
})

watch(() => route.params.id, () => {
  currentPage.value = 1
  loadCategoryDetail()
  window.scrollTo({ top: 0, behavior: 'smooth' })
})
</script>

<template>
  <div class="category-page" :class="{ 'dark-mode': themeStore.isDark }">
    <div class="container">
      <el-skeleton :loading="loading" animated>
        <template #template>
          <div>
            <el-skeleton-item variant="h1" style="width: 40%; margin-bottom: 30px;" />
            <div class="article-grid">
              <div v-for="i in 6" :key="i" style="margin-bottom: 20px;">
                <el-skeleton-item variant="image" style="width: 100%; height: 200px" />
                <div style="padding: 20px;">
                  <el-skeleton-item variant="h3" style="width: 70%" />
                  <el-skeleton-item variant="text" style="margin-top: 10px" />
                </div>
              </div>
            </div>
          </div>
        </template>
        <template #default>
          <div v-if="category">
            <div class="category-header">
              <el-icon><FolderOpened /></el-icon>
              <h1>{{ category.name }}</h1>
              <p>{{ category.description || `查看 ${category.name} 分类下的所有文章` }}</p>
              <span class="article-count">共{{ total }} 篇文章</span>
            </div>

            <div class="article-grid" v-if="articles.length">
              <ArticleCard
                v-for="article in articles"
                :key="article.id"
                :article="article"
              />
            </div>
            <el-empty v-else description="该分类下暂无文章"/>
          </div>
          <el-empty v-else description="分类不存在" />
        </template>
      </el-skeleton>

      <Pagination
        v-if="category"
        :current-page="currentPage"
        :total="total"
        :page-size="pageSize"
        :total-pages="totalPages"
        @change="handlePageChange"
      />
    </div>
  </div>
</template>

<style scoped>
.category-page {
  min-height: calc(100vh - 140px);
  padding: 50px 0;
  transition: background 0.3s ease;
}

.dark-mode .category-page {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
}

.category-header {
  text-align: center;
  margin-bottom: 30px;
  padding: 25px 30px;
  background: linear-gradient(135deg, #38bdf8 0%, #0ea5e9 50%, #0284c7 100%);
  border-radius: 16px;
  color: #fff;
  box-shadow: 0 6px 20px rgba(14, 165, 233, 0.25);
  position: relative;
  overflow: hidden;
}

.category-header .el-icon {
  font-size: 32px;
  margin-bottom: 12px;
  position: relative;
  z-index: 1;
}

.category-header h1 {
  font-size: 24px;
  margin: 0 0 12px;
  font-weight: 800;
  letter-spacing: -0.5px;
  position: relative;
  z-index: 1;
}

.category-header p {
  font-size: 14px;
  opacity: 0.95;
  margin: 0 0 16px;
  position: relative;
  z-index: 1;
}

.article-count {
  display: inline-block;
  padding: 6px 16px;
  background: rgba(255, 255, 255, 0.25);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
  position: relative;
  z-index: 1;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

@media (max-width: 768px) {
  .category-page {
    padding: 30px 0;
  }

  .category-header {
    padding: 40px 20px;
    border-radius: 16px;
  }

  .category-header .el-icon {
    font-size: 44px;
  }

  .category-header h1 {
    font-size: 28px;
  }

  .article-grid {
    grid-template-columns: 1fr;
  }
}
</style>
