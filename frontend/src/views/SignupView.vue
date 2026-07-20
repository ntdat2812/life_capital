<template>
  <div class="max-w-md mx-auto my-12 glass-card rounded-2xl p-8 space-y-6">
    <div>
      <h1 class="text-2xl font-bold font-outfit text-slate-100 text-center">Đăng Ký</h1>
      <p class="text-xs text-slate-400 text-center mt-1">Tạo tài khoản Hệ điều hành Quản lý Tài sản Cá nhân</p>
    </div>

    <form @submit.prevent="handleSignup" class="space-y-4 text-sm">
      <div v-if="error" class="p-3 bg-red-900/30 border border-red-500/30 rounded-xl text-red-400 text-xs">
        {{ error }}
      </div>
      <div>
        <label class="block text-slate-400 mb-1 font-bold">Họ và Tên</label>
        <input v-model="name" type="text" required class="w-full bg-slate-900 border border-slate-700 rounded-xl px-3 py-2 focus:outline-none focus:border-indigo-500 text-slate-200">
      </div>
      <div>
        <label class="block text-slate-400 mb-1 font-bold">Email</label>
        <input v-model="email" type="email" required class="w-full bg-slate-900 border border-slate-700 rounded-xl px-3 py-2 focus:outline-none focus:border-indigo-500 text-slate-200">
      </div>
      <div>
        <label class="block text-slate-400 mb-1 font-bold">Mật khẩu</label>
        <input v-model="password" type="password" required class="w-full bg-slate-900 border border-slate-700 rounded-xl px-3 py-2 focus:outline-none focus:border-indigo-500 text-slate-200">
      </div>
      <button :disabled="loading" type="submit" class="w-full py-2.5 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white rounded-xl font-semibold shadow-md transition">
        <span v-if="loading">Đang xử lý...</span>
        <span v-else>Đăng ký</span>
      </button>
    </form>

    <div class="relative flex py-2 items-center">
        <div class="flex-grow border-t border-slate-700"></div>
        <span class="flex-shrink-0 mx-4 text-slate-500 text-xs">Hoặc</span>
        <div class="flex-grow border-t border-slate-700"></div>
    </div>

    <div class="flex justify-center">
      <GoogleLogin :callback="handleGoogleCallback" />
    </div>

    <p class="text-center text-xs text-slate-400 mt-4">
      Đã có tài khoản? 
      <router-link to="/login" class="text-indigo-400 hover:text-indigo-300 font-semibold">Đăng nhập</router-link>
    </p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/authStore'

const router = useRouter()
const authStore = useAuthStore()

const name = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const handleSignup = async () => {
  error.value = ''
  loading.value = true
  try {
    await authStore.signup(name.value, email.value, password.value)
    // Auto login after signup
    await authStore.login(email.value, password.value)
    router.push('/')
  } catch (err) {
    error.value = err.response?.data?.message || 'Đăng ký thất bại. Email có thể đã tồn tại.'
  } finally {
    loading.value = false
  }
}

const handleGoogleCallback = async (response) => {
  if (response.credential) {
    try {
      await authStore.loginWithGoogle(response.credential)
      router.push('/')
    } catch (err) {
      error.value = 'Đăng nhập bằng Google thất bại.'
    }
  }
}
</script>
