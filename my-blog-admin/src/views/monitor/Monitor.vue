<template>
  <div class="monitor-container">
    <!-- 顶部统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card" v-for="stat in stats" :key="stat.key">
        <div class="stat-icon" :style="{ background: stat.iconBg }">
          <el-icon :size="24" :color="stat.iconColor">
            <component :is="stat.icon" />
          </el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-value">{{ stat.value }}</div>
          <div class="stat-trend" :class="stat.trendClass">
            <el-icon :size="12">
              <component :is="stat.trendIcon" />
            </el-icon>
            <span>{{ stat.trend }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 左侧：服务器列表 -->
      <div class="server-panel">
        <div class="panel-header">
          <div class="header-title">
            <el-icon><Monitor /></el-icon>
            <span>服务器列表</span>
          </div>
          <el-button type="primary" size="small" :icon="Plus">添加服务器</el-button>
        </div>
        <div class="server-list">
          <div 
            v-for="server in servers" 
            :key="server.id"
            class="server-item"
            :class="{ active: selectedServer?.id === server.id }"
            @click="selectServer(server)"
          >
            <div class="server-status">
              <span class="status-dot" :class="server.status"></span>
            </div>
            <div class="server-info">
              <div class="server-name">{{ server.name }}</div>
              <div class="server-ip">{{ server.ip }}</div>
            </div>
            <div class="server-metrics">
              <div class="metric-item">
                <el-icon :size="14" color="#67c23a"><Cpu /></el-icon>
                <span>{{ server.cpu }}%</span>
              </div>
              <div class="metric-item">
                <el-icon :size="14" color="#409eff"><DataLine /></el-icon>
                <span>{{ server.memory }}%</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：监控详情 -->
      <div class="detail-panel" v-if="selectedServer">
        <!-- 服务器信息 -->
        <div class="server-info-card">
          <div class="info-header">
            <div class="info-title">
              <el-icon><Monitor /></el-icon>
              <span>{{ selectedServer.name }}</span>
              <el-tag :type="selectedServer.status === 'online' ? 'success' : 'danger'" size="small">
                {{ selectedServer.status === 'online' ? '在线' : '离线' }}
              </el-tag>
            </div>
            <div class="info-actions">
              <el-button size="small" :icon="Position" @click="openTerminal">终端</el-button>
              <el-button size="small" :icon="Switch">刷新</el-button>
            </div>
          </div>
          <div class="info-details">
            <div class="detail-item">
              <span class="detail-label">IP 地址:</span>
              <span class="detail-value">{{ selectedServer.ip }}:{{ selectedServer.port }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">操作系统:</span>
              <span class="detail-value">{{ selectedServer.os }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">运行时间:</span>
              <span class="detail-value">{{ selectedServer.uptime }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">最后更新:</span>
              <span class="detail-value">{{ selectedServer.lastUpdate }}</span>
            </div>
          </div>
        </div>

        <!-- 资源使用率 -->
        <div class="metrics-grid">
          <!-- CPU 使用率 -->
          <div class="metric-card">
            <div class="metric-header">
              <div class="metric-title">
                <el-icon color="#67c23a"><Cpu /></el-icon>
                <span>CPU 使用率</span>
              </div>
              <span class="metric-value">{{ metrics.cpu }}%</span>
            </div>
            <div class="metric-chart">
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: metrics.cpu + '%', background: getMetricColor(metrics.cpu) }"></div>
              </div>
            </div>
            <div class="metric-history">
              <div 
                v-for="(value, index) in metrics.cpuHistory" 
                :key="index"
                class="history-bar"
                :style="{ height: value + '%', background: getMetricColor(value) }"
              ></div>
            </div>
          </div>

          <!-- 内存使用率 -->
          <div class="metric-card">
            <div class="metric-header">
              <div class="metric-title">
                <el-icon color="#409eff"><DataLine /></el-icon>
                <span>内存使用率</span>
              </div>
              <span class="metric-value">{{ metrics.memory }}%</span>
            </div>
            <div class="metric-chart">
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: metrics.memory + '%', background: getMetricColor(metrics.memory) }"></div>
              </div>
            </div>
            <div class="memory-details">
              <div class="memory-info">
                <span class="info-label">已用: {{ metrics.memoryUsed }} GB</span>
                <span class="info-label">可用: {{ metrics.memoryAvailable }} GB</span>
              </div>
              <div class="memory-info">
                <span class="info-label">总计: {{ metrics.memoryTotal }} GB</span>
              </div>
            </div>
          </div>

          <!-- 磁盘使用率 -->
          <div class="metric-card">
            <div class="metric-header">
              <div class="metric-title">
                <el-icon color="#e6a23c"><Coin /></el-icon>
                <span>磁盘使用率</span>
              </div>
              <span class="metric-value">{{ metrics.disk }}%</span>
            </div>
            <div class="disk-list">
              <div v-for="disk in metrics.disks" :key="disk.path" class="disk-item">
                <div class="disk-info">
                  <span class="disk-path">{{ disk.path }}</span>
                  <span class="disk-usage">{{ disk.used }} / {{ disk.total }}</span>
                </div>
                <div class="progress-bar">
                  <div class="progress-fill" :style="{ width: disk.percent + '%', background: getMetricColor(disk.percent) }"></div>
                </div>
              </div>
            </div>
          </div>

          <!-- 网络流量 -->
          <div class="metric-card">
            <div class="metric-header">
              <div class="metric-title">
                <el-icon color="#f56c6c"><Connection /></el-icon>
                <span>网络流量</span>
              </div>
            </div>
            <div class="network-stats">
              <div class="network-item upload">
                <div class="network-icon">
                  <el-icon :size="20"><Upload /></el-icon>
                </div>
                <div class="network-info">
                  <div class="network-label">上传</div>
                  <div class="network-value">{{ metrics.network.upload }}</div>
                </div>
              </div>
              <div class="network-item download">
                <div class="network-icon">
                  <el-icon :size="20"><Download /></el-icon>
                </div>
                <div class="network-info">
                  <div class="network-label">下载</div>
                  <div class="network-value">{{ metrics.network.download }}</div>
                </div>
              </div>
            </div>
            <div class="network-history">
              <div class="history-legend">
                <span class="legend-item"><span class="legend-color upload"></span>上传</span>
                <span class="legend-item"><span class="legend-color download"></span>下载</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 进程列表 -->
        <div class="process-panel">
          <div class="panel-title">
            <el-icon><List /></el-icon>
            <span>系统进程</span>
          </div>
          <div class="process-table">
            <div class="table-header">
              <div class="header-cell">PID</div>
              <div class="header-cell">用户</div>
              <div class="header-cell">进程名</div>
              <div class="header-cell">CPU</div>
              <div class="header-cell">内存</div>
              <div class="header-cell">状态</div>
            </div>
            <div class="table-body">
              <div v-for="process in processes" :key="process.pid" class="table-row">
                <div class="body-cell">{{ process.pid }}</div>
                <div class="body-cell">{{ process.user }}</div>
                <div class="body-cell">{{ process.name }}</div>
                <div class="body-cell">{{ process.cpu }}%</div>
                <div class="body-cell">{{ process.memory }}%</div>
                <div class="body-cell">
                  <el-tag :type="process.status === 'running' ? 'success' : 'warning'" size="small">
                    {{ process.status }}
                  </el-tag>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Monitor, Plus, Position, Switch, Cpu, Coin,
  Connection, Upload, Download, List, ArrowUp,
  ArrowDown, TrendCharts, DataLine
} from '@element-plus/icons-vue'

// 统计数据
const stats = ref([
  {
    key: 'servers',
    label: '服务器总数',
    value: '5',
    trend: '+2',
    trendClass: 'up',
    trendIcon: ArrowUp,
    icon: Monitor,
    iconBg: '#ecf5ff',
    iconColor: '#409eff'
  },
  {
    key: 'online',
    label: '在线服务器',
    value: '4',
    trend: '+1',
    trendClass: 'up',
    trendIcon: ArrowUp,
    icon: Connection,
    iconBg: '#f0f9ff',
    iconColor: '#67c23a'
  },
  {
    key: 'alerts',
    label: '告警数量',
    value: '3',
    trend: '-1',
    trendClass: 'down',
    trendIcon: ArrowDown,
    icon: TrendCharts,
    iconBg: '#fef0f0',
    iconColor: '#f56c6c'
  },
  {
    key: 'uptime',
    label: '平均运行时间',
    value: '99.8%',
    trend: '+0.2%',
    trendClass: 'up',
    trendIcon: ArrowUp,
    icon: Monitor,
    iconBg: '#f4f4f5',
    iconColor: '#909399'
  }
])

// 服务器列表
const servers = ref([
  {
    id: 1,
    name: 'Web服务器-01',
    ip: '192.168.1.100',
    port: 22,
    status: 'online',
    cpu: 35,
    memory: 45,
    os: 'Ubuntu 22.04 LTS',
    uptime: '30天 12小时',
    lastUpdate: '2026-01-29 10:30:00'
  },
  {
    id: 2,
    name: '数据库服务器',
    ip: '192.168.1.101',
    port: 22,
    status: 'online',
    cpu: 68,
    memory: 78,
    os: 'CentOS 7.9',
    uptime: '45天 8小时',
    lastUpdate: '2026-01-29 10:30:00'
  },
  {
    id: 3,
    name: '应用服务器',
    ip: '192.168.1.102',
    port: 22,
    status: 'online',
    cpu: 42,
    memory: 55,
    os: 'Ubuntu 22.04 LTS',
    uptime: '15天 20小时',
    lastUpdate: '2026-01-29 10:29:00'
  },
  {
    id: 4,
    name: '缓存服务器',
    ip: '192.168.1.103',
    port: 22,
    status: 'online',
    cpu: 15,
    memory: 30,
    os: 'Debian 11',
    uptime: '60天 5小时',
    lastUpdate: '2026-01-29 10:30:00'
  },
  {
    id: 5,
    name: '备份服务器',
    ip: '192.168.1.104',
    port: 22,
    status: 'offline',
    cpu: 0,
    memory: 0,
    os: 'CentOS 7.9',
    uptime: '-',
    lastUpdate: '2026-01-28 15:20:00'
  }
])

const selectedServer = ref(servers.value[0])

// 监控指标
const metrics = ref({
  cpu: 35,
  cpuHistory: [28, 32, 35, 38, 35, 30, 28, 32, 35, 38, 35, 30],
  memory: 45,
  memoryTotal: 16,
  memoryUsed: 7.2,
  memoryAvailable: 8.8,
  disk: 65,
  disks: [
    { path: '/', used: '120 GB', total: '200 GB', percent: 60 },
    { path: '/data', used: '320 GB', total: '500 GB', percent: 64 }
  ],
  network: {
    upload: '2.5 MB/s',
    download: '8.3 MB/s'
  }
})

// 进程列表
const processes = ref([
  { pid: 1234, user: 'root', name: 'nginx', cpu: 5.2, memory: 2.1, status: 'running' },
  { pid: 5678, user: 'mysql', name: 'mysqld', cpu: 15.8, memory: 35.2, status: 'running' },
  { pid: 9012, user: 'root', name: 'dockerd', cpu: 8.5, memory: 12.3, status: 'running' },
  { pid: 3456, user: 'www-data', name: 'php-fpm', cpu: 3.2, memory: 4.5, status: 'running' },
  { pid: 7890, user: 'root', name: 'sshd', cpu: 0.5, memory: 1.2, status: 'running' }
])

// 定时器
let refreshTimer = null

onMounted(() => {
  startRefresh()
})

onUnmounted(() => {
  stopRefresh()
})

// 选择服务器
const selectServer = (server) => {
  selectedServer.value = server
  refreshMetrics()
}

// 打开终端
const openTerminal = () => {
  const url = window.location.origin + '/terminal/' + selectedServer.value.id
  window.open(url, '_blank')
}

// 刷新监控指标
const refreshMetrics = () => {
  // 模拟数据更新
  metrics.value.cpu = Math.floor(Math.random() * 60) + 20
  metrics.value.memory = Math.floor(Math.random() * 40) + 40
  metrics.value.memoryUsed = (metrics.value.memoryTotal * metrics.value.memory / 100).toFixed(1)
  metrics.value.memoryAvailable = (metrics.value.memoryTotal - metrics.value.memoryUsed).toFixed(1)
  metrics.value.network.upload = (Math.random() * 5 + 1).toFixed(1) + ' MB/s'
  metrics.value.network.download = (Math.random() * 10 + 5).toFixed(1) + ' MB/s'
  
  // 更新CPU历史
  metrics.value.cpuHistory.shift()
  metrics.value.cpuHistory.push(metrics.value.cpu)
}

// 获取指标颜色
const getMetricColor = (value) => {
  if (value < 50) return '#67c23a'
  if (value < 80) return '#e6a23c'
  return '#f56c6c'
}

// 开始自动刷新
const startRefresh = () => {
  refreshTimer = setInterval(() => {
    refreshMetrics()
  }, 5000)
}

// 停止刷新
const stopRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}
</script>

