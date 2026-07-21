<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-gray-100">Dòng thời gian sự kiện</h1>
      <button 
        @click="$router.push('/profile/timeline/new')"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors shadow-[0_0_15px_rgba(37,99,235,0.3)]"
      >
        <span class="mr-2">✨</span> Thêm sự kiện mới
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="timelineStore.loading" class="glass-panel p-8 text-center">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-4"></div>
      <p class="text-gray-400">Đang tải dữ liệu...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="timelineStore.events.length === 0" class="glass-panel p-12 text-center">
      <div class="text-4xl mb-4">🌱</div>
      <h3 class="text-xl font-semibold text-gray-200 mb-2">Chưa có sự kiện nào</h3>
      <p class="text-gray-400 max-w-md mx-auto mb-6">
        Hãy bắt đầu xây dựng hồ sơ tài chính của bạn bằng cách ghi lại các sự kiện quan trọng trong cuộc sống. AI sẽ tự động phân tích và điều chỉnh kế hoạch tài chính cho bạn.
      </p>
      <button 
        @click="$router.push('/profile/timeline/new')"
        class="bg-blue-600/20 text-blue-400 hover:bg-blue-600/30 px-6 py-2 rounded-lg transition-colors"
      >
        Thêm sự kiện đầu tiên
      </button>
    </div>

    <!-- Timeline -->
    <div v-else class="relative border-l border-gray-700 ml-4 space-y-8 pb-8">
      <div 
        v-for="(event, index) in timelineStore.events" 
        :key="event.id"
        class="relative pl-8"
      >
        <!-- Node point -->
        <div class="absolute -left-2.5 top-1.5 w-5 h-5 rounded-full bg-gray-800 border-2 border-blue-500 flex items-center justify-center">
          <div class="w-1.5 h-1.5 rounded-full bg-blue-500 shadow-[0_0_8px_rgba(59,130,246,0.8)]"></div>
        </div>
        
        <div class="glass-panel p-5 transition-transform hover:-translate-y-1 hover:shadow-lg">
          <div class="flex justify-between items-start mb-2">
            <h3 class="text-lg font-bold text-blue-400">{{ event.title }}</h3>
            <span class="text-xs text-gray-400 bg-gray-800 px-2 py-1 rounded-md">
              {{ formatDate(event.event_date) }}
            </span>
          </div>
          
          <div class="flex flex-wrap gap-2 mb-4">
            <span class="text-xs bg-blue-900/40 text-blue-300 px-2 py-1 rounded border border-blue-800/50">
              {{ formatCategory(event.category) }}
            </span>
            <span v-if="event.triggered_profile_version" class="text-xs bg-purple-900/40 text-purple-300 px-2 py-1 rounded border border-purple-800/50">
              Profile v{{ event.triggered_profile_version }}
            </span>
          </div>
          
          <div v-if="event.ai_impact_analysis" class="bg-gray-900/50 border border-gray-700/50 p-3 rounded-lg mb-3">
            <p class="text-sm text-gray-300"><span class="mr-1">🤖</span> {{ event.ai_impact_analysis }}</p>
          </div>
          
          <div class="grid grid-cols-2 gap-4 mt-2">
            <div v-if="event.income_impact !== 0">
              <span class="text-xs text-gray-500">Tác động thu nhập:</span>
              <p :class="event.income_impact > 0 ? 'text-green-400' : 'text-red-400'" class="font-medium">
                {{ event.income_impact > 0 ? '+' : '' }}{{ formatCurrency(event.income_impact) }}
              </p>
            </div>
            <div v-if="event.expense_impact !== 0">
              <span class="text-xs text-gray-500">Tác động chi phí:</span>
              <p :class="event.expense_impact > 0 ? 'text-red-400' : 'text-green-400'" class="font-medium">
                {{ event.expense_impact > 0 ? '+' : '' }}{{ formatCurrency(event.expense_impact) }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useTimelineStore } from '../stores/timelineStore'

const timelineStore = useTimelineStore()

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('vi-VN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  }).format(date)
}

const formatCurrency = (value) => {
  if (value === null || value === undefined) return '0 ₫'
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(value)
}

const formatCategory = (cat) => {
  const map = {
    'income_change': 'Thay đổi thu nhập',
    'family_change': 'Gia đình',
    'housing': 'Nhà ở',
    'dependent_change': 'Người phụ thuộc',
    'inheritance': 'Thừa kế',
    'health': 'Sức khỏe',
    'education': 'Giáo dục',
    'career': 'Sự nghiệp',
    'windfall': 'Lộc bất ngờ',
    'major_expense': 'Chi tiêu lớn',
    'other': 'Khác'
  }
  return map[cat] || cat
}

onMounted(() => {
  timelineStore.fetchTimeline()
})
</script>
