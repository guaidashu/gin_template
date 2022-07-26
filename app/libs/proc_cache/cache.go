package proc_cache

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var Cache *Group

func init() {
	Cache = &Group{
		cache: cache.New(5*time.Minute, 10*time.Minute),
		lock:  sync.Mutex{},
		data:  make(map[string]*cacheData),
	}
}
