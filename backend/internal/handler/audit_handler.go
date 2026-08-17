package handler

import (
	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/service"
	"github.com/flowershop/backend/internal/util"
	"github.com/gin-gonic/gin"
)

// AuditHandler 审计日志处理器。
type AuditHandler struct {
	svc *service.AuditQueryService
}

func NewAuditHandler(svc *service.AuditQueryService) *AuditHandler {
	return &AuditHandler{svc: svc}
}

func (h *AuditHandler) List(c *gin.Context) {
	var q dto.AuditLogQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "audit query invalid: "+err.Error()))
		return
	}
	page := util.ParsePage(c)
	q.Page = page.Page
	q.PageSize = page.PageSize
	list, total, err := h.svc.List(q)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OKPage(c, list, total, page.Page, page.PageSize)
}
