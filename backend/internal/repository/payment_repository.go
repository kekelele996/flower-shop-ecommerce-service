package repository

import (
	"time"

	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
)

// PaymentRepository 支付记录仓储。
type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(tx *gorm.DB, p *model.Payment) error {
	if err := tx.Create(p).Error; err != nil {
		return util.WrapAppError(50001, "create payment failed", err)
	}
	return nil
}

func (r *PaymentRepository) FindByOrderID(orderID uint) ([]model.Payment, error) {
	var list []model.Payment
	if err := r.db.Where("order_id = ?", orderID).Order("id desc").Find(&list).Error; err != nil {
		return nil, util.WrapAppError(50001, "find payments failed", err)
	}
	return list, nil
}

func (r *PaymentRepository) MarkSuccess(tx *gorm.DB, payNo, transactionID string) error {
	now := time.Now()
	if err := tx.Model(&model.Payment{}).Where("pay_no = ?", payNo).
		Updates(map[string]interface{}{"status": "SUCCESS", "transaction_id": transactionID, "paid_at": &now}).Error; err != nil {
		return util.WrapAppError(50001, "mark payment success failed", err)
	}
	return nil
}
