<template>
  <div class="space-y-6 pb-20 max-w-5xl mx-auto">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold font-outfit text-white">Chiến Lược Đầu Tư (IPS)</h1>
        <p class="text-xs text-slate-400 mt-1">Bản hiến pháp định hướng quản lý tài sản của bạn</p>
      </div>
      <button 
        v-if="ipsStore.latestIps"
        @click="openRegenerateConfirm"
        :disabled="ipsStore.isGenerating"
        class="px-4 py-2 text-sm bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white rounded-xl transition shadow-lg shadow-indigo-500/20 font-medium flex items-center"
      >
        <span v-if="ipsStore.isGenerating" class="mr-2">
          <svg class="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </span>
        ✨ Nhờ AI phân bổ lại
      </button>
    </div>

    <div v-if="ipsStore.isLoading" class="text-center py-20 text-slate-400">
      Đang tải dữ liệu...
    </div>

    <!-- Generating State (When generating or regenerating) -->
    <div v-else-if="ipsStore.isGenerating" class="text-center py-24 bg-slate-900/50 rounded-3xl border border-slate-800 flex flex-col items-center justify-center">
      <div class="relative w-24 h-24 mb-6">
        <div class="absolute inset-0 border-4 border-indigo-500/20 rounded-full"></div>
        <div class="absolute inset-0 border-4 border-indigo-500 rounded-full border-t-transparent animate-spin"></div>
        <div class="absolute inset-0 flex items-center justify-center text-indigo-400">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z"></path></svg>
        </div>
      </div>
      <h3 class="text-xl font-bold font-outfit text-white mb-2">Đang phân tích dữ liệu...</h3>
      <p class="text-slate-400 text-sm max-w-sm">AI Wealth Manager đang xây dựng chiến lược cá nhân hóa dựa trên hồ sơ của bạn. Vui lòng đợi trong giây lát.</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="!ipsStore.latestIps && !ipsStore.isGenerating" class="text-center py-16 bg-slate-900/50 rounded-3xl border border-slate-800">
      <div class="w-20 h-20 bg-indigo-500/10 rounded-full flex items-center justify-center mx-auto mb-6 relative">
        <svg class="w-10 h-10 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
        </svg>
      </div>
      <h3 class="text-2xl font-bold font-outfit text-white mb-3">Chưa có chiến lược đầu tư (IPS)</h3>
      <p class="text-slate-400 text-sm max-w-md mx-auto mb-4 leading-relaxed">
        Để hệ thống tạo ra một bản hiến pháp đầu tư chuẩn xác nhất, bạn nên chắc chắn mình đã:
        <br/><span class="text-indigo-400 font-medium">1. Lập Hồ sơ Rủi ro</span> và 
        <span class="text-indigo-400 font-medium">2. Có Mục tiêu tài chính rõ ràng</span>.
      </p>
      
      <button @click="openPreSurvey" class="mt-4 px-8 py-3 bg-indigo-600 hover:bg-indigo-500 shadow-lg shadow-indigo-500/25 text-white rounded-xl transition font-semibold">
        Tạo Chiến Lược Mới Ngay
      </button>
    </div>

    <!-- Active State -->
    <div v-else-if="ipsStore.latestIps" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Cột Trái: Target Allocation & Trạng thái -->
      <div class="lg:col-span-1 space-y-6">
        <div class="premium-card rounded-3xl p-6 relative overflow-hidden group">
          <div class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-indigo-500 to-purple-500"></div>
          <div class="flex justify-between items-center mb-6">
            <h3 class="text-white font-medium flex items-center">
              <svg class="w-5 h-5 mr-2 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 3.055A9.001 9.001 0 1020.945 13H11V3.055z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.488 9H15V3.512A9.025 9.025 0 0120.488 9z"></path></svg>
              Tỷ Trọng (Target)
            </h3>
            <button @click="openEditModal" class="p-2 text-slate-400 hover:text-white bg-slate-800/50 hover:bg-slate-700 rounded-lg transition" title="Chỉnh sửa tỷ trọng bằng tay">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
            </button>
          </div>
          
          <div class="space-y-4">
            <div v-for="(val, key) in validAllocations" :key="key" class="bg-slate-900/40 p-3.5 rounded-xl border border-slate-800/60 hover:border-slate-700 transition">
              <div class="flex justify-between items-center mb-2">
                <span class="text-sm font-medium text-slate-300 capitalize">{{ formatAssetName(key) }}</span>
                <span class="text-sm font-bold text-white">{{ val }}%</span>
              </div>
              <div class="w-full bg-slate-800 rounded-full h-2">
                <div class="bg-indigo-500 h-2 rounded-full" :style="{ width: `${val}%` }"></div>
              </div>
            </div>
          </div>
          
          <div class="mt-5 p-3 rounded-xl bg-indigo-500/10 border border-indigo-500/20 text-xs text-indigo-200 leading-relaxed">
            <span class="font-bold">Lưu ý:</span> Việc tự ý thay đổi tỷ trọng sẽ ảnh hưởng trực tiếp đến Cảnh báo (Alerts) của hệ thống khi tái cân bằng định kỳ.
          </div>
        </div>

        <div class="premium-card rounded-3xl p-6">
           <h3 class="text-slate-400 text-xs font-medium uppercase tracking-wider mb-2">Trạng thái Bản Hiến Pháp</h3>
           <div class="text-emerald-400 text-sm font-bold flex items-center mb-4">
             <span class="relative flex h-3 w-3 mr-2">
               <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
               <span class="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
             </span>
             Đang Áp Dụng (Active)
           </div>
           
           <div class="space-y-2 text-sm text-slate-300">
             <div class="flex justify-between">
               <span class="text-slate-500">Phiên bản:</span>
               <span>v{{ ipsStore.latestIps.version }}</span>
             </div>
             <div class="flex justify-between">
               <span class="text-slate-500">Nguồn tạo:</span>
               <span>{{ ipsStore.latestIps.is_ai_recommended ? 'AI Wealth Manager' : 'Chỉnh sửa thủ công' }}</span>
             </div>
             <div class="flex justify-between">
               <span class="text-slate-500">Cập nhật lúc:</span>
               <span>{{ new Date(ipsStore.latestIps.updated_at).toLocaleDateString('vi-VN') }}</span>
             </div>
           </div>
        </div>
      </div>

      <!-- Cột Phải: Bài viết chi tiết Markdown -->
      <div class="lg:col-span-2">
        <div class="premium-card rounded-3xl p-8 h-full">
          <div class="flex items-center mb-6 border-b border-slate-800 pb-5">
             <div class="w-12 h-12 rounded-2xl bg-indigo-500/10 flex items-center justify-center mr-4 border border-indigo-500/20">
               <svg class="w-6 h-6 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"></path></svg>
             </div>
             <div>
               <h2 class="text-xl font-bold font-outfit text-white">Tư Vấn Chiến Lược Chuyên Sâu</h2>
               <p class="text-sm text-slate-400 mt-1">Goal-Aware Unified Strategy</p>
             </div>
             <div class="ml-auto">
               <button @click="openEditStrategyModal" class="px-4 py-2 text-sm bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white rounded-xl transition flex items-center border border-slate-700">
                 <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
                 Sửa văn bản (Nâng cao)
               </button>
             </div>
          </div>
          
          <div class="prose prose-invert prose-slate prose-sm md:prose-base max-w-none markdown-body" v-html="formattedStrategy">
          </div>
        </div>
      </div>
    </div>

    <!-- Modals -->

    <!-- 1. Pre-survey Modal (Lựa chọn lớp tài sản) -->
    <div v-if="showSurveyModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-950/80 backdrop-blur-sm" @click="showSurveyModal = false"></div>
      <div class="bg-slate-900 border border-slate-700 p-8 rounded-3xl shadow-2xl relative z-10 w-full max-w-md animate-fade-in-up">
        <h3 class="text-xl font-bold text-white mb-2 font-outfit">Sở Thích Đầu Tư</h3>
        <p class="text-slate-400 text-sm mb-6 leading-relaxed">
          Đánh dấu vào những loại tài sản mà bạn HIỂU và MUỐN AI đưa vào trong danh mục phân bổ của mình. 
          <br/><span class="text-indigo-400 font-medium">Nếu bạn chưa rõ, hãy bỏ trống để AI tự động tối ưu giúp bạn!</span>
        </p>
        
        <div class="space-y-3 mb-8">
          <label v-for="asset in availableAssetClasses" :key="asset.id" class="flex items-center p-3 border border-slate-700 rounded-xl cursor-pointer hover:bg-slate-800 transition" :class="{'bg-indigo-900/20 border-indigo-500/50': selectedAssets.includes(asset.id)}">
            <input type="checkbox" :value="asset.id" v-model="selectedAssets" class="w-5 h-5 rounded border-slate-600 text-indigo-600 bg-slate-800 focus:ring-indigo-500 focus:ring-offset-slate-900">
            <div class="ml-3">
              <span class="block text-sm font-medium text-slate-200">{{ asset.name }}</span>
              <span class="block text-xs text-slate-500">{{ asset.desc }}</span>
            </div>
          </label>
        </div>

        <div class="flex justify-end gap-3">
          <button @click="showSurveyModal = false" class="px-5 py-2.5 text-sm font-medium text-slate-300 hover:text-white transition">Hủy</button>
          <button @click="executeGenerate" class="px-5 py-2.5 bg-indigo-600 hover:bg-indigo-500 text-white rounded-xl font-medium transition shadow-lg shadow-indigo-500/20">Bắt đầu Phân Tích</button>
        </div>
      </div>
    </div>

    <!-- 2. Danger Confirm Modal -->
    <div v-if="showConfirmModal" class="fixed inset-0 z-[60] flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-950/80 backdrop-blur-sm" @click="showConfirmModal = false"></div>
      <div class="bg-slate-900 border border-red-900/50 p-8 rounded-3xl shadow-2xl relative z-10 w-full max-w-md animate-fade-in-up">
        <div class="w-16 h-16 bg-red-500/10 rounded-full flex items-center justify-center mb-6 border border-red-500/20">
          <svg class="w-8 h-8 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
        </div>
        <h3 class="text-xl font-bold text-white mb-2 font-outfit">Bạn có chắc chắn?</h3>
        <p class="text-slate-400 text-sm mb-4 leading-relaxed">
          IPS là **Chiến lược xuyên suốt dài hạn**. Việc thay đổi này sẽ định hình lại toàn bộ mục tiêu và các quyết định đầu tư của bạn trong tương lai.
        </p>
        <p class="text-slate-400 text-sm mb-6">
          Vui lòng gõ chữ <strong class="text-red-400">OK</strong> để tiếp tục:
        </p>
        
        <input 
          v-model="confirmText" 
          type="text" 
          placeholder="OK"
          class="w-full bg-slate-950 border border-slate-700 rounded-xl px-4 py-3 text-white mb-8 focus:border-red-500 focus:ring-1 focus:ring-red-500 outline-none uppercase font-mono tracking-widest text-center"
        >

        <div class="flex justify-end gap-3">
          <button @click="showConfirmModal = false" class="px-5 py-2.5 text-sm font-medium text-slate-300 hover:text-white transition">Giữ nguyên</button>
          <button 
            @click="proceedConfirmAction" 
            :disabled="confirmText.toUpperCase() !== 'OK'"
            class="px-5 py-2.5 bg-red-600 hover:bg-red-500 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-xl font-medium transition"
          >
            Tiếp tục sửa đổi
          </button>
        </div>
      </div>
    </div>

    <!-- 3. Edit Allocation Modal -->
    <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-950/80 backdrop-blur-sm" @click="showEditModal = false"></div>
      <div class="bg-slate-900 border border-slate-700 p-8 rounded-3xl shadow-2xl relative z-10 w-full max-w-lg animate-fade-in-up">
        <h3 class="text-xl font-bold text-white mb-2 font-outfit">Chỉnh Sửa Tỷ Trọng</h3>
        <p class="text-slate-400 text-sm mb-6">Bạn có thể can thiệp bằng tay nếu cảm thấy đề xuất của AI chưa sát với ý muốn.</p>
        
        <div class="space-y-5 mb-8">
          <div v-for="(val, key) in editAllocations" :key="key">
            <div class="flex justify-between items-center mb-2">
              <label class="text-sm font-medium text-slate-300 capitalize">{{ formatAssetName(key) }}</label>
            </div>
            <div class="flex items-center">
              <input 
                type="number" 
                min="0" 
                max="100" 
                v-model.number="editAllocations[key]" 
                class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-2 text-white focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 outline-none"
              >
              <span class="ml-3 text-slate-400 font-bold">%</span>
            </div>
          </div>
        </div>

        <div class="p-4 rounded-xl mb-8 flex justify-between items-center border" :class="totalAllocation === 100 ? 'bg-emerald-500/10 border-emerald-500/30' : 'bg-amber-500/10 border-amber-500/30'">
          <span class="text-sm font-medium" :class="totalAllocation === 100 ? 'text-emerald-400' : 'text-amber-400'">
            {{ totalAllocation === 100 ? 'Tổng Tỷ Trọng Hợp Lệ' : (totalAllocation > 100 ? 'Vượt quá 100%' : 'Chưa phân bổ hết 100%') }}:
          </span>
          <span class="text-lg font-bold" :class="totalAllocation === 100 ? 'text-emerald-400' : 'text-amber-400'">{{ totalAllocation }}%</span>
        </div>

        <div class="flex justify-end gap-3">
          <button @click="showEditModal = false" class="px-5 py-2.5 text-sm font-medium text-slate-300 hover:text-white transition">Hủy</button>
          <button 
            @click="triggerSaveAllocation" 
            :disabled="totalAllocation !== 100"
            class="px-5 py-2.5 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white rounded-xl font-medium transition"
          >
            Lưu thay đổi
          </button>
        </div>
      </div>
    </div>

    <!-- 4. Edit Strategy Modal -->
    <div v-if="showEditStrategyModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-950/80 backdrop-blur-sm" @click="showEditStrategyModal = false"></div>
      <div class="bg-slate-900 border border-slate-700 p-8 rounded-3xl shadow-2xl relative z-10 w-full max-w-4xl max-h-[90vh] flex flex-col animate-fade-in-up">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-xl font-bold text-white font-outfit">Sửa Văn Bản Chiến Lược (Markdown)</h3>
          <button @click="showEditStrategyModal = false" class="text-slate-400 hover:text-white"><svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg></button>
        </div>
        
        <div class="p-3 mb-4 rounded-xl bg-amber-500/10 border border-amber-500/20 text-sm text-amber-200">
          <span class="font-bold">Lưu ý:</span> Hỗ trợ định dạng Markdown. Vui lòng giữ nguyên các tiêu đề (##) để đảm bảo cấu trúc bài phân tích.
        </div>
        
        <textarea 
          v-model="editStrategyText" 
          class="w-full flex-1 bg-slate-950 border border-slate-700 rounded-xl p-4 text-slate-300 font-mono text-sm leading-relaxed focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 outline-none resize-none min-h-[400px] mb-6"
          placeholder="Viết chiến lược của bạn vào đây..."
        ></textarea>

        <div class="flex justify-end gap-3 mt-auto">
          <button @click="showEditStrategyModal = false" class="px-5 py-2.5 text-sm font-medium text-slate-300 hover:text-white transition">Hủy</button>
          <button 
            @click="triggerSaveStrategy" 
            :disabled="!editStrategyText.trim()"
            class="px-5 py-2.5 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white rounded-xl font-medium transition flex items-center"
          >
            Lưu Chiến Lược
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useIpsStore } from '../stores/ips'
import { marked } from 'marked'

