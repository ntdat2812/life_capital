<template>
  <div class="space-y-6 pb-20">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold font-outfit text-white">Tổng Quan Tài Sản</h1>
        <p class="text-xs text-slate-400">Bức tranh toàn cảnh về sức khỏe tài chính của bạn</p>
      </div>
      <router-link to="/review/active" class="px-4 py-2 text-sm bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl transition shadow-lg shadow-indigo-500/20 font-medium">
        Chốt Sổ Tháng Này
      </router-link>
    </div>

    <!-- Key Metrics Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="premium-card rounded-2xl p-6 relative overflow-hidden group">
        <div class="absolute inset-0 bg-gradient-to-br from-emerald-500/5 to-emerald-500/0"></div>
        <span class="text-xs text-slate-400 uppercase font-bold tracking-wider relative z-10">Tài Sản Ròng (Net Worth)</span>
        <h2 class="text-3xl font-bold mt-2 text-emerald-400 relative z-10">{{ formatCurrency(wealthStore.netWorthSummary.net_worth) }}</h2>
        <div v-if="wealthStore.netWorthSummary.previous_net_worth > 0" class="mt-4 flex items-center text-xs text-slate-400 relative z-10">
          <span :class="wealthStore.netWorthSummary.net_worth >= wealthStore.netWorthSummary.previous_net_worth ? 'text-emerald-400' : 'text-red-400'" class="mr-2 flex items-center font-medium">
            <svg v-if="wealthStore.netWorthSummary.net_worth >= wealthStore.netWorthSummary.previous_net_worth" class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path></svg>
            <svg v-else class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0v-8m0 8l-8-8-4 4-6-6"></path></svg>
            {{ (((wealthStore.netWorthSummary.net_worth - wealthStore.netWorthSummary.previous_net_worth) / wealthStore.netWorthSummary.previous_net_worth) * 100).toFixed(2) }}%
          </span>
          so với tháng trước
        </div>
        <div v-else class="mt-4 flex items-center text-xs text-slate-500 relative z-10">
          Chưa có dữ liệu tháng trước
        </div>
      </div>
      
      <div class="premium-card rounded-2xl p-6 relative group hover:border-indigo-500/50 transition-colors cursor-pointer" @click="$router.push('/assets')">
        <span class="text-xs text-slate-400 uppercase font-bold tracking-wider relative z-10">Tổng Tài Sản</span>
        <h2 class="text-3xl font-bold mt-2 text-white relative z-10">{{ formatCurrency(wealthStore.netWorthSummary.total_assets) }}</h2>
        <div class="mt-4 flex items-center justify-between text-xs text-slate-400 relative z-10">
          <span>Gồm {{ wealthStore.allAssets.length }} danh mục đầu tư</span>
          <span class="text-indigo-400 opacity-0 group-hover:opacity-100 transition-opacity flex items-center">
            Chi tiết <svg class="w-3 h-3 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
          </span>
        </div>
      </div>
      
      <div class="premium-card rounded-2xl p-6 relative overflow-hidden group hover:border-red-500/50 transition-colors cursor-pointer" @click="$router.push('/assets')">
        <div class="absolute inset-0 bg-gradient-to-br from-red-500/5 to-red-500/0"></div>
        <span class="text-xs text-slate-400 uppercase font-bold tracking-wider relative z-10">Tổng Nợ</span>
        <h2 class="text-3xl font-bold mt-2 text-red-400 relative z-10">{{ formatCurrency(wealthStore.netWorthSummary.total_liabilities) }}</h2>
        <div class="mt-4 flex items-center justify-between text-xs text-slate-400 relative z-10">
          <span>Tỷ lệ đòn bẩy: {{ leverageRatio }}%</span>
          <span class="text-red-400 opacity-0 group-hover:opacity-100 transition-opacity flex items-center">
            Chi tiết <svg class="w-3 h-3 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
          </span>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      
      <!-- Cột Trái: Chart Phân Bổ (2/3) -->
      <div class="lg:col-span-2 space-y-6">
        <div class="premium-card rounded-2xl p-6 h-[400px] flex flex-col">
          <h3 class="text-lg font-bold text-white mb-4">Cơ cấu Tài sản (Asset Allocation)</h3>
          <div class="flex-1 w-full flex items-center justify-center relative">
            <template v-if="wealthStore.loading">
              <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
            </template>
            <template v-else-if="Object.keys(wealthStore.assetAllocation).length === 0">
              <div class="text-slate-500 text-sm">Chưa có dữ liệu tài sản.</div>
            </template>
            <template v-else>
              <Doughnut :data="chartData" :options="chartOptions" />
            </template>
          </div>
        </div>
      </div>
      
      <!-- Cột Phải: FI Progress & Action Items (1/3) -->
      <div class="space-y-6">
        <div class="premium-card rounded-2xl p-6 bg-gradient-to-b from-indigo-900/40 to-slate-900/40 border border-indigo-500/20 relative overflow-hidden">
          <div class="absolute top-0 right-0 p-4 opacity-10">
            <svg class="w-24 h-24 text-indigo-400" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z"></path></svg>
          </div>
          
          <h3 class="text-lg font-bold text-white mb-2 relative z-10">Tiến trình Tự do tài chính (FI)</h3>
          <div class="text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 to-emerald-400 mb-6 relative z-10">
            {{ fiProgressPercent }}%
          </div>
          
          <div class="relative w-full bg-slate-800 rounded-full h-3 mb-2 z-10">
            <div class="h-3 rounded-full bg-gradient-to-r from-indigo-500 to-emerald-400 transition-all duration-1000" :style="{ width: `${Math.min(fiProgressPercent, 100)}%` }"></div>
          </div>
          
          <div class="flex justify-between text-xs text-slate-400 relative z-10">
            <span>Hiện tại: {{ formatCurrency(wealthStore.netWorthSummary.net_worth) }}</span>
            <span>Mục tiêu: {{ formatCurrency(fiTarget) }}</span>
          </div>
        </div>

        <div class="premium-card rounded-2xl p-6">
          <h3 class="text-sm font-bold text-slate-300 uppercase tracking-wider mb-4 flex items-center">
            Chỉ số Sức khỏe
            <div class="ml-2 relative group cursor-help">
              <svg class="w-4 h-4 text-slate-400 hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
              <div class="absolute right-0 w-64 p-3 bg-slate-800 text-xs text-slate-300 rounded-lg shadow-xl opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-20 top-full mt-2 border border-slate-700">
                <p class="mb-2"><strong>Tỷ lệ Tiền mặt:</strong> = (Tiền mặt + Gửi tiết kiệm) / Tổng Tài sản.</p>
                <p><strong>Tỷ lệ Nợ:</strong> = Tổng Nợ / Tổng Tài sản.</p>
              </div>
            </div>
          </h3>
          <ul class="space-y-4">
            <li>
              <div class="flex justify-between text-xs mb-1">
                <span class="text-slate-400">Tỷ lệ Tiền mặt (Khuyến nghị: 10-20%)</span>
                <span :class="cashRatio > 50 ? 'text-amber-400' : 'text-emerald-400'">{{ cashRatio }}%</span>
              </div>
              <div class="w-full bg-slate-800 rounded-full h-1.5">
                <div class="bg-indigo-400 h-1.5 rounded-full" :style="{ width: `${Math.min(cashRatio, 100)}%` }"></div>
              </div>
            </li>
            <li>
              <div class="flex justify-between text-xs mb-1">
                <span class="text-slate-400">Tỷ lệ Nợ (Đòn bẩy) (Khuyến nghị: <40%)</span>
                <span :class="leverageRatio > 50 ? 'text-red-400' : 'text-emerald-400'">{{ leverageRatio }}%</span>
              </div>
              <div class="w-full bg-slate-800 rounded-full h-1.5">
                <div class="bg-red-400 h-1.5 rounded-full" :style="{ width: `${Math.min(leverageRatio, 100)}%` }"></div>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useWealthStore } from '../stores/wealthStore'
