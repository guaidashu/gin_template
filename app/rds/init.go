/**
  create by yy on 2019-09-25
*/

package rds

import (
	"gin_template/app/config"
	"gin_template/app/data_struct/_interface"
	"gin_template/app/enum"
	"gin_template/app/libs"
	"github.com/go-redis/redis"
)


// redis 模块包
type RedisInit struct{}

func NewRedisInit() *RedisInit {
	return &RedisInit{}
}

func (r *RedisInit) Init(*_interface.ServiceParam) error {
	return InitRedis()
}

func (r *RedisInit) ComponentName() enum.BootModuleType {
	return enum.RedisInit
}

func (r *RedisInit) Close() error {
	if Redis != nil {
		libs.Logger.Info("Close rds")
		if err := Redis.Close(); err != nil {
			libs.Logger.Info("Close rds failed, error: %v", err)
			return nil
		}
	}

	return nil
}

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
