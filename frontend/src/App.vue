<template>
  <div class="min-h-screen flex flex-col md:flex-row relative">
    <!-- Decorative background glows -->
    <div class="absolute top-[-10%] left-[-10%] w-[500px] h-[500px] rounded-full bg-indigo-900/10 blur-[120px] pointer-events-none"></div>
    <div class="absolute bottom-[-10%] right-[-10%] w-[500px] h-[500px] rounded-full bg-purple-900/10 blur-[120px] pointer-events-none"></div>

    <!-- Navigation Sidebar -->
    <aside class="w-full md:w-64 shrink-0 glass-card md:min-h-screen p-6 flex flex-col justify-between border-b md:border-b-0 md:border-r border-slate-800 relative z-20">
      <div class="space-y-8">
        <!-- Logo -->
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-indigo-500 to-purple-600 flex items-center justify-center font-outfit font-bold text-xl shadow-lg shadow-indigo-500/20 text-white">
            LC
          </div>
          <div>
            <h1 class="font-outfit text-base font-bold tracking-tight text-slate-100">Life Capital</h1>
            <span class="text-[10px] text-slate-400 font-mono">WealthOS v0.1</span>
          </div>
        </div>

        <!-- Sidebar Navigation -->
        <nav v-if="authStore.isAuthenticated" class="space-y-1">
          <router-link 
            to="/" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'dashboard' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            📊 Tổng quan
          </router-link>
          
          <router-link 
            to="/onboarding/interview" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'onboarding-interview' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            🤖 Lập hồ sơ (AI)
          </router-link>

          <router-link 
            to="/profile" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'profile' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            👤 Hồ sơ cá nhân
          </router-link>

          <router-link 
            to="/profile/timeline" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'timeline' || $route.name === 'log-event' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            ⏳ Dòng thời gian
          </router-link>

          <router-link 
            to="/assets" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'assets' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            💰 Tài sản & Nợ
          </router-link>

          <router-link 
            to="/goals" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'goals' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            🎯 Mục tiêu tài chính
          </router-link>

          <router-link 
            to="/portfolio" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'portfolio' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            📈 Danh mục đầu tư
          </router-link>

          <router-link 
            to="/ips" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'ips' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            🧭 Chiến lược đầu tư
          </router-link>

          <router-link 
            to="/thesis" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'thesis' || $route.name === 'thesis-detail' || $route.name === 'thesis-new' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            📝 Luận điểm đầu tư
          </router-link>

          <router-link 
            to="/review" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'review' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            🗓️ Đánh giá hàng tháng
          </router-link>
        </nav>
      </div>

      <!-- Footer / Auth trigger mock -->
      <div class="mt-8 border-t border-slate-800/60 pt-4">
        <template v-if="authStore.isAuthenticated">
          <div class="flex flex-col gap-2">
            <div class="px-4 py-2 text-xs text-slate-400 font-semibold truncate">
              👤 {{ authStore.user?.name }}
            </div>
            <button 
              @click="handleLogout"
              class="flex w-full items-center gap-3 px-4 py-2 text-xs font-medium text-red-400 hover:text-red-300 hover:bg-red-900/20 rounded-xl transition text-left"
            >
              🚪 Đăng xuất
            </button>
          </div>
        </template>
        <template v-else>
          <router-link 
            to="/login"
            class="flex items-center gap-3 px-4 py-2 text-xs font-medium text-slate-400 hover:text-slate-200 transition"
          >
            🔐 Đăng nhập hệ thống
          </router-link>
        </template>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-grow relative z-10 flex flex-col h-screen overflow-hidden">
      
      <!-- Top Header -->
      <header v-if="authStore.isAuthenticated" class="h-16 border-b border-slate-800 bg-slate-900/50 backdrop-blur-md flex items-center justify-end px-6 sticky top-0 z-40">
        <div class="flex items-center gap-4">
          <NotificationsPopover />
          
          <div class="w-8 h-8 rounded-full bg-indigo-600 flex items-center justify-center text-sm font-bold text-white shadow-lg border border-indigo-400">
            {{ authStore.user?.name ? authStore.user.name.charAt(0).toUpperCase() : 'U' }}
          </div>
        </div>
      </header>

      <div class="flex-grow p-4 md:p-8 overflow-y-auto custom-scrollbar">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </main>
  </div>
</template>

<script setup>
import { onMounted, watch } from 'vue'
import { useAuthStore } from './stores/authStore'
import { useNotificationStore } from './stores/notificationStore'
import { useRouter } from 'vue-router'
import NotificationsPopover from './components/NotificationsPopover.vue'

const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const router = useRouter()

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

watch(() => authStore.isAuthenticated, (isAuth) => {
  if (isAuth) {
    notificationStore.fetchUnreadCount()
  }
})

onMounted(() => {
  if (authStore.isAuthenticated) {
    notificationStore.fetchUnreadCount()
  }
})
</script>

<style>
/* Vue route transition animations */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
