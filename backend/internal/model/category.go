package model

import "time"

// Category 商品分类实体：支持多级分类（level 1/2）。
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ParentID  uint      `gorm:"index;default:0" json:"parent_id"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Level     int       `gorm:"default:1" json:"level"`
	Sort      int       `gorm:"default:0" json:"sort"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
