<template>
  <div class="search-page">
    <div class="container">
        <div class="search-header">
          <div class="search-box">
            <el-input
              v-model="keyword"
              placeholder="搜索文章..."
              size="large"
              @keyup.enter="performSearch"
              clearable
            >
              <template #suffix>
                <el-icon @click="performSearch"><Search /></el-icon>
              </template>
            </el-input>
          </div>
          <div class="search-info" v-if="keyword">
            <p>搜索 "{{ keyword }}" 的结果</p>
            <span class="result-count">找到 {{ total }} 篇文章</span>
          </div>
        </div>

      <el-skeleton :loading="loading" animated :count="6">
        <template #template>
          <div class="article-grid">
            <div v-for="i in 6" :key="i" style="margin-bottom: 20px;">
              <el-skeleton-item variant="image" style="width: 100%; height: 200px" />
              <div style="padding: 20px;">
                <el-skeleton-item variant="h3" style="width: 70%" />
                <el-skeleton-item variant="text" style="margin-top: 10px" />
              </div>
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
          <div v-else-if="keyword" class="no-result">
            <el-empty description="未找到相关文章">
              <el-button type="primary" @click="keyword = ''; performSearch()">清除搜索</el-button>
            </el-empty>
          </div>
          <div v-else class="search-tips">
            <el-icon class="search-icon"><Search /></el-icon>
            <h3>输入关键词搜索文章</h3>
            <p>支持搜索文章标题、摘要、标签等内容</p>
          </div>
        </template>
      </el-skeleton>

      <Pagination
        :current-page="currentPage"
        :total="total"
        :page-size="pageSize"
        :total-pages="totalPages"
        @change="handlePageChange"
      />
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { searchArticles } from '@/api'
import ArticleCard from '@/components/ArticleCard.vue'
import Pagination from '@/components/Pagination.vue'
import { Search } from '@element-plus/icons-vue'

const route = useRoute()

const keyword = ref('')
const articles = ref([])
const loading = ref(true)
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const performSearch = async () => {
  if (!keyword.value.trim()) {
    articles.value = []
    loading.value = false
    total.value = 0
    return
  }

  loading.value = true
  articles.value = []
  try {
    const res = await searchArticles({
      keyword: keyword.value,
      page: currentPage.value,
      pageSize: pageSize.value
    })
    articles.value = res.data?.items || []
    total.value = res.data?.total || 0
    currentPage.value = res.data?.page || 1
  } catch (error) {
    console.error('搜索失败:', error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  currentPage.value = page
  performSearch()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  keyword.value = route.query.keyword || ''
  currentPage.value = Number(route.query.page) || 1
  if (keyword.value) {
    performSearch()
  }
})

watch(() => route.query.keyword, (newVal) => {
  keyword.value = newVal || ''
  currentPage.value = 1
  performSearch()
  window.scrollTo({ top: 0, behavior: 'smooth' })
})
</script>


<style scoped>
.search-page {
  min-height: calc(100vh - 140px);
  padding: 40px 0;
}

.search-header {
  margin-bottom: 40px;
  text-align: center;
}

.search-box {
  max-width: 600px;
  margin: 0 auto 20px;
}

.search-box .el-input {
  font-size: 16px;
}

.search-info {
  color: #606266;
}

.search-info p {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 10px;
}

.result-count {
  color: #909399;
  font-size: 14px;
}

.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.search-tips {
  text-align: center;
  padding: 100px 20px;
  color: #909399;
}

.search-icon {
  font-size: 64px;
  margin-bottom: 20px;
  color: #c0c4cc;
}

.search-tips h3 {
  font-size: 24px;
  margin: 0 0 10px;
  color: #606266;
}

.search-tips p {
  font-size: 14px;
  margin: 0;
}

@media (max-width: 768px) {
  .search-page {
    padding: 20px 0;
  }

  .article-grid {
    grid-template-columns: 1fr;
  }

  .search-tips {
    padding: 60px 20px;
  }
}
</style>
