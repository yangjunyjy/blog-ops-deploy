/**
 * Mock API
 */

import { articles, categories, tags, series, comments, users } from './data/blog'

const mockDB = {
  articles: [...articles],
  categories: [...categories],
  tags: [...tags],
  series: [...series],
  comments: [...comments],
  users: [...users],
  hosts: [
    { id: 1, name: 'Web服务器', ip: '192.168.1.100', port: 22, auth_type: 'password', username: 'root', password: 'password123', private_key: '', description: '主要Web服务器', status: 1 },
    { id: 2, name: '数据库服务器', ip: '192.168.1.101', port: 22, auth_type: 'key', username: 'admin', password: '', private_key: '-----BEGIN RSA PRIVATE KEY-----\nMIIEow...', description: 'MySQL数据库服务器', status: 1 },
    { id: 3, name: '测试服务器', ip: '192.168.1.102', port: 22, auth_type: 'password', username: 'test', password: 'test123', private_key: '', description: '测试环境', status: 0 },
    { id: 4, name: '备份服务器', ip: '192.168.1.103', port: 22, auth_type: 'password', username: 'backup', password: 'backup456', private_key: '', description: '备份存储服务器', status: 1 },
    { id: 5, name: '监控服务器', ip: '192.168.1.104', port: 22, auth_type: 'key', username: 'monitor', password: '', private_key: '-----BEGIN RSA PRIVATE KEY-----\nMIIEpA...', description: 'Zabbix监控服务器', status: 1 }
  ]
}

const delay = (ms = 300) => new Promise(resolve => setTimeout(resolve, ms))
const sleep = (ms = 300) => new Promise(resolve => setTimeout(resolve, ms))

const successResponse = (data, message = 'Success') => ({
  code: 200,
  message,
  data
})

const errorResponse = (message = 'Error', code = 500) => ({
  code,
  message,
  data: null
})

const paginate = (data, page, size) => {
  const start = (page - 1) * size
  const end = start + size
  return {
    items: data.slice(start, end),
    total: data.length,
    page,
    size,
    totalPages: Math.ceil(data.length / size)
  }
}

