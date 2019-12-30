/**
  create by yy on 2019-07-29
*/

package config

import (
	"fmt"
	"github.com/guaidashu/go_helper/configor"
)

type CustomConfig struct {
	Mysql      MysqlConf
	PostGreSql PostGreSql
	App        AppConf
	Redis      RedisConf
	Mongodb    MongodbConf
}

type MysqlConf struct {
	Database   string `json:"database"`
	DbHost     string `json:"db_host"`
	DbPassword string `json:"db_password"`
	DbUsername string `json:"db_username"`
	DbPort     string `json:"db_port"`
	DbPoolSize int    `json:"db_pool_size"`
}

type AppConf struct {
	LogDir     string      `json:"log_dir"`
	RunAddress string      `json:"run_address"`
	RunPort    interface{} `json:"run_port"`
	DEBUG      bool        `json:"debug"`
}

type RedisConf struct {
	RedisHost     string `json:"redis_host"`
	RedisPort     string `json:"redis_port"`
	RedisPassword string `json:"redis_password"`
}

type PostGreSql struct {
	Database string `json:"database"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Username string `json:"username"`
	Port     string `json:"port"`
	PoolSize int    `json:"pool_size"`
}

type MongodbConf struct {
	Database string      `json:"database"`
	Host     string      `json:"host"`
	Password string      `json:"password"`
	Username string      `json:"username"`
	Port     interface{} `json:"port"`
	PoolSize int         `json:"pool_size"`
}

var Config CustomConfig

func init() {
	fmt.Println("开始加载开发配置文件")
	err := configor.Load(&Config, "app/config/config_dev.yml")
	if err != nil || Config.App.LogDir == "" {
		fmt.Println("开发环境配置文件加载失败")
		err = configor.Load(&Config, "app/config/config_product.yml")
		if err != nil || Config.App.LogDir == "" {
			fmt.Println("线上环境配置文件加载失败")
			fmt.Println("配置文件加载失败")
		} else {
			fmt.Println("线上环境配置文件加载成功")
		}
	} else {
		fmt.Println("开发环境配置文件加载完成")
	}
}
