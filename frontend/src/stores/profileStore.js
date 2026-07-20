import { defineStore } from 'pinia';
import api from '../lib/api';

export const useProfileStore = defineStore('profile', {
  state: () => ({
    profile: null,
    incomes: [],
    dependents: [],
    loading: false,
    error: null,
  }),
  getters: {
    totalIncome: (state) => {
      return state.incomes.filter(i => i.is_active).reduce((sum, i) => sum + (i.frequency === 'yearly' ? i.amount / 12 : i.amount), 0)
    },
    essentialExpense: (state) => {
      return state.profile?.essential_monthly_expense || 0
    },
    discretionaryExpense: (state) => {
      return state.profile?.discretionary_monthly_expense || 0
    },
    dependentsExpense: (state) => {
      return state.dependents.filter(d => d.is_active).reduce((sum, d) => sum + d.monthly_cost, 0)
    }
  },
  actions: {
    async fetchProfile() {
      this.loading = true;
      this.error = null;
      try {
        const response = await api.get('/profile/me');
        this.profile = response.data;
        return this.profile;
      } catch (error) {
        if (error.response && error.response.status !== 404) {
          this.error = error.response.data.message || 'Lỗi khi tải hồ sơ';
        }
        this.profile = null;
      } finally {
        this.loading = false;
      }
    },
    async submitOnboarding(chatHistory) {
      this.loading = true;
      this.error = null;
      try {
        const response = await api.post('/profile/onboarding', { chat_history: chatHistory });
        this.profile = response.data;
        return this.profile;
      } catch (error) {
        this.error = error.response?.data?.message || 'Lỗi khi tạo hồ sơ';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    async updateProfile(data) {
      try {
        const response = await api.put('/profile/me', data);
        this.profile = response.data;
        return this.profile;
      } catch (error) {
        throw error.response?.data?.message || 'Lỗi khi cập nhật hồ sơ';
      }
    },
    async fetchIncomes() {
      try {
        const response = await api.get('/cashflow/income');
        this.incomes = response.data || [];
      } catch (error) {
        console.error('Failed to fetch incomes', error);
      }
    },
    async createIncome(data) {
      try {
        await api.post('/cashflow/income', data);
        await this.fetchIncomes();
      } catch (error) {
        throw error.response?.data?.message || 'Lỗi khi thêm thu nhập';
      }
    },
    async updateIncome(id, data) {
      try {
        await api.put(`/cashflow/income/${id}`, data);
        await this.fetchIncomes();
      } catch (error) {
        throw error.response?.data?.message || 'Lỗi khi cập nhật thu nhập';
      }
    },
    async deleteIncome(id) {
      try {
        await api.delete(`/cashflow/income/${id}`);
        await this.fetchIncomes();
      } catch (error) {
        throw error.response?.data?.message || 'Lỗi khi xoá thu nhập';
      }
    },
    async fetchDependents() {
      try {
        const response = await api.get('/cashflow/dependents');
        this.dependents = response.data || [];
      } catch (error) {
        console.error('Failed to fetch dependents', error);
      }
    },
    async createDependent(data) {
      try {
        await api.post('/cashflow/dependents', data);
        await this.fetchDependents();
      } catch (error) {
        throw error.response?.data?.message || 'Lỗi khi thêm người phụ thuộc';
      }
    },
    async updateDependent(id, data) {
      try {
        await api.put(`/cashflow/dependents/${id}`, data);
        await this.fetchDependents();
      } catch (error) {
        throw error.response?.data?.message || 'Lỗi khi cập nhật người phụ thuộc';
      }
    },
    async deleteDependent(id) {
      try {
        await api.delete(`/cashflow/dependents/${id}`);
        await this.fetchDependents();
      } catch (error) {
        throw error.response?.data?.message || 'Lỗi khi xoá người phụ thuộc';
      }
    }
  }
});
