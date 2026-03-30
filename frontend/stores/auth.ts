import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    isLoggedIn: false,
  }),

  actions: {
    setLoggedIn(value: boolean) {
      this.isLoggedIn = value
    },

    logout() {
      this.isLoggedIn = false
    },
  },
})
