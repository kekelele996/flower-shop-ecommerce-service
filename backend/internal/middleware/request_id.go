package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const ContextRequestIDKey = "request_id"

// RequestID 为每个请求注入 request_id，并写入响应头。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Set(ContextRequestIDKey, rid)
		c.Header("X-Request-ID", rid)
		c.Next()
	}
}

// GetRequestID 从上下文取 request_id。
func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get(ContextRequestIDKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
