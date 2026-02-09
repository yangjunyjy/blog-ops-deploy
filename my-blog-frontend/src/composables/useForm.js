import { ref } from 'vue'
import { ElMessage } from 'element-plus'

/**
 * 表单验证composable
 * @param {Object} rules - 验证规则
 * @returns {Object} 包含validate, validateField, resetFields的对象
 */
export function useForm(rules = {}) {
  const errors = ref({})
  const isSubmitting = ref(false)

  const validate = (formData) => {
    errors.value = {}
    let isValid = true

    for (const field in rules) {
      const rule = rules[field]
      const value = formData[field]

      if (rule.required && !value) {
        errors.value[field] = rule.message || '此项为必填项'
        isValid = false
      } else if (rule.pattern && !rule.pattern.test(value)) {
        errors.value[field] = rule.message || '格式不正确'
        isValid = false
      } else if (rule.minLength && value.length < rule.minLength) {
        errors.value[field] = rule.message || `至少需要${rule.minLength}个字符`
        isValid = false
      } else if (rule.validator) {
        const result = rule.validator(value, formData)
        if (result !== true) {
          errors.value[field] = result || '验证失败'
          isValid = false
        }
      }
    }

    return isValid
  }

  const validateField = (field, value) => {
    if (!rules[field]) return true

    const rule = rules[field]
    errors.value[field] = ''

    if (rule.required && !value) {
      errors.value[field] = rule.message || '此项为必填项'
      return false
    } else if (rule.pattern && !rule.pattern.test(value)) {
      errors.value[field] = rule.message || '格式不正确'
      return false
    } else if (rule.minLength && value.length < rule.minLength) {
      errors.value[field] = rule.message || `至少需要${rule.minLength}个字符`
      return false
    } else if (rule.validator) {
      const result = rule.validator(value)
      if (result !== true) {
        errors.value[field] = result || '验证失败'
        return false
      }
    }

    return true
  }

  const resetFields = () => {
    errors.value = {}
    isSubmitting.value = false
  }

  return {
    errors,
    isSubmitting,
    validate,
    validateField,
    resetFields
  }
}
