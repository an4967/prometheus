package promql

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/prometheus/prometheus/promql/parser"
)

type cacheMap struct {
	sync.RWMutex
	m map[uint64]parser.Value
}

var once sync.Once
var instance *cacheMap
var cache_list *list.List

const defaultMapSize = 10000000
const maxMapSize = 100000000

// Get is a wrapper for getting the value from the underlying map
func (r *cacheMap) Get(key uint64) (parser.Value, bool) {
	r.RLock()
	defer r.RUnlock()

	v, ok := r.m[key]
	return v, ok
}

// Set is a wrapper for setting the value of a key in the underlying map
func (r *cacheMap) Set(key uint64, val parser.Value) {
	r.Lock()
	defer r.Unlock()
	if maxMapSize <= len(r.m) {
		fmt.Println("Remove 5% cache : ", len(r.m))
		for i := 0; i < maxMapSize/20; i++ {
			old_item := cache_list.Front()
			delete(r.m, old_item.Value.(uint64))
			cache_list.Remove(old_item)
		}
		fmt.Println("Current Map Size : ", len(r.m))
	}
	cache_list.PushBack(key)
	r.m[key] = val
}

func NewCacheMap() *cacheMap {

	once.Do(func() {
		cache_list = list.New()
		instance = &cacheMap{m: make(map[uint64]parser.Value, defaultMapSize)}
	})

	return instance
}
