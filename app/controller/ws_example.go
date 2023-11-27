package controller

import "gin_template/app/ws"

type (
	WsExampleController struct {
	}
)

func (c *WsExampleController) TestWsRouter(ctx *ws.Context) {
	ctx.Send("ok")
}
