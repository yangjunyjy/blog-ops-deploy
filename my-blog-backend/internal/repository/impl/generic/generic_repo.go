package generic

import (
	"gorm.io/gorm"
)

// GenericRepository 通用仓储接口
type GenericRepository[T any] interface {
	Create(entity *T) error
	Update(entity *T) error
	Delete(id uint) error
	GetByID(id uint) (*T, error)
	List(page, pageSize int) ([]*T, int64, error)
}

// BaseRepository 基础仓储实现
type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *BaseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *BaseRepository[T]) Delete(id uint) error {
	return r.db.Delete(new(T), id).Error
}

func (r *BaseRepository[T]) GetByID(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) List(page, pageSize int) ([]*T, int64, error) {
	var entities []*T
	var total int64

	offset := (page - 1) * pageSize
	r.db.Model(new(T)).Count(&total)
	err := r.db.Offset(offset).Limit(pageSize).Find(&entities).Error

	return entities, total, err
}
