package middleware

import (
	"log/slog"
	"time"

	"github.com/flowershop/backend/internal/constants"
	"github.com/gin-gonic/gin"
)

// RequestLogger 请求日志中间件：包含 request_id/method/path/status/latency_ms。
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info(constants.LogRequestHandled,
			"request_id", GetRequestID(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
		)
	}
}
