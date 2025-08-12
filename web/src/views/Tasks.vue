<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <div class="d-flex justify-space-between align-center mb-4">
          <h1 class="text-h4">任务管理</h1>
          <v-btn
            color="primary"
            prepend-icon="mdi-plus"
            @click="openTaskDialog()"
          >
            新建任务
          </v-btn>
        </div>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-text-field
              v-model="search"
              prepend-inner-icon="mdi-magnify"
              label="搜索任务"
              single-line
              hide-details
              clearable
            ></v-text-field>
          </v-card-title>
          <v-data-table
            :headers="headers"
            :items="tasks"
            :search="search"
            :loading="loading"
            :page="page"
            :items-per-page="itemsPerPage"
            :server-items-length="totalItems"
            @update:page="handlePageChange"
            @update:items-per-page="handleItemsPerPageChange"
          >
            <template v-slot:item.enabled="{ item }">
              <v-switch
                v-model="item.enabled"
                @change="updateTaskStatus(item)"
                hide-details
              ></v-switch>
            </template>
            <template v-slot:item.actions="{ item }">
              <v-btn
                icon="mdi-play"
                size="small"
                color="primary"
                @click="executeTask(item.id)"
                class="mr-2"
              ></v-btn>
              <v-btn
                icon="mdi-pencil"
                size="small"
                color="info"
                @click="openTaskDialog(item)"
                class="mr-2"
              ></v-btn>
              <v-btn
                icon="mdi-delete"
                size="small"
                color="error"
                @click="deleteTask(item)"
              ></v-btn>
            </template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>

    <!-- 任务编辑对话框 -->
    <v-dialog v-model="taskDialog" max-width="600px">
      <v-card>
        <v-card-title>
          {{ editingTask ? '编辑任务' : '新建任务' }}
        </v-card-title>
        <v-card-text>
          <v-form ref="taskForm" v-model="taskFormValid">
            <v-text-field
              v-model="taskFormData.name"
              label="任务名称"
              :rules="[v => !!v || '任务名称不能为空']"
              required
            ></v-text-field>
            
            <v-text-field
              v-model="taskFormData.schedule"
              label="Cron 表达式"
              :rules="[v => !!v || 'Cron表达式不能为空']"
              hint="例如: 0 */5 * * * * (每5分钟执行一次)"
              persistent-hint
              required
            ></v-text-field>
            
            <v-text-field
              v-model="taskFormData.command"
              label="命令或URL"
              :rules="[v => !!v || '命令不能为空']"
              required
            ></v-text-field>
            
            <v-select
              v-model="taskFormData.method"
              label="HTTP方法"
              :items="httpMethods"
              hint="仅当命令为URL时需要"
              persistent-hint
            ></v-select>
            
            <v-textarea
              v-model="taskFormData.headers"
              label="请求头 (JSON格式)"
              hint='例如: {"Content-Type": "application/json"}'
              persistent-hint
              rows="3"
            ></v-textarea>
            
            <v-textarea
              v-model="taskFormData.description"
              label="任务描述"
              rows="2"
            ></v-textarea>
            
            <v-switch
              v-model="taskFormData.enabled"
              label="启用任务"
            ></v-switch>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="taskDialog = false">取消</v-btn>
          <v-btn
            color="primary"
            :disabled="!taskFormValid"
            :loading="saving"
            @click="saveTask"
          >
            保存
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar
      v-model="snackbar"
      :color="snackbarColor"
      timeout="3000"
    >
      {{ snackbarText }}
    </v-snackbar>
  </v-container>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import api from '../services/api'

const loading = ref(false)
const saving = ref(false)
const search = ref('')
const tasks = ref([])
const page = ref(1)
const itemsPerPage = ref(10)
const totalItems = ref(0)

const taskDialog = ref(false)
const taskForm = ref(null)
const taskFormValid = ref(false)
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

const snackbar = ref(false)
const snackbarText = ref('')
const snackbarColor = ref('success')

const headers = [
  { title: 'ID', key: 'id', width: '80px' },
  { title: '任务名称', key: 'name' },
  { title: 'Cron表达式', key: 'schedule' },
  { title: '命令', key: 'command' },
  { title: '状态', key: 'enabled', width: '100px' },
  { title: '操作', key: 'actions', sortable: false, width: '150px' }
]

const httpMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH']

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
    showSnackbar('获取任务列表失败', 'error')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (newPage) => {
  page.value = newPage
  fetchTasks()
}

const handleItemsPerPageChange = (newItemsPerPage) => {
  itemsPerPage.value = newItemsPerPage
  page.value = 1
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
  if (!taskForm.value.validate()) return

  saving.value = true
  try {
    if (editingTask.value) {
      await api.put(`/tasks/${editingTask.value.id}`, taskFormData.value)
      showSnackbar('任务更新成功', 'success')
    } else {
      await api.post('/tasks', taskFormData.value)
      showSnackbar('任务创建成功', 'success')
    }
    taskDialog.value = false
    fetchTasks()
  } catch (error) {
    showSnackbar(error.response?.data?.message || '保存失败', 'error')
  } finally {
    saving.value = false
  }
}

const updateTaskStatus = async (task) => {
  try {
    await api.put(`/tasks/${task.id}`, task)
    showSnackbar('任务状态更新成功', 'success')
  } catch (error) {
    showSnackbar('更新失败', 'error')
    // 恢复原状态
    task.enabled = !task.enabled
  }
}

const executeTask = async (taskId) => {
  try {
    await api.post(`/tasks/${taskId}/execute`)
    showSnackbar('任务执行成功', 'success')
  } catch (error) {
    showSnackbar('任务执行失败', 'error')
  }
}

const deleteTask = async (task) => {
  if (!confirm(`确定要删除任务 "${task.name}" 吗？`)) return

  try {
    await api.delete(`/tasks/${task.id}`)
    showSnackbar('任务删除成功', 'success')
    fetchTasks()
  } catch (error) {
    showSnackbar('删除失败', 'error')
  }
}

const showSnackbar = (text, color = 'success') => {
  snackbarText.value = text
  snackbarColor.value = color
  snackbar.value = true
}

onMounted(() => {
  fetchTasks()
})
</script>