const ipsStore = useIpsStore()

onMounted(() => {
  ipsStore.fetchLatestIps()
})

// --- Data & Modals State ---
const showConfirmModal = ref(false)
const confirmText = ref('')
const confirmAction = ref('')

const showSurveyModal = ref(false)
const selectedAssets = ref(['stock', 'cash'])
const availableAssetClasses = [
  { id: 'stock', name: 'Cổ phiếu', desc: 'Rủi ro cao, sinh lời cao dài hạn' },
  { id: 'bond', name: 'Trái phiếu', desc: 'An toàn, sinh lời cố định, ổn định' },
  { id: 'fund', name: 'Chứng chỉ quỹ', desc: 'Đầu tư qua quỹ mở/ETF, rủi ro đa dạng' },
  { id: 'gold', name: 'Vàng', desc: 'Phòng ngừa lạm phát và khủng hoảng' },
  { id: 'real_estate', name: 'Bất động sản', desc: 'Tích sản, cần vốn lớn, thanh khoản chậm' },
  { id: 'crypto', name: 'Tiền số (Crypto)', desc: 'Biến động cực mạnh, rủi ro cực cao' },
  { id: 'cash', name: 'Tiền mặt (Gửi TK)', desc: 'Thanh khoản lập tức, rủi ro bằng 0' }
]

