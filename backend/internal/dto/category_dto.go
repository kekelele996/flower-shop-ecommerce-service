package dto

// CategoryCreateRequest 创建分类请求。
type CategoryCreateRequest struct {
	ParentID uint   `json:"parent_id" binding:"omitempty,min=0"`
	Name     string `json:"name" binding:"required,max=64"`
	Sort     int    `json:"sort" binding:"omitempty,min=0"`
}

// CategoryUpdateRequest 更新分类请求。
type CategoryUpdateRequest struct {
	ParentID *uint  `json:"parent_id" binding:"omitempty,min=0"`
	Name     string `json:"name" binding:"omitempty,max=64"`
	Sort     *int   `json:"sort" binding:"omitempty,min=0"`
}

// CategoryVO 分类视图对象。
type CategoryVO struct {
	ID       uint         `json:"id"`
	ParentID uint         `json:"parent_id"`
	Name     string       `json:"name"`
	Level    int          `json:"level"`
	Sort     int          `json:"sort"`
	Children []CategoryVO `json:"children,omitempty"`
}
