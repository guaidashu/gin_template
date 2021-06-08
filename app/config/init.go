/**
  create by yy on 2019-07-29
*/

package config

import (
	"fmt"
	"github.com/guaidashu/go_helper/configor"
	"log"
	"os"
	"regexp"
	"strings"
)

type CustomConfig struct {
	Mysql       MysqlConf
	PostGreSql  PostGreSql
	App         AppConf
	Redis       RedisConf
	Mongodb     MongodbConf
	MiniProgram MiniProgramConf
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
	LogDir          string      `json:"log_dir"`
	RunAddress      string      `json:"run_address"`
	RunPort         interface{} `json:"run_port"`
	DEBUG           bool        `json:"debug"`
	TokenKey        string      `json:"token_key"`
	TokenExpireTime int64       `json:"token_expire_time"`
	Mode            string      `json:"mode"` // 配置文件环境
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

type MiniProgramConf struct {
	Appid           string `json:"appid"`
	Secret          string `json:"secret"`
	Token           string `json:"token"`
	TokenSecretKey  string `json:"token_secret_key"`
	TokenExpireTime int    `json:"token_expire_time"`
	Tokenissuer     string `json:"tokenissuer"`
}

var Config CustomConfig

func InitConf() {
	var (
		err  error
		pwd  string
		conf string
	)

	conf = os.Getenv("GIN_CONFIG")

	devMap := map[string]string{
		"debug":   "config_dev",
		"release": "config_product",
	}

	if conf == "" {
		conf = "debug"
	}

	if pwd, err = os.Getwd(); err != nil {
		log.Println("get config pwd error: ", err.Error())
		pwd = "."
	} else {
		pwd = strings.Replace(pwd, "\\", "/", -1)
		re3, _ := regexp.Compile("gin_template(.?)*([a-zA-Z])*")
		rep := re3.ReplaceAllStringFunc(pwd, func(s string) string {
			return ""
		})
		pwd = rep + "gin_template"
	}

	fmt.Println("开始加载开发配置文件")
	err = configor.Load(&Config, fmt.Sprintf("%s/app/config/%v.yml", pwd, devMap[conf]))
	if err != nil || Config.App.LogDir == "" {
		fmt.Println("配置文件加载失败")
	} else {
		fmt.Println("配置文件加载完成")
	}
}

func InitConfForTest(conf string) {
	var (
		err error
		pwd string
	)

	devMap := map[string]string{
		"debug":   "config_dev",
		"release": "config_product",
	}

	if conf == "" {
		conf = "debug"
	}

	if pwd, err = os.Getwd(); err != nil {
		log.Println("get config pwd error: ", err.Error())
		pwd = "."
	} else {
		pwd = strings.Replace(pwd, "\\", "/", -1)
		re3, _ := regexp.Compile("gin_template(.?)*([a-zA-Z])*")
		rep := re3.ReplaceAllStringFunc(pwd, func(s string) string {
			return ""
		})
		pwd = rep + "gin_template"
	}

	fmt.Println("开始加载开发配置文件")
	err = configor.Load(&Config, fmt.Sprintf("%s/app/config/%v.yml", pwd, devMap[conf]))
	if err != nil || Config.App.LogDir == "" {
		fmt.Println("配置文件加载失败")
	} else {
		fmt.Println("配置文件加载完成")
	}
}
