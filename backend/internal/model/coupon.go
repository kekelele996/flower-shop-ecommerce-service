package model

import "time"

// Coupon 优惠券实体：满减/折扣，用户领取后结算抵扣。
type Coupon struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `gorm:"index;not null" json:"user_id"`
	TemplateID   uint       `gorm:"index;default:0" json:"template_id"`
	Name         string     `gorm:"size:64;not null" json:"name"`
	Type         string     `gorm:"size:24;not null" json:"type"`
	Threshold    float64    `gorm:"type:decimal(10,2);default:0" json:"threshold"`
	Amount       float64    `gorm:"type:decimal(10,2);default:0" json:"amount"`
	DiscountRate float64    `gorm:"type:decimal(4,2);default:0" json:"discount_rate"`
	Status       string     `gorm:"size:16;default:UNUSED;index" json:"status"`
	UsedAt       *time.Time `json:"used_at"`
	OrderID      uint       `gorm:"default:0" json:"order_id"`
	ValidFrom    time.Time  `json:"valid_from"`
	ValidTo      time.Time  `json:"valid_to"`
	CreatedAt    time.Time  `json:"created_at"`
}
