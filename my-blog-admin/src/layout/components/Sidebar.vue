<template>
  <div class="sidebar" :class="{ collapsed: isCollapse }">
    <div class="logo">
      <div class="logo-icon">
        <el-icon>
          <Promotion />
        </el-icon>
      </div>
      <h2 v-show="!isCollapse">博客管理系统</h2>
    </div>
    <el-menu :default-active="activeMenu" :collapse="isCollapse" background-color="#001529" active-text-color="#1890ff"
      :unique-opened="true" router>
      <template v-for="menu in menus" :key="menu.id">
        <el-menu-item v-if="menu.type === 'menu'" :index="menu.path" v-show="menu.isVisible">
          <el-icon>
            <component :is="getIcon(menu.icon)" />
          </el-icon>
          <template #title>{{ menu.title }}</template>
        </el-menu-item>
        <el-sub-menu v-else :index="String(menu.id)">
          <template #title>
            <el-icon>
              <component :is="getIcon(menu.icon)" />
            </el-icon>
            <span>{{ menu.title }}</span>
          </template>
          <el-menu-item v-for="submenu in menu.children" :key="submenu.id" :index="submenu.path"
            v-show="menu.isVisible">
            <el-icon>
              <component :is="getIcon(submenu.icon)" />
            </el-icon>
            <template #title>{{ submenu.title }}</template>
          </el-menu-item>
        </el-sub-menu>
      </template>
    </el-menu>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  HomeFilled,
  Document,
  ChatDotRound,
  User,
  DataLine,
  Folder,
  PriceTag,
  Collection,
  TrendCharts,
  UserFilled,
  Setting,
  Menu,
  Avatar,
  CircleCheck,
  Promotion,
  FolderOpened,
  Platform,
  Management
} from '@element-plus/icons-vue'
import { useRbacStore } from '../../store/rbac'

const props = defineProps({
  isCollapse: {
    type: Boolean,
    default: false
  }
})

const route = useRoute()
const activeMenu = computed(() => route.path)
const rbacStore = useRbacStore()

// 图标映射
const iconMap = {
  HomeFilled,
  Document,
  ChatDotRound,
  User,
  DataLine,
  Folder,
  PriceTag,
  Collection,
  TrendCharts,
  UserFilled,
  Setting,
  Menu,
  Avatar,
  CircleCheck,
  Promotion,
  FolderOpened,
  Platform,
  Management
}

// 获取图标组件
const getIcon = (iconName) => {
  return iconMap[iconName] || Document
}

// 从 RBAC store 获取菜单（只显示可见的目录和菜单，不显示按钮）
const menus = computed(() => {
  const menuList = rbacStore.currentMenus || []
  console.log('Sidebar 获取菜单数据:', menuList)

  // 递归过滤菜单，只显示 is_visible 为 true 的菜单
  const filterVisibleMenus = (menus) => {
    return menus.filter(menu => {
      // 过滤掉按钮类型和不可见的菜单
      if (menu.type === 'button' || !menu.isVisible) {
        return false
      }
      // 如果有子菜单，递归过滤
      if (menu.children && menu.children.length > 0) {
        menu.children = filterVisibleMenus(menu.children)
      }
      return true
    })
  }

  return filterVisibleMenus(menuList)
})
</script>

