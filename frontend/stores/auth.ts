import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    isLoggedIn: false,
    username: '' as string,
  }),

  actions: {
    setLoggedIn(value: boolean, username = '') {
      this.isLoggedIn = value
      this.username = username
    },

    /** 调用服务端 check 接口验证 cookie token，返回是否登录 */
    async checkAuth(): Promise<boolean> {
      try {
        // SSR 阶段需要转发客户端 Cookie，否则服务端 $fetch 不会自动携带 auth_token
        const headers = import.meta.server ? useRequestHeaders(['cookie']) : {}
        const res = await $fetch<{ ok: boolean; username: string }>('/api/auth/check', {
          headers,
        })
        this.isLoggedIn = true
        this.username = res.username
        return true
      } catch {
        this.isLoggedIn = false
        this.username = ''
        return false
      }
    },

    logout() {
      this.isLoggedIn = false
      this.username = ''
    },
  },
})
