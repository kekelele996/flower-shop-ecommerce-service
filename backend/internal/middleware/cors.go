package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件。
func CORS(origins string) gin.HandlerFunc {
	cfg := cors.DefaultConfig()
	if origins == "*" {
		cfg.AllowAllOrigins = true
	} else {
		cfg.AllowOrigins = []string{origins}
	}
	cfg.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	cfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"}
	cfg.ExposeHeaders = []string{"X-Request-ID"}
	return cors.New(cfg)
}
