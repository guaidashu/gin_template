package serror

var (
	ErrUserNotFound = NewError(10001, "用户不存在")
	ErrParamsError  = NewError(10002, "参数错误")
)
