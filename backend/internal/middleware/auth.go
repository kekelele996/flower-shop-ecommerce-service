package middleware

import (
	"strings"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/util"
	"github.com/gin-gonic/gin"
)

const ContextUserKey = "auth_user"

// AuthUser 认证后的用户上下文。
type AuthUser struct {
	ID       uint
	Username string
	Role     string
}

// Auth JWT 认证中间件。
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			util.Fail(c, 401, constants.CodeUnauthorized, constants.MsgUnauthorized)
			c.Abort()
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := util.ParseToken(jwtSecret, token)
		if err != nil {
			util.Fail(c, 401, constants.CodeTokenExpired, "token invalid or expired")
			c.Abort()
			return
		}
		c.Set(ContextUserKey, &AuthUser{ID: claims.UserID, Username: claims.Username, Role: claims.Role})
		c.Next()
	}
}

// CurrentUser 从上下文取当前用户。
func CurrentUser(c *gin.Context) *AuthUser {
	v, ok := c.Get(ContextUserKey)
	if !ok {
		return nil
	}
	u, ok := v.(*AuthUser)
	if !ok {
		return nil
	}
	return u
}
