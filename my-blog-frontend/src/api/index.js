import request from '@/utils/request'
import * as mockApi from '@/mock/api'

// 使用 Mock 数据的开关
const USE_MOCK = false

// API 路径前缀
const API_PREFIX = '/api/v1'

// ============= Public API (公开接口，无需认证) =============


// 获取文章列表
export const getArticles = (params) => {
  if (USE_MOCK) return mockApi.getArticles(params)
  return request({
    url: `${API_PREFIX}/public/articles`,
    method: 'get',
    params
  })
}

// 获取文章详情
export const getArticleDetail = (id) => {
  if (USE_MOCK) return mockApi.getArticleDetail(id)
  return request({
    url: `${API_PREFIX}/public/articles/${id}`,
    method: 'get'
  })
}

// 获取分类列表
export const getCategories = (params = {}) => {
  if (USE_MOCK) return mockApi.getCategories()
  return request({
    url: `${API_PREFIX}/public/categories`,
    method: 'get',
    params
  })
}

// 获取分类详情
export const getCategoryDetail = (params) => {
  if (USE_MOCK) return mockApi.getCategoryDetail(params)
  return request({
    url: `${API_PREFIX}/public/categories/${typeof params === 'object' ? params.id : params}`,
    method: 'get',
    params: typeof params === 'object' ? {
      page: params.page,
      pageSize: params.pageSize
    } : undefined
  })
}

// 获取标签列表
export const getTags = (params = {}) => {
  if (USE_MOCK) return mockApi.getTags()
  return request({
    url: `${API_PREFIX}/public/tags`,
    method: 'get',
    params
  })
}

// 获取标签详情
export const getTagDetail = (params) => {
  if (USE_MOCK) return mockApi.getTagDetail(params)
  return request({
    url: `${API_PREFIX}/public/tags/${typeof params === 'object' ? params.id : params}`,
    method: 'get',
    params: typeof params === 'object' ? {
      page: params.page,
      pageSize: params.pageSize
    } : undefined
  })
}

// 搜索文章
export const searchArticles = (params) => {
  if (USE_MOCK) return mockApi.searchArticles(params)
  return request({
    url: `${API_PREFIX}/public/articles/search`,
    method: 'get',
    params
  })
}

// 获取关于信息
export const getAbout = () => {
  if (USE_MOCK) return mockApi.getAbout()
  return request({
    url: `${API_PREFIX}/public/about`,
    method: 'get'
  })
}

// 获取热门文章
export const getHotArticles = () => {
  if (USE_MOCK) return mockApi.getHotArticles()
  return request({
    url: `${API_PREFIX}/public/articles/hot`,
    method: 'get'
  })
}

// 获取最新文章
export const getLatestArticles = () => {
  if (USE_MOCK) return mockApi.getLatestArticles()
  return request({
    url: `${API_PREFIX}/public/articles/recent`,
    method: 'get'
  })
}

// 获取文章评论
export const getComments = (articleId) => {
  if (USE_MOCK) return mockApi.getComments(articleId)
  return request({
    url: `${API_PREFIX}/public/comments/article/${articleId}`,
    method: 'get'
  })
}

// 获取子评论列表
export const getCommentReplies = (commentId, page = 1, pageSize = 10) => {
  return request({
    url: `${API_PREFIX}/public/comments/${commentId}/replies`,
    method: 'get',
    params: {
      page,
      pageSize
    }
  })
}

// 获取子评论数量
export const getReplyCount = (commentId) => {
  return request({
    url: `${API_PREFIX}/public/comments/${commentId}/count`,
    method: 'get'
  })
}

// 获取专栏列表
export const getSeries = () => {
  if (USE_MOCK) return mockApi.getSeries()
  return request({
    url: `${API_PREFIX}/public/series`,
    method: 'get'
  })
}

// 获取专栏详情
export const getSeriesDetail = (id) => {
  if (USE_MOCK) return mockApi.getSeriesDetail(id)
  return request({
    url: `${API_PREFIX}/public/series/${id}`,
    method: 'get'
  })
}

