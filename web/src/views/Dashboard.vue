<template>
  <div>
    <h1 style="margin-bottom: 24px;">仪表板</h1>

    <a-row :gutter="16" style="margin-bottom: 24px;">
      <a-col :xs="24" :sm="12" :md="6">
        <a-card>
          <a-statistic
            title="总任务数"
            :value="stats.totalTasks"
            :value-style="{ color: '#1890ff' }"
          />
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card>
          <a-statistic
            title="启用任务"
            :value="stats.enabledTasks"
            :value-style="{ color: '#52c41a' }"
          />
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card>
          <a-statistic
            title="今日执行"
            :value="stats.todayExecutions"
            :value-style="{ color: '#13c2c2' }"
          />
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card>
          <a-statistic
            title="成功率"
            :value="stats.successRate"
            suffix="%"
            :value-style="{ color: '#faad14' }"
          />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16">
      <a-col :xs="24" :lg="16">
        <a-card title="最近任务">
          <a-table
            :columns="taskColumns"
            :data-source="recentTasks"
            :loading="loading"
            :pagination="false"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'enabled'">
                <a-tag :color="record.enabled ? 'success' : 'error'">
                  {{ record.enabled ? '启用' : '禁用' }}
                </a-tag>
              </template>
              <template v-else-if="column.key === 'actions'">
                <a-button
                  type="primary"
                  size="small"
                  @click="executeTask(record.id)"
                >
                  <PlayCircleOutlined />
                </a-button>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
      <a-col :xs="24" :lg="8">
        <a-card title="系统状态">
          <div v-if="systemStats">
            <div style="margin-bottom: 16px;">
              <div style="margin-bottom: 8px;">CPU 使用率</div>
              <a-progress
                :percent="systemStats.CPUUsage"
                :show-info="true"
                stroke-color="#1890ff"
              />
            </div>
            <div style="margin-bottom: 16px;">
              <div style="margin-bottom: 8px;">内存使用率</div>
              <a-progress
                :percent="systemStats.MemoryUsage"
                :show-info="true"
                stroke-color="#52c41a"
              />
            </div>
            <div>
              <div style="margin-bottom: 8px;">磁盘使用率</div>
              <a-progress
                :percent="systemStats.DiskUsage"
                :show-info="true"
                stroke-color="#faad14"
              />
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { PlayCircleOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import api from '../services/api'

const loading = ref(false)
const stats = ref({
  totalTasks: 0,
  enabledTasks: 0,
  todayExecutions: 0,
  successRate: 0
})
const recentTasks = ref([])
const systemStats = ref(null)

const taskColumns = [
  { title: '任务名称', dataIndex: 'name', key: 'name' },
  { title: '计划', dataIndex: 'schedule', key: 'schedule' },
  { title: '状态', dataIndex: 'enabled', key: 'enabled' },
  { title: '操作', key: 'actions' }
]

const fetchDashboardData = async () => {
  loading.value = true
  try {
    // 获取任务列表
    const tasksResponse = await api.get('/tasks')
    const tasks = tasksResponse.data
    
    stats.value.totalTasks = tasks.length
    stats.value.enabledTasks = tasks.filter(task => task.enabled).length
    recentTasks.value = tasks.slice(0, 5)

    // 获取系统状态
    const systemResponse = await api.get('/system/stats')
    systemStats.value = systemResponse.data

    // 模拟统计数据（实际应该从API获取）
    stats.value.todayExecutions = Math.floor(Math.random() * 100)
    stats.value.successRate = Math.floor(Math.random() * 20 + 80)
  } catch (error) {
    console.error('获取仪表板数据失败:', error)
    message.error('获取仪表板数据失败')
  } finally {
    loading.value = false
  }
}

const executeTask = async (taskId) => {
  try {
    await api.post(`/tasks/${taskId}/execute`)
    message.success('任务执行成功')
  } catch (error) {
    console.error('执行任务失败:', error)
    message.error('任务执行失败')
  }
}

onMounted(() => {
  fetchDashboardData()
  // 每30秒刷新一次数据
  setInterval(fetchDashboardData, 30000)
})
</script>