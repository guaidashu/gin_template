/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 02/09/2021
 * @Desc: 购物车
 */

package internal

import "sync"

type (
	ShoppingCartCache interface {
	}

	defaultShoppingCartCache struct {
		*abstract
		keyMap map[string]int64 // 所有的key存在此map里, 初始化的时候从redis获取
	}
)

var (
	_shoppingCartCache     ShoppingCartCache
	_shoppingCartCacheOnce sync.Once
)

func NewShoppingCartCache() ShoppingCartCache {
	_shoppingCartCacheOnce.Do(func() {
		_shoppingCartCache = &defaultShoppingCartCache{
			abstract: &abstract{
				keyMap:   NewKeyMapCache().GetKeyMaps("shoppingCartKeyMapCache"),
				cacheKey: "shoppingCartKeyMapCache",
			},
		}
	})

	return _shoppingCartCache
}

func (c *defaultShoppingCartCache) AddMerchandise() {
	c.setKeys()
}
