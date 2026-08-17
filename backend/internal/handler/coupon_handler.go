package handler

import (
	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/middleware"
	"github.com/flowershop/backend/internal/service"
	"github.com/flowershop/backend/internal/util"
	"github.com/gin-gonic/gin"
)

// CouponHandler 优惠券处理器。
type CouponHandler struct {
	svc *service.CouponService
}

func NewCouponHandler(svc *service.CouponService) *CouponHandler {
	return &CouponHandler{svc: svc}
}

func (h *CouponHandler) CreateTemplate(c *gin.Context) {
	var req dto.CouponTemplateCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "coupon template params invalid: "+err.Error()))
		return
	}
	created, err := h.svc.CreateTemplate(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": "coupon template created", "created": created})
}

func (h *CouponHandler) Claim(c *gin.Context) {
	user := middleware.CurrentUser(c)
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "coupon id invalid"))
		return
	}
	coupon, err := h.svc.Claim(user.ID, id.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": constants.MsgCouponClaimed, "coupon": coupon})
}

func (h *CouponHandler) Mine(c *gin.Context) {
	user := middleware.CurrentUser(c)
	list, err := h.svc.ListByUser(user.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, list)
}

func (h *CouponHandler) Available(c *gin.Context) {
	user := middleware.CurrentUser(c)
	list, err := h.svc.ListAvailable(user.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, list)
}
