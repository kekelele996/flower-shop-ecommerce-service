package router

import "github.com/flowershop/backend/internal/handler"

// ProductRoutes 商品模块路由注册说明。
type ProductRoutes struct {
	Handler *handler.ProductHandler
}

func NewProductRoutes(h *handler.ProductHandler) *ProductRoutes {
	return &ProductRoutes{Handler: h}
}