<style scoped lang="scss">
.monitor-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 20px;
  background: #f5f7fa;
  min-height: 100%;
}

// 统计卡片
.stats-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;

  .stat-card {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

    .stat-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 56px;
      height: 56px;
      border-radius: 12px;
      flex-shrink: 0;
    }

    .stat-content {
      flex: 1;

      .stat-label {
        font-size: 13px;
        color: #909399;
        margin-bottom: 8px;
      }

      .stat-value {
        font-size: 24px;
        font-weight: 600;
        color: #303133;
        margin-bottom: 4px;
      }

      .stat-trend {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 12px;

        &.up {
          color: #67c23a;
        }

        &.down {
          color: #f56c6c;
        }
      }
    }
  }
}

// 主内容区
.main-content {
  display: flex;
  gap: 20px;
  flex: 1;
}

// 服务器面板
.server-panel {
  width: 320px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  display: flex;
  flex-direction: column;

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border-bottom: 1px solid #ebeef5;

    .header-title {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 15px;
      font-weight: 600;
      color: #303133;
    }
  }

  .server-list {
    flex: 1;
    overflow-y: auto;
    padding: 12px;

    .server-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 16px;
      margin-bottom: 8px;
      border: 2px solid transparent;
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        background: #f5f7fa;
      }

      &.active {
        border-color: #409eff;
        background: #ecf5ff;
      }

      .server-status {
        flex-shrink: 0;

        .status-dot {
          display: block;
          width: 10px;
          height: 10px;
          border-radius: 50%;

          &.online {
            background: #67c23a;
            box-shadow: 0 0 6px rgba(103, 194, 58, 0.4);
          }

          &.offline {
            background: #f56c6c;
          }
        }
      }

      .server-info {
        flex: 1;
        min-width: 0;

        .server-name {
          font-size: 14px;
          font-weight: 600;
          color: #303133;
          margin-bottom: 4px;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .server-ip {
          font-size: 12px;
          color: #909399;
        }
      }

      .server-metrics {
        display: flex;
        flex-direction: column;
        gap: 4px;

        .metric-item {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 11px;
          color: #606266;
          font-weight: 600;
        }
      }
    }
  }
}

