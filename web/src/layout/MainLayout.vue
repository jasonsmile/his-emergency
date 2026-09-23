<template>
  <el-container class="layout-container">
    <el-aside width="200px" class="sidebar">
      <div class="logo">应急HIS</div>
      <el-menu :default-active="route.path" router background-color="#304156"
        text-color="#bfcbd9" active-text-color="#409eff">
        <el-menu-item index="/registration">
          <el-icon><DocumentAdd /></el-icon><span>挂号工作台</span>
        </el-menu-item>
        <el-menu-item index="/registration/list">
          <el-icon><Document /></el-icon><span>就诊列表</span>
        </el-menu-item>
        <el-menu-item index="/charge">
          <el-icon><Money /></el-icon><span>收费工作台</span>
        </el-menu-item>
        <el-menu-item index="/charge/report">
          <el-icon><DataAnalysis /></el-icon><span>日结报表</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container class="content-container">
      <el-header class="header">
        <h1>{{ route.meta.title || '应急HIS' }}</h1>
        <div class="user-actions">
          <span>{{ userStore.userName || userStore.userInfo.user_id }}</span>
          <el-button @click="handleLogout">退出登录</el-button>
        </div>
      </el-header>
      <el-main><router-view /></el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { DocumentAdd, Document, Money, DataAnalysis } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

async function handleLogout() {
  userStore.logout()
  await router.replace({ name: 'Login' })
}
</script>

<style scoped>
.layout-container { min-height: 100vh; background: #f5f7fa; }
.sidebar { background: #304156; }
.logo { height: 60px; display: flex; align-items: center; justify-content: center;
  color: white; font-size: 22px; font-weight: 600; }
.el-menu { border-right: 0; }
.content-container { min-width: 0; }
.header { display: flex; align-items: center; justify-content: space-between;
  gap: 16px; background: white; border-bottom: 1px solid #e4e7ed; }
.header h1 { margin: 0; font-size: 18px; color: #303133; }
.user-actions { display: flex; align-items: center; gap: 16px; }
</style>
