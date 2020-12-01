//
//	redis list 的所有服务
//
package service

import (
	"fmt"
	"unicode/utf8"

	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/mange_redis_manage/common"
	manredis "github.com/mangenotwork/mange_redis_manage/common/redis"
)

type RedisList struct{}

//获取数据
func (this *RedisList) Get(rc redis.Conn, key string) (value interface{}, keytype, size string, err error) {
	keytype = "list"
	list_value := manredis.ListLRANGEALL(rc, key)
	list_value_map := make(map[int]string, 0)
	list_size := 0
	for v, k := range list_value {
		kStr := common.Uint82Str(k.([]uint8))
		list_value_map[v] = kStr
		list_size = list_size + utf8.RuneCount([]byte(kStr))
	}
	value = list_value_map
	size = fmt.Sprintf("%dByte", list_size)
	err = nil
	return
}

//创建,修改
func (this *RedisList) Create(rc redis.Conn, key string, value interface{}) error {
	values := value.([]interface{})
	return manredis.ListLPUSH(rc, key, values)
}

//追加
func (this *RedisList) Append(rc redis.Conn, key string, value interface{}) error {
	return manredis.ListRPUSH(rc, key, value.([]interface{}))
}

//删除
func (this *RedisList) Del() {}

func (this *RedisList) Update() {}

//返回key value的大小 单位:b
func (this *RedisList) ValueSize(rc redis.Conn, key string) (size int64) {
	list_value := manredis.ListLRANGEALL(rc, key)
	size = 0
	for _, v := range list_value {
		kStr := common.Uint82Str(v.([]uint8))
		size = size + int64(utf8.RuneCount([]byte(kStr)))
	}
	return
}