const showEditModal = ref(false)
const editAllocations = ref({})

const showEditStrategyModal = ref(false)
const editStrategyText = ref('')

// --- Logic ---
const openRegenerateConfirm = () => {
  if (ipsStore.latestIps) {
    confirmText.value = ''
    confirmAction.value = 'regenerate'
    showConfirmModal.value = true
  } else {
    openPreSurvey()
  }
}

const proceedConfirmAction = () => {
  showConfirmModal.value = false
  if (confirmAction.value === 'regenerate') {
    openPreSurvey()
  } else if (confirmAction.value === 'save_allocation') {
    saveAllocation()
  } else if (confirmAction.value === 'save_strategy') {
    saveStrategy()
  }
}

const openPreSurvey = () => {
  showSurveyModal.value = true
}

const executeGenerate = async () => {
  showSurveyModal.value = false
  try {
    await ipsStore.generateIps(selectedAssets.value)
  } catch (error) {
    alert("Đã có lỗi xảy ra: " + error)
  }
}

// --- Edit Allocation Logic ---
const openEditModal = () => {
  if (!ipsStore.latestIps) return
  try {
    const alloc = typeof ipsStore.latestIps.target_allocation === 'string' 
      ? JSON.parse(ipsStore.latestIps.target_allocation)
      : ipsStore.latestIps.target_allocation
      
    // Đảm bảo hiển thị cả các lớp tài sản có sẵn nhưng AI đang set là 0%
    const fullAlloc = { ...alloc }
    availableAssetClasses.forEach(asset => {
      if (fullAlloc[asset.id] === undefined) {
        fullAlloc[asset.id] = 0
      }
    })
      
    editAllocations.value = fullAlloc
    showEditModal.value = true
  } catch (e) {
    console.error("Invalid allocation format", e)
  }
}

