package model

import "time"

// OrderItem 订单明细实体。
type OrderItem struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderID      uint      `gorm:"index;not null" json:"order_id"`
	ProductID    uint      `gorm:"index;not null" json:"product_id"`
	ProductName  string    `gorm:"size:128;not null" json:"product_name"`
	ProductImage string    `gorm:"size:255" json:"product_image"`
	Price        float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Quantity     int       `gorm:"not null" json:"quantity"`
	TotalPrice   float64   `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Reviewed     bool      `gorm:"default:false" json:"reviewed"`
	CreatedAt    time.Time `json:"created_at"`
}
