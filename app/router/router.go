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
	ginServer.GET("/test", controller.Test)
	ginServer.GET("/", controller.Index)
	ginServer.GET("/index", controller.Index)
	ginServer.GET("/init_table", controller.InitTables)
	fmt.Println("router初始化成功")
}
