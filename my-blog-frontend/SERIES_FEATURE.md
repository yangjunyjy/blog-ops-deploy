# 专栏功能使用文档

## 📚 功能概述

专栏功能是博客系统的一个特色模块,用于组织系列文章的学习路径。

### 主要特性
- **专栏导航**: 顶部导航栏新增"专栏"入口
- **专栏详情**: 展示专栏信息和章节列表
- **章节文章**: 每个章节包含多篇文章
- **新窗口阅读**: 点击文章在新窗口打开阅读

## 🎯 功能结构

### 1. 顶部导航

**位置**: `Header.vue`

**实现方式**:
- 悬停"专栏"链接显示下拉菜单
- 列出所有可用专栏
- 每个专栏显示: 图标、名称、描述

**代码位置**:
```vue
<div class="series-nav" @mouseenter="seriesDropdownVisible = true">
  <span>专栏</span>
  <div class="series-dropdown" v-show="seriesDropdownVisible">
    <!-- 专栏列表 -->
  </div>
</div>
```

### 2. 专栏详情页

**路由**: `/series/:id`

**组件**: `src/views/SeriesDetail.vue`

**页面布局**:
```
┌─────────────────────────────────┐
│  专栏头部                      │
│  ├─ 专栏图标                  │
│  ├─ 专栏名称                  │
│  ├─ 描述信息                  │
│  └─ 统计数据                  │
├─────────────────────────────────┤
│  章节列表                      │
│  ├─ 章节1 (激活)              │
│  ├─ 章节2                    │
│  └─ 章节3                    │
├─────────────────────────────────┤
│  文章卡片 (当前章节)            │
│  ├─ 文章1                    │
│  ├─ 文章2                    │
│  └─ 文章3                    │
└─────────────────────────────────┘
```

### 3. 章节文章阅读页

**路由**: `/series/:id/article/:articleId`

**组件**: `src/views/ChapterArticle.vue`

**页面布局**:
```
┌─────────────────────────────────┐
│  文章头部                      │
│  ├─ 标题                    │
│  ├─ 元数据 (分类、日期、浏览)    │
│  └─ 操作按钮 (点赞、收藏)       │
├─────────────────────────────────┤
│  文章内容 (Markdown渲染)        │
│  ├─ 标题、段落、代码块         │
│  └─ 样式美化                  │
├─────────────────────────────────┤
│  上一篇/下一篇导航              │
├─────────────────────────────────┤
│  评论区                        │
└─────────────────────────────────┘
```

## 📊 数据结构

### 专栏数据 (Series)
```javascript
{
  id: 1,
  name: 'Python基础知识',
  icon: '🐍',
  description: '从入门到精通的Python学习路线',
  sections: [
    {
      id: 1,
      name: '基础篇',
      description: 'Python入门必备知识',
      subchapters: [
        {
          id: 1,
          name: 'Python3面向对象',
          description: '面向对象编程详解',
          articleIds: [1, 2]
        },
        {
          id: 2,
          name: '数据结构',
          description: 'Python内置数据结构',
          articleIds: [3, 4]
        }
      ]
    },
    {
      id: 2,
      name: '进阶篇',
      description: '进阶编程技巧',
      subchapters: [
        {
          id: 3,
          name: '装饰器',
          description: 'Python装饰器详解',
          articleIds: [5, 6]
        }
      ]
    }
  ]
}
```

### 章节标题数据 (Section)
```javascript
{
  id: 1,
  name: '基础篇',
  description: 'Python入门必备知识',
  subchapters: [
    {
      id: 1,
      name: 'Python3面向对象',
      description: '面向对象编程详解',
      articleIds: [1, 2]
    }
  ]
}
```

### 子章节数据 (Subchapter)
```javascript
{
  id: 1,
  name: 'Python3面向对象',
  description: '面向对象编程详解',
  articleIds: [1, 2],
  articles: [
    {
      id: 1,
      title: 'Python面向对象编程详解',
      summary: '详细讲解Python的面向对象特性',
      content: '# Python面向对象...',
      cover: 'https://example.com/cover.jpg',
      views: 100,
      createdAt: '2024-01-01T00:00:00Z'
    }
  ]
}
```

## 🔌 API接口

### 获取专栏列表
```javascript
GET /api/series

Response:
{
  code: 200,
  message: 'success',
  data: [
    {
      id: 1,
      name: 'Python基础知识',
      icon: '🐍',
      description: '从入门到精通的Python学习路线',
      sections: [
        {
          id: 1,
          name: '基础篇',
          description: 'Python入门必备知识',
          subchapters: [...]
        }
      ]
    }
  ]
}
```

### 获取专栏详情
```javascript
GET /api/series/:id

Response:
{
  code: 200,
  message: 'success',
  data: {
    id: 1,
    name: 'Python基础知识',
    icon: '🐍',
    description: '从入门到精通的Python学习路线',
    sections: [
      {
        id: 1,
        name: '基础篇',
        description: 'Python入门必备知识',
        subchapters: [
          {
            id: 1,
            name: 'Python3面向对象',
            description: '面向对象编程详解',
            articleIds: [1, 2],
            articles: [...]
          }
        ]
      }
    ]
  }
}
```

