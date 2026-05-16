package request

// CreateUserReq 创建用户请求
type CreateUserReq struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Password    string `json:"password" binding:"required,min=6,max=20"`
	DisplayName string `json:"display_name" binding:"omitempty,max=100"`
	Role        string `json:"role" binding:"required,oneof=admin operator viewer"`
}

// UpdateUserReq 编辑用户请求
type UpdateUserReq struct {
	DisplayName *string `json:"display_name" binding:"omitempty,max=100"`
	Role        *string `json:"role" binding:"omitempty,oneof=admin operator viewer"`
	Status      *int8   `json:"status" binding:"omitempty,oneof=0 1"`
}

// ResetPasswordReq 重置密码请求
type ResetPasswordReq struct {
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
}

// UserQuery 用户列表查询
type UserQuery struct {
	PageQuery
	Keyword string `form:"keyword" binding:"omitempty,max=100"`
	Role    string `form:"role" binding:"omitempty"`
	Status  *int8  `form:"status" binding:"omitempty,oneof=0 1"`
}
