<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-4">系统监控</h1>
      </v-col>
    </v-row>

    <v-row v-if="systemStats">
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>CPU 使用率</v-card-title>
          <v-card-text>
            <div class="text-h3 mb-2">{{ systemStats.CPUUsage.toFixed(1) }}%</div>
            <v-progress-linear
              :model-value="systemStats.CPUUsage"
              color="primary"
              height="20"
            ></v-progress-linear>
          </v-card-text>
        </v-card>
      </v-col>
      
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>内存使用</v-card-title>
          <v-card-text>
            <div class="text-h3 mb-2">{{ systemStats.MemoryUsage.toFixed(1) }}%</div>
            <v-progress-linear
              :model-value="systemStats.MemoryUsage"
              color="success"
              height="20"
            ></v-progress-linear>
            <div class="text-caption mt-2">
              已用: {{ formatBytes(systemStats.MemoryUsed * 1024 * 1024) }} / 
              总计: {{ formatBytes(systemStats.MemoryTotal * 1024 * 1024) }}
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>磁盘使用</v-card-title>
          <v-card-text>
            <div class="text-h3 mb-2">{{ systemStats.DiskUsage.toFixed(1) }}%</div>
            <v-progress-linear
              :model-value="systemStats.DiskUsage"
              color="warning"
              height="20"
            ></v-progress-linear>
            <div class="text-caption mt-2">
              已用: {{ formatBytes(systemStats.DiskUsed * 1024 * 1024 * 1024) }} / 
              总计: {{ formatBytes(systemStats.DiskTotal * 1024 * 1024 * 1024) }}
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>系统负载</v-card-title>
          <v-card-text>
            <div class="text-h3 mb-2">{{ systemStats.SystemLoad.toFixed(2) }}</div>
            <v-progress-linear
              :model-value="Math.min(systemStats.SystemLoad * 25, 100)"
              color="info"
              height="20"
            ></v-progress-linear>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row v-if="systemStats">
      <v-col cols="12" md="4">
        <v-card>
          <v-card-title>网络流量</v-card-title>
          <v-card-text>
            <div class="mb-2">
              <div class="text-subtitle2">接收</div>
              <div class="text-h6">{{ formatBytes(systemStats.NetworkRxBytes) }}</div>
            </div>
            <div>
              <div class="text-subtitle2">发送</div>
              <div class="text-h6">{{ formatBytes(systemStats.NetworkTxBytes) }}</div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      
      <v-col cols="12" md="4">
        <v-card>
          <v-card-title>进程信息</v-card-title>
          <v-card-text>
            <div class="mb-2">
              <div class="text-subtitle2">进程数</div>
              <div class="text-h6">{{ systemStats.ProcessCount }}</div>
            </div>
            <div>
              <div class="text-subtitle2">协程数</div>
              <div class="text-h6">{{ systemStats.GoroutineCount }}</div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      
      <v-col cols="12" md="4">
        <v-card>
          <v-card-title>运行时间</v-card-title>
          <v-card-text>
            <div class="text-h6">{{ formatUptime(systemStats.Uptime) }}</div>
            <div class="text-caption">
              更新时间: {{ formatDateTime(systemStats.Timestamp) }}
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <div class="d-flex justify-space-between align-center">
              <span>系统状态历史</span>
              <v-btn
                icon="mdi-refresh"
                @click="fetchSystemStats"
                :loading="loading"
              ></v-btn>
            </div>
          </v-card-title>
          <v-card-text>
            <div class="text-center" v-if="loading">
              <v-progress-circular indeterminate></v-progress-circular>
            </div>
            <div v-else-if="!systemStats" class="text-center text-grey">
              暂无数据
            </div>
            <div v-else>
              <v-alert
                type="info"
                variant="tonal"
                class="mb-4"
              >
                系统监控数据每30秒自动刷新一次
              </v-alert>
              
              <!-- 这里可以添加图表组件来显示历史数据 -->
              <div class="text-center text-grey">
                历史趋势图表功能待开发
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import api from '../services/api'

const loading = ref(false)
const systemStats = ref(null)
let refreshInterval = null

const fetchSystemStats = async () => {
  loading.value = true
  try {
    const response = await api.get('/system/stats')
    systemStats.value = response.data
  } catch (error) {
    console.error('获取系统状态失败:', error)
  } finally {
    loading.value = false
  }
}

const formatBytes = (bytes) => {
  if (bytes === 0) return '0 B'
  
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatUptime = (seconds) => {
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  
  if (days > 0) {
    return `${days}天 ${hours}小时 ${minutes}分钟`
  } else if (hours > 0) {
    return `${hours}小时 ${minutes}分钟`
  } else {
    return `${minutes}分钟`
  }
}

const formatDateTime = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchSystemStats()
  // 每30秒刷新一次数据
  refreshInterval = setInterval(fetchSystemStats, 30000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>