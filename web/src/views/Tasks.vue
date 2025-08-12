<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px;">
      <h1>任务管理</h1>
      <a-button type="primary" @click="openTaskDialog()">
        <PlusOutlined />
        新建任务
      </a-button>
    </div>

    <a-card>
      <template #title>
        <a-input-search
          v-model:value="search"
          placeholder="搜索任务"
          style="width: 300px;"
          allow-clear
        />
      </template>
      
      <a-table
        :columns="columns"
        :data-source="filteredTasks"
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
          <template v-if="column.key === 'enabled'">
            <a-switch
              v-model:checked="record.enabled"
              @change="updateTaskStatus(record)"
            />
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button
                type="primary"
                size="small"
                @click="executeTask(record.id)"
              >
                <PlayCircleOutlined />
              </a-button>
              <a-button
                size="small"
                @click="openTaskDialog(record)"
              >
                <EditOutlined />
              </a-button>
              <a-popconfirm
                title="确定要删除这个任务吗？"
                @confirm="deleteTask(record)"
              >
                <a-button
                  danger
                  size="small"
                >
                  <DeleteOutlined />
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 任务编辑对话框 -->
    <a-modal
      v-model:open="taskDialog"
      :title="editingTask ? '编辑任务' : '新建任务'"
      width="600px"
      @ok="saveTask"
      @cancel="taskDialog = false"
      :confirm-loading="saving"
    >
      <a-form
        ref="taskFormRef"
        :model="taskFormData"
        :rules="taskFormRules"
        layout="vertical"
      >
        <a-form-item label="任务名称" name="name">
          <a-input v-model:value="taskFormData.name" />
        </a-form-item>
        
        <a-form-item label="Cron 表达式" name="schedule">
          <a-input
            v-model:value="taskFormData.schedule"
            placeholder="例如: 0 */5 * * * * (每5分钟执行一次)"
          />
        </a-form-item>
        
        <a-form-item label="命令或URL" name="command">
          <a-input v-model:value="taskFormData.command" />
        </a-form-item>
        
        <a-form-item label="HTTP方法">
          <a-select
            v-model:value="taskFormData.method"
            placeholder="仅当命令为URL时需要"
          >
            <a-select-option
              v-for="method in httpMethods"
              :key="method"
              :value="method"
            >
              {{ method }}
            </a-select-option>
          </a-select>
        </a-form-item>
        
        <a-form-item label="请求头 (JSON格式)">
          <a-textarea
            v-model:value="taskFormData.headers"
            :rows="3"
            placeholder='例如: {"Content-Type": "application/json"}'
          />
        </a-form-item>
        
        <a-form-item label="任务描述">
          <a-textarea
            v-model:value="taskFormData.description"
            :rows="2"
          />
        </a-form-item>
        
        <a-form-item label="启用任务">
          <a-switch v-model:checked="taskFormData.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import {
  PlusOutlined,
  PlayCircleOutlined,
  EditOutlined,
  DeleteOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import api from '../services/api'

const loading = ref(false)
const saving = ref(false)
const search = ref('')
const tasks = ref([])
const page = ref(1)
const itemsPerPage = ref(10)
const totalItems = ref(0)

const taskDialog = ref(false)
const taskFormRef = ref(null)
const editingTask = ref(null)
const taskFormData = ref({
  name: '',
  schedule: '',
  command: '',
  method: 'GET',
  headers: '',
  description: '',
  enabled: true
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '任务名称', dataIndex: 'name', key: 'name' },
  { title: 'Cron表达式', dataIndex: 'schedule', key: 'schedule' },
  { title: '命令', dataIndex: 'command', key: 'command', ellipsis: true },
  { title: '状态', key: 'enabled', width: 100 },
  { title: '操作', key: 'actions', width: 150 }
]

const httpMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH']

const taskFormRules = {
  name: [{ required: true, message: '请输入任务名称' }],
  schedule: [{ required: true, message: '请输入Cron表达式' }],
  command: [{ required: true, message: '请输入命令或URL' }]
}

const filteredTasks = computed(() => {
  if (!search.value) return tasks.value
  return tasks.value.filter(task =>
    task.name.toLowerCase().includes(search.value.toLowerCase()) ||
    task.command.toLowerCase().includes(search.value.toLowerCase())
  )
})

const fetchTasks = async () => {
  loading.value = true
  try {
    const response = await api.get('/tasks/paginated', {
      params: {
        page: page.value,
        page_size: itemsPerPage.value
      }
    })
    tasks.value = response.data.data
    totalItems.value = response.data.total
  } catch (error) {
    message.error('获取任务列表失败')
  } finally {
    loading.value = false
  }
}

const handleTableChange = (pagination) => {
  page.value = pagination.current
  itemsPerPage.value = pagination.pageSize
  fetchTasks()
}

const openTaskDialog = (task = null) => {
  editingTask.value = task
  if (task) {
    taskFormData.value = { ...task }
  } else {
    taskFormData.value = {
      name: '',
      schedule: '',
      command: '',
      method: 'GET',
      headers: '',
      description: '',
      enabled: true
    }
  }
  taskDialog.value = true
}

const saveTask = async () => {
  try {
    await taskFormRef.value.validate()
    saving.value = true
    
    if (editingTask.value) {
      await api.put(`/tasks/${editingTask.value.id}`, taskFormData.value)
      message.success('任务更新成功')
    } else {
      await api.post('/tasks', taskFormData.value)
      message.success('任务创建成功')
    }
    taskDialog.value = false
    fetchTasks()
  } catch (error) {
    if (error.response?.data?.message) {
      message.error(error.response.data.message)
    } else if (error.errorFields) {
      // 表单验证错误
      return
    } else {
      message.error('保存失败')
    }
  } finally {
    saving.value = false
  }
}

const updateTaskStatus = async (task) => {
  try {
    await api.put(`/tasks/${task.id}`, task)
    message.success('任务状态更新成功')
  } catch (error) {
    message.error('更新失败')
    // 恢复原状态
    task.enabled = !task.enabled
  }
}

const executeTask = async (taskId) => {
  try {
    await api.post(`/tasks/${taskId}/execute`)
    message.success('任务执行成功')
  } catch (error) {
    message.error('任务执行失败')
  }
}

const deleteTask = async (task) => {
  try {
    await api.delete(`/tasks/${task.id}`)
    message.success('任务删除成功')
    fetchTasks()
  } catch (error) {
    message.error('删除失败')
  }
}

onMounted(() => {
  fetchTasks()
})
</script>