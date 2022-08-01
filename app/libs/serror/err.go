/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 22/07/2021
 * @Desc: 自定义error 返回结构体
 */

package serror

import (
	"fmt"
	"gin_template/app/config"
	"runtime"
	"strings"
)

type (
	Error interface {
		error
		// 绑定err
		SetErr(value interface{}) Error
		// 设置返回code
		SetCode(code int64) Error
		// 获取设置的code
		Code() int64
		// 设置返回提示msg
		SetMsg(msg string) Error
		// 获取返回提示msg
		Msg() string
		String() string
		// 获取返回的所有提示msg
		ErrMsg() []string
	}

	defaultErr struct {
		msg    string   // 报错信息
		err    error    // 报错error
		code   int64    // 状态码
		errMsg []string // 报错信息
		isLog  bool     // 是否打印日志
	}
)

func NewErr(isLog ...bool) Error {
	var (
		logFlag bool
	)

	if len(isLog) > 0 {
		logFlag = isLog[0]
	}
	return &defaultErr{
		isLog: logFlag,
	}
}

func (e *defaultErr) SetErr(value interface{}) Error {
	// 这里要考虑到多层的问题
	tmpErr, ok := value.(Error)
	if !ok {
		e.msg = fmt.Sprintf("%v", value)
		e.errMsg = append(e.errMsg, fmt.Sprintf("%s: %v", getCaller(value), value))
		return e
	}

	if e.Msg() == "" {
		_ = e.SetMsg(tmpErr.Msg())
	}
	for _, v := range tmpErr.ErrMsg() {
		e.errMsg = append(e.errMsg, v)
	}

	e.errMsg = append(e.errMsg, fmt.Sprintf("%s", getCaller(value)))

	return e
}

func (e *defaultErr) Code() int64 {
	return e.code
}

func (e *defaultErr) SetCode(code int64) Error {
	e.code = code
	return e
}

func (e *defaultErr) SetMsg(msg string) Error {
	e.msg = msg
	return e
}

func (e *defaultErr) Msg() string {
	return e.msg
}

func (e *defaultErr) Error() string {
	return e.String()
}

func (e *defaultErr) String() string {
	var s strings.Builder
	for _, v := range e.errMsg {
		s.WriteString(v)
	}
	return s.String()
}

func (e *defaultErr) ErrMsg() []string {
	return e.errMsg
}

func getCaller(value interface{}) string {
	if !config.Config.App.DEBUG {
		return ""
	}
	// 这里用于打印报错文件及行数
	_, fileName, line, _ := runtime.Caller(2)
	return fmt.Sprintf("report in: %v: in line %v", fileName, line)
}
