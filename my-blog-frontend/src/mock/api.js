import { categories, tags, articles, comments, series } from './data'

// 模拟延迟
const delay = (ms = 500) => new Promise(resolve => setTimeout(resolve, ms))

// 通用响应格式
const success = (data) => ({
  code: 200,
  message: 'success',
  data
})

const error = (message) => ({
  code: 500,
  message,
  data: null
})

// 获取文章列表
export const getArticles = async (params = {}) => {
  await delay()

  let result = [...articles]

  // 按分类筛选
  if (params.categoryId) {
    result = result.filter(a => a.category.id === params.categoryId)
  }

  // 按标签筛选
  if (params.tag) {
    result = result.filter(a => a.tags.some(t => t.name === params.tag))
  }

  // 搜索
  if (params.keyword) {
    const keyword = params.keyword.toLowerCase()
    result = result.filter(a =>
      a.title.toLowerCase().includes(keyword) ||
      a.summary.toLowerCase().includes(keyword)
    )
  }

  // 排序（最新）
  if (params.sort === 'latest') {
    result.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
  }

  // 排序（热门）
  if (params.sort === 'hot') {
    result.sort((a, b) => b.views - a.views)
  }

  // 分页
  const page = params.page || 1
  const pageSize = params.pageSize || 10
  const start = (page - 1) * pageSize
  const end = start + pageSize

  return success({
    list: result.slice(start, end),
    total: result.length,
    page,
    pageSize
  })
}

// 获取文章详情
export const getArticleDetail = async (id) => {
  await delay()

  const article = articles.find(a => a.id === parseInt(id))
  if (!article) {
    return error('文章不存在')
  }

  // 查找上一篇和下一篇
  const index = articles.findIndex(a => a.id === parseInt(id))
  const prevArticle = index > 0 ? articles[index - 1] : null
  const nextArticle = index < articles.length - 1 ? articles[index + 1] : null

  // 获取评论
  const articleComments = comments.filter(c => c.articleId === parseInt(id))

  return success({
    ...article,
    prevArticle,
    nextArticle,
    comments: articleComments.length
  })
}

// 获取分类列表
export const getCategories = async () => {
  await delay()
  return success(categories)
}

// 获取分类详情
export const getCategoryDetail = async (params) => {
  await delay()

  // 兼容两种调用方式：getCategoryDetail(id) 和 getCategoryDetail({ id, page, pageSize })
  const id = typeof params === 'object' ? params.id : params
  const page = typeof params === 'object' ? (params.page || 1) : 1
  const pageSize = typeof params === 'object' ? (params.pageSize || 12) : 12

  const category = categories.find(c => c.id === parseInt(id))
  if (!category) {
    return error('分类不存在')
  }

  const categoryArticles = articles.filter(a => a.category.id === parseInt(id))

  // 分页
  const start = (page - 1) * pageSize
  const end = start + pageSize

  return success({
    category,
    articles: {
      list: categoryArticles.slice(start, end),
      total: categoryArticles.length,
      page,
      pageSize
    }
  })
}

// 获取标签列表
export const getTags = async () => {
  await delay()
  return success(tags)
}

// 获取标签详情
export const getTagDetail = async (params) => {
  await delay()

  // 兼容两种调用方式：getTagDetail(name) 和 getTagDetail({ name, page, pageSize })
  const name = typeof params === 'object' ? params.name : params
  const page = typeof params === 'object' ? (params.page || 1) : 1
  const pageSize = typeof params === 'object' ? (params.pageSize || 12) : 12

  const tag = tags.find(t => t.name === name)
  if (!tag) {
    return error('标签不存在')
  }

  const tagArticles = articles.filter(a => a.tags.some(t => t.name === name))

  // 分页
  const start = (page - 1) * pageSize
  const end = start + pageSize

  return success({
    tag,
    articles: {
      list: tagArticles.slice(start, end),
      total: tagArticles.length,
      page,
      pageSize
    }
  })
}

// 搜索文章
export const searchArticles = async (params) => {
  await delay()

  // 兼容两种调用方式：searchArticles(keyword) 和 searchArticles({ keyword, page, pageSize })
  const keyword = typeof params === 'string' ? params : (params.keyword || '')
  const page = typeof params === 'object' ? (params.page || 1) : 1
  const pageSize = typeof params === 'object' ? (params.pageSize || 12) : 12

  if (!keyword) {
    return success({
      list: [],
      total: 0,
      page,
      pageSize
    })
  }

  const result = await getArticles({ keyword, page, pageSize })
  return success({
    list: result.data.list || [],
    total: result.data.total || 0,
    page: result.data.page || 1,
    pageSize: result.data.pageSize || 12
  })
}

// 获取热门文章
export const getHotArticles = async () => {
  await delay()
  const hot = [...articles].sort((a, b) => b.views - a.views).slice(0, 5)
  return success(hot)
}

