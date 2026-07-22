<template>
  <div class="space-y-6 max-w-4xl mx-auto">
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <div class="flex items-center gap-4">
        <button @click="$router.back()" class="p-2 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
        </button>
        <div>
          <h1 class="text-3xl font-bold font-outfit text-white">
            {{ isEditing ? 'Cập nhật Luận điểm' : 'Tạo Luận điểm mới' }}
          </h1>
          <p class="text-slate-400 mt-1">
            {{ isEditing ? `Đang chỉnh sửa luận điểm cho ${thesisForm.ticker || thesisForm.asset_name}` : 'Ghi lại lý do đầu tư để tránh quyết định theo cảm xúc.' }}
          </p>
        </div>
      </div>
      <div class="flex items-center gap-3">
        <button @click="generateAI" v-if="!isEditing || !thesisForm.why_i_own" class="bg-indigo-600/20 text-indigo-400 border border-indigo-500/30 hover:bg-indigo-600/30 px-4 py-2 rounded-lg font-medium transition-all flex items-center gap-2" :disabled="isGenerating">
          <svg v-if="isGenerating" class="animate-spin -ml-1 mr-2 h-4 w-4 text-indigo-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span v-else>✨</span>
          {{ isGenerating ? 'AI Đang viết...' : 'Nhờ AI viết hộ' }}
        </button>
        <button @click="saveThesis" class="btn-primary" :disabled="isSaving || isGenerating">
          {{ isSaving ? 'Đang lưu...' : 'Lưu Luận Điểm' }}
        </button>
      </div>
    </div>

    <!-- Error state -->
    <div v-if="error" class="p-4 bg-rose-500/10 border border-rose-500/20 rounded-xl flex items-start gap-3">
      <svg class="w-5 h-5 text-rose-400 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      <p class="text-rose-400">{{ error }}</p>
    </div>

    <!-- Loading state for edit -->
    <div v-if="isLoading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500"></div>
    </div>

    <form v-else @submit.prevent="saveThesis" class="space-y-6">
      
      <!-- Basic Info Section -->
      <div class="glass-card rounded-2xl p-6 space-y-4">
        <h2 class="text-lg font-bold text-white mb-4 border-b border-slate-800/50 pb-2">Thông Tin Cơ Bản</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Mã Tài Sản (Ticker)</label>
            <input type="text" v-model="thesisForm.ticker" placeholder="VD: FPT, HPG, BTC" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2.5 text-white focus:outline-none focus:border-indigo-500 uppercase disabled:opacity-50 disabled:cursor-not-allowed" :disabled="isAssetLocked" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Tên Tài Sản / Công Ty</label>
            <input type="text" v-model="thesisForm.asset_name" placeholder="VD: Tập đoàn FPT" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2.5 text-white focus:outline-none focus:border-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed" :disabled="isAssetLocked" />
          </div>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Loại Tài Sản</label>
            <select v-model="thesisForm.asset_class" class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2.5 text-white focus:outline-none focus:border-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed" :disabled="isAssetLocked">
              <option value="stock">Cổ phiếu</option>
              <option value="crypto">Tiền điện tử</option>
              <option value="real_estate">Bất động sản</option>
              <option value="gold">Vàng</option>
              <option value="fund">Chứng chỉ quỹ</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Trạng thái</label>
            <select v-model="thesisForm.status" class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2.5 text-white focus:outline-none focus:border-indigo-500">
              <option value="active">Đang nắm giữ (Active)</option>
              <option value="watching">Đang theo dõi (Watching)</option>
              <option value="sold">Đã bán (Sold)</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1 flex justify-between">
              <span>Độ tự tin (1-10)</span>
              <span class="text-indigo-400 font-bold">{{ thesisForm.conviction_score }}</span>
            </label>
            <input type="range" min="1" max="10" v-model="thesisForm.conviction_score" class="w-full mt-2 accent-indigo-500" />
          </div>
        </div>
      </div>

      <!-- Core Thesis Section -->
      <div class="glass-card rounded-2xl p-6 space-y-6">
        <h2 class="text-lg font-bold text-white mb-4 border-b border-slate-800/50 pb-2">Nội Dung Luận Điểm (The Thesis)</h2>
        
        <div>
          <label class="flex items-center gap-2 text-sm font-bold text-emerald-400 mb-2">
            <span>🎯</span> Why I Own (Lý do cốt lõi)
          </label>
          <p class="text-xs text-slate-400 mb-2">Tóm tắt 1-2 câu lý do tại sao bạn quyết định mua tài sản này.</p>
          <textarea v-model="thesisForm.why_i_own" rows="2" required placeholder="VD: FPT là doanh nghiệp công nghệ số 1 VN, đang hưởng lợi lớn từ làn sóng AI và bán dẫn toàn cầu." class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-emerald-500 resize-none"></textarea>
        </div>

        <div>
          <label class="flex items-center gap-2 text-sm font-bold text-blue-400 mb-2">
            <span>🏰</span> Moats (Lợi thế cạnh tranh)
          </label>
          <p class="text-xs text-slate-400 mb-2">Điều gì bảo vệ doanh nghiệp/tài sản này khỏi các đối thủ? (Thương hiệu, Chi phí chuyển đổi, Độc quyền...)</p>
          <textarea v-model="thesisForm.moat" rows="4" placeholder="Nhập các lợi thế cạnh tranh..." class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-blue-500"></textarea>
        </div>

        <div>
          <label class="flex items-center gap-2 text-sm font-bold text-amber-400 mb-2">
            <span>🚀</span> Catalysts (Động lực tăng trưởng)
          </label>
          <p class="text-xs text-slate-400 mb-2">Sự kiện hoặc xu hướng nào trong 1-3 năm tới sẽ làm giá tài sản này tăng vọt?</p>
          <textarea v-model="thesisForm.catalysts" rows="4" placeholder="Nhập động lực tăng trưởng..." class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-amber-500"></textarea>
        </div>

        <div>
          <label class="flex items-center gap-2 text-sm font-bold text-rose-400 mb-2">
            <span>⚠️</span> Risks (Rủi ro)
          </label>
          <p class="text-xs text-slate-400 mb-2">Điều tồi tệ nhất có thể xảy ra và phá hỏng mọi kỳ vọng của bạn là gì?</p>
          <textarea v-model="thesisForm.risks" rows="4" placeholder="Nhập rủi ro tiềm ẩn..." class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-rose-500"></textarea>
        </div>
      </div>

      <!-- Sell Rules Section -->
      <div class="glass-card rounded-2xl p-6 border-l-4 border-rose-500">
        <h2 class="text-lg font-bold text-white mb-2 flex items-center gap-2">
          <span>🛑</span> Nguyên Tắc Bán (Sell Rules)
        </h2>
        <p class="text-sm text-slate-400 mb-4">Các quy tắc cứng bạn tự đặt ra. Nếu vi phạm, BẮT BUỘC PHẢI BÁN không cảm xúc.</p>
        <textarea v-model="thesisForm.sell_conditions" rows="4" required placeholder="VD: 1. Biên lợi nhuận giảm 2 quý liên tiếp. 2. Ban lãnh đạo bán chui cổ phiếu. 3. Mất hợp đồng trọng điểm." class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-rose-500"></textarea>
      </div>
      
      <!-- Hidden submit button to allow form submission via Enter or JS -->
      <button type="submit" class="hidden"></button>
    </form>
    
    <!-- Delete button (only in edit mode) -->
    <div v-if="isEditing" class="flex justify-end pt-4 border-t border-slate-800">
      <button @click="deleteThesis" class="text-rose-400 hover:text-rose-300 text-sm font-medium transition-colors" :disabled="isSaving">
        Xóa Luận điểm này
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePortfolioStore } from '../stores/portfolioStore'

