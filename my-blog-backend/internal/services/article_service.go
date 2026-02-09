package services

import (
	"errors"
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
	"time"
)

type ArticleService interface {
	CreateArticle(article *models.Article, tagIDs []uint) error
	UpdateArticle(id uint, article *models.Article, tagIDs []uint) error
	DeleteArticle(id uint) error
	GetArticle(id uint) (*models.Article, error)
	ListArticles(page, pageSize int, status *int, categoryID *uint) ([]*models.Article, int64, error)
	SearchArticles(keyword string, page, pageSize int) ([]*models.Article, int64, error)
	ViewArticle(id uint) error
	GetHotArticles(limit int) ([]*models.Article, error)
	GetRecentArticles(limit int) ([]*models.Article, error)
	LikeArticle(articleID, userID uint) error
	UnlikeArticle(articleID, userID uint) error
	CheckArticleLiked(articleID, userID uint) (bool, error)
	FavoriteArticle(articleID, userID uint) error
	UnfavoriteArticle(articleID, userID uint) error
}

type articleService struct {
	articleRepo     repository.ArticleRepository
	tagRepo         repository.TagRepository
	commentRepo     repository.CommentRepository
	commentLikeRepo repository.CommentLikeRepository
	favoriteRepo    repository.FavoriteRepository
	articleLikeRepo repository.ArticleLikeRepository
	activityService *UserActivityService
}

func NewArticleService(articleRepo repository.ArticleRepository,
	tagRepo repository.TagRepository,
	commentRepo repository.CommentRepository,
	commentLikeRepo repository.CommentLikeRepository,
	favoriteRepo repository.FavoriteRepository,
	articleLikeRepo repository.ArticleLikeRepository,
	activityService *UserActivityService) ArticleService {
	return &articleService{
		articleRepo:     articleRepo,
		tagRepo:         tagRepo,
		commentRepo:     commentRepo,
		commentLikeRepo: commentLikeRepo,
		favoriteRepo:    favoriteRepo,
		articleLikeRepo: articleLikeRepo,
		activityService: activityService,
	}
}

func (s *articleService) CreateArticle(article *models.Article, tagIDs []uint) error {
	now := time.Now()
	article.CreatedAt = now
	article.UpdatedAt = now

	if article.Status == 1 {
		article.PublishedAt = &now
	}

	if err := s.articleRepo.Create(article); err != nil {
		return err
	}

	// TODO: 关联标签
	return nil
}

func (s *articleService) UpdateArticle(id uint, article *models.Article, tagIDs []uint) error {
	existing, err := s.articleRepo.GetByID(id)
	if err != nil {
		return errors.New("article not found")
	}

	article.ID = uint64(id)
	article.CreatedAt = existing.CreatedAt
	article.UpdatedAt = time.Now()

	if article.Status == 1 && existing.Status != 1 {
		now := time.Now()
		article.PublishedAt = &now
	}

	if err := s.articleRepo.Update(article); err != nil {
		return err
	}

	// TODO: 更新标签关联
	return nil
}

func (s *articleService) DeleteArticle(id uint) error {
	// TODO: 删除标签关联
	if err := s.articleRepo.Delete(id); err != nil {
		return err
	}
	// TODO: 删除文章所有评论
	if err := s.commentRepo.DeleteCommentByArticleID(id); err != nil {
		return err
	}

	// TODO: 删除文章所有收藏记录
	if err := s.favoriteRepo.DeleteByArticleID(id); err != nil {
		return err
	}
	return s.articleRepo.Delete(id)
}

func (s *articleService) GetArticle(id uint) (*models.Article, error) {
	return s.articleRepo.GetByID(id)
}

func (s *articleService) ListArticles(page, pageSize int, status *int, categoryID *uint) ([]*models.Article, int64, error) {
	return s.articleRepo.List(page, pageSize, status, categoryID)
}

func (s *articleService) SearchArticles(keyword string, page, pageSize int) ([]*models.Article, int64, error) {
	return s.articleRepo.Search(keyword, page, pageSize)
}

func (s *articleService) ViewArticle(id uint) error {
	return s.articleRepo.IncrementViewCount(id)
}

func (s *articleService) GetHotArticles(limit int) ([]*models.Article, error) {
	return s.articleRepo.GetHotArticles(limit)
}

func (s *articleService) GetRecentArticles(limit int) ([]*models.Article, error) {
	return s.articleRepo.GetRecentArticles(limit)
}

func (s *articleService) LikeArticle(articleID, userID uint) error {
	// 创建点赞记录
	like := &models.ArticleLike{
		ArticleID: uint64(articleID),
		UserID:    uint64(userID),
	}
	if err := s.articleLikeRepo.Create(like); err != nil {
		return err
	}

	// 更新文章点赞数
	if err := s.articleRepo.IncrementLikeCount(articleID); err != nil {
		return err
	}

	// 记录用户活动（点赞）
	go func() {
		_ = s.activityService.RecordLike(userID, articleID)
	}()

	return nil
}

func (s *articleService) UnlikeArticle(articleID, userID uint) error {
	// 删除点赞记录
	if err := s.articleLikeRepo.Delete(articleID, userID); err != nil {
		return err
	}

	// 更新文章点赞数
	return s.articleRepo.DecrementLikeCount(articleID)
}

func (s *articleService) CheckArticleLiked(articleID, userID uint) (bool, error) {
	return s.articleLikeRepo.CheckUserLiked(articleID, userID)
}

func (s *articleService) FavoriteArticle(articleID, userID uint) error {
	// 添加到默认收藏夹
	favorite := &models.ArticleFavorite{
		ArticleID: articleID,
		UserID:    userID,
		FolderID:  0, // 默认收藏夹
	}
	if err := s.favoriteRepo.AddFavorite(favorite); err != nil {
		return err
	}

	// 更新文章收藏数
	if err := s.favoriteRepo.IncrementFavoriteCount(articleID); err != nil {
		return err
	}

	// 记录用户活动（收藏到默认收藏夹）
	go func() {
		_ = s.activityService.RecordFavorite(userID, articleID, "默认收藏夹")
	}()

	return nil
}

func (s *articleService) UnfavoriteArticle(articleID, userID uint) error {
	// 删除收藏记录
	if err := s.favoriteRepo.RemoveFavorite(articleID, userID); err != nil {
		return err
	}

	// 更新文章收藏数
	return s.favoriteRepo.DecrementFavoriteCount(articleID)
}
