//
//	redis Hash相关的所有操作
//

package redis

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	_ "github.com/mangenotwork/mange_redis_manage/common/manlog"
)

//获取Hash value
func HashHGETALL(rc redis.Conn, keyname string) map[string]string {
	fmt.Println("执行redis : ", "HGETALL", keyname)
	res, err := redis.StringMap(rc.Do("HGETALL", keyname))
	if err != nil {
		fmt.Println("GET error", err.Error())
	}
	fmt.Println(res)
	return res
}

//新建Hash 单个field
// 如果 key 不存在，一个新的哈希表被创建并进行 HSET 操作。
// 如果域 field 已经存在于哈希表中，旧值将被覆盖。
func HashHSET(rc redis.Conn, keyname, field string, value interface{}) string {
	fmt.Println("执行redis : ", "HSET", keyname, field, value)
	res, err := redis.String(rc.Do("HSET", keyname, field, value))
	if err != nil {
		fmt.Println("GET error", err.Error())
	}
	fmt.Println(res)
	return res
}

// 新建Hash 多个field
// HMSET key field value [field value ...]
// 同时将多个 field-value (域-值)对设置到哈希表 key 中。
// 此命令会覆盖哈希表中已存在的域。
func HashHMSET(rc redis.Conn, keyname string, values []interface{}) error {
	args := redis.Args{}.Add(keyname)
	for _, value := range values {
		fmt.Println(value)
		for k, v := range value.(map[string]interface{}) {
			args = args.Add(k)
			args = args.Add(v)
		}
	}
	fmt.Println("执行redis : ", "HMSET", args)
	res, err := rc.Do("HMSET", args...)
	if err != nil {
		fmt.Println("GET error", err.Error())
		return err
	}
	fmt.Println(res)
	return nil
}

//HSETNX key field value
// 给hash追加field value
//将哈希表 key 中的域 field 的值设置为 value ，当且仅当域 field 不存在。
func HashHSETNX(rc redis.Conn, keyname, field string, value interface{}) error {
	fmt.Println("执行redis : ", "HSETNX", keyname, field, value)
	res, err := rc.Do("HSETNX", keyname, field, value)
	if err != nil {
		fmt.Println("GET error", err.Error())
		return err
	}
	fmt.Println(res)
	return nil
}

// HDEL key field [field ...]
// 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。
func HashHDEL(rc redis.Conn, keyname string, fields []string) error {
	args := redis.Args{}.Add(keyname)
	for _, v := range fields {
		args = args.Add(v)
	}
	fmt.Println("执行redis : ", "HDEL", args)
	res, err := rc.Do("HDEL", args)
	if err != nil {
		fmt.Println("GET error", err.Error())
		return err
	}
	fmt.Println(res)
	return nil
}

// HEXISTS key field
// 查看哈希表 key 中，给定域 field 是否存在。
func HashHEXISTS(rc redis.Conn, keyname, fields string) bool {
	fmt.Println("执行redis : ", "HEXISTS", keyname, fields)
	res, err := redis.Int(rc.Do("HEXISTS", keyname, fields))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return false
	}

	if res == 0 {
		return false
	}

	fmt.Println(res)
	return true
}

// HGET key field
// 返回哈希表 key 中给定域 field 的值。
func HashHGET(rc redis.Conn, keyname, fields string) (res string, err error) {
	fmt.Println("执行redis : ", "HGET", keyname, fields)
	res, err = redis.String(rc.Do("HGET", keyname, fields))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// HINCRBY key field increment
// 为哈希表 key 中的域 field 的值加上增量 increment 。
// 增量也可以为负数，相当于对给定域进行减法操作。
// 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
// 如果域 field 不存在，那么在执行命令前，域的值被初始化为 0
func HashHINCRBY(rc redis.Conn, keyname, field string, increment int64) (res int64, err error) {
	fmt.Println("执行redis : ", "HINCRBY", keyname, field, increment)
	res, err = redis.Int64(rc.Do("HINCRBY", keyname, field, increment))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// HINCRBYFLOAT key field increment
// 为哈希表 key 中的域 field 加上浮点数增量 increment 。
// 如果哈希表中没有域 field ，那么 HINCRBYFLOAT 会先将域 field 的值设为 0 ，然后再执行加法操作。
// 如果键 key 不存在，那么 HINCRBYFLOAT 会先创建一个哈希表，再创建域 field ，最后再执行加法操作。
func HashHINCRBYFLOAT(rc redis.Conn, keyname, field string, increment float64) (res float64, err error) {
	fmt.Println("执行redis : ", "HINCRBYFLOAT", keyname, field, increment)
	res, err = redis.Float64(rc.Do("HINCRBYFLOAT", keyname, field, increment))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// HKEYS key
// 返回哈希表 key 中的所有域。
func HashHKEYS(rc redis.Conn, keyname string) (res []string, err error) {
	fmt.Println("执行redis : ", "HKEYS", keyname)
	res, err = redis.Strings(rc.Do("HKEYS", keyname))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// HLEN key
// 返回哈希表 key 中域的数量。
func HashHLEN(rc redis.Conn, keyname string) (res int64, err error) {
	fmt.Println("执行redis : ", "HLEN", keyname)
	res, err = redis.Int64(rc.Do("HLEN", keyname))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// HMGET key field [field ...]
// 返回哈希表 key 中，一个或多个给定域的值。
// 如果给定的域不存在于哈希表，那么返回一个 nil 值。
func HashHMGET(rc redis.Conn, keyname string, fields []string) (res []string, err error) {
	args := redis.Args{}.Add(keyname)
	for _, v := range fields {
		args = args.Add(v)
	}
	fmt.Println("执行redis : ", "HMGET", args)
	res, err = redis.Strings(rc.Do("HMGET", args))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// HVALS key
// 返回哈希表 key 中所有域的值。
func HashHVALS(rc redis.Conn, keyname string) (res []string, err error) {
	fmt.Println("执行redis : ", "HVALS", keyname)
	res, err = redis.Strings(rc.Do("HVALS", keyname))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

//HSCAN
//搜索value hscan test4 0 match *b*
