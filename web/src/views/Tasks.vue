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

        <!-- 通知配置 -->
        <a-divider>通知配置</a-divider>
        
        <a-form-item label="通知时机">
          <a-checkbox-group v-model:value="notifyTiming">
            <a-checkbox value="success">执行成功时通知</a-checkbox>
            <a-checkbox value="failure">执行失败时通知</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        
        <a-form-item label="通知方式" v-if="notifyTiming.length > 0">
          <a-checkbox-group v-model:value="notificationTypes">
            <a-checkbox value="email">邮件通知</a-checkbox>
            <a-checkbox value="dingtalk">钉钉通知</a-checkbox>
            <a-checkbox value="wechat">企业微信通知</a-checkbox>
          </a-checkbox-group>
        </a-form-item>

        <!-- 邮件配置 -->
        <a-collapse v-if="notificationTypes.includes('email')" style="margin-bottom: 16px;">
          <a-collapse-panel key="email" header="邮件通知配置">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="SMTP服务器">
                  <a-input v-model:value="emailConfig.smtp_host" placeholder="smtp.gmail.com" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="SMTP端口">
                  <a-input-number v-model:value="emailConfig.smtp_port" :min="1" :max="65535" placeholder="587" style="width: 100%" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="用户名">
                  <a-input v-model:value="emailConfig.username" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="密码">
                  <a-input-password v-model:value="emailConfig.password" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item label="发件人">
              <a-input v-model:value="emailConfig.from" />
            </a-form-item>
            <a-form-item label="收件人">
              <a-select
                v-model:value="emailConfig.to"
                mode="tags"
                placeholder="输入邮箱地址，按回车添加"
                style="width: 100%"
              />
            </a-form-item>
            <a-form-item label="邮件主题">
              <a-input v-model:value="emailConfig.subject" placeholder="留空使用默认主题" />
            </a-form-item>
            <a-form-item>
              <a-checkbox v-model:checked="emailConfig.enable_tls">启用TLS</a-checkbox>
            </a-form-item>
          </a-collapse-panel>
        </a-collapse>

        <!-- 钉钉配置 -->
        <a-collapse v-if="notificationTypes.includes('dingtalk')" style="margin-bottom: 16px;">
          <a-collapse-panel key="dingtalk" header="钉钉通知配置">
            <a-form-item label="Webhook URL">
              <a-input v-model:value="dingtalkConfig.webhook_url" placeholder="https://oapi.dingtalk.com/robot/send?access_token=..." />
            </a-form-item>
            <a-form-item label="签名密钥">
              <a-input v-model:value="dingtalkConfig.secret" placeholder="可选，用于签名验证" />
            </a-form-item>
            <a-form-item label="@手机号">
              <a-select
                v-model:value="dingtalkConfig.at_mobiles"
                mode="tags"
                placeholder="输入手机号，按回车添加"
                style="width: 100%"
              />
            </a-form-item>
            <a-form-item>
              <a-checkbox v-model:checked="dingtalkConfig.at_all">@所有人</a-checkbox>
            </a-form-item>
          </a-collapse-panel>
        </a-collapse>

        <!-- 企业微信配置 -->
        <a-collapse v-if="notificationTypes.includes('wechat')" style="margin-bottom: 16px;">
          <a-collapse-panel key="wechat" header="企业微信通知配置">
            <a-form-item label="Webhook URL">
              <a-input v-model:value="wechatConfig.webhook_url" placeholder="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=..." />
            </a-form-item>
            <a-form-item label="@用户ID">
              <a-select
                v-model:value="wechatConfig.at_user_ids"
                mode="tags"
                placeholder="输入用户ID，按回车添加"
                style="width: 100%"
              />
            </a-form-item>
            <a-form-item>
              <a-checkbox v-model:checked="wechatConfig.at_all">@所有人</a-checkbox>
            </a-form-item>
          </a-collapse-panel>
        </a-collapse>

        <!-- 测试通知按钮 -->
        <a-form-item v-if="notificationTypes.length > 0">
          <a-button @click="testNotification" :loading="testingNotification">
            测试通知
          </a-button>
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
  enabled: true,
  notify_on_success: false,
  notify_on_failure: true,
  notification_types: '',
  notification_config: ''
})

