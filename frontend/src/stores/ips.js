import { defineStore } from 'pinia'
import api from '../lib/api'

export const useIpsStore = defineStore('ips', {
  state: () => ({
    latestIps: null,
    isLoading: false,
    error: null,
    isGenerating: false
  }),

  actions: {
    async fetchLatestIps() {
      this.isLoading = true
      this.error = null
      try {
        const response = await api.get('/ips/latest')
        this.latestIps = response.data
      } catch (error) {
        if (error.response?.status !== 404) {
          this.error = error.response?.data?.error || 'Không thể tải IPS'
        } else {
          this.latestIps = null
        }
      } finally {
        this.isLoading = false
      }
    },

    async generateIps(preferredAssets = []) {
      this.isGenerating = true
      this.error = null
      try {
        const response = await api.post('/ips/generate', { preferred_asset_classes: preferredAssets })
        this.latestIps = response.data
        return this.latestIps
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi khi tạo IPS'
        throw this.error
      } finally {
        this.isGenerating = false
      }
    },

    async updateIps(updatedAllocation) {
      this.isLoading = true
      try {
        const payload = {
          ...this.latestIps,
          target_allocation: JSON.stringify(updatedAllocation)
        }
        await api.put('/ips/latest', payload)
        this.latestIps = payload
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi khi cập nhật IPS'
        throw this.error
      } finally {
        this.isLoading = false
      }
    }
  }
})
