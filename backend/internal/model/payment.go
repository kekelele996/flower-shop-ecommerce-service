package model

import "time"

// Payment 支付记录实体：模拟支付宝沙箱支付。
type Payment struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	OrderID       uint       `gorm:"index;not null" json:"order_id"`
	OrderNo       string     `gorm:"size:32;index;not null" json:"order_no"`
	PayNo         string     `gorm:"size:48;uniqueIndex;not null" json:"pay_no"`
	UserID        uint       `gorm:"index;not null" json:"user_id"`
	Channel       string     `gorm:"size:32;default:ALIPAY_SANDBOX" json:"channel"`
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status        string     `gorm:"size:16;default:PENDING" json:"status"`
	TransactionID string     `gorm:"size:64" json:"transaction_id"`
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAt     time.Time  `json:"created_at"`
}
