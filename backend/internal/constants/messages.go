package constants

// 消息文案集中管理：接口返回文案、日志文案、错误提示文案。
const (
	MsgOK                  = "ok"
	MsgInvalidParams       = "invalid parameters"
	MsgUnauthorized        = "unauthorized, please login first"
	MsgForbidden           = "forbidden, no permission"
	MsgNotFound            = "resource not found"
	MsgInternalError       = "internal server error"
	MsgLoginSuccess        = "login success"
	MsgRegisterSuccess     = "register success"
	MsgOrderCreated        = "order created, please pay"
	MsgOrderPaid           = "order paid success"
	MsgOrderCancelled      = "order cancelled"
	MsgOrderShipped        = "order shipped"
	MsgOrderCompleted      = "order confirmed received"
	MsgReviewSuccess       = "review submitted"
	MsgReviewReplySuccess  = "review replied"
	MsgCouponClaimed       = "coupon claimed"
	MsgUploadSuccess       = "upload success"
	MsgInsufficientStock   = "insufficient stock for product=%s, available=%d"
	MsgOrderStatusError    = "order status=%s not allowed for current operation, order_no=%s"
	MsgDuplicateUsername   = "username=%s already exists"
	MsgInvalidCredential   = "invalid username or password for user=%s"
	MsgCouponNotApplicable = "coupon=%s not applicable for order amount=%.2f"
	MsgReviewNotAllow      = "user=%s cannot review order item=%d, order status=%s"
)
