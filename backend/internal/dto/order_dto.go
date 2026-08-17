package dto

// CheckoutRequest 下单请求。
type CheckoutRequest struct {
	CartItemIDs   []uint `json:"cart_item_ids" binding:"required,min=1,dive,gt=0"`
	CouponID      uint   `json:"coupon_id" binding:"omitempty,gt=0"`
	ReceiverName  string `json:"receiver_name" binding:"required,max=64"`
	ReceiverPhone string `json:"receiver_phone" binding:"required,max=32"`
	ReceiverAddr  string `json:"receiver_addr" binding:"required,max=255"`
	Remark        string `json:"remark" binding:"omitempty,max=255"`
}

// OrderQuery 订单查询。
type OrderQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Status   string `form:"status"`
	OrderNo  string `form:"order_no"`
	Keyword  string `form:"keyword"`
}

// CancelOrderRequest 取消订单请求。
type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"omitempty,max=255"`
}

// ShipOrderRequest 发货请求（管理员）。
type ShipOrderRequest struct {
	Carrier string `json:"carrier" binding:"omitempty,max=32"`
}

// OrderItemVO 订单明细视图对象。
type OrderItemVO struct {
	ID           uint    `json:"id"`
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
	TotalPrice   float64 `json:"total_price"`
	Reviewed     bool    `json:"reviewed"`
}

// OrderVO 订单视图对象。
type OrderVO struct {
	ID             uint          `json:"id"`
	OrderNo        string        `json:"order_no"`
	UserID         uint          `json:"user_id"`
	Status         string        `json:"status"`
	StatusText     string        `json:"status_text"`
	TotalAmount    float64       `json:"total_amount"`
	DiscountAmount float64       `json:"discount_amount"`
	ShippingFee    float64       `json:"shipping_fee"`
	PayAmount      float64       `json:"pay_amount"`
	CouponID       uint          `json:"coupon_id"`
	ReceiverName   string        `json:"receiver_name"`
	ReceiverPhone  string        `json:"receiver_phone"`
	ReceiverAddr   string        `json:"receiver_addr"`
	Remark         string        `json:"remark"`
	PaidAt         string        `json:"paid_at"`
	ShippedAt      string        `json:"shipped_at"`
	CompletedAt    string        `json:"completed_at"`
	CancelledAt    string        `json:"cancelled_at"`
	CreatedAt      string        `json:"created_at"`
	Items          []OrderItemVO `json:"items"`
}

// LogisticsEventVO 物流轨迹节点。
type LogisticsEventVO struct {
	Time   string `json:"time"`
	Status string `json:"status"`
	Desc   string `json:"desc"`
}

// LogisticsVO 物流视图对象。
type LogisticsVO struct {
	TrackingNo string             `json:"tracking_no"`
	Carrier    string             `json:"carrier"`
	Status     string             `json:"status"`
	StatusText string             `json:"status_text"`
	Events     []LogisticsEventVO `json:"events"`
}

// PayRequest 模拟支付请求。
type PayRequest struct {
	Channel string `json:"channel" binding:"omitempty,oneof=ALIPAY_SANDBOX"`
}
