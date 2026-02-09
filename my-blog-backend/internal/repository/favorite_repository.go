package repository

import (
	models "my-blog-backend/internal/models/frontendModel"
)

// FavoriteRepository 收藏仓储接口
type FavoriteRepository interface {
	// 文件夹相关
	CreateFolder(folder *models.FavoriteFolder) error
	UpdateFolder(folder *models.FavoriteFolder) error
	DeleteFolder(id uint, userID uint) error
	GetFolderByID(id uint, userID uint) (*models.FavoriteFolder, error)
	GetUserFolders(userID uint) ([]*models.FavoriteFolder, error)

	// 收藏相关
	AddFavorite(favorite *models.ArticleFavorite) error
	RemoveFavorite(articleID uint, userID uint) error
	MoveFavorite(articleID uint, userID uint, folderID uint) error
	GetUserFavorites(userID uint, folderID *uint, page, pageSize int) ([]*models.ArticleFavorite, int64, error)
	CheckFavorite(articleID uint, userID uint) (*models.ArticleFavorite, error)

	// 统计相关
	GetFolderArticleCount(folderID uint) (int64, error)
	GetUserFolderCount(userID uint) (int64, error)

	// 批量统计 - 优化版
	GetUserFoldersWithCounts(userID uint) ([]*models.FavoriteFolder, map[uint]int64, error)

	// 分页获取用户的文件夹及数量
	GetUserFoldersWithCountsPaged(userID uint, page, pageSize int) ([]*models.FavoriteFolder, map[uint]int64, int64, error)

	// 删除文章所有收藏
	DeleteByArticleID(ArticleId uint) error

	// 更新文章收藏数
	IncrementFavoriteCount(articleID uint) error
	DecrementFavoriteCount(articleID uint) error
}
