<template>
  <div class="space-y-8 animate-fade-in">
    <!-- Header: Liquid Wealth Summary -->
    <div class="relative overflow-hidden rounded-3xl premium-card p-8 border-l-4 border-indigo-500">
      <div class="absolute -right-20 -top-20 w-64 h-64 bg-indigo-500/20 blur-3xl rounded-full pointer-events-none"></div>
      <div class="relative z-10">
        <h1 class="text-3xl font-bold font-outfit text-white mb-2 flex items-center gap-3">
          <span>🌊</span> Thác Nước Tích Sản (Waterfall)
          <button @click="showInfoModal = true" class="p-1.5 rounded-full text-indigo-300 hover:text-white hover:bg-indigo-500/30 transition-colors shadow-lg" title="Cách hoạt động">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
          </button>
        </h1>
        <p class="text-slate-400 max-w-2xl text-lg">
          Nguồn lực thanh khoản của bạn sẽ tự động "rót" vào các mục tiêu theo thứ tự ưu tiên. Hãy sắp xếp các lọ (buckets) bên dưới.
        </p>

        <div class="mt-6 flex flex-wrap gap-6">
          <div class="bg-slate-800/50 rounded-2xl p-5 border border-slate-700/50 min-w-[250px]">
            <p class="text-slate-400 text-sm mb-1">Tổng tài sản thanh khoản</p>
            <div class="text-3xl font-bold text-emerald-400">
              {{ formatCurrency(liquidWealth) }}
            </div>
            <p class="text-xs text-slate-500 mt-2">Bao gồm: Tiền mặt, Cổ phiếu, Vàng, Quỹ</p>
          </div>
          
          <div class="bg-slate-800/50 rounded-2xl p-5 border border-slate-700/50 min-w-[250px]">
            <p class="text-slate-400 text-sm mb-1">Nguồn lực còn lại chưa phân bổ</p>
            <div class="text-3xl font-bold text-indigo-400">
              {{ formatCurrency(remainingWealth) }}
            </div>
            <p class="text-xs text-slate-500 mt-2">Sẵn sàng cho mục tiêu tiếp theo</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Actions -->
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold text-white font-outfit">Danh sách Mục tiêu</h2>
      <button @click="openAddModal" class="px-4 py-2 bg-indigo-500 hover:bg-indigo-600 text-white rounded-lg font-medium transition-colors disabled:opacity-50 flex items-center gap-2">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        Thêm mục tiêu mới
      </button>
    </div>

    <!-- Waterfall Goals List -->
    <div v-if="loading" class="flex justify-center p-12">
      <div class="w-12 h-12 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
    </div>
    
    <div v-else-if="processedGoals.length === 0" class="glass-card p-12 text-center rounded-3xl border-dashed border-2 border-slate-700">
      <div class="text-5xl mb-4">🎯</div>
      <h3 class="text-xl font-bold text-white mb-2">Chưa có mục tiêu nào</h3>
      <p class="text-slate-400 mb-6">Hãy thiết lập các chiếc lọ tài chính để quản lý dòng tiền hiệu quả hơn.</p>
      <button @click="openAddModal" class="px-4 py-2 bg-indigo-500 hover:bg-indigo-600 text-white rounded-lg font-medium transition-colors disabled:opacity-50">Tạo mục tiêu đầu tiên</button>
    </div>

    <div v-else class="space-y-4" @dragover.prevent @drop="onDrop">
      <div 
        v-for="(goal, index) in processedGoals" 
        :key="goal.id"
        draggable="true"
        @dragstart="onDragStart($event, index)"
        @dragenter.prevent="onDragEnter(index)"
        @dragend="onDragEnd"
        class="glass-card rounded-2xl overflow-hidden cursor-move transition-transform duration-200 group border"
        :class="[
          draggedOverIndex === index ? 'border-indigo-500 scale-[1.02] shadow-xl shadow-indigo-500/20' : 'border-slate-800/50 hover:border-slate-700',
          goal.isFilled ? 'border-l-4 border-l-emerald-500' : (goal.allocated > 0 ? 'border-l-4 border-l-indigo-500' : 'border-l-4 border-l-slate-700')
        ]"
      >
        <div class="p-6">
          <div class="flex justify-between items-start mb-4">
            <div class="flex items-center gap-4">
              <div class="text-slate-500 cursor-grab active:cursor-grabbing p-1">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h16M4 16h16"></path></svg>
              </div>
              <div>
                <div class="flex items-center gap-2">
                  <span class="bg-slate-800 text-slate-300 text-xs px-2 py-1 rounded-md font-mono">Ưu tiên {{ index + 1 }}</span>
                  <h3 class="text-xl font-bold text-white font-outfit">{{ goal.name }}</h3>
                </div>
                <div class="text-slate-400 text-sm mt-1" v-if="goal.target_date">
                  Mục tiêu: {{ new Date(goal.target_date).toLocaleDateString('vi-VN') }}
                </div>
              </div>
            </div>
            
            <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <button @click="openLinkModal(goal)" class="p-2 text-slate-400 hover:text-emerald-400 hover:bg-slate-800 rounded-lg transition-colors" title="Gắn tài sản">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path></svg>
              </button>
              <button @click="openEditModal(goal)" class="p-2 text-slate-400 hover:text-indigo-400 hover:bg-slate-800 rounded-lg transition-colors" title="Chỉnh sửa">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
              </button>
              <button @click="confirmDelete(goal)" class="p-2 text-slate-400 hover:text-rose-400 hover:bg-slate-800 rounded-lg transition-colors" title="Xóa">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
              </button>
            </div>
          </div>

          <!-- Progress Bar -->
          <div class="space-y-2">
            <div class="flex justify-between text-sm">
              <span :class="goal.allocated > 0 ? 'text-indigo-400 font-bold' : 'text-slate-500'">
                {{ formatCurrency(goal.allocated) }}
              </span>
              <span class="text-slate-400 font-medium">Mục tiêu: {{ formatCurrency(goal.target_amount) }}</span>
            </div>
            <div class="w-full bg-slate-800 rounded-full h-4 overflow-hidden relative border border-slate-700/50 flex">
              <!-- Dedicated Asset progress -->
              <div 
                v-if="goal.dedicatedAmount > 0"
                class="h-full bg-emerald-500 transition-all duration-1000 ease-out border-r border-emerald-600/50"
                :style="`width: ${goal.dedicatedPercent}%`"
                title="Từ tài sản đã gắn"
              ></div>
              <!-- Waterfall progress -->
              <div 
                v-if="goal.waterfallAmount > 0"
                class="h-full bg-indigo-500 transition-all duration-1000 ease-out"
                :style="`width: ${goal.waterfallPercent}%`"
                title="Từ nguồn thanh khoản chung"
              ></div>
              <!-- Animated stripes for partially filled -->
              <div v-if="!goal.isFilled && goal.allocated > 0" class="absolute inset-0 bg-white/20 progress-stripes"></div>
            </div>
            
            <div class="flex flex-col gap-1 mt-2 text-xs">
              <div class="flex justify-between items-center">
                <span class="text-slate-400">Tiến độ tổng:</span>
                <span class="font-bold whitespace-nowrap" :class="goal.isFilled ? 'text-emerald-400' : 'text-indigo-400'">
                  {{ goal.percentage.toFixed(1) }}%
                </span>
              </div>
              <div class="flex justify-between items-center" v-if="goal.waterfallAmount > 0">
                <span class="text-indigo-400 flex items-center gap-1.5"><div class="w-2 h-2 rounded-full bg-indigo-500"></div> Nguồn thanh khoản chung:</span>
                <span class="text-slate-300">{{ formatCurrency(goal.waterfallAmount) }} đ</span>
              </div>
            </div>
            
            <!-- Linked Assets Section -->
            <div v-if="goal.linkedAssets && goal.linkedAssets.length > 0" class="mt-4 border-t border-slate-700/50 pt-4">
              <h4 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Tài sản đã gắn ({{ formatCurrency(goal.dedicatedAmount) }})</h4>
              <div class="space-y-2">
                <div v-for="asset in goal.linkedAssets" :key="asset.id" class="flex justify-between items-center bg-slate-900/60 p-2.5 rounded-lg border border-slate-700/50 hover:border-indigo-500/30 transition-colors">
                  <div class="flex items-center gap-2">
                    <div class="w-2 h-2 rounded-full bg-emerald-500"></div>
                    <div>
                      <div class="text-sm font-medium text-white truncate max-w-[150px] sm:max-w-[200px]" :title="asset.name">{{ asset.name }}</div>
                      <div class="text-xs text-slate-500 capitalize">{{ asset.category === 'deposit' ? 'Gửi tiết kiệm' : asset.category === 'stock' ? 'Cổ phiếu' : asset.category === 'gold' ? 'Vàng' : asset.category === 'fund' ? 'Chứng chỉ quỹ' : 'Tiền mặt' }}</div>
                    </div>
                  </div>
                  <div class="flex items-center gap-3">
                    <div class="text-sm font-bold text-emerald-400">{{ formatCurrency(asset.current_value) }}</div>
                    <button @click="confirmUnlink(goal, asset.id)" class="p-1 text-slate-500 hover:text-rose-400 hover:bg-rose-500/10 rounded transition-colors" title="Bỏ gắn tài sản">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>

          </div>
          
        </div>
      </div>
    </div>

    <!-- Modal Form -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-fade-in">
      <div class="glass-card premium-card rounded-2xl w-full max-w-md overflow-hidden animate-scale-in border border-slate-700">
        <div class="p-6 border-b border-slate-700/50 bg-slate-800/30">
          <h3 class="text-2xl font-bold text-white font-outfit">{{ isEditing ? 'Cập nhật Mục tiêu' : 'Thêm Mục tiêu mới' }}</h3>
        </div>
        
        <form @submit.prevent="saveGoal" class="p-6 space-y-5 bg-slate-900/50">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Tên mục tiêu</label>
            <input v-model="formData.name" type="text" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" placeholder="VD: Mua nhà, Quỹ khẩn cấp...">
          </div>
          
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Số tiền mục tiêu (VND)</label>
            <CurrencyInput v-model="formData.target_amount" required class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500" />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Ngày dự kiến đạt được (Không bắt buộc)</label>
            <input v-model="formData.target_date" type="date" class="w-full bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500">
          </div>
          
          <div class="pt-4 flex justify-end gap-3 border-t border-slate-700/50">
            <button type="button" @click="closeModal" class="px-4 py-2 text-slate-300 hover:text-white transition-colors">Hủy</button>
            <button type="submit" class="px-4 py-2 bg-indigo-500 hover:bg-indigo-600 text-white rounded-lg font-medium transition-colors disabled:opacity-50" :disabled="saving">
              <span v-if="saving" class="flex items-center gap-2">
                <svg class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
                Đang lưu...
              </span>
              <span v-else>{{ isEditing ? 'Cập nhật' : 'Thêm mới' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Link Asset Modal -->
    <div v-if="showLinkModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-fade-in">
      <div class="glass-card premium-card rounded-2xl w-full max-w-md overflow-hidden animate-scale-in border border-slate-700">
        <div class="p-6 border-b border-slate-700/50 bg-slate-800/30">
          <h3 class="text-2xl font-bold text-white font-outfit">Gắn Tài Sản</h3>
        </div>
        <div class="p-6 space-y-4 bg-slate-900/50 max-h-[60vh] overflow-y-auto">
          <p v-if="availableAssetsGrouped.length === 0" class="text-sm text-amber-500 bg-amber-500/10 p-3 rounded-lg border border-amber-500/20">
            Không còn tài sản thanh khoản nào chưa được gắn. Hãy thêm tài sản mới trong danh mục Tài sản.
          </p>

          <div v-for="group in availableAssetsGrouped" :key="group.label" class="space-y-2 relative">
            <h4 class="text-[11px] font-bold text-slate-400 uppercase tracking-widest sticky top-0 bg-slate-900/95 py-2 z-10 backdrop-blur-sm">{{ group.label }}</h4>
            <div class="space-y-2">
              <button 
                v-for="asset in group.items" 
                :key="asset.id"
                @click="selectedAssetId = asset.id"
                class="w-full text-left p-3.5 rounded-xl border transition-all flex justify-between items-center group focus:outline-none focus:ring-2 focus:ring-indigo-500/50"
                :class="selectedAssetId === asset.id ? 'bg-indigo-500/10 border-indigo-500 shadow-[0_0_15px_rgba(99,102,241,0.15)]' : 'bg-slate-800/40 border-slate-700/50 hover:bg-slate-800 hover:border-slate-600'"
              >
                <div class="flex items-center gap-3 overflow-hidden">
                  <div class="w-2.5 h-2.5 rounded-full flex-shrink-0 transition-all duration-300" :class="selectedAssetId === asset.id ? 'bg-indigo-400 shadow-[0_0_8px_rgba(129,140,248,0.8)] scale-110' : 'bg-slate-600 group-hover:bg-slate-400'"></div>
                  <div class="truncate text-sm font-medium transition-colors" :class="selectedAssetId === asset.id ? 'text-indigo-200' : 'text-slate-200 group-hover:text-white'" :title="asset.name">{{ asset.name }}</div>
                </div>
                <div class="text-sm font-bold pl-3 whitespace-nowrap transition-colors" :class="selectedAssetId === asset.id ? 'text-emerald-300' : 'text-emerald-400/80 group-hover:text-emerald-400'">{{ formatCurrency(asset.current_value) }}</div>
              </button>
            </div>
          </div>
        </div>
        
        <div class="p-5 flex justify-end gap-3 border-t border-slate-700/50 bg-slate-800/30">
          <button type="button" @click="showLinkModal = false" class="px-5 py-2.5 text-sm text-slate-300 hover:text-white transition-colors rounded-xl hover:bg-slate-700/50">Hủy</button>
          <button @click="linkAsset" class="px-6 py-2.5 bg-indigo-500 hover:bg-indigo-600 text-white text-sm rounded-xl font-medium transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-indigo-500/20 hover:shadow-indigo-500/40" :disabled="saving || !selectedAssetId">
            {{ saving ? 'Đang xử lý...' : 'Xác nhận gắn' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Confirm Delete Modal -->
    <ConfirmModal 
      :show="deleteConfirm.show" 
      :title="'Xóa Mục Tiêu'" 
      :message="'Bạn có chắc chắn muốn xóa mục tiêu này không? Tài sản đã gắn sẽ được chuyển lại vào nguồn thanh khoản chung.'" 
      @confirm="executeDelete" 
      @cancel="deleteConfirm.show = false" 
    />

    <!-- Confirm Unlink Modal -->
    <ConfirmModal 
      :show="unlinkConfirm.show" 
      :title="'Bỏ Gắn Tài Sản'" 
      :message="'Bạn có chắc chắn muốn bỏ gắn tài sản này khỏi mục tiêu không? Tài sản sẽ được hoàn trả lại vào nguồn thanh khoản chung.'" 
      @confirm="executeUnlink" 
      @cancel="unlinkConfirm.show = false" 
    />

    <!-- Info Modal -->
    <div v-if="showInfoModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-md animate-fade-in">
      <div class="glass-card premium-card rounded-2xl w-full max-w-2xl overflow-hidden animate-scale-in border border-indigo-500/30 shadow-[0_0_40px_rgba(99,102,241,0.15)]">
        <div class="p-6 border-b border-slate-700/50 bg-slate-800/80 flex justify-between items-center">
          <h3 class="text-2xl font-bold text-white font-outfit flex items-center gap-2">
            <span>🌊</span> Giải thích Thuật toán Hybrid Waterfall
          </h3>
          <button @click="showInfoModal = false" class="text-slate-400 hover:text-white transition-colors bg-slate-800 p-2 rounded-xl hover:bg-slate-700">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>
        <div class="p-8 space-y-8 bg-slate-900/80 max-h-[70vh] overflow-y-auto text-slate-300 leading-relaxed text-sm custom-scrollbar">
          <p class="text-base text-slate-200">
            <strong class="text-white">Life Capital (WealthOS)</strong> không quản lý mục tiêu theo cách thông thường (chia nhỏ thu nhập mỗi tháng). Thay vào đó, hệ thống sử dụng thuật toán <strong class="text-indigo-400">Hybrid Waterfall (Dòng chảy thác nước)</strong> để trả lời câu hỏi: <em>"Với tổng tài sản hiện có, tôi đang đứng ở đâu trên hành trình tài chính?"</em>
          </p>

          <div class="space-y-4">
            <h4 class="text-white font-bold text-base border-b border-slate-700/80 pb-2">1. Nguồn tiền được tính như thế nào?</h4>
            <p>Hệ thống tự động gom toàn bộ <strong>Tài sản thanh khoản</strong> của bạn (Tiền mặt, Gửi tiết kiệm, Cổ phiếu, Vàng, Chứng chỉ quỹ) thành một "Hồ chứa" duy nhất.</p>
          </div>

          <div class="space-y-4">
            <h4 class="text-white font-bold text-base border-b border-slate-700/80 pb-2">2. Tiền chảy vào Mục tiêu ra sao? (Ý nghĩa màu sắc)</h4>
            <ul class="space-y-4 list-none pl-0">
              <li class="flex gap-4 items-start bg-slate-800/40 p-4 rounded-xl border border-slate-700/50">
                <div class="w-4 h-4 rounded-full mt-0.5 bg-emerald-500 shrink-0 shadow-[0_0_10px_rgba(16,185,129,0.5)]"></div>
                <div>
                  <strong class="text-emerald-400 text-base block mb-1">Tài sản đã gắn (Màu Xanh Ngọc)</strong> 
                  Nếu bạn chỉ định một tài sản (ví dụ: Sổ tiết kiệm A) cho một mục tiêu cụ thể, toàn bộ giá trị tài sản đó sẽ được ưu tiên khóa lại chỉ dành riêng cho mục tiêu này.
                </div>
              </li>
              <li class="flex gap-4 items-start bg-slate-800/40 p-4 rounded-xl border border-slate-700/50">
                <div class="w-4 h-4 rounded-full mt-0.5 bg-indigo-500 shrink-0 shadow-[0_0_10px_rgba(99,102,241,0.5)]"></div>
                <div>
                  <strong class="text-indigo-400 text-base block mb-1">Nguồn thanh khoản chung (Màu Xanh Tím)</strong> 
                  Nếu mục tiêu vẫn chưa đạt đủ 100% sau khi đã tính tài sản gắn, hệ thống sẽ tự động "rót" phần tiền còn dư trong Hồ chứa chung vào mục tiêu đó theo thứ tự ưu tiên (Mục tiêu Ưu tiên 1 -> Mục tiêu Ưu tiên 2...).
                </div>
              </li>
            </ul>
          </div>

          <div class="bg-indigo-500/10 border border-indigo-500/30 p-5 rounded-2xl relative overflow-hidden">
            <div class="absolute -right-10 -top-10 w-32 h-32 bg-indigo-500/20 blur-2xl rounded-full pointer-events-none"></div>
            <h4 class="text-indigo-300 font-bold mb-3 flex items-center gap-2"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg> Ví dụ minh họa:</h4>
            <p class="text-slate-200">Bạn đang có tổng <strong class="text-white">1.5 Tỷ VNĐ</strong> nhàn rỗi. Bạn tạo mục tiêu "Mua ô tô" (500 Triệu) và gắn Sổ tiết kiệm 12.5 Triệu vào.</p>
            <ul class="list-disc pl-5 mt-3 space-y-2 text-indigo-200/90">
              <li>Màu Xanh Ngọc sẽ hiển thị: 12.5 Triệu (Tài sản gắn).</li>
              <li>Màu Xanh Tím sẽ tự động rót vào: 487.5 Triệu còn thiếu (lấy từ 1.5 Tỷ thanh khoản chung).</li>
            </ul>
            <p class="mt-4 font-bold text-indigo-300 bg-indigo-900/30 p-3 rounded-lg border border-indigo-500/20">=> Vì tổng tài sản của bạn lớn hơn giá trị ô tô, hệ thống đánh giá bạn dư sức mua ô tô và báo tiến độ đạt 100%!</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useGoalStore } from '../stores/goalStore'
import { useWealthStore } from '../stores/wealthStore'
import { useProfileStore } from '../stores/profileStore'
import CurrencyInput from '../components/common/CurrencyInput.vue'
import ConfirmModal from '../components/common/ConfirmModal.vue'

const goalStore = useGoalStore()
const wealthStore = useWealthStore()
const profileStore = useProfileStore()

const loading = ref(true)
const showModal = ref(false)
const showInfoModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)
const currentGoalId = ref(null)

const formData = ref({
  name: '',
  target_amount: '',
  target_date: ''
})

// Drag and drop state
const draggedIndex = ref(null)
const draggedOverIndex = ref(null)

onMounted(async () => {
  try {
    await Promise.all([
      goalStore.fetchGoals(),
      wealthStore.fetchAllForDashboard(),
      profileStore.fetchProfile()
    ])
  } catch (error) {
    console.error("Failed to load goals data:", error)
  } finally {
    loading.value = false
  }
})

const liquidWealth = computed(() => {
  // Cash, Deposit, Stock, Gold, Fund are liquid. Real Estate is not.
  const liquidCategories = ['cash', 'deposit', 'stock', 'gold', 'fund']
  return wealthStore.allAssets
    .filter(a => liquidCategories.includes(a.category))
    .reduce((sum, a) => sum + a.current_value, 0)
})

const processedGoals = computed(() => {
  const liquidCategories = ['cash', 'deposit', 'stock', 'gold', 'fund']
  
  // First calculate total liquid wealth
  let totalLiquid = wealthStore.allAssets
    .filter(a => liquidCategories.includes(a.category))
    .reduce((sum, a) => sum + a.current_value, 0)
    
  // Find all assets that are allocated to ANY goal
  const allocatedAssetIds = new Set()
  goalStore.goals.forEach(goal => {
    if (goal.allocations) {
      goal.allocations.forEach(alloc => allocatedAssetIds.add(alloc.asset_id))
    }
  })
  
  // Calculate total dedicated wealth (wealth from assets assigned to goals)
  let dedicatedWealth = wealthStore.allAssets
    .filter(a => allocatedAssetIds.has(a.id))
    .reduce((sum, a) => sum + a.current_value, 0)
    
  // Unallocated wealth is the pool for the waterfall
  let unallocatedRemaining = totalLiquid - dedicatedWealth
  if (unallocatedRemaining < 0) unallocatedRemaining = 0
  
  return goalStore.goals.map(goal => {
    // 1. Calculate dedicated amount for this specific goal
    let dedicatedAmount = 0
    let linkedAssets = []
    if (goal.allocations) {
      goal.allocations.forEach(alloc => {
        const asset = wealthStore.allAssets.find(a => a.id === alloc.asset_id)
        if (asset) {
          dedicatedAmount += asset.current_value
          linkedAssets.push(asset)
        }
      })
    }
    
    // 2. Add waterfall amount
    let waterfallAmount = 0
    let deficit = goal.target_amount - dedicatedAmount
    
    if (deficit > 0 && unallocatedRemaining > 0) {
      if (unallocatedRemaining >= deficit) {
        waterfallAmount = deficit
        unallocatedRemaining -= deficit
      } else {
        waterfallAmount = unallocatedRemaining
        unallocatedRemaining = 0
      }
    }
    
    const totalAllocated = dedicatedAmount + waterfallAmount
    let percentage = (totalAllocated / goal.target_amount) * 100
    if (percentage > 100) percentage = 100
    
    let dedicatedPercent = (dedicatedAmount / goal.target_amount) * 100
    if (dedicatedPercent > 100) dedicatedPercent = 100
    
    let waterfallPercent = (waterfallAmount / goal.target_amount) * 100
    if (dedicatedPercent + waterfallPercent > 100) waterfallPercent = 100 - dedicatedPercent
    
    return {
      ...goal,
      dedicatedAmount,
      waterfallAmount,
      dedicatedPercent,
      waterfallPercent,
      allocated: totalAllocated,
      percentage,
      isFilled: totalAllocated >= goal.target_amount,
      linkedAssets
    }
  })
})

const remainingWealth = computed(() => {
  // Get the last unallocated amount remaining after all goals processed
  const lastGoal = processedGoals.value[processedGoals.value.length - 1]
  if (!lastGoal) return liquidWealth.value
  
  let totalAllocatedToGoals = processedGoals.value.reduce((sum, g) => sum + g.allocated, 0)
  let remaining = liquidWealth.value - totalAllocatedToGoals
  return remaining > 0 ? remaining : 0
})

const openAddModal = () => {
  isEditing.value = false
  formData.value = {
    name: '',
    target_amount: '',
    target_date: ''
  }
  showModal.value = true
}

const openEditModal = (goal) => {
  isEditing.value = true
  currentGoalId.value = goal.id
  formData.value = {
    name: goal.name,
    target_amount: goal.target_amount,
    target_date: goal.target_date ? goal.target_date.split('T')[0] : ''
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveGoal = async () => {
  saving.value = true
  try {
    const payload = {
      name: formData.value.name,
      target_amount: Number(formData.value.target_amount),
      target_date: formData.value.target_date ? new Date(formData.value.target_date).toISOString() : null,
      priority: isEditing.value ? goalStore.goals.find(g => g.id === currentGoalId.value).priority : goalStore.goals.length + 1
    }

    if (isEditing.value) {
      await goalStore.updateGoal(currentGoalId.value, payload)
    } else {
      await goalStore.createGoal(payload)
    }
    closeModal()
  } catch (error) {
    alert('Có lỗi xảy ra: ' + error.message)
  } finally {
    saving.value = false
  }
}

const deleteConfirm = ref({
  show: false,
  goal: null
})

const confirmDelete = (goal) => {
  deleteConfirm.value.goal = goal
  deleteConfirm.value.show = true
}

const executeDelete = async () => {
  if (deleteConfirm.value.goal) {
    try {
      await goalStore.deleteGoal(deleteConfirm.value.goal.id)
      const updatedGoals = goalStore.goals.map((g, idx) => ({ ...g, priority: idx + 1 }))
      await goalStore.updatePriorities(updatedGoals)
    } catch (error) {
      alert('Có lỗi xảy ra: ' + error.message)
    } finally {
      deleteConfirm.value.show = false
      deleteConfirm.value.goal = null
    }
  }
}

// Link Asset Logic
const showLinkModal = ref(false)
const linkingGoalId = ref(null)
const selectedAssetId = ref('')

const openLinkModal = (goal) => {
  linkingGoalId.value = goal.id
  selectedAssetId.value = ''
  showLinkModal.value = true
}

const linkAsset = async () => {
  if (!selectedAssetId.value) return
  saving.value = true
  try {
    await goalStore.linkAsset(linkingGoalId.value, selectedAssetId.value)
    showLinkModal.value = false
  } catch (error) {
    alert('Có lỗi xảy ra: ' + error.message)
  } finally {
    saving.value = false
  }
}

const unlinkConfirm = ref({
  show: false,
  goal: null,
  assetId: null
})

const confirmUnlink = (goal, assetId) => {
  unlinkConfirm.value.goal = goal
  unlinkConfirm.value.assetId = assetId
  unlinkConfirm.value.show = true
}

const executeUnlink = async () => {
  if (unlinkConfirm.value.goal && unlinkConfirm.value.assetId) {
    try {
      await goalStore.unlinkAsset(unlinkConfirm.value.goal.id, unlinkConfirm.value.assetId)
    } catch (error) {
      alert('Có lỗi xảy ra: ' + error.message)
    } finally {
      unlinkConfirm.value.show = false
      unlinkConfirm.value.goal = null
      unlinkConfirm.value.assetId = null
    }
  }
}

// Get available assets (not linked to any goal)
const availableAssetsGrouped = computed(() => {
  const liquidCategories = ['cash', 'deposit', 'stock', 'gold', 'fund']
  const allocatedAssetIds = new Set()
  goalStore.goals.forEach(goal => {
    if (goal.allocations) {
      goal.allocations.forEach(alloc => allocatedAssetIds.add(alloc.asset_id))
    }
  })
  
  const available = wealthStore.allAssets.filter(a => 
    liquidCategories.includes(a.category) && !allocatedAssetIds.has(a.id)
  )

  // Sort available by value descending
  available.sort((a, b) => b.current_value - a.current_value)

  // Group them
  const groups = {
    deposit: { label: 'Gửi tiết kiệm', items: [] },
    stock: { label: 'Cổ phiếu', items: [] },
    gold: { label: 'Vàng', items: [] },
    fund: { label: 'Chứng chỉ quỹ', items: [] },
    cash: { label: 'Tiền mặt', items: [] }
  }

  available.forEach(asset => {
    if (groups[asset.category]) {
      groups[asset.category].items.push(asset)
    }
  })

  // Return only non-empty groups
  return Object.values(groups).filter(g => g.items.length > 0)
})

// Drag & Drop handlers
const onDragStart = (e, index) => {
  draggedIndex.value = index
  e.dataTransfer.effectAllowed = 'move'
  // Custom transparent drag image to look better
  const crt = e.target.cloneNode(true)
  crt.style.opacity = '0.5'
  crt.style.position = 'absolute'
  crt.style.top = '-1000px'
  document.body.appendChild(crt)
  e.dataTransfer.setDragImage(crt, 0, 0)
  setTimeout(() => document.body.removeChild(crt), 0)
}

const onDragEnter = (index) => {
  draggedOverIndex.value = index
}

const onDragEnd = () => {
  draggedOverIndex.value = null
}

const onDrop = async () => {
  const from = draggedIndex.value
  const to = draggedOverIndex.value
  
  draggedOverIndex.value = null
  
  if (from !== null && to !== null && from !== to) {
    const newGoalsList = [...goalStore.goals]
    const [movedGoal] = newGoalsList.splice(from, 1)
    newGoalsList.splice(to, 0, movedGoal)
    
    // Update priorities locally and send to server
    const orderedGoals = newGoalsList.map((g, idx) => ({
      ...g,
      priority: idx + 1
    }))
    
    try {
      await goalStore.updatePriorities(orderedGoals)
    } catch (e) {
      alert('Không thể lưu thay đổi thứ tự ưu tiên')
    }
  }
  draggedIndex.value = null
}

const formatCurrency = (value) => {
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(value)
}
</script>

<style scoped>
.progress-stripes {
  background-image: linear-gradient(
    45deg,
    rgba(255, 255, 255, 0.15) 25%,
    transparent 25%,
    transparent 50%,
    rgba(255, 255, 255, 0.15) 50%,
    rgba(255, 255, 255, 0.15) 75%,
    transparent 75%,
    transparent
  );
  background-size: 1rem 1rem;
  animation: progress-stripes 1s linear infinite;
}

input[type="date"]::-webkit-calendar-picker-indicator {
  cursor: pointer;
  filter: invert(1);
  opacity: 0.8;
  transition: opacity 0.2s;
}

input[type="date"]::-webkit-calendar-picker-indicator:hover {
  opacity: 1;
}

@keyframes progress-stripes {
  from {
    background-position: 1rem 0;
  }
  to {
    background-position: 0 0;
  }
}
</style>
