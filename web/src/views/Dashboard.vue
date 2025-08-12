<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-4">仪表板</h1>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12" md="3">
        <v-card>
          <v-card-text>
            <div class="text-h6">总任务数</div>
            <div class="text-h4 primary--text">{{ stats.totalTasks }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card>
          <v-card-text>
            <div class="text-h6">启用任务</div>
            <div class="text-h4 success--text">{{ stats.enabledTasks }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card>
          <v-card-text>
            <div class="text-h6">今日执行</div>
            <div class="text-h4 info--text">{{ stats.todayExecutions }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card>
          <v-card-text>
            <div class="text-h6">成功率</div>
            <div class="text-h4 warning--text">{{ stats.successRate }}%</div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12" md="8">
        <v-card>
          <v-card-title>最近任务</v-card-title>
          <v-card-text>
            <v-data-table
              :headers="taskHeaders"
              :items="recentTasks"
              :loading="loading"
              hide-default-footer
              density="compact"
            >
              <template v-slot:item.enabled="{ item }">
                <v-chip
                  :color="item.enabled ? 'success' : 'error'"
                  size="small"
                >
                  {{ item.enabled ? '启用' : '禁用' }}
                </v-chip>
              </template>
              <template v-slot:item.actions="{ item }">
                <v-btn
                  icon="mdi-play"
                  size="small"
                  color="primary"
                  @click="executeTask(item.id)"
                ></v-btn>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card>
          <v-card-title>系统状态</v-card-title>
          <v-card-text>
            <div v-if="systemStats">
              <div class="mb-2">
                <div class="text-subtitle2">CPU 使用率</div>
                <v-progress-linear
                  :model-value="systemStats.CPUUsage"
                  color="primary"
                  height="20"
                >
                  <template v-slot:default="{ value }">
                    <strong>{{ Math.ceil(value) }}%</strong>
                  </template>
                </v-progress-linear>
              </div>
              <div class="mb-2">
                <div class="text-subtitle2">内存使用率</div>
                <v-progress-linear
                  :model-value="systemStats.MemoryUsage"
                  color="success"
                  height="20"
                >
                  <template v-slot:default="{ value }">
                    <strong>{{ Math.ceil(value) }}%</strong>
                  </template>
                </v-progress-linear>
              </div>
              <div class="mb-2">
                <div class="text-subtitle2">磁盘使用率</div>
                <v-progress-linear
                  :model-value="systemStats.DiskUsage"
                  color="warning"
                  height="20"
                >
                  <template v-slot:default="{ value }">
                    <strong>{{ Math.ceil(value) }}%</strong>
                  </template>
                </v-progress-linear>
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
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

const taskHeaders = [
  { title: '任务名称', key: 'name' },
  { title: '计划', key: 'schedule' },
  { title: '状态', key: 'enabled' },
  { title: '操作', key: 'actions', sortable: false }
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
  } finally {
    loading.value = false
  }
}

const executeTask = async (taskId) => {
  try {
    await api.post(`/tasks/${taskId}/execute`)
    // 可以添加成功提示
  } catch (error) {
    console.error('执行任务失败:', error)
  }
}

onMounted(() => {
  fetchDashboardData()
  // 每30秒刷新一次数据
  setInterval(fetchDashboardData, 30000)
})
</script>