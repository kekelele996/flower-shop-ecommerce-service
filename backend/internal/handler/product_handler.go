package handler

import (
	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/middleware"
	"github.com/flowershop/backend/internal/service"
	"github.com/flowershop/backend/internal/util"
	"github.com/gin-gonic/gin"
)

// ProductHandler 商品处理器。
type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) List(c *gin.Context) {
	var q dto.ProductQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "product query invalid: "+err.Error()))
		return
	}
	page := util.ParsePage(c)
	list, total, err := h.svc.List(q, false)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OKPage(c, list, total, page.Page, page.PageSize)
}

func (h *ProductHandler) Detail(c *gin.Context) {
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "product id invalid"))
		return
	}
	p, err := h.svc.Detail(id.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, p)
}

func (h *ProductHandler) View(c *gin.Context) {
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "product id invalid"))
		return
	}
	userID := uint(0)
	if u := middleware.CurrentUser(c); u != nil {
		userID = u.ID
	}
	h.svc.RecordView(c.Request.Context(), userID, id.ID)
	util.OK(c, gin.H{"message": "view recorded"})
}

func (h *ProductHandler) Recommendations(c *gin.Context) {
	userID := uint(0)
	if u := middleware.CurrentUser(c); u != nil {
		userID = u.ID
	}
	list, err := h.svc.Recommendations(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, list)
}

func (h *ProductHandler) AdminList(c *gin.Context) {
	var q dto.ProductQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "product query invalid: "+err.Error()))
		return
	}
	page := util.ParsePage(c)
	list, total, err := h.svc.List(q, true)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OKPage(c, list, total, page.Page, page.PageSize)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "product params invalid: "+err.Error()))
		return
	}
	p, err := h.svc.Create(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, p)
}

func (h *ProductHandler) Update(c *gin.Context) {
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "product id invalid"))
		return
	}
	var req dto.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "product params invalid: "+err.Error()))
		return
	}
	p, err := h.svc.Update(id.ID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, p)
}

func (h *ProductHandler) ChangeStatus(c *gin.Context) {
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "product id invalid"))
		return
	}
	var req dto.ProductStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "status invalid: "+err.Error()))
		return
	}
	p, err := h.svc.ChangeStatus(id.ID, req.Status)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, p)
}
