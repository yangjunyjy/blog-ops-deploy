package mysql

import (
	"gorm.io/gorm"

	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
)

// FavoriteRepositoryImpl 收藏仓储实现
type FavoriteRepositoryImpl struct {
	db *gorm.DB
}

// NewFavoriteRepositoryImpl 创建收藏仓储实例
func NewFavoriteRepositoryImpl(db *gorm.DB) repository.FavoriteRepository {
	return &FavoriteRepositoryImpl{db: db}
}

// CreateFolder 创建收藏文件夹
func (r *FavoriteRepositoryImpl) CreateFolder(folder *models.FavoriteFolder) error {
	return r.db.Create(folder).Error
}

// UpdateFolder 更新收藏文件夹
func (r *FavoriteRepositoryImpl) UpdateFolder(folder *models.FavoriteFolder) error {
	return r.db.Save(folder).Error
}

// DeleteFolder 删除收藏文件夹（软删除）
func (r *FavoriteRepositoryImpl) DeleteFolder(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.FavoriteFolder{}).Error
}

// GetFolderByID 根据ID和用户ID获取文件夹
func (r *FavoriteRepositoryImpl) GetFolderByID(id uint, userID uint) (*models.FavoriteFolder, error) {
	var folder models.FavoriteFolder
	err := r.db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).First(&folder).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

// GetUserFolders 获取用户的所有收藏文件夹
func (r *FavoriteRepositoryImpl) GetUserFolders(userID uint) ([]*models.FavoriteFolder, error) {
	var folders []*models.FavoriteFolder
	err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("sort_order ASC, created_at DESC").
		Find(&folders).Error
	return folders, err
}

// AddFavorite 添加收藏
func (r *FavoriteRepositoryImpl) AddFavorite(favorite *models.ArticleFavorite) error {
	return r.db.Create(favorite).Error
}

// RemoveFavorite 移除收藏
func (r *FavoriteRepositoryImpl) RemoveFavorite(articleID uint, userID uint) error {
	return r.db.Where("article_id = ? AND user_id = ?", articleID, userID).Delete(&models.ArticleFavorite{}).Error
}

// MoveFavorite 移动收藏到其他文件夹
func (r *FavoriteRepositoryImpl) MoveFavorite(articleID uint, userID uint, folderID uint) error {
	return r.db.Model(&models.ArticleFavorite{}).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Update("folder_id", folderID).Error
}

