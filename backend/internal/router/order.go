package router

import "github.com/flowershop/backend/internal/handler"

// OrderRoutes 订单模块路由注册说明。
type OrderRoutes struct {
	Handler *handler.OrderHandler
}

func NewOrderRoutes(h *handler.OrderHandler) *OrderRoutes {
	return &OrderRoutes{Handler: h}
}
