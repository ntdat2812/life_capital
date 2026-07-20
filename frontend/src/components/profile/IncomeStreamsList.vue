<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h3 class="text-xl font-semibold text-white">Nguồn Thu Nhập</h3>
      <button @click="openModal()" class="px-4 py-2 text-sm bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition">
        + Thêm Thu Nhập
      </button>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div class="bg-slate-800/50 p-4 rounded-xl border border-white/5">
        <p class="text-xs text-slate-400 mb-1">Tổng Thu Nhập / tháng</p>
        <p class="text-lg font-bold text-emerald-400">{{ formatCurrency(totalIncome) }}</p>
      </div>
      <div class="bg-slate-800/50 p-4 rounded-xl border border-white/5">
        <p class="text-xs text-slate-400 mb-1">Thu nhập Chủ động / tháng</p>
        <p class="text-lg font-bold text-indigo-400">{{ formatCurrency(activeIncome) }}</p>
      </div>
      <div class="bg-slate-800/50 p-4 rounded-xl border border-white/5">
        <p class="text-xs text-slate-400 mb-1">Thu nhập Thụ động / tháng</p>
        <p class="text-lg font-bold text-violet-400">{{ formatCurrency(passiveIncome) }}</p>
      </div>
      <div class="bg-slate-800/50 p-4 rounded-xl border border-white/5">
        <p class="text-xs text-slate-400 mb-1">Tỷ lệ Thụ động</p>
        <p class="text-lg font-bold text-white">{{ passiveRatio }}%</p>
      </div>
    </div>

    <!-- List -->
    <div v-if="profileStore.incomes.length === 0" class="text-center py-12 bg-slate-800/30 rounded-2xl border border-white/5 border-dashed">
      <p class="text-slate-400 text-sm">Chưa có nguồn thu nhập nào.</p>
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div v-for="income in profileStore.incomes" :key="income.id" 
           class="premium-card p-5 relative group">
        <div class="flex justify-between items-start mb-2">
          <div>
            <h4 class="text-lg font-semibold text-white">{{ income.name }}</h4>
            <span class="text-xs px-2 py-0.5 rounded-full" 
                  :class="income.is_passive ? 'bg-violet-500/20 text-violet-400' : 'bg-indigo-500/20 text-indigo-400'">
              {{ income.is_passive ? 'Thụ động' : 'Chủ động' }}
            </span>
            <span v-if="!income.is_active" class="text-xs px-2 py-0.5 rounded-full bg-slate-500/20 text-slate-400 ml-2">
              Đã dừng
            </span>
          </div>
          <div class="opacity-0 group-hover:opacity-100 transition-opacity flex space-x-2">
            <button @click="openModal(income)" class="text-slate-400 hover:text-white p-1">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
            </button>
            <button @click="handleDelete(income.id)" class="text-red-400 hover:text-red-300 p-1">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
            </button>
          </div>
        </div>
        <div class="text-2xl font-bold text-emerald-400 mb-1">{{ formatCurrency(income.amount) }} <span class="text-sm font-normal text-slate-400">/ {{ getFrequencyLabel(income.frequency) }}</span></div>
        <p v-if="income.notes" class="text-sm text-slate-400 truncate">{{ income.notes }}</p>
      </div>
    </div>

    <!-- Modal Form -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="premium-card w-full max-w-md p-6">
        <h2 class="text-xl font-bold text-white mb-4">{{ form.id ? 'Sửa' : 'Thêm' }} Nguồn Thu Nhập</h2>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Tên nguồn thu (VD: Lương Techcombank)</label>
            <input v-model="form.name" required class="w-full bg-slate-900/50 border border-white/10 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Loại</label>
              <select v-model="form.type" class="w-full bg-slate-900/50 border border-white/10 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500">
                <option value="salary">Tiền lương</option>
                <option value="business">Kinh doanh</option>
                <option value="rental">Cho thuê nhà</option>
                <option value="dividend">Cổ tức</option>
                <option value="interest">Lãi tiết kiệm</option>
                <option value="other">Khác</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Tần suất</label>
              <select v-model="form.frequency" class="w-full bg-slate-900/50 border border-white/10 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500">
                <option value="monthly">Hàng tháng</option>
                <option value="yearly">Hàng năm</option>
                <option value="one_time">Một lần</option>
              </select>
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Số tiền</label>
            <CurrencyInput v-model="form.amount" required min="1" />
          </div>
          <div class="flex items-center space-x-2 mt-2">
            <input type="checkbox" id="isPassive" v-model="form.is_passive" class="w-4 h-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500">
            <label for="isPassive" class="text-sm text-slate-300">Đây là nguồn thu nhập thụ động</label>
          </div>
          <div v-if="form.id" class="flex items-center space-x-2 mt-2">
            <input type="checkbox" id="isActive" v-model="form.is_active" class="w-4 h-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500">
            <label for="isActive" class="text-sm text-slate-300">Đang hoạt động</label>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Ghi chú</label>
            <input v-model="form.notes" class="w-full bg-slate-900/50 border border-white/10 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
          </div>
          <div v-if="error" class="text-red-400 text-sm">{{ error }}</div>
          <div class="pt-4 flex justify-end space-x-3">
            <button type="button" @click="showModal = false" class="px-4 py-2 text-sm text-slate-300 hover:text-white transition">Hủy</button>
            <button type="submit" :disabled="loading" class="px-4 py-2 text-sm bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition disabled:opacity-50">
              Lưu
            </button>
          </div>
        </form>
      </div>
    </div>
    <!-- Confirm Delete Modal -->
    <ConfirmModal 
      :show="deleteConfirm.show" 
      :title="'Xóa Nguồn Thu Nhập'" 
      :message="'Bạn có chắc muốn xóa nguồn thu nhập này không? Thao tác này không thể hoàn tác.'" 
      @confirm="executeDelete" 
      @cancel="deleteConfirm.show = false" 
    />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useProfileStore } from '../../stores/profileStore'
