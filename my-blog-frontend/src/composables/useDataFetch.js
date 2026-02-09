import { ref, onMounted } from 'vue'

/**
 * 通用数据获取composable
 * @param {Function} apiCall - API调用函数
 * @param {Array|Object} initialData - 初始数据
 * @returns {Object} 包含data, loading, fetchData的对象
 */
export function useDataFetch(apiCall, initialData = []) {
  const data = ref(initialData)
  const loading = ref(true)
  const error = ref(null)

  const fetchData = async () => {
    loading.value = true
    error.value = null
    try {
      const res = await apiCall()
      data.value = res.data || initialData
    } catch (err) {
      error.value = err
      console.error('数据加载失败:', err)
    } finally {
      loading.value = false
    }
  }

  onMounted(fetchData)

  return {
    data,
    loading,
    error,
    fetchData
  }
}
