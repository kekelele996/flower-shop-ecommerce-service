package router

import "github.com/flowershop/backend/internal/handler"

// CartRoutes 购物车模块路由注册说明。
type CartRoutes struct {
	Handler *handler.CartHandler
}

func NewCartRoutes(h *handler.CartHandler) *CartRoutes {
	return &CartRoutes{Handler: h}
}
