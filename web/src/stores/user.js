import { defineStore } from 'pinia'
import { login as loginApi } from '@/api/auth'

function readUserInfo() {
  try {
    const user = JSON.parse(localStorage.getItem('user') || '{}')
    if (user && typeof user === 'object' && !Array.isArray(user)) {
      const { user_id, user_name, dept_code } = user
      return { user_id, user_name, dept_code }
    }
  } catch {
    // 忽略损坏的缓存，避免阻止页面初始化。
  }
  return {}
}

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: readUserInfo(),
  }),

  getters: {
    isLoggedIn: state => !!state.token,
    userName: state => state.userInfo.user_name || '',
  },

  actions: {
    async login(loginName, password) {
      const res = await loginApi({ login_name: loginName, password })
      const { token, user_id, user_name, dept_code } = res.data || {}
      if (typeof token !== 'string' || !token.trim()) {
        throw new Error('登录响应缺少有效 Token')
      }

      this.token = token
      this.userInfo = { user_id, user_name, dept_code }
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify(this.userInfo))
      return res
    },

    logout() {
      this.token = ''
      this.userInfo = {}
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    },
  },
})
