import { defineStore } from 'pinia'
import api from '../lib/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user')) || null,
    token: localStorage.getItem('token') || null,
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
  },
  actions: {
    setAuth(user, token) {
      this.user = user
      this.token = token
      localStorage.setItem('user', JSON.stringify(user))
      localStorage.setItem('token', token)
    },
    logout() {
      this.user = null
      this.token = null
      localStorage.removeItem('user')
      localStorage.removeItem('token')
    },
    async login(email, password) {
      try {
        const response = await api.post('/auth/login', { email, password })
        this.setAuth(response.data.user, response.data.token)
        return true
      } catch (error) {
        console.error('Login failed:', error)
        throw error
      }
    },
    async signup(name, email, password) {
      try {
        const response = await api.post('/auth/signup', { name, email, password })
        return response.data
      } catch (error) {
        console.error('Signup failed:', error)
        throw error
      }
    },
    async loginWithGoogle(credential) {
      try {
        const response = await api.post('/auth/google', { id_token: credential })
        this.setAuth(response.data.user, response.data.token)
        return true
      } catch (error) {
        console.error('Google login failed:', error)
        throw error
      }
    }
  }
})
