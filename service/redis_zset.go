//
//	redis zset 的所有服务
//
package service

import (
	"fmt"
	"unicode/utf8"

	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/mange_redis_manage/common"
	manredis "github.com/mangenotwork/mange_redis_manage/common/redis"
)

type RedisZSet struct{}

//获取数据
func (this *RedisZSet) Get(rc redis.Conn, key string) (value interface{}, keytype, size string, err error) {
	keytype = "zset"
	zset_value := manredis.ZSetZRANGEALL(rc, key)
	zset_value_map := make(map[int]string, 0)
	zset_size := 0
	for v, k := range zset_value {
		kStr := common.Uint82Str(k.([]uint8))
		zset_value_map[v] = kStr
		zset_size = zset_size + utf8.RuneCount([]byte(kStr))
	}
	value = zset_value_map
	size = fmt.Sprintf("%dByte", zset_size)
	err = nil
	return
}

//创建,修改
func (this *RedisZSet) Create(rc redis.Conn, key string, value interface{}) error {
	values := value.([]interface{})
	return manredis.ZSetZADD(rc, key, values)
}

//追加
func (this *RedisZSet) Append(rc redis.Conn, key string, value interface{}) error {
	return manredis.ZSetZADD(rc, key, value.([]interface{}))
}

//删除
func (this *RedisZSet) Del() {}

func (this *RedisZSet) Update() {}

//返回key value的大小 单位:b
func (this *RedisZSet) ValueSize(rc redis.Conn, key string) (size int64) {
	zset_value := manredis.ZSetZRANGEALL(rc, key)
	size = 0
	for _, v := range zset_value {
		kStr := common.Uint82Str(v.([]uint8))
		size = size + int64(utf8.RuneCount([]byte(kStr)))
	}
	return
}
