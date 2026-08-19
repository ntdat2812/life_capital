import { defineStore } from 'pinia'
import axios from 'axios'
import { useAuthStore } from './authStore'

export const useGoalStore = defineStore('goal', {
  state: () => ({
    goals: [],
    loading: false,
    error: null,
  }),

  actions: {
    async fetchGoals() {
      this.loading = true
      this.error = null
      try {
        const authStore = useAuthStore()
        const response = await axios.get(`${import.meta.env.VITE_API_URL}/goals`, {
          headers: {
            Authorization: `Bearer ${authStore.token}`
          }
        })
        
        this.goals = response.data || []
        
        // Sort by priority ASC
        this.goals.sort((a, b) => a.priority - b.priority)
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi khi tải mục tiêu'
        console.error('Error fetching goals:', error)
      } finally {
        this.loading = false
      }
    },

    async createGoal(goalData) {
      try {
        const authStore = useAuthStore()
        const response = await axios.post(`${import.meta.env.VITE_API_URL}/goals`, goalData, {
          headers: {
            Authorization: `Bearer ${authStore.token}`
          }
        })
        this.goals.push(response.data)
        this.goals.sort((a, b) => a.priority - b.priority)
        return response.data
      } catch (error) {
        console.error('Error creating goal:', error)
        throw error
      }
    },

    async updateGoal(id, goalData) {
      try {
        const authStore = useAuthStore()
        const response = await axios.put(`${import.meta.env.VITE_API_URL}/goals/${id}`, goalData, {
          headers: {
            Authorization: `Bearer ${authStore.token}`
          }
        })
        
        const index = this.goals.findIndex(g => g.id === id)
        if (index !== -1) {
          this.goals[index] = response.data
          this.goals.sort((a, b) => a.priority - b.priority)
        }
        return response.data
      } catch (error) {
        console.error('Error updating goal:', error)
        throw error
      }
    },

    async updatePriorities(orderedGoals) {
      // Optimistic update
      const originalGoals = [...this.goals]
      this.goals = orderedGoals

      try {
        const promises = orderedGoals.map((goal, index) => {
          const original = originalGoals.find(g => g.id === goal.id)
          if (original && original.priority !== index + 1) {
            return this.updateGoal(goal.id, {
              name: goal.name,
              target_amount: goal.target_amount,
              target_date: goal.target_date,
              priority: index + 1
            })
          }
          return Promise.resolve()
        })
        await Promise.all(promises)
      } catch (error) {
        console.error('Error updating priorities:', error)
        this.goals = originalGoals // Rollback
        throw error
      }
    },

    async deleteGoal(id) {
      try {
        const authStore = useAuthStore()
        await axios.delete(`${import.meta.env.VITE_API_URL}/goals/${id}`, {
          headers: {
            Authorization: `Bearer ${authStore.token}`
          }
        })
        
        this.goals = this.goals.filter(g => g.id !== id)
      } catch (error) {
        console.error('Error deleting goal:', error)
        throw error
      }
    },

    async linkAsset(goalId, assetId) {
      try {
        const authStore = useAuthStore()
        await axios.post(`${import.meta.env.VITE_API_URL}/goals/${goalId}/allocations`, { asset_id: assetId }, {
          headers: {
            Authorization: `Bearer ${authStore.token}`
          }
        })
        await this.fetchGoals() // Refresh goals to get updated allocations
      } catch (error) {
        console.error('Error linking asset:', error)
        throw error
      }
    },

    async unlinkAsset(goalId, assetId) {
      try {
        const authStore = useAuthStore()
        await axios.delete(`${import.meta.env.VITE_API_URL}/goals/${goalId}/allocations/${assetId}`, {
          headers: {
            Authorization: `Bearer ${authStore.token}`
          }
        })
        await this.fetchGoals() // Refresh goals to get updated allocations
      } catch (error) {
        console.error('Error unlinking asset:', error)
        throw error
      }
    }
  }
})
