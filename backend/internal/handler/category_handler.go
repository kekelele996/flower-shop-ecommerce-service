package handler

import (
	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/service"
	"github.com/flowershop/backend/internal/util"
	"github.com/gin-gonic/gin"
)

// CategoryHandler 分类处理器。
type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) List(c *gin.Context) {
	tree, err := h.svc.ListTree()
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, tree)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "category params invalid: "+err.Error()))
		return
	}
	cat, err := h.svc.Create(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, cat)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "category id invalid"))
		return
	}
	var req dto.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "category params invalid: "+err.Error()))
		return
	}
	cat, err := h.svc.Update(id.ID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, cat)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "category id invalid"))
		return
	}
	if err := h.svc.Delete(id.ID); err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": "category deleted"})
}
