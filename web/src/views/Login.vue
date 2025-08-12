<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="elevation-12">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title>登录</v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <v-form ref="form" v-model="valid" lazy-validation>
              <v-text-field
                v-model="credentials.username"
                :rules="usernameRules"
                label="用户名"
                prepend-icon="mdi-account"
                required
              ></v-text-field>

              <v-text-field
                v-model="credentials.password"
                :rules="passwordRules"
                label="密码"
                prepend-icon="mdi-lock"
                type="password"
                required
              ></v-text-field>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              color="primary"
              :disabled="!valid"
              :loading="loading"
              @click="handleLogin"
            >
              登录
            </v-btn>
          </v-card-actions>
          <v-divider></v-divider>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              text
              color="primary"
              @click="$router.push('/register')"
            >
              没有账号？注册
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

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
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const form = ref(null)
const valid = ref(false)
const loading = ref(false)
const snackbar = ref(false)
const snackbarText = ref('')
const snackbarColor = ref('success')

const credentials = ref({
  username: '',
  password: ''
})

const usernameRules = [
  v => !!v || '用户名不能为空'
]

const passwordRules = [
  v => !!v || '密码不能为空'
]

const handleLogin = async () => {
  if (!form.value.validate()) return

  loading.value = true
  try {
    await userStore.login(credentials.value)
    snackbarText.value = '登录成功'
    snackbarColor.value = 'success'
    snackbar.value = true
    router.push('/dashboard')
  } catch (error) {
    snackbarText.value = error.response?.data?.message || '登录失败'
    snackbarColor.value = 'error'
    snackbar.value = true
  } finally {
    loading.value = false
  }
}
</script>