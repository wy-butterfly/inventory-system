package handler

import (
	"strconv"

	"inventory-system/internal/dto/request"
	"inventory-system/internal/service"
	"inventory-system/pkg/errcode"
	apiResp "inventory-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// Create 创建用户
// POST /api/v1/users
func (h *UserHandler) Create(c *gin.Context) {
	var req request.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage(err.Error()))
		return
	}

	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	user, errCode := h.userService.Create(&req, userID, userName)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.SuccessWithMessage(c, "创建成功", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

// Update 编辑用户
// PUT /api/v1/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("无效的用户ID"))
		return
	}

	var req request.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage(err.Error()))
		return
	}

	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	errCode := h.userService.Update(uint(id), &req, userID, userName)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.SuccessWithMessage(c, "编辑成功", nil)
}

// ResetPassword 重置密码
// PUT /api/v1/users/:id/reset-password
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("无效的用户ID"))
		return
	}

	var req request.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("请输入新密码（6-20位）"))
		return
	}

	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	errCode := h.userService.ResetPassword(uint(id), req.NewPassword, userID, userName)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.SuccessWithMessage(c, "密码重置成功", nil)
}

// List 用户列表
// GET /api/v1/users?page=1&page_size=20
func (h *UserHandler) List(c *gin.Context) {
	var query request.UserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage(err.Error()))
		return
	}

	users, total, errCode := h.userService.List(&query)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.PageSuccess(c, users, total, query.GetPage(), query.GetPageSize())
}
