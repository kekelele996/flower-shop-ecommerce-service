package model

import "time"

// CartItem 购物车条目实体。
type CartItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	Selected  bool      `gorm:"default:true" json:"selected"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
