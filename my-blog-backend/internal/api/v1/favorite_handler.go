package api

import (
	"strconv"

	"my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	favoriteService services.FavoriteService
}

func NewFavoriteHandler(favoriteService services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
	}
}

// GetUserFoldersWithCount 获取用户的所有收藏文件夹及文章数量
// GET /api/v1/front/favorite/userFolder?page=1&pageSize=10
func (h *FavoriteHandler) GetUserFoldersWithCount(c *gin.Context) {
	// 获取当前用户ID（从中间件设置的 Context 中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 确保用户有默认文件夹
	if err := h.favoriteService.EnsureDefaultFolder(userID.(uint)); err != nil {
		// 只记录错误，不影响后续流程
		// log.Error("Ensure default folder failed", err)
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 调用 Service 获取文件夹列表
	folders, total, err := h.favoriteService.GetUserFoldersWithCountPaged(userID.(uint), page, pageSize)
	if err != nil {
		response.InternalError(c, "获取文件夹失败", err)
		return
	}

	// 返回分页数据
	response.Success(c, response.FolderListResponse{
		List:     folders,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "")
}

// CreateFolder 创建收藏文件夹
// POST /api/v1/front/favorite/create
func (h *FavoriteHandler) CreateFolder(c *gin.Context) {
	// 获取当前用户ID（从中间件设置的 Context 中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 解析请求参数
	var req request.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), err)
		return
	}

	// 创建文件夹
	folder := &models.FavoriteFolder{
		UserID:      userID.(uint),
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}

	if err := h.favoriteService.CreateFolder(folder); err != nil {
		response.InternalError(c, "创建文件夹失败", err)
		return
	}

	response.Success(c, folder, "创建成功")
}

// UpdateFolder 更新收藏文件夹
// PUT /api/v1/front/favorite/userFolder/:id
func (h *FavoriteHandler) UpdateFolder(c *gin.Context) {
	// 获取当前用户ID（从中间件设置的 Context 中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 获取文件夹ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文件夹ID", err)
		return
	}

	// 解析请求参数
	var req request.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), err)
		return
	}

	// 获取原文件夹
	folder, err := h.favoriteService.GetFolderByID(uint(id), userID.(uint))
	if err != nil {
		response.BadRequest(c, "文件夹不存在或无权访问", err)
		return
	}

	// 更新文件夹
	folder.Name = req.Name
	folder.Description = req.Description
	folder.SortOrder = req.SortOrder

	if err := h.favoriteService.UpdateFolder(folder); err != nil {
		response.InternalError(c, "更新文件夹失败", err)
		return
	}

	response.Success(c, folder, "更新成功")
}

// DeleteFolder 删除收藏文件夹
// DELETE /api/v1/front/favorite/userFolder/:id
func (h *FavoriteHandler) DeleteFolder(c *gin.Context) {
	// 获取当前用户ID（从中间件设置的 Context 中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 获取文件夹ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文件夹ID", err)
		return
	}

	// 删除文件夹
	if err := h.favoriteService.DeleteFolder(uint(id), userID.(uint)); err != nil {
		response.InternalError(c, "删除文件夹失败", err)
		return
	}

	response.Success(c, nil, "删除成功")
}

// GetUserFavorites 获取用户的收藏列表
// GET /api/v1/front/favorite/list
func (h *FavoriteHandler) GetUserFavorites(c *gin.Context) {
	// 获取当前用户ID（从中间件设置的 Context 中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 获取文件夹ID（可选）
	var folderID *uint
	if folderIDStr := c.Query("folderId"); folderIDStr != "" {
		id, err := strconv.ParseUint(folderIDStr, 10, 32)
		if err == nil {
			fid := uint(id)
			folderID = &fid
		}
	}

	// 获取收藏列表
	favorites, total, err := h.favoriteService.GetUserFavorites(userID.(uint), folderID, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取收藏列表失败", err)
		return
	}

	// 转换为 DTO
	var list []response.FavoriteArticleItem
	for _, fav := range favorites {
		if fav.Article != nil {
			article := fav.Article
			authorInfo := response.AuthorInfo{}
			if article.Author != nil {
				var avatar *string
				if article.Author.Avatar != "" {
					avatar = &article.Author.Avatar
				}
				authorInfo = response.AuthorInfo{
					ID:     article.Author.ID,
					Name:   article.Author.Nickname,
					Avatar: avatar,
				}
			}

			var cover *string
			if article.Cover != "" {
				cover = &article.Cover
			}

			list = append(list, response.FavoriteArticleItem{
				ID:        uint(article.ID),
				Title:     article.Title,
				Summary:   article.Summary,
				Cover:     cover,
				Author:    authorInfo,
				CreatedAt: fav.CreatedAt,
			})
		}
	}

	response.Success(c, response.FavoriteListResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "")
}