// 获取专栏章节详情
export const getChapterDetail = (seriesId, chapterId) => {
  if (USE_MOCK) return mockApi.getChapterDetail(seriesId, chapterId)
  return request({
    url: `${API_PREFIX}/public/series/${seriesId}/sections/${chapterId}`,
    method: 'get'
  })
}

// 获取图形验证码
export const getCaptcha = (type = 1) => {
  return request({
    url: `${API_PREFIX}/public/auth/captcha`,
    method: 'get',
    params: { type }
  })
}

// 密码登录
export const login = (data) => {
  return request({
    url: `${API_PREFIX}/public/auth/login`,
    method: 'post',
    data
  })
}

// 刷新令牌
export const refreshToken = (data) => {
  return request({
    url: `${API_PREFIX}/public/auth/refresh`,
    method: 'post',
    data
  })
}

// 发送邮箱验证码
export const sendEmailCode = (email) => {
  return request({
    url: `${API_PREFIX}/public/auth/send-email-code`,
    method: 'post',
    data: { email }
  })
}

// 邮箱验证码注册
export const registerWithEmailCode = (data) => {
  return request({
    url: `${API_PREFIX}/public/auth/register`,
    method: 'post',
    data
  })
}

// 验证邮箱验证码
export const verifyEmailCaptcha = (data) => {
  return request({
    url: `${API_PREFIX}/public/auth/verify`,
    method: 'post',
    data
  })
}

// ============= Front API (前台接口，需要登录) =============

// 点赞文章
export const likeArticle = (id) => {
  if (USE_MOCK) return mockApi.likeArticle(id)
  return request({
    url: `${API_PREFIX}/front/articles/${id}/like`,
    method: 'post'
  })
}

// 取消点赞文章
export const unlikeArticle = (id) => {
  if (USE_MOCK) return mockApi.unlikeArticle(id)
  return request({
    url: `${API_PREFIX}/front/articles/${id}/like`,
    method: 'delete'
  })
}

// 检查文章点赞状态
export const checkArticleLikeStatus = (id) => {
  if (USE_MOCK) {
    return Promise.resolve({
      code: 200,
      data: {
        isLiked: false
      }
    })
  }
  return request({
    url: `${API_PREFIX}/front/articles/${id}/like/status`,
    method: 'get'
  })
}

// 收藏文章
export const favoriteArticle = (data) => {
  if (USE_MOCK) return mockApi.favoriteArticle(data.articleId)
  return request({
    url: `${API_PREFIX}/front/favorite/add`,
    method: 'post',
    data: {
      article_id: data.articleId,
      folder_id: data.folderId
    }
  })
}

// 取消收藏文章
export const unfavoriteArticle = (data) => {
  if (USE_MOCK) return mockApi.unfavoriteArticle(data.articleId)
  return request({
    url: `${API_PREFIX}/front/articles/${data.articleId}/favorite`,
    method: 'delete'
  })
}

// 增加文章阅读量
export const incrementArticleView = (id) => {
  return request({
    url: `${API_PREFIX}/front/articles/${id}/view`,
    method: 'post'
  })
}

// 添加评论
export const addComment = (data) => {
  if (USE_MOCK) return mockApi.addComment(data)
  return request({
    url: `${API_PREFIX}/front/comments`,
    method: 'post',
    data
  })
}

// 点赞评论
export const likeComment = (commentId) => {
  if (USE_MOCK) return Promise.resolve({ code: 200, message: '点赞成功' })
  return request({
    url: `${API_PREFIX}/front/comments/${commentId}/like`,
    method: 'post'
  })
}

// 取消点赞评论
export const unlikeComment = (commentId) => {
  if (USE_MOCK) return Promise.resolve({ code: 200, message: '取消点赞成功' })
  return request({
    url: `${API_PREFIX}/front/comments/${commentId}/like`,
    method: 'delete'
  })
}


// 更新评论
export const updateComment = (id, data) => {
  return request({
    url: `${API_PREFIX}/front/comments/${id}`,
    method: 'put',
    data
  })
}

// 删除评论
export const deleteComment = (id) => {
  if (USE_MOCK) return mockApi.deleteComment(id)
  return request({
    url: `${API_PREFIX}/front/comments/${id}`,
    method: 'delete'
  })
}

