package middlewares

import (
	"gin_template/app/libs"
	"gin_template/app/libs/serror"
	"gin_template/app/ws"
)

// 注册client到用户map
// 调用此中间件前请先调用 ValidateWsToken 方法设置userId
func RegisterWsClient(ctx *ws.Context) error {
	userIdIt, err := ctx.Get("userId")
	if err != nil {
		libs.Logger.Error("用户id不存在，请检查是否登录")
		return err
	}

	if userIdIt == nil {
		return serror.ErrUserNotFound
	}

	if ctx.Client().GetUserId() != 0 {
		return nil
	}

	ws.NewClientPool().SetClientByUserId(ctx.ClientId(), userIdIt.(int64))

	return nil
}
