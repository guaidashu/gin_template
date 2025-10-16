/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 02/08/2021
 * @Desc: 键值维护模块
 */

package internal

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

// 键值维护模块
// 由于redis 集群Keys函数很可能导致redis卡死,
// 并且Scan函数很可能会不起效 除非花钱升级云redis集群 但是成本是不允许的
// 而且会造成其他一系列的复杂操作, 所以考虑自行维护 一些必要的键值
// 比如定时任务等 后续需要删除的键值, 就可以通过此模块进行维护
// PS: 有定时任务等相关要删除键值操作的才建议用此模块进行操作,否则会多出一些多余的redis读写
type (
	KeyMapCache interface {
		// 通过key 获取原有的值
		GetKeyMaps() (data map[string]string)
		// 设置key map
		SetKeyMaps(key string) error
		// 删除key map
		DelKeyMaps(key ...string) error
		// key是否存在
		Exists(key string) (bool, error)
	}

	defaultKeyMapCache struct {
		cacheKey string // 服务对应的键值对数据 key
		client   *redis.Client
	}
)

var (
	_keyMapCache     KeyMapCache
	_keyMapCacheOnce sync.Once
)

func NewKeyMapCache(cacheKey string, client *redis.Client) KeyMapCache {
	_keyMapCacheOnce.Do(func() {
		_keyMapCache = &defaultKeyMapCache{
			client:   client,
			cacheKey: cacheKey,
		}
	})

	return _keyMapCache
}

func (c *defaultKeyMapCache) GetKeyMaps() (data map[string]string) {
	data = make(map[string]string)
	ctx := context.Background()
	val, err := c.client.HGetAll(ctx, c.cacheKey).Result()
	if err != nil || val == nil {
		return
	}

	data = val
	return
}

func (c *defaultKeyMapCache) SetKeyMaps(key string) error {
	ctx := context.Background()
	return c.client.HSet(ctx, c.cacheKey, key, "1").Err()
}

func (c *defaultKeyMapCache) DelKeyMaps(key ...string) error {
	ctx := context.Background()
	return c.client.HDel(ctx, c.cacheKey, key...).Err()
}

func (c *defaultKeyMapCache) Exists(key string) (bool, error) {
	ctx := context.Background()
	return c.client.HExists(ctx, c.cacheKey, key).Result()
}