// 获取用户评论列表
export const getUserComments = (userId) => {
  return request({
    url: `${API_PREFIX}/front/comments/user/${userId}`,
    method: 'get'
  })
}


// 获取当前用户信息
export const getCurrentUser = () => {
  return request({
    url: `${API_PREFIX}/front/users/me`,
    method: 'get'
  })
}

// 更新当前用户信息
export const updateCurrentUser = (data) => {
  return request({
    url: `${API_PREFIX}/front/users/me`,
    method: 'put',
    data
  })
}

// ============= 收藏文件夹 API =============

// 获取用户的收藏文件夹列表
export const getFavoriteFolders = (params) => {
  if (USE_MOCK) {
    return Promise.resolve({
      code: 200,
      data: [
        { id: 1, name: '默认收藏夹', isDefault: true, articleCount: 2, description: '系统默认收藏夹' },
        { id: 2, name: '技术干货', isDefault: false, articleCount: 3, description: '实用的技术文章' },
        { id: 3, name: '学习笔记', isDefault: false, articleCount: 1, description: '日常学习记录' }
      ]
    })
  }
  return request({
    url: `${API_PREFIX}/front/favorite/userFolder`,
    method: 'get',
    params
  })
}

// 创建收藏文件夹
export const createFavoriteFolder = (data) => {
  if (USE_MOCK) {
    return Promise.resolve({
      code: 200,
      data: { id: Date.now(), ...data, isDefault: false, articleCount: 0 }
    })
  }
  return request({
    url: `${API_PREFIX}/front/favorite/create`,
    method: 'post',
    data
  })
}

// 更新收藏文件夹
export const updateFavoriteFolder = (id, data) => {
  if (USE_MOCK) {
    return Promise.resolve({ code: 200, data: { id, ...data } })
  }
  return request({
    url: `${API_PREFIX}/front/favorite/userFolder/${id}`,
    method: 'put',
    data
  })
}

// 删除收藏文件夹
export const deleteFavoriteFolder = (id) => {
  if (USE_MOCK) {
    return Promise.resolve({ code: 200, message: '删除成功' })
  }
  return request({
    url: `${API_PREFIX}/front/favorite/userFolder/${id}`,
    method: 'delete'
  })
}

// 获取文件夹下的收藏文章列表
export const getFolderArticles = (folderId, params = {}) => {
  if (USE_MOCK) {
    const { page = 1, pageSize = 5 } = params
    const allArticles = [
      { id: 2, title: 'JavaScript 异步编程深入理解', views: 987, likes: 28, createdAt: '2024-01-12' },
      { id: 3, title: 'React Hooks 完全指南', views: 856, likes: 22, createdAt: '2024-01-08' },
      { id: 4, title: 'TypeScript 高级类型实战', views: 765, likes: 25, createdAt: '2024-01-10' },
      { id: 5, title: 'CSS Grid 布局详解', views: 543, likes: 19, createdAt: '2024-01-05' },
      { id: 6, title: 'Vue 3 Composition API 最佳实践', views: 1234, likes: 35, createdAt: '2024-01-15' },
      { id: 7, title: 'Webpack 性能优化指南', views: 678, likes: 21, createdAt: '2024-01-09' },
      { id: 8, title: 'Node.js 事件循环机制', views: 892, likes: 27, createdAt: '2024-01-11' },
      { id: 9, title: '前端性能监控方案', views: 456, likes: 16, createdAt: '2024-01-03' },
      { id: 10, title: 'React 源码解析', views: 1567, likes: 42, createdAt: '2024-01-18' },
      { id: 11, title: 'JavaScript 设计模式', views: 1098, likes: 31, createdAt: '2024-01-14' }
    ]

    const startIndex = (page - 1) * pageSize
    const endIndex = startIndex + pageSize
    const list = allArticles.slice(startIndex, endIndex)

    return Promise.resolve({
      code: 200,
      data: {
        list,
        total: allArticles.length
      }
    })
  }
  return request({
    url: `${API_PREFIX}/front/favorite/list`,
    method: 'get',
    params: folderId ? { ...params, folderId } : params
  })
}