// 获取最新文章
export const getLatestArticles = async () => {
  await delay()
  const latest = [...articles].sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
  return success(latest)
}

// 获取关于信息
export const getAbout = async () => {
  await delay()
  return success({
    name: '博客作者',
    avatar: 'https://i.pravatar.cc/150?img=5',
    bio: '热爱技术，分享知识，持续学习',
    email: 'example@email.com',
    github: 'https://github.com/username',
    website: 'https://example.com',
    skills: ['JavaScript', 'Vue.js', 'React', 'Node.js', 'Python', 'TypeScript'],
    description: `你好！欢迎来到我的博客。

我是一名热爱技术的开发者，专注于前端开发和用户体验设计。

这个博客是我分享技术心得、记录学习笔记、传播有价值内容的平台。

希望我的文章能够帮助到你，也欢迎与我交流讨论！`
  })
}

// 点赞文章
export const likeArticle = async (id) => {
  await delay()
  const article = articles.find(a => a.id === parseInt(id))
  if (article) {
    article.likes += 1
    return success({ likes: article.likes })
  }
  return error('文章不存在')
}

// 取消点赞文章
export const unlikeArticle = async (id) => {
  await delay()
  const article = articles.find(a => a.id === parseInt(id))
  if (article && article.likes > 0) {
    article.likes -= 1
    return success({ likes: article.likes })
  }
  return error('文章不存在')
}

// 收藏文章
export const favoriteArticle = async (id) => {
  await delay()
  const article = articles.find(a => a.id === parseInt(id))
  if (article) {
    article.favorites += 1
    return success({ favorites: article.favorites })
  }
  return error('文章不存在')
}

// 取消收藏文章
export const unfavoriteArticle = async (id) => {
  await delay()
  const article = articles.find(a => a.id === parseInt(id))
  if (article && article.favorites > 0) {
    article.favorites -= 1
    return success({ favorites: article.favorites })
  }
  return error('文章不存在')
}

// 获取文章评论
export const getComments = async (articleId) => {
  await delay()
  const articleComments = comments.filter(c => c.articleId === parseInt(articleId))
  return success(articleComments)
}

// 添加评论
export const addComment = async (data) => {
  await delay()
  const newComment = {
    id: comments.length + 1,
    articleId: parseInt(data.articleId),
    content: data.content,
    author: {
      id: 99,
      name: '当前用户',
      avatar: 'https://i.pravatar.cc/100?img=99'
    },
    createdAt: new Date().toISOString(),
    parentId: data.parentId || null,
    replies: []
  }

  if (data.parentId) {
    // 如果是回复，添加到对应评论的回复中
    const parentComment = comments.find(c => c.id === parseInt(data.parentId))
    if (parentComment) {
      parentComment.replies.push(newComment)
    }
  } else {
    comments.push(newComment)
  }

  return success(newComment)
}

// 删除评论
export const deleteComment = async (id) => {
  await delay()
  const index = comments.findIndex(c => c.id === parseInt(id))
  if (index !== -1) {
    comments.splice(index, 1)
    return success({ message: '删除成功' })
  }
  return error('评论不存在')
}

// 获取专栏列表
export const getSeries = async () => {
  await delay()
  return success(series)
}

// 获取专栏详情
export const getSeriesDetail = async (id) => {
  await delay()
  const seriesData = series.find(s => s.id === parseInt(id))
  if (!seriesData) {
    return error('专栏不存在')
  }

  // 获取章节标题下的子章节对应的文章
  const sectionsWithArticles = seriesData.sections.map(section => ({
    ...section,
    subchapters: section.subchapters.map(subchapter => ({
      ...subchapter,
      articles: subchapter.articleIds.map(articleId => articles.find(a => a.id === articleId)).filter(Boolean)
    }))
  }))

  return success({
    ...seriesData,
    sections: sectionsWithArticles
  })
}

// 获取章节文章
export const getChapterDetail = async (seriesId, chapterId) => {
  await delay()
  const seriesData = series.find(s => s.id === parseInt(seriesId))
  if (!seriesData) {
    return error('专栏不存在')
  }

  // 在所有章节标题的子章节中查找
  let foundSubchapter = null
  for (const section of seriesData.sections) {
    const subchapter = section.subchapters.find(sc => sc.id === parseInt(chapterId))
    if (subchapter) {
      foundSubchapter = subchapter
      break
    }
  }

  if (!foundSubchapter) {
    return error('章节不存在')
  }

  const subchapterArticles = foundSubchapter.articleIds.map(articleId => articles.find(a => a.id === articleId)).filter(Boolean)

  return success({
    ...foundSubchapter,
    articles: subchapterArticles,
    series: seriesData
  })
}

