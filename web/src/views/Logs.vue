<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-4">执行日志</h1>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-row>
              <v-col cols="12" md="6">
                <v-select
                  v-model="selectedTask"
                  :items="tasks"
                  item-title="name"
                  item-value="id"
                  label="选择任务"
                  clearable
                  @update:model-value="fetchLogs"
                ></v-select>
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="statusFilter"
                  :items="statusOptions"
                  label="执行状态"
                  clearable
                  @update:model-value="fetchLogs"
                ></v-select>
              </v-col>
            </v-row>
          </v-card-title>
          
          <v-data-table
            :headers="headers"
            :items="logs"
            :loading="loading"
            :page="page"
            :items-per-page="itemsPerPage"
            :server-items-length="totalItems"
            @update:page="handlePageChange"
            @update:items-per-page="handleItemsPerPageChange"
          >
            <template v-slot:item.Success="{ item }">
              <v-chip
                :color="item.Success ? 'success' : 'error'"
                size="small"
              >
                {{ item.Success ? '成功' : '失败' }}
              </v-chip>
            </template>
            
            <template v-slot:item.StartTime="{ item }">
              {{ formatDateTime(item.StartTime) }}
            </template>
            
            <template v-slot:item.EndTime="{ item }">
              {{ formatDateTime(item.EndTime) }}
            </template>
            
            <template v-slot:item.duration="{ item }">
              {{ calculateDuration(item.StartTime, item.EndTime) }}
            </template>
            
            <template v-slot:item.actions="{ item }">
              <v-btn
                icon="mdi-eye"
                size="small"
                color="primary"
                @click="viewLogDetail(item)"
              ></v-btn>
            </template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>

    <!-- 日志详情对话框 -->
    <v-dialog v-model="logDetailDialog" max-width="800px">
      <v-card v-if="selectedLog">
        <v-card-title>
          执行日志详情 - {{ selectedLog.TaskName }}
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="6">
              <strong>任务ID:</strong> {{ selectedLog.TaskID }}
            </v-col>
            <v-col cols="6">
              <strong>执行状态:</strong>
              <v-chip
                :color="selectedLog.Success ? 'success' : 'error'"
                size="small"
                class="ml-2"
              >
                {{ selectedLog.Success ? '成功' : '失败' }}
              </v-chip>
            </v-col>
            <v-col cols="6">
              <strong>开始时间:</strong> {{ formatDateTime(selectedLog.StartTime) }}
            </v-col>
            <v-col cols="6">
              <strong>结束时间:</strong> {{ formatDateTime(selectedLog.EndTime) }}
            </v-col>
            <v-col cols="6">
              <strong>执行时长:</strong> {{ calculateDuration(selectedLog.StartTime, selectedLog.EndTime) }}
            </v-col>
          </v-row>
          
          <v-divider class="my-4"></v-divider>
          
          <div v-if="selectedLog.Output">
            <strong>执行输出:</strong>
            <v-card class="mt-2" variant="outlined">
              <v-card-text>
                <pre class="text-body-2">{{ selectedLog.Output }}</pre>
              </v-card-text>
            </v-card>
          </div>
          
          <div v-if="selectedLog.Error" class="mt-4">
            <strong>错误信息:</strong>
            <v-card class="mt-2" variant="outlined" color="error">
              <v-card-text>
                <pre class="text-body-2">{{ selectedLog.Error }}</pre>
              </v-card-text>
            </v-card>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="logDetailDialog = false">关闭</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../services/api'

const loading = ref(false)
const logs = ref([])
const tasks = ref([])
const selectedTask = ref(null)
const statusFilter = ref(null)
const page = ref(1)
const itemsPerPage = ref(10)
const totalItems = ref(0)

const logDetailDialog = ref(false)
const selectedLog = ref(null)

const headers = [
  { title: 'ID', key: 'ID', width: '80px' },
  { title: '任务名称', key: 'TaskName' },
  { title: '开始时间', key: 'StartTime' },
  { title: '结束时间', key: 'EndTime' },
  { title: '执行时长', key: 'duration' },
  { title: '状态', key: 'Success', width: '100px' },
  { title: '操作', key: 'actions', sortable: false, width: '80px' }
]

const statusOptions = [
  { title: '成功', value: true },
  { title: '失败', value: false }
]

const fetchTasks = async () => {
  try {
    const response = await api.get('/tasks')
    tasks.value = response.data
  } catch (error) {
    console.error('获取任务列表失败:', error)
  }
}

const fetchLogs = async () => {
  loading.value = true
  try {
    let url = '/tasks'
    let params = {
      page: page.value,
      page_size: itemsPerPage.value
    }

    if (selectedTask.value) {
      url = `/tasks/${selectedTask.value}/logs/paginated`
    } else {
      // 如果没有选择特定任务，需要获取所有日志
      // 这里假设有一个获取所有日志的接口
      url = '/logs/paginated'
    }

    const response = await api.get(url, { params })
    
    if (response.data.data) {
      logs.value = response.data.data
      totalItems.value = response.data.total
    } else {
      logs.value = response.data || []
      totalItems.value = logs.value.length
    }

    // 如果有状态过滤器，在前端过滤
    if (statusFilter.value !== null) {
      logs.value = logs.value.filter(log => log.Success === statusFilter.value)
    }
  } catch (error) {
    console.error('获取日志失败:', error)
    logs.value = []
  } finally {
    loading.value = false
  }
}

const handlePageChange = (newPage) => {
  page.value = newPage
  fetchLogs()
}

const handleItemsPerPageChange = (newItemsPerPage) => {
  itemsPerPage.value = newItemsPerPage
  page.value = 1
  fetchLogs()
}

const viewLogDetail = (log) => {
  selectedLog.value = log
  logDetailDialog.value = true
}

const formatDateTime = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}

const calculateDuration = (startTime, endTime) => {
  if (!startTime || !endTime) return '-'
  const start = new Date(startTime)
  const end = new Date(endTime)
  const duration = end - start
  
  if (duration < 1000) {
    return `${duration}ms`
  } else if (duration < 60000) {
    return `${(duration / 1000).toFixed(1)}s`
  } else {
    return `${(duration / 60000).toFixed(1)}min`
  }
}

onMounted(() => {
  fetchTasks()
  fetchLogs()
})
</script>

<style scoped>
pre {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>