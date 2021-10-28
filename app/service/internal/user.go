/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/07/2021
 * @Desc: desc
 */

package internal

import "sync"

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
			abstract: abstract{
				keyMap:   NewKeyMapCache().GetKeyMaps("userKeyMapCache"),
				cacheKey: "userKeyMapCache",
			},
		}
	})

	return _userCache
}
