<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px;">
      <h1>任务执行统计</h1>
      <a-space>
        <a-select
          v-model:value="selectedDays"
          style="width: 120px"
          @change="refreshData"
        >
          <a-select-option value="7">最近7天</a-select-option>
          <a-select-option value="30">最近30天</a-select-option>
          <a-select-option value="90">最近90天</a-select-option>
        </a-select>
        <a-button @click="refreshData" :loading="loading">
          <ReloadOutlined />
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 概览卡片 -->
    <a-row :gutter="16" style="margin-bottom: 24px;">
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="总任务数"
            :value="report.total_tasks"
            :value-style="{ color: '#1890ff' }"
          >
            <template #prefix>
              <AppstoreOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="活跃任务数"
            :value="report.active_tasks"
            :value-style="{ color: '#52c41a' }"
          >
            <template #prefix>
              <PlayCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="总执行次数"
            :value="report.total_executions"
            :value-style="{ color: '#722ed1' }"
          >
            <template #prefix>
              <ThunderboltOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="整体成功率"
            :value="report.success_rate"
            suffix="%"
            :precision="2"
            :value-style="{ color: report.success_rate >= 90 ? '#52c41a' : report.success_rate >= 70 ? '#faad14' : '#f5222d' }"
          >
            <template #prefix>
              <CheckCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 图表区域 -->
    <a-row :gutter="16" style="margin-bottom: 24px;">
      <!-- 执行趋势图 -->
      <a-col :span="12">
        <a-card title="执行趋势" :loading="loading">
          <div ref="trendChart" style="height: 300px;"></div>
        </a-card>
      </a-col>
      <!-- 小时分布图 -->
      <a-col :span="12">
        <a-card title="小时分布" :loading="loading">
          <div ref="hourlyChart" style="height: 300px;"></div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 任务统计表格 -->
    <a-row :gutter="16" style="margin-bottom: 24px;">
      <a-col :span="24">
        <a-card title="任务执行统计">
          <a-table
            :columns="statisticsColumns"
            :data-source="taskStatistics"
            :loading="loading"
            :pagination="{ pageSize: 10 }"
            size="middle"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'success_rate'">
                <a-progress
                  :percent="record.success_rate"
                  :status="record.success_rate >= 90 ? 'success' : record.success_rate >= 70 ? 'normal' : 'exception'"
                  :show-info="true"
                  size="small"
                />
              </template>
              <template v-else-if="column.key === 'average_execution_time'">
                {{ formatDuration(record.average_execution_time) }}
              </template>
              <template v-else-if="column.key === 'last_execution_time'">
                <span v-if="record.last_execution_time">
                  {{ formatDateTime(record.last_execution_time) }}
                </span>
                <span v-else style="color: #999;">从未执行</span>
              </template>
              <template v-else-if="column.key === 'last_execution_status'">
                <a-tag :color="record.last_execution_status ? 'success' : 'error'">
                  {{ record.last_execution_status ? '成功' : '失败' }}
                </a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <!-- 性能指标 -->
    <a-row :gutter="16">
      <a-col :span="24">
        <a-card title="性能指标">
          <a-table
            :columns="performanceColumns"
            :data-source="performanceMetrics"
            :loading="loading"
            :pagination="{ pageSize: 10 }"
            size="middle"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key.includes('time')">
                {{ formatDuration(record[column.key]) }}
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import {
  ReloadOutlined,
  AppstoreOutlined,
  PlayCircleOutlined,
  ThunderboltOutlined,
  CheckCircleOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import * as echarts from 'echarts'
import api from '../services/api'

const loading = ref(false)
const selectedDays = ref(30)
const report = ref({
  total_tasks: 0,
  active_tasks: 0,
  total_executions: 0,
  success_rate: 0
})
const taskStatistics = ref([])
const performanceMetrics = ref([])
const executionTrends = ref([])
const hourlyStats = ref([])

// 图表实例
const trendChart = ref(null)
const hourlyChart = ref(null)
let trendChartInstance = null
let hourlyChartInstance = null

// 表格列定义
const statisticsColumns = [
  { title: '任务名称', dataIndex: 'task_name', key: 'task_name', width: 200 },
  { title: '总执行次数', dataIndex: 'total_executions', key: 'total_executions', width: 120, sorter: true },
  { title: '成功次数', dataIndex: 'success_executions', key: 'success_executions', width: 100 },
  { title: '失败次数', dataIndex: 'failure_executions', key: 'failure_executions', width: 100 },
  { title: '成功率', key: 'success_rate', width: 150 },
  { title: '平均执行时间', key: 'average_execution_time', width: 120 },
  { title: '最后执行时间', key: 'last_execution_time', width: 160 },
  { title: '最后执行状态', key: 'last_execution_status', width: 120 }
]

const performanceColumns = [
  { title: '任务名称', dataIndex: 'task_name', key: 'task_name', width: 200 },
  { title: '最短时间', dataIndex: 'min_execution_time', key: 'min_execution_time', width: 100 },
  { title: '最长时间', dataIndex: 'max_execution_time', key: 'max_execution_time', width: 100 },
  { title: '平均时间', dataIndex: 'average_execution_time', key: 'average_execution_time', width: 100 },
  { title: '中位数时间', dataIndex: 'median_execution_time', key: 'median_execution_time', width: 100 },
  { title: '标准差', dataIndex: 'execution_time_std_dev', key: 'execution_time_std_dev', width: 100 }
]

// 获取所有统计数据
const fetchAllData = async () => {
  loading.value = true
  try {
    const params = { days: selectedDays.value }
    
    // 并行获取所有数据
    const [
      reportRes,
      statisticsRes,
      performanceRes,
      trendsRes,
      hourlyRes
    ] = await Promise.all([
      api.get('/statistics/report', { params }),
      api.get('/statistics/tasks', { params }),
      api.get('/statistics/performance', { params }),
      api.get('/statistics/trends', { params }),
      api.get('/statistics/hourly', { params })
    ])

    report.value = reportRes.data || {
      total_tasks: 0,
      active_tasks: 0,
      total_executions: 0,
      success_rate: 0
    }
    taskStatistics.value = statisticsRes.data || []
    performanceMetrics.value = performanceRes.data || []
    executionTrends.value = trendsRes.data || []
    hourlyStats.value = hourlyRes.data || []

    // 更新图表
    await nextTick()
    setTimeout(() => {
      updateCharts()
    }, 100) // 延迟一点时间确保DOM完全渲染
  } catch (error) {
    message.error('获取统计数据失败')
    console.error('Error fetching statistics:', error)
  } finally {
    loading.value = false
  }
}

// 刷新数据
const refreshData = () => {
  fetchAllData()
}

// 更新图表
const updateCharts = () => {
  try {
    updateTrendChart()
    updateHourlyChart()
  } catch (error) {
    console.error('Error updating charts:', error)
  }
}

// 更新趋势图表
const updateTrendChart = () => {
  if (!trendChart.value) {
    console.warn('Trend chart container not ready')
    return
  }
  
  if (!trendChartInstance) {
    trendChartInstance = echarts.init(trendChart.value)
  }

  if (!executionTrends.value || executionTrends.value.length === 0) {
    console.warn('No execution trends data available')
    return
  }

  const dates = executionTrends.value.map(item => item.date || '')
  const totalExecutions = executionTrends.value.map(item => item.total_executions || 0)
  const successRates = executionTrends.value.map(item => item.success_rate || 0)

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      }
    },
    legend: {
      data: ['执行次数', '成功率']
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLabel: {
        rotate: 45
      }
    },
    yAxis: [
      {
        type: 'value',
        name: '执行次数',
        position: 'left'
      },
      {
        type: 'value',
        name: '成功率 (%)',
        position: 'right',
        max: 100
      }
    ],
    series: [
      {
        name: '执行次数',
        type: 'bar',
        data: totalExecutions,
        itemStyle: {
          color: '#1890ff'
        }
      },
      {
        name: '成功率',
        type: 'line',
        yAxisIndex: 1,
        data: successRates,
        itemStyle: {
          color: '#52c41a'
        },
        lineStyle: {
          width: 2
        }
      }
    ]
  }

  trendChartInstance.setOption(option)
}

