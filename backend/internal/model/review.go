package model

import "time"

// Review 评价实体：订单完成后用户评价，商家可回复。
type Review struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	OrderID     uint       `gorm:"index;not null" json:"order_id"`
	OrderItemID uint       `gorm:"index;not null" json:"order_item_id"`
	ProductID   uint       `gorm:"index;not null" json:"product_id"`
	Rating      int        `gorm:"not null" json:"rating"`
	Content     string     `gorm:"type:text" json:"content"`
	Images      string     `gorm:"type:text" json:"images"`
	Reply       string     `gorm:"type:text" json:"reply"`
	RepliedAt   *time.Time `json:"replied_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
