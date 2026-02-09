package services

import (
	"errors"
	"my-blog-backend/internal/config"
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/repository"
	"time"
)

type CommentService interface {
	CreateComment(comment *models.Comment) error
	CreateCommentWithAudit(comment *models.Comment, userID uint) (*models.Comment, error)
	UpdateComment(id uint, comment *models.Comment) error
	DeleteComment(id uint) error
	DeleteCommentWithChildren(id uint) error               // 删除评论及其所有子评论
	DeleteCommentWithChildrenByUser(id, userID uint) error // 删除评论及其所有子评论（需要验证用户权限）
	GetComment(id uint) (*models.Comment, error)
	ListComments(articleID uint, page, pageSize int) ([]*models.Comment, int64, error)
	ListCommentsWithPending(articleID, userID uint, page, pageSize int) ([]*models.Comment, int64, error) // 包含当前用户的待审核评论
	ListReplies(parentID uint, page, pageSize int) ([]*models.Comment, int64, error)                      // 获取子评论列表
	GetReplyCount(parentID uint) (int64, error)                                                           // 获取子评论数量
	ListUserComments(userID uint, page, pageSize int) ([]*models.Comment, int64, error)
	ListPendingComments(page, pageSize int) ([]*models.Comment, int64, error)
	ListAllComments(page, pageSize int) ([]*models.Comment, int64, error)
	ApproveComment(id uint) error
	RejectComment(id uint) error
	IsCommentEnabled() bool
	// 评论点赞相关
	LikeComment(commentID, userID uint) error
	UnlikeComment(commentID, userID uint) error
	CheckUserLiked(commentID, userID uint) (bool, error)
	RefreshCommentLikeCount(commentID uint) error
}

type commentService struct {
	commentRepo          repository.CommentRepository
	commentLikeRepo      repository.CommentLikeRepository
	config               *config.Config
	sensitiveWordService SensitiveWordService
	activityService      *UserActivityService
}

func NewCommentService(commentRepo repository.CommentRepository, commentLikeRepo repository.CommentLikeRepository, cfg *config.Config, swService SensitiveWordService, activityService *UserActivityService) CommentService {
	return &commentService{
		commentRepo:          commentRepo,
		commentLikeRepo:      commentLikeRepo,
		config:               cfg,
		sensitiveWordService: swService,
		activityService:      activityService,
	}
}

// IsCommentEnabled 检查评论功能是否启用
func (s *commentService) IsCommentEnabled() bool {
	return s.config.Comment.Enabled
}

// CreateCommentWithAudit 创建评论并执行审核逻辑
func (s *commentService) CreateCommentWithAudit(comment *models.Comment, userID uint) (*models.Comment, error) {
	// 检查评论功能是否启用
	if !s.config.Comment.Enabled {
		return nil, errors.New("评论功能未启用")
	}

	// 验证评论长度
	if len(comment.Content) < s.config.Comment.MinLength {
		return nil, errors.New("评论内容太短")
	}
	if len(comment.Content) > s.config.Comment.MaxLength {
		logger.Info("评论太长,支持的长度为", logger.Any("MaxLength", s.config.Comment.MaxLength))
		return nil, errors.New("评论内容太长")
	}

	// 检查多级评论层级
	if s.config.Comment.ReplyEnabled && comment.ParentID != nil {
		level, err := s.getCommentLevel(*comment.ParentID)
		if err != nil {
			return nil, errors.New("无法获取父评论层级")
		}
		if level >= s.config.Comment.MaxLevel {
			return nil, errors.New("回复层级过深")
		}
	}

	// 设置用户ID和IP
	comment.UserID = userID

	// 审核逻辑
	if s.config.Comment.AutoApprove {
		// 自动通过模式：检查敏感词
		if s.config.Comment.SensitiveWordCheck {
			hasSensitive, sensitiveWords := s.sensitiveWordService.CheckSensitiveWords(comment.Content)
			if hasSensitive {
				// 包含敏感词，进入审核队列
				comment.Status = 0 // 待审核
				logger.Info("评论包含敏感词，进入审核队列",
					logger.String("content", comment.Content),
					logger.Any("sensitive_words", sensitiveWords))
			} else {
				// 无敏感词，自动通过
				comment.Status = 1 // 已通过
			}
		} else {
			// 不检测敏感词，直接通过
			comment.Status = 1
		}
	} else {
		// 手动审核模式：所有评论都进入审核队列
		comment.Status = 0 // 待审核
	}

	// 保存原始内容（用于日志记录）
	originalContent := comment.Content

	// 创建评论
	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	// 记录用户活动（评论）- 使用原始内容
	go func() {
		_ = s.activityService.RecordComment(comment.UserID, comment.ArticleID, originalContent)
	}()

	// 过滤返回给前端的评论内容中的敏感词
	if s.config.Comment.SensitiveWordCheck {
		comment.Content = s.sensitiveWordService.FilterSensitiveWords(comment.Content)
	}

	return comment, nil
}

