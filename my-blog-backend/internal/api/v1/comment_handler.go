package api

import (
	"net/http"
	"strconv"
	"strings"

	dtoRequest "my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/pkg/claims"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/pkg/middleware"
	"my-blog-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService services.CommentService
}

func NewCommentHandler(commentService services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateComment 创建评论
// @Summary 创建评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dtoRequest.CreateCommentRequest true "创建评论请求"
// @Success 200 {object} response.Response
// @Router /comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req dtoRequest.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("创建评论参数错误", logger.Any("err", err))
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 检查评论功能是否启用
	if !h.commentService.IsCommentEnabled() {
		response.Error(c, http.StatusForbidden, "评论功能未启用", nil)
		return
	}

	comment := &models.Comment{
		ArticleID: req.ArticleID,
		Content:   req.Content,
		ParentID:  req.ParentID,
	}

	// 从上下文获取用户ID
	userClaims, exists := c.Get("user_claims")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	claims := userClaims.(*claims.UserClaims)

	// 使用带审核逻辑的创建方法
	newComment, err := h.commentService.CreateCommentWithAudit(comment, claims.UserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	// 从数据库重新加载评论以获取关联的用户信息
	// 注意：newComment 已经在 CreateCommentWithAudit 中过滤了敏感词
	// 这里重新加载后需要再次过滤
	newComment, err = h.commentService.GetComment(newComment.ID)
	if err != nil {
		logger.Error("获取评论详情失败", logger.Any("err", err))
	}

	// 构建返回数据
	result := gin.H{
		"comment": newComment,
		"status":  newComment.Status,
		"message": "",
	}

	// 添加作者信息
	if newComment.User != nil {
		result["user"] = gin.H{
			"id":       newComment.User.ID,
			"nickname": newComment.User.Nickname,
			"username": newComment.User.Username,
			"avatar":   newComment.User.Avatar,
		}
	}

	// 如果是回复评论，添加被回复的用户信息
	if newComment.Parent != nil && newComment.Parent.User != nil {
		result["reply_to"] = newComment.ParentID
		result["reply_to_user"] = gin.H{
			"id":       newComment.Parent.User.ID,
			"nickname": newComment.Parent.User.Nickname,
			"username": newComment.Parent.User.Username,
			"avatar":   newComment.Parent.User.Avatar,
		}
	}

	if newComment.Status == 0 {
		result["message"] = "评论已提交，等待审核"
	} else {
		result["message"] = "评论发布成功"
	}

	response.Success(c, result, "")
}

// UpdateComment 更新评论
// @Summary 更新评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Param request body dtoRequest.CreateCommentRequest true "更新评论请求"
// @Success 200 {object} response.Response
// @Router /comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID", err)
		return
	}

	var req dtoRequest.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	comment := &models.Comment{
		ArticleID: req.ArticleID,
		Content:   req.Content,
		ParentID:  req.ParentID,
	}

	if err := h.commentService.UpdateComment(uint(id), comment); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新评论失败", err)
		return
	}

	response.Success(c, comment, "更新成功")
}

// DeleteComment 删除评论
// @Summary 删除评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Router /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID", err)
		return
	}

	// 从上下文获取用户ID
	userClaims, exists := c.Get("user_claims")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	claims := userClaims.(*claims.UserClaims)

	// 删除评论及其所有子评论（带权限验证）
	if err := h.commentService.DeleteCommentWithChildrenByUser(uint(id), claims.UserID); err != nil {
		response.Error(c, http.StatusForbidden, err.Error(), err)
		return
	}

	response.Success(c, nil, "删除成功")
}

// GetReplies 获取子评论列表
// @Summary 获取子评论列表
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param id path int true "父评论ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /comments/{id}/replies [get]
func (h *CommentHandler) GetReplies(c *gin.Context) {
	parentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	replies, total, err := h.commentService.ListReplies(uint(parentID), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取子评论失败", err)
		return
	}

	// 尝试从token获取当前用户ID（可选的，用于判断是否点赞）
	var currentUserID uint
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			token := parts[1]
			// 验证token但不强制要求
			if userClaims, err := middleware.ValidateToken(token); err == nil {
				currentUserID = userClaims.UserID
			}
		}
	}

	// 转换为前端需要的响应格式
	items := make([]response.CommentItemFull, 0, len(replies))
	for _, comment := range replies {
		item := response.CommentItemFull{
			ID:        comment.ID,
			ArticleID: comment.ArticleID,
			UserID:    comment.UserID,
			ParentID:  comment.ParentID,
			Content:   comment.Content,
			Status:    comment.Status,
			Likes:     comment.Likes,
			IPAddress: comment.IPAddress,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}

		// 检查当前用户是否已点赞
		if currentUserID > 0 {
			isLiked, _ := h.commentService.CheckUserLiked(comment.ID, currentUserID)
			item.IsLiked = isLiked
		}

		// 设置作者信息
		if comment.User != nil {
			name := comment.User.Nickname
			if name == "" {
				name = comment.User.Username
			}

			item.Author = &response.CommentAuthorWithInfo{
				ID:     comment.User.ID,
				Name:   name,
				Avatar: comment.User.Avatar,
			}
		}

		// 设置被回复的用户信息（如果有）
		if comment.Parent != nil && comment.Parent.User != nil {
			item.ReplyTo = comment.Parent.ID
			replyToName := comment.Parent.User.Nickname
			if replyToName == "" {
				replyToName = comment.Parent.User.Username
			}
			item.ReplyToUser = &response.CommentAuthorWithInfo{
				ID:     comment.Parent.User.ID,
				Name:   replyToName,
				Avatar: comment.Parent.User.Avatar,
			}
		}

		items = append(items, item)
	}

	response.Success(c, gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
		"parent_id": parentID,
	}, "")
}