const totalAllocation = computed(() => {
  return Object.values(editAllocations.value).reduce((sum, val) => sum + Number(val), 0)
})

const triggerSaveAllocation = () => {
  if (totalAllocation.value !== 100) return
  confirmText.value = ''
  confirmAction.value = 'save_allocation'
  showConfirmModal.value = true
}

const saveAllocation = async () => {
  if (totalAllocation.value !== 100) return
  
  try {
    await ipsStore.updateIps(editAllocations.value)
    showEditModal.value = false
  } catch (e) {
    alert(e)
  }
}

const openEditStrategyModal = () => {
  if (!ipsStore.latestIps) return
  editStrategyText.value = ipsStore.latestIps.detailed_strategy
  showEditStrategyModal.value = true
}

const triggerSaveStrategy = () => {
  if (!editStrategyText.value.trim()) return
  confirmText.value = ''
  confirmAction.value = 'save_strategy'
  showConfirmModal.value = true
}

const saveStrategy = async () => {
  if (!editStrategyText.value.trim()) return
  
  try {
    let finalStrategy = editStrategyText.value
    // Replace old signatures if exist
    finalStrategy = finalStrategy.replace(/\n\n---\n\*(Phân tích và tạo tự động|Chỉnh sửa thủ công).*?\*/g, '')
    // Add new manual signature
    const now = new Date().toLocaleString('vi-VN')
    finalStrategy += `\n\n---\n*Chỉnh sửa thủ công bởi người dùng vào lúc ${now}*`

    const alloc = typeof ipsStore.latestIps.target_allocation === 'string'
      ? JSON.parse(ipsStore.latestIps.target_allocation)
      : ipsStore.latestIps.target_allocation
      
    await ipsStore.updateIps(alloc, finalStrategy)
    showEditStrategyModal.value = false
  } catch (e) {
    alert(e)
  }
}

