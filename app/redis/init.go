/**
  create by yy on 2019-09-25
*/

package redis

import (
	"gin_template/app/config"
	"gin_template/app/enum"
	"gin_template/app/libs"
	"github.com/go-redis/redis"
	"time"
)

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

func AcquireLock(resource string, acquireTimeout, lockTimeout int) (string, bool) {
	r := Redis
	if acquireTimeout <= 0 {
		acquireTimeout = enum.DEFAULT_LOCK_ACQUIRE_TIMEOUT
	}
	if lockTimeout <= 0 {
		lockTimeout = enum.DEFAULT_LOCK_KEY_TIMEOUT
	}

	lockResource := enum.LOCK_PREFIX + resource
	val := libs.GenerateDataId()
	lockTimeoutD := time.Duration(lockTimeout) * time.Second
	endTime := time.Now().Add(time.Duration(acquireTimeout) * time.Second)
	for time.Now().Unix() < endTime.Unix() {
		ok, err := r.SetNX(lockResource, val, lockTimeoutD).Result()
		if err != nil {
			libs.Logger.Errorf("设置[%s]的锁失败, %s", resource, err.Error())
			return "", false
		}

		if ok {
			return val, true
		} else {
			time.Sleep(10 * time.Millisecond)
			continue
		}
	}
	return "", false
}

func ReleaseLock(resource, val string) bool {
	r := Redis
	lockResource := enum.LOCK_PREFIX + resource
	v, err := r.Get(lockResource).Result()
	if err != nil && err != redis.Nil {
		libs.Logger.Errorf("释放[%s]的锁失败, %s", resource, err.Error())
		return false
	}

	if err == redis.Nil {
		return true
	}
	redis.NewStringStructMapCmd()
	if v == val {
		r.Del(lockResource)
		return true
	} else {
		// 数据已被其他人加锁，那么此处可以认为是 ok 的
		return true
	}
}
