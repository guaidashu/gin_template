/**
  create by yy on 2019-07-02
*/

package router

import (
	"fmt"
	"gin_template/app/controller"
	"gin_template/app/gin_server"
)

func init() {
	fmt.Println("开始初始化router")

	// 添加 html template
	// gin_server.LoadHTMLGlob("app/views/**/*")

	// 添加 静态资源路由
	// gin_server.StaticFS("/asset", http.Dir("app/asset"))

	// 开启跨域
	// gin_server.Router.Use(middlewares.Cors())

	// 创建websocket服务
	// gin_server.GET("/ws", (&controller.WsController{}).WsHandler)
	// wsSrv := ws.NewWsSrv()
	//
	// wsExampleController := &controller.WsExampleController{}
	// wsSrv.Register(&ws.WsHandler{
	// 	EventName: enum.WsExampleEvent,
	// 	Handler:   wsExampleController.TestWsRouter,
	// })

	test := gin_server.Group("/test")
	{
		test.GET("/", controller.Test)
	}
	gin_server.GET("/", controller.Index)
	gin_server.GET("/index", controller.Index)
	gin_server.GET("/init_table", controller.InitTables)
	fmt.Println("router初始化成功")
}