### 获取子章节文章
```javascript
GET /api/series/:seriesId/chapters/:chapterId

Response:
{
  code: 200,
  message: 'success',
  data: {
    id: 1,
    name: 'Python3面向对象',
    description: '面向对象编程详解',
    articleIds: [1, 2],
    articles: [
      {
        id: 1,
        title: 'Python面向对象编程详解',
        summary: '详细讲解Python的面向对象特性',
        content: '# Python面向对象...',
        // ...其他字段
      }
    ],
    series: {
      id: 1,
      name: 'Python基础知识',
      icon: '🐍'
    }
  }
}
```

## 🎨 样式设计

### 配色方案
- **主色调**: 蓝紫色渐变 `#667eea` → `#764ba2`
- **背景色**: 白色卡片 + 渐变背景
- **暗黑模式**: 深色卡片 + 对应背景

### 卡片设计
- **圆角**: 12px - 20px
- **阴影**: 柔和的投影效果
- **悬停**: 上浮动画 + 阴影增强
- **激活**: 渐变边框 + 高亮背景

### 动画效果
```css
/* 章节卡片悬停 */
.chapter-item:hover {
  transform: translateX(8px);
  border-color: rgba(102, 126, 234, 0.3);
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.2);
}

/* 文章卡片悬停 */
.chapter-article-card:hover {
  transform: translateY(-4px);
  border-color: rgba(102, 126, 234, 0.3);
}
```

## 📱 响应式设计

### 桌面端 (> 768px)
- 专栏头部: 横向布局
- 章节列表: 完整展示
- 文章网格: 多列布局

### 移动端 (≤ 768px)
- 专栏头部: 纵向布局
- 章节列表: 压缩间距
- 文章网格: 单列布局

## 🚀 使用示例

### 1. 创建新专栏
在 `src/mock/data.js` 中添加:
```javascript
const series = [
  {
    id: 4,
    name: 'JavaScript系列',
    icon: '💛',
    description: 'JavaScript深入学习',
    sections: [
      {
        id: 1,
        name: '基础篇',
        description: 'JS核心概念',
        subchapters: [
          {
            id: 1,
            name: 'JS基础语法',
            description: 'JavaScript基础语法',
            articleIds: [10, 11, 12]
          }
        ]
      }
    ]
  }
]
```

### 2. 添加章节文章
确保文章ID存在于 `articles` 数组中:
```javascript
const articles = [
  {
    id: 10,
    title: 'JS基础语法',
    summary: 'JavaScript基础语法学习',
    content: '# JS基础语法...',
    // ...
  }
]
```

### 3. 浏览专栏
1. 点击顶部导航"专栏"
2. 悬停选择目标专栏
3. 默认显示第一个章节
4. 点击章节切换
5. 点击文章卡片在新窗口打开

## 🎯 交互流程

```
用户操作:
  1. 鼠标悬停顶部导航"专栏"
  2. 点击下拉菜单中的专栏
  3. 进入专栏详情页
  4. 点击章节切换
  5. 点击文章卡片
  6. 在新窗口打开文章阅读页

系统响应:
  1. 显示专栏下拉菜单
  2. 跳转到专栏详情页
  3. 加载章节文章列表
  4. 更新右侧文章列表
  5. 打开新窗口显示文章
  6. 渲染Markdown内容
```

## 💡 扩展建议

### 1. 学习进度追踪
```javascript
// 记录用户阅读进度
const readingProgress = {
  seriesId: 1,
  sectionId: 1,
  subchapterId: 1,
  completedArticles: [1, 2]
}
```

### 2. 章节导航优化
- 添加"上一章/下一章"导航
- 标记已读章节
- 显示阅读进度条

### 3. 文章关联推荐
- 在文章底部推荐同专栏其他文章
- 显示"继续学习"按钮

## 🔧 开发注意事项

1. **数据结构**: 专栏包含 sections(章节标题),section 包含 subchapters(子章节),subchapter 包含 articles(文章)
2. **Mock数据**: 当前使用Mock数据,需要替换为真实API
3. **文章ID关联**: 确保 `articleIds` 与 `articles` 数组中的ID匹配
4. **路由参数**: 系列ID和子章节ID都是动态路由参数
5. **新窗口打开**: 使用 `window.open()` 在新标签页打开文章
6. **Markdown渲染**: 使用 `marked` 库渲染Markdown内容
7. **子章节定位**: 在 `series.sections` 数组中遍历查找对应的子章节

## 📦 依赖包

```json
{
  "marked": "^17.0.1",
  "@element-plus/icons-vue": "^2.3.2"
}
```

## 🐛 已知问题

- [ ] Mock数据中文章ID与articleIds需要手动关联
- [ ] 暂不支持学习进度保存
- [ ] 暂不支持章节排序

## ✨ 未来优化

- [ ] 添加专栏搜索功能
- [ ] 支持专栏收藏
- [ ] 添加专栏阅读统计
- [ ] 支持章节自定义排序
- [ ] 添加学习笔记功能
