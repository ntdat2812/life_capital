<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h3 class="text-xl font-semibold text-white">Người Phụ Thuộc</h3>
      <button @click="openModal()" class="px-4 py-2 text-sm bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition">
        + Thêm Người Phụ Thuộc
      </button>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 gap-4">
      <div class="bg-slate-800/50 p-4 rounded-xl border border-white/5">
        <p class="text-xs text-slate-400 mb-1">Số lượng người phụ thuộc</p>
        <p class="text-lg font-bold text-white">{{ activeDependentsCount }}</p>
      </div>
      <div class="bg-slate-800/50 p-4 rounded-xl border border-white/5">
        <p class="text-xs text-slate-400 mb-1">Tổng chi phí nuôi dưỡng / tháng</p>
        <p class="text-lg font-bold text-amber-400">{{ formatCurrency(totalCost) }}</p>
      </div>
    </div>

    <!-- List -->
    <div v-if="profileStore.dependents.length === 0" class="text-center py-12 bg-slate-800/30 rounded-2xl border border-white/5 border-dashed">
      <p class="text-slate-400 text-sm">Chưa có người phụ thuộc nào.</p>
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div v-for="dep in profileStore.dependents" :key="dep.id" 
           class="premium-card p-5 relative group">
        <div class="flex justify-between items-start mb-2">
          <div>
            <h4 class="text-lg font-semibold text-white">{{ dep.name }}</h4>
            <span class="text-xs px-2 py-0.5 rounded-full bg-slate-500/20 text-slate-300">
              {{ dep.relationship }}
            </span>
            <span v-if="!dep.is_active" class="text-xs px-2 py-0.5 rounded-full bg-red-500/20 text-red-400 ml-2">
              Đã dừng
            </span>
          </div>
          <div class="opacity-0 group-hover:opacity-100 transition-opacity flex space-x-2">
            <button @click="openModal(dep)" class="text-slate-400 hover:text-white p-1">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
            </button>
            <button @click="handleDelete(dep.id)" class="text-red-400 hover:text-red-300 p-1">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
            </button>
          </div>
        </div>
        <div v-if="dep.date_of_birth" class="text-sm text-slate-400 mb-2">Tuổi: {{ calculateAge(dep.date_of_birth) }}</div>
        <div class="text-xl font-bold text-amber-400 mb-1">{{ formatCurrency(dep.monthly_cost) }} <span class="text-sm font-normal text-slate-400">/ tháng</span></div>
        <p v-if="dep.notes" class="text-sm text-slate-400 truncate">{{ dep.notes }}</p>
      </div>
    </div>

    <!-- Modal Form -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="premium-card w-full max-w-md p-6">
        <h2 class="text-xl font-bold text-white mb-4">{{ form.id ? 'Sửa' : 'Thêm' }} Người Phụ Thuộc</h2>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Họ tên</label>
            <input v-model="form.name" required class="w-full bg-slate-900/50 border border-white/10 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Mối quan hệ</label>
            <select v-model="form.relationship" class="w-full bg-slate-900/50 border border-white/10 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500">
              <option value="Con cái">Con cái</option>
              <option value="Bố mẹ">Bố mẹ</option>
              <option value="Vợ/Chồng">Vợ/Chồng</option>
              <option value="Khác">Khác</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Ngày sinh (tùy chọn)</label>
            <DateInput v-model="form.date_of_birth" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Chi phí nuôi dưỡng / tháng (VND)</label>
            <CurrencyInput v-model="form.monthly_cost" min="0" />
          </div>
          <div v-if="form.id" class="flex items-center space-x-2 mt-2">
            <input type="checkbox" id="isActiveDep" v-model="form.is_active" class="w-4 h-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500">
            <label for="isActiveDep" class="text-sm text-slate-300">Vẫn đang phụ thuộc</label>
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
      :title="'Xóa Người Phụ Thuộc'" 
      :message="'Bạn có chắc muốn xóa người phụ thuộc này không? Thao tác này không thể hoàn tác.'" 
      @confirm="executeDelete" 
      @cancel="deleteConfirm.show = false" 
    />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useProfileStore } from '../../stores/profileStore'
import CurrencyInput from '../common/CurrencyInput.vue'
import DateInput from '../common/DateInput.vue'
import ConfirmModal from '../common/ConfirmModal.vue'

const profileStore = useProfileStore()

const formatCurrency = (val) => new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(val || 0)

const calculateAge = (dateString) => {
  if (!dateString) return ''
  const today = new Date()
  const birthDate = new Date(dateString)
  let age = today.getFullYear() - birthDate.getFullYear()
  const m = today.getMonth() - birthDate.getMonth()
  if (m < 0 || (m === 0 && today.getDate() < birthDate.getDate())) {
    age--
  }
  return age
}

const activeDependentsCount = computed(() => profileStore.dependents.filter(d => d.is_active).length)
const totalCost = computed(() => profileStore.dependents.filter(d => d.is_active).reduce((sum, d) => sum + d.monthly_cost, 0))

const showModal = ref(false)
const loading = ref(false)
const error = ref('')
const form = ref({
  id: null,
  name: '',
  relationship: 'Con cái',
  date_of_birth: '',
  monthly_cost: 0,
  is_active: true,
  notes: ''
})

const openModal = (dep = null) => {
  if (dep) {
    let dob = ''
    if (dep.date_of_birth) {
      const dObj = new Date(dep.date_of_birth)
      dob = `${dObj.getFullYear()}-${String(dObj.getMonth() + 1).padStart(2, '0')}-${String(dObj.getDate()).padStart(2, '0')}`
    }
    form.value = { ...dep, date_of_birth: dob }
  } else {
    form.value = { id: null, name: '', relationship: 'Con cái', date_of_birth: '', monthly_cost: 0, is_active: true, notes: '' }
  }
  error.value = ''
  showModal.value = true
}

const handleSubmit = async () => {
  loading.value = true
  error.value = ''
  try {
    const payload = { ...form.value }
    if (payload.date_of_birth) {
      payload.date_of_birth = new Date(payload.date_of_birth).toISOString()
    } else {
      payload.date_of_birth = null
    }

    if (form.value.id) {
      await profileStore.updateDependent(form.value.id, payload)
    } else {
      await profileStore.createDependent(payload)
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
    await profileStore.deleteDependent(deleteConfirm.value.id)
    deleteConfirm.value.show = false
  }
}
</script>
