<template>
  <div>
    <h1 style="margin-bottom: 24px;">系统监控</h1>

    <div v-if="systemStats">
      <a-row :gutter="16" style="margin-bottom: 24px;">
        <a-col :xs="24" :md="12">
          <a-card title="CPU 使用率">
            <div style="font-size: 32px; font-weight: bold; margin-bottom: 16px; color: #1890ff;">
              {{ systemStats.CPUUsage.toFixed(1) }}%
            </div>
            <a-progress
              :percent="systemStats.CPUUsage"
              stroke-color="#1890ff"
              :show-info="false"
            />
          </a-card>
        </a-col>
        
        <a-col :xs="24" :md="12">
          <a-card title="内存使用">
            <div style="font-size: 32px; font-weight: bold; margin-bottom: 16px; color: #52c41a;">
              {{ systemStats.MemoryUsage.toFixed(1) }}%
            </div>
            <a-progress
              :percent="systemStats.MemoryUsage"
              stroke-color="#52c41a"
              :show-info="false"
            />
            <div style="margin-top: 8px; color: #666; font-size: 12px;">
              已用: {{ formatBytes(systemStats.MemoryUsed * 1024 * 1024) }} / 
              总计: {{ formatBytes(systemStats.MemoryTotal * 1024 * 1024) }}
            </div>
          </a-card>
        </a-col>
        
        <a-col :xs="24" :md="12">
          <a-card title="磁盘使用">
            <div style="font-size: 32px; font-weight: bold; margin-bottom: 16px; color: #faad14;">
              {{ systemStats.DiskUsage.toFixed(1) }}%
            </div>
            <a-progress
              :percent="systemStats.DiskUsage"
              stroke-color="#faad14"
              :show-info="false"
            />
            <div style="margin-top: 8px; color: #666; font-size: 12px;">
              已用: {{ formatBytes(systemStats.DiskUsed * 1024 * 1024 * 1024) }} / 
              总计: {{ formatBytes(systemStats.DiskTotal * 1024 * 1024 * 1024) }}
            </div>
          </a-card>
        </a-col>
        
        <a-col :xs="24" :md="12">
          <a-card title="系统负载">
            <div style="font-size: 32px; font-weight: bold; margin-bottom: 16px; color: #13c2c2;">
              {{ systemStats.SystemLoad.toFixed(2) }}
            </div>
            <a-progress
              :percent="Math.min(systemStats.SystemLoad * 25, 100)"
              stroke-color="#13c2c2"
              :show-info="false"
            />
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="16" style="margin-bottom: 24px;">
        <a-col :xs="24" :md="8">
          <a-card title="网络流量">
            <div style="margin-bottom: 16px;">
              <div style="color: #666; margin-bottom: 4px;">接收</div>
              <div style="font-size: 18px; font-weight: 500;">{{ formatBytes(systemStats.NetworkRxBytes) }}</div>
            </div>
            <div>
              <div style="color: #666; margin-bottom: 4px;">发送</div>
              <div style="font-size: 18px; font-weight: 500;">{{ formatBytes(systemStats.NetworkTxBytes) }}</div>
            </div>
          </a-card>
        </a-col>
        
        <a-col :xs="24" :md="8">
          <a-card title="进程信息">
            <div style="margin-bottom: 16px;">
              <div style="color: #666; margin-bottom: 4px;">进程数</div>
              <div style="font-size: 18px; font-weight: 500;">{{ systemStats.ProcessCount }}</div>
            </div>
            <div>
              <div style="color: #666; margin-bottom: 4px;">协程数</div>
              <div style="font-size: 18px; font-weight: 500;">{{ systemStats.GoroutineCount }}</div>
            </div>
          </a-card>
        </a-col>
        
        <a-col :xs="24" :md="8">
          <a-card title="运行时间">
            <div style="font-size: 18px; font-weight: 500; margin-bottom: 8px;">
              {{ formatUptime(systemStats.Uptime) }}
            </div>
            <div style="color: #666; font-size: 12px;">
              更新时间: {{ formatDateTime(systemStats.Timestamp) }}
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <a-card>
      <template #title>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span>系统状态历史</span>
          <a-button
            type="primary"
            @click="fetchSystemStats"
            :loading="loading"
          >
            <ReloadOutlined />
            刷新
          </a-button>
        </div>
      </template>
      
      <div v-if="loading" style="text-align: center; padding: 40px;">
        <a-spin size="large" />
      </div>
      <div v-else-if="!systemStats" style="text-align: center; color: #999; padding: 40px;">
        暂无数据
      </div>
      <div v-else>
        <a-alert
          type="info"
          message="系统监控数据每30秒自动刷新一次"
          style="margin-bottom: 16px;"
          show-icon
        />
        
        <!-- 这里可以添加图表组件来显示历史数据 -->
        <div style="text-align: center; color: #999; padding: 40px;">
          历史趋势图表功能待开发
        </div>
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
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
    message.error('获取系统状态失败')
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