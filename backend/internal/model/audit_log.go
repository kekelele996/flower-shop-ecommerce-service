package model

import "time"

// AuditLog 操作审计日志实体。
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"size:64" json:"username"`
	Action    string    `gorm:"size:24;index" json:"action"`
	Method    string    `gorm:"size:16" json:"method"`
	Path      string    `gorm:"size:255;index" json:"path"`
	Entity    string    `gorm:"size:64" json:"entity"`
	EntityID  string    `gorm:"size:64" json:"entity_id"`
	Detail    string    `gorm:"type:text" json:"detail"`
	RequestID string    `gorm:"size:64;index" json:"request_id"`
	IP        string    `gorm:"size:64" json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}
