package model

import "time"

// LogisticsEvent 物流轨迹节点。
type LogisticsEvent struct {
	Time   time.Time `json:"time"`
	Status string    `json:"status"`
	Desc   string    `json:"desc"`
}

// Logistics 物流实体：Events 为 JSON 数组（轨迹节点）。
type Logistics struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    uint      `gorm:"uniqueIndex;not null" json:"order_id"`
	OrderNo    string    `gorm:"size:32;index;not null" json:"order_no"`
	TrackingNo string    `gorm:"size:64;uniqueIndex;not null" json:"tracking_no"`
	Carrier    string    `gorm:"size:32;default:SF_EXPRESS" json:"carrier"`
	Status     string    `gorm:"size:24;default:PENDING" json:"status"`
	Events     string    `gorm:"type:text" json:"events"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
