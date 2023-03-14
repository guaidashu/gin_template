package controller

import (
	"gin_template/app/libs/serror"
	"github.com/gin-gonic/gin"
)

type (
	BaseController struct {
	}
)

func (c *BaseController) GetUserId(ctx *gin.Context) (int64, error) {
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
