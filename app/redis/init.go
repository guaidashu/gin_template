/**
  create by yy on 2019-09-25
*/

package redis

import (
	"gin_template/app/config"
	"github.com/guaidashu/go_helper/redis_tool"
)

var Redis *redis_tool.RedisConn

func InitRedis() {
	Redis = getConnect()
}

func getConnect() *redis_tool.RedisConn {
	rds := &redis_tool.RedisConn{}
	conf := &redis_tool.RedisConfig{
		Host:     config.Config.Redis.RedisHost,
		Port:     config.Config.Redis.RedisPort,
		Password: config.Config.Redis.RedisPassword,
	}
	rds.Init(conf)
	return rds
}