// 通知相关的响应式数据
const testingNotification = ref(false)
const notifyTiming = ref([])
const notificationTypes = ref([])
const emailConfig = ref({
  smtp_host: '',
  smtp_port: 587,
  username: '',
  password: '',
  from: '',
  to: [],
  subject: '',
  enable_tls: true
})
const dingtalkConfig = ref({
  webhook_url: '',
  secret: '',
  at_mobiles: [],
  at_all: false
})
const wechatConfig = ref({
  webhook_url: '',
  at_user_ids: [],
  at_all: false
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
    
    // 解析通知配置
    parseNotificationConfig(task)
  } else {
    taskFormData.value = {
      name: '',
      schedule: '',
      command: '',
      method: 'GET',
      headers: '',
      description: '',
      enabled: true,
      notify_on_success: false,
      notify_on_failure: true,
      notification_types: '',
      notification_config: ''
    }
    
    // 重置通知配置
    resetNotificationConfig()
  }
  taskDialog.value = true
}

// 解析通知配置
const parseNotificationConfig = (task) => {
  // 解析通知时机
  notifyTiming.value = []
  if (task.notify_on_success) notifyTiming.value.push('success')
  if (task.notify_on_failure) notifyTiming.value.push('failure')
  
  // 解析通知类型
  try {
    notificationTypes.value = task.notification_types ? JSON.parse(task.notification_types) : []
  } catch (e) {
    notificationTypes.value = []
  }
  
  // 解析通知配置
  try {
    const config = task.notification_config ? JSON.parse(task.notification_config) : {}
    emailConfig.value = config.email || {
      smtp_host: '',
      smtp_port: 587,
      username: '',
      password: '',
      from: '',
      to: [],
      subject: '',
      enable_tls: true
    }
    dingtalkConfig.value = config.dingtalk || {
      webhook_url: '',
      secret: '',
      at_mobiles: [],
      at_all: false
    }
    wechatConfig.value = config.wechat || {
      webhook_url: '',
      at_user_ids: [],
      at_all: false
    }
  } catch (e) {
    resetNotificationConfig()
  }
}

// 重置通知配置
const resetNotificationConfig = () => {
  notifyTiming.value = []
  notificationTypes.value = []
  emailConfig.value = {
    smtp_host: '',
    smtp_port: 587,
    username: '',
    password: '',
    from: '',
    to: [],
    subject: '',
    enable_tls: true
  }
  dingtalkConfig.value = {
    webhook_url: '',
    secret: '',
    at_mobiles: [],
    at_all: false
  }
  wechatConfig.value = {
    webhook_url: '',
    at_user_ids: [],
    at_all: false
  }
}

const saveTask = async () => {
  try {
    await taskFormRef.value.validate()
    saving.value = true
    
    // 构建通知配置
    buildNotificationConfig()
    
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

// 构建通知配置
const buildNotificationConfig = () => {
  // 设置通知时机
  taskFormData.value.notify_on_success = notifyTiming.value.includes('success')
  taskFormData.value.notify_on_failure = notifyTiming.value.includes('failure')
  
  // 设置通知类型
  taskFormData.value.notification_types = JSON.stringify(notificationTypes.value)
  
  // 构建通知配置
  const config = {}
  if (notificationTypes.value.includes('email')) {
    config.email = { ...emailConfig.value }
  }
  if (notificationTypes.value.includes('dingtalk')) {
    config.dingtalk = { ...dingtalkConfig.value }
  }
  if (notificationTypes.value.includes('wechat')) {
    config.wechat = { ...wechatConfig.value }
  }
  
  taskFormData.value.notification_config = JSON.stringify(config)
}

// 测试通知
const testNotification = async () => {
  if (notificationTypes.value.length === 0) {
    message.warning('请先选择通知方式')
    return
  }
  
  testingNotification.value = true
  try {
    const config = {}
    if (notificationTypes.value.includes('email')) {
      config.email = { ...emailConfig.value }
    }
    if (notificationTypes.value.includes('dingtalk')) {
      config.dingtalk = { ...dingtalkConfig.value }
    }
    if (notificationTypes.value.includes('wechat')) {
      config.wechat = { ...wechatConfig.value }
    }
    
    await api.post('/notifications/test', {
      notification_types: notificationTypes.value,
      notification_config: config
    })
    
    message.success('测试通知已发送，请检查相应的通知渠道')
  } catch (error) {
    message.error('测试通知发送失败')
  } finally {
    testingNotification.value = false
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