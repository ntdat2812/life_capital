import { defineStore } from 'pinia'
import api from '../lib/api'

export const useReviewStore = defineStore('review', {
  state: () => ({
    history: [],
    currentReview: null,
    loading: false,
    generating: false,
    error: null
  }),

  actions: {
    async fetchHistory() {
      this.loading = true
      this.error = null
      try {
        const response = await api.get('/reviews')
        this.history = response.data || []
      } catch (err) {
        this.error = err.response?.data?.message || err.message
      } finally {
        this.loading = false
      }
    },

    async fetchReviewByMonth(month) {
      this.loading = true
      this.error = null
      try {
        const response = await api.get(`/reviews/${month}`)
        this.currentReview = response.data
      } catch (err) {
        this.error = err.response?.data?.message || err.message
      } finally {
        this.loading = false
      }
    },

    async generateReview(newInvestmentAmount) {
      this.generating = true
      this.error = null
      try {
        const response = await api.post('/reviews/generate', {
          new_investment_amount: newInvestmentAmount
        }, {
          timeout: 60000 // AI might take a while
        })
        this.currentReview = response.data
        return this.currentReview
      } catch (err) {
        this.error = err.response?.data?.message || err.message
        throw err
      } finally {
        this.generating = false
      }
    },
    
    async saveReview(reviewData) {
      this.loading = true
      this.error = null
      try {
        const response = await api.post('/reviews', reviewData)
        this.currentReview = response.data
        await this.fetchHistory()
        return this.currentReview
      } catch (err) {
        this.error = err.response?.data?.message || err.message
        throw err
      } finally {
        this.loading = false
      }
    },
    
    clearCurrentReview() {
      this.currentReview = null
    }
  }
})