import { useProfileStore } from '../stores/profileStore'
import { useIpsStore } from '../stores/ips'
import { Doughnut } from 'vue-chartjs'
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js'

ChartJS.register(ArcElement, Tooltip, Legend)

const wealthStore = useWealthStore()
const profileStore = useProfileStore()
const ipsStore = useIpsStore()

onMounted(async () => {
  await wealthStore.fetchAllForDashboard()
  if (!profileStore.profile) {
    await profileStore.fetchProfile()
  }
  await ipsStore.fetchLatestIps()
})

const formatCurrency = (value) => {
  if (value === undefined || value === null) return '0'
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: wealthStore.netWorthSummary?.base_currency || 'VND' }).format(value)
}

const leverageRatio = computed(() => {
  if (wealthStore.netWorthSummary.total_assets === 0) return 0
  return ((wealthStore.netWorthSummary.total_liabilities / wealthStore.netWorthSummary.total_assets) * 100).toFixed(1)
})

const cashRatio = computed(() => {
  if (wealthStore.netWorthSummary.total_assets === 0) return 0
  const cash = wealthStore.assetAllocation['cash'] || 0
  const deposit = wealthStore.assetAllocation['deposit'] || 0
  return (((cash + deposit) / wealthStore.netWorthSummary.total_assets) * 100).toFixed(1)
})

