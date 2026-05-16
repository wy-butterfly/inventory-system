package handler

import (
	"inventory-system/internal/dto/request"
	"inventory-system/internal/service"
	"inventory-system/pkg/errcode"
	apiResp "inventory-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

// Login 用户登录
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("请输入用户名和密码"))
		return
	}

	ip := getClientIP(c)
	token, user, errCode := h.authService.Login(req.Username, req.Password, ip)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"display_name": user.DisplayName,
			"role":         user.Role,
		},
	})
}

// GetMe 获取当前登录用户信息
// GET /api/v1/auth/me
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, errCode := h.authService.GetCurrentUser(userID)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.Success(c, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"display_name": user.DisplayName,
		"role":         user.Role,
		"status":       user.Status,
	})
}

// Logout 用户登出
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	apiResp.SuccessWithMessage(c, "登出成功", nil)
}
