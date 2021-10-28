/**
  create by yy on 2019-09-25
*/

package rds

import (
	"gin_template/app/config"
	"github.com/go-redis/redis"
)

// redis 模块包
var Redis *redis.Client

func InitRedis() (err error) {
	Redis, err = getConnect()
	return
}

func getConnect() (*redis.Client, error) {
	rds := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.RedisHost + ":" + config.Config.Redis.RedisPort,
		Password: config.Config.Redis.RedisPassword,
	})
	return rds, nil
}
