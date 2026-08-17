package repository

import (
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
)

// OrderItemRepository 订单明细仓储。
type OrderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (r *OrderItemRepository) CreateBatch(tx *gorm.DB, items []model.OrderItem) error {
	if len(items) == 0 {
		return nil
	}
	if err := tx.Create(&items).Error; err != nil {
		return util.WrapAppError(50001, "create order items failed", err)
	}
	return nil
}

func (r *OrderItemRepository) FindByID(id uint) (*model.OrderItem, error) {
	var item model.OrderItem
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, util.WrapAppError(50001, "find order item failed", err)
	}
	return &item, nil
}

func (r *OrderItemRepository) MarkReviewed(tx *gorm.DB, id uint) error {
	if err := tx.Model(&model.OrderItem{}).Where("id = ?", id).Update("reviewed", true).Error; err != nil {
		return util.WrapAppError(50001, "mark order item reviewed failed", err)
	}
	return nil
}
