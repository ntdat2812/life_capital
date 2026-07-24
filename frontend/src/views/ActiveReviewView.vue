<template>
  <div class="space-y-6 animate-fade-in pb-20 max-w-4xl mx-auto">
    
    <div class="flex items-center gap-4 mb-8">
      <router-link to="/review" class="p-2 hover:bg-slate-800 rounded-full transition-colors text-slate-400 hover:text-white">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
      </router-link>
      <div>
        <h1 class="text-3xl font-bold text-white tracking-tight">Chốt Sổ Tháng {{ currentMonth }}</h1>
        <p class="text-slate-400">Phân tích danh mục và đề xuất rải ngân bằng AI.</p>
      </div>
    </div>

    <!-- Step 1: Input -->
    <div v-if="step === 1" class="premium-card p-8 text-center animate-fade-in-up">
      <div class="w-20 h-20 bg-indigo-500/20 rounded-full flex items-center justify-center mx-auto mb-6">
        <span class="text-4xl">💰</span>
      </div>
      <h2 class="text-2xl font-bold text-white mb-4">Vốn Rải Ngân Mới</h2>
      <p class="text-slate-400 mb-8 max-w-lg mx-auto">
        Nhập số tiền vốn mới mà bạn dự định tiết kiệm hoặc đầu tư trong tháng này. AI sẽ phân tích và đề xuất phân bổ tối ưu nhất cho bạn.
      </p>

      <form @submit.prevent="generateReview" class="max-w-md mx-auto space-y-6">
        <div>
          <CurrencyInput v-model="newInvestmentAmount" required />
          <p class="text-sm text-slate-500 mt-2 text-left">Gợi ý: Nhập 0 nếu bạn không có vốn mới và chỉ muốn nhận đánh giá danh mục.</p>
        </div>

        <button type="submit" class="w-full btn-primary py-4 text-lg font-bold" :disabled="reviewStore.generating">
          {{ reviewStore.generating ? 'AI đang phân tích danh mục...' : 'Bắt đầu phân tích' }}
        </button>
      </form>
      
      <div v-if="reviewStore.error" class="mt-4 p-4 bg-rose-500/10 border border-rose-500/20 rounded-xl text-rose-400 text-left">
        {{ reviewStore.error }}
      </div>
    </div>

    <!-- Step 2: Generating -->
    <div v-else-if="step === 2 && reviewStore.generating" class="premium-card p-12 text-center animate-pulse">
      <div class="w-16 h-16 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin mx-auto mb-6"></div>
      <h2 class="text-2xl font-bold text-white mb-2">Đang Phân Tích...</h2>
      <p class="text-slate-400">AI đang quét qua IPS, Danh mục đầu tư, Khoản nợ và Danh sách theo dõi để tìm ra giải pháp tối ưu nhất cho bạn.</p>
    </div>

    <!-- Step 3: Result -->
    <div v-else-if="step === 3 && reviewStore.currentReview" class="space-y-6 animate-fade-in-up">
      <!-- Summary Card -->
      <div class="premium-card p-8">
        <h2 class="text-xl font-bold text-white mb-4 flex items-center gap-2">
          <span class="text-indigo-400">🧠</span> Tổng Quan Đánh Giá AI
        </h2>
        <div class="prose prose-invert max-w-none text-slate-300">
          <p class="whitespace-pre-wrap leading-relaxed text-lg">{{ reviewStore.currentReview.ai_overall_summary }}</p>
        </div>
      </div>

      <!-- Action Plan -->
      <h2 class="text-xl font-bold text-white mt-8 mb-4">Đề Xuất Hành Động</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div 
          v-for="(rec, index) in reviewStore.currentReview.ai_recommendations" 
          :key="index"
          class="glass-card p-6 rounded-xl border-l-4"
          :class="{
            'border-emerald-500': rec.action === 'buy' || rec.action === 'add',
            'border-rose-500': rec.action === 'sell' || rec.action === 'reduce',
            'border-amber-500': rec.action === 'hold',
            'border-indigo-500': rec.action === 'pay_debt'
          }"
        >
          <div class="flex justify-between items-start mb-2">
            <h3 class="font-bold text-lg text-white">{{ rec.ticker }}</h3>
            <span 
              class="px-2 py-1 rounded-md text-xs font-bold uppercase"
              :class="{
                'bg-emerald-500/20 text-emerald-400': rec.action === 'buy' || rec.action === 'add',
                'bg-rose-500/20 text-rose-400': rec.action === 'sell' || rec.action === 'reduce',
                'bg-amber-500/20 text-amber-400': rec.action === 'hold',
                'bg-indigo-500/20 text-indigo-400': rec.action === 'pay_debt'
              }"
            >
              {{ getActionLabel(rec.action) }}
            </span>
          </div>
          
          <div v-if="rec.amount > 0" class="text-2xl font-bold mb-3"
            :class="{
              'text-emerald-400': rec.action === 'buy' || rec.action === 'add',
              'text-rose-400': rec.action === 'sell' || rec.action === 'reduce',
              'text-indigo-400': rec.action === 'pay_debt'
            }"
          >
            {{ formatCurrency(rec.amount) }}
          </div>
          
          <p class="text-sm text-slate-400 bg-slate-900/50 p-3 rounded-lg">{{ rec.reason }}</p>
        </div>
      </div>
      
      <!-- Personal Note Section -->
      <div v-if="reviewStore.currentReview.status === 'draft'" class="premium-card p-8 mt-6">
        <h2 class="text-xl font-bold text-white mb-4">📝 Ghi Chú Cá Nhân</h2>
        <textarea
          v-model="userNote"
          rows="3"
          class="w-full bg-slate-900 border border-slate-700 rounded-xl p-4 text-white placeholder-slate-500 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-colors"
          placeholder="Bạn có ghi chú gì thêm cho tháng này không? (Ví dụ: Thưởng dự án 20tr, chuẩn bị mua xe...)"
        ></textarea>
      </div>
      
      <div v-if="reviewStore.currentReview.status === 'completed' && reviewStore.currentReview.user_note" class="premium-card p-8 mt-6">
        <h2 class="text-xl font-bold text-white mb-4">📝 Ghi Chú Cá Nhân</h2>
        <p class="text-slate-300 whitespace-pre-wrap">{{ reviewStore.currentReview.user_note }}</p>
      </div>

      <!-- Action Buttons -->
      <div class="flex justify-end gap-4 pt-8">
        <template v-if="reviewStore.currentReview.status === 'draft'">
          <button @click="cancelDraft" class="px-6 py-3 rounded-xl font-bold transition-all text-slate-300 hover:text-white hover:bg-slate-800">
            Hủy Bỏ
          </button>
          <button @click="saveDraft" class="btn-primary" :disabled="reviewStore.loading">
            {{ reviewStore.loading ? 'Đang Lưu...' : 'Lưu Chốt Sổ' }}
          </button>
        </template>
        <template v-else>
          <router-link to="/review" class="btn-primary">
            Quay Lại Lịch Sử
          </router-link>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useReviewStore } from '../stores/reviewStore'
