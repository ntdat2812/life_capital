<template>
  <div class="relative" ref="selectContainer">
    <!-- Nút bấm hiển thị lựa chọn hiện tại -->
    <button 
      type="button"
      @click="toggleOpen"
      class="w-full bg-slate-900/50 border rounded-lg px-3 h-[42px] text-left flex justify-between items-center transition-colors focus:outline-none focus:border-indigo-500"
      :class="[
        disabled ? 'opacity-50 cursor-not-allowed border-slate-700/50' : (isOpen ? 'border-indigo-500 shadow-[0_0_10px_rgba(99,102,241,0.2)]' : 'border-slate-700 hover:border-slate-600')
      ]"
      :disabled="disabled"
    >
      <span class="truncate pr-2" :class="modelValue ? 'text-white' : 'text-slate-400'">
        {{ selectedLabel }}
      </span>
      <svg class="w-4 h-4 shrink-0 text-slate-400 transition-transform duration-200" :class="{ 'rotate-180': isOpen }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
    </button>

    <!-- Danh sách dropdown -->
    <transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="transform scale-95 opacity-0"
      enter-to-class="transform scale-100 opacity-100"
      leave-active-class="transition duration-75 ease-in"
      leave-from-class="transform scale-100 opacity-100"
      leave-to-class="transform scale-95 opacity-0"
    >
      <div 
        v-if="isOpen" 
        class="absolute z-[100] w-full mt-1 bg-slate-800 border border-slate-600 rounded-lg shadow-xl max-h-60 overflow-y-auto custom-scrollbar"
      >
        <ul class="py-1">
          <li 
            v-for="option in options" 
            :key="option.value"
            @click="selectOption(option)"
            class="px-4 py-2 cursor-pointer transition-colors text-sm"
            :class="option.value === modelValue ? 'bg-indigo-500/20 text-indigo-300 font-medium' : 'text-slate-300 hover:bg-slate-700 hover:text-white'"
          >
            {{ option.label }}
          </li>
        </ul>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  options: {
    type: Array,
    required: true // Format: [{ label: 'Name', value: 'id' }]
  },
  placeholder: {
    type: String,
    default: 'Vui lòng chọn...'
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false)
const selectContainer = ref(null)

const selectedLabel = computed(() => {
  const selected = props.options.find(opt => opt.value === props.modelValue)
  return selected ? selected.label : props.placeholder
})

const toggleOpen = () => {
  if (!props.disabled) {
    isOpen.value = !isOpen.value
  }
}

const selectOption = (option) => {
  emit('update:modelValue', option.value)
  isOpen.value = false
}

const handleClickOutside = (event) => {
  if (selectContainer.value && !selectContainer.value.contains(event.target)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
