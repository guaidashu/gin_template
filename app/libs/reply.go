/**
  create by yy on 2019-07-29
*/

package libs

import (
	"fmt"
	"gin_template/app/libs/serror"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Reply struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Err  interface{} `json:"err"`
	Data interface{} `json:"data"`
}

func Success(ctx *gin.Context, data interface{}) {
	r := &Reply{
		Code: 0,
		Msg:  "",
		Data: data,
	}
	ctx.JSON(http.StatusOK, r)
	ctx.Abort()
}

func Error(ctx *gin.Context, err error, code ...int) {
	customErr := err.(serror.Error)

	r := &Reply{
		Msg: customErr.Msg(),
		Err: customErr.ErrMsg(),
	}
	if len(code) > 0 {
		r.Code = code[0]
	} else {
		r.Code = 1
	}

	ctx.JSON(http.StatusOK, r)
	ctx.Abort()
}

func ErrorWithCode(ctx *gin.Context, msg string, code int, err ...error) {
	r := &Reply{
		Code: code,
		Msg:  msg,
		Data: "",
	}
	if len(err) > 0 {
		r.Err = fmt.Sprintf("%v", err[0])
	}

	ctx.JSON(http.StatusOK, r)
	ctx.Abort()
}

func CustomReply(ctx *gin.Context, code int, msg string, err string, data ...interface{}) {
	var (
		replyData interface{}
	)

	if len(data) > 0 {
		replyData = data[0]
	} else {
		replyData = ""
	}

	r := &Reply{
		Code: code,
		Msg:  msg,
		Data: replyData,
		Err:  err,
	}

	ctx.JSON(http.StatusOK, r)
	ctx.Abort()
}