export const mockApi = {
  // ========== 认证 ==========
  uploadImage: async (formData) => {
    await sleep(1000)
    const file = formData.get('file')
    const url = URL.createObjectURL(file)
    return {
      code: 200,
      data: {
        url,
        name: file.name,
        size: file.size,
        type: file.type
      },
      message: '上传成功'
    }
  },
   getArticleDetail: async (id) => {
    await sleep(500)
    return {
      code: 200,
      data: {
        id,
        title: '示例文章标题',
        content: '# 示例文章\n\n这是一个示例文章内容...',
        cover: 'https://picsum.photos/800/400',
        categoryId: 1,
        tagIds: [1, 2],
        summary: '这是文章的摘要内容',
        keywords: 'Vue,Element Plus,Markdown',
        isTop: false,
        isRecommended: true,
        status: 1,
        allowComment: true,
        createdAt: new Date().toISOString(),
        publishTime: new Date().toISOString()
      }
    }
  },
   // 创建文章
  createArticle: async (article) => {
    await sleep(1000)
    return {
      code: 200,
      data: {
        id: Date.now(),
        ...article
      },
      message: '创建成功'
    }
  },
  updateArticle: async (id, article) => {
    await sleep(800)
    return {
      code: 200,
      data: {
        id,
        ...article
      },
      message: '更新成功'
    }
  },
  async login(data) {
    await delay()
    const user = mockDB.users.find(u => u.username === data.username && u.password === data.password)
    if (user) {
      const token = 'mock-token-' + Date.now()
      return successResponse({
        token,
        userInfo: user
      }, '登录成功')
    }
    return errorResponse('用户名或密码错误', 401)
  },

  async logout() {
    await delay()
    return successResponse(null, '登出成功')
  },

  async getUserInfo() {
    await delay()
    return successResponse(mockDB.users[0])
  },

  // ========== 统计 ==========
  async getDashboardStats() {
    await delay()
    return successResponse({
      overview: {
        articles: mockDB.articles.length,
        published: mockDB.articles.filter(a => a.status === 1).length,
        draft: mockDB.articles.filter(a => a.status === 0).length,
        views: mockDB.articles.reduce((sum, a) => sum + a.views, 0),
        likes: mockDB.articles.reduce((sum, a) => sum + a.likes, 0),
        comments: mockDB.comments.length
      },
      content: {
        categories: mockDB.categories.length,
        tags: mockDB.tags.length,
        series: mockDB.series.length
      },
      trend: {
        views: [120, 132, 101, 134, 90, 230, 210],
        articles: [5, 8, 6, 9, 7, 10, 8]
      }
    })
  },

  // ========== 文章 ==========
  async getArticleList(params) {
    await delay()
    const { title, categoryId, status, page = 1, size = 10 } = params || {}
    let filtered = [...mockDB.articles]
    
    if (title) filtered = filtered.filter(a => a.title.includes(title))
    if (categoryId) filtered = filtered.filter(a => a.categoryId === Number(categoryId))
    if (status !== '') filtered = filtered.filter(a => a.status === Number(status))
    
    filtered.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
    const paginated = paginate(filtered, page, size)
    
    return successResponse({
      list: paginated.items,
      ...paginated
    })
  },

  async getArticleDetail(id) {
    await delay()
    const article = mockDB.articles.find(a => a.id === Number(id))
    return article ? successResponse(article) : errorResponse('文章不存在', 404)
  },

  async createArticle(data) {
    await delay()
    const newArticle = {
      id: mockDB.articles.length + 1,
      ...data,
      views: 0,
      likes: 0,
      favorites: 0,
      commentCount: 0,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    }
    mockDB.articles.push(newArticle)
    return successResponse(newArticle, '创建成功')
  },

  async updateArticle(id, data) {
    await delay()
    const index = mockDB.articles.findIndex(a => a.id === Number(id))
    if (index === -1) return errorResponse('文章不存在', 404)
    
    mockDB.articles[index] = { ...mockDB.articles[index], ...data, updatedAt: new Date().toISOString() }
    return successResponse(mockDB.articles[index], '更新成功')
  },

  async deleteArticle(id) {
    await delay()
    const index = mockDB.articles.findIndex(a => a.id === Number(id))
    if (index === -1) return errorResponse('文章不存在', 404)

    mockDB.articles.splice(index, 1)
    return successResponse(null, '删除成功')
  },

  async batchDeleteArticles(ids) {
    await delay()
    ids.forEach(id => {
      const index = mockDB.articles.findIndex(a => a.id === Number(id))
      if (index !== -1) mockDB.articles.splice(index, 1)
    })
    return successResponse(null, `成功删除${ids.length}篇文章`)
  },

  async publishArticle(id) {
    await delay()
    const article = mockDB.articles.find(a => a.id === Number(id))
    if (article) article.status = 1
    return successResponse(article, '发布成功')
  },

  async withdrawArticle(id) {
    await delay()
    const article = mockDB.articles.find(a => a.id === Number(id))
    if (article) article.status = 0
    return successResponse(article, '撤回成功')
  },

  // ========== 分类 ==========
  async getCategoryList(params) {
    await delay()
    const { name, page = 1, size = 10 } = params || {}
    let filtered = [...mockDB.categories]

    if (name) filtered = filtered.filter(c => c.name.includes(name))

    const paginated = paginate(filtered, page, size)
    return successResponse({
      list: paginated.items,
      ...paginated
    })
  },

  async getCategoryDetail(id) {
    await delay()
    const category = mockDB.categories.find(c => c.id === Number(id))
    return category ? successResponse(category) : errorResponse('分类不存在', 404)
  },

  async createCategory(data) {
    await delay()
    const newCategory = {
      id: mockDB.categories.length + 1,
      ...data,
      articleCount: 0,
      createdAt: new Date().toISOString()
    }
    mockDB.categories.push(newCategory)
    return successResponse(newCategory, '创建成功')
  },

  async updateCategory(id, data) {
    await delay()
    const index = mockDB.categories.findIndex(c => c.id === Number(id))
    if (index === -1) return errorResponse('分类不存在', 404)
    
    mockDB.categories[index] = { ...mockDB.categories[index], ...data }
    return successResponse(mockDB.categories[index], '更新成功')
  },

  async deleteCategory(id) {
    await delay()
    const index = mockDB.categories.findIndex(c => c.id === Number(id))
    if (index === -1) return errorResponse('分类不存在', 404)

    mockDB.categories.splice(index, 1)
    return successResponse(null, '删除成功')
  },

  async getCategories(params) {
    await delay()
    const { page = 1, size = 100 } = params || {}
    const paginated = paginate(mockDB.categories, page, size)
    return successResponse({
      list: paginated.items,
      ...paginated
    })
  },

  // ========== 标签 ==========
  async getTagList(params) {
    await delay()
    const { name, page = 1, size = 10 } = params || {}
    let filtered = [...mockDB.tags]

    if (name) filtered = filtered.filter(t => t.name.includes(name))

    const paginated = paginate(filtered, page, size)
    return successResponse({
      list: paginated.items,
      ...paginated
    })
  },

  async getTagDetail(id) {
    await delay()
    const tag = mockDB.tags.find(t => t.id === Number(id))
    return tag ? successResponse(tag) : errorResponse('标签不存在', 404)
  },

  async createTag(data) {
    await delay()
    const newTag = {
      id: mockDB.tags.length + 1,
      ...data,
      articleCount: 0,
      createdAt: new Date().toISOString()
    }
    mockDB.tags.push(newTag)
    return successResponse(newTag, '创建成功')
  },

  async updateTag(id, data) {
    await delay()
    const index = mockDB.tags.findIndex(t => t.id === Number(id))
    if (index === -1) return errorResponse('标签不存在', 404)
    
    mockDB.tags[index] = { ...mockDB.tags[index], ...data }
    return successResponse(mockDB.tags[index], '更新成功')
  },

  async deleteTag(id) {
    await delay()
    const index = mockDB.tags.findIndex(t => t.id === Number(id))
    if (index === -1) return errorResponse('标签不存在', 404)

    mockDB.tags.splice(index, 1)
    return successResponse(null, '删除成功')
  },

  async batchDeleteTags(ids) {
    await delay()
    ids.forEach(id => {
      const index = mockDB.tags.findIndex(t => t.id === Number(id))
      if (index !== -1) mockDB.tags.splice(index, 1)
    })
    return successResponse(null, `成功删除${ids.length}个标签`)
  },

  async getTags(params) {
    await delay()
    const { page = 1, size = 100 } = params || {}
    const paginated = paginate(mockDB.tags, page, size)
    return successResponse({
      list: paginated.items,
      ...paginated
    })
  },

  // ========== 系列 ==========
  async getSeriesList(params) {
    await delay()
    const { name, status, page = 1, size = 10 } = params || {}
    let filtered = [...mockDB.series]
    
    if (name) filtered = filtered.filter(s => s.name.includes(name))
    if (status !== '') filtered = filtered.filter(s => s.status === Number(status))
    
    const paginated = paginate(filtered, page, size)
    return successResponse({
      list: paginated.items,
      ...paginated
    })
  },

  async getSeriesDetail(id) {
    await delay()
    const series = mockDB.series.find(s => s.id === Number(id))
    return series ? successResponse(series) : errorResponse('系列不存在', 404)
  },

  async createSeries(data) {
    await delay()
    const newSeries = {
      id: mockDB.series.length + 1,
      ...data,
      sections: [],
      createdAt: new Date().toISOString()
    }
    mockDB.series.push(newSeries)
    return successResponse(newSeries, '创建成功')
  },

  async updateSeries(id, data) {
    await delay()
    const index = mockDB.series.findIndex(s => s.id === Number(id))
    if (index === -1) return errorResponse('系列不存在', 404)
    
    mockDB.series[index] = { ...mockDB.series[index], ...data }
    return successResponse(mockDB.series[index], '更新成功')
  },

  async deleteSeries(id) {
    await delay()
    const index = mockDB.series.findIndex(s => s.id === Number(id))
    if (index === -1) return errorResponse('系列不存在', 404)

    mockDB.series.splice(index, 1)
    return successResponse(null, '删除成功')
  },

  async createSection(seriesId, data) {
    await delay()
    const series = mockDB.series.find(s => s.id === Number(seriesId))
    if (!series) return errorResponse('系列不存在', 404)

    const newSection = {
      id: Date.now(),
      ...data,
      sections: []
    }
    series.sections.push(newSection)
    return successResponse(newSection, '创建章节成功')
  },

  async updateSection(seriesId, sectionId, data) {
    await delay()
    const series = mockDB.series.find(s => s.id === Number(seriesId))
    if (!series) return errorResponse('系列不存在', 404)

    const section = series.sections.find(s => s.id === Number(sectionId))
    if (!section) return errorResponse('章节不存在', 404)

    Object.assign(section, data)
    return successResponse(section, '更新章节成功')
  },

  async deleteSection(seriesId, sectionId) {
    await delay()
    const series = mockDB.series.find(s => s.id === Number(seriesId))
    if (!series) return errorResponse('系列不存在', 404)

    const index = series.sections.findIndex(s => s.id === Number(sectionId))
    if (index === -1) return errorResponse('章节不存在', 404)

    series.sections.splice(index, 1)
    return successResponse(null, '删除章节成功')
  },

  // ========== 评论 ==========
  async getCommentList(params) {
    await delay()
    const { articleId, status, page = 1, size = 10 } = params || {}
    let filtered = [...mockDB.comments]

    if (articleId) filtered = filtered.filter(c => c.articleId === Number(articleId))
    if (status !== '') filtered = filtered.filter(c => c.status === Number(status))

    const paginated = paginate(filtered, page, size)
    return successResponse({
      list: paginated.items,
      ...paginated
    })
  },

  async getCommentDetail(id) {
    await delay()
    const comment = mockDB.comments.find(c => c.id === Number(id))
    return comment ? successResponse(comment) : errorResponse('评论不存在', 404)
  },

  async approveComment(id) {
    await delay()
    const comment = mockDB.comments.find(c => c.id === Number(id))
    if (comment) comment.status = 1
    return successResponse(comment, '审核通过')
  },

  async rejectComment(id) {
    await delay()
    const comment = mockDB.comments.find(c => c.id === Number(id))
    if (comment) comment.status = 2
    return successResponse(comment, '审核拒绝')
  },

  async deleteComment(id) {
    await delay()
    const index = mockDB.comments.findIndex(c => c.id === Number(id))
    if (index === -1) return errorResponse('评论不存在', 404)

    mockDB.comments.splice(index, 1)
    return successResponse(null, '删除成功')
  },

  async batchDeleteComments(ids) {
    await delay()
    ids.forEach(id => {
      const index = mockDB.comments.findIndex(c => c.id === Number(id))
      if (index !== -1) mockDB.comments.splice(index, 1)
    })
    return successResponse(null, `成功删除${ids.length}条评论`)
  },

  async batchApproveComments(ids) {
    await delay()
    ids.forEach(id => {
      const comment = mockDB.comments.find(c => c.id === Number(id))
      if (comment) comment.status = 1
    })
    return successResponse(null, `成功审核${ids.length}条评论`)
  },

  // ========== 用户 ==========
  async getUserList(params) {
    await delay()
    const { username, role, page = 1, size = 10 } = params || {}
    let filtered = [...mockDB.users]

    if (username) filtered = filtered.filter(u => u.username.includes(username))
    if (role !== '') filtered = filtered.filter(u => u.role === Number(role))

    const paginated = paginate(filtered, page, size)
    return successResponse({
      list: paginated.items,
      ...paginated
    })
  },

  async getUserDetail(id) {
    await delay()
    const user = mockDB.users.find(u => u.id === Number(id))
    return user ? successResponse(user) : errorResponse('用户不存在', 404)
  },

  async createUser(data) {
    await delay()
    const newUser = {
      id: mockDB.users.length + 1,
      ...data,
      createdAt: new Date().toISOString()
    }
    mockDB.users.push(newUser)
    return successResponse(newUser, '创建成功')
  },

  async updateUser(id, data) {
    await delay()
    const index = mockDB.users.findIndex(u => u.id === Number(id))
    if (index === -1) return errorResponse('用户不存在', 404)
    
    mockDB.users[index] = { ...mockDB.users[index], ...data }
    return successResponse(mockDB.users[index], '更新成功')
  },

  async deleteUser(id) {
    await delay()
    const index = mockDB.users.findIndex(u => u.id === Number(id))
    if (index === -1) return errorResponse('用户不存在', 404)

    mockDB.users.splice(index, 1)
    return successResponse(null, '删除成功')
  },

  async batchDeleteUsers(ids) {
    await delay()
    ids.forEach(id => {
      const index = mockDB.users.findIndex(u => u.id === Number(id))
      if (index !== -1) mockDB.users.splice(index, 1)
    })
    return successResponse(null, `成功删除${ids.length}个用户`)
  },

  async resetPassword(id) {
    await delay()
    const user = mockDB.users.find(u => u.id === Number(id))
    if (user) user.password = '123456'
    return successResponse({ password: '123456' }, '密码已重置为123456')
  },

  async updateUserStatus(id, status) {
    await delay()
    const user = mockDB.users.find(u => u.id === Number(id))
    if (user) user.status = status
    return successResponse(user, '状态更新成功')
  },

  // ========== 统计 ==========
  async getArticleStats(params) {
    await delay()
    return successResponse({
      total: mockDB.articles.length,
      published: mockDB.articles.filter(a => a.status === 1).length,
      draft: mockDB.articles.filter(a => a.status === 0).length,
      totalViews: mockDB.articles.reduce((sum, a) => sum + a.views, 0),
      totalLikes: mockDB.articles.reduce((sum, a) => sum + a.likes, 0),
      trend: Array.from({ length: 7 }, (_, i) => Math.floor(Math.random() * 20) + 5)
    })
  },

  async getViewTrend(params) {
    await delay()
    const { days = 7 } = params || {}
    return successResponse({
      dates: Array.from({ length: days }, (_, i) => {
        const date = new Date()
        date.setDate(date.getDate() - (days - 1 - i))
        return date.toISOString().split('T')[0]
      }),
      views: Array.from({ length: days }, () => Math.floor(Math.random() * 200) + 50)
    })
  },

  async getUserGrowth(params) {
    await delay()
    const { days = 7 } = params || {}
    return successResponse({
      dates: Array.from({ length: days }, (_, i) => {
        const date = new Date()
        date.setDate(date.getDate() - (days - 1 - i))
        return date.toISOString().split('T')[0]
      }),
      users: Array.from({ length: days }, () => Math.floor(Math.random() * 10))
    })
  },

  async getHotArticles(params) {
    await delay()
    const { limit = 10 } = params || {}
    const sorted = [...mockDB.articles].sort((a, b) => b.views - a.views).slice(0, limit)
    return successResponse({
      list: sorted.map(a => ({
        id: a.id,
        title: a.title,
        views: a.views,
        likes: a.likes,
        commentCount: a.commentCount
      }))
    })
  },

  async getHotTags(params) {
    await delay()
    const { limit = 10 } = params || {}
    const sorted = [...mockDB.tags].sort((a, b) => b.articleCount - a.articleCount).slice(0, limit)
    return successResponse({
      list: sorted
    })
  },

  // ========== 主机管理 ==========
  async getHostList(params) {
    await delay()
    const { name, ip, status, page = 1, size = 10 } = params || {}
    let filtered = [...mockDB.hosts]
    
    if (name) filtered = filtered.filter(h => h.name.includes(name))
    if (ip) filtered = filtered.filter(h => h.ip.includes(ip))
    if (status !== '') filtered = filtered.filter(h => h.status === Number(status))
    
    const paginated = paginate(filtered, page, size)
    return successResponse({
      list: paginated.items,
      ...paginated
    })
  },

  async getHostDetail(id) {
    await delay()
    const host = mockDB.hosts.find(h => h.id === Number(id))
    return host ? successResponse(host) : errorResponse('主机不存在', 404)
  },

  async createHost(data) {
    await delay()
    const newHost = {
      id: mockDB.hosts.length + 1,
      ...data,
      createdAt: new Date().toISOString()
    }
    mockDB.hosts.push(newHost)
    return successResponse(newHost, '创建成功')
  },

  async updateHost(data) {
    await delay()
    const index = mockDB.hosts.findIndex(h => h.id === Number(data.id))
    if (index === -1) return errorResponse('主机不存在', 404)
    
    mockDB.hosts[index] = { ...mockDB.hosts[index], ...data, updatedAt: new Date().toISOString() }
    return successResponse(mockDB.hosts[index], '更新成功')
  },

  async deleteHost(id) {
    await delay()
    const index = mockDB.hosts.findIndex(h => h.id === Number(id))
    if (index === -1) return errorResponse('主机不存在', 404)

    mockDB.hosts.splice(index, 1)
    return successResponse(null, '删除成功')
  },

  async batchDeleteHosts(ids) {
    await delay()
    ids.forEach(id => {
      const index = mockDB.hosts.findIndex(h => h.id === Number(id))
      if (index !== -1) mockDB.hosts.splice(index, 1)
    })
    return successResponse(null, `成功删除${ids.length}个主机`)
  },

  async testHostConnection(id) {
    await delay()
    const host = mockDB.hosts.find(h => h.id === Number(id))
    if (!host) return errorResponse('主机不存在', 404)
    
    // 模拟连接测试，随机成功/失败
    const success = Math.random() > 0.3
    return successResponse({
      connected: success,
      message: success ? '连接成功' : '连接失败'
    })
  },

  async remoteConnectHost(id, auth_type) {
    await delay()
    const host = mockDB.hosts.find(h => h.id === Number(id))
    if (!host) return errorResponse('主机不存在', 404)
    
    // 模拟返回连接URL
    return successResponse({
      connect_url: `http://localhost:8080/ssh?host=${host.ip}&port=${host.port}&auth=${auth_type}`
    })
  },

  async transferFile(id, data) {
    await delay()
    const host = mockDB.hosts.find(h => h.id === Number(id))
    if (!host) return errorResponse('主机不存在', 404)
    
    // 模拟传输结果
    return successResponse({
      success: true,
      message: `文件已传输到 ${host.name}`,
      download_url: data.mode === 'download' ? `http://localhost:8080/download/${Date.now()}` : null
    })
  }
}
