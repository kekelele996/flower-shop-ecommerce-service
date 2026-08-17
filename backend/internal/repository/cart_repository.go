package repository

import (
	"errors"

	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
)

// CartRepository 购物车仓储。
type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) Create(item *model.CartItem) error {
	if err := r.db.Create(item).Error; err != nil {
		return util.WrapAppError(50001, "create cart item failed", err)
	}
	return nil
}

// FindByUserAndProduct 查询用户购物车中指定商品条目。
func (r *CartRepository) FindByUserAndProduct(userID, productID uint) (*model.CartItem, error) {
	var item model.CartItem
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "find cart item failed", err)
	}
	return &item, nil
}

func (r *CartRepository) FindByID(id uint) (*model.CartItem, error) {
	var item model.CartItem
	err := r.db.First(&item, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "find cart item by id failed", err)
	}
	return &item, nil
}

// FindByUserID 用户购物车（复用：列表、结算汇总、下单）。
func (r *CartRepository) FindByUserID(userID uint) ([]model.CartItem, error) {
	var items []model.CartItem
	if err := r.db.Preload("Product").Where("user_id = ?", userID).Order("id desc").Find(&items).Error; err != nil {
		return nil, util.WrapAppError(50001, "list cart items failed", err)
	}
	return items, nil
}

func (r *CartRepository) Update(item *model.CartItem) error {
	if err := r.db.Save(item).Error; err != nil {
		return util.WrapAppError(50001, "update cart item failed", err)
	}
	return nil
}

func (r *CartRepository) UpdateQuantity(id uint, quantity int) error {
	if err := r.db.Model(&model.CartItem{}).Where("id = ?", id).Update("quantity", quantity).Error; err != nil {
		return util.WrapAppError(50001, "update cart quantity failed", err)
	}
	return nil
}

func (r *CartRepository) Delete(id uint) error {
	res := r.db.Delete(&model.CartItem{}, id)
	if res.Error != nil {
		return util.WrapAppError(50001, "delete cart item failed", res.Error)
	}
	if res.RowsAffected == 0 {
		return util.ErrNotFound
	}
	return nil
}

// DeleteByIDs 批量删除（下单后清理已购条目）。
func (r *CartRepository) DeleteByIDs(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	if err := r.db.Delete(&model.CartItem{}, "id IN ?", ids).Error; err != nil {
		return util.WrapAppError(50001, "delete cart items failed", err)
	}
	return nil
}
