<template>
  <div class="code-block-wrapper">
    <slot></slot>
    <CopyButton :content="codeContent" class="copy-button-absolute" @copy="handleCopy" />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import CopyButton from './CopyButton.vue'

const props = defineProps({
  codeElement: {
    type: HTMLElement,
    required: true
  }
})

const emit = defineEmits(['copy'])

const codeContent = computed(() => {
  return props.codeElement?.textContent || ''
})

const handleCopy = (success) => {
  emit('copy', success)
}
</script>

<style scoped>
.code-block-wrapper {
  position: relative;
  margin: 28px 0;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.code-block-wrapper:hover {
  box-shadow: 0 6px 24px rgba(102, 126, 234, 0.15);
  transform: translateY(-2px);
}

.copy-button-absolute {
  position: absolute !important;
  top: -20px !important;
  right: 10px !important;
  z-index: 50 !important;
  opacity: 0 !important;
  transform: translateY(-10px) !important;
  transition: all 0.4s ease !important;
}

.code-block-wrapper:hover .copy-button-absolute {
  opacity: 1 !important;
  transform: translateY(0) !important;
}

/* 移动端始终显�?*/
@media (max-width: 768px) {
  .copy-button-absolute {
    opacity: 1 !important;
    top: -18px !important;
    right: 8px !important;
  }
}
</style>