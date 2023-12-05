/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 02/08/2021
 * @Desc: 键值维护模块
 */

package internal

import (
	"sync"
)

// 键值维护模块
// 由于redis 集群Keys函数很可能导致redis卡死,
// 并且Scan函数很可能会不起效 除非花钱升级云redis集群 但是成本是不允许的
// 而且会造成其他一系列的复杂操作, 所以考虑自行维护 一些必要的键值
// 比如定时任务等 后续需要删除的键值, 就可以通过此模块进行维护
// PS: 有定时任务等相关要删除键值操作的才建议用此模块进行操作,否则会多出一些多余的redis读写
// 暂时只是单机版本, 如果后续业务有拓展微服务需求,则要通过redis全局锁进行加锁实现事务操作
type (
	KeyMapCache interface {
		// 通过key 获取原有的值
		GetKeyMaps(key string) (data map[string]string)
		// 设置key map
		SetKeyMaps(key string) error
		// 删除key map
		DelKeyMaps(key ...string) error
	}

	defaultKeyMapCache struct {
		abstract
	}
)

var (
	_keyMapCache     KeyMapCache
	_keyMapCacheOnce sync.Once
)

func NewKeyMapCache() KeyMapCache {
	_keyMapCacheOnce.Do(func() {
		_keyMapCache = &defaultKeyMapCache{}
	})

	return _keyMapCache
}

func (c *defaultKeyMapCache) GetKeyMaps(key string) (data map[string]string) {
	data = make(map[string]string)
	val, err := c.client().HGetAll(key).Result()
	if err != nil || val == nil {
		return
	}

	data = val
	return
}

func (c *defaultKeyMapCache) SetKeyMaps(key string) error {
	return c.client().HSet(c.cacheKey, key, "1").Err()
}

func (c *defaultKeyMapCache) DelKeyMaps(key ...string) error {
	return c.client().HDel(c.cacheKey, key...).Err()
}
