<script setup>
import { ref, onMounted, onActivated, computed } from 'vue'
import { getCategories } from '@/api'
import Pagination from '@/components/Pagination.vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const categories = ref([])
const loading = ref(true)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const loadCategories = async () => {
  loading.value = true
  categories.value = []
  try {
    const res = await getCategories({
      page: currentPage.value,
      pageSize: pageSize.value
    })
    console.log("分类数据",res.data)
    categories.value = res.data.items || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('加载分类失败:', error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadCategories()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  loadCategories()
})

onActivated(() => {
  loadCategories()
})
</script>

<template>
  <div class="categories-page" :class="{ 'dark-mode': themeStore.isDark }">
    <div class="container">
      <el-skeleton :loading="loading" animated>
        <template #template>
          <div class="skeleton-wrapper">
            <el-skeleton-item variant="h1" style="width: 200px; margin-bottom: 30px;" />
            <div class="categories-grid">
              <div v-for="i in 6" :key="i" class="skeleton-card">
                <el-skeleton-item variant="circle" style="width: 60px; height: 60px; margin-bottom: 15px;" />
                <el-skeleton-item variant="h3" style="width: 60%;" />
              </div>
            </div>
          </div>
        </template>
        <template #default>
          <div class="page-content">
            <div class="page-title">
              <h1>文章分类</h1>
              <p>探索不同领域的精彩内容</p>
            </div>

            <div class="categories-grid" v-if="categories.length">
              <router-link
                v-for="category in categories"
                :key="category.id"
                :to="`/category/${category.id}`"
                class="category-card"
              >
                <div class="category-icon">
                  <span class="emoji">📁</span>
                </div>
                <h3 class="category-name">{{ category.name }}</h3>
                <span class="category-count">{{ category.articleCount || 0 }} </span>
              </router-link>
            </div>
            <div v-else class="empty-state">
              <span class="empty-emoji">🍃</span>
              <p>暂无分类</p>
            </div>

            <div class="page-now">
              <Pagination
              v-if="total > 0"
              :current-page="currentPage"
              :total="total"
              :page-size="pageSize"
              :total-pages="totalPages"
              @change="handlePageChange"
            />
            </div>
          </div>
        </template>
      </el-skeleton>
    </div>
  </div>
</template>

<style scoped>
.categories-page {
  min-height: calc(100vh - 140px);
  padding: 40px 0;
  background: linear-gradient(135deg, #fef7f0 0%, #f0e7ff 100%);
  transition: background 0.3s ease;
}

.categories-page.dark-mode {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
}


.page-content {
  max-width: 900px;
  margin: 0 auto;
}

.page-now {
  margin-top: 60px !important;
}

.page-title {
  text-align: center;
  margin-bottom: 50px;
}

.title-emoji {
  font-size: 48px;
  display: block;
  margin-bottom: 15px;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.page-title h1 {
  font-size: 36px;
  margin: 0 0 10px;
  color: #4a4a4a;
  font-weight: 700;
  transition: color 0.3s ease;
}

.dark-mode .page-title h1 {
  color: #f7fafc;
}

.page-title p {
  font-size: 16px;
  color: #999;
  margin: 0;
  transition: color 0.3s ease;
}

.dark-mode .page-title p {
  color: #a0aec0;
}

.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 24px;
  padding: 0 20px;
}

.category-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32px 24px;
  background: #f8fafc;
  border-radius: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  text-decoration: none;
  border: 2px solid transparent;
}

.dark-mode .category-card {
  background: rgba(45, 55, 72, 0.9);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.category-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 8px 30px rgba(102, 126, 234, 0.2);
  border-color: #d4c4ff;
}

.category-icon {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  background: linear-gradient(135deg, #ffeaa7 0%, #dfe6e9 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 18px;
  transition: all 0.3s;
}

.category-card:hover .category-icon {
  background: linear-gradient(135deg, #a29bfe 0%, #fdcb6e 100%);
  transform: rotate(10deg) scale(1.1);
}

.emoji {
  font-size: 32px;
}

.category-name {
  font-size: 18px;
  font-weight: 600;
  color: #4a4a4a;
  margin: 0 0 12px;
  text-align: center;
  transition: color 0.3s;
}

.dark-mode .category-name {
  color: #e2e8f0;
}

.category-card:hover .category-name {
  color: #667eea;
}

.category-count {
  font-size: 14px;
  color: #999;
  background: #f5f5f5;
  padding: 6px 16px;
  border-radius: 20px;
  transition: all 0.3s;
}

.dark-mode .category-count {
  color: #a0aec0;
  background: rgba(0, 0, 0, 0.2);
}

.category-card:hover .category-count {
  background: linear-gradient(135deg, #a29bfe 0%, #fdcb6e 100%);
  color: #fff;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-emoji {
  font-size: 64px;
  display: block;
  margin-bottom: 20px;
  animation: bounce 2s ease-in-out infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-15px); }
}

.empty-state p {
  font-size: 18px;
  color: #999;
  margin: 0;
  transition: color 0.3s ease;
}

.dark-mode .empty-state p {
  color: #a0aec0;
}

.skeleton-wrapper {
  max-width: 900px;
  margin: 0 auto;
  padding: 0 20px;
}

.skeleton-card {
  background: #f8fafc;
  border-radius: 20px;
  padding: 32px 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.dark-mode .skeleton-card {
  background: rgba(45, 55, 72, 0.9);
}

@media (max-width: 768px) {
  .categories-page {
    padding: 30px 0;
  }

  .title-emoji {
    font-size: 40px;
  }

  .page-title h1 {
    font-size: 28px;
  }

  .page-title p {
    font-size: 14px;
  }

  .categories-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 16px;
    padding: 0 15px;
  }

  .category-card {
    padding: 24px 20px;
  }

  .category-icon {
    width: 60px;
    height: 60px;
  }

  .emoji {
    font-size: 28px;
  }

  .category-name {
    font-size: 16px;
  }

  .category-count {
    font-size: 13px;
    padding: 5px 14px;
  }
}
</style>
