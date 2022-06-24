// Package singleflight provides a duplicate function call suppression
// mechanism.
package proc_cache

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

const (
	EmptyMark           = "*"
	EmptyMarkExpireTime = time.Second * 60 * 1
)

type Group struct {
	m     map[string]*sync.Mutex // lazily initialized
	cache *cache.Cache           // lazily initialized
	lock  *sync.Mutex
	mw    map[string]*sync.WaitGroup
}

// 注解：
// 这个操作方式会阻塞后来的请求，同一个key来10个请求，只会执行一次函数体(内可以为mysql查询等)，其他的请求等待第一个执行的结果
// 所有请求将会得到同一个结果，更详细看代码
func (g *Group) Do(key string, fn func() (interface{}, error), expireTime ...time.Duration) (interface{}, error) {
	// 加大锁防止内部锁未初始化争抢
	g.lock.Lock()
	if _, ok := g.m[key]; !ok {
		g.m[key] = &sync.Mutex{}
	}
	// 解大锁
	g.lock.Unlock()

	expire := time.Duration(0)
	g.m[key].Lock()

	if data, ok := g.cache.Get(key); ok {
		g.m[key].Unlock()
		// wait group存在则等待，不存在则直接返回
		if _, exist := g.mw[key]; exist {
			g.mw[key].Wait()
		}
		return data, nil
	}

	g.mw[key] = &sync.WaitGroup{}
	g.mw[key].Add(1)

	data, err := fn()
	if len(expireTime) > 0 {
		expire = expireTime[0]
	}
	if err != nil || data == nil {
		g.cache.Set(key, EmptyMark, EmptyMarkExpireTime)
	}
	g.Set(key, data, expire)
	g.m[key].Unlock()
	g.mw[key].Done()

	// 移除锁
	g.lock.Lock()
	delete(g.m, key)
	delete(g.mw, key)
	g.lock.Unlock()

	return data, err
}

// 包装一层set
func (g *Group) Set(k string, x interface{}, d time.Duration) {
	g.cache.Set(k, x, d)
}

// 封装的Get方法(单独Get)，带防缓存穿透
// 如果空数据，传入的方法fn需要返回空数据并且err返回nil
func (g *Group) Get(key string, fn func() (interface{}, error), expireTime ...time.Duration) (interface{}, bool) {
	data, err := g.Do(key, fn, expireTime...)

	// 先判断是否为空数据标记
	if s, ok := data.(string); ok && s == EmptyMark {
		return nil, false
	}

	if err == nil && data != nil {
		return data, true
	}

	return nil, false
}

// 直接返回go-cache实例(为了更自由灵活的操作)
func (g *Group) Cache() *cache.Cache {
	return g.cache
}
