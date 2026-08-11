<template>
  <div v-if="show" class="fixed inset-0 bg-slate-900/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
    <div class="bg-slate-800 rounded-2xl w-full max-w-2xl border border-slate-700 overflow-hidden flex flex-col max-h-[90vh]">
      <div class="p-6 border-b border-slate-700 flex justify-between items-center bg-slate-800/50">
        <h2 class="text-xl font-bold text-white">Cập nhật giá Vàng hàng loạt</h2>
        <button @click="closeModal" class="text-slate-400 hover:text-white transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>
      
      <div class="p-6 overflow-y-auto custom-scrollbar flex-1">
        <p class="text-sm text-slate-400 mb-4">Chọn các tài sản vàng muốn cập nhật và nhập giá mới (tính theo Lượng).</p>
        
        <div class="mb-4">
          <label class="flex items-center text-slate-300 font-medium cursor-pointer">
            <input type="checkbox" v-model="selectAll" @change="toggleSelectAll" class="mr-2 rounded border-slate-600 bg-slate-900/50 text-indigo-500 focus:ring-indigo-500 focus:ring-offset-slate-800">
            Chọn tất cả
          </label>
        </div>

        <div class="space-y-2 mb-6">
          <div v-for="asset in goldAssets" :key="asset.id" class="flex items-center p-3 rounded-lg bg-slate-900/50 border border-slate-700/50 hover:bg-slate-800/80 transition-colors">
            <input type="checkbox" :id="asset.id" :value="asset.id" v-model="selectedAssetIds" class="mr-3 rounded border-slate-600 bg-slate-900/50 text-indigo-500 focus:ring-indigo-500 focus:ring-offset-slate-800">
            <label :for="asset.id" class="flex-1 cursor-pointer flex justify-between items-center">
              <div>
                <span class="text-slate-200 font-medium">{{ asset.name }}</span>
                <span class="text-xs text-slate-500 ml-2">SL: {{ asset.quantity }}</span>
              </div>
              <span class="text-sm text-slate-400">Giá hiện tại: {{ formatCurrency(asset.current_price) }}</span>
            </label>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-300 mb-1">Giá mới (VNĐ/Lượng) <span class="text-red-400">*</span></label>
          <input type="text" v-model="displayNewPrice" @input="formatInputPrice" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-indigo-500" placeholder="VD: 78,000,000">
        </div>
      </div>
      
      <div class="p-6 border-t border-slate-700 bg-slate-800/50 flex justify-end gap-3">
        <button type="button" @click="closeModal" class="px-4 py-2 text-slate-300 hover:text-white transition-colors">Hủy</button>
        <button type="button" @click="submit" :disabled="isSubmitting || selectedAssetIds.length === 0 || !newPrice" class="px-6 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2">
          <span v-if="isSubmitting" class="animate-spin inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full"></span>
          Cập nhật
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import api from '../lib/api'
import { getGoldMultiplier, extractGoldUnit } from '../utils/goldCalculations'

const props = defineProps({
  show: Boolean,
  goldAssets: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['close', 'updated'])

const selectedAssetIds = ref([])
const newPrice = ref(0)
const displayNewPrice = ref('')
const isSubmitting = ref(false)

const selectAll = computed({
  get: () => goldAssets.value && selectedAssetIds.value.length === goldAssets.value.length && goldAssets.value.length > 0,
  set: (val) => {
    if (val) {
      selectedAssetIds.value = props.goldAssets.map(a => a.id)
    } else {
      selectedAssetIds.value = []
    }
  }
})

const goldAssets = computed(() => props.goldAssets)

watch(() => props.show, (val) => {
  if (val) {
    selectedAssetIds.value = []
    newPrice.value = 0
    displayNewPrice.value = ''
    isSubmitting.value = false
  }
})

const toggleSelectAll = (e) => {
  if (e.target.checked) {
    selectedAssetIds.value = props.goldAssets.map(a => a.id)
  } else {
    selectedAssetIds.value = []
  }
}

const formatCurrency = (value) => {
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(value || 0)
}

const formatInputPrice = (e) => {
  let val = e.target.value.replace(/\D/g, '')
  if (!val) {
    newPrice.value = 0
    displayNewPrice.value = ''
    return
  }
  newPrice.value = parseInt(val, 10)
  displayNewPrice.value = new Intl.NumberFormat('vi-VN').format(newPrice.value)
}

const closeModal = () => {
  emit('close')
}

const submit = async () => {
  if (selectedAssetIds.value.length === 0 || newPrice.value <= 0) return

  isSubmitting.value = true
  try {
    const updates = []
    
    // Calculate the per-unit price and total value for each selected asset
    selectedAssetIds.value.forEach(id => {
      const asset = props.goldAssets.find(a => a.id === id)
      if (asset) {
        const unit = extractGoldUnit(asset.name)
        const multiplier = getGoldMultiplier(unit)
        const pricePerUnit = newPrice.value * multiplier
        const currentValue = Math.round((asset.quantity || 0) * pricePerUnit)
        
        updates.push({
          asset_id: asset.id,
          current_price: pricePerUnit,
          current_value: currentValue
        })
      }
    })

    const res = await api.patch('/wealth/assets/gold/bulk-update-price', {
      updates: updates
    })

    emit('updated')
    closeModal()
  } catch (error) {
    alert(error.response?.data?.message || error.message || 'Có lỗi xảy ra')
  } finally {
    isSubmitting.value = false
  }
}
</script>
