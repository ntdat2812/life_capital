<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h1 class="text-3xl font-bold font-outfit text-white">Investment Thesis</h1>
        <p class="text-slate-400 mt-1">Ghi lại lý do nắm giữ cho từng tài sản để giữ vững kỷ luật và tránh cảm xúc.</p>
      </div>
      <router-link to="/thesis/new" class="btn-primary flex items-center gap-2">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        <span>Tạo Luận điểm mới</span>
      </router-link>
    </div>

    <!-- Error State -->
    <div v-if="portfolioStore.error" class="p-4 bg-rose-500/10 border border-rose-500/20 rounded-xl flex items-start gap-3">
      <svg class="w-5 h-5 text-rose-400 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      <p class="text-rose-400">{{ portfolioStore.error }}</p>
    </div>

    <!-- Loading State -->
    <div v-if="portfolioStore.loading && portfolioStore.theses.length === 0" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500"></div>
    </div>

    <!-- Empty State -->
    <div v-else-if="portfolioStore.theses.length === 0" class="glass-card rounded-2xl p-12 text-center border-dashed border-2 border-slate-700">
      <div class="w-16 h-16 bg-slate-800 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
      </div>
      <h3 class="text-xl font-bold text-white mb-2">Chưa có Luận điểm nào</h3>
      <p class="text-slate-400 mb-6 max-w-md mx-auto">Bạn chưa ghi lại lý do mua tài sản nào. Hãy bắt đầu viết Thesis đầu tiên để định hình chiến lược dài hạn nhé.</p>
      <router-link to="/thesis/new" class="btn-primary inline-flex">
        Tạo Thesis Đầu Tiên
      </router-link>
    </div>

    <!-- Grid View -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="thesis in portfolioStore.theses" :key="thesis.id" class="premium-card rounded-2xl p-6 group hover:border-indigo-500/50 transition-all flex flex-col h-full cursor-pointer relative" @click="$router.push(`/thesis/${thesis.ticker}`)">
        
        <!-- Header -->
        <div class="flex justify-between items-start mb-4">
          <div class="flex-1 min-w-0 mr-4">
            <h3 class="text-2xl font-bold text-white truncate" :title="thesis.ticker || thesis.asset_name">{{ thesis.ticker || thesis.asset_name }}</h3>
            <p class="text-sm text-slate-400">{{ thesis.asset_class }}</p>
          </div>
          <div class="flex-shrink-0 flex flex-col items-center justify-center w-12 h-12 rounded-full border-2 border-indigo-500/30 bg-slate-800/50">
            <span class="text-lg font-bold" :class="getScoreClass(thesis.conviction_score)">{{ thesis.conviction_score }}</span>
          </div>
        </div>

        <!-- Content -->
        <div class="flex-1 space-y-4">
          <div>
            <h4 class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Why I Own</h4>
            <p class="text-sm text-slate-300 line-clamp-3">{{ thesis.why_i_own }}</p>
          </div>
          <div v-if="thesis.catalysts">
            <h4 class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Catalysts (Động lực)</h4>
            <p class="text-sm text-slate-400 line-clamp-2">{{ thesis.catalysts }}</p>
          </div>
        </div>

        <!-- Footer -->
        <div class="mt-6 pt-4 border-t border-slate-800/50 flex justify-between items-center">
          <span class="text-xs text-slate-500">
            Cập nhật: {{ new Date(thesis.updated_at).toLocaleDateString('vi-VN') }}
          </span>
          <span class="text-sm font-medium text-indigo-400 group-hover:text-indigo-300 flex items-center gap-1">
            Xem chi tiết 
            <svg class="w-4 h-4 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { usePortfolioStore } from '../stores/portfolioStore'
import { useRouter } from 'vue-router'

const portfolioStore = usePortfolioStore()
const router = useRouter()

onMounted(() => {
  portfolioStore.fetchTheses()
})

const getScoreClass = (score) => {
  if (score >= 8) return 'text-emerald-400'
  if (score >= 5) return 'text-amber-400'
  return 'text-rose-400'
}
</script>
