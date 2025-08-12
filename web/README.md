# Crontab Go 前端界面

基于 Vue 3 + Vuetify 3 开发的定时任务管理系统前端界面。

## 功能特性

- 🔐 用户认证（登录/注册）
- 📊 仪表板概览
- ⏰ 任务管理（增删改查）
- 📝 执行日志查看
- 📈 系统监控
- 🌙 深色/浅色主题切换
- 📱 响应式设计

## 技术栈

- Vue 3 - 渐进式 JavaScript 框架
- Vuetify 3 - Material Design 组件库
- Vue Router - 路由管理
- Pinia - 状态管理
- Axios - HTTP 客户端
- Vite - 构建工具

## 安装依赖

```bash
cd web
npm install
```

## 开发运行

```bash
npm run dev
```

访问 http://localhost:3000

## 构建生产版本

```bash
npm run build
```

## 项目结构

```
web/
├── src/
│   ├── components/          # 公共组件
│   ├── views/              # 页面组件
│   │   ├── Login.vue       # 登录页面
│   │   ├── Register.vue    # 注册页面
│   │   ├── Dashboard.vue   # 仪表板
│   │   ├── Tasks.vue       # 任务管理
│   │   ├── Logs.vue        # 执行日志
│   │   └── System.vue      # 系统监控
│   ├── stores/             # Pinia 状态管理
│   │   └── user.js         # 用户状态
│   ├── services/           # API 服务
│   │   └── api.js          # HTTP 客户端配置
│   ├── router/             # 路由配置
│   │   └── index.js
│   ├── App.vue             # 根组件
│   └── main.js             # 入口文件
├── index.html              # HTML 模板
├── vite.config.js          # Vite 配置
└── package.json            # 项目配置
```

## API 代理配置

开发环境下，Vite 会将 `/api` 路径的请求代理到后端服务器（默认 http://localhost:8080）。

如需修改后端地址，请编辑 `vite.config.js` 文件中的 proxy 配置。

## 页面说明

### 登录页面 (`/login`)
- 用户名/密码登录
- 跳转到注册页面

### 注册页面 (`/register`)
- 用户注册功能
- 表单验证

### 仪表板 (`/dashboard`)
- 系统概览统计
- 最近任务列表
- 系统状态监控

### 任务管理 (`/tasks`)
- 任务列表（分页）
- 创建/编辑任务
- 启用/禁用任务
- 立即执行任务
- 删除任务

### 执行日志 (`/logs`)
- 任务执行日志查看
- 按任务筛选
- 按状态筛选
- 日志详情查看

### 系统监控 (`/system`)
- CPU、内存、磁盘使用率
- 网络流量统计
- 进程信息
- 系统运行时间

## 默认账号

根据 API 文档，默认管理员账号：
- 用户名: admin
- 密码: admin123

## 注意事项

1. 确保后端 API 服务已启动并运行在 8080 端口
2. 所有 API 请求都需要 JWT 认证（除登录和注册接口）
3. Token 会自动存储在 localStorage 中
4. 页面会自动处理认证状态和路由守卫