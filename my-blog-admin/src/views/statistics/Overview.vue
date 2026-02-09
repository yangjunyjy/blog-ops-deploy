<template>
  <div class="statistics-overview">
    <el-row :gutter="20">
      <el-col :xs="24" :sm="12" :lg="6" v-for="(stat, index) in statCards" :key="index">
        <el-card shadow="never" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" :style="{ background: stat.color }">
              <el-icon :size="32">
                <component :is="stat.icon" />
              </el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stat.value }}</div>
              <div class="stat-label">{{ stat.label }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :xs="24" :lg="16">
        <el-card shadow="never">
          <template #header>
            <span>访问趋势</span>
          </template>
          <div ref="chartRef" style="height: 350px"></div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :lg="8">
        <el-card shadow="never">
          <template #header>
            <span>内容分布</span>
          </template>
          <div ref="pieChartRef" style="height: 350px"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
defineOptions({
  name: 'Overview'
})

import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { mockApi } from '../../mock'
import * as echarts from 'echarts'
import { Document, View, Star, ChatDotRound } from '@element-plus/icons-vue'

const stats = ref({})
const chartRef = ref(null)
const pieChartRef = ref(null)
let chart = null
let pieChart = null

const statCards = computed(() => [
  {
    label: '文章总数',
    value: stats.value.overview?.articles || 0,
    icon: Document,
    color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  },
  {
    label: '总浏览量',
    value: stats.value.overview?.views || 0,
    icon: View,
    color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)'
  },
  {
    label: '点赞数',
    value: stats.value.overview?.likes || 0,
    icon: Star,
    color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)'
  },
  {
    label: '评论数',
    value: stats.value.overview?.comments || 0,
    icon: ChatDotRound,
    color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)'
  }
])

const loadStats = async () => {
  try {
    const res = await mockApi.getDashboardStats()
    stats.value = res.data
    initCharts()
  } catch (error) {
    console.error('加载统计失败', error)
  }
}

const initCharts = () => {
  // 折线图
  if (chartRef.value) {
    chart = echarts.init(chartRef.value)
    const option = {
      tooltip: {
        trigger: 'axis'
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
          data: stats.value.trend?.views || [],
          itemStyle: { color: '#409eff' },
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
          data: stats.value.trend?.articles || [],
          itemStyle: { color: '#67c23a' }
        }
      ]
    }
    chart.setOption(option)
  }

  // 饼图
  if (pieChartRef.value) {
    pieChart = echarts.init(pieChartRef.value)
    const option = {
      tooltip: {
        trigger: 'item'
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
            { value: stats.value.content?.categories || 0, name: '分类', itemStyle: { color: '#5470c6' } },
            { value: stats.value.content?.tags || 0, name: '标签', itemStyle: { color: '#91cc75' } },
            { value: stats.value.content?.series || 0, name: '系列', itemStyle: { color: '#fac858' } }
          ]
        }
      ]
    }
    pieChart.setOption(option)
  }
}

onMounted(() => {
  loadStats()
  
  window.addEventListener('resize', () => {
    chart?.resize()
    pieChart?.resize()
  })
})

onBeforeUnmount(() => {
  chart?.dispose()
  pieChart?.dispose()
})
</script>

<style scoped lang="scss">
.statistics-overview {
  .stat-card {
    margin-bottom: 20px;
    
    .stat-content {
      display: flex;
      align-items: center;
      gap: 20px;
      
      .stat-icon {
        width: 64px;
        height: 64px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
      }
      
      .stat-info {
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
}
</style>
