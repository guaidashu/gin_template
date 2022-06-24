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
		m:     make(map[string]*sync.Mutex),
		cache: cache.New(5*time.Minute, 10*time.Minute),
		lock:  &sync.Mutex{},
		mw:    make(map[string]*sync.WaitGroup),
	}

	data, ok := g.Get("test", func() (interface{}, error) {
		return nil, nil
	})
	fmt.Println("ok ==>", ok, " data ====>", data)

	data, ok = g.Get("test", func() (interface{}, error) {
		fmt.Println("走到这里了吗")
		return "1111111", nil
	})
	fmt.Println("ok ==>", ok, " data ====>", data)
}
