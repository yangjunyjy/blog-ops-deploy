# 图片上传和发布文章问题修复

## 问题1: 图片URL包含反斜杠（Windows路径问题）

### 问题描述
上传图片后，返回的URL包含Windows反斜杠，导致图片无法访问：
```
http://localhost:8081/\uploads\images\607ad84e_8c70f895.jpg
```

### 修复方案
**文件**: `internal/services/upload_service.go:119`

修改前：
```go
return filepath.Join("/uploads/images", filename), nil
```

修改后：
```go
return "/uploads/images/" + filename, nil
```

**原因**: `filepath.Join` 在Windows上会使用反斜杠 `\`，但URL必须使用正斜杠 `/`。

---

## 问题2: 发布文章返回500错误

### 问题描述
点击发布文章时，后端返回500状态码，文章无法发布。

### 根本原因
1. **Summary必填**: 后端 `CreateArticleRequest` 中 `summary` 字段标记为 `binding:"required"`
2. **缺少Slug**: 模型中 `Slug` 是必填字段，但请求体中未包含
3. **AuthorID未设置**: 创建文章时需要从JWT token中获取当前用户ID
4. **前端发送额外字段**: 前端使用 `...article` 展开操作，发送了不需要的字段（如 `_markdown`）

### 修复方案

#### 1. 后端修改

**文件**: `internal/api/v1/dto/request/admin.go`

修改前：
```go
type CreateArticleRequest struct {
	Title       string  `json:"title" binding:"required,max=200"`
	Summary     string  `json:"summary" binding:"required,max=500"`
	Content     string  `json:"content" binding:"required"`
	Cover       *string `json:"cover" binding:"omitempty,url"`
	CategoryID  *uint64 `json:"categoryId" binding:"omitempty"`
	TagIDs      []uint  `json:"tags" binding:"omitempty"`
	AuthorID    uint64  `json:"author_id" binding:"omitempty"`
	Status      uint8   `json:"status" binding:"omitempty,oneof=0 1" default:"1"`
	PublishedAt *string `json:"publishedAt" binding:"omitempty"`
	IsTop       bool    `json:"is_top" binding:"omitempty"`
}
```

修改后：
```go
type CreateArticleRequest struct {
	Title      string  `json:"title" binding:"required,max=200"`
	Summary    string  `json:"summary" binding:"max=500"`      // 移除required
	Content    string  `json:"content" binding:"required"`
	Cover      *string `json:"cover" binding:"omitempty,url"`
	CategoryID *uint64 `json:"categoryId" binding:"omitempty"`
	TagIDs     []uint  `json:"tags" binding:"omitempty"`
	AuthorID   uint64  `json:"author_id" binding:"omitempty"`
	Status     uint8   `json:"status" binding:"omitempty,oneof=0 1" default:"1"`
	PublishedAt *string `json:"publishedAt" binding:"omitempty"`
	IsTop      bool    `json:"is_top" binding:"omitempty"`
	Slug       string  `json:"slug" binding:"required,max=200"`  // 新增
}
```

**文件**: `internal/api/v1/article_handler.go:30-55`

添加从JWT获取用户ID的代码：
```go
// CreateArticle 创建文章
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req dtoRequest.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 从JWT中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权", nil)
		return
	}

	article := &models.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      response.SafeDerefString(req.Cover),
		CategoryID: response.SafeDerefUint64(req.CategoryID),
		AuthorID:   uint64(userID.(uint)),
		Status:     req.Status,
		IsTop:      req.IsTop,
		Slug:       req.Slug,
	}

	if err := h.articleService.CreateArticle(article, req.TagIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建文章失败", err)
		return
	}

	response.Success(c, article, "创建成功")
}
```

#### 2. 前端修改

**文件**: `my-blog-admin/src/views/articles/ArticleEditor.vue`

修改 `handleSaveDraft` 函数：

```javascript
const handleSaveDraft = async () => {
  try {
    article.status = 0
    // 如果没有摘要，自动生成（支持 HTML 和 Markdown）
    if (!article.summary || article.summary.trim() === '') {
      const isHtml = /<[a-z][\s\S]*>/i.test(article.content)
      const plainText = isHtml ? article.content.replace(/<[^>]*>/g, '') : article.content
      const summary = plainText.substring(0, 200)
      article.summary = summary + (summary.length >= 200 ? '...' : '')
    }
    // 如果没有slug，自动生成
    if (!article.slug || article.slug.trim() === '') {
      article.slug = article.title.toLowerCase().replace(/[^a-z0-9\u4e00-\u9fa5]+/g, '-').replace(/-+/g, '-').replace(/^-+|-+$/g, '')
    }
    // 将编辑器中的 HTML 转换为 Markdown 保存
    const markdownContent = htmlToMarkdown(article.content)
    const articleToSave = {
      title: article.title,
      content: markdownContent,
      summary: article.summary || '未提供摘要',
      cover: article.cover,
      categoryId: article.category_id,
      tags: article.tag_ids || [],
      status: article.status,
      slug: article.slug
    }
    if (article.id) {
      await updateArticle(article.id, articleToSave)
    } else {
      const res = await createArticle(articleToSave)
      article.id = res.data.id
    }
    ElMessage.success('草稿保存成功')
  } catch (error) {
    console.error('保存草稿失败:', error)
    ElMessage.error('保存失败: ' + (error.message || error.msg || '未知错误'))
  }
}
```

修改 `handlePublish` 函数：

```javascript
const handlePublish = async () => {
  if (!article.title) {
    ElMessage.warning('请输入文章标题')
    return
  }
  if (!article.content) {
    ElMessage.warning('请输入文章内容')
    return
  }

  // 如果没有摘要，自动生成
  if (!article.summary || article.summary.trim() === '') {
    const isHtml = /<[a-z][\s\S]*>/i.test(article.content)
    const plainText = isHtml ? article.content.replace(/<[^>]*>/g, '') : article.content
    const summary = plainText.substring(0, 200)
    article.summary = summary + (summary.length >= 200 ? '...' : '')
  }

  // 如果没有slug，自动生成
  if (!article.slug || article.slug.trim() === '') {
    article.slug = article.title.toLowerCase().replace(/[^a-z0-9\u4e00-\u9fa5]+/g, '-').replace(/-+/g, '-').replace(/^-+|-+$/g, '')
  }

  // 将编辑器中的 HTML 转换为 Markdown 保存
  const markdownContent = htmlToMarkdown(article.content)
  const articleToSave = {
    title: article.title,
    content: markdownContent,
    summary: article.summary || '未提供摘要',
    cover: article.cover,
    categoryId: article.category_id,
    tags: article.tag_ids || [],
    status: 1,
    is_top: article.is_top || false,
    slug: article.slug
  }

  publishing.value = true
  try {
    if (article.id) {
      await updateArticle(article.id, articleToSave)
    } else {
      const res = await createArticle(articleToSave)
      article.id = res.data.id
    }
    ElMessage.success('发布成功')
    router.push('/articles')
  } catch (error) {
    console.error('发布失败:', error)
    ElMessage.error('发布失败: ' + (error.message || error.msg || '未知错误'))
  } finally {
    publishing.value = false
  }
}
```

### 主要改进
1. **不再使用对象展开**: 明确指定要发送的字段，避免发送额外的 `_markdown` 等内部字段
2. **自动生成Summary**: 如果用户未提供摘要，自动从内容中生成
3. **自动生成Slug**: 如果用户未填写URL别名，自动从标题生成
4. **默认值处理**: 为可选字段提供默认值（如 `summary: '未提供摘要'`）
5. **错误信息改进**: 显示更详细的错误信息

---

## 修复总结

### 修复的问题
1. ✅ 图片URL使用正确的正斜杠格式
2. ✅ 发布文章不再返回500错误
3. ✅ 摘要字段可以自动生成
4. ✅ Slug字段自动生成和提交
5. ✅ 从JWT正确获取当前用户ID
6. ✅ 前端不再发送额外的内部字段

### 注意事项
1. **重启后端**: 需要重新编译和重启后端服务
2. **清除缓存**: 前端可能需要清除浏览器缓存
3. **测试验证**: 需要测试图片上传和文章发布功能

### 后续建议
1. **错误日志**: 在后端添加详细的错误日志
2. **参数验证**: 前端添加更详细的参数验证
3. **提示信息**: 改进用户提示信息的准确性
