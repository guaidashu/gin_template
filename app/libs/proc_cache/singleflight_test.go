package proc_cache

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"sync"
	"testing"
	"time"
)

func TestGroup_Get(t *testing.T) {
	g := &Group{
		cache: cache.New(5*time.Minute, 10*time.Minute),
		lock:  sync.Mutex{},
		data:  make(map[string]*cacheData),
	}

	// data, ok := g.Get("test", func() (interface{}, error) {
	// 	return "我是测试", nil
	// }, time.Second*2)
	// fmt.Println("ok ==>", ok, " data ====>", data)

	for i := 0; i < 1000; i++ {
		go func() {
			data, ok := g.Get("test", func() (interface{}, error) {
				return "我是测试", nil
			}, time.Second*2)
			fmt.Println("ok ==>", ok, " data ====>", data)
			// data, ok = g.Get("test", func() (interface{}, error) {
			// 	fmt.Println("走到这里了吗")
			// 	return "1111111", nil
			// })
			// fmt.Println("ok ==>", ok, " data ====>", data)
		}()
	}

	data, ok := g.Get("test_empty", func() (interface{}, error) {
		fmt.Println("空数据测试")
		return nil, nil
	})
	fmt.Println("ok ==>", ok, " data ====>", data)

	data, ok = g.Get("test_empty", func() (interface{}, error) {
		fmt.Println("空数据测试")
		return "空数据测试第二次尝试有数据没", nil
	})
	fmt.Println("ok ==>", ok, " data ====>", data)

	time.Sleep(time.Second * 3)
}