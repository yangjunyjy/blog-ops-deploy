<template>
  <div class="navbar">
    <div class="navbar-left">
      <div class="collapse-btn" @click="handleCollapse">
        <el-icon>
          <Fold v-if="!isCollapse" />
          <Expand v-else />
        </el-icon>
      </div>
      <breadcrumb />
    </div>
    <div class="navbar-right">
      <el-dropdown @command="handleCommand" trigger="click">
        <div class="user-info">
          <el-avatar :size="40" :src="userInfo.avatar" class="user-avatar" />
          <span class="username">{{ userInfo.nickname }}</span>
          <el-icon class="dropdown-icon">
            <ArrowDown />
          </el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">
              <el-icon>
                <User />
              </el-icon>
              个人中心
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon>
                <SwitchButton />
              </el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { logout } from '@/api/auth'
import { useRbacStore } from '@/store/rbac'
import {
  Fold,
  Expand,
  User,
  SwitchButton,
  ArrowDown,
  ArrowDownBold
} from '@element-plus/icons-vue'
import Breadcrumb from './Breadcrumb.vue'

const props = defineProps({
  isCollapse: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['toggle-collapse'])

const router = useRouter()
const rbacStore = useRbacStore()

const userInfo = rbacStore.currentUser

const handleCollapse = () => {
  emit('toggle-collapse')
}

const handleCommand = async (command) => {
  switch (command) {
    case 'profile':
      router.push('/dashboard/profile')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        // 清理 RBAC store 中的所有数据
        rbacStore.resetAllData()
        sessionStorage.removeItem('token')
        sessionStorage.removeItem('userInfo')
        sessionStorage.setItem('routeInitialized', 'false')

        // 清理动态路由
        const { resetRoutes } = await import('@/router/index.js')
        resetRoutes()

        const data = await logout()
        if (data.code == 200) {
          ElMessage.success("退出登录成功")
        }
      } catch (error) {
        ElMessage.error("清除会话发生错误")
      } finally {
        router.push('/login')
      }
      break
  }
}
</script>

<style scoped lang="scss">
.navbar {
  height: 60px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  z-index: 10;

  .navbar-left {
    display: flex;
    align-items: center;
    gap: 16px;

    .collapse-btn {
      font-size: 20px;
      cursor: pointer;
      color: #595959;
      transition: all 0.3s;
      padding: 8px;
      border-radius: 4px;
      display: flex;
      align-items: center;
      justify-content: center;

      &:hover {
        color: #1890ff;
        background: rgba(24, 144, 255, 0.08);
      }
    }
  }

  .navbar-right {
    .user-info {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      padding: 6px 12px;
      border-radius: 4px;
      transition: all 0.3s;

      &:hover {
        background: rgba(24, 144, 255, 0.08);
      }

      .user-avatar {
        border: 1px solid #d9d9d9;
      }

      .username {
        font-size: 14px;
        font-weight: 500;
        color: #262626;
      }

      .dropdown-icon {
        font-size: 12px;
        color: #8c8c8c;
        transition: transform 0.3s;
      }
    }

    .user-info:hover .dropdown-icon {
      transform: rotate(180deg);
    }
  }
}
</style>
