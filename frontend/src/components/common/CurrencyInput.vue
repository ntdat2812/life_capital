<template>
  <input
    type="text"
    :value="displayValue"
    @input="handleInput"
    @blur="handleBlur"
    @focus="handleFocus"
    class="w-full bg-slate-900/50 border border-white/10 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500"
  />
</template>

<script setup>
import { ref, watch, computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: [Number, String],
    default: null
  }
})

const emit = defineEmits(['update:modelValue'])

const isFocused = ref(false)
const localValue = ref('')

// Format number with en-US (comma for thousands, dot for decimals)
// We use en-US because it's standard in tech and avoids issues with dot/comma swapping on keyboards.
const formatNumber = (val) => {
  if (val === null || val === undefined || val === '') return ''
  return new Intl.NumberFormat('en-US', { maximumFractionDigits: 6 }).format(val)
}

const displayValue = computed(() => {
  if (isFocused.value) {
    return localValue.value
  }
  return formatNumber(props.modelValue)
})

watch(() => props.modelValue, (newVal) => {
  if (!isFocused.value) {
    localValue.value = formatNumber(newVal)
  }
}, { immediate: true })

const handleInput = (event) => {
  // Allow digits and dot
  let rawValue = event.target.value.replace(/[^0-9.]/g, '')
  
  if (rawValue === '') {
    localValue.value = ''
    emit('update:modelValue', null)
    event.target.value = ''
    return
  }
  
  // Handle multiple dots by keeping only the first one
  const parts = rawValue.split('.')
  if (parts.length > 2) {
    rawValue = parts[0] + '.' + parts.slice(1).join('')
  }
  
  localValue.value = rawValue

  // If the string ends with a dot or has trailing zeros after dot, we don't emit a clean float yet 
  // to avoid breaking the user's typing flow (e.g. "1." or "1.0").
  // But we DO need to emit something so the parent knows.
  const numValue = parseFloat(rawValue)
  
  if (!isNaN(numValue)) {
    emit('update:modelValue', numValue)
  }
  
  // Apply formatting dynamically if there's no decimal point to add commas
  if (!rawValue.includes('.')) {
    localValue.value = new Intl.NumberFormat('en-US').format(numValue)
    event.target.value = localValue.value
  } else {
    // Let the user type their decimal freely
    event.target.value = localValue.value
  }
}

const handleBlur = (e) => {
  isFocused.value = false
  if (props.modelValue === '' || props.modelValue === null) {
    emit('update:modelValue', null)
  }
}

const handleFocus = () => {
  isFocused.value = true
  // When focused, if it's a clean number, just show it normally without commas so they can edit easily?
  // Actually, keeping commas is fine, our replace(/[^0-9.]/g, '') handles it.
}
</script>
