/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 28/10/2021
 * @Desc: redis分布式锁
 */

package rds

import (
	"context"
	"errors"
	"gin_template/app/enum"
	"gin_template/app/libs"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	Locker interface {
		// 加锁
		Acquire() (bool, error)
		// 解锁
		Release() bool
		// 设置基础数据
	}

	RedisLock struct {
		key            string
		val            string
		acquireTimeout int
		lockTimeout    int
		cache          *redis.Client
	}
)

func NewRedisLock(key string) Locker {
	return &RedisLock{
		key:            key,
		acquireTimeout: enum.DefaultLockAcquireTimeout,
		lockTimeout:    enum.DefaultLockKeyTimeout,
		cache:          Redis,
	}
}

// 分布式锁
func (r *RedisLock) Acquire() (bool, error) {
	ctx := context.Background()
	r.val = libs.GenerateDataId()
	lockResource := enum.LockPrefix + r.key
	lockTimeoutD := time.Duration(r.lockTimeout) * time.Second
	endTime := time.Now().Add(time.Duration(r.acquireTimeout) * time.Second)

	for time.Now().Unix() < endTime.Unix() {
		ok, err := r.cache.SetNX(ctx, lockResource, r.val, lockTimeoutD).Result()
		if err != nil {
			libs.Logger.Errorf("设置[%s]的锁失败, %s", r.key, err.Error())
			return false, err
		}

		if ok {
			return true, nil
		} else {
			time.Sleep(10 * time.Millisecond)
			continue
		}
	}
	return false, errors.New("加锁失败")
}

func (r *RedisLock) Release() bool {
	ctx := context.Background()
	lockResource := enum.LockPrefix + r.key
	v, err := r.cache.Get(ctx, lockResource).Result()
	if err != nil && err != redis.Nil {
		libs.Logger.Errorf("释放[%s]的锁失败, %s", r.key, err.Error())
		return false
	}

	if err == redis.Nil {
		return true
	}
	if v == r.val {
		r.cache.Del(ctx, lockResource)
	}

	// 数据已被其他人加锁，那么此处可以认为是 ok 的
	return true
}
