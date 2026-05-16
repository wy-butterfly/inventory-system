package middleware

import (
	"inventory-system/pkg/errcode"
	apiResp "inventory-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequireRoles 创建角色检查中间件
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")

		allowed := false
		for _, role := range roles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			apiResp.Error(c, errcode.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin 仅管理员可访问
func RequireAdmin() gin.HandlerFunc {
	return RequireRoles("admin")
}