// GetUserFavorites 获取用户的收藏列表
func (r *FavoriteRepositoryImpl) GetUserFavorites(userID uint, folderID *uint, page, pageSize int) ([]*models.ArticleFavorite, int64, error) {
	var favorites []*models.ArticleFavorite
	var total int64

	query := r.db.Model(&models.ArticleFavorite{}).
		Where("user_id = ?", userID).
		Where("article_id IN (SELECT id FROM articles)") // 只返回存在的文章

	if folderID != nil {
		query = query.Where("folder_id = ?", *folderID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&favorites).Error

	if err != nil {
		return nil, 0, err
	}

	// 手动加载 Article 基本信息
	articleIDs := make([]uint64, len(favorites))
	for i, fav := range favorites {
		articleIDs[i] = uint64(fav.ArticleID)
	}

	var articles []models.Article
	if len(articleIDs) > 0 {
		err = r.db.Where("id IN ?", articleIDs).Find(&articles).Error
		if err != nil {
			return favorites, total, nil // 即使文章加载失败，也返回收藏列表
		}
	}

	// 将文章数据填充到收藏中
	articleMap := make(map[uint64]*models.Article)
	for i := range articles {
		articleMap[articles[i].ID] = &articles[i]
	}

	for _, fav := range favorites {
		if article, ok := articleMap[uint64(fav.ArticleID)]; ok {
			fav.Article = article
		}
	}

	return favorites, total, err
}

// CheckFavorite 检查是否已收藏
func (r *FavoriteRepositoryImpl) CheckFavorite(articleID uint, userID uint) (*models.ArticleFavorite, error) {
	var favorite models.ArticleFavorite
	err := r.db.Where("article_id = ? AND user_id = ?", articleID, userID).First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

// GetFolderArticleCount 获取文件夹的文章数量
func (r *FavoriteRepositoryImpl) GetFolderArticleCount(folderID uint) (int64, error) {
	var count int64
	// 使用子查询只统计仍然存在的文章收藏
	err := r.db.Model(&models.ArticleFavorite{}).
		Where("folder_id = ?", folderID).
		Where("article_id IN (SELECT id FROM articles)"). // 只统计存在的文章
		Count(&count).Error
	return count, err
}

// GetUserFolderCount 获取用户的文件夹数量
func (r *FavoriteRepositoryImpl) GetUserFolderCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.FavoriteFolder{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

// GetUserFoldersWithCounts 获取用户的所有文件夹及每个文件夹的文章数量（优化版，避免 N+1 查询）
func (r *FavoriteRepositoryImpl) GetUserFoldersWithCounts(userID uint) ([]*models.FavoriteFolder, map[uint]int64, error) {
	var folders []*models.FavoriteFolder

	// 获取所有文件夹
	err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("sort_order ASC, created_at DESC").
		Find(&folders).Error
	if err != nil {
		return nil, nil, err
	}

	// 批量获取所有文件夹的文章数量
	folderIDs := make([]uint, len(folders))
	for i, folder := range folders {
		folderIDs[i] = folder.ID
	}

	// 使用子查询一次性获取所有文件夹的文章数量，只统计存在的文章
	type FolderCount struct {
		FolderID uint  `json:"folder_id"`
		Count    int64 `json:"count"`
	}

	var counts []FolderCount
	err = r.db.Table("article_favorites").
		Select("article_favorites.folder_id, COUNT(*) as count").
		Where("article_favorites.folder_id IN ?", folderIDs).
		Where("article_favorites.article_id IN (SELECT id FROM articles)"). // 只统计存在的文章
		Group("article_favorites.folder_id").
		Find(&counts).Error
	if err != nil {
		return nil, nil, err
	}

	// 将结果转换为 map
	countMap := make(map[uint]int64)
	for _, c := range counts {
		countMap[c.FolderID] = c.Count
	}

	return folders, countMap, nil
}

// GetUserFoldersWithCountsPaged 分页获取用户的文件夹及数量
func (r *FavoriteRepositoryImpl) GetUserFoldersWithCountsPaged(userID uint, page, pageSize int) ([]*models.FavoriteFolder, map[uint]int64, int64, error) {
	var folders []*models.FavoriteFolder
	var total int64

	// 获取总数
	err := r.db.Model(&models.FavoriteFolder{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&total).Error
	if err != nil {
		return nil, nil, 0, err
	}

	// 分页查询文件夹
	offset := (page - 1) * pageSize
	err = r.db.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&folders).Error
	if err != nil {
		return nil, nil, 0, err
	}

	// 如果没有文件夹，直接返回
	if len(folders) == 0 {
		return folders, make(map[uint]int64), total, nil
	}

	// 批量获取当前页文件夹的文章数量
	folderIDs := make([]uint, len(folders))
	for i, folder := range folders {
		folderIDs[i] = folder.ID
	}

	// 使用子查询一次性获取所有文件夹的文章数量，只统计存在的文章
	type FolderCount struct {
		FolderID uint  `json:"folder_id"`
		Count    int64 `json:"count"`
	}

	var counts []FolderCount
	err = r.db.Table("article_favorites").
		Select("article_favorites.folder_id, COUNT(*) as count").
		Where("article_favorites.folder_id IN ?", folderIDs).
		Where("article_favorites.article_id IN (SELECT id FROM articles)"). // 只统计存在的文章
		Group("article_favorites.folder_id").
		Find(&counts).Error
	if err != nil {
		return nil, nil, 0, err
	}

	// 将结果转换为 map
	countMap := make(map[uint]int64)
	for _, c := range counts {
		countMap[c.FolderID] = c.Count
	}

	return folders, countMap, total, nil
}

func (r *FavoriteRepositoryImpl) DeleteByArticleID(ArticleId uint) error {
	return r.db.Where("article_id = ?", ArticleId).Delete(&models.ArticleFavorite{}).Error
}

// IncrementFavoriteCount 增加文章的收藏数
func (r *FavoriteRepositoryImpl) IncrementFavoriteCount(articleID uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", articleID).UpdateColumn("favorites", gorm.Expr("favorites + 1")).Error
}

// DecrementFavoriteCount 减少文章的收藏数
func (r *FavoriteRepositoryImpl) DecrementFavoriteCount(articleID uint) error {
	return r.db.Model(&models.Article{}).Where("id = ? AND favorites > 0", articleID).UpdateColumn("favorites", gorm.Expr("favorites - 1")).Error
}
