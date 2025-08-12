<template>
  <div>
    <h1 style="margin-bottom: 24px;">执行日志</h1>

    <a-card>
      <template #title>
        <a-row :gutter="16">
          <a-col :xs="24" :md="12">
            <a-select
              v-model:value="selectedTask"
              placeholder="选择任务"
              allow-clear
              style="width: 100%"
              @change="fetchLogs"
            >
              <a-select-option
                v-for="task in tasks"
                :key="task.id"
                :value="task.id"
              >
                {{ task.name }}
              </a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :md="12">
            <a-select
              v-model:value="statusFilter"
              placeholder="执行状态"
              allow-clear
              style="width: 100%"
              @change="fetchLogs"
            >
              <a-select-option
                v-for="option in statusOptions"
                :key="option.value"
                :value="option.value"
              >
                {{ option.title }}
              </a-select-option>
            </a-select>
          </a-col>
        </a-row>
      </template>
      
      <a-table
        :columns="columns"
        :data-source="logs"
        :loading="loading"
        :pagination="{
          current: page,
          pageSize: itemsPerPage,
          total: totalItems,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`
        }"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'Success'">
            <a-tag :color="record.Success ? 'success' : 'error'">
              {{ record.Success ? '成功' : '失败' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'StartTime'">
            {{ formatDateTime(record.StartTime) }}
          </template>
          <template v-else-if="column.key === 'EndTime'">
            {{ formatDateTime(record.EndTime) }}
          </template>
          <template v-else-if="column.key === 'duration'">
            {{ calculateDuration(record.StartTime, record.EndTime) }}
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-button
              type="primary"
              size="small"
              @click="viewLogDetail(record)"
            >
              <EyeOutlined />
            </a-button>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 日志详情对话框 -->
    <a-modal
      v-model:open="logDetailDialog"
      title="执行日志详情"
      width="800px"
      :footer="null"
    >
      <div v-if="selectedLog">
        <h3>{{ selectedLog.TaskName }}</h3>
        
        <a-descriptions :column="2" bordered style="margin: 16px 0;">
          <a-descriptions-item label="任务ID">
            {{ selectedLog.TaskID }}
          </a-descriptions-item>
          <a-descriptions-item label="执行状态">
            <a-tag :color="selectedLog.Success ? 'success' : 'error'">
              {{ selectedLog.Success ? '成功' : '失败' }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="开始时间">
            {{ formatDateTime(selectedLog.StartTime) }}
          </a-descriptions-item>
          <a-descriptions-item label="结束时间">
            {{ formatDateTime(selectedLog.EndTime) }}
          </a-descriptions-item>
          <a-descriptions-item label="执行时长" :span="2">
            {{ calculateDuration(selectedLog.StartTime, selectedLog.EndTime) }}
          </a-descriptions-item>
        </a-descriptions>
        
        <div v-if="selectedLog.Output" style="margin-bottom: 16px;">
          <h4>执行输出:</h4>
          <a-card size="small">
            <pre style="white-space: pre-wrap; word-break: break-word; margin: 0;">{{ selectedLog.Output }}</pre>
          </a-card>
        </div>
        
        <div v-if="selectedLog.Error">
          <h4>错误信息:</h4>
          <a-alert
            type="error"
            show-icon
          >
            <template #message>
              <pre style="white-space: pre-wrap; word-break: break-word; margin: 0;">{{ selectedLog.Error }}</pre>
            </template>
          </a-alert>
        </div>
        
        <div style="text-align: right; margin-top: 16px;">
          <a-button @click="logDetailDialog = false">关闭</a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { EyeOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
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

const columns = [
  { title: 'ID', dataIndex: 'ID', key: 'ID', width: 80 },
  { title: '任务名称', dataIndex: 'TaskName', key: 'TaskName' },
  { title: '开始时间', key: 'StartTime' },
  { title: '结束时间', key: 'EndTime' },
  { title: '执行时长', key: 'duration' },
  { title: '状态', key: 'Success', width: 100 },
  { title: '操作', key: 'actions', width: 80 }
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
    message.error('获取任务列表失败')
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
    message.error('获取日志失败')
    logs.value = []
  } finally {
    loading.value = false
  }
}

const handleTableChange = (pagination) => {
  page.value = pagination.current
  itemsPerPage.value = pagination.pageSize
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