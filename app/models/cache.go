/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 28/10/2021
 * @Desc: 缓存相关
 */

package models

import (
	"encoding/json"
	"gin_template/app/rds"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"time"
)

// 默认过期时间
const (
	defaultExpiry         = time.Hour * 24 * 7
	defaultNotFoundExpiry = time.Minute
)

// db 缓存层, 主要用于单数据全量缓存
// 列表不建议用缓存, 个别可以使用
type (
	// QueryFn defines the query method.
	QueryFn func(conn *gorm.DB, v interface{}) error

	// CacheConn 自定义 model 缓存层
	CacheConn struct {
		cache *redis.Client
	}
)

func NewCacheConn() CacheConn {
	return CacheConn{cache: rds.Redis}
}

// 获取列表
func (c *CacheConn) QueryRows(v interface{}, key string, db *gorm.DB, fn QueryFn) error {
	return c.GetResult(v, key, func(v interface{}) error {
		return fn(db, v)
	})
}

// 获取单条数据
func (c *CacheConn) QueryRow(v interface{}, key string, db *gorm.DB, fn QueryFn) error {
	return c.GetResult(v, key, func(v interface{}) error {
		return fn(db, v)
	})
}

// 具体的redis 操作
func (c *CacheConn) GetResult(v interface{}, key string, query func(v interface{}) error) error {
	// 加分布式锁, 防止缓存击穿
	lock := rds.NewRedisLock("get-" + key)
	if ok, err := lock.Acquire(); !ok {
		return err
	}
	defer func() {
		_ = lock.Release()
	}()

	// 先从redis获取数据
	val, err := c.cache.Get(key).Result()
	switch err {
	case redis.Nil:
		// 从数据库获取
		err = query(v)
		if err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				val = "*"
			default:
				return err
			}
		}

		// 设置缓存
		if val == "*" {
			// not found的缓存应该未一分钟
			err = c.cache.Set(key, val, defaultNotFoundExpiry).Err()
			if err != nil {
				return err
			}
			return gorm.ErrRecordNotFound
		}

		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return c.cache.Set(key, string(b), defaultExpiry).Err()
	case nil:
		if val == "*" {
			v = nil
			return gorm.ErrRecordNotFound
		}
		// 构造数据并返回
		return json.Unmarshal([]byte(val), v)
	default:
		return err
	}
}

// DelCache deletes cache with keys.
func (c *CacheConn) DelCache(keys ...string) error {
	return c.cache.Del(keys...).Err()
}
