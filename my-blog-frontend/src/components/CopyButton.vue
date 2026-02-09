<template>
  <button 
    class="copy-button"
    :class="{ 'copied': isCopied }"
    @click="handleCopy"
    :title="isCopied ? '已复制' : '复制代码'"
    aria-label="复制代码"
  >
    <div class="papers-icon">
      <div class="paper front">
        <div class="paper-content">
          <div class="line line-1"></div>
          <div class="line line-2"></div>
          <div class="line line-3"></div>
        </div>
      </div>
      <div class="paper back">
        <div class="paper-content">
          <div class="line line-1"></div>
          <div class="line line-2"></div>
          <div class="line line-3"></div>
        </div>
      </div>
    </div>
    
    <span class="copy-text">
      {{ isCopied ? '已复制' : '复制' }}
    </span>
    
    <div class="success-check" v-if="isCopied">
      <svg viewBox="0 0 24 24" width="16" height="16">
        <path fill="currentColor" d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"></path>
      </svg>
    </div>
  </button>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  content: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['copy'])

const isCopied = ref(false)

const handleCopy = async () => {
  try {
    await navigator.clipboard.writeText(props.content)
    emit('copy', true)
    
    // 显示成功状态
    isCopied.value = true
    setTimeout(() => {
      isCopied.value = false
    }, 2000)
  } catch (error) {
    // 降级方案
    try {
      const textArea = document.createElement('textarea')
      textArea.value = props.content
      textArea.style.position = 'fixed'
      textArea.style.opacity = '0'
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
      
      emit('copy', true)
      isCopied.value = true
      setTimeout(() => {
        isCopied.value = false
      }, 2000)
    } catch (err) {
      emit('copy', false)
    }
  }
}
</script>

<style scoped>
.copy-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(255, 255, 255, 0.85) 100%);
  backdrop-filter: blur(20px);
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  color: #4a5568;
  transition: all 0.4s ease;
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.1),
    inset 0 0 0 1px rgba(255, 255, 255, 0.8);
  position: relative;
  overflow: hidden;
  min-width: 80px;
  justify-content: center;
}

.copy-button:hover {
  background: linear-gradient(135deg, rgba(255, 255, 255, 1) 0%, rgba(255, 255, 255, 0.9) 100%);
  transform: translateY(-2px);
  box-shadow: 
    0 6px 20px rgba(102, 126, 234, 0.2),
    inset 0 0 0 1px rgba(255, 255, 255, 0.9);
  color: #667eea;
}

.copy-button:active {
  transform: translateY(0);
  box-shadow: 
    0 2px 4px rgba(0, 0, 0, 0.1),
    inset 0 0 0 1px rgba(255, 255, 255, 0.8);
}

.copy-button.copied {
  background: linear-gradient(135deg, rgba(72, 187, 120, 0.95) 0%, rgba(56, 161, 105, 0.85) 100%);
  color: white;
  box-shadow: 
    0 2px 8px rgba(72, 187, 120, 0.3),
    inset 0 0 0 1px rgba(255, 255, 255, 0.3);
}

.copy-button.copied:hover {
  background: linear-gradient(135deg, rgba(72, 187, 120, 1) 0%, rgba(56, 161, 105, 0.95) 100%);
}

/* 两页纸图�?*/
.papers-icon {
  position: relative;
  width: 20px;
  height: 24px;
  perspective: 600px;
}

.paper {
  position: absolute;
  width: 100%;
  height: 100%;
  background: #f8fafc;
  border-radius: 2px;
  box-shadow: 
    0 1px 3px rgba(0, 0, 0, 0.1),
    inset 0 0 0 1px rgba(0, 0, 0, 0.05);
  transition: all 0.4s ease;
  transform-origin: bottom center;
}

.paper.front {
  z-index: 2;
  transform: rotateX(0deg);
}

.paper.back {
  z-index: 1;
  transform: rotateX(-10deg) translateY(-2px);
  opacity: 0.8;
  filter: blur(0.5px);
}

.copy-button:hover .paper.front {
  transform: rotateX(5deg) translateY(-1px);
}

.copy-button:hover .paper.back {
  transform: rotateX(-15deg) translateY(-3px);
  opacity: 0.6;
}

.copy-button.copied .paper.front {
  transform: rotateX(20deg) translateY(-2px);
}

.copy-button.copied .paper.back {
  transform: rotateX(-25deg) translateY(-4px);
  opacity: 0.4;
}

.paper-content {
  padding: 4px;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.line {
  height: 2px;
  background: rgba(102, 126, 234, 0.2);
  border-radius: 1px;
}

.line-1 { width: 100%; }
.line-2 { width: 80%; align-self: flex-end; }
.line-3 { width: 60%; }

.copy-button.copied .line {
  background: rgba(255, 255, 255, 0.6);
}

.copy-text {
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.5px;
  transition: color 0.3s;
}

.copy-button.copied .copy-text {
  color: white;
}

.success-check {
  display: flex;
  align-items: center;
  justify-content: center;
  animation: checkIn 0.3s ease-out;
}

@keyframes checkIn {
  0% {
    opacity: 0;
    transform: scale(0);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

/* 暗黑模式适配 */
html.dark .copy-button {
  background: linear-gradient(135deg, rgba(45, 55, 72, 0.95) 0%, rgba(30, 41, 59, 0.85) 100%);
  color: #a0aec0;
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.3),
    inset 0 0 0 1px rgba(255, 255, 255, 0.1);
}

html.dark .copy-button:hover {
  background: linear-gradient(135deg, rgba(56, 178, 172, 0.95) 0%, rgba(49, 151, 149, 0.85) 100%);
  color: white;
  box-shadow: 
    0 6px 20px rgba(56, 178, 172, 0.2),
    inset 0 0 0 1px rgba(255, 255, 255, 0.2);
}

html.dark .paper {
  background: #2d3748;
  box-shadow: 
    0 1px 3px rgba(0, 0, 0, 0.3),
    inset 0 0 0 1px rgba(255, 255, 255, 0.05);
}

html.dark .line {
  background: rgba(160, 174, 192, 0.2);
}

html.dark .copy-button.copied {
  background: linear-gradient(135deg, rgba(72, 187, 120, 0.95) 0%, rgba(56, 161, 105, 0.85) 100%);
}

/* 响应式调�?*/
@media (max-width: 768px) {
  .copy-button {
    padding: 6px 12px;
    min-width: 70px;
  }
  
  .papers-icon {
    width: 16px;
    height: 20px;
  }
  
  .copy-text {
    font-size: 11px;
  }
}
</style>