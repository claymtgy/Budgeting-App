import { defineStore } from 'pinia'
import api from '@/api/client'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user') || 'null')
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.token),
    joinCode: (state) => state.user?.join_code || ''
  },
  actions: {
    setSession(token, user) {
      this.token = token
      this.user = user
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify(user))
    },
    clearSession() {
      this.token = ''
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    },
    async register(email, password, joinCode = '') {
      const { data } = await api.post('/api/auth/register', {
        email,
        password,
        join_code: joinCode.trim().toUpperCase()
      })
      this.setSession(data.token, data.user)
    },
    async login(email, password) {
      const { data } = await api.post('/api/auth/login', { email, password })
      this.setSession(data.token, data.user)
    },
    async fetchMe() {
      const { data } = await api.get('/api/auth/me')
      this.user = data
      localStorage.setItem('user', JSON.stringify(data))
    },
    logout() {
      this.clearSession()
    }
  }
})
