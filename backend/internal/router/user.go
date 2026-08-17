package router

import "github.com/flowershop/backend/internal/handler"

// UserRoutes 用户模块路由注册说明（实际注册见 router.go）。
type UserRoutes struct {
	Handler *handler.UserHandler
}

func NewUserRoutes(h *handler.UserHandler) *UserRoutes {
	return &UserRoutes{Handler: h}
}
