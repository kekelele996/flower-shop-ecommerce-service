package router

import "github.com/flowershop/backend/internal/handler"

// CategoryRoutes 分类模块路由注册说明。
type CategoryRoutes struct {
	Handler *handler.CategoryHandler
}

func NewCategoryRoutes(h *handler.CategoryHandler) *CategoryRoutes {
	return &CategoryRoutes{Handler: h}
}
