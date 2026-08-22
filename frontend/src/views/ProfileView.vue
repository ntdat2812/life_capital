<template>
  <div class="min-h-screen p-4 md:p-8">
    <div class="max-w-5xl mx-auto">
      <div class="flex flex-col md:flex-row md:items-center justify-between mb-8 gap-4">
        <h1 class="text-3xl font-bold text-white">Hồ Sơ Của Tôi</h1>
        <div class="flex items-center space-x-4">
          <button @click="showEditModal = true" class="px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white text-sm rounded-lg transition border border-white/10">
            Chỉnh sửa Hồ Sơ
          </button>
          <button @click="reTakeOnboarding" class="text-indigo-400 hover:text-indigo-300 text-sm flex items-center transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Làm Lại Khảo Sát
          </button>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
      
      <div v-else-if="error" class="premium-card p-8 text-center text-red-400">
        {{ error }}
      </div>

      <div v-else-if="profile">
        <!-- Tabs -->
        <div class="flex space-x-1 bg-slate-800/50 p-1 rounded-xl mb-6 overflow-x-auto">
          <button @click="activeTab = 'overview'" 
                  :class="['px-6 py-2.5 rounded-lg text-sm font-medium transition-all whitespace-nowrap', 
                          activeTab === 'overview' ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400 hover:text-white hover:bg-white/5']">
            Tổng Quan
          </button>
          <button @click="activeTab = 'incomes'" 
                  :class="['px-6 py-2.5 rounded-lg text-sm font-medium transition-all whitespace-nowrap', 
                          activeTab === 'incomes' ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400 hover:text-white hover:bg-white/5']">
            Tiền Vào (Incomes)
          </button>
          <button @click="activeTab = 'dependents'" 
                  :class="['px-6 py-2.5 rounded-lg text-sm font-medium transition-all whitespace-nowrap', 
                          activeTab === 'dependents' ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400 hover:text-white hover:bg-white/5']">
            Người Phụ Thuộc
          </button>
          <button @click="activeTab = 'settings'" 
                  :class="['px-6 py-2.5 rounded-lg text-sm font-medium transition-all whitespace-nowrap', 
                          activeTab === 'settings' ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400 hover:text-white hover:bg-white/5']">
            Cài Đặt
          </button>
        </div>

        <!-- Tab Content: Overview -->
        <div v-if="activeTab === 'overview'" class="space-y-6">
          <!-- Hero Banner: Mục Tiêu Tự Do Tài Chính (FI) -->
          <div class="premium-card p-8 bg-gradient-to-br from-slate-800/80 to-indigo-900/20 border-indigo-500/20 text-center relative overflow-hidden">
            <div class="absolute inset-0 bg-[url('/grid.svg')] bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))] opacity-10"></div>
            <div class="relative z-10">
              <div class="flex items-center justify-center mb-4">
                <span class="text-3xl mr-3">🎯</span>
                <h3 class="text-xl font-semibold text-white tracking-wide uppercase">Mục Tiêu Tự Do Tài Chính (FI)</h3>
              </div>
              <div class="text-5xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 via-violet-400 to-emerald-400 mb-4 drop-shadow-sm">
                {{ formatCurrency(profile.fi_target_amount) }}
              </div>
              <p class="text-slate-400 text-base max-w-2xl mx-auto leading-relaxed">
                Dựa trên mức chi tiêu hiện tại, đây là con số tài sản bạn cần tích lũy để có thể nghỉ hưu an nhàn và đạt trạng thái hoàn toàn tự do tài chính.
              </p>
            </div>
          </div>

          <!-- Cột Trái & Phải -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Cột Trái: Thông tin cá nhân & Rủi ro -->
            <div class="premium-card p-6 h-full flex flex-col justify-center">
              <div class="flex items-center space-x-5 mb-8">
                <div class="w-20 h-20 rounded-full bg-gradient-to-tr from-indigo-500 to-violet-500 flex items-center justify-center text-3xl text-white font-bold shadow-lg shadow-indigo-500/20">
                  {{ userInitials }}
                </div>
                <div>
                  <h3 class="text-2xl font-bold text-white mb-1">{{ authStore.user?.name || 'Nhà đầu tư' }}</h3>
                  <p class="text-slate-400 text-sm">Cập nhật: {{ formatDate(profile.updated_at) }}</p>
                </div>
              </div>
              
              <div class="space-y-3">
                <p class="text-slate-400 text-sm uppercase tracking-wider font-medium">Khẩu vị rủi ro</p>
                <div class="flex items-baseline space-x-3">
                  <span class="text-3xl font-extrabold" :class="riskColor(profile.risk_score)">{{ profile.risk_score }}/100</span>
                  <span class="text-lg font-medium px-3 py-1 rounded-full bg-slate-800/50 border border-white/5" :class="riskColor(profile.risk_score)">{{ profile.risk_tolerance }}</span>
                </div>
                <!-- Thanh tiến trình rủi ro -->
                <div class="w-full bg-slate-800 rounded-full h-3 mt-3 shadow-inner">
                  <div class="h-3 rounded-full transition-all duration-1000 shadow-[0_0_10px_rgba(0,0,0,0.5)]" :class="riskBgColor(profile.risk_score)" :style="`width: ${profile.risk_score}%`"></div>
                </div>
              </div>
            </div>

            <!-- Cột Phải: Thông tin chung & Dòng tiền -->
            <div class="premium-card p-6 h-full flex flex-col">
              <h4 class="text-sm font-semibold text-white mb-5 uppercase tracking-wider flex items-center">
                <svg class="w-4 h-4 mr-2 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                Chỉ số tổng quan
              </h4>
              <ul class="space-y-4 flex-1">
                <li class="flex justify-between items-center border-b border-white/5 pb-3">
                  <span class="text-slate-400">Tuổi hiện tại</span>
                  <span class="text-white font-medium bg-slate-800/50 px-3 py-1 rounded-lg border border-white/5">{{ calculateAge(profile.date_of_birth) }}</span>
                </li>
                <li class="flex justify-between items-center border-b border-white/5 pb-3">
                  <span class="text-slate-400">Tình trạng hôn nhân</span>
                  <span class="text-white font-medium bg-slate-800/50 px-3 py-1 rounded-lg border border-white/5">{{ profile.marital_status || 'Chưa cập nhật' }}</span>
                </li>
                <li class="flex justify-between items-center border-b border-white/5 pb-3">
                  <span class="text-slate-400">Tổng thu nhập / tháng</span>
                  <span class="text-emerald-400 font-bold text-lg">{{ formatCurrency(profileStore.totalIncome) }}</span>
                </li>
                <li class="flex flex-col border-b border-white/5 pb-3">
                  <div class="flex justify-between items-center mb-2">
                    <span class="text-slate-400">Tổng chi phí / tháng</span>
                    <span class="text-amber-400 font-bold text-lg">{{ formatCurrency(profileStore.essentialExpense + profileStore.discretionaryExpense + profileStore.dependentsExpense) }}</span>
                  </div>
                  <div class="bg-slate-800/30 p-3 rounded-xl border border-white/5 space-y-2">
                    <div class="flex justify-between items-center">
                      <span class="text-slate-500 text-sm flex items-center"><span class="w-1.5 h-1.5 rounded-full bg-slate-400 mr-2"></span>Thiết yếu</span>
                      <span class="text-slate-300 text-sm font-medium">{{ formatCurrency(profileStore.essentialExpense) }}</span>
                    </div>
                    <div class="flex justify-between items-center">
                      <span class="text-slate-500 text-sm flex items-center"><span class="w-1.5 h-1.5 rounded-full bg-slate-400 mr-2"></span>Hưởng thụ</span>
                      <span class="text-slate-300 text-sm font-medium">{{ formatCurrency(profileStore.discretionaryExpense) }}</span>
                    </div>
                    <div class="flex justify-between items-center" v-if="profileStore.dependentsExpense > 0">
                      <span class="text-slate-500 text-sm flex items-center"><span class="w-1.5 h-1.5 rounded-full bg-slate-400 mr-2"></span>Người phụ thuộc</span>
                      <span class="text-slate-300 text-sm font-medium">{{ formatCurrency(profileStore.dependentsExpense) }}</span>
                    </div>
                  </div>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Tab Content: Incomes -->
        <div v-else-if="activeTab === 'incomes'">
          <IncomeStreamsList />
        </div>

        <!-- Tab Content: Dependents -->
        <div v-else-if="activeTab === 'dependents'">
          <DependentsList />
        </div>

        <!-- Tab Content: Settings -->
        <div v-else-if="activeTab === 'settings'" class="max-w-xl mx-auto space-y-6">
          <div class="premium-card p-8 rounded-2xl">
            <h3 class="text-xl font-bold text-white mb-6">Đổi mật khẩu</h3>
            <form @submit.prevent="handleChangePassword" class="space-y-4">
              <div v-if="passwordError" class="p-4 bg-red-500/10 border border-red-500/20 rounded-xl text-red-400 text-sm">
                {{ passwordError }}
              </div>
              <div v-if="passwordSuccess" class="p-4 bg-emerald-500/10 border border-emerald-500/20 rounded-xl text-emerald-400 text-sm">
                Đổi mật khẩu thành công!
              </div>

              <div>
                <label class="block text-sm font-medium text-slate-300 mb-1">Mật khẩu hiện tại</label>
                <input type="password" v-model="passwordForm.oldPassword" required
                  class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-indigo-500 transition-colors">
              </div>
              
              <div>
                <label class="block text-sm font-medium text-slate-300 mb-1">Mật khẩu mới</label>
                <input type="password" v-model="passwordForm.newPassword" required minlength="6"
                  class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-indigo-500 transition-colors">
                <p class="text-xs text-slate-500 mt-1">Ít nhất 6 ký tự</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-slate-300 mb-1">Xác nhận mật khẩu mới</label>
                <input type="password" v-model="passwordForm.confirmPassword" required
                  class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-indigo-500 transition-colors">
              </div>

              <div class="pt-4">
                <button type="submit" :disabled="isSubmittingPassword"
                  class="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-3 rounded-lg transition-colors flex justify-center items-center">
                  <svg v-if="isSubmittingPassword" class="animate-spin h-5 w-5 mr-2 text-white" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Đổi mật khẩu
                </button>
              </div>
            </form>
          </div>
        </div>

      </div>

      <div v-else class="premium-card p-12 flex flex-col items-center justify-center text-center">
        <div class="w-20 h-20 bg-indigo-500/20 rounded-full flex items-center justify-center mb-6">
          <span class="text-4xl">📝</span>
        </div>
        <h2 class="text-2xl font-bold text-white mb-2">Bạn chưa có Hồ Sơ Đầu Tư</h2>
        <p class="text-slate-400 mb-8 max-w-md">
          Hồ sơ đầu tư giúp AI của WealthOS hiểu rõ mức độ rủi ro và tình hình tài chính của bạn để đưa ra những lời khuyên phù hợp nhất.
        </p>
        <button @click="reTakeOnboarding" class="bg-indigo-600 hover:bg-indigo-700 text-white px-8 py-3 rounded-xl font-semibold shadow-lg shadow-indigo-500/30 transition-all hover:-translate-y-1">
          Bắt Đầu Khảo Sát AI
        </button>
      </div>

      <!-- Edit Modal -->
      <ProfileEditModal 
        :show="showEditModal" 
        :profile="profile" 
        @close="showEditModal = false" 
        @saved="handleProfileSaved" 
      />

      <!-- Confirm Retake Modal -->
      <ConfirmModal 
        :show="showRetakeConfirm"
        title="Làm lại khảo sát"
        message="Bạn có chắc chắn muốn làm lại khảo sát từ đầu? Các dữ liệu Hồ Sơ rủi ro hiện tại sẽ bị ghi đè."
        @confirm="proceedRetakeOnboarding"
        @cancel="showRetakeConfirm = false"
      />

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useProfileStore } from '../stores/profileStore';
import { useAuthStore } from '../stores/authStore';