// AddFavorite 添加收藏
// POST /api/v1/front/favorite/add
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	// 获取当前用户ID（从中间件设置的 Context 中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req request.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err)
		return
	}

	// 检查是否已收藏
	_, err2 := h.favoriteService.CheckFavorite(req.ArticleID, userID.(uint))
	if err2 == nil {
		response.BadRequest(c, "已收藏该文章", nil)
		return
	}

	// 创建收藏
	favorite := &models.ArticleFavorite{
		ArticleID: req.ArticleID,
		UserID:    userID.(uint),
		FolderID:  req.FolderID,
	}

	if err := h.favoriteService.AddFavorite(favorite); err != nil {
		response.InternalError(c, "收藏失败", err)
		return
	}

	// 更新文章的收藏数
	if err := h.favoriteService.IncrementFavoriteCount(req.ArticleID); err != nil {
		// 只记录错误，不影响收藏操作
		// log.Error("Increment favorite count failed", err)
	}

	response.Success(c, nil, "收藏成功")
}

// RemoveFavorite 移除收藏
// POST /api/v1/front/favorite/remove
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	// 获取当前用户ID（从中间件设置的 Context 中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 获取文章ID（从请求体获取）
	var req struct {
		ArticleID uint `json:"articleId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的文章ID", err)
		return
	}

	// 移除收藏
	if err := h.favoriteService.RemoveFavorite(req.ArticleID, userID.(uint)); err != nil {
		response.InternalError(c, "取消收藏失败", err)
		return
	}

	// 更新文章的收藏数
	if err := h.favoriteService.DecrementFavoriteCount(req.ArticleID); err != nil {
		// 只记录错误，不影响移除操作
		// log.Error("Decrement favorite count failed", err)
	}

	response.Success(c, nil, "取消收藏成功")
}

// MoveFavorite 移动收藏到其他文件夹
// POST /api/v1/front/favorite/move
func (h *FavoriteHandler) MoveFavorite(c *gin.Context) {
	// 获取当前用户ID（从中间件设置的 Context 中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 解析请求参数
	var req struct {
		ArticleID uint `json:"articleId" binding:"required"`
		FolderID  uint `json:"folderId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), err)
		return
	}

	// 移动收藏
	if err := h.favoriteService.MoveFavorite(req.ArticleID, userID.(uint), req.FolderID); err != nil {
		response.InternalError(c, "移动失败", err)
		return
	}

	response.Success(c, nil, "移动成功")
}

// CheckFavoriteStatus 检查文章收藏状态
// GET /api/v1/front/favorite/check/:articleId
func (h *FavoriteHandler) CheckFavoriteStatus(c *gin.Context) {
	// 获取当前用户ID（从中间件设置的 Context 中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	// 获取文章ID
	articleIDStr := c.Param("articleId")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章ID", err)
		return
	}

	// 检查是否已收藏
	favorite, err := h.favoriteService.CheckFavorite(uint(articleID), userID.(uint))
	if err != nil {
		// 未收藏不是错误，返回未收藏状态
		response.Success(c, map[string]interface{}{
			"isFavorited": false,
			"folderId":    nil,
		}, "未收藏")
		return
	}

	response.Success(c, map[string]interface{}{
		"isFavorited": true,
		"folderId":    favorite.FolderID,
	}, "已收藏")
}
