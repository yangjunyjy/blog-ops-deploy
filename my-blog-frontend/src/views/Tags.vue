<script setup>
import { ref, onMounted, computed } from 'vue'
import { getTags } from '@/api'
import Pagination from '@/components/Pagination.vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const tags = ref([])
const loading = ref(true)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

onMounted(async () => {
  loading.value = true
  try {
    const res = await getTags({
      page: currentPage.value,
      pageSize: pageSize.value
    })
    tags.value = res.data.items || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('加载标签失败:', error)
  } finally {
    loading.value = false
  }
})

const handlePageChange = (page) => {
  currentPage.value = page
  loading.value = true
  getTags({
    page: currentPage.value,
    pageSize: pageSize.value
  }).then(res => {
    tags.value = res.data.items || []
    total.value = res.data.total || 0
  }).catch(error => {
    console.error('加载标签失败:', error)
  }).finally(() => {
    loading.value = false
  })
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<template>
  <div class="tags-page" :class="{ 'dark-mode': themeStore.isDark }">
    <div class="container">
      <el-skeleton :loading="loading" animated>
        <template #template>
          <div class="skeleton-wrapper">
            <el-skeleton-item variant="h1" style="width: 200px; margin-bottom: 30px;" />
            <div class="tags-grid">
              <div v-for="i in 12" :key="i" class="skeleton-tag">
                <el-skeleton-item variant="text" style="width: 100%;" />
              </div>
            </div>
          </div>
        </template>
        <template #default>
          <div class="page-content">
            <div class="page-title">
              <span class="title-emoji">🏷</span>
              <h1>文章标签</h1>
              <p>用标签发现感兴趣的内容</p>
            </div>

            <div class="tags-grid" v-if="tags.length">
              <router-link
                v-for="tag in tags"
                :key="tag.id"
                :to="`/tag/${tag.id}`"
                class="tag-card"
              >
                <span class="tag-name">{{ tag.name }}</span>
                <span class="tag-count">{{ tag.articleCount || 0 }}</span>
              </router-link>
            </div>
            <div v-else class="empty-state">
              <span class="empty-emoji">🌸</span>
              <p>暂无标签</p>
            </div>

            <Pagination class="page-now"
              v-if="total > 0"
              :current-page="currentPage"
              :total="total"
              :page-size="pageSize"
              :total-pages="totalPages"
              @change="handlePageChange"
            />
          </div>
        </template>
      </el-skeleton>
    </div>
  </div>
</template>

<style scoped>
.tags-page {
  min-height: calc(100vh - 140px);
  padding: 40px 0;
  background: linear-gradient(135deg, #e8f5e9 0%, #e3f2fd 100%);
  transition: background 0.3s ease;
}

.tags-page.dark-mode {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
}

.page-content {
  max-width: 900px;
  margin: 0 auto;
  .page-now{
    margin-top: 60px;
  }
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

.tags-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  justify-content: center;
  padding: 0 20px;
}

.tag-card {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 14px 24px;
  background: #f8fafc;
  border-radius: 30px;
  box-shadow: 0 3px 12px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  text-decoration: none;
  border: 2px solid transparent;
  position: relative;
  overflow: hidden;
}

.dark-mode .tag-card {
  background: rgba(45, 55, 72, 0.9);
  box-shadow: 0 3px 12px rgba(0, 0, 0, 0.3);
}

.tag-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.4), transparent);
  transition: left 0.5s;
}

.tag-card:hover::before {
  left: 100%;
}

.tag-card:hover {
  transform: translateY(-4px) scale(1.05);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.2);
  border-color: #c4dfff;
}

.tag-name {
  font-size: 16px;
  font-weight: 600;
  color: #4a4a4a;
  transition: color 0.3s;
}

.dark-mode .tag-name {
  color: #e2e8f0;
}

.tag-card:hover .tag-name {
  color: #667eea;
}

.tag-count {
  font-size: 13px;
  color: #999;
  background: linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%);
  padding: 4px 12px;
  border-radius: 12px;
  min-width: 28px;
  text-align: center;
  font-weight: 600;
}

.dark-mode .tag-count {
  color: #a0aec0;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.3) 0%, rgba(118, 75, 162, 0.3) 100%);
}

.tag-card:hover .tag-count {
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

.skeleton-tag {
  background: #f8fafc;
  border-radius: 30px;
  padding: 14px 24px;
  min-width: 120px;
}

.dark-mode .skeleton-tag {
  background: rgba(45, 55, 72, 0.9);
}

@media (max-width: 768px) {
  .tags-page {
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

  .tags-grid {
    gap: 10px;
    padding: 0 15px;
  }

  .tag-card {
    padding: 12px 18px;
  }

  .tag-name {
    font-size: 14px;
  }

  .tag-count {
    font-size: 12px;
    padding: 3px 10px;
    min-width: 24px;
  }
}
</style>
