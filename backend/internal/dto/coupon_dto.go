package dto

// CouponTemplateCreateRequest 创建优惠券模板请求（管理员）。
type CouponTemplateCreateRequest struct {
	Name         string  `json:"name" binding:"required,max=64"`
	Type         string  `json:"type" binding:"required,oneof=FULL_REDUCTION DISCOUNT"`
	Threshold    float64 `json:"threshold" binding:"required,gte=0"`
	Amount       float64 `json:"amount" binding:"omitempty,gte=0"`
	DiscountRate float64 `json:"discount_rate" binding:"omitempty,gt=0,lte=1"`
	Total        int     `json:"total" binding:"required,min=1"`
	ValidDays    int     `json:"valid_days" binding:"required,min=1"`
}

// CouponVO 优惠券视图对象。
type CouponVO struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	TypeText     string  `json:"type_text"`
	Threshold    float64 `json:"threshold"`
	Amount       float64 `json:"amount"`
	DiscountRate float64 `json:"discount_rate"`
	Status       string  `json:"status"`
	ValidFrom    string  `json:"valid_from"`
	ValidTo      string  `json:"valid_to"`
}
