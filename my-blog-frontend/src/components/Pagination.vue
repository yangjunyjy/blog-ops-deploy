<script setup>
const props = defineProps({
  currentPage: {
    type: Number,
    required: true
  },
  total: {
    type: Number,
    required: true
  },
  pageSize: {
    type: Number,
    default: 12
  },
  totalPages: {
    type: Number,
    required: true
  },
  layout: {
    type: String,
    default: 'prev, pager, next'
  }
})

const emit = defineEmits(['change'])

const handlePageChange = (page) => {
  emit('change', page)
}
</script>

<template>
  <div class="pagination" v-if="total > 0">
    <div class="pagination-container">
      <button
        class="pagination-btn prev"
        :disabled="currentPage === 1"
        @click="handlePageChange(currentPage - 1)"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="15 18 9 12 15 6"></polyline>
        </svg>
      </button>

      <div class="pagination-numbers">
        <button
          v-for="page in totalPages"
          :key="page"
          class="pagination-number"
          :class="{ active: page === currentPage }"
          @click="handlePageChange(page)"
        >
          {{ page }}
        </button>
      </div>

      <button
        class="pagination-btn next"
        :disabled="currentPage === totalPages"
        @click="handlePageChange(currentPage + 1)"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="9 18 15 12 9 6"></polyline>
        </svg>
      </button>
    </div>

    <div class="pagination-info">
      <span>第 {{ currentPage }} 页/ 第 {{ totalPages }} 页</span>
      <span>总计 {{ total }} 个</span>
    </div>
  </div>
</template>

<style scoped>
.pagination {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 20px 0;
}

.pagination-container {
  display: flex;
  align-items: center;
  gap: 6px;
  background: linear-gradient(135deg, #fff 0%, #f8fafc 100%);
  padding: 8px 16px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(102, 126, 234, 0.1);
  border: 1px solid rgba(102, 126, 234, 0.08);
  transition: all 0.3s;
}

html.dark .pagination-container {
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
  border-color: rgba(102, 126, 234, 0.15);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.25);
}

.pagination-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  background: #3b82f6;
  color: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.25);
}

html.dark .pagination-btn {
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.4);
}

.pagination-btn:hover:not(:disabled) {
  transform: translateY(-1px) scale(1.05);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.pagination-btn:disabled {
  background: #e2e8f0;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

html.dark .pagination-btn:disabled {
  background: #4a5568;
}

.pagination-numbers {
  display: flex;
  align-items: center;
  gap: 4px;
}

.pagination-number {
  min-width: 32px;
  height: 32px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  background: #f8fafc;
  color: #4a5568;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
}

html.dark .pagination-number {
  background: #1a202c;
  border-color: #4a5568;
  color: #cbd5e0;
}

.pagination-number:hover {
  border-color: #667eea;
  color: #667eea;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.15);
}

.pagination-number.active {
  background: #3b82f6;
  border-color: transparent;
  color: #fff;
  transform: scale(1.08);
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.4);
}

.pagination-info {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 12px;
  color: #718096;
  padding: 6px 12px;
  background: rgba(102, 126, 234, 0.04);
  border-radius: 12px;
}

html.dark .pagination-info {
  color: #a0aec0;
  background: rgba(102, 126, 234, 0.08);
}

.pagination-info span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.pagination-info span::before {
  content: '';
  display: inline-block;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #3b82f6;
}

/* 动画效果 */
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.8;
  }
}

.pagination-number.active {
  animation: pulse 2s ease-in-out infinite;
}

/* 响应�?*/
@media (max-width: 768px) {
  .pagination {
    padding: 16px 0;
    gap: 10px;
  }

  .pagination-container {
    padding: 6px 12px;
    gap: 4px;
  }

  .pagination-btn {
    width: 28px;
    height: 28px;
  }

  .pagination-btn svg {
    width: 14px;
    height: 14px;
  }

  .pagination-number {
    min-width: 28px;
    height: 28px;
    font-size: 12px;
  }

  .pagination-info {
    flex-direction: column;
    gap: 6px;
    text-align: center;
    font-size: 11px;
    padding: 5px 10px;
  }
}
</style>
