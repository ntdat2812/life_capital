import { defineStore } from 'pinia'
import api from '../lib/api'

export const useWealthStore = defineStore('wealth', {
  state: () => ({
    netWorthSummary: {
      total_assets: 0,
      total_liabilities: 0,
      net_worth: 0,
      base_currency: 'VND'
    },
    assets: [],
    assetPage: 1,
    assetTotalPages: 1,
    liabilities: [],
    liabilityPage: 1,
    liabilityTotalPages: 1,
    loading: false,
    error: null,
    assetCategoryFilter: '',
    liabilityCategoryFilter: '',
    assetSort: 'value_desc',
    liabilitySort: 'value_desc'
  }),

  actions: {
    async fetchNetWorthSummary() {
      try {
        const response = await api.get('/wealth/net-worth')
        this.netWorthSummary = response.data
      } catch (err) {
        console.error('Failed to fetch net worth summary', err)
      }
    },

    async fetchAssets(append = false, limit = 20) {
      if (!append) {
        this.loading = true
        this.assetPage = 1
      } else {
        this.assetPage++
      }
      this.error = null
      
      try {
        const offset = (this.assetPage - 1) * limit
        let url = `/wealth/assets?limit=${limit}&offset=${offset}&sort=${this.assetSort}`
        if (this.assetCategoryFilter) {
          url += `&category=${this.assetCategoryFilter}`
        }
        
        const response = await api.get(url)
        
        if (append) {
          this.assets = [...this.assets, ...(response.data.data || [])]
        } else {
          this.assets = response.data.data || []
        }
        this.assetTotalPages = response.data.total_pages || 1
      } catch (err) {
        this.error = err.response?.data?.message || 'Failed to load assets'
      } finally {
        if (!append) this.loading = false
      }
    },

    async createAsset(assetData) {
      try {
        await api.post('/wealth/assets', assetData)
        await this.fetchAssets()
        await this.fetchNetWorthSummary()
      } catch (err) {
        throw err.response?.data?.message || 'Failed to create asset'
      }
    },

    async updateAsset(id, assetData) {
      try {
        await api.put(`/wealth/assets/${id}`, assetData)
        // Optionally update the item in the list directly or refetch. We'll refetch.
        await this.fetchAssets()
        await this.fetchNetWorthSummary()
      } catch (err) {
        throw err.response?.data?.message || 'Failed to update asset'
      }
    },

    async deleteAsset(id) {
      try {
        await api.delete(`/wealth/assets/${id}`)
        await this.fetchAssets()
        await this.fetchNetWorthSummary()
      } catch (err) {
        throw err.response?.data?.message || 'Failed to delete asset'
      }
    },

    async fetchLiabilities(append = false, limit = 20) {
      if (!append) {
        this.loading = true
        this.liabilityPage = 1
      } else {
        this.liabilityPage++
      }
      this.error = null
      
      try {
        const offset = (this.liabilityPage - 1) * limit
        let url = `/wealth/liabilities?limit=${limit}&offset=${offset}&sort=${this.liabilitySort}`
        if (this.liabilityCategoryFilter) {
          url += `&category=${this.liabilityCategoryFilter}`
        }
        
        const response = await api.get(url)
        
        if (append) {
          this.liabilities = [...this.liabilities, ...(response.data.data || [])]
        } else {
          this.liabilities = response.data.data || []
        }
        this.liabilityTotalPages = response.data.total_pages || 1
      } catch (err) {
        this.error = err.response?.data?.message || 'Failed to load liabilities'
      } finally {
        if (!append) this.loading = false
      }
    },

    async createLiability(liabilityData) {
      try {
        await api.post('/wealth/liabilities', liabilityData)
        await this.fetchLiabilities()
        await this.fetchNetWorthSummary()
      } catch (err) {
        throw err.response?.data?.message || 'Failed to create liability'
      }
    },

    async updateLiability(id, liabilityData) {
      try {
        await api.put(`/wealth/liabilities/${id}`, liabilityData)
        await this.fetchLiabilities()
        await this.fetchNetWorthSummary()
      } catch (err) {
        throw err.response?.data?.message || 'Failed to update liability'
      }
    },

    async deleteLiability(id) {
      try {
        await api.delete(`/wealth/liabilities/${id}`)
        await this.fetchLiabilities()
        await this.fetchNetWorthSummary()
      } catch (err) {
        throw err.response?.data?.message || 'Failed to delete liability'
      }
    },

    async fetchAll() {
      this.loading = true
      this.assetPage = 1
      this.liabilityPage = 1
      this.assetCategoryFilter = ''
      this.liabilityCategoryFilter = ''
      this.assetSort = 'value_desc'
      this.liabilitySort = 'value_desc'
      try {
        await Promise.all([
          this.fetchNetWorthSummary(),
          this.fetchAssets(),
          this.fetchLiabilities()
        ])
      } finally {
        this.loading = false
      }
    }
  }
})
