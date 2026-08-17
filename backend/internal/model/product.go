package model

import "time"

// Product 商品实体：Images 为 JSON 数组字符串，Detail 为富文本。
type Product struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CategoryID    uint      `gorm:"index;not null" json:"category_id"`
	Name          string    `gorm:"size:128;not null;index" json:"name"`
	SubTitle      string    `gorm:"size:255" json:"sub_title"`
	Description   string    `gorm:"type:text" json:"description"`
	Detail        string    `gorm:"type:text" json:"detail"`
	CoverImage    string    `gorm:"size:255" json:"cover_image"`
	Images        string    `gorm:"type:text" json:"images"`
	Price         float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice float64   `gorm:"type:decimal(10,2)" json:"original_price"`
	Stock         int       `gorm:"not null;default:0" json:"stock"`
	Sales         int       `gorm:"default:0" json:"sales"`
	Rating        float64   `gorm:"type:decimal(3,2);default:5" json:"rating"`
	RatingCount   int       `gorm:"default:0" json:"rating_count"`
	ShippingFrom  string    `gorm:"size:64" json:"shipping_from"`
	FreeShipping  bool      `gorm:"default:false" json:"free_shipping"`
	Status        string    `gorm:"size:16;default:ON_SALE;index" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