// --- Formatters ---
const formattedStrategy = computed(() => {
  if (!ipsStore.latestIps?.detailed_strategy) return ''
  return marked(ipsStore.latestIps.detailed_strategy)
})

const validAllocations = computed(() => {
  if (!ipsStore.latestIps?.target_allocation) return {}
  const alloc = typeof ipsStore.latestIps.target_allocation === 'string' 
                ? JSON.parse(ipsStore.latestIps.target_allocation) 
                : ipsStore.latestIps.target_allocation
  
  return Object.fromEntries(
    Object.entries(alloc).filter(([_, val]) => val > 0).sort((a, b) => b[1] - a[1])
  )
})

const formatAssetName = (key) => {
  const map = {
    'stock': 'Cổ phiếu',
    'fund': 'Chứng chỉ quỹ',
    'bond': 'Trái phiếu',
    'gold': 'Vàng',
    'cash': 'Tiền mặt',
    'cash_equivalent': 'Tiền mặt / Tương đương tiền',
    'cash_equivalents': 'Tiền mặt / Tương đương tiền',
    'real_estate': 'Bất động sản',
    'crypto': 'Tiền số (Crypto)'
  }
  if (map[key]) return map[key]
  
  // Fallback for random AI generated keys
  return key.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')
}
</script>

