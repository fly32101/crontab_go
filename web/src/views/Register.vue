<template>
  <div style="min-height: 100vh; display: flex; align-items: center; justify-content: center; background-color: #f0f2f5;">
    <a-card style="width: 400px;" title="注册">
      <a-form
        ref="formRef"
        :model="userData"
        :rules="formRules"
        layout="vertical"
        @finish="handleRegister"
      >
        <a-form-item label="用户名" name="username">
          <a-input
            v-model:value="userData.username"
            size="large"
            placeholder="请输入用户名"
          >
            <template #prefix>
              <UserOutlined />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item label="邮箱" name="email">
          <a-input
            v-model:value="userData.email"
            size="large"
            placeholder="请输入邮箱"
          >
            <template #prefix>
              <MailOutlined />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item label="密码" name="password">
          <a-input-password
            v-model:value="userData.password"
            size="large"
            placeholder="请输入密码"
          >
            <template #prefix>
              <LockOutlined />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item label="确认密码" name="confirmPassword">
          <a-input-password
            v-model:value="confirmPassword"
            size="large"
            placeholder="请再次输入密码"
          >
            <template #prefix>
              <LockOutlined />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            block
            :loading="loading"
          >
            注册
          </a-button>
        </a-form-item>
      </a-form>

      <a-divider />

      <div style="text-align: center;">
        <a-button type="link" @click="$router.push('/login')">
          已有账号？登录
        </a-button>
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { UserOutlined, MailOutlined, LockOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref(null)
const loading = ref(false)

const userData = ref({
  username: '',
  email: '',
  password: ''
})

const confirmPassword = ref('')

const formRules = {
  username: [
    { required: true, message: '请输入用户名' },
    { min: 3, message: '用户名至少3个字符' }
  ],
  email: [
    { required: true, message: '请输入邮箱' },
    { type: 'email', message: '邮箱格式不正确' }
  ],
  password: [
    { required: true, message: '请输入密码' },
    { min: 6, message: '密码至少6个字符' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码' },
    {
      validator: (rule, value) => {
        if (value !== userData.value.password) {
          return Promise.reject('两次密码不一致')
        }
        return Promise.resolve()
      }
    }
  ]
}

const handleRegister = async () => {
  loading.value = true
  try {
    await userStore.register(userData.value)
    message.success('注册成功，请登录')
    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } catch (error) {
    message.error(error.response?.data?.message || '注册失败')
  } finally {
    loading.value = false
  }
}
</script>