<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)">
            <el-icon :size="32"><Document /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.overview?.articles || 0 }}</div>
            <div class="stat-label">文章总数</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%)">
            <el-icon :size="32"><View /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.overview?.views || 0 }}</div>
            <div class="stat-label">总浏览量</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)">
            <el-icon :size="32"><Star /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.overview?.likes || 0 }}</div>
            <div class="stat-label">点赞数</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)">
            <el-icon :size="32"><ChatDotRound /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.overview?.comments || 0 }}</div>
            <div class="stat-label">评论数</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="charts-row">
      <el-col :xs="24" :lg="16">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>访问趋势</span>
            </div>
          </template>
          <div ref="viewsChartRef" style="height: 350px"></div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :lg="8">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>内容分布</span>
            </div>
          </template>
          <div ref="contentChartRef" style="height: 350px"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最新文章 -->
    <el-card shadow="never" class="latest-articles">
      <template #header>
        <div class="card-header">
          <span>最新文章</span>
          <el-button type="primary" link @click="$router.push('/articles')">查看全部</el-button>
        </div>
      </template>
      <el-table :data="latestArticles" style="width: 100%">
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="category.name" label="分类" width="120" />
        <el-table-column prop="views" label="浏览" width="80" align="center" />
        <el-table-column prop="likes" label="点赞" width="80" align="center" />
        <el-table-column prop="commentCount" label="评论" width="80" align="center" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
defineOptions({
  name: 'Dashboard'
})

import { ref, onMounted, onBeforeUnmount } from 'vue'
import { mockApi } from '../mock'
import * as echarts from 'echarts'
import { Document, View, Star, ChatDotRound } from '@element-plus/icons-vue'

const stats = ref({})
const latestArticles = ref([])
const viewsChartRef = ref(null)
const contentChartRef = ref(null)
let viewsChart = null
let contentChart = null

const loadStats = async () => {
  try {
    const res = await mockApi.getDashboardStats()
    stats.value = res.data
    latestArticles.value = res.data?.latest || []
  } catch (error) {
    console.error('加载统计失败', error)
  }
}

const initCharts = () => {
  // 访问趋势图
  if (viewsChartRef.value) {
    viewsChart = echarts.init(viewsChartRef.value)
    const option = {
      tooltip: {
        trigger: 'axis'
      },
      legend: {
        data: ['浏览量', '文章数']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '浏览量',
          type: 'line',
          smooth: true,
          data: [120, 132, 101, 134, 90, 230, 210],
          itemStyle: {
            color: '#409eff'
          },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
              { offset: 1, color: 'rgba(64, 158, 255, 0.1)' }
            ])
          }
        },
        {
          name: '文章数',
          type: 'line',
          smooth: true,
          data: [5, 8, 6, 9, 7, 10, 8],
          itemStyle: {
            color: '#67c23a'
          }
        }
      ]
    }
    viewsChart.setOption(option)
  }

  // 内容分布图
  if (contentChartRef.value) {
    contentChart = echarts.init(contentChartRef.value)
    const option = {
      tooltip: {
        trigger: 'item'
      },
      legend: {
        orient: 'vertical',
        left: 'left'
      },
      series: [
        {
          name: '内容分布',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: true,
            formatter: '{b}: {c}'
          },
          data: [
            { value: stats.value.content?.categories || 5, name: '分类' },
            { value: stats.value.content?.tags || 10, name: '标签' },
            { value: stats.value.content?.series || 3, name: '系列' }
          ]
        }
      ]
    }
    contentChart.setOption(option)
  }
}

onMounted(() => {
  loadStats()
  initCharts()

  window.addEventListener('resize', () => {
    viewsChart?.resize()
    contentChart?.resize()
  })
})

onBeforeUnmount(() => {
  viewsChart?.dispose()
  contentChart?.dispose()
})
</script>

<style scoped lang="scss">
.dashboard {
  .stats-row {
    margin-bottom: 20px;
    
    .stat-card {
      background: #fff;
      border-radius: 8px;
      padding: 20px;
      display: flex;
      align-items: center;
      gap: 20px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
      
      .stat-icon {
        width: 64px;
        height: 64px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
      }
      
      .stat-content {
        flex: 1;
        
        .stat-value {
          font-size: 28px;
          font-weight: bold;
          color: #303133;
          margin-bottom: 8px;
        }
        
        .stat-label {
          font-size: 14px;
          color: #909399;
        }
      }
    }
  }
  
  .charts-row {
    margin-bottom: 20px;
  }
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 600;
  }
  
  .latest-articles {
    :deep(.el-card__body) {
      padding: 20px 0;
    }
  }
}
</style>
