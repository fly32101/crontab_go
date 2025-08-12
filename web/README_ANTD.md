# Crontab Go 前端 - Ant Design Vue 版本

## 项目概述

这是 Crontab Go 管理系统的前端界面，已从 Vuetify 成功迁移到 Ant Design Vue。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **Ant Design Vue 4.x** - 企业级 UI 设计语言和组件库
- **Vue Router 4** - Vue.js 官方路由管理器
- **Pinia** - Vue 状态管理库
- **Axios** - HTTP 客户端
- **Vite** - 现代前端构建工具

## 功能特性

### 🎨 现代化 UI 设计
- 采用 Ant Design 设计规范
- 响应式布局，支持多种屏幕尺寸
- 深色/浅色主题切换（可扩展）

### 📊 仪表板
- 系统概览统计
- 任务执行状态监控
- 系统资源使用情况
- 实时数据更新

### ⚙️ 任务管理
- 创建、编辑、删除定时任务
- 支持 Cron 表达式配置
- 任务启用/禁用切换
- 手动执行任务
- 分页和搜索功能

### 📝 执行日志
- 查看任务执行历史
- 按任务和状态筛选
- 详细的执行信息展示
- 错误日志查看

### 🖥️ 系统监控
- CPU、内存、磁盘使用率
- 网络流量统计
- 进程和协程信息
- 系统运行时间
- 自动刷新监控数据

### 👤 用户管理
- 用户注册和登录
- 会话管理
- 路由守卫

## 项目结构

```
web/
├── src/
│   ├── components/          # 可复用组件
│   ├── views/              # 页面组件
│   │   ├── Dashboard.vue   # 仪表板
│   │   ├── Tasks.vue       # 任务管理
│   │   ├── Logs.vue        # 执行日志
│   │   ├── System.vue      # 系统监控
│   │   ├── Login.vue       # 登录页面
│   │   └── Register.vue    # 注册页面
│   ├── router/             # 路由配置
│   ├── stores/             # 状态管理
│   ├── services/           # API 服务
│   ├── App.vue             # 根组件
│   └── main.js             # 应用入口
├── public/                 # 静态资源
├── package.json            # 项目配置
├── vite.config.js          # Vite 配置
└── README_ANTD.md          # 项目说明
```

## 开发指南

### 安装依赖
```bash
cd web
npm install
```

### 启动开发服务器
```bash
npm run dev
```

### 构建生产版本
```bash
npm run build
```

### 预览生产构建
```bash
npm run preview
```

## 组件使用示例

### 表格组件
```vue
<a-table
  :columns="columns"
  :data-source="dataSource"
  :loading="loading"
  :pagination="pagination"
  @change="handleTableChange"
>
  <template #bodyCell="{ column, record }">
    <template v-if="column.key === 'actions'">
      <a-button @click="handleEdit(record)">编辑</a-button>
    </template>
  </template>
</a-table>
```

### 表单组件
```vue
<a-form
  :model="formData"
  :rules="formRules"
  @finish="handleSubmit"
>
  <a-form-item label="名称" name="name">
    <a-input v-model:value="formData.name" />
  </a-form-item>
  <a-form-item>
    <a-button type="primary" html-type="submit">提交</a-button>
  </a-form-item>
</a-form>
```

### 消息提示
```javascript
import { message } from 'ant-design-vue'

// 成功消息
message.success('操作成功')

// 错误消息
message.error('操作失败')

// 警告消息
message.warning('请注意')
```

## API 接口

前端通过 Axios 与后端 API 通信，主要接口包括：

- `GET /api/tasks` - 获取任务列表
- `POST /api/tasks` - 创建任务
- `PUT /api/tasks/:id` - 更新任务
- `DELETE /api/tasks/:id` - 删除任务
- `POST /api/tasks/:id/execute` - 执行任务
- `GET /api/logs` - 获取执行日志
- `GET /api/system/stats` - 获取系统状态

## 部署说明

1. 构建生产版本：
```bash
npm run build
```

2. 将 `dist` 目录中的文件部署到 Web 服务器

3. 配置反向代理，将 `/api` 请求转发到后端服务

### Nginx 配置示例
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    root /path/to/dist;
    index index.html;
    
    # 前端路由
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # API 代理
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 自定义主题

可以通过修改 `main.js` 来自定义 Ant Design 主题：

```javascript
import { ConfigProvider } from 'ant-design-vue'

app.use(ConfigProvider, {
  theme: {
    token: {
      colorPrimary: '#1890ff',
      borderRadius: 6,
      // 更多主题配置...
    }
  }
})
```

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证。