package errcode

import "fmt"

// ErrCode 错误码类型
type ErrCode struct {
	Code    int
	Message string
}

func (e *ErrCode) Error() string {
	return fmt.Sprintf("错误码: %d, 信息: %s", e.Code, e.Message)
}

func (e *ErrCode) WithMessage(msg string) *ErrCode {
	return &ErrCode{Code: e.Code, Message: msg}
}

var (
	Success = &ErrCode{Code: 0, Message: "success"}

	ErrInvalidParam   = &ErrCode{Code: 40001, Message: "参数校验失败"}
	ErrDuplicate      = &ErrCode{Code: 40002, Message: "数据重复"}
	ErrNotFound       = &ErrCode{Code: 40003, Message: "数据不存在"}
	ErrOptimisticLock = &ErrCode{Code: 40004, Message: "数据已被其他人修改，请刷新后重试"}

	ErrUnauthorized    = &ErrCode{Code: 40101, Message: "未登录或Token已过期"}
	ErrForbidden       = &ErrCode{Code: 40102, Message: "权限不足"}
	ErrLoginFailed     = &ErrCode{Code: 40103, Message: "用户名或密码错误"}
	ErrAccountDisabled = &ErrCode{Code: 40104, Message: "账号已被停用"}

	ErrFileFormat = &ErrCode{Code: 40301, Message: "不支持的文件格式"}
	ErrFileSize   = &ErrCode{Code: 40302, Message: "文件大小超过限制"}
	ErrFileUpload = &ErrCode{Code: 40303, Message: "文件上传失败"}

	ErrServer   = &ErrCode{Code: 50001, Message: "服务器内部错误"}
	ErrDatabase = &ErrCode{Code: 50002, Message: "数据库错误"}
)
