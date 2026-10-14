<template>
  <div class="h-full flex flex-col justify-between">
    <!-- Header Controls & Quick Stats -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4">
      <!-- Quick Highlights -->
      <div class="flex flex-wrap items-center gap-3 text-xs">
        <!-- Chế độ Theo Tháng -->
        <template v-if="timeframe === 'month'">
          <div v-if="latestMoMChange" class="flex items-center gap-1.5 bg-slate-800/80 px-2.5 py-1.5 rounded-lg border border-slate-700/50">
            <span class="text-slate-400">Tháng gần nhất:</span>
            <span :class="latestMoMChange.amount >= 0 ? 'text-emerald-400 font-bold' : 'text-rose-400 font-bold'" class="flex items-center">
              <svg v-if="latestMoMChange.amount >= 0" class="w-3 h-3 mr-0.5 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path></svg>
              <svg v-else class="w-3 h-3 mr-0.5 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0v-8m0 8l-8-8-4 4-6-6"></path></svg>
              {{ latestMoMChange.amount >= 0 ? '+' : '' }}{{ formatCurrency(latestMoMChange.amount) }}
              ({{ latestMoMChange.percent >= 0 ? '+' : '' }}{{ latestMoMChange.percent.toFixed(1) }}%)
            </span>
          </div>

          <div v-if="yearGrowth" class="flex items-center gap-1.5 bg-slate-800/80 px-2.5 py-1.5 rounded-lg border border-slate-700/50">
            <span class="text-slate-400">Cả năm {{ selectedYear === 'all' ? 'lũy kế' : selectedYear }}:</span>
            <span :class="yearGrowth.amount >= 0 ? 'text-emerald-400 font-bold' : 'text-rose-400 font-bold'" class="flex items-center">
              <svg v-if="yearGrowth.amount >= 0" class="w-3 h-3 mr-0.5 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path></svg>
              <svg v-else class="w-3 h-3 mr-0.5 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0v-8m0 8l-8-8-4 4-6-6"></path></svg>
              {{ yearGrowth.amount >= 0 ? '+' : '' }}{{ formatCurrency(yearGrowth.amount) }}
              ({{ yearGrowth.percent >= 0 ? '+' : '' }}{{ yearGrowth.percent.toFixed(1) }}%)
            </span>
          </div>
        </template>

        <!-- Chế độ Theo Năm -->
        <template v-else>
          <div v-if="latestYearItem" class="flex items-center gap-1.5 bg-slate-800/80 px-2.5 py-1.5 rounded-lg border border-slate-700/50">
            <span class="text-slate-400">Năm {{ latestYearItem.year }}:</span>
            <span :class="latestYearItem.growthAmount >= 0 ? 'text-emerald-400 font-bold' : 'text-rose-400 font-bold'" class="flex items-center">
              <svg v-if="latestYearItem.growthAmount >= 0" class="w-3 h-3 mr-0.5 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path></svg>
              <svg v-else class="w-3 h-3 mr-0.5 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0v-8m0 8l-8-8-4 4-6-6"></path></svg>
              {{ latestYearItem.growthAmount >= 0 ? '+' : '' }}{{ formatCurrency(latestYearItem.growthAmount) }}
              ({{ latestYearItem.growthPercent >= 0 ? '+' : '' }}{{ latestYearItem.growthPercent.toFixed(1) }}%)
            </span>
          </div>

          <div class="flex items-center gap-1.5 bg-slate-800/80 px-2.5 py-1.5 rounded-lg border border-slate-700/50 text-slate-400">
            <span>Theo dõi: <strong class="text-white">{{ yearlyTimeline.length }} năm</strong></span>
          </div>
        </template>
      </div>

      <!-- Action Buttons (Timeframe / Year filter & Table toggle) -->
      <div class="flex flex-wrap items-center gap-2 self-end sm:self-auto">
        <!-- Timeframe Toggle: Theo Tháng vs Theo Năm -->
        <div class="flex items-center bg-slate-800/90 p-0.5 rounded-lg border border-slate-700/50 text-xs">
          <button 
            type="button"
            @click="timeframe = 'month'"
            :class="timeframe === 'month' ? 'bg-indigo-600 text-white font-medium shadow' : 'text-slate-400 hover:text-white'"
            class="px-2.5 py-1 rounded transition"
          >
            Theo Tháng
          </button>
          <button 
            type="button"
            @click="timeframe = 'year'"
            :class="timeframe === 'year' ? 'bg-indigo-600 text-white font-medium shadow' : 'text-slate-400 hover:text-white'"
            class="px-2.5 py-1 rounded transition"
          >
            Theo Năm
          </button>
        </div>

        <!-- Year filter pills (chỉ hiện khi xem Theo Tháng và có nhiều năm) -->
        <div v-if="timeframe === 'month' && availableYears.length > 1" class="flex items-center bg-slate-800/80 p-0.5 rounded-lg border border-slate-700/50 text-xs">
          <button 
            type="button"
            @click="selectedYear = 'all'"
            :class="selectedYear === 'all' ? 'bg-slate-700 text-white font-medium shadow' : 'text-slate-400 hover:text-white'"
            class="px-2 py-1 rounded transition"
          >
            Tất cả
          </button>
          <button 
            v-for="year in availableYears" 
            :key="year"
            type="button"
            @click="selectedYear = year"
            :class="selectedYear === year ? 'bg-slate-700 text-white font-medium shadow' : 'text-slate-400 hover:text-white'"
            class="px-2 py-1 rounded transition"
          >
            {{ year }}
          </button>
        </div>

        <!-- View mode: Chart vs Table -->
        <div class="flex items-center bg-slate-800/80 p-0.5 rounded-lg border border-slate-700/50 text-xs">
          <button 
            type="button"
            @click="viewMode = 'chart'"
            :class="viewMode === 'chart' ? 'bg-slate-700 text-white font-medium shadow' : 'text-slate-400 hover:text-white'"
            class="px-2 py-1 rounded transition flex items-center gap-1"
            title="Xem dạng biểu đồ"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"></path></svg>
            <span class="hidden md:inline">Biểu đồ</span>
          </button>
          <button 
            type="button"
            @click="viewMode = 'table'"
            :class="viewMode === 'table' ? 'bg-slate-700 text-white font-medium shadow' : 'text-slate-400 hover:text-white'"
            class="px-2 py-1 rounded transition flex items-center gap-1"
            title="Xem dạng bảng"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M3 14h18m-9-4v8m-7 4h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
            <span class="hidden md:inline">Bảng</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Main Chart / Table Area -->
    <div class="flex-1 w-full relative min-h-[260px] flex items-center justify-center">
      <!-- Loading State -->
      <div v-if="loading" class="flex flex-col items-center justify-center">
        <div class="animate-spin rounded-full h-10 w-10 border-t-2 border-b-2 border-indigo-500"></div>
        <p class="text-xs text-slate-400 mt-2">Đang tải lịch sử tăng trưởng...</p>
      </div>

      <!-- Empty State (chưa đủ mốc dữ liệu) -->
      <div v-else-if="isEmpty" class="text-center p-6 bg-slate-800/40 border border-slate-700/50 rounded-xl max-w-md mx-auto">
        <div class="w-12 h-12 bg-indigo-500/10 rounded-full flex items-center justify-center mx-auto mb-3 text-indigo-400 text-2xl">
          📈
        </div>
        <h4 class="text-sm font-bold text-white mb-1">Cần thêm dữ liệu chốt sổ</h4>
        <p class="text-xs text-slate-400 mb-4">
          Hệ thống cần ít nhất 2 kỳ chốt sổ để vẽ đường cong tăng trưởng. Hãy thực hiện chốt sổ định kỳ nhé!
        </p>
        <router-link to="/review/active" class="inline-flex items-center text-xs px-3 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition font-medium">
          Chốt sổ tháng ngay
          <svg class="w-3.5 h-3.5 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </router-link>
      </div>

      <!-- Chart View -->
      <div v-else-if="viewMode === 'chart'" class="w-full h-full relative">
        <Line :data="chartData" :options="chartOptions" />
      </div>

      <!-- Table View: Theo Tháng -->
      <div v-else-if="timeframe === 'month'" class="w-full h-full overflow-y-auto max-h-[300px] rounded-xl border border-slate-700/50">
        <table class="w-full text-left text-xs">
          <thead class="bg-slate-800/90 text-slate-400 sticky top-0 uppercase tracking-wider text-[10px]">
            <tr>
              <th class="py-2.5 px-4 font-semibold">Kỳ Chốt Sổ</th>
              <th class="py-2.5 px-4 font-semibold text-right">Tài Sản Ròng</th>
              <th class="py-2.5 px-4 font-semibold text-right">Biến Động (VND)</th>
              <th class="py-2.5 px-4 font-semibold text-right">% Tăng Trưởng</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800">
            <tr 
              v-for="(item, idx) in tableRowsMonth" 
              :key="idx"
              class="hover:bg-slate-800/50 transition-colors"
            >
              <td class="py-2.5 px-4 font-medium text-slate-200">
                {{ item.label }}
                <span v-if="item.isCurrent" class="ml-1.5 px-1.5 py-0.5 bg-indigo-500/20 text-indigo-400 rounded text-[10px]">Hiện tại</span>
              </td>
              <td class="py-2.5 px-4 text-right font-bold text-white">
                {{ formatCurrency(item.netWorth) }}
              </td>
              <td class="py-2.5 px-4 text-right font-medium">
                <span v-if="item.diffAmount !== null" :class="item.diffAmount >= 0 ? 'text-emerald-400' : 'text-rose-400'">
                  {{ item.diffAmount >= 0 ? '+' : '' }}{{ formatCurrency(item.diffAmount) }}
                </span>
                <span v-else class="text-slate-500">—</span>
              </td>
              <td class="py-2.5 px-4 text-right font-medium">
                <span v-if="item.diffPercent !== null" :class="item.diffPercent >= 0 ? 'text-emerald-400' : 'text-rose-400'">
                  {{ item.diffPercent >= 0 ? '+' : '' }}{{ item.diffPercent.toFixed(2) }}%
                </span>
                <span v-else class="text-slate-500">—</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Table View: Theo Năm -->
      <div v-else class="w-full h-full overflow-y-auto max-h-[300px] rounded-xl border border-slate-700/50">
        <table class="w-full text-left text-xs">
          <thead class="bg-slate-800/90 text-slate-400 sticky top-0 uppercase tracking-wider text-[10px]">
            <tr>
              <th class="py-2.5 px-4 font-semibold">Năm</th>
              <th class="py-2.5 px-4 font-semibold text-right">Đầu Năm</th>
              <th class="py-2.5 px-4 font-semibold text-right">Cuối Năm</th>
              <th class="py-2.5 px-4 font-semibold text-right">Tăng Trưởng Cả Năm</th>
              <th class="py-2.5 px-4 font-semibold text-right">% Cả Năm</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800">
            <tr 
              v-for="item in tableRowsYear" 
              :key="item.year"
              class="hover:bg-slate-800/50 transition-colors"
            >
              <td class="py-2.5 px-4 font-medium text-slate-200">
                Năm {{ item.year }}
                <span v-if="item.isCurrentYear" class="ml-1.5 px-1.5 py-0.5 bg-indigo-500/20 text-indigo-400 rounded text-[10px]">Đến nay</span>
              </td>
              <td class="py-2.5 px-4 text-right text-slate-400">
                {{ formatCurrency(item.startNetWorth) }}
              </td>
              <td class="py-2.5 px-4 text-right font-bold text-white">
                {{ formatCurrency(item.endNetWorth) }}
              </td>
              <td class="py-2.5 px-4 text-right font-medium">
                <span :class="item.growthAmount >= 0 ? 'text-emerald-400' : 'text-rose-400'">
                  {{ item.growthAmount >= 0 ? '+' : '' }}{{ formatCurrency(item.growthAmount) }}
                </span>
              </td>
              <td class="py-2.5 px-4 text-right font-medium">
                <span :class="item.growthPercent >= 0 ? 'text-emerald-400' : 'text-rose-400'">
                  {{ item.growthPercent >= 0 ? '+' : '' }}{{ item.growthPercent.toFixed(2) }}%
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
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
  history: {
    type: Array,
    default: () => []
  },
  currentNetWorth: {
    type: Number,
    default: 0
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const timeframe = ref('month') // 'month' | 'year'
const viewMode = ref('chart') // 'chart' | 'table'
const selectedYear = ref('all') // 'all' | number

const formatCurrency = (value) => {
  if (value === undefined || value === null) return '0 ₫'
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(value)
}

// 1. Chuẩn hóa chuỗi thời gian chi tiết theo tháng (Oldest -> Newest)
const fullTimeline = computed(() => {
  const points = []

  if (props.history && props.history.length > 0) {
    props.history.forEach(rev => {
      if (rev.review_month && rev.net_worth_at_review > 0) {
        const d = new Date(rev.review_month)
        points.push({
          date: d,
          year: d.getFullYear(),
          month: d.getMonth() + 1,
          label: `T${d.getMonth() + 1}/${String(d.getFullYear()).slice(-2)}`,
          fullLabel: `Tháng ${d.getMonth() + 1}/${d.getFullYear()}`,
          netWorth: rev.net_worth_at_review,
          isCurrent: false
        })
      }
    })
  }

  // Sắp xếp theo ngày tăng dần
  points.sort((a, b) => a.date - b.date)

  // Bổ sung điểm mốc hiện tại nếu mốc hiện tại khác tháng với điểm chốt gần nhất
  if (props.currentNetWorth > 0) {
    const now = new Date()
    const currentMonth = now.getMonth() + 1
    const currentYear = now.getFullYear()

    const hasCurrentMonth = points.some(p => p.year === currentYear && p.month === currentMonth)
    if (!hasCurrentMonth) {
      points.push({
        date: now,
        year: currentYear,
        month: currentMonth,
        label: `T${currentMonth}/${String(currentYear).slice(-2)} (Nay)`,
        fullLabel: `Hiện tại (Tháng ${currentMonth}/${currentYear})`,
        netWorth: props.currentNetWorth,
        isCurrent: true
      })
    }
  }

  return points
})

// Danh sách các năm có dữ liệu
const availableYears = computed(() => {
  const years = fullTimeline.value.map(p => p.year)
  return Array.from(new Set(years)).sort((a, b) => b - a)
})

// Dữ liệu theo năm được lọc (dành cho chế độ Theo Tháng)
const filteredTimeline = computed(() => {
  if (selectedYear.value === 'all') {
    return fullTimeline.value
  }
  return fullTimeline.value.filter(p => p.year === selectedYear.value)
})

// 2. Chuẩn hóa chuỗi thời gian gom theo từng Năm (Year-over-Year)
const yearlyTimeline = computed(() => {
  const mapByYear = {}
  fullTimeline.value.forEach(p => {
    if (!mapByYear[p.year]) {
      mapByYear[p.year] = []
    }
    mapByYear[p.year].push(p)
  })

  const years = Object.keys(mapByYear).map(Number).sort((a, b) => a - b)
  const result = []

  years.forEach((yr, idx) => {
    const pointsInYear = mapByYear[yr]
    // Mốc đầu năm: nếu có năm trước, lấy mốc cuối của năm trước để đảm bảo tính liên tục tài chính
    let startPoint = pointsInYear[0]
    if (idx > 0) {
      const prevYearPoints = mapByYear[years[idx - 1]]
      startPoint = prevYearPoints[prevYearPoints.length - 1]
    }
    // Mốc cuối năm: mốc gần nhất trong năm
    const endPoint = pointsInYear[pointsInYear.length - 1]

    const startNetWorth = startPoint.netWorth
    const endNetWorth = endPoint.netWorth
    const growthAmount = endNetWorth - startNetWorth
    const growthPercent = startNetWorth > 0 ? (growthAmount / startNetWorth) * 100 : 0
    const isCurrentYear = yr === new Date().getFullYear()

    result.push({
      year: yr,
      label: isCurrentYear ? `${yr} (Nay)` : `${yr}`,
      fullLabel: isCurrentYear ? `Năm ${yr} (Đến nay)` : `Năm ${yr}`,
      startNetWorth,
      endNetWorth,
      growthAmount,
      growthPercent,
      isCurrentYear
    })
  })

  return result
})

const latestYearItem = computed(() => {
  if (yearlyTimeline.value.length === 0) return null
  return yearlyTimeline.value[yearlyTimeline.value.length - 1]
})

const isEmpty = computed(() => {
  if (timeframe.value === 'month') {
    return filteredTimeline.value.length < 2
  }
  return yearlyTimeline.value.length === 0 || fullTimeline.value.length < 2
})

// Dữ liệu bảng - Theo Tháng (Mới nhất lên đầu)
const tableRowsMonth = computed(() => {
  const items = [...filteredTimeline.value]
  items.reverse()

  return items.map(item => {
    const originalIdx = fullTimeline.value.findIndex(p => p.label === item.label)
    let diffAmount = null
    let diffPercent = null

    if (originalIdx > 0) {
      const prev = fullTimeline.value[originalIdx - 1]
      diffAmount = item.netWorth - prev.netWorth
      diffPercent = prev.netWorth > 0 ? (diffAmount / prev.netWorth) * 100 : 0
    }

    return {
      ...item,
      diffAmount,
      diffPercent
    }
  })
})

// Dữ liệu bảng - Theo Năm (Mới nhất lên đầu)
const tableRowsYear = computed(() => {
  const items = [...yearlyTimeline.value]
  items.reverse()
  return items
})

// Thống kê tháng gần nhất (MoM)
const latestMoMChange = computed(() => {
  const list = fullTimeline.value
  if (list.length < 2) return null
  const current = list[list.length - 1]
  const prev = list[list.length - 2]
  const amount = current.netWorth - prev.netWorth
  const percent = prev.netWorth > 0 ? (amount / prev.netWorth) * 100 : 0
  return { amount, percent }
})

// Thống kê cả năm (YTD)
const yearGrowth = computed(() => {
  const list = filteredTimeline.value
  if (list.length < 2) return null
  const first = list[0]
  const latest = list[list.length - 1]
  const amount = latest.netWorth - first.netWorth
  const percent = first.netWorth > 0 ? (amount / first.netWorth) * 100 : 0
  return { amount, percent }
})

// Cấu hình Chart.js linh hoạt theo timeframe
const chartData = computed(() => {
  if (timeframe.value === 'year') {
    return {
      labels: yearlyTimeline.value.map(p => p.label),
      datasets: [
        {
          label: 'Tài sản cuối năm',
          data: yearlyTimeline.value.map(p => p.endNetWorth),
          borderColor: '#10b981', // emerald-500
          borderWidth: 3,
          tension: 0.25,
          fill: true,
          backgroundColor: (context) => {
            const ctx = context.chart.ctx
            const gradient = ctx.createLinearGradient(0, 0, 0, 300)
            gradient.addColorStop(0, 'rgba(16, 185, 129, 0.3)')
            gradient.addColorStop(1, 'rgba(16, 185, 129, 0.0)')
            return gradient
          },
          pointBackgroundColor: '#6366f1',
          pointBorderColor: '#ffffff',
          pointBorderWidth: 2,
          pointRadius: 6,
          pointHoverRadius: 9
        }
      ]
    }
  }

  // Theo Tháng
  return {
    labels: filteredTimeline.value.map(p => p.label),
    datasets: [
      {
        label: 'Tài sản ròng',
        data: filteredTimeline.value.map(p => p.netWorth),
        borderColor: '#6366f1', // indigo-500
        borderWidth: 3,
        tension: 0.35,
        fill: true,
        backgroundColor: (context) => {
          const ctx = context.chart.ctx
          const gradient = ctx.createLinearGradient(0, 0, 0, 300)
          gradient.addColorStop(0, 'rgba(99, 102, 241, 0.35)')
          gradient.addColorStop(1, 'rgba(99, 102, 241, 0.0)')
          return gradient
        },
        pointBackgroundColor: '#10b981', // emerald-500
        pointBorderColor: '#ffffff',
        pointBorderWidth: 2,
        pointRadius: 4,
        pointHoverRadius: 7,
        pointHoverBackgroundColor: '#6366f1'
      }
    ]
  }
})

const chartOptions = computed(() => {
  return {
    responsive: true,
    maintainAspectRatio: false,
    layout: {
      padding: {
        top: 10,
        bottom: 5,
        left: 5,
        right: 15
      }
    },
    scales: {
      x: {
        grid: {
          display: false
        },
        ticks: {
          color: '#94a3b8',
          font: {
            family: 'Outfit, sans-serif',
            size: 11
          }
        }
      },
      y: {
        grid: {
          color: 'rgba(51, 65, 85, 0.4)'
        },
        ticks: {
          color: '#94a3b8',
          font: {
            family: 'Outfit, sans-serif',
            size: 11
          },
          callback: function(value) {
            if (value >= 1000000000) {
              return (value / 1000000000).toFixed(1) + ' tỷ'
            }
            if (value >= 1000000) {
              return (value / 1000000).toFixed(0) + ' tr'
            }
            return value
          }
        }
      }
    },
    plugins: {
      legend: {
        display: false
      },
      tooltip: {
        backgroundColor: 'rgba(15, 23, 42, 0.95)',
        titleColor: '#ffffff',
        bodyColor: '#cbd5e1',
        borderColor: 'rgba(99, 102, 241, 0.3)',
        borderWidth: 1,
        padding: 12,
        boxPadding: 4,
        callbacks: {
          title: function(context) {
            const idx = context[0].dataIndex
            if (timeframe.value === 'year') {
              const item = yearlyTimeline.value[idx]
              return item ? item.fullLabel : context[0].label
            }
            const item = filteredTimeline.value[idx]
            return item ? item.fullLabel : context[0].label
          },
          label: function(context) {
            const val = context.parsed.y
            if (timeframe.value === 'year') {
              const idx = context.dataIndex
              const item = yearlyTimeline.value[idx]
              let str = `Tài sản cuối năm: ${formatCurrency(val)}`
              if (item) {
                const sign = item.growthAmount >= 0 ? '+' : ''
                str += ` (Tăng cả năm: ${sign}${formatCurrency(item.growthAmount)} / ${sign}${item.growthPercent.toFixed(1)}%)`
              }
              return str
            }

            let str = `Tài sản: ${formatCurrency(val)}`
            const idx = context.dataIndex
            if (idx > 0) {
              const prevVal = filteredTimeline.value[idx - 1].netWorth
              const diff = val - prevVal
              const pct = prevVal > 0 ? (diff / prevVal) * 100 : 0
              const sign = diff >= 0 ? '+' : ''
              str += ` (${sign}${formatCurrency(diff)} / ${sign}${pct.toFixed(1)}%)`
            }
            return str
          }
        }
      }
    }
  }
})
</script>
