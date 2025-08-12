<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="elevation-12">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title>注册</v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <v-form ref="form" v-model="valid" lazy-validation>
              <v-text-field
                v-model="userData.username"
                :rules="usernameRules"
                label="用户名"
                prepend-icon="mdi-account"
                required
              ></v-text-field>

              <v-text-field
                v-model="userData.email"
                :rules="emailRules"
                label="邮箱"
                prepend-icon="mdi-email"
                required
              ></v-text-field>

              <v-text-field
                v-model="userData.password"
                :rules="passwordRules"
                label="密码"
                prepend-icon="mdi-lock"
                type="password"
                required
              ></v-text-field>

              <v-text-field
                v-model="confirmPassword"
                :rules="confirmPasswordRules"
                label="确认密码"
                prepend-icon="mdi-lock-check"
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
              @click="handleRegister"
            >
              注册
            </v-btn>
          </v-card-actions>
          <v-divider></v-divider>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              text
              color="primary"
              @click="$router.push('/login')"
            >
              已有账号？登录
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

const userData = ref({
  username: '',
  email: '',
  password: ''
})

const confirmPassword = ref('')

const usernameRules = [
  v => !!v || '用户名不能为空',
  v => v.length >= 3 || '用户名至少3个字符'
]

const emailRules = [
  v => !!v || '邮箱不能为空',
  v => /.+@.+\..+/.test(v) || '邮箱格式不正确'
]

const passwordRules = [
  v => !!v || '密码不能为空',
  v => v.length >= 6 || '密码至少6个字符'
]

const confirmPasswordRules = [
  v => !!v || '请确认密码',
  v => v === userData.value.password || '两次密码不一致'
]

const handleRegister = async () => {
  if (!form.value.validate()) return

  loading.value = true
  try {
    await userStore.register(userData.value)
    snackbarText.value = '注册成功，请登录'
    snackbarColor.value = 'success'
    snackbar.value = true
    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } catch (error) {
    snackbarText.value = error.response?.data?.message || '注册失败'
    snackbarColor.value = 'error'
    snackbar.value = true
  } finally {
    loading.value = false
  }
}
</script>