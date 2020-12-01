//
//	redis hash 的所有服务
//
package service

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/garyburd/redigo/redis"
	manredis "github.com/mangenotwork/mange_redis_manage/common/redis"
)

type RedisHash struct{}

//获取数据
func (this *RedisHash) Get(rc redis.Conn, key string) (value interface{}, keytype, size string, err error) {
	keytype = "hash"
	hash_value := manredis.HashHGETALL(rc, key)
	hash_value_size := unsafe.Sizeof(hash_value)
	hash_value_result, _ := json.Marshal(hash_value)
	value = string(hash_value_result)
	size = fmt.Sprintf("%dByte", hash_value_size)
	err = nil
	return
}

//创建,修改
func (this *RedisHash) Create(rc redis.Conn, key string, value interface{}) error {
	values := value.([]interface{})
	return manredis.HashHMSET(rc, key, values)
}

//追加
func (this *RedisHash) Append(rc redis.Conn, key string, value interface{}) error {
	var field string
	var data interface{}
	for k, v := range value.(map[string]interface{}) {
		field = k
		data = v
	}
	return manredis.HashHSETNX(rc, key, field, data)
}

//删除
func (this *RedisHash) Del() {}

func (this *RedisHash) Update() {}

//返回key value的大小 单位:b
func (this *RedisHash) ValueSize(rc redis.Conn, key string) (size int64) {
	hash_value := manredis.HashHGETALL(rc, key)
	size = int64(unsafe.Sizeof(hash_value))
	return
}
