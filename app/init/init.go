/**
  create by yy on 2019-08-23
*/

package init

import (
	"fmt"
	_ "gin_template/app/config"
	"gin_template/app/libs"
	"gin_template/app/models"
	"gin_template/app/redis"
	"log"
)

func init() {
	// 初始化日志
	err, _ := libs.InitLogger()
	if err != nil {
		log.Println(fmt.Sprintf("init logger failed, error: %v", err))
		return
	}
	libs.Logger.Info("======= 初始化日志系统 ======")
	// 初始化redis
	libs.Logger.Info("====== 初始化redis系统 ======")
	redis.InitRedis()
	// 初始化mysql
	libs.Logger.Info("====== 初始化mysql系统 ======")
	err = models.InitDB()
	if err != nil {
		libs.Logger.Info(fmt.Sprintf("init db failed, error: %v", err))
	}
	libs.Logger.Info("====== 初始化postgresql系统 ======")
	err = models.InitPostGreDB()
	if err != nil {
		libs.Logger.Info(fmt.Sprintf("init db failed, error: %v", err))
	}
	// 自动建表， 第一次运行了之后可以注释掉，或者通过 router里配置的init_table 可视化访问来创建
	// models.CreateTable()
}
