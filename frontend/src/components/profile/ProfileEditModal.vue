<template>
  <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
    <div class="premium-card w-full max-w-lg p-6 max-h-[90vh] overflow-y-auto">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-bold text-white font-outfit">Chỉnh sửa Hồ sơ</h2>
        <button @click="$emit('close')" class="text-slate-400 hover:text-white">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-1">Ngày sinh</label>
          <DateInput v-model="form.date_of_birth" />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-300 mb-1">Tình trạng hôn nhân</label>
          <select v-model="form.marital_status" class="w-full bg-slate-900/50 border border-white/10 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500">
            <option value="Độc thân">Độc thân</option>
            <option value="Đã kết hôn">Đã kết hôn</option>
            <option value="Ly hôn">Ly hôn</option>
          </select>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Khẩu vị rủi ro</label>
            <select v-model="form.risk_tolerance" class="w-full bg-slate-900/50 border border-white/10 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500">
              <option value="Thấp">Thấp</option>
              <option value="Trung bình">Trung bình</option>
              <option value="Cao">Cao</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Điểm rủi ro (1-100)</label>
            <input v-model.number="form.risk_score" type="number" min="1" max="100" class="w-full bg-slate-900/50 border border-white/10 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Chi phí thiết yếu / tháng</label>
            <CurrencyInput v-model="form.essential_monthly_expense" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Chi phí hưởng thụ / tháng</label>
            <CurrencyInput v-model="form.discretionary_monthly_expense" />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-300 mb-1">Mục tiêu Tự do tài chính (VND)</label>
          <CurrencyInput v-model="form.fi_target_amount" />
        </div>

        <div v-if="error" class="text-red-400 text-sm mt-2">{{ error }}</div>

        <div class="pt-4 flex justify-end space-x-3">
          <button type="button" @click="$emit('close')" class="px-4 py-2 text-sm text-slate-300 hover:text-white transition">Hủy</button>
          <button type="submit" :disabled="loading" class="px-4 py-2 text-sm bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition disabled:opacity-50">
            {{ loading ? 'Đang lưu...' : 'Lưu thay đổi' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useProfileStore } from '../../stores/profileStore'
import CurrencyInput from '../common/CurrencyInput.vue'
import DateInput from '../common/DateInput.vue'

const props = defineProps({
  show: Boolean,
  profile: Object
})

const emit = defineEmits(['close', 'saved'])
const profileStore = useProfileStore()

const loading = ref(false)
const error = ref('')
const form = ref({
  date_of_birth: '',
  marital_status: '',
  risk_tolerance: '',
  risk_score: 50,
  essential_monthly_expense: 0,
  discretionary_monthly_expense: 0,
  fi_target_amount: 0
})

watch(() => props.show, (newVal) => {
  if (newVal && props.profile) {
    // Format date for input type="date"
    let dob = ''
    if (props.profile.date_of_birth) {
      // Dữ liệu từ API về thường là dạng Date string chuẩn YYYY-MM-DD
      const dateObj = new Date(props.profile.date_of_birth)
      dob = `${dateObj.getFullYear()}-${String(dateObj.getMonth() + 1).padStart(2, '0')}-${String(dateObj.getDate()).padStart(2, '0')}`
    }
    
    form.value = {
      date_of_birth: dob,
      marital_status: props.profile.marital_status || 'Độc thân',
      risk_tolerance: props.profile.risk_tolerance || 'Trung bình',
      risk_score: props.profile.risk_score || 50,
      essential_monthly_expense: props.profile.essential_monthly_expense || 0,
      discretionary_monthly_expense: props.profile.discretionary_monthly_expense || 0,
      fi_target_amount: props.profile.fi_target_amount || 0
    }
  }
})

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
    
    await profileStore.updateProfile(payload)
    emit('saved')
    emit('close')
  } catch (err) {
    error.value = err
  } finally {
    loading.value = false
  }
}
</script>