import CurrencyInput from '../components/common/CurrencyInput.vue'

const reviewStore = useReviewStore()
const route = useRoute()
const router = useRouter()
const step = ref(1)
const newInvestmentAmount = ref(null)
const userNote = ref('')

const currentMonth = computed(() => {
  if (route.params.month) {
    const d = new Date(route.params.month)
    return `${d.getMonth() + 1}/${d.getFullYear()}`
  }
  const d = new Date()
  return `${d.getMonth() + 1}/${d.getFullYear()}`
})

onMounted(async () => {
  if (route.params.month) {
    step.value = 2
    await reviewStore.fetchReviewByMonth(route.params.month)
    step.value = 3
  }
})

const generateReview = async () => {
  step.value = 2
  try {
    const amount = newInvestmentAmount.value || 0
    await reviewStore.generateReview(amount)
    userNote.value = '' // Reset note on new generation
    step.value = 3
  } catch (e) {
    step.value = 1
  }
}

const saveDraft = async () => {
  if (!reviewStore.currentReview) return
  try {
    const reviewData = { ...reviewStore.currentReview }
    if (userNote.value.trim()) {
      reviewData.user_note = userNote.value.trim()
    }
    await reviewStore.saveReview(reviewData)
    router.push('/review')
  } catch (e) {
    // Error is handled in store
  }
}

const cancelDraft = () => {
  reviewStore.clearCurrentReview()
  router.push('/review')
}

const getActionLabel = (action) => {
  const map = {
    'buy': 'Mua Mới',
    'add': 'Mua Thêm',
    'sell': 'Bán Toàn Bộ',
    'reduce': 'Bán Bớt',
    'hold': 'Nắm Giữ',
    'pay_debt': 'Trả Nợ'
  }
  return map[action] || action
}

const formatCurrency = (value) => {
  if (value === null || value === undefined) return '0 ₫'
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(Math.round(value))
}
</script>
