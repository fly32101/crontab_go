<template>
  <a-layout style="min-height: 100vh">
    <a-layout-sider
      v-model:collapsed="collapsed"
      :trigger="null"
      collapsible
      :breakpoint="'lg'"
      @breakpoint="onBreakpoint"
    >
      <div class="user-info">
        <a-avatar
          :size="collapsed ? 32 : 64"
          src="http://q.qlogo.cn/headimg_dl?dst_uin=944219401&spec=640&img_type=jpg"
        />
        <div v-if="!collapsed" class="user-details">
          <div class="username">{{ userStore.user?.username || '未登录' }}</div>
          <div class="email">{{ userStore.user?.email || '' }}</div>
        </div>
      </div>
      
      <a-divider style="margin: 16px 0" />
      
      <a-menu
        v-model:selectedKeys="selectedKeys"
        theme="dark"
        mode="inline"
        @click="handleMenuClick"
      >
        <a-menu-item
          v-for="item in menuItems"
          :key="item.value"
        >
          <component :is="item.icon" />
          <span>{{ item.title }}</span>
        </a-menu-item>
      </a-menu>

      <div class="logout-btn" :style="{ padding: collapsed ? '8px' : '16px' }">
        <a-button
          type="primary"
          danger
          :block="!collapsed"
          @click="logout"
        >
          <LogoutOutlined />
          <span v-if="!collapsed">退出登录</span>
        </a-button>
      </div>
    </a-layout-sider>

    <a-layout>
      <a-layout-header style="background: #fff; padding: 0; display: flex; align-items: center; justify-content: space-between;">
        <div style="display: flex; align-items: center;">
          <a-button
            type="text"
            @click="collapsed = !collapsed"
            style="margin-left: 16px;"
          >
            <MenuUnfoldOutlined v-if="collapsed" />
            <MenuFoldOutlined v-else />
          </a-button>
          <h1 style="margin: 0 0 0 16px; font-size: 18px;">Crontab Go 管理系统</h1>
        </div>
        
        <div style="margin-right: 16px;">
          <a-button type="text" @click="toggleTheme">
            <BulbOutlined v-if="isDark" />
            <BulbFilled v-else />
          </a-button>
        </div>
      </a-layout-header>

      <a-layout-content style="margin: 24px 16px; padding: 24px; background: #fff; min-height: 280px;">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from './stores/user'
import {
  DashboardOutlined,
  ClockCircleOutlined,
  FileTextOutlined,
  MonitorOutlined,
  LogoutOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  BulbOutlined,
  BulbFilled
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const collapsed = ref(false)
const isDark = ref(false)
const selectedKeys = ref(['dashboard'])

const menuItems = [
  { title: '仪表板', icon: DashboardOutlined, to: '/dashboard', value: 'dashboard' },
  { title: '任务管理', icon: ClockCircleOutlined, to: '/tasks', value: 'tasks' },
  { title: '执行日志', icon: FileTextOutlined, to: '/logs', value: 'logs' },
  { title: '系统监控', icon: MonitorOutlined, to: '/system', value: 'system' }
]

const onBreakpoint = (broken) => {
  collapsed.value = broken
}

const handleMenuClick = ({ key }) => {
  const item = menuItems.find(item => item.value === key)
  if (item) {
    router.push(item.to)
  }
}

const toggleTheme = () => {
  isDark.value = !isDark.value
  // 这里可以添加主题切换逻辑
}

const logout = () => {
  userStore.logout()
  router.push('/login')
}

// 监听路由变化，更新选中的菜单项
watch(() => route.path, (newPath) => {
  const currentItem = menuItems.find(item => item.to === newPath)
  if (currentItem) {
    selectedKeys.value = [currentItem.value]
  }
}, { immediate: true })
</script>
<style scoped>
.user-info {
  padding: 16px;
  text-align: center;
  color: rgba(255, 255, 255, 0.85);
}

.user-details {
  margin-top: 8px;
}

.username {
  font-weight: 500;
  font-size: 14px;
}

.email {
  font-size: 12px;
  opacity: 0.7;
}

.logout-btn {
  position: absolute;
  bottom: 16px;
  left: 0;
  right: 0;
}

:deep(.ant-layout-sider) {
  position: relative;
}
</style>