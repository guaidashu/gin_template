package ws

import "gin_template/app/enum"

type WsHandlerMiddleware func(ctx *Context) error

type WsHandlerFunc func(ctx *Context)

type WsHandler struct {
	// 事件名
	EventName enum.WsEventEnum
	// 处理方法
	Handler WsHandlerFunc
}

type WsEventHandler struct {
	Handler     *WsHandler
	Middlewares []WsHandlerMiddleware
}

func (w *WsEventHandler) Use(middlewares ...WsHandlerMiddleware) *WsEventHandler {
	w.Middlewares = append(w.Middlewares, middlewares...)
	return w
}
