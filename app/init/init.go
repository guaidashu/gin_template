/**
  create by yy on 2019-08-23
*/

package init

import (
	"fmt"
	"gin_template/app/config"
	_ "gin_template/app/config"
	"gin_template/app/libs"
	"gin_template/app/models"
	"log"
)

func init() {
	config.InitConf()
	// 初始化日志
	err, _ := libs.InitLogger()
	if err != nil {
		log.Println(fmt.Sprintf("init logger failed, error: %v", libs.NewReportError(err)))
		return
	}
	libs.Logger.Info("======= 初始化日志系统 ======")
	// 初始化redis
	// libs.Logger.Info("====== 初始化redis系统 ======")
	// err = rds.InitRedis()
	// if err != nil {
	// 	libs.Logger.Info(fmt.Sprintf("init rds failed, error: %v", libs.NewReportError(err)))
	// }
	// 初始化mysql
	libs.Logger.Info("====== 初始化mysql系统 ======")
	err = models.InitDB()
	if err != nil {
		libs.Logger.Info(fmt.Sprintf("init db failed, error: %v", libs.NewReportError(err)))
	}

	// 初始化postgresql系统
	// libs.Logger.Info("====== 初始化postgresql系统 ======")
	// err = models.InitPostGreDB()
	// if err != nil {
	// 	libs.Logger.Info(fmt.Sprintf("init db failed, error: %v", libs.NewReportError(err)))
	// }
	// 自动建表(目前仅针对于 mysql 和 postgresql 可开启此功能)， 或者通过 router里配置的init_table 可视化访问来创建
	// models.CreateTable()

	// 初始化mongodb
	// libs.Logger.Info("====== 初始化mongodb系统 ======")
	// mongodb.InitMongoDB()

	// 初始化 kafka
	// mq.InitKafka()
}
