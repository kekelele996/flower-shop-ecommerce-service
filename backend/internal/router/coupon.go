package router

import "github.com/flowershop/backend/internal/handler"

// CouponRoutes 优惠券模块路由注册说明。
type CouponRoutes struct {
	Handler *handler.CouponHandler
}

func NewCouponRoutes(h *handler.CouponHandler) *CouponRoutes {
	return &CouponRoutes{Handler: h}
}
