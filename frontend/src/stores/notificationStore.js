import { defineStore } from 'pinia';
import api from '../lib/api';

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    notifications: [],
    unreadCount: 0,
    loading: false,
    error: null,
  }),
  actions: {
    async fetchNotifications() {
      this.loading = true;
      this.error = null;
      try {
        const response = await api.get('/notifications'); // Assumes baseURL is /api/v1
        this.notifications = response.data || [];
      } catch (err) {
        this.error = err.response?.data?.error || 'Failed to fetch notifications';
        console.error('Error fetching notifications:', err);
      } finally {
        this.loading = false;
      }
    },
    async fetchUnreadCount() {
      try {
        const response = await api.get('/notifications/unread-count');
        this.unreadCount = response.data.count;
      } catch (err) {
        console.error('Error fetching unread count:', err);
      }
    },
    async markAsRead(id) {
      try {
        await api.put(`/notifications/${id}/read`);
        const notif = this.notifications.find((n) => n.id === id);
        if (notif && !notif.is_read) {
          notif.is_read = true;
          this.unreadCount = Math.max(0, this.unreadCount - 1);
        }
      } catch (err) {
        console.error('Error marking notification as read:', err);
      }
    }
  }
});
