package repository

import (
	"errors"

	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
)

// LogisticsRepository 物流仓储。
type LogisticsRepository struct {
	db *gorm.DB
}

func NewLogisticsRepository(db *gorm.DB) *LogisticsRepository {
	return &LogisticsRepository{db: db}
}

func (r *LogisticsRepository) Create(tx *gorm.DB, l *model.Logistics) error {
	if err := tx.Create(l).Error; err != nil {
		return util.WrapAppError(50001, "create logistics failed", err)
	}
	return nil
}

func (r *LogisticsRepository) Update(l *model.Logistics) error {
	if err := r.db.Save(l).Error; err != nil {
		return util.WrapAppError(50001, "update logistics failed", err)
	}
	return nil
}

func (r *LogisticsRepository) FindByOrderID(orderID uint) (*model.Logistics, error) {
	var l model.Logistics
	err := r.db.Where("order_id = ?", orderID).First(&l).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "find logistics failed", err)
	}
	return &l, nil
}
