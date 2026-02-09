// 假数据生成器

const generateId = () => Math.floor(Math.random() * 100000)

const now = new Date()

const categories = [
  { id: 1, name: '前端开发', articleCount: 15 },
  { id: 2, name: '后端开发', articleCount: 10 },
  { id: 3, name: '数据库', articleCount: 8 },
  { id: 4, name: '运维部署', articleCount: 5 },
  { id: 5, name: '设计', articleCount: 6 },
  { id: 6, name: '算法', articleCount: 7 }
]

const tags = [
  { name: 'JavaScript', articleCount: 20 },
  { name: 'Vue.js', articleCount: 15 },
  { name: 'React', articleCount: 12 },
  { name: 'Node.js', articleCount: 10 },
  { name: 'Python', articleCount: 8 },
  { name: 'MySQL', articleCount: 9 },
  { name: 'MongoDB', articleCount: 6 },
  { name: 'Docker', articleCount: 7 },
  { name: 'Git', articleCount: 5 },
  { name: 'Webpack', articleCount: 4 },
  { name: 'TypeScript', articleCount: 8 },
  { name: 'CSS', articleCount: 12 },
  { name: 'HTML', articleCount: 10 },
  { name: 'Linux', articleCount: 6 },
  { name: 'Nginx', articleCount: 5 }
]

const series = [
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
          },
          {
            id: 4,
            name: '生成器',
            description: '生成器与迭代器',
            articleIds: [7, 8]
          }
        ]
      }
    ]
  },
  {
    id: 2,
    name: 'Vue.js系列',
    icon: '💚',
    description: 'Vue框架全方位学习',
    sections: [
      {
        id: 1,
        name: '基础篇',
        description: 'Vue核心概念',
        subchapters: [
          {
            id: 1,
            name: 'Vue3基础',
            description: 'Vue3核心特性',
            articleIds: [9, 10]
          },
          {
            id: 2,
            name: '组合式API',
            description: 'Composition API详解',
            articleIds: [11, 12]
          }
        ]
      },
      {
        id: 2,
        name: '进阶篇',
        description: '高级特性与最佳实践',
        subchapters: [
          {
            id: 3,
            name: '状态管理',
            description: 'Pinia状态管理',
            articleIds: [13, 14]
          },
          {
            id: 4,
            name: '路由管理',
            description: 'Vue Router进阶',
            articleIds: [15, 16]
          }
        ]
      }
    ]
  },
  {
    id: 3,
    name: '算法系列',
    icon: '🧮',
    description: '数据结构与算法',
    sections: [
      {
        id: 1,
        name: '基础算法',
        description: '排序与查找',
        subchapters: [
          {
            id: 1,
            name: '排序算法',
            description: '常见排序算法',
            articleIds: [17, 18]
          },
          {
            id: 2,
            name: '查找算法',
            description: '二分查找等',
            articleIds: [19, 20]
          }
        ]
      },
      {
        id: 2,
        name: '进阶算法',
        description: '动态规划与图论',
        subchapters: [
          {
            id: 3,
            name: '动态规划',
            description: 'DP经典问题',
            articleIds: [21, 22]
          },
          {
            id: 4,
            name: '图论算法',
            description: '图的遍历与最短路',
            articleIds: [23, 24]
          }
        ]
      }
    ]
  }
]

