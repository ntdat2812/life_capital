<template>
  <div class="min-h-screen p-4 md:p-8">
    <div class="max-w-4xl mx-auto">
      <div class="flex items-center justify-between mb-8">
        <h1 class="text-3xl font-bold text-white">Hồ Sơ Của Tôi</h1>
        <button @click="reTakeOnboarding" class="text-indigo-400 hover:text-indigo-300 text-sm flex items-center transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Làm Lại Khảo Sát
        </button>
      </div>

      <div v-if="loading" class="flex justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
      
      <div v-else-if="error" class="premium-card p-8 text-center text-red-400">
        {{ error }}
      </div>

      <div v-else-if="profile" class="grid grid-cols-1 md:grid-cols-3 gap-6">
        
        <!-- Cột Trái (Tóm tắt cá nhân & Rủi ro) -->
        <div class="col-span-1 space-y-6">
          <div class="premium-card p-6">
            <div class="flex items-center space-x-4 mb-6">
              <div class="w-16 h-16 rounded-full bg-gradient-to-tr from-indigo-500 to-violet-500 flex items-center justify-center text-2xl text-white font-bold">
                {{ userInitials }}
              </div>
              <div>
                <h3 class="text-xl font-semibold text-white">{{ authStore.user?.name || 'Nhà đầu tư' }}</h3>
                <p class="text-slate-400 text-sm">Cập nhật: {{ formatDate(profile.updated_at) }}</p>
              </div>
            </div>
            
            <div class="space-y-4">
              <div>
                <p class="text-slate-400 text-xs uppercase tracking-wider mb-1">Mức độ rủi ro</p>
                <div class="flex items-baseline space-x-2">
                  <span class="text-2xl font-bold" :class="riskColor(profile.risk_score)">{{ profile.risk_score }}/100</span>
                  <span class="text-sm font-medium" :class="riskColor(profile.risk_score)">{{ profile.risk_tolerance }}</span>
                </div>
                <!-- Thanh tiến trình rủi ro -->
                <div class="w-full bg-slate-700 rounded-full h-2 mt-2">
                  <div class="h-2 rounded-full transition-all duration-1000" :class="riskBgColor(profile.risk_score)" :style="`width: ${profile.risk_score}%`"></div>
                </div>
              </div>
            </div>
          </div>

          <div class="premium-card p-6">
            <h4 class="text-sm font-medium text-white mb-4 uppercase tracking-wider">Thông tin chung</h4>
            <ul class="space-y-3">
              <li class="flex justify-between items-center border-b border-white/5 pb-2">
                <span class="text-slate-400 text-sm">Tuổi</span>
                <span class="text-white text-sm font-medium">{{ calculateAge(profile.date_of_birth) }}</span>
              </li>
              <li class="flex justify-between items-center border-b border-white/5 pb-2">
                <span class="text-slate-400 text-sm">Hôn nhân</span>
                <span class="text-white text-sm font-medium">{{ profile.marital_status || 'Chưa cập nhật' }}</span>
              </li>
              <li class="flex justify-between items-center border-b border-white/5 pb-2">
                <span class="text-slate-400 text-sm">Thu nhập / tháng</span>
                <span class="text-emerald-400 text-sm font-medium">{{ formatCurrency(profile.total_monthly_income) }}</span>
              </li>
              <li class="flex justify-between items-center pb-2">
                <span class="text-slate-400 text-sm">Chi phí / tháng</span>
                <span class="text-amber-400 text-sm font-medium">{{ formatCurrency(profile.total_monthly_expense) }}</span>
              </li>
            </ul>
          </div>
        </div>

        <!-- Cột Phải (AI Insights & Constraints) -->
        <div class="col-span-1 md:col-span-2 space-y-6">
          <div class="premium-card p-6 bg-gradient-to-br from-slate-800/80 to-indigo-900/20 border-indigo-500/20">
            <div class="flex items-center mb-4">
              <span class="text-2xl mr-2">🎯</span>
              <h3 class="text-lg font-semibold text-white">Mục Tiêu Tự Do Tài Chính (FI)</h3>
            </div>
            <div class="text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 to-violet-400 mb-2">
              {{ formatCurrency(profile.fi_target_amount) }}
            </div>
            <p class="text-slate-400 text-sm leading-relaxed">
              Dựa trên mức chi tiêu hiện tại và các yếu tố đời sống, AI ước tính đây là con số bạn cần đạt được để có thể nghỉ hưu sớm và hoàn toàn tự do tài chính.
            </p>
          </div>

          <div class="premium-card p-6">
            <div class="flex items-center mb-6">
              <span class="text-xl mr-2">🧠</span>
              <h3 class="text-lg font-semibold text-white">Phân Tích & Rào Cản (AI)</h3>
            </div>
            
            <div v-if="profile.life_constraints && Object.keys(profile.life_constraints).length > 0" class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div v-for="(val, key) in profile.life_constraints" :key="key" class="bg-slate-900/50 p-4 rounded-xl border border-white/5 hover:border-white/10 transition-colors">
                <p class="text-slate-400 text-xs uppercase tracking-wider mb-1">{{ formatKey(key) }}</p>
                <p class="text-white text-sm font-medium">{{ val }}</p>
              </div>
            </div>
            <div v-else class="text-slate-400 text-sm italic">
              AI chưa trích xuất được rào cản tài chính nào từ cuộc hội thoại.
            </div>
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
    </div>
  </div>
</template>

<script setup>
import { onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useProfileStore } from '../stores/profileStore';
import { useAuthStore } from '../stores/authStore';

const router = useRouter();
const profileStore = useProfileStore();
const authStore = useAuthStore();

const profile = computed(() => profileStore.profile);
const loading = computed(() => profileStore.loading);
const error = computed(() => profileStore.error);

const userInitials = computed(() => {
  const name = authStore.user?.name || 'User';
  return name.charAt(0).toUpperCase();
});

onMounted(async () => {
  await profileStore.fetchProfile();
});

const reTakeOnboarding = () => {
  router.push('/onboarding/interview');
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

const formatKey = (key) => {
  return key.replace(/_/g, ' ');
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
