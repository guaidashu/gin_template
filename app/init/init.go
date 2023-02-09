/**
  create by yy on 2019-08-23
*/

package init

import (
	"fmt"
	"gin_template/app/data_struct/_interface"
	"gin_template/app/enum"
	"gin_template/app/libs"
	"gin_template/app/models"
	"gin_template/app/mongodb"
	"gin_template/app/mq/kafka"
	"gin_template/app/nacos"
	"gin_template/app/rds"
	"math/rand"
	"time"
)

func init() {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 初始化配置文件，如果走的是配置文件的话
	// config.InitConf()

	err, _ := libs.InitLogger("logs")
	if err != nil {
		panic("初始化日志系统失败")
	}

	initModules := InitModules()
	initModuleTemp := make(map[enum.BootModuleType]struct{})
	for _, moduleType := range initModules {
		initModuleTemp[moduleType] = struct{}{}
	}

	// 初始化各种组件
	for _, v := range InitList() {
		if _, ok := initModuleTemp[v.ComponentName()]; ok {
			err := v.Init(&_interface.ServiceParam{})
			if err != nil {
				libs.Logger.Panic(fmt.Sprintf("%v: 初始化失败, err: %v", v.ComponentName(), err))
			}
		}
	}

	// 自动建表(目前仅针对于 mysql 和 postgresql 可开启此功能)， 或者通过 router里配置的init_table 可视化访问来创建
	// models.CreateTable()
}

func InitList() []_interface.ComponentsInit {
	// nacos必须放在第一位
	return []_interface.ComponentsInit{
		nacos.NewNacosInit(),
		nacos.NewConfigInit(),
		rds.NewRedisInit(),
		models.NewMysqlInit(),
		models.NewPsqlInit(),
		kafka.NewKqInit(),
		mongodb.NewMDBInit(),
	}
}

func InitModules() []enum.BootModuleType {
	return []enum.BootModuleType{
		enum.NacosInit,
		enum.ConfigInit,
		enum.RedisInit,
		enum.MysqlInit,
	}
}
