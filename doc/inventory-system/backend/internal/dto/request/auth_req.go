package request

// LoginReq 登录请求参数
type LoginReq struct {
	Username string `json:"username" binding:"required,min=1,max=50"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}
