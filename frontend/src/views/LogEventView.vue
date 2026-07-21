<template>
  <div class="max-w-3xl mx-auto">
    <div class="mb-6 flex items-center">
      <button @click="$router.push('/profile/timeline')" class="text-gray-400 hover:text-white mr-4 transition-colors">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
      </button>
      <h1 class="text-2xl font-bold text-gray-100">Khai báo sự kiện</h1>
    </div>

    <!-- Step 1: Input -->
    <div v-if="!draft" class="glass-panel p-8 slide-up">
      <div class="mb-6">
        <div class="w-12 h-12 bg-blue-600/20 rounded-full flex items-center justify-center mb-4 border border-blue-500/30 shadow-[0_0_15px_rgba(37,99,235,0.2)]">
          <span class="text-2xl">🤖</span>
        </div>
        <h2 class="text-xl font-bold text-white mb-2">Gần đây bạn có biến động gì không?</h2>
        <p class="text-gray-400 text-sm">Hãy kể cho tôi nghe về những thay đổi trong cuộc sống của bạn (VD: đổi việc, lập gia đình, sinh con, mua nhà...). Hệ thống AI sẽ tự động phân tích tác động tài chính.</p>
      </div>

      <textarea
        v-model="eventText"
        class="w-full bg-gray-900/50 border border-gray-700 rounded-xl p-4 text-white focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all min-h-[120px]"
        placeholder="Ví dụ: Mình chuẩn bị lấy vợ vào tháng tới, dự kiến chi phí sinh hoạt sẽ tăng lên..."
        :disabled="timelineStore.loading"
      ></textarea>

      <div class="mt-6 flex justify-end">
        <button 
          @click="analyze"
          :disabled="!eventText.trim() || timelineStore.loading"
          class="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white px-6 py-2.5 rounded-lg transition-colors flex items-center shadow-[0_0_15px_rgba(37,99,235,0.3)]"
        >
          <span v-if="timelineStore.loading" class="animate-spin -ml-1 mr-2 h-4 w-4 border-b-2 border-white rounded-full"></span>
          <span>Phân tích bằng AI</span>
        </button>
      </div>
    </div>

    <!-- Step 2: Preview & Edit -->
    <div v-else class="space-y-6 slide-up">
      <div class="glass-panel p-6 border border-blue-500/30 shadow-[0_0_30px_rgba(37,99,235,0.1)] relative overflow-hidden">
        <div class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-blue-500 to-purple-500"></div>
        
        <h2 class="text-xl font-bold text-white mb-4 flex items-center">
          <span class="mr-2">✨</span> Bản nháp đề xuất từ AI
        </h2>
        
        <div class="bg-blue-900/20 border border-blue-800/50 rounded-lg p-4 mb-6">
          <p class="text-sm text-blue-200">
            <strong>Phân tích:</strong> {{ draft.ai_impact_analysis }}
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- General -->
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-400 mb-1">Tiêu đề sự kiện</label>
              <input v-model="draft.title" type="text" class="w-full bg-gray-800 border border-gray-700 rounded-lg p-2.5 text-white focus:border-blue-500 outline-none" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-400 mb-1">Danh mục</label>
              <select v-model="draft.category" class="w-full bg-gray-800 border border-gray-700 rounded-lg p-2.5 text-white focus:border-blue-500 outline-none">
                <option value="income_change">Thay đổi thu nhập</option>
                <option value="family_change">Gia đình</option>
                <option value="housing">Nhà ở</option>
                <option value="dependent_change">Người phụ thuộc</option>
                <option value="major_expense">Chi tiêu lớn</option>
                <option value="career">Sự nghiệp</option>
                <option value="other">Khác</option>
              </select>
            </div>
          </div>

          <!-- Impacts -->
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-400 mb-1">Tác động thu nhập (VND/tháng)</label>
              <CurrencyInput v-model="draft.income_impact" class="w-full bg-gray-800 border border-gray-700 rounded-lg p-2.5 text-white focus:border-green-500 outline-none" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-400 mb-1">Tác động chi phí (VND/tháng)</label>
              <CurrencyInput v-model="draft.expense_impact" class="w-full bg-gray-800 border border-gray-700 rounded-lg p-2.5 text-white focus:border-red-500 outline-none" />
            </div>
          </div>
        </div>

        <div class="mt-8 flex justify-end space-x-4">
          <button @click="draft = null" class="px-6 py-2.5 rounded-lg text-gray-300 hover:text-white hover:bg-gray-800 transition-colors">
            Hủy bỏ
          </button>
          <button 
            @click="confirm"
            :disabled="timelineStore.loading"
            class="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white px-8 py-2.5 rounded-lg transition-colors shadow-[0_0_15px_rgba(37,99,235,0.4)]"
          >
            <span v-if="timelineStore.loading" class="animate-spin -ml-1 mr-2 h-4 w-4 border-b-2 border-white rounded-full inline-block"></span>
            Xác nhận & Lưu
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useTimelineStore } from '../stores/timelineStore'
import CurrencyInput from '../components/common/CurrencyInput.vue'

const router = useRouter()
const timelineStore = useTimelineStore()

const eventText = ref('')
const draft = ref(null)

const analyze = async () => {
  try {
    const result = await timelineStore.analyzeEvent(eventText.value)
    draft.value = result
  } catch (err) {
    alert(err)
  }
}

const confirm = async () => {
  try {
    await timelineStore.confirmEvent(draft.value)
    router.push('/profile/timeline')
  } catch (err) {
    alert(err)
  }
}
</script>
