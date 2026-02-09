/**
 * 博客 Mock 数据
 */

const mockImages = [
  'https://picsum.photos/800/400?random=1',
  'https://picsum.photos/800/400?random=2',
  'https://picsum.photos/800/400?random=3',
  'https://picsum.photos/800/400?random=4',
  'https://picsum.photos/800/400?random=5'
]

const mockAvatars = [
  'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
  'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
  'https://cube.elemecdn.com/9/c2/f0ee8a3c7c9638a54940382568c9dpng.png'
]

// 用户数据
export const users = [
  {
    id: 1,
    username: 'admin',
    email: 'admin@example.com',
    nickname: '管理员',
    avatar: mockAvatars[0],
    role: 1,
    status: 1,
    createdAt: '2024-01-01T00:00:00Z'
  },
  {
    id: 2,
    username: 'editor',
    email: 'editor@example.com',
    nickname: '编辑',
    avatar: mockAvatars[1],
    role: 0,
    status: 1,
    createdAt: '2024-01-15T00:00:00Z'
  }
]

// 分类数据
export const categories = [
  { id: 1, name: '前端开发', slug: 'frontend', icon: '💻', sortOrder: 1, status: 1, articleCount: 12 },
  { id: 2, name: '后端开发', slug: 'backend', icon: '🔧', sortOrder: 2, status: 1, articleCount: 8 },
  { id: 3, name: '人工智能', slug: 'ai', icon: '🤖', sortOrder: 3, status: 1, articleCount: 5 },
  { id: 4, name: '数据库', slug: 'database', icon: '🗄️', sortOrder: 4, status: 1, articleCount: 6 },
  { id: 5, name: '运维部署', slug: 'devops', icon: '🚀', sortOrder: 5, status: 1, articleCount: 4 }
]

// 标签数据
export const tags = [
  { id: 1, name: 'Vue.js', slug: 'vuejs', articleCount: 8 },
  { id: 2, name: 'React', slug: 'react', articleCount: 6 },
  { id: 3, name: 'JavaScript', slug: 'javascript', articleCount: 10 },
  { id: 4, name: 'TypeScript', slug: 'typescript', articleCount: 7 },
  { id: 5, name: 'Node.js', slug: 'nodejs', articleCount: 5 },
  { id: 6, name: 'Python', slug: 'python', articleCount: 4 },
  { id: 7, name: 'Java', slug: 'java', articleCount: 6 },
  { id: 8, name: 'Spring Boot', slug: 'spring-boot', articleCount: 4 },
  { id: 9, name: 'MySQL', slug: 'mysql', articleCount: 5 },
  { id: 10, name: 'Docker', slug: 'docker', articleCount: 4 }
]

// 系列数据（四层结构）
export const series = [
  {
    id: 1,
    name: 'Vue3 源码解析',
    slug: 'vue3-source',
    icon: '📚',
    description: '深入解析 Vue3 源码',
    cover: mockImages[0],
    sortOrder: 1,
    status: 1,
    createdAt: '2024-01-01T00:00:00Z',
    sections: [
      {
        id: 1,
        seriesId: 1,
        name: '响应式系统',
        description: 'Vue3 的响应式原理',
        sortOrder: 1,
        subchapters: [
          {
            id: 1,
            sectionId: 1,
            name: 'Proxy 代理',
            description: 'Proxy 实现响应式',
            sortOrder: 1,
            articleIds: [1, 2]
          }
        ]
      }
    ]
  }
]

// 文章数据
export const articles = [
  {
    id: 1,
    title: 'Vue3 响应式系统：Proxy 代理原理详解',
    slug: 'vue3-reactive-proxy',
    summary: '深入解析 Vue3 的响应式系统',
    content: '# Vue3 响应式系统\n\nProxy 是 JavaScript 中强大的元编程特性...',
    cover: mockImages[0],
    categoryId: 1,
    category: { id: 1, name: '前端开发' },
    authorId: 1,
    author: { id: 1, username: 'admin', nickname: '管理员', avatar: mockAvatars[0] },
    tagIds: [1, 3],
    tags: [{ id: 1, name: 'Vue.js' }, { id: 3, name: 'JavaScript' }],
    seriesId: 1,
    views: 1234,
    likes: 56,
    favorites: 23,
    commentCount: 12,
    status: 1,
    isTop: true,
    isRecommended: true,
    keywords: 'Vue3,响应式',
    publishedAt: '2024-01-15T10:00:00Z',
    createdAt: '2024-01-15T00:00:00Z'
  },
  {
    id: 2,
    title: 'TypeScript 泛型编程完全指南',
    slug: 'typescript-generics',
    summary: '全面介绍 TypeScript 泛型编程',
    content: '# TypeScript 泛型编程\n\n泛型是 TypeScript 中强大的类型工具...',
    cover: mockImages[1],
    categoryId: 1,
    category: { id: 1, name: '前端开发' },
    authorId: 1,
    author: { id: 1, username: 'admin', nickname: '管理员', avatar: mockAvatars[0] },
    tagIds: [4],
    tags: [{ id: 4, name: 'TypeScript' }],
    seriesId: null,
    views: 987,
    likes: 45,
    favorites: 18,
    commentCount: 8,
    status: 1,
    isTop: false,
    isRecommended: true,
    keywords: 'TypeScript,泛型',
    publishedAt: '2024-02-10T10:00:00Z',
    createdAt: '2024-02-10T00:00:00Z'
  },
  {
    id: 3,
    title: 'Spring Boot 3.0 新特性',
    slug: 'spring-boot-3-features',
    summary: 'Spring Boot 3.0 带来了许多新特性',
    content: '# Spring Boot 3.0 新特性\n\n基于 Jakarta EE...',
    cover: mockImages[2],
    categoryId: 2,
    category: { id: 2, name: '后端开发' },
    authorId: 1,
    author: { id: 1, username: 'admin', nickname: '管理员', avatar: mockAvatars[0] },
    tagIds: [7, 8],
    tags: [{ id: 7, name: 'Java' }, { id: 8, name: 'Spring Boot' }],
    seriesId: null,
    views: 756,
    likes: 34,
    favorites: 12,
    commentCount: 6,
    status: 0,
    isTop: false,
    isRecommended: false,
    keywords: 'Spring Boot,Java',
    publishedAt: null,
    createdAt: '2024-03-01T00:00:00Z'
  }
]

// 评论数据
export const comments = [
  {
    id: 1,
    articleId: 1,
    article: { id: 1, title: 'Vue3 响应式系统：Proxy 代理原理详解' },
    userId: 2,
    user: { id: 2, username: 'editor', nickname: '编辑', avatar: mockAvatars[1] },
    content: '文章写得很好，学到了很多！',
    status: 1,
    createdAt: '2024-01-16T10:30:00Z'
  },
  {
    id: 2,
    articleId: 1,
    article: { id: 1, title: 'Vue3 响应式系统：Proxy 代理原理详解' },
    userId: 2,
    user: { id: 2, username: 'editor', nickname: '编辑', avatar: mockAvatars[1] },
    content: 'Proxy 确实比 Object.defineProperty 强很多',
    status: 1,
    createdAt: '2024-01-17T14:20:00Z'
  },
  {
    id: 3,
    articleId: 2,
    article: { id: 2, title: 'TypeScript 泛型编程完全指南' },
    userId: 2,
    user: { id: 2, username: 'editor', nickname: '编辑', avatar: mockAvatars[1] },
    content: '泛型确实强大，希望能多出一些高级用法的文章',
    status: 0,
    createdAt: '2024-02-11T09:00:00Z'
  }
]
