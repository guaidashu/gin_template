// Package singleflight provides a duplicate function call suppression
// mechanism.
package proc_cache

import (
	"context"
	"gin_template/app/libs"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

const (
	EmptyMark           emptyType = "*"
	EmptyMarkExpireTime           = time.Second * 5 * 1
)

type (
	emptyType string

	cacheData struct {
		val       interface{}
		waitGroup sync.WaitGroup
		err       error
	}

	Group struct {
		cache *cache.Cache // lazily initialized
		lock  sync.Mutex
		data  map[string]*cacheData
	}
)

// 注解：
// 这个操作方式会阻塞后来的请求，同一个key来10个请求，只会执行一次函数体(内可以为mysql查询等)，其他的请求等待第一个执行的结果
// 所有请求将会得到同一个结果，更详细看代码
func (g *Group) Do(key string, fn func() (interface{}, error), expireTime ...time.Duration) (interface{}, error) {
	// 加大锁防止内部锁未初始化争抢
	g.lock.Lock()

	if data, ok := g.data[key]; ok {
		// 解大锁
		g.lock.Unlock()
		// wait group等待
		data.waitGroup.Wait()

		return data.val, data.err
	}

	c := new(cacheData)
	c.waitGroup.Add(1)
	g.data[key] = c
	// 解锁后用waitGroup在等待了
	g.lock.Unlock()

	// 先从 go-cache 获取数据，如果存在，则直接跳过
	if data, ok := g.cache.Get(key); ok {
		c.val = data
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		ch := make(chan int)
		libs.RunSafe(func() {
			c.val, c.err = fn()
			ch <- 1
		})
		select {
		case <-ch:
		case <-ctx.Done():
		}
		cancel()

		expire := time.Duration(0)
		if c.val == nil {
			g.cache.Set(key, EmptyMark, EmptyMarkExpireTime)
		} else {
			if len(expireTime) > 0 {
				expire = expireTime[0]
			}
			g.Set(key, c.val, expire)
		}
	}

	// 数据获取完毕，结束
	c.waitGroup.Done()

	// 移除map数据，这里不用担心，因为如果有后续请求，会从头走流程，不影响获取
	g.lock.Lock()
	delete(g.data, key)
	g.lock.Unlock()

	return c.val, c.err
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
	if s, ok := data.(emptyType); ok && s == EmptyMark {
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