// 从文件夹中删除收藏文章（实际上是从收藏表中删除该文章的所有收藏记录）
export const removeArticleFromFolder = (_folderId, articleId) => {
  if (USE_MOCK) {
    return Promise.resolve({ code: 200, message: '删除成功' })
  }
  return request({
    url: `${API_PREFIX}/front/favorite/remove`,
    method: 'post',
    data: {
      articleId: articleId
    }
  })
}

// 移动文章到另一个文件夹
export const moveArticleToFolder = (articleId, _fromFolderId, toFolderId) => {
  if (USE_MOCK) {
    return Promise.resolve({ code: 200, message: '移动成功' })
  }
  return request({
    url: `${API_PREFIX}/front/favorite/move`,
    method: 'post',
    data: {
      articleId: articleId,
      folderId: toFolderId
    }
  })
}

// 检查文章收藏状态
export const checkFavoriteStatus = (articleId) => {
  if (USE_MOCK) {
    return Promise.resolve({
      code: 200,
      data: {
        isFavorited: false,
        folderId: null
      }
    })
  }
  return request({
    url: `${API_PREFIX}/front/favorite/check/${articleId}`,
    method: 'get'
  })
}

// 获取用户详情
export const getUserDetail = (id) => {
  return request({
    url: `${API_PREFIX}/front/users/${id}`,
    method: 'get'
  })
}

// 获取用户活动记录
export const getUserActivities = (userId, params = {}) => {
  if (USE_MOCK) {
    return Promise.resolve({
      code: 200,
      data: {
        list: [],
        total: 0,
        stats: {
          total: 0,
          like: 0,
          comment: 0,
          share: 0,
          favorite: 0
        }
      }
    })
  }
  return request({
    url: `${API_PREFIX}/front/users/${userId}/activities`,
    method: 'get',
    params
  })
}

// 获取用户活动统计
export const getUserActivityStats = (userId) => {
  if (USE_MOCK) {
    return Promise.resolve({
      code: 200,
      data: {
        total: 0,
        like: 0,
        comment: 0,
        share: 0,
        favorite: 0
      }
    })
  }
  return request({
    url: `${API_PREFIX}/front/users/${userId}/activities/stats`,
    method: 'get'
  })
}

// 创建分类
export const createCategory = (data) => {
  return request({
    url: `${API_PREFIX}/front/categories`,
    method: 'post',
    data
  })
}

// 更新分类
export const updateCategory = (id, data) => {
  return request({
    url: `${API_PREFIX}/front/categories/${id}`,
    method: 'put',
    data
  })
}

// 删除分类
export const deleteCategory = (id) => {
  return request({
    url: `${API_PREFIX}/front/categories/${id}`,
    method: 'delete'
  })
}

// 创建标签
export const createTag = (data) => {
  return request({
    url: `${API_PREFIX}/front/tags`,
    method: 'post',
    data
  })
}

// 更新标签
export const updateTag = (id, data) => {
  return request({
    url: `${API_PREFIX}/front/tags/${id}`,
    method: 'put',
    data
  })
}

// 删除标签
export const deleteTag = (id) => {
  return request({
    url: `${API_PREFIX}/front/tags/${id}`,
    method: 'delete'
  })
}

// 创建系列
export const createSeries = (data) => {
  return request({
    url: `${API_PREFIX}/front/series`,
    method: 'post',
    data
  })
}

// 更新系列
export const updateSeries = (id, data) => {
  return request({
    url: `${API_PREFIX}/front/series/${id}`,
    method: 'put',
    data
  })
}

// 删除系列
export const deleteSeries = (id) => {
  return request({
    url: `${API_PREFIX}/front/series/${id}`,
    method: 'delete'
  })
}

// 创建章节
export const createSection = (seriesId, data) => {
  return request({
    url: `${API_PREFIX}/front/series/${seriesId}/sections`,
    method: 'post',
    data
  })
}

// 邮箱验证码登录
export const loginWithEmailCode = (data) => {
  return request({
    url: `${API_PREFIX}/public/auth/loginByEmail`,
    method: 'post',
    data
  })
}





