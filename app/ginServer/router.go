/**
  create by yy on 2019-07-02
*/

package ginServer

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// If you want to use AutoRouter, please read README.md.

// func AutoRouter(controller interface{}) {
//	libs.AutoRoute(Router, controller)
// }

func GET(pattern string, function gin.HandlerFunc) {
	Router.GET(pattern, function)
}

func POST(pattern string, function gin.HandlerFunc) {
	Router.POST(pattern, function)
}

func PUT(pattern string, function gin.HandlerFunc) {
	Router.PUT(pattern, function)
}

func Group(pattern string, function ...gin.HandlerFunc) *gin.RouterGroup {
	return Router.Group(pattern, function...)
}

func DELETE(pattern string, function gin.HandlerFunc) {
	Router.DELETE(pattern, function)
}

func LoadHTMLGlob(pattern string) {
	Router.LoadHTMLGlob(pattern)
}

func StaticFS(pattern string, fs http.FileSystem) {
	Router.StaticFS(pattern, fs)
}
