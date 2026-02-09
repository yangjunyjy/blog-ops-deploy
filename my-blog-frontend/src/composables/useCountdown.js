import { ref } from 'vue'
import { ElMessage } from 'element-plus'

/**
 * 倒计时composable
 * @param {number} seconds - 倒计时秒数
 * @param {Function} callback - 倒计时结束时的回调
 * @returns {Object} 包含countdown, start, reset的对象
 */
export function useCountdown(seconds = 60, callback = null) {
  const countdown = ref(0)
  let timer = null

  const start = () => {
    if (countdown.value > 0) return

    countdown.value = seconds
    timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
        if (callback) callback()
      }
    }, 1000)
  }

  const reset = () => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
    countdown.value = 0
  }

  return {
    countdown,
    start,
    reset
  }
}
