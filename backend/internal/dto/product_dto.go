package dto

import "mime/multipart"

// ProductCreateRequest 创建商品请求（管理员）。
type ProductCreateRequest struct {
	CategoryID    uint     `json:"category_id" binding:"required,gt=0"`
	Name          string   `json:"name" binding:"required,max=128"`
	SubTitle      string   `json:"sub_title" binding:"omitempty,max=255"`
	Description   string   `json:"description" binding:"omitempty"`
	Detail        string   `json:"detail" binding:"omitempty"`
	CoverImage    string   `json:"cover_image" binding:"omitempty,max=255"`
	Images        []string `json:"images"`
	Price         float64  `json:"price" binding:"required,gt=0"`
	OriginalPrice float64  `json:"original_price" binding:"omitempty,gte=0"`
	Stock         int      `json:"stock" binding:"required,gte=0"`
	ShippingFrom  string   `json:"shipping_from" binding:"omitempty,max=64"`
	FreeShipping  bool     `json:"free_shipping"`
	Status        string   `json:"status" binding:"omitempty,oneof=ON_SALE OFF_SALE"`
}

// ProductUpdateRequest 更新商品请求（管理员）。
type ProductUpdateRequest struct {
	CategoryID    *uint    `json:"category_id" binding:"omitempty,gt=0"`
	Name          *string  `json:"name" binding:"omitempty,max=128"`
	SubTitle      *string  `json:"sub_title" binding:"omitempty,max=255"`
	Description   *string  `json:"description" binding:"omitempty"`
	Detail        *string  `json:"detail" binding:"omitempty"`
	CoverImage    *string  `json:"cover_image" binding:"omitempty,max=255"`
	Images        []string `json:"images"`
	Price         *float64 `json:"price" binding:"omitempty,gt=0"`
	OriginalPrice *float64 `json:"original_price" binding:"omitempty,gte=0"`
	Stock         *int     `json:"stock" binding:"omitempty,gte=0"`
	ShippingFrom  *string  `json:"shipping_from" binding:"omitempty,max=64"`
	FreeShipping  *bool    `json:"free_shipping"`
	Status        *string  `json:"status" binding:"omitempty,oneof=ON_SALE OFF_SALE"`
}

// ProductQuery 商品查询条件。
type ProductQuery struct {
	Page         int     `form:"page"`
	PageSize     int     `form:"page_size"`
	Keyword      string  `form:"keyword"`
	CategoryID   uint    `form:"category_id"`
	MinPrice     float64 `form:"min_price"`
	MaxPrice     float64 `form:"max_price"`
	FreeShipping *bool   `form:"free_shipping"`
	ShippingFrom string  `form:"shipping_from"`
	Sort         string  `form:"sort"`
	Status       string  `form:"status"`
}

// ProductStatusRequest 上下架请求。
type ProductStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=ON_SALE OFF_SALE"`
}

// ProductVO 商品视图对象。
type ProductVO struct {
	ID            uint     `json:"id"`
	CategoryID    uint     `json:"category_id"`
	CategoryName  string   `json:"category_name"`
	Name          string   `json:"name"`
	SubTitle      string   `json:"sub_title"`
	Description   string   `json:"description"`
	Detail        string   `json:"detail"`
	CoverImage    string   `json:"cover_image"`
	Images        []string `json:"images"`
	Price         float64  `json:"price"`
	OriginalPrice float64  `json:"original_price"`
	Stock         int      `json:"stock"`
	Sales         int      `json:"sales"`
	Rating        float64  `json:"rating"`
	RatingCount   int      `json:"rating_count"`
	ShippingFrom  string   `json:"shipping_from"`
	FreeShipping  bool     `json:"free_shipping"`
	Status        string   `json:"status"`
	CreatedAt     string   `json:"created_at"`
}

// UploadFile 上传文件表单。
type UploadFile struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	Kind string                `form:"kind" binding:"omitempty,oneof=product review avatar"`
}