import IncomeStreamsList from '../components/profile/IncomeStreamsList.vue';
import DependentsList from '../components/profile/DependentsList.vue';
import ProfileEditModal from '../components/profile/ProfileEditModal.vue';
import ConfirmModal from '../components/common/ConfirmModal.vue';
import api from '../lib/api';

const router = useRouter();
const profileStore = useProfileStore();
const authStore = useAuthStore();

const activeTab = ref('overview');
const showEditModal = ref(false);
const showRetakeConfirm = ref(false);

const profile = computed(() => profileStore.profile);
const loading = computed(() => profileStore.loading);
const error = computed(() => profileStore.error);

const userInitials = computed(() => {
  const name = authStore.user?.name || 'User';
  return name.charAt(0).toUpperCase();
});

onMounted(async () => {
  await profileStore.fetchProfile();
  await profileStore.fetchIncomes();
  await profileStore.fetchDependents();
});

const handleProfileSaved = () => {
  // refresh if needed, but the store updates the profile state directly
};

const reTakeOnboarding = () => {
  if (profileStore.profile) {
    showRetakeConfirm.value = true;
  } else {
    router.push('/onboarding/interview');
  }
};

const proceedRetakeOnboarding = () => {
  showRetakeConfirm.value = false;
  router.push('/onboarding/interview');
};

