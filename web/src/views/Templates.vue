<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px;">
      <h1>任务模板</h1>
      <a-space>
        <a-button type="primary" @click="openTemplateDialog()">
          <PlusOutlined />
          新建模板
        </a-button>
        <a-button @click="refreshData" :loading="loading">
          <ReloadOutlined />
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 搜索和筛选 -->
    <a-card style="margin-bottom: 24px;">
      <a-row :gutter="16">
        <a-col :span="6">
          <a-input-search
            v-model:value="searchForm.keyword"
            placeholder="搜索模板名称或描述"
            @search="handleSearch"
            allow-clear
          />
        </a-col>
        <a-col :span="4">
          <a-select
            v-model:value="searchForm.category"
            placeholder="选择分类"
            allow-clear
            @change="handleSearch"
          >
            <a-select-option
              v-for="category in categories"
              :key="category.name"
              :value="category.name"
            >
              {{ category.description }}
            </a-select-option>
          </a-select>
        </a-col>
        <a-col :span="4">
          <a-select
            v-model:value="searchForm.scope"
            placeholder="模板范围"
            @change="handleSearch"
          >
            <a-select-option value="">全部</a-select-option>
            <a-select-option value="public">公共模板</a-select-option>
            <a-select-option value="my">我的模板</a-select-option>
          </a-select>
        </a-col>
        <a-col :span="4">
          <a-select
            v-model:value="searchForm.sortBy"
            placeholder="排序方式"
            @change="handleSearch"
          >
            <a-select-option value="usage">使用次数</a-select-option>
            <a-select-option value="created">创建时间</a-select-option>
            <a-select-option value="name">名称</a-select-option>
          </a-select>
        </a-col>
      </a-row>
    </a-card>

    <!-- 热门模板 -->
    <a-card title="热门模板" style="margin-bottom: 24px;" v-if="popularTemplates.length > 0">
      <a-row :gutter="16">
        <a-col :span="6" v-for="template in popularTemplates" :key="template.id">
          <a-card
            size="small"
            :title="template.name"
            style="margin-bottom: 16px; cursor: pointer;"
            @click="useTemplate(template)"
            hoverable
          >
            <template #extra>
              <a-tag color="blue">{{ template.usage_count }} 次使用</a-tag>
            </template>
            <p style="margin: 0; color: #666; font-size: 12px;">{{ template.description }}</p>
            <a-tag style="margin-top: 8px;">{{ template.category }}</a-tag>
          </a-card>
        </a-col>
      </a-row>
    </a-card>

    <!-- 模板列表 -->
    <a-card title="模板列表">
      <a-table
        :columns="columns"
        :data-source="templates"
        :loading="loading"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: pagination.total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`
        }"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <div>
              <strong>{{ record.name }}</strong>
              <div style="color: #666; font-size: 12px;">{{ record.description }}</div>
            </div>
          </template>
          <template v-else-if="column.key === 'category'">
            <a-tag :color="getCategoryColor(record.category)">
              {{ getCategoryName(record.category) }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'tags'">
            <a-tag
              v-for="tag in parseTagList(record.tag_list)"
              :key="tag"
              style="margin-bottom: 4px;"
            >
              {{ tag }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'is_public'">
            <a-tag :color="record.is_public ? 'green' : 'orange'">
              {{ record.is_public ? '公共' : '私有' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'usage_count'">
            <a-statistic
              :value="record.usage_count"
              :value-style="{ fontSize: '14px' }"
            />
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button
                type="primary"
                size="small"
                @click="useTemplate(record)"
              >
                <ThunderboltOutlined />
                使用
              </a-button>
              <a-button
                size="small"
                @click="openTemplateDialog(record)"
              >
                <EditOutlined />
                编辑
              </a-button>
              <a-popconfirm
                title="确定要删除这个模板吗？"
                @confirm="deleteTemplate(record)"
              >
                <a-button
                  danger
                  size="small"
                >
                  <DeleteOutlined />
                  删除
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 模板编辑对话框 -->
    <a-modal
      v-model:open="templateDialog"
      :title="editingTemplate ? '编辑模板' : '新建模板'"
      width="800px"
      @ok="saveTemplate"
      @cancel="templateDialog = false"
      :confirm-loading="saving"
    >
      <a-form
        ref="templateFormRef"
        :model="templateFormData"
        :rules="templateFormRules"
        layout="vertical"
      >
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="模板名称" name="name">
              <a-input v-model:value="templateFormData.name" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="分类" name="category">
              <a-select v-model:value="templateFormData.category">
                <a-select-option
                  v-for="category in categories"
                  :key="category.name"
                  :value="category.name"
                >
                  {{ category.description }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="模板描述" name="description">
          <a-textarea
            v-model:value="templateFormData.description"
            :rows="2"
          />
        </a-form-item>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="Cron 表达式" name="schedule">
              <a-input
                v-model:value="templateFormData.schedule"
                placeholder="例如: 0 */5 * * * * (每5分钟执行一次)"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="HTTP方法" v-if="isHttpCommand">
              <a-select v-model:value="templateFormData.method">
                <a-select-option
                  v-for="method in httpMethods"
                  :key="method"
                  :value="method"
                >
                  {{ method }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="命令或URL" name="command">
          <a-input v-model:value="templateFormData.command" />
        </a-form-item>

        <a-form-item label="请求头 (JSON格式)" v-if="isHttpCommand">
          <a-textarea
            v-model:value="templateFormData.headers"
            :rows="3"
            placeholder='例如: {"Content-Type": "application/json"}'
          />
        </a-form-item>

        <a-form-item label="标签">
          <a-select
            v-model:value="templateTags"
            mode="tags"
            placeholder="输入标签，按回车添加"
            style="width: 100%"
          />
        </a-form-item>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item>
              <a-checkbox v-model:checked="templateFormData.is_public">
                公共模板（其他用户可见）
              </a-checkbox>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item>
              <a-checkbox v-model:checked="templateFormData.notify_on_success">
                成功时通知
              </a-checkbox>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <!-- 使用模板对话框 -->
    <a-modal
      v-model:open="useTemplateDialog"
      title="使用模板创建任务"
      width="600px"
      @ok="createTaskFromTemplate"
      @cancel="useTemplateDialog = false"
      :confirm-loading="creating"
    >
      <a-form
        ref="useTemplateFormRef"
        :model="useTemplateFormData"
        :rules="useTemplateFormRules"
        layout="vertical"
      >
        <a-form-item label="任务名称" name="task_name">
          <a-input v-model:value="useTemplateFormData.task_name" />
        </a-form-item>

        <a-form-item>
          <a-checkbox v-model:checked="useTemplateFormData.enabled">
            创建后立即启用
          </a-checkbox>
        </a-form-item>

        <a-divider>模板信息</a-divider>
        <div v-if="selectedTemplate">
          <p><strong>模板名称：</strong>{{ selectedTemplate.name }}</p>
          <p><strong>描述：</strong>{{ selectedTemplate.description }}</p>
          <p><strong>Cron表达式：</strong>{{ selectedTemplate.schedule }}</p>
          <p><strong>命令：</strong>{{ selectedTemplate.command }}</p>
        </div>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import {
  PlusOutlined,
  ReloadOutlined,
  ThunderboltOutlined,
  EditOutlined,
  DeleteOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import api from '../services/api'

const loading = ref(false)
const saving = ref(false)
const creating = ref(false)

// 数据
const templates = ref([])
const categories = ref([])
const popularTemplates = ref([])
const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0
})

// 搜索表单
const searchForm = ref({
  keyword: '',
  category: '',
  scope: '',
  sortBy: 'usage'
})

// 模板对话框
const templateDialog = ref(false)
const templateFormRef = ref(null)
const editingTemplate = ref(null)
const templateFormData = ref({
  name: '',
  description: '',
  category: 'general',
  schedule: '',
  command: '',
  method: 'GET',
  headers: '',
  is_public: false,
  notify_on_success: false,
  notify_on_failure: true
})
const templateTags = ref([])

// 使用模板对话框
const useTemplateDialog = ref(false)
const useTemplateFormRef = ref(null)
const selectedTemplate = ref(null)
const useTemplateFormData = ref({
  template_id: 0,
  task_name: '',
  enabled: true
})

// 表格列定义
const columns = [
  { title: '模板名称', key: 'name', width: 250 },
  { title: '分类', key: 'category', width: 100 },
  { title: '标签', key: 'tags', width: 150 },
  { title: '类型', key: 'is_public', width: 80 },
  { title: '使用次数', key: 'usage_count', width: 100 },
  { title: '创建者', dataIndex: 'creator_name', key: 'creator_name', width: 100 },
  { title: '操作', key: 'actions', width: 200, fixed: 'right' }
]

const httpMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH']

// 表单验证规则
const templateFormRules = {
  name: [{ required: true, message: '请输入模板名称' }],
  schedule: [{ required: true, message: '请输入Cron表达式' }],
  command: [{ required: true, message: '请输入命令或URL' }]
}

const useTemplateFormRules = {
  task_name: [{ required: true, message: '请输入任务名称' }]
}

// 计算属性
const isHttpCommand = computed(() => {
  return templateFormData.value.command.startsWith('http://') || 
         templateFormData.value.command.startsWith('https://')
})

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.value.current,
      page_size: pagination.value.pageSize,
      ...searchForm.value
    }

    // 处理搜索范围
    if (searchForm.value.scope === 'public') {
      params.is_public = true
    } else if (searchForm.value.scope === 'my') {
      delete params.is_public
      // 使用专门的接口获取我的模板
    }

    const response = await api.get('/templates/search', { params })
    templates.value = response.data.data || []
    pagination.value.total = response.data.total || 0
  } catch (error) {
    message.error('获取模板列表失败')
  } finally {
    loading.value = false
  }
}

// 获取分类
const fetchCategories = async () => {
  try {
    const response = await api.get('/template-categories')
    categories.value = response.data || []
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

// 获取热门模板
const fetchPopularTemplates = async () => {
  try {
    const response = await api.get('/templates/popular', { params: { limit: 4 } })
    popularTemplates.value = response.data || []
  } catch (error) {
    console.error('获取热门模板失败:', error)
  }
}

// 刷新数据
const refreshData = () => {
  fetchData()
  fetchPopularTemplates()
}

// 搜索处理
const handleSearch = () => {
  pagination.value.current = 1
  fetchData()
}

// 表格变化处理
const handleTableChange = (pag) => {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchData()
}

// 打开模板对话框
const openTemplateDialog = (template = null) => {
  editingTemplate.value = template
  if (template) {
    templateFormData.value = { ...template }
    templateTags.value = parseTagList(template.tag_list) || []
  } else {
    templateFormData.value = {
      name: '',
      description: '',
      category: 'general',
      schedule: '',
      command: '',
      method: 'GET',
      headers: '',
      is_public: false,
      notify_on_success: false,
      notify_on_failure: true
    }
    templateTags.value = []
  }
  templateDialog.value = true
}

// 保存模板
const saveTemplate = async () => {
  try {
    await templateFormRef.value.validate()
    saving.value = true

    // 构建标签
    templateFormData.value.tags = JSON.stringify(templateTags.value)

    if (editingTemplate.value) {
      await api.put(`/templates/${editingTemplate.value.id}`, templateFormData.value)
      message.success('模板更新成功')
    } else {
      await api.post('/templates', templateFormData.value)
      message.success('模板创建成功')
    }

    templateDialog.value = false
    fetchData()
  } catch (error) {
    if (error.response?.data?.error) {
      message.error(error.response.data.error)
    } else if (error.errorFields) {
      return
    } else {
      message.error('保存失败')
    }
  } finally {
    saving.value = false
  }
}

// 删除模板
const deleteTemplate = async (template) => {
  try {
    await api.delete(`/templates/${template.id}`)
    message.success('模板删除成功')
    fetchData()
  } catch (error) {
    message.error('删除失败')
  }
}

// 使用模板
const useTemplate = (template) => {
  selectedTemplate.value = template
  useTemplateFormData.value = {
    template_id: template.id,
    task_name: `${template.name}_${Date.now()}`,
    enabled: true
  }
  useTemplateDialog.value = true
}

// 从模板创建任务
const createTaskFromTemplate = async () => {
  try {
    await useTemplateFormRef.value.validate()
    creating.value = true

    await api.post('/templates/create-task', useTemplateFormData.value)
    message.success('任务创建成功')
    useTemplateDialog.value = false
  } catch (error) {
    if (error.response?.data?.error) {
      message.error(error.response.data.error)
    } else {
      message.error('创建任务失败')
    }
  } finally {
    creating.value = false
  }
}

// 辅助方法
const getCategoryColor = (category) => {
  const colorMap = {
    general: 'blue',
    backup: 'green',
    monitoring: 'orange',
    cleanup: 'red',
    notification: 'purple',
    api: 'cyan'
  }
  return colorMap[category] || 'default'
}

const getCategoryName = (category) => {
  const category_obj = categories.value.find(c => c.name === category)
  return category_obj ? category_obj.description : category
}

const parseTagList = (tagList) => {
  if (Array.isArray(tagList)) {
    return tagList
  }
  if (typeof tagList === 'string') {
    try {
      return JSON.parse(tagList)
    } catch {
      return []
    }
  }
  return []
}

onMounted(() => {
  fetchCategories()
  fetchData()
  fetchPopularTemplates()
})
</script>

<style scoped>
.ant-card {
  margin-bottom: 16px;
}

.ant-statistic {
  text-align: center;
}
</style>