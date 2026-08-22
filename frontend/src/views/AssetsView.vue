<template>
  <div class="assets-view max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
    <header class="mb-8 flex justify-between items-end">
      <div>
        <h1 class="text-3xl font-bold text-white mb-2">Quản lý Tài sản & Nợ</h1>
        <p class="text-slate-400">Theo dõi toàn bộ bức tranh tài chính cá nhân của bạn.</p>
      </div>
      <div class="text-sm text-slate-400">
        Đơn vị tiền tệ: <span class="text-indigo-400 font-bold ml-1">{{ wealthStore.netWorthSummary?.base_currency || 'VND' }}</span>
      </div>
    </header>

    <!-- Khu vực Top: Summary Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
      <div class="p-6 rounded-2xl bg-gradient-to-br from-emerald-900/40 to-slate-900 border border-emerald-500/30 shadow-[0_0_15px_rgba(16,185,129,0.1)] relative overflow-hidden">
        <div class="absolute -top-4 -right-4 opacity-10 text-emerald-400 text-8xl">💰</div>
        <h3 class="text-emerald-200/80 text-sm font-semibold uppercase tracking-wider mb-2 relative z-10">Tổng Tài Sản</h3>
        <p class="text-4xl font-bold text-emerald-400 relative z-10">
          {{ formatCurrency(wealthStore.netWorthSummary?.total_assets) }}
        </p>
      </div>
      
      <div class="p-6 rounded-2xl bg-gradient-to-br from-amber-900/40 to-slate-900 border border-amber-500/30 shadow-[0_0_15px_rgba(245,158,11,0.1)] relative overflow-hidden">
        <div class="absolute -top-4 -right-4 opacity-10 text-amber-500 text-8xl">💳</div>
        <h3 class="text-amber-200/80 text-sm font-semibold uppercase tracking-wider mb-2 flex justify-between relative z-10">
          <span>Tổng Nợ</span>
          <span class="text-amber-300 bg-amber-500/20 px-2 py-0.5 rounded text-xs" v-if="wealthStore.netWorthSummary?.total_liabilities > 0">
            {{ debtToAssetRatio }}% Đòn bẩy
          </span>
        </h3>
        <p class="text-4xl font-bold text-amber-500 relative z-10">
          {{ formatCurrency(wealthStore.netWorthSummary?.total_liabilities) }}
        </p>
      </div>

      <div class="p-6 rounded-2xl bg-gradient-to-br from-indigo-900/50 to-violet-900/20 border border-indigo-500/40 shadow-[0_0_20px_rgba(99,102,241,0.15)] relative overflow-hidden group">
        <div class="absolute inset-0 bg-gradient-to-br from-indigo-500/20 to-violet-500/20 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
        <div class="relative z-10">
          <h3 class="text-indigo-200 text-sm font-semibold uppercase tracking-wider mb-2">Tài Sản Ròng (Net Worth)</h3>
          <p class="text-5xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-indigo-300 to-violet-300">
            {{ formatCurrency(wealthStore.netWorthSummary?.net_worth) }}
          </p>
        </div>
      </div>
    </div>

    <!-- Toolbar Filters -->
    <div class="flex flex-col sm:flex-row justify-between items-center mb-6 gap-4 bg-slate-800/40 p-4 rounded-xl border border-slate-700/50">
      <div class="flex flex-wrap items-center gap-4 w-full sm:w-auto">
        <div class="flex items-center gap-3">
          <label class="text-sm font-medium text-slate-300 whitespace-nowrap">Lọc Tài Sản:</label>
          <CustomSelect v-model="wealthStore.assetCategoryFilter" :options="assetCategoryFilterOptions" @update:modelValue="applyAssetFilter" class="w-[180px]" />
        </div>
        <div class="flex items-center gap-3">
          <label class="text-sm font-medium text-slate-300 whitespace-nowrap">Sắp xếp:</label>
          <CustomSelect v-model="wealthStore.assetSort" :options="assetSortOptions" @update:modelValue="applyAssetFilter" class="w-[180px]" />
        </div>
      </div>
      <div class="flex flex-wrap items-center gap-4 w-full sm:w-auto">
        <div class="flex items-center gap-3">
          <label class="text-sm font-medium text-slate-300 whitespace-nowrap">Lọc Nợ:</label>
          <CustomSelect v-model="wealthStore.liabilityCategoryFilter" :options="liabilityCategoryFilterOptions" @update:modelValue="applyLiabilityFilter" class="w-[180px]" />
        </div>
        <div class="flex items-center gap-3">
          <label class="text-sm font-medium text-slate-300 whitespace-nowrap">Sắp xếp:</label>
          <CustomSelect v-model="wealthStore.liabilitySort" :options="liabilitySortOptions" @update:modelValue="applyLiabilityFilter" class="w-[180px]" />
        </div>
      </div>
    </div>

    <!-- Khu vực Bottom: 2 Columns Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      
      <!-- CỘT TRÁI: TÀI SẢN -->
      <div>
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-bold text-white flex items-center">
            Tài sản của bạn
          </h2>
          <div class="flex gap-2">
            <button @click="triggerPriceSync" :disabled="syncingPrices" class="px-4 py-2 bg-emerald-500/20 text-emerald-400 rounded-lg hover:bg-emerald-500/30 transition-colors text-sm font-medium border border-emerald-500/30 flex items-center gap-2">
              <span v-if="syncingPrices" class="animate-spin">↻</span>
              <span v-else>↻</span>
              Đồng bộ giá
            </button>
            <button @click="openAddAsset" class="px-4 py-2 bg-indigo-500/20 text-indigo-400 rounded-lg hover:bg-indigo-500/30 transition-colors text-sm font-medium border border-indigo-500/30">
              + Thêm
            </button>
          </div>
        </div>

        <div class="premium-card overflow-hidden flex flex-col h-[550px]">
          <div v-if="wealthStore.loading && wealthStore.assets.length === 0" class="p-8 text-center text-slate-400 flex-1 flex items-center justify-center">
            Đang tải dữ liệu...
          </div>
          <div v-else-if="groupedAssets.length === 0" class="p-8 text-center text-slate-500 italic flex-1 flex items-center justify-center">
            Không tìm thấy tài sản nào.
          </div>
          <div v-else class="divide-y divide-slate-700/50 overflow-y-auto custom-scrollbar flex-1" @scroll="handleAssetScroll">
            <!-- Group by Category -->
            <div v-for="group in groupedAssets" :key="group.category" class="group-container">
              <!-- Header của nhóm -->
              <div class="bg-slate-800/60 px-4 py-3 flex justify-between items-center sticky top-0 z-10 cursor-pointer hover:bg-slate-700/60 transition-colors" @click="toggleAssetGroup(group.category)">
                <div class="flex items-center gap-3">
                  <svg class="w-4 h-4 text-slate-500 transition-transform duration-200" :class="{ '-rotate-90': !expandedAssetGroups[group.category] }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
                  <span class="text-xl">{{ getCategoryIcon(group.category) }}</span>
                  <h3 class="text-white font-semibold">{{ group.categoryName }}</h3>
                </div>
                <div class="text-right flex items-center gap-3">
                  <div v-if="group.category === 'gold'" class="mr-4 text-right pr-4 border-r border-slate-700 hidden sm:block">
                    <div class="px-2.5 py-0.5 bg-amber-500/10 border border-amber-500/20 rounded-md">
                      <span class="text-amber-400 font-medium text-xs whitespace-nowrap">{{ formatGoldVolume(group.totalGoldVolume) }}</span>
                    </div>
                  </div>
                  <div class="text-right flex-1 sm:flex-none">
                    <p class="text-emerald-400 font-bold leading-tight">{{ formatCurrency(group.total_value) }}</p>
                    <p class="text-[11px] text-slate-400">Chiếm {{ group.percentage }}%</p>
                  </div>
                  <div class="ml-2 pl-2 border-l border-slate-700 flex items-center justify-end gap-1 w-[72px]">
                    <div v-if="group.category === 'gold'" class="relative group/tooltip">
                      <button @click.stop="openBulkUpdateGold(group.items)" class="p-1.5 text-slate-400 hover:text-amber-400 hover:bg-amber-500/10 rounded transition-colors">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
                      </button>
                      <div class="absolute top-full right-0 mt-2 hidden group-hover/tooltip:block whitespace-nowrap bg-slate-900 text-xs text-slate-200 px-2.5 py-1.5 rounded-lg border border-slate-700 shadow-xl z-50">
                        Cập nhật giá vàng
                      </div>
                    </div>
                    <div class="relative group/tooltip">
                      <button @click.stop="openGroupDetail(group, 'asset')" class="p-1.5 text-slate-400 hover:text-indigo-400 hover:bg-indigo-500/10 rounded transition-colors">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 3.055A9.001 9.001 0 1020.945 13H11V3.055z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.488 9H15V3.512A9.025 9.025 0 0120.488 9z"></path></svg>
                      </button>
                      <div class="absolute top-full right-0 mt-2 hidden group-hover/tooltip:block whitespace-nowrap bg-slate-900 text-xs text-slate-200 px-2.5 py-1.5 rounded-lg border border-slate-700 shadow-xl z-50">
                        Xem biểu đồ cơ cấu
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Danh sách items trong nhóm -->
              <ul v-show="expandedAssetGroups[group.category]" class="divide-y divide-slate-700/30">
                <li v-for="asset in group.items" :key="asset.id" class="p-4 pl-12 hover:bg-slate-800/50 transition-colors flex justify-between items-center group/item">
                  <div>
                    <h4 class="text-slate-200 font-medium">{{ asset.name }}</h4>
                    <p class="text-xs text-slate-400 mt-1">
                      <span v-if="asset.ticker">{{ asset.ticker }}</span>
                      <span v-if="asset.ticker && asset.quantity"> • </span>
                      <span v-if="asset.quantity">SL: {{ formatInputNumber(asset.quantity) }}</span>
                    </p>
                  </div>
                  <div class="text-right">
                    <p class="text-slate-300">{{ formatCurrency(asset.current_value) }}</p>
                    <div class="flex items-center justify-end gap-3 mt-1">
                      <span class="text-xs text-slate-500 w-12 text-right">{{ asset.percentage }}%</span>
                      <div class="flex gap-2 opacity-0 group-hover/item:opacity-100 transition-opacity">
                        <button @click.stop="openEditAsset(asset)" class="text-xs text-indigo-400 hover:text-indigo-300">Sửa</button>
                        <button @click.stop="deleteAsset(asset.id)" class="text-xs text-red-400 hover:text-red-300">Xóa</button>
                      </div>
                    </div>
                  </div>
                </li>
              </ul>
            </div>
            <!-- Loading indicator for infinite scroll -->
            <div v-if="wealthStore.loading && wealthStore.assets.length > 0" class="p-4 text-center text-slate-500 text-sm">
              Đang tải thêm...
            </div>
          </div>
        </div>
      </div>

      <!-- CỘT PHẢI: KHOẢN NỢ -->
      <div>
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-bold text-white flex items-center">
            Các khoản nợ
          </h2>
          <button @click="openAddLiability" class="px-4 py-2 bg-amber-500/20 text-amber-500 rounded-lg hover:bg-amber-500/30 transition-colors text-sm font-medium border border-amber-500/30">
            + Thêm Khoản Nợ
          </button>
        </div>

        <div class="premium-card overflow-hidden flex flex-col h-[550px]">
          <div v-if="wealthStore.loading && wealthStore.liabilities.length === 0" class="p-8 text-center text-slate-400 flex-1 flex items-center justify-center">
            Đang tải dữ liệu...
          </div>
          <div v-else-if="groupedLiabilities.length === 0" class="p-8 text-center text-slate-500 italic flex-1 flex items-center justify-center">
            Không tìm thấy khoản nợ nào.
          </div>
          <div v-else class="divide-y divide-slate-700/50 overflow-y-auto custom-scrollbar flex-1" @scroll="handleLiabilityScroll">
            <!-- Group by Category -->
            <div v-for="group in groupedLiabilities" :key="group.category" class="group-container">
              <div class="bg-slate-800/60 px-4 py-3 flex justify-between items-center sticky top-0 z-10 cursor-pointer hover:bg-slate-700/60 transition-colors" @click="toggleLiabilityGroup(group.category)">
                <div class="flex items-center gap-3">
                  <svg class="w-4 h-4 text-slate-500 transition-transform duration-200" :class="{ '-rotate-90': !expandedLiabilityGroups[group.category] }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
                  <span class="text-xl">{{ getLiabilityIcon(group.category) }}</span>
                  <h3 class="text-white font-semibold">{{ group.categoryName }}</h3>
                </div>
                <div class="text-right flex items-center gap-3">
                  <div class="text-right flex-1 sm:flex-none">
                    <p class="text-amber-500 font-bold leading-tight">{{ formatCurrency(group.total_value) }}</p>
                    <p class="text-[11px] text-slate-400">Chiếm {{ group.percentage }}%</p>
                  </div>
                  <div class="ml-2 pl-2 border-l border-slate-700 flex items-center justify-end gap-1 w-[40px]">
                    <div class="relative group/tooltip">
                      <button @click.stop="openGroupDetail(group, 'liability')" class="p-1.5 text-slate-400 hover:text-amber-400 hover:bg-amber-500/10 rounded transition-colors">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 3.055A9.001 9.001 0 1020.945 13H11V3.055z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.488 9H15V3.512A9.025 9.025 0 0120.488 9z"></path></svg>
                      </button>
                      <div class="absolute top-full right-0 mt-2 hidden group-hover/tooltip:block whitespace-nowrap bg-slate-900 text-xs text-slate-200 px-2.5 py-1.5 rounded-lg border border-slate-700 shadow-xl z-50">
                        Xem chi tiết nhóm
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              
              <ul v-show="expandedLiabilityGroups[group.category]" class="divide-y divide-slate-700/30">
                <li v-for="liability in group.items" :key="liability.id" class="p-4 pl-12 hover:bg-slate-800/50 transition-colors flex justify-between items-center group/item">
                  <div>
                    <h4 class="text-slate-200 font-medium">{{ liability.name }}</h4>
                    <p class="text-xs text-slate-400 mt-1">
                      <span v-if="liability.interest_rate">Lãi: {{ (liability.interest_rate * 100).toFixed(1) }}%/năm</span>
                    </p>
                  </div>
                  <div class="text-right">
                    <p class="text-slate-300">{{ formatCurrency(liability.remaining_balance) }}</p>
                    <div class="flex items-center justify-end gap-3 mt-1">
                      <span class="text-xs text-slate-500 w-12 text-right">{{ liability.percentage }}%</span>
                      <div class="flex gap-2 opacity-0 group-hover/item:opacity-100 transition-opacity">
                        <button @click.stop="openEditLiability(liability)" class="text-xs text-indigo-400 hover:text-indigo-300">Sửa</button>
                        <button @click.stop="deleteLiability(liability.id)" class="text-xs text-red-400 hover:text-red-300">Xóa</button>
                      </div>
                    </div>
                  </div>
                </li>
              </ul>
            </div>
            <!-- Loading indicator for infinite scroll -->
            <div v-if="wealthStore.loading && wealthStore.liabilities.length > 0" class="p-4 text-center text-slate-500 text-sm">
              Đang tải thêm...
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- GROUP DETAIL MODAL -->
    <div v-if="showGroupDetailModal" class="fixed inset-0 bg-slate-900/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="premium-card w-full max-w-lg p-0 relative max-h-[85vh] flex flex-col">
        <div class="p-6 pb-4 border-b border-slate-700/50 flex justify-between items-start">
          <div>
            <h2 class="text-2xl font-bold text-white mb-1 flex items-center gap-2">
              <span>{{ selectedGroupType === 'asset' ? getCategoryIcon(selectedGroup.category) : getLiabilityIcon(selectedGroup.category) }}</span>
              Phân bổ: {{ selectedGroup.categoryName }}
            </h2>
            <p class="text-slate-400 text-sm">
              Tổng giá trị: <span :class="selectedGroupType === 'asset' ? 'text-emerald-400' : 'text-amber-500'" class="font-bold">{{ formatCurrency(selectedGroup.total_value) }}</span>
            </p>
          </div>
          <button @click="showGroupDetailModal = false" class="text-slate-400 hover:text-white">✕</button>
        </div>
        
        <div class="p-6 overflow-y-auto custom-scrollbar flex-1">
          <ul class="space-y-5">
            <li v-for="item in selectedGroup.items" :key="item.id">
              <div class="flex justify-between items-end mb-2">
                <h4 class="text-slate-200 font-medium">{{ item.name }}</h4>
                <p class="text-sm font-bold" :class="selectedGroupType === 'asset' ? 'text-emerald-300' : 'text-amber-400'">
                  {{ formatCurrency(selectedGroupType === 'asset' ? item.current_value : item.remaining_balance) }}
                </p>
              </div>
              <!-- Tính % trong nội bộ nhóm -->
              <div class="flex items-center gap-3">
                <div class="flex-1 h-2.5 bg-slate-800 rounded-full overflow-hidden border border-slate-700">
                  <div 
                    class="h-full rounded-full transition-all duration-1000" 
                    :class="selectedGroupType === 'asset' ? 'bg-gradient-to-r from-indigo-500 to-emerald-400' : 'bg-gradient-to-r from-orange-500 to-amber-400'"
                    :style="{ width: `${((selectedGroupType === 'asset' ? item.current_value : item.remaining_balance) / selectedGroup.total_value * 100)}%` }"
                  ></div>
                </div>
                <span class="text-xs font-medium text-slate-300 w-12 text-right">
                  {{ ((selectedGroupType === 'asset' ? item.current_value : item.remaining_balance) / selectedGroup.total_value * 100).toFixed(1) }}%
                </span>
              </div>
            </li>
          </ul>
        </div>
        <div class="p-4 border-t border-slate-700/50 bg-slate-800/30 text-center">
          <button @click="showGroupDetailModal = false" class="text-sm text-slate-400 hover:text-white transition-colors">Đóng</button>
        </div>
      </div>
    </div>

    <!-- MODAL THÊM/SỬA TÀI SẢN -->
    <div v-if="showAddAssetModal" class="fixed inset-0 bg-slate-900/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="premium-card w-full max-w-md p-6 relative max-h-[90vh] overflow-y-auto">
        <button @click="showAddAssetModal = false" class="absolute top-4 right-4 text-slate-400 hover:text-white">✕</button>
        <h2 class="text-2xl font-bold text-white mb-6">{{ isEditingAsset ? 'Cập Nhật Tài Sản' : 'Thêm Tài Sản Mới' }}</h2>
        
        <form @submit.prevent="submitAssetForm" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Loại Tài Sản</label>
            <CustomSelect v-model="assetForm.category" :options="categoryOptions" :disabled="isEditingAsset" />
          </div>
          
          <div v-if="showNameField">
            <label class="block text-sm font-medium text-slate-300 mb-1">Tên Tài Sản</label>
            <input type="text" v-model="assetForm.name" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
          </div>

          <div v-if="assetForm.category === 'deposit'" class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Ngân hàng</label>
              <CustomSelect v-model="depositBank" :options="bankOptions" />
              <input v-if="depositBank === 'Khác'" type="text" v-model="customDepositBank" placeholder="Nhập tên NH..." class="mt-2 w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Kỳ hạn</label>
              <CustomSelect v-model="depositTerm" :options="termOptions" />
              <input v-if="depositTerm === 'Khác'" type="text" v-model="customDepositTerm" placeholder="Nhập kỳ hạn..." class="mt-2 w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Lãi suất (%/năm)</label>
              <input type="number" step="0.1" v-model="depositInterestRate" placeholder="Ví dụ: 5.5" class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
            </div>
          </div>

          <div v-if="assetForm.category === 'gold'" class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Loại vàng</label>
              <CustomSelect v-model="goldType" :options="goldTypeOptions" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Thương hiệu</label>
              <CustomSelect v-model="goldBrand" :options="goldBrandOptions" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Độ tinh khiết</label>
              <CustomSelect v-model="goldPurity" :options="goldPurityOptions" />
              <input v-if="goldPurity === 'Khác'" type="text" v-model="customGoldPurity" placeholder="Nhập loại..." class="mt-2 w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
            </div>
          </div>

          <div v-if="assetForm.category === 'fund'" class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Mã Quỹ (Ticker)</label>
              <CustomSelect v-model="fundCode" :options="fundCodeOptions" />
              <input v-if="fundCode === 'Khác'" type="text" v-model="customFundCode" placeholder="Nhập mã quỹ..." class="mt-2 w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500 uppercase" />
            </div>
          </div>

          <div v-if="isFluctuatingAsset(assetForm.category)" class="space-y-4">
            <div v-if="showTickerField" class="flex flex-col justify-end">
              <label class="block text-sm font-medium text-slate-300 mb-1">Mã (Ticker)</label>
              <input type="text" v-model="assetForm.ticker" :placeholder="tickerPlaceholder" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 h-[42px] text-white focus:outline-none focus:border-indigo-500" />
            </div>
            <div class="flex flex-col justify-end">
              <label class="block text-sm font-medium text-slate-300 mb-1">Số lượng <span v-if="assetForm.category === 'gold'" class="text-slate-500 text-xs font-normal">(Nhập số lượng theo Đơn vị tính)</span></label>
              <div class="flex gap-2 items-center">
                <CurrencyInput v-model="assetForm.quantity" @update:modelValue="calculateCurrentValue()" class="flex-1 min-w-0" />
                <CustomSelect v-if="assetForm.category === 'gold'" v-model="goldUnit" :options="goldUnitOptions" @update:modelValue="calculateCurrentValue()" class="shrink-0 w-[140px]" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div class="flex flex-col justify-end">
                <label class="block text-sm font-medium text-slate-300 mb-1 flex-1">Giá Vốn <span v-if="assetForm.category === 'gold'" class="text-slate-500 text-xs font-normal">(Nhập theo giá Lượng, tự quy đổi)</span><span v-else class="text-slate-500 text-xs font-normal">(Tùy chọn)</span></label>
                <CurrencyInput v-model="assetForm.avg_price" placeholder="Bỏ qua" />
              </div>
              <div class="flex flex-col justify-end">
                <label class="block text-sm font-medium text-slate-300 mb-1 flex-1">Giá Hiện Tại <span v-if="assetForm.category === 'gold'" class="text-slate-500 text-xs font-normal">(Nhập theo giá Lượng, tự quy đổi)</span></label>
                <CurrencyInput v-model="assetForm.current_price" @update:modelValue="calculateCurrentValue()" />
              </div>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Tổng Giá Trị Hiện Tại</label>
            <CurrencyInput v-model="assetForm.current_value" required class="!text-emerald-400 font-bold text-lg" :readonly="isFluctuatingAsset(assetForm.category)" :class="{ 'opacity-50 cursor-not-allowed': isFluctuatingAsset(assetForm.category) }" />
          </div>

          <div v-if="assetError" class="text-red-400 text-sm">{{ assetError }}</div>
          
          <div class="pt-4 flex justify-end gap-3">
            <button type="button" @click="showAddAssetModal = false" class="px-4 py-2 text-slate-300 hover:text-white transition-colors">Hủy</button>
            <button type="submit" class="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-500/20" :disabled="submitting">
              {{ submitting ? 'Đang lưu...' : (isEditingAsset ? 'Cập Nhật' : 'Lưu Tài Sản') }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL THÊM/SỬA KHOẢN NỢ -->
    <div v-if="showAddLiabilityModal" class="fixed inset-0 bg-slate-900/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="premium-card w-full max-w-md p-6 relative max-h-[90vh] overflow-y-auto">
        <button @click="showAddLiabilityModal = false" class="absolute top-4 right-4 text-slate-400 hover:text-white">✕</button>
        <h2 class="text-2xl font-bold text-white mb-6">{{ isEditingLiability ? 'Cập Nhật Khoản Nợ' : 'Thêm Khoản Nợ' }}</h2>
        
        <form @submit.prevent="submitLiabilityForm" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Loại Nợ</label>
            <CustomSelect v-model="liabilityForm.category" :options="liabilityCategoryOptions" :disabled="isEditingLiability" />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Tên Khoản Nợ</label>
            <input type="text" v-model="liabilityForm.name" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">Dư Nợ Còn Lại</label>
            <CurrencyInput v-model="liabilityForm.remaining_balance" required class="!text-amber-500 font-bold text-lg" />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Lãi Suất (%)</label>
              <CurrencyInput v-model="liabilityForm.interest_rate_percent" placeholder="VD: 7.5" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Trả Góp Hàng Tháng</label>
              <CurrencyInput v-model="liabilityForm.monthly_payment" />
            </div>
          </div>

          <div v-if="liabilityError" class="text-red-400 text-sm">{{ liabilityError }}</div>
          
          <div class="pt-4 flex justify-end gap-3">
            <button type="button" @click="showAddLiabilityModal = false" class="px-4 py-2 text-slate-300 hover:text-white transition-colors">Hủy</button>
            <button type="submit" class="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-500/20" :disabled="submitting">
              {{ submitting ? 'Đang lưu...' : (isEditingLiability ? 'Cập Nhật' : 'Lưu Khoản Nợ') }}
            </button>
          </div>
        </form>
      </div>
    </div>
    <!-- Confirm Delete Modal -->
    <ConfirmModal 
      :show="deleteConfirm.show" 
      :title="deleteConfirm.type === 'asset' ? 'Xóa Tài Sản' : 'Xóa Khoản Nợ'" 
      :message="'Bạn có chắc chắn muốn xóa mục này không? Thao tác này không thể hoàn tác.'" 
      @confirm="executeDelete" 
      @cancel="deleteConfirm.show = false" 
    />

    <!-- Modals -->
    <BulkUpdateGoldModal 
      :show="showBulkUpdateGoldModal" 
      :gold-assets="goldAssetsToUpdate" 
      @close="showBulkUpdateGoldModal = false" 
      @updated="handleBulkUpdateGoldSuccess" 
    />

    <!-- Sync Result Modal -->
    <div v-if="showSyncResultModal" class="fixed inset-0 bg-slate-900/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="premium-card w-full max-w-sm p-6 relative text-center">
        <div v-if="syncSuccess" class="w-16 h-16 bg-emerald-500/20 text-emerald-400 rounded-full flex items-center justify-center mx-auto mb-4 text-3xl">
          ✓
        </div>
        <div v-else class="w-16 h-16 bg-red-500/20 text-red-400 rounded-full flex items-center justify-center mx-auto mb-4 text-3xl">
          ✕
        </div>
        
        <h2 class="text-xl font-bold text-white mb-2">{{ syncSuccess ? 'Đồng bộ thành công!' : 'Đồng bộ thất bại' }}</h2>
        <p class="text-slate-400 text-sm mb-6">{{ syncMessage }}</p>
        
        <button @click="showSyncResultModal = false" class="w-full px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors font-medium">
          Đóng
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { getCategoryName, getCategoryIcon, getLiabilityName, getLiabilityIcon } from '../utils/assetUtils'
import { useWealthStore } from '../stores/wealthStore'
import { getGoldMultiplier, extractGoldUnit } from '../utils/goldCalculations'
import BulkUpdateGoldModal from '../components/BulkUpdateGoldModal.vue'
import CurrencyInput from '../components/common/CurrencyInput.vue'
import ConfirmModal from '../components/common/ConfirmModal.vue'
import CustomSelect from '../components/common/CustomSelect.vue'

const categoryOptions = [
  { label: 'Tiền mặt', value: 'cash' },
  { label: 'Tiền gửi ngân hàng', value: 'deposit' },
  { label: 'Vàng', value: 'gold' },
  { label: 'Cổ phiếu', value: 'stock' },
  { label: 'Chứng chỉ quỹ', value: 'fund' },
  { label: 'Tiền điện tử (Crypto)', value: 'crypto' },
  { label: 'Bất động sản', value: 'real_estate' }
]

const bankOptions = [
  'Vietcombank', 'BIDV', 'VietinBank', 'Agribank', 'Techcombank', 'MBBank', 'VPBank', 'ACB', 'TPBank', 'Sacombank', 'VIB', 'HDBank', 'MSB', 'SHB', 'SeABank', 'OCB', 'Eximbank', 'LPBank', 'Nam A Bank', 'Khác'
].map(name => ({ label: name, value: name }))

const termOptions = [
  'Không kỳ hạn', '1 tháng', '3 tháng', '6 tháng', '12 tháng', '13 tháng', '18 tháng', '24 tháng', 'Khác'
].map(t => ({ label: t, value: t }))

const goldTypeOptions = [
  { label: 'Vàng miếng', value: 'Vàng miếng' },
  { label: 'Vàng nhẫn tròn trơn', value: 'Vàng nhẫn' },
  { label: 'Vàng trang sức', value: 'Vàng trang sức' },
  { label: 'Khác', value: 'Khác' }
]

const goldBrandOptions = [
  'SJC', 'DOJI', 'PNJ', 'Phú Quý', 'Mi Hồng', 'Bảo Tín Minh Châu', 'Bảo Tín Mạnh Hải', 'Tiệm vàng tư nhân', 'Khác'
].map(b => ({ label: b, value: b === 'Tiệm vàng tư nhân' ? 'Tư nhân' : b }))

const goldPurityOptions = [
  { label: '9999 (24K)', value: '9999' },
  { label: '999', value: '999' },
  { label: '99', value: '99' },
  { label: '18K', value: '18K' },
  { label: '14K', value: '14K' },
  { label: '10K', value: '10K' },
  { label: 'Khác', value: 'Khác' }
]

const fundCodeOptions = [
  { label: 'VESAF - VinaCapital', value: 'VESAF' },
  { label: 'VEOF - VinaCapital', value: 'VEOF' },
  { label: 'VIBF - VinaCapital', value: 'VIBF' },
  { label: 'VLBF - VinaCapital', value: 'VLBF' },
  { label: 'SSISCA - SSIAM', value: 'SSISCA' },
  { label: 'VLGF - SSIAM', value: 'VLGF' },
  { label: 'SSIPDF - SSIAM', value: 'SSIPDF' },
  { label: 'VFMVF1 - Dragon Capital', value: 'VFMVF1' },
  { label: 'VFMVF4 - Dragon Capital', value: 'VFMVF4' },
  { label: 'VNDAF - VNDirect', value: 'VNDAF' },
  { label: 'MAFEQI - Manulife', value: 'MAFEQI' },
  { label: 'Khác...', value: 'Khác' }
]

const goldUnitOptions = [
  { label: 'Loại 1 Lượng', value: 'Loại 1 Lượng' },
  { label: 'Loại 5 Chỉ', value: 'Loại 5 Chỉ' },
  { label: 'Loại 2 Chỉ', value: 'Loại 2 Chỉ' },
  { label: 'Loại 1 Chỉ', value: 'Loại 1 Chỉ' },
  { label: 'Loại 0.5 Chỉ', value: 'Loại 0.5 Chỉ' },
  { label: 'Tùy chỉnh (Lượng)', value: 'Lượng' },
  { label: 'Tùy chỉnh (Chỉ)', value: 'Chỉ' },
  { label: 'Tùy chỉnh (Phân)', value: 'Phân' }
]

const liabilityCategoryOptions = [
  { label: 'Vay mua nhà', value: 'mortgage' },
  { label: 'Vay mua xe', value: 'auto_loan' },
  { label: 'Vay học tập', value: 'student_loan' },
  { label: 'Thẻ tín dụng', value: 'credit_card' },
  { label: 'Vay tín chấp', value: 'personal_loan' },
  { label: 'Khác', value: 'other' }
]

const assetCategoryFilterOptions = [
  { label: 'Tất cả các loại', value: '' },
  { label: 'Tiền mặt', value: 'cash' },
  { label: 'Tiền gửi ngân hàng', value: 'deposit' },
  { label: 'Vàng', value: 'gold' },
  { label: 'Cổ phiếu', value: 'stock' },
  { label: 'Chứng chỉ quỹ', value: 'fund' },
  { label: 'Tiền điện tử', value: 'crypto' },
  { label: 'Bất động sản', value: 'real_estate' }
]

const assetSortOptions = [
  { label: 'Giá trị giảm dần', value: 'value_desc' },
  { label: 'Giá trị tăng dần', value: 'value_asc' },
  { label: 'Tên (A-Z)', value: 'name_asc' },
  { label: 'Tên (Z-A)', value: 'name_desc' }
]

const liabilityCategoryFilterOptions = [
  { label: 'Tất cả các loại', value: '' },
  { label: 'Vay mua nhà', value: 'mortgage' },
  { label: 'Vay mua xe', value: 'auto_loan' },
  { label: 'Vay học tập', value: 'student_loan' },
  { label: 'Thẻ tín dụng', value: 'credit_card' },
  { label: 'Vay tín chấp', value: 'personal_loan' },
  { label: 'Khác', value: 'other' }
]

const liabilitySortOptions = [
  { label: 'Dư nợ giảm dần', value: 'value_desc' },
  { label: 'Dư nợ tăng dần', value: 'value_asc' },
  { label: 'Tên (A-Z)', value: 'name_asc' },
  { label: 'Tên (Z-A)', value: 'name_desc' }
]


const wealthStore = useWealthStore()

const showAddAssetModal = ref(false)
const showAddLiabilityModal = ref(false)
const showGroupDetailModal = ref(false)
const selectedGroup = ref(null)
const selectedGroupType = ref('asset') // 'asset' or 'liability'

const expandedAssetGroups = ref({})
const toggleAssetGroup = (category) => {
  expandedAssetGroups.value[category] = !expandedAssetGroups.value[category]
}

const expandedLiabilityGroups = ref({})
const toggleLiabilityGroup = (category) => {
  expandedLiabilityGroups.value[category] = !expandedLiabilityGroups.value[category]
}

const isEditingAsset = ref(false)
const editingAssetId = ref(null)
const isEditingLiability = ref(false)
const editingLiabilityId = ref(null)

const submitting = ref(false)
const assetError = ref('')
const liabilityError = ref('')

const goldType = ref('Vàng miếng')
const goldBrand = ref('SJC')
const goldPurity = ref('9999')
const customGoldPurity = ref('')
const goldUnit = ref('Lượng')
const depositBank = ref('Vietcombank')
const customDepositBank = ref('')
const depositTerm = ref('6 tháng')
const customDepositTerm = ref('')
const depositInterestRate = ref('')
const fundCode = ref('VESAF')
const customFundCode = ref('')

const assetForm = ref({
  category: 'cash',
  name: '',
  ticker: '',
  quantity: null,
  avg_price: null,
  current_price: null,
  current_value: null
})

const liabilityForm = ref({
  category: 'mortgage',
  name: '',
  remaining_balance: null,
  interest_rate_percent: null,
  monthly_payment: null
})

const syncingPrices = ref(false)
const showSyncResultModal = ref(false)
const syncSuccess = ref(true)
const syncMessage = ref('')

const showBulkUpdateGoldModal = ref(false)
const goldAssetsToUpdate = ref([])

const triggerPriceSync = async () => {
  if (syncingPrices.value) return
  syncingPrices.value = true
  try {
    const res = await wealthStore.syncPrices()
    syncSuccess.value = true
    syncMessage.value = `Đã cập nhật thành công giá mới cho ${res.total_updated} tài sản. Thất bại: ${res.total_failed}. Bỏ qua: ${res.total_skipped}.`
  } catch (err) {
    syncSuccess.value = false
    syncMessage.value = err.toString()
  } finally {
    syncingPrices.value = false
    showSyncResultModal.value = true
  }
}

onMounted(() => {
  wealthStore.fetchAll()
})

const applyAssetFilter = () => {
  wealthStore.fetchAssets(false)
}

const applyLiabilityFilter = () => {
  wealthStore.fetchLiabilities(false)
}

const handleAssetScroll = (e) => {
  const { scrollTop, clientHeight, scrollHeight } = e.target
  if (scrollTop + clientHeight >= scrollHeight - 20) {
    if (wealthStore.assetPage < wealthStore.assetTotalPages && !wealthStore.loading) {
      wealthStore.fetchAssets(true)
    }
  }
}

const handleLiabilityScroll = (e) => {
  const { scrollTop, clientHeight, scrollHeight } = e.target
  if (scrollTop + clientHeight >= scrollHeight - 20) {
    if (wealthStore.liabilityPage < wealthStore.liabilityTotalPages && !wealthStore.loading) {
      wealthStore.fetchLiabilities(true)
    }
  }
}

const showNameField = computed(() => {
  return !['stock', 'crypto', 'gold', 'deposit'].includes(assetForm.value.category)
})

const showTickerField = computed(() => {
  return ['stock', 'crypto'].includes(assetForm.value.category)
})

// --- COMPUTED PROPERTIES FOR UI ---
// ... (omitting irrelevant parts, let's just replace the exact submit logic)

const tickerPlaceholder = computed(() => {
  if (assetForm.value.category === 'crypto') return 'BTC, ETH, USDT...'
  if (assetForm.value.category === 'stock') return 'HPG, VCB, SSI...'
  return 'Nhập mã...'
})

// --- COMPUTED PROPERTIES FOR UI ---
const debtToAssetRatio = computed(() => {
  const totalAssets = wealthStore.netWorthSummary?.total_assets || 0
  const totalLiabilities = wealthStore.netWorthSummary?.total_liabilities || 0
  if (totalAssets === 0) return 0
  return ((totalLiabilities / totalAssets) * 100).toFixed(1)
})

const groupedAssets = computed(() => {
  const groups = {}
  let total = 0
  
  wealthStore.assets.forEach(asset => {
    total += asset.current_value
    if (!groups[asset.category]) {
      groups[asset.category] = {
        category: asset.category,
        categoryName: getCategoryName(asset.category),
        total_value: 0,
        totalGoldVolume: 0,
        items: []
      }
    }
    groups[asset.category].items.push(asset)
    groups[asset.category].total_value += asset.current_value

    if (asset.category === 'gold') {
      const unit = extractGoldUnit(asset.name);
      groups[asset.category].totalGoldVolume += (asset.quantity || 0) * getGoldMultiplier(unit);
    }
  })

  // Avoid div by zero
  if (total === 0) total = 1

  const sortedGroups = Object.values(groups).sort((a, b) => b.total_value - a.total_value)
  
  sortedGroups.forEach(g => {
    g.percentage = ((g.total_value / total) * 100).toFixed(1)
    g.items.forEach(i => {
      i.percentage = ((i.current_value / total) * 100).toFixed(1)
    })
  })

  return sortedGroups
})

const groupedLiabilities = computed(() => {
  const groups = {}
  let total = 0

  wealthStore.liabilities.forEach(liability => {
    total += liability.remaining_balance
    if (!groups[liability.category]) {
      groups[liability.category] = {
        category: liability.category,
        categoryName: getLiabilityName(liability.category),
        total_value: 0,
        items: []
      }
    }
    groups[liability.category].items.push(liability)
    groups[liability.category].total_value += liability.remaining_balance
  })

  if (total === 0) total = 1

  const sortedGroups = Object.values(groups).sort((a, b) => b.total_value - a.total_value)
  
  sortedGroups.forEach(g => {
    g.percentage = ((g.total_value / total) * 100).toFixed(1)
    g.items.forEach(i => {
      i.percentage = ((i.remaining_balance / total) * 100).toFixed(1)
    })
  })

  return sortedGroups
})

// --- FORMATTING LOGIC ---
const getPriceClass = (current, avg) => {
  if (!avg || avg === 0) return 'text-slate-300'
  if (current > avg) return 'text-emerald-400'
  if (current < avg) return 'text-red-400'
  return 'text-slate-300'
}

const formatCurrency = (value) => {
  if (value === undefined || value === null) return '0'
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: wealthStore.netWorthSummary?.base_currency || 'VND' }).format(Math.round(value))
}

