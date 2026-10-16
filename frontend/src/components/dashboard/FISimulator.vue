<template>
  <div class="premium-card rounded-2xl p-6 bg-slate-900 border border-slate-700/50">
    <div class="flex flex-col md:flex-row gap-8">
      
      <!-- Cột Trái: Control Panel (Inputs) -->
      <div class="w-full md:w-1/3 flex flex-col justify-between">
        <div>
          <div class="flex items-center mb-6">
            <span class="text-3xl mr-3">🔮</span>
            <div>
              <h3 class="text-xl font-bold text-white">FI Simulator</h3>
              <p class="text-xs text-slate-400">Mô phỏng bản đồ Tự do tài chính</p>
            </div>
          </div>
          
          <div class="space-y-6">
            <!-- Lợi nhuận kỳ vọng -->
            <div>
              <div class="flex justify-between items-end mb-2">
                <label class="text-sm font-medium text-slate-300">Lợi nhuận kỳ vọng</label>
                <span class="text-lg font-bold text-indigo-400">{{ expectedReturn }}%/năm</span>
              </div>
              <input 
                type="range" 
                v-model.number="expectedReturn" 
                min="0" max="30" step="1"
                class="w-full h-2 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-indigo-500"
              />
              <div class="flex justify-between text-xs text-slate-500 mt-1">
                <span>0%</span>
                <span>30%</span>
              </div>
              
              <!-- Tính toán sinh lời hiện tại -->
              <div v-if="currentNetWorth > 0" class="mt-3 text-xs text-slate-400 bg-indigo-900/20 p-2.5 rounded-lg border border-indigo-500/20">
                <span class="block mb-1">💰 Sinh lời trên tổng tài sản hiện tại:</span>
                <span class="text-indigo-300 font-bold text-sm">{{ formatCurrency((currentNetWorth * (expectedReturn / 100)) / 12) }}/tháng</span>
                <span class="text-slate-500 ml-1">({{ formatCurrency(currentNetWorth * (expectedReturn / 100)) }}/năm)</span>
              </div>
            </div>

            <!-- Dòng tiền tiết kiệm hàng tháng -->
            <div>
              <div class="flex items-center mb-2">
                <label class="block text-sm font-medium text-slate-300">Dòng tiền đầu tư mỗi tháng</label>
                <div class="ml-2 relative group cursor-help">
                  <svg class="w-4 h-4 text-slate-400 hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                  <div class="absolute left-0 w-64 p-3 bg-slate-800 text-xs text-slate-300 rounded-lg shadow-xl opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50 top-full mt-2 border border-slate-700 pointer-events-none">
                    <p>💡 Con số mặc định được đặt ở mức 10 triệu để làm ví dụ. Bạn có thể tự do nhập số tiền tiết kiệm thực tế hàng tháng của mình vào ô này để xem biểu đồ thay đổi ra sao nhé!</p>
                  </div>
                </div>
              </div>
              <CurrencyInput v-model="monthlyContribution" class="mb-3" />
            </div>
          </div>
        </div>

        <!-- Card 1: Cột mốc Tự do tài chính (Hero Card riêng biệt) -->
        <div class="mt-6 p-4 rounded-2xl bg-gradient-to-r from-emerald-950/40 via-slate-900 to-indigo-950/40 border border-emerald-500/30 relative overflow-hidden shadow-lg shadow-emerald-500/5">
          <!-- Subtle decorative background icon -->
          <div class="absolute -right-2 -bottom-2 opacity-10 pointer-events-none">
            <svg class="w-20 h-20 text-emerald-400" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 14.5v-9l6 4.5-6 4.5z"/></svg>
          </div>

          <div class="flex items-center gap-2 mb-2">
            <span class="text-base">🎯</span>
            <span class="text-xs font-bold uppercase tracking-wider text-emerald-400">Cột mốc Tự do tài chính (FI)</span>
          </div>

          <div v-if="yearsToFI > 0 && yearsToFI <= 40">
            <div class="flex items-baseline gap-2">
              <span class="text-3xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-teal-200 font-outfit">
                {{ yearsToFI }} năm nữa
              </span>
            </div>
            <div class="mt-2 inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-500/10 text-emerald-300 text-xs border border-emerald-500/20">
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
              Dự kiến cán đích vào năm <strong class="font-bold text-white ml-0.5">{{ new Date().getFullYear() + yearsToFI }}</strong>
            </div>
          </div>
          <div v-else-if="yearsToFI === 0" class="text-xl font-bold text-emerald-400 flex items-center gap-2">
            <span>🎉</span> Đã đạt mục tiêu FI!
          </div>
          <div v-else class="text-sm font-medium text-amber-400 bg-amber-500/10 p-2.5 rounded-lg border border-amber-500/20 mt-1">
            ⚠️ Cần tăng thêm vốn tích lũy tháng hoặc kỳ vọng lợi nhuận để đạt mục tiêu.
          </div>
        </div>

        <!-- Card 2: Mô phỏng tài sản theo thanh kéo -->
        <div class="mt-4 p-4 bg-slate-800/40 border border-slate-700/50 rounded-xl">
          <div class="flex justify-between items-end mb-2">
            <label class="text-xs font-medium text-slate-400">Ước tính tài sản sau:</label>
            <span class="text-xs font-bold text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded border border-emerald-500/20">{{ targetYears }} năm</span>
          </div>
          <input 
            type="range" 
            v-model.number="targetYears" 
            min="1" max="30" step="1"
            class="w-full h-2 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-emerald-500"
          />
          <div class="mt-3 flex justify-between items-center pt-2.5 border-t border-slate-700/40">
            <span class="text-xs text-slate-400">Tài sản dự kiến:</span>
            <span class="text-base font-bold text-white font-outfit">{{ formatCurrency(projectedAmountAtTargetYear) }}</span>
          </div>
        </div>
      </div>

      <!-- Cột Phải: Biểu đồ (Chart) -->
      <div class="w-full md:w-2/3 h-[400px] relative">
        <Line :data="chartData" :options="chartOptions" />
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'
import CurrencyInput from '../common/CurrencyInput.vue'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const props = defineProps({
  currentNetWorth: {
    type: Number,
    required: true
  },
  fiTarget: {
    type: Number,
    required: true
  },
  defaultMonthlySavings: {
    type: Number,
    default: 0
  }
})