import CurrencyInput from '../common/CurrencyInput.vue'
import ConfirmModal from '../common/ConfirmModal.vue'

const profileStore = useProfileStore()

const formatCurrency = (val) => new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(val || 0)

const getFrequencyLabel = (freq) => {
  const map = { 'monthly': 'tháng', 'yearly': 'năm', 'one_time': 'lần' }
  return map[freq] || freq
}

const totalIncome = computed(() => profileStore.incomes.filter(i => i.is_active).reduce((sum, i) => sum + (i.frequency === 'yearly' ? i.amount / 12 : i.amount), 0))
const passiveIncome = computed(() => profileStore.incomes.filter(i => i.is_active && i.is_passive).reduce((sum, i) => sum + (i.frequency === 'yearly' ? i.amount / 12 : i.amount), 0))
const activeIncome = computed(() => totalIncome.value - passiveIncome.value)
const passiveRatio = computed(() => totalIncome.value > 0 ? ((passiveIncome.value / totalIncome.value) * 100).toFixed(1) : 0)

const showModal = ref(false)
const loading = ref(false)
const error = ref('')
const form = ref({
  id: null,
  name: '',
  type: 'salary',
  frequency: 'monthly',
  amount: 0,
  is_passive: false,
  is_active: true,
  notes: ''
})

const openModal = (income = null) => {
  if (income) {
    form.value = { ...income }
  } else {
    form.value = { id: null, name: '', type: 'salary', frequency: 'monthly', amount: 0, is_passive: false, is_active: true, notes: '' }
  }
  error.value = ''
  showModal.value = true
}

const handleSubmit = async () => {
  loading.value = true
  error.value = ''
  try {
    if (form.value.id) {
      await profileStore.updateIncome(form.value.id, form.value)
    } else {
      await profileStore.createIncome(form.value)
    }
    showModal.value = false
  } catch (err) {
    error.value = err
  } finally {
    loading.value = false
  }
}

const deleteConfirm = ref({ show: false, id: null })

const handleDelete = (id) => {
  deleteConfirm.value = { show: true, id }
}

const executeDelete = async () => {
  if (deleteConfirm.value.id) {
    await profileStore.deleteIncome(deleteConfirm.value.id)
    deleteConfirm.value.show = false
  }
}
</script>
