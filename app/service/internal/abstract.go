/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/07/2021
 * @Desc: desc
 */

package internal

type (
	abstract struct {
		keyMap      map[string]string // 所有的key存在此map里, 初始化的时候从redis获取
		keyMapCache KeyMapCache
	}
)

// 设置键值后同步到redis
func (c *abstract) setKeys(keys ...string) {
	for _, v := range keys {
		if _, ok := c.keyMap[v]; ok {
			continue
		}

		c.keyMap[v] = "1"
		// 内存不存在，则去redis尝试获取
		exists, _ := c.keyMapCache.Exists(v)
		if exists {
			continue
		}

		_ = c.keyMapCache.SetKeyMaps(v)
	}
}

// 删除键值后同步到redis
func (c *abstract) delKeys(keys ...string) {
	for _, v := range keys {
		delete(c.keyMap, v)
	}

	_ = c.keyMapCache.DelKeyMaps(keys...)
}
