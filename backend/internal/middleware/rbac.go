package middleware

import (
	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/util"
	"github.com/gin-gonic/gin"
)

// RBAC 角色权限校验中间件：仅允许指定角色访问。
func RBAC(roles ...string) gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		user := CurrentUser(c)
		if user == nil {
			util.Fail(c, 401, constants.CodeUnauthorized, constants.MsgUnauthorized)
			c.Abort()
			return
		}
		if !allowed[user.Role] {
			util.Fail(c, 403, constants.CodeForbidden, "forbidden: role="+user.Role+" not allowed")
			c.Abort()
			return
		}
		c.Next()
	}
}