// 更新小时分布图表
const updateHourlyChart = () => {
  if (!hourlyChart.value) {
    console.warn('Hourly chart container not ready')
    return
  }
  
  if (!hourlyChartInstance) {
    hourlyChartInstance = echarts.init(hourlyChart.value)
  }

  if (!hourlyStats.value || hourlyStats.value.length === 0) {
    console.warn('No hourly stats data available')
    return
  }

  const hours = hourlyStats.value.map(item => `${item.hour || 0}:00`)
  const executions = hourlyStats.value.map(item => item.total_executions || 0)

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: hours,
      name: '小时'
    },
    yAxis: {
      type: 'value',
      name: '执行次数'
    },
    series: [
      {
        name: '执行次数',
        type: 'bar',
        data: executions,
        itemStyle: {
          color: '#722ed1'
        }
      }
    ]
  }

  hourlyChartInstance.setOption(option)
}

// 格式化时间
const formatDuration = (seconds) => {
  if (!seconds || seconds === 0) return '0s'
  
  if (seconds < 60) {
    return `${seconds.toFixed(2)}s`
  } else if (seconds < 3600) {
    const minutes = Math.floor(seconds / 60)
    const remainingSeconds = (seconds % 60).toFixed(0)
    return `${minutes}m${remainingSeconds}s`
  } else {
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    return `${hours}h${minutes}m`
  }
}

// 格式化日期时间
const formatDateTime = (dateTime) => {
  if (!dateTime) return ''
  return new Date(dateTime).toLocaleString('zh-CN')
}

// 窗口大小变化时重新调整图表
const handleResize = () => {
  if (trendChartInstance) {
    trendChartInstance.resize()
  }
  if (hourlyChartInstance) {
    hourlyChartInstance.resize()
  }
}

onMounted(async () => {
  await nextTick() // 确保DOM已经渲染
  fetchAllData()
  window.addEventListener('resize', handleResize)
})

// 组件卸载时清理
import { onUnmounted } from 'vue'
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (trendChartInstance) {
    trendChartInstance.dispose()
  }
  if (hourlyChartInstance) {
    hourlyChartInstance.dispose()
  }
})
</script>

<style scoped>
.ant-statistic {
  text-align: center;
}

.ant-card {
  margin-bottom: 16px;
}

.ant-table {
  background: white;
}
</style>