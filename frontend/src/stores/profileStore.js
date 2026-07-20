import { defineStore } from 'pinia';
import api from '../lib/api';

export const useProfileStore = defineStore('profile', {
  state: () => ({
    profile: null,
    loading: false,
    error: null,
  }),
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
    }
  }
});
