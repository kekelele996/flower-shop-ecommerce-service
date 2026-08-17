package router

import "github.com/flowershop/backend/internal/handler"

// ReviewRoutes 评价模块路由注册说明。
type ReviewRoutes struct {
	Handler *handler.ReviewHandler
}

func NewReviewRoutes(h *handler.ReviewHandler) *ReviewRoutes {
	return &ReviewRoutes{Handler: h}
}
