# Vuetify 到 Ant Design Vue 迁移指南

## 已完成的迁移工作

### 1. 依赖更新
- 移除了 `vuetify`, `@mdi/font`, `vite-plugin-vuetify`
- 添加了 `ant-design-vue`, `@ant-design/icons-vue`

### 2. 主要文件更新
- `main.js`: 替换 Vuetify 为 Ant Design Vue
- `vite.config.js`: 移除 Vuetify 插件
- `App.vue`: 完全重写布局，使用 Ant Design 的 Layout 组件
- 所有视图文件已更新：
  - `Dashboard.vue`: 使用 Ant Design 的统计卡片、表格、进度条
  - `Tasks.vue`: 使用 Ant Design 的表格、模态框、表单
  - `Logs.vue`: 使用 Ant Design 的表格、模态框、描述列表
  - `System.vue`: 使用 Ant Design 的卡片、进度条、统计组件
  - `Login.vue`: 使用 Ant Design 的表单、输入框
  - `Register.vue`: 使用 Ant Design 的表单、输入框

### 3. 组件映射
| Vuetify 组件 | Ant Design Vue 组件 |
|-------------|-------------------|
| v-app | a-layout |
| v-navigation-drawer | a-layout-sider |
| v-app-bar | a-layout-header |
| v-main | a-layout-content |
| v-card | a-card |
| v-btn | a-button |
| v-text-field | a-input |
| v-select | a-select |
| v-data-table | a-table |
| v-dialog | a-modal |
| v-form | a-form |
| v-snackbar | message (全局方法) |
| v-chip | a-tag |
| v-progress-linear | a-progress |
| v-switch | a-switch |

## 安装步骤

1. 删除旧的 node_modules 和 package-lock.json：
```bash
cd web
rm -rf node_modules package-lock.json
```

2. 安装新的依赖：
```bash
npm install
```

3. 启动开发服务器：
```bash
npm run dev
```

## 主要变化

### 1. 主题系统
- Vuetify 的主题系统已移除
- 可以通过 Ant Design 的 ConfigProvider 来自定义主题

### 2. 图标系统
- 从 MDI 图标切换到 Ant Design 图标
- 需要从 `@ant-design/icons-vue` 导入具体的图标组件

### 3. 消息提示
- 从 v-snackbar 改为使用 `message` 全局方法
- 更简洁的 API，无需管理状态

### 4. 表单验证
- Ant Design 的表单验证更加强大
- 支持异步验证和更复杂的验证规则

### 5. 响应式布局
- 使用 Ant Design 的栅格系统 (a-row, a-col)
- 更灵活的断点配置

## 注意事项

1. 确保后端 API 接口正常工作
2. 检查所有路由是否正确配置
3. 测试所有功能是否正常
4. 如需自定义主题，请参考 Ant Design Vue 文档

## 可能需要的额外配置

如果需要自定义主题，可以在 `main.js` 中添加：

```javascript
import { ConfigProvider } from 'ant-design-vue'

// 在创建应用时配置主题
app.use(ConfigProvider, {
  theme: {
    token: {
      colorPrimary: '#1890ff',
      // 其他主题配置
    }
  }
})
```