<style scoped lang="scss">
.sidebar {
  width: 240px;
  height: 100%;
  background: #fff;
  transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
  box-shadow: 2px 0 12px rgba(0, 0, 0, 0.15);
  overflow-x: hidden;
  overflow-y: auto;
  position: relative;
  color: #000000; // 设置默认文字颜色为黑色

  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 1px;
    background: linear-gradient(180deg, transparent, rgba(183, 102, 234, 0.3), transparent); // 改为Logo颜色
  }

  &::-webkit-scrollbar {
    width: 5px;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.15); // 滚动条改为黑色系
    border-radius: 3px;
    transition: background 0.2s;

    &:hover {
      background: rgba(0, 0, 0, 0.25);
    }
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }

  &.collapsed {
    width: 64px;

    .logo h2 {
      opacity: 0;
      width: 0;
    }
  }

  .logo {
    height: 61px;
    display: flex;
    align-items: center;
    gap: 12px;
    color: #a072ee; // Logo文字颜色
    padding: 0 20px;
    transition: all 0.3s;
    border-bottom: 1px solid rgba(0, 0, 0, 0.08); // 边框改为黑色系
    position: relative;
    background: white; // 改为Logo颜色

    &::after {
      content: '';
      position: absolute;
      bottom: -1px;
      left: 20px;
      right: 20px;
      height: 2px;
      background: linear-gradient(90deg, transparent, #a072ee, transparent); // 改为Logo颜色
    }

    .logo-icon {
      font-size: 28px;
      color: #a072ee;
      filter: drop-shadow(0 0 12px rgba(183, 102, 234, 0.6)); // 改为Logo颜色
      animation: iconFloat 3s ease-in-out infinite;
    }

    h2 {
      margin: 0;
      font-size: 18px;
      font-weight: 700;
      color: #a072ee; // 系统名称改为Logo颜色
      white-space: nowrap;
      overflow: hidden;
      transition: all 0.3s;
      letter-spacing: 0.8px;
      // 移除渐变效果，直接使用Logo颜色
    }
  }

  @keyframes iconFloat {

    0%,
    100% {
      transform: translateY(0);
    }

    50% {
      transform: translateY(-3px);
    }
  }

  :deep(.el-menu) {
    border-right: none;
    background: transparent;
    padding: 8px 0;

    // 设置图标颜色
    .el-icon {
      color: #000000;
    }
  }

  :deep(.el-menu-item),
  :deep(.el-sub-menu__title) {
    margin: 3px 10px;
    border-radius: 8px;
    transition: all 0.25s cubic-bezier(0.645, 0.045, 0.355, 1);
    height: 48px;
    line-height: 48px;
    position: relative;
    color: #000000 !important; // 设置菜单文字颜色为黑色

    .el-icon {
      font-size: 18px;
      margin-right: 8px;
      transition: transform 0.2s;
      color: #000000; // 设置图标颜色为黑色
    }

    &:hover {
      background: rgba(183, 102, 234, 0.08) !important; // hover背景改为Logo颜色
      transform: translateX(3px);
      color: #000000 !important; // hover文字颜色保持黑色

      .el-icon {
        transform: scale(1.1);
        color: #a072ee !important; // hover时图标改为Logo颜色
      }
    }
  }

  :deep(.el-menu-item.is-active) {
    background: linear-gradient(90deg, rgba(183, 102, 234, 0.15) 0%, rgba(183, 102, 234, 0.08) 100%) !important; // 激活背景改为Logo颜色
    color: #a072ee !important; // 激活文字颜色改为Logo颜色
    font-weight: 500;
    box-shadow: 0 2px 8px rgba(183, 102, 234, 0.2); // 阴影改为Logo颜色

    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 50%;
      transform: translateY(-50%);
      width: 4px;
      height: 60%;
      background: linear-gradient(180deg, #a072ee, #a072ee); // 改为Logo颜色
      border-radius: 0 4px 4px 0;
      box-shadow: 0 0 12px rgba(183, 102, 234, 0.8);
    }

    .el-icon {
      color: #a072ee; // 激活图标颜色改为Logo颜色
      filter: drop-shadow(0 0 8px rgba(183, 102, 234, 0.6));
    }
  }

  :deep(.el-sub-menu__title) {
    color: #000000 !important; // 子菜单标题颜色设为黑色

    &:hover {
      background: rgba(183, 102, 234, 0.08) !important; // hover背景改为Logo颜色
      transform: translateX(3px);
      color: #000000 !important; // hover文字颜色保持黑色

      .el-icon {
        color: #a072ee !important; // hover图标改为Logo颜色
      }
    }
  }

  :deep(.el-sub-menu.is-opened > .el-sub-menu__title) {
    color: #000000 !important; // 打开的子菜单标题颜色设为黑色
    font-weight: 500;

    .el-sub-menu__icon-arrow {
      color: #a072ee; // 箭头颜色改为Logo颜色
      transform: rotateZ(180deg);
    }
  }

  :deep(.el-sub-menu .el-menu-item) {
    padding-left: 56px !important;
    height: 44px;
    line-height: 44px;
    font-size: 14px;
    color: #000000 !important; // 子菜单项文字颜色设为黑色

    .el-icon {
      color: #000000; // 子菜单项图标颜色设为黑色
    }

    &:hover {
      background: rgba(183, 102, 234, 0.08) !important; // hover背景改为Logo颜色
      color: #000000 !important; // hover文字颜色保持黑色

      .el-icon {
        color: #a072ee !important; // hover图标改为Logo颜色
      }
    }

    &.is-active {
      background: linear-gradient(90deg, rgba(183, 102, 234, 0.12) 0%, rgba(183, 102, 234, 0.06) 100%) !important; // 激活背景改为Logo颜色
      color: #a072ee !important; // 激活文字颜色改为Logo颜色

      .el-icon {
        color: #a072ee; // 激活图标颜色改为Logo颜色
      }
    }
  }
}
</style>