const articles = [
  {
    id: 1,
    title: 'Vue 3 组合式 API 完全指南',
    summary: '深入了解 Vue 3 的组合式 API，包括 setup、ref、reactive、computed 等核心概念的使用方法和最佳实践。',
    content: `# Vue 3 组合式 API 完全指南

Vue 3 引入的组合式 API 是一个全新的编程范式，它提供了更灵活的逻辑复用方式。

## 什么是组合式 API？

组合式 API 是一组基于函数的 API，允许我们使用函数来组合组件的逻辑。

### setup 函数

setup 函数是组合式 API 的入口点。

\`\`\`javascript
import { ref, reactive } from 'vue'

export default {
  setup() {
    const count = ref(0)
    const state = reactive({ name: 'Vue' })

    return { count, state }
  }
}
\`\`\`

### ref 和 reactive

- ref 用于创建响应式的基本类型
- reactive 用于创建响应式对象

### computed 和 watch

computed 用于计算属性，watch 用于监听变化。

## 最佳实践

1. 合理使用 ref 和 reactive
2. 避免过度解构
3. 使用 composables 复用逻辑`,
    cover: 'https://images.unsplash.com/photo-1633356122544-f134324a6cee?w=800',
    category: categories[0],
    tags: [tags[1], tags[0]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 1250,
    likes: 86,
    favorites: 42,
    comments: 15,
    createdAt: new Date(now - 2 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 2 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 2,
    title: 'JavaScript 异步编程深度解析',
    summary: '从回调函数到 Promise，再到 async/await，全面掌握 JavaScript 异步编程的演进历程和最佳实践。',
    content: `# JavaScript 异步编程深度解析

异步编程是 JavaScript 的核心特性之一。

## 回调函数

最早的异步处理方式。

\`\`\`javascript
function fetchData(callback) {
  setTimeout(() => {
    callback('数据')
  }, 1000)
}
\`\`\`

## Promise

Promise 提供了更优雅的异步处理方式。

\`\`\`javascript
const promise = new Promise((resolve, reject) => {
  setTimeout(() => {
    resolve('数据')
  }, 1000)
})
\`\`\`

## async/await

async/await 是基于 Promise 的语法糖。

\`\`\`javascript
async function getData() {
  const data = await promise
  return data
}
\`\`\``,
    cover: 'https://images.unsplash.com/photo-1579468118864-1b9ea3c0db4a?w=800',
    category: categories[0],
    tags: [tags[0]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 980,
    likes: 65,
    favorites: 28,
    comments: 12,
    createdAt: new Date(now - 5 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 5 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 3,
    title: 'Node.js 性能优化实战',
    summary: '分享 Node.js 应用性能优化的实用技巧，包括内存管理、事件循环优化、并发控制等方面。',
    content: `# Node.js 性能优化实战

Node.js 性能优化是后端开发的重要话题。

## 事件循环

理解事件循环是优化 Node.js 性能的基础。

### 宏任务和微任务

- 宏任务：setTimeout、setInterval、I/O
- 微任务：Promise.then、queueMicrotask

## 内存管理

### 内存泄漏

常见的内存泄漏原因：
- 未清理的定时器
- 全局变量
- 事件监听器未移除

## 并发控制

使用 p-limit 等库控制并发数。

\`\`\`javascript
import pLimit from 'p-limit'

const limit = pLimit(10)

const tasks = urls.map(url =>
  limit(() => fetch(url))
)
\`\`\``,
    cover: 'https://images.unsplash.com/photo-1627398242454-45a1465c2479?w=800',
    category: categories[1],
    tags: [tags[3], tags[7]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 856,
    likes: 52,
    favorites: 21,
    comments: 8,
    createdAt: new Date(now - 7 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 7 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 4,
    title: 'MySQL 性能调优指南',
    summary: '深入探讨 MySQL 数据库的性能优化策略，包括索引优化、查询优化、配置调优等实用技巧。',
    content: `# MySQL 性能调优指南

MySQL 性能优化是数据库管理的核心技能。

## 索引优化

### 索引类型

- 主键索引
- 唯一索引
- 普通索引
- 全文索引

### 索引设计原则

1. 选择合适的字段建立索引
2. 避免过多的索引
3. 使用复合索引

## 查询优化

### EXPLAIN 分析

使用 EXPLAIN 分析查询执行计划。

### 慢查询日志

开启慢查询日志，定位性能瓶颈。

## 配置调优

### my.cnf 配置

\`\`\`ini
[mysqld]
innodb_buffer_pool_size = 4G
max_connections = 500
query_cache_size = 256M
\`\`\``,
    cover: 'https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=800',
    category: categories[2],
    tags: [tags[5]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 720,
    likes: 45,
    favorites: 18,
    comments: 6,
    createdAt: new Date(now - 10 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 10 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 5,
    title: 'Docker 容器化部署实践',
    summary: '从零开始学习 Docker，掌握容器化部署的核心概念和实际应用，包括镜像构建、容器编排等。',
    content: `# Docker 容器化部署实践

Docker 是现代应用部署的重要工具。

## Docker 核心概念

### 镜像

镜像是容器的只读模板。

### 容器

容器是镜像的运行实例。

### 仓库

仓库用于存储和分发镜像。

## Dockerfile 编写

\`\`\`dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000
CMD ["npm", "start"]
\`\`\`

## Docker Compose

使用 Docker Compose 管理多容器应用。

\`\`\`yaml
version: '3'
services:
  app:
    build: .
    ports:
      - "3000:3000"
    depends_on:
      - db
  db:
    image: mysql:8
\`\`\``,
    cover: 'https://images.unsplash.com/photo-1605745341112-85968b19335b?w=800',
    category: categories[3],
    tags: [tags[7], tags[13], tags[14]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 645,
    likes: 38,
    favorites: 15,
    comments: 9,
    createdAt: new Date(now - 12 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 12 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 6,
    title: 'React Hooks 最佳实践',
    summary: '全面掌握 React Hooks 的使用方法，包括常用 Hooks 自定义 Hooks 的编写技巧。',
    content: `# React Hooks 最佳实践

Hooks 改变了 React 组件的编写方式。

## 基础 Hooks

### useState

useState 用于管理组件状态。

\`\`\`javascript
const [count, setCount] = useState(0)
\`\`\`

### useEffect

useEffect 用于处理副作用。

\`\`\`javascript
useEffect(() => {
  document.title = \`Count: \${count}\`
}, [count])
\`\`\`

## 自定义 Hooks

自定义 Hooks 用于复用逻辑。

\`\`\`javascript
function useWindowSize() {
  const [size, setSize] = useState({
    width: window.innerWidth,
    height: window.innerHeight
  })

  useEffect(() => {
    const handleResize = () => {
      setSize({
        width: window.innerWidth,
        height: window.innerHeight
      })
    }

    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  return size
}
\`\`\``,
    cover: 'https://images.unsplash.com/photo-1633356122102-3fe601e05e49?w=800',
    category: categories[0],
    tags: [tags[2], tags[10]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 1120,
    likes: 78,
    favorites: 36,
    comments: 14,
    createdAt: new Date(now - 15 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 15 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 7,
    title: 'TypeScript 进阶技巧',
    summary: '深入 TypeScript 高级特性，掌握泛型、类型守卫、条件类型等进阶用法。',
    content: `# TypeScript 进阶技巧

TypeScript 的强大在于其类型系统。

## 泛型

泛型提供了创建可复用组件的能力。

\`\`\`typescript
function identity<T>(arg: T): T {
  return arg
}

const result = identity<string>('hello')
\`\`\`

## 类型守卫

类型守卫用于在运行时确定类型。

\`\`\`typescript
function isString(value: unknown): value is string {
  return typeof value === 'string'
}
\`\`\`

## 条件类型

条件类型根据条件选择类型。

\`\`\`typescript
type NonNullable<T> = T extends null | undefined ? never : T
\`\`\``,
    cover: 'https://images.unsplash.com/photo-1516116216624-53e697fedbea?w=800',
    category: categories[0],
    tags: [tags[10], tags[0]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 920,
    likes: 58,
    favorites: 24,
    comments: 7,
    createdAt: new Date(now - 18 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 18 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 8,
    title: 'Webpack 5 配置优化',
    summary: 'Webpack 5 新特性详解，包括模块联邦、持久化缓存、Tree Shaking 等优化技巧。',
    content: `# Webpack 5 配置优化

Webpack 5 带来了很多新特性。

## 模块联邦

模块联邦允许多个构建之间共享代码。

\`\`\`javascript
const ModuleFederationPlugin = require('webpack/lib/container/ModuleFederationPlugin')

module.exports = {
  plugins: [
    new ModuleFederationPlugin({
      name: 'app1',
      filename: 'remoteEntry.js',
      exposes: {
        './Button': './src/Button'
      }
    })
  ]
}
\`\`\`

## 持久化缓存

Webpack 5 内置了缓存功能。

\`\`\`javascript
module.exports = {
  cache: {
    type: 'filesystem'
  }
}
\`\`\`

## Tree Shaking

Tree Shaking 可以移除未使用的代码。

\`\`\`javascript
module.exports = {
  optimization: {
    usedExports: true,
    sideEffects: false
  }
}
\`\`\``,
    cover: 'https://images.unsplash.com/photo-1507721999472-8ed4421c4af2?w=800',
    category: categories[0],
    tags: [tags[9], tags[0]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 780,
    likes: 48,
    favorites: 22,
    comments: 5,
    createdAt: new Date(now - 20 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 20 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 9,
    title: 'Python 数据分析入门',
    summary: '使用 Python 进行数据分析的基础教程，包括 Pandas、NumPy、Matplotlib 等库的使用。',
    content: `# Python 数据分析入门

Python 是数据分析的首选语言。

## Pandas 基础

Pandas 是 Python 数据分析的核心库。

\`\`\`python
import pandas as pd

# 读取数据
df = pd.read_csv('data.csv')

# 查看数据
print(df.head())
\`\`\`

## NumPy 数组操作

NumPy 提供了高效的数组操作。

\`\`\`python
import numpy as np

arr = np.array([1, 2, 3])
print(arr.mean())
\`\`\`

## 数据可视化

使用 Matplotlib 绘制图表。

\`\`\`python
import matplotlib.pyplot as plt

plt.plot([1, 2, 3, 4])
plt.show()
\`\`\``,
    cover: 'https://images.unsplash.com/photo-1526379095098-d400fd0bf935?w=800',
    category: categories[1],
    tags: [tags[4]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 680,
    likes: 42,
    favorites: 19,
    comments: 4,
    createdAt: new Date(now - 22 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 22 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 10,
    title: 'MongoDB 集群部署方案',
    summary: 'MongoDB 集群部署的完整指南，包括副本集、分片集群的配置和管理。',
    content: `# MongoDB 集群部署方案

MongoDB 提供了强大的集群功能。

## 副本集

副本集提供数据冗余和高可用性。

\`\`\`javascript
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "mongodb1:27017" },
    { _id: 1, host: "mongodb2:27017" }
  ]
})
\`\`\`

## 分片集群

分片集群支持大规模数据存储。

## 性能优化

使用索引提升查询性能。

\`\`\`javascript
db.collection.createIndex({ "field": 1 })
\`\`\``,
    cover: 'https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=800',
    category: categories[2],
    tags: [tags[6]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 590,
    likes: 36,
    favorites: 14,
    comments: 3,
    createdAt: new Date(now - 25 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 25 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 11,
    title: 'Nginx 反向代理配置',
    summary: 'Nginx 反向代理的详细配置教程，包括负载均衡、SSL 配置、缓存优化等。',
    content: `# Nginx 反向代理配置

Nginx 是高性能的 Web 服务器和反向代理。

## 基础反向代理

\`\`\`nginx
server {
    listen 80;
    server_name example.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
    }
}
\`\`\`

## 负载均衡

\`\`\`nginx
upstream backend {
    server localhost:3000;
    server localhost:3001;
    server localhost:3002;
}

server {
    location / {
        proxy_pass http://backend;
    }
}
\`\`\`

## SSL 配置

\`\`\`nginx
server {
    listen 443 ssl;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
}
\`\`\``,
    cover: 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=800',
    category: categories[3],
    tags: [tags[14]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 520,
    likes: 31,
    favorites: 12,
    comments: 5,
    createdAt: new Date(now - 28 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 28 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 12,
    title: 'UI 设计原则与实践',
    summary: '深入探讨用户界面设计的核心原则，包括一致性、可用性、反馈机制等。',
    content: `# UI 设计原则与实践

优秀的 UI 设计遵循一定的原则。

## 一致性

保持界面元素的一致性。

## 可用性

让用户能够轻松完成目标任务。

## 反馈机制

及时给予用户操作反馈。

## 视觉层次

通过大小、颜色、位置建立清晰的视觉层次。

## 留白

合理使用留白，提升可读性。

## 色彩理论

了解色彩理论和配色方案。`,
    cover: 'https://images.unsplash.com/photo-1561070791-2526d30994b5?w=800',
    category: categories[4],
    tags: [tags[11]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 750,
    likes: 49,
    favorites: 20,
    comments: 8,
    createdAt: new Date(now - 30 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 30 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 13,
    title: '快速排序算法详解',
    summary: '详细讲解快速排序算法的原理、实现和时间复杂度分析。',
    content: `# 快速排序算法详解

快速排序是一种高效的排序算法。

## 算法原理

1. 选择一个基准元素
2. 将数组分为两部分
3. 递归排序两部分

## 代码实现

\`\`\`javascript
function quickSort(arr) {
  if (arr.length <= 1) return arr

  const pivot = arr[0]
  const left = []
  const right = []

  for (let i = 1; i < arr.length; i++) {
    if (arr[i] < pivot) {
      left.push(arr[i])
    } else {
      right.push(arr[i])
    }
  }

  return [...quickSort(left), pivot, ...quickSort(right)]
}
\`\`\`

## 时间复杂度

- 平均：O(n log n)
- 最坏：O(n^2)
- 空间：O(log n)`,
    cover: 'https://images.unsplash.com/photo-1509228468518-180dd4864904?w=800',
    category: categories[5],
    tags: [tags[8]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 610,
    likes: 38,
    favorites: 16,
    comments: 6,
    createdAt: new Date(now - 33 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 33 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 14,
    title: 'CSS Grid 布局指南',
    summary: '全面掌握 CSS Grid 布局系统，创建复杂的二维布局。',
    content: `# CSS Grid 布局指南

Grid 是强大的二维布局系统。

## 基础语法

\`\`\`css
.container {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 20px;
}
\`\`\`

## Grid 属性

- grid-template-columns
- grid-template-rows
- grid-gap
- grid-area

## 响应式布局

\`\`\`css
.container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
}
\`\`\``,
    cover: 'https://images.unsplash.com/photo-1507721999472-8ed4421c4af2?w=800',
    category: categories[0],
    tags: [tags[11]],
    author: { id: 1, name: '技术博主', avatar: 'https://i.pravatar.cc/100?img=2', bio: '专注前端技术分享' },
    views: 880,
    likes: 56,
    favorites: 23,
    comments: 7,
    createdAt: new Date(now - 35 * 24 * 60 * 60 * 1000).toISOString(),
    updatedAt: new Date(now - 35 * 24 * 60 * 60 * 1000).toISOString()
  }
]

const comments = [
  {
    id: 1,
    articleId: 1,
    content: '这篇文章写得非常好，对组合式 API 讲解很清晰！',
    author: {
      id: 1,
      name: '张三',
      avatar: 'https://i.pravatar.cc/100?img=1'
    },
    createdAt: new Date(now - 1 * 24 * 60 * 60 * 1000).toISOString(),
    parentId: null,
    replies: [
      {
        id: 2,
        articleId: 1,
        content: '感谢支持！',
        author: {
          id: 2,
          name: '博主',
          avatar: 'https://i.pravatar.cc/100?img=2'
        },
        createdAt: new Date(now - 1 * 24 * 60 * 60 * 1000 + 3600000).toISOString(),
        parentId: 1,
        replies: []
      }
    ]
  },
  {
    id: 3,
    articleId: 1,
    content: 'setup 函数和 Options API 可以混用吗？',
    author: {
      id: 3,
      name: '李四',
      avatar: 'https://i.pravatar.cc/100?img=3'
    },
    createdAt: new Date(now - 2 * 24 * 60 * 60 * 1000).toISOString(),
    parentId: null,
    replies: []
  },
  {
    id: 4,
    articleId: 2,
    content: 'Promise.all 和 Promise.race 有什么区别？',
    author: {
      id: 4,
      name: '王五',
      avatar: 'https://i.pravatar.cc/100?img=4'
    },
    createdAt: new Date(now - 3 * 24 * 60 * 60 * 1000).toISOString(),
    parentId: null,
    replies: []
  }
]

export { categories, tags, articles, comments, series }
