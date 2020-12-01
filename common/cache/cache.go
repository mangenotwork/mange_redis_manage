package cache

import (
	"os"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type Caches struct {
	C    *cache.Cache
	lock sync.Mutex
}

//	AllCaches　所有缓存
var (
	caches = new(Caches)
)

func init() {
	caches.C = cache.New(60*24*7*time.Minute, 10*time.Minute)
	caches.lock = sync.Mutex{}
}

func Set(key string, value interface{}) {
	caches.lock.Lock()
	caches.C.Set(key, value, cache.DefaultExpiration)
	caches.lock.Unlock()
	return
}

func SetAlways(key string, value interface{}) {
	caches.lock.Lock()
	caches.C.Set(key, value, -1)
	caches.lock.Unlock()
	return
}

func Get(key string) (value interface{}, isOk bool) {
	caches.lock.Lock()
	value, isOk = caches.C.Get(key)
	caches.lock.Unlock()
	return
}

func Save2File() {
	file_name := "temp"
	pwd, _ := os.Getwd()
	file_path := pwd + "/cache/" + file_name

	caches.C.SaveFile(file_path)
}