// GetReplyCount 获取子评论数量
// @Summary 获取子评论数量
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param id path int true "父评论ID"
// @Success 200 {object} response.Response
// @Router /comments/{id}/replies/count [get]
func (h *CommentHandler) GetReplyCount(c *gin.Context) {
	parentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID", err)
		return
	}

	count, err := h.commentService.GetReplyCount(uint(parentID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取子评论数量失败", err)
		return
	}

	response.Success(c, gin.H{
		"count":     count,
		"parent_id": parentID,
	}, "")
}

// GetComment 获取评论详情
// @Summary 获取评论详情
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Router /comments/{id} [get]
func (h *CommentHandler) GetComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID", err)
		return
	}

	comment, err := h.commentService.GetComment(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "评论不存在", err)
		return
	}

	response.Success(c, comment, "")
}

// ListComments 获取评论列表
// @Summary 获取文章评论列表
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /articles/{id}/comments [get]
func (h *CommentHandler) ListComments(c *gin.Context) {
	var req dtoRequest.CommentListRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 尝试从token获取当前用户ID（可选的，用于判断是否点赞和显示自己的待审核评论）
	var currentUserID uint
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			token := parts[1]
			// 验证token但不强制要求
			if userClaims, err := middleware.ValidateToken(token); err == nil {
				currentUserID = userClaims.UserID
			}
		}
	}

	// 使用新方法：包含当前用户的待审核评论
	comments, total, err := h.commentService.ListCommentsWithPending(req.ArticleID, currentUserID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取评论列表失败", err)
		return
	}

	// 转换为前端需要的响应格式（包含完整的评论数据）
	items := make([]response.CommentItemFull, 0, len(comments))
	for _, comment := range comments {
		item := response.CommentItemFull{
			ID:        comment.ID,
			ArticleID: comment.ArticleID,
			UserID:    comment.UserID,
			ParentID:  comment.ParentID,
			Content:   comment.Content,
			Status:    comment.Status,
			Likes:     comment.Likes,
			IPAddress: comment.IPAddress,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}

		// 检查当前用户是否已点赞
		if currentUserID > 0 {
			isLiked, _ := h.commentService.CheckUserLiked(comment.ID, currentUserID)
			item.IsLiked = isLiked
			logger.Info("isLiked", logger.Any("isLiked", isLiked))
		}
		logger.Info("item", logger.Any("currentUserID", currentUserID))
		// 设置作者信息
		if comment.User != nil {
			// 处理用户名：优先使用 nickname，其次使用 username
			name := comment.User.Nickname
			if name == "" {
				name = comment.User.Username
			}

			item.Author = &response.CommentAuthorWithInfo{
				ID:     comment.User.ID,
				Name:   name,
				Avatar: comment.User.Avatar,
			}
		}

		items = append(items, item)
	}

	commentListResponse := response.CommentListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	response.Success(c, commentListResponse, "")
}

// ListUserComments 获取用户评论列表
// @Summary 获取用户评论列表
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /users/{user_id}/comments [get]
func (h *CommentHandler) ListUserComments(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	comments, total, err := h.commentService.ListUserComments(uint(userID), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户评论列表失败", err)
		return
	}

	response.Success(c, gin.H{
		"items":     comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, "")
}

// ListPendingComments 获取待审核评论列表（管理员）
// @Summary 获取待审核评论列表（管理员）
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /admin/comments/pending [get]
func (h *CommentHandler) ListPendingComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	comments, total, err := h.commentService.ListPendingComments(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取待审核评论列表失败", err)
		return
	}

	response.Success(c, gin.H{
		"items":     comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, "")
}

// ListAllComments 获取所有评论列表（管理员）
// @Summary 获取所有评论列表（管理员）
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /admin/comments [get]
func (h *CommentHandler) ListAllComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	comments, total, err := h.commentService.ListAllComments(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取评论列表失败", err)
		return
	}

	response.Success(c, gin.H{
		"items":     comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, "")
}

// ApproveComment 审核通过评论
// @Summary 审核通过评论（管理员）
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Router /admin/comments/{id}/approve [put]
func (h *CommentHandler) ApproveComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID", err)
		return
	}

	if err := h.commentService.ApproveComment(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "审核失败", err)
		return
	}

	response.Success(c, nil, "审核通过")
}

// RejectComment 审核拒绝评论
// @Summary 审核拒绝评论（管理员）
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Router /admin/comments/{id}/reject [put]
func (h *CommentHandler) RejectComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID", err)
		return
	}

	if err := h.commentService.RejectComment(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "审核失败", err)
		return
	}

	response.Success(c, nil, "审核拒绝")
}

// LikeComment 点赞评论
// @Summary 点赞评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Router /comments/{id}/like [post]
func (h *CommentHandler) LikeComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID", err)
		return
	}

	// 从上下文获取用户ID
	userClaims, exists := c.Get("user_claims")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	claims := userClaims.(*claims.UserClaims)

	if err := h.commentService.LikeComment(uint(commentID), claims.UserID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	response.Success(c, nil, "点赞成功")
}

// UnlikeComment 取消点赞评论
// @Summary 取消点赞评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Router /comments/{id}/unlike [post]
func (h *CommentHandler) UnlikeComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID", err)
		return
	}

	// 从上下文获取用户ID
	userClaims, exists := c.Get("user_claims")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	claims := userClaims.(*claims.UserClaims)

	if err := h.commentService.UnlikeComment(uint(commentID), claims.UserID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	response.Success(c, nil, "取消点赞成功")
}
