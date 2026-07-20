<template>
  <input 
    type="text" 
    :value="displayValue" 
    @input="handleInput"
    @blur="handleBlur"
    placeholder="DD/MM/YYYY" 
    class="w-full bg-slate-900/50 border border-white/10 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500 font-mono" 
  />
</template>

<script setup>
import { ref, watch } from 'vue';

const props = defineProps({
  modelValue: {
    type: [String, Date],
    default: ''
  }
});

const emit = defineEmits(['update:modelValue']);

const formatToDisplay = (val) => {
  if (!val) return '';
  const d = new Date(val);
  if (isNaN(d.getTime())) return '';
  
  // Xử lý bù trừ múi giờ để hiển thị đúng ngày
  const day = String(d.getDate()).padStart(2, '0');
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const year = d.getFullYear();
  return `${day}/${month}/${year}`;
};

const displayValue = ref(formatToDisplay(props.modelValue));
let lastValidValue = props.modelValue;

watch(() => props.modelValue, (newVal) => {
  if (newVal !== lastValidValue) {
    displayValue.value = formatToDisplay(newVal);
    lastValidValue = newVal;
  }
});

const handleInput = (e) => {
  // Lọc chỉ giữ lại số
  let val = e.target.value.replace(/\D/g, ''); 
  
  if (val.length > 8) {
    val = val.slice(0, 8);
  }

  // Format mask DD/MM/YYYY
  let formatted = '';
  if (val.length > 4) {
    formatted = `${val.slice(0, 2)}/${val.slice(2, 4)}/${val.slice(4)}`;
  } else if (val.length > 2) {
    formatted = `${val.slice(0, 2)}/${val.slice(2)}`;
  } else {
    formatted = val;
  }
  
  displayValue.value = formatted;
  
  // Khi đủ 8 số, thử parse sang ISO
  if (val.length === 8) {
    const day = parseInt(val.slice(0, 2), 10);
    const month = parseInt(val.slice(2, 4), 10) - 1;
    const year = parseInt(val.slice(4), 10);
    
    // Kiểm tra tính hợp lệ của ngày
    const d = new Date(year, month, day);
    if (d.getFullYear() === year && d.getMonth() === month && d.getDate() === day) {
      // Offset timezone to avoid UTC shifts affecting the stored date
      const isoDate = `${year}-${String(month+1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
      lastValidValue = isoDate;
      emit('update:modelValue', isoDate);
    } else {
      emit('update:modelValue', null);
    }
  } else if (val.length === 0) {
    emit('update:modelValue', null);
  }
};

const handleBlur = (e) => {
  // Có thể thêm logic validate khi blur
};
</script>
