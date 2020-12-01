//
//	redis set 的所有服务
//
package service

import (
	"fmt"
	"unicode/utf8"

	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/mange_redis_manage/common"
	manredis "github.com/mangenotwork/mange_redis_manage/common/redis"
)

type RedisSet struct{}

//获取数据
func (this *RedisSet) Get(rc redis.Conn, key string) (value interface{}, keytype, size string, err error) {
	keytype = "set"
	set_value := manredis.SetSMEMBERS(rc, key)
	set_value_map := make(map[int]string, 0)
	set_size := 0
	for v, k := range set_value {
		kStr := common.Uint82Str(k.([]uint8))
		set_value_map[v] = kStr
		set_size = set_size + utf8.RuneCount([]byte(kStr))
	}
	value = set_value_map
	size = fmt.Sprintf("%dByte", set_size)
	err = nil
	return
}

//创建,修改
func (this *RedisSet) Create(rc redis.Conn, key string, value interface{}) error {
	values := value.([]interface{})
	return manredis.SetSADD(rc, key, values)
}

//追加
func (this *RedisSet) Append(rc redis.Conn, key string, value interface{}) error {
	return manredis.SetSADD(rc, key, value.([]interface{}))
}

//删除
func (this *RedisSet) Del() {}

//修改
func (this *RedisSet) Update() {}

//返回key value的大小 单位:b
func (this *RedisSet) ValueSize(rc redis.Conn, key string) (size int64) {
	set_value := manredis.SetSMEMBERS(rc, key)
	size = 0
	for _, v := range set_value {
		kStr := common.Uint82Str(v.([]uint8))
		size = size + int64(utf8.RuneCount([]byte(kStr)))
	}
	return
}
