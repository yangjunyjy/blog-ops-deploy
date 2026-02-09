<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getSeries, getSeriesDetail } from '@/api'

const route = useRoute()
const router = useRouter()

const seriesList = ref([])
const activeSeriesId = ref(null)
const currentSeries = ref(null)
const loading = ref(true)

const loadSeriesList = async () => {
  loading.value = true
  try {
    const res = await getSeries()
    seriesList.value = res.data.list || []
    console.log(res.data);
    // 默认选中第一个系�?
    if (seriesList.value.length > 0) {
      // 如果URL中有id参数，则使用该id
      const seriesId = route.params.id
      if (seriesId) {
        const found = seriesList.value.find(s => s.id === parseInt(seriesId))
        if (found) {
          activeSeriesId.value = found.id
          await loadSeriesDetail(found.id)
        }
      } else {
        activeSeriesId.value = seriesList.value[0].id
        await loadSeriesDetail(seriesList.value[0].id)
      }
    }
  } catch (error) {
    console.error('加载专栏列表失败:', error)
  } finally {
    loading.value = false
  }
}

const loadSeriesDetail = async (seriesId) => {
  try {
    const res = await getSeriesDetail(seriesId)
    currentSeries.value = res.data
    console.log("系列详情",res.data);
    
  } catch (error) {
    console.error('加载专栏详情失败:', error)
  }
}

const handleSeriesClick = async (seriesId) => {
  activeSeriesId.value = seriesId
  await loadSeriesDetail(seriesId)
  router.replace(`/series/${seriesId}`)
}

const handleSubChapterClick = ({ seriesId, subchapterId }) => {
  console.log(seriesId,subchapterId);
  const url = router.resolve({
    name: 'ChapterDetail',
    params: {
      seriesId,
      chapterId: subchapterId
    }
  }).href
  window.open(url, '_blank')
}

onMounted(() => {
  loadSeriesList()
})
</script>

