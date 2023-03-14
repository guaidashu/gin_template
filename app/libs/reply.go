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
	Msg  string      `json:"msg,omitempty"`
	Err  interface{} `json:"err,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func Success(ctx *gin.Context, data interface{}) {
	r := &Reply{
		Code: 0,
		Data: data,
	}
	ctx.JSON(http.StatusOK, r)
	ctx.Abort()
}

func Error(ctx *gin.Context, err error, code ...int) {
	customErr, ok := err.(serror.Error)
	if !ok {
		customErr = serror.NewErr().SetErr(err)
	}

	r := &Reply{
		Msg:  customErr.Msg(),
		Err:  customErr.ErrMsg(),
		Code: int(customErr.Code()),
	}
	if len(code) > 0 {
		r.Code = code[0]
	}
	if r.Code == 0 {
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

// 某些地方可能需要直接获取到回复结构体， 比如websocket
func GetSuccessReply(data interface{}) *Reply {
	return &Reply{
		Code: 0,
		Msg:  "",
		Data: data,
	}
}

func GetErrorReply(err error, code ...int) *Reply {
	customErr, ok := err.(serror.Error)
	if !ok {
		customErr = serror.NewErr().SetErr(err)
	}

	r := &Reply{
		Msg: customErr.Msg(),
		Err: customErr.ErrMsg(),
	}
	if len(code) > 0 {
		r.Code = code[0]
	} else {
		r.Code = 1
	}

	return r
}