const expectedReturn = ref(10) // 10% annual
const monthlyContribution = ref(0)
const yearsToFI = ref(-1)
const targetYears = ref(5)

const projectedAmountAtTargetYear = computed(() => {
  let p = props.currentNetWorth
  const m = monthlyContribution.value
  const r = expectedReturn.value / 100
  const rm = Math.pow(1 + r, 1 / 12) - 1 // Monthly compound rate
  
  for (let year = 1; year <= targetYears.value; year++) {
    for(let month = 0; month < 12; month++) {
        p = p * (1 + rm) + m
    }
  }
  return p
})

// Cập nhật giá trị default ban đầu
onMounted(() => {
  if (props.defaultMonthlySavings > 0) {
    monthlyContribution.value = props.defaultMonthlySavings
  }
})

watch(() => props.defaultMonthlySavings, (newVal) => {
  if (monthlyContribution.value === 0 && newVal > 0) {
    monthlyContribution.value = newVal
  }
})

const formatCurrency = (val) => {
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(val || 0)
}

// Hàm tính toán Net Worth sau N năm
const calculateProjection = () => {
  const currentYear = new Date().getFullYear()
  const years = []
  const dataNW = []
  const dataFI = []
  
  let p = props.currentNetWorth
  const target = props.fiTarget
  const m = monthlyContribution.value
  const r = expectedReturn.value / 100
  const rm = Math.pow(1 + r, 1 / 12) - 1 // Monthly compound rate
  
  // Nếu đã đạt FI từ đầu
  if (p >= target && target > 0) {
    yearsToFI.value = 0
    return {
      labels: [currentYear],
      nwData: [p],
      fiData: [target]
    }
  }

  let reachedYear = -1
  const MAX_YEARS = 30 // Giới hạn biểu đồ tối đa 30 năm
  
  for (let year = 0; year <= MAX_YEARS; year++) {
    years.push(currentYear + year)
    dataNW.push(p)
    dataFI.push(target)
    
    if (reachedYear === -1 && p >= target) {
      reachedYear = year
    }
    
    // Tính toán cho năm tiếp theo (cộng dồn 12 tháng)
    for(let month = 0; month < 12; month++) {
        p = p * (1 + rm) + m
    }
  }

  yearsToFI.value = reachedYear

  // Cắt biểu đồ sau khi đạt FI vài năm cho đẹp, không vẽ hết 30 năm nếu đạt sớm
  const displayLength = reachedYear !== -1 ? Math.min(reachedYear + 3, MAX_YEARS + 1) : MAX_YEARS + 1
  
  return {
    labels: years.slice(0, displayLength),
    nwData: dataNW.slice(0, displayLength),
    fiData: dataFI.slice(0, displayLength)
  }
}

// Cấu hình Chart.js
const chartData = computed(() => {
  const projection = calculateProjection()
  
  return {
    labels: projection.labels,
    datasets: [
      {
        label: 'Dự phóng Tài sản',
        data: projection.nwData,
        borderColor: '#6366f1', // indigo-500
        backgroundColor: 'rgba(99, 102, 241, 0.1)',
        borderWidth: 3,
        pointBackgroundColor: '#8b5cf6',
        pointRadius: 0,
        pointHoverRadius: 6,
        fill: true,
        tension: 0.4
      },
      {
        label: 'Mục tiêu Tự do Tài chính',
        data: projection.fiData,
        borderColor: '#10b981', // emerald-500
        borderWidth: 2,
        borderDash: [5, 5],
        pointRadius: 0,
        fill: false,
        tension: 0
      }
    ]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    mode: 'index',
    intersect: false,
  },
  plugins: {
    legend: {
      position: 'top',
      labels: {
        color: '#cbd5e1',
        font: { family: 'Outfit, sans-serif' }
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
          let label = context.dataset.label || '';
          if (label) {
            label += ': ';
          }
          if (context.parsed.y !== null) {
            label += formatCurrency(context.parsed.y);
          }
          return label;
        }
      }
    }
  },
  scales: {
    y: {
      grid: {
        color: 'rgba(255, 255, 255, 0.05)'
      },
      ticks: {
        color: '#94a3b8',
        callback: function(value) {
          if (value >= 1000000000) {
            return (value / 1000000000).toFixed(1) + ' Tỷ';
          }
          if (value >= 1000000) {
            return (value / 1000000).toFixed(0) + ' Tr';
          }
          return value;
        }
      }
    },
    x: {
      grid: {
        display: false
      },
      ticks: {
        color: '#94a3b8',
        maxTicksLimit: 8
      }
    }
  }
}
</script>