// Change Password logic
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
});
const passwordError = ref('');
const passwordSuccess = ref(false);
const isSubmittingPassword = ref(false);

const handleChangePassword = async () => {
  passwordError.value = '';
  passwordSuccess.value = false;

  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    passwordError.value = 'Mật khẩu xác nhận không khớp.';
    return;
  }

  isSubmittingPassword.value = true;
  try {
    await api.put('/auth/change-password', {
      old_password: passwordForm.value.oldPassword,
      new_password: passwordForm.value.newPassword
    });
    
    passwordSuccess.value = true;
    passwordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' };
    
    setTimeout(() => {
      passwordSuccess.value = false;
    }, 3000);
  } catch (err) {
    passwordError.value = err.response?.data?.message || err.message || 'Có lỗi xảy ra khi đổi mật khẩu.';
  } finally {
    isSubmittingPassword.value = false;
  }
};

const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString('vi-VN');
};

const formatCurrency = (value) => {
  if (!value) return '0 ₫';
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(value);
};

const calculateAge = (dobString) => {
  if (!dobString) return 'Chưa cập nhật';
  const dob = new Date(dobString);
  const diff_ms = Date.now() - dob.getTime();
  const age_dt = new Date(diff_ms); 
  return Math.abs(age_dt.getUTCFullYear() - 1970);
};

const riskColor = (score) => {
  if (score < 40) return 'text-emerald-400';
  if (score < 70) return 'text-amber-400';
  return 'text-red-400';
};

const riskBgColor = (score) => {
  if (score < 40) return 'bg-emerald-400';
  if (score < 70) return 'bg-amber-400';
  return 'bg-red-400';
};
</script>
