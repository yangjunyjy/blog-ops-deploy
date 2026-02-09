<template>
  <div class="user-statistics">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户统计</span>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            @change="handleDateChange"
          />
        </div>
      </template>

      <!-- 用户统计卡片 -->
      <el-row :gutter="20" class="stats-cards">
        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-item">
              <div class="stat-value">{{ stats.totalUsers }}</div>
              <div class="stat-label">总用户数</div>
              <div class="stat-icon total-icon">
                <el-icon><User /></el-icon>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-item">
              <div class="stat-value">{{ stats.activeUsers }}</div>
              <div class="stat-label">活跃用户</div>
              <div class="stat-icon active-icon">
                <el-icon><UserFilled /></el-icon>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-item">
              <div class="stat-value">{{ stats.newUsers }}</div>
              <div class="stat-label">今日新增</div>
              <div class="stat-icon new-icon">
                <el-icon><Plus /></el-icon>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-item">
              <div class="stat-value">{{ stats.onlineUsers }}</div>
              <div class="stat-label">在线用户</div>
              <div class="stat-icon online-icon">
                <el-icon><Connection /></el-icon>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 用户增长趋势 -->
      <el-row :gutter="20" class="charts-row">
        <el-col :span="16">
          <el-card>
            <template #header>用户增长趋势</template>
            <div ref="growthChartRef" style="height: 350px"></div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card>
            <template #header>用户分布</template>
            <div ref="distChartRef" style="height: 350px"></div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 用户活跃度 -->
      <el-row :gutter="20" class="activity-row">
        <el-col :span="12">
          <el-card>
            <template #header>用户活跃度</template>
            <div ref="activityChartRef" style="height: 300px"></div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card>
            <template #header>访问时段分布</template>
            <div ref="timeChartRef" style="height: 300px"></div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 最新用户 -->
      <el-card class="recent-users-card">
        <template #header>最新注册用户</template>
        <el-table :data="recentUsers" style="width: 100%">
          <el-table-column type="index" label="#" width="60" />
          <el-table-column label="用户" min-width="200">
            <template #default="{ row }">
              <div class="user-cell">
                <el-avatar :size="32" :src="row.avatar" />
                <div class="user-info">
                  <div class="username">{{ row.nickname }}</div>
                  <div class="email">{{ row.email }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="role" label="角色" width="100">
            <template #default="{ row }">
              <el-tag :type="row.role === 1 ? 'danger' : 'primary'">
                {{ row.role === 1 ? '管理员' : '普通用户' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? '正常' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="注册时间" width="180" />
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small">查看</el-button>
              <el-button type="danger" link size="small">编辑</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </el-card>
  </div>
</template>

<script setup>
defineOptions({
  name: 'User'
})

import { ref, onMounted, onBeforeUnmount } from 'vue'
import * as echarts from 'echarts'
import { User, UserFilled, Plus, Connection } from '@element-plus/icons-vue'
import { getUserList } from '@/api'

const dateRange = ref([])
const stats = ref({
  totalUsers: 0,
  activeUsers: 0,
  newUsers: 0,
  onlineUsers: 0
})
const recentUsers = ref([])

let growthChart = null
let distChart = null
let activityChart = null
let timeChart = null

const growthChartRef = ref(null)
const distChartRef = ref(null)
const activityChartRef = ref(null)
const timeChartRef = ref(null)

const loadData = async () => {
  try {
    const res = await getUserList({ page: 1, size: 10 })
    recentUsers.value = res.data?.list || []

    // 模拟统计数据
    stats.value = {
      totalUsers: 1245,
      activeUsers: 856,
      newUsers: 23,
      onlineUsers: 156
    }

    initCharts()
  } catch (error) {
    console.error('加载数据失败:', error)
    // 即使API失败也初始化图表
    initCharts()
  }
}

const handleDateChange = () => {
  loadData()
}

const initCharts = () => {
  // 用户增长趋势图
  if (growthChartRef.value) {
    growthChart = echarts.init(growthChartRef.value)
    growthChart.setOption({
      tooltip: { trigger: 'axis' },
      legend: { data: ['新增用户', '活跃用户'] },
      xAxis: {
        type: 'category',
        data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
      },
      yAxis: { type: 'value' },
      series: [
        {
          name: '新增用户',
          type: 'line',
          data: [12, 15, 18, 22, 25, 30, 28],
          smooth: true,
          itemStyle: { color: '#409eff' }
        },
        {
          name: '活跃用户',
          type: 'line',
          data: [80, 95, 110, 125, 140, 155, 150],
          smooth: true,
          itemStyle: { color: '#67c23a' }
        }
      ]
    })
  }

  // 用户分布饼图
  if (distChartRef.value) {
    distChart = echarts.init(distChartRef.value)
    distChart.setOption({
      tooltip: { trigger: 'item' },
      legend: { orient: 'vertical', left: 'left' },
      series: [
        {
          name: '用户类型',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
          label: { show: false },
          emphasis: { label: { show: true, fontSize: 20, fontWeight: 'bold' } },
          data: [
            { value: 124, name: '管理员' },
            { value: 1121, name: '普通用户' }
          ]
        }
      ]
    })
  }

  // 用户活跃度雷达图
  if (activityChartRef.value) {
    activityChart = echarts.init(activityChartRef.value)
    activityChart.setOption({
      tooltip: {},
      radar: {
        indicator: [
          { name: '登录次数', max: 100 },
          { name: '评论数', max: 50 },
          { name: '点赞数', max: 100 },
          { name: '收藏数', max: 50 },
          { name: '分享数', max: 30 }
        ]
      },
      series: [
        {
          name: '用户活跃度',
          type: 'radar',
          data: [
            {
              value: [80, 45, 70, 35, 20],
              name: '平均值'
            },
            {
              value: [95, 48, 85, 40, 25],
              name: '活跃用户'
            }
          ]
        }
      ]
    })
  }

  // 访问时段分布图
  if (timeChartRef.value) {
    timeChart = echarts.init(timeChartRef.value)
    timeChart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: {
        type: 'category',
        data: ['0-4', '4-8', '8-12', '12-16', '16-20', '20-24']
      },
      yAxis: { type: 'value' },
      series: [
        {
          name: '访问次数',
          type: 'bar',
          data: [45, 30, 120, 150, 180, 100],
          itemStyle: { color: '#e6a23c' }
        }
      ]
    })
  }
}

const handleResize = () => {
  growthChart?.resize()
  distChart?.resize()
  activityChart?.resize()
  timeChart?.resize()
}

onMounted(() => {
  loadData()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  growthChart?.dispose()
  distChart?.dispose()
  activityChart?.dispose()
  timeChart?.dispose()
})
</script>

<style scoped lang="scss">
.user-statistics {
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
      color: #fff;

      &.total-icon {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      }

      &.active-icon {
        background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
      }

      &.new-icon {
        background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
      }

      &.online-icon {
        background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
      }
    }
  }

  .charts-row,
  .activity-row {
    margin-top: 20px;
  }

  .recent-users-card {
    margin-top: 20px;
  }

  .user-cell {
    display: flex;
    align-items: center;
    gap: 10px;

    .user-info {
      .username {
        font-size: 14px;
        font-weight: 500;
        color: #303133;
      }

      .email {
        font-size: 12px;
        color: #909399;
        margin-top: 2px;
      }
    }
  }
}
</style>