const fiTarget = computed(() => {
  return profileStore.profile?.fi_target_amount || 0
})

const fiProgressPercent = computed(() => {
  if (fiTarget.value === 0) return 0
  const nw = wealthStore.netWorthSummary.net_worth
  if (nw <= 0) return 0
  return ((nw / fiTarget.value) * 100).toFixed(1)
})

// --- CHART LOGIC ---
const categoryColors = {
  'cash': '#10b981', // emerald-500
  'deposit': '#059669', // emerald-600
  'gold': '#fbbf24', // amber-400
  'stock': '#6366f1', // indigo-500
  'fund': '#8b5cf6', // violet-500
  'crypto': '#f59e0b', // amber-500
  'real_estate': '#64748b' // slate-500
}

const categoryLabels = {
  'cash': 'Tiền mặt & Tiết kiệm',
  'gold': 'Vàng',
  'stock': 'Cổ phiếu',
  'fund': 'Chứng chỉ quỹ',
  'crypto': 'Crypto',
  'real_estate': 'Bất động sản'
}

const chartData = computed(() => {
  const alloc = { ...wealthStore.assetAllocation }
  
  if (alloc['deposit']) {
    alloc['cash'] = (alloc['cash'] || 0) + alloc['deposit']
    delete alloc['deposit']
  }

  const keys = Object.keys(alloc).filter(k => alloc[k] > 0)
  const total = keys.reduce((sum, k) => sum + alloc[k], 0)
  
  return {
    labels: keys.map(k => {
      const name = categoryLabels[k] || k
      const percent = total > 0 ? ((alloc[k] / total) * 100).toFixed(1) : 0
      
      let ipsTarget = 0
      if (ipsStore.latestIps?.target_allocation) {
        const targetAlloc = typeof ipsStore.latestIps.target_allocation === 'string' 
          ? JSON.parse(ipsStore.latestIps.target_allocation) 
          : ipsStore.latestIps.target_allocation
        ipsTarget = targetAlloc[k] || 0
      }
      
      return `${name} (${percent}% / IPS ${ipsTarget}%)`
    }),
    datasets: [
      {
        data: keys.map(k => alloc[k]),
        backgroundColor: keys.map(k => categoryColors[k] || '#94a3b8'),
        borderWidth: 0,
        hoverOffset: 4
      }
    ]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  cutout: '70%',
  plugins: {
    legend: {
      position: 'right',
      labels: {
        color: '#cbd5e1',
        padding: 20,
        font: {
          family: 'Outfit, sans-serif',
          size: 13
        }
      }
    },
    tooltip: {
      backgroundColor: 'rgba(15, 23, 42, 0.9)',
      titleColor: '#fff',
      bodyColor: '#cbd5e1',
      padding: 12,
      borderColor: 'rgba(99, 102, 241, 0.2)',
      borderWidth: 1,
      callbacks: {
        label: function(context) {
          let label = context.label || '';
          if (label) {
            label += ': ';
          }
          if (context.parsed !== null) {
            label += formatCurrency(context.parsed);
          }
          return label;
        }
      }
    }
  }
}
</script>
