<template>
  <router-view />
</template>

<script setup>
// import { onMounted } from 'vue'
// import { useRbacStore } from './store/rbac'

// 应用初始化：检查token并加载菜单
// onMounted(async () => {
//   const token = localStorage.getItem('token')
//   if (token) {
//     try {
//       const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
//       const rbacStore = useRbacStore()
//       // 如果用户信息存在，则初始化权限
//       if (userInfo.email) {
//         console.log('应用初始化，开始加载用户菜单...')
//         await rbacStore.initUserPermission(userInfo)
//         console.log('应用初始化菜单加载完成')
//       }
//     } catch (error) {
//       console.error('应用初始化菜单加载失败:', error)
//     }
//   }
// })
</script>

<style>
@import './styles/variables.scss';
@import './styles/common.scss';
</style>
