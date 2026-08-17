package handler

import (
	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/middleware"
	"github.com/flowershop/backend/internal/service"
	"github.com/flowershop/backend/internal/util"
	"github.com/gin-gonic/gin"
)

// ReviewHandler 评价处理器。
type ReviewHandler struct {
	svc *service.ReviewService
}

func NewReviewHandler(svc *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{svc: svc}
}

func (h *ReviewHandler) Create(c *gin.Context) {
	user := middleware.CurrentUser(c)
	var req dto.ReviewCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "review params invalid: "+err.Error()))
		return
	}
	review, err := h.svc.Create(user.ID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": constants.MsgReviewSuccess, "review": review})
}

func (h *ReviewHandler) List(c *gin.Context) {
	var q dto.ReviewQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "review query invalid: "+err.Error()))
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

func (h *ReviewHandler) Reply(c *gin.Context) {
	var id dto.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "review id invalid"))
		return
	}
	var req dto.ReviewReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "reply params invalid: "+err.Error()))
		return
	}
	review, err := h.svc.Reply(id.ID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": constants.MsgReviewReplySuccess, "review": review})
}
