package router

import "github.com/flowershop/backend/internal/handler"

// AuditRoutes 审计模块路由注册说明。
type AuditRoutes struct {
	Handler *handler.AuditHandler
}

func NewAuditRoutes(h *handler.AuditHandler) *AuditRoutes {
	return &AuditRoutes{Handler: h}
}
