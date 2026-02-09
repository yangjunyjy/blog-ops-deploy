package services

import (
	"my-blog-backend/internal/api/v1/dto/response"
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
)

type FavoriteService interface {
	CreateFolder(folder *models.FavoriteFolder) error
	UpdateFolder(folder *models.FavoriteFolder) error
	DeleteFolder(id uint, userID uint) error
	GetFolderByID(id uint, userID uint) (*models.FavoriteFolder, error)
	GetUserFolders(userID uint) ([]*models.FavoriteFolder, error)
	GetUserFoldersWithCount(userID uint) ([]response.FolderWithCount, error)
	GetUserFoldersWithCountPaged(userID uint, page, pageSize int) ([]response.FolderWithCount, int64, error)
	EnsureDefaultFolder(userID uint) error // 确保用户有默认文件夹

	AddFavorite(favorite *models.ArticleFavorite) error
	RemoveFavorite(articleID uint, userID uint) error
	MoveFavorite(articleID uint, userID uint, folderID uint) error
	GetUserFavorites(userID uint, folderID *uint, page, pageSize int) ([]*models.ArticleFavorite, int64, error)
	CheckFavorite(articleID uint, userID uint) (*models.ArticleFavorite, error)
	IncrementFavoriteCount(articleID uint) error
	DecrementFavoriteCount(articleID uint) error
}

type favoriteService struct {
	repo repository.FavoriteRepository
}

func NewFavoriteService(repo repository.FavoriteRepository) FavoriteService {
	return &favoriteService{repo: repo}
}

func (s *favoriteService) CreateFolder(folder *models.FavoriteFolder) error {
	return s.repo.CreateFolder(folder)
}

func (s *favoriteService) UpdateFolder(folder *models.FavoriteFolder) error {
	return s.repo.UpdateFolder(folder)
}

func (s *favoriteService) DeleteFolder(id uint, userID uint) error {
	return s.repo.DeleteFolder(id, userID)
}

func (s *favoriteService) GetFolderByID(id uint, userID uint) (*models.FavoriteFolder, error) {
	return s.repo.GetFolderByID(id, userID)
}

func (s *favoriteService) GetUserFolders(userID uint) ([]*models.FavoriteFolder, error) {
	return s.repo.GetUserFolders(userID)
}

// GetUserFoldersWithCount 获取用户的所有文件夹及每个文件夹的文章数量（优化版）
func (s *favoriteService) GetUserFoldersWithCount(userID uint) ([]response.FolderWithCount, error) {
	// 一次性获取文件夹和文章数量
	folders, countMap, err := s.repo.GetUserFoldersWithCounts(userID)
	if err != nil {
		return nil, err
	}

	// 组装带文章数量的文件夹列表
	result := make([]response.FolderWithCount, 0, len(folders))
	for _, folder := range folders {
		result = append(result, response.FolderWithCount{
			ID:           folder.ID,
			Name:         folder.Name,
			Description:  folder.Description,
			SortOrder:    folder.SortOrder,
			ArticleCount: countMap[folder.ID], // 直接从 map 获取，无需额外查询
			CreatedAt:    folder.CreatedAt,
			UpdatedAt:    folder.UpdatedAt,
		})
	}

	return result, nil
}

// GetUserFoldersWithCountPaged 分页获取用户的所有文件夹及每个文件夹的文章数量
func (s *favoriteService) GetUserFoldersWithCountPaged(userID uint, page, pageSize int) ([]response.FolderWithCount, int64, error) {
	// 分页获取文件夹和文章数量
	folders, countMap, total, err := s.repo.GetUserFoldersWithCountsPaged(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 组装带文章数量的文件夹列表
	result := make([]response.FolderWithCount, 0, len(folders))
	for _, folder := range folders {
		result = append(result, response.FolderWithCount{
			ID:           folder.ID,
			Name:         folder.Name,
			Description:  folder.Description,
			SortOrder:    folder.SortOrder,
			ArticleCount: countMap[folder.ID],
			CreatedAt:    folder.CreatedAt,
			UpdatedAt:    folder.UpdatedAt,
		})
	}

	return result, total, nil
}

// AddFavorite 添加收藏
func (s *favoriteService) AddFavorite(favorite *models.ArticleFavorite) error {
	return s.repo.AddFavorite(favorite)
}

// RemoveFavorite 移除收藏
func (s *favoriteService) RemoveFavorite(articleID uint, userID uint) error {
	return s.repo.RemoveFavorite(articleID, userID)
}

// MoveFavorite 移动收藏到其他文件夹
func (s *favoriteService) MoveFavorite(articleID uint, userID uint, folderID uint) error {
	return s.repo.MoveFavorite(articleID, userID, folderID)
}

// GetUserFavorites 获取用户的收藏列表
func (s *favoriteService) GetUserFavorites(userID uint, folderID *uint, page, pageSize int) ([]*models.ArticleFavorite, int64, error) {
	return s.repo.GetUserFavorites(userID, folderID, page, pageSize)
}

// CheckFavorite 检查是否已收藏
func (s *favoriteService) CheckFavorite(articleID uint, userID uint) (*models.ArticleFavorite, error) {
	return s.repo.CheckFavorite(articleID, userID)
}

// EnsureDefaultFolder 确保用户有默认文件夹
func (s *favoriteService) EnsureDefaultFolder(userID uint) error {
	// 检查用户是否有任何文件夹
	folders, err := s.repo.GetUserFolders(userID)
	if err != nil {
		return err
	}

	// 如果用户没有文件夹，创建默认文件夹
	if len(folders) == 0 {
		defaultFolder := &models.FavoriteFolder{
			UserID:      userID,
			Name:        "默认收藏夹",
			Description: "系统默认收藏夹，用于存放未分类的收藏",
			SortOrder:   0,
		}
		return s.repo.CreateFolder(defaultFolder)
	}

	return nil
}

// IncrementFavoriteCount 增加文章的收藏数
func (s *favoriteService) IncrementFavoriteCount(articleID uint) error {
	return s.repo.IncrementFavoriteCount(articleID)
}

// DecrementFavoriteCount 减少文章的收藏数
func (s *favoriteService) DecrementFavoriteCount(articleID uint) error {
	return s.repo.DecrementFavoriteCount(articleID)
}
