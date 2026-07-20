<template>
  <div class="max-w-md mx-auto my-12 glass-card rounded-2xl p-8 space-y-6">
    <div>
      <h1 class="text-2xl font-bold font-outfit text-slate-100 text-center">Đăng Nhập</h1>
      <p class="text-xs text-slate-400 text-center mt-1">Đăng nhập vào Hệ điều hành Quản lý Tài sản Cá nhân</p>
    </div>

    <!-- Form layout placeholder -->
    <div class="space-y-4 text-sm">
      <div>
        <label class="block text-slate-400 mb-1 font-bold">Email</label>
        <input type="email" class="w-full bg-slate-900 border border-slate-700 rounded-xl px-3 py-2 focus:outline-none focus:border-indigo-500 text-slate-200">
      </div>
      <div>
        <label class="block text-slate-400 mb-1 font-bold">Mật khẩu</label>
        <input type="password" class="w-full bg-slate-900 border border-slate-700 rounded-xl px-3 py-2 focus:outline-none focus:border-indigo-500 text-slate-200">
      </div>
      <button class="w-full py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-semibold shadow-md transition">
        Đăng nhập
      </button>
    </div>
  </div>
</template>

<script setup>
</script>
