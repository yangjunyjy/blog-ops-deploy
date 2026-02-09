package models

import "time"

// Comment 文章评论模型，存储用户对文章的评论和回复
type Comment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`                    // 评论唯一标识ID
	ArticleID uint      `gorm:"not null;index:idx_article_id" json:"article_id"`       // 被评论文章的ID
	UserID    uint      `gorm:"not null;index:idx_user_id" json:"user_id"`             // 发表评论用户的ID
	ParentID  *uint     `gorm:"index:idx_parent_id" json:"parent_id"`                  // 父评论ID，用于回复评论，NULL表示顶级评论
	Content   string    `gorm:"type:text;not null" json:"content" validate:"required"` // 评论内容，必填
	Status    uint8     `gorm:"not null;default:1" json:"status"`                      // 评论审核状态：0-待审核，1-已通过，2-已拒绝
	Likes     uint32    `gorm:"not null;default:0" json:"likes"`                       // 评论点赞数统计
	IPAddress string    `gorm:"size:50" json:"ip_address"`                             // 发表评论的IP地址
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                      // 评论发表时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`                      // 评论最后更新时间

	// 关联关系（不使用外键约束）
	Article *Article  `gorm:"-" json:"article,omitempty"` // 被评论的文章，通过手动查询获取
	User    *User     `gorm:"-" json:"user,omitempty"`    // 发表评论的用户，通过手动查询获取
	Parent  *Comment  `gorm:"-" json:"parent,omitempty"`  // 父评论，用于构建评论回复树
	Replies []Comment `gorm:"-" json:"replies,omitempty"` // 该评论的所有回复子评论列表，通过手动查询获取
}

func (Comment) TableName() string {
	return "comment"
}

// UserCommentArticle 用户评论文章关联模型，用于查询用户参与评论的所有文章
type UserCommentArticle struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`                                             // 记录唯一标识ID
	UserID    uint      `gorm:"not null;index:idx_user_id,index:idx_user_article_comment" json:"user_id"`       // 用户ID
	ArticleID uint      `gorm:"not null;index:idx_article_id,index:idx_user_article_comment" json:"article_id"` // 文章ID
	CommentID uint      `gorm:"not null;index:idx_comment_id,index:idx_user_article_comment" json:"comment_id"` // 评论ID
	IsAuthor  uint8     `gorm:"not null;default:0" json:"is_author"`                                            // 是否为文章作者：0-否，1-是
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                                               // 创建时间
}

func (UserCommentArticle) TableName() string {
	return "user_comment_article"
}
