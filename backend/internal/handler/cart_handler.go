package handler

import (
	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/middleware"
	"github.com/flowershop/backend/internal/service"
	"github.com/flowershop/backend/internal/util"
	"github.com/gin-gonic/gin"
)

// CartHandler 购物车处理器。
type CartHandler struct {
	svc *service.CartService
}

func NewCartHandler(svc *service.CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

func (h *CartHandler) Add(c *gin.Context) {
	user := middleware.CurrentUser(c)
	var req dto.CartAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "cart params invalid: "+err.Error()))
		return
	}
	item, err := h.svc.Add(user.ID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, item)
}

func (h *CartHandler) List(c *gin.Context) {
	user := middleware.CurrentUser(c)
	summary, err := h.svc.Summary(user.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, summary)
}

func (h *CartHandler) Update(c *gin.Context) {
	user := middleware.CurrentUser(c)
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "cart item id invalid"))
		return
	}
	var req dto.CartUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "cart params invalid: "+err.Error()))
		return
	}
	if err := h.svc.Update(user.ID, id.ID, req); err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": "cart updated"})
}

func (h *CartHandler) Remove(c *gin.Context) {
	user := middleware.CurrentUser(c)
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "cart item id invalid"))
		return
	}
	if err := h.svc.Remove(user.ID, id.ID); err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": "cart item removed"})
}
