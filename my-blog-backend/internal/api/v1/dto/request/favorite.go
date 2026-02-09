package request

type CreateFolderRequest struct {
	Name        string `form:"name" binding:"required,max=32"`
	Description string `form:"description" binding:"max=255"`
	SortOrder   int    `form:"sortOrder" binding:"min=0"`
}

type AddFavoriteRequest struct {
	ArticleID uint `form:"article_id" json:"article_id" binding:"required"`
	FolderID  uint `form:"folder_id" json:"folder_id" binding:"required"`
}
