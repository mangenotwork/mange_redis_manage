//
//	redis string 的所有服务
//
package service

import (
	"fmt"
	"unicode/utf8"

	"github.com/garyburd/redigo/redis"
	manredis "github.com/mangenotwork/mange_redis_manage/common/redis"
)

type RedisString struct{}

//获取string数据
func (this *RedisString) Get(rc redis.Conn, key string) (value interface{}, keytype, size string, err error) {
	keytype = "string"
	string_value := manredis.StringGet(rc, key)
	value = string_value
	//计算字节大小,单位b,  1kb=1024b
	string_value_size := utf8.RuneCount([]byte(string_value))
	size = fmt.Sprintf("%dByte", string_value_size)
	err = nil
	return
}

//创建,修改,string数据
func (this *RedisString) Create(rc redis.Conn, key string, value interface{}) error {
	return manredis.StringSET(rc, key, value)
}

//追加 string数据
func (this *RedisString) Append(rc redis.Conn, key string, value interface{}) error {
	return manredis.StringAPPEND(rc, key, value)
}

func (this *RedisString) Del() {}

func (this *RedisString) Update() {}

//返回key value的大小 单位:b
func (this *RedisString) ValueSize(rc redis.Conn, key string) (size int64) {
	string_value := manredis.StringGet(rc, key)
	size = int64(utf8.RuneCount([]byte(string_value)))
	return
}
