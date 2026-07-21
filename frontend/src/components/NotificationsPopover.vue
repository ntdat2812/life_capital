<template>
  <div class="relative" ref="popoverRef">
    <button 
      @click="toggle"
      class="relative p-2 text-slate-400 hover:text-white transition rounded-full hover:bg-slate-800/50"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"></path>
      </svg>
      <span v-if="notificationStore.unreadCount > 0" class="absolute top-1 right-1 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-[9px] font-bold text-white shadow-lg border border-slate-900">
        {{ notificationStore.unreadCount > 9 ? '9+' : notificationStore.unreadCount }}
      </span>
    </button>

    <transition name="popover-fade">
      <div v-if="isOpen" class="absolute right-0 mt-2 w-80 sm:w-96 bg-slate-900 border border-slate-700/60 rounded-2xl shadow-2xl z-50 overflow-hidden backdrop-blur-xl bg-slate-900/95">
        <div class="px-4 py-3 border-b border-slate-800 flex justify-between items-center bg-slate-800/20">
          <h3 class="text-sm font-semibold text-white font-outfit">Thông báo</h3>
          <span v-if="notificationStore.unreadCount > 0" class="text-xs text-indigo-400 font-medium bg-indigo-500/10 px-2 py-0.5 rounded-md">
            {{ notificationStore.unreadCount }} chưa đọc
          </span>
        </div>

        <div class="max-h-96 overflow-y-auto custom-scrollbar p-2 space-y-1">
          <div v-if="notificationStore.loading" class="flex justify-center py-8">
            <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-500"></div>
          </div>
          
          <div v-else-if="notificationStore.error" class="text-red-400 text-xs p-4 text-center">
            {{ notificationStore.error }}
          </div>
          
          <div v-else-if="notificationStore.notifications.length === 0" class="text-slate-400 text-sm p-8 text-center flex flex-col items-center">
            <svg class="w-8 h-8 text-slate-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path></svg>
            Không có thông báo mới
          </div>

          <div 
            v-else
            v-for="notif in notificationStore.notifications" 
            :key="notif.id"
            @click="handleClick(notif)"
            class="p-3 rounded-xl transition cursor-pointer flex gap-3 relative overflow-hidden group"
            :class="notif.is_read ? 'hover:bg-slate-800/40 opacity-70' : 'bg-slate-800/60 hover:bg-slate-800 border border-slate-700/50'"
          >
            <div v-if="!notif.is_read" class="absolute left-0 top-0 bottom-0 w-0.5 bg-indigo-500"></div>
            
            <div class="text-xl flex-shrink-0 mt-0.5">{{ getIconForType(notif.type) }}</div>
            <div class="flex-1 min-w-0">
              <div class="flex justify-between items-start mb-0.5">
                <h4 class="text-sm truncate pr-2 font-medium" :class="notif.is_read ? 'text-slate-300' : 'text-white'">{{ notif.title }}</h4>
                <span class="text-[10px] text-slate-500 whitespace-nowrap mt-0.5">{{ formatDate(notif.created_at) }}</span>
              </div>
              <p class="text-xs text-slate-400 line-clamp-2 leading-relaxed">{{ notif.message }}</p>
            </div>
          </div>
        </div>
        
        <div class="p-2 border-t border-slate-800 bg-slate-900/50 text-center">
          <button @click="markAllAsRead" v-if="notificationStore.unreadCount > 0" class="text-xs text-indigo-400 hover:text-indigo-300 font-medium transition px-3 py-1.5 hover:bg-indigo-500/10 rounded-lg w-full">
            Đánh dấu tất cả đã đọc
          </button>
          <span v-else class="text-[10px] text-slate-600 uppercase font-bold tracking-wider">Hệ thống luôn đồng hành cùng bạn</span>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationStore } from '../stores/notificationStore'

const notificationStore = useNotificationStore()
const router = useRouter()
const isOpen = ref(false)
const popoverRef = ref(null)

const toggle = async () => {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    await notificationStore.fetchNotifications()
  }
}

const getIconForType = (type) => {
  switch (type) {
    case 'IPS_RECOMMENDATION': return '🤖'
    case 'PRICE_ALERT': return '📈'
    case 'EARNINGS_ALERT': return '⚠️'
    default: return '🔔'
  }
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('vi-VN', { month: 'short', day: 'numeric' })
}

const handleClick = async (notif) => {
  if (!notif.is_read) {
    await notificationStore.markAsRead(notif.id)
  }
  if (notif.action_link) {
    router.push(notif.action_link)
    isOpen.value = false
  }
}

const markAllAsRead = async () => {
  for (const notif of notificationStore.notifications) {
    if (!notif.is_read) {
      await notificationStore.markAsRead(notif.id)
    }
  }
}

const handleClickOutside = (event) => {
  if (popoverRef.value && !popoverRef.value.contains(event.target)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.popover-fade-enter-active,
.popover-fade-leave-active {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: top right;
}

.popover-fade-enter-from,
.popover-fade-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(-5px);
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #475569;
}
</style>
