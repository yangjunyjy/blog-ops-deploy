# 代码块显示问题修复总结

## 问题描述

1. 代码块背景没有变成 IDE 中的黑色背景效果
2. 代码块内嵌代码块内容还是存在

## 修复内容

### 1. 更换代码高亮主题

**文件**: `src/utils/format.js`

```javascript
// 从浅色主题改为深色主题
import 'highlight.js/styles/atom-one-dark.css'  // 替换 github.css
```

### 2. 修复 marked.js 渲染逻辑

**文件**: `src/utils/format.js`

使用自定义渲染器替代已弃用的 `highlight` 选项：

```javascript
const renderer = new Renderer()
renderer.code = function(code, language) {
  const validLang = language && hljs.getLanguage(language)
  const highlighted = validLang
    ? hljs.highlight(code, { language }).value
    : hljs.highlightAuto(code).value

  return `<pre><code class="hljs language-${language || 'plaintext'}">${highlighted}</code></pre>`
}

marked.setOptions({
  renderer: renderer,
  breaks: true,
  gfm: true
})
```

### 3. 更新代码块样式

**文件**: `src/views/ArticleDetail.vue`

```css
/* 代码块样式 - 使用 Atom One Dark 主题 */
.article-content :deep(pre) {
  position: relative;
  padding: 20px 24px !important;
  margin: 28px 0 !important;
  border: none !important;
  background: #282c34 !important;  /* IDE 深色背景 */
  border-radius: 12px;
  overflow-x: auto;
  border: 1px solid #3e4451;  /* 边框颜色 */
}

.article-content :deep(pre code.hljs) {
  background: transparent !important;
  padding: 0 !important;
  font-size: 14px;
  line-height: 1.6;
  font-family: 'Fira Code', 'Monaco', 'Consolas', monospace !important;
}

/* 确保代码块内的代码不再被样式化 */
.article-content :deep(pre code.hljs *) {
  background: transparent !important;
  padding: 0 !important;
  border: none !important;
  font-family: inherit !important;
}

/* 行内代码样式 - 排除代码块内的代码 */
.article-content :deep(code):not(.hljs) {
  background: #f1f5f9;
  color: #e53e3e;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Fira Code', 'Monaco', 'Consolas', monospace;
  font-size: 14px;
}

html.dark .article-content :deep(code):not(.hljs) {
  background: #3d4759;
  color: #f87171;
}
```

### 4. 清理重复样式

删除了文件中的重复样式定义，特别是：

- 删除了文件末尾的旧代码块样式（使用浅色背景的版本）
- 删除了媒体查询中重复的 `pre` 样式

## 效果

### 代码块背景

- **背景色**: `#282c34` (Atom One Dark 深色)
- **文字颜色**: `#abb2bf` (默认高亮颜色)
- **边框**: `#3e4451`
- **圆角**: 12px

### 代码高亮

- 根据编程语言自动高亮（JavaScript, Python, Go, Java 等）
- 使用 highlight.js 的语法高亮规则
- 支持 190+ 种编程语言

### 解决嵌套问题

- 使用 `:not(.hljs)` 选择器排除代码块内的代码
- 使用 `pre code.hljs *` 强制代码块内所有元素继承样式
- 行内代码只匹配直接子元素：`p > code`, `li > code`

## 测试要点

- [ ] 代码块显示深色背景（类似 IDEA）
- [ ] 代码根据语言正确高亮
- [ ] 代码块内没有嵌套的代码块样式
- [ ] 行内代码保持原有样式
- [ ] 复制按钮正常显示和工作
- [ ] 响应式布局正常

## 注意事项

1. **复制按钮**: 位于代码块右上角，悬停时显示
2. **字体**: 使用 'Fira Code' 等等宽字体
3. **滚动**: 代码块横向滚动支持长代码行
4. **主题**: Atom One Dark 是最接近 IDEA Dark 的主题之一
