<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold font-outfit text-slate-100">Tổng Quan Dashboard</h1>
        <p class="text-xs text-slate-400">Tổng quan tài sản và mục tiêu</p>
      </div>
      <button class="px-4 py-2 text-xs bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl transition duration-150 font-semibold">
        Khởi tạo tháng mới
      </button>
    </div>

    <!-- Basic grid placeholder -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="glass-card rounded-2xl p-6">
        <span class="text-xs text-slate-400 uppercase font-bold tracking-wider">Tài Sản Ròng</span>
        <h2 class="text-3xl font-mono font-bold mt-2 text-emerald-400">₫0</h2>
      </div>
      <div class="glass-card rounded-2xl p-6">
        <span class="text-xs text-slate-400 uppercase font-bold tracking-wider">Tài Sản</span>
        <h2 class="text-3xl font-mono font-bold mt-2 text-slate-100">₫0</h2>
      </div>
      <div class="glass-card rounded-2xl p-6">
        <span class="text-xs text-slate-400 uppercase font-bold tracking-wider">Khoản Nợ</span>
        <h2 class="text-3xl font-mono font-bold mt-2 text-red-400">₫0</h2>
      </div>
    </div>
  </div>
</template>

<script setup>
// Dashboard view logic here
</script>
