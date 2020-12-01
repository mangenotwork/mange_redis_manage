//
//	本地内存
//
package service

import (
	"github.com/mangenotwork/mange_redis_manage/common/cache"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
)

type CacheService interface {
	GetAll()                       //查看所有缓存
	Get(key string)                //查看指定缓存
	Update(Key string)             //修改指定缓存
	Del(key string)                //删除指定缓存
	Add(key string, v interface{}) //新增缓存
	Backup()                       //备份缓存
	SetConf()                      //设置本机缓存
}

type HostCache struct {
}

func (this *HostCache) GetAll() {
	c := new(cache.Caches)
	manlog.Debug(c, c.C.Items())
}

func (this *HostCache) Get(key string) {
}

func (this *HostCache) Update(key string) {
}

func (this *HostCache) Del(key string) {
}

func (this *HostCache) Add(key string, v interface{}) {
}

func (this *HostCache) Backup() {}

func (this *HostCache) SetConf() {}
