package serror

var (
	ErrUserNotFound            = NewError(10001, "用户不存在")
	ErrParamsError             = NewError(10002, "参数错误")
	ErrContextKeyNotExistError = NewError(10003, "上下文参数不存在")
	ErrMethodNotExist          = NewError(10004, "方法不存在")
	ErrDataIsNil               = NewError(10005, "数据为空")
	ErrNoAuth                  = NewError(10006, "未登录")
)
