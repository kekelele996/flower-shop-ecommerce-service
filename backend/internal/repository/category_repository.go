package repository

import (
	"errors"

	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
)

// CategoryRepository 分类仓储。
type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(c *model.Category) error {
	if err := r.db.Create(c).Error; err != nil {
		return util.WrapAppError(50001, "create category failed", err)
	}
	return nil
}

func (r *CategoryRepository) Update(c *model.Category) error {
	if err := r.db.Save(c).Error; err != nil {
		return util.WrapAppError(50001, "update category failed", err)
	}
	return nil
}

func (r *CategoryRepository) Delete(id uint) error {
	res := r.db.Delete(&model.Category{}, id)
	if res.Error != nil {
		return util.WrapAppError(50001, "delete category failed", res.Error)
	}
	if res.RowsAffected == 0 {
		return util.ErrNotFound
	}
	return nil
}

func (r *CategoryRepository) FindByID(id uint) (*model.Category, error) {
	var c model.Category
	err := r.db.First(&c, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "find category by id failed", err)
	}
	return &c, nil
}

// ListAll 返回全部分类，用于前端构建树与筛选（复用：handler 列表/商品查询共用）。
func (r *CategoryRepository) ListAll() ([]model.Category, error) {
	var list []model.Category
	if err := r.db.Order("level asc, sort asc, id asc").Find(&list).Error; err != nil {
		return nil, util.WrapAppError(50001, "list categories failed", err)
	}
	return list, nil
}
