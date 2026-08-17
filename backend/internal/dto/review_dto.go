package dto

// ReviewCreateRequest 提交评价请求。
type ReviewCreateRequest struct {
	OrderItemID uint     `json:"order_item_id" binding:"required,gt=0"`
	Rating      int      `json:"rating" binding:"required,min=1,max=5"`
	Content     string   `json:"content" binding:"required,min=2,max=1000"`
	Images      []string `json:"images"`
}

// ReviewReplyRequest 商家回复请求。
type ReviewReplyRequest struct {
	Reply string `json:"reply" binding:"required,min=1,max=500"`
}

// ReviewQuery 评价查询。
type ReviewQuery struct {
	Page      int  `form:"page"`
	PageSize  int  `form:"page_size"`
	ProductID uint `form:"product_id"`
	OrderID   uint `form:"order_id"`
	UserID    uint `form:"user_id"`
}

// ReviewVO 评价视图对象。
type ReviewVO struct {
	ID          uint     `json:"id"`
	UserID      uint     `json:"user_id"`
	Username    string   `json:"username"`
	OrderID     uint     `json:"order_id"`
	OrderItemID uint     `json:"order_item_id"`
	ProductID   uint     `json:"product_id"`
	ProductName string   `json:"product_name"`
	Rating      int      `json:"rating"`
	Content     string   `json:"content"`
	Images      []string `json:"images"`
	Reply       string   `json:"reply"`
	RepliedAt   string   `json:"replied_at"`
	CreatedAt   string   `json:"created_at"`
}
