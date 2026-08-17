package constants

// 业务枚举集中定义：角色、商品状态、订单状态、支付状态、物流状态、优惠券类型/状态、用户状态。
const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"

	UserStatusActive   = "ACTIVE"
	UserStatusDisabled = "DISABLED"

	ProductStatusOnSale  = "ON_SALE"
	ProductStatusOffSale = "OFF_SALE"

	OrderStatusPendingPayment  = "PENDING_PAYMENT"
	OrderStatusPendingShipment = "PENDING_SHIPMENT"
	OrderStatusShipped         = "SHIPPED"
	OrderStatusCompleted       = "COMPLETED"
	OrderStatusCancelled       = "CANCELLED"

	PaymentStatusPending = "PENDING"
	PaymentStatusSuccess = "SUCCESS"
	PaymentStatusFailed  = "FAILED"

	PaymentChannelAlipaySandbox = "ALIPAY_SANDBOX"

	LogisticsStatusPending        = "PENDING"
	LogisticsStatusPickedUp       = "PICKED_UP"
	LogisticsStatusInTransit      = "IN_TRANSIT"
	LogisticsStatusOutForDelivery = "OUT_FOR_DELIVERY"
	LogisticsStatusDelivered      = "DELIVERED"

	CouponTypeFullReduction = "FULL_REDUCTION"
	CouponTypeDiscount      = "DISCOUNT"

	CouponStatusUnused  = "UNUSED"
	CouponStatusUsed    = "USED"
	CouponStatusExpired = "EXPIRED"
)

// ValidOrderStatuses 合法的订单状态集合，用于 handler 校验与前端筛选。
var ValidOrderStatuses = map[string]bool{
	OrderStatusPendingPayment:  true,
	OrderStatusPendingShipment: true,
	OrderStatusShipped:         true,
	OrderStatusCompleted:       true,
	OrderStatusCancelled:       true,
}

// ValidProductStatuses 合法的商品状态集合。
var ValidProductStatuses = map[string]bool{
	ProductStatusOnSale:  true,
}
