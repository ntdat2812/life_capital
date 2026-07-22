import { defineStore } from 'pinia'
import api from '../lib/api'

export const usePortfolioStore = defineStore('portfolio', {
  state: () => ({
    holdings: [],
    watchlist: [],
    theses: [],
    loading: false,
    error: null,
  }),

  actions: {
    async fetchHoldings() {
      this.loading = true
      this.error = null
      try {
        const response = await api.get('/portfolio')
        this.holdings = response.data || []
      } catch (err) {
        this.error = err.response?.data?.message || err.message
        console.error('Failed to fetch holdings:', err)
      } finally {
        this.loading = false
      }
    },

    async fetchWatchlist() {
      this.loading = true
      this.error = null
      try {
        const response = await api.get('/watchlist')
        this.watchlist = response.data || []
      } catch (err) {
        this.error = err.response?.data?.message || err.message
        console.error('Failed to fetch watchlist:', err)
      } finally {
        this.loading = false
      }
    },

    async addWatchlist(item) {
      this.loading = true
      this.error = null
      try {
        const response = await api.post('/watchlist', item)
        this.watchlist.push(response.data)
        return response.data
      } catch (err) {
        this.error = err.response?.data?.message || err.message
        throw err
      } finally {
        this.loading = false
      }
    },

    async deleteWatchlist(id) {
      this.loading = true
      this.error = null
      try {
        await api.delete(`/watchlist/${id}`)
        this.watchlist = this.watchlist.filter(item => item.id !== id)
      } catch (err) {
        this.error = err.response?.data?.message || err.message
        throw err
      } finally {
        this.loading = false
      }
    },

    async fetchTheses() {
      this.loading = true
      this.error = null
      try {
        const response = await api.get('/theses')
        this.theses = response.data || []
      } catch (err) {
        this.error = err.response?.data?.message || err.message
        console.error('Failed to fetch theses:', err)
      } finally {
        this.loading = false
      }
    },

    async getThesisByTicker(ticker) {
      this.loading = true
      this.error = null
      try {
        const response = await api.get(`/theses/${ticker}`)
        return response.data
      } catch (err) {
        if (err.response?.status === 404) return null
        this.error = err.response?.data?.message || err.message
        throw err
      } finally {
        this.loading = false
      }
    },

    async createThesis(thesis) {
      this.loading = true
      this.error = null
      try {
        const response = await api.post('/theses', thesis)
        this.theses.push(response.data)
        return response.data
      } catch (err) {
        this.error = err.response?.data?.message || err.message
        throw err
      } finally {
        this.loading = false
      }
    },

    async updateThesis(id, thesis) {
      this.loading = true
      this.error = null
      try {
        const response = await api.put(`/theses/${id}`, thesis)
        const index = this.theses.findIndex(t => t.id === id)
        if (index !== -1) {
          this.theses[index] = response.data
        }
        return response.data
      } catch (err) {
        this.error = err.response?.data?.message || err.message
        throw err
      } finally {
        this.loading = false
      }
    },

    async deleteThesis(id) {
      this.loading = true
      this.error = null
      try {
        await api.delete(`/theses/${id}`)
        this.theses = this.theses.filter(t => t.id !== id)
      } catch (err) {
        this.error = err.response?.data?.message || err.message
        throw err
      } finally {
        this.loading = false
      }
    },

    async generateThesisAI(ticker, assetType, companyName) {
      this.loading = true
      this.error = null
      try {
        const response = await api.post('/theses/generate', {
          ticker: ticker,
          asset_type: assetType,
          company_name: companyName
        })
        return response.data
      } catch (err) {
        this.error = err.response?.data?.message || err.message
        throw err
      } finally {
        this.loading = false
      }
    }
  }
})
