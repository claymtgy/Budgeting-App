import axios from 'axios'
import router from '@/router'

function resolveApiBaseURL() {
  const configured = import.meta.env.VITE_API_URL
  if (configured) return configured
  // Vite dev server proxies /api to the backend (same origin, no CORS).
  if (import.meta.env.DEV) return ''
  return 'http://localhost:8080'
}

const api = axios.create({
  baseURL: resolveApiBaseURL()
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      if (router.currentRoute.value.name !== 'login') {
        router.push({ name: 'login' })
      }
    }
    return Promise.reject(error)
  }
)

export default api