<style>
/* Markdown styles */
.markdown-body h1 {
  font-size: 1.5rem;
  color: #fff;
  font-family: 'Outfit', sans-serif;
  margin-top: 1.5em;
  margin-bottom: 1em;
  border-bottom: 1px solid #334155;
  padding-bottom: 0.5em;
}
.markdown-body h2 {
  font-size: 1.25rem;
  color: #818cf8;
  font-family: 'Outfit', sans-serif;
  margin-top: 1.5em;
  margin-bottom: 0.75em;
}
.markdown-body h3 {
  font-size: 1.1rem;
  color: #cbd5e1;
  font-family: 'Outfit', sans-serif;
  margin-top: 1.25em;
  margin-bottom: 0.5em;
}
.markdown-body ul, .markdown-body ol {
  padding-left: 1.5em;
  color: #cbd5e1;
  margin-bottom: 1em;
}
.markdown-body li {
  margin-bottom: 0.5em;
}
.markdown-body strong {
  color: #e2e8f0;
}
.markdown-body p {
  color: #94a3b8;
  line-height: 1.7;
  margin-bottom: 1em;
}
.markdown-body blockquote {
  border-left: 4px solid #6366f1;
  padding-left: 1em;
  color: #cbd5e1;
  font-style: italic;
  background: rgba(99, 102, 241, 0.05);
  padding: 0.5em 1em;
  border-radius: 0 8px 8px 0;
}

/* Animations */
.animate-fade-in-up {
  animation: fadeInUp 0.3s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
