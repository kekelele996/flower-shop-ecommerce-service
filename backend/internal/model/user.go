package model

import "time"

// User 用户实体：角色字段驱动 JWT + RBAC 权限。
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:128;not null" json:"-"`
	Nickname  string    `gorm:"size:64" json:"nickname"`
	Email     string    `gorm:"size:128" json:"email"`
	Phone     string    `gorm:"size:32" json:"phone"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Role      string    `gorm:"size:16;default:USER;index" json:"role"`
	Status    string    `gorm:"size:16;default:ACTIVE;index" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
