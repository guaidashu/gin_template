/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/07/2021
 * @Desc: desc
 */

package internal

import (
	"gin_template/app/rds"

	"github.com/go-redis/redis"
)

type (
	abstract struct {
		keyMap   map[string]int64 // 所有的key存在此map里, 初始化的时候从redis获取
		cacheKey string           // 服务对应的键值对数据 key
	}
)

// client 获取客户端
func (c *abstract) client() *redis.Client {
	return rds.Redis
}

// 设置键值后同步到redis
func (c *abstract) setKeys(keys ...string) {
	for _, v := range keys {
		if c.keyMap[v] == 0 {
			c.keyMap[v] = 1
		}
	}

	_ = NewKeyMapCache().SetKeyMaps(c.cacheKey, c.keyMap)
}

// 删除键值后同步到redis
func (c *abstract) delKeys(keys ...string) {
	for _, v := range keys {
		delete(c.keyMap, v)
	}

	_ = NewKeyMapCache().SetKeyMaps(c.cacheKey, c.keyMap)
}
