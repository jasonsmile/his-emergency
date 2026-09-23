import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { useUserStore } from '@/stores/user'

const request = axios.create({
  baseURL: '/api',
  timeout: 30000,
})

request.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token && !config.skipAuth) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  response => {
    const res = response.data
    if (res?.code !== 0) {
      const message = res?.message || '请求失败'
      ElMessage.error(message)
      return Promise.reject(new Error(message))
    }
    // 保留完整响应，兼容分页接口的 total、page、page_size。
    return res
  },
  error => {
    if (axios.isCancel(error)) return Promise.reject(error)

    if (error.response) {
      const { status, data } = error.response
      if (status === 401 && !error.config?.skipAuth) {
        ElMessage.error('登录已过期，请重新登录')
        useUserStore().logout()
        // 避免在登录页重复跳转。
        if (router.currentRoute.value.path !== '/login' &&
            router.getRoutes().some(route => route.path === '/login')) {
          router.replace('/login')
        }
      } else {
        ElMessage.error(data?.message || `请求失败(${status})`)
      }
    } else if (error.code === 'ECONNABORTED' || error.code === 'ETIMEDOUT') {
      ElMessage.error('请求超时，请稍后重试')
    } else {
      ElMessage.error('网络错误，请检查网络连接')
    }
    return Promise.reject(error)
  },
)

export default request
