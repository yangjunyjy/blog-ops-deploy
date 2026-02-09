# 分页器组件使用文档

## 📦 组件位置
`src/components/Pagination.vue`

## 🎯 组件特性

### 设计亮点
- **美观**: 采用蓝紫色渐变主题，与整体设计风格统一
- **动画**: 悬停时上浮、缩放等流畅动画效果
- **响应式**: 完美适配桌面端和移动端
- **信息完整**: 显示当前页、总页数、总条数
- **暗黑模式**: 支持暗黑主题

### 交互体验
- 禁用状态清晰可见
- 激活状态有脉冲动画
- 悬停时有上浮和阴影效果
- 点击后自动滚动到页面顶部

## 📖 使用方法

### 基础用法

```vue
<template>
  <Pagination
    :current-page="currentPage"
    :total="total"
    :page-size="pageSize"
    @change="handlePageChange"
  />
</template>

<script setup>
import { ref } from 'vue'
import Pagination from '@/components/Pagination.vue'

const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

const handlePageChange = (page) => {
  currentPage.value = page
  // 重新加载数据
  loadData()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>
```

### Props 参数

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| currentPage | Number | 是 | - | 当前页码 |
| total | Number | 是 | - | 总数据条数 |
| pageSize | Number | 否 | 12 | 每页显示条数 |
| layout | String | 否 | 'prev, pager, next' | 布局方式 |

### Events 事件

| 事件名 | 参数 | 说明 |
|--------|------|------|
| change | (page: Number) | 页码改变时触发 |

## 🎨 样式定制

### 主题色
组件使用 CSS 变量，如需修改主题色，可在组件中修改:

```css
/* 渐变色背景 */
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

/* 边框色 */
border-color: rgba(102, 126, 234, 0.1);

/* 阴影 */
box-shadow: 0 4px 20px rgba(102, 126, 234, 0.15);
```

### 尺寸调整
如需调整按钮大小，修改以下数值:

```css
.pagination-btn {
  width: 40px;
  height: 40px;
}

.pagination-number {
  min-width: 40px;
  height: 40px;
  font-size: 14px;
}
```

## 📱 响应式断点

### 桌面端 (> 768px)
- 按钮尺寸: 40x40px
- 信息横排显示
- 正常间距

### 移动端 (≤ 768px)
- 按钮尺寸: 36x36px
- 信息竖排显示
- 压缩间距
- 字号缩小

## 🌓 暗黑模式支持

组件自动检测 `html.dark` 类并应用暗黑样式:

```css
html.dark .pagination-container {
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
  border-color: rgba(102, 126, 234, 0.2);
}
```

## 🔄 已集成的页面

目前已在以下页面使用:

1. **Articles.vue** - 文章列表页
2. **Search.vue** - 搜索页
3. **Category.vue** - 分类详情页
4. **Tag.vue** - 标签详情页

## 💡 使用建议

### 1. 配合 computed 计算总页数
```javascript
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))
```

### 2. 切换页码时重置搜索条件
```javascript
const handlePageChange = (page) => {
  currentPage.value = page
  // 如果需要，重置其他条件
  // searchKeyword.value = ''
  loadData()
}
```

### 3. 切换筛选条件时重置页码
```javascript
const handleCategoryChange = (categoryId) => {
  selectedCategory.value = categoryId
  currentPage.value = 1  // 重置到第一页
  loadData()
}
```

## 🎯 最佳实践

1. **条件显示**: 只在 `total > 0` 时显示分页器
   ```vue
   <Pagination v-if="total > 0" ... />
   ```

2. **平滑滚动**: 页面切换后滚动到顶部
   ```javascript
   window.scrollTo({ top: 0, behavior: 'smooth' })
   ```

3. **加载状态**: 切换页码时显示加载状态
   ```javascript
   const handlePageChange = async (page) => {
     loading.value = true
     currentPage.value = page
     await loadData()
     loading.value = false
     window.scrollTo({ top: 0, behavior: 'smooth' })
   }
   ```

## 🐛 常见问题

### Q: 如何禁用上一页/下一页按钮?
A: 组件会根据 `currentPage` 和 `totalPages` 自动计算禁用状态，无需手动处理。

### Q: 如何自定义按钮样式?
A: 可以通过深度选择器覆盖样式:
```vue
<style scoped>
.pagination :deep(.pagination-btn) {
  /* 自定义样式 */
}
</style>
```

### Q: 如何实现跳转到指定页码?
A: 组件目前只支持点击页码按钮，如需跳转功能可以扩展组件。

## 📦 依赖

该组件不依赖任何第三方库，纯 Vue 3 + CSS 实现。

## 🎬 动画效果

- **悬停**: `transform: translateY(-2px)` + `scale(1.05)`
- **激活**: `scale(1.1)` + 脉冲动画 `pulse`
- **禁用**: `opacity: 0.4` + `cursor: not-allowed`

所有动画都有 `transition: all 0.3s` 平滑过渡。