// getCommentLevel 获取评论层级
func (s *commentService) getCommentLevel(commentID uint) (int, error) {
	level := 0
	currentID := commentID

	for i := 0; i < s.config.Comment.MaxDepth; i++ {
		comment, err := s.commentRepo.GetByID(currentID)
		if err != nil {
			return -1, err
		}

		if comment.ParentID == nil {
			return level, nil
		}

		level++
		currentID = *comment.ParentID
	}

	return -1, errors.New("评论层级过深")
}

func (s *commentService) CreateComment(comment *models.Comment) error {
	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now
	return s.commentRepo.Create(comment)
}

// ListPendingComments 获取待审核评论列表（管理员使用）
func (s *commentService) ListPendingComments(page, pageSize int) ([]*models.Comment, int64, error) {
	comments, total, err := s.commentRepo.ListByStatus(0, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 注意：待审核评论不过滤敏感词，让管理员看到原始内容
	// 这是为了让管理员能够完整审查评论内容

	return comments, total, nil
}

func (s *commentService) UpdateComment(id uint, comment *models.Comment) error {
	existing, err := s.commentRepo.GetByID(id)
	if err != nil {
		return errors.New("comment not found")
	}

	comment.ID = id
	comment.CreatedAt = existing.CreatedAt
	comment.UpdatedAt = time.Now()
	return s.commentRepo.Update(comment)
}

func (s *commentService) DeleteComment(id uint) error {
	return s.commentRepo.Delete(id)
}

func (s *commentService) DeleteCommentWithChildren(id uint) error {
	// 获取评论信息
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 删除评论及其子评论
	if err := s.commentRepo.DeleteWithChildren(id); err != nil {
		return err
	}

	// 更新文章评论数（减去删除的评论数）
	// 注意：这里使用原始评论数减1，因为删除评论后无法准确统计子评论数量
	if err := s.commentRepo.DecrementArticleCommentCount(comment.ArticleID); err != nil {
		logger.Error("更新文章评论数失败", logger.Err("err", err))
	}

	return nil
}

// DeleteCommentWithChildrenByUser 删除评论及其所有子评论（需要验证用户权限）
func (s *commentService) DeleteCommentWithChildrenByUser(id, userID uint) error {
	// 获取评论信息
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return errors.New("评论不存在")
	}

	// 验证用户是否为评论作者
	if comment.UserID != userID {
		return errors.New("无权删除该评论")
	}

	// 删除评论及其子评论
	if err := s.commentRepo.DeleteWithChildren(id); err != nil {
		return err
	}

	// 更新文章评论数
	if err := s.commentRepo.DecrementArticleCommentCount(comment.ArticleID); err != nil {
		logger.Error("更新文章评论数失败", logger.Err("err", err))
	}

	return nil
}

