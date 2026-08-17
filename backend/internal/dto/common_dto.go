package dto

// IDRequest 通用 ID 路径参数。
type IDRequest struct {
	ID uint `uri:"id" binding:"required,gt=0"`
}

// PageQuery 通用分页查询。
type PageQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}
