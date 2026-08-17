package handler

import (
	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/service"
	"github.com/flowershop/backend/internal/util"
	"github.com/gin-gonic/gin"
)

// UploadHandler 文件上传处理器。
type UploadHandler struct {
	svc *service.UploadService
}

func NewUploadHandler(svc *service.UploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	var form dto.UploadFile
	if err := c.ShouldBind(&form); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "upload params invalid: "+err.Error()))
		return
	}
	if form.Kind == "" {
		form.Kind = "product"
	}
	url, err := h.svc.Upload(c.Request.Context(), form.File, form.Kind)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": constants.MsgUploadSuccess, "url": url})
}
