package dto

// CartAddRequest 添加购物车请求。
type CartAddRequest struct {
	ProductID uint `json:"product_id" binding:"required,gt=0"`
	Quantity  int  `json:"quantity" binding:"required,min=1,max=99"`
}

// CartUpdateRequest 更新购物车请求。
type CartUpdateRequest struct {
	Quantity *int  `json:"quantity" binding:"omitempty,min=1,max=99"`
	Selected *bool `json:"selected"`
}

// CartItemVO 购物车条目视图对象。
type CartItemVO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Selected  bool      `json:"selected"`
	Product   ProductVO `json:"product"`
	Subtotal  float64   `json:"subtotal"`
}

// CartSummaryVO 购物车汇总（复用 cartRepository.FindByUserID）。
type CartSummaryVO struct {
	Items         []CartItemVO `json:"items"`
	TotalQuantity int          `json:"total_quantity"`
	TotalPrice    float64      `json:"total_price"`
}