func (s *commentService) ListReplies(parentID uint, page, pageSize int) ([]*models.Comment, int64, error) {
	replies, total, err := s.commentRepo.ListReplies(parentID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 过滤回复内容中的敏感词
	if s.config.Comment.SensitiveWordCheck {
		for _, reply := range replies {
			reply.Content = s.sensitiveWordService.FilterSensitiveWords(reply.Content)
		}
	}

	return replies, total, nil
}

func (s *commentService) GetReplyCount(parentID uint) (int64, error) {
	return s.commentRepo.GetReplyCount(parentID)
}

func (s *commentService) GetComment(id uint) (*models.Comment, error) {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 过滤评论内容中的敏感词
	if s.config.Comment.SensitiveWordCheck {
		comment.Content = s.sensitiveWordService.FilterSensitiveWords(comment.Content)
	}

	return comment, nil
}

func (s *commentService) ListComments(articleID uint, page, pageSize int) ([]*models.Comment, int64, error) {
	comments, total, err := s.commentRepo.List(articleID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 过滤评论内容中的敏感词
	if s.config.Comment.SensitiveWordCheck {
		for _, comment := range comments {
			comment.Content = s.sensitiveWordService.FilterSensitiveWords(comment.Content)
		}
	}

	return comments, total, nil
}

// ListCommentsWithPending 包含当前用户的待审核评论
func (s *commentService) ListCommentsWithPending(articleID, userID uint, page, pageSize int) ([]*models.Comment, int64, error) {
	comments, total, err := s.commentRepo.List(articleID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 如果有用户ID，额外查询该用户的待审核评论
	if userID > 0 {
		pendingComments, err := s.commentRepo.ListPendingByUserAndArticle(userID, articleID)
		if err == nil && len(pendingComments) > 0 {
			// 合并已通过和待审核的评论
			allComments := make([]*models.Comment, 0, len(comments)+len(pendingComments))
			allComments = append(allComments, comments...)
			allComments = append(allComments, pendingComments...)
			// 更新总数（这里简单处理，实际可能需要更精确的计算）
			total += int64(len(pendingComments))
			comments = allComments
		}
	}

	// 过滤评论内容中的敏感词
	if s.config.Comment.SensitiveWordCheck {
		for _, comment := range comments {
			comment.Content = s.sensitiveWordService.FilterSensitiveWords(comment.Content)
		}
	}

	return comments, total, nil
}

func (s *commentService) ListUserComments(userID uint, page, pageSize int) ([]*models.Comment, int64, error) {
	comments, total, err := s.commentRepo.ListByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 过滤评论内容中的敏感词
	if s.config.Comment.SensitiveWordCheck {
		for _, comment := range comments {
			comment.Content = s.sensitiveWordService.FilterSensitiveWords(comment.Content)
		}
	}

	return comments, total, nil
}

func (s *commentService) ListAllComments(page, pageSize int) ([]*models.Comment, int64, error) {
	comments, total, err := s.commentRepo.ListAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 过滤评论内容中的敏感词
	if s.config.Comment.SensitiveWordCheck {
		for _, comment := range comments {
			comment.Content = s.sensitiveWordService.FilterSensitiveWords(comment.Content)
		}
	}

	return comments, total, nil
}

func (s *commentService) ApproveComment(id uint) error {
	return s.commentRepo.Approve(id)
}

func (s *commentService) RejectComment(id uint) error {
	return s.commentRepo.Reject(id)
}

// LikeComment 点赞评论
func (s *commentService) LikeComment(commentID, userID uint) error {
	// 检查评论是否存在
	_, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}

	// 检查是否已点赞
	hasLiked, err := s.commentLikeRepo.CheckUserLiked(commentID, userID)
	if err != nil {
		return err
	}
	if hasLiked {
		return errors.New("已经点赞过该评论")
	}

	// 创建点赞记录
	like := &models.CommentLike{
		CommentID: commentID,
		UserID:    userID,
	}
	if err := s.commentLikeRepo.Create(like); err != nil {
		return err
	}

	// 更新评论的点赞数
	return s.RefreshCommentLikeCount(commentID)
}

// UnlikeComment 取消点赞评论
func (s *commentService) UnlikeComment(commentID, userID uint) error {
	// 检查是否已点赞
	hasLiked, err := s.commentLikeRepo.CheckUserLiked(commentID, userID)
	if err != nil {
		return err
	}
	if !hasLiked {
		return errors.New("未点赞该评论")
	}

	// 删除点赞记录
	if err := s.commentLikeRepo.Delete(commentID, userID); err != nil {
		return err
	}

	// 更新评论的点赞数
	return s.RefreshCommentLikeCount(commentID)
}

// CheckUserLiked 检查用户是否已点赞某评论
func (s *commentService) CheckUserLiked(commentID, userID uint) (bool, error) {
	return s.commentLikeRepo.CheckUserLiked(commentID, userID)
}

// RefreshCommentLikeCount 刷新评论的点赞数
func (s *commentService) RefreshCommentLikeCount(commentID uint) error {
	count, err := s.commentLikeRepo.GetLikeCount(commentID)
	if err != nil {
		return err
	}

	// 获取评论并更新点赞数
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return err
	}

	comment.Likes = uint32(count)
	comment.UpdatedAt = time.Now()
	return s.commentRepo.Update(comment)
}
