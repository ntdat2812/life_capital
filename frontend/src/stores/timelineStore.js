import { defineStore } from 'pinia'
import api from '../lib/api'

export const useTimelineStore = defineStore('timeline', {
  state: () => ({
    events: [],
    loading: false,
    error: null,
    draftEvent: null,
  }),
  
  actions: {
    async fetchTimeline() {
      this.loading = true
      this.error = null
      try {
        const response = await api.get('/profile/timeline')
        this.events = response.data || []
      } catch (err) {
        this.error = err.response?.data?.error || 'Lỗi khi tải lịch sử sự kiện'
        throw err
      } finally {
        this.loading = false
      }
    },
    
    async analyzeEvent(eventText) {
      this.loading = true
      this.error = null
      this.draftEvent = null
      try {
        const response = await api.post('/profile/timeline/analyze', {
          event_text: eventText
        })
        this.draftEvent = {
          ...response.data,
          event_text: eventText,
          title: eventText,
        }
        return this.draftEvent
      } catch (err) {
        this.error = err.response?.data?.error || 'Lỗi khi phân tích sự kiện'
        throw err
      } finally {
        this.loading = false
      }
    },
    
    async confirmEvent(payload) {
      this.loading = true
      this.error = null
      try {
        await api.post('/profile/timeline/confirm', payload)
        this.draftEvent = null
        await this.fetchTimeline()
      } catch (err) {
        this.error = err.response?.data?.error || 'Lỗi khi lưu sự kiện'
        throw err
      } finally {
        this.loading = false
      }
    }
  }
})
