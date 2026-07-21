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
          <select v-model="wealthStore.assetCategoryFilter" @change="applyAssetFilter" class="bg-slate-900/80 border border-slate-600 rounded-lg px-3 py-1.5 text-white text-sm focus:outline-none focus:border-indigo-500">
            <option value="">Tất cả các loại</option>
            <option value="cash">Tiền mặt</option>
            <option value="deposit">Tiền gửi ngân hàng</option>
            <option value="gold">Vàng</option>
            <option value="stock">Cổ phiếu</option>
            <option value="fund">Chứng chỉ quỹ</option>
            <option value="crypto">Tiền điện tử</option>
            <option value="real_estate">Bất động sản</option>
          </select>
        </div>
        <div class="flex items-center gap-3">
          <label class="text-sm font-medium text-slate-300 whitespace-nowrap">Sắp xếp:</label>
          <select v-model="wealthStore.assetSort" @change="applyAssetFilter" class="bg-slate-900/80 border border-slate-600 rounded-lg px-3 py-1.5 text-white text-sm focus:outline-none focus:border-indigo-500">
            <option value="value_desc">Giá trị giảm dần</option>
            <option value="value_asc">Giá trị tăng dần</option>
            <option value="name_asc">Tên (A-Z)</option>
            <option value="name_desc">Tên (Z-A)</option>
          </select>
        </div>
      </div>
      <div class="flex flex-wrap items-center gap-4 w-full sm:w-auto">
        <div class="flex items-center gap-3">
          <label class="text-sm font-medium text-slate-300 whitespace-nowrap">Lọc Nợ:</label>
          <select v-model="wealthStore.liabilityCategoryFilter" @change="applyLiabilityFilter" class="bg-slate-900/80 border border-slate-600 rounded-lg px-3 py-1.5 text-white text-sm focus:outline-none focus:border-indigo-500">
            <option value="">Tất cả các loại</option>
            <option value="mortgage">Vay mua nhà</option>
            <option value="auto_loan">Vay mua xe</option>
            <option value="student_loan">Vay học tập</option>
            <option value="credit_card">Thẻ tín dụng</option>
            <option value="personal_loan">Vay tín chấp</option>
            <option value="other">Khác</option>
          </select>
        </div>
        <div class="flex items-center gap-3">
          <label class="text-sm font-medium text-slate-300 whitespace-nowrap">Sắp xếp:</label>
          <select v-model="wealthStore.liabilitySort" @change="applyLiabilityFilter" class="bg-slate-900/80 border border-slate-600 rounded-lg px-3 py-1.5 text-white text-sm focus:outline-none focus:border-indigo-500">
            <option value="value_desc">Dư nợ giảm dần</option>
            <option value="value_asc">Dư nợ tăng dần</option>
            <option value="name_asc">Tên (A-Z)</option>
            <option value="name_desc">Tên (Z-A)</option>
          </select>
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
          <button @click="openAddAsset" class="px-4 py-2 bg-indigo-500/20 text-indigo-400 rounded-lg hover:bg-indigo-500/30 transition-colors text-sm font-medium border border-indigo-500/30">
            + Thêm Tài sản
          </button>
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
              <div class="bg-slate-800/60 px-4 py-3 flex justify-between items-center sticky top-0 z-10 cursor-pointer hover:bg-slate-700/60 transition-colors" @click="openGroupDetail(group, 'asset')">
                <div class="flex items-center gap-3">
                  <span class="text-xl">{{ getAssetIcon(group.category) }}</span>
                  <h3 class="text-white font-semibold">{{ group.categoryName }}</h3>
                </div>
                <div class="text-right flex items-center gap-3">
                  <div>
                    <p class="text-emerald-400 font-bold">{{ formatCurrency(group.total_value) }}</p>
                    <p class="text-xs text-slate-400">Chiếm {{ group.percentage }}%</p>
                  </div>
                  <span class="text-slate-500 text-xs">▶</span>
                </div>
              </div>
              
              <!-- Danh sách items trong nhóm -->
              <ul class="divide-y divide-slate-700/30">
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
              <div class="bg-slate-800/60 px-4 py-3 flex justify-between items-center sticky top-0 z-10 cursor-pointer hover:bg-slate-700/60 transition-colors" @click="openGroupDetail(group, 'liability')">
                <div class="flex items-center gap-3">
                  <span class="text-xl">{{ getLiabilityIcon(group.category) }}</span>
                  <h3 class="text-white font-semibold">{{ group.categoryName }}</h3>
                </div>
                <div class="text-right flex items-center gap-3">
                  <div>
                    <p class="text-amber-500 font-bold">{{ formatCurrency(group.total_value) }}</p>
                    <p class="text-xs text-slate-400">Chiếm {{ group.percentage }}%</p>
                  </div>
                  <span class="text-slate-500 text-xs">▶</span>
                </div>
              </div>
              
              <ul class="divide-y divide-slate-700/30">
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
              <span>{{ selectedGroupType === 'asset' ? getAssetIcon(selectedGroup.category) : getLiabilityIcon(selectedGroup.category) }}</span>
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
            <select v-model="assetForm.category" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" :disabled="isEditingAsset">
              <option value="cash">Tiền mặt</option>
              <option value="deposit">Tiền gửi ngân hàng</option>
              <option value="gold">Vàng</option>
              <option value="stock">Cổ phiếu</option>
              <option value="fund">Chứng chỉ quỹ</option>
              <option value="crypto">Tiền điện tử (Crypto)</option>
              <option value="real_estate">Bất động sản</option>
            </select>
          </div>
          
          <div v-if="showNameField">
            <label class="block text-sm font-medium text-slate-300 mb-1">Tên Tài Sản</label>
            <input type="text" v-model="assetForm.name" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
          </div>

          <div v-if="assetForm.category === 'deposit'">
            <label class="block text-sm font-medium text-slate-300 mb-1">Lãi suất (%/năm)</label>
            <input type="number" step="0.1" v-model="depositInterestRate" placeholder="Ví dụ: 5.5" class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
          </div>

          <div v-if="assetForm.category === 'gold'" class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Loại vàng</label>
              <select v-model="goldType" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500">
                <option value="Vàng miếng">Vàng miếng</option>
                <option value="Vàng nhẫn">Vàng nhẫn tròn trơn</option>
                <option value="Vàng trang sức">Vàng trang sức</option>
                <option value="Khác">Khác</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Thương hiệu</label>
              <select v-model="goldBrand" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500">
                <option value="SJC">SJC</option>
                <option value="DOJI">DOJI</option>
                <option value="PNJ">PNJ</option>
                <option value="Phú Quý">Phú Quý</option>
                <option value="Mi Hồng">Mi Hồng</option>
                <option value="Bảo Tín Minh Châu">Bảo Tín Minh Châu</option>
                <option value="Bảo Tín Mạnh Hải">Bảo Tín Mạnh Hải</option>
                <option value="Tư nhân">Tiệm vàng tư nhân</option>
                <option value="Khác">Khác</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Độ tinh khiết</label>
              <select v-model="goldPurity" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500">
                <option value="9999">9999 (24K)</option>
                <option value="999">999</option>
                <option value="99">99</option>
                <option value="18K">18K</option>
                <option value="14K">14K</option>
                <option value="10K">10K</option>
                <option value="Khác">Khác</option>
              </select>
              <input v-if="goldPurity === 'Khác'" type="text" v-model="customGoldPurity" placeholder="Nhập loại..." class="mt-2 w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
            </div>
          </div>

          <div v-if="isFluctuatingAsset(assetForm.category)" class="grid grid-cols-2 gap-4">
            <div v-if="showTickerField">
              <label class="block text-sm font-medium text-slate-300 mb-1">Mã (Ticker)</label>
              <input type="text" v-model="assetForm.ticker" :placeholder="tickerPlaceholder" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Số lượng <span v-if="assetForm.category === 'gold'" class="text-slate-500 text-xs font-normal">(Tính theo Lượng)</span></label>
              <div class="flex gap-2">
                <CurrencyInput v-model="assetForm.quantity" @update:modelValue="calculateCurrentValue(); quickGoldWeight = ''" />
                <select v-if="assetForm.category === 'gold'" v-model="quickGoldWeight" @change="applyQuickGoldWeight" class="bg-slate-900/50 border border-slate-700 rounded-lg px-2 py-2 text-white text-sm focus:outline-none focus:border-indigo-500 shrink-0 w-[110px]">
                  <option value="">Tùy chỉnh</option>
                  <option value="1">1 Lượng</option>
                  <option value="0.5">5 Chỉ</option>
                  <option value="0.2">2 Chỉ</option>
                  <option value="0.1">1 Chỉ</option>
                  <option value="0.05">5 Phân</option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Giá Vốn <span class="text-slate-500 text-xs font-normal">(Tùy chọn)</span></label>
              <CurrencyInput v-model="assetForm.avg_price" placeholder="Bỏ qua" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Giá Hiện Tại</label>
              <CurrencyInput v-model="assetForm.current_price" @update:modelValue="calculateCurrentValue()" />
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
            <select v-model="liabilityForm.category" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" :disabled="isEditingLiability">
              <option value="mortgage">Vay mua nhà</option>
              <option value="auto_loan">Vay mua xe</option>
              <option value="student_loan">Vay học tập</option>
              <option value="credit_card">Thẻ tín dụng</option>
              <option value="personal_loan">Vay tín chấp</option>
              <option value="other">Khác</option>
            </select>
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
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useWealthStore } from '../stores/wealthStore'
import CurrencyInput from '../components/common/CurrencyInput.vue'
import ConfirmModal from '../components/common/ConfirmModal.vue'

