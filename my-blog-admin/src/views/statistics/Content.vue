<template>
  <div class="content-statistics">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>内容统计</span>
          <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期"
            end-placeholder="结束日期" @change="handleDateChange" />
        </div>
      </template>

      <!-- 文章统计卡片 -->
      <el-row :gutter="20" class="stats-cards">
        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-item">
              <div class="stat-value">{{ stats.total }}</div>
              <div class="stat-label">文章总数</div>
              <div class="stat-icon article-icon">
                <el-icon>
                  <Document />
                </el-icon>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-item">
              <div class="stat-value">{{ stats.published }}</div>
              <div class="stat-label">已发布</div>
              <div class="stat-icon published-icon">
                <el-icon>
                  <CircleCheck />
                </el-icon>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-item">
              <div class="stat-value">{{ stats.draft }}</div>
              <div class="stat-label">草稿</div>
              <div class="stat-icon draft-icon">
                <el-icon>
                  <EditPen />
                </el-icon>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-item">
              <div class="stat-value">{{ stats.totalViews }}</div>
              <div class="stat-label">总浏览量</div>
              <div class="stat-icon views-icon">
                <el-icon>
                  <View />
                </el-icon>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 文章发布趋势 -->
      <el-row :gutter="20" class="charts-row">
        <el-col :span="12">
          <el-card>
            <template #header>文章发布趋势</template>
            <div ref="publishChartRef" style="height: 300px"></div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card>
            <template #header>浏览量趋势</template>
            <div ref="viewChartRef" style="height: 300px"></div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 热门文章 -->
      <el-card class="hot-articles-card">
        <template #header>热门文章 TOP 10</template>
        <el-table :data="hotArticles" style="width: 100%">
          <el-table-column type="index" label="排名" width="60" align="center" />
          <el-table-column prop="title" label="文章标题" min-width="300" />
          <el-table-column prop="views" label="浏览量" width="120" sortable>
            <template #default="{ row }">
              <el-icon>
                <View />
              </el-icon>
              {{ row.views }}
            </template>
          </el-table-column>
          <el-table-column prop="likes" label="点赞数" width="120" sortable>
            <template #default="{ row }">
              <el-icon>
                <Star />
              </el-icon>
              {{ row.likes }}
            </template>
          </el-table-column>
          <el-table-column prop="commentCount" label="评论数" width="120" sortable>
            <template #default="{ row }">
              <el-icon>
                <ChatDotRound />
              </el-icon>
              {{ row.commentCount }}
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <!-- 分类和标签统计 -->
      <el-row :gutter="20" class="category-tag-row">
        <el-col :span="12">
          <el-card>
            <template #header>分类文章数</template>
            <div ref="categoryChartRef" style="height: 300px"></div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card>
            <template #header>标签使用率</template>
            <div ref="tagChartRef" style="height: 300px"></div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
defineOptions({
  name: 'Content'
})

import { ref, onMounted, onBeforeUnmount } from 'vue'
import * as echarts from 'echarts'
import { Document, CircleCheck, EditPen, View, Star, ChatDotRound } from '@element-plus/icons-vue'
import { getArticleStats, getHotArticles, getCategoryList, getTagList } from '@/api'


const dateRange = ref([])
const stats = ref({
  total: 0,
  published: 0,
  draft: 0,
  totalViews: 0,
  trend: []
})
const hotArticles = ref([])
const categories = ref([])
const tags = ref([])

let publishChart = null
let viewChart = null
let categoryChart = null
let tagChart = null

const publishChartRef = ref(null)
const viewChartRef = ref(null)
const categoryChartRef = ref(null)
const tagChartRef = ref(null)

const loadData = async () => {
  try {
    const [statsRes, hotRes, catRes, tagRes] = await Promise.all([
      getArticleStats(),
      getHotArticles({ limit: 10 }),
      getCategoryList({ page: 1, size: 100 }),
      getTagList({ page: 1, size: 100 })
    ])

    stats.value = statsRes.data || stats.value
    hotArticles.value = hotRes.data?.list || []
    categories.value = catRes.data?.list || []
    tags.value = tagRes.data?.list || []

    initCharts()
  } catch (error) {
    console.error('加载数据失败:', error)
  }
}

const handleDateChange = () => {
  loadData()
}

const initCharts = () => {
  // 文章发布趋势图
  if (publishChartRef.value) {
    publishChart = echarts.init(publishChartRef.value)
    publishChart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: {
        type: 'category',
        data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
      },
      yAxis: { type: 'value' },
      series: [
        {
          name: '发布数',
          type: 'line',
          data: stats.value.trend,
          smooth: true,
          areaStyle: { opacity: 0.3 }
        }
      ]
    })
  }

  // 浏览量趋势图
  if (viewChartRef.value) {
    viewChart = echarts.init(viewChartRef.value)
    viewChart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: {
        type: 'category',
        data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
      },
      yAxis: { type: 'value' },
      series: [
        {
          name: '浏览量',
          type: 'bar',
          data: Array.from({ length: 7 }, () => Math.floor(Math.random() * 500) + 100),
          itemStyle: { color: '#67c23a' }
        }
      ]
    })
  }

  // 分类文章数饼图
  if (categoryChartRef.value) {
    categoryChart = echarts.init(categoryChartRef.value)
    categoryChart.setOption({
      tooltip: { trigger: 'item' },
      legend: { orient: 'vertical', left: 'left' },
      series: [
        {
          name: '分类文章数',
          type: 'pie',
          radius: '50%',
          data: categories.value.map(cat => ({
            value: cat.articleCount,
            name: cat.name
          }))
        }
      ]
    })
  }

  // 标签使用率柱状图
  if (tagChartRef.value) {
    tagChart = echarts.init(tagChartRef.value)
    tagChart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: {
        type: 'category',
        data: tags.value.slice(0, 10).map(tag => tag.name)
      },
      yAxis: { type: 'value' },
      series: [
        {
          name: '文章数',
          type: 'bar',
          data: tags.value.slice(0, 10).map(tag => tag.articleCount),
          itemStyle: { color: '#409eff' }
        }
      ]
    })
  }
}

const handleResize = () => {
  publishChart?.resize()
  viewChart?.resize()
  categoryChart?.resize()
  tagChart?.resize()
}

onMounted(() => {
  loadData()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  publishChart?.dispose()
  viewChart?.dispose()
  categoryChart?.dispose()
  tagChart?.dispose()
})
</script>

<style scoped lang="scss">
.content-statistics {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .stats-cards {
    margin-bottom: 20px;
  }

  .stat-item {
    position: relative;
    padding: 20px;

    .stat-value {
      font-size: 32px;
      font-weight: bold;
      color: #303133;
      margin-bottom: 8px;
    }

    .stat-label {
      font-size: 14px;
      color: #909399;
    }

    .stat-icon {
      position: absolute;
      right: 20px;
      top: 50%;
      transform: translateY(-50%);
      width: 50px;
      height: 50px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 24px;

      &.article-icon {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: #fff;
      }

      &.published-icon {
        background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
        color: #fff;
      }

      &.draft-icon {
        background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
        color: #fff;
      }

      &.views-icon {
        background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
        color: #fff;
      }
    }
  }

  .charts-row,
  .category-tag-row {
    margin-top: 20px;
  }

  .hot-articles-card {
    margin-top: 20px;
  }
}
</style>