const formatInputNumber = (value) => {
  if (value === null || value === undefined || value === '') return ''
  return new Intl.NumberFormat('en-US', { maximumFractionDigits: 6 }).format(value)
}

// --- FORM LOGIC ---
const isFluctuatingAsset = (category) => {
  return ['stock', 'fund', 'crypto', 'gold'].includes(category)
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

const calculateCurrentValue = () => {
  if (assetForm.value.quantity && assetForm.value.current_price) {
    if (assetForm.value.category === 'gold') {
      const multiplier = getGoldMultiplier(goldUnit.value)
      assetForm.value.current_value = Math.round(assetForm.value.quantity * (assetForm.value.current_price * multiplier))
    } else {
      assetForm.value.current_value = Math.round(assetForm.value.quantity * assetForm.value.current_price)
    }
  }
}

const openAddAsset = () => {
  isEditingAsset.value = false
  editingAssetId.value = null
  goldUnit.value = 'Lượng'
  depositBank.value = 'Vietcombank'
  customDepositBank.value = ''
  depositTerm.value = '6 tháng'
  customDepositTerm.value = ''
  depositInterestRate.value = ''
  fundCode.value = 'VESAF'
  customFundCode.value = ''
  assetForm.value = { category: 'cash', name: '', ticker: '', quantity: null, avg_price: null, current_price: null, current_value: null }
  showAddAssetModal.value = true
}

const openEditAsset = (asset) => {
  isEditingAsset.value = true
  editingAssetId.value = asset.id
  goldUnit.value = 'Lượng'
  assetForm.value = { ...asset }
  if (assetForm.value.current_value) {
    assetForm.value.current_value = Math.round(assetForm.value.current_value)
  }
  
  if (asset.category === 'gold') {
    const parts = asset.name.split(' - ')
    if (parts.length >= 2) {
      goldBrand.value = parts[0].trim()
      const rest = parts.slice(1).join(' - ').trim()
      
      const unitMatch = rest.match(/\((.+?)\)$/)
      let purityStr = rest
      if (unitMatch) {
        goldUnit.value = unitMatch[1]
        purityStr = rest.replace(unitMatch[0], '').trim()
      } else {
        goldUnit.value = 'Lượng'
      }

      const multiplier = getGoldMultiplier(goldUnit.value)
      
      if (assetForm.value.current_price) assetForm.value.current_price = assetForm.value.current_price / multiplier
      if (assetForm.value.avg_price) assetForm.value.avg_price = assetForm.value.avg_price / multiplier
      
      if (purityStr.includes('9999')) { goldPurity.value = '9999'; goldType.value = purityStr.replace('9999', '').trim() }
      else if (purityStr.includes('999')) { goldPurity.value = '999'; goldType.value = purityStr.replace('999', '').trim() }
      else if (purityStr.includes('99')) { goldPurity.value = '99'; goldType.value = purityStr.replace('99', '').trim() }
      else if (purityStr.includes('18K')) { goldPurity.value = '18K'; goldType.value = purityStr.replace('18K', '').trim() }
      else if (purityStr.includes('14K')) { goldPurity.value = '14K'; goldType.value = purityStr.replace('14K', '').trim() }
      else if (purityStr.includes('10K')) { goldPurity.value = '10K'; goldType.value = purityStr.replace('10K', '').trim() }
      else { goldPurity.value = 'Khác'; customGoldPurity.value = ''; goldType.value = purityStr }
    }
  } else if (asset.category === 'deposit') {
    const match = asset.name.match(/(.+) - Kỳ hạn (.+) - Lãi ([\d.]+)%$/)
    if (match) {
      const bank = match[1].trim()
      const term = match[2].trim()
      
      const banks = ['Vietcombank', 'BIDV', 'VietinBank', 'Agribank', 'Techcombank', 'MBBank', 'VPBank', 'ACB', 'TPBank', 'Sacombank', 'VIB', 'HDBank', 'MSB', 'SHB', 'SeABank', 'OCB', 'Eximbank', 'LPBank', 'Nam A Bank']
      if (banks.includes(bank)) depositBank.value = bank; else { depositBank.value = 'Khác'; customDepositBank.value = bank }
      
      const terms = ['Không kỳ hạn', '1 tháng', '3 tháng', '6 tháng', '12 tháng', '13 tháng', '18 tháng', '24 tháng']
      if (terms.includes(term)) depositTerm.value = term; else { depositTerm.value = 'Khác'; customDepositTerm.value = term }
      
      depositInterestRate.value = match[3]
    } else {
      const oldMatch = asset.name.match(/(.+) - ([\d.]+)%$/)
      if (oldMatch) {
         depositBank.value = 'Khác'; customDepositBank.value = oldMatch[1].trim()
         depositInterestRate.value = oldMatch[2]
      } else {
         depositBank.value = 'Khác'; customDepositBank.value = asset.name
         depositInterestRate.value = ''
      }
    }
  } else if (asset.category === 'fund') {
    if (asset.ticker) {
      const knownFunds = ['VESAF', 'VEOF', 'VIBF', 'VLBF', 'SSISCA', 'VLGF', 'SSIPDF', 'VFMVF1', 'VFMVF4', 'VNDAF', 'MAFEQI']
      if (knownFunds.includes(asset.ticker)) {
        fundCode.value = asset.ticker
      } else {
        fundCode.value = 'Khác'
        customFundCode.value = asset.ticker
      }
    }
  }
  showAddAssetModal.value = true
}

// Bulk Update Gold
const openBulkUpdateGold = (assets) => {
  goldAssetsToUpdate.value = assets
  showBulkUpdateGoldModal.value = true
}

const handleBulkUpdateGoldSuccess = () => {
  showBulkUpdateGoldModal.value = false
  wealthStore.fetchAssets()
}

const submitAssetForm = async () => {
  submitting.value = true
  assetError.value = ''
  try {
    const payload = {
      ...assetForm.value,
      quantity: assetForm.value.quantity ? Number(assetForm.value.quantity) : undefined,
      avg_price: assetForm.value.avg_price ? Number(assetForm.value.avg_price) : undefined,
      current_price: assetForm.value.current_price ? Number(assetForm.value.current_price) : undefined,
      current_value: Math.round(Number(assetForm.value.current_value))
    }
    
    if (payload.category === 'stock' || payload.category === 'crypto') {
      payload.name = payload.ticker.toUpperCase()
      payload.ticker = payload.ticker.toUpperCase()
    } else if (payload.category === 'gold') {
      const purity = goldPurity.value === 'Khác' ? customGoldPurity.value : goldPurity.value
      payload.name = `${goldBrand.value} - ${goldType.value} ${purity} (${goldUnit.value})`.trim()
      
      // Auto-generate ticker for gold sync
      let typeCode = "NHAN"
      if (goldType.value === "Vàng miếng") typeCode = "MIENG"
      else if (goldType.value === "Vàng trang sức") typeCode = "TS"
      
      let brandCode = "KHAC"
      if (goldBrand.value === "Bảo Tín Minh Châu") brandCode = "BTMC"
      else if (goldBrand.value === "Bảo Tín Mạnh Hải") brandCode = "BTMH"
      else if (goldBrand.value === "Phú Quý") brandCode = "PHUQUY"
      else if (goldBrand.value === "Mi Hồng") brandCode = "MIHONG"
      else if (goldBrand.value !== "Tư nhân" && goldBrand.value !== "Khác") brandCode = goldBrand.value.toUpperCase().replace(/\s+/g, '')
      
      payload.ticker = `GOLD_${brandCode}_${typeCode}`
      if (payload.ticker.length > 20) {
        payload.ticker = payload.ticker.substring(0, 20) // Ensure it fits the db
      }
      
      const multiplier = getGoldMultiplier(goldUnit.value)
      
      if (payload.current_price) payload.current_price = payload.current_price * multiplier
      if (payload.avg_price) payload.avg_price = payload.avg_price * multiplier
    } else if (payload.category === 'fund') {
      const fc = fundCode.value === 'Khác' ? customFundCode.value : fundCode.value
      payload.ticker = fc.toUpperCase()
      payload.name = payload.ticker + " - Chứng chỉ quỹ"
    } else if (payload.category === 'deposit') {
      const bank = depositBank.value === 'Khác' ? customDepositBank.value : depositBank.value
      const term = depositTerm.value === 'Khác' ? customDepositTerm.value : depositTerm.value
      
      if (depositInterestRate.value) {
        payload.name = `${bank} - Kỳ hạn ${term} - Lãi ${depositInterestRate.value}%`
      } else {
        payload.name = `${bank} - Kỳ hạn ${term}`
      }
      payload.ticker = undefined
      payload.quantity = undefined
      payload.avg_price = undefined
      payload.current_price = undefined
    } else {
      payload.ticker = undefined
      payload.quantity = undefined
      payload.avg_price = undefined
      payload.current_price = undefined
    }

    if (isEditingAsset.value) {
      await wealthStore.updateAsset(editingAssetId.value, payload)
    } else {
      await wealthStore.createAsset(payload)
    }
    
    showAddAssetModal.value = false
  } catch (err) {
    assetError.value = err
  } finally {
    submitting.value = false
  }
}

const openAddLiability = () => {
  isEditingLiability.value = false
  editingLiabilityId.value = null
  liabilityForm.value = { category: 'mortgage', name: '', remaining_balance: null, interest_rate_percent: null, monthly_payment: null }
  showAddLiabilityModal.value = true
}

const openEditLiability = (liability) => {
  isEditingLiability.value = true
  editingLiabilityId.value = liability.id
  liabilityForm.value = { 
    ...liability,
    interest_rate_percent: liability.interest_rate ? liability.interest_rate * 100 : null
  }
  showAddLiabilityModal.value = true
}

const submitLiabilityForm = async () => {
  submitting.value = true
  liabilityError.value = ''
  try {
    const payload = {
      ...liabilityForm.value,
      remaining_balance: Number(liabilityForm.value.remaining_balance),
      interest_rate: liabilityForm.value.interest_rate_percent ? Number(liabilityForm.value.interest_rate_percent) / 100 : undefined,
      monthly_payment: liabilityForm.value.monthly_payment ? Number(liabilityForm.value.monthly_payment) : 0
    }
    
    if (isEditingLiability.value) {
      await wealthStore.updateLiability(editingLiabilityId.value, payload)
    } else {
      await wealthStore.createLiability(payload)
    }
    showAddLiabilityModal.value = false
  } catch (err) {
    liabilityError.value = err
  } finally {
    submitting.value = false
  }
}

const deleteConfirm = ref({ show: false, id: null, type: '' })

const deleteAsset = (id) => {
  deleteConfirm.value = { show: true, id, type: 'asset' }
}

const deleteLiability = (id) => {
  deleteConfirm.value = { show: true, id, type: 'liability' }
}

const executeDelete = async () => {
  if (deleteConfirm.value.id) {
    if (deleteConfirm.value.type === 'asset') {
      await wealthStore.deleteAsset(deleteConfirm.value.id)
    } else {
      await wealthStore.deleteLiability(deleteConfirm.value.id)
    }
    deleteConfirm.value.show = false
  }
}

const openGroupDetail = (group, type) => {
  selectedGroup.value = group
  selectedGroupType.value = type
  showGroupDetailModal.value = true
}




</script>

<style scoped>
/* Tweak scrollbar for virtual lists */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: rgba(30, 41, 59, 0.5); 
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(99, 102, 241, 0.3); 
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(99, 102, 241, 0.5); 
}
</style>