const wealthStore = useWealthStore()

const showAddAssetModal = ref(false)
const showAddLiabilityModal = ref(false)
const showGroupDetailModal = ref(false)
const selectedGroup = ref(null)
const selectedGroupType = ref('asset') // 'asset' or 'liability'

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
const quickGoldWeight = ref('')
const depositInterestRate = ref('')

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

// removed formattedInputs

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
  return !['stock', 'crypto', 'gold'].includes(assetForm.value.category)
})

const showTickerField = computed(() => {
  return ['stock', 'crypto'].includes(assetForm.value.category)
})

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
        items: []
      }
    }
    groups[asset.category].items.push(asset)
    groups[asset.category].total_value += asset.current_value
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
        categoryName: getLiabilityCategoryName(liability.category),
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
const formatCurrency = (value) => {
  if (value === undefined || value === null) return '0'
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: wealthStore.netWorthSummary?.base_currency || 'VND' }).format(value)
}

const formatInputNumber = (value) => {
  if (value === null || value === undefined || value === '') return ''
  return new Intl.NumberFormat('en-US', { maximumFractionDigits: 6 }).format(value)
}

// --- FORM LOGIC ---
const isFluctuatingAsset = (category) => {
  return ['stock', 'fund', 'crypto', 'gold'].includes(category)
}

