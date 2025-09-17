/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/07/2021
 * @Desc: desc
 */

package internal

import (
	"gin_template/app/rds"
	"sync"

	"github.com/redis/go-redis/v9"
)

type (
	UserCache interface {
	}

	defaultUserCache struct {
		abstract
	}
)

var (
	_userCache     UserCache
	_userCacheOnce sync.Once
)

func NewUserCache() UserCache {
	_userCacheOnce.Do(func() {
		_userCache = &defaultUserCache{
			abstract: abstract{},
		}

		_userCache.(*defaultUserCache).init()
	})

	return _userCache
}

func (c *defaultUserCache) init() {
	// 初始化 key map
	c.abstract.keyMapCache = NewKeyMapCache("defaultUserCache", c.client())
	c.abstract.keyMap = c.abstract.keyMapCache.GetKeyMaps()
}

func (c *defaultUserCache) client() *redis.Client {
	return rds.Redis
}