const route = useRoute()
const router = useRouter()
const portfolioStore = usePortfolioStore()

const isEditing = ref(route.name === 'thesis-detail')
const tickerParam = computed(() => route.params.ticker)

const isLoading = ref(false)
const isSaving = ref(false)
const isGenerating = ref(false)
const error = ref(null)

const thesisForm = ref({
  ticker: '',
  asset_name: '',
  asset_class: 'stock',
  why_i_own: '',
  moat: '',
  catalysts: '',
  risks: '',
  sell_conditions: '',
  conviction_score: 7,
  status: 'active'
})

let originalThesisId = null

const isAssetLocked = computed(() => !!tickerParam.value)

onMounted(async () => {
  if (portfolioStore.holdings.length === 0) {
    await portfolioStore.fetchHoldings()
  }
  if (portfolioStore.watchlist.length === 0) {
    await portfolioStore.fetchWatchlist()
  }

  if (tickerParam.value) {
    isLoading.value = true
    try {
      const data = await portfolioStore.getThesisByTicker(tickerParam.value)
      if (data) {
        originalThesisId = data.id
        let derivedAssetClass = 'stock'
        const holding = portfolioStore.holdings.find(a => a.ticker === tickerParam.value || a.name === tickerParam.value)
        const watching = portfolioStore.watchlist.find(w => w.ticker === tickerParam.value)
        if (holding && holding.category) {
          derivedAssetClass = holding.category
        }

        thesisForm.value = { 
          ...data,
          asset_class: derivedAssetClass,
          moat: Array.isArray(data.moat) ? data.moat.join('\n') : (data.moat || ''),
          catalysts: Array.isArray(data.catalysts) ? data.catalysts.join('\n') : (data.catalysts || ''),
          risks: Array.isArray(data.risks) ? data.risks.join('\n') : (data.risks || ''),
          sell_conditions: Array.isArray(data.sell_conditions) ? data.sell_conditions.join('\n') : (data.sell_conditions || '')
        }
      } else {
        // Not found - switch to create mode and prefill
        isEditing.value = false
        thesisForm.value.ticker = tickerParam.value
        
        // Cố gắng tìm trong danh mục hoặc watchlist để prefill Tên & Loại
        const holding = portfolioStore.holdings.find(a => a.ticker === tickerParam.value || a.name === tickerParam.value)
        const watching = portfolioStore.watchlist.find(w => w.ticker === tickerParam.value)
        
        if (holding) {
          thesisForm.value.asset_name = holding.name || holding.ticker
          thesisForm.value.asset_class = holding.category || 'stock'
        } else if (watching) {
          thesisForm.value.asset_name = watching.company_name || watching.ticker
          thesisForm.value.asset_class = 'stock'
        } else {
          thesisForm.value.asset_name = tickerParam.value
        }
      }
    } catch (err) {
      error.value = 'Lỗi tải dữ liệu. Vui lòng thử lại.'
    } finally {
      isLoading.value = false
    }
  }
})

