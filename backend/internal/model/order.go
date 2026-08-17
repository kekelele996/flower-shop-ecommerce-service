package model

import "time"

// Order 订单实体：核心状态机 PENDING_PAYMENT → PENDING_SHIPMENT → SHIPPED → COMPLETED / CANCELLED。
type Order struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	OrderNo        string     `gorm:"size:32;uniqueIndex;not null" json:"order_no"`
	UserID         uint       `gorm:"index;not null" json:"user_id"`
	Status         string     `gorm:"size:24;default:PENDING_PAYMENT;index" json:"status"`
	TotalAmount    float64    `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	DiscountAmount float64    `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	ShippingFee    float64    `gorm:"type:decimal(10,2);default:0" json:"shipping_fee"`
	PayAmount      float64    `gorm:"type:decimal(10,2);not null" json:"pay_amount"`
	CouponID       uint       `gorm:"default:0" json:"coupon_id"`
	ReceiverName   string     `gorm:"size:64;not null" json:"receiver_name"`
	ReceiverPhone  string     `gorm:"size:32;not null" json:"receiver_phone"`
	ReceiverAddr   string     `gorm:"size:255;not null" json:"receiver_addr"`
	Remark         string     `gorm:"size:255" json:"remark"`
	PaidAt         *time.Time `json:"paid_at"`
	ShippedAt      *time.Time `json:"shipped_at"`
	CompletedAt    *time.Time `json:"completed_at"`
	CancelledAt    *time.Time `json:"cancelled_at"`
	CancelledBy    string     `gorm:"size:32" json:"cancelled_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}
