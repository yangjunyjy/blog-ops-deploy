import { useRbacStore } from '../store/rbac'

/**
 * 权限指令
 * 使用方法：v-permission="'system:user:create'"
 * 如果没有权限，元素会被移除
 */
export default {
  mounted(el, binding) {
    const permission = binding.value
    if (!permission) return

    const rbacStore = useRbacStore()
    if (!rbacStore.hasPermission(permission)) {
      el.parentNode?.removeChild(el)
    }
  },
  updated(el, binding) {
    const permission = binding.value
    if (!permission) return

    const rbacStore = useRbacStore()
    if (!rbacStore.hasPermission(permission)) {
      el.parentNode?.removeChild(el)
    }
  }
}
