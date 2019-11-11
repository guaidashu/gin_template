/**
  create by yy on 2019-07-29
*/

package libs

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Reply struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(ctx *gin.Context, data interface{}) {
	r := &Reply{
		Code: 0,
		Msg:  "",
		Data: data,
	}
	ctx.JSON(http.StatusOK, r)
}

func Error(ctx *gin.Context, msg string) {
	r := &Reply{
		Code: 1,
		Msg:  msg,
		Data: "",
	}
	ctx.JSON(http.StatusOK, r)
}
