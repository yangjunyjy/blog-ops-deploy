package mysql

import (
	"gorm.io/gorm"

	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/repository"
)

// CommentRepositoryImpl 评论仓储实现（更新版：不使用外键，使用中间表）
type CommentRepositoryImpl struct {
	db *gorm.DB
}

// NewCommentRepositoryImpl 创建评论仓储实例
func NewCommentRepositoryImpl(db *gorm.DB) repository.CommentRepository {
	return &CommentRepositoryImpl{db: db}
}

// Create 创建评论
func (r *CommentRepositoryImpl) Create(comment *models.Comment) error {
	// 开始事务
	tx := r.db.Begin()

	// 创建评论
	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 获取文章作者ID
	var authorID uint64
	tx.Model(&models.Article{}).Where("id = ?", comment.ArticleID).Select("author_id").Scan(&authorID)

	// 插入到中间表（评论者）
	if err := tx.Exec(`
		INSERT IGNORE INTO user_comment_article (user_id, article_id, comment_id, is_author)
		VALUES (?, ?, ?, 0)
	`, comment.UserID, comment.ArticleID, comment.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果评论者不是作者，插入作者记录
	if uint64(comment.UserID) != authorID {
		if err := tx.Exec(`
			INSERT IGNORE INTO user_comment_article (user_id, article_id, comment_id, is_author)
			VALUES (?, ?, ?, 1)
		`, authorID, comment.ArticleID, comment.ID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// Update 更新评论
func (r *CommentRepositoryImpl) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

// DeleteCommentByArticleID 根据文章ID，删除文章下的所有评论
func (r *CommentRepositoryImpl) DeleteCommentByArticleID(articleId uint) error {
	tx := r.db.Begin()
	// 删除评论的点赞记录
	if err := tx.Exec("DELETE FROM comment_like WHERE comment_id IN (SELECT id FROM comment WHERE article_id = ?)", articleId).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 删除中间表记录
	if err := tx.Exec("DELETE FROM user_comment_article WHERE article_id = ?", articleId).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 删除评论
	if err := tx.Delete(&models.Comment{}, "article_id = ?", articleId).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// Delete 删除评论
func (r *CommentRepositoryImpl) Delete(id uint) error {
	// 开始事务
	tx := r.db.Begin()

	// 删除评论的点赞记录
	if err := tx.Exec("DELETE FROM comment_like WHERE comment_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除中间表记录
	if err := tx.Exec("DELETE FROM user_comment_article WHERE comment_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除评论
	if err := tx.Delete(&models.Comment{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// DeleteWithChildren 删除评论及其所有子评论（级联删除）
func (r *CommentRepositoryImpl) DeleteWithChildren(id uint) error {
	// 开始事务
	tx := r.db.Begin()

	// 获取所有需要删除的评论ID（递归查找所有子评论）
	commentIDs, err := r.getAllCommentIDs(tx, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除所有相关评论的点赞记录
	if err := tx.Exec("DELETE FROM comment_like WHERE comment_id IN ?", commentIDs).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除中间表记录
	if err := tx.Exec("DELETE FROM user_comment_article WHERE comment_id IN ?", commentIDs).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除所有评论
	if err := tx.Delete(&models.Comment{}, "id IN ?", commentIDs).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// getAllCommentIDs 递归获取所有评论ID（包括自身）
func (r *CommentRepositoryImpl) getAllCommentIDs(tx *gorm.DB, parentID uint) ([]uint, error) {
	var allIDs []uint
	allIDs = append(allIDs, parentID)

	// 查找直接子评论
	var childIDs []uint
	if err := tx.Model(&models.Comment{}).Where("parent_id = ?", parentID).Pluck("id", &childIDs).Error; err != nil {
		return nil, err
	}

	// 递归获取子评论的ID
	for _, childID := range childIDs {
		grandChildIDs, err := r.getAllCommentIDs(tx, childID)
		if err != nil {
			return nil, err
		}
		allIDs = append(allIDs, grandChildIDs...)
	}

	return allIDs, nil
}

// ListReplies 获取子评论列表（抖音模式：返回所有层级的回复）
func (r *CommentRepositoryImpl) ListReplies(parentID uint, page, pageSize int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	// 1. 先查询直接子评论（二级评论）
	var directChildren []models.Comment
	if err := r.db.Where("parent_id = ? AND status = 1", parentID).
		Find(&directChildren).Error; err != nil {
		return nil, 0, err
	}

	// 2. 收集所有需要查询的父评论ID（包括顶级评论ID和所有直接子评论ID）
	parentIDs := []uint{parentID}
	for _, c := range directChildren {
		parentIDs = append(parentIDs, c.ID)
	}

	// 3. 查询所有以这些ID为父ID的评论（包括二级和三级评论）
	var allComments []models.Comment
	if err := r.db.Where("parent_id IN ? AND status = 1", parentIDs).
		Order("created_at ASC").
		Find(&allComments).Error; err != nil {
		return nil, 0, err
	}

	// 总数应该是所有层级的评论总数（不包括顶级评论本身）
	total = int64(len(allComments))

	// 4. 分页：由于抖音是扁平化显示，直接对所有回复进行分页
	offset := (page - 1) * pageSize
	pagedCommentsSlice := allComments

	// 如果当前页的数据少于全部数据，则返回当前页
	if offset < len(allComments) {
		end := offset + pageSize
		if end > len(allComments) {
			end = len(allComments)
		}
		pagedCommentsSlice = allComments[offset:end]
	} else {
		pagedCommentsSlice = []models.Comment{}
	}

	// 转换为指针切片
	for i := range pagedCommentsSlice {
		comments = append(comments, &pagedCommentsSlice[i])
	}

	// 5. 收集所有用户ID
	userIDs := make([]uint, 0, len(comments))
	commentIDs := make([]uint, 0, len(comments))

	for _, c := range comments {
		commentIDs = append(commentIDs, c.ID)
		if c.UserID > 0 {
			userIDs = append(userIDs, c.UserID)
		}
		// 收集被回复的评论ID
		if c.ParentID != nil {
			commentIDs = append(commentIDs, *c.ParentID)
		}
	}

	// 6. 查询所有被引用的评论
	if len(commentIDs) > 0 {
		commentIDMap := make(map[uint]bool)
		for _, id := range commentIDs {
			commentIDMap[id] = true
		}

		uniqueCommentIDs := make([]uint, 0, len(commentIDMap))
		for id := range commentIDMap {
			uniqueCommentIDs = append(uniqueCommentIDs, id)
		}

		var referencedComments []models.Comment
		if err := r.db.Where("id IN ?", uniqueCommentIDs).Find(&referencedComments).Error; err != nil {
			logger.Error("批量查询评论失败", logger.Any("err", err))
			return comments, total, nil
		}

		// 构建评论映射
		commentMap := make(map[uint]*models.Comment)
		for i := range referencedComments {
			commentMap[referencedComments[i].ID] = &referencedComments[i]
		}

		// 设置每条评论的父评论
		for _, c := range comments {
			if c.ParentID != nil {
				if parentComment, ok := commentMap[*c.ParentID]; ok {
					c.Parent = parentComment
					if parentComment.UserID > 0 {
						userIDs = append(userIDs, parentComment.UserID)
					}
				}
			}
		}
	}

	// 7. 批量加载用户信息
	if len(userIDs) > 0 {
		userIDMap := make(map[uint]bool)
		for _, id := range userIDs {
			userIDMap[id] = true
		}

		uniqueUserIDs := make([]uint, 0, len(userIDMap))
		for id := range userIDMap {
			uniqueUserIDs = append(uniqueUserIDs, id)
		}

		var users []models.User
		if err := r.db.Where("id IN ?", uniqueUserIDs).Find(&users).Error; err != nil {
			logger.Error("批量查询用户失败", logger.Any("err", err))
			return comments, total, nil
		}

		userMap := make(map[uint]*models.User)
		for i := range users {
			userMap[users[i].ID] = &users[i]
		}

		// 设置评论作者
		for _, c := range comments {
			if c.UserID > 0 {
				c.User = userMap[c.UserID]
			}
			// 设置父评论的用户
			if c.Parent != nil && c.Parent.UserID > 0 {
				c.Parent.User = userMap[c.Parent.UserID]
			}
		}
	}

	return comments, total, nil
}

// GetReplyCount 获取子评论数量（抖音模式：统计所有层级的回复）
func (r *CommentRepositoryImpl) GetReplyCount(parentID uint) (int64, error) {
	var count int64

	// 收集所有可能的父评论ID（递归向下查找）
	allParentIDs := []uint{parentID}

	maxDepth := 10 // 防止无限递归
	for depth := 0; depth < maxDepth; depth++ {
		var childIDs []uint
		if err := r.db.Model(&models.Comment{}).
			Where("parent_id IN ? AND status = 1", allParentIDs).
			Pluck("id", &childIDs).Error; err != nil {
			return 0, err
		}

		if len(childIDs) == 0 {
			break
		}

		// 将这些子评论也作为父评论ID添加进去
		allParentIDs = append(allParentIDs, childIDs...)
	}

	// 统计所有以这些ID为父ID的评论数量
	if err := r.db.Model(&models.Comment{}).
		Where("parent_id IN ? AND status = 1", allParentIDs).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// DecrementArticleCommentCount 减少文章评论数
func (r *CommentRepositoryImpl) DecrementArticleCommentCount(articleID uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", articleID).
		UpdateColumn("comment_count", gorm.Expr("comment_count - 1")).Error
}

// IncrementArticleCommentCount 增加文章评论数
func (r *CommentRepositoryImpl) IncrementArticleCommentCount(articleID uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", articleID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error
}

// GetByID 根据ID获取评论（包含关联的用户和文章）
func (r *CommentRepositoryImpl) GetByID(id uint) (*models.Comment, error) {
	return r.GetByIDWithUser(id)
}

// GetByIDWithUser 根据ID获取评论及完整用户信息
func (r *CommentRepositoryImpl) GetByIDWithUser(id uint) (*models.Comment, error) {
	var comment models.Comment
	if err := r.db.Where("id = ?", id).First(&comment).Error; err != nil {
		return nil, err
	}

	// 手动加载关联数据
	if comment.UserID > 0 {
		var user models.User
		if err := r.db.Where("id = ?", comment.UserID).First(&user).Error; err == nil {
			comment.User = &user
		}
	}

	if comment.ArticleID > 0 {
		var article models.Article
		if err := r.db.Where("id = ?", comment.ArticleID).First(&article).Error; err == nil {
			comment.Article = &article
		}
	}

	if comment.ParentID != nil {
		var parent models.Comment
		if err := r.db.Where("id = ?", *comment.ParentID).First(&parent).Error; err == nil {
			comment.Parent = &parent
			// 加载父评论的用户信息
			if parent.UserID > 0 {
				var parentUser models.User
				if err := r.db.Where("id = ?", parent.UserID).First(&parentUser).Error; err == nil {
					parent.User = &parentUser
				}
			}
		}
	}

	return &comment, nil
}

// GetByArticleID 根据文章ID获取评论列表（包含用户信息）
func (r *CommentRepositoryImpl) GetByArticleID(articleID uint, page, pageSize int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	offset := (page - 1) * pageSize
	r.db.Model(&models.Comment{}).Where("article_id = ? AND status = 1", articleID).Count(&total)

	// 获取评论列表
	if err := r.db.Where("article_id = ? AND status = 1", articleID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// 批量加载用户信息
	userIDs := make([]uint, 0, len(comments))
	for _, c := range comments {
		if c.UserID > 0 {
			userIDs = append(userIDs, c.UserID)
		}
	}

	if len(userIDs) > 0 {
		var users []models.User
		r.db.Where("id IN ?", userIDs).Find(&users)

		userMap := make(map[uint]*models.User)
		for i := range users {
			userMap[users[i].ID] = &users[i]
		}

		for _, c := range comments {
			if c.UserID > 0 {
				c.User = userMap[c.UserID]
			}
		}
	}

	return comments, total, nil
}

// List 分页获取评论列表（只获取顶级评论，parent_id IS NULL）
func (r *CommentRepositoryImpl) List(articleID uint, page, pageSize int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Comment{}).Where("status = 1").Where("parent_id IS NULL")
	if articleID > 0 {
		query = query.Where("article_id = ?", articleID)
	}
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// 批量加载用户信息
	userIDs := make([]uint, 0, len(comments))
	for _, c := range comments {
		if c.UserID > 0 {
			userIDs = append(userIDs, c.UserID)
		}
	}

	if len(userIDs) > 0 {
		var users []models.User
		if err := r.db.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
			logger.Error("批量查询用户失败", logger.Any("err", err), logger.Any("userIDs", userIDs))
			return comments, total, nil // 不影响评论列表返回，只是没有用户信息
		}

		logger.Info("查询到用户", logger.Any("user_count", len(users)), logger.Any("userIDs", userIDs))

		userMap := make(map[uint]*models.User)
		for i := range users {
			userMap[users[i].ID] = &users[i]
		}

		for _, c := range comments {
			if c.UserID > 0 {
				c.User = userMap[c.UserID]
			}
		}
	}

	return comments, total, nil
}

// ListByUserID 根据用户ID获取评论列表（使用中间表）
func (r *CommentRepositoryImpl) ListByUserID(userID uint, page, pageSize int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	offset := (page - 1) * pageSize

	// 使用中间表查询用户参与的评论
	query := r.db.Table("comment c").
		Joins("INNER JOIN user_comment_article uca ON c.id = uca.comment_id").
		Where("uca.user_id = ? AND c.status = 1", userID)

	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Order("c.created_at DESC").Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// ListUserCommentedArticles 获取用户评论过的所有文章（使用中间表）
func (r *CommentRepositoryImpl) ListUserCommentedArticles(userID uint, page, pageSize int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	offset := (page - 1) * pageSize

	// 使用中间表查询用户评论过的文章
	query := r.db.Table("articles a").
		Joins("INNER JOIN user_comment_article uca ON a.id = uca.article_id").
		Where("uca.user_id = ? AND a.status = 1", userID).
		Group("a.id")

	query.Count(&total)

	// 选择文章字段，排除 content
	if err := query.Select("a.id, a.title, a.slug, a.summary, a.cover, a.category_id, a.author_id, a.views, a.likes, a.favorites, a.comment_count, a.status, a.is_top, a.sort_order, a.published_at, a.created_at, a.updated_at").
		Offset(offset).
		Limit(pageSize).
		Order("uca.created_at DESC").
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// Approve 审核通过评论
func (r *CommentRepositoryImpl) Approve(id uint) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", id).Update("status", 1).Error
}

// Reject 审核拒绝评论
func (r *CommentRepositoryImpl) Reject(id uint) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", id).Update("status", 2).Error
}

// ListByStatus 根据状态获取评论列表（管理员使用，包含用户和文章信息）
func (r *CommentRepositoryImpl) ListByStatus(status uint8, page, pageSize int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Comment{}).Where("status = ?", status)
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// 批量加载用户信息
	userIDs := make([]uint, 0, len(comments))
	articleIDs := make([]uint, 0, len(comments))
	for _, c := range comments {
		if c.UserID > 0 {
			userIDs = append(userIDs, c.UserID)
		}
		if c.ArticleID > 0 {
			articleIDs = append(articleIDs, c.ArticleID)
		}
	}

	if len(userIDs) > 0 {
		var users []models.User
		r.db.Where("id IN ?", userIDs).Find(&users)

		userMap := make(map[uint]*models.User)
		for i := range users {
			userMap[users[i].ID] = &users[i]
		}

		for _, c := range comments {
			if c.UserID > 0 {
				c.User = userMap[c.UserID]
			}
		}
	}

	if len(articleIDs) > 0 {
		var articles []models.Article
		r.db.Where("id IN ?", articleIDs).Find(&articles)

		articleMap := make(map[uint64]*models.Article)
		for i := range articles {
			articleMap[articles[i].ID] = &articles[i]
		}

		for _, c := range comments {
			if c.ArticleID > 0 {
				c.Article = articleMap[uint64(c.ArticleID)]
			}
		}
	}

	return comments, total, nil
}

// ListAll 获取所有评论列表（管理员使用）
func (r *CommentRepositoryImpl) ListAll(page, pageSize int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Comment{})
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// 批量加载用户信息
	userIDs := make([]uint, 0, len(comments))
	articleIDs := make([]uint, 0, len(comments))
	for _, c := range comments {
		if c.UserID > 0 {
			userIDs = append(userIDs, c.UserID)
		}
		if c.ArticleID > 0 {
			articleIDs = append(articleIDs, c.ArticleID)
		}
	}

	if len(userIDs) > 0 {
		var users []models.User
		r.db.Where("id IN ?", userIDs).Find(&users)

		userMap := make(map[uint]*models.User)
		for i := range users {
			userMap[users[i].ID] = &users[i]
		}

		for _, c := range comments {
			if c.UserID > 0 {
				c.User = userMap[c.UserID]
			}
		}
	}

	if len(articleIDs) > 0 {
		var articles []models.Article
		r.db.Where("id IN ?", articleIDs).Find(&articles)

		articleMap := make(map[uint64]*models.Article)
		for i := range articles {
			articleMap[articles[i].ID] = &articles[i]
		}

		for _, c := range comments {
			if c.ArticleID > 0 {
				c.Article = articleMap[uint64(c.ArticleID)]
			}
		}
	}

	return comments, total, nil
}

// ListPendingByUserAndArticle 获取某用户在某文章下的待审核评论
func (r *CommentRepositoryImpl) ListPendingByUserAndArticle(userID, articleID uint) ([]*models.Comment, error) {
	var comments []*models.Comment

	query := r.db.Model(&models.Comment{}).
		Where("status = 0").
		Where("user_id = ?", userID).
		Where("parent_id IS NULL")

	if articleID > 0 {
		query = query.Where("article_id = ?", articleID)
	}

	if err := query.Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}

	// 批量加载用户信息
	userIDs := make([]uint, 0, len(comments))
	for _, c := range comments {
		if c.UserID > 0 {
			userIDs = append(userIDs, c.UserID)
		}
	}

	if len(userIDs) > 0 {
		var users []models.User
		if err := r.db.Where("id IN ?", userIDs).Find(&users).Error; err == nil {
			userMap := make(map[uint]*models.User)
			for i := range users {
				userMap[users[i].ID] = &users[i]
			}

			for _, c := range comments {
				if c.UserID > 0 {
					c.User = userMap[c.UserID]
				}
			}
		}
	}

	return comments, nil
}
