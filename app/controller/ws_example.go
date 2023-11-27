package controller

import (
	"gin_template/app/data_struct/requests"
	"gin_template/app/ws"
)

type (
	WsExampleController struct {
		BaseController
	}
)

func (c *WsExampleController) TestWsRouter(ctx *ws.Context) {
	req := &requests.TestWsReq{}
	err := ctx.BindJson(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	userId, err := c.getWsUserId(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Send(userId)
}
