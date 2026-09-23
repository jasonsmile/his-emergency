<template>
  <div class="login-container">
    <el-card class="login-card">
      <div class="login-title">
        <h2>应急HIS</h2>
        <p>请登录您的账号</p>
      </div>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top"
        @submit.prevent="handleLogin">
        <el-form-item label="登录名" prop="loginName">
          <el-input v-model="form.loginName" :prefix-icon="User" placeholder="请输入登录名"
            autocomplete="username" :disabled="loading" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" :prefix-icon="Lock" type="password" show-password
            placeholder="请输入密码" autocomplete="current-password" :disabled="loading" />
        </el-form-item>
        <el-button class="login-button" type="primary" native-type="submit" :loading="loading">
          登录
        </el-button>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const formRef = ref()
const loading = ref(false)
const form = reactive({ loginName: '', password: '' })
const rules = {
  loginName: [{ required: true, whitespace: true, message: '请输入登录名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function handleLogin() {
  if (loading.value || !formRef.value) return
  loading.value = true
  try {
    const valid = await formRef.value.validate().catch(() => false)
    if (!valid) return
    await userStore.login(form.loginName.trim(), form.password)
    ElMessage.success('登录成功')
    const redirect = route.query.redirect
    // 只返回已注册的受保护业务路由。
    const target = typeof redirect === 'string' && redirect.startsWith('/') &&
      !redirect.startsWith('//') && !redirect.includes('\\') &&
      router.resolve(redirect).meta.requiresAuth ? redirect : '/registration'
    await router.replace(target)
  } catch (error) {
    // 接口错误由 Axios 统一提示；本地响应校验错误在此提示。
    if (error.message === '登录响应缺少有效 Token') ElMessage.error(error.message)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container { min-height: 100vh; padding: 24px; display: flex;
  justify-content: center; align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.login-card { width: 400px; max-width: 100%; }
.login-title { text-align: center; margin-bottom: 30px; }
.login-title h2 { margin: 0; color: #303133; }
.login-title p { margin: 10px 0 0; color: #909399; font-size: 14px; }
.login-button { width: 100%; }
</style>
