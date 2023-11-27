package controller

import (
	"gin_template/app/libs/serror"
	"gin_template/app/ws"
	"github.com/gin-gonic/gin"
)

type (
	BaseController struct {
	}
)

func (c *BaseController) getUserId(ctx *gin.Context) (int64, error) {
	userIdStr, ok := ctx.Get("userId")
	if !ok {
		return 0, serror.ErrUserNotFound
	}

	userId, ok := userIdStr.(int64)
	if !ok {
		return 0, serror.ErrUserNotFound
	}

	return userId, nil
}

func (c *BaseController) getWsUserId(ctx *ws.Context) (int64, error) {
	userIdStr, err := ctx.Get("userId")
	if err != nil {
		return 0, serror.ErrUserNotFound
	}

	userId, ok := userIdStr.(int64)
	if !ok {
		return 0, serror.ErrUserNotFound
	}

	return userId, nil
}
