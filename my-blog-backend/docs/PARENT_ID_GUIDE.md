# parentID 获取和传递说明

## 前端实现

### 数据流

```
用户点击"回复"按钮
  ↓
handleReply(comment) 被触发
  ↓
存储当前回复的评论对象到 replyingTo
  ↓
用户输入回复内容
  ↓
handleSubmitReply(comment) 被调用
  ↓
从 comment 对象获取 id 作为 parentId
  ↓
发送请求到后端
```

### 关键代码

#### 1. 用户点击回复按钮（第 58-61 行）

```javascript
const handleReply = (comment) => {
  replyingTo.value = comment  // 保存整个评论对象
  replyText.value = `@${comment.author.name} `  // 在输入框中预填充 @用户名
}
```

#### 2. 模板中渲染回复按钮（第 158-166 行）

```vue
<el-button
  type="primary"
  text
  size="small"
  :icon="CirclePlus"
  @click="handleReply(comment)"  <!-- 点击时传入完整的 comment 对象 -->
>
  回复
</el-button>
```

#### 3. 用户提交回复（第 68-89 行）

```javascript
const handleSubmitReply = async (comment) => {
  if (!replyText.value.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }

  submitting.value = true
  try {
    await addComment({
      articleId: props.articleId,
      content: replyText.value,
      parentId: comment.id  // ← 从 comment 对象中获取 id 作为 parentId
    })
    ElMessage.success('回复成功')
    handleCancelReply()
    await loadComments()
  } catch (error) {
    ElMessage.error('回复失败')
  } finally {
    submitting.value = false
  }
}
```

## 后端接收

### DTO 结构

```go
// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
    ArticleID uint   `json:"articleId" binding:"required"`
    Content   string `json:"content" binding:"required,min=1,max=1000"`
    ParentID  *uint  `json:"parentId" binding:"omitempty"`  // ← 指针类型，允许为空
}
```

### Handler 处理

```go
// CreateComment 创建评论（comment_handler.go 第 27-75 行）
func (h *CommentHandler) CreateComment(c *gin.Context) {
    var req dtoRequest.CreateCommentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "参数错误", err)
        return
    }

    comment := &models.Comment{
        ArticleID: req.ArticleID,
        Content:   req.Content,
        ParentID:  req.ParentID,  // ← 这里接收 parentId
    }
    
    // ... 其他逻辑
}
```

## 数据说明

### ParentID 的类型

- **前端传递**：`comment.id`（数字类型）
- **后端接收**：`*uint`（指针类型，允许为 nil）
- **顶级评论**：不传 `parentId` 或传 `null`
- **回复评论**：传递父评论的 `id`

### 前端传参示例

```json
// 顶级评论
{
  "articleId": 123,
  "content": "这是一条评论"
}

// 回复评论
{
  "articleId": 123,
  "content": "@张三 这是一条回复",
  "parentId": 456  // ← 被回复评论的 ID
}
```

## 完整交互流程

1. **显示评论列表**：`loadComments()` 从后端获取评论
2. **用户点击回复**：点击某个评论的"回复"按钮
3. **记录被回复评论**：`replyingTo.value = comment`
4. **显示回复框**：展开该评论下的回复输入框
5. **用户输入内容**：在 `replyText` 中输入回复内容
6. **提交回复**：调用 `handleSubmitReply(comment)`
7. **构建请求参数**：从 `comment.id` 获取 `parentId`
8. **发送请求**：`addComment({ articleId, content, parentId })`
9. **后端处理**：创建评论，`ParentID` 指向父评论
10. **刷新列表**：重新加载评论列表显示新回复

## 注意事项

1. **评论层级**：后端配置有最大回复层级限制（`MaxLevel`）
2. **ParentID 为 nil**：表示顶级评论
3. **重复回复**：可以对已回复过的评论继续回复（但受层级限制）
4. **删除父评论**：删除父评论时，子评论的处理逻辑（取决于业务需求）

## 数据库存储

```sql
-- 顶级评论
INSERT INTO comment (article_id, user_id, content, parent_id, ...)
VALUES (123, 1, '顶级评论', NULL, ...);

-- 回复评论
INSERT INTO comment (article_id, user_id, content, parent_id, ...)
VALUES (123, 2, '回复评论', 456, ...);  -- parent_id = 456
```
