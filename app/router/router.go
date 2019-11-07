/**
  create by yy on 2019-07-02
*/

package router

import (
	"fmt"
	"gin_template/app/controller"
	"gin_template/app/ginServer"
)

func init() {
	fmt.Println("开始初始化router")

	// 添加 html template
	// ginServer.LoadHTMLGlob("app/views/**/*")

	// 添加 静态资源路由
	// ginServer.StaticFS("/asset", http.Dir("app/asset"))

	test := ginServer.Group("/test")
	{
		test.GET("/", controller.Test)
	}
	ginServer.GET("/", controller.Index)
	ginServer.GET("/index", controller.Index)
	ginServer.GET("/init_table", controller.InitTables)
	fmt.Println("router初始化成功")
}
