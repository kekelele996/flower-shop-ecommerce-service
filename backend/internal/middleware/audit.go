package middleware

import (
	"encoding/json"
	"strings"

	"github.com/flowershop/backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Audit 操作审计中间件：记录写操作（POST/PUT/PATCH/DELETE）到审计表。
func Audit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		body := ""
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			bs := make([]byte, c.Request.ContentLength)
			_, _ = c.Request.Body.Read(bs)
			body = string(bs)
		}
		c.Next()

		user := CurrentUser(c)
		userID := uint(0)
		username := ""
		if user != nil {
			userID = user.ID
			username = user.Username
		}
		entity := inferEntity(c.FullPath())
		action := inferAction(c.Request.Method)
		_ = db.Create(&model.AuditLog{
			UserID:    userID,
			Username:  username,
			Action:    action,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Entity:    entity,
			EntityID:  c.Param("id"),
			Detail:    truncateDetail(body),
			RequestID: GetRequestID(c),
			IP:        c.ClientIP(),
		})
	}
}

func inferEntity(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for _, p := range parts {
		if p == "api" || p == "v1" || p == "admin" || p == "auth" || strings.HasPrefix(p, ":") || p == "uploads" {
			continue
		}
		return p
	}
	return "unknown"
}

func inferAction(method string) string {
	switch method {
	case "POST":
		return "CREATE"
	case "PUT", "PATCH":
		return "UPDATE"
	case "DELETE":
		return "DELETE"
	default:
		return method
	}
}

func truncateDetail(s string) string {
	_ = json.Valid([]byte(s))
	if len(s) > 500 {
		return s[:500]
	}
	return s
}
