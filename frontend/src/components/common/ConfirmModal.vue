<template>
  <div v-if="show" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
    <div class="premium-card w-full max-w-sm p-6 transform transition-all shadow-2xl relative">
      <div class="flex items-center space-x-4 mb-4">
        <div class="w-12 h-12 rounded-full bg-red-500/20 flex items-center justify-center shrink-0">
          <svg class="w-6 h-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
          </svg>
        </div>
        <div>
          <h3 class="text-lg font-bold text-white">{{ title }}</h3>
          <p class="text-sm text-slate-400 mt-1 leading-relaxed">{{ message }}</p>
        </div>
      </div>
      
      <div class="flex space-x-3 mt-6">
        <button @click="$emit('cancel')" class="flex-1 py-2 px-4 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors border border-white/5 text-sm font-medium">
          Hủy
        </button>
        <button @click="$emit('confirm')" class="flex-1 py-2 px-4 bg-red-500 hover:bg-red-600 text-white rounded-lg transition-colors shadow-lg shadow-red-500/20 text-sm font-medium">
          Xác nhận
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  show: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: 'Xác nhận xóa'
  },
  message: {
    type: String,
    default: 'Bạn có chắc chắn muốn xóa mục này không? Hành động này không thể hoàn tác.'
  }
})

defineEmits(['confirm', 'cancel'])
</script>