<template>
  <div class="series-page" v-if="!loading">
    <!-- 子导航栏 -->
    <div class="series-sub-nav">
      <div class="container">
        <div class="nav-items">
          <div
            v-for="series in seriesList"
            :key="series.id"
            class="nav-item"
            :class="{ active: series.id === activeSeriesId }"
            @click="handleSeriesClick(series.id)"
          >
            <!-- <span class="nav-icon">{{ series.name.charAt(0) }}</span> -->
            <span class="nav-name">{{ series.name }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 专栏内容 -->
    <div class="container" v-if="currentSeries">
      <div class="series-content">
        <!-- 专栏头部 -->
        <div class="series-header">
          <div class="series-icon">{{ currentSeries.name.charAt(0) }}</div>
          <div class="series-info">
            <h1 class="series-title">{{ currentSeries.name }}</h1>
            <p class="series-description">{{ currentSeries.description }}</p>
            <div class="series-stats">
              <span class="stat">
                <span class="stat-number">{{ currentSeries.sections?.length || 0 }}</span>
                <span class="stat-label">章节标题</span>
              </span>
            </div>
          </div>
        </div>

        <!-- 章节标题和子章节列表 -->
        <div class="sections-container">
          <div
            v-for="section in currentSeries.sections"
            :key="section.id"
            class="section-block"
          >
            <h2 class="section-heading">{{ section.name }}</h2>
            <p class="section-description">{{ section.description || '暂无描述' }}</p>

            <div class="subchapter-cards">
              <div
                v-for="subchapter in section.subchapters"
                :key="subchapter.id"
                class="subchapter-card"
                @click="handleSubChapterClick({ seriesId: activeSeriesId, subchapterId: subchapter.id })"
              >
                <div class="subchapter-card-icon">📂</div>
                <div class="subchapter-card-content">
                  <h3 class="subchapter-card-title">{{ subchapter.name }}</h3>
                  <p class="subchapter-card-summary">{{ subchapter.description || '暂无描述' }}</p>
                  <div class="subchapter-card-meta">
                    <span class="subchapter-card-articles">📄 {{ subchapter.articleIds?.length || 0 }} 篇文章</span>
                  </div>
                </div>
              </div>
            </div>
            <el-empty v-if="!section.subchapters?.length" description="该章节暂无目录" />
          </div>
        </div>
      </div>
    </div>
  </div>

  <div class="loading-container" v-else>
    <el-skeleton animated>
      <template #template>
        <div class="skeleton-nav">
          <el-skeleton-item variant="rect" style="height: 60px; margin-right: 16px;" />
          <el-skeleton-item variant="rect" style="height: 60px; margin-right: 16px;" />
          <el-skeleton-item variant="rect" style="height: 60px;" />
        </div>
        <div class="skeleton-header">
          <el-skeleton-item variant="circle" style="width: 80px; height: 80px;" />
          <el-skeleton-item variant="h3" style="width: 60%;" />
          <el-skeleton-item variant="text" style="width: 80%;" />
        </div>
        <div class="skeleton-content">
          <el-skeleton-item variant="h1" style="width: 40%; margin-bottom: 24px;" />
          <div class="skeleton-subchapters">
            <el-skeleton-item variant="rect" style="height: 120px; margin-bottom: 24px;" />
            <el-skeleton-item variant="rect" style="height: 120px; margin-bottom: 24px;" />
            <el-skeleton-item variant="rect" style="height: 120px;" />
          </div>
        </div>
      </template>
    </el-skeleton>
  </div>
</template>

<style scoped>
.series-page {
  min-height: calc(100vh - 140px);
  padding-top: 20px;
}

/* 子导航栏 */
.series-sub-nav {
  position: sticky;
  top: 80px;
  z-index: 90;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid #e2e8f0;
  margin-bottom: 30px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

html.dark .series-sub-nav {
  background: rgba(26, 32, 44, 0.95);
  border-bottom-color: #4a5568;
}

.nav-items {
  display: flex;
  gap: 12px;
  padding: 12px 0;
  overflow-x: auto;
}

.nav-items::-webkit-scrollbar {
  height: 4px;
}

.nav-items::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 2px;
}

.nav-items::-webkit-scrollbar-track {
  background: transparent;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s;
  white-space: nowrap;
  background: #f8fafc;
  border: 2px solid transparent;
}

html.dark .nav-item {
  background: #2d3748;
}

.nav-item:hover {
  transform: translateY(-2px);
  border-color: rgba(75, 131, 155, 0.3);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
}

.nav-item.active {
  background: skyblue;
  color: #fff;
  border-color: transparent;
}

.nav-icon {
  font-size: 20px;
}

.nav-name {
  font-size: 15px;
  font-weight: 600;
}

/* 专栏内容 */
.series-content {
  padding-bottom: 60px;
}

.series-header {
  background: linear-gradient(135deg, #f8fafc 0%, #f0f4f8 100%);
  border-radius: 12px;
  padding: 20px;
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
  background: skyblue;
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
  /* flex-direction: column; */
  gap: 4px;
}

.stat-number {
  font-size: 24px;
  font-weight: 700;
  color: #667eea;
}

.stat-label {
  font-size: 12px;
  margin-top: 12px;
  color: #a0aec0;
  text-transform: uppercase;
  letter-spacing: 0.5px;

}

/* 章节标题和子章节列表 */
.sections-container {
  display: flex;
  flex-direction: column;
  gap: 50px;
}

.section-block {
  margin-bottom: 20px;
}

.section-heading {
  font-size: 28px;
  font-weight: 700;
  color: #1a202c;
  margin: 0 0 12px;
  padding-bottom: 16px;
  border-bottom: 3px solid #e2e8f0;
}

html.dark .section-heading {
  color: #f7fafc;
  border-bottom-color: #4a5568;
}

.section-description {
  font-size: 15px;
  color: #718096;
  margin: 0 0 24px;
  line-height: 1.6;
}

html.dark .section-description {
  color: #a0aec0;
}

.subchapter-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

.subchapter-card {
  background: linear-gradient(135deg, #f8fafc 0%, #f0f4f8 100%);
  border-radius: 16px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
  display: flex;
  align-items: flex-start;
  gap: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

html.dark .subchapter-card {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.3);
}

.subchapter-card:hover {
  transform: translateY(-6px);
  border-color: rgba(102, 126, 234, 0.4);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.25);
}

.subchapter-card-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  background: skyblue;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.subchapter-card-content {
  flex: 1;
  min-width: 0;
}

.subchapter-card-title {
  font-size: 18px;
  font-weight: 700;
  color: #1a202c;
  margin: 0 0 10px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

html.dark .subchapter-card-title {
  color: #f7fafc;
}

.subchapter-card-summary {
  font-size: 14px;
  color: #718096;
  margin: 0 0 16px;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

html.dark .subchapter-card-summary {
  color: #a0aec0;
}

.subchapter-card-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 13px;
  color: #a0aec0;
}

.subchapter-card-articles {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  color: #667eea;
}

html.dark .subchapter-card-articles {
  color: #a5b4fc;
}

.loading-container {
  padding: 40px 0;
}

.skeleton-nav {
  background: #f8fafc;
  border-radius: 12px;
  padding: 12px;
  margin-bottom: 30px;
  display: flex;
  gap: 12px;
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

.skeleton-content {
  margin-bottom: 30px;
}

.skeleton-subchapters {
  background: #f8fafc;
  border-radius: 16px;
  padding: 32px;
}

@media (max-width: 768px) {
  .series-page {
    padding-top: 10px;
  }

  .series-sub-nav {
    top: 70px;
  }

  .nav-items {
    padding: 8px 0;
  }

  .nav-item {
    padding: 8px 16px;
  }

  .series-icon {
    width: 64px;
    height: 64px;
    font-size: 36px;
  }

  .series-title {
    font-size: 24px;
  }

  .section-heading {
    font-size: 24px;
  }

  .subchapter-cards {
    grid-template-columns: 1fr;
  }
}
</style>
