<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h1 class="text-3xl font-bold font-outfit text-white">Quản lý Danh mục & Theo dõi</h1>
        <p class="text-slate-400 mt-1">Danh sách tài sản sinh lời đang nắm giữ và các tài sản đưa vào tầm ngắm chờ mua.</p>
      </div>
    </div>

    <!-- Tabs Component -->
    <div class="border-b border-slate-800">
      <nav class="-mb-px flex space-x-8" aria-label="Tabs">
        <button 
          v-for="tab in tabs" 
          :key="tab.id"
          @click="currentTab = tab.id"
          :class="[
            currentTab === tab.id
              ? 'border-indigo-500 text-indigo-400'
              : 'border-transparent text-slate-400 hover:text-slate-300 hover:border-slate-300',
            'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm transition-colors'
          ]"
        >
          {{ tab.name }}
          <span v-if="tab.count" class="ml-2 bg-slate-800 text-slate-300 py-0.5 px-2.5 rounded-full text-xs">
            {{ tab.count }}
          </span>
        </button>
      </nav>
    </div>

    <!-- TAB: HOLDINGS -->
    <div v-if="currentTab === 'holdings'" class="space-y-6">
      
      <!-- Chart & Summary -->
      <div v-if="!portfolioStore.loading && portfolioStore.holdings.length > 0" class="glass-card rounded-2xl p-6 flex flex-col md:flex-row gap-8 items-center">
        <div class="w-full md:w-1/3 h-64 relative">
          <Doughnut :data="chartData" :options="chartOptions" />
        </div>
        <div class="w-full md:w-2/3 grid grid-cols-2 md:grid-cols-3 gap-4">
          <div v-for="(group, cat) in groupedHoldings" :key="cat" class="bg-slate-900/50 rounded-xl p-4 border border-slate-800">
            <h3 class="text-sm font-medium text-slate-400 mb-1">{{ getCategoryName(cat) }}</h3>
            <p class="text-lg font-bold text-white mb-2">{{ formatCurrency(group.totalValue) }}</p>
            <div class="flex justify-between items-end">
              <div>
                <p class="text-xs text-slate-500">Tỷ trọng TT</p>
                <p class="text-sm font-bold text-indigo-400">{{ ((group.totalValue / totalInvestableValue) * 100).toFixed(1) }}%</p>
              </div>
              <div class="text-right">
                <p class="text-xs text-slate-500">IPS Target</p>
                <p class="text-sm font-medium text-slate-300">{{ getTargetAllocation({category: cat}) }}%</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="portfolioStore.loading" class="glass-card rounded-2xl p-8 text-center text-slate-500">
        Đang tải dữ liệu...
      </div>
      <div v-else-if="portfolioStore.holdings.length === 0" class="glass-card rounded-2xl p-8 text-center text-slate-500">
        Chưa có tài sản sinh lời nào. Vui lòng thêm tại mục Tài sản.
      </div>
      
      <!-- Grouped Tables -->
      <div v-else v-for="(group, cat) in groupedHoldings" :key="'table-'+cat" class="glass-card rounded-2xl overflow-hidden">
        <div class="p-4 border-b border-slate-800/50 flex justify-between items-center bg-slate-900/30">
          <h2 class="text-lg font-bold text-slate-100 flex items-center gap-2">
            <span class="text-indigo-400">{{ getCategoryIcon(cat) }}</span> Nhóm {{ getCategoryName(cat) }}
          </h2>
          <span class="text-sm font-medium text-slate-400">
            <span v-if="cat === 'gold'" class="mr-4 border-r border-slate-700 pr-4">
              Khối lượng: <span class="text-amber-400 font-bold">{{ formatGoldVolume(group.totalGoldVolume) }}</span>
            </span>
            Tổng: <span class="text-white">{{ formatCurrency(group.totalValue) }}</span>
          </span>
        </div>
        
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="border-b border-slate-800 bg-slate-900/50">
                <th class="p-4 text-sm font-semibold text-slate-400 cursor-pointer hover:text-white transition-colors select-none group/th" @click="setSort('name')">
                  Mã / Tên <span class="text-indigo-400 ml-1">{{ sortCol === 'name' ? (sortDir === 'asc' ? '↑' : '↓') : '' }}</span>
                </th>
                <th class="p-4 text-sm font-semibold text-slate-400 text-right cursor-pointer hover:text-white transition-colors select-none group/th" @click="setSort('quantity')">
                  Số lượng <span class="text-indigo-400 ml-1">{{ sortCol === 'quantity' ? (sortDir === 'asc' ? '↑' : '↓') : '' }}</span>
                </th>
                <th class="p-4 text-sm font-semibold text-slate-400 text-right cursor-pointer hover:text-white transition-colors select-none group/th" @click="setSort('avg_price')">
                  Giá vốn <span class="text-indigo-400 ml-1">{{ sortCol === 'avg_price' ? (sortDir === 'asc' ? '↑' : '↓') : '' }}</span>
                </th>
                <th class="p-4 text-sm font-semibold text-slate-400 text-right cursor-pointer hover:text-white transition-colors select-none group/th" @click="setSort('current_price')">
                  Giá hiện tại <span class="text-indigo-400 ml-1">{{ sortCol === 'current_price' ? (sortDir === 'asc' ? '↑' : '↓') : '' }}</span>
                </th>
                <th class="p-4 text-sm font-semibold text-slate-400 text-right cursor-pointer hover:text-white transition-colors select-none group/th" @click="setSort('current_value')">
                  Tổng giá trị <span class="text-indigo-400 ml-1">{{ sortCol === 'current_value' ? (sortDir === 'asc' ? '↑' : '↓') : '' }}</span>
                </th>
                <th class="p-4 text-sm font-semibold text-slate-400 text-right cursor-pointer hover:text-white transition-colors select-none group/th" @click="setSort('current_value')">
                  Tỷ trọng % <span class="text-indigo-400 ml-1">{{ sortCol === 'current_value' ? (sortDir === 'asc' ? '↑' : '↓') : '' }}</span>
                </th>
                <th class="p-4 text-sm font-semibold text-slate-400 text-center">Luận điểm (Thesis)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800/50">
              <tr v-for="asset in group.items" :key="asset.id" class="bg-slate-900/20 hover:bg-slate-800/30 transition-colors">
                <td class="p-4">
                  <div class="flex flex-col">
                    <span class="font-bold text-slate-200">{{ asset.ticker || asset.name }}</span>
                    <span v-if="asset.ticker" class="text-xs text-slate-500">{{ asset.name }}</span>
                  </div>
                </td>
                <td class="p-4 text-right text-slate-300 font-medium">{{ formatNumber(asset.quantity) }}</td>
                <td class="p-4 text-right text-slate-300">{{ formatCurrency(asset.avg_price) }}</td>
                <td class="p-4 text-right">
                  <div class="flex flex-col items-end">
                    <span :class="getPriceClass(asset.current_price, asset.avg_price)">
                      {{ formatCurrency(asset.current_price) }}
                    </span>
                    <span v-if="asset.avg_price > 0 && asset.current_price > 0" class="text-[11px] font-medium mt-0.5" :class="getPriceClass(asset.current_price, asset.avg_price)">
                      {{ asset.current_price >= asset.avg_price ? '+' : '' }}{{ (((asset.current_price - asset.avg_price) / asset.avg_price) * 100).toFixed(2) }}%
                    </span>
                  </div>
                </td>
                <td class="p-4 text-right font-bold text-white">{{ formatCurrency(asset.current_value) }}</td>
                <td class="p-4 text-right font-medium text-emerald-400">
                  {{ totalInvestableValue > 0 ? ((asset.current_value / totalInvestableValue) * 100).toFixed(1) : '0.0' }}%
                </td>
                <td class="p-4 text-center">
                  <router-link v-if="asset.ticker || asset.category !== 'cash'" :to="`/thesis/${asset.ticker || encodeURIComponent(asset.name)}`" class="inline-flex items-center justify-center p-1.5 bg-slate-800 hover:bg-indigo-600/20 text-slate-400 hover:text-indigo-400 rounded-lg transition-colors group" title="Xem Luận điểm đầu tư">
                    <svg class="w-5 h-5 group-hover:scale-110 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
                  </router-link>
                  <span v-else class="text-slate-600">-</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- TAB: WATCHLIST -->
    <div v-else-if="currentTab === 'watchlist'" class="space-y-6">
      <div class="glass-card rounded-2xl overflow-hidden">
        <div class="p-6 border-b border-slate-800/50 flex justify-between items-center bg-slate-900/30">
          <h2 class="text-xl font-bold text-slate-100 flex items-center gap-2">
            <span class="text-amber-400">👀</span> Danh sách theo dõi
          </h2>
          <button @click="openAddWatchlistModal" class="btn-primary text-sm px-4 py-2">
            + Thêm Mã Theo Dõi
          </button>
        </div>
        
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="border-b border-slate-800 bg-slate-900/50">
                <th class="p-4 text-sm font-semibold text-slate-400">Mã / Công ty</th>
                <th class="p-4 text-sm font-semibold text-slate-400 text-center">Độ tự tin (1-10)</th>
                <th class="p-4 text-sm font-semibold text-slate-400 text-right">Giá mục tiêu mua</th>
                <th class="p-4 text-sm font-semibold text-slate-400 text-right">Giá thị trường</th>
                <th class="p-4 text-sm font-semibold text-slate-400 text-center">Trạng thái</th>
                <th class="p-4 text-sm font-semibold text-slate-400 text-right">Hành động</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800/50">
              <tr v-if="portfolioStore.loading" class="bg-slate-900/20">
                <td colspan="6" class="p-8 text-center text-slate-500">Đang tải dữ liệu...</td>
              </tr>
              <tr v-else-if="!portfolioStore.watchlist || portfolioStore.watchlist.length === 0" class="bg-slate-900/20">
                <td colspan="6" class="p-8 text-center text-slate-500">Chưa có mã nào trong danh sách theo dõi. Bấm nút Thêm để bắt đầu.</td>
              </tr>
              <tr v-else v-for="item in portfolioStore.watchlist" :key="item.id" class="bg-slate-900/20 hover:bg-slate-800/30 transition-colors group">
                <td class="p-4">
                  <div class="flex items-center gap-3">
                    <router-link :to="`/thesis/${item.ticker}`" class="font-bold text-indigo-400 hover:text-indigo-300 underline-offset-4 hover:underline">
                      {{ item.ticker }}
                    </router-link>
                    <span class="text-sm text-slate-400 truncate max-w-[200px]" :title="item.company_name">{{ item.company_name }}</span>
                  </div>
                </td>
                <td class="p-4 text-center">
                  <div class="flex justify-center">
                    <div class="flex items-center justify-center w-8 h-8 rounded-full font-bold text-xs" :class="getScoreClass(item.priority)">
                      {{ item.priority }}
                    </div>
                  </div>
                </td>
                <td class="p-4 text-right font-medium text-emerald-400">{{ formatCurrency(item.target_price) }}</td>
                <td class="p-4 text-right text-slate-300">{{ item.current_price > 0 ? formatCurrency(item.current_price) : '-' }}</td>
                <td class="p-4 text-center">
                  <span v-if="item.current_price > 0 && item.current_price <= item.target_price" class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 animate-pulse-slow">
                    🔥 Sẵn sàng mua
                  </span>
                  <span v-else class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-slate-800 text-slate-400 border border-slate-700">
                    Đang theo dõi
                  </span>
                </td>
                <td class="p-4 text-right">
                  <button @click="confirmDeleteWatchlist(item.id, item.ticker)" class="text-slate-500 hover:text-rose-400 transition-colors p-2 opacity-0 group-hover:opacity-100">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- MODAL THÊM WATCHLIST -->
    <div v-if="showWatchlistModal" class="fixed inset-0 bg-slate-900/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="premium-card w-full max-w-md p-6 relative">
        <button @click="showWatchlistModal = false" class="absolute top-4 right-4 text-slate-400 hover:text-white">✕</button>
        <h2 class="text-2xl font-bold text-white mb-6">Thêm Mã Theo Dõi</h2>
        
        <form @submit.prevent="submitWatchlist" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Mã Tài Sản (Ticker)</label>
            <input type="text" v-model="watchlistForm.ticker" required placeholder="VD: HPG, VCB, BTC" class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500 uppercase" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Tên Công Ty / Dự Án</label>
            <input type="text" v-model="watchlistForm.company_name" placeholder="VD: Tập đoàn Hòa Phát" class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Giá mục tiêu (Vùng giá mua an toàn)</label>
            <input type="number" step="0.01" v-model="watchlistForm.target_price" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Độ tự tin (Conviction Score: 1-10)</label>
            <div class="flex items-center gap-4">
              <input type="range" min="1" max="10" v-model="watchlistForm.priority" class="flex-1 accent-indigo-500" />
              <span class="font-bold text-lg text-indigo-400 w-8 text-center">{{ watchlistForm.priority }}</span>
            </div>
            <p class="text-xs text-slate-500 mt-1">10: Cơ hội không thể bỏ lỡ | 1: Rủi ro cao, theo dõi cho vui</p>
          </div>
          
          <div v-if="portfolioStore.error" class="p-3 bg-rose-500/10 border border-rose-500/20 rounded-lg">
            <p class="text-sm text-rose-400">{{ portfolioStore.error }}</p>
          </div>
          
          <div class="pt-4 flex gap-3">
            <button type="button" @click="showWatchlistModal = false" class="flex-1 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors font-medium">Hủy</button>
            <button type="submit" class="flex-1 btn-primary" :disabled="portfolioStore.loading">
              {{ portfolioStore.loading ? 'Đang lưu...' : 'Thêm vào danh sách' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { usePortfolioStore } from '../stores/portfolioStore'
import { useIpsStore } from '../stores/ips'
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js'
import { Doughnut } from 'vue-chartjs'
import { getCategoryName, getCategoryIcon } from '../utils/assetUtils'

ChartJS.register(ArcElement, Tooltip, Legend)

const portfolioStore = usePortfolioStore()
const ipsStore = useIpsStore()

const currentTab = ref('holdings')
const tabs = computed(() => [
  { id: 'holdings', name: 'Danh mục thực tế', count: portfolioStore.holdings.length },
  { id: 'watchlist', name: 'Danh sách theo dõi', count: portfolioStore.watchlist.length }
])

onMounted(async () => {
  await Promise.all([
    portfolioStore.fetchHoldings(),
    portfolioStore.fetchWatchlist(),
    ipsStore.fetchLatestIps()
  ])
})

const sortCol = ref('current_value')
const sortDir = ref('desc')

const setSort = (col) => {
  if (sortCol.value === col) {
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortCol.value = col
    sortDir.value = 'desc'
  }
}

const formatCurrency = (value) => {
  if (value === null || value === undefined) return '-'
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(value)
}

const formatNumber = (value) => {
  if (!value) return '-'
  return new Intl.NumberFormat('en-US', { maximumFractionDigits: 4 }).format(value)
}

const getPriceClass = (current, avg) => {
  if (!current || !avg) return 'text-slate-300'
  if (current > avg) return 'text-emerald-400 font-medium'
  if (current < avg) return 'text-rose-400 font-medium'
  return 'text-slate-300'
}

const getScoreClass = (score) => {
  if (score >= 8) return 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
  if (score >= 5) return 'bg-amber-500/20 text-amber-400 border border-amber-500/30'
  return 'bg-rose-500/20 text-rose-400 border border-rose-500/30'
}

const getGoldMultiplier = (unit) => {
  if (unit === 'Loại 1 Lượng' || unit === 'Lượng') return 1;
  if (unit === 'Loại 5 Chỉ') return 0.5;
  if (unit === 'Loại 2 Chỉ') return 0.2;
  if (unit === 'Loại 1 Chỉ' || unit === 'Chỉ') return 0.1;
  if (unit === 'Loại 0.5 Chỉ') return 0.05;
  if (unit === 'Phân') return 0.01;
  return 1;
}

const formatGoldVolume = (volumeInLuong) => {
  if (!volumeInLuong || volumeInLuong <= 0) return '0 Lượng';
  
  let totalPhan = Math.round(volumeInLuong * 100);
  const luong = Math.floor(totalPhan / 100);
  totalPhan %= 100;
  
  const chi = Math.floor(totalPhan / 10);
  const phan = totalPhan % 10;
  
  let parts = [];
  if (luong > 0) parts.push(`${luong} Lượng`);
  if (chi > 0) parts.push(`${chi} Chỉ`);
  if (phan > 0) parts.push(`${phan} Phân`);
  
  return parts.length > 0 ? parts.join(' ') : '0 Lượng';
}

const groupedHoldings = computed(() => {
  const groups = {}
  portfolioStore.holdings.forEach(asset => {
    if (!groups[asset.category]) {
      groups[asset.category] = {
        items: [],
        totalValue: 0,
        totalGoldVolume: 0
      }
    }
    groups[asset.category].items.push(asset)
    groups[asset.category].totalValue += (asset.current_value || 0)
    
    if (asset.category === 'gold') {
      const match = asset.name?.match(/\((.+?)\)$/);
      const unit = match ? match[1] : 'Lượng';
      groups[asset.category].totalGoldVolume += (asset.quantity || 0) * getGoldMultiplier(unit);
    }
  })
  
  // Sort items within each group
  for (const cat in groups) {
    groups[cat].items.sort((a, b) => {
      let valA = a[sortCol.value]
      let valB = b[sortCol.value]
      
      if (sortCol.value === 'name') {
        valA = (a.ticker || a.name || '').toLowerCase()
        valB = (b.ticker || b.name || '').toLowerCase()
        return sortDir.value === 'asc' ? valA.localeCompare(valB) : valB.localeCompare(valA)
      }
      
      valA = valA || 0
      valB = valB || 0
      
      if (valA < valB) return sortDir.value === 'asc' ? -1 : 1
      if (valA > valB) return sortDir.value === 'asc' ? 1 : -1
      return 0
    })
  }
  
  return groups
})

const totalInvestableValue = computed(() => {
  return portfolioStore.holdings.reduce((sum, item) => sum + (item.current_value || 0), 0)
})

const getTargetAllocation = (asset) => {
  if (!ipsStore.latestIps || !ipsStore.latestIps.target_allocation) return 0
  try {
    let allocation = ipsStore.latestIps.target_allocation
    if (typeof allocation === 'string') {
      allocation = JSON.parse(allocation)
    }
    return allocation[asset.category] || 0
  } catch (e) {
    return 0
  }
}

// Chart Data
const chartData = computed(() => {
  const groups = groupedHoldings.value
  return {
    labels: Object.keys(groups).map(k => getCategoryName(k)),
    datasets: [{
      data: Object.values(groups).map(g => g.totalValue),
      backgroundColor: [
        '#6366f1', // indigo-500
        '#10b981', // emerald-500
        '#f59e0b', // amber-500
        '#ef4444', // rose-500
        '#8b5cf6', // violet-500
        '#0ea5e9'  // sky-500
      ],
      borderWidth: 0,
      hoverOffset: 4
    }]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { 
      position: 'right', 
      labels: { color: '#94a3b8', font: { family: 'Inter', size: 12 } } 
    },
    tooltip: {
      callbacks: {
        label: function(context) {
          let label = context.label || '';
          if (label) label += ': ';
          if (context.parsed !== null) {
            label += new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(context.parsed);
          }
          return label;
        }
      }
    }
  },
  cutout: '70%'
}

// Watchlist Form Logic
const showWatchlistModal = ref(false)
const watchlistForm = ref({
  ticker: '',
  company_name: '',
  target_price: null,
  priority: 5,
  status: 'watching'
})

const openAddWatchlistModal = () => {
  watchlistForm.value = {
    ticker: '',
    company_name: '',
    target_price: null,
    priority: 5,
    status: 'watching'
  }
  showWatchlistModal.value = true
}

const submitWatchlist = async () => {
  try {
    const payload = {
      ...watchlistForm.value,
      ticker: watchlistForm.value.ticker.toUpperCase(),
      target_price: Number(watchlistForm.value.target_price)
    }
    await portfolioStore.addWatchlist(payload)
    showWatchlistModal.value = false
  } catch (err) {
    // Error handled in store
  }
}

const confirmDeleteWatchlist = async (id, ticker) => {
  if (confirm(`Bạn có chắc muốn xóa mã ${ticker} khỏi danh sách theo dõi?`)) {
    await portfolioStore.deleteWatchlist(id)
  }
}
</script>