const calculateCurrentValue = () => {
  if (assetForm.value.quantity && assetForm.value.current_price) {
    assetForm.value.current_value = assetForm.value.quantity * assetForm.value.current_price
  }
}

const openAddAsset = () => {
  isEditingAsset.value = false
  editingAssetId.value = null
  quickGoldWeight.value = ''
  depositInterestRate.value = ''
  assetForm.value = { category: 'cash', name: '', ticker: '', quantity: null, avg_price: null, current_price: null, current_value: null }
  showAddAssetModal.value = true
}

const openEditAsset = (asset) => {
  isEditingAsset.value = true
  editingAssetId.value = asset.id
  quickGoldWeight.value = ''
  assetForm.value = { ...asset }
  
  if (asset.category === 'gold') {
    const parts = asset.name.split(' - ')
    if (parts.length >= 2) {
      goldBrand.value = parts[0].trim()
      const rest = parts.slice(1).join(' - ').trim()
      if (rest.includes('9999')) { goldPurity.value = '9999'; goldType.value = rest.replace('9999', '').trim() }
      else if (rest.includes('999')) { goldPurity.value = '999'; goldType.value = rest.replace('999', '').trim() }
      else if (rest.includes('99')) { goldPurity.value = '99'; goldType.value = rest.replace('99', '').trim() }
      else if (rest.includes('18K')) { goldPurity.value = '18K'; goldType.value = rest.replace('18K', '').trim() }
      else if (rest.includes('14K')) { goldPurity.value = '14K'; goldType.value = rest.replace('14K', '').trim() }
      else if (rest.includes('10K')) { goldPurity.value = '10K'; goldType.value = rest.replace('10K', '').trim() }
      else { goldPurity.value = 'Khác'; customGoldPurity.value = ''; goldType.value = rest }
    }
  } else if (asset.category === 'deposit') {
    const match = asset.name.match(/(.+) - ([\d.]+)%$/)
    if (match) {
      assetForm.value.name = match[1].trim()
      depositInterestRate.value = match[2]
    } else {
      depositInterestRate.value = ''
    }
  }
  showAddAssetModal.value = true
}

const applyQuickGoldWeight = () => {
  if (quickGoldWeight.value) {
    assetForm.value.quantity = parseFloat(quickGoldWeight.value)
    calculateCurrentValue()
  }
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
      current_value: Number(assetForm.value.current_value)
    }
    
    if (payload.category === 'stock' || payload.category === 'crypto') {
      payload.name = payload.ticker.toUpperCase()
      payload.ticker = payload.ticker.toUpperCase()
    } else if (payload.category === 'gold') {
      const purity = goldPurity.value === 'Khác' ? customGoldPurity.value : goldPurity.value
      payload.name = `${goldBrand.value} - ${goldType.value} ${purity}`.trim()
      payload.ticker = undefined
    } else if (payload.category === 'fund') {
      payload.ticker = undefined
    } else if (payload.category === 'deposit') {
      if (depositInterestRate.value) {
        payload.name = `${payload.name} - ${depositInterestRate.value}%`
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

// --- UTILS ---
const getCategoryName = (category) => {
  const names = {
    cash: 'Tiền mặt', deposit: 'Tiền gửi ngân hàng', gold: 'Vàng', stock: 'Cổ phiếu', fund: 'Chứng chỉ quỹ', crypto: 'Tiền điện tử', real_estate: 'Bất động sản'
  }
  return names[category] || category
}

const getLiabilityCategoryName = (category) => {
  const names = {
    mortgage: 'Vay mua nhà', auto_loan: 'Vay mua xe', student_loan: 'Vay học tập', credit_card: 'Thẻ tín dụng', personal_loan: 'Vay tín chấp', other: 'Khác'
  }
  return names[category] || category
}

const getAssetIcon = (category) => {
  const icons = {
    cash: '💵', deposit: '🏦', gold: '🪙', stock: '📈', fund: '📊', crypto: '₿', real_estate: '🏠'
  }
  return icons[category] || '💰'
}

const getLiabilityIcon = (category) => {
  const icons = {
    mortgage: '🏘️', auto_loan: '🚗', student_loan: '🎓', credit_card: '💳', personal_loan: '👤', other: '🧾'
  }
  return icons[category] || '💸'
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
