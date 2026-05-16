package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"inventory-system/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// Recovery 异常恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[PANIC] %v\n%s\n", err, debug.Stack())

				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    errcode.ErrServer.Code,
					"message": "服务器内部错误",
					"data":    nil,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