// 详情面板
.detail-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;

  // 服务器信息卡片
  .server-info-card {
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    padding: 20px;

    .info-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
      padding-bottom: 16px;
      border-bottom: 1px solid #ebeef5;

      .info-title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 600;
        color: #303133;
      }

      .info-actions {
        display: flex;
        gap: 8px;
      }
    }

    .info-details {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 12px;

      .detail-item {
        display: flex;
        gap: 12px;

        .detail-label {
          min-width: 80px;
          font-size: 13px;
          color: #909399;
        }

        .detail-value {
          font-size: 13px;
          color: #303133;
          font-weight: 500;
        }
      }
    }
  }

  // 指标网格
  .metrics-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;

    .metric-card {
      background: #fff;
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
      padding: 20px;

      .metric-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;

        .metric-title {
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 14px;
          font-weight: 600;
          color: #303133;
        }

        .metric-value {
          font-size: 24px;
          font-weight: 700;
          color: #303133;
        }
      }

      .metric-chart {
        margin-bottom: 16px;

        .progress-bar {
          height: 8px;
          background: #f0f0f0;
          border-radius: 4px;
          overflow: hidden;

          .progress-fill {
            height: 100%;
            border-radius: 4px;
            transition: width 0.3s ease;
          }
        }
      }

      .metric-history {
        display: flex;
        gap: 4px;
        height: 40px;

        .history-bar {
          flex: 1;
          background: #67c23a;
          border-radius: 2px 2px 0 0;
          transition: height 0.3s ease;
        }
      }

      .memory-details {
        display: flex;
        flex-direction: column;
        gap: 8px;

        .memory-info {
          display: flex;
          justify-content: space-between;
          font-size: 13px;

          .info-label {
            color: #606266;

            &:last-child {
              color: #909399;
            }
          }
        }
      }

      .disk-list {
        display: flex;
        flex-direction: column;
        gap: 12px;

        .disk-item {
          .disk-info {
            display: flex;
            justify-content: space-between;
            margin-bottom: 6px;
            font-size: 13px;

            .disk-path {
              color: #303133;
              font-weight: 500;
            }

            .disk-usage {
              color: #909399;
            }
          }

          .progress-bar {
            height: 6px;
            background: #f0f0f0;
            border-radius: 3px;
            overflow: hidden;

            .progress-fill {
              height: 100%;
              border-radius: 3px;
              transition: width 0.3s ease;
            }
          }
        }
      }

      .network-stats {
        display: flex;
        gap: 16px;
        margin-bottom: 16px;

        .network-item {
          flex: 1;
          display: flex;
          align-items: center;
          gap: 12px;
          padding: 16px;
          border-radius: 8px;

          &.upload {
            background: linear-gradient(135deg, #fff7e6 0%, #ffe7ba 100%);
          }

          &.download {
            background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
          }

          .network-icon {
            color: #606266;
          }

          .network-info {
            flex: 1;

            .network-label {
              font-size: 12px;
              color: #909399;
              margin-bottom: 4px;
            }

            .network-value {
              font-size: 18px;
              font-weight: 700;
              color: #303133;
            }
          }
        }
      }

      .network-history {
        .history-legend {
          display: flex;
          gap: 16px;
          font-size: 12px;
          color: #909399;

          .legend-item {
            display: flex;
            align-items: center;
            gap: 6px;

            .legend-color {
              display: block;
              width: 12px;
              height: 12px;
              border-radius: 2px;

              &.upload {
                background: #e6a23c;
              }

              &.download {
                background: #409eff;
              }
            }
          }
        }
      }
    }
  }

  // 进程面板
  .process-panel {
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    padding: 20px;

    .panel-title {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 15px;
      font-weight: 600;
      color: #303133;
      margin-bottom: 16px;
    }

    .process-table {
      .table-header {
        display: grid;
        grid-template-columns: 80px 100px 1fr 80px 80px 80px;
        gap: 8px;
        padding: 12px;
        background: #f5f7fa;
        border-radius: 6px;
        font-size: 13px;
        font-weight: 600;
        color: #303133;

        .header-cell {
          display: flex;
          align-items: center;
        }
      }

      .table-body {
        .table-row {
          display: grid;
          grid-template-columns: 80px 100px 1fr 80px 80px 80px;
          gap: 8px;
          padding: 12px;
          border-bottom: 1px solid #ebeef5;
          font-size: 13px;
          color: #606266;
          transition: background 0.2s;

          &:hover {
            background: #f5f7fa;
          }

          .body-cell {
            display: flex;
            align-items: center;
          }
        }
      }
    }
  }
}
</style>
