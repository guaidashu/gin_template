/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 02/09/2021
 * @Desc: 购物车
 */

package internal

import (
	"gin_template/app/rds"
	"sync"

	"github.com/redis/go-redis/v9"
)

type (
	ShoppingCartCache interface {
	}

	defaultShoppingCartCache struct {
		abstract *abstract
	}
)

var (
	_shoppingCartCache     ShoppingCartCache
	_shoppingCartCacheOnce sync.Once
)

func NewShoppingCartCache() ShoppingCartCache {
	_shoppingCartCacheOnce.Do(func() {
		_shoppingCartCache = &defaultShoppingCartCache{
			abstract: &abstract{},
		}

		_shoppingCartCache.(*defaultShoppingCartCache).init()
	})

	return _shoppingCartCache
}

func (c *defaultShoppingCartCache) init() {
	// 初始化 key map
	c.abstract.keyMapCache = NewKeyMapCache("defaultShoppingCartCache", c.client())
	c.abstract.keyMap = c.abstract.keyMapCache.GetKeyMaps()
}

func (c *defaultShoppingCartCache) client() *redis.Client {
	return rds.Redis
}

// keyMap:   NewKeyMapCache().GetKeyMaps("shoppingCartKeyMapCache"),
// cacheKey: "shoppingCartKeyMapCache",
func (c *defaultShoppingCartCache) AddMerchandise() {
	c.abstract.setKeys()
}
