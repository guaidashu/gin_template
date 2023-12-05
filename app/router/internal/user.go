/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 25/10/2021
 * @Desc: 用户路由
 */

package internal

import (
	"gin_template/app/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRouter(engine *gin.RouterGroup) {
	// 普通端用户接口
	userRouter := engine.Group("/user").Use(middlewares.ValidateUser())
	userRouterFunc(userRouter)
	userRouterNoLogin := engine.Group("/user")
	userRouterFuncNoLogin(userRouterNoLogin)
	tokenRouter := engine.Group("/user").Use(middlewares.ValidateRefreshToken())
	tokenRouterFunc(tokenRouter)

	// 管理后台用户接口
	adminRouter := engine.Group("/admin/user").Use(middlewares.ValidateAdmin())
	adminUserRouter(adminRouter)
	adminRouterNoLogin := engine.Group("/admin/user")
	adminUserRouterNoLogin(adminRouterNoLogin)
}

// 小程序端用户相关路由
func userRouterFunc(user gin.IRoutes) {
	// userController := new(controller.UserController)
	// user.GET("/user-info", userController.UserInfo) // 用户信息
}

func tokenRouterFunc(user gin.IRoutes) {
	// loginController := new(controller.LoginController)
	// user.POST("/get-token", loginController.GetToken)                // 获取token
	// user.POST("/get-refresh-token", loginController.GetRefreshToken) // 获取刷新token
}

func userRouterFuncNoLogin(user gin.IRoutes) {
	// loginController := new(controller.LoginController)
	// userController := new(controller.UserController)
	// user.POST("/login", loginController.Login)      // 登录
	// user.POST("/register", userController.Register) // 注册
}

func adminUserRouter(user gin.IRoutes) {
	// adminUserController := new(controller.AdminUserController)
	// userController := new(controller.UserController)
	//
	// user.GET("/get-user-info", adminUserController.AdminUserInfo) // 获取管理员用户信息
	// user.GET("/get-user-list", userController.MgrUserList)        // 获取用户列表
}

func adminUserRouterNoLogin(user gin.IRoutes) {
	// loginController := new(controller.LoginController)
	// user.POST("/login", loginController.AdminLogin) // 登录
}
