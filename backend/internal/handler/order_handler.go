package handler

import (
	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/middleware"
	"github.com/flowershop/backend/internal/service"
	"github.com/flowershop/backend/internal/util"
	"github.com/gin-gonic/gin"
)

// OrderHandler 订单处理器。
type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	user := middleware.CurrentUser(c)
	var req dto.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "checkout params invalid: "+err.Error()))
		return
	}
	order, err := h.svc.Checkout(user.ID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": constants.MsgOrderCreated, "order": order})
}

func (h *OrderHandler) List(c *gin.Context) {
	user := middleware.CurrentUser(c)
	var q dto.OrderQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "order query invalid: "+err.Error()))
		return
	}
	page := util.ParsePage(c)
	var filterUser uint
	if user.Role != "ADMIN" {
		filterUser = user.ID
	}
	list, total, err := h.svc.List(q, filterUser)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OKPage(c, list, total, page.Page, page.PageSize)
}

func (h *OrderHandler) Detail(c *gin.Context) {
	user := middleware.CurrentUser(c)
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "order id invalid"))
		return
	}
	order, err := h.svc.Detail(user.ID, id.ID, false)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, order)
}

func (h *OrderHandler) Pay(c *gin.Context) {
	user := middleware.CurrentUser(c)
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "order id invalid"))
		return
	}
	var req dto.PayRequest
	_ = c.ShouldBindJSON(&req)
	order, err := h.svc.Pay(user.ID, id.ID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": constants.MsgOrderPaid, "order": order})
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	user := middleware.CurrentUser(c)
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "order id invalid"))
		return
	}
	var req dto.CancelOrderRequest
	_ = c.ShouldBindJSON(&req)
	order, err := h.svc.Cancel(user.ID, id.ID, req.Reason)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": constants.MsgOrderCancelled, "order": order})
}

func (h *OrderHandler) Complete(c *gin.Context) {
	user := middleware.CurrentUser(c)
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "order id invalid"))
		return
	}
	order, err := h.svc.Complete(user.ID, id.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": constants.MsgOrderCompleted, "order": order})
}

func (h *OrderHandler) Logistics(c *gin.Context) {
	user := middleware.CurrentUser(c)
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "order id invalid"))
		return
	}
	admin := user.Role == "ADMIN"
	logis, err := h.svc.Logistics(user.ID, id.ID, admin)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, logis)
}

func (h *OrderHandler) Ship(c *gin.Context) {
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "order id invalid"))
		return
	}
	var req dto.ShipOrderRequest
	_ = c.ShouldBindJSON(&req)
	order, err := h.svc.Ship(id.ID, req.Carrier)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": constants.MsgOrderShipped, "order": order})
}

// CompleteAdmin 管理员代用户确认收货（用于演示状态流转）。
func (h *OrderHandler) CompleteAdmin(c *gin.Context) {
	user := middleware.CurrentUser(c)
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "order id invalid"))
		return
	}
	order, err := h.svc.AdminComplete(user.ID, id.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": constants.MsgOrderCompleted, "order": order})
}
