<template>
  <div class="space-y-6 animate-fade-in pb-20">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-3xl font-bold text-white mb-2 tracking-tight">Chốt Sổ Tháng</h1>
        <p class="text-slate-400">Xem lại lịch sử quản lý gia sản và thực hiện báo cáo định kỳ.</p>
      </div>
      
      <router-link to="/review/active" class="btn-primary group shadow-lg shadow-indigo-500/20">
        <span class="mr-2">🚀</span> Bắt đầu chốt sổ tháng này
        <svg class="w-4 h-4 ml-1 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
      </router-link>
    </div>

    <div v-if="reviewStore.loading" class="flex justify-center p-12">
      <div class="w-10 h-10 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
    </div>
    
    <div v-else-if="reviewStore.error" class="p-4 bg-rose-500/10 border border-rose-500/20 rounded-xl text-rose-400">
      {{ reviewStore.error }}
    </div>

    <div v-else-if="reviewStore.history.length === 0" class="premium-card p-12 text-center">
      <div class="w-20 h-20 bg-slate-800 rounded-full flex items-center justify-center mx-auto mb-4">
        <span class="text-4xl">📅</span>
      </div>
      <h3 class="text-xl font-bold text-white mb-2">Chưa có bản chốt sổ nào</h3>
      <p class="text-slate-400 max-w-md mx-auto mb-6">Bạn chưa thực hiện kỳ chốt sổ (Monthly Review) nào. Hãy bắt đầu ngay để nhận được tư vấn phân bổ vốn từ AI.</p>
      <router-link to="/review/active" class="btn-primary inline-flex">
        Bắt đầu chốt sổ đầu tiên
      </router-link>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- History Cards -->
      <router-link 
        v-for="review in sortedHistory" 
        :key="review.id"
        :to="`/review/${review.review_month.split('T')[0]}`"
        class="glass-card p-6 rounded-2xl hover:-translate-y-1 hover:shadow-xl hover:shadow-indigo-500/10 transition-all cursor-pointer group"
      >
        <div class="flex justify-between items-start mb-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-lg bg-indigo-500/20 text-indigo-400 flex items-center justify-center font-bold text-lg">
              T{{ new Date(review.review_month).getMonth() + 1 }}
            </div>
            <div>
              <h3 class="font-bold text-slate-200 group-hover:text-indigo-300 transition-colors">
                Tháng {{ new Date(review.review_month).getMonth() + 1 }}/{{ new Date(review.review_month).getFullYear() }}
              </h3>
              <span class="text-xs text-slate-500">{{ new Date(review.created_at).toLocaleDateString('vi-VN') }}</span>
            </div>
          </div>
          <span v-if="review.status === 'completed'" class="px-2 py-1 bg-emerald-500/20 text-emerald-400 text-xs font-medium rounded-full">
            Hoàn tất
          </span>
          <span v-else class="px-2 py-1 bg-amber-500/20 text-amber-400 text-xs font-medium rounded-full">
            Bản nháp
          </span>
        </div>
        
        <div class="space-y-3 mb-6">
          <div class="flex justify-between items-center text-sm border-b border-slate-700/50 pb-2">
            <span class="text-slate-400">Vốn rải ngân:</span>
            <span class="font-bold text-white">{{ formatCurrency(review.new_investment_amount) }}</span>
          </div>
          <div class="flex justify-between items-center text-sm">
            <span class="text-slate-400">Tổng tài sản ròng:</span>
            <span class="font-bold text-emerald-400">{{ formatCurrency(review.net_worth_at_review) }}</span>
          </div>
        </div>
        
        <p class="text-sm text-slate-400 line-clamp-2 italic">
          "{{ review.ai_overall_summary || 'Không có tóm tắt.' }}"
        </p>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useReviewStore } from '../stores/reviewStore'

const reviewStore = useReviewStore()

onMounted(() => {
  reviewStore.fetchHistory()
})

const sortedHistory = computed(() => {
  if (!reviewStore.history) return []
  return [...reviewStore.history].sort((a, b) => new Date(b.review_month) - new Date(a.review_month))
})

const formatCurrency = (value) => {
  if (value === null || value === undefined) return '0 ₫'
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(Math.round(value))
}
</script>
