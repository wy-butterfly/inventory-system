package response

import (
	"net/http"

	resp "inventory-system/internal/dto/response"
	"inventory-system/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, resp.Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 返回成功响应（带自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, resp.Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, err *errcode.ErrCode) {
	c.JSON(http.StatusOK, resp.Response{
		Code:    err.Code,
		Message: err.Message,
		Data:    nil,
	})
}

// ErrorWithMessage 返回错误响应（带自定义消息）
func ErrorWithMessage(c *gin.Context, err *errcode.ErrCode, message string) {
	c.JSON(http.StatusOK, resp.Response{
		Code:    err.Code,
		Message: message,
		Data:    nil,
	})
}

// PageSuccess 返回分页成功响应
func PageSuccess(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, resp.Response{
		Code:    0,
		Message: "success",
		Data: resp.PageResult{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}