const generateAI = async () => {
  if (!thesisForm.value.ticker && !thesisForm.value.asset_name) {
    error.value = "Vui lòng nhập Ticker hoặc Tên tài sản trước khi nhờ AI viết."
    return
  }
  
  error.value = null
  isGenerating.value = true
  
  try {
    const data = await portfolioStore.generateThesisAI(
      thesisForm.value.ticker, 
      thesisForm.value.asset_class, 
      thesisForm.value.asset_name
    )
    
    // Auto fill form with AI data
    thesisForm.value.why_i_own = data.why_i_own || ''
    thesisForm.value.thesis_summary = data.thesis_summary || ''
    thesisForm.value.thesis_detail = data.thesis_detail || ''
    thesisForm.value.moat = Array.isArray(data.moat) ? data.moat.join('\n') : (data.moat || '')
    thesisForm.value.catalysts = Array.isArray(data.catalysts) ? data.catalysts.join('\n') : (data.catalysts || '')
    thesisForm.value.risks = Array.isArray(data.risks) ? data.risks.join('\n') : (data.risks || '')
    thesisForm.value.sell_conditions = Array.isArray(data.sell_conditions) ? data.sell_conditions.join('\n') : (data.sell_conditions || '')
    
    if (data.conviction_score) {
      thesisForm.value.conviction_score = data.conviction_score
    }
    
  } catch (err) {
    error.value = "AI gặp lỗi trong quá trình phân tích. Vui lòng thử lại sau."
  } finally {
    isGenerating.value = false
  }
}

const saveThesis = async () => {
  if (!thesisForm.value.ticker || !thesisForm.value.why_i_own) {
    error.value = "Vui lòng điền đủ Ticker và Why I Own."
    return
  }

  error.value = null
  isSaving.value = true
  
  try {
    // Upper case ticker
    thesisForm.value.ticker = thesisForm.value.ticker.toUpperCase()
    
    // Prepare payload for backend (convert string to array for JSONB fields)
    const payload = {
      ...thesisForm.value,
      moat: thesisForm.value.moat ? thesisForm.value.moat.split('\n').filter(i => i.trim()) : [],
      catalysts: thesisForm.value.catalysts ? thesisForm.value.catalysts.split('\n').filter(i => i.trim()) : [],
      risks: thesisForm.value.risks ? thesisForm.value.risks.split('\n').filter(i => i.trim()) : [],
      sell_conditions: thesisForm.value.sell_conditions ? thesisForm.value.sell_conditions.split('\n').filter(i => i.trim()) : []
    }
    
    if (isEditing.value && originalThesisId) {
      await portfolioStore.updateThesis(originalThesisId, payload)
    } else {
      await portfolioStore.createThesis(payload)
    }
    router.push('/thesis')
  } catch (err) {
    error.value = portfolioStore.error || "Lỗi khi lưu luận điểm."
  } finally {
    isSaving.value = false
  }
}

const deleteThesis = async () => {
  if (!originalThesisId) return
  if (confirm(`Bạn có chắc muốn xóa vĩnh viễn luận điểm này?`)) {
    isSaving.value = true
    try {
      await portfolioStore.deleteThesis(originalThesisId)
      router.push('/thesis')
    } catch (err) {
      error.value = "Lỗi khi xóa luận điểm."
    } finally {
      isSaving.value = false
    }
  }
}
</script>
