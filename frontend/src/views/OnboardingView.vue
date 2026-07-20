<template>
  <div class="min-h-screen flex items-center justify-center p-4">
    <div class="premium-card w-full max-w-2xl flex flex-col h-[600px]">
      <!-- Header -->
      <div class="p-4 border-b border-white/5 flex items-center">
        <div class="w-10 h-10 rounded-full bg-indigo-500/20 flex items-center justify-center mr-3">
          <span class="text-indigo-400 text-xl">🤖</span>
        </div>
        <div>
          <h2 class="text-lg font-semibold text-white">WealthOS Advisor</h2>
          <p class="text-xs text-slate-400">AI Onboarding</p>
        </div>
      </div>
      
      <!-- Chat Area -->
      <div class="flex-1 overflow-y-auto p-4 space-y-4" ref="chatContainer">
        <div v-for="(msg, index) in messages" :key="index" class="flex" :class="msg.role === 'user' ? 'justify-end' : 'justify-start'">
          <div 
            class="max-w-[80%] rounded-2xl px-4 py-2"
            :class="msg.role === 'user' ? 'bg-indigo-600 text-white' : 'bg-slate-700/50 text-slate-200'"
          >
            {{ msg.content }}
          </div>
        </div>
        <div v-if="loading" class="flex justify-start">
          <div class="bg-slate-700/50 text-slate-400 rounded-2xl px-4 py-2 flex space-x-1 items-center">
            <div class="w-2 h-2 bg-slate-400 rounded-full animate-bounce"></div>
            <div class="w-2 h-2 bg-slate-400 rounded-full animate-bounce" style="animation-delay: 0.1s"></div>
            <div class="w-2 h-2 bg-slate-400 rounded-full animate-bounce" style="animation-delay: 0.2s"></div>
          </div>
        </div>
      </div>

      <!-- Input Area -->
      <div class="p-4 border-t border-white/5 bg-slate-800/30">
        <form @submit.prevent="sendMessage" class="flex gap-2" v-if="!onboardingComplete">
          <input 
            v-model="input" 
            type="text" 
            placeholder="Nhập câu trả lời của bạn..." 
            class="flex-1 bg-slate-900/50 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500"
            :disabled="loading"
          >
          <button 
            type="submit" 
            class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg transition-colors"
            :disabled="!input.trim() || loading"
          >
            Gửi
          </button>
        </form>
        <div v-else class="text-center">
          <button @click="goToProfile" class="bg-emerald-600 hover:bg-emerald-700 text-white px-6 py-2 rounded-lg font-medium transition-colors">
            Xem Hồ Sơ Của Tôi
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { useProfileStore } from '../stores/profileStore';

const router = useRouter();
const profileStore = useProfileStore();

const messages = ref([]);
const input = ref('');
const loading = ref(false);
const chatContainer = ref(null);
const onboardingComplete = ref(false);

const scriptSteps = [
  "Chào mừng bạn đến với WealthOS! Để tôi có thể tư vấn tốt nhất, hãy cho tôi biết bạn sinh năm bao nhiêu?",
  "Cảm ơn! Tình trạng hôn nhân hiện tại của bạn là gì? Bạn có người phụ thuộc không?",
  "Thu nhập bình quân mỗi tháng của bạn (sau thuế) khoảng bao nhiêu?",
  "Còn chi phí sinh hoạt hàng tháng của bạn thì sao?",
  "Công việc hiện tại của bạn có ổn định không? Bạn có dự định lớn nào tốn nhiều tiền trong 3-5 năm tới không? (ví dụ: Mua nhà, sinh con, khởi nghiệp...)",
  "Hiện tại bạn có khoản nợ nào cần thanh toán hàng tháng không? (Vay mua nhà, thẻ tín dụng...)",
  "Mức độ chấp nhận rủi ro của bạn như thế nào? (Thấp/Trung bình/Cao)",
  "Cuối cùng, khi nghỉ hưu bạn mong muốn có mức sống như thế nào? (ví dụ: Chi tiêu khoảng 20tr/tháng, hay chỉ cần sống giản dị...). Tôi sẽ dựa vào đó để tự tính ra con số Tự Do Tài Chính cần thiết cho bạn."
];

let currentStep = 0;

const scrollToBottom = async () => {
  await nextTick();
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
  }
};

onMounted(() => {
  addBotMessage(scriptSteps[0]);
});

const addBotMessage = (content) => {
  messages.value.push({ role: 'assistant', content });
  scrollToBottom();
};

const sendMessage = async () => {
  if (!input.value.trim() || loading.value) return;
  
  messages.value.push({ role: 'user', content: input.value });
  input.value = '';
  scrollToBottom();
  
  currentStep++;
  
  if (currentStep < scriptSteps.length) {
    loading.value = true;
    setTimeout(() => {
      addBotMessage(scriptSteps[currentStep]);
      loading.value = false;
    }, 800);
  } else {
    loading.value = true;
    try {
      await profileStore.submitOnboarding(messages.value);
      onboardingComplete.value = true;
      addBotMessage("Tuyệt vời! Tôi đã tổng hợp xong hồ sơ đầu tư của bạn bằng AI.");
    } catch (error) {
      addBotMessage("Xin lỗi, có lỗi xảy ra khi lưu hồ sơ. Vui lòng kiểm tra API Key.");
    } finally {
      loading.value = false;
    }
  }
};

const goToProfile = () => {
  router.push('/profile');
};
</script>
