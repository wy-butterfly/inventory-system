package middleware

import (
	"strings"

	"inventory-system/pkg/errcode"
	"inventory-system/pkg/jwt"
	apiResp "inventory-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apiResp.Error(c, errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			apiResp.Error(c, errcode.ErrUnauthorized)
			c.Abort()
			return
		}
		tokenString := parts[1]

		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			apiResp.Error(c, errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
