<template>
  <div class="min-h-screen flex flex-col md:flex-row relative">
    <!-- Decorative background glows -->
    <div class="absolute top-[-10%] left-[-10%] w-[500px] h-[500px] rounded-full bg-indigo-900/10 blur-[120px] pointer-events-none"></div>
    <div class="absolute bottom-[-10%] right-[-10%] w-[500px] h-[500px] rounded-full bg-purple-900/10 blur-[120px] pointer-events-none"></div>

    <!-- Navigation Sidebar -->
    <aside class="w-full md:w-64 glass-card md:min-h-screen p-6 flex flex-col justify-between border-b md:border-b-0 md:border-r border-slate-800 relative z-20">
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
        <nav class="space-y-1">
          <router-link 
            to="/" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'dashboard' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            📊 Dashboard
          </router-link>
          
          <router-link 
            to="/onboarding" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'onboarding' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            🤖 AI Onboarding
          </router-link>

          <router-link 
            to="/profile" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'profile' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            👤 Profile
          </router-link>

          <router-link 
            to="/assets" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'assets' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            💰 Assets & Debts
          </router-link>

          <router-link 
            to="/goals" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'goals' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            🎯 Goals
          </router-link>

          <router-link 
            to="/portfolio" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'portfolio' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            📈 Portfolio
          </router-link>

          <router-link 
            to="/review" 
            class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition duration-150"
            :class="$route.name === 'review' ? 'bg-indigo-600/20 text-indigo-400 border border-indigo-500/30' : 'text-slate-400 hover:bg-slate-800/30 hover:text-slate-200'"
          >
            🗓️ Monthly Review
          </router-link>
        </nav>
      </div>

      <!-- Footer / Auth trigger mock -->
      <div class="mt-8 border-t border-slate-800/60 pt-4">
        <router-link 
          to="/login"
          class="flex items-center gap-3 px-4 py-2 text-xs font-medium text-slate-400 hover:text-slate-200 transition"
        >
          🔐 Đăng nhập hệ thống
        </router-link>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-grow p-4 md:p-8 relative z-10 overflow-y-auto">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup>
// Layout wrapper
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
