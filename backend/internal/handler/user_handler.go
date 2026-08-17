package handler

import (
	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/middleware"
	"github.com/flowershop/backend/internal/service"
	"github.com/flowershop/backend/internal/util"
	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器。
type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "register params invalid: "+err.Error()))
		return
	}
	user, err := h.svc.Register(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": constants.MsgRegisterSuccess, "user": user})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "login params invalid: "+err.Error()))
		return
	}
	token, err := h.svc.Login(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, token)
}

func (h *UserHandler) Profile(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		_ = c.Error(util.NewAppError(constants.CodeUnauthorized, constants.MsgUnauthorized))
		return
	}
	profile, err := h.svc.GetProfile(user.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, profile)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		_ = c.Error(util.NewAppError(constants.CodeUnauthorized, constants.MsgUnauthorized))
		return
	}
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "profile params invalid: "+err.Error()))
		return
	}
	profile, err := h.svc.UpdateProfile(user.ID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, profile)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		_ = c.Error(util.NewAppError(constants.CodeUnauthorized, constants.MsgUnauthorized))
		return
	}
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(constants.CodeValidationFailed, "password params invalid: "+err.Error()))
		return
	}
	if err := h.svc.ChangePassword(user.ID, req); err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"message": "password changed"})
}
