package middleware

import (
	"errors"
	"net/http"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/util"
	"github.com/gin-gonic/gin"
)

// ErrorHandler 统一错误处理：将 service 层错误转换为标准 JSON 响应。
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		ae := util.AsAppError(err)
		if errors.Is(err, util.ErrNotFound) {
			ae = util.NewAppError(constants.CodeNotFound, constants.MsgNotFound)
		}
		httpStatus := http.StatusInternalServerError
		switch {
		case ae.Code >= 40000 && ae.Code < 40100:
			httpStatus = http.StatusBadRequest
		case ae.Code >= 40100 && ae.Code < 40300:
			httpStatus = http.StatusUnauthorized
		case ae.Code >= 40300 && ae.Code < 40400:
			httpStatus = http.StatusForbidden
		case ae.Code >= 40400 && ae.Code < 40900:
			httpStatus = http.StatusNotFound
		case ae.Code >= 40900 && ae.Code < 42200:
			httpStatus = http.StatusConflict
		case ae.Code >= 42200 && ae.Code < 50000:
			httpStatus = http.StatusUnprocessableEntity
		}
		util.Fail(c, httpStatus, ae.Code, ae.Message)
		c.Abort()
	}
}
