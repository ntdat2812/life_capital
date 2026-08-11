<template>
  <div v-if="show" class="fixed inset-0 bg-slate-900/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
    <div class="bg-slate-800 rounded-2xl w-full max-w-xl border border-slate-700 overflow-hidden flex flex-col max-h-[90vh]">
      <div class="p-6 border-b border-slate-700 flex justify-between items-center bg-slate-800/50">
        <h2 class="text-xl font-bold text-white">Quản lý cảnh báo - {{ asset?.ticker || asset?.name }}</h2>
        <button @click="closeModal" class="text-slate-400 hover:text-white transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>
      
      <div class="p-6 overflow-y-auto custom-scrollbar flex-1">
        <!-- List existing alerts -->
        <div v-if="alerts.length > 0" class="mb-6 space-y-3">
          <h3 class="text-sm font-semibold text-slate-400 uppercase tracking-wider">Cảnh báo hiện tại</h3>
          <div v-for="alert in alerts" :key="alert.id" class="p-4 rounded-xl bg-slate-900/50 border border-slate-700 flex justify-between items-center group">
            <div>
              <div class="flex items-center gap-2 mb-1">
                <span class="px-2 py-0.5 rounded text-xs font-medium bg-indigo-500/20 text-indigo-400">
                  {{ formatAlertType(alert.alert_type) }}
                </span>
                <span v-if="alert.is_triggered" class="px-2 py-0.5 rounded text-xs font-medium bg-emerald-500/20 text-emerald-400">Đã kích hoạt</span>
                <span v-else-if="!alert.is_active" class="px-2 py-0.5 rounded text-xs font-medium bg-slate-500/20 text-slate-400">Đã tắt</span>
              </div>
              <div class="flex items-center gap-2">
                <p class="text-white font-medium">Mục tiêu: {{ formatCurrency(alert.target_value) }}</p>
                <span v-if="props.asset?.current_price" 
                      :class="['text-xs font-medium px-2 py-0.5 rounded-full', alert.target_value >= props.asset.current_price ? 'bg-emerald-500/10 text-emerald-400' : 'bg-red-500/10 text-red-400']"
                      title="So với giá hiện tại">
                  {{ alert.target_value > props.asset.current_price ? '+' : '' }}{{ (((alert.target_value - props.asset.current_price) / props.asset.current_price) * 100).toFixed(1) }}% (so với giá HT)
                </span>
              </div>
              <p v-if="alert.notes" class="text-xs text-slate-400 mt-1">{{ alert.notes }}</p>
            </div>
            <button @click="confirmDelete(alert.id)" class="text-slate-500 hover:text-red-400 transition-colors opacity-0 group-hover:opacity-100">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
            </button>
          </div>
        </div>

        <!-- Add new alert form -->
        <div>
          <h3 class="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-3">Thêm cảnh báo mới</h3>
          <div class="space-y-4 bg-slate-900/30 p-4 rounded-xl border border-slate-700/50">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Loại cảnh báo <span class="text-red-400">*</span></label>
              <select v-model="form.alert_type" class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2.5 text-white focus:outline-none focus:border-indigo-500">
                <option value="take_profit">Chốt lời (Take Profit)</option>
                <option value="stop_loss">Cắt lỗ (Stop Loss)</option>
                <option value="stop_buying">Ngừng mua (Stop Buying)</option>
                <option value="custom">Khác</option>
              </select>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Giá mục tiêu (VNĐ) <span class="text-red-400">*</span></label>
              <input type="text" v-model="displayTargetValue" @input="formatInputPrice" class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2.5 text-white focus:outline-none focus:border-indigo-500" placeholder="VD: 50,000">
            </div>

            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Ghi chú (Tùy chọn)</label>
              <input type="text" v-model="form.notes" class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2.5 text-white focus:outline-none focus:border-indigo-500" placeholder="Lý do đặt mục tiêu này">
            </div>

            <button @click="createAlert" :disabled="isSubmitting || form.target_value <= 0" class="w-full py-2.5 bg-indigo-600/20 text-indigo-400 border border-indigo-600/50 hover:bg-indigo-600 hover:text-white rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
              Thêm cảnh báo
            </button>
          </div>
        </div>
      </div>
      
      <!-- Confirm Delete Modal -->
      <ConfirmModal 
        :show="deleteConfirm.show" 
        title="Xóa cảnh báo" 
        message="Bạn có chắc chắn muốn xóa cảnh báo này không? Thao tác này không thể hoàn tác." 
        @confirm="executeDelete" 
        @cancel="deleteConfirm.show = false" 
      />
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import api from '../lib/api'
import ConfirmModal from './common/ConfirmModal.vue'

const props = defineProps({
  show: Boolean,
  asset: Object
})

const emit = defineEmits(['close'])

const alerts = ref([])
const isSubmitting = ref(false)

const form = ref({
  alert_type: 'take_profit',
  target_value: 0,
  notes: ''
})
const displayTargetValue = ref('')

watch(() => props.show, (val) => {
  if (val && props.asset) {
    fetchAlerts()
    resetForm()
  }
})

const formatCurrency = (value) => {
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(value || 0)
}

const formatInputPrice = (e) => {
  let val = e.target.value.replace(/\D/g, '')
  if (!val) {
    form.value.target_value = 0
    displayTargetValue.value = ''
    return
  }
  form.value.target_value = parseInt(val, 10)
  displayTargetValue.value = new Intl.NumberFormat('vi-VN').format(form.value.target_value)
}

const formatAlertType = (type) => {
  const map = {
    'take_profit': 'Chốt lời',
    'stop_loss': 'Cắt lỗ',
    'stop_buying': 'Ngừng mua',
    'custom': 'Khác'
  }
  return map[type] || type
}

const resetForm = () => {
  form.value = {
    alert_type: 'take_profit',
    target_value: 0,
    notes: ''
  }
  displayTargetValue.value = ''
}

const fetchAlerts = async () => {
  try {
    const res = await api.get(`/wealth/assets/${props.asset.id}/alerts`)
    alerts.value = res.data || []
  } catch (error) {
    console.error('Failed to fetch alerts:', error)
  }
}

const createAlert = async () => {
  if (form.value.target_value <= 0) return
  isSubmitting.value = true

  try {
    await api.post(`/wealth/assets/${props.asset.id}/alerts`, {
      alert_type: form.value.alert_type,
      target_value: form.value.target_value,
      notes: form.value.notes
    })

    await fetchAlerts()
    resetForm()
  } catch (error) {
    alert(error.response?.data?.message || error.message || 'Failed to create alert')
  } finally {
    isSubmitting.value = false
  }
}

const deleteConfirm = ref({
  show: false,
  id: null
})

const confirmDelete = (id) => {
  deleteConfirm.value = {
    show: true,
    id: id
  }
}

const executeDelete = async () => {
  const id = deleteConfirm.value.id
  if (!id) return

  try {
    await api.delete(`/wealth/alerts/${id}`)
    await fetchAlerts()
  } catch (error) {
    alert(error.response?.data?.message || error.message || 'Failed to delete alert')
  } finally {
    deleteConfirm.value.show = false
  }
}

const closeModal = () => {
  emit('close')
}
</script>
