package mysql

import (
	"gorm.io/gorm"

	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
)

// CategoryRepositoryImpl 分类仓储实现
type CategoryRepositoryImpl struct {
	db *gorm.DB
}

// NewCategoryRepositoryImpl 创建分类仓储实例
func NewCategoryRepositoryImpl(db *gorm.DB) repository.CategoryRepository {
	return &CategoryRepositoryImpl{db: db}
}

// Create 创建分类
func (r *CategoryRepositoryImpl) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// Update 更新分类
func (r *CategoryRepositoryImpl) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类
func (r *CategoryRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

// GetByID 根据ID获取分类
func (r *CategoryRepositoryImpl) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetBySlug 根据slug获取分类
func (r *CategoryRepositoryImpl) GetBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// List 分页获取分类列表
func (r *CategoryRepositoryImpl) List(page, pageSize int) ([]*models.Category, int64, error) {
	var categories []*models.Category
	var total int64

	offset := (page - 1) * pageSize
	r.db.Model(&models.Category{}).Count(&total)

	// 预加载文章数量
	err := r.db.Offset(offset).Limit(pageSize).Find(&categories).Error

	// 为每个分类填充文章数量
	for _, category := range categories {
		var count int64
		r.db.Model(&models.Article{}).Where("category_id = ? AND status = ?", category.ID, 1).Count(&count)
		category.ArticleCount = count
	}

	return categories, total, err
}

// GetTree 获取分类树
func (r *CategoryRepositoryImpl) GetTree() ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.Where("parent_id IS NULL OR parent_id = 0").Find(&categories).Error
	return categories, err
}

// GetChildren 获取子分类
func (r *CategoryRepositoryImpl) GetChildren(parentID uint) ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.Where("parent_id = ?", parentID).Find(&categories).Error
	return categories, err
}

// GetCategoryArticles 获取分类下的文章
func (r *CategoryRepositoryImpl) GetCategoryArticles(categoryID uint, page, pageSize int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	offset := (page - 1) * pageSize

	// 获取该分类及其所有子分类ID
	categoryIDs := []uint{categoryID}
	var children []*models.Category
	r.db.Where("parent_id = ?", categoryID).Find(&children)
	for _, child := range children {
		categoryIDs = append(categoryIDs, child.ID)
	}

	// 获取文章列表和总数
	query := r.db.Model(&models.Article{}).Where("status = ?", 1) // 只获取已发布的文章
	if len(categoryIDs) > 0 {
		query = query.Where("category_id IN ?", categoryIDs)
	}

	query.Count(&total)
	err := query.Preload("Tags").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&articles).Error

	return articles, total, err
}
