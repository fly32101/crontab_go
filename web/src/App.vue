<template>
  <v-app>
    <v-navigation-drawer
      v-model="drawer"
      app
      :permanent="$vuetify.display.lgAndUp"
    >
      <v-list>
        <v-list-item
          prepend-avatar="https://randomuser.me/api/portraits/men/85.jpg"
          :title="userStore.user?.username || '未登录'"
          :subtitle="userStore.user?.email || ''"
        ></v-list-item>
      </v-list>

      <v-divider></v-divider>

      <v-list density="compact" nav>
        <v-list-item
          v-for="item in menuItems"
          :key="item.title"
          :prepend-icon="item.icon"
          :title="item.title"
          :to="item.to"
          :value="item.value"
        ></v-list-item>
      </v-list>

      <template v-slot:append>
        <div class="pa-2">
          <v-btn
            block
            color="error"
            @click="logout"
            prepend-icon="mdi-logout"
          >
            退出登录
          </v-btn>
        </div>
      </template>
    </v-navigation-drawer>

    <v-app-bar app>
      <v-app-bar-nav-icon
        v-if="!$vuetify.display.lgAndUp"
        @click="drawer = !drawer"
      ></v-app-bar-nav-icon>

      <v-app-bar-title>Crontab Go 管理系统</v-app-bar-title>

      <v-spacer></v-spacer>

      <v-btn
        icon="mdi-theme-light-dark"
        @click="toggleTheme"
      ></v-btn>
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useTheme } from 'vuetify'
import { useRouter } from 'vue-router'
import { useUserStore } from './stores/user'

const theme = useTheme()
const router = useRouter()
const userStore = useUserStore()

const drawer = ref(false)

const menuItems = [
  { title: '仪表板', icon: 'mdi-view-dashboard', to: '/dashboard', value: 'dashboard' },
  { title: '任务管理', icon: 'mdi-clock-outline', to: '/tasks', value: 'tasks' },
  { title: '执行日志', icon: 'mdi-file-document-outline', to: '/logs', value: 'logs' },
  { title: '系统监控', icon: 'mdi-monitor', to: '/system', value: 'system' }
]

const toggleTheme = () => {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}

const logout = () => {
  userStore.logout()
  router.push('/login')
}
</